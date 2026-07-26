package draft

import (
	"strings"
	"testing"
)

func validDraft() string {
	return "# A real title\n\n*A subtitle that sets the stake.*\n\nFirst paragraph, front-loaded[^src-a].\n\n## A section\n\nBody text here.\n"
}

func TestValidateRequiresPiecePath(t *testing.T) {
	b := Bundle{DraftMD: validDraft()}
	err := Validate(b, map[string]bool{"src-a": true})
	if err == nil {
		t.Fatal("expected error when piece_path is empty")
	}
}

func TestValidateRequiresDraftMD(t *testing.T) {
	b := Bundle{PiecePath: "x"}
	if err := Validate(b, nil); err == nil {
		t.Fatal("expected error when draft_md is empty")
	}
}

func TestValidateAcceptsWellFormedDraft(t *testing.T) {
	b := Bundle{PiecePath: "x", DraftMD: validDraft()}
	if err := Validate(b, map[string]bool{"src-a": true}); err != nil {
		t.Fatalf("expected valid draft to pass, got %v", err)
	}
}

func TestValidateDoesNotEnforceHeadingStructure(t *testing.T) {
	// A letter-surface draft has no heading hierarchy at all -- Validate
	// must not impose a blog-shaped "exactly one H1, nothing past H3"
	// constraint. Zero headings, multiple "#" lines, and a deeply nested
	// heading must all be accepted; only citation integrity is mechanical.
	cases := []string{
		"Dear reader,\n\nNo heading at all, just prose, like a letter.\n",
		"# Title one\n\nbody\n\n# Title two\n\nmore body\n",
		"# Title\n\nbody\n\n#### As deep as the piece needs\n\nmore\n",
	}
	for _, draftMD := range cases {
		if err := Validate(Bundle{PiecePath: "x", DraftMD: draftMD}, nil); err != nil {
			t.Fatalf("expected heading shape to be unconstrained, got %v for:\n%s", err, draftMD)
		}
	}
}

func TestValidateRejectsRawURL(t *testing.T) {
	b := Bundle{PiecePath: "x", DraftMD: "# Title\n\nSee https://example.com for more.\n"}
	if err := Validate(b, nil); err == nil {
		t.Fatal("expected error for a raw URL in prose")
	}
}

func TestValidateAllowsURLInsideCodeFence(t *testing.T) {
	b := Bundle{PiecePath: "x", DraftMD: "# Title\n\n```text\nhttps://example.com\n```\n\nbody\n"}
	if err := Validate(b, nil); err != nil {
		t.Fatalf("expected a URL inside a code fence to be allowed, got %v", err)
	}
}

func TestValidateAllowsURLInsideInlineCode(t *testing.T) {
	b := Bundle{PiecePath: "x", DraftMD: "# Title\n\nRun `curl https://api.example.com/health` to check.\n"}
	if err := Validate(b, nil); err != nil {
		t.Fatalf("expected a URL inside inline code to be allowed, got %v", err)
	}
}

func TestValidateIgnoresCitationSyntaxShownInCodeFence(t *testing.T) {
	// A draft that illustrates citation syntax as a literal example inside a
	// code fence must not have that bracket token checked against real
	// shortnames -- the same exclusion the raw-URL check already gets.
	b := Bundle{PiecePath: "x", DraftMD: "# Title\n\nCite claims like this:\n\n```\n[^example]\n```\n\nActual claim[^real].\n"}
	if err := Validate(b, map[string]bool{"real": true}); err != nil {
		t.Fatalf("expected citation syntax shown inside a code fence to be ignored, got %v", err)
	}
}

func TestValidateRejectsUnknownCitation(t *testing.T) {
	b := Bundle{PiecePath: "x", DraftMD: "# Title\n\nA claim that needs backing[^does-not-exist].\n"}
	err := Validate(b, map[string]bool{"other-source": true})
	if err == nil {
		t.Fatal("expected error for a citation with no matching source")
	}
	if !strings.Contains(err.Error(), "does-not-exist") {
		t.Fatalf("expected error to name the bad shortname, got %v", err)
	}
}

func TestValidateRejectsCitationWhenNoResearchRan(t *testing.T) {
	// validShortnames is nil/empty when research never ran -- any citation
	// marker is correctly invalid, not silently accepted.
	b := Bundle{PiecePath: "x", DraftMD: "# Title\n\nA claim[^src-a].\n"}
	if err := Validate(b, nil); err == nil {
		t.Fatal("expected error: no research ran, so no shortname can be valid")
	}
}

func TestWordCount(t *testing.T) {
	if got := WordCount("one two three"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}
