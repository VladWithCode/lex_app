package controllers

import (
	"database/sql"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/db"
	"github.com/vladwithcode/lex_app/internal/readers"
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

func (ctl *CaseController) FindAllCases() ([]*db.LexCase, error) {
	return db.FindAllCases(ctl.ctx, ctl.appDb.Db)
}

func (ctl *CaseController) FindCases(opts *db.FindCaseOptions) ([]*db.LexCase, error) {
	if opts != nil {
		return db.FindFilteredCases(ctl.ctx, ctl.appDb.Db, opts)
	}

	return db.FindAllCases(ctl.ctx, ctl.appDb.Db)
}

func (ctl *CaseController) FindCaseById(id string) (*db.LexCase, error) {
	return db.FindCaseById(ctl.ctx, ctl.appDb.Db, id)
}

func (ctl *CaseController) FindCase(caseId, caseType string) (*db.LexCase, error) {
	caseKey := caseId + readers.CaseKeySeparator + caseType
	return db.FindCase(ctl.ctx, ctl.appDb.Db, caseKey)
}

func (ctl *CaseController) FindCaseWithAccords(id string, accordCount int) (*db.LexCase, error) {
	return db.FindCaseWithAccords(ctl.ctx, ctl.appDb.Db, id, accordCount)
}

func (ctl *CaseController) CreateCase(caseId, caseType, alias string) (*db.LexCase, error) {
	newCase, err := db.NewCase(caseId, caseType)
	if err != nil {
		return nil, err
	}

	newCase.Alias = alias
	if err := db.InsertCase(ctl.ctx, ctl.appDb.Db, newCase); err != nil {
		return nil, err
	}

	return newCase, nil
}

func (ctl *CaseController) UpdateCase(id string, caseData *db.LexCase) error {
	return db.UpdateCaseById(ctl.ctx, ctl.appDb.Db, id, caseData)
}
