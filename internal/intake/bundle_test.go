package intake

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDegenerate(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"  ", true},
		{"ai", true},
		{"xyz", true},
		{"asdfasdfasdfasdf", true},
		{"write something about AI security for parents", false},
	}
	for _, tc := range cases {
		if got := Degenerate(tc.in); got != tc.want {
			t.Fatalf("Degenerate(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestAllocateUniqueNameCapsLengthWithoutShorteningTheSuffix(t *testing.T) {
	root := t.TempDir()
	longIdea := "Docker for Developers Who Just Need It to Work: Containerize Your First App Without Cargo-Culting a Twelve Layer Configuration You Copied From a Stack Overflow Answer That Was Solving a Completely Different Problem"
	name, err := AllocateUniqueName(root, longIdea)
	if err != nil {
		t.Fatal(err)
	}
	if len(name) > MaxPieceNameLength {
		t.Fatalf("name length %d exceeds MaxPieceNameLength %d: %s", len(name), MaxPieceNameLength, name)
	}
	if strings.HasSuffix(name, "-") {
		t.Fatalf("truncation left a dangling trailing hyphen before the suffix: %s", name)
	}
	// The last 12 characters (the suffix) must be untouched hex -- only the
	// readable prefix is ever shortened to make room.
	suffix := name[len(name)-suffixHexLen:]
	if _, err := hex.DecodeString(suffix); err != nil {
		t.Fatalf("expected the trailing %d chars to be an intact hex suffix, got %q", suffixHexLen, suffix)
	}
}

func TestTruncatePrefixCutsAtWordBoundaryNotMidWord(t *testing.T) {
	// Budget (MaxPieceNameLength - 1 - suffixHexLen) does not land on a
	// hyphen here -- a naive character cut would slice "containerize" into
	// a fragment. truncatePrefix must back up to the last whole word.
	long := "docker-for-developers-who-just-need-it-to-work-containerize-your-first-app"
	got := truncatePrefix(long)
	budget := MaxPieceNameLength - 1 - suffixHexLen
	if len(got) > budget {
		t.Fatalf("truncated prefix %q (%d chars) exceeds budget %d", got, len(got), budget)
	}
	if strings.HasSuffix(got, "-") {
		t.Fatalf("truncated prefix has a dangling trailing hyphen: %q", got)
	}
	for _, frag := range []string{"-co", "-cont", "-conta", "-contai", "-contain", "-containe"} {
		if strings.HasSuffix(got, frag) {
			t.Fatalf("truncated prefix ends on a fragment of 'containerize', not a whole word: %q", got)
		}
	}
	if !strings.HasPrefix(long, got) {
		t.Fatalf("truncated prefix %q is not a prefix of the original %q", got, long)
	}
}

func TestValidateRejectsBareKillerSectionsAndCounterintuitive(t *testing.T) {
	base := func() Bundle {
		return Bundle{
			Idea:        "write something about AI safety for parents",
			EntryMode:   EntryIdea,
			POVSource:   POVAuthorPractitioner,
			Restatement: "x", Audience: "x", Register: "x", ExpertLens: "x", TopicScope: "x", Thesis: "x", Promise: "x",
			OutlineMD: "## x\n", AnglesMD: "## Contrarian\n- x\n",
		}
	}

	b := base()
	b.KillerSections = []KillerSection{{Title: "A title with no edge or example"}}
	b.Counterintuitive = []CounterintuitiveFinding{{Claim: "A claim with no contradiction"}}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected Validate to reject killer_sections/counterintuitive missing edge/example/contradiction")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	joined := strings.Join(ve.Missing, " ")
	for _, want := range []string{"killer_sections[0].edge", "killer_sections[0].example", "counterintuitive[0].contradiction"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected missing list to name %q, got %v", want, ve.Missing)
		}
	}
}

func TestValidateRequiresKeyQuestionsAndNonGoals(t *testing.T) {
	b := Bundle{
		Idea:        "write something about AI safety for parents",
		EntryMode:   EntryIdea,
		POVSource:   POVAuthorPractitioner,
		Restatement: "x", Audience: "x", Register: "x", ExpertLens: "x", TopicScope: "x", Thesis: "x", Promise: "x",
		OutlineMD: "## x\n", AnglesMD: "## Contrarian\n- x\n",
		KillerSections:   []KillerSection{{Title: "t", Edge: "e", Example: "ex"}},
		Counterintuitive: []CounterintuitiveFinding{{Claim: "c", Contradiction: "d"}},
		// KeyQuestions and NonGoals deliberately left empty.
	}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected Validate to reject an empty key_questions/non_goals -- they are required, not optional")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	joined := strings.Join(ve.Missing, " ")
	for _, want := range []string{"key_questions", "non_goals"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected missing list to name %q, got %v", want, ve.Missing)
		}
	}
}

