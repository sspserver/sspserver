package adsource

import (
	"time"
)

// RequestStrategy defines politics of request sending
type RequestStrategy int

const (
	// AsynchronousRequestStrategy is default strategy implies
	// requesting all auction participants and choising the most
	// profitable variant of all
	AsynchronousRequestStrategy RequestStrategy = iota

	// SingleRequestStrategy tells that if response was
	// received it should be performed
	SingleRequestStrategy
)

// DefaultSourceID constant
const (
	DefaultSourceID = 0
)

// MinimalSource contains only minimal set of methods
type MinimalSource interface {
	// Bid request for standart system filter
	Bid(request *BidRequest) Responser

	// ProcessResponseItem result or error
	ProcessResponseItem(Responser, ResponserItem)
}

// Source of advertisement and where will be selled the traffic
type Source interface {
	MinimalSource

	// ID of the source driver
	ID() uint64

	// Test request before processing
	Test(request *BidRequest) bool

	// RevenueShareReduceFactor which is a potential
	RevenueShareReduceFactor() float64

	// RequestStrategy description
	RequestStrategy() RequestStrategy
}

// SourceIterator returns next source from the scope
type SourceIterator interface {
	// Next returns source interface according to strategy
	Next() Source
}

// SourceAccessor preoritise the source access
type SourceAccessor interface {
	// Iterator returns the configured queue accessor
	Iterator(request *BidRequest) SourceIterator

	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}

type minimalSourceWrapper struct {
	source MinimalSource
}

// ID of the source driver
func (wp minimalSourceWrapper) ID() uint64 { return 0 }

// Test request before processing
func (wp minimalSourceWrapper) Test(request *BidRequest) bool { return true }

// Bid request for standart system filter
func (wp minimalSourceWrapper) Bid(request *BidRequest) Responser { return wp.source.Bid(request) }

// ProcessResponseItem result or error
func (wp minimalSourceWrapper) ProcessResponseItem(response Responser, item ResponserItem) {
	wp.source.ProcessResponseItem(response, item)
}

// RevenueShareReduceFactor which is a potential
func (wp minimalSourceWrapper) RevenueShareReduceFactor() float64 { return 0 }

// RequestStrategy description
func (wp minimalSourceWrapper) RequestStrategy() RequestStrategy {
	return AsynchronousRequestStrategy
}

// ToSource interface from different types of interfaces with the implementation of unsupported methods
func ToSource(val MinimalSource) Source {
	switch v := val.(type) {
	case Source:
		return v
	case MinimalSource:
		return minimalSourceWrapper{source: v}
	}
	return nil
}
