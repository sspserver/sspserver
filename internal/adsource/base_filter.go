//
// @project GeniusRabbit rotator 2016, 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2018
//

// @TODO add format types

package adsource

import (
	"sort"

	"github.com/cznic/sortutil"
	"github.com/sspserver/udetect"
)

// BaseFilter object
type BaseFilter struct {
	Secure               int   // 0 - any, 1 - only secure, 2 - no secure
	Adblock              int   // 0 - any, 1 - only adblock, 2 - no adblock
	PrivateBrowsing      int   // 0 - any, 1 - only private, 2 - no private
	Devices              []int // Devices type
	OS                   []uint
	OSExclude            []uint
	Browsers             []uint
	BrowsersExclude      []uint
	Categories           []uint
	Countries            []string
	Applications         []uint
	ApplicationsExclude  []uint
	Domains              []string
	DomainsExclude       []string
	Zones                []uint
	ZonesExclude         []uint
	ExternalZones        []string
	ExternalZonesExclude []string
}

// Normalise params
func (f *BaseFilter) Normalise() {
	sort.Ints(f.Devices)
	sort.Sort(sortutil.UintSlice(f.Categories))
	sort.Strings(f.Countries)
	sort.Sort(sortutil.UintSlice(f.Applications))
	sort.Sort(sortutil.UintSlice(f.ApplicationsExclude))
	sort.Strings(f.Domains)
	sort.Strings(f.DomainsExclude)
	sort.Sort(sortutil.UintSlice(f.Zones))
	sort.Sort(sortutil.UintSlice(f.ZonesExclude))
	sort.Strings(f.ExternalZones)
	sort.Strings(f.ExternalZonesExclude)
}

// Test base from search request
func (f *BaseFilter) Test(request *BidRequest) bool {
	switch {
	case (request.IsSecure() && f.Secure == 2) || (!request.IsSecure() && f.Secure == 1):
		return false
	case (request.IsAdblock() && f.Adblock == 2) || (!request.IsAdblock() && f.Adblock == 1):
		return false
	case (request.IsPrivateBrowsing() && f.PrivateBrowsing == 2) || (!request.IsPrivateBrowsing() && f.PrivateBrowsing == 1):
		return false
	}

	var (
		deviceType  udetect.DeviceType
		countryCode string
	)

	if request.Device != nil {
		deviceType = request.Device.DeviceType
	}
	if request.User != nil {
		if request.User.Geo != nil {
			countryCode = request.User.Geo.Country
		}
	}

	return true &&
		(len(f.Devices) /* ***** */ < 1 || hasInIntArr(int(deviceType), f.Devices)) &&
		(len(f.Countries) /* *** */ < 1 || hasInStringArr(countryCode, f.Countries)) &&
		(len(f.Categories) /* ** */ < 1 || intersecUintArr(request.Categories(), f.Categories)) &&
		(len(f.ApplicationsExclude) < 1 || !hasInUintArr(uint(request.AppID()), f.ApplicationsExclude)) &&
		(len(f.Applications) /*  */ < 1 || hasInUintArr(uint(request.AppID()), f.Applications)) &&
		(len(f.DomainsExclude) /**/ < 1 || !hasOneInStringArr(request.Domain(), f.DomainsExclude)) &&
		(len(f.Domains) /* ***** */ < 1 || hasOneInStringArr(request.Domain(), f.Domains)) &&
		(len(f.ZonesExclude) /*  */ < 1 || !intersecUintArr(request.TargetIDs(), f.ZonesExclude)) &&
		(len(f.Zones) /* ******* */ < 1 || intersecUintArr(request.TargetIDs(), f.Zones)) &&
		(len(f.ExternalZonesExclude) < 1 || !hasOneInStringArr(request.ExtTargetIDs(), f.ExternalZonesExclude)) &&
		(len(f.ExternalZones) /* */ < 1 || hasOneInStringArr(request.ExtTargetIDs(), f.ExternalZones))
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func hasInIntArr(v int, arr []int) bool {
	i := sort.SearchInts(arr, v)
	return i >= 0 && i < len(arr) && v == arr[i]
}

func hasInInt32Arr(v int32, arr []int32) bool {
	i := sortutil.SearchInt32s(arr, v)
	return i >= 0 && i < len(arr) && v == arr[i]
}

func hasInUintArr(v uint, arr []uint) bool {
	i := sortutil.SearchUints(arr, v)
	return i >= 0 && i < len(arr) && v == arr[i]
}

func hasInStringArr(v string, arr []string) bool {
	i := sort.SearchStrings(arr, v)
	return i >= 0 && i < len(arr) && v == arr[i]
}

func hasOneInStringArr(arr1, arr2 []string) bool {
	for _, v := range arr1 {
		if hasInStringArr(v, arr2) {
			return true
		}
	}
	return false
}

func intersecInt32Arr(cat1, cat2 []int32) bool {
	if len(cat1) < 1 && len(cat2) < 1 {
		return true
	}
	if len(cat1) < 1 || len(cat2) < 1 {
		return false
	}

	for _, c1 := range cat1 {
		if hasInInt32Arr(c1, cat2) {
			return true
		}
	}
	return false
}

func intersecIntArr(cat1, cat2 []int) bool {
	if len(cat1) < 1 && len(cat2) < 1 {
		return true
	}
	if len(cat1) < 1 || len(cat2) < 1 {
		return false
	}

	for _, c1 := range cat1 {
		if hasInIntArr(c1, cat2) {
			return true
		}
	}
	return false
}

func intersecUintArr(cat1, cat2 []uint) bool {
	if len(cat1) < 1 && len(cat2) < 1 {
		return true
	}
	if len(cat1) < 1 || len(cat2) < 1 {
		return false
	}

	for _, c1 := range cat1 {
		if hasInUintArr(c1, cat2) {
			return true
		}
	}
	return false
}
