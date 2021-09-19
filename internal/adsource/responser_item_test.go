package adsource

import (
	"context"
	"testing"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

func TestPrices(t *testing.T) {

}

type priceItem struct {
	pricingModel types.PricingModel
	maxBid       billing.Money
	price        billing.Money
	leadPrice    billing.Money
	ecpm         billing.Money
}

func (it priceItem) ID() string                                     { return "" }
func (it priceItem) Impression() *Impression                        { return nil }
func (it priceItem) ImpressionID() string                           { return "" }
func (it priceItem) ExtImpressionID() string                        { return "" }
func (it priceItem) ExtTargetID() string                            { return "" }
func (it priceItem) PriorityFormatType() types.FormatType           { return types.FormatUndefinedType }
func (it priceItem) Validate() error                                { return nil }
func (it priceItem) Context(ctx ...context.Context) context.Context { return nil }
func (it priceItem) Get(key string) interface{}                     { return nil }

func (it priceItem) MaxBid() billing.Money    { return it.maxBid }
func (it priceItem) Price() billing.Money     { return it.price }
func (it priceItem) LeadPrice() billing.Money { return it.leadPrice }
func (it priceItem) ActionBid() billing.Money { return it.price }
func (it priceItem) AuctionCPMBid() billing.Money {
	switch it.pricingModel {
	case types.PricingModelCPM:
		return it.price * 1000
	}
	return it.ecpm
}

func newResponseItem() ResponserItemCommon {
	return priceItem{}
}
