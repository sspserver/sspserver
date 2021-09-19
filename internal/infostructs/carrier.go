//
// @project Geniusrabbit::corelib 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package infostructs

// Carrier base information structure
type Carrier struct {
	ID   uint   `json:"id,omitempty"`      // Internal carrier ID
	Name string `json:"name,omitempty"`    //
	Code string `json:"carrier,omitempty"` //
}

// CarrierDefault value
var CarrierDefault Carrier
