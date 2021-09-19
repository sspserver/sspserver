package adsource

import (
	"context"

	"github.com/bsm/openrtb"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// ResponseEmpty object
type ResponseEmpty struct {
	ItemID  string
	Req     *BidRequest
	Src     Sourcer
	Imp     *Impression
	Bid     *openrtb.Bid
	Err     error
	context context.Context
}

// NewEmptyResponse by bid request
func NewEmptyResponse(request *BidRequest, err error) *ResponseEmpty {
	return &ResponseEmpty{Req: request, Err: err}
}

// ID of current response item (unique code of current response)
func (r ResponseEmpty) ID() string {
	return r.ItemID
}

// AuctionID response
func (r ResponseEmpty) AuctionID() string {
	if r.Req == nil {
		return ""
	}
	return r.Req.ID
}

// AuctionType of request
func (r ResponseEmpty) AuctionType() AuctionType {
	if r.Req == nil {
		return UndefinedAuctionType
	}
	return r.Req.AuctionType
}

// Source of response
func (r ResponseEmpty) Source() Sourcer {
	return r.Src
}

// AccessPoint of response
func (r ResponseEmpty) AccessPoint() AccessPoint {
	if r.Req == nil || r.Req.AccessPoint == nil {
		return nil
	}
	return r.Req.AccessPoint
}

// PriorityFormatType from current Ad
func (r ResponseEmpty) PriorityFormatType() types.FormatType {
	return types.FormatUndefinedType
}

// Request information
func (r ResponseEmpty) Request() *BidRequest {
	return r.Req
}

// Impression place object
func (r ResponseEmpty) Impression() *Impression {
	return r.Imp
}

// ImpressionID unique code string
func (r ResponseEmpty) ImpressionID() string {
	if r.Imp == nil {
		return ""
	}
	return r.Imp.ID
}

// ExtImpressionID it's unique code of the auction bid impression
func (r ResponseEmpty) ExtImpressionID() string {
	if r.Bid == nil {
		return ""
	}
	return r.Bid.ImpID
}

// ExtTargetID of the external network
func (r ResponseEmpty) ExtTargetID() string {
	if r.Imp == nil {
		return ""
	}
	return r.Imp.ExtTargetID
}

// Ads list
func (r ResponseEmpty) Ads() []ResponserItemCommon {
	return nil
}

// Item by impression code
func (r ResponseEmpty) Item(impid string) ResponserItemCommon {
	return nil
}

// Price of whole response
func (r ResponseEmpty) Price() billing.Money {
	return 0
}

// LeadPrice summ
func (r ResponseEmpty) LeadPrice() billing.Money {
	return 0
}

// CPMBid value
func (r ResponseEmpty) CPMBid() billing.Money {
	return 0
}

// AuctionCPMBid value price without any comission
func (r ResponseEmpty) AuctionCPMBid() billing.Money {
	return 0
}

// Count of response items
func (r ResponseEmpty) Count() int {
	return 0
}

// Validate response
func (r ResponseEmpty) Validate() error {
	if r.Err != nil {
		return r.Err
	}
	return ErrResponseEmpty
}

// Error of the response
func (r ResponseEmpty) Error() error {
	return r.Err
}

// Context value
func (r *ResponseEmpty) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 && ctx[0] != nil {
		r.context = ctx[0]
	}
	return r.context
}

// Get ext field
func (r *ResponseEmpty) Get(key string) (res interface{}) {
	if r.context != nil {
		res = r.context.Value(key)
	}
	return
}

var (
	_ ResponserItemCommon = (*ResponseEmpty)(nil)
	_ Responser           = (*ResponseEmpty)(nil)
)
