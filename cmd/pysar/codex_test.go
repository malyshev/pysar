package main

import (
	"strings"
	"testing"
)

func TestTransformCodexSkillRewritesInvocations(t *testing.T) {
	in := []byte("Run /ps-intake then /ps. Prefer slash commands here.\n")
	got := string(transformCodexSkill(in))
	if strings.Contains(got, "/ps") {
		t.Fatalf("left /ps invocations unrewritten: %q", got)
	}
	if !strings.Contains(got, "$ps-intake") || !strings.Contains(got, "$ps.") {
		t.Fatalf("missing $ps rewrites: %q", got)
	}
	if !strings.Contains(got, "explicit skill invocations") {
		t.Fatalf("missing slash-command wording rewrite: %q", got)
	}
}

func TestCodexAllowImplicitOnlyOrchestrator(t *testing.T) {
	if !codexAllowImplicit("ps") {
		t.Fatal("ps must allow implicit invocation")
	}
	for _, name := range []string{"ps-intake", "ps-draft", "ps-onboard", "ps-voice"} {
		if codexAllowImplicit(name) {
			t.Fatalf("%s must be explicit-only", name)
		}
	}
}

func TestCodexPolicyYAML(t *testing.T) {
	if got := string(codexPolicyYAML("ps")); got != "policy:\n  allow_implicit_invocation: true\n" {
		t.Fatalf("ps policy = %q", got)
	}
	if got := string(codexPolicyYAML("ps-intake")); got != "policy:\n  allow_implicit_invocation: false\n" {
		t.Fatalf("ps-intake policy = %q", got)
	}
}
