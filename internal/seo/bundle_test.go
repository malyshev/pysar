package seo

import (
	"strings"
	"testing"
)

func validChecklist() Bundle {
	return Bundle{
		PiecePath:       "x",
		RevisedMD:       "# A concrete title\n\n*A subtitle that sets the stake.*\n\nFirst paragraph, citing [the retry budget](https://example.com/retry-budget).\n",
		Checks:          []string{"[citation] resolved [^src-a] to an inline link"},
		Tags:            []string{"reliability"},
		MetaTitle:       "A concrete title for search",
		MetaDescription: "A concrete, specific description written for a reader who hasn't opened the piece yet.",
		URLSlug:         "a-concrete-title",
	}
}

var validURLs = map[string]bool{"https://example.com/retry-budget": true}

func TestValidateRequiresAtLeastOneCheck(t *testing.T) {
	b := validChecklist()
	b.Checks = nil
	if err := Validate(b, validURLs); err == nil {
		t.Fatal("expected error when checks is empty")
	}
}

func TestValidateRejectsUnresolvedCitationMarker(t *testing.T) {
	b := validChecklist()
	b.RevisedMD = "# Title\n\nA claim that still cites [^src-a] unresolved.\n"
	err := Validate(b, validURLs)
	if err == nil {
		t.Fatal("expected error for an unresolved [^shortname] marker")
	}
	if !strings.Contains(err.Error(), "src-a") {
		t.Fatalf("expected error to name the unresolved marker, got %v", err)
	}
}

func TestValidateRejectsFabricatedLink(t *testing.T) {
	b := validChecklist()
	b.RevisedMD = "# Title\n\nA claim with an [invented link](https://not-a-real-source.example.com).\n"
	err := Validate(b, validURLs)
	if err == nil {
		t.Fatal("expected error for a link that doesn't match any recorded source")
	}
	if !strings.Contains(err.Error(), "not-a-real-source") {
		t.Fatalf("expected error to name the bad URL, got %v", err)
	}
}

func TestValidateRejectsBareURLOutsideLink(t *testing.T) {
	b := validChecklist()
	b.RevisedMD = "# Title\n\nSee https://example.com/retry-budget directly.\n"
	if err := Validate(b, validURLs); err == nil {
		t.Fatal("expected error for a bare URL not wrapped in a markdown link")
	}
}

func TestValidateAcceptsResolvedLinkMatchingSource(t *testing.T) {
	b := validChecklist()
	if err := Validate(b, validURLs); err != nil {
		t.Fatalf("expected a well-formed resolved revision to pass, got %v", err)
	}
}

func TestValidateRequiresTags(t *testing.T) {
	b := validChecklist()
	b.Tags = nil
	if err := Validate(b, validURLs); err == nil {
		t.Fatal("expected error when tags is empty")
	}
}

func TestValidateRequiresMetaFields(t *testing.T) {
	b := validChecklist()
	b.MetaTitle = ""
	b.MetaDescription = ""
	err := Validate(b, validURLs)
	if err == nil {
		t.Fatal("expected error when meta_title/meta_description are empty")
	}
	if !strings.Contains(err.Error(), "meta_title") || !strings.Contains(err.Error(), "meta_description") {
		t.Fatalf("expected error to name both missing meta fields, got %v", err)
	}
}

func TestValidateRejectsNonKebabSlug(t *testing.T) {
	b := validChecklist()
	b.URLSlug = "Not A Slug!"
	err := Validate(b, validURLs)
	if err == nil {
		t.Fatal("expected error for a non-kebab-case url_slug")
	}
	if !strings.Contains(err.Error(), "url_slug") {
		t.Fatalf("expected error to name url_slug, got %v", err)
	}
}

func TestValidatePreservesCodeFenceExamples(t *testing.T) {
	// A code fence showing citation/link syntax as a literal example must
	// not be checked against real content rules, same as internal/draft's
	// own stripCode reasoning.
	b := validChecklist()
	b.RevisedMD += "\n```\ncite like this: [^example]\n```\n"
	if err := Validate(b, validURLs); err != nil {
		t.Fatalf("expected code-fenced example syntax to be ignored, got %v", err)
	}
}
