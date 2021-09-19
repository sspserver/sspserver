package msgpack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	// "github.com/alecthomas/binary"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/geniusrabbit/billing"
	billingthrift "bitbucket.org/geniusrabbit/billing/thrift"
	"bitbucket.org/geniusrabbit/corelib/msgpack/thrift"
)

type msgStd struct {
	Time        int64         `thrift:",1" json:"tm,omitempty"`
	ID          string        `thrift:",2" json:"id,omitempty"`
	Source      uint          `thrift:",3" json:"src,omitempty"`
	Project     uint          `thrift:",4" json:"project,omitempty"`
	IPString    string        `thrift:",5" json:"ip,omitempty"`
	Country     string        `thrift:",6" json:"cc,omitempty"`
	Request     string        `thrift:",7" json:"request,omitempty"`
	Response    string        `thrift:",8" json:"response,omitempty"`
	Error       string        `thrift:",9" json:"err,omitempty"`
	LatencyTime uint64        `thrift:",10" json:"ltm,omitempty"`
	Money       billing.Money `thrift:",11" json:"money,omitempty"`
}

var thriftapi = thrift.NewAPI(billingthrift.MoneyExt{})

func TestMessagePack(t *testing.T) {
	var (
		msg = msgStd{
			Time:        time.Now().UnixNano(),
			ID:          "1233-dda22t-e3oeqq-1233",
			Source:      1,
			Project:     212,
			IPString:    "127.0.0.1",
			Country:     "US",
			Request:     "message",
			Response:    "adss",
			Error:       "error",
			LatencyTime: 12312312,
			Money:       billing.MoneyFloat(11.2),
		}
		tests = []struct {
			name   string
			pack   func(obj interface{}) ([]byte, error)
			unpack func(data []byte, r interface{}) error
		}{
			{name: "lz4json", pack: StdPack, unpack: StdUnpack},
			{name: "json", pack: json.Marshal, unpack: json.Unmarshal},
			{name: "thrift", pack: thriftapi.Marshal, unpack: thriftapi.Unmarshal},
			// {name: "binary", pack: binary.Marshal, unpack: binary.Unmarshal},
		}
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				data, err = test.pack(&msg)
				msg2      msgStd
			)
			assert.NoError(t, err, test.name+" pack")
			assert.NoError(t, test.unpack(data, &msg2), test.name+" unpack")

			if !reflect.DeepEqual(msg, msg2) {
				t.Error("message must be equal")
			}
		})
	}
}

func BenchmarkMessagePack(b *testing.B) {
	var (
		msg = msgStd{
			Time:        time.Now().UnixNano(),
			ID:          "1233-dda22t-e3oeqq-1233",
			Source:      1,
			Project:     21322,
			IPString:    "127.0.0.1",
			Country:     "US",
			Request:     "message",
			Response:    "adss",
			Error:       "error",
			LatencyTime: 12312312,
			Money:       billing.MoneyFloat(11.2),
		}
		benches = []struct {
			name string
			pack func(obj interface{}) ([]byte, error)
		}{
			{name: "lz4json", pack: StdPack},
			{name: "json", pack: json.Marshal},
			{name: "thrift", pack: thriftapi.Marshal},
			// {name: "binary", pack: binary.Marshal},
		}
	)

	for _, bench := range benches {
		fmt.Println("")
		b.Run(bench.name, func(b *testing.B) {
			var (
				count    = 0
				dataSize = 0
			)

			b.ResetTimer()
			b.ReportAllocs()
			b.RunParallel(func(b *testing.PB) {
				for b.Next() {
					if data, err := bench.pack(&msg); err == nil {
						dataSize += len(data)
					}
					count++
				}
			})

			fmt.Println("MSG:", count, "AVG:", dataSize/count, "DATA:", dataSize)
		})
	}
}
