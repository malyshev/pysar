package intake

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pysar/internal/onboarding"
)

// PiecesDir is the provisional on-disk root for piece scaffolds
// (dec-20260725-35fa2d24). Marked provisional — topology may change when
// prob-20260718-f6ac8127 resolves.
const PiecesDir = ".pysar/pieces"

// StorageMarkerFilename is written into every piece directory so migration
// tooling and humans can see the layout is not final.
const StorageMarkerFilename = "STORAGE.md"

const storageMarkerBody = `# Provisional piece storage

This directory layout is **provisional** (dec-20260725-35fa2d24).

Final piece/project topology is still open under problem
` + "`prob-20260718-f6ac8127`" + `. Expect a migration (path rename or regrouping)
when that decision lands. Author prose in brief/outline/angles is preserved
across migration; only paths/grouping may change.
`

// PieceDir returns the absolute piece directory for name under projectRoot.
func PieceDir(projectRoot, name string) string {
	return filepath.Join(projectRoot, PiecesDir, onboarding.Slug(name))
}

// randomSuffixHex returns n random bytes hex-encoded, e.g. "0e4b95f2d211"
// for n=6.
func randomSuffixHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random suffix: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

const (
	// MaxPieceNameLength caps the whole piece directory name (readable
	// prefix + random suffix). 60 total (47 for the prefix) matches
	// real-world slug convention -- Medium's own "docker-for-developers-
	// who-just-need-it-to-work" is 46 characters -- rather than the prior
	// 120, which was technically correct but still an unreadable wall of
	// text for what's supposed to be a recognizable directory name.
	MaxPieceNameLength = 60
	suffixHexLen       = 12 // randomSuffixHex(6) always produces 12 hex chars
)

// truncatePrefix caps base so base + "-" + <12-hex-char suffix> never
// exceeds MaxPieceNameLength. Cuts at the last whole-word (hyphen)
// boundary within budget, never mid-word -- "docker-for-developers-who-
// just-need-it-to-work", never a fragment like "...to-work-co" cut out of
// "containerize". Falls back to a hard character cut only when a single
// word alone exceeds the entire budget. Only the readable prefix is ever
// shortened — the suffix, which is what actually guarantees uniqueness,
// is never touched.
func truncatePrefix(base string) string {
	budget := MaxPieceNameLength - 1 - suffixHexLen // hyphen + suffix
	if len(base) <= budget {
		return base
	}
	cut := base[:budget]
	if i := strings.LastIndex(cut, "-"); i > 0 {
		return cut[:i]
	}
	return strings.TrimRight(cut, "-")
}

// AllocateUniqueName finds an unused .pysar/pieces/<base>-<suffix> name
// under projectRoot. Every piece directory gets a random suffix, always —
// called exactly once per write, from WriteBundle, so a name never gets
// suffixed twice regardless of whether base came from an explicit author
// override, a propose_intake_name preview, or the idea text itself.
// Directory naming and collisions are storage mechanics, never a decision
// surfaced to the author. base may be empty; a bare suffix is still unique.
func AllocateUniqueName(projectRoot, base string) (string, error) {
	base = truncatePrefix(onboarding.Slug(base))
	for attempt := 0; attempt < 8; attempt++ {
		suffix, err := randomSuffixHex(6)
		if err != nil {
			return "", err
		}
		candidate := suffix
		if base != "" {
			candidate = base + "-" + suffix
		}
		if _, statErr := os.Stat(filepath.Join(PieceDir(projectRoot, candidate), "brief.md")); os.IsNotExist(statErr) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not allocate a unique piece name after several attempts")
}

// WriteBundle validates b, allocates a fresh unique piece directory, and
// writes brief/outline/angles/sources + storage marker + changelog append.
// The directory always gets a random suffix — a name collision is never
// surfaced as an error to the caller, since it's storage mechanics the
// author cannot meaningfully act on.
func WriteBundle(projectRoot string, b Bundle) (pieceDir string, err error) {
	if err := Validate(b); err != nil {
		return "", err
	}

	base := strings.TrimSpace(b.Name)
	if base == "" {
		base = b.Idea
	}
	name, aerr := AllocateUniqueName(projectRoot, base)
	if aerr != nil {
		return "", aerr
	}
	dir := PieceDir(projectRoot, name)
	b.Name = name

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create piece dir: %w", err)
	}

	files := map[string]string{
		StorageMarkerFilename: storageMarkerBody,
		"brief.md":            renderBrief(b),
		"outline.md":          strings.TrimSpace(b.OutlineMD) + "\n",
		"angles.md":           strings.TrimSpace(b.AnglesMD) + "\n",
		"sources.md":          renderSourcesStub(name, b.Sources),
	}
	for fname, body := range files {
		path := filepath.Join(dir, fname)
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return "", fmt.Errorf("write %s: %w", fname, err)
		}
	}

	if err := appendChangelog(dir, b); err != nil {
		return "", err
	}
	return dir, nil
}

