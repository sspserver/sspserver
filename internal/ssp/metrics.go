package ssp

import (
	"time"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

// Metrics wrapper
type Metrics struct{}

// IncrementBidRequestCount metric
func (m Metrics) IncrementBidRequestCount(source adsource.Source, duration time.Duration) {

}

// IncrementBidErrorCount metric
func (m Metrics) IncrementBidErrorCount(source adsource.Source, err error) {

}
