package accupdter

import (
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
	CaseType internal.CaseType
	CaseId   string
	Content  string
	Date     time.Time
	Nature   string
	OthIds   []string
}

type BasicAccUpdter struct {
	fetch fetchers.Fetcher
	read  readers.Reader
	opts  *AccUpdterOpts
}

func NewBasicUpdter(opts *AccUpdterOpts) *BasicAccUpdter {
	return &BasicAccUpdter{
		fetch: fetchers.NewFetcher(opts.Region),
		read:  readers.NewReader(opts.Region),
		opts:  opts,
	}
}

// Doesn't return but updates DB with accords found
func (updter *BasicAccUpdter) Update() {
}
