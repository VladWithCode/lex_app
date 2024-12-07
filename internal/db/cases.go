package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	otherIdsSeparator  = ","
	casePartsSeparator = "/"
	caseTrailSeparator = "-"
)

var (
	ErrorInvalidCaseId = errors.New("caseId is not in a valid format. Should be formatted as '123/2024[-I]'")
)

type Case struct {
	Id            string    `json:"id" db:"id"`
	CaseId        string    `json:"caseId" db:"case_id"`
	CaseType      string    `json:"caseType" db:"case_type"`
	CaseYear      *string   `json:"caseYear" db:"case_year"`
	CaseNo        *string   `json:"caseNo" db:"case_no"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt" db:"last_updated_at"`
	Alias         *string   `json:"alias" db:"alias"`
	OtherIds      []string  `json:"otherIds" db:"other_ids"`
}

func NewEmptyCase() *Case {
	return &Case{
		CaseYear: new(string),
		CaseNo:   new(string),
		Alias:    new(string),
		OtherIds: []string{},
	}
}

func NewCase(caseId, caseType string) (*Case, error) {
	if !isValidCaseId(caseId) {
		return nil, fmt.Errorf("%s is not a valid caseId value:\n\t%w", caseId, ErrorInvalidCaseId)
	}

	c := &Case{
		CaseId:   caseId,
		CaseType: caseType,
		CaseYear: new(string),
		CaseNo:   new(string),
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("Error creating a new UUID for case %v:\n\t%w", caseId, err)
	}

	c.Id = id.String()

	parts := strings.Split(
		strings.Split(caseId, caseTrailSeparator)[0],
		casePartsSeparator,
	)
	c.CaseNo = &parts[0]
	c.CaseYear = &parts[1]

	c.AddOtherId(caseId)

	return c, nil
}

func (c *Case) SetIdsFromStr(str string) error {
	ids := strings.Split(str, otherIdsSeparator)
	if len(ids) == 0 {
		return errors.New("String generated no id candidates")
	}

	for _, candidate := range ids {
		candidate = strings.TrimSpace(candidate)
		if !isValidCaseId(candidate) {
			return fmt.Errorf("String %s is not a valid id", candidate)
		}

		c.AddOtherId(candidate)
	}

	return nil
}

func (c *Case) AddOtherId(candidate string) error {
	if !isValidCaseId(candidate) {
		return fmt.Errorf("%s is not a valid caseId value:\n\t%w", candidate, ErrorInvalidCaseId)
	}

	c.OtherIds = append(c.OtherIds, candidate)

	return nil
}

func isValidCaseId(candidate string) bool {
	parts := strings.Split(candidate, casePartsSeparator)

	if parts[0] == "" {
		return false
	}

	yearParts := strings.Split(parts[1], caseTrailSeparator)

	if _, err := strconv.Atoi(yearParts[0]); err != nil {
		return false
	}

	return true
}

func InsertCase(caseData *Case) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	otherIds := strings.Join(caseData.OtherIds, otherIdsSeparator)
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO cases (id, case_id, case_type, case_year, case_no, alias, other_ids)
		VALUES (@Id, @CaseId, @CaseType, @CaseYear, @CaseNo, @Alias, @OtherIds)`,
		sql.Named("Id", caseData.Id),
		sql.Named("CaseId", caseData.CaseId),
		sql.Named("CaseType", caseData.CaseType),
		sql.Named("CaseYear", *caseData.CaseYear),
		sql.Named("CaseNo", *caseData.CaseNo),
		sql.Named("Alias", *caseData.Alias),
		sql.Named("OtherIds", otherIds),
	)

	if err != nil {
		return err
	}

	return nil
}

func FindAllCases() ([]*Case, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.QueryContext(
		ctx,
		`SELECT (id, case_id, case_type, case_year, case_no, alias, other_ids) FROM cases`,
	)
	if err != nil {
		return nil, err
	}

	cases := []*Case{}
	for rows.Next() {
		c := &Case{}
		otherIds := new(string)
		err = rows.Scan(
			&c.Id,
			&c.CaseId,
			&c.CaseType,
			c.CaseYear,
			c.CaseNo,
			c.Alias,
			otherIds,
		)
		if err != nil {
			return nil, err
		}

		if len(*otherIds) > 0 {
			c.SetIdsFromStr(*otherIds)
		}

		cases = append(cases, c)
	}

	return cases, nil
}

func FindCaseById(id string) (*Case, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := db.QueryRowContext(
		ctx,
		`SELECT (id, case_id, case_type, case_year, case_no, alias, other_ids) FROM cases WHERE id = @Id`,
		sql.Named("Id", id),
	)

	c := &Case{}
	otherIds := new(string)
	err := row.Scan(
		&c.Id,
		&c.CaseId,
		&c.CaseType,
		c.CaseYear,
		c.CaseNo,
		c.Alias,
		otherIds,
	)
	if err != nil {
		return nil, err
	}

	if len(*otherIds) > 0 {
		c.SetIdsFromStr(*otherIds)
	}

	return c, nil
}

func FindCase(caseId, caseType string) (*Case, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := db.QueryRowContext(
		ctx,
		"SELECT (id, case_id, case_type, case_year, case_no, alias, other_ids) FROM cases WHERE case_id = @CaseId AND case_type = @CaseType",
		sql.Named("CaseId", caseId),
		sql.Named("CaseType", caseType),
	)

	c := NewEmptyCase()
	otherIds := new(string)
	err := row.Scan(
		&c.Id,
		&c.CaseId,
		&c.CaseType,
		c.CaseYear,
		c.CaseNo,
		c.Alias,
		otherIds,
	)

	if err != nil {
		return nil, err
	}

	if len(*otherIds) > 0 {
		c.SetIdsFromStr(*otherIds)
	}

	return c, nil
}

func UpdateCaseById(id string, newCaseData *Case) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cols := make([]string, 0)
	args := make([]interface{}, 0)

	if newCaseData.CaseId != "" {
		if !isValidCaseId(newCaseData.CaseId) {
			return fmt.Errorf("Can't insert/update case with invalid id: %s\n  %w", newCaseData.CaseId, ErrorInvalidCaseId)
		}

		cols = append(cols, "case_id = @CaseId")
		args = append(args, sql.Named("CaseId", newCaseData.CaseId))
	}

	if newCaseData.CaseType != "" {
		cols = append(cols, "case_type = @CaseType")
		args = append(args, sql.Named("CaseType", newCaseData.CaseType))
	}

	if newCaseData.Alias != nil {
		cols = append(cols, "alias = @Alias")
		args = append(args, sql.Named("Alias", newCaseData.Alias))
	}

	if newCaseData.OtherIds != nil {
		cols = append(cols, "other_ids = @OtherIds")
		args = append(args, sql.Named("OtherIds", newCaseData.OtherIds))
	}

	if len(cols) == 0 {
		return errors.New("No fields to update")
	}

	args = append(args, sql.Named("Id", id))

	query := fmt.Sprintf("UPDATE cases SET %s WHERE id = @Id", strings.Join(cols, ", "))

	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCaseById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "DELETE FROM cases WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
