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

type LexCase struct {
	Id             string    `json:"id" db:"id"`
	CaseId         string    `json:"caseId" db:"case_id"`
	CaseType       string    `json:"caseType" db:"case_type"`
	CaseYear       string    `json:"caseYear" db:"case_year"`
	CaseNo         string    `json:"caseNo" db:"case_no"`
	LastUpdatedAt  time.Time `json:"lastUpdatedAt" db:"last_updated_at"`
	LastAccessedAt time.Time `json:"lastAccessedAt" db:"last_accessed_at"`
	Alias          string    `json:"alias" db:"alias"`
	OtherIds       []string  `json:"otherIds" db:"other_ids"`
	Accords        []*Accord `json:"accords"`
}

func NewEmptyCase() *LexCase {
	return &LexCase{
		OtherIds: []string{},
		Accords:  []*Accord{},
	}
}

func NewCase(caseId, caseType string) (*LexCase, error) {
	if !isValidCaseId(caseId) {
		return nil, fmt.Errorf("%s is not a valid caseId value:\n\t%w", caseId, ErrorInvalidCaseId)
	}

	c := &LexCase{
		CaseId:   caseId,
		CaseType: caseType,
		OtherIds: []string{},
		Accords:  []*Accord{},
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
	c.CaseNo = parts[0]
	c.CaseYear = parts[1]

	c.AddOtherId(caseId)

	return c, nil
}

func (c *LexCase) SetIdsFromStr(str string) error {
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

func (c *LexCase) AddOtherId(candidate string) error {
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

type FindCaseOptions struct {
	Limit          int
	CaseId         string
	CaseType       string
	CaseYear       string
	CaseNo         string
	LastUpdatedAt  string
	IncludeAccords bool
	MaxAccords     int
}

var DefaultFindCaseOptions = FindCaseOptions{
	Limit:          0,
	CaseId:         "",
	CaseType:       "",
	CaseYear:       "",
	CaseNo:         "",
	LastUpdatedAt:  "",
	IncludeAccords: false,
	MaxAccords:     1,
}

func InsertCase(ctx context.Context, appDb *sql.DB, caseData *LexCase) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	otherIds := strings.Join(caseData.OtherIds, otherIdsSeparator)
	_, err := appDb.ExecContext(
		ctx,
		`INSERT INTO cases (id, case_id, case_type, case_year, case_no, alias, other_ids)
		VALUES (:Id, :CaseId, :CaseType, :CaseYear, :CaseNo, :Alias, :OtherIds)`,
		sql.Named("Id", caseData.Id),
		sql.Named("CaseId", caseData.CaseId),
		sql.Named("CaseType", caseData.CaseType),
		sql.Named("CaseYear", caseData.CaseYear),
		sql.Named("CaseNo", caseData.CaseNo),
		sql.Named("Alias", caseData.Alias),
		sql.Named("OtherIds", otherIds),
	)

	if err != nil {
		return err
	}

	return nil
}

func FindAllCases(ctx context.Context, appDb *sql.DB) ([]*LexCase, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rows, err := appDb.QueryContext(
		ctx,
		"SELECT id, case_id, case_type, case_year, case_no, alias, other_ids FROM cases",
	)
	if err != nil {
		return nil, err
	}

	cases := []*LexCase{}
	nOtherIds := sql.NullString{}
	nCaseYear := sql.NullString{}
	nCaseNo := sql.NullString{}
	nAlias := sql.NullString{}

	for rows.Next() {
		nOtherIds.Valid = false
		nCaseYear.Valid = false
		nCaseNo.Valid = false
		nAlias.Valid = false
		c := &LexCase{}

		rows.Scan(
			&c.Id,
			&c.CaseId,
			&c.CaseType,
			&nCaseYear,
			&nCaseNo,
			&nAlias,
			&nOtherIds,
		)
		if nCaseYear.Valid {
			c.CaseYear = nCaseYear.String
		}
		if nCaseNo.Valid {
			c.CaseNo = nCaseNo.String
		}
		if nAlias.Valid {
			c.Alias = nAlias.String
		}

		if nOtherIds.Valid {
			c.SetIdsFromStr(nOtherIds.String)
		}

		cases = append(cases, c)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	return cases, nil
}

func FindFilteredCases(ctx context.Context, appDb *sql.DB, opts *FindCaseOptions) ([]*LexCase, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	baseQuery := ""
	args := []interface{}{}
	conditions := []string{}
	if opts.IncludeAccords {
		baseQuery = `SELECT 
cases.id, cases.case_id, cases.case_type, cases.case_year, cases.case_no, cases.alias,
cases.other_ids, accords.accord_id, accords.content, unixepoch(accords.date) as date, accords.raw_data
FROM cases LEFT JOIN (
	SELECT
		id as accord_id, for_case, content, date, raw_data,
		ROW_NUMBER() OVER (PARTITION BY for_case ORDER BY date DESC NULLS LAST) as rn
	FROM accords
) accords ON cases.id = accords.for_case AND accords.rn <= :accordCount`
		args = append(args, sql.Named("accordCount", opts.MaxAccords))
	} else {
		baseQuery = "SELECT * FROM cases"
	}

	if opts.CaseId != "" {
		conditions = append(conditions, "case_id LIKE '%%:caseId%'")
		args = append(args, sql.Named("caseId", opts.CaseId))
	}
	if opts.CaseType != "" {
		conditions = append(conditions, "case_type LIKE '%%:caseType%'")
		args = append(args, sql.Named("caseType", opts.CaseType))
	}
	if opts.CaseYear != "" {
		conditions = append(conditions, "case_year LIKE '%%:caseYear%'")
		args = append(args, sql.Named("caseYear", opts.CaseYear))
	}
	if opts.CaseNo != "" {
		conditions = append(conditions, "case_no LIKE '%%:caseNo%'")
		args = append(args, sql.Named("caseNo", opts.CaseNo))
	}
	if opts.LastUpdatedAt != "" {
		conditions = append(conditions, "julianday(accords.date) >= julianday(:lastUpdated)")
		args = append(args, sql.Named("lastUpdated", opts.LastUpdatedAt))
	}

	if len(conditions) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(conditions, " AND "))
	}

	if opts.IncludeAccords {
		baseQuery = fmt.Sprintf("%s ORDER BY accords.date DESC NULLS LAST", baseQuery)
	}
	if opts.Limit > 0 {
		baseQuery = fmt.Sprintf("%s LIMIT :limit", baseQuery)
		args = append(args, sql.Named("limit", opts.Limit))
	}

	rows, err := appDb.QueryContext(
		ctx,
		baseQuery,
		args...,
	)
	if err != nil {
		return nil, err
	}

	cases := []*LexCase{}
	caseMap := map[string]int{}
	for rows.Next() {
		var (
			id       = ""
			caseId   = ""
			caseType = ""
			caseYear = sql.NullString{}
			caseNo   = sql.NullString{}
			alias    = sql.NullString{}
			othIds   = sql.NullString{}
			accord   = Accord{}
			accDate  = 0
		)

		rows.Scan(
			&id,
			&caseId,
			&caseType,
			&caseYear,
			&caseNo,
			&alias,
			&othIds,
			&accord.Id,
			&accord.Content,
			&accDate,
			&accord.rawData,
		)
		if accDate != 0 {
			accord.Date = time.Unix(int64(accDate), 0)
		}

		if cIdx, ok := caseMap[id]; ok {
			cases[cIdx].Accords = append(cases[cIdx].Accords, &accord)
		} else {
			c := NewEmptyCase()
			c.Id = id
			c.CaseId = caseId
			c.CaseType = caseType
			c.CaseYear = caseYear.String
			c.CaseNo = caseNo.String
			c.Alias = alias.String
			c.Accords = []*Accord{&accord}
			if othIds.Valid {
				c.SetIdsFromStr(othIds.String)
			}

			// Curent length is next idx
			caseMap[c.Id] = len(cases)
			cases = append(cases, c)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cases, nil
}

func FindCasesById(ctx context.Context, appDb *sql.DB, ids []string) ([]*LexCase, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := appDb.QueryContext(ctx, "SELECT id, case_id, case_type FROM cases WHERE id IN :Ids", sql.Named("Ids", ids))
	if err != nil {
		return nil, err
	}

	cases := []*LexCase{}
	for rows.Next() {
		c := NewEmptyCase()
		rows.Scan(
			&c.Id,
			&c.CaseId,
			&c.CaseType,
		)

		cases = append(cases, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cases, nil
}

func FindCaseById(ctx context.Context, appDb *sql.DB, id string) (*LexCase, error) {
	row := appDb.QueryRowContext(
		ctx,
		`SELECT id, case_id, case_type, case_year, case_no, alias, other_ids FROM cases WHERE id = :Id`,
		sql.Named("Id", id),
	)

	c := &LexCase{}
	otherIds := new(string)
	err := row.Scan(
		&c.Id,
		&c.CaseId,
		&c.CaseType,
		&c.CaseYear,
		&c.CaseNo,
		&c.Alias,
		&otherIds,
	)
	if err != nil {
		return nil, err
	}

	if len(*otherIds) > 0 {
		c.SetIdsFromStr(*otherIds)
	}

	return c, nil
}

func FindCases(ctx context.Context, appDb *sql.DB, caseKeys []string) ([]*LexCase, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := appDb.QueryContext(
		ctx,
		"SELECT id, case_id, case_type, (case_id || ':' || case_type) as case_key FROM cases WHERE case_key IN :CaseKeys",
		sql.Named("CaseKeys", caseKeys),
	)
	if err != nil {
		return nil, err
	}

	cases := []*LexCase{}
	for rows.Next() {
		c := NewEmptyCase()
		rows.Scan(
			&c.Id,
			&c.CaseId,
			&c.CaseType,
		)

		cases = append(cases, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cases, nil
}

func FindCase(ctx context.Context, appDb *sql.DB, caseKey string) (*LexCase, error) {
	row := appDb.QueryRowContext(
		ctx,
		"SELECT id, case_id, case_type, case_year, case_no, alias, other_ids, (case_id || ':' || case_type) as case_key FROM cases WHERE case_key = :CaseKey",
		sql.Named("CaseKey", caseKey),
	)

	c := NewEmptyCase()
	otherIds := new(string)
	err := row.Scan(
		&c.Id,
		&c.CaseId,
		&c.CaseType,
		&c.CaseYear,
		&c.CaseNo,
		&c.Alias,
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

func FindCaseWithAccords(ctx context.Context, appDb *sql.DB, id string, accordCount int) (*LexCase, error) {
	c := NewEmptyCase()
	rows, err := appDb.QueryContext(
		ctx,
		`SELECT
			cases.id,
			cases.case_id,
			cases.case_type,
			cases.alias,
			accords.id,
			accords.content,
			date(accords.date) as date,
			accords.raw_data
		FROM cases
		LEFT JOIN accords
		ON cases.id = accords.for_case
		WHERE cases.id = $1
		LIMIT $2`,
		id,
		accordCount,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			acId      sql.NullString
			acContent sql.NullString
			acDate    sql.NullTime
			acRawData sql.NullString
		)
		rows.Scan(
			&c.Id,
			&c.CaseId,
			&c.CaseType,
			&c.Alias,
			&acId,
			&acContent,
			&acDate,
			&acRawData,
		)

		if acId.Valid {
			c.Accords = append(c.Accords, &Accord{
				Id:      acId.String,
				Content: acContent.String,
				Date:    acDate.Time,
				rawData: acRawData.String,
				ForCase: c.Id,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return c, nil
}

func UpdateCaseById(ctx context.Context, appDb *sql.DB, id string, newCaseData *LexCase) error {
	cols := make([]string, 0)
	args := make([]interface{}, 0)

	if newCaseData.CaseId != "" {
		if !isValidCaseId(newCaseData.CaseId) {
			return fmt.Errorf("Can't insert/update case with invalid id: %s\n  %w", newCaseData.CaseId, ErrorInvalidCaseId)
		}

		cols = append(cols, "case_id = :CaseId")
		args = append(args, sql.Named("CaseId", newCaseData.CaseId))
	}

	if newCaseData.CaseType != "" {
		cols = append(cols, "case_type = :CaseType")
		args = append(args, sql.Named("CaseType", newCaseData.CaseType))
	}

	if newCaseData.Alias != "" {
		cols = append(cols, "alias = :Alias")
		args = append(args, sql.Named("Alias", newCaseData.Alias))
	}

	if newCaseData.OtherIds != nil {
		cols = append(cols, "other_ids = :OtherIds")
		args = append(args, sql.Named("OtherIds", newCaseData.OtherIds))
	}

	if len(cols) == 0 {
		return errors.New("No fields to update")
	}

	args = append(args, sql.Named("Id", id))

	query := fmt.Sprintf("UPDATE cases SET %s WHERE id = :Id", strings.Join(cols, ", "))

	_, err := appDb.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCaseById(ctx context.Context, appDb *sql.DB, id string) error {
	_, err := appDb.ExecContext(ctx, "DELETE FROM cases WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
