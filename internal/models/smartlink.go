//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package models

import (
	"github.com/geniusrabbit/gosql"

	"geniusrabbit.dev/sspserver/internal/billing"
)

// Smartlink model
type Smartlink struct {
	id                uint64
	StringID          string
	Price             billing.Money // The cost of single view
	Comp              *Company
	CompID            uint64
	AllowedTypes      gosql.NullableOrderedIntArray
	AllowedSources    gosql.NullableOrderedIntArray
	DisallowedSources gosql.NullableOrderedIntArray
	Campaigns         gosql.NullableOrderedIntArray
	DefaultCode       map[string]string
}

// // SmartlinkFromModel convert database model to specified model
// func SmartlinkFromModel(zone models.Zone) *Smartlink {
// 	var code map[string]string
// 	zone.DefaultCode.UnmarshalTo(&code)

// 	return &Smartlink{
// 		id:                zone.ID,
// 		StringID:          strconv.FormatUint(zone.ID, 10),
// 		Price:             zone.Price,
// 		Comp:              nil,
// 		CompID:            zone.CompanyID,
// 		AllowedTypes:      zone.AllowedTypes,
// 		AllowedSources:    zone.AllowedSources,
// 		DisallowedSources: zone.DisallowedSources,
// 		Campaigns:         zone.Campaigns,
// 		DefaultCode:       code,
// 	}
// }

// ID of object
func (l *Smartlink) ID() uint64 {
	return l.id
}

// Codename of the target (equal to tagid)
func (l *Smartlink) Codename() string {
	return l.StringID
}

// AlternativeAdCode returns URL or any code (HTML, XML, etc)
func (l *Smartlink) AlternativeAdCode(key string) (val string) {
	if l.DefaultCode == nil {
		return ""
	}
	val, _ = l.DefaultCode[key]
	return val
}

// PurchasePrice gives the price of view from external resource
func (l *Smartlink) PurchasePrice(action Action) billing.Money {
	if action.IsImpression() {
		return l.Price
	}
	return 0
}

// Company object
func (l *Smartlink) Company() *Company {
	return l.Comp
}

// CompanyID of current target
func (l *Smartlink) CompanyID() uint64 {
	return l.CompID
}

// RevenueShareFactor amount %
func (l *Smartlink) RevenueShareFactor() float64 {
	return l.Comp.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (l *Smartlink) ComissionShareFactor() float64 {
	return l.Comp.ComissionShareFactor()
}

// IsAllowedSource for targeting
func (l *Smartlink) IsAllowedSource(id uint64, types []int) bool {
	if len(l.AllowedSources) > 0 {
		return l.AllowedSources.IndexOf(int(id)) >= 0
	}
	if len(l.DisallowedSources) > 0 {
		return l.DisallowedSources.IndexOf(int(id)) < 0
	}
	return true
}
