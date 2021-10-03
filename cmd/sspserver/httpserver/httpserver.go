package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/demdxx/gocast"
	fastp "github.com/flf2ko/fasthttp-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sspserver/udetect"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/debugtool"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/internal/eventtraking/events"
	"geniusrabbit.dev/sspserver/internal/eventtraking/eventstream"
	"geniusrabbit.dev/sspserver/internal/gtracing"
	"geniusrabbit.dev/sspserver/internal/middleware"
	"geniusrabbit.dev/sspserver/internal/models/types"
	"geniusrabbit.dev/sspserver/internal/personification"
	"geniusrabbit.dev/sspserver/internal/ssp/openlatency"
	"geniusrabbit.dev/sspserver/internal/ssp/platform"
)

type (
	personalizedHandler func(person personification.Person, ctx *fasthttp.RequestCtx)
	customRouterFnk     func(router *fasthttprouter.Router)
)

// Errors list...
var (
	ErrInvalidZoneAccessor          = errors.New("[HTTPServer] invalid zone accessor")
	ErrInvalidSourceOfAdvertisement = errors.New("[HTTPServer] invalid source of advertisement")
	ErrInvalidURLGenerator          = errors.New("[HTTPServer] invalid URL generator")
	ErrInvalidEventstream           = errors.New("[HTTPServer] invalid event stream")
	ErrInvalidSpyWrapper            = errors.New("[HTTPServer] invalid spy wrapper")
)

type sourceAccessor interface {
	Sources() adsource.SourceAccessor
}

type metricsAccessor interface {
	Metrics() *openlatency.MetricsInfo
}

// Server implements basic HTTP infostructure and routing
type Server struct {
	// Debug mode of the server
	debug bool

	// service name
	serviceName string

	// zoneAccessor to the paticular target by ID
	zoneAccessor targetAccessor

	// formatAccessor to the suitable formats by size and types
	formatAccessor types.FormatsAccessor

	// Source of the advertisement
	source Source

	// net connection to the listen port
	httpConnection net.Listener

	// httpServer object
	httpServer *fasthttp.Server

	// prepareRequestCtx params
	prepareRequestCtx func(ctx *fasthttp.RequestCtx)

	// URL generator used for initialisation of submodules
	urlGenerator adsource.URLGenerator

	// Event stream interface sends data into the queue
	eventStream eventstream.Stream

	// Endpoint list of factories
	endpoints map[string]endpoint.Factory

	// Custom router function
	customRouter customRouterFnk

	// Logger base object
	logger *zap.Logger

	// Wrapper with user grabbing information
	spyMiddleware middleware.Spy

	// Metrics prepared methods
	metrics Metrics

	// tracer interface helping to trace application
	tracer opentracing.Tracer

	// Setup into the 1 when server is shuting down
	shutdownMode uint32
}

// NewServer http server object
func NewServer(options ...Option) (*Server, error) {
	srv := &Server{}
	for _, opt := range options {
		opt(srv)
	}

	switch {
	case srv.zoneAccessor == nil:
		return nil, ErrInvalidZoneAccessor
	case srv.source == nil:
		return nil, ErrInvalidSourceOfAdvertisement
	case srv.urlGenerator == nil:
		return nil, ErrInvalidURLGenerator
	case srv.eventStream == nil:
		return nil, ErrInvalidEventstream
	case srv.spyMiddleware == nil:
		return nil, ErrInvalidSpyWrapper
	case srv.logger == nil:
		srv.logger = zap.L().With(zap.String("module", "httpserver"))
	}

	if err := srv.initTracer(); err != nil {
		return nil, err
	}

	return srv, nil
}

// Listen server address
func (srv *Server) Listen(ctx context.Context, address string) (err error) {
	if srv.httpServer == nil {
		srv.httpServer = &fasthttp.Server{ReadBufferSize: 1 << 20}
	}

	p := fastp.NewPrometheus("fasthttp")
	srv.httpServer.Handler = p.WrapHandler(srv.newRouter(ctx))

	srv.httpConnection, err = net.Listen("tcp4", address)
	if err != nil {
		return err
	}

	return srv.httpServer.Serve(srv.httpConnection)
}

// RegisterEdpoint factory object
func (srv *Server) RegisterEdpoint(name string, factory endpoint.Factory) {
	if srv.endpoints == nil {
		srv.endpoints = map[string]endpoint.Factory{}
	}
	srv.endpoints[name] = factory
}

// Shutdown server gracefully
func (srv *Server) Shutdown() {
	srv.logger.Debug("Shutdown the HTTP server", zap.String("method", "Shutdown"))

	if !atomic.CompareAndSwapUint32(&srv.shutdownMode, 0, 1) {
		return
	}

	if srv.httpConnection != nil {
		srv.httpConnection.Close()
		srv.httpConnection = nil
	}
}

// IsShutdownMode on or off
func (srv *Server) IsShutdownMode() bool {
	return atomic.LoadUint32(&srv.shutdownMode) == 1
}

