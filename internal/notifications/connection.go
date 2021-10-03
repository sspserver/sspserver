package notifications

import (
	"context"
	"errors"
	"strings"

	"github.com/geniusrabbit/notificationcenter"
	"github.com/geniusrabbit/notificationcenter/dummy"
	"github.com/geniusrabbit/notificationcenter/kafka"
	"github.com/geniusrabbit/notificationcenter/nats"
	"github.com/geniusrabbit/notificationcenter/redis"
)

var (
	errInvalidConnection = errors.New("[notifications] invalid connection")
)

// Connection to the stream
func Connection(ctx context.Context, name, connection string) (err error) {
	if connection == "" || connection == "dummy" {
		return notificationcenter.Register(name, dummy.Publisher{})
	}
	idx := strings.Index(connection, ":")
	if idx < 0 {
		return errInvalidConnection
	}
	var (
		conn   notificationcenter.Publisher
		scheme = connection[:idx]
	)
	switch scheme {
	case "nats":
		if conn, err = nats.NewPublisher(nats.WithNatsURL(connection)); err == nil {
			err = notificationcenter.Register(name, conn)
		}
	case "kafka":
		if conn, err = kafka.NewPublisher(ctx, kafka.WithKafkaURL(connection)); err == nil {
			err = notificationcenter.Register(name, conn)
		}
	case "redis":
		if conn, err = redis.NewPublisher(redis.WithRedisURL(connection)); err == nil {
			err = notificationcenter.Register(name, conn)
		}
	default:
		err = errInvalidConnection
	}
	return err
}
