package regdrive

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	errInvalidResponseStatus      = errors.New("regdrive: invalid response status")
	errInvalidRTBSource           = errors.New("regdrive: invalid RTB source")
	errInvalidWinEventStream      = errors.New("regdrive: invalid win event stream")
	errInvalidEventStream         = errors.New("regdrive: invalid event stream")
	errInvalidLogger              = errors.New("regdrive: invalid logger object")
	errInvalidRedirectBody        = errors.New("regdrive: invalid redirect body")
	errInvalidResponseContentType = errors.New("regdrive: invalid response content type; Supports: application/json, application/xml, plain/text")
	errInvalidResponseMapping     = errors.New("regdrive: invalid response mapping of data")
	errInvalidResponseData        = errors.New("regdrive: invalid response data")
)

type redirectError struct {
	TargetURL string
}

func newRedirectError(url string) *redirectError {
	return &redirectError{TargetURL: url}
}

func (err *redirectError) Error() string {
	return protocol + ": prevent second redirect to: " + err.TargetURL
}

type dataErrorWrapper struct {
	data interface{}
	err  error
}

func newDataError(data interface{}, err error) dataErrorWrapper {
	return dataErrorWrapper{
		data: data,
		err:  err,
	}
}

func (er dataErrorWrapper) Error() string {
	data, _ := json.Marshal(er.data)
	return fmt.Sprintf("%s: %s", er.err.Error(), string(data))
}

func (er dataErrorWrapper) Data() interface{} {
	return er.data
}
