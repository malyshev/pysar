package research

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validSource(shortname string) Source {
	return Source{
		Shortname:  shortname,
		URL:        "https://example.com/" + shortname,
		Tier:       TierPrimary,
		Accessed:   "2026-07-25",
		KeyClaims:  []string{"a real claim this source backs"},
		Notes:      "a real note about this source",
		RawExcerpt: "a verbatim quoted sentence from the source",
	}
}

func TestValidateRequiresTopicOrPiecePath(t *testing.T) {
	b := Bundle{ExpertLens: "x", Sources: []Source{validSource("a")}}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected error when neither topic nor piece_path is set")
	}
	ve := err.(*ValidationError)
	joined := strings.Join(ve.Missing, " ")
	if !strings.Contains(joined, "topic") {
		t.Fatalf("expected missing list to name topic, got %v", ve.Missing)
	}
}

func TestValidateRequiresAtLeastOneSource(t *testing.T) {
	b := Bundle{Topic: "x", ExpertLens: "x"}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected error for zero sources")
	}
}

func TestValidateRejectsDuplicateShortnames(t *testing.T) {
	b := Bundle{
		Topic:      "x",
		ExpertLens: "x",
		Sources:    []Source{validSource("dup"), validSource("dup")},
	}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected error for duplicate shortnames")
	}
	ve := err.(*ValidationError)
	joined := strings.Join(ve.Missing, " ")
	if !strings.Contains(joined, "duplicate") {
		t.Fatalf("expected duplicate shortname in missing list, got %v", ve.Missing)
	}
}

func TestValidateRejectsPathTraversalShortname(t *testing.T) {
	s := validSource("x")
	s.Shortname = "../../../../etc/cron.d/evil"
	b := Bundle{Topic: "x", ExpertLens: "x", Sources: []Source{s}}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected error for a shortname containing path-traversal characters")
	}
	ve := err.(*ValidationError)
	joined := strings.Join(ve.Missing, " ")
	if !strings.Contains(joined, "kebab-case") {
		t.Fatalf("expected kebab-case rejection in missing list, got %v", ve.Missing)
	}
}

func TestValidateRejectsUppercaseShortname(t *testing.T) {
	// Uppercase is rejected outright rather than merely deduped case-
	// insensitively -- this is what prevents "Docker-Docs" and
	// "docker-docs" from both validating and then colliding as the same
	// raw/<shortname>.md file on a case-insensitive filesystem.
	s := validSource("x")
	s.Shortname = "Docker-Docs"
	b := Bundle{Topic: "x", ExpertLens: "x", Sources: []Source{s}}
	if err := Validate(b); err == nil {
		t.Fatal("expected error for an uppercase shortname")
	}
}

func TestValidateEnforcesAuthorityFloor(t *testing.T) {
	// 1 primary out of 3 = 33% recognized -- fails the 60% floor.
	b := Bundle{
		Topic:      "x",
		ExpertLens: "x",
		Sources: []Source{
			func() Source { s := validSource("s1"); s.Tier = TierPrimary; return s }(),
			func() Source { s := validSource("s2"); s.Tier = TierCommunity; return s }(),
			func() Source { s := validSource("s3"); s.Tier = TierCommunity; return s }(),
		},
	}
	err := Validate(b)
	if err == nil {
		t.Fatal("expected authority-floor rejection")
	}
	ve := err.(*ValidationError)
	joined := strings.Join(ve.Missing, " ")
	if !strings.Contains(joined, "authority floor") {
		t.Fatalf("expected authority floor message, got %v", ve.Missing)
	}
}

func TestValidateAcceptsFloorExactlyMet(t *testing.T) {
	// 3 primary/secondary out of 5 = 60% -- exactly meets the floor.
	b := Bundle{
		Topic:      "x",
		ExpertLens: "x",
		Sources: []Source{
			func() Source { s := validSource("s1"); s.Tier = TierPrimary; return s }(),
			func() Source { s := validSource("s2"); s.Tier = TierSecondary; return s }(),
			func() Source { s := validSource("s3"); s.Tier = TierPrimary; return s }(),
			func() Source { s := validSource("s4"); s.Tier = TierCommunity; return s }(),
			func() Source { s := validSource("s5"); s.Tier = TierCommunity; return s }(),
		},
	}
	if err := Validate(b); err != nil {
		t.Fatalf("expected floor-exactly-met bundle to validate, got: %v", err)
	}
}

