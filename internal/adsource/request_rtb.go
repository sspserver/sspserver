//
// @project geniusrabbit::rotator 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//

package adsource

import (
	"context"
	"encoding/json"
	"io"
	"net"

	"github.com/bsm/openrtb"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/infostructs"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
	"geniusrabbit.dev/sspserver/internal/personification"
	"geniusrabbit.dev/sspserver/internal/rand"
	openrtbnreq "github.com/bsm/openrtb/native/request"
)

var (
	emptyRtbGeo      = openrtb.Geo{}
	emptyRtbDevice   = openrtb.Device{}
	emptyRtbUser     = openrtb.User{Geo: &emptyRtbGeo}
	emptyOsInfo      = infostructs.OS{}
	emptyBrowserInfo = infostructs.Browser{}
	emptyCompany     = models.Company{}
	emptyZone        = models.Zone{Comp: &emptyCompany}
)

// RTBRequest object
type RTBRequest struct {
	ID         string
	RTBRequest openrtb.BidRequest
	Context    context.Context
}

// Load from reader
func (r *RTBRequest) Load(reader io.Reader) error {
	err := json.NewDecoder(reader).Decode(&r.RTBRequest)
	if err == nil {
		r.UpdateState()
	}
	return err
}

// LoadBytes from array
func (r *RTBRequest) LoadBytes(data []byte) (err error) {
	if err = json.Unmarshal(data, &r.RTBRequest); nil == err {
		r.UpdateState()
	}
	return
}

