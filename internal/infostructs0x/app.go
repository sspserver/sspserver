//
// @project Geniusrabbit::corelib 2016 – 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2019
//

package infostructs

import "github.com/bsm/openrtb"

// App information
type App struct {
	ExtID         string   `json:"eid,omitempty"`          // External ID
	Keywords      string   `json:"keywords,omitempty"`     // Comma separated list of keywords about the site.
	Cat           []string `json:"cat,omitempty"`          // Array of categories
	Bundle        string   `json:"bundle,omitempty"`       // App bundle or package name
	StoreURL      string   `json:"storeurl,omitempty"`     // App store URL for an installed app
	Ver           string   `json:"ver,omitempty"`          // App version
	Paid          int      `json:"paid,omitempty"`         // "1": Paid, "2": Free
	PrivacyPolicy int      `json:"pivacypolicy,omitempty"` // Default: 1 ("1": has a privacy policy)
}

// AppDefault object
var AppDefault App

// RTBObject of App
func (a *App) RTBObject() *openrtb.App {
	if a == nil {
		return nil
	}
	return &openrtb.App{
		Inventory: openrtb.Inventory{
			ID:            a.ExtID,               // External ID
			Keywords:      a.Keywords,            // Comma separated list of keywords about the site.
			Cat:           a.Cat,                 // Array of IAB content categories
			PrivacyPolicy: intO(a.PrivacyPolicy), // Default: 1 ("1": has a privacy policy)
		},
		Bundle:   a.Bundle,   // App bundle or package name
		StoreURL: a.StoreURL, // App store URL for an installed app
		Ver:      a.Ver,      // App version
		Paid:     a.Paid,     // "1": Paid, "2": Free
	}
}

// DomainPrepared value
func (a *App) DomainPrepared() []string {
	return PrepareDomain(a.Bundle)
}
