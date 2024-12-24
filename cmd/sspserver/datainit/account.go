package datainit

import (
	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/billing"
)

type Account struct {
	ID uint64 `json:"id" yaml:"id"`

	Balance  billing.Money `json:"balance" yaml:"balance"`
	MaxDaily billing.Money `json:"max_daily" yaml:"max_daily"`
	Spent    billing.Money `json:"spent" yaml:"spent"`

	// RevenueShare it's amount of percent of the raw incode which will be shared with the publisher company
	// For example:
	//   Displayed ads for 100$
	//   Account revenue share 60%
	//   In such case the ad network have 40$
	//   The publisher have 60$
	RevenueShare float64 `json:"revenue_share" yaml:"revenue_share"` // % 100_00, 1 -> 100%, 0.655 -> 65.5%
}

func AdModelAccount(a *Account) (*admodels.Account, bool) {
	return &admodels.Account{
		IDval:        a.ID,
		MaxDaily:     a.MaxDaily,
		RevenueShare: min(a.RevenueShare, 1.0),
	}, true
}

func (a *Account) GetID() uint64 {
	return a.ID
}

func (a *Account) IsApproved() bool {
	return true
}

func (a *Account) TableName() string {
	return `account_base`
}
