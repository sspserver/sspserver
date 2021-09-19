//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package models

import (
	"testing"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

func Test_AdModel(t *testing.T) {
	ad := Ad{
		ID:          1,
		Content:     nil,
		Format:      &types.Format{},
		Bids:        nil,
		BidPrice:    billing.MoneyFloat(0.5),
		Price:       billing.MoneyFloat(1.5),
		LeadPrice:   billing.MoneyFloat(32.5),
		DailyBudget: billing.MoneyFloat(100),
		Budget:      billing.MoneyFloat(1000),
		Campaign:    &Campaign{},
		Hours:       nil,
	}

	ad.SetPricingModel(types.PricingModelCPA)
	ad.SetWeight(10)
	ad.SetFrequencyCapping(10)
	ad.AddFlag(AdFlagActive)
	// ad.State.SetSpent(10)

	t.Run("PricingModel", func(t *testing.T) {
		if ad.PricingModel != types.PricingModelCPA {
			t.Error("Wrong pricing model")
		}
	})

	t.Run("Weight", func(t *testing.T) {
		if ad.Weight != 10 {
			t.Error("Wrong weight of model")
		}
	})

	t.Run("FrequencyCapping", func(t *testing.T) {
		if ad.FrequencyCapping != 10 {
			t.Error("Wrong frequency capping of model")
		}
	})

	// t.Run("PricePerModel", func(t *testing.T) {
	// 	if ad.PricePerModel() != ad.LeadPrice {
	// 		t.Errorf("This is CPA campaign and price should be %.3f not %.3f",
	// 			ad.LeadPrice.Float64(), ad.PricePerModel().Float64())
	// 	}
	// })

	// t.Run("Budget", func(t *testing.T) {
	// 	if !ad.TestBudgetValues() {
	// 		t.Error("Budget must be valid")
	// 	}
	// })
}
