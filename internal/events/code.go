//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package events

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"

	// "github.com/pierrec/lz4"
	lz4 "github.com/bkaradzic/go-lz4"
	basethrift "github.com/thrift-iterator/go"

	"bitbucket.org/geniusrabbit/corelib/msgpack/thrift"
)

// LZ4BlockMaxSize constant
const LZ4BlockMaxSize = 64 << 10

// ErrEmptyData error
var ErrEmptyData = errors.New(`data is empty`)

// Code structure can contains and conver data for URL
type Code struct {
	data []byte
	err  error
}

// CodeObj by values
func CodeObj(data []byte, err error) Code {
	return Code{data: data, err: err}
}

// ThriftObjectCode converts object to code object with thrift parser
func ThriftObjectCode(obj interface{}, apis ...basethrift.API) Code {
	var (
		buff bytes.Buffer
		enc  *basethrift.Encoder
	)

	if len(apis) > 0 && apis[0] != nil {
		enc = apis[0].NewEncoder(&buff)
	} else {
		enc = thrift.NewEncoder(&buff)
	}

	if err := enc.Encode(obj); err != nil {
		return CodeObj(nil, err)
	}
	return CodeObj(buff.Bytes(), nil)
}

func (c Code) String() string {
	return string(c.data)
}

// Data value
func (c Code) Data() []byte {
	return c.data
}

// Compress code.data
func (c Code) Compress() Code {
	if c.err != nil {
		return c
	}

	return CodeObj(lz4.Encode(nil, c.data))

	// var (
	// 	newCode Code
	// 	lzwr    = lz4.NewWriter(&newCode)
	// )
	// lzwr.Header.BlockMaxSize = LZ4BlockMaxSize
	// defer lzwr.Close()

	// if _, newCode.err = lzwr.Write(c.data); newCode.err == nil {
	// 	newCode.err = lzwr.Flush()
	// }

	// if newCode.err != nil {
	// 	newCode.ResetData()
	// }

	// // data, err := lz4.Encode(nil, c.data)
	// return newCode
}

// Decompress current code.data
func (c Code) Decompress() Code {
	if c.err != nil {
		return c
	}

	return CodeObj(lz4.Decode(nil, c.data))

	// 	if len(c.data) < 1 {
	// 		return CodeObj(nil, ErrEmptyData)
	// 	}

	// 	lzrd := lz4.NewReader(bytes.NewBuffer(c.data))
	// 	data, err := ioutil.ReadAll(lzrd)

	// 	// data, err := lz4.Decode(nil, c.data)
	// 	return CodeObj(data, err)
}

func (c Code) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

// ErrorObj response
func (c Code) ErrorObj() error {
	return c.err
}

// URLEncode data for URL using
func (c Code) URLEncode() Code {
	if len(c.data) < 1 {
		return CodeObj(nil, ErrEmptyData)
	}
	return CodeObj([]byte(base64.URLEncoding.EncodeToString(c.data)), nil)
}

// URLDecode value
func (c Code) URLDecode() Code {
	data, err := base64.URLEncoding.DecodeString(string(c.data))
	return CodeObj(data, err)
}

// DecodeThriftObject to target
func (c Code) DecodeThriftObject(target interface{}, apis ...basethrift.API) error {
	if c.err != nil {
		return c.err
	}

	var dec *basethrift.Decoder
	if len(apis) > 0 && apis[0] != nil {
		dec = apis[0].NewDecoder(nil, c.data)
	} else {
		dec = thrift.NewDecoder(nil, c.data)
	}
	return dec.Decode(target)
}

// ResetData object
func (c *Code) ResetData() {
	if len(c.data) > 0 {
		c.data = c.data[:]
	}
}

// Write method is implementation of io.Writer interface{}
func (c *Code) Write(p []byte) (n int, err error) {
	c.data = append(c.data, p...)
	return len(p), nil
}

var _ io.Writer = (*Code)(nil)
