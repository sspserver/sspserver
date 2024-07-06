//go:build nats || allstreams
// +build nats allstreams

package stream

import (
	"context"
	"time"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/nats"
	natsio "github.com/nats-io/nats.go"
)

func init() {
	subscriberConnectors["nats"] = func(ctx context.Context, url string) (nc.Subscriber, error) {
		return nats.NewSubscriber(nats.WithNatsURL(url), nats.WithNatsOptions(natsio.ReconnectWait(time.Second*5)))
	}
	publisherConnectors["nats"] = func(ctx context.Context, url string) (nc.Publisher, error) {
		return nats.NewPublisher(nats.WithNatsURL(url))
	}
}
