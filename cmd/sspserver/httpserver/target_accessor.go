package httpserver

import "geniusrabbit.dev/sspserver/internal/models"

type targetAccessor interface {
	TargetByID(id uint64) models.Target
}
