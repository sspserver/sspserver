package adsource

import (
	"encoding/json"

	"github.com/bsm/openrtb"
	openrtbnreq "github.com/bsm/openrtb/native/request"

	"geniusrabbit.dev/sspserver/internal/models/types"
)

// BidRequestRTBOptions of request build
type BidRequestRTBOptions struct {
	OpenNative struct {
		Ver string
	}
	FormatFilter func(f *types.Format) bool
	Currency     []string
}

func (opts *BidRequestRTBOptions) openNativeVer() string {
	return opts.OpenNative.Ver
}

func (opts *BidRequestRTBOptions) currencies() []string {
	if len(opts.Currency) > 0 {
		return opts.Currency
	}
	return []string{"USD"}
}

// BidRequestRTBOption set function
type BidRequestRTBOption func(opts *BidRequestRTBOptions)

// WithRTBOpenNativeVersion set version
func WithRTBOpenNativeVersion(ver string) BidRequestRTBOption {
	return func(opts *BidRequestRTBOptions) {
		opts.OpenNative.Ver = ver
	}
}

// WithFormatFilter set custom method
func WithFormatFilter(f func(f *types.Format) bool) BidRequestRTBOption {
	return func(opts *BidRequestRTBOptions) {
		opts.FormatFilter = f
	}
}

///////////////////////////////////////////////////////////////////////////////
/// BidRequest method extend
///////////////////////////////////////////////////////////////////////////////

// RTBBidRequest object
func (r *BidRequest) RTBBidRequest(opts ...BidRequestRTBOption) *RTBRequest {
	var opt BidRequestRTBOptions

	for _, fn := range opts {
		fn(&opt)
	}

	return &RTBRequest{
		RTBRequest: openrtb.BidRequest{
			ID:          r.ID,
			Imp:         r._RTBImpressions(&opt),
			Site:        r.SiteInfo().RTBObject(),
			App:         r.AppInfo().RTBObject(),
			Device:      r.DeviceInfo().RTBObject(r.UserInfo().Geo),
			User:        r.UserInfo().RTBObject(),
			AuctionType: 1,                // 1 = First Price, 2 = Second Price Plus
			TMax:        0,                // Maximum amount of time in milliseconds to submit a bid
			WSeat:       nil,              // Array of buyer seats allowed to bid on this auction
			AllImps:     0,                //
			Cur:         opt.currencies(), // Array of allowed currencies
			Bcat:        nil,              // Blocked Advertiser Categories
			BAdv:        nil,              // Array of strings of blocked toplevel domains of advertisers
			Regs:        nil,
			Ext:         nil,
		},
	}
}

func (r *BidRequest) _RTBImpressions(opts *BidRequestRTBOptions) (list []openrtb.Impression) {
	for _, imp := range r.Imps {
		for _, format := range imp.Formats() {
			if openRTBImp := r.impressionByFormat(&imp, format, opts); openRTBImp != nil {
				list = append(list, *openRTBImp)
			}
		}
	}
	return
}

func (r *BidRequest) impressionByFormat(imp *Impression, format *types.Format, opts *BidRequestRTBOptions) *openrtb.Impression {
	var (
		banner *openrtb.Banner
		video  *openrtb.Video
		native *openrtb.Native
		ext    openrtb.Extension
	)

	switch {
	case format.IsBanner() || format.IsProxy():
		w, h := imp.W, imp.H
		wm, wh := imp.WMax, imp.HMax
		if w < 1 && h < 1 {
			w, h = format.Width, format.Height
		}
		if !format.IsStretch() {
			wm, wh = 0, 0
		}
		banner = &openrtb.Banner{
			ID:       "",
			W:        w,
			H:        h,
			WMax:     wm,
			HMax:     wh,
			WMin:     0,
			HMin:     0,
			Pos:      imp.Pos,
			BType:    nil,
			BAttr:    nil,
			Mimes:    nil,
			TopFrame: 0,
			ExpDir:   nil,
			Api:      nil,
			Ext:      nil,
		}
	case format.IsNative():
		native = &openrtb.Native{
			Request: r.nativeRTBRequest(imp, format, opts),
			Ver:     opts.openNativeVer(),
			API:     nil,
			BAttr:   nil,
			Ext:     nil,
		}
	case format.IsDirect():
		ext = openrtb.Extension(`{"type":"pop"}`)
	default:
		return nil
	}

	tagid := imp.Target.Codename() + "_" + format.Codename
	return &openrtb.Impression{
		ID:                imp.IDByFormat(format),
		Banner:            banner,
		Video:             video,
		Native:            native,
		DisplayManager:    "",                                        // Name of ad mediation partner, SDK technology, etc
		DisplayManagerVer: "",                                        // Version of the above
		Instl:             b2i(imp.IsDirect()),                       // Interstitial, Default: 0 ("1": Interstitial, "0": Something else)
		TagID:             tagid,                                     // IDentifier for specific ad placement or ad tag
		BidFloor:          imp.BidFloor.Float64(),                    // Bid floor for this impression in CPM
		BidFloorCurrency:  "",                                        // Currency of bid floor
		Secure:            openrtb.NumberOrString(b2i(r.IsSecure())), // Flag to indicate whether the impression requires secure HTTPS URL creative assets and markup.
		IFrameBuster:      nil,                                       // Array of names for supportediframe busters.
		Pmp:               nil,                                       // A reference to the PMP object containing any Deals eligible for the impression object.
		Ext:               ext,
	}
}

