package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/jinzhu/gorm"
	ugrpc "github.com/sspserver/udetect/transport/grpc"
	uhttp "github.com/sspserver/udetect/transport/http"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"geniusrabbit.dev/sspserver/cmd/sspserver/appcontext"
	"geniusrabbit.dev/sspserver/cmd/sspserver/httpserver"
	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/data/adsourceaccessor"
	"geniusrabbit.dev/sspserver/internal/data/loaders"
	"geniusrabbit.dev/sspserver/internal/data/targetaccessor"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	_ "geniusrabbit.dev/sspserver/internal/endpoint/init"
	"geniusrabbit.dev/sspserver/internal/eventtraking/eventgenerator"
	"geniusrabbit.dev/sspserver/internal/eventtraking/eventstream"
	"geniusrabbit.dev/sspserver/internal/eventtraking/pixelgenerator"
	"geniusrabbit.dev/sspserver/internal/middleware"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
	"geniusrabbit.dev/sspserver/internal/notifications"
	"geniusrabbit.dev/sspserver/internal/personification"
	"geniusrabbit.dev/sspserver/internal/ssp"
	_ "geniusrabbit.dev/sspserver/internal/ssp/platform/init"
	"geniusrabbit.dev/sspserver/internal/urlgenerator"
	"geniusrabbit.dev/sspserver/private/templates"
)

const (
	eventsStreamName   = "events"
	userInfoStreamName = "user_info"
	winStreamName      = "wins"
)

var (
	datasource interface{}
	config     appcontext.Config
	commit     = ""
	appVersion = "develop"
)

func init() {
	// Init random
	rand.Seed(time.Now().UnixNano())
	fatalError(config.Load(), "config loading")

	// Init new logger object
	loggerObj, err := newLogger(config.IsDebug(), config.LogLevel, zap.Fields(
		zap.String("commit", commit),
		zap.String("version", appVersion),
	))
	fatalError(err, "configure logger")

	zap.ReplaceGlobals(loggerObj)

	// if config.IsDebug() {
	fmt.Println(config.String())
	// }
}

func main() {
	var (
		sspConf      = &config.AdServer.SSP
		adServerConf = &config.AdServer
		storageConf  = &config.AdServer.Storage
		logger       = zap.L()
		ctx, cancel  = context.WithCancel(context.Background())
	)

	defer cancel()

	// Init database connection
	urlData, err := url.Parse(adServerConf.Storage.Connection)
	fatalError(err, "connection URL parse")

	if urlData.Scheme == "fs" {
		datasource = adServerConf.Storage.Connection[3:]
	} else {
		datasource, err = gorm.Open(urlData.Scheme, adServerConf.Storage.Connection)
		fatalError(err, "open database")
	}

	// Register event streams
	{
		// Register events logger
		eventQueue := adServerConf.EventPipeline.EventQueue
		{
			err = notifications.Connection(ctx, eventsStreamName, eventQueue.Connection)
			fatalError(err, "connect to '"+eventQueue.Connection+"' topics")
		}

		// Register user info loop
		userQueue := adServerConf.EventPipeline.UserInfoQueue
		{
			err = notifications.Connection(ctx, userInfoStreamName, userQueue.Connection)
			fatalError(err, "connect to '"+userQueue.Connection+"' topics")
		}

		// Register wins info loop
		winQueue := adServerConf.EventPipeline.WinQueue
		{
			err = notifications.Connection(ctx, winStreamName, winQueue.Connection)
			fatalError(err, "connect to '"+winQueue.Connection+"' topics")
		}
	}

	// Event flow
	eventGenerator := eventgenerator.New(config.ServiceName)
	eventStream := eventstream.New(
		nc.PublisherByName(eventsStreamName),
		nc.PublisherByName(userInfoStreamName),
		eventGenerator,
	)

	// URL generator object
	urlGenerator := &urlgenerator.Generator{
		EventGenerator: eventGenerator,
		PixelGenerator: pixelgenerator.NewPixelGenerator(adServerConf.TrackerHost),
		CDNDomain:      adServerConf.CDNDomain,
		ClickPattern:   "/click?c={code}",
		DirectPattern:  "/direct?c={code}",
	}

	// Init templates
	templates.Debug = config.IsDebug()
	templates.ServiceDomain = ""
	templates.URLGen = urlGenerator

	// Run notification listener
	go nc.Listen(ctx)

	// Close notification processors
	defer nc.Close()

	// Create new source accessor which returns preorety iterator to every request
	sourceAccessor := adsourceaccessor.MustNewAccessor(
		loaders.SourceReloader(
			zap.L(),
			datasourceFrom(datasource, storageConf.Sources),
			companyGetter,
			eventStream,
			nc.PublisherByName(winStreamName),
		),
	)

	targetAccessor := targetaccessor.NewSimpleTargetAccessor(
		loaders.TargetReloader(datasourceFrom(datasource, storageConf.Zones)),
	)
	formatAccessor := types.NewSimpleFormatAccessorWithLoader(
		loaders.FormatLoader(datasourceFrom(datasource, storageConf.Formats)),
	)

	// Connect to the detector
	detector, err := newDetector(ctx, config.Person.Connect,
		config.Person.MaxConn, config.Person.RequestTimeout, config.Person.KeepAliveTimeout)
	fatalError(err, "Connect detector")

	// Configure advertisement service
	sspServer, err := ssp.NewServer(
		ssp.WithBaseSource(nil),
		ssp.WithSourceAccessor(sourceAccessor),
		ssp.WithTimeout(time.Duration(sspConf.RequestTimeout)*time.Millisecond),
		ssp.WithMaxParallelRequests(sspConf.MaxParallelRequests),
	)
	fatalError(err, "Create server")

	server, err := httpserver.NewServer(
		httpserver.WithDebugMode(config.IsDebug()),
		httpserver.WithServiceName(config.ServiceName),
		httpserver.WithZoneAccessor(targetAccessor),
		httpserver.WithFormatAccessor(formatAccessor),
		httpserver.WithAdvertisementSource(sspServer),
		httpserver.WithURLGenerator(urlGenerator),
		httpserver.WithEventStream(eventStream),
		httpserver.WithLogger(logger.With(zap.String("module", "httpserver"))),
		httpserver.WithSpy(newSpy(ctx, &config, detector)),
		httpserver.WithCustomHTTPServer(&fasthttp.Server{
			ReadBufferSize: 1 << 20,
			ReadTimeout:    time.Duration(config.Server.HTTP.ReadTimeout) * time.Millisecond,
			WriteTimeout:   time.Duration(config.Server.HTTP.WriteTimeout) * time.Millisecond,
		}),
	)
	fatalError(err)

	endpoint.Each(func(name string, factory endpoint.Factory) {
		server.RegisterEdpoint(name, func(source endpoint.Sourcer, options ...interface{}) (endpoint.Endpoint, error) {
			return factory(source, append(options, urlGenerator, eventStream)...)
		})
	})

	// RUN reloader loop
	timer := reloadLoop(time.Second*30, func() {
		logError(sourceAccessor.Reload(), "source accessor")
		logError(targetAccessor.Reload(), "target accessor")
		logError(formatAccessor.Reload(), "format accessor")
	})

	onKillApplication(func() {
		timer.Stop()
		nc.Close()
		server.Shutdown()
	})

	fmt.Println("Run HTTP server", config.Server.HTTP.Listen)
	fatalError(server.Listen(ctx, config.Server.HTTP.Listen))
}

