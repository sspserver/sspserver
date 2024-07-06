package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fasthttp/router"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/sspserver/udetect/transport/http"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adsource"
	"github.com/geniusrabbit/adcorelib/adsource/openrtb"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/eventtraking/pixelgenerator"
	"github.com/geniusrabbit/adcorelib/httpserver"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/endpoint"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/pixel"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/trakeraction"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
	"github.com/geniusrabbit/adcorelib/net/fasthttp/middleware"
	"github.com/geniusrabbit/adcorelib/personification"
	"github.com/geniusrabbit/adcorelib/simplepersondetector"
	"github.com/geniusrabbit/adcorelib/storage/accessors/adsourceaccessor"
	"github.com/geniusrabbit/adcorelib/storage/accessors/companyaccessor"
	"github.com/geniusrabbit/adcorelib/storage/accessors/formataccessor"
	"github.com/geniusrabbit/adcorelib/storage/accessors/zoneaccessor"
	"github.com/geniusrabbit/adcorelib/urlgenerator"
	"github.com/geniusrabbit/adcorelib/zlogger"

	"github.com/sspserver/sspserver/cmd/sspserver/appcontext"
	"github.com/sspserver/sspserver/cmd/sspserver/datainit"
	"github.com/sspserver/sspserver/internal/endpoint/direct"
	"github.com/sspserver/sspserver/internal/endpoint/dynamic"
	"github.com/sspserver/sspserver/internal/endpoint/proxy"
	"github.com/sspserver/sspserver/internal/stream"
)

const (
	eventsStreamName   = "events"
	userInfoStreamName = "user_info"
	winStreamName      = "wins"
	adInfoStreamName   = "ad_info"
)

var (
	config       appcontext.Config
	buildCommit  = ""
	buildVersion = "develop"
	buildDate    = "unknown"
)

func init() {
	fatalError(config.Load(), "config loading")

	// Init new logger object
	loggerObj, err := zlogger.New(config.ServiceName, config.LogEncoder,
		config.LogLevel, config.LogAddr, zap.Fields(
			zap.String("commit", buildCommit),
			zap.String("version", buildVersion),
			zap.String("build_date", buildDate),
		))
	fatalError(err, "configure logger")

	zap.ReplaceGlobals(loggerObj)

	if config.IsDebug() {
		fmt.Println(config.String())
	}
}

func main() {
	var (
		logger       = zap.L()
		adServerConf = &config.AdServer
		adSSPConf    = &config.AdServer.SSP
		ctx, cancel  = context.WithCancel(context.Background())
	)
	defer cancel()

	// Connect to advertisement data source
	dataAccessor, err := datainit.Connect(ctx, config.AdServer.Storage.Connection)
	fatalError(err, "advertisement data")

	// Register event streams
	{
		configureEventPipeline(ctx, adServerConf)

		// Run notification listener
		go func() {
			if err := nc.Listen(ctx); err != nil {
				logger.Error("notification listener", zap.Error(err))
			}
		}()

		// Close notification processors
		defer nc.Close()
	}

	// Event flow processor
	eventGenerator := eventgenerator.New(config.ServiceName)
	eventStream := eventstream.New(
		nc.PublisherByName(eventsStreamName),
		nc.PublisherByName(userInfoStreamName),
		eventGenerator,
	)
	ctx = eventstream.WithStream(ctx, eventStream)

	// Win processor store into the context of requests
	winStream := eventstream.WinNotifications(nc.PublisherByName(winStreamName))
	ctx = eventstream.WithWins(ctx, winStream)

	// URL generator object
	urlGenerator := (&urlgenerator.Generator{
		EventGenerator: eventGenerator,
		PixelGenerator: pixelgenerator.NewPixelGenerator(adServerConf.TrackerHost),
		CDNDomain:      adServerConf.CDNDomain,
		LibDomain:      adServerConf.LibDomain,
		ClickPattern:   "/click?c={code}",
		DirectPattern:  "/direct?c={code}",
		WinPattern:     "/win?c={code}",
	}).Init()

	// Init side modules
	datainit.Initialize(config.IsDebug(), urlGenerator)

	// Init format accessor (format types of advertisement)
	formatAccessor := mustFormatAccessor(ctx, dataAccessor)
	ctx = formataccessor.WithContext(ctx, formatAccessor)

	// Init company (similar to client account) data accessor
	companyDataAccessor, err := dataAccessor(ctx, "company")
	fatalError(err, "company data accessor")

	companyAccessor := companyaccessor.NewCompanyAccessor(companyDataAccessor)

	// Init source data accessor (ad sources like: RTB, direct, etc)
	sourceDataAccessor, err := dataAccessor(ctx, "source")
	fatalError(err, "RTB source data accessor")

	// Init advertisement source accessor (provides multiple sources of advertisement access as container)
	sourceAccessor, err := adsourceaccessor.NewAccessor(ctx,
		sourceDataAccessor, companyAccessor, openrtb.NewFactory())
	fatalError(err, "RTB source accessor")

	// Init target data accessor (targeting zones where advertisement will be shown)
	targetDataAccessor, err := dataAccessor(ctx, "zone")
	fatalError(err, "target(zone) data accessor")

	targetAccessor := zoneaccessor.NewZoneAccessor(targetDataAccessor, companyAccessor)

	// Configure advertisement source accessor (provides multiple sources of advertisement access as one source)
	adsourceWrapper, err := adsource.NewMultisourceWrapper(
		adsource.WithSourceAccessor(sourceAccessor),
		adsource.WithTimeout(time.Duration(adSSPConf.RequestTimeout)*time.Millisecond),
		adsource.WithMaxParallelRequests(adSSPConf.MaxParallelRequests),
	)
	fatalError(err, "create adsource accessor")

	// HTTP wrapper for extended handlers
	httpHandlerWrapper := httphandler.NewHTTPHandlerWrapper(
		nil, // func(*fasthttp.RequestCtx) {},
		func(*fasthttp.RequestCtx) context.Context { return ctx },
		logger.With(zap.String("module", "httpserver")),
	)

	// Init personification client
	personDetector := personification.Client(&simplepersondetector.SimpleClient{})
	if config.Person.Connect != "" {
		personDetector = personification.Connect(http.NewTransport(
			config.Person.Connect,
			http.WithTimeout(config.Person.RequestTimeout, config.Person.KeepAliveTimeout),
		))
	}

	// Init signature
	signature := personification.Signeture{
		UUIDName:       config.Person.UUIDCookieName,
		SessidName:     config.Person.SessiCookiedName,
		SessidLifetime: config.Person.SessionLifetime,
		Detector:       personDetector,
	}

	// Prepare spy middleware
	spyMiddleware := middleware.NewSpy(ctx, signature.Whois, signature.SignCookie)

	// HTTP server initialyze
	server, err := httpserver.NewServer(
		httpserver.WithDebugMode(config.IsDebug()),
		httpserver.WithServiceName(config.ServiceName),
		httpserver.WithLogger(logger.With(zap.String("module", "httpserver"))),
		httpserver.WithCustomHTTPServer(&fasthttp.Server{
			ReadBufferSize: 1 << 20,
			ReadTimeout:    config.Server.HTTP.ReadTimeout,
			WriteTimeout:   config.Server.HTTP.WriteTimeout,
		}),
		httpserver.WithCustomRouter(func(router *router.Router) {
			router.ServeFiles("/public/{filepath:*}", "/public")
		}),
		httpserver.WithExtensions(
			pixel.NewExtension(
				pixel.WithEventStream(eventStream),
				pixel.WithHTTPHandlerWrapper(httpHandlerWrapper),
			),
			trakeraction.NewExtension(
				trakeraction.WithEventStream(eventStream),
				trakeraction.WithHTTPHandlerWrapper(httpHandlerWrapper),
				trakeraction.WithURLGenerator(urlGenerator),
			),
			// Register HTTP endpoints extensions
			endpoint.NewExtension(
				endpoint.WithAdvertisementSource(adsourceWrapper),
				endpoint.WithHTTPHandlerWrapper(httpHandlerWrapper),
				endpoint.WithFormatAccessor(formatAccessor),
				endpoint.WithZoneAccessor(targetAccessor),
				endpoint.WithSpy(spyMiddleware),
				endpoint.WithSendpoints(
					direct.New(formatAccessor, config.AdServer.Logic.Direct.DefaultURL),
					dynamic.New(urlGenerator),
					proxy.New(),
				),
			),
		),
	)
	fatalError(err, "new HTTP server")

	fmt.Println("Run HTTP server", config.Server.HTTP.Listen)
	fatalError(server.Listen(ctx, config.Server.HTTP.Listen))
}

