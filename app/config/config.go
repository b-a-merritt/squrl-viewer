package config

import (
	"database/sql"
)

type Config struct {
	ConnStr string
	DB      *sql.DB
}
