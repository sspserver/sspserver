package infostructs

import (
	"reflect"
	"testing"
)

func TestPrepareDomain(t *testing.T) {
	var tests = []struct {
		domain string
		result []string
	}{
		{
			domain: "example.com",
			result: []string{
				"example.com",
				"*.example.com",
				"*.com",
			},
		},
		{
			domain: "super.puper.example.com",
			result: []string{
				"super.puper.example.com",
				"*.super.puper.example.com",
				"*.puper.example.com",
				"*.example.com",
				"*.com",
			},
		},
		{
			domain: "",
			result: []string{
				"*.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.domain, func(t *testing.T) {
			if domain := PrepareDomain(test.domain); !reflect.DeepEqual(test.result, domain) {
				t.Errorf("Invalid domain [%s] prepare: %v", test.domain, domain)
			}
		})
	}
}
