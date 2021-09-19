package dataencode

import (
	"encoding/json"
	"encoding/xml"
	"net/url"

	"github.com/demdxx/gocast"

	"geniusrabbit.dev/sspserver/internal/models"
)

type marshalFnk func(interface{}) ([]byte, error)

// Encoder of the request tupe data
type Encoder interface {
	// ContentType of request
	ContentType() string

	// Marshal for RTB
	Marshal(request interface{}) ([]byte, error)
}

type commonEncoder struct {
	contentType string
	marshalel   marshalFnk
}

// NewEncoder endcoder by RTB request type
func NewEncoder(requestType models.RTBRequestType) Encoder {
	return &commonEncoder{
		contentType: contentTypeByRequestType(requestType),
		marshalel:   marshalelFuncByRequestType(requestType),
	}
}

func (enc *commonEncoder) ContentType() string                         { return enc.contentType }
func (enc *commonEncoder) Marshal(request interface{}) ([]byte, error) { return enc.marshalel(request) }

func marshalelFuncByRequestType(requestType models.RTBRequestType) marshalFnk {
	switch requestType {
	case models.RTBRequestTypeJSON:
		return json.Marshal
	case models.RTBRequestTypeXML:
		return xml.Marshal
	case models.RTBRequestTypeProtoBUFF:
		return nil
	case models.RTBRequestTypePOSTFormEncoded:
		return formUrlencoded
	}
	// * There is cant be plain text marshalel of data, only UNMASHALEL
	return json.Marshal
}

func contentTypeByRequestType(requestType models.RTBRequestType) string {
	switch requestType {
	case models.RTBRequestTypeJSON:
		return "application/json"
	case models.RTBRequestTypeXML:
		return "application/xml"
	case models.RTBRequestTypeProtoBUFF:
		return "application/protobuff"
	case models.RTBRequestTypePOSTFormEncoded:
		return "application/x-www-form-urlencoded"
	}
	return "plain/text"
}

func formUrlencoded(data interface{}) ([]byte, error) {
	var (
		values    = url.Values{}
		vals, err = gocast.ToStringMap(data, "", false)
	)
	if err != nil {
		return nil, err
	}
	for key, val := range vals {
		values.Add(key, val)
	}
	return []byte(values.Encode()), nil
}
