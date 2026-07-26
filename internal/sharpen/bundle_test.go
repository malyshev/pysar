package sharpen

import (
	"strings"
	"testing"
)

func validRevisedDraft() string {
	return "# A sharper title\n\n*A subtitle that sets the stake.*\n\nFirst paragraph, front-loaded[^src-a].\n\n## A section\n\nRevised body text here.\n"
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
	// sharpen.Validate must not re-implement citation checking -- proven by
	// confirming the exact same failure draft.Validate would catch (an
	// unmatched shortname) is caught here too.
	b := Bundle{
		PiecePath: "x",
		RevisedMD: "# Title\n\nA claim that cites a source that was never fetched[^ghost-source].\n",
		Checks:    []string{"[opener] tightened the hook"},
	}
	err := Validate(b, map[string]bool{"other-source": true})
	if err == nil {
		t.Fatal("expected error for an unmatched citation")
	}
	if !strings.Contains(err.Error(), "ghost-source") {
		t.Fatalf("expected error to name the bad shortname, got %v", err)
	}
}

func TestValidateReusesDraftRawURLCheck(t *testing.T) {
	b := Bundle{
		PiecePath: "x",
		RevisedMD: "# Title\n\nSee https://example.com directly.\n",
		Checks:    []string{"[elevate] promoted a buried finding"},
	}
	if err := Validate(b, nil); err == nil {
		t.Fatal("expected error for a raw URL in the revised prose")
	}
}

func TestValidateAcceptsWellFormedRevision(t *testing.T) {
	b := Bundle{
		PiecePath: "x",
		RevisedMD: validRevisedDraft(),
		Checks:    []string{"[opener] tightened the hook", "[arc] closed the loop on the opening promise"},
	}
	if err := Validate(b, map[string]bool{"src-a": true}); err != nil {
		t.Fatalf("expected a well-formed revision to pass, got %v", err)
	}
}
