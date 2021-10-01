//
// @project GeniusRabbit AdNet
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package models

import (
	"time"

	"bitbucket.org/geniusrabbit/billing"
	"github.com/guregu/null"
)

// ```pg
// CREATE TABLE company_m2m_member
// ( company_id          BIGINT                      NOT NULL
// , user_id             BIGINT                      NOT NULL
//
// , is_admin            BOOL                        NOT NULL      DEFAULT FALSE -- for current company
// , acl                 JSONB                                     DEFAULT NULL  -- {model:flags,@custom:value}
// , roles               BIGINT[]                                  DEFAULT NULL
//
// , created_at          TIMESTAMPTZ                 NOT NULL      DEFAULT NOW()
// , updated_at          TIMESTAMPTZ                 NOT NULL      DEFAULT NOW()
// , deleted_at          TIMESTAMPTZ
//
// , PRIMARY KEY (user_id, company_id)
// );
// ```

// Company model
type Company struct {
	ID          uint64        `json:"id"`                                           //
	Name        string        `json:"name"`                                         // Unique project name. Like login
	Title       string        `json:"title"`                                        //
	Description string        `json:"description"`                                  //
	Status      ApproveStatus `json:"status"`                                       //
	Members     []*User       `gorm:"many2many:company_m2m_member;" json:"members"` // Members of project
	CompanyName string        `json:"company_name"`                                 // Company info
	Country     string        `json:"country"`                                      // - // -
	City        string        `json:"city"`                                         // - // -
	Address     string        `json:"address"`                                      // - // -
	Phone       string        `json:"phone"`                                        // Contacts
	Email       string        `json:"email"`                                        // - // -
	Messanger   string        `json:"messanger"`                                    // - // -

	MaxDaily     billing.Money `json:"max_daily,omitempty"`
	RevenueShare float64       `json:"revenue_share,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (c *Company) TableName() string {
	return "company"
}
