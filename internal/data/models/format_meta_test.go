package models

import (
	"encoding/json"
	"os"
	"testing"
)

func TestMeta(t *testing.T) {
	f, err := os.Open("./assets/format.config.json")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	var config FormatConfig
	if err = json.NewDecoder(f).Decode(&config); err != nil {
		t.Error(err)
	}

	if len(config.Assets) != 2 {
		t.Errorf("Invalid count of assets: %d", len(config.Assets))
	}

	if len(config.Fields) != 3 {
		t.Errorf("Invalid count of fields: %d", len(config.Fields))
	}
}

func TestConfigIntersec(t *testing.T) {
	var tests = []struct {
		name   string
		c1     *FormatConfig
		c2     *FormatConfig
		result bool
	}{
		{
			name:   "empty",
			c1:     &FormatConfig{},
			c2:     &FormatConfig{},
			result: true,
		},
		{
			name: "basic_stretch",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{{}},
			},
			result: true,
		},
		{
			name: "basic_assets_similar_sizes",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1000, Height: 1000, MinHeight: 100, MinWidth: 100},
					{Required: false, Name: "icon", Width: 80, Height: 80},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1100, Height: 900, MinHeight: 200, MinWidth: 150},
					{Required: true, Name: "icon", Width: 100, Height: 100, MinHeight: 50, MinWidth: 50},
				},
			},
			result: true,
		},
		{
			name: "basic_assets_similar_sizes_negative",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1000, Height: 1000, MinHeight: 100, MinWidth: 100},
					{Required: true, Name: "icon", Width: 100, Height: 100, MinHeight: 50, MinWidth: 50},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1100, Height: 900, MinHeight: 200, MinWidth: 150},
					{Required: false, Name: "icon", Width: 80, Height: 80},
				},
			},
			result: false,
		},
		{
			name: "basic_ext",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Name: "main", Width: 300, Height: 100},
				},
				Fields: []FormatField{
					{Required: true, Name: "title"},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Width: 300, Height: 100},
					{Required: false, Name: "icon", Width: 30, Height: 30},
				},
				Fields: []FormatField{
					{Required: false, Name: "title"},
				},
			},
			result: true,
		},
		{
			name: "basic_field_negative",
			c1: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
				},
			},
			c2:     &FormatConfig{},
			result: false,
		},
		{
			name: "basic_field_negative2",
			c1: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
				},
			},
			c2: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
					{Required: true, Name: "icon"},
				},
			},
			result: false,
		},
		{
			name: "basic_field_positive",
			c1: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
					{Required: false, Name: "icon"},
				},
			},
			c2: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
					{Required: true, Name: "icon"},
				},
			},
			result: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.c1.Intersec(test.c2) != test.result {
				t.Errorf("Intersec must be %t", test.result)
			}
		})
	}
}
