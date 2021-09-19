//
// @project GeniusRabbit rotator 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package adsource

import "errors"

// Set of errors
var (
	// For bidding validation
	ErrInvalidCur                = errors.New("bid currency is not valid")
	ErrInvalidCreativeSize       = errors.New("creative size is invalid")
	ErrInvalidViewType           = errors.New("view type is invalid")
	ErrLowPrice                  = errors.New("bid Price is lower than floor price")
	ErrResponseEmpty             = errors.New("response is empty")
	ErrResponseInvalidType       = errors.New("invalid response type")
	ErrResponseInvalidGroup      = errors.New("system not support group winners")
	ErrInvalidItemInitialisation = errors.New("invalid item initialisation")
)

// NoSupportError object
type NoSupportError struct {
	NSField string
}

// Error text
func (e NoSupportError) Error() string {
	return e.NSField + " is not supported"
}
