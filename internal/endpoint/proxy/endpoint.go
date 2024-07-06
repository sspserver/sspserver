//go:build htmltemplates
// +build htmltemplates

//
// @project GeniusRabbit rotator 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2022
//

package proxy

import (
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/private/templates"
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