func renderBrief(b Bundle) string {
	var killer strings.Builder
	for i, k := range b.KillerSections {
		fmt.Fprintf(&killer, "%d. **%s**\n   - Edge: %s\n   - Example: %s\n",
			i+1, strings.TrimSpace(k.Title), strings.TrimSpace(k.Edge), strings.TrimSpace(k.Example))
	}
	var counter strings.Builder
	for i, c := range b.Counterintuitive {
		fmt.Fprintf(&counter, "%d. **%s**\n   - Contradiction: %s\n",
			i+1, strings.TrimSpace(c.Claim), strings.TrimSpace(c.Contradiction))
	}
	var questions strings.Builder
	if len(b.KeyQuestions) == 0 {
		questions.WriteString("(none listed — expand while drafting if needed)\n")
	} else {
		for i, q := range b.KeyQuestions {
			fmt.Fprintf(&questions, "%d. %s\n", i+1, strings.TrimSpace(q))
		}
	}
	var nonGoals strings.Builder
	if len(b.NonGoals) == 0 {
		nonGoals.WriteString("- (none listed)\n")
	} else {
		for _, n := range b.NonGoals {
			fmt.Fprintf(&nonGoals, "- %s\n", strings.TrimSpace(n))
		}
	}

	return fmt.Sprintf(`---
piece: %s
research_mode: %s
pov_source: %s
entry_mode: %s
audience: %q
register: %q
expert_lens: %q
topic_scope: %q
storage_status: provisional
---

# Brief

**Restatement:** %s

**Thesis:** %s

## Promise to the reader

%s

## Killer sections / competitive edge

%s
## Counterintuitive findings to elevate

%s
## Key questions to answer

%s
## Non-goals

%s
`,
		b.Name,
		researchMode(b),
		b.POVSource,
		b.EntryMode,
		b.Audience,
		b.Register,
		b.ExpertLens,
		b.TopicScope,
		strings.TrimSpace(b.Restatement),
		strings.TrimSpace(b.Thesis),
		strings.TrimSpace(b.Promise),
		killer.String(),
		counter.String(),
		questions.String(),
		nonGoals.String(),
	)
}

// researchMode reports "targeted" when Step 4 (SKILL.md) actually fetched
// real sources to ground a technical claim, "none" otherwise. Both are
// legitimate outcomes -- most topics don't need grounding.
func researchMode(b Bundle) string {
	if len(b.Sources) > 0 {
		return "targeted"
	}
	return "none"
}

func renderSourcesStub(name string, sources []Source) string {
	if len(sources) == 0 {
		return fmt.Sprintf(`# Sources

Research mode: none — intake only for piece %q. No technical claims needed grounding beyond the author's own expertise.
Intake never fabricates citations; every source it does list here was actually fetched.
`, name)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Sources\n\nResearch mode: targeted — the following were fetched to ground specific claims in %q. Intake never fabricates citations.\n\n", name)
	for _, s := range sources {
		fmt.Fprintf(&b, "- %s\n  %s\n", strings.TrimSpace(s.URL), strings.TrimSpace(s.Note))
	}
	return b.String()
}

func appendChangelog(dir string, b Bundle) error {
	path := filepath.Join(dir, "intake-changelog.md")
	line := fmt.Sprintf("- %s — entry=%s piece=%s thesis=%q\n",
		time.Now().UTC().Format(time.RFC3339), b.EntryMode, b.Name, strings.TrimSpace(b.Thesis))
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open changelog: %w", err)
	}
	defer f.Close()
	if _, err := f.WriteString(line); err != nil {
		return fmt.Errorf("append changelog: %w", err)
	}
	return nil
}
