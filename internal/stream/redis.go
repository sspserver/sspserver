//go:build redisps || allstreams
// +build redisps allstreams

package stream

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/redis"
)

func init() {
	subscriberConnectors["redis"] = func(ctx context.Context, url string) (nc.Subscriber, error) {
		return redis.NewSubscriber(redis.WithRedisURL(url))
	}
	publisherConnectors["redis"] = func(ctx context.Context, url string) (nc.Publisher, error) {
		return redis.NewPublisher(redis.WithRedisURL(url))
	}
}
