//
// @project GeniusRabbit::corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package models

import (
	"fmt"
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"

	"geniusrabbit.dev/sspserver/internal/billing"
)

// Zone Types enum
const (
	ZoneTypeDefault   = iota // 0
	ZoneTypeSmartlink        // 1
)

// Zone model
type Zone struct {
	ID                uint64                        `json:"id"`
	Title             string                        `json:"title"`
	Type              int                           `json:"type,omitempty"`
	Company           *Company                      `json:"company,omitempty"`
	CompanyID         uint64                        `json:"company_id,omitempty"`
	Status            ApproveStatus                 `json:"status"`
	Active            uint                          `json:"active"`
	DefaultCode       gosql.NullableJSON            `json:"default_code,omitempty"`
	Context           gosql.NullableJSON            `json:"context,omitempty"`            //
	MinECPM           float64                       `json:"min_ecpm,omitempty"`           // Default
	MinECPMByGeo      gosql.NullableJSON            `json:"min_ecpm_by_geo,omitempty"`    // {"CODE": <ecpm>, ...}
	Price             billing.Money                 `json:"price,omitempty"`              // The cost of single view
	AllowedFormats    gosql.NullableOrderedIntArray `json:"allowed_formats,omitempty"`    //
	AllowedTypes      gosql.NullableOrderedIntArray `json:"allowed_types,omitempty"`      //
	AllowedSources    gosql.NullableOrderedIntArray `json:"allowed_sources,omitempty"`    //
	DisallowedSources gosql.NullableOrderedIntArray `json:"disallowed_sources,omitempty"` //
	Campaigns         gosql.NullableOrderedIntArray `json:"campaigns,omitempty"`          // Strict campaigns targeting (smartlinks only)

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (z *Zone) TableName() string {
	return "adv_zone"
}

// RevenueShare amount %
func (z *Zone) RevenueShare() float64 {
	return 0
}

// SetCompany object
func (z *Zone) SetCompany(c interface{}) error {
	switch v := c.(type) {
	case *Company:
		z.Company = v
		z.CompanyID = v.ID
	case uint64:
		z.Company = nil
		z.CompanyID = v
	default:
		return fmt.Errorf("undefined value type: %t", c)
	}
	return nil
}
