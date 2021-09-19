package optimizer

import (
	"fmt"
	"testing"
)

func Test_atomicFloatItem(t *testing.T) {
	var (
		item  = newAtomicFloatItem()
		tests = []struct {
			name    string
			initial float64
			value   float64
			target  float64
		}{
			{
				name:    "decrease_to_0_5",
				initial: 1.0,
				value:   -1,
				target:  0.5,
			},
			{
				name:    "decrease_to_0",
				initial: 0.5,
				value:   -1,
				target:  0,
			},
			{
				name:    "rise_from_0",
				initial: 0,
				value:   0.001,
				target:  0.5,
			},
			{
				name:    "rise_from_0_5",
				initial: 0.5,
				value:   0.001,
				target:  1,
			},
		}
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				v    = item.Set(test.initial)
				iter = 0
			)
			for ; iter < 1000; iter++ {
				v = item.Inc(test.value)
				if test.target > test.initial {
					if v >= test.target {
						break
					}
				} else if v <= test.target {
					break
				}
			}

			if test.target > test.initial {
				if v < test.target {
					t.Errorf("Target value %.3f must be greater then %.3f", v, test.target)
				}
			} else if v > test.target {
				t.Errorf("Target value %.3f must be less then %.3f", v, test.target)
			}

			fmt.Println("Iterations", iter)
			fmt.Println("Result", v)
		})
	}
}
