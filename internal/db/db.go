package db

import (
	"database/sql"
	"os"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "The application requires a version SQLite that allows for foreign key constraint usage.")
	}

	return db, nil
}
