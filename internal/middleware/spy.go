//
// @project geniusrabbit::rotator 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package middleware

import (
	"context"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"geniusrabbit.dev/sspserver/internal/gtracing"
	"geniusrabbit.dev/sspserver/internal/personification"
)

type (
	whoisFn func(ctx context.Context, req *fasthttp.RequestCtx) (personification.Person, error)
	signFn  func(resp personification.Person, ctx *fasthttp.RequestCtx)

	// Spy function wrapper
	Spy func(next func(p personification.Person, ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler
)

// NewSpy wrapper looking information about user and pass it in to
func NewSpy(ctx context.Context, whois whoisFn, sign signFn) Spy {
	return func(next func(p personification.Person, ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler {
		newLog := zap.L().With(zap.String("middleware", "spy"))
		return func(req *fasthttp.RequestCtx) {
			if span, _ := gtracing.StartSpanFromFastContext(req, "middleware.spy"); span != nil {
				gtracing.FastContextWithSpan(req, span)
				defer span.Finish()
			}
			person, err := whois(ctx, req)
			doLogError(newLog, err)
			next(person, req)
			sign(person, req)
		}
	}
}

func doLogError(log *zap.Logger, err error) {
	if err != nil {
		log.Error("", zap.Error(err))
	}
}