func TestSourcesDefaultToNoneUnlessActuallyFetched(t *testing.T) {
	root := t.TempDir()
	b := Bundle{
		Idea: "write something about AI safety for parents", EntryMode: EntryIdea, POVSource: POVAuthorPractitioner,
		Restatement: "x", Audience: "x", Register: "x", ExpertLens: "x", TopicScope: "x", Thesis: "x", Promise: "x",
		OutlineMD: "## x\n", AnglesMD: "## Contrarian\n- x\n",
		KillerSections:   []KillerSection{{Title: "t", Edge: "e", Example: "ex"}},
		Counterintuitive: []CounterintuitiveFinding{{Claim: "c", Contradiction: "d"}},
		KeyQuestions:     []string{"q"}, NonGoals: []string{"n"},
		// Sources deliberately left empty -- the common, legitimate case.
	}
	dir, err := WriteBundle(root, b)
	if err != nil {
		t.Fatalf("WriteBundle: %v", err)
	}
	brief, _ := os.ReadFile(filepath.Join(dir, "brief.md"))
	if !strings.Contains(string(brief), "research_mode: none") {
		t.Fatalf("expected research_mode: none with no sources:\n%s", brief)
	}
	sources, _ := os.ReadFile(filepath.Join(dir, "sources.md"))
	if !strings.Contains(string(sources), "Research mode: none") {
		t.Fatalf("expected sources.md stub for the no-sources case:\n%s", sources)
	}
}

func TestSourcesTargetedWhenActuallyFetched(t *testing.T) {
	root := t.TempDir()
	b := Bundle{
		Idea: "write something about AI safety for parents", EntryMode: EntryIdea, POVSource: POVAuthorPractitioner,
		Restatement: "x", Audience: "x", Register: "x", ExpertLens: "x", TopicScope: "x", Thesis: "x", Promise: "x",
		OutlineMD: "## x\n", AnglesMD: "## Contrarian\n- x\n",
		KillerSections:   []KillerSection{{Title: "t", Edge: "e", Example: "ex"}},
		Counterintuitive: []CounterintuitiveFinding{{Claim: "c", Contradiction: "d"}},
		KeyQuestions:     []string{"q"}, NonGoals: []string{"n"},
		Sources: []Source{{URL: "https://example.com/docs", Note: "confirmed the default timeout value"}},
	}
	dir, err := WriteBundle(root, b)
	if err != nil {
		t.Fatalf("WriteBundle: %v", err)
	}
	brief, _ := os.ReadFile(filepath.Join(dir, "brief.md"))
	if !strings.Contains(string(brief), "research_mode: targeted") {
		t.Fatalf("expected research_mode: targeted when sources were fetched:\n%s", brief)
	}
	sources, _ := os.ReadFile(filepath.Join(dir, "sources.md"))
	if !strings.Contains(string(sources), "https://example.com/docs") || !strings.Contains(string(sources), "confirmed the default timeout value") {
		t.Fatalf("expected the real fetched source in sources.md:\n%s", sources)
	}
}

