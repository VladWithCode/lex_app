package main

import (
	"context"
	"database/sql"

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
}
