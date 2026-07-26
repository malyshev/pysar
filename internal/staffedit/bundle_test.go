package staffedit

import (
	"strings"
	"testing"
)

func validRevisedDraft() string {
	return "# A tighter title\n\n*A subtitle that sets the stake.*\n\nFirst paragraph, front-loaded[^src-a].\n\n## A section\n\nRevised body text here.\n"
}

func TestValidateRequiresAtLeastOneCheck(t *testing.T) {
	b := Bundle{PiecePath: "x", RevisedMD: validRevisedDraft()}
	err := Validate(b, map[string]bool{"src-a": true})
	if err == nil {
		t.Fatal("expected error when checks is empty")
	}
	if !strings.Contains(err.Error(), "checks") {
		t.Fatalf("expected missing list to name checks, got %v", err)
	}
}

func TestValidateRejectsEmptyCheckEntry(t *testing.T) {
	b := Bundle{PiecePath: "x", RevisedMD: validRevisedDraft(), Checks: []string{"   "}}
	if err := Validate(b, map[string]bool{"src-a": true}); err == nil {
		t.Fatal("expected error for a blank check entry")
	}
}

func TestValidateAcceptsNoChangesNeededAsOneEntry(t *testing.T) {
	b := Bundle{PiecePath: "x", RevisedMD: validRevisedDraft(), Checks: []string{"no changes needed -- all checks already satisfied"}}
	if err := Validate(b, map[string]bool{"src-a": true}); err != nil {
		t.Fatalf("expected a single 'no changes needed' entry to be sufficient, got %v", err)
	}
}

func TestValidateReusesDraftCitationIntegrity(t *testing.T) {
	// staffedit.Validate must not re-implement citation checking -- proven
	// by confirming the exact same failure draft.Validate would catch
	// (an unmatched shortname) is caught here too.
	b := Bundle{
		PiecePath: "x",
		RevisedMD: "# Title\n\nA claim that cites a source that was never fetched[^ghost-source].\n",
		Checks:    []string{"[stakes] tightened the opener"},
	}
	err := Validate(b, map[string]bool{"other-source": true})
	if err == nil {
		t.Fatal("expected error for an unmatched citation")
	}
	if !strings.Contains(err.Error(), "ghost-source") {
		t.Fatalf("expected error to name the bad shortname, got %v", err)
	}
}

func TestValidateErrorNamesRevisedMDNotDraftMD(t *testing.T) {
	// save_staff_edit_bundle's schema has revised_md, not draft_md -- an
	// error telling the agent to fix "draft_md" would point at a field
	// that doesn't exist in the tool it's calling.
	b := Bundle{
		PiecePath: "x",
		RevisedMD: "# Title\n\nSee https://example.com directly.\n",
		Checks:    []string{"[stakes] tightened the opener"},
	}
	err := Validate(b, nil)
	if err == nil {
		t.Fatal("expected error for a raw URL")
	}
	if !strings.Contains(err.Error(), "revised_md") {
		t.Fatalf("expected error to name revised_md, got %v", err)
	}
	if strings.Contains(err.Error(), "draft_md") {
		t.Fatalf("error must not reference draft_md -- that field doesn't exist in save_staff_edit_bundle's schema, got %v", err)
	}
}

func TestValidateReusesDraftRawURLCheck(t *testing.T) {
	b := Bundle{
		PiecePath: "x",
		RevisedMD: "# Title\n\nSee https://example.com directly.\n",
		Checks:    []string{"[readability] smoothed the transition"},
	}
	if err := Validate(b, nil); err == nil {
		t.Fatal("expected error for a raw URL in the revised prose")
	}
}

func TestValidateAcceptsWellFormedRevision(t *testing.T) {
	b := Bundle{
		PiecePath: "x",
		RevisedMD: validRevisedDraft(),
		Checks:    []string{"[readability] split a long sentence into two", "[stakes] named the specific failure mode in paragraph 2"},
	}
	if err := Validate(b, map[string]bool{"src-a": true}); err != nil {
		t.Fatalf("expected a well-formed revision to pass, got %v", err)
	}
}
