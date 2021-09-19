package counter

import (
	"sync/atomic"
)

// Counter value
type Counter int64

// Inc - rement counter value
func (cn *Counter) Inc(v int64) {
	atomic.AddInt64(cn.IntPtr(), v)
}

// Get value
func (cn *Counter) Get() int64 {
	return atomic.LoadInt64(cn.IntPtr())
}

// Set value
func (cn *Counter) Set(v int64) {
	atomic.StoreInt64(cn.IntPtr(), v)
}

// IntPtr returns the value type int64 pointer
func (cn *Counter) IntPtr() *int64 {
	return (*int64)(cn)
}
