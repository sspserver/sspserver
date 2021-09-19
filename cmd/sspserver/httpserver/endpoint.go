package httpserver

import (
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

// the private description of endpoint driver
type endpointer interface {
	// PrepareRequest before send to the handler
	PrepareRequest(ctx *fasthttp.RequestCtx) error

	// Handle the advertisement direct request
	Handle(request *adsource.BidRequest) adsource.Responser

	// Render response
	Render(ctx *fasthttp.RequestCtx, response adsource.Responser) error
}

type endpointRouting interface {
	Routing() map[string][]string
}
