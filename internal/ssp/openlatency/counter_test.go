package openlatency

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkCounter(b *testing.B) {
	var counter = NewMetricsCounter()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			now := time.Now()
			counter.StartQuery()
			switch rand.Intn(4) {
			case 0:
				counter.IncNobid()
			case 1:
				counter.IncTimeout()
			case 2:
				counter.IncError(MetricErrorHTTP, "500")
			default:
				counter.IncSuccess()
			}
			counter.UpdateQueryLatency(time.Now().Sub(now))
		}
	})

	var metrics MetricsInfo
	counter.FillMetrics(&metrics)
	b.Log(metrics.QPS, metrics.Success, metrics.Timeouts, metrics.NoBids,
		metrics.AvgLatency, metrics.MaxLatency, metrics.MinLatency)
}
