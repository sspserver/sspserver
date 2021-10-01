//
// @project GeniusRabbit rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package adsource

import (
	"context"
)

// AuctionType data
type AuctionType uint8

// Auction types
const (
	UndefinedAuctionType   AuctionType = 0
	FirstPriceAuctionType  AuctionType = 1
	SecondPriceAuctionType AuctionType = 2
)

// IsFirtsPrice auction type
func (at AuctionType) IsFirtsPrice() bool {
	return at == FirstPriceAuctionType
}

// IsSecondPrice auction type
func (at AuctionType) IsSecondPrice() bool {
	return at == SecondPriceAuctionType
}

// Responser type
type Responser interface {
	// AuctionID response
	AuctionID() string

	// AuctionType of request
	AuctionType() AuctionType

	// Source of response
	Source() Sourcer

	// Request information
	Request() *BidRequest

	// Ads list
	Ads() []ResponserItemCommon

	// Item by impression code
	Item(impid string) ResponserItemCommon

	// Count of response items
	Count() int

	// Validate response
	Validate() error

	// Error of the response
	Error() error

	// Context value
	Context(ctx ...context.Context) context.Context

	// Get context item by key
	Get(key string) interface{}
}
