package platform

import (
	"io"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/ssp/openlatency"
)

// Platformer accessor
type Platformer interface {
	adsource.Source
}

// State of the platform
type State interface {
	// GetState of the driver object
	GetState() interface{}

	// SetState from driver object
	SetState(state interface{})

	// GetStateData of the platform
	GetStateData() (io.Reader, error)

	// LoadStateData of the platform
	LoadStateData(data io.Reader) error
}

// Metrics of the platform
type Metrics interface {
	Metrics() *openlatency.MetricsInfo
}
