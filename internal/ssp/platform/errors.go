//
// @project GeniusRabbit rotator 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package platform

import "errors"

// Error list
var (
	ErrInvalidResponseStatus = errors.New("Invalid response status")
	ErrNoCampaignsStatus     = errors.New("No campaigns")
)
