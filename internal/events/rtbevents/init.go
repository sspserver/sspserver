package rtbevents

import msgjson "github.com/geniusrabbit/adcorelib/msgpack/json"

var (
	streamCodeEncoder = &msgjson.EncodeGenerator{}
	streamCodeDecoder = &msgjson.DecodeGenerator{}
)
