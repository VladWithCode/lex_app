package controllers

import (
	"database/sql"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/db"
	"golang.org/x/net/context"
)

type CaseController struct {
	ctx   context.Context
	appDb *internal.AppDb
}

func NewCaseControler() *CaseController {
	return &CaseController{}
}

func (ctl *CaseController) Startup(ctx context.Context, db *sql.DB) {
	ctl.ctx = ctx
	ctl.appDb = internal.NewAppDb(db)
}

func (ctl *CaseController) FindAllCases() ([]*db.Case, error) {
	return db.FindAllCases(ctl.ctx, ctl.appDb.Db)
}

func (ctl *CaseController) FindCases(opts *db.FindCaseOptions) ([]*db.Case, error) {
	if opts != nil {
		return db.FindFilteredCases(ctl.ctx, ctl.appDb.Db, opts)
	}

	return db.FindAllCases(ctl.ctx, ctl.appDb.Db)
}

func (ctl *CaseController) FindCaseById(id string) (*db.Case, error) {
	return db.FindCaseById(ctl.ctx, ctl.appDb.Db, id)
}

func (ctl *CaseController) FindCase(caseId, caseType string) (*db.Case, error) {
	return db.FindCase(ctl.ctx, ctl.appDb.Db, caseId, caseType)
}

func (ctl *CaseController) FindCaseWithAccords(id string, accordCount int) (*db.Case, error) {
	return db.FindCaseWithAccords(ctl.ctx, ctl.appDb.Db, id, accordCount)
}

func (ctl *CaseController) CreateCase(caseId, caseType string) (*db.Case, error) {
	newCase, err := db.NewCase(caseId, caseType)
	if err != nil {
		return nil, err
	}

	if err := db.InsertCase(ctl.ctx, ctl.appDb.Db, newCase); err != nil {
		return nil, err
	}

	return newCase, nil
}

func (ctl *CaseController) UpdateCase(id string, caseData *db.Case) error {
	return db.UpdateCaseById(ctl.ctx, ctl.appDb.Db, id, caseData)
}
