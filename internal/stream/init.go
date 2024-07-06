package stream

import (
	"context"
	"net/url"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/pkg/errors"
)

// ErrUnsupportedScheme in case if scheme is not defined
var ErrUnsupportedScheme = errors.New(`unsupported scheme`)

type (
	subscribeConnector func(ctx context.Context, url string) (nc.Subscriber, error)
	publisherConnector func(ctx context.Context, url string) (nc.Publisher, error)
)

var (
	subscriberConnectors = map[string]subscribeConnector{}
	publisherConnectors  = map[string]publisherConnector{}
)

// ConnectSubscriber from URL
func ConnectSubscriber(ctx context.Context, urlStr string) (nc.Subscriber, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	conn := subscriberConnectors[parsedURL.Scheme]
	if conn == nil {
		return nil, errors.Wrap(ErrUnsupportedScheme, parsedURL.Scheme)
	}
	return conn(ctx, urlStr)
}

// ConnectPublisher from URL
func ConnectPublisher(ctx context.Context, urlStr string) (nc.Publisher, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	conn := publisherConnectors[parsedURL.Scheme]
	if conn == nil {
		return nil, errors.Wrap(ErrUnsupportedScheme, parsedURL.Scheme)
	}
	pub, err := conn(ctx, urlStr)
	return pub, err
}
