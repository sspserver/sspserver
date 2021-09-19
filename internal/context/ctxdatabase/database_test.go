package ctxdatabase

import (
	"context"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	ctx := context.Background()
	ctx = WithDatabase(ctx, &gorm.DB{}, &gorm.DB{})
	assert.NotNil(t, Master(ctx))
	assert.NotNil(t, Readonly(ctx))
}
