//go:build kafka || allstreams
// +build kafka allstreams

package stream

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/kafka"
)

func init() {
	subscriberConnectors["kafka"] = func(ctx context.Context, url string) (nc.Subscriber, error) {
		return kafka.NewSubscriber(kafka.WithKafkaURL(url))
	}
	publisherConnectors["kafka"] = func(ctx context.Context, url string) (nc.Publisher, error) {
		return kafka.NewPublisher(ctx, kafka.WithKafkaURL(url))
	}
}
