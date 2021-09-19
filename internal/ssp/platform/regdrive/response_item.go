package regdrive

import (
	"context"

	"github.com/demdxx/gocast"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/billing"
	imodels "geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// ResponseItem of the driver
type ResponseItem struct {
	imp      *adsource.Impression
	link     string
	data     map[string]interface{}
	bidPrice billing.Money
	source   adsource.Sourcer
	context  context.Context
}

// ID of current response item (unique code of current response)
func (r *ResponseItem) ID() string {
	return r.imp.ID
}

// Impression place object
func (r *ResponseItem) Impression() *adsource.Impression {
	return r.imp
}

// ImpressionID unique code string
func (r *ResponseItem) ImpressionID() string {
	return r.imp.ID
}

// ExtImpressionID it's unique code of the auction bid impression
func (r *ResponseItem) ExtImpressionID() string {
	return r.imp.ExtID
}

// ExtTargetID of the external network
func (r *ResponseItem) ExtTargetID() string {
	return r.imp.ExtTargetID
}

// AuctionCPMBid value price without any comission
// ! Important: There is no price for now
func (r *ResponseItem) AuctionCPMBid() billing.Money {
	return r.bidPrice
}

// PriorityFormatType from current Ad
func (r *ResponseItem) PriorityFormatType() types.FormatType {
	return types.FormatAutoregisterType
}

// Validate item
func (r *ResponseItem) Validate() error {
	return nil
}

// Context of response
func (r *ResponseItem) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		r.context = ctx[0]
	}
	return r.context
}

// Get context value
func (r *ResponseItem) Get(key string) interface{} {
	if r.context != nil {
		return r.context.Value(key)
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Implementation of adsource.ResponserItem interface{}
///////////////////////////////////////////////////////////////////////////////

// AdID number
func (r *ResponseItem) AdID() uint64 {
	return 0
}

// CompanyID number
func (r *ResponseItem) CompanyID() uint64 {
	type sourceCompanyIDGetr interface {
		CompanyID() uint64
	}
	return r.source.(sourceCompanyIDGetr).CompanyID()
}

// CampaignID number
func (r *ResponseItem) CampaignID() uint64 {
	return 0
}

// Format object
func (r *ResponseItem) Format() *types.Format {
	return r.imp.FormatByType(r.PriorityFormatType())
}

// PricingModel of advertisement
func (r *ResponseItem) PricingModel() types.PricingModel {
	return types.PricingModelCPM
}

// ContentItem returns the ad response data
func (r *ResponseItem) ContentItem(name string) (v interface{}) {
	switch name {
	case adsource.ContentItemLink:
		v = r.link
	default:
		if r.data != nil {
			v, _ = r.data[name]
		}
	}
	return
}

// ContentItemString from the ad
func (r *ResponseItem) ContentItemString(name string) (v string) {
	switch name {
	case adsource.ContentItemLink:
		v = r.link
	default:
		if r.data != nil {
			if vl, _ := r.data[name]; vl != nil {
				v = gocast.ToString(vl)
			}
		}
	}
	return
}

// ContentFields from advertisement object
func (r *ResponseItem) ContentFields() map[string]interface{} {
	return r.data
}

// ViewTrackerLinks returns traking links for view action
func (r *ResponseItem) ViewTrackerLinks() []string {
	return nil
}

// ClickTrackerLinks returns traking links for click action
func (r *ResponseItem) ClickTrackerLinks() []string {
	return nil
}

// MainAsset from response
func (r *ResponseItem) MainAsset() *imodels.AdFile { return nil }

// Asset by name
func (r *ResponseItem) Asset(name string) *imodels.AdFile { return nil }

// Assets list
func (r *ResponseItem) Assets() []imodels.AdFile { return nil }

// Source of response
func (r *ResponseItem) Source() adsource.Sourcer { return r.source }

// NetworkName by source
func (r *ResponseItem) NetworkName() string { return "" }

// Price per action
func (r *ResponseItem) Price(action imodels.Action) billing.Money {
	if action.IsImpression() {
		return r.bidPrice
	}
	return 0
}

// SetCPMPrice update of DSP auction value
func (r *ResponseItem) SetCPMPrice(price billing.Money, includeFactors ...bool) {}

// CPMPrice value price value for DSP auction
func (r *ResponseItem) CPMPrice(removeFactors ...bool) billing.Money { return r.bidPrice }

// Second campaigns
func (r *ResponseItem) Second() *adsource.SecondAd { return nil }

// Revenue money (in percents)
func (r *ResponseItem) Revenue() float64 {
	return r.RevenueShareFactor() * 100
}

// Potential money (in percents)
func (r *ResponseItem) Potential() float64 {
	return r.Source().RevenueShareReduceFactor() * 100
}

// ECPM item value
func (r *ResponseItem) ECPM() billing.Money { return r.bidPrice }

// RevenueShareFactor returns the multipler for company
// revenue calculation per action from 0 to 1
func (r *ResponseItem) RevenueShareFactor() float64 {
	return r.imp.RevenueShareFactor()
}

// ComissionShareFactor returns the multipler for commission
// calculation which system get from user revenue from 0 to 1
func (r *ResponseItem) ComissionShareFactor() float64 {
	return r.imp.ComissionShareFactor()
}

// IsDirect ad
func (r *ResponseItem) IsDirect() bool {
	return false
}

// ActionURL returns target resource link for direct and banner click as well
func (r *ResponseItem) ActionURL() string {
	return r.link
}

// Width of item
func (r *ResponseItem) Width() int { return 0 }

// Height of item
func (r *ResponseItem) Height() int { return 0 }

var (
	_ adsource.ResponserItem = (*ResponseItem)(nil)
)
