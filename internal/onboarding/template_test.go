package onboarding

import (
	"strings"
	"testing"
)

func TestWrapUnwrapTemplateRoundTrips(t *testing.T) {
	profile := completeProfile()
	rendered, err := Render(profile)
	if err != nil {
		t.Fatalf("render profile: %v", err)
	}

	const name = "Measured plain English -- speakable, understated, general audience"
	wrapped, err := WrapTemplate(name, rendered)
	if err != nil {
		t.Fatalf("wrap template: %v", err)
	}

	gotName, gotContent, err := UnwrapTemplate(wrapped)
	if err != nil {
		t.Fatalf("unwrap template: %v", err)
	}
	if gotName != name {
		t.Fatalf("unwrapped name = %q, want %q", gotName, name)
	}
	if gotContent != rendered {
		t.Fatalf("unwrapped content does not match original Render() output:\ngot:  %q\nwant: %q", gotContent, rendered)
	}
}

func TestWrapTemplatePreservesUnusualNames(t *testing.T) {
	profile := completeProfile()
	rendered, err := Render(profile)
	if err != nil {
		t.Fatalf("render profile: %v", err)
	}

	cases := []string{
		"Simple",
		"Name: with a colon",
		"Name with \"quotes\"",
		"Em dash — and comma, here",
	}
	for _, name := range cases {
		wrapped, err := WrapTemplate(name, rendered)
		if err != nil {
			t.Fatalf("wrap template %q: %v", name, err)
		}
		gotName, gotContent, err := UnwrapTemplate(wrapped)
		if err != nil {
			t.Fatalf("unwrap template %q: %v", name, err)
		}
		if gotName != name {
			t.Fatalf("unwrapped name = %q, want %q", gotName, name)
		}
		if gotContent != rendered {
			t.Fatalf("unwrapped content mismatch for name %q", name)
		}
	}
}

func TestUnwrapTemplateRejectsMissingWrapper(t *testing.T) {
	profile := completeProfile()
	rendered, err := Render(profile)
	if err != nil {
		t.Fatalf("render profile: %v", err)
	}

	// A plain onboarding-style profile, never wrapped -- must not be
	// silently accepted or have a name guessed for it.
	_, _, err = UnwrapTemplate(rendered)
	if err == nil {
		t.Fatal("expected an error for content with no wrapper, got none")
	}
}

func TestUnwrapTemplateRejectsUnterminatedWrapper(t *testing.T) {
	_, _, err := UnwrapTemplate("---\nname: Broken\nno closing delimiter here")
	if err == nil {
		t.Fatal("expected an error for an unterminated wrapper block, got none")
	}
}

func TestWrapTemplateOutputStripsToOriginal(t *testing.T) {
	profile := completeProfile()
	rendered, err := Render(profile)
	if err != nil {
		t.Fatalf("render profile: %v", err)
	}
	wrapped, err := WrapTemplate("Generic", rendered)
	if err != nil {
		t.Fatalf("wrap template: %v", err)
	}
	if !strings.HasPrefix(wrapped, "---\nname: Generic\n---\n\n---\nkind: voice\n") {
		t.Fatalf("unexpected wrapped shape, got: %q", wrapped[:min(len(wrapped), 60)])
	}
}
