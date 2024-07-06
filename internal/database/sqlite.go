//go:build sqlite || alldb
// +build sqlite alldb

package database

import (
	"gorm.io/driver/sqlite"
)

func init() {
	dialectors["sqlite"] = sqlite.Open
}
