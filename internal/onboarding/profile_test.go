package onboarding

import (
	"errors"
	"strings"
	"testing"

	"pysar/internal/editorial"

	"gopkg.in/yaml.v3"
)

func completeProfile() Profile {
	return Profile{
		Kind:           KindVoice,
		Tone:           "warm, direct, a little dry",
		Formality:      "neutral",
		SentenceLength: "varied, mostly short",
		Register:       "international standard English",
		Goldens: []GoldenExample{
			{Label: "opening", Text: "Bring your take. Pysar helps you shape it."},
			{Label: "explanation", Text: "The gap is going from idea to shaped piece, not tooling."},
			{Label: "closing", Text: "You decide when it's ready. We never post for you."},
		},
	}
}

func TestValidateAcceptsCompleteProfile(t *testing.T) {
	if err := Validate(completeProfile()); err != nil {
		t.Fatalf("expected complete profile to validate, got %v", err)
	}
}

func TestValidateRejectsMissingFields(t *testing.T) {
	p := completeProfile()
	p.Tone = ""
	p.Register = ""

	err := Validate(p)
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %v", err)
	}
	if len(ve.Missing) != 2 {
		t.Fatalf("expected 2 missing fields, got %v", ve.Missing)
	}
}

func TestValidateRejectsTooFewGoldens(t *testing.T) {
	p := completeProfile()
	p.Goldens = p.Goldens[:1]

	err := Validate(p)
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %v", err)
	}
	if len(ve.Missing) != 1 {
		t.Fatalf("expected exactly 1 missing entry (golden count), got %v", ve.Missing)
	}
}

func TestValidateRejectsEmptyGoldenText(t *testing.T) {
	p := completeProfile()
	p.Goldens[1].Text = ""

	err := Validate(p)
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %v", err)
	}
}

func TestRenderProducesParseableFrontmatterAndGoldenSections(t *testing.T) {
	p := completeProfile()
	p.Notes = `contains: a colon, "quotes", and — an em dash`
	p.BannedPhrases = []string{"in today's fast-paced world"}

	rendered, err := Render(p)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	parts := strings.SplitN(rendered, "---\n", 3)
	if len(parts) != 3 {
		t.Fatalf("expected exactly one frontmatter block delimited by ---, got: %q", rendered)
	}
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		t.Fatalf("frontmatter did not parse as valid YAML: %v\ncontent: %s", err, parts[1])
	}
	if fm.Tone != p.Tone || fm.Notes != p.Notes {
		t.Fatalf("round-tripped frontmatter mismatch: got %+v", fm)
	}

	for _, g := range p.Goldens {
		if !strings.Contains(rendered, "## Golden: "+g.Label) || !strings.Contains(rendered, g.Text) {
			t.Fatalf("expected golden %q to appear as its own section, got: %s", g.Label, rendered)
		}
	}
}

func TestVoicePassPreconditionDeclinesWhenAlreadyProduced(t *testing.T) {
	s := editorial.NewState(editorial.SurfaceBlog, editorial.ArtifactVoiceProfile)
	if _, err := editorial.Run(voicePass{}, s); err == nil {
		t.Fatal("expected voicePass to decline when voice profile already produced")
	}
}

func TestVoicePassRunProducesArtifactOnEmptyState(t *testing.T) {
	s := editorial.NewState(editorial.SurfaceBlog)
	next, err := editorial.Run(voicePass{}, s)
	if err != nil {
		t.Fatalf("expected voicePass to run on empty state, got %v", err)
	}
	if !next.Has(editorial.ArtifactVoiceProfile) {
		t.Fatal("expected resulting state to have ArtifactVoiceProfile")
	}
}

func TestVoicePassIsRegistered(t *testing.T) {
	for _, p := range editorial.Registered() {
		if p.Name() == "ps-voice" {
			return
		}
	}
	t.Fatal("expected ps-voice to be registered in the editorial Pass registry")
}

