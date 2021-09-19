//
// @project Geniusrabbit::corelib 2016 – 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2019
//

package infostructs

import (
	"net"

	"github.com/bsm/openrtb"
)

// Geo base information structure
type Geo struct {
	ID            uint     `json:"id,omitempty"`            // Internal geo ID
	IP            net.IP   `json:"ip,omitempty"`            // IPv4/6
	Carrier       *Carrier `json:"carrier,omitempty"`       // Carrier or ISP derived from the IP address
	Lat           float64  `json:"lat,omitempty"`           // Latitude from -90 to 90
	Lon           float64  `json:"lon,omitempty"`           // Longitude from -180 to 180
	Country       string   `json:"country,omitempty"`       // Country using ISO 3166-1 Alpha 2
	Region        string   `json:"region,omitempty"`        // Region using ISO 3166-2
	RegionFIPS104 string   `json:"regionFIPS104,omitempty"` // Region of a country using FIPS 10-4
	Metro         string   `json:"metro,omitempty"`         //
	City          string   `json:"city,omitempty"`          //
	Zip           string   `json:"zip,omitempty"`           //
	UTCOffset     int      `json:"utcoffset,omitempty"`     // Local time as the number +/- of minutes from UTC
}

// GeoDefault value
var GeoDefault = Geo{Country: "**", Carrier: &CarrierDefault}

// IsIPv6 format
func (g Geo) IsIPv6() bool {
	return g.IP != nil && g.IP.To4() == nil
}

// IPv4String string value
func (g Geo) IPv4String() string {
	if g.IP == nil {
		return ""
	}
	if g.IsIPv6() {
		return ""
	}
	return g.IP.String()
}

// IPv6String string value
func (g Geo) IPv6String() string {
	if g.IP == nil {
		return ""
	}
	if !g.IsIPv6() {
		return ""
	}
	return g.IP.String()
}

// RTBObject of Geo
func (g Geo) RTBObject() *openrtb.Geo {
	return &openrtb.Geo{
		Lat:           g.Lat,           // Latitude from -90 to 90
		Lon:           g.Lon,           // Longitude from -180 to 180
		Type:          0,               // Indicate the source of the geo data
		Accuracy:      0,               // Estimated location accuracy in meters; recommended when lat/lon are specified and derived from a device’s location services
		LastFix:       0,               // Number of seconds since this geolocation fix was established.
		IPService:     0,               // Service or provider used to determine geolocation from IP address if applicable
		Country:       g.Country,       // Country using ISO 3166-1 Alpha 3
		Region:        g.Region,        // Region using ISO 3166-2
		RegionFIPS104: g.RegionFIPS104, // Region of a country using FIPS 10-4
		Metro:         g.Metro,         //
		City:          g.City,          //
		Zip:           g.Zip,           //
		UTCOffset:     g.UTCOffset,     // Local time as the number +/- of minutes from UTC
		Ext:           nil,             //
	}
}
