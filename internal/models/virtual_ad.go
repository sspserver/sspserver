//
// @project GeniusRabbit rotator 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package models

import (
	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// VirtualAd extract for targeting
type VirtualAd struct {
	Ad       *Ad
	Campaign *Campaign
	Bid      TargetBid
	Weight   float64
}

// CampaignObject reference
func (ad *VirtualAd) CampaignObject() *Campaign {
	return ad.Campaign
}

// ID value
func (ad *VirtualAd) ID() uint64 {
	return ad.Ad.ID
}

// PricingModel value
func (ad *VirtualAd) PricingModel() types.PricingModel {
	return ad.Ad.PricingModel
}

// ECPM value
func (ad *VirtualAd) ECPM() billing.Money {
	return ad.Bid.ECPM
}
