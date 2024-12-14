package fetchers

import (
	"time"

	"github.com/vladwithcode/lex_app/internal"
)

type Fetcher func(time.Time, internal.CaseType) (*[]byte, error)

func NewFetcher(region internal.Region) Fetcher {
	switch region {
	case internal.RegionDgo:
		return DgoFetch
	default:
		return DgoFetch
	}
}
