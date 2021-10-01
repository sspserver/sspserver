//
// @project GeniusRabbit rotator 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//

package events

import "encoding/json"

// UserInfo collects the general information about users
type UserInfo struct {
	Time      int64  `json:"tm,omitempty"`  // Timestamp
	AuctionID string `json:"auc,omitempty"` // Internal Auction ID

	// User IDENTITY
	UDID      string `json:"udi,omitempty"` // Unique Device ID (IDFA)
	UUID      string `json:"uui,omitempty"` // User
	SessionID string `json:"ses,omitempty"` // -- // --

	// Personal information
	Age           byte   `age:"age,omitempty"`
	Gender        byte   `json:"gender,omitempty"`
	SearchGender  byte   `json:"search_gender,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	MessangerType string `json:"messanger_type,omitempty"`
	Messanger     string `json:"messanger,omitempty"`
	Postcode      string `json:"zip,omitempty"`
	Facebook      string `json:"facebook,omitempty"`
	Twitter       string `json:"twitter,omitempty"`
	Linkedin      string `json:"linkedin,omitempty"`

	// Location info
	Country   string `json:"cc,omitempty"`  // Country Code ISO-2
	City      string `json:"ct,omitempty"`  // City Code
	Latitude  string `json:"lt,omitempty"`  // -- // --
	Longitude string `json:"lg,omitempty"`  // -- // --
	Language  string `json:"lng,omitempty"` // en-US
}

// String imprlementation of Stringer interface
func (info *UserInfo) String() string {
	if json, err := json.MarshalIndent(info, "", "  "); err == nil {
		return string(json)
	}
	return ""
}
