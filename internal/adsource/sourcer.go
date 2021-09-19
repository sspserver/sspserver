//
// @project GeniusRabbit rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package adsource

import (
	"time"
)

// Sourcer – source for DSP platform
type Sourcer interface {
	// ID of source
	ID() uint64

	// Bid request for standart system filter
	Bid(request *BidRequest) Responser

	// ProcessResponseItem result
	ProcessResponseItem(Responser, ResponserItem)

	// RevenueShareReduceFactor which is a potential
	RevenueShareReduceFactor() float64
}

// Tester checker
type Tester interface {
	// Test current request for compatibility
	Test(request *BidRequest) bool
}

// TimeoutSetter interface
type TimeoutSetter interface {
	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}

// SourceTester interface
type SourceTester interface {
	Sourcer
	Tester
}
