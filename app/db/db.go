package db

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once  sync.Once
	db    *sql.DB
	dbErr error
)

func DbOnce(connStr string) (*sql.DB, error) {
	once.Do(func() {
		db, dbErr = sql.Open("postgres", connStr)
		if dbErr != nil {
			return
		}
		dbErr = db.Ping()
	})
	return db, dbErr
}
