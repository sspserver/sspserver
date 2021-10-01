//
// @project geniusrabbit::rotator 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package middleware

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/gtracing"
	"geniusrabbit.dev/sspserver/internal/personification"
)

type whoisFn func(ctx *fasthttp.RequestCtx) (personification.Person, error)
type signFn func(resp personification.Person, ctx *fasthttp.RequestCtx)

// Spy function wrapper
type Spy func(next func(p personification.Person, ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler

// NewSpy wrapper looking information about user and pass it in to
func NewSpy(whois whoisFn, sign signFn) Spy {
	return func(next func(p personification.Person, ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler {
		var newLog = logrus.WithField("middleware", "spy")

		return func(ctx *fasthttp.RequestCtx) {
			if span, _ := gtracing.StartSpanFromFastContext(ctx, "middleware.spy"); span != nil {
				gtracing.FastContextWithSpan(ctx, span)
				defer span.Finish()
			}

			person, err := whois(ctx)
			doLogError(newLog, err)
			next(person, ctx)
			sign(person, ctx)
		}
	}
}

func doLogError(log *logrus.Entry, err error) {
	if err != nil {
		log.Error(err)
	}
}
