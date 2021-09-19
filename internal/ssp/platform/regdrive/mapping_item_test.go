package regdrive

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_MappingItemJSONDecode(t *testing.T) {
	var (
		it    MappingItem
		tests = []struct {
			name   string
			code   string
			target MappingItem
		}{
			{
				name: "target1",
				code: `{"target":"target1"}`,
				target: MappingItem{
					Target: "target1",
				},
			},
			{
				name: "target2",
				code: `{"target":"target2","to_lower":true}`,
				target: MappingItem{
					Target:  "target2",
					ToLower: true,
				},
			},
			{
				name: "target3",
				code: `{"target":"target3","to_lower":true,"map":{"f":"Female","m":"Male"}}`,
				target: MappingItem{
					Target:  "target3",
					ToLower: true,
					Map: map[string]string{
						"f": "Female",
						"m": "Male",
					},
				},
			},
			{
				name: "target4",
				code: `"target4"`,
				target: MappingItem{
					Target: "target4",
				},
			},
		}
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(test.code), &it); err != nil {
				t.Errorf("error decode item: %s", err.Error())
			}
			if !reflect.DeepEqual(it, test.target) {
				t.Errorf("decode object must be equal")
			}
		})
	}
}
