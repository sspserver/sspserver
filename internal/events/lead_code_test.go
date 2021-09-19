//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package events

import (
	"testing"
	"time"

	lz4 "github.com/bkaradzic/go-lz4"
)

func Test_LeadCode(t *testing.T) {
	for i := uint64(0); i < 10; i++ {
		var (
			lc = &LeadCode{
				AuctionID:         "8781fffd-6cd8-4fd0-99f2-a58962c27751",
				ImpAdID:           "fe7cb3f1-2e37-443c-bac0-9df2444aea67",
				SourceID:          1,
				ProjectID:         1,
				PublisherCompany:  1,
				AdvertiserCompany: 1,
				CampaignID:        i,
				AdID:              i * i,
				Price:             1000,
				Timestamp:         time.Now().Unix(),
			}
			err error
			res Code
		)

		// New check
		{
			if res = lc.Pack(); res.ErrorObj() != nil {
				t.Error("Invalid lead code:", res.ErrorObj())
				return
			}

			// Check compression/decompression
			if res.Compress().Decompress().URLEncode().String() != res.URLEncode().String() {
				t.Error("Compression/Decompression is not works")
			}

			if err = lc.Unpack(res.Data()); err != nil {
				t.Error(err)
				return
			}

			t.Logf("message size: %d, %d, %d", len(res.Data()),
				len(res.Compress().Data()), len(res.Compress().URLEncode().Data()))
		}
	}
}

func Benchmark_LeadCode(b *testing.B) {
	var benchmarks = []struct {
		name string
		fnk  func(lc *LeadCode)
	}{
		{name: "encode_lz4", fnk: func(lc *LeadCode) {
			_ = lc.Pack().Compress()
		}},
		{name: "encode_go-lz4", fnk: func(lc *LeadCode) {
			data, _ := lz4.Encode(nil, lc.Pack().Data())
			_ = CodeObj(data, nil)
		}},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for _, bench := range benchmarks {
		b.Run(bench.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for i := uint64(0); pb.Next(); i++ {
					lc := &LeadCode{
						AuctionID:         "8781fffd-6cd8-4fd0-99f2-a58962c27751",
						ImpAdID:           "fe7cb3f1-2e37-443c-bac0-9df2444aea67",
						SourceID:          1,
						ProjectID:         1,
						PublisherCompany:  1,
						AdvertiserCompany: 1,
						CampaignID:        i,
						AdID:              i * i,
						Price:             1000,
						Timestamp:         time.Now().Unix(),
					}
					bench.fnk(lc)
				}
			})
		})
	}
}
