//
// @project GeniusRabbit rotator 2016 – 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2021
//

package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/demdxx/gocast"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"bitbucket.org/geniusrabbit/bigbrother/client"
	fasthttpext "bitbucket.org/geniusrabbit/corelib/net/fasthttp"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/infostructs"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
	"geniusrabbit.dev/sspserver/internal/rand"
)

// Errors
var (
	ErrInvalidTargetZone = errors.New("Invalid target zone")
)

// RequestOptions prepare
type RequestOptions struct {
	Debug   bool
	Request *fasthttp.RequestCtx
	Count   int
	X, Y    int
	W, WMax int
	H, HMax int
	Page    string
	SubID1  string
	SubID2  string
	SubID3  string
	SubID4  string
	SubID5  string
}

// NewRequestOptions prepare
func NewRequestOptions(ctx *fasthttp.RequestCtx) *RequestOptions {
	var (
		queryArgs        = ctx.QueryArgs()
		w, h, minW, minH = getSizeByCtx(ctx)
		debug, _         = strconv.ParseBool(string(queryArgs.Peek("debug")))
	)

	return &RequestOptions{
		Debug:   debug,
		Request: ctx,
		X:       gocast.ToInt(string(queryArgs.Peek("x"))),
		Y:       gocast.ToInt(string(queryArgs.Peek("y"))),
		W:       minW,
		WMax:    ifPositiveNumber(w, -1),
		H:       minH,
		HMax:    ifPositiveNumber(h, -1),
		SubID1:  queryParam(queryArgs, "subid1", "subid", "s1"),
		SubID2:  queryParam(queryArgs, "subid2", "s2"),
		SubID3:  queryParam(queryArgs, "subid3", "s3"),
		SubID4:  queryParam(queryArgs, "subid4", "s4"),
		SubID5:  queryParam(queryArgs, "subid5", "s5"),
	}
}

// NewDirectRequestOptions prepare
func NewDirectRequestOptions(ctx *fasthttp.RequestCtx) *RequestOptions {
	var (
		queryArgs = ctx.QueryArgs()
		debug, _  = strconv.ParseBool(string(queryArgs.Peek("debug")))
	)

	return &RequestOptions{
		Debug:   debug,
		Request: ctx,
		X:       gocast.ToInt(queryArgs.Peek("x")),
		Y:       gocast.ToInt(queryArgs.Peek("y")),
		W:       -1,
		H:       -1,
		SubID1:  queryParam(queryArgs, "subid1", "subid", "s1"),
		SubID2:  queryParam(queryArgs, "subid2", "s2"),
		SubID3:  queryParam(queryArgs, "subid3", "s3"),
		SubID4:  queryParam(queryArgs, "subid4", "s4"),
		SubID5:  queryParam(queryArgs, "subid5", "s5"),
	}
}

// NewRequestFor person
func NewRequestFor(ctx context.Context, requestID string, target models.Target, person client.Person, opt *RequestOptions, formatAccessor types.FormatsAccessor) (req *adsource.BidRequest) {
	var (
		userInfo         = person.UserInfo()
		ageStart, ageEnd = userInfo.Ages()
		referer          = string(opt.Request.Referer())
	)
	if requestID == "" {
		requestID = rand.UUID()
	}

	req = &adsource.BidRequest{
		ID:         requestID,
		Debug:      opt.Debug,
		RequestCtx: opt.Request,
		Secure:     b2i(fasthttpext.IsSecureCF(opt.Request)),
		Device:     userInfo.DeviceInfo(),
		Imps: []adsource.Impression{
			{
				ID:          rand.UUID(), // Impression ID
				ExtTargetID: "",
				Target:      target,
				FormatTypes: directTypeMask(opt.W == -1 && opt.H == -1),
				Count:       minInt(opt.Count, 1),
				X:           opt.X,
				Y:           opt.Y,
				W:           opt.W,
				H:           opt.H,
				WMax:        opt.WMax,
				HMax:        opt.HMax,
				SubID1:      opt.SubID1,
				SubID2:      opt.SubID2,
				SubID3:      opt.SubID3,
				SubID4:      opt.SubID4,
				SubID5:      opt.SubID5,
			},
		},
		User: &adsource.User{
			ID:            userInfo.UUID(),                     // Unique User ID
			SessionID:     userInfo.SessionID(),                // Unique session ID
			FingerPrintID: userInfo.Fingerprint(),              //
			ETag:          userInfo.ETag(),                     //
			AgeStart:      ageStart,                            // Year of birth from
			AgeEnd:        ageEnd,                              // Year of birth from
			Gender:        sexFrom(userInfo.MostPossibleSex()), // Gender ("M": male, "F" female, "O" Other)
			Keywords:      userInfo.Keywords(),                 // Comma separated list of keywords, interests, or intent
			Geo:           userInfo.GeoInfo(),
		},
		Site: &infostructs.Site{
			ExtID:         "",              // External ID
			Domain:        domain(referer), //
			Cat:           nil,             // Array of categories
			PrivacyPolicy: 1,               // Default: 1 ("1": has a privacy policy)
			Keywords:      "",              // Comma separated list of keywords about the site.
			Page:          referer,         // URL of the page
			Ref:           referer,         // Referrer URL
			Search:        "",              // Search string that caused naviation
			Mobile:        0,               // Mobile ("1": site is mobile optimised)
		},
		Person:   person,
		Context:  ctx,
		Timemark: time.Now(),
	}
	req.Init(formatAccessor)
	return
}

