package personification

import (
	"time"

	"github.com/sspserver/udetect"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/gtracing"
	fasthttpext "geniusrabbit.dev/sspserver/internal/net/fasthttp"
)

// Signeture provides the builder of cookie assigned to the user by HTTP
type Signeture struct {
	uuidName       string
	sessidName     string
	sessidLifetime time.Duration
	detector       Client
}

// Whois user information
func (sign *Signeture) Whois(ctx *fasthttp.RequestCtx) (Person, error) {
	var (
		uuidCookie   fasthttp.Cookie
		sessidCookie fasthttp.Cookie
	)

	if span, _ := gtracing.StartSpanFromFastContext(ctx, "personification.whois"); span != nil {
		defer span.Finish()
	}

	uuidCookie.ParseBytes(
		ctx.Request.Header.Cookie(sign.uuidName),
	)

	sessidCookie.ParseBytes(
		ctx.Request.Header.Cookie(sign.sessidName),
	)

	uuidObj, _ := udetect.UUIDFromString(string(uuidCookie.Value()))
	sessidObj, _ := udetect.UUIDFromString(string(sessidCookie.Value()))
	request := &udetect.Request{
		Uid:    uuidObj,
		Sessid: sessidObj,
		Ip:     fasthttpext.IPAdressByRequestCF(ctx),
		Ua:     string(ctx.UserAgent()),
		Url:    string(ctx.Referer()),
	}

	_, err := sign.detector.Detect(request)
	return &person{request: request}, err
}

// SignCookie do sign request by traking response
func (sign *Signeture) SignCookie(resp Person, ctx *fasthttp.RequestCtx) {
	if span, _ := gtracing.StartSpanFromFastContext(ctx, "personification.sign"); span != nil {
		defer span.Finish()
	}

	if resp == nil {
		return
	}

	// if len(resp.UserInfo().UUID()) > 0 {
	// 	c := &fasthttp.Cookie{}
	// 	c.SetKey(sign.uuidName)
	// 	c.SetValue(resp.UserInfo().UUID())
	// 	c.SetHTTPOnly(true)
	// 	c.SetExpire(time.Now().Add(365 * 24 * time.Hour))
	// 	ctx.Response.Header.SetCookie(c)
	// }

	// if len(resp.UserInfo().SessionID()) > 0 {
	// 	c := &fasthttp.Cookie{}
	// 	c.SetKey(sign.sessidName)
	// 	c.SetValue(resp.UserInfo().SessionID())
	// 	c.SetHTTPOnly(true)
	// 	c.SetExpire(time.Now().Add(sign.sessidLifetime))
	// 	ctx.Response.Header.SetCookie(c)
	// }
}
