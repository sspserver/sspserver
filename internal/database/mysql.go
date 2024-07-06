// +build mysql alldb

package database

import (
	// _ "gorm.io/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
)

func init() {
	dialectors["mysql"] = mysql.Open
}
