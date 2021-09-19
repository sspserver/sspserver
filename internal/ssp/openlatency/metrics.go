package openlatency

// MetricErrorType values
type MetricErrorType string

// Error type list...
const (
	MetricErrorHTTP    MetricErrorType = "http"
	MetricErrorNetwork MetricErrorType = "network"
)

type MetricErrorRate struct {
	Type MetricErrorType `json:"type"`
	Code string          `json:"code"`
	Rate float64         `json:"rate"`
}

type MetricsGeoRate struct {
	Country string  `json:"country"`
	Rate    float64 `json:"rate"`
}

// MetricsInfo describes basic metric information of AdNetworks integration
// All counters it's numbers per second
type MetricsInfo struct {
	ID         uint64            `json:"id"`
	Protocol   string            `json:"protocol"`
	Traceroute string            `json:"traceroute,omitempty"`
	MinLatency int64             `json:"min_latency"` // Minimal request delay in Milliseconds
	MaxLatency int64             `json:"max_latency"` // Maximal request delay in Milliseconds
	AvgLatency int64             `json:"avg_latency"`
	QPSLimit   int               `json:"qps_limit,omitempty"`
	QPS        int32             `json:"qps"`
	Success    int32             `json:"success"`
	Timeouts   int32             `json:"timeouts"`
	NoBids     int32             `json:"no_bids"`
	Errors     int32             `json:"errors"`
	ErrorRates []MetricErrorRate `json:"error_rates,omitempty"`
	GeoRates   []MetricsGeoRate  `json:"geo_rates,omitempty"`
}