func (r *BidRequest) nativeRTBRequest(imp *Impression, format *types.Format, opts *BidRequestRTBOptions) openrtb.Extension {
	var (
		nativePrepared []byte
		native         *openrtbnreq.Request
	)

	if native = imp.RTBNativeRequest(); native == nil {
		native = &openrtbnreq.Request{
			Ver:              opts.openNativeVer(),           // Version of the Native Markup
			LayoutID:         0,                              // DEPRECATED The Layout ID of the native ad
			AdUnitID:         0,                              // DEPRECATED The Ad unit ID of the native ad
			ContextTypeID:    imp.ContextType(),              // The context in which the ad appears
			ContextSubTypeID: imp.ContextSubType(),           // A more detailed context in which the ad appears
			PlacementTypeID:  imp.PlacementType(),            // The design/format/layout of the ad unit being offered
			PlacementCount:   imp.Count,                      // The number of identical placements in this Layout
			Sequence:         0,                              // 0 for the first ad, 1 for the second ad, and so on
			Assets:           r.nativeRTBAssets(imp, format), // An array of Asset Objects
			Ext:              nil,
		}
	}

	nativePrepared, _ = json.Marshal(native)

	// We have to encode it as a JSON string
	nativePrepared, _ = json.Marshal(`{"native":` + string(nativePrepared) + `}`)

	return openrtb.Extension(nativePrepared)
}

func (r *BidRequest) nativeRTBAssets(imp *Impression, format *types.Format) []openrtbnreq.Asset {
	assets := make([]openrtbnreq.Asset, 0, len(format.Config.Assets)+len(format.Config.Fields))
	for _, asset := range format.Config.Assets {
		if !asset.IsVideoSupport() || asset.IsImageSupport() {
			// By default we suppose that this is image
			var typeid openrtbnreq.ImageTypeID
			switch asset.Name {
			case types.FormatAssetMain:
				typeid = openrtbnreq.ImageTypeMain
			case types.FormatAssetIcon:
				typeid = openrtbnreq.ImageTypeIcon
			case "logo":
				typeid = openrtbnreq.ImageTypeLogo
			}
			assets = append(assets, openrtbnreq.Asset{
				ID:       int(asset.ID),
				Required: b2i(asset.Required),
				Image: &openrtbnreq.Image{
					TypeID:    typeid,
					WidthMin:  asset.MinWidth,
					HeightMin: asset.MinHeight,
					Mimes:     asset.AllowedTypes,
				},
			})
		} else {
			// TODO add video tag support
		}
	}

	for _, field := range format.Config.Fields {
		switch field.Name {
		case types.FormatFieldTitle:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Title:    &openrtbnreq.Title{Length: field.MaxLength()},
			})
		case types.FormatFieldDescription:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypeDesc,
					Length: field.MaxLength(),
				},
			})
		case types.FormatFieldBrandname:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypeSponsored,
					Length: field.MaxLength(),
				},
			})
		case types.FormatFieldPhone:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypePhone,
					Length: field.MaxLength(),
				},
			})
		case types.FormatFieldURL:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypeDisplayURL,
					Length: field.MaxLength(),
				},
			})
		case types.FormatFieldRating:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypeRating,
					Length: field.MaxLength(),
				},
			})
		case types.FormatFieldLikes:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypeLikes,
					Length: field.MaxLength(),
				},
			})
		case types.FormatFieldAddress:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypeAddress,
					Length: field.MaxLength(),
				},
			})
		case types.FormatFieldSponsored:
			assets = append(assets, openrtbnreq.Asset{
				ID:       field.ID,
				Required: b2i(field.Required),
				Data: &openrtbnreq.Data{
					TypeID: openrtbnreq.DataTypeSponsored,
					Length: field.MaxLength(),
				},
			})
		}
	}
	return assets
}
