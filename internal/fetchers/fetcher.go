package fetchers

import "time"

func NewFetcher(region REGION) func(time.Time, string) (*[]byte, error) {
	switch region {
	case REGION_DGO:
		return DgoFetch
	default:
		return DgoFetch
	}
}
