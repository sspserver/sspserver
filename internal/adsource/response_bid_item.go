//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adsource

import (
	"context"
	"strings"

	"github.com/demdxx/gocast"

	"github.com/bsm/openrtb"
	natresp "github.com/bsm/openrtb/native/response"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

type responseDataAccessor interface {
	Get(key string) interface{}
}

// ResponseBidItem value
type ResponseBidItem struct {
	ItemID      string
	Src         Sourcer
	Req         *BidRequest
	Imp         *Impression
	FormatType  types.FormatType
	RespFormat  *types.Format
	Bid         *openrtb.Bid
	Native      *natresp.Response
	Data        responseDataAccessor
	BidPrice    billing.Money
	CPMBidPrice billing.Money // This param can update only price predictor
	SecondAd    SecondAd
	assets      []models.AdFile
	context     context.Context
}

// ID of current response item (unique code of current response)
func (it *ResponseBidItem) ID() string {
	return it.ItemID
}

// Source of response
func (it *ResponseBidItem) Source() Sourcer {
	return it.Src
}

// NetworkName by source
func (it *ResponseBidItem) NetworkName() string {
	return ""
}

// AccessPoint of response
func (it *ResponseBidItem) AccessPoint() AccessPoint {
	if it.Req == nil || it.Req.AccessPoint == nil {
		return nil
	}
	return it.Req.AccessPoint
}

// ContentItemString from the ad
func (it *ResponseBidItem) ContentItemString(name string) string {
	val := it.ContentItem(name)
	if val != nil {
		return gocast.ToString(val)
	}
	return ""
}

// ContentItem returns the ad response data
func (it *ResponseBidItem) ContentItem(name string) interface{} {
	if it.Data != nil {
		return it.Data.Get(name)
	}

	formatType := it.PriorityFormatType()

	switch name {
	case ContentItemContent, ContentItemIFrameURL:
		if formatType.IsBanner() {
			switch name {
			case ContentItemIFrameURL:
				if strings.HasPrefix(it.Bid.AdMarkup, "http://") ||
					strings.HasPrefix(it.Bid.AdMarkup, "https://") ||
					(strings.HasPrefix(it.Bid.AdMarkup, "//") && !strings.ContainsAny(it.Bid.AdMarkup, "\n\t")) {
					return it.Bid.AdMarkup
				}
			case ContentItemContent:
				return it.Bid.AdMarkup
			}
		}
	case ContentItemLink:
		switch {
		case it.Native != nil:
			return it.Native.Link.URL
		case formatType.IsDirect():
			// In this case here have to be the advertisement link
			return it.Bid.AdMarkup
		}
	case ContentItemNotifyWinURL:
		if it.Bid != nil {
			return it.Bid.NURL
		}
	case ContentItemNotifyDisplayURL:
		if it.Bid != nil {
			return it.Bid.BURL
		}
	case types.FormatFieldTitle:
		if it.Native != nil {
			for _, asset := range it.Native.Assets {
				if asset.Title != nil {
					return asset.Title.Text
				}
			}
		}
	default:
		if it.Native != nil {
			for _, asset := range it.Native.Assets {
				if asset.Data != nil && asset.Data.Label == name {
					return asset.Data.Value
				}
			}
		}
	}
	return nil
}

// ContentFields from advertisement object
func (it *ResponseBidItem) ContentFields() map[string]interface{} {
	fields := map[string]interface{}{}
	config := it.Format().Config
	for _, field := range config.Fields {
		for _, asset := range it.Native.Assets {
			if field.ID != asset.ID {
				continue
			}
			switch {
			case asset.Title != nil:
				fields[field.Name] = asset.Title.Text
			case asset.Link != nil:
				fields[field.Name] = asset.Link.URL
			case asset.Data != nil:
				fields[field.Name] = asset.Data.Value
			}
			break
		}
	}
	return fields
}

// ViewTrackerLinks returns traking links for view action
func (it *ResponseBidItem) ViewTrackerLinks() []string {
	if it.Native == nil {
		return nil
	}
	return it.Native.ImpTrackers
}

// ClickTrackerLinks returns third-party tracker URLs to be fired on click of the URL
func (it *ResponseBidItem) ClickTrackerLinks() []string {
	if it.Native == nil {
		return nil
	}
	return it.Native.Link.ClickTrackers
}

// MainAsset from response
func (it *ResponseBidItem) MainAsset() *models.AdFile {
	mainAsset := it.Format().Config.MainAsset()
	if mainAsset == nil {
		return nil
	}
	for _, asset := range it.Assets() {
		if int(asset.ID) == mainAsset.ID {
			return &asset
		}
	}
	return nil
}

// Asset by name
func (it *ResponseBidItem) Asset(name string) *models.AdFile {
	for _, asset := range it.Assets() {
		if asset.Name == name {
			return &asset
		}
	}
	return nil
}

// Assets list
func (it *ResponseBidItem) Assets() (assets []models.AdFile) {
	if it.assets != nil {
		return it.assets
	}

	config := it.Format().Config
	for _, configAsset := range config.Assets {
		for _, asset := range it.Native.Assets {
			if asset.ID != configAsset.ID {
				continue
			}
			newAsset := models.AdFile{
				ID:   uint64(asset.ID),
				Name: configAsset.GetName(),
			}
			switch {
			case asset.Image != nil:
				newAsset.Path = asset.Image.URL
				newAsset.Type = models.AdFileImageType
				newAsset.ContentType = ""
				newAsset.Width = asset.Image.Width
				newAsset.Height = asset.Image.Height
			// case asset.Video != nil:
			// 	newAsset.Path = asset.Video.URL
			// 	newAsset.Type = models.AdFileVideoType
			default:
				// TODO error generation
			}
			it.assets = append(it.assets, newAsset)
			break
		}
	}
	return it.assets
}

// Format object model
func (it *ResponseBidItem) Format() *types.Format {
	if it == nil {
		return nil
	}
	return it.RespFormat
}

// PriorityFormatType from current Ad
func (it *ResponseBidItem) PriorityFormatType() types.FormatType {
	if it.FormatType != types.FormatUndefinedType {
		return it.FormatType
	}
	format := it.Imp.FormatTypes
	if formatType := format.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}
	return format.FirstType()
}

