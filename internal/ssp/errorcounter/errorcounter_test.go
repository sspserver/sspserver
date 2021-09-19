package counter

import (
	"testing"
)

func Test_ErrorErrorCounter(t *testing.T) {
	var cnt ErrorCounter
	for i := 0; i < 1000; i++ {
		cnt.Inc()
	}
	t.Logf("value: %.3f = %.3f", cnt.getValue(), cnt.factor())
	for i := 0; i < 2000; i++ {
		cnt.Dec()
	}
	t.Logf("value: %.3f = %.3f", cnt.getValue(), cnt.factor())
}
