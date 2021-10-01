package adsource

import (
	"reflect"
	"testing"

	"github.com/bsm/openrtb"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

func Test_ItemPricing(t *testing.T) {
	var (
		comp = &models.Company{
			ID:           1,
			RevenueShare: 90,
		}
		imp   = Impression{Target: &models.Smartlink{Comp: comp}}
		items = []ResponserItem{newRTBResponse(comp, imp), newAdResponse(comp, imp)}
	)

	for _, item := range items {
		prefix := reflect.TypeOf(item).String()

		t.Run(prefix+"_empty_lead_price", func(t *testing.T) {
			if item.Price(models.ActionLead) != leadPrice(item) {
				t.Error("lead_price should be empty")
			}
		})

		t.Run(prefix+"_bid_price", func(t *testing.T) {
			if item.Price(models.ActionImpression) != billing.MoneyFloat(10) {
				t.Error("target price must be 10")
			}
		})

		t.Run(prefix+"_revenue_value", func(t *testing.T) {
			rev := item.RevenueShareFactor() * item.Price(models.ActionImpression).Float64()
			if rev != 9 {
				t.Errorf("wrong_revenue value: %.3f", rev)
			}
		})

		t.Run(prefix+"_comission_value", func(t *testing.T) {
			com := item.ComissionShareFactor() * item.Price(models.ActionImpression).Float64()
			if com != 1 {
				t.Errorf("wrong_comission value: %.3f", com)
			}
		})

		t.Run(prefix+"_potential_value", func(t *testing.T) {
			pot := item.Potential() * item.Price(models.ActionImpression).Float64()
			if pot != 0 {
				t.Errorf("wrong_potential value: %.3f", pot)
			}
		})

		t.Run(prefix+"_cpm_price", func(t *testing.T) {
			if item.CPMPrice() != billing.MoneyFloat(5) {
				t.Errorf("cpm_price value: 5 != %.3f", item.CPMPrice().Float64())
			}
		})
	}
}

func newRTBResponse(comp *models.Company, imp Impression) *ResponseBidItem {
	return &ResponseBidItem{
		ItemID:      "1",
		Src:         nil,
		Req:         &BidRequest{ID: "xxx", Imps: []Impression{imp}},
		Imp:         &imp,
		Bid:         &openrtb.Bid{Price: 60},
		BidPrice:    billing.MoneyFloat(10),
		CPMBidPrice: billing.MoneyFloat(5),
		SecondAd:    SecondAd{},
	}
}

func newAdResponse(comp *models.Company, imp Impression) *ResponseAdItem {
	return &ResponseAdItem{
		ItemID: "1",
		Src:    nil,
		Req:    &BidRequest{ID: "xxx", Imps: []Impression{imp}},
		Imp:    &imp,
		Ad: &models.Ad{
			Format:      &types.Format{Width: 250, Height: 250},
			BidPrice:    billing.MoneyFloat(5),
			Price:       billing.MoneyFloat(60),
			LeadPrice:   billing.MoneyFloat(120),
			DailyBudget: billing.MoneyFloat(1200),
			Budget:      billing.MoneyFloat(10000),
			Hours:       nil,
		},
		BidPrice: billing.MoneyFloat(10),
		SecondAd: SecondAd{},
	}
}

func leadPrice(item ResponserItem) billing.Money {
	switch it := item.(type) {
	case *ResponseAdItem:
		return it.Ad.LeadPrice
	}
	return 0
}
