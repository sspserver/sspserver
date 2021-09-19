package endpoint

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

// Sourcer interface of data accessor
type Sourcer interface {
	// Bid request for standart system filter
	Bid(request *adsource.BidRequest) adsource.Responser
}

// Endpoint describes advertisement accessor
type Endpoint interface {
	Version() string
	Handle(request *adsource.BidRequest) adsource.Responser
	Render(ctx *fasthttp.RequestCtx, response adsource.Responser) error
	PrepareRequest(ctx *fasthttp.RequestCtx) error
}

// Routing returns the acceptable royting postfix
type Routing interface {
	Routing() map[string][]string
}

// Factory function
type Factory func(source Sourcer, options ...interface{}) (Endpoint, error)

var factories = map[string]Factory{}

// Register endpoint factory
func Register(name string, f Factory) {
	if _, ok := factories[name]; ok {
		panic(fmt.Errorf("factory [%s] already registered", name))
	}
	factories[name] = f
}

// Each every factory
func Each(f func(name string, f Factory)) {
	for name, fk := range factories {
		f(name, fk)
	}
}
