//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package models

import "geniusrabbit.dev/sspserver/internal/billing"

// Target type object
type Target interface {
	// ID of object (Zone OR SmartLink only)
	ID() uint64

	// Codename of the target (equal to tagid)
	Codename() string

	// AlternativeAdCode returns URL or any code (HTML, XML, etc)
	AlternativeAdCode(key string) string

	// PurchasePrice gives the price of view from external resource
	PurchasePrice(action Action) billing.Money

	// RevenueShareFactor of current target
	RevenueShareFactor() float64

	// ComissionShareFactor of current target
	ComissionShareFactor() float64

	// Company object
	Company() *Company

	// CompanyID of current target
	CompanyID() uint64
}

// // TargetFromModel convert datavase model specified model
// // which implements Target interface
// func TargetFromModel(zone models.Zone) Target {
// 	if zone.Type == models.ZoneTypeSmartlink {
// 		return SmartlinkFromModel(zone)
// 	}
// 	return ZoneFromModel(zone)
// }
