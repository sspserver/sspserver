//go:build htmltemplates
// +build htmltemplates

//
// @project GeniusRabbit sspserver 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2022
//

package proxy

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/sspserver/sspserver/internal/endpoint"
	"github.com/sspserver/sspserver/private/templates"
)

type _endpoint struct{}

func New() *_endpoint { return &_endpoint{} }

func (e *_endpoint) Codename() string {
	return "proxy"
}

func (e *_endpoint) Handle(source endpoint.Source, request *adtype.BidRequest) adtype.Responser {
	request.RequestCtx.SetContentType("text/html; charset=UTF-8")
	templates.WriteAdRenderDinamicProxyBanner(request.RequestCtx, request)
	return nil
}
