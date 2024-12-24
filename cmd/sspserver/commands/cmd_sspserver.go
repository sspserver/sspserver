package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/demdxx/cloudregistry"
	"github.com/fasthttp/router"
	"github.com/geniusrabbit/adcorelib/adsource"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/httpserver"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/actiontracker"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/endpoint"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/pixeltracker"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/version"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
	openrtbsrc "github.com/geniusrabbit/adsource-openrtb"
	"github.com/geniusrabbit/adstdendpoints/direct"
	"github.com/geniusrabbit/adstdendpoints/dynamic"
	"github.com/geniusrabbit/adstdendpoints/proxy"
	"github.com/geniusrabbit/adstorage"
	"github.com/geniusrabbit/adstorage/accessors/formataccessor"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/sspserver/sspserver/cmd/sspserver/appcontext"
	"github.com/sspserver/sspserver/cmd/sspserver/datainit"
	"github.com/sspserver/sspserver/internal/events/rtbevents"
	"github.com/sspserver/sspserver/internal/netdriver"
)

type sspserverConfig struct {
	ServiceName string `json:"service_name" yaml:"service_name" env:"SERVICE_NAME" default:"adserver"`

	LogLevel string `json:"log_level" default:"error" env:"LOG_LEVEL"`

	// Server config
	Server appcontext.ServerConfig `field:"server" json:"server" yaml:"server"`

	// Configuration of Advertisement server
	AdServer appcontext.AdServerConfig `field:"adserver" yaml:"adserver" json:"adserver"`

	// Person data extraction service
	Person appcontext.PersonConfig `field:"person" yaml:"person" json:"person"`
}

func (cfg *sspserverConfig) IsDebug() bool {
	return strings.EqualFold(cfg.LogLevel, "debug")
}

var SSPServerCommand = &Command[sspserverConfig]{
	Name:     "sspserver",
	HelpDesc: "Run SSP server",
	Exec:     sspServerCommand,
}

func init() {
	Commands = append(Commands, SSPServerCommand)
}

func sspServerCommand(ctx context.Context, args []string, config *sspserverConfig, numberOfAdServers *cloudregistry.SyncUInt64Value) error {
	type eventType = rtbevents.Event

	var (
		eventAllocator = rtbevents.AllocateEvent
		logger         = ctxlogger.Get(ctx)
		adServerConf   = &config.AdServer
		urlGenerator   adtype.URLGenerator
		eventStream    eventstream.Stream
	)

	// Register advertisement data accessor
	adstorage.Register("fs", adstorage.FSDataAccessor[datainit.Account])
	// Register all supported database schemas for all supported dialects
	adstorage.RegisterAllSchemas[datainit.Account]()

	// Connect to advertisement data source
	storageDataAccessor, err := adstorage.ConnectAllAccessors(ctx,
		adServerConf.Storage.Connection,
		datainit.AdModelAccount,
	)
	if err != nil {
		return errors.Wrap(err, "connect to advertisement data")
	}

	// Register event streams
	{
		configureEventPipeline(ctx, adServerConf, config.IsDebug())

		// Run notification listener
		go func() {
			if err := nc.Listen(ctx); err != nil {
				logger.Error("notification listener", zap.Error(err))
			}
		}()

		// Close notification processors
		defer func() { _ = nc.Close() }()
	}

	// Init events generator
	ctx, urlGenerator, eventStream, err = initEventsGenerator(ctx,
		config.ServiceName, config.Server.Hostname, *adServerConf, eventAllocator)
	if err != nil {
		return errors.Wrap(err, "init events generator")
	}

	// Init side modules
	datainit.Initialize(config.IsDebug(), urlGenerator)

	// =================== Init formats ========================
	// Format is a specific type of advertisement
	// ========================================================

	// Init format accessor
	formatAccessor, err := storageDataAccessor.Formats()
	if err != nil {
		return errors.Wrap(err, "format accessor")
	}
	ctx = formataccessor.WithContext(ctx, formatAccessor)

	// =================== Init Ads sources ===================
	// Ad source is a specific driver for advertisement access
	// ========================================================

	// Init RTB source accessor
	sourceAccessor, err := storageDataAccessor.Sources(openrtbsrc.NewFactory(netdriver.NewDriver))
	if err != nil {
		return errors.Wrap(err, "RTB source accessor")
	}

	// Init applicaion data accessor (sites/mobile apps/desktop apps/console/etc)
	appAccessor, err := storageDataAccessor.Apps()
	if err != nil {
		return errors.Wrap(err, "app accessor")
	}

	// Init target data accessor (targeting zones where advertisement will be shown)
	targetAccessor, err := storageDataAccessor.Zones()
	if err != nil {
		return errors.Wrap(err, "target accessor")
	}

	// Configure advertisement source accessor
	adsourceWrapper, err := adsource.NewMultisourceWrapper(
		adsource.WithSourceAccessor(sourceAccessor),
		adsource.WithTimeout(time.Duration(adServerConf.AdSource.RequestTimeout)*time.Millisecond),
		adsource.WithMaxParallelRequests(adServerConf.AdSource.MaxParallelRequests),
	)
	if err != nil {
		return errors.Wrap(err, "create adsource accessor")
	}

	// =================== Init HTTP server ===================
	// HTTP server is a main entry point for external ad requests
	// ========================================================

	// HTTP wrapper for extended handlers
	httpHandlerWrapper := httphandler.NewHTTPHandlerWrapper(
		nil, // func(*fasthttp.RequestCtx) {},
		func(*fasthttp.RequestCtx) context.Context { return ctx },
		logger.With(zap.String("module", "httpserver")),
	)

	// Init personification client
	_, spyMiddleware, err := initPersonification(ctx, &config.Person)
	if err != nil {
		return errors.Wrap(err, "init personification")
	}

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
			router.GET("/info", infoHandler)
			router.GET("/health", func(ctx *fasthttp.RequestCtx) { ctx.SetStatusCode(fasthttp.StatusOK) })
			router.ServeFiles("/public/{filepath:*}", "/public")
		}),
		httpserver.WithExtensions(
			&version.Extension{},
			// Pixel action tracker
			pixeltracker.NewExtension(
				pixeltracker.WithEventStream[*eventType](eventStream),
				pixeltracker.WithHTTPHandlerWrapper[*eventType](httpHandlerWrapper),
				pixeltracker.WithEventAllocator(eventAllocator),
			),
			// Postback actions (click, direct, win, etc)
			actiontracker.NewExtension(
				actiontracker.WithEventStream[*eventType](eventStream),
				actiontracker.WithHTTPHandlerWrapper[*eventType](httpHandlerWrapper),
				actiontracker.WithURLGenerator[*eventType](urlGenerator),
				actiontracker.WithEventAllocator(eventAllocator),
			),
			// Register HTTP endpoints extensions
			endpoint.NewExtension(
				endpoint.WithAdvertisementSource(adsourceWrapper),
				endpoint.WithHTTPHandlerWrapper(httpHandlerWrapper),
				endpoint.WithFormatAccessor(formatAccessor),
				endpoint.WithAppAccessor(appAccessor),
				endpoint.WithZoneAccessor(targetAccessor),
				endpoint.WithSpy(spyMiddleware),
				endpoint.WithSendpoints(
					direct.New(formatAccessor, adServerConf.Logic.Direct.DefaultURL),
					dynamic.New(urlGenerator),
					proxy.New(),
				),
			),
		),
	)
	if err != nil {
		return errors.Wrap(err, "new HTTP server")
	}

	fmt.Println("Run HTTP server", config.Server.HTTP.Listen)
	return server.Listen(ctx, config.Server.HTTP.Listen)
}
