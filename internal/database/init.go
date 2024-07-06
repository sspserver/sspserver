package database

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"geniusrabbit.dev/adcorelib/context/ctxdatabase"
)

type openFnk func(dsn string) gorm.Dialector

var dialectors = map[string]openFnk{}

// Connect to database
func Connect(ctx context.Context, connection string, debug bool) (*gorm.DB, error) {
	var (
		i      = strings.Index(connection, "://")
		driver = connection[:i]
	)
	if driver == "mysql" {
		connection = connection[i+3:]
	}
	openDriver := dialectors[driver]
	if openDriver == nil {
		return nil, fmt.Errorf(`unsupported database driver %s`, driver)
	}
	db, err := gorm.Open(openDriver(connection), &gorm.Config{SkipDefaultTransaction: true})
	if err == nil && debug {
		db = db.Debug()
	}
	return db, err
}

// WithDatabase puts databases to context
func WithDatabase(ctx context.Context, master, slave *gorm.DB) context.Context {
	return ctxdatabase.WithDatabase(ctx, master, slave)
}

// ListOfDialects returns list of available DB drivers
func ListOfDialects() []string {
	list := make([]string, 0, len(dialectors))
	for d := range dialectors {
		list = append(list, d)
	}
	return list
}
