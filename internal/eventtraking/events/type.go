package events

import (
	"reflect"

	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/spi"
)

// Type of event
type Type string

func (t Type) String() string { return string(t) }

// Event types
const (
	Undefined  Type = ""
	Request    Type = "request"
	Impression Type = "impression"
	View       Type = "view"
	Direct     Type = "direct"
	Click      Type = "click"
	Lead       Type = "lead"
	// Source types
	SourceNoBid Type = "src.nobid"
	SourceBid   Type = "src.bid"
	SourceWin   Type = "src.win"
	SourceFail  Type = "src.fail"
	// Access Point types
	AccessPointNoBid Type = "ap.nobid"
	AccessPointBid   Type = "ap.bid"
	AccessPointWin   Type = "ap.win"
	AccessPointFail  Type = "ap.fail"
)

// typeThriftExy extension
type typeThriftExy struct{}

func (e typeThriftExy) EncoderOf(valType reflect.Type) spi.ValEncoder {
	switch valType {
	case reflect.TypeOf(Undefined):
		return e
	}
	return nil
}

func (e typeThriftExy) DecoderOf(valType reflect.Type) spi.ValDecoder {
	switch valType {
	case reflect.TypeOf((*Type)(nil)):
		return e
	}
	return nil
}

/// spi.ValDecoder ...
func (e typeThriftExy) Decode(val interface{}, iter spi.Iterator) {
	*val.(*Type) = Type(iter.ReadString())
}

/// spi.ValEncoder ...
func (e typeThriftExy) Encode(val interface{}, stream spi.Stream) {
	stream.WriteString(val.(Type).String())
}

func (e typeThriftExy) ThriftType() protocol.TType {
	return protocol.TypeString
}