// UpdateState params
func (r *RTBRequest) UpdateState() {
	if len(r.ID) < 1 {
		r.ID = rand.UUID()
	}

	// Fix reuest omission
	if r.RTBRequest.Device == nil {
		r.RTBRequest.Device = &emptyRtbDevice
	}

	if r.RTBRequest.User == nil {
		r.RTBRequest.User = &emptyRtbUser
	}

	if r.RTBRequest.User.Geo == nil {
		r.RTBRequest.User.Geo = r.RTBRequest.Device.Geo
	}

	// TODO: check IP, even if we have geo targeting.
	// Because IP database could be different
	if r.RTBRequest.User.Geo == nil {
		r.RTBRequest.User.Geo = &emptyRtbGeo
		r.RTBRequest.Device.Geo = r.RTBRequest.User.Geo
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Getters
///////////////////////////////////////////////////////////////////////////////

// UserID string
func (r *RTBRequest) UserID() string {
	var uuid = r.RTBRequest.User.BuyerUID
	if len(uuid) < 1 {
		uuid = r.RTBRequest.User.ID
	}

	// @NOTE: We have restriction for 36 chars
	if len(uuid) > 36 {
		uuid = uuid[:36]
	}
	return uuid
}

// Keywords for request
func (r *RTBRequest) Keywords() (words string) {
	if r.RTBRequest.User != nil {
		words = r.RTBRequest.User.Keywords
	}
	if len(words) < 1 && r.RTBRequest.Site != nil && r.RTBRequest.Site.Content != nil {
		words = r.RTBRequest.Site.Content.Keywords
	}
	return
}

// ImpressionByID object
func (r *RTBRequest) ImpressionByID(id string) *openrtb.Impression {
	for i, im := range r.RTBRequest.Imp {
		if id == im.ID {
			return &r.RTBRequest.Imp[i]
		}
	}
	return nil
}

// BidRequest from RTB request
func (r *RTBRequest) BidRequest(formats types.FormatsAccessor, person personification.Person, defaultTarget models.Target, accessPoint ...*models.RTBAccessPoint) (request *BidRequest, err error) {
	if err = r.Validate([]string{"USD"}); err != nil {
		return nil, err
	}

	if person == nil {
		person = personification.EmptyPerson
	}

	request = &BidRequest{
		ID:     r.ID,
		ExtID:  r.RTBRequest.ID,
		Person: person,
		App:    r.app(),
		Site:   r.site(),
		Device: r.device(person),
		User:   r.user(person),
	}

	for _, imp := range r.RTBRequest.Imp {
		if _imp := r.imp(defaultTarget, &imp); _imp != nil {
			request.Imps = append(request.Imps, *_imp)
		}
	}

	if len(accessPoint) > 0 && accessPoint[0] != nil {
		request.AccessPoint = accessPoint[0]
	}
	request.Init(formats)
	return request, nil
}

///////////////////////////////////////////////////////////////////////////////
/// Converters
///////////////////////////////////////////////////////////////////////////////

func (r *RTBRequest) app() *infostructs.App {
	if r.RTBRequest.App == nil {
		return nil
	}

	return &infostructs.App{
		ExtID:         r.RTBRequest.App.ID,                 // External ID
		Bundle:        r.RTBRequest.App.Bundle,             // App bundle or package name
		Cat:           r.RTBRequest.App.Cat,                // Array of categories
		PrivacyPolicy: r.RTBRequest.App.GetPrivacyPolicy(), // Default: 1 ("1": has a privacy policy)
		StoreURL:      r.RTBRequest.App.StoreURL,           // App store URL for an installed app
		Ver:           r.RTBRequest.App.Ver,                // App version
		Paid:          r.RTBRequest.App.Paid,               // "1": Paid, "2": Free
	}
}

func (r *RTBRequest) site() *infostructs.Site {
	if r.RTBRequest.Site == nil {
		return nil
	}

	return &infostructs.Site{
		ExtID:         r.RTBRequest.Site.ID,                 // External ID
		Domain:        r.RTBRequest.Site.Domain,             //
		Cat:           r.RTBRequest.Site.Cat,                // Array of categories
		PrivacyPolicy: r.RTBRequest.Site.GetPrivacyPolicy(), // Default: 1 ("1": has a privacy policy)
		Keywords:      r.RTBRequest.Site.Keywords,           // Comma separated list of keywords about the site.
		Page:          r.RTBRequest.Site.Page,               // URL of the page
		Ref:           r.RTBRequest.Site.Ref,                // Referrer URL
		Search:        r.RTBRequest.Site.Search,             // Search string that caused naviation
		Mobile:        r.RTBRequest.Site.Mobile,             // Mobile ("1": site is mobile optimised)
	}
}

func (r *RTBRequest) device(person personification.Person) *infostructs.Device {
	var (
		os      = r.os(person)
		browser = r.browser(person)
		// info    = person.UserInfo()
	)

	if r.RTBRequest.Device == nil {
		return &infostructs.Device{
			// Make:       info.Device.Make,       // Device make
			// Model:      info.Device.Model,      // Device model
			OS:      os,      // Device OS
			Browser: browser, // Device OS version
			// ConnType:   info.Device.ConnType,   //
			// DeviceType: info.Device.DeviceType, //
			// IFA:        info.Device.IFA,        // Native identifier for advertisers
			Height: 0, // Physical height of the screen in pixels.
			Width:  0, // Physical width of the screen in pixels.
			// PPI:        info.Device.PPI,        // Screen size as pixels per linear inch.
			// PxRatio:    info.Device.PxRatio,    // The ratio of physical pixels to device independent pixels.
			// HwVer:      info.Device.HwVer,      // Hardware version of the device (e.g., "5S" for iPhone 5S).
		}
	}

	devType := infostructs.DeviceType(r.RTBRequest.Device.DeviceType)
	return &infostructs.Device{
		Make:       r.RTBRequest.Device.Make,     // Device make
		Model:      r.RTBRequest.Device.Model,    // Device model
		OS:         os,                           // Device OS
		Browser:    browser,                      // Device OS version
		ConnType:   r.RTBRequest.Device.ConnType, //
		DeviceType: devType,                      //
		IFA:        r.RTBRequest.Device.IFA,      // Native identifier for advertisers
		Height:     r.RTBRequest.Device.H,        // Physical height of the screen in pixels.
		Width:      r.RTBRequest.Device.W,        // Physical width of the screen in pixels.
		PPI:        r.RTBRequest.Device.PPI,      // Screen size as pixels per linear inch.
		PxRatio:    r.RTBRequest.Device.PxRatio,  // The ratio of physical pixels to device independent pixels.
		HwVer:      r.RTBRequest.Device.HwVer,    // Hardware version of the device (e.g., "5S" for iPhone 5S).
	}
}

func (r *RTBRequest) browser(person personification.Person) *infostructs.Browser {
	// var info = person.UserInfo()

	// if r.RTBRequest.Device == nil {
	// 	return info.Device.Browser
	// }

	return &infostructs.Browser{
		// ID:              info.Device.Browser.ID,       // Internal system ID
		// Name:            info.Device.Browser.Name,     //
		// Version:         info.Device.Browser.Version,  //
		DNT:             r.RTBRequest.Device.DNT,      // "1": Do not track
		LMT:             r.RTBRequest.Device.LMT,      // "1": Limit Ad Tracking
		UA:              r.RTBRequest.Device.UA,       // User agent
		JS:              r.RTBRequest.Device.JS,       //
		PrimaryLanguage: r.RTBRequest.Device.Language, // Browser language (en-US)
		FlashVer:        r.RTBRequest.Device.FlashVer, // Flash version
		Width:           r.RTBRequest.Device.W,        // Window in pixels
		Height:          r.RTBRequest.Device.H,        // Window in pixels
	}
}

func (r *RTBRequest) os(person personification.Person) *infostructs.OS {
	// var info = person.UserInfo()

	// if info.Device.OS.ID > 0 || r.RTBRequest.Device == nil {
	// 	return info.Device.OS
	// }

	return &infostructs.OS{
		ID:      0,
		Name:    r.RTBRequest.Device.OS,
		Version: r.RTBRequest.Device.OSVer,
	}
}

func (r *RTBRequest) user(person personification.Person) *User {
	var (
		user = r.RTBRequest.User
		// info             = person.UserInfo()
		// ageStart, ageEnd = info.Ages()
	)

	if user == nil {
		user = &emptyRtbUser
	}

	return &User{
		// ID:            info.UUID(),                // Unique User ID
		// SessionID:     info.SessionID(),           // Unique session ID
		// FingerPrintID: info.Fingerprint(),         //
		// ETag:          info.ETag(),                //
		// AgeStart:      ageStart,                   // Year of birth from
		// AgeEnd:        ageEnd,                     // Year of birth from
		Gender:   user.Gender,   // Gender ("M": male, "F" female, "O" Other)
		Keywords: user.Keywords, // Comma separated list of keywords, interests, or intent
		Geo:      r.geo(person),
	}
}

func (r *RTBRequest) geo(person personification.Person) *infostructs.Geo {
	var (
		ip net.IP
		// info = person.UserInfo()
		geo infostructs.Geo
	)

	if r.RTBRequest.Device != nil {
		ip = net.ParseIP(r.RTBRequest.Device.IP)
	}
	// if ip == nil {
	// 	ip = info.Geo.IP
	// }

	if r.RTBRequest.User != nil && r.RTBRequest.User.Geo != nil {
		geo = infostructs.Geo{
			// ID:            info.Geo.ID,                                                         // Internal geo ID
			IP:      ip,                // IPv4/6
			Carrier: r.carrier(person), // Carrier or ISP derived from the IP address
			// Lat:           defFloat(info.Geo.Lat, r.RTBRequest.User.Geo.Lat),                   // Latitude from -90 to 90
			// Lon:           defFloat(info.Geo.Lon, r.RTBRequest.User.Geo.Lon),                   // Longitude from -180 to 180
			// Country:       defStr(info.Geo.Country, r.RTBRequest.User.Geo.Country),             // Country using ISO 3166-1 Alpha 2
			// Region:        defStr(info.Geo.Region, r.RTBRequest.User.Geo.Region),               // Region using ISO 3166-2
			// RegionFIPS104: defStr(info.Geo.RegionFIPS104, r.RTBRequest.User.Geo.RegionFIPS104), // Region of a country using FIPS 10-4
			// Metro:         defStr(info.Geo.Metro, r.RTBRequest.User.Geo.Metro),                 //
			// City:          defStr(info.Geo.City, r.RTBRequest.User.Geo.City),                   //
			// Zip:           defStr(info.Geo.Zip, r.RTBRequest.User.Geo.Zip),                     //
			// UTCOffset:     defInt(info.Geo.UTCOffset, r.RTBRequest.User.Geo.UTCOffset),         // Local time as the number +/- of minutes from UTC
		}
	} else {
		geo = infostructs.Geo{
			// ID:            info.Geo.ID,            // Internal geo ID
			IP:      ip,                // IPv4/6
			Carrier: r.carrier(person), // Carrier or ISP derived from the IP address
			// Lat:           info.Geo.Lat,           // Latitude from -90 to 90
			// Lon:           info.Geo.Lon,           // Longitude from -180 to 180
			// Country:       info.Geo.Country,       // Country using ISO 3166-1 Alpha 2
			// Region:        info.Geo.Region,        // Region using ISO 3166-2
			// RegionFIPS104: info.Geo.RegionFIPS104, // Region of a country using FIPS 10-4
			// Metro:         info.Geo.Metro,         //
			// City:          info.Geo.City,          //
			// Zip:           info.Geo.Zip,           //
			// UTCOffset:     info.Geo.UTCOffset,     // Local time as the number +/- of minutes from UTC
		}
	}

	return &geo
}

func (r *RTBRequest) carrier(person personification.Person) *infostructs.Carrier {
	// carrier := person.UserInfo().CarrierInfo()
	// return &infostructs.Carrier{
	// 	ID:   carrier.ID,
	// 	Name: carrier.Name,
	// 	Code: carrier.Code,
	// }
	return &infostructs.Carrier{}
}

func (r *RTBRequest) imp(target models.Target, imp *openrtb.Impression) *Impression {
	var (
		ext  map[string]string
		nimp = &Impression{
			ID:          rand.UUID(),
			ExtID:       imp.ID,
			ExtTargetID: imp.TagID,
			Target:      target,
			BidFloor:    billing.MoneyFloat(imp.BidFloor),
		}
	)

	if nimp.Target == nil {
		nimp.Target = &emptyZone
	}

	if len(imp.Ext) > 2 {
		_ = json.Unmarshal(imp.Ext, &ext)
	}

	// Prepare params by ext
	if ext != nil && len(ext) > 0 {
		if nimp.SubID1, _ = ext["subid1"]; nimp.SubID1 == "" {
			nimp.SubID1, _ = ext["subid"]
		}
		nimp.SubID2, _ = ext["subid2"]
		nimp.SubID3, _ = ext["subid3"]
		nimp.SubID4, _ = ext["subid4"]
		nimp.SubID5, _ = ext["subid5"]
	}

	switch {
	// Banner Ad information
	case imp.Banner != nil:
		nimp.W, nimp.H = imp.Banner.W, imp.Banner.H
		nimp.WMax, nimp.HMax = imp.Banner.WMax, imp.Banner.HMax
		nimp.Pos = imp.Banner.Pos

		// @NOTE Hack for direct Ad!!!
		if nimp.W <= 0 && nimp.H <= 0 {
			nimp.FormatTypes.Reset().Set(types.FormatDirectType)
		} else {
			nimp.FormatTypes.Reset().Set(types.FormatBannerType).
				Set(types.FormatBannerHTML5Type)
		}

		// Rewrite type if it's passed
		if !nimp.IsDirect() && ext != nil && len(ext) > 0 {
			if tp, _ := ext["type"]; tp == "pop" || tp == "direct" {
				nimp.FormatTypes.Reset().Set(types.FormatDirectType)
			}
		}

	// Native Ad information
	case imp.Native != nil:
		if len(imp.Native.Request) < 1 {
			return nil
		}

		var (
			req         = new(openrtbnreq.Request)
			requestBody = imp.Native.Request
		)
		if imp.Native.Request[0] == '"' {
			var respStr string
			if err := json.Unmarshal(requestBody, &respStr); err != nil {
				return nil
			}
			requestBody = []byte(respStr)
		}
		if err := json.Unmarshal(requestBody, req); err != nil {
			return nil
		}

		nimp.Count = req.PlacementCount
		nimp.FormatTypes.Reset().Set(types.FormatNativeType)
		nimp.Request = req

		for _, a := range req.Assets {
			if a.Image != nil {
				if a.Image.TypeID == openrtbnreq.ImageTypeMain {
					nimp.W, nimp.H = a.Image.Width, a.Image.Height
					nimp.WMax, nimp.HMax = a.Image.WidthMin, a.Image.HeightMin
					break
				}
			}
		}

	// Video Ad information
	case imp.Video != nil:
		// It's not supported at the moment
		return nil
	default:
		nimp.FormatTypes.Reset().Set(types.FormatBannerType)
	}

	return nimp
}

func (r *RTBRequest) reset() {
	r.RTBRequest.Imp = r.RTBRequest.Imp[:0]
	r.RTBRequest.BSeat = r.RTBRequest.BSeat[:0]
	r.RTBRequest.WSeat = r.RTBRequest.WSeat[:0]
	r.RTBRequest.WLang = r.RTBRequest.WLang[:0]
	r.RTBRequest.Cur = r.RTBRequest.Cur[:0]
	r.RTBRequest.Bcat = r.RTBRequest.Bcat[:0]
	r.RTBRequest.BAdv = r.RTBRequest.BAdv[:0]
	r.RTBRequest.BApp = r.RTBRequest.BApp[:0]
	if r.RTBRequest.Source != nil {
		r.RTBRequest.Source.TransactionID = ""
		r.RTBRequest.Source.PaymentChain = ""
		r.RTBRequest.Source.Ext = r.RTBRequest.Source.Ext[:0]
	}
	if r.RTBRequest.User != nil {
		r.RTBRequest.User.Data = r.RTBRequest.User.Data[:0]
		if r.RTBRequest.User.Geo != nil {
			r.RTBRequest.User.Geo.Ext = r.RTBRequest.User.Geo.Ext[:0]
		}
	}
	if r.RTBRequest.Site != nil {
		r.RTBRequest.Site.Cat = r.RTBRequest.Site.Cat[:0]
		r.RTBRequest.Site.SectionCat = r.RTBRequest.Site.SectionCat[:0]
		r.RTBRequest.Site.PageCat = r.RTBRequest.Site.PageCat[:0]
		if r.RTBRequest.Site.Publisher != nil {
			r.RTBRequest.Site.Publisher.Cat = r.RTBRequest.Site.Publisher.Cat[:0]
			r.RTBRequest.Site.Publisher.Ext = r.RTBRequest.Site.Publisher.Ext[:0]
		}
		if r.RTBRequest.Site.Content != nil {
			r.RTBRequest.Site.Content.Cat = r.RTBRequest.Site.Content.Cat[:0]
			r.RTBRequest.Site.Content.Ext = r.RTBRequest.Site.Content.Ext[:0]
			r.RTBRequest.Site.Content.Data = r.RTBRequest.Site.Content.Data[:0]
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Validation
///////////////////////////////////////////////////////////////////////////////

// Validate request by currency
func (r *RTBRequest) Validate(currency []string) (err error) {
	// Check currency
	// @NOTE USD currency in hardcore
	if len(r.RTBRequest.Cur) > 0 && indexOfStringArray("USD", r.RTBRequest.Cur) < 0 {
		return ErrInvalidCur
	}

	for _, imp := range r.RTBRequest.Imp {
		// Video Ad is not supported
		if imp.Video != nil {
			return NoSupportError{"video"}
		}

		// Check bid currency
		if imp.BidFloorCurrency != "" {
			if len(currency) > 0 {
				for _, cur := range currency {
					if cur != imp.BidFloorCurrency {
						err = ErrInvalidCur
					}
				} // end for
			} // end if
		}
	}
	return
}
