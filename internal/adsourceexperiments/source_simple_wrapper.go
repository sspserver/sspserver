package adsourceexperiments

import (
	"time"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

type sourceSimpleWrapper struct {
	source adsource.Source
}

// NewSimpleWrapper object
func NewSimpleWrapper(source adsource.Source) SourceWrapper {
	return &sourceSimpleWrapper{source: source}
}

// Next returns source interface according to strategy
func (w *sourceSimpleWrapper) Next() adsource.Source {
	return w.source
}

// SetTimeout for sourcer
func (w *sourceSimpleWrapper) SetTimeout(timeout time.Duration) {
	if src, _ := w.source.(adsource.TimeoutSetter); src != nil {
		src.SetTimeout(timeout)
	}
}

var _ SourceWrapper = (*sourceSimpleWrapper)(nil)
