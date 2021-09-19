package optimizer

import (
	"math"
	"math/rand"
	"sync/atomic"
)

const (
	probabilityInertness = 1000
	growStateBase        = 3.0  // Speed of growing
	reduseStateBase      = -1.0 // Speed of reducing
	initialState         = float64(1.0)
)

type atomicFloatItem uint64

func (i *atomicFloatItem) Get() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(i)))
}

func (i *atomicFloatItem) Set(n float64) float64 {
	n = math.Max(math.Min(n, 1.0), 0)
	atomic.StoreUint64((*uint64)(i), math.Float64bits(n))
	return n
}

func (i *atomicFloatItem) Inc(value float64) float64 {
	return i.Set(
		(i.Get()*probabilityInertness + growValue(value)) / (probabilityInertness + 1),
	)
}

// Test current state
func (i *atomicFloatItem) Test(min float64) bool {
	state := math.Max(i.Get(), math.Max(min, 0.005))
	return state >= 0.95 || rand.Float64() <= state
}

func newAtomicFloatItem() atomicFloatItem {
	return atomicFloatItem(math.Float64bits(initialState))
}

// in a range of -1 : 1
func sugmoid(x float64) float64 {
	if x >= 1 || x <= -1 {
		return x
	}
	return 1.0 / (1.0 + math.Exp(-x*10))
}

func growValue(x float64) float64 {
	if x > 0 {
		return growStateBase + x + sugmoid(x)
	} else if x == 0 {
		return 0
	}
	return reduseStateBase + sugmoid(x)
}
