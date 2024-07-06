//go:build postgres || alldb
// +build postgres alldb

package database

import (
	"fmt"
	"net/url"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	dialectors["postgres"] = openPostgres
}

func openPostgres(dsn string) gorm.Dialector {
	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}
	sslmode := "disable"
	if sslmodeVar := parsedDSN.Query().Get("sslmode"); sslmodeVar != "" {
		sslmode = sslmodeVar
	}
	password, _ := parsedDSN.User.Password()
	newDSN := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		parsedDSN.Hostname(), parsedDSN.Port(),
		parsedDSN.User.Username(), password,
		strings.TrimLeft(parsedDSN.Path, "/"),
		sslmode,
	)
	return postgres.Open(newDSN)
}
