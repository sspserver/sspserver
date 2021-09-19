package middleware

import "github.com/valyala/fasthttp"

// RequestID passing internal request ID
func RequestID(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		handler(ctx)
	})
}
