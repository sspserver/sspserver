//
// @project geniusrabbit::sspserver 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package models

import "sort"

// Zone object descrbes the place where will be displayd the advertisement
type Zone struct {
	ID uint64

	// Zone hashed codename
	Codename string

	// Tags of the zone
	Tags []string

	// Minimal bid by country
	GeoMinimalBids GeoBidSlice

	// Global minimal bid
	MinBid float64

	// Source list of available resources
	SourceList []ZoneSourceRule

	// Default advertisement code if no bids or less then minimal bid
	Default string
}

// Prepare params
func (it *Zone) Prepare() {
	if len(it.Tags) > 0 {
		sort.Strings(it.Tags)
	}
}
