package fetchers

import (
	"time"

	"github.com/vladwithcode/lex_app/internal"
)

func NewFetcher(region internal.Region) func(time.Time, string) (*[]byte, error) {
	switch region {
	case internal.RegionDgo:
		return DgoFetch
	default:
		return DgoFetch
	}
}
