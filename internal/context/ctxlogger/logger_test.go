package ctxlogger

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	ctx = WithLogger(ctx, logger)
	assert.NotNil(t, Get(ctx))
}