// NewRequestByContext from request
func NewRequestByContext(ctx context.Context, req *fasthttp.RequestCtx) (*adsource.BidRequest, error) {
	request := &adsource.BidRequest{RequestCtx: req, Timemark: time.Now(), Context: ctx}
	if err := json.NewDecoder(bytes.NewBuffer(req.Request.Body())).Decode(request); err != nil {
		return nil, err
	}
	return request, nil
}

///////////////////////////////////////////////////////////////////////////////
/// Helper methods
///////////////////////////////////////////////////////////////////////////////

func directTypeMask(is bool) types.FormatTypeBitset {
	if is {
		return *types.NewFormatTypeBitset(types.FormatDirectType)
	}
	return types.FormatTypeBitsetEmpty
}

func domain(surl string) (name string) {
	if len(surl) < 1 {
		name = ""
	} else if len(surl) < 7 {
		name = strings.Split(surl, ",")[0]
	} else {
		switch strings.ToLower(surl[:7]) {
		case "http://", "https:/":
			if u, err := url.Parse(surl); nil == err {
				name = u.Host
			}
		}
	}
	return name
}

func queryParam(args *fasthttp.Args, names ...string) (v string) {
	for _, n := range names {
		if qv := args.Peek(n); len(qv) > 0 {
			v = string(qv)
		}
	}
	return
}

func sexFrom(v int) string {
	switch v {
	case 1:
		return "M"
	case 2:
		return "F"
	}
	return "?"
}

func getSizeByCtx(ctx *fasthttp.RequestCtx) (sw, sh, minSW, minSH int) {
	var (
		queryArgs = ctx.QueryArgs()
		w         = string(queryArgs.Peek("w"))
		h         = string(queryArgs.Peek("h"))
		minW      = string(queryArgs.Peek("mw"))
		minH      = string(queryArgs.Peek("mh"))
	)

	if isEmptyNumString(w) && isEmptyNumString(h) {
		if s := strings.Split(string(queryArgs.Peek("fmt")), "x"); len(s) > 0 {
			if len(s) == 2 {
				w, h = s[0], s[1]
			} else {
				w = s[0]
			}
		}
	}

	sw, sh, minSW, minSH = gocast.ToInt(w), gocast.ToInt(h),
		gocast.ToInt(minW), gocast.ToInt(minH)

	if sw < minSW {
		sw, minSW = minSW, sw
	}
	if sh < minSH {
		sh, minSH = minSH, sh
	}
	return sw, sh, minSW, minSH
}

func logger(id string, ctx *fasthttp.RequestCtx) (log *zap.Logger) {
	if lg := ctx.UserValue("_log"); lg != nil {
		if log, _ = lg.(*zap.Logger); log != nil {
			log = log.With(zap.String("bidrequest_id", id))
		}
	}
	return log
}

func minInt(v1, v2 int) int {
	if v1 > 0 {
		return v1
	}
	return v2
}

func ifPositiveNumber(v1, v2 int) int {
	if v1 > 0 {
		return v1
	}
	return v2
}

func isEmptyNumString(s string) bool {
	return s == "" || s == "0"
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
