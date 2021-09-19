//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package models

// VirtualAds extract for targeting
type VirtualAds struct {
	Campaign *Campaign
	Bids     []TargetBid
}
