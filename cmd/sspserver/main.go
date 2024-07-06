package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/fasthttp/router"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/sspserver/udetect/transport/http"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"geniusrabbit.dev/adcorelib/admodels/types"
	"geniusrabbit.dev/adcorelib/adsource"
	"geniusrabbit.dev/adcorelib/adsource/openrtb"
	"geniusrabbit.dev/adcorelib/eventtraking/eventgenerator"
	"geniusrabbit.dev/adcorelib/eventtraking/eventstream"
	"geniusrabbit.dev/adcorelib/eventtraking/pixelgenerator"
	"geniusrabbit.dev/adcorelib/httpserver"
	"geniusrabbit.dev/adcorelib/httpserver/extensions/endpoint"
	"geniusrabbit.dev/adcorelib/httpserver/extensions/pixel"
	"geniusrabbit.dev/adcorelib/httpserver/extensions/trakeraction"
	"geniusrabbit.dev/adcorelib/httpserver/wrappers/httphandler"
	"geniusrabbit.dev/adcorelib/net/fasthttp/middleware"
	"geniusrabbit.dev/adcorelib/personification"
	"geniusrabbit.dev/adcorelib/simplepersondetector"
	"geniusrabbit.dev/adcorelib/storage/accessors/adsourceaccessor"
	"geniusrabbit.dev/adcorelib/storage/accessors/companyaccessor"
	"geniusrabbit.dev/adcorelib/storage/accessors/formataccessor"
	"geniusrabbit.dev/adcorelib/storage/accessors/zoneaccessor"
	"geniusrabbit.dev/adcorelib/urlgenerator"
	"geniusrabbit.dev/adcorelib/zlogger"

	"geniusrabbit.dev/sspserver/cmd/sspserver/appcontext"
	"geniusrabbit.dev/sspserver/cmd/sspserver/datainit"
	"geniusrabbit.dev/sspserver/internal/endpoint/direct"
	"geniusrabbit.dev/sspserver/internal/endpoint/dynamic"
	"geniusrabbit.dev/sspserver/internal/endpoint/proxy"
	"geniusrabbit.dev/sspserver/internal/stream"
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
)

func init() {
	// Init random
	rand.Seed(time.Now().UnixNano())
	fatalError(config.Load(), "config loading")

	// Init new logger object
	loggerObj, err := zlogger.New(config.ServiceName, config.LogEncoder,
		config.LogLevel, config.LogAddr, zap.Fields(
			zap.String("commit", buildCommit),
			zap.String("version", buildVersion),
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
		go nc.Listen(ctx)

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

	// Init format accessor
	formatAccessor := mustFormatAccessor(ctx, dataAccessor)
	ctx = formataccessor.WithContext(ctx, formatAccessor)

	// Init company data accessor
	companyDataAccessor, err := dataAccessor(ctx, "company")
	fatalError(err, "company data accessor")

	companyAccessor := companyaccessor.NewCompanyAccessor(companyDataAccessor)

	// Init source data accessor
	sourceDataAccessor, err := dataAccessor(ctx, "source")
	fatalError(err, "RTB source data accessor")

	// Init advertisement source accessor
	sourceAccessor, err := adsourceaccessor.NewAccessor(ctx,
		sourceDataAccessor, companyAccessor, openrtb.NewFactory())
	fatalError(err, "RTB source accessor")

	// Init target data accessor
	targetDataAccessor, err := dataAccessor(ctx, "zone")
	fatalError(err, "target(zone) data accessor")

	targetAccessor := zoneaccessor.NewZoneAccessor(targetDataAccessor, companyAccessor)

	// Configure advertisement source accessor
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
	if eventQueue := adServerConf.EventPipeline.EventQueue; eventQueue.Connection != "" {
		pub, err := stream.ConnectPublisher(ctx, eventQueue.Connection)
		fatalError(err, "connect to '"+eventQueue.Connection+"' topics")
		if config.IsDebug() {
			pub = stream.WrapPublisherWithLog(eventsStreamName, pub)
		}
		nc.Register(eventsStreamName, pub)
	} else {
		zap.L().Info("register new events dummy publisher")
		nc.Register(eventsStreamName, notificationMessageLog("events"))
	}

	// Register user info data stream
	if userQueue := adServerConf.EventPipeline.UserInfoQueue; userQueue.Connection != "" {
		pub, err := stream.ConnectPublisher(ctx, userQueue.Connection)
		fatalError(err, "connect to '"+userQueue.Connection+"' topics")
		if config.IsDebug() {
			pub = stream.WrapPublisherWithLog(userInfoStreamName, pub)
		}
		nc.Register(userInfoStreamName, pub)
	} else {
		zap.L().Info("register new userInfo dummy publisher")
		nc.Register(userInfoStreamName, notificationMessageLog("user_info"))
	}

	// Register wins info data stream
	if winQueue := adServerConf.EventPipeline.WinQueue; winQueue.Connection != "" {
		pub, err := stream.ConnectPublisher(ctx, winQueue.Connection)
		fatalError(err, "connect to '"+winQueue.Connection+"' topics")
		if config.IsDebug() {
			pub = stream.WrapPublisherWithLog(winStreamName, pub)
		}
		nc.Register(winStreamName, pub)
	} else {
		zap.L().Info("register new wins dummy publisher")
		nc.Register(winStreamName, notificationMessageLog("win"))
	}

	// Register adinfo data stream
	if adInfoQueue := adServerConf.EventPipeline.AdInfoQueue; adInfoQueue.Connection != "" {
		pub, err := stream.ConnectPublisher(ctx, adInfoQueue.Connection)
		fatalError(err, "connect to '"+adInfoQueue.Connection+"' topics")
		nc.Register(adInfoStreamName, pub)
	} else {
		zap.L().Info("register new adInfo dummy publisher")
		nc.Register(adInfoStreamName, notificationMessageLog("adinfo"))
	}
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
