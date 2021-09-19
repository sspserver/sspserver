package searchtypes

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestUIntBitset(t *testing.T) {
	bits := NewUIntBitset(1, 3, 5, 7)

	bits.Set(2)
	if !bits.Has(2) {
		t.Error("bits Set/Has not working")
	}

	bits = bits.Unset(3, 1).Unset(3, 7)
	if bits.Has(3) {
		t.Error("bits Unet/Has not working")
	}

	if !reflect.DeepEqual([]uint{2, 5, 7}, bits.Values()) {
		t.Errorf("bits Types problems: %v", bits.Values())
	}
}

func BenchmarkUIntBitset(b *testing.B) {
	var (
		cursor   int
		variants []uint
		bits     = NewUIntBitset(1, 3)
	)

	// Init random numbers
	for i := 0; i < 10000; i++ {
		variants = append(variants, uint(rand.Uint64()%1000))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch {
		case i%2 == 0 || i%37 == 0:
			bits.Set(variants[cursor])
		case i%3 == 0 || i%5 == 0 || i%11 == 0:
			bits = bits.Unset(variants[cursor])
		default:
			bits.Has(variants[cursor])
		}
		if cursor++; cursor >= len(variants) {
			cursor = 0
		}
	}
}