func TestValidateRejectsSourceMissingURLOrNote(t *testing.T) {
	b := Bundle{
		Idea: "write something about AI safety for parents", EntryMode: EntryIdea, POVSource: POVAuthorPractitioner,
		Restatement: "x", Audience: "x", Register: "x", ExpertLens: "x", TopicScope: "x", Thesis: "x", Promise: "x",
		OutlineMD: "## x\n", AnglesMD: "## Contrarian\n- x\n",
		KillerSections:   []KillerSection{{Title: "t", Edge: "e", Example: "ex"}},
		Counterintuitive: []CounterintuitiveFinding{{Claim: "c", Contradiction: "d"}},
		KeyQuestions:     []string{"q"}, NonGoals: []string{"n"},
		Sources: []Source{{URL: "", Note: ""}},
	}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected Validate to reject a source with no URL/note -- it wasn't actually fetched")
	}
	ve := err.(*ValidationError)
	joined := strings.Join(ve.Missing, " ")
	for _, want := range []string{"sources[0].url", "sources[0].note"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected missing list to name %q, got %v", want, ve.Missing)
		}
	}
}

func TestWriteBundleProvisional(t *testing.T) {
	root := t.TempDir()
	b := Bundle{
		Name:        "ai-for-parents",
		Idea:        "write something about AI safety for parents",
		EntryMode:   EntryIdea,
		POVSource:   POVAuthorPractitioner,
		Restatement: "For parents, in plain language, focusing on practical AI safety at home",
		Audience:    "parents",
		Register:    "plain language",
		ExpertLens:  "child-safety / parenting practitioner",
		TopicScope:  "practical AI safety at home",
		Thesis:      "Parents need concrete defaults, not model-name trivia",
		Promise:     "Leave with three household rules you can apply tonight",
		KillerSections: []KillerSection{
			{Title: "The one setting that matters more than the brand", Edge: "Competitors compare brands; this compares actual default-safety settings across any brand", Example: "Screenshots of the exact router-level DNS filter setting on a $50 router"},
		},
		Counterintuitive: []CounterintuitiveFinding{
			{Claim: "More filters can make kids less safe", Contradiction: "hiding the conversation behind a filter removes the one thing that actually protects them: talking to a parent when something goes wrong"},
		},
		KeyQuestions: []string{"What should a parent do in the first week?"},
		NonGoals:     []string{"Model benchmarks"},
		OutlineMD:    "## Why this matters\n### A concrete scene\n## Three household rules\n",
		AnglesMD:     "## Misconceptions\n- Kids will not find workarounds\n\n## Contrarian / under-discussed\n- Defaults beat lectures\n\n## Trade-offs\n- Friction vs trust\n\n## Edge cases\n- Shared family accounts\n",
	}
	dir, err := WriteBundle(root, b)
	if err != nil {
		t.Fatalf("WriteBundle: %v", err)
	}
	if filepath.Base(dir) == "ai-for-parents" {
		t.Fatalf("expected a random suffix appended to the directory name, got bare %q", dir)
	}
	if !strings.HasPrefix(filepath.Base(dir), "ai-for-parents-") {
		t.Fatalf("expected dir to start with the readable prefix, got %q", filepath.Base(dir))
	}
	for _, name := range []string{"STORAGE.md", "brief.md", "outline.md", "angles.md", "sources.md", "intake-changelog.md"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	brief, _ := os.ReadFile(filepath.Join(dir, "brief.md"))
	if !strings.Contains(string(brief), "storage_status: provisional") {
		t.Fatalf("brief missing provisional marker:\n%s", brief)
	}
	if !strings.Contains(string(brief), "research_mode: none") {
		t.Fatalf("brief missing research_mode none")
	}
	if !strings.Contains(string(brief), `expert_lens: "child-safety / parenting practitioner"`) {
		t.Fatalf("brief missing expert_lens frontmatter:\n%s", brief)
	}
	if strings.Contains(string(brief), "slug:") {
		t.Fatalf("brief.md must not use the word 'slug':\n%s", brief)
	}
	for _, want := range []string{"Edge: Competitors compare brands", "Example: Screenshots of the exact router", "Contradiction: hiding the conversation"} {
		if !strings.Contains(string(brief), want) {
			t.Fatalf("brief missing killer-section/counterintuitive structure (%q):\n%s", want, brief)
		}
	}

	// Re-running with the identical Bundle (same idea, same requested name)
	// must never be blocked or error -- it silently lands in a second,
	// differently-suffixed directory. Directory collisions are storage
	// mechanics, never a decision surfaced to the author.
	dir2, err := WriteBundle(root, b)
	if err != nil {
		t.Fatalf("expected second WriteBundle with the same Bundle to succeed, got: %v", err)
	}
	if dir2 == dir {
		t.Fatalf("expected a different directory on re-run, got the same: %s", dir)
	}
}

func TestLoadAuthorDefaultsFallback(t *testing.T) {
	d, err := LoadAuthorDefaults(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if d.Source != "generalist-fallback" {
		t.Fatalf("source=%q", d.Source)
	}
}

func writeProfile(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, ".pysar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestLoadAuthorDefaultsProfileDrivenAudience covers dec-20260725-35fa2d24's
// prediction #2: a children's-book author and a cybersecurity practitioner
// must resolve to materially different defaults for an identical broad idea,
// with no clarification question involved.
func TestLoadAuthorDefaultsProfileDrivenAudience(t *testing.T) {
	kidsRoot := t.TempDir()
	writeProfile(t, kidsRoot, "voice.md", `---
kind: voice
tone: playful, warm
formality: casual
register: simple words, picture-book cadence
notes: writes children's picture books, read aloud to kids age 4-8
---
`)
	dKids, err := LoadAuthorDefaults(kidsRoot)
	if err != nil {
		t.Fatal(err)
	}

	secRoot := t.TempDir()
	writeProfile(t, secRoot, "voice.md", `---
kind: voice
tone: precise, matter-of-fact
formality: formal
register: technical, practitioner-to-practitioner
notes: 12 years in appsec, writes for backend/security engineers
---
`)
	dSec, err := LoadAuthorDefaults(secRoot)
	if err != nil {
		t.Fatal(err)
	}

	if dKids.Source != "profile" || dSec.Source != "profile" {
		t.Fatalf("expected source=profile for both, got kids=%q sec=%q", dKids.Source, dSec.Source)
	}
	if dKids.AudienceHint == dSec.AudienceHint {
		t.Fatalf("expected different audience hints, both got %q", dKids.AudienceHint)
	}
	if dKids.AudienceHint != "children / young readers" {
		t.Fatalf("kids audience_hint=%q", dKids.AudienceHint)
	}
	if dSec.AudienceHint != "peers / practitioners in the author's field" {
		t.Fatalf("security audience_hint=%q", dSec.AudienceHint)
	}
}

// TestLoadAuthorDefaultsStyleOnlyProfile covers the style.md-without-voice.md
// case: register/tone/formality must come from the author's own profile, not
// silently stay at the generalist-fallback text just because voice.md is absent.
func TestLoadAuthorDefaultsStyleOnlyProfile(t *testing.T) {
	root := t.TempDir()
	writeProfile(t, root, "style.md", `---
kind: style
tone: precise, matter-of-fact
formality: formal
register: technical, practitioner-to-practitioner
notes: appsec practitioner voice
---
`)
	d, err := LoadAuthorDefaults(root)
	if err != nil {
		t.Fatal(err)
	}
	if d.Source != "profile" {
		t.Fatalf("source=%q, want profile", d.Source)
	}
	if d.Register != "technical, practitioner-to-practitioner" {
		t.Fatalf("register=%q, want style.md's own register, not the generalist fallback", d.Register)
	}
	if d.Tone != "precise, matter-of-fact" {
		t.Fatalf("tone=%q, want style.md's own tone, not the generalist fallback", d.Tone)
	}
	if d.AudienceHint != "peers / practitioners in the author's field" {
		t.Fatalf("audience_hint=%q", d.AudienceHint)
	}
}
