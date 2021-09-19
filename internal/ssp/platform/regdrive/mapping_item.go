package regdrive

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/demdxx/gocast"
)

type mappingItem struct {
	Target   string            `json:"target"`
	Default  string            `json:"default,omitempty"`
	Required bool              `json:"required,omitempty"`
	ToLower  bool              `json:"to_lower,omitempty"`
	Format   string            `json:"format,omitempty"`
	Map      map[string]string `json:"map,omitempty"`
}

// MappingItem contains information about target item value
type MappingItem struct {
	Target   string
	Default  string
	Required bool
	ToLower  bool
	Format   string
	Map      map[string]string
}

// MarshalJSON implements the json.Marshaler
func (it *MappingItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(mappingItem{
		Target:   it.Target,
		Default:  it.Default,
		Required: it.Required,
		ToLower:  it.ToLower,
		Format:   it.Format,
		Map:      it.Map,
	})
}

// UnmarshalJSON implements the json.Unmarshaller
func (it *MappingItem) UnmarshalJSON(b []byte) (err error) {
	if len(b) > 1 && b[0] == '"' {
		var str string
		if err = json.Unmarshal(b, &str); err != nil {
			return
		}
		it.Target = str
		it.Default = ""
		it.Required = false
		it.ToLower = false
		it.Format = ""
		it.Map = nil
		return
	}

	var item mappingItem
	if err = json.Unmarshal(b, &item); err != nil {
		return
	}
	it.Target = item.Target
	it.Default = item.Default
	it.Required = item.Required
	it.ToLower = item.ToLower
	it.Format = item.Format
	it.Map = item.Map
	return
}

// PrepareValue target value
func (it *MappingItem) PrepareValue(vl interface{}) (string, error) {
	var v string
	switch vt := vl.(type) {
	case time.Time:
		if it.Format != "" {
			v = vt.Format("2006-01-02")
		} else {
			v = vt.Format("2006-01-02")
		}
	case string:
		v = vt
	default:
		v = gocast.ToString(v)
	}

	v = strings.TrimSpace(v)
	if it.ToLower {
		v = strings.ToLower(v)
	}
	if it.Map != nil {
		if v2, ok := it.Map[v]; ok {
			v = v2
		}
	}

	if v == "" {
		v = it.Default
	}

	if it.Required && v == "" {
		return "", fmt.Errorf(`field [%s] is required`, it.Target)
	}
	return v, nil
}
