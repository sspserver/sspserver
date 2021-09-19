//
// @project GeniusRabbit
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 - 2018
//

package useragent

// Browser agent type
type Browser struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

// Browsers list
var Browsers = []Browser{
	{ID: 0, Name: "Other"},
	{ID: 1, Name: "Chrome"},
	{ID: 2, Name: "Chrome Mobile"},
	{ID: 3, Name: "Firefox"},
	{ID: 4, Name: "Firefox Mobile"},
	{ID: 5, Name: "IE"},
	{ID: 6, Name: "IE Mobile"},
	{ID: 7, Name: "Opera"},
	{ID: 8, Name: "Opera Mini"},
	{ID: 9, Name: "Safari"},
	{ID: 10, Name: "Mobile Safari"},
}

var (
	browserCodes map[string]uint8
)

func init() {
	browserCodes = map[string]uint8{}
	for _, b := range Browsers {
		browserCodes[b.Name] = b.ID
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Global methods
///////////////////////////////////////////////////////////////////////////////

// GetBrowserByCode returns browser from code
func GetBrowserByCode(code string) *Browser {
	id := browserCodes[code]
	return GetBrowserByID(id)
}

// GetBrowserByID returns browser from ID
func GetBrowserByID(id uint8) *Browser {
	if id < 1 || len(Browsers) <= int(id) {
		return &Browsers[0]
	}
	return &Browsers[id]
}
