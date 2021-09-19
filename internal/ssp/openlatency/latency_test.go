package openlatency

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func Benchmark_RequestInitTime(b *testing.B) {
	tm := time.Now()
	b.Logf("milliseconds %d", RequestInitTime(tm))

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RequestInitTime(tm)
		}
	})
}

func Benchmark_Format1(b *testing.B) {
	var (
		idx     int32
		numbers = []uint64{0, 1, 598, 10833, 400300, 2310210, 1549107735603, 1549107735604, 1549107735605, 1549107735606, 1549107735624}
	)

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("std", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				i := atomic.AddInt32(&idx, 1) % int32(len(numbers))
				_ = strconv.FormatUint(numbers[i], 10)
			}
		})
	})

	b.Run("my", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				i := atomic.AddInt32(&idx, 1) % int32(len(numbers))
				_ = formatUint(numbers[i])
			}
		})
	})

	for _, v := range numbers {
		fmt.Println(">>>>", v, formatUint(v))
	}
}

var numbers []string

func init() {
	for i := 0; i < 1000; i++ {
		numbers = append(numbers, fmt.Sprintf("%03d", i))
	}
}

func formatUint(val uint64) string {
	switch {
	case val < 10:
		return numbers[val][2:]
	case val > 9 && val < 100:
		return numbers[val][1:]
	case val < 1000:
		return numbers[val]
	}

	var (
		buf  [128]byte
		size = len(buf)
	)

	for val != 0 {
		chunk := val % 1000
		val = val / 1000

		nm := numbers[chunk]
		if val > 0 {
			buf[size-4] = nm[0]
			buf[size-3] = nm[1]
			buf[size-2] = nm[2]
			size -= 3
		} else {
			buf[size-4] = nm[0]
			if chunk > 9 && chunk < 100 {
				buf[size-3] = nm[1]
				size -= 2
			} else if chunk > 100 {
				buf[size-3] = nm[1]
				buf[size-2] = nm[2]
				size -= 3
			} else {
				size--
			}
		}
	}
	return string(buf[size-1:])
}