// Impression place object
func (it *ResponseBidItem) Impression() *Impression {
	return it.Imp
}

// ImpressionID unique code string
func (it *ResponseBidItem) ImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ID
}

// ExtImpressionID unique code of RTB response
func (it *ResponseBidItem) ExtImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ExtID
}

// ExtTargetID of the external network
func (it *ResponseBidItem) ExtTargetID() string {
	return it.Imp.ExtTargetID
}

// AdID number
func (it *ResponseBidItem) AdID() uint64 {
	return 0
}

// CompanyID number
func (it *ResponseBidItem) CompanyID() uint64 {
	if it.Src != nil {
		type companyIDGetter interface {
			CompanyID() uint64
		}
		if src, _ := it.Src.(companyIDGetter); src != nil {
			return src.CompanyID()
		}
	}
	return 0
}

// CampaignID number
func (it *ResponseBidItem) CampaignID() uint64 {
	return 0
}

// PricingModel of advertisement
func (it *ResponseBidItem) PricingModel() types.PricingModel {
	return types.PricingModelCPM
}

// ECPM item value
func (it *ResponseBidItem) ECPM() billing.Money {
	if it == nil || it.Bid == nil {
		return 0
	}
	if it.PricingModel().IsCPM() && it.CPMBidPrice > 0 {
		return it.CPMBidPrice
	}
	return billing.MoneyFloat(it.Bid.Price)
}

// Price summ
func (it *ResponseBidItem) Price(action models.Action) billing.Money {
	if it == nil || it.Bid == nil {
		return 0
	}
	if action.IsImpression() {
		if it.BidPrice > 0 {
			return it.BidPrice
		}
		return billing.MoneyFloat(it.Bid.Price / 1000)
	}
	return 0
}

// SetCPMPrice update of DSP auction value
func (it *ResponseBidItem) SetCPMPrice(price billing.Money, includeFactors ...bool) {
	if len(includeFactors) > 0 && includeFactors[0] {
		price = it.PreparePrice(price, false)
	}
	if it != nil && price < it.ECPM() {
		it.CPMBidPrice = price
	}
}

// CPMPrice value price value for DSP auction
func (it *ResponseBidItem) CPMPrice(removeFactors ...bool) (price billing.Money) {
	if it.PricingModel() == types.PricingModelCPM {
		price = it.Price(models.ActionImpression) * 1000
	} else {
		price = it.ECPM()
	}
	if it.CPMBidPrice > 0 && it.CPMBidPrice < price {
		price = it.CPMBidPrice
	}
	if len(removeFactors) > 0 && removeFactors[0] {
		price = it.PreparePrice(price, true)
	}
	return price
}

// AuctionCPMBid value price without any comission
func (it *ResponseBidItem) AuctionCPMBid() billing.Money {
	return it.CPMPrice()
}

// Second campaigns
func (it *ResponseBidItem) Second() *SecondAd {
	return &it.SecondAd
}

// Revenue money (in percents)
func (it *ResponseBidItem) Revenue() float64 {
	return it.RevenueShareFactor() * 100
}

// Potential money (in percents)
func (it *ResponseBidItem) Potential() float64 {
	return it.Source().RevenueShareReduceFactor() * 100
}

// RevenueShareFactor value
func (it *ResponseBidItem) RevenueShareFactor() float64 {
	return it.Imp.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher 0..1
func (it *ResponseBidItem) ComissionShareFactor() float64 {
	return it.Imp.ComissionShareFactor()
}

// IsDirect ad
func (it *ResponseBidItem) IsDirect() bool {
	return it.Imp.IsDirect()
}

// ActionURL for direct ads
func (it *ResponseBidItem) ActionURL() string {
	if it.Native == nil {
		return ""
	}
	return it.Native.Link.URL
}

// Validate item
func (it *ResponseBidItem) Validate() error {
	if it.Src == nil || it.Req == nil || it.Imp == nil || it.Bid == nil {
		return ErrInvalidItemInitialisation
	}
	return it.Bid.Validate()
}

// Width of item
func (it *ResponseBidItem) Width() int {
	return 0
}

// Height of item
func (it *ResponseBidItem) Height() int {
	return 0
}

// Markup advertisement
func (it *ResponseBidItem) Markup() (string, error) {
	return "", nil
}

// Context value
func (it *ResponseBidItem) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		it.context = ctx[0]
	}
	return it.context
}

// Get ext field
func (it *ResponseBidItem) Get(key string) (res interface{}) {
	if it.context == nil {
		return res
	}
	return it.context.Value(key)
}

// PreparePrice value
func (it *ResponseBidItem) PreparePrice(price billing.Money, removeFactors bool) billing.Money {
	fc := it.ComissionShareFactor() + it.Source().RevenueShareReduceFactor()
	if removeFactors {
		price -= price * billing.MoneyFloat(fc)
	} else {
		price += price * billing.MoneyFloat(fc)
	}
	return price
}

var (
	_ ResponserItem = &ResponseBidItem{}
)
