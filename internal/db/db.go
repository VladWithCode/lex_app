package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const DEFAULT_CONN_STR = "data/lex_app.db"

var connStr = os.Getenv("LEX_DB_URL")

var db *sql.DB

func Connect() (db *sql.DB, err error) {
	if connStr == "" {
		wd, _ := os.Getwd()
		connStr = wd + "/" + DEFAULT_CONN_STR
	}
	fmt.Printf("connStr: %v\n", connStr)
	db, err = sql.Open("sqlite3", connStr)
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
