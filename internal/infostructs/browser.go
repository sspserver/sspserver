//
// @project Geniusrabbit::corelib 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package infostructs

// Browser base information structure
type Browser struct {
	ID              uint        `json:"id,omitempty"`    // Internal system ID
	Name            string      `json:"name,omitempty"`  //
	Version         string      `json:"ver,omitempty"`   //
	DNT             int         `json:"dnt,omitempty"`   // "1": Do not track
	LMT             int         `json:"lmt,omitempty"`   // "1": Limit Ad Tracking
	Adblock         int         `json:"ab,omitempty"`    // "1": AdBlock is ON
	PrivateBrowsing int         `json:"pb,omitempty"`    // "1": Private Browsing mode ON
	UA              string      `json:"ua,omitempty"`    // User agent
	Ref             string      `json:"r,omitempty"`     // Referer
	JS              int         `json:"js,omitempty"`    //
	Languages       []string    `json:"langs,omitempty"` //
	PrimaryLanguage string      `json:"lang,omitempty"`  // Browser language (en-US)
	FlashVer        string      `json:"flver,omitempty"` // Flash version
	Width           int         `json:"w,omitempty"`     // Window in pixels
	Height          int         `json:"h,omitempty"`     // Window in pixels
	IsRobot         int         `json:"rb,omitempty"`
	Extensions      []Extension `json:"extensions,omitempty"`
}

// Extension of some Browser/OS
type Extension struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"ver,omitempty"`
}

// BrowserDefault value
var BrowserDefault Browser