func completeStyleProfile() Profile {
	return Profile{
		Kind: KindStyle,
		Rules: []string{
			"Put the main point first",
			"Prefer active voice",
			"Cut what doesn't earn its place",
		},
		Goldens: []GoldenExample{
			{Label: "opening", Text: "The deadline moved to Friday. Ship what's ready; hold the rest."},
			{Label: "structure", Text: "Three things changed. First, ... Second, ... Third, ..."},
			{Label: "cutting", Text: "We tested it. It broke. We fixed it."},
		},
	}
}

func TestValidateAcceptsCompleteStyleProfileWithoutVoiceFields(t *testing.T) {
	if err := Validate(completeStyleProfile()); err != nil {
		t.Fatalf("expected a complete style profile (rules + goldens, no tone/formality/etc) to validate, got %v", err)
	}
}

func TestValidateRejectsStyleProfileWithTooFewRules(t *testing.T) {
	p := completeStyleProfile()
	p.Rules = p.Rules[:1]

	err := Validate(p)
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %v", err)
	}
	found := false
	for _, m := range ve.Missing {
		if strings.HasPrefix(m, "rules") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected missing entries to name 'rules', got %v", ve.Missing)
	}
}

func TestValidateStillRequiresVoiceFieldsForVoiceKind(t *testing.T) {
	// A voice-kind profile with plenty of Rules but no tone/formality/etc
	// must still fail -- Rules doesn't substitute for voice's own fields.
	p := completeProfile()
	p.Tone = ""
	p.Rules = completeStyleProfile().Rules

	err := Validate(p)
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %v", err)
	}
}

func TestRenderStyleProfileOmitsVoiceFieldsAndShowsRulesHeading(t *testing.T) {
	p := completeStyleProfile()
	rendered, err := Render(p)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if strings.Contains(rendered, "tone:") || strings.Contains(rendered, "formality:") {
		t.Fatalf("expected an unpopulated voice field to be omitted from a style profile's frontmatter, got: %q", rendered)
	}
	if !strings.Contains(rendered, "# Style Profile") {
		t.Fatalf("expected the H1 heading to say 'Style Profile', got: %q", rendered[:min(len(rendered), 80)])
	}
	if !strings.Contains(rendered, "## Rules") {
		t.Fatalf("expected a Rules section, got: %q", rendered)
	}
	for _, r := range p.Rules {
		if !strings.Contains(rendered, "- "+r) {
			t.Fatalf("expected rule %q rendered as a list item, got: %q", r, rendered)
		}
	}
}

func TestRenderVoiceProfileStillSaysVoiceProfile(t *testing.T) {
	rendered, err := Render(completeProfile())
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(rendered, "# Voice Profile") {
		t.Fatalf("expected the H1 heading to remain 'Voice Profile' for Kind=voice, got: %q", rendered[:min(len(rendered), 80)])
	}
}

func TestTemplatesDirIsHostAgnosticAndKindScoped(t *testing.T) {
	got := TemplatesDir("/home/author", KindVoice)
	want := "/home/author/.pysar/templates/voice"
	if got != want {
		t.Fatalf("TemplatesDir(voice) = %q, want %q", got, want)
	}
	if style := TemplatesDir("/home/author", KindStyle); style == got {
		t.Fatalf("expected voice and style template dirs to differ, both got %q", style)
	}
}

func TestSlugNormalizesArbitraryNames(t *testing.T) {
	cases := map[string]string{
		"Measured plain English -- speakable, understated, general audience": "measured-plain-english-speakable-understated-general-audience",
		"  leading and trailing spaces  ":                                    "leading-and-trailing-spaces",
		"ALLCAPS":                                                            "allcaps",
		"already-slugged":                                                    "already-slugged",
		"":                                                                   "template",
		"!!!":                                                                "template",
	}
	for input, want := range cases {
		if got := Slug(input); got != want {
			t.Fatalf("Slug(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestSlugNeverProducesLeadingOrTrailingDash(t *testing.T) {
	got := Slug("-- wrapped --")
	if strings.HasPrefix(got, "-") || strings.HasSuffix(got, "-") {
		t.Fatalf("Slug produced a leading/trailing dash: %q", got)
	}
}
