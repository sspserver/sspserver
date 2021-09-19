//
// @project GeniusRabbit rotator 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package models

import (
	"github.com/geniusrabbit/gosql"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// RTBAccessPoint for DSP connect.
// It means that this is entry point which contains
// information for access and search data
//
// As and Zone or Smartlink this object implements target interface
type RTBAccessPoint struct {
	id             uint64
	codename       string
	CompanyObj     *Company
	SecretKey      string
	Protocol       string
	DomainDefault  string
	Timeout        int
	Headers        gosql.Hstore
	Filter         types.BaseFilter
	RevenueShare   float64 // % 100.00, 80%, 65.5%
	ExcludeSources bool
	Sources        gosql.NullableOrderedUintArray
}

// RTBAccessPointFromModel convert database model to specified model
// func RTBAccessPointFromModel(cl *models.RTBAccessPoint, comp *Company) *RTBAccessPoint {
// 	if comp == nil {
// 		return nil
// 	}

// 	var (
// 		filter = types.BaseFilter{
// 			Secure:          cl.Secure,
// 			Adblock:         cl.AdBlock,
// 			PrivateBrowsing: cl.PrivateBrowsing,
// 		}
// 	)

// 	filter.Set(types.FieldFormat, cl.Formats)
// 	filter.Set(types.FieldDeviceTypes, cl.DeviceTypes)
// 	filter.Set(types.FieldDevices, cl.Devices)
// 	filter.Set(types.FieldOS, cl.OS)
// 	filter.Set(types.FieldBrowsers, cl.Browsers)
// 	filter.Set(types.FieldCategories, cl.Categories)
// 	filter.Set(types.FieldCountries, cl.Countries)
// 	filter.Set(types.FieldLanguages, cl.Languages)
// 	filter.Set(types.FieldZones, cl.Zones)
// 	filter.Set(types.FieldDomains, cl.Domains)

// 	sources, exclude := types.IDArrayFilter(cl.Sources)

// 	return &RTBAccessPoint{
// 		id:             cl.ID,
// 		codename:       strings.TrimSpace(cl.Codename),
// 		CompanyObj:     comp,
// 		Protocol:       strings.ToLower(cl.Protocol),
// 		DomainDefault:  cl.DomainDefault,
// 		Timeout:        cl.Timeout,
// 		Headers:        cl.Headers,
// 		Filter:         filter,
// 		RevenueShare:   cl.RevenueShare,
// 		ExcludeSources: exclude,
// 		Sources:        sources,
// 	}
// }

// TestSource for access point
func (s RTBAccessPoint) TestSource(sourceID uint64) bool {
	return s.Sources.Len() < 1 || (s.Sources.IndexOf(uint(sourceID)) > 0) != s.ExcludeSources
}

// ID of object (Zone OR SmartLink only)
func (s RTBAccessPoint) ID() uint64 {
	return s.id
}

// Codename of the target (equal to tagid)
func (s RTBAccessPoint) Codename() string {
	return s.codename
}

// PurchasePrice gives the price of view from external resource
func (s RTBAccessPoint) PurchasePrice() billing.Money {
	return 0
}

// Size default of target item
func (s RTBAccessPoint) Size() (w, h int) {
	return
}

// RevenueShareFactor of current target
func (s RTBAccessPoint) RevenueShareFactor() float64 {
	return s.RevenueShare / 100.0
}

// ComissionShareFactor of current target
func (s RTBAccessPoint) ComissionShareFactor() float64 {
	return 1.0 - s.RevenueShareFactor()
}

// Company object
func (s RTBAccessPoint) Company() *Company {
	return s.CompanyObj
}

// CompanyID of current target
func (s RTBAccessPoint) CompanyID() uint64 {
	if s.CompanyObj == nil {
		return 0
	}
	return s.CompanyObj.ID
}

func (s *RTBAccessPoint) reset() {
	if s != nil {
		s.Sources = s.Sources[:0]
		s.CompanyObj = nil
		s.Filter.Reset()
	}
}
