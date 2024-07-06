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
	"go.uber.org/zap"

	"geniusrabbit.dev/adcorelib/admodels/types"
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/context/ctxlogger"
	"geniusrabbit.dev/adcorelib/eventtraking/events"
	"geniusrabbit.dev/adcorelib/eventtraking/eventstream"
	"geniusrabbit.dev/adcorelib/gtracing"
	"geniusrabbit.dev/sspserver/internal/endpoint"
)

// Error list...
var (
	ErrMultipleDirectNotSupported = errors.New("direct: multiple direct responses not supported")
	ErrInvalidResponseType        = errors.New("direct: invalid response type")
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
	formats          types.FormatsAccessor
	superFailoverURL string
}

func New(formats types.FormatsAccessor, superFailoverURL string) *_endpoint {
	return &_endpoint{
		formats:          formats,
		superFailoverURL: superFailoverURL,
	}
}

func (e *_endpoint) Codename() string {
	return "direct"
}

func (e *_endpoint) Handle(source endpoint.Source, request *adtype.BidRequest) adtype.Responser {
	request.ImpressionUpdate(func(imp *adtype.Impression) bool {
		imp.W, imp.H = -1, -1
		imp.FormatTypes.Reset().Set(types.FormatDirectType)
		return true
	})
	request.Init(e.formats)
	response := source.Bid(request)
	if err := e.execDirect(request.RequestCtx, response); err != nil {
		ctxlogger.Get(request.Ctx).Error("exec direct", zap.Error(err))
	} else {
		e.sendViewEvent(response)
	}
	return response
}

func (e *_endpoint) execDirect(req *fasthttp.RequestCtx, response adtype.Responser) (err error) {
	var (
		id              uint64
		zoneID          uint64
		impID           string
		link            string
		alternativeLink = false
	)

	if span, _ := gtracing.StartSpanFromFastContext(req, "render"); span != nil {
		ext.Component.Set(span, "endpoint.direct")
		defer span.Finish()
	}

	if response == nil || response.Count() == 0 {
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
			case adtype.ResponserItem:
				id = ad.AdID()
				if !ad.IsDirect() {
					err = ErrInvalidResponseType
				} else {
					link = adtype.PrepareURL(ad.ActionURL(), response, ad)
				}
			case adtype.ResponserMultipleItem:
				err = ErrMultipleDirectNotSupported
			default:
				// ...
			}
		}
	}

	switch {
	case response != nil && response.Request().Debug && req.QueryArgs().Has("noredirect"):
		req.SetStatusCode(http.StatusOK)
		req.SetContentType("application/json")
		json.NewEncoder(req).Encode(debugResponse{
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
		req.Response.Header.Set("X-Status-Alternative", "1")
		req.Redirect(link, http.StatusFound)
	case e.superFailoverURL == "":
		req.Success("text/plain", []byte("Please add superfailover link"))
	default:
		req.Response.Header.Set("X-Status-Failover", "1")
		req.Redirect(e.superFailoverURL, http.StatusFound)
	}
	return err
}

func (e *_endpoint) sendViewEvent(response adtype.Responser) {
	if response == nil || response.Error() != nil || len(response.Ads()) == 0 {
		return
	}
	if response.Request().Debug && response.Request().RequestCtx.QueryArgs().Has("noredirect") {
		ctxlogger.Get(response.Context()).Info("skip event log", zap.String("request_id", response.Request().ID))
		return
	}
	var (
		err    error
		stream = eventstream.StreamFromContext(response.Context())
	)

	switch ad := response.Ads()[0].(type) {
	case adtype.ResponserItem:
		err = stream.Send(events.Direct, events.StatusSuccess, response, ad)
	case adtype.ResponserMultipleItem:
		if len(ad.Ads()) > 0 {
			err = stream.Send(events.Direct, events.StatusSuccess, response, ad.Ads()[0])
		}
	default:
		// Invalid ad type
	}
	if err != nil {
		ctxlogger.Get(response.Context()).Error("send direct event", zap.Error(err))
	}
}
