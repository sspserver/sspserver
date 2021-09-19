package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/geniusrabbit/notificationcenter"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
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
	"geniusrabbit.dev/sspserver/internal/eventgenerator"
	"geniusrabbit.dev/sspserver/internal/eventstream"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
	"geniusrabbit.dev/sspserver/internal/notifications"
	"geniusrabbit.dev/sspserver/internal/ssp"
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

	if config.IsDebug() {
		fmt.Println(config.String())
	}
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
		notificationcenter.PublisherByName(eventsStreamName),
		notificationcenter.PublisherByName(userInfoStreamName),
		eventGenerator,
	)

	// URL generator object
	urlGenerator := &urlgenerator.Generator{
		EventGenerator: eventGenerator,
		PixelGenerator: nil, //client.NewPixelGenerator(adServerConf.TrackerHost),
		CDNDomain:      adServerConf.CDNDomain,
		ClickPattern:   "/click?c={code}",
		DirectPattern:  "/direct?c={code}",
	}

	// Init templates
	templates.Debug = config.IsDebug()
	templates.ServiceDomain = ""
	templates.URLGen = urlGenerator

	// Run notification listener
	go notificationcenter.Listen(ctx)

	// Close notification processors
	defer notificationcenter.Close()

	// Create new source accessor which returns preorety iterator to every request
	sourceAccessor := adsourceaccessor.MustNewAccessor(
		loaders.SourceReloader(
			zap.L(),
			datasourceFrom(datasource, storageConf.Sources),
			func(id uint64) *models.Company {
				return &models.Company{
					ID:       id,
					Balance:  billing.MoneyFloat(9999999),
					MaxDaily: billing.MoneyFloat(9999999),
				}
			},
			eventStream,
		),
	)

	targetAccessor := targetaccessor.NewSimpleTargetAccessor(
		loaders.TargetReloader(datasourceFrom(datasource, storageConf.Zones)),
	)
	formatAccessor := types.NewSimpleFormatAccessorWithLoader(
		loaders.FormatLoader(datasourceFrom(datasource, storageConf.Formats)),
	)

	// Configure advertisement service
	sspServer, err := ssp.NewServer(
		ssp.WithBaseSource(nil),
		ssp.WithSourceAccessor(sourceAccessor),
		ssp.WithTimeout(time.Duration(sspConf.RequestTimeout)*time.Millisecond),
		ssp.WithMaxParallelRequests(sspConf.MaxParallelRequests),
	)
	fatalError(err)

	server, err := httpserver.NewServer(
		httpserver.WithDebugMode(config.IsDebug()),
		httpserver.WithServiceName(config.ServiceName),
		httpserver.WithZoneAccessor(targetAccessor),
		httpserver.WithFormatAccessor(formatAccessor),
		httpserver.WithAdvertisementSource(sspServer),
		httpserver.WithURLGenerator(urlGenerator),
		httpserver.WithEventStream(eventStream),
		httpserver.WithLogger(logger.With(zap.String("module", "httpserver"))),
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
		notificationcenter.Close()
		server.Shutdown()
	})

	fmt.Println("Run HTTP server", config.Server.HTTP.Listen)
	fatalError(server.Listen(ctx, config.Server.HTTP.Listen))
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
			case <-notificationcenter.OnClose():
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
		log.Fatal(fmt.Sprint(append(message, err)...))
	}
}

func logError(err error, message ...interface{}) error {
	if err != nil {
		log.Error(append(message, err)...)
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

func datasourceFrom(datasource interface{}, postfix string) interface{} {
	switch ds := datasource.(type) {
	case string:
		return filepath.Join(ds, postfix)
	}
	return datasource
}
