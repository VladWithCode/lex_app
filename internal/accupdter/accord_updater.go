package accupdter

import (
	"errors"
	"fmt"
	"time"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/fetchers"
	"github.com/vladwithcode/lex_app/internal/readers"
)

type AccUpdter interface {
	Update()
}

type AccUpdterOpts struct {
	CaseType        internal.CaseType
	CaseIds         []string
	Region          internal.Region
	MaxSearchBack   int
	SearchStartDate time.Time
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

type BasicAccUpdter struct {
	Fetch fetchers.Fetcher `json:"-"`
	Read  readers.Reader   `json:"-"`
	opts  *AccUpdterOpts
}

func NewBasicUpdter(opts *AccUpdterOpts) *BasicAccUpdter {
	return &BasicAccUpdter{
		Fetch: fetchers.NewFetcher(opts.Region),
		Read:  readers.NewReader(opts.Region),
		opts:  opts,
	}
}

// Only fetches and reads, returning results
func (updter *BasicAccUpdter) FindUpdates() (updatedAccords []*UpdatedAccord, notFoundIds []string, err error) {
	searchDate := time.Now()
	updatedAccords = []*UpdatedAccord{}
	notFoundIds = make([]string, len(updter.opts.CaseIds))
	copy(notFoundIds, updter.opts.CaseIds)

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

// Doesn't return but updates DB with accords found
func (updter *BasicAccUpdter) Update() {
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