func (srv *Server) newRouter(ctx context.Context) *fasthttprouter.Router {
	srv.logger.Debug("Initialise the router", zap.String("method", "newRouter"))

	router := fasthttprouter.New()
	router.PanicHandler = srv.panicCallback
	routeWrapper := traceRouterWrapper{router: router, tracer: srv.tracer}

	options := []interface{}{
		srv.zoneAccessor,
		srv.urlGenerator,
		srv.formatAccessor,
		srv.tracer,
	}

	// * Init request endpoints of direct traffic
	if srv.endpoints != nil {
		for name, factory := range srv.endpoints {
			if endpoint, err := factory(srv.source, options...); err == nil {
				ver := endpoint.Version()
				if ver != "" {
					ver = "/" + ver
				}
				srv.registerEndpointHandler(ctx, ver, name, endpoint, routeWrapper)
			} else {
				panic(err)
			}
		}
	}

	platform.Each(func(protocol string, platFact platform.Factory) {
		router.GET("/protocol/"+protocol+"/info",
			srv.metricsWrapper("ssp."+protocol+".info", func(p personification.Person, ctx *fasthttp.RequestCtx) {
				ctx.Response.SetStatusCode(http.StatusOK)
				ctx.Response.Header.SetContentType("text/json")
				json.NewEncoder(ctx.Response.BodyWriter()).Encode(platFact.Info())
			}),
		)
	})

	// Get source status
	router.GET("/dsp/:id/status",
		srv.metricsWrapper("dsp.status", func(p personification.Person, ctx *fasthttp.RequestCtx) {
			var (
				id                 = gocast.ToUint64(ctx.UserValue("id"))
				source interface{} = srv.source
			)
			if id > 0 {
				if srcacc := srv.source.(sourceAccessor); srcacc != nil {
					source = srcacc.Sources().SourceByID(id)
				}
			}
			if source == nil {
				ctx.Response.SetStatusCode(http.StatusNotFound)
				ctx.Response.Header.SetContentType("text/json")
				ctx.Response.BodyWriter().Write([]byte(`{"status":"not_found"}`))
			} else if metrics := source.(metricsAccessor); metrics != nil {
				ctx.Response.SetStatusCode(http.StatusServiceUnavailable)
				ctx.Response.Header.SetContentType("text/json")
				json.NewEncoder(ctx.Response.BodyWriter()).Encode(metrics.Metrics())
			} else {
				ctx.Response.SetStatusCode(http.StatusServiceUnavailable)
				ctx.Response.Header.SetContentType("text/json")
				ctx.Response.BodyWriter().Write([]byte(`{"status":"unsupported"}`))
			}
		}),
	)

	// Click/Direct links
	router.GET(srv.urlGenerator.ClickRouterURL(),
		srv.metricsWrapper("ssp.click", srv.eventHandler(ctx, events.Click)),
	)
	router.GET(srv.urlGenerator.DirectRouterURL(),
		srv.metricsWrapper("ssp.direct", srv.eventHandler(ctx, events.Direct)),
	)

	router.GET("/healthcheck", srv.spyMiddleware(srv.healthCheck))
	router.GET("/check", srv.check)

	if srv.customRouter != nil {
		srv.customRouter(router)
	}

	return router
}

