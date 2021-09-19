//
// @project Geniusrabbit::corelib 2016 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2018
//

package infostructs

import "github.com/bsm/openrtb"

// DeviceType object declaration
type DeviceType int

// RTB 5.17 Device Type
const (
	DeviceTypeUnknown   DeviceType = 0
	DeviceTypeMobile    DeviceType = 1
	DeviceTypePC        DeviceType = 2
	DeviceTypeTV        DeviceType = 3
	DeviceTypePhone     DeviceType = 4
	DeviceTypeTablet    DeviceType = 5
	DeviceTypeConnected DeviceType = 6
	DeviceTypeSetTopBox DeviceType = 7
	DeviceTypeWatch     DeviceType = 8
	DeviceTypeGlasses   DeviceType = 9
)

// Device base information structure
type Device struct {
	ID         uint       `json:"id,omitempty"`             // Device ID
	Make       string     `json:"make,omitempty"`           // Device make
	Model      string     `json:"model,omitempty"`          // Device model
	OS         *OS        `json:"os,omitempty"`             // Device OS
	Browser    *Browser   `json:"browser,omitempty"`        // Device OS version
	ConnType   int        `json:"connectiontype,omitempty"` //
	DeviceType DeviceType `json:"devicetype,omitempty"`     //
	IFA        string     `json:"ifa,omitempty"`            // Native identifier for advertisers
	Height     int        `json:"h,omitempty"`              // Physical height of the screen in pixels.
	Width      int        `json:"w,omitempty"`              // Physical width of the screen in pixels.
	PPI        int        `json:"ppi,omitempty"`            // Screen size as pixels per linear inch.
	PxRatio    float64    `json:"pxratio,omitempty"`        // The ratio of physical pixels to device independent pixels.
	HwVer      string     `json:"hwv,omitempty"`            // Hardware version of the device (e.g., "5S" for iPhone 5S).
}

// DeviceDefault value
var DeviceDefault = Device{Browser: &BrowserDefault, OS: &OSDefault}

// RTBObject of Device
func (d Device) RTBObject(geo *Geo) *openrtb.Device {
	var (
		browser = d.Browser
		os      = d.OS
		carrier *Carrier
		ipV4    = geo.IPv4String()
	)

	{
		if browser == nil {
			browser = &BrowserDefault
		}

		if os == nil {
			os = &OSDefault
		}

		if geo == nil {
			geo = &GeoDefault
		}

		if carrier = geo.Carrier; carrier == nil {
			carrier = &CarrierDefault
		}

		// IP by default
		if ipV4 == "" && geo.IPv6String() == "" {
			ipV4 = "0.0.0.0"
		}
	}

	return &openrtb.Device{
		UA:         browser.UA,              // User agent
		Geo:        geo.RTBObject(),         // Location of the device assumed to be the user’s current location
		DNT:        browser.DNT,             // "1": Do not track
		LMT:        browser.LMT,             // "1": Limit Ad Tracking
		IP:         ipV4,                    // IPv4
		IPv6:       geo.IPv6String(),        // IPv6
		DeviceType: int(d.DeviceType),       // The general type of d.
		Make:       d.Make,                  // Device make
		Model:      d.Model,                 // Device model
		OS:         os.Name,                 // Device OS
		OSVer:      os.Version,              // Device OS version
		HwVer:      d.HwVer,                 // Hardware version of the device (e.g., "5S" for iPhone 5S).
		H:          d.Height,                // Physical height of the screen in pixels.
		W:          d.Width,                 // Physical width of the screen in pixels.
		PPI:        d.PPI,                   // Screen size as pixels per linear inch.
		PxRatio:    d.PxRatio,               // The ratio of physical pixels to device independent pixels.
		JS:         browser.JS,              // Javascript status ("0": Disabled, "1": Enabled)
		GeoFetch:   0,                       // Indicates if the geolocation API will be available to JavaScript code running in the banner,
		FlashVer:   browser.FlashVer,        // Flash version
		Language:   browser.PrimaryLanguage, // Browser language
		Carrier:    carrier.Name,            // Carrier or ISP derived from the IP address
		MCCMNC:     "",                      // Mobile carrier as the concatenated MCC-MNC code (e.g., "310-005" identifies Verizon Wireless CDMA in the USA).
		ConnType:   d.ConnType,              // Network connection type.
		IFA:        d.IFA,                   // Native identifier for advertisers
		IDSHA1:     "",                      // SHA1 hashed device ID
		IDMD5:      "",                      // MD5 hashed device ID
		PIDSHA1:    "",                      // SHA1 hashed platform device ID
		PIDMD5:     "",                      // MD5 hashed platform device ID
		MacSHA1:    "",                      // SHA1 hashed device ID; IMEI when available, else MEID or ESN
		MacMD5:     "",                      // MD5 hashed device ID; IMEI when available, else MEID or ESN
	}
}