func TestWriteStandaloneCreatesResearchDirWithSuffix(t *testing.T) {
	root := t.TempDir()
	b := Bundle{
		Topic:                 "tomato blossom end rot",
		ExpertLens:            "horticulture",
		Sources:               []Source{validSource("uconn-ber")},
		KeyQuestionsAdditions: []string{"Does calcium supplementation fix it?"},
		AnglesMisconceptions:  []string{"Soil lacks calcium -- usually false"},
	}
	dir, err := WriteStandalone(root, b)
	if err != nil {
		t.Fatalf("WriteStandalone: %v", err)
	}
	if !strings.HasPrefix(filepath.Base(dir), "tomato-blossom-end-rot-") {
		t.Fatalf("expected suffixed dir, got %q", filepath.Base(dir))
	}
	for _, f := range []string{"sources.md", "research-summary.md", "raw/uconn-ber.md"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
	summary, _ := os.ReadFile(filepath.Join(dir, "research-summary.md"))
	if !strings.Contains(string(summary), "Does calcium supplementation fix it?") {
		t.Fatalf("summary missing key question:\n%s", summary)
	}
	if !strings.Contains(string(summary), "/ps-intake --from-draft=") {
		t.Fatalf("summary missing hand-off hint to /ps-intake:\n%s", summary)
	}
}

func TestWriteToPieceAppendsWithoutTouchingThesis(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-abc123")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	brief := `---
piece: example-abc123
research_mode: none
pov_source: author-practitioner
entry_mode: idea
audience: "x"
register: "x"
expert_lens: "horticulture"
topic_scope: "x"
storage_status: provisional
---

# Brief

**Restatement:** x

**Thesis:** The operator's own thesis, must survive untouched.

## Promise to the reader

x

## Killer sections / competitive edge

1. **A killer section**
   - Edge: x
   - Example: x

## Counterintuitive findings to elevate

1. **A counterintuitive claim**
   - Contradiction: x

## Key questions to answer

1. An existing question

## Non-goals

- x
`
	angles := `## Misconceptions
- An existing misconception

## Contrarian / under-discussed
- An existing contrarian angle

## Trade-offs
- x

## Edge cases
- x
`
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte(brief), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "angles.md"), []byte(angles), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Bundle{
		PiecePath:             ".pysar/pieces/example-abc123",
		ExpertLens:            "horticulture",
		Sources:               []Source{validSource("s1")},
		KeyQuestionsAdditions: []string{"A new research-backed question"},
		AnglesMisconceptions:  []string{"A new research-backed misconception"},
		AnglesContrarian:      []string{"A new research-backed contrarian angle"},
	}
	if err := WriteToPiece(root, b); err != nil {
		t.Fatalf("WriteToPiece: %v", err)
	}

	gotBrief, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	briefStr := string(gotBrief)
	if !strings.Contains(briefStr, "research_mode: full") {
		t.Fatalf("expected research_mode: full:\n%s", briefStr)
	}
	if !strings.Contains(briefStr, "The operator's own thesis, must survive untouched.") {
		t.Fatalf("thesis was touched:\n%s", briefStr)
	}
	if !strings.Contains(briefStr, "1. An existing question") || !strings.Contains(briefStr, "2. A new research-backed question") {
		t.Fatalf("key questions not correctly continued:\n%s", briefStr)
	}
	if !strings.Contains(briefStr, "A killer section") {
		t.Fatalf("killer sections were touched/lost:\n%s", briefStr)
	}
	if !strings.Contains(briefStr, "## Non-goals") {
		t.Fatalf("non-goals section lost:\n%s", briefStr)
	}

	gotAngles, _ := os.ReadFile(filepath.Join(pieceDir, "angles.md"))
	anglesStr := string(gotAngles)
	if !strings.Contains(anglesStr, "An existing misconception") || !strings.Contains(anglesStr, "A new research-backed misconception") {
		t.Fatalf("misconceptions not correctly appended:\n%s", anglesStr)
	}
	if !strings.Contains(anglesStr, "An existing contrarian angle") || !strings.Contains(anglesStr, "A new research-backed contrarian angle") {
		t.Fatalf("contrarian angles not correctly appended:\n%s", anglesStr)
	}
	if !strings.Contains(anglesStr, "## Trade-offs") {
		t.Fatalf("trade-offs section lost:\n%s", anglesStr)
	}

	for _, f := range []string{"sources.md", "raw/s1.md", "research-changelog.md"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
}

func TestWriteToPieceReplacesNoneListedPlaceholder(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-xyz789")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	brief := `---
research_mode: none
---

# Brief

## Key questions to answer

(none listed — expand while drafting if needed)

## Non-goals

- x
`
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte(brief), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "angles.md"), []byte("## Misconceptions\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Bundle{
		PiecePath:             ".pysar/pieces/example-xyz789",
		ExpertLens:            "x",
		Sources:               []Source{validSource("s1")},
		KeyQuestionsAdditions: []string{"First real question"},
	}
	if err := WriteToPiece(root, b); err != nil {
		t.Fatalf("WriteToPiece: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	gotStr := string(got)
	if strings.Contains(gotStr, "none listed") {
		t.Fatalf("placeholder should have been replaced:\n%s", gotStr)
	}
	if !strings.Contains(gotStr, "1. First real question") {
		t.Fatalf("expected numbering to start at 1:\n%s", gotStr)
	}
}

func TestWriteToPieceMergesSourcesAcrossRepeatRuns(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-merge")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte("---\nresearch_mode: none\n---\n# Brief\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "angles.md"), []byte("## Misconceptions\n\n## Contrarian / under-discussed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	first := Bundle{PiecePath: ".pysar/pieces/example-merge", ExpertLens: "x", Sources: []Source{validSource("first-source")}}
	if err := WriteToPiece(root, first); err != nil {
		t.Fatalf("first WriteToPiece: %v", err)
	}
	second := Bundle{PiecePath: ".pysar/pieces/example-merge", ExpertLens: "x", Sources: []Source{validSource("second-source")}}
	if err := WriteToPiece(root, second); err != nil {
		t.Fatalf("second WriteToPiece: %v", err)
	}

	sourcesMD, _ := os.ReadFile(filepath.Join(pieceDir, "sources.md"))
	s := string(sourcesMD)
	if !strings.Contains(s, "first-source") {
		t.Fatalf("second research run lost the first run's source:\n%s", s)
	}
	if !strings.Contains(s, "second-source") {
		t.Fatalf("second research run's own source missing:\n%s", s)
	}
	if !strings.Contains(s, "2 sources") {
		t.Fatalf("expected sources.md to report the merged count of 2:\n%s", s)
	}
	for _, f := range []string{"raw/first-source.md", "raw/second-source.md"} {
		if _, err := os.Stat(filepath.Join(pieceDir, f)); err != nil {
			t.Fatalf("missing %s after merge: %v", f, err)
		}
	}
}

func TestWriteToPieceAppendIsIdempotentOnRetry(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-retry")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte("---\nresearch_mode: none\n---\n# Brief\n\n## Key questions to answer\n\n(none listed — expand while drafting if needed)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "angles.md"), []byte("## Misconceptions\n\n## Contrarian / under-discussed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Bundle{
		PiecePath:             ".pysar/pieces/example-retry",
		ExpertLens:            "x",
		Sources:               []Source{validSource("s1")},
		KeyQuestionsAdditions: []string{"A question a client might resubmit after a retry"},
	}
	// Simulate a client retrying the same save_research_bundle call after a
	// failure downstream of WriteToPiece (e.g. the run-log append) -- the
	// append must not duplicate the entry the second time.
	if err := WriteToPiece(root, b); err != nil {
		t.Fatalf("first WriteToPiece: %v", err)
	}
	if err := WriteToPiece(root, b); err != nil {
		t.Fatalf("retried WriteToPiece: %v", err)
	}

	got, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	count := strings.Count(string(got), "A question a client might resubmit after a retry")
	if count != 1 {
		t.Fatalf("expected the retried question to appear exactly once, got %d:\n%s", count, got)
	}
}

func TestWriteToPieceCreatesMissingAnglesHeading(t *testing.T) {
	// intake's own Validate() only guarantees a case-insensitive
	// "contrarian" substring somewhere in angles_md, never this exact
	// heading -- so a Validate-passing angles.md can legitimately lack
	// "## Misconceptions" entirely.
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-no-heading")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte("---\nresearch_mode: none\n---\n# Brief\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "angles.md"), []byte("## Contrarian takes\n- an existing angle\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Bundle{
		PiecePath:            ".pysar/pieces/example-no-heading",
		ExpertLens:           "x",
		Sources:              []Source{validSource("s1")},
		AnglesMisconceptions: []string{"A new research-backed misconception"},
	}
	if err := WriteToPiece(root, b); err != nil {
		t.Fatalf("WriteToPiece should tolerate a missing heading, got: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(pieceDir, "angles.md"))
	s := string(got)
	if !strings.Contains(s, "## Misconceptions") || !strings.Contains(s, "A new research-backed misconception") {
		t.Fatalf("expected a freshly-created Misconceptions section:\n%s", s)
	}
	if !strings.Contains(s, "## Contrarian takes") || !strings.Contains(s, "an existing angle") {
		t.Fatalf("existing content should survive untouched:\n%s", s)
	}
}

func TestWriteToPieceLeavesBriefUntouchedWhenSourcesWriteFails(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-partial-fail")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte("---\nresearch_mode: none\n---\n# Brief\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "angles.md"), []byte("## Misconceptions\n\n## Contrarian / under-discussed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Force writeSourcesAndRaw's os.MkdirAll(".../raw", ...) to fail by
	// pre-creating "raw" as a plain file instead of a directory.
	if err := os.WriteFile(filepath.Join(pieceDir, "raw"), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Bundle{PiecePath: ".pysar/pieces/example-partial-fail", ExpertLens: "x", Sources: []Source{validSource("s1")}}
	if err := WriteToPiece(root, b); err == nil {
		t.Fatal("expected WriteToPiece to fail when the research write fails")
	}

	got, _ := os.ReadFile(filepath.Join(pieceDir, "brief.md"))
	if strings.Contains(string(got), "research_mode: full") {
		t.Fatalf("brief.md must not claim research_mode: full when sources were never written:\n%s", got)
	}
	if strings.Contains(string(got), "research_mode: none") == false {
		t.Fatalf("expected brief.md to remain at its original research_mode:\n%s", got)
	}
}

func TestWriteToPieceRejectsMissingPiece(t *testing.T) {
	root := t.TempDir()
	b := Bundle{
		PiecePath:  ".pysar/pieces/does-not-exist",
		ExpertLens: "x",
		Sources:    []Source{validSource("s1")},
	}
	if err := WriteToPiece(root, b); err == nil {
		t.Fatal("expected error for a piece path that doesn't exist")
	}
}

func TestResolvePieceDirAcceptsDirectoryOrFileInside(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-abc")
	if err := os.MkdirAll(filepath.Join(pieceDir, "raw"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte("# Brief\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "raw", "uconn-ber.md"), []byte("excerpt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// A piece directory name that legitimately contains a dot -- not
	// guaranteed impossible, since piece_path is free-form input, not
	// required to have come through onboarding.Slug.
	dottedDir := filepath.Join(root, ".pysar", "pieces", "my-v1.0-piece")
	if err := os.MkdirAll(dottedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dottedDir, "brief.md"), []byte("# Brief\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"bare directory", ".pysar/pieces/example-abc", pieceDir},
		{"directory with trailing slash", ".pysar/pieces/example-abc/", pieceDir},
		{"a Claude Code @-reference resolving to brief.md", ".pysar/pieces/example-abc/brief.md", pieceDir},
		{"a file reference to angles.md (need not exist)", ".pysar/pieces/example-abc/angles.md", pieceDir},
		{"a file nested two levels deep, e.g. this feature's own raw/ output", ".pysar/pieces/example-abc/raw/uconn-ber.md", pieceDir},
		{"a piece directory name that itself contains a dot", ".pysar/pieces/my-v1.0-piece", dottedDir},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ResolvePieceDir(root, tc.input)
			if got != tc.want {
				t.Fatalf("ResolvePieceDir(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestWriteToPieceAcceptsAtReferenceToBriefFile is the integration-level
// proof: research still lands in the right directory when piece_path is a
// file inside it (what an @-reference commonly resolves to), not just the
// directory itself.
func TestWriteToPieceAcceptsAtReferenceToBriefFile(t *testing.T) {
	root := t.TempDir()
	pieceDir := filepath.Join(root, ".pysar", "pieces", "example-atref")
	if err := os.MkdirAll(pieceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "brief.md"), []byte("---\nresearch_mode: none\n---\n# Brief\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pieceDir, "angles.md"), []byte("## Misconceptions\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	b := Bundle{
		PiecePath:  ".pysar/pieces/example-atref/brief.md", // as if @-referenced
		ExpertLens: "x",
		Sources:    []Source{validSource("s1")},
	}
	if err := WriteToPiece(root, b); err != nil {
		t.Fatalf("WriteToPiece with a file-inside-piece path: %v", err)
	}
	if _, err := os.Stat(filepath.Join(pieceDir, "sources.md")); err != nil {
		t.Fatalf("expected sources.md written to the piece directory, not a sibling of brief.md: %v", err)
	}
}
