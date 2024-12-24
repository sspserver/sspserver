package jobs

import (
	"context"
	"time"

	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"go.uber.org/zap"
)

func RunIntervalJob(ctx context.Context, name string, interval time.Duration, job func(context.Context) error) {
	if err := job(ctx); err != nil {
		ctxlogger.Get(ctx).Error("first run job failed",
			zap.String("job", name),
			zap.Error(err),
		)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := job(ctx); err != nil {
				ctxlogger.Get(ctx).Error("job failed",
					zap.Duration("interval", interval),
					zap.String("job", name),
					zap.Error(err),
				)
			}
		}
	}
}
