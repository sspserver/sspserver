//
// @project GeniusRabbit sspserver 2018 - 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2024
//

package dynamic

import (
	"encoding/json"
	"math/rand"

	"github.com/demdxx/gocast/v2"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/sspserver/sspserver/internal/endpoint"
)

type _endpoint struct {
	urlGen adtype.URLGenerator
}

func New(urlGen adtype.URLGenerator) *_endpoint {
	return &_endpoint{urlGen: urlGen}
}

func (e *_endpoint) Codename() string {
	return "dynamic"
}

func (e _endpoint) Handle(source endpoint.Source, request *adtype.BidRequest) adtype.Responser {
	response := source.Bid(request)
	if err := e.render(request.RequestCtx, response); err != nil {
		response = adtype.NewErrorResponse(request, err)
	}
	return response
}

func (e _endpoint) render(ctx *fasthttp.RequestCtx, response adtype.Responser) error {
	resp := Response{Version: "1"}

	// Process response ad items
	for _, ad := range response.Ads() {
		var (
			assets       []asset
			aditm        = ad.(adtype.ResponserItem)
			url          string
			impPixel, _  = e.urlGen.PixelURL(events.Impression, events.StatusSuccess, aditm, response, false)
			viewPixel, _ = e.urlGen.PixelURL(events.View, events.StatusSuccess, aditm, response, false)
			trackerBlock tracker
		)

		// Generate click URL
		if !aditm.Format().IsProxy() {
			url, _ = e.urlGen.ClickURL(aditm, response)
		}

		trackerBlock = tracker{
			Impressions: []string{impPixel},
			Views:       []string{viewPixel},
		}

		// Third-party trackers pixels
		if item, _ := ad.(adtype.ResponserItem); item != nil {
			trackerBlock.Clicks = item.ClickTrackerLinks()
			if links := item.ViewTrackerLinks(); len(links) > 0 {
				trackerBlock.Views = append(trackerBlock.Views, links...)
			}
		}

		// Process assets if provided
		if baseAssets := aditm.Assets(); len(baseAssets) > 0 {
			assets = make([]asset, 0, len(baseAssets))
			processed := map[string]int{}
			for _, as := range baseAssets {
				if idx, ok := processed[as.Name]; !ok || rand.Float64() > 0.5 {
					nas := asset{
						Name:   as.Name,
						Path:   e.urlGen.CDNURL(as.Path),
						Type:   as.Type.Code(),
						Width:  as.Width,
						Height: as.Height,
						Thumbs: e.thumbsPrepare(as.Thumbs),
					}
					if !ok {
						processed[as.Name] = len(assets)
						assets = append(assets, nas)
					} else {
						assets[idx] = nas
					}
				}
			}
		}

		// Add item to response group by impression ID
		resp.getGroupOrCreate(ad.ImpressionID()).addItem(&item{
			ID:         ad.ID(),
			Type:       ad.PriorityFormatType().Name(),
			URL:        url,
			Content:    aditm.ContentItemString(adtype.ContentItemContent),
			ContentURL: aditm.ContentItemString(adtype.ContentItemIFrameURL),
			Fields:     aditm.ContentFields(),
			Assets:     assets,
			Tracker:    trackerBlock,
			Debug:      gocast.IfThen(response.Request().Debug, ad, nil),
		})
	}

	// Render response to the client as JSONP
	format := string(ctx.QueryArgs().Peek("format"))
	if format == "jsonp" {
		callback := string(ctx.QueryArgs().Peek("callback"))
		if callback == "" {
			callback = "callback"
		}
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("application/javascript")
		_, _ = ctx.Write([]byte(callback + "("))
		_ = json.NewEncoder(ctx).Encode(resp)
		_, _ = ctx.Write([]byte(")"))
		return nil
	}

	// Default JSON response
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	return json.NewEncoder(ctx).Encode(resp)
}

func (e _endpoint) thumbsPrepare(thumbs []admodels.AdAssetThumb) []assetThumb {
	nthumbs := make([]assetThumb, 0, len(thumbs))
	for _, th := range thumbs {
		nthumbs = append(nthumbs, assetThumb{
			Path:   e.urlGen.CDNURL(th.Path),
			Type:   th.Type.Code(),
			Width:  th.Width,
			Height: th.Height,
		})
	}
	return nthumbs
}
