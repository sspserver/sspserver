//
// @project geniusrabbit::rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package adsource

import (
	openrtbnreq "github.com/bsm/openrtb/native/request"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
	"geniusrabbit.dev/sspserver/internal/searchtypes"
)

// Impression target
type Impression struct {
	ID           string                 `json:"id,omitempty"`           // Internal impression ID
	ExtID        string                 `json:"extid,omitempty"`        // External impression ID (ImpID)
	ExtTargetID  string                 `json:"exttrgid"`               // External zone ID (tagid)
	Request      interface{}            `json:"request,omitempty"`      // Contains subrequest from RTB or another protocol
	Target       models.Target          `json:"target,omitempty"`       //
	BidFloor     billing.Money          `json:"bid_floor,omitempty"`    //
	SourcePrice  billing.Money          `json:"source_price,omitempty"` //
	Pos          int                    `json:"pos,omitempty"`          // 5.4 Ad Position
	Count        int                    `json:"cnt,omitempty"`          // Count of places for multiple banners
	X            int                    `json:"x,omitempty"`            // Position on the site screen
	Y            int                    `json:"y,omitempty"`            //
	W            int                    `json:"w,omitempty"`            //
	H            int                    `json:"h,omitempty"`            //
	WMax         int                    `json:"wm,omitempty"`           //
	HMax         int                    `json:"hm,omitempty"`           //
	SubID1       string                 `json:"subid1,omitempty"`       //
	SubID2       string                 `json:"subid2,omitempty"`       //
	SubID3       string                 `json:"subid3,omitempty"`       //
	SubID4       string                 `json:"subid4,omitempty"`       //
	SubID5       string                 `json:"subid5,omitempty"`       //
	Ext          map[string]interface{} `json:"ext,omitempty"`
	FormatTypes  types.FormatTypeBitset `json:"-"`
	formats      []*types.Format
	formatBitset *searchtypes.UIntBitset
}

// Init internal information
func (i *Impression) Init(formats types.FormatsAccessor) {
	var w, h, minw, minh = i.WMax, i.HMax, i.W, i.H
	if w <= 0 && h <= 0 {
		w, h = minw, minh
		minw, minh = minw-(minw/3), minh/3
	}
	if minw == 0 {
		minw = w - (w / 8)
	}
	if minh == 0 {
		minh = h - (h / 5)
	}

	i.formats = formats.FormatsBySize(w+10, h+10, minw, minh, i.FormatTypes)

	i.formatBitset = searchtypes.NewUIntBitset()
	for _, f := range i.formats {
		i.formatBitset.Set(uint(f.ID))
	}

	if i.FormatTypes.IsEmpty() {
		i.FormatTypes = *types.NewFormatTypeBitset().SetFromFormats(i.formats...)
	}
}

// Formats models
func (i *Impression) Formats() (f []*types.Format) {
	return i.formats
}

// FormatByType of formats
func (i *Impression) FormatByType(tp types.FormatType) *types.Format {
	for _, f := range i.formats {
		if f.Types.Is(tp) {
			return f
		}
	}
	return nil
}

// FormatBitset of IDs
func (i *Impression) FormatBitset() *searchtypes.UIntBitset {
	return i.formatBitset
}

// IDByFormat return specific ID to link format
func (i *Impression) IDByFormat(format *types.Format) string {
	return i.ID + "_" + format.Codename
}

// TargetID value
func (i Impression) TargetID() uint {
	if i.Target == nil {
		return 0
	}
	return uint(i.Target.ID())
}

// CompanyID number
func (i *Impression) CompanyID() uint64 {
	if i.Target != nil {
		return i.Target.CompanyID()
	}
	return 0
}

// IsDirect value
func (i *Impression) IsDirect() bool {
	return i.FormatTypes.Is(types.FormatDirectType)
}

// IsNative target support
func (i *Impression) IsNative() bool {
	return i.FormatTypes.Is(types.FormatNativeType)
}

// IsStandart target support
func (i *Impression) IsStandart() bool {
	return false ||
		i.FormatTypes.Is(types.FormatBannerType) ||
		i.FormatTypes.Is(types.FormatBannerHTML5Type)
}

// RevenueShareFactor value
func (i *Impression) RevenueShareFactor() float64 {
	if i == nil || i.Target == nil {
		return 0
	}
	return i.Target.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (i *Impression) ComissionShareFactor() float64 {
	if i == nil || i.Target == nil {
		return 0
	}
	return i.Target.ComissionShareFactor()
}

// PurchasePrice gives the price of view from external resource.
// The cost of this request.
func (i *Impression) PurchasePrice(action models.Action) billing.Money {
	if action.IsImpression() && i.SourcePrice > 0 {
		return i.SourcePrice
	}
	if i.Target != nil {
		return i.Target.PurchasePrice(action)
	}
	return 0
}

///////////////////////////////////////////////////////////////////////////////
/// OpenRTB methods
///////////////////////////////////////////////////////////////////////////////

// ContextType IDs 7.3
// @link https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-Native-Ads-Specification-1-1_2016.pdf
func (i *Impression) ContextType() openrtbnreq.ContextTypeID {
	return openrtbnreq.ContextTypeContent
}

// ContextSubType IDs 7.4
// @link https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-Native-Ads-Specification-1-1_2016.pdf
func (i *Impression) ContextSubType() openrtbnreq.ContextSubTypeID {
	return openrtbnreq.ContextSubTypeGeneral
}

// PlacementType IDs 7.5
// @link https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-Native-Ads-Specification-1-1_2016.pdf
func (i *Impression) PlacementType() openrtbnreq.PlacementTypeID {
	return openrtbnreq.PlacementTypeRecommendation
}

// RTBNativeRequest object
func (i *Impression) RTBNativeRequest() (r *openrtbnreq.Request) {
	r, _ = i.Request.(*openrtbnreq.Request)
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Ext data methods
///////////////////////////////////////////////////////////////////////////////

// Get context item by key
func (i *Impression) Get(key string) interface{} {
	if i.Ext == nil {
		return nil
	}
	v, _ := i.Ext[key]
	return v
}

// Set context item with key
func (i *Impression) Set(key string, val interface{}) {
	if i.Ext == nil {
		i.Ext = map[string]interface{}{}
	}
	i.Ext[key] = val
}

// Unset context item with keys
func (i *Impression) Unset(keys ...string) {
	if i.Ext == nil {
		return
	}

	for _, key := range keys {
		delete(i.Ext, key)
	}
}
