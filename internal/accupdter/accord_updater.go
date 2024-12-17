package accupdter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/db"
	"github.com/vladwithcode/lex_app/internal/fetchers"
	"github.com/vladwithcode/lex_app/internal/readers"
)

type CaseStore interface {
	FindAll(ids []string) ([]*db.LexCase, error)
	FindAllKeys(keys []string) ([]*db.LexCase, error)
	FindById(id string) (*db.LexCase, error)
	FindByKey(key string) (*db.LexCase, error)

	Save(updates []*UpdatedAccord) error
}

type AccUpdter interface {
	FindUpdates(keys []string, ids *[]string) (updates []*UpdatedAccord, notFoundKeys []string, err error)
	Update(keys []string, ids *[]string) (notFoundKeys []string, err error)

	getStore() *CaseStore
}

// TODO: Implement update queue
/* type AccUpdterQueue struct {
	updaters []*AccUpdter
	mu       sync.Mutex
} */

type UpdatedAccord struct {
	CaseKey  string
	CaseType internal.CaseType
	CaseId   string
	Content  string
	Date     time.Time
	Nature   string
	OthIds   []string
}

type AccUpdterOpts struct {
	Store           *CaseStore
	CaseType        internal.CaseType
	Region          internal.Region
	MaxSearchBack   int
	SearchStartDate time.Time
	FetchFn         func(time.Time, internal.CaseType) (*[]byte, error)
	ReadFn          func(*[]byte) (*readers.CaseTable, error)
}

// Basic implementation of AccUpdater
//
// Meant to be a demostration and not to be used
type basicAccUpdter struct {
	Fetch fetchers.Fetcher `json:"-"`
	Read  readers.Reader   `json:"-"`
	opts  *AccUpdterOpts
}

// Only fetches and reads, returning results
func (updter *basicAccUpdter) FindUpdates(keys []string, ids *[]string) (updatedAccords []*UpdatedAccord, notFoundIds []string, err error) {
	searchDate := time.Now()
	updatedAccords = []*UpdatedAccord{}
	notFoundIds = make([]string, len(keys))
	for _, k := range keys {
		parts := strings.Split(k, readers.CaseKeySeparator)
		notFoundIds = append(notFoundIds, parts[0])
	}

	if ids != nil && len(*ids) > 0 {
		if store := updter.getStore(); store != nil {
			cases, _ := (*store).FindAll(*ids)
			for _, c := range cases {
				notFoundIds = append(notFoundIds, c.CaseId)
			}
		}
	}

	if updter.opts.SearchStartDate != (time.Time{}) {
		searchDate, _ = time.Parse(
			"2006-01-02",
			updter.opts.SearchStartDate.Format("2006-01-02"),
		)
	}

	for i := 0; i <= updter.opts.MaxSearchBack; i++ {
		data, err := updter.Fetch(searchDate, updter.opts.CaseType)
		if err != nil {
			if errors.Is(err, fetchers.ErrDocNotFound) {
				continue
			}

			return nil, nil, fmt.Errorf(
				"Error FetchFail: with SearchStartDate: %s; CaseType %s; Region: %s\n%w",
				updter.opts.SearchStartDate.Format("2006-01-02"),
				updter.opts.CaseType,
				updter.opts.Region,
				err,
			)
		}

		caseTable, err := updter.Read(data)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"Error ReadFail: with SearchStartDate: %s; CaseType %s; Region: %s\n  %w",
				updter.opts.SearchStartDate.Format("2006-01-02"),
				updter.opts.CaseType,
				updter.opts.Region,
				err,
			)
		}

		nextSearchIds := []string{}
		for _, id := range notFoundIds {
			if caseRow := caseTable.Find(id); caseRow != nil {
				acc := UpdatedAccord{
					CaseKey:  caseRow.GetCaseKey(),
					CaseType: internal.CaseType(caseRow.CaseType),
					CaseId:   caseRow.CaseId,
					Content:  caseRow.Accord,
					Date:     searchDate,
					Nature:   caseRow.Nature,
					OthIds:   caseRow.AllIds,
				}
				updatedAccords = append(updatedAccords, &acc)
				continue
			}
			nextSearchIds = append(nextSearchIds, id)
		}

		if len(nextSearchIds) == 0 {
			break
		}

		notFoundIds = nextSearchIds
		searchDate = searchDate.Add(-24 * time.Hour)
	}

	return updatedAccords, notFoundIds, nil
}

// Not implemented. BasicAccUpdter is meant to be a demostration
//
// BasicAccUpdter should not be used in production
func (updter *basicAccUpdter) Update(keys []string, ids *[]string) (notFoundKeys []string, err error) {
	return nil, nil
}

// Not implemented. BasicAccUpdter is meant to be a demostration
//
// BasicAccUpdter should not be used in production
func (updter *basicAccUpdter) getStore() *CaseStore {
	return nil
}

type DefaultCaseStore struct {
	db *sql.DB
}

func NewDefaultCaseStore(db *sql.DB) *DefaultCaseStore {
	return &DefaultCaseStore{db}
}

func (st *DefaultCaseStore) FindAll(ids []string) ([]*db.LexCase, error) {
	return db.FindCasesById(context.TODO(), st.db, ids)
}
func (st *DefaultCaseStore) FindAllKeys(keys []string) ([]*db.LexCase, error) {
	return db.FindCases(context.TODO(), st.db, keys)
}
func (st *DefaultCaseStore) FindById(id string) (*db.LexCase, error) {
	return db.FindCaseById(context.TODO(), st.db, id)
}
func (st *DefaultCaseStore) FindByKey(key string) (*db.LexCase, error) {
	return db.FindCase(context.TODO(), st.db, key)
}

// Appends CaseRows and index entries from the mergingTable into the targetTable
// only if the targetTable does not contain an existing row for the caseId
//
// The merge occurs in place
func mergeCaseTables(targetTable, mergingTable *readers.CaseTable) {
	for _, mRow := range mergingTable.Cases {
		if tRow := targetTable.Find(mRow.CaseId); tRow == nil {
			targetTable.Add(mRow)
		}
	}
}
