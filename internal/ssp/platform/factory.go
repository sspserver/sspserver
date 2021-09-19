package platform

import (
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/ssp/platform/info"
)

// Factory of platform
type Factory interface {
	New(source *models.RTBSource, opts ...interface{}) (Platformer, error)
	Info() info.Platform
}
