package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/vladwithcode/lex_app/internal"
)

// App struct
type App struct {
	ctx   context.Context
	appDb *internal.AppDb
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context, db *sql.DB) {
	a.ctx = ctx
	a.appDb = internal.NewAppDb(db)

	// Migrate DB
	err := goose.SetDialect("sqlite3")
	if err != nil {
		log.Fatalf("couldn't set dialect: %v\n", err)
	}

	goose.SetBaseFS(migrations)
	err = goose.Up(db, "data/migrations")
	if err != nil {
		log.Fatalf("couldn't migrate DB: %v\n", err)
	}
}
