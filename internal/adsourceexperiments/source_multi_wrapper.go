package adsourceexperiments

import (
	"time"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

// SourceMultiWrapper implements source functionality with source choicer
type sourceMultiWrapper struct {
	sources       []adsource.Source
	rotateStategy IndexStrategy
}

// NewSourceWrapper with roundrobin strategy
func NewSourceWrapper(sources ...adsource.Source) SourceWrapper {
	return &sourceMultiWrapper{
		sources:       sources,
		rotateStategy: NewCounterStrategy(),
	}
}

// NewRandomSourceWrapper with roundrobin strategy
func NewRandomSourceWrapper(sources ...adsource.Source) SourceWrapper {
	return &sourceMultiWrapper{
		sources:       sources,
		rotateStategy: NewRandomStrategy(),
	}
}

// Next returns source interface according to strategy
func (w *sourceMultiWrapper) Next() adsource.Source {
	return w.sources[w.rotateStategy.GetIndex()%uint(len(w.sources))]
}

// SetTimeout for sourcer
func (w *sourceMultiWrapper) SetTimeout(timeout time.Duration) {
	for _, source := range w.sources {
		if src, _ := source.(adsource.TimeoutSetter); src != nil {
			src.SetTimeout(timeout)
		}
	}
}

var _ SourceWrapper = (*sourceMultiWrapper)(nil)
