package targetaccessor

import "geniusrabbit.dev/sspserver/internal/models"

// Accessor interface
type Accessor interface {
	TargetByID(id uint64) models.Target
}
