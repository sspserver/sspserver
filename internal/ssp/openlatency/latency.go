package openlatency

import "time"

const (
	HTTPHeaderRequestTimemark = "X-Request-Ts" // In milliseconds
	HTTPHeaderRequestNode     = "X-Request-Node"
)

// RequestInitTime converts UnixNano into Milliseconds
func RequestInitTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
