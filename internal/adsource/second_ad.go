//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adsource

import "geniusrabbit.dev/sspserver/internal/billing"

// SecondAd description
type SecondAd struct {
	ID         string        `json:"id,omitempty"`
	Network    string        `json:"network,omitempty"`
	SourceID   uint64        `json:"source_id,omitempty"`
	CampaignID uint64        `json:"campaign_id,omitempty"`
	Price      billing.Money `json:"price,omitempty"`
}

// GetCampaignID value
func (s *SecondAd) GetCampaignID() uint64 {
	if s == nil {
		return 0
	}
	return s.CampaignID
}

// GetSourceID value
func (s *SecondAd) GetSourceID() uint64 {
	if s == nil {
		return 0
	}
	return s.SourceID
}

// GetPrice value
func (s *SecondAd) GetPrice() billing.Money {
	if s == nil {
		return 0
	}
	return s.Price
}

// GetECPM value
func (s *SecondAd) GetECPM() billing.Money {
	return 0
}