func configureEventPipeline(ctx context.Context, adServerConf *appcontext.AdServerConfig) {
	// Register events data stream
	fatalError(nc.Register(
		eventsStreamName,
		connectPublisherOrLog(ctx,
			eventsStreamName,
			adServerConf.EventPipeline.EventQueue.Connection,
			config.IsDebug(),
		),
	), "register events stream")

	// Register user info data stream
	fatalError(nc.Register(
		userInfoStreamName,
		connectPublisherOrLog(ctx,
			userInfoStreamName,
			adServerConf.EventPipeline.UserInfoQueue.Connection,
			config.IsDebug(),
		),
	), "register user info stream")

	// Register wins info data stream
	fatalError(nc.Register(
		winStreamName,
		connectPublisherOrLog(ctx,
			winStreamName,
			adServerConf.EventPipeline.WinQueue.Connection,
			config.IsDebug(),
		),
	), "register win stream")

	// Register adinfo data stream
	fatalError(nc.Register(
		adInfoStreamName,
		connectPublisherOrLog(ctx,
			adInfoStreamName,
			adServerConf.EventPipeline.AdInfoQueue.Connection,
			config.IsDebug(),
		),
	), "register ad info stream")
}

func notificationMessageLog(streamName string) nc.FuncPublisher {
	return func(ctx context.Context, msgs ...any) error {
		for _, msg := range msgs {
			zap.L().Debug("notification message",
				zap.String("stream", streamName),
				zap.Any("msg", msg),
			)
		}
		return nil
	}
}

func connectPublisherOrLog(ctx context.Context, name, connection string, debug bool) nc.Publisher {
	if connection != "" {
		pub, err := stream.ConnectPublisher(ctx, connection)
		fatalError(err, "connect to '"+connection+"' topics")
		if debug {
			pub = stream.WrapPublisherWithLog(name, pub)
		}
		return pub
	}
	zap.L().Info("register new dummy publisher", zap.String("name", name))
	return notificationMessageLog(name)
}

func mustFormatAccessor(ctx context.Context, dataAccessor datainit.DataLoaderAccessorFnk) types.FormatsAccessor {
	formatDataAccessor, err := dataAccessor(ctx, "format")
	fatalError(err, "format data accessor")
	return formataccessor.NewFormatAccessor(formatDataAccessor)
}

func fatalError(err error, message ...any) {
	if err != nil {
		zap.L().Fatal(fmt.Sprint(message...), zap.Error(err))
	}
}
