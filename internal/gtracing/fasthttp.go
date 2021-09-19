package gtracing

import (
	"fmt"
	"net/http"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type otErrorWrapper struct {
	logger *zap.Logger
}

func (l otErrorWrapper) Error(s string) {
	l.logger.Error(s)
}

func (l otErrorWrapper) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

type fasthttpHeaderWrapper struct {
	headers *fasthttp.RequestHeader
}

// Set conforms to the TextMapWriter interface.
func (c fasthttpHeaderWrapper) Set(key, val string) {
	c.headers.Set(key, val)
}

// ForeachKey conforms to the TextMapReader interface.
func (c fasthttpHeaderWrapper) ForeachKey(handler func(key, val string) error) (err error) {
	c.headers.VisitAll(func(k, v []byte) {
		if err != nil {
			return
		}
		err = handler(string(k), string(v))
	})
	return nil
}

// InitTracer object
func InitTracer(serviceName string, logger *zap.Logger) (tracer opentracing.Tracer, _ error) {
	if os.Getenv("JAEGER_ENDPOINT") == "" && os.Getenv("JAEGER_AGENT_HOST") == "" {
		// Skip the tracer initialisation
		return nil, nil
	}

	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("cannot parse Jaeger env vars: %s", err.Error())
	}

	cfg.ServiceName = serviceName

	tracer, _, err = cfg.NewTracer(
		config.Logger(otErrorWrapper{logger.With(zap.String("module", "opentracing"))}),
	)

	if err != nil {
		return nil, fmt.Errorf("jager server connect error: %s", err.Error())
	}

	opentracing.SetGlobalTracer(tracer)
	return nil, nil
}

// FastHTTPTraceWrapper must be the first in the chain of wrappers if we wan't to test them too
func FastHTTPTraceWrapper(tracer opentracing.Tracer, h fasthttp.RequestHandler) fasthttp.RequestHandler {
	if tracer == nil {
		// The tracer is not inited
		return h
	}

	return func(rctx *fasthttp.RequestCtx) {
		ctx, _ := tracer.Extract(opentracing.HTTPHeaders, fasthttpHeaderWrapper{&rctx.Request.Header})
		sp := tracer.StartSpan(fmt.Sprintf("HTTP %s", rctx.Method()), ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, string(rctx.Method()))
		ext.HTTPUrl.Set(sp, string(rctx.RequestURI()))
		ext.Component.Set(sp, "net/http")

		FastContextWithSpan(rctx, sp)

		h(rctx)

		ext.HTTPStatusCode.Set(sp, uint16(rctx.Response.StatusCode()))
		if rctx.Response.StatusCode() >= http.StatusInternalServerError {
			ext.Error.Set(sp, true)
		}
		sp.Finish()
	}
}
