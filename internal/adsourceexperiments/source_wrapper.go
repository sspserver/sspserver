package adsourceexperiments

import (
	"time"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

// SourceWrapper advertisement accessor interface
type SourceWrapper interface {
	// Next returns source interface according to strategy
	Next() adsource.Source

	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}
