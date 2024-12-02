package fetchers

import (
	"time"

	"github.com/vladwithcode/lex_app/internal"
)

func NewFetcher(region internal.REGION) func(time.Time, string) (*[]byte, error) {
	switch region {
	case internal.REGION_DGO:
		return DgoFetch
	default:
		return DgoFetch
	}
}
