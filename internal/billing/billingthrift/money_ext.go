package billingthrift

import (
	"reflect"

	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/spi"

	"bitbucket.org/geniusrabbit/billing"
)

// MoneyExt extension
type MoneyExt struct{}

func (e MoneyExt) EncoderOf(valType reflect.Type) spi.ValEncoder {
	switch valType {
	case reflect.TypeOf(billing.Money(0)):
		return e
	}
	return nil
}

func (e MoneyExt) DecoderOf(valType reflect.Type) spi.ValDecoder {
	switch valType {
	case reflect.TypeOf((*billing.Money)(nil)):
		return e
	}
	return nil
}

/// spi.ValDecoder ...
func (e MoneyExt) Decode(val interface{}, iter spi.Iterator) {
	*val.(*billing.Money) = billing.Money(iter.ReadInt64())
}

/// spi.ValEncoder ...
func (e MoneyExt) Encode(val interface{}, stream spi.Stream) {
	stream.WriteInt64(val.(billing.Money).Int64())
}

func (e MoneyExt) ThriftType() protocol.TType {
	return protocol.TypeI64
}
