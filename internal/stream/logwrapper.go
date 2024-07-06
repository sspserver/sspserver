package stream

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter/v2"
	"go.uber.org/zap"
)

type publisherLogWrapper struct {
	nc.Publisher
	name string
}

// WrapPublisherWithLog of the requests
func WrapPublisherWithLog(name string, pub nc.Publisher) nc.Publisher {
	return &publisherLogWrapper{name: name, Publisher: pub}
}

func (p *publisherLogWrapper) Publish(ctx context.Context, messages ...any) error {
	zap.L().Debug("push message", zap.String("stream", p.name), zap.Any("messages", messages))
	return p.Publisher.Publish(ctx, messages...)
}
