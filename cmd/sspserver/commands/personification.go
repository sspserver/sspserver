package commands

import (
	"context"

	"github.com/geniusrabbit/adcorelib/net/fasthttp/middleware"
	"github.com/geniusrabbit/adcorelib/personification"
	"github.com/geniusrabbit/udetect"
	"github.com/geniusrabbit/udetect/transport/http"

	"github.com/sspserver/sspserver/cmd/sspserver/appcontext"
)

func initPersonification(ctx context.Context, conf *appcontext.PersonConfig) (personification.Client, middleware.Spy, error) {
	var personDetector personification.Client
	if conf.Connect != "" {
		personDetector = udetect.NewClient(http.NewTransport(
			conf.Connect,
			http.WithTimeout(conf.RequestTimeout, conf.KeepAliveTimeout),
		))
	} else {
		personDetector = &personification.SimpleClient{}
	}

	// Init signature
	signature := personification.Signature{
		UUIDName:       conf.UUIDCookieName,
		SessidName:     conf.SessiCookiedName,
		SessidLifetime: conf.SessionLifetime,
		Detector:       personDetector,
	}

	// Prepare spy middleware
	spyMiddleware := middleware.NewSpy(ctx, signature.Whois, signature.SignCookie)

	return personDetector, spyMiddleware, nil
}
