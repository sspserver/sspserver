package personification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sspserver/udetect"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/gtracing"
	fasthttpext "geniusrabbit.dev/sspserver/internal/net/fasthttp"
)

// Signeture provides the builder of cookie assigned to the user by HTTP
type Signeture struct {
	UUIDName       string
	SessidName     string
	SessidLifetime time.Duration
	Detector       Client
}

// Whois user information
func (sign *Signeture) Whois(ctx context.Context, req *fasthttp.RequestCtx) (Person, error) {
	var (
		uuidCookie   fasthttp.Cookie
		sessidCookie fasthttp.Cookie
	)

	if span, _ := gtracing.StartSpanFromFastContext(req, "personification.whois"); span != nil {
		defer span.Finish()
	}

	uuidCookie.ParseBytes(
		req.Request.Header.Cookie(sign.UUIDName),
	)

	sessidCookie.ParseBytes(
		req.Request.Header.Cookie(sign.SessidName),
	)

	uuidObj, _ := uuid.Parse(string(uuidCookie.Value()))
	sessidObj, _ := uuid.Parse(string(sessidCookie.Value()))
	request := &udetect.Request{
		UID:             uuidObj,
		SessID:          sessidObj,
		IP:              fasthttpext.IPAdressByRequestCF(req),
		UA:              string(req.UserAgent()),
		URL:             string(req.Referer()),
		Ref:             "", // TODO: add additional information
		DNT:             0,
		LMT:             0,
		Adblock:         0,
		PrivateBrowsing: 0,
		JS:              0,
		Languages:       nil,
		PrimaryLanguage: "",
		FlashVer:        "",
		Width:           0,
		Height:          0,
		Extensions:      nil,
	}

	_, err := sign.Detector.Detect(ctx, request)
	return &person{request: request}, err
}

// SignCookie do sign request by traking response
func (sign *Signeture) SignCookie(resp Person, req *fasthttp.RequestCtx) {
	if span, _ := gtracing.StartSpanFromFastContext(req, "personification.sign"); span != nil {
		defer span.Finish()
	}

	if resp == nil {
		return
	}

	if _uuid := resp.UserInfo().UUID(); len(_uuid) > 0 {
		c := &fasthttp.Cookie{}
		c.SetKey(sign.UUIDName)
		c.SetValue(_uuid)
		c.SetHTTPOnly(true)
		c.SetExpire(time.Now().Add(365 * 24 * time.Hour))
		req.Response.Header.SetCookie(c)
	}

	if sessid := resp.UserInfo().SessionID(); len(sessid) > 0 {
		c := &fasthttp.Cookie{}
		c.SetKey(sign.SessidName)
		c.SetValue(sessid)
		c.SetHTTPOnly(true)
		c.SetExpire(time.Now().Add(sign.SessidLifetime))
		req.Response.Header.SetCookie(c)
	}
}
