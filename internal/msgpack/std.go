package msgpack

import (
	"bytes"
	"encoding/json"

	"github.com/pierrec/lz4"
)

// StdPack message
func StdPack(msg interface{}) ([]byte, error) {
	var (
		buff   bytes.Buffer
		writer = lz4.NewWriter(&buff)
		err    = json.NewEncoder(writer).Encode(msg)
	)
	if err != nil {
		return nil, err
	}
	if err = writer.Flush(); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// StdUnpack message
func StdUnpack(data []byte, msg interface{}) error {
	return json.NewDecoder(lz4.NewReader(bytes.NewReader(data))).Decode(msg)
}
