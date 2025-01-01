package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const DEFAULT_CONN_STR = "data/lex_app.db"

var connStr = os.Getenv("LEX_DB_URL")

var db *sql.DB

func Connect() (db *sql.DB, err error) {
	if connStr == "" {
		connStr = DEFAULT_CONN_STR
	}

	db, err = sql.Open("sqlite", "file:"+connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = 1")

	if err != nil {
		return nil, fmt.Errorf("failed PRAGMA exec: %w\n  %w", ErrForeignKeysUnsupported, err)
	}

	return db, nil
}

// Common errors
var (
	ErrGenUUID                = errors.New("error generating UUID")
	ErrForeignKeysUnsupported = errors.New("the application requires a version SQLite that allows for foreign key constraint usage.")
)
