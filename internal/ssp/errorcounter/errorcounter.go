package counter

import (
	"math"
	"math/rand"
	"sync/atomic"
)

const treshold = 1000.0

// ErrorCounter of errors
type ErrorCounter struct {
	value int32
}

// Inc - rement counter value
func (cn *ErrorCounter) Inc() {
	if newVal := atomic.AddInt32(&cn.value, 1); newVal > treshold {
		atomic.StoreInt32(&cn.value, treshold)
	}
}

// Dec - rement counter value
func (cn *ErrorCounter) Dec() {
	if newVal := atomic.AddInt32(&cn.value, -1); newVal < -treshold {
		atomic.StoreInt32(&cn.value, -treshold)
	}
}

// Next value
func (cn *ErrorCounter) Next() bool {
	v := cn.factor()
	return v <= 0.01 || rand.Float64()*1.03 > v
}

func (cn *ErrorCounter) factor() float64 {
	return notmalise(cn.getValue() / treshold * 7.8)
}

func (cn *ErrorCounter) getValue() float64 {
	return float64(atomic.LoadInt32(&cn.value))
}

func notmalise(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