func newDetector(ctx context.Context, url string, maxConn int, timeout, keepAlive time.Duration) (cli personification.Client, err error) {
	switch {
	case strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"):
		tr := uhttp.NewTransport(url,
			uhttp.WithGET,
			uhttp.WithMaxConn(maxConn),
			uhttp.WithTimeout(timeout, keepAlive))
		cli = personification.Connect(tr)
	case strings.HasPrefix(url, "grpc://"):
		tr, err := ugrpc.NewTransport(ctx, url, ugrpc.WithSecure(false))
		if err != nil {
			return nil, err
		}
		cli = personification.Connect(tr)
	default:
		cli = personification.DummyClient{}
	}
	return cli, err
}

func newSpy(ctx context.Context, cfg *appcontext.Config, detector personification.Client) middleware.Spy {
	signature := personification.Signeture{
		UUIDName:       cfg.Person.UUIDCookieName,
		SessidName:     cfg.Person.SessiCookiedName,
		SessidLifetime: cfg.Person.SessionLifetime,
		Detector:       detector,
	}
	return middleware.NewSpy(ctx, signature.Whois, signature.SignCookie)
}

func reloadLoop(duration time.Duration, fnk func()) *time.Ticker {
	fnk()

	timer := time.NewTicker(duration)
	go func() {
	loop:
		for {
			select {
			case <-timer.C:
				fnk()
			case <-nc.OnClose():
				break loop
			}
		}
	}()
	return timer
}

func onKillApplication(fnk func()) {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		<-gracefulStop
		fnk()
	}()
}

func fatalError(err error, message ...interface{}) {
	if err != nil {
		zap.L().Fatal(fmt.Sprint(message...), zap.Error(err))
	}
}

func logError(err error, message ...interface{}) error {
	if err != nil {
		zap.L().Error(fmt.Sprint(message...), zap.Error(err))
	}
	return err
}

func newLogger(debug bool, loglevel string, options ...zap.Option) (logger *zap.Logger, err error) {
	if debug {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return config.Build(options...)
	}

	var (
		level         zapcore.Level
		loggerEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		})
	)
	if err := level.UnmarshalText([]byte(loglevel)); err != nil {
		logger.Error("parse log level error", zap.Error(err))
	}
	core := zapcore.NewCore(loggerEncoder, os.Stdout, level)
	logger = zap.New(core, options...)

	return logger, nil
}

// fake companyGetter with unlimited balance
func companyGetter(id uint64) *models.Company {
	return &models.Company{
		ID:       id,
		Balance:  billing.MoneyFloat(9999999),
		MaxDaily: billing.MoneyFloat(9999999),
	}
}

func datasourceFrom(datasource interface{}, postfix string) interface{} {
	switch ds := datasource.(type) {
	case string:
		return filepath.Join(ds, postfix)
	}
	return datasource
}
