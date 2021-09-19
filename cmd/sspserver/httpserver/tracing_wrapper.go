package httpserver

import (
	"github.com/buaazp/fasthttprouter"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/gtracing"
)

type traceRouterWrapper struct {
	router *fasthttprouter.Router
	tracer opentracing.Tracer
}

func (w traceRouterWrapper) GET(path string, h fasthttp.RequestHandler) {
	w.router.GET(path, gtracing.FastHTTPTraceWrapper(w.tracer, h))
}

func (w traceRouterWrapper) POST(path string, h fasthttp.RequestHandler) {
	w.router.POST(path, gtracing.FastHTTPTraceWrapper(w.tracer, h))
}

func (w traceRouterWrapper) Handle(method, path string, h fasthttp.RequestHandler) {
	w.router.Handle(method, path, gtracing.FastHTTPTraceWrapper(w.tracer, h))
}
