//
// @project geniusrabbit::sspserver 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package models

import "sort"

// CountryBid targeting
type CountryBid struct {
	CountryCode CountryCode
	MinBid      float64
}

// GeoBidSlice object
type GeoBidSlice []CountryBid

// Sort the slice of geo bids
func (s GeoBidSlice) Sort() {
	sort.Slice(s, func(i, j int) bool { return s[i].CountryCode < s[j].CountryCode })
}

// Test the minimal bid in specific country
func (s GeoBidSlice) Test(cc CountryCode, bid float64) bool {
	return len(s) < 1 || bid >= s.Bid(cc)
}

// Bid by country
func (s GeoBidSlice) Bid(cc CountryCode) float64 {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].CountryCode >= cc
	})
	if i >= 0 && i < len(s) && s[i].CountryCode == cc {
		return s[i].MinBid
	}
	return 0
}
