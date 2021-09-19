//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package proxy

import (
	"net/http"

	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/private/templates"
)

func init() {
	endpoint.Register("proxy", func(source endpoint.Sourcer, options ...interface{}) (endpoint.Endpoint, error) {
		return &_endpoint{source: source}, nil
	})
}

type _endpoint struct {
	source endpoint.Sourcer
}

func (e _endpoint) Version() string {
	return "v1"
}

func (e _endpoint) Handle(request *adsource.BidRequest) adsource.Responser {
	return &adsource.ResponseEmpty{Req: request}
}

func (e _endpoint) Render(ctx *fasthttp.RequestCtx, response adsource.Responser) error {
	ctx.SetContentType("text/html")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody([]byte(templates.AdRenderDinamicProxyBanner(response.Request())))
	return nil
}

func (e _endpoint) PrepareRequest(ctx *fasthttp.RequestCtx) (err error) {
	return
}
