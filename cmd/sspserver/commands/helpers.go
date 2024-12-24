package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/httpserver/extensions/version"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/sspserver/sspserver/cmd/sspserver/appcontext"
	"github.com/sspserver/sspserver/internal/stream"
)

const (
	eventsStreamName = "events"
	winStreamName    = "wins"
)

func configureEventPipeline(ctx context.Context, adServerConf *appcontext.AdServerConfig, debug bool) {
	// Register events data stream
	fatalError(nc.Register(
		eventsStreamName,
		connectPublisherOrLog(ctx,
			eventsStreamName,
			adServerConf.EventPipeline.EventQueue.Connection,
			debug,
		),
	), "register events stream")

	// Register wins info data stream
	fatalError(nc.Register(
		winStreamName,
		connectPublisherOrLog(ctx,
			winStreamName,
			adServerConf.EventPipeline.WinQueue.Connection,
			debug,
		),
	), "register win stream")
}

func notificationMessageLog(streamName string) nc.FuncPublisher {
	return func(ctx context.Context, msgs ...any) error {
		for _, msg := range msgs {
			zap.L().Debug("notification message",
				zap.String("stream", streamName),
				zap.Any("msg", msg),
			)
		}
		return nil
	}
}

func connectPublisherOrLog(ctx context.Context, name, connection string, debug bool) nc.Publisher {
	if connection != "" {
		pub, err := stream.ConnectPublisher(ctx, connection)
		fatalError(err, "connect to '"+connection+"' topics")
		if debug {
			pub = stream.WrapPublisherWithLog(name, pub)
		}
		return pub
	}
	zap.L().Info("register new dummy publisher", zap.String("name", name))
	return notificationMessageLog(name)
}

func infoHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	headers := make(map[string]any)
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	ver := version.Get(ctx)

	_ = json.NewEncoder(ctx).Encode(map[string]any{
		"commit":  ver.Commit,
		"version": ver.Version,
		"date":    ver.Date,
		"request": map[string]any{
			"uri":     string(ctx.Request.URI().FullURI()),
			"method":  string(ctx.Method()),
			"headers": headers,
		},
	})
}

func fatalError(err error, message ...any) {
	if err != nil {
		zap.L().Fatal(fmt.Sprint(message...), zap.Error(err))
	}
}
