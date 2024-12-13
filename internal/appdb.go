package internal

import (
	"database/sql"
	"sync"
)

type AppDb struct {
	Db *sql.DB
	mu sync.Mutex
}

func NewAppDb(db *sql.DB) *AppDb {
	return &AppDb{
		Db: db,
	}
}
