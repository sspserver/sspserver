//
// @project RevolvEads
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015
//

package useragent

import (
	"strings"
	"sync"

	"github.com/ua-parser/uap-go/uaparser"
)

var (
	parser       *uaparser.Parser
	uaCache      map[string]*uaparser.Client
	uaCacheLimit int = 0
	uaCacheSync  sync.Mutex
)

var (
	MOBILE_DEVICE_FAMILIES = []string{
		`iPhone`,
		`iPod`,
		`Generic Smartphone`,
		`Generic Feature Phone`,
		`PlayStation Vita`,
	}

	PC_OS_FAMILIES = []string{
		`Windows 95`,
		`Windows 98`,
		`Windows ME`,
		`Solaris`,
	}

	MOBILE_OS_FAMILIES = []string{
		`Windows Phone`,
		`Windows Phone OS`, // Earlier versions of ua-parser returns Windows Phone OS
		`Symbian OS`,
		`Bada`,
		`Windows CE`,
		`Windows Mobile`,
	}

	MOBILE_BROWSER_FAMILIES = []string{
		`Opera Mobile`,
		`Opera Mini`,
	}

	TABLET_DEVICE_FAMILIES = []string{
		`iPad`,
		`BlackBerry Playbook`,
		`Blackberry Playbook`, // Earlier versions of ua-parser returns "Blackberry" instead of "BlackBerry"
		`Kindle`,
		`Kindle Fire`,
		`Kindle Fire HD`,
		`Galaxy Tab`,
		`Xoom`,
		`Dell Streak`,
	}

	TOUCH_CAPABLE_OS_FAMILIES = []string{
		`iOS`,
		`Android`,
		`Windows Phone`,
		`Windows Phone OS`,
		`Windows RT`,
		`Windows CE`,
		`Windows Mobile`,
	}

	TOUCH_CAPABLE_DEVICE_FAMILIES = []string{
		`BlackBerry Playbook`,
		`Blackberry Playbook`,
		`Kindle Fire`,
	}
)

type UserAgent struct {
	*uaparser.Client
	device  Device
	browser int
	os      int
	Raw     string
}

func (ua *UserAgent) GetBrowser() *Browser {
	if ua.browser < 0 {
		browser := GetBrowserByCode(ua.UserAgent.Family)
		ua.browser = int(browser.ID)
		return browser
	}
	return GetBrowserByID(uint8(ua.browser))
}

func (ua *UserAgent) GetOs() *OS {
	if ua.os < 0 {
		os := GetOSByName(ua.Os.Family+" "+ua.Device.Family, ua.Raw)
		ua.os = int(os.ID)
		return os
	}
	return GetOSByID(uint(ua.os))
}

func (ua *UserAgent) GetDevice() Device {
	if ua.device < 0 {
		if ua.IsTab() {
			ua.device = DEVICE_TABLET
		} else if ua.IsMobile() {
			ua.device = DEVICE_MOBILE
		} else if ua.IsTV() {
			ua.device = DEVICE_TV
		} else {
			ua.device = DEVICE_WEB
		}
	}
	return ua.device
}

///////////////////////////////////////////////////////////////////////////////
/// Check
///////////////////////////////////////////////////////////////////////////////

func (ua *UserAgent) IsMobile() bool {
	// First check for mobile device and mobile browser families
	return aHas(ua.Device.Family, MOBILE_DEVICE_FAMILIES) ||
		aHas(ua.UserAgent.Family, MOBILE_BROWSER_FAMILIES) ||
		// Device is considered Mobile OS is Android and not tablet
		// This is not fool proof but would have to suffice for now
		(ua.Os.Family == `Android` && !ua.IsTab()) ||
		(ua.Os.Family == `BlackBerry OS` && ua.Device.Family == `Blackberry Playbook`) ||
		aHas(ua.Os.Family, MOBILE_OS_FAMILIES) ||
		// TODO: remove after https://github.com/tobie/ua-parser/issues/126 is closed
		(strings.Contains(ua.Raw, `J2ME`) || strings.Contains(ua.Raw, `MIDP`)) ||
		// This is here mainly to detect Google's Mobile Spider
		strings.Contains(ua.Raw, `iPhone`) ||
		// Mobile Spiders should be identified as mobile
		(ua.Device.Family == `Spider` && strings.Contains(ua.UserAgent.Family, `Mobile`))

}

func (ua *UserAgent) IsTab() bool {
	return aHas(ua.Device.Family, TABLET_DEVICE_FAMILIES) ||
		(ua.Os.Family == `Android` && ua.IsAndroidTablet()) ||
		strings.Contains(ua.Os.Family, `Windows RT`)
}

func (ua *UserAgent) IsPC() bool {
	// Returns True for "PC" devices (Windows, Mac and Linux)
	return strings.Contains(ua.Raw, `Windows NT`) ||
		aHas(ua.Os.Family, PC_OS_FAMILIES) ||
		// TODO: remove after https://github.com/tobie/ua-parser/issues/127 is closed
		(ua.Os.Family == `Mac OS X` && strings.Contains(ua.Raw, `Silk`)) ||
		(strings.Contains(ua.Raw, `Linux`) && strings.Contains(ua.Raw, `X11`))
}

func (ua *UserAgent) IsBot() bool {
	return ua.Device.Family == `Spider`
}

func (ua *UserAgent) IsAndroidTablet() bool {
	// Newer Android tablets don't have "Mobile" in their user agent string,
	// older ones like Galaxy Tab still have "Mobile" though they're not
	return !strings.Contains(ua.Raw, "Mobile Safari") &&
		ua.UserAgent.Family != "Firefox Mobile"
}

func (ua *UserAgent) IsBlackberryTouchCapableDevice() bool {
	// A helper to determine whether a BB phone has touch capabilities
	// Blackberry Bold Touch series begins with 99XX
	return strings.Contains(ua.Device.Family, `Blackberry 99`) ||
		strings.Contains(ua.Device.Family, `Blackberry 95`) ||
		strings.Contains(ua.Device.Family, `Blackberry 99`)
}

func (ua *UserAgent) IsTV() bool {
	return strings.Contains(ua.Device.Family, "tv")
}

///////////////////////////////////////////////////////////////////////////////
// Global Methods
///////////////////////////////////////////////////////////////////////////////

func InitParser(regexFile string) (*uaparser.Parser, error) {
	return uaparser.New(regexFile)
}

func InitParserWithCache(regexFile string, limit int) (*uaparser.Parser, error) {
	uaCacheLimit = limit
	return InitParser(regexFile)
}

// Parse user-agent
//
// @param ua UseAgent raw string
// @return UserAgent struct
func ParseUserAgent(ua string) *UserAgent {
	var uaObject *uaparser.Client

	// Check cache and lock point
	if uaCacheLimit > 0 {
		uaCacheSync.Lock()
		defer uaCacheSync.Unlock()
		if uaCache == nil || len(uaCache) > uaCacheLimit {
			uaCache = make(map[string]*uaparser.Client)
		}
		uaObject = uaCache[ua]
	}

	// Parse user-agent and put to cache if necessary
	if uaObject == nil {
		uaObject = parser.Parse(ua)
		if uaCacheLimit > 0 {
			uaCache[ua] = uaObject
		}
	}

	return &UserAgent{
		Client:  uaObject,
		device:  -1,
		browser: -1,
		os:      -1,
		Raw:     ua,
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func indexOfs(v string, arr []string) int {
	for i, s := range arr {
		if s == v {
			return i
		}
	}
	return -1
}

func aHas(v string, arr []string) bool {
	return indexOfs(v, arr) >= 0
}
