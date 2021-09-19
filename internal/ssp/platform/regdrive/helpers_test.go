package regdrive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_elementByKey(t *testing.T) {
	v, _ := elementByKey("a", map[string]interface{}{"a": "x"})
	assert.Equal(t, "x", v, "invalid level 1")

	v, _ = elementByKey("a->b", map[string]interface{}{
		"a": map[interface{}]interface{}{"b": "x"},
	})
	assert.Equal(t, "x", v, "invalid level 2")

	v, _ = elementByKey("a->b->c", map[string]interface{}{
		"a": map[interface{}]interface{}{
			"b": map[interface{}]interface{}{"c": "x"},
		},
	})
	assert.Equal(t, "x", v, "invalid level 3")
}
