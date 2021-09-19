//
// @project GeniusRabbit rotator 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package models

import (
	"sort"

	"geniusrabbit.dev/sspserver/internal/billing"
)

// GeoBid item
type GeoBid struct {
	Code [2]byte
	Bid  billing.Money
}

// GeoLess campare
func (b *GeoBid) GeoLess(bid *GeoBid) bool {
	return b.GeoLessCode(bid.Code)
}

// GeoLessCode campare
func (b *GeoBid) GeoLessCode(code [2]byte) bool {
	if b.Code[0] < code[0] {
		return true
	}
	return b.Code[0] == code[0] && b.Code[1] < code[1]
}

// GeoBidSlice type
type GeoBidSlice []GeoBid

func (a GeoBidSlice) Len() int           { return len(a) }
func (a GeoBidSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a GeoBidSlice) Less(i, j int) bool { return a[i].GeoLess(&a[j]) }

// IndexOf by country code (only for ordered)
func (a GeoBidSlice) IndexOf(code [2]byte) int {
	if len(a) < 1 {
		return -1
	}
	idx := sort.Search(a.Len(), func(i int) bool {
		return !a[i].GeoLessCode(code)
	})
	if idx >= len(a) {
		return -1
	}
	return idx
}
