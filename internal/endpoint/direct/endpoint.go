//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package direct

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/internal/gtracing"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// Error list...
var (
	ErrMultipleDirectNotSupported = errors.New("direct: multiple direct responses not supported")
	ErrInvalidResponseType        = errors.New("direct: invalid response type")
	ErrInvalidFormatsAccessor     = errors.New("direct: FormatsAccessor object is not inited")
)

type debugResponse struct {
	ID                uint64 `json:"id,omitempty"`
	ZoneID            uint64 `json:"zone_id,omitempty"`
	AuctionID         string `json:"auction_id,omitempty"`
	ImpressionID      string `json:"impression_id,omitempty"`
	IsAlternativeLink bool   `json:"is_alternative_link,omitempty"`
	Link              string `json:"link,omitempty"`
	Superfailover     string `json:"superfailover,omitempty"`
	Error             error  `json:"error,omitempty"`
	IsEmpty           bool   `json:"is_empty,omitempty"`
}

type _endpoint struct {
	superFailoverURL string
	source           endpoint.Sourcer
	formats          types.FormatsAccessor
}

func (e _endpoint) Version() string {
	return "v1"
}

func (e _endpoint) Handle(request *adsource.BidRequest) adsource.Responser {
	request.ImpressionUpdate(func(imp *adsource.Impression) bool {
		imp.W, imp.H = -1, -1
		imp.FormatTypes.Reset().Set(types.FormatDirectType)
		return true
	})
	request.Init(e.formats)
	return e.source.Bid(request)
}

func (e _endpoint) Render(ctx *fasthttp.RequestCtx, response adsource.Responser) (err error) {
	var (
		id              uint64
		zoneID          uint64
		impID           string
		link            string
		alternativeLink = false
	)

	if span, _ := gtracing.StartSpanFromContext(response.Request().Context, "render"); span != nil {
		ext.Component.Set(span, "endpoint.direct")
		defer span.Finish()
	}

	if response == nil || response.Count() < 1 {
		if response != nil {
			if imps := response.Request().Imps; len(imps) > 0 {
				impID = imps[0].ID
				if imps[0].Target != nil {
					link = imps[0].Target.AlternativeAdCode("direct")
					zoneID = uint64(imps[0].TargetID())
					alternativeLink = link != ""
				}
			}
		}
	} else if err = response.Validate(); err == nil {
		if response.Count() > 1 {
			err = ErrMultipleDirectNotSupported
		} else {
			adv := response.Ads()[0]
			impID = adv.ImpressionID()
			if adv.Impression() != nil {
				zoneID = uint64(adv.Impression().TargetID())
			}

			switch ad := adv.(type) {
			case adsource.ResponserItem:
				id = ad.AdID()
				if !ad.IsDirect() {
					err = ErrInvalidResponseType
				} else {
					link = ad.ActionURL()
				}
			case adsource.ResponserMultipleItem:
				err = ErrMultipleDirectNotSupported
			default:
				// ...
			}
		}
	}

	switch {
	case response != nil && response.Request().Debug && ctx.QueryArgs().Has("noredirect"):
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetContentType("application/json")
		json.NewEncoder(ctx).Encode(debugResponse{
			ID:                id,
			ZoneID:            zoneID,
			ImpressionID:      impID,
			AuctionID:         response.Request().ID,
			IsAlternativeLink: alternativeLink,
			Link:              link,
			Superfailover:     e.superFailoverURL,
			Error:             err,
			IsEmpty:           response.Count() < 1,
		})
	case link != "":
		ctx.Response.Header.Set("X-Status-Alternative", "1")
		ctx.Redirect(link, http.StatusFound)
	case e.superFailoverURL == "":
		ctx.Success("text/plain", []byte("Please add superfailover link"))
	default:
		ctx.Response.Header.Set("X-Status-Failover", "1")
		ctx.Redirect(e.superFailoverURL, http.StatusFound)
	}
	return
}

func (e _endpoint) PrepareRequest(ctx *fasthttp.RequestCtx) (err error) {
	return
}
