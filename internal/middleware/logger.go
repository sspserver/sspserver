package middleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Logger wraps http request context with the logger object
func Logger(basicLogger *zap.Logger, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return HTTPRequestID(func(ctx *fasthttp.RequestCtx) {
		requestID := ctx.UserValue("_id").(string)
		ctx.SetUserValue("_log", basicLogger.With(zap.String("id", requestID)))
		handler(ctx)
	})
}
