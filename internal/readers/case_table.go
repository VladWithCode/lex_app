package readers

import (
	"errors"
	"fmt"
	"strings"
)

type CaseData struct {
	CaseType string
	CaseId   string
	Nature   string
	Accord   string
}

func NewCaseData() *CaseData {
	return &CaseData{}
}

func (cd *CaseData) Clear() {
	cd.CaseType = ""
	cd.CaseId = ""
	cd.Nature = ""
	cd.Accord = ""
}

func (cd *CaseData) Clone() CaseData {
	return CaseData{
		CaseType: cd.CaseType,
		CaseId:   cd.CaseId,
		Nature:   cd.Nature,
		Accord:   cd.Accord,
	}
}

type CaseRow struct {
	CaseType string
	CaseId   string
	IdNo     string
	IdYear   string
	IdTrail  string
	Nature   string
	Accord   string
	AllIds   []string
}

func NewCaseRow(caseData *CaseData) (*CaseRow, error) {
	caseRow := CaseRow{
		CaseType: caseData.CaseType,
		Nature:   caseData.Nature,
		Accord:   caseData.Accord,
		AllIds:   []string{},
	}

	ids := strings.Split(
		caseData.CaseId,
		"\n",
	)

	if len(ids) == 0 {
		return nil, errors.New("Can't create a case row with now ID")
	}

	for _, id := range ids {
		parts := strings.Split(id, "/")
		if len(parts) < 2 {
			continue
		}

		caseRow.AllIds = append(caseRow.AllIds, id)

		if caseRow.CaseId == "" && isValidNumericStr(parts[0]) {
			subParts := strings.Split(parts[1], "-")
			caseRow.CaseId = id
			caseRow.IdNo = parts[0]
			caseRow.IdYear = subParts[0]
			if len(subParts) == 2 {
				caseRow.IdTrail = subParts[1]
			}
		}
	}

	if len(caseRow.AllIds) == 0 {
		return nil, fmt.Errorf("The data provided didn't produce any valid caseId:\n  %s", caseRow.CaseId)
	}

	if caseRow.CaseId == "" {
		caseRow.CaseId = caseRow.AllIds[0]
	}

	return &caseRow, nil
}

type CaseTable struct {
	Cases         []*CaseRow
	index         map[string]int
	UnparsedCases []*CaseData
}

func NewCaseTable() *CaseTable {
	return &CaseTable{
		Cases:         []*CaseRow{},
		index:         map[string]int{},
		UnparsedCases: []*CaseData{},
	}
}

func (ct *CaseTable) Add(caseRow *CaseRow) {
	ct.Cases = append(ct.Cases, caseRow)
	caseIdx := len(ct.Cases) - 1
	for _, id := range caseRow.AllIds {
		ct.index[id] = caseIdx
	}
}

func (ct *CaseTable) Find(caseId string) *CaseRow {
	idx, ok := ct.index[caseId]

	if !ok {
		return nil
	}

	ct.Cases[idx].CaseId = caseId
	return ct.Cases[idx]
}
