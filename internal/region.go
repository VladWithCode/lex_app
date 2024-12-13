package internal

type Region string

// Region identifiers of the form
// COUNTRY_STATE_REGION
// Where `REGION` does not match directly to a municipality
// but to a group of municipalities available in an specific
// website
const (
	RegionDgo     Region = "MX_DGO_DGO"
	RegionDefault Region = RegionDgo
)

var AllRegions = []struct {
	Value  Region
	TSName string
}{
	{RegionDgo, "REGION_DGO"},
	{RegionDefault, "REGION_DEFAULT"},
}

type CaseType string

const (
	CaseTypeAux1  CaseType = "aux1"
	CaseTypeAux2  CaseType = "aux2"
	CaseTypeCiv2  CaseType = "civ2"
	CaseTypeCiv3  CaseType = "civ3"
	CaseTypeCiv4  CaseType = "civ4"
	CaseTypeFam1  CaseType = "fam1"
	CaseTypeFam2  CaseType = "fam2"
	CaseTypeFam3  CaseType = "fam3"
	CaseTypeFam4  CaseType = "fam4"
	CaseTypeFam5  CaseType = "fam5"
	CaseTypeMer1  CaseType = "mer1"
	CaseTypeMer2  CaseType = "mer2"
	CaseTypeMer3  CaseType = "mer3"
	CaseTypeMer4  CaseType = "mer4"
	CaseTypeSECCC CaseType = "seccc"
	CaseTypeSECCU CaseType = "seccu"
	CaseTypeCJMF  CaseType = "cjmf"
	CaseTypeCJMF2 CaseType = "cjmf2"
	CaseTypeTRIBL CaseType = "tribl"
)

var AllCaseTypes = []struct {
	Value  CaseType
	TSName string
}{
	{CaseTypeAux1, "aux1"},
	{CaseTypeAux2, "aux2"},
	{CaseTypeCiv2, "civ2"},
	{CaseTypeCiv3, "civ3"},
	{CaseTypeCiv4, "civ4"},
	{CaseTypeFam1, "fam1"},
	{CaseTypeFam2, "fam2"},
	{CaseTypeFam3, "fam3"},
	{CaseTypeFam4, "fam4"},
	{CaseTypeFam5, "fam5"},
	{CaseTypeMer1, "mer1"},
	{CaseTypeMer2, "mer2"},
	{CaseTypeMer3, "mer3"},
	{CaseTypeMer4, "mer4"},
	{CaseTypeSECCC, "seccc"},
	{CaseTypeSECCU, "seccu"},
	{CaseTypeCJMF, "cjmf"},
	{CaseTypeCJMF2, "cjmf2"},
	{CaseTypeTRIBL, "tribl"},
}
