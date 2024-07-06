//go:build redisps || allstreams
// +build redisps allstreams

package stream

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/redis"
)

func init() {
	subscriberConnectors["redis"] = func(ctx context.Context, url string) (nc.Subscriber, error) {
		return redis.NewSubscriber(redis.WithRedisURL(url))
	}
	publisherConnectors["redis"] = func(ctx context.Context, url string) (nc.Publisher, error) {
		return redis.NewPublisher(redis.WithRedisURL(url))
	}
}
