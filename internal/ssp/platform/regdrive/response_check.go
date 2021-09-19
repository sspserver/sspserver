package regdrive

import (
	"bytes"
	"fmt"

	"github.com/demdxx/gocast"
)

type validationErrorItem struct {
	Code    string
	Message string
}

type validationError struct {
	items []validationErrorItem
}

func (e *validationError) Error() string {
	var buff bytes.Buffer
	for _, it := range e.items {
		if it.Code != "" {
			buff.WriteString("\n\t" + it.Code + ": " + it.Message)
		} else {
			buff.WriteString("\n\t" + it.Message)
		}
	}
	return "[redrive] validation failed: " + buff.String()
}

// ResponseCheck rules of checking
type ResponseCheck struct {
	// Field response which related to response status
	Field string `json:"field"`

	// Comparation of response status Field with response constant
	Equal string `json:"eq"`

	// Error field in response
	Error string `json:"error"`
}

// Validate response object
func (r ResponseCheck) Validate(data interface{}) error {
	if r.Field == "" || r.Equal == "" {
		return nil
	}
	iv, _ := elementByKey(r.Field, data)
	v := gocast.ToString(iv)

	if r.Equal == v {
		return nil
	}

	if r.Error != "" {
		errors, _ := elementByKey(r.Error, data)
		switch e := errors.(type) {
		case string:
			return fmt.Errorf(`[redrive] validation failed: %s`, e)
		case []interface{}:
			err := &validationError{}
			for _, it := range e {
				err.items = append(err.items, validationErrorItemBy(it))
			}
			return err
		case []map[interface{}]interface{}:
			err := &validationError{}
			for _, it := range e {
				err.items = append(err.items, validationErrorItemBy(it))
			}
			return err
		default:
			// Not supportable response type
		}
	}

	return fmt.Errorf(`[redrive] validation failed '%s' must be equal to '%s'`, v, r.Equal)
}

func validationErrorItemBy(it interface{}) validationErrorItem {
	switch v := it.(type) {
	case string:
		return validationErrorItem{Message: v}
	default:
		m, _ := gocast.ToSiMap(v, "json", false)
		if m == nil {
			return validationErrorItem{}
		}
		var (
			msg     interface{}
			ok      bool
			code, _ = m["code"]
		)
		if msg, ok = m["text"]; !ok {
			msg, _ = m["msg"]
		}
		return validationErrorItem{
			Code:    gocast.ToString(code),
			Message: gocast.ToString(msg),
		}
	}
}
