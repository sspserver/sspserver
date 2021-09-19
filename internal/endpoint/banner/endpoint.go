//
// @project GeniusRabbit rotator 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019, 2021
//

package banner

import (
	"net/http"

	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/private/templates"
)

type _endpoint struct {
	// Data source
	source endpoint.Sourcer

	// URL generator
	urlGen adsource.URLGenerator
}

func (e _endpoint) Version() string {
	return "v1"
}

func (e _endpoint) Handle(request *adsource.BidRequest) adsource.Responser {
	return e.source.Bid(request)
}

func (e _endpoint) Render(ctx *fasthttp.RequestCtx, response adsource.Responser) error {
	ctx.SetContentType("text/html")
	ctx.SetStatusCode(http.StatusOK)
	templates.WriteAdRenderComplexIFrame(ctx, response)
	return nil
}

func (e _endpoint) PrepareRequest(ctx *fasthttp.RequestCtx) error {
	return nil
}
