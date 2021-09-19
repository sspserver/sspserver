//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package types

import (
	"fmt"
	"testing"
)

func Test_Formats(t *testing.T) {
	var tests = []struct {
		Format       Format
		TargetFormat Format
		Target       int
	}{
		{
			Format:       Format{Width: 300, Height: 250},
			TargetFormat: Format{Width: 300, Height: 250},
			Target:       0,
		},
		{
			Format:       Format{Width: 300, Height: 300},
			TargetFormat: Format{Width: 300, Height: 250},
			Target:       1,
		},
		{
			Format:       Format{Width: 300, Height: 300},
			TargetFormat: Format{Width: 500, Height: 250},
			Target:       -1,
		},
		{
			Format:       Format{Width: 300, Height: 300, MinWidth: 100, MinHeight: 100},
			TargetFormat: Format{Width: 150, Height: 150},
			Target:       0,
		},
		{
			Format:       Format{Width: 300, Height: 300, MinWidth: 100, MinHeight: 100},
			TargetFormat: Format{Width: 50, Height: 150},
			Target:       1,
		},
		{
			Format:       Format{Width: 300, Height: 300, MinWidth: 100, MinHeight: 100},
			TargetFormat: Format{Width: 300, Height: 350},
			Target:       -1,
		},
		{
			Format:       Format{Width: 300, Height: 300, MinWidth: -1, MinHeight: -1},
			TargetFormat: Format{Width: 150, Height: 150},
			Target:       0,
		},
		{
			Format:       Format{Width: 300, Height: 300},
			TargetFormat: Format{Width: 500, Height: 300, MinWidth: -1, MinHeight: -1},
			Target:       0,
		},
		{
			Format:       Format{Width: 300, Height: 400, MinWidth: 0, MinHeight: -1},
			TargetFormat: Format{Width: 500, Height: 300, MinWidth: -1, MinHeight: -1},
			Target:       0,
		},
		{
			Format:       Format{Width: 300, Height: 250},
			TargetFormat: Format{Width: 300, Height: 100, MinWidth: -1, MinHeight: -1},
			Target:       1,
		},
	}

	for _, tt := range tests {
		t.Run("format_"+tt.TargetFormat.String(), func(t *testing.T) {
			if tt.Format.SuitsCompare(&tt.TargetFormat) != tt.Target {
				t.Errorf("Invalid test format SuitsCompare: %s != %s : %d => %d",
					tt.Format, tt.TargetFormat, tt.Format.SuitsCompare(&tt.TargetFormat), tt.Target)
			}
		})
	}
}

func Test_Aspect(t *testing.T) {
	var tests = []struct {
		V1, MinV1 int
		V2, MinV2 int
		target    int
	}{
		{
			V1:     0,
			MinV1:  -1,
			V2:     1000,
			MinV2:  400,
			target: 0,
		},
		{
			V1:     100,
			MinV1:  0,
			V2:     100,
			MinV2:  0,
			target: 0,
		},
		{
			V1:     150,
			MinV1:  50,
			V2:     100,
			MinV2:  50,
			target: 0,
		},
		{
			V1:     150,
			MinV1:  0,
			V2:     100,
			MinV2:  50,
			target: 1,
		},
		{
			V1:     150,
			MinV1:  110,
			V2:     100,
			MinV2:  50,
			target: 1,
		},
		{
			V1:     50,
			MinV1:  0,
			V2:     100,
			MinV2:  0,
			target: -1,
		},
		{
			V1:     90,
			MinV1:  50,
			V2:     100,
			MinV2:  0,
			target: -1,
		},
		{
			V1:     200,
			MinV1:  50,
			V2:     100,
			MinV2:  20,
			target: 0,
		},
		{
			V1:     200,
			MinV1:  50,
			V2:     300,
			MinV2:  150,
			target: 0,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d_%d___%d_%d", test.V1, test.MinV1, test.V2, test.MinV2), func(t *testing.T) {
			if compareAspect(test.V1, test.MinV1, test.V2, test.MinV2) != test.target {
				t.Errorf("Fail compare, should be %d", test.target)
			}
		})
	}
}
