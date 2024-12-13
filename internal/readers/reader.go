package readers

import (
	"strconv"

	"github.com/vladwithcode/lex_app/internal"
)

// Returns a reader func that takes an pointer to an byte
// slice and creates a CaseTable
func NewReader(region internal.Region) func(*[]byte) (*CaseTable, error) {
	switch region {
	default:
		return dgoReader
	case internal.RegionDgo:
		return dgoReader
	}
}

func isValidNumericStr(candidate string) bool {
	_, err := strconv.Atoi(candidate)

	return err == nil
}
