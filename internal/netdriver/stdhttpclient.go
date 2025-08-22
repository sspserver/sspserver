package netdriver

import (
	"context"
	"net/http"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
	"github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

type Driver = *stdhttpclient.Driver

func NewDriver(ctx context.Context, timeout time.Duration) (httpclient.Driver, error) {
	return stdhttpclient.NewDriverWithHTTPClient(
		&http.Client{Timeout: timeout},
	), nil
}
