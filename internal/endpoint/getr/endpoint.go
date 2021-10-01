//
// @project GeniusRabbit rotator 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//

// Automatic register auction

package getr

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/internal/eventtraking/events"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

const (
	requestPreparedKey     = "$_req"
	requestHeaderXRealIP   = "X-Real-Ip"
	requestHeaderUserAgent = "User-Agent"
)

// Error list...
var (
	ErrInvalidRequest              = errors.New("[getr]: invalid request")
	ErrInvalidResponseType         = errors.New("[getr]: invalid response type")
	ErrInvalidFormatsAccessor      = errors.New("[getr]: FormatsAccessor object is not inited")
	ErrInvalidInvalidRequestMethod = errors.New("[getr]: invald request method")
)

type _endpoint struct {
	// Superfailover link URL
	superFailoverURL string

	// Base advertisement source
	source endpoint.Sourcer

	// Formats accessor
	formats types.FormatsAccessor

	// URL generator object
	urlGen adsource.URLGenerator
}

func (e _endpoint) Version() string {
	return "v1"
}

func (e _endpoint) PrepareRequest(ctx *fasthttp.RequestCtx) (err error) {
	var (
		inReq = &registerRequest{}
		data  = map[string]interface{}{}
	)

	// Collect and bind information into the request object
	// ! if HTTP method is wrong will be returns error
	switch {
	case ctx.IsGet():
		ctx.QueryArgs().VisitAll(func(key, value []byte) { data[string(key)] = string(value) })
		err = inReq.QueryMapDecode(data)
	case ctx.IsPost() || ctx.IsPut():
		ctx.PostArgs().VisitAll(func(key, value []byte) { data[string(key)] = string(value) })
		err = inReq.FormMapDecode(data)
	default:
		err = ErrInvalidInvalidRequestMethod
	}
	if inReq.IP != "" {
		ctx.Request.Header.Set(requestHeaderXRealIP, inReq.IP)
	}
	if inReq.UserAgent != "" {
		ctx.Request.Header.Set(requestHeaderUserAgent, inReq.UserAgent)
	}
	ctx.SetUserValue(requestPreparedKey, inReq)
	return
}

func (e _endpoint) Handle(request *adsource.BidRequest) adsource.Responser {
	var (
		inReq *registerRequest
		err   error
	)

	// Collect and bind information into the request object
	// ! if HTTP method is wrong will be returns error
	{
		httpReq := request.HTTPRequest()
		inReq, _ = httpReq.UserValue(requestPreparedKey).(*registerRequest)
		if inReq == nil {
			err = ErrInvalidRequest
		}

		if err == nil {
			var (
				id       uint64
				driverID = httpReq.UserValue("driver")
			)
			if driverID == nil || driverID.(string) == "main" {
				request.SetSourceFilter(adsource.DefaultSourceID)
			} else if id, err = strconv.ParseUint(driverID.(string), 10, 64); err == nil {
				request.SetSourceFilter(id)
			}
		} // end if
	}

	if err == nil {
		err = inReq.Validate()
	}

	if err != nil {
		return adsource.NewErrorResponse(request, err)
	}

	// Format reset of impression
	request.ImpressionUpdate(func(imp *adsource.Impression) bool {
		imp.FormatTypes.Reset().Set(types.FormatAutoregisterType)
		imp.Count = 1
		return true
	})

	// Update request with data of request
	inReq.FillBidRequest(request)

	// Reinit the impression information
	request.Init(e.formats)

	// * Send request to the source object
	return e.source.Bid(request)
}

func (e _endpoint) Render(ctx *fasthttp.RequestCtx, response adsource.Responser) (err error) {
	var (
		adID uint64
		adIt adsource.ResponserItem
	)

	if err = response.Validate(); err == nil {
		for _, adc := range response.Ads() {
			switch ad := adc.(type) {
			case adsource.ResponserItem:
				if ad.ActionURL() == "" {
					// TODO: send fail metric
					continue
				}
				adID = ad.AdID()
				adIt = ad
			default:
				// TODO: send fail metric
				continue
			}
		}
	}

	if response.Request().Debug {
		var (
			targetLink   string
			originalLink string
		)
		if adIt != nil {
			targetLink, _ = e.urlGen.DirectURL(events.View, adIt, response)
			originalLink = adIt.ActionURL()
		}
		err = json.NewEncoder(ctx).Encode(struct {
			AdID      uint64 `json:"aid,omitempty"`
			Link      string `json:"link,omitempty"`
			PixelLink string `json:"pixel_link,omitempty"`
			Error     string `json:"error,omitempty"`
		}{
			AdID:      adID,
			Link:      originalLink,
			PixelLink: targetLink,
			Error:     errorMessage(err),
		})
	} else if adIt != nil && adIt.ActionURL() != "" && err == nil {
		var targetLink string
		targetLink, err = e.urlGen.DirectURL(events.View, adIt, response)

		if err == nil {
			ctx.SetStatusCode(http.StatusOK)
			err = json.NewEncoder(ctx).Encode(struct {
				Link string `json:"link"`
			}{
				Link: targetLink,
			})
		} else {
			e.noAdResponse(ctx, err)
		}
	} else {
		// TODO: error log and metrics
		e.noAdResponse(ctx, err)
	}
	return err
}

func (e _endpoint) Routing() map[string][]string {
	return map[string][]string{
		":zone":         {http.MethodGet, http.MethodPost, http.MethodPut},
		":zone/:driver": {http.MethodGet, http.MethodPost, http.MethodPut},
	}
}

func (e _endpoint) noAdResponse(ctx *fasthttp.RequestCtx, err error) {
	if err != nil {
		ctx.Request.Header.Set("X-Response-Error", err.Error())
	}
	ctx.SetStatusCode(http.StatusNoContent)
}

func errorMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
