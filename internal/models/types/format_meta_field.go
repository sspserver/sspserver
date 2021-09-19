//
// @project GeniusRabbit AdNet 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package types

import (
	"errors"
	"fmt"

	"github.com/demdxx/gocast"
)

// Errors list...
var (
	ErrFormatFieldIsRequired = errors.New("field is required")
)

// FormatFieldType format
type FormatFieldType string

// Format Field types
const (
	FormatFieldStringType  FormatFieldType = "string"
	FormatFieldIntType                     = "int"
	FormatFieldFloatType                   = "float"
	FormatFieldBoolType                    = "bool"
	FormatFieldPhoneType                   = "phone"
	FormatFieldDefaultType                 = FormatFieldStringType
)

// Format field defaults
const (
	FormatFieldTitle        = "title"
	FormatFieldDescription  = "description"
	FormatFieldBrandname    = "brandname"
	FormatFieldPhone        = "phone"
	FormatFieldURL          = "url"
	FormatFieldStartDisplay = "start"
	FormatFieldContent      = "content"
	FormatFieldRating       = "rating"
	FormatFieldLikes        = "likes"
	FormatFieldAddress      = "address"
	FormatFieldSponsored    = "sponsored"
)

// FormatField description
type FormatField struct {
	// ID of the field
	ID int `json:"id,omitempty"`

	// Is required field
	Required bool `json:"required,omitempty"`

	// Title of yje field for interface
	Title string `json:"title,omitempty"`

	// Name of the field (required)
	Name string `json:"name" validation:"required"`

	// Type of the field data
	Type FormatFieldType `json:"type,omitempty"`

	// Exclude those fields if this field is filled
	Exclude []string `json:"exclude,omitempty"`

	// Select is the choice of the available values
	Select []interface{} `json:"select,omitempty"`

	// Minimum length for string or minimum value for int
	Min float64 `json:"min,omitempty"`

	// Maximum length for string or maximum value for int
	Max float64 `json:"max,omitempty"`

	// Mask pattern of data
	Mask string `json:"mask,omitempty"`

	// RegExp validation pattern
	RegExp string `json:"regexp,omitempty"`
}

// GetType of the field
func (f *FormatField) GetType() FormatFieldType {
	if f.Type == "" {
		f.Type = FormatFieldDefaultType
	}
	return f.Type
}

// IsRequired field value
func (f FormatField) IsRequired() bool {
	return f.GetType() != FormatFieldBoolType && f.Required
}

// MaxLength of the field limit
func (f FormatField) MaxLength() int {
	if f.Type != FormatFieldStringType && f.Type != FormatFieldPhoneType {
		return 0
	}
	return int(f.Max)
}

// SoftEqual compare
func (f FormatField) SoftEqual(field *FormatField) bool {
	return f.Name == field.Name && (f.GetType() == FormatFieldStringType || f.GetType() == field.GetType())
}

// Prepare and validate value
func (f FormatField) Prepare(value interface{}) (result interface{}, err error) {
	if f.IsRequired() && gocast.IsEmpty(value) {
		return nil, ErrFormatFieldIsRequired
	}

	switch f.GetType() {
	case FormatFieldStringType:
		v := gocast.ToString(value)
		if f.Min > 0 && len(v) < int(f.Min) {
			err = fmt.Errorf("min length is %d", int(f.Min))
		} else if f.Max > 0 && len(v) > int(f.Max) {
			err = fmt.Errorf("max length is %d", int(f.Max))
		}
		result = v
	case FormatFieldIntType:
		v := gocast.ToInt64(value)
		if f.Min > 0 && v < int64(f.Min) {
			err = fmt.Errorf("min value is %d", int(f.Min))
		} else if f.Max > 0 && v > int64(f.Max) {
			err = fmt.Errorf("max value is %d", int(f.Max))
		}
		result = v
	case FormatFieldFloatType:
		v := gocast.ToFloat64(value)
		if f.Min > 0 && v < f.Min {
			err = fmt.Errorf("min value is %f.3", f.Min)
		} else if f.Max > 0 && v > f.Min {
			err = fmt.Errorf("max value is %f.3", f.Max)
		}
		result = v
	case FormatFieldBoolType:
	case FormatFieldPhoneType:
	}

	return
}
