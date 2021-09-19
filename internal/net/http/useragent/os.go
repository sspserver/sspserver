//
// @project geniusrabbit 2015 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 - 2018
//

package useragent

import (
	"strings"
)

// OSType value
const (
	OSTypeDesktop = iota
	OSTypeMobile
	OSTypeTV
	OSTypeConsole
)

// OS object
type OS struct {
	ID      uint     `json:"id"`
	Type    uint     `json:"type"`
	Name    string   `json:"name"`
	Match   []string `json:"-"`
	MatchUA []string `json:"-"`
}

// UndefinedOS value
var UndefinedOS = OS{ID: 0, Type: OSTypeDesktop, Name: "?"}

// OSes list of OS
var OSes = []OS{
	{ID: 0, Type: OSTypeDesktop, Name: "Undefined"},

	// TV
	{ID: 60, Type: OSTypeTV, Name: "GoogleTV", MatchUA: []string{"GoogleTV"}},
	{ID: 61, Type: OSTypeTV, Name: "WebTV"},
	{ID: 62, Type: OSTypeTV, Name: "Google Wireless Transcoder", MatchUA: []string{"Google Wireless Transcoder"}},

	// Desctop
	{ID: 1, Type: OSTypeDesktop, Name: "Chrome OS"},
	{ID: 2, Type: OSTypeDesktop, Name: "Linux & Unix", Match: []string{"Debian", "Gentoo", "Solaris", "Linux", "Unix"}},
	{ID: 3, Type: OSTypeDesktop, Name: "Mac"},
	{ID: 4, Type: OSTypeDesktop, Name: "Windows 10", MatchUA: []string{"Windows NT 10"}},
	{ID: 5, Type: OSTypeDesktop, Name: "Windows 8/8.1", Match: []string{"Windows RT", "Windows 8"}},
	{ID: 6, Type: OSTypeDesktop, Name: "Windows 7"},
	{ID: 7, Type: OSTypeDesktop, Name: "Windows Vista"},
	{ID: 8, Type: OSTypeDesktop, Name: "Windows XP"},
	{ID: 9, Type: OSTypeDesktop, Name: "Windows NT/2000", Match: []string{"Windows NT", "Windows 2000"}},
	{ID: 10, Type: OSTypeDesktop, Name: "Windows 98/ME", Match: []string{"Windows 98", "Windows ME"}},
	{ID: 11, Type: OSTypeDesktop, Name: "Windows 3.1/95", Match: []string{"Windows 3.1", "Windows 95"}},

	// Console
	{ID: 21, Type: OSTypeConsole, Name: "Nintendo Wii U"},
	{ID: 22, Type: OSTypeConsole, Name: "Nintendo Wii"},
	{ID: 23, Type: OSTypeConsole, Name: "Nintendo 3DS"},
	{ID: 24, Type: OSTypeConsole, Name: "Nintendo DS"},
	{ID: 25, Type: OSTypeConsole, Name: "PlayStation 3"},
	{ID: 26, Type: OSTypeConsole, Name: "PlayStation 4"},
	{ID: 27, Type: OSTypeConsole, Name: "PlayStation Portable"},
	{ID: 28, Type: OSTypeConsole, Name: "PlayStation Vita"},

	// Mobile
	{ID: 41, Type: OSTypeMobile, Name: "iOS"},
	{ID: 42, Type: OSTypeMobile, Name: "Tizen", MatchUA: []string{"Tizen"}},
	{ID: 43, Type: OSTypeMobile, Name: "Android"},
	{ID: 44, Type: OSTypeMobile, Name: "Bada"},
	{ID: 45, Type: OSTypeMobile, Name: "Nokia Series 40", MatchUA: []string{"Series40", "Series 40"}},
	{ID: 46, Type: OSTypeMobile, Name: "Symbian Series 60", MatchUA: []string{"Series60", "Series 60"}},
	{ID: 47, Type: OSTypeMobile, Name: "Symbian OS", Match: []string{"Symbian"}},
	{ID: 48, Type: OSTypeMobile, Name: "Firefox OS", Match: []string{"Firefox"}},
	{ID: 49, Type: OSTypeMobile, Name: "BlackBerry Playbook"},
	{ID: 50, Type: OSTypeMobile, Name: "BlackBerry OS 10", MatchUA: []string{"BlackBerry 10"}},
	{ID: 51, Type: OSTypeMobile, Name: "BlackBerry OS", Match: []string{"BlackBerry"}},
	{ID: 52, Type: OSTypeMobile, Name: "Windows Phone"},
	{ID: 53, Type: OSTypeMobile, Name: "MeeGo"},
	{ID: 54, Type: OSTypeMobile, Name: "SpreadTrum", MatchUA: []string{"SpreadTrum"}},
	{ID: 55, Type: OSTypeMobile, Name: "Palm OS", MatchUA: []string{"Palm OS"}},
	{ID: 56, Type: OSTypeMobile, Name: "Brew", MatchUA: []string{"Brew"}},
	{ID: 57, Type: OSTypeMobile, Name: "Java ME", MatchUA: []string{"J2ME"}},
	{ID: 58, Type: OSTypeMobile, Name: "MAUI Runtime", MatchUA: []string{"MAUI Runtime"}},
}

func init() {
	for _, os := range OSes {
		os.matchingUpdate()
	}
}

///////////////////////////////////////////////////////////////////////////////
// Methods
///////////////////////////////////////////////////////////////////////////////

// Matching by name or UA
func (os *OS) Matching(name, userAgent string) bool {
	for _, m := range os.Match {
		if strings.Contains(name, m) {
			return true
		}
	}
	for _, m := range os.MatchUA {
		if strings.Contains(userAgent, m) {
			return true
		}
	}
	return len(os.Match) < 1 && len(os.MatchUA) < 1
}

///////////////////////////////////////////////////////////////////////////////
// Internal Methods
///////////////////////////////////////////////////////////////////////////////

func (os *OS) matchingUpdate() {
	for i, m := range os.Match {
		os.Match[i] = strings.ToLower(m)
	}
	for i, m := range os.MatchUA {
		os.MatchUA[i] = strings.ToLower(m)
	}
	if len(os.Match) < 1 && len(os.MatchUA) < 1 {
		os.Match = []string{strings.ToLower(os.Name)}
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Global methods
///////////////////////////////////////////////////////////////////////////////

// GetOSByName value
func GetOSByName(name, userAgent string) *OS {
	name = strings.ToLower(name)
	userAgent = strings.ToLower(userAgent)
	for i, os := range OSes {
		if os.Matching(name, userAgent) {
			return &OSes[i]
		}
	}
	return &UndefinedOS
}

// GetOSByID value
func GetOSByID(id uint) *OS {
	if uint(len(OSes)) < id {
		return &OSes[id]
	}
	return &UndefinedOS
}
