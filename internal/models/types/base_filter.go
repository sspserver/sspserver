//
// @project GeniusRabbit rotator 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package types

import "github.com/geniusrabbit/gosql"

// Base filter fields
const (
	FieldFormat = iota
	FieldDeviceTypes
	FieldDevices
	FieldOS
	FieldBrowsers
	FieldCategories
	FieldCountries
	FieldLanguages
	FieldZones
	FieldDomains
)

// BaseFilter object
type BaseFilter struct {
	excludeMask     uint64
	Formats         gosql.StringArray
	DeviceTypes     gosql.NullableOrderedUintArray //
	Devices         gosql.NullableOrderedUintArray //
	OS              gosql.NullableOrderedUintArray //
	Browsers        gosql.NullableOrderedUintArray //
	Categories      gosql.NullableOrderedUintArray //
	Countries       gosql.NullableOrderedUintArray //
	Languages       gosql.NullableOrderedUintArray //
	Zones           gosql.NullableOrderedUintArray //
	Domains         gosql.StringArray              //
	Secure          int                            // 0 - any, 1 - only, 2 - exclude
	Adblock         int                            // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int                            // 0 - any, 1 - only, 2 - exclude
	IP              int                            // 0 - any, 1 - IPv4, 2 - IPv6
}

// Set filter item
func (fl *BaseFilter) Set(field uint64, data interface{}) {
	var positive bool
	switch field {
	case FieldFormat:
		fl.Formats, _ = data.(gosql.StringArray)
	case FieldDeviceTypes:
		fl.DeviceTypes, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
	case FieldDevices:
		fl.Devices, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
	case FieldOS:
		fl.OS, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
	case FieldBrowsers:
		fl.Browsers, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
	case FieldCategories:
		fl.Categories, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
	case FieldCountries:
		switch data.(type) {
		case gosql.NullableOrderedIntArray:
			fl.Countries, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
		case gosql.StringArray:
			fl.Countries, positive = CountryFilter(data.(gosql.StringArray))
		}
	case FieldLanguages:
		switch data.(type) {
		case gosql.NullableOrderedIntArray:
			fl.Languages, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
		case gosql.StringArray:
			fl.Languages, positive = LanguageFilter(data.(gosql.StringArray))
		}
	case FieldZones:
		fl.Zones, positive = IDArrayFilter(data.(gosql.NullableOrderedIntArray))
	case FieldDomains:
		fl.Domains, positive = StringArrayFilter(data.(gosql.StringArray))
	}
	fl.SetPositive(field, positive)
}

// SetPositive field state
func (fl *BaseFilter) SetPositive(field uint64, positive bool) {
	if positive {
		fl.excludeMask |= 1 << field
	} else {
		fl.excludeMask &= ^(1 << field)
	}
}

// Test filter items
func (fl *BaseFilter) Test(t TargetPointer) bool {
	found := len(t.Formats()) < 1
	for _, f := range t.Formats() {
		if found = fl.TestFormat(f); found {
			break
		}
	}
	return found &&
		(fl.Secure == 0 || (fl.Secure == 1) == t.IsSecure()) &&
		(fl.Adblock == 0 || (fl.Adblock == 1) == t.IsAdblock()) &&
		(fl.PrivateBrowsing == 0 || (fl.PrivateBrowsing == 1) == t.IsPrivateBrowsing()) &&
		(fl.IP == 0 || (fl.IP == 2) == t.IsIPv6()) &&
		fl.checkUintArr(t.DeviceType(), FieldDeviceTypes, fl.DeviceTypes) &&
		fl.checkUintArr(t.DeviceID(), FieldDevices, fl.Devices) &&
		fl.checkUintArr(t.OSID(), FieldOS, fl.OS) &&
		fl.checkUintArr(t.BrowserID(), FieldBrowsers, fl.Browsers) &&
		fl.multyCheckUintArr(t.Categories(), FieldCategories, fl.Categories) &&
		fl.checkUintArr(t.GeoID(), FieldCountries, fl.Countries) &&
		fl.checkUintArr(t.LanguageID(), FieldLanguages, fl.Languages) &&
		fl.checkUintArr(uint(t.TargetID()), FieldZones, fl.Zones)
}

// TestFormat available in filter
func (fl *BaseFilter) TestFormat(f *Format) bool {
	found := len(fl.Formats) < 1
	for _, code := range fl.Formats {
		if f.Codename == code {
			return true
		}
	}
	return found
}

func (fl *BaseFilter) checkUintArr(v uint, off uint64, arr gosql.NullableOrderedUintArray) bool {
	return arr.Len() < 1 || (arr.IndexOf(v) >= 0) == (fl.excludeMask&(1<<off) == 0)
}

func (fl *BaseFilter) multyCheckUintArr(v []uint, off uint64, arr gosql.NullableOrderedUintArray) bool {
	return arr.Len() < 1 || arr.OneOf(v) == (fl.excludeMask&(1<<off) == 0)
}

// Reset filter object
func (fl *BaseFilter) Reset() {
	fl.excludeMask = 0
	fl.Formats = fl.Formats[:0]
	fl.DeviceTypes = fl.DeviceTypes[:0]
	fl.Devices = fl.Devices[:0]
	fl.OS = fl.OS[:0]
	fl.Browsers = fl.Browsers[:0]
	fl.Categories = fl.Categories[:0]
	fl.Countries = fl.Countries[:0]
	fl.Languages = fl.Languages[:0]
	fl.Zones = fl.Zones[:0]
	fl.Domains = fl.Domains[:0]
	fl.Secure = 0
	fl.Adblock = 0
	fl.PrivateBrowsing = 0
	fl.IP = 0
}