// eventHandler by ad. This method works only for BaseSource
func (srv *Server) eventHandler(ctx context.Context, eventName events.Type) func(person personification.Person, ctx *fasthttp.RequestCtx) {
	handlerCode := "ssp." + eventName.String()
	return func(person personification.Person, ctx *fasthttp.RequestCtx) {
		defer debugtool.Trace()

		var (
			clickData = ctx.QueryArgs().Peek("c")
			event     events.Event
			err       = event.Unpack(clickData, func(code events.Code) events.Code {
				return code.URLDecode().Decompress()
			})
			span, _ = gtracing.StartSpanFromFastContext(ctx, handlerCode)
		)

		if span != nil {
			defer span.Finish()
		}

		if err != nil || eventName != event.Event {
			// TODO: Log event as unappropriate
			// s.Render.DirectBidResponse(nil, ctx)
			// return
		}

		ctx.Redirect(event.URL, http.StatusFound)
		event.SetDateTime(time.Now())
		srv.eventStream.SendEvent(ctx, &event)
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Handlers
///////////////////////////////////////////////////////////////////////////////

func (srv *Server) handlerEndpoint(ctx context.Context, endpointName string, endpoint endpointer) personalizedHandler {
	log := srv.logger.With(
		zap.String("endpoint", endpointName),
		zap.String("method", "handlerEndpoint"),
	)

	log.Debug("Initialise the new endpoint")
	traceName := "ssp." + endpointName

	return func(person personification.Person, req *fasthttp.RequestCtx) {
		defer debugtool.Trace()

		var (
			response adsource.Responser
			request  = srv.requestByHTTPRequest(ctx, person, req)
			span, _  = gtracing.StartSpanFromFastContext(req, traceName)
		)

		if span != nil {
			ext.Component.Set(span, endpointName)
			request.Context = opentracing.ContextWithSpan(request.Context, span)
			defer span.Finish()
		}

		if response = endpoint.Handle(request); response == nil {
			fmt.Println(">>>> EMPTY RESPONSE")
			response = adsource.NewEmptyResponse(request, nil)
		}

		fmt.Println(">>>> RESPONSE", response)

		srv.logError(endpoint.Render(req, response))
		srv.source.ProcessResponse(response)
	}
}

func (srv *Server) healthCheck(p personification.Person, ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.Header.SetContentType("text/json")
	headers := strings.TrimSpace(ctx.Request.Header.String())

	json.NewEncoder(ctx.Response.BodyWriter()).Encode(&struct {
		Status  string      `json:"status"`
		Headers interface{} `json:"headers"`
		Info    interface{} `json:"info"`
	}{
		Status:  "ok",
		Headers: strings.Split(headers, "\r\n"),
		Info:    p.UserInfo(),
	})
}

func (srv *Server) check(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.Header.SetContentType("text/json")
	fmt.Fprint(ctx.Response.BodyWriter(), `{"status":"ok"}`)
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func (srv *Server) requestByHTTPRequest(ctx context.Context, person personification.Person, req *fasthttp.RequestCtx) *adsource.BidRequest {
	var (
		spotInfo  = strings.Split(req.UserValue("zone").(string), ".")
		zoneID, _ = strconv.ParseUint(spotInfo[0], 10, 64)
		request   = NewRequestFor(
			ctx,
			req.UserValue("_id").(string),
			srv.zoneAccessor.TargetByID(zoneID),
			person,
			NewRequestOptions(req),
			srv.formatAccessor,
		)
	)
	if request != nil {
		srv.prepareRequest(request)
	}
	return request
}

func (srv *Server) prepareRequest(request *adsource.BidRequest) {
	// BigBrother, tell me! Who is it?
	var (
		query    = request.RequestCtx.QueryArgs()
		keywords = peekOneFromQuery(query, "keywords", "keyword", "kw")
	)

	if request.Person != nil {
		// Fill user info
		// var ui = request.Person.UserInfo()
		request.User = &adsource.User{
			// ID:            ui.UUID(),
			// SessionID:     ui.SessionID(),
			// FingerPrintID: ui.Fingerprint(),
			// ETag:          ui.ETag(),
			Keywords: keywords,
			// Geo:           ui.GeoInfo(),
		}
		// request.Device = ui.DeviceInfo()

		// if ui != nil && ui.User != nil {
		// 	request.User.AgeStart = ui.User.AgeStart
		// 	request.User.AgeEnd = ui.User.AgeEnd
		// }
	} else {
		request.User = &adsource.User{
			FingerPrintID: "",
			Keywords:      keywords,
		}
	}

	// Fill GEO info
	if request.User.Geo == nil {
		request.User.Geo = &udetect.Geo{
			IP: request.RequestCtx.RemoteIP(),
		}
	}
}

func (srv *Server) registerEndpointHandler(ctx context.Context, urlPrefix, name string, endpoint endpointer, router traceRouterWrapper) {
	handler := srv.handlerEndpoint(ctx, name, endpoint)
	sspFirstHandler := srv.metricsWrapper("ssp."+name, handler)
	sspSecondHandler := middleware.Logger(srv.logger, func(ctx *fasthttp.RequestCtx) {
		if srv.prepareRequestCtx != nil {
			srv.prepareRequestCtx(ctx)
		}
		endpoint.PrepareRequest(ctx)
		sspFirstHandler(ctx)
	})

	for postfix, methods := range getRoutes(endpoint) {
		for _, method := range methods {
			router.Handle(method, urlPrefix+"/"+name+"/"+postfix, sspSecondHandler)
		}
	}
}

func (srv *Server) prepareHandler(name string, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	handlerWrapper := middleware.CollectSimpleMetrics(name, handler)
	if srv.prepareRequestCtx == nil {
		middleware.Logger(srv.logger, handlerWrapper)
	}
	return middleware.Logger(srv.logger, func(ctx *fasthttp.RequestCtx) {
		srv.prepareRequestCtx(ctx)
		handlerWrapper(ctx)
	})
}

func (srv *Server) panicCallback(ctx *fasthttp.RequestCtx, rcv interface{}) {
	srv.logError(fmt.Errorf("server panic: %+v\n%s", rcv, debug.Stack()))
	if srv.debug {
		msg := fmt.Sprintf("server painc: %v\n", rcv)
		ctx.Write([]byte(msg))
		ctx.Write(debug.Stack())
	}
	ctx.SetStatusCode(http.StatusInternalServerError)
}

func (srv *Server) initTracer() (err error) {
	if srv.tracer, err = gtracing.InitTracer(srv.serviceName, srv.logger); err != nil {
		srv.logError(err)
	}
	return nil
}

func (srv *Server) metricsWrapper(metric string, next func(p personification.Person, ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler {
	return srv.spyMiddleware(func(p personification.Person, ctx *fasthttp.RequestCtx) {
		// TODO detect country and timings
		next(p, ctx)
	})
}

func (srv *Server) logError(err error) error {
	if err != nil {
		srv.logger.Error("", zap.Error(err))
	}
	return err
}
