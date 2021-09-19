//
// @project GeniusRabbit rotator 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package urlgenerator

import (
	"net/url"
	"strings"
	"time"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/client"
	"geniusrabbit.dev/sspserver/internal/eventgenerator"
	"geniusrabbit.dev/sspserver/internal/events"
	"geniusrabbit.dev/sspserver/internal/models"
)

// Generator of URLs
type Generator struct {
	EventGenerator eventgenerator.Generator
	PixelGenerator client.PixelGenerator
	CDNDomain      string
	ClickPattern   string
	DirectPattern  string
}

// CDNURL returns full URL to path
func (g *Generator) CDNURL(path string) string {
	if path == "" {
		return ""
	}
	if path[0] == '/' {
		return "//" + g.CDNDomain + path
	}
	return "//" + g.CDNDomain + "/" + path
}

// PixelURL generator from response of item
func (g *Generator) PixelURL(event events.Type, status uint8, item adsource.ResponserItem, response adsource.Responser, js bool) (string, error) {
	ev, err := g.EventGenerator.Event(event, status, response, item)
	if err != nil {
		return "", err
	}
	return g.PixelGenerator.Event(ev, js)
}

// PixelLead url
func (g *Generator) PixelLead(item adsource.ResponserItem, response adsource.Responser, js bool) (string, error) {
	var sourceID uint64
	if item.Source() != nil {
		sourceID = item.Source().ID()
	}
	return g.PixelGenerator.Lead(&events.LeadCode{
		AuctionID:  response.Request().ID,
		ImpAdID:    item.ID(),
		SourceID:   sourceID,
		ProjectID:  response.Request().ProjectID(),
		CampaignID: item.CampaignID(),
		AdID:       item.AdID(),
		Price:      item.Price(models.ActionLead).Int64(),
		Timestamp:  time.Now().Unix(),
	})
}

// PixelDirectURL generator from response of item
func (g *Generator) PixelDirectURL(event events.Type, status uint8, item adsource.ResponserItem, response adsource.Responser, direct string) (string, error) {
	ev, err := g.EventGenerator.Event(event, status, response, item)
	if err != nil {
		return "", err
	}
	return g.PixelGenerator.EventDirect(ev, direct)
}

// ClickURL generator from respponse of item
func (g *Generator) ClickURL(item adsource.ResponserItem, response adsource.Responser) (string, error) {
	return g.encodeURL(g.ClickPattern, events.Click, events.StatusSuccess, item, response)
}

// MustClickURL generator from respponse of item
func (g *Generator) MustClickURL(item adsource.ResponserItem, response adsource.Responser) string {
	res, _ := g.ClickURL(item, response)
	return res
}

// ClickRouterURL returns router pattern
func (g *Generator) ClickRouterURL() string {
	urls := strings.Split(g.ClickPattern, "?")
	return urls[0]
}

// DirectURL generator from respponse of item
func (g *Generator) DirectURL(event events.Type, item adsource.ResponserItem, response adsource.Responser) (string, error) {
	if event == events.Undefined {
		event = events.Direct
	}
	return g.encodeURL(g.DirectPattern, event, events.StatusSuccess, item, response)
}

// DirectRouterURL returns router pattern
func (g *Generator) DirectRouterURL() string {
	urls := strings.Split(g.DirectPattern, "?")
	return urls[0]
}

// EventCode generator
func (g *Generator) EventCode(event events.Type, status uint8, item adsource.ResponserItem, response adsource.Responser) (string, error) {
	ev, err := g.EventGenerator.Event(event, status, response, item)
	if err != nil {
		return "", err
	}

	code := ev.Pack().Compress().URLEncode()
	return code.String(), code.ErrorObj()
}

func (g *Generator) encodeURL(pattern string, event events.Type, status uint8, item adsource.ResponserItem, response adsource.Responser) (string, error) {
	var (
		ctx       = response.Request().HTTPRequest()
		code, err = g.EventCode(event, status, item, response)
	)

	if err != nil {
		return "", err
	}

	code = url.QueryEscape(code)
	if !strings.Contains(pattern, "{hostname}") {
		if '/' == pattern[0] {
			return "//" + string(ctx.Host()) + strings.Replace(pattern, "{code}", code, -1), nil
		}
		return "//" + string(ctx.Host()) + "/" + strings.Replace(pattern, "{code}", code, -1), nil
	}

	return strings.NewReplacer(
		"{code}", code,
		"{hostname}", string(ctx.Host()),
	).Replace(pattern), nil
}
