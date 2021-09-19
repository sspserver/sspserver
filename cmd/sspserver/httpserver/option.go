package httpserver

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/eventstream"
	"geniusrabbit.dev/sspserver/internal/middleware"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// Option type
type Option func(srv *Server)

// WithServiceName which represents the internal name of the service
func WithServiceName(name string) Option {
	return func(srv *Server) {
		srv.serviceName = name
	}
}

// WithDebugMode of the server
func WithDebugMode(debug bool) Option {
	return func(srv *Server) {
		srv.debug = debug
	}
}

// WithCustomHTTPServer setup customly configured server
func WithCustomHTTPServer(server *fasthttp.Server) Option {
	return func(srv *Server) {
		srv.httpServer = server
	}
}

// WithZoneAccessor server option
func WithZoneAccessor(accessor interface{}) Option {
	return func(srv *Server) {
		switch acc := accessor.(type) {
		case targetAccessor:
			srv.zoneAccessor = acc
		default:
			panic(fmt.Sprintf("invalid zone accessor %T, have to implement targetAccessor interface", accessor))
		}
	}
}

// WithAdvertisementSource accessor
func WithAdvertisementSource(source Source) Option {
	return func(srv *Server) {
		srv.source = source
	}
}

// WithFormatAccessor interface
func WithFormatAccessor(formatAccessor types.FormatsAccessor) Option {
	return func(srv *Server) {
		srv.formatAccessor = formatAccessor
	}
}

// WithURLGenerator interface
func WithURLGenerator(urlGenerator adsource.URLGenerator) Option {
	return func(srv *Server) {
		srv.urlGenerator = urlGenerator
	}
}

// WithEventStream interface
func WithEventStream(eventStream eventstream.Stream) Option {
	return func(srv *Server) {
		srv.eventStream = eventStream
	}
}

// WithCustomRouter registrator
func WithCustomRouter(fouterFnk customRouterFnk) Option {
	return func(srv *Server) {
		srv.customRouter = fouterFnk
	}
}

// WithLogger interface
func WithLogger(logger *zap.Logger) Option {
	return func(srv *Server) {
		srv.logger = logger
	}
}

// WithPrepareRequest callback
func WithPrepareRequest(prepare func(ctx *fasthttp.RequestCtx)) Option {
	return func(srv *Server) {
		srv.prepareRequestCtx = prepare
	}
}

// WithSpy middleware wrapper
func WithSpy(spy middleware.Spy) Option {
	return func(srv *Server) {
		srv.spyMiddleware = spy
	}
}
