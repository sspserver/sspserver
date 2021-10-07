package stream

import (
	"context"
	"net/url"

	"github.com/geniusrabbit/notificationcenter"
	"github.com/pkg/errors"
)

// ErrUnsupportedScheme in case if scheme is not defined
var ErrUnsupportedScheme = errors.New(`unsupported scheme`)

type (
	subscribeConnector func(ctx context.Context, url string) (notificationcenter.Subscriber, error)
	publisherConnector func(ctx context.Context, url string) (notificationcenter.Publisher, error)
)

var (
	subscriberConnectors = map[string]subscribeConnector{}
	publisherConnectors  = map[string]publisherConnector{}
)

// ConnectSubscriber from URL
func ConnectSubscriber(ctx context.Context, urlStr string) (notificationcenter.Subscriber, error) {
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
func ConnectPublisher(ctx context.Context, urlStr string) (notificationcenter.Publisher, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	conn := publisherConnectors[parsedURL.Scheme]
	if conn == nil {
		return nil, errors.Wrap(ErrUnsupportedScheme, parsedURL.Scheme)
	}
	return conn(ctx, urlStr)
}
