package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vladwithcode/lex_app/internal"
	_ "modernc.org/sqlite"
)

const DEFAULT_CONN_STR = "data/lex_app.db"

var connStr = os.Getenv("LEX_DB_URL")
var debug = os.Getenv("LEX_DEBUG") == "true"

var db *sql.DB

func Connect() (db *sql.DB, err error) {
	if debug && connStr == "" {
		connStr = DEFAULT_CONN_STR
	} else if !debug {
		appDir, err := internal.GetAppDataDir()
		if err != nil {
			return nil, err
		}
		connStr = filepath.Join(appDir, "data/lex_app.db")

		if err := EnsureDBFileExists(connStr); err != nil {
			return nil, err
		}
	}

	db, err = sql.Open("sqlite", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = 1")

	if err != nil {
		return nil, fmt.Errorf("failed PRAGMA exec: %w\n  %w", ErrForeignKeysUnsupported, err)
	}

	return db, nil
}

func EnsureDBFileExists(dbPath string) error {
	_, err := os.Stat(dbPath)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(dbPath), 0775)
		if err != nil && !os.IsExist(err) {
			return err
		}
		f, err := os.Create(dbPath)
		if err != nil {
			return err
		}
		defer f.Close()
		err = f.Chmod(0775)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

// Common errors
var (
	ErrGenUUID                = errors.New("error generating UUID")
	ErrForeignKeysUnsupported = errors.New("the application requires a version SQLite that allows for foreign key constraint usage")
)
