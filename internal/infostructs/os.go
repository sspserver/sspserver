//
// @project Geniusrabbit::BigBrother 2016, 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2018
//

package infostructs

// OS base information structure
type OS struct {
	ID      uint   `json:"id"` // Internal system ID
	Name    string `json:"name,omitempty"`
	Version string `json:"ver,omitempty"`
}

// OSDefault value
var OSDefault OS
