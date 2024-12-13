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

func (ctrl *CaseController) Startup(ctx context.Context, db *sql.DB) {
	ctrl.ctx = ctx
	ctrl.appDb = internal.NewAppDb(db)
}

func (ctrl *CaseController) FindAllCases() ([]*db.Case, error) {
	return db.FindAllCases(ctrl.ctx, ctrl.appDb.Db)
}

func (ctrl *CaseController) FindCaseById(id string) (*db.Case, error) {
	return db.FindCaseById(ctrl.ctx, ctrl.appDb.Db, id)
}

func (ctrl *CaseController) FindCase(caseId, caseType string) (*db.Case, error) {
	return db.FindCase(ctrl.ctx, ctrl.appDb.Db, caseId, caseType)
}

func (ctrl *CaseController) FindCaseWithAccords(id string, accordCount int) (*db.Case, error) {
	return db.FindCaseWithAccords(ctrl.ctx, ctrl.appDb.Db, id, accordCount)
}

func (ctrl *CaseController) CreateCase(caseId, caseType string) (*db.Case, error) {
	newCase, err := db.NewCase(caseId, caseType)
	if err != nil {
		return nil, err
	}

	if err := db.InsertCase(ctrl.ctx, ctrl.appDb.Db, newCase); err != nil {
		return nil, err
	}

	return newCase, nil
}

func (ctrl *CaseController) UpdateCase(id string, caseData *db.Case) error {
	return db.UpdateCaseById(ctrl.ctx, ctrl.appDb.Db, id, caseData)
}
