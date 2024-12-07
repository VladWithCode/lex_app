package controllers

import (
	"github.com/vladwithcode/lex_app/internal/db"
	"golang.org/x/net/context"
)

type CaseController struct {
	ctx context.Context
}

func NewCaseControler() *CaseController {
	return &CaseController{}
}

func (ctrl *CaseController) Startup(ctx context.Context) {
	ctrl.ctx = ctx
}

func (ctrl *CaseController) FindAllCases() ([]*db.Case, error) {
	return db.FindAllCases()
}

func (ctrl *CaseController) FindCaseById(id string) (*db.Case, error) {
	return db.FindCaseById(id)
}

func (ctrl *CaseController) FindCase(caseId, caseType string) (*db.Case, error) {
	return db.FindCase(caseId, caseType)
}

func (ctrl *CaseController) CreateCase(caseId, caseType string) error {
	newCase, err := db.NewCase(caseId, caseType)
	if err != nil {
		return err
	}

	return db.InsertCase(newCase)
}

func (ctrl *CaseController) UpdateCase(id string, caseData *db.Case) error {
	return db.UpdateCaseById(id, caseData)
}
