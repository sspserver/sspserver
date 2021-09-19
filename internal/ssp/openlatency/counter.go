package openlatency

import (
	"sync/atomic"
	"time"
)

// MetricsCounter implements several counters of request metrics
type MetricsCounter struct {
	minLatency   int64 // In Milliseconds
	maxLatency   int64 // In Milliseconds
	avgLatency   int64 // In Milliseconds
	queries      int32
	success      int32
	timeouts     int32
	noBids       int32
	errors       int32
	startingTime int64
}

// NewMetricsCounter object
func NewMetricsCounter() *MetricsCounter {
	counter := &MetricsCounter{}
	counter.refresh()
	return counter
}

// UpdateQueryLatency of request
func (cnt *MetricsCounter) UpdateQueryLatency(latency time.Duration) {
	duration := int64(latency / time.Millisecond)

	atomic.StoreInt64(&cnt.avgLatency, (atomic.LoadInt64(&cnt.avgLatency)+duration)/2)

	if minLatency := atomic.LoadInt64(&cnt.minLatency); minLatency <= 0 || minLatency > duration {
		atomic.StoreInt64(&cnt.minLatency, duration)
	}

	if atomic.LoadInt64(&cnt.maxLatency) < duration {
		atomic.StoreInt64(&cnt.maxLatency, duration)
	}
}

// StartQuery new query counter
func (cnt *MetricsCounter) StartQuery() int32 {
	return atomic.AddInt32(&cnt.queries, 1)
}

// IncTimeout counter
func (cnt *MetricsCounter) IncTimeout() int32 {
	return atomic.AddInt32(&cnt.timeouts, 1)
}

// IncNobid counter
func (cnt *MetricsCounter) IncNobid() int32 {
	return atomic.AddInt32(&cnt.noBids, 1)
}

// IncSuccess counter
func (cnt *MetricsCounter) IncSuccess() int32 {
	return atomic.AddInt32(&cnt.success, 1)
}

// IncError counter
func (cnt *MetricsCounter) IncError(etype MetricErrorType, code string) {
	atomic.AddInt32(&cnt.errors, 1)
}

// FillMetrics info object
func (cnt *MetricsCounter) FillMetrics(info *MetricsInfo) {
	seconds := float64(time.Now().Sub(cnt.getStartingTime())) / float64(time.Second)
	if seconds <= 0 {
		seconds = 1
	}

	info.MinLatency = atomic.LoadInt64(&cnt.minLatency)
	info.MaxLatency = atomic.LoadInt64(&cnt.maxLatency)
	info.AvgLatency = atomic.LoadInt64(&cnt.avgLatency)
	info.QPS = counter(&cnt.queries, seconds)
	info.Success = counter(&cnt.success, seconds)
	info.Timeouts = counter(&cnt.timeouts, seconds)
	info.NoBids = counter(&cnt.noBids, seconds)
	info.Errors = counter(&cnt.errors, seconds)
}

func (cnt *MetricsCounter) refresh() {
	var (
		now     = time.Now()
		seconds = float64(now.Sub(cnt.getStartingTime())) / float64(time.Second)
	)
	if seconds <= 0 {
		seconds = 1
	}

	cnt.setStartingTime(now.Add(-time.Second))
	atomic.StoreInt64(&cnt.minLatency, atomic.LoadInt64(&cnt.avgLatency))
	atomic.StoreInt64(&cnt.maxLatency, atomic.LoadInt64(&cnt.avgLatency))
	atomic.StoreInt32(&cnt.queries, counter(&cnt.queries, seconds))
	atomic.StoreInt32(&cnt.success, counter(&cnt.success, seconds))
	atomic.StoreInt32(&cnt.timeouts, counter(&cnt.timeouts, seconds))
	atomic.StoreInt32(&cnt.noBids, counter(&cnt.noBids, seconds))
	atomic.StoreInt32(&cnt.errors, counter(&cnt.errors, seconds))
}

func (cnt *MetricsCounter) setStartingTime(tm time.Time) {
	atomic.StoreInt64(&cnt.startingTime, tm.UnixNano())
}

func (cnt *MetricsCounter) getStartingTime() time.Time {
	return time.Unix(0, atomic.LoadInt64(&cnt.startingTime))
}

func counter(cnt *int32, seconds float64) int32 {
	return int32(float64(atomic.LoadInt32(cnt)) / seconds)
}
