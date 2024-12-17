package accupdter

import (
	"os"
	"testing"
	"time"

	"github.com/vladwithcode/lex_app/internal"
	"github.com/vladwithcode/lex_app/internal/readers"
)

func TestFindUpdates(t *testing.T) {
	expectAccords := []*UpdatedAccord{
		{
			CaseId:  "84/2003",
			Nature:  "Some sample nature",
			Content: "Some sample accord content",
			OthIds:  []string{"84/2003"},
		},
		{
			CaseId:  "264/2018",
			Nature:  "Second Nature\nin 2 lines",
			Content: "Valid content for an accord\nwith two lines of height?",
			OthIds:  []string{"264/2018"},
		},
		{
			CaseId:  "13/1998",
			Nature:  "Third Nature",
			Content: "This content is longer in width that previous contents, but should be stored correctly in the content col",
			OthIds:  []string{"13/1998"},
		},
		{
			CaseId:  "45/3000",
			Nature:  "Forth Nature non\nstandard too",
			Content: "This accord has a nature that doesn't start\non the first line",
			OthIds:  []string{"45/3000"},
		},
		{
			CaseId:  "60/1234",
			Nature:  "Fifth with multi ID",
			Content: "This accord has multiple caseIds",
			OthIds:  []string{"60/1234", "60/1234-I"},
		},
		{
			CaseId:  "1024/2048",
			Nature:  "Sixth with both\nvalid/invalid id",
			Content: "This accord has a valid and an invalid caseId",
			OthIds:  []string{"1024/2048", "eee/wrong"},
		},
	}
	expectNotFound := []string{"12/12"}

	updtr := basicAccUpdter{
		Fetch: mockFetch,
		Read:  readers.NewReader(internal.RegionDefault),
		opts: &AccUpdterOpts{
			Region:   internal.RegionDefault,
			CaseType: internal.CaseTypeAux1,
			CaseIds: []string{
				"84/2003",
				"264/2018",
				"13/1998",
				"45/3000",
				"60/1234",
				"1024/2048",
				"12/12",
			},
		},
	}
	accords, notFound, err := updtr.FindUpdates()
	if err != nil {
		t.Fatalf("errored with %v", err)
	}
	if len(notFound) != len(expectNotFound) {
		t.Fatalf("Expected notFoundIds to have %d ids, got %d", len(expectNotFound), len(notFound))
	}
	for i, expectId := range expectNotFound {
		if expectId != notFound[i] {
			t.Errorf("expected notFound[%d] to be %s, got %s", i, expectId, notFound[i])
			t.Fail()
		}
	}
	if len(accords) != len(expectAccords) {
		t.Fatalf("Expected updatedAccords to have %d accords, got %d", len(expectAccords), len(accords))
	}
	for i, expectedAccord := range expectAccords {
		actualAccord := accords[i]
		if expectedAccord.CaseId != actualAccord.CaseId {
			t.Errorf("expectedAccord[%d].CaseId is not %q, got %q", i, expectedAccord.CaseId, actualAccord.CaseId)
			t.Fail()
		}
		if expectedAccord.Content != actualAccord.Content {
			t.Errorf("expectedAccord[%d].Content is not %q, got %q", i, expectedAccord.Content, actualAccord.Content)
			t.Fail()
		}
		if expectedAccord.Nature != actualAccord.Nature {
			t.Errorf("expectedAccord[%d].Nature is not %q, got %q", i, expectedAccord.Nature, actualAccord.Nature)
			t.Fail()
		}

		if len(expectedAccord.OthIds) != len(actualAccord.OthIds) {
			t.Errorf("expectedAccord[%d].OthIds is not %q, got %q", i, expectedAccord.OthIds, actualAccord.OthIds)
			t.Fail()
		} else {
			for j, othId := range expectedAccord.OthIds {
				if othId != actualAccord.OthIds[j] {
					t.Errorf(
						"expectedAccord[%d].OthIds[%d] is not %q, got %q",
						i,
						j,
						othId,
						actualAccord.OthIds[j],
					)
					t.Fail()
				}
			}
		}
	}
}

func NewTestUpdater() *basicAccUpdter {
	return &basicAccUpdter{
		Fetch: mockFetch,
		Read:  readers.NewReader(internal.RegionDefault),
		opts: &AccUpdterOpts{
			Region: internal.RegionDefault,
			CaseIds: []string{
				"84/2003",
				"264/2018",
				"13/1998",
				"45/3000",
				"60/1234",
				"1024/2048",
				"12/12",
			},
		},
	}
}

func mockFetch(_ time.Time, _ internal.CaseType) (*[]byte, error) {
	out, _ := os.ReadFile("./test_accord_file.txt")

	return &out, nil
}
