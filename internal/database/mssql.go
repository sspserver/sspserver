//go:build mssql || alldb
// +build mssql alldb

package database

import (
	"gorm.io/driver/sqlserver"
)

func init() {
	dialectors["mssql"] = sqlserver.Open
	dialectors["sqlserver"] = sqlserver.Open
}
