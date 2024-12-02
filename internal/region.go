package internal

type REGION string

// Region identifiers of the form
// COUNTRY_STATE_REGION
// Where `REGION` does not match directly to a municipality
// but to a group of municipalities available in an specific
// website
const (
	REGION_DGO     REGION = "MX_DGO_DGO"
	REGION_DEFAULT REGION = REGION_DGO
)

var AllRegions = []struct {
	Value  REGION
	TSName string
}{
	{REGION_DGO, "REGION_DGO"},
	{REGION_DEFAULT, "REGION_DEFAULT"},
}
