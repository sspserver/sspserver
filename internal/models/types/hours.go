package types

import (
	"github.com/geniusrabbit/hourstable"
)

// Hours SQL type declaration
type Hours = hourstable.Hours

// HoursByString returns hours object by string pattern
func HoursByString(s string) (Hours, error) {
	return hourstable.HoursByString(s)
}
