package middleware

import (
	"github.com/valyala/fasthttp"
)

// HTTPRequestID collects the request ID from X-Request-Id header or Cloudflare analog Cf-Ray
func HTTPRequestID(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var (
			headers   = ctx.Request.Header
			requestID = string(headers.Peek("X-Request-Id"))
		)
		if requestID == "" {
			requestID = string(headers.Peek("Cf-Ray"))
		}
		ctx.SetUserValue("_id", requestID)
		handler(ctx)
	}
}
