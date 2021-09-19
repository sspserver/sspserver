package optimizer

import (
	"encoding/json"
	"math/rand"
	"testing"
)

func Test_OptimizerEncode(t *testing.T) {
	var optimizer Optimizer
	optimizer.Init()

	for i := 0; i < 10; i++ {
		var (
			i          = rand.Intn(10000)
			formatID   = uint(i%50 + 1)
			countryID  = byte(i % 254)
			languageID = byte(i % 50)
			deviceID   = uint(i % 1000)
			osID       = uint(i % 100)
			browserID  = uint(i % 1000)
			v          = (rand.Float64() - 0.5) / 100
		)
		if v < 0 {
			v = -1
		}
		optimizer.Inc(formatID, countryID, uint(languageID), deviceID, osID, browserID, v)
	}

	if _, err := json.Marshal(&optimizer); err != nil {
		t.Errorf("marshaling error: %v", err.Error())
	}
}

func Test_OptimizerDecode(t *testing.T) {
	var (
		optimizer Optimizer
		data      = []byte(`{"1":[{"cc":1,"ln":1,"dv":1,"os":1,"br":1,"vl":1},{"vl":0.5}]}`)
	)
	optimizer.Init()

	if err := json.Unmarshal(data, &optimizer); err != nil {
		t.Errorf("decode error: %s", err.Error())
	}

	if !optimizer.Test(1, 1, 1, 1, 1, 1, 1.0) {
		t.Error("format '1' must be in the highest preorety")
		return
	}

	if v := optimizer.getOrCreate(1, 0, 0, 0, 0, 0).Get(); v != 0.5 {
		t.Errorf("invalid format '1' value %.1f must be 0.5", v)
		return
	}

	if l := optimizer.Len(); l != 2 {
		t.Errorf("invalid elements count %d, must be 2 elements in the list", l)
		return
	}
}

func Benchmark_Optimizer(b *testing.B) {
	var optimizer Optimizer
	optimizer.Init()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var (
				i          = rand.Intn(10000)
				formatID   = uint(i%50 + 1)
				countryID  = byte(i % 254)
				languageID = byte(i % 50)
				deviceID   = uint(i % 1000)
				osID       = uint(i % 100)
				browserID  = uint(i % 1000)
				v          = (rand.Float64() - 0.5) / 100
			)
			if v < 0 {
				v = -1
			}
			optimizer.Inc(formatID, countryID, uint(languageID), deviceID, osID, browserID, v)
		}
	})
}
