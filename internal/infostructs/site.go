//
// @project Geniusrabbit::corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package infostructs

import (
	"strings"

	"github.com/bsm/openrtb"
)

// Site information
type Site struct {
	ExtID         string   `json:"eid,omitempty"`          // External ID
	Domain        string   `json:"domain,omitempty"`       //
	Cat           []string `json:"cat,omitempty"`          // Array of categories
	PrivacyPolicy int      `json:"pivacypolicy,omitempty"` // Default: 1 ("1": has a privacy policy)
	Keywords      string   `json:"keywords,omitempty"`     // Comma separated list of keywords about the site.
	Page          string   `json:"page,omitempty"`         // URL of the page
	Ref           string   `json:"ref,omitempty"`          // Referrer URL
	Search        string   `json:"search,omitempty"`       // Search string that caused naviation
	Mobile        int      `json:"mobile,omitempty"`       // Mobile ("1": site is mobile optimised)
}

// SiteDefault info
var SiteDefault Site

// DomainPrepared value
func (s Site) DomainPrepared() []string {
	return PrepareDomain(s.Domain)
}

// RTBObject of Site
func (s *Site) RTBObject() *openrtb.Site {
	if s == nil {
		return nil
	}
	return &openrtb.Site{
		Inventory: openrtb.Inventory{
			ID:            s.ExtID,               // External ID
			Keywords:      s.Keywords,            // Comma separated list of keywords about the site.
			Cat:           s.Cat,                 // Array of IAB content categories
			Domain:        s.Domain,              //
			PrivacyPolicy: intO(s.PrivacyPolicy), // Default: 1 ("1": has a privacy policy)
		},
		Page:   s.Page,   // URL of the page
		Ref:    s.Ref,    // Referrer URL
		Search: s.Search, // Search string that caused naviation
		Mobile: s.Mobile, // Mobile ("1": site is mobile optimised)
	}
}

// PrepareDomain parts
func PrepareDomain(domain string) (list []string) {
	domain = strings.ToLower(domain)
	if domain == "" {
		return []string{"*."}
	}

	list = append(list, domain)
	if strings.HasPrefix(domain, "www.") {
		list = append(list, "*."+domain)
		domain = domain[4:]
	}

	list = append(list, "*."+domain)
	arr := strings.Split(domain, ".")
	for i := 1; i < len(arr); i++ {
		list = append(list, "*."+strings.Join(arr[i:], "."))
	}

	return
}
