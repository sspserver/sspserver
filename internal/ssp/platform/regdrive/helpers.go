package regdrive

import (
	"strings"
)

func elementByKey(key string, data interface{}) (interface{}, bool) {
	var (
		v    = (interface{})(data)
		keys = strings.Split(key, "->")
	)
	for _, k := range keys {
		switch data := v.(type) {
		case map[string]interface{}:
			v, _ = data[k]
		case map[interface{}]interface{}:
			v, _ = data[k]
		case nil:
			return "", false
		default:
			return "", false
		}
	}
	return v, true
}
