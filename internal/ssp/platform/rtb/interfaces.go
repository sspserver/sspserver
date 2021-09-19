//
// @project GeniusRabbit rotator 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

// Platform implements interface Sourcer from api::v1::ssp::server.go
// Initialisation of platform based on models.RTBClient object which
// 	contains all options for configuration of this object.

package rtb

import (
	"errors"

	openrtb "github.com/bsm/openrtb"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/models"
)

// Request type enum
const (
	RequestTypeJSON      = models.RTBRequestTypeJSON
	RequestTypeXML       = models.RTBRequestTypeXML
	RequestTypeProtobuff = models.RTBRequestTypeProtoBUFF
)

// Errors set
var (
	ErrResponseAreNotSecure  = errors.New("Response are not secure")
	ErrInvalidResponseStatus = errors.New("Invalid response status")
)

type platformExt interface {
	// InitPlatform extension
	InitPlatform(platform *Platform) error
}

type platformRequestPreparer interface {
	// PrepareRequest RTB object
	PrepareRequest(request *openrtb.BidRequest) error
}

type platformResponsePreparer interface {
	// PrepareResponse of RTB
	PrepareResponse(bid *openrtb.BidResponse, request *adsource.BidRequest) error
}

type requestMarshalel interface {
	// ContentType of request
	ContentType() string

	// Marshal for RTB
	Marshal(request interface{}) ([]byte, error)
}
