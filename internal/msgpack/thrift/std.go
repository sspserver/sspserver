package thrift

import (
	"io"

	thrift "github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/spi"
)

// NewAPI of thrift encode/decode
func NewAPI(ext ...spi.Extension) thrift.API {
	return thrift.Config{
		Protocol:      thrift.ProtocolBinary,
		StaticCodegen: false,
		Extensions:    ext,
	}.Froze()
}

// Marshal to []byte
func Marshal(obj interface{}) ([]byte, error) {
	return thrift.Marshal(obj)
}

// Unmarshal message
func Unmarshal(buf []byte, obj interface{}) error {
	return thrift.Unmarshal(buf, obj)
}

// NewDecoder object
func NewDecoder(reader io.Reader, buf []byte) *thrift.Decoder {
	return thrift.NewDecoder(reader, buf)
}

// NewEncoder object
func NewEncoder(writer io.Writer) *thrift.Encoder {
	return thrift.NewEncoder(writer)
}
