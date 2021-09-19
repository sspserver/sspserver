//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package dynamic

import (
	"encoding/json"

	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/internal/events"
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
	resp := Response{Version: "1"}

	for _, ad := range response.Ads() {
		var (
			assets       []asset
			aditm        = ad.(adsource.ResponserItem)
			url          string
			impPixel, _  = e.urlGen.PixelURL(events.Impression, events.StatusSuccess, aditm, response, false)
			viewPixel, _ = e.urlGen.PixelURL(events.View, events.StatusSuccess, aditm, response, false)
			trackerBlock tracker
		)

		if !aditm.Format().IsProxy() {
			url, _ = e.urlGen.ClickURL(aditm, response)
		}

		trackerBlock = tracker{
			Impressions: []string{impPixel},
			Views:       []string{viewPixel},
		}

		if item, _ := ad.(adsource.ResponserItem); item != nil {
			trackerBlock.Clicks = item.ClickTrackerLinks()
			if links := item.ViewTrackerLinks(); len(links) > 0 {
				trackerBlock.Views = append(trackerBlock.Views, links...)
			}
		}

		for _, as := range aditm.Assets() {
			assets = append(assets, asset{
				Name: as.Name,
				Path: as.Path,
				Type: as.Type.String(),
			})
		}

		group := resp.getGroupOrCreate(ad.ImpressionID())
		item := &item{
			ID:         ad.ID(),
			Type:       ad.PriorityFormatType().Name(),
			URL:        url,
			Content:    aditm.ContentItemString(adsource.ContentItemContent),
			ContentURL: aditm.ContentItemString(adsource.ContentItemIFrameURL),
			Fields:     aditm.ContentFields(),
			Assets:     assets,
			Tracker:    trackerBlock,
		}
		group.Items = append(group.Items, item)

		if response.Request().Debug {
			item.Debug = ad
		}
	}

	format := string(ctx.QueryArgs().Peek("format"))
	if format == "jsonp" {
		callback := string(ctx.QueryArgs().Peek("callback"))
		if callback == "" {
			callback = "callback"
		}
		ctx.Write([]byte(callback + "("))
		json.NewEncoder(ctx).Encode(resp)
		ctx.Write([]byte(")"))
		return nil
	}

	return json.NewEncoder(ctx).Encode(resp)
}

func (e _endpoint) PrepareRequest(ctx *fasthttp.RequestCtx) (err error) {
	return
}

func defString(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
