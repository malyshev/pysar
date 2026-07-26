package research

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"pysar/internal/intake"
)

// ResearchDir is the standalone research output root — a sibling of
// intake's PiecesDir, not nested under it, since standalone research
// exists before any piece does.
const ResearchDir = ".pysar/research"

// standaloneMarkerFile is what AllocateUniqueNameUnder checks for to
// detect an existing standalone research directory (research has no
// brief.md the way a piece does).
const standaloneMarkerFile = "sources.md"

// ResolvePieceDir returns the piece directory for a given piece_path input,
// tolerating the directory itself or a path to any file inside it -- a
// direct child like brief.md, or something nested deeper like
// raw/<shortname>.md. Referencing a piece via Claude Code's native @
// file-picker is normal input, not a special case -- @ commonly resolves to
// a specific file rather than the directory itself, and that file isn't
// always a direct child.
//
// Real piece directories are identified by containing brief.md, not by
// string-shape guesses like "does the last path segment have a dot in it"
// (that heuristic both misresolves anything nested more than one level deep,
// and misfires on a directory name that legitimately contains a dot). So
// this walks up from the given path looking for the real marker on disk.
//
// found reports whether brief.md actually exists at the returned dir --
// callers that only need to know "is this piece real" (e.g. building an
// editorial.State) should use this instead of re-stat'ing brief.md
// themselves; this function already paid for that exact check internally.
func ResolvePieceDir(projectRoot, piecePath string) (dir string, found bool) {
	dir = filepath.Clean(filepath.Join(projectRoot, filepath.FromSlash(strings.TrimSpace(piecePath))))
	root := filepath.Clean(projectRoot)
	for {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			if _, err := os.Stat(filepath.Join(dir, "brief.md")); err == nil {
				return dir, true
			}
		}
		if dir == root {
			return dir, false
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return dir, false // reached filesystem root without finding brief.md
		}
		dir = parent
	}
}

// WriteStandalone validates b, allocates a fresh unique directory under
// ResearchDir, and writes sources.md, raw/, competitors.md (if any), and
// a research-summary.md — no piece exists yet, so there's nothing to
// append to. research-summary.md is plain Markdown, consumable directly
// by /ps-intake --from-draft= if the operator wants to seed a piece from
// what this research found.
func WriteStandalone(projectRoot string, b Bundle) (dir string, err error) {
	if err := Validate(b); err != nil {
		return "", err
	}
	if strings.TrimSpace(b.PiecePath) != "" {
		return "", fmt.Errorf("WriteStandalone called with piece_path set — use WriteToPiece instead")
	}

	name, aerr := intake.AllocateUniqueNameUnder(projectRoot, ResearchDir, standaloneMarkerFile, b.Topic)
	if aerr != nil {
		return "", aerr
	}
	dir = filepath.Join(projectRoot, ResearchDir, name)
	if err := os.MkdirAll(filepath.Join(dir, "raw"), 0o755); err != nil {
		return "", fmt.Errorf("create research dir: %w", err)
	}

	if err := writeSourcesAndRaw(dir, b); err != nil {
		return "", err
	}
	if len(b.Competitors) > 0 {
		if err := os.WriteFile(filepath.Join(dir, "competitors.md"), []byte(renderCompetitors(b.Competitors)), 0o644); err != nil {
			return "", fmt.Errorf("write competitors.md: %w", err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "research-summary.md"), []byte(renderStandaloneSummary(b)), 0o644); err != nil {
		return "", fmt.Errorf("write research-summary.md: %w", err)
	}
	return dir, nil
}

// WriteToPiece validates b and applies its research to an existing piece
// directory: sources.md and raw/ are (re)written wholesale — they're
// research-owned, no author prose to preserve, same as intake's own
// sources.md stub-to-real transition. competitors.md is written only if
// competitors were supplied. brief.md's research_mode is set to "full" and
// key_questions gets new items appended; angles.md's Misconceptions and
// Contrarian sections get new items appended. Everything else in both
// files — thesis, killer_sections, counterintuitive, the operator's own
// angles — is never touched. "We found sources; your take stayed yours."
func WriteToPiece(projectRoot string, b Bundle) error {
	if err := Validate(b); err != nil {
		return err
	}
	if strings.TrimSpace(b.PiecePath) == "" {
		return fmt.Errorf("WriteToPiece called with no piece_path — use WriteStandalone instead")
	}
	dir, _ := ResolvePieceDir(projectRoot, b.PiecePath)

	briefPath := filepath.Join(dir, "brief.md")
	briefRaw, err := os.ReadFile(briefPath)
	if err != nil {
		return fmt.Errorf("no piece found at %s (missing brief.md) — check the path; it's relative to the project root, e.g. .pysar/pieces/<name>/: %w", b.PiecePath, err)
	}
	anglesPath := filepath.Join(dir, "angles.md")
	anglesRaw, err := os.ReadFile(anglesPath)
	if err != nil {
		return fmt.Errorf("piece at %s has no angles.md: %w", b.PiecePath, err)
	}

	// Compute the new brief/angles content in memory first, but write the
	// actual research (sources/raw/competitors) to disk BEFORE touching
	// brief.md/angles.md. If writing the research fails partway through, the
	// piece must never end up claiming research_mode: full with no sources
	// to back it -- that ordering is what makes this failure unreachable.
	brief := setResearchModeFull(string(briefRaw))
	if len(b.KeyQuestionsAdditions) > 0 {
		brief = appendNumberedSection(brief, "## Key questions to answer", b.KeyQuestionsAdditions)
	}

	angles := string(anglesRaw)
	if len(b.AnglesMisconceptions) > 0 {
		angles = appendBulletSection(angles, "## Misconceptions", b.AnglesMisconceptions)
	}
	if len(b.AnglesContrarian) > 0 {
		angles = appendBulletSection(angles, "## Contrarian / under-discussed", b.AnglesContrarian)
	}

	if err := writeSourcesAndRaw(dir, b); err != nil {
		return err
	}
	if len(b.Competitors) > 0 {
		if err := os.WriteFile(filepath.Join(dir, "competitors.md"), []byte(renderCompetitors(b.Competitors)), 0o644); err != nil {
			return fmt.Errorf("write competitors.md: %w", err)
		}
	}

	if err := os.WriteFile(briefPath, []byte(brief), 0o644); err != nil {
		return fmt.Errorf("write brief.md: %w", err)
	}
	if err := os.WriteFile(anglesPath, []byte(angles), 0o644); err != nil {
		return fmt.Errorf("write angles.md: %w", err)
	}

	return appendResearchChangelog(dir, b)
}

// LoadShortnames returns the set of Shortnames a piece's research has
// recorded, for callers outside this package (e.g. internal/draft, which
// validates [^shortname] citation markers against real sources) that need
// to know what's citable without depending on this package's internal
// manifest format. Empty (not an error) when no research has run yet for
// this piece -- that's a legitimate state, not a fault.
func LoadShortnames(pieceDir string) (map[string]bool, error) {
	sources, err := loadSourcesManifest(pieceDir)
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(sources))
	for _, s := range sources {
		set[s.Shortname] = true
	}
	return set, nil
}

// sourcesManifestFile stores the piece's full accumulated source list as
// JSON -- the source of truth mergeSources reads and writes -- so a second
// /ps-research pass against the same piece adds to sources.md instead of
// silently overwriting it with only that pass's own batch. sources.md
// itself stays a pure rendered projection of this manifest.
const sourcesManifestFile = ".sources.json"

func writeSourcesAndRaw(dir string, b Bundle) error {
	if err := os.MkdirAll(filepath.Join(dir, "raw"), 0o755); err != nil {
		return fmt.Errorf("create raw dir: %w", err)
	}
	existing, err := loadSourcesManifest(dir)
	if err != nil {
		return err
	}
	merged := mergeSources(existing, b.Sources)
	if err := saveSourcesManifest(dir, merged); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "sources.md"), []byte(renderSources(b.TopicFamilyNote, merged)), 0o644); err != nil {
		return fmt.Errorf("write sources.md: %w", err)
	}
	for _, s := range merged {
		// Shortname is validated kebab-case ([a-z0-9-]) before Validate lets
		// a bundle reach here, so it's a safe bare filename component -- no
		// "/" or ".." can arrive in this path.
		path := filepath.Join(dir, "raw", s.Shortname+".md")
		body := fmt.Sprintf("<!-- source: %s (%s) -->\n\n%s\n", s.Shortname, s.URL, strings.TrimSpace(s.RawExcerpt))
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return fmt.Errorf("write raw/%s.md: %w", s.Shortname, err)
		}
	}
	return nil
}

func loadSourcesManifest(dir string) ([]Source, error) {
	raw, err := os.ReadFile(filepath.Join(dir, sourcesManifestFile))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil // first research pass against this directory -- nothing to merge yet
	}
	if err != nil {
		return nil, fmt.Errorf("read sources manifest: %w", err)
	}
	var sources []Source
	if err := json.Unmarshal(raw, &sources); err != nil {
		return nil, fmt.Errorf("parse sources manifest: %w", err)
	}
	return sources, nil
}

func saveSourcesManifest(dir string, sources []Source) error {
	raw, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		return fmt.Errorf("encode sources manifest: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, sourcesManifestFile), raw, 0o644); err != nil {
		return fmt.Errorf("write sources manifest: %w", err)
	}
	return nil
}

// mergeSources combines a piece's previously-accumulated sources with a new
// pass's batch. A shortname reused across passes supersedes the prior entry
// for that shortname -- the same citation-hygiene rule that keeps shortnames
// unique within a single pass, extended across passes -- rather than
// duplicating or erroring.
func mergeSources(existing, incoming []Source) []Source {
	merged := append([]Source{}, existing...)
	index := make(map[string]int, len(merged))
	for i, s := range merged {
		index[s.Shortname] = i
	}
	for _, s := range incoming {
		if i, ok := index[s.Shortname]; ok {
			merged[i] = s
		} else {
			merged = append(merged, s)
			index[s.Shortname] = len(merged) - 1
		}
	}
	return merged
}

func renderSources(topicFamilyNote string, sources []Source) string {
	var out strings.Builder
	out.WriteString("# Sources\n\n")
	if topicFamilyNote != "" {
		fmt.Fprintf(&out, "_Topic-family mapping: %s_\n\n", strings.TrimSpace(topicFamilyNote))
	}
	recognized := 0
	for _, s := range sources {
		if s.Tier != TierCommunity {
			recognized++
		}
	}
	fmt.Fprintf(&out, "%d sources, %d/%d primary+secondary (authority floor: >=60%%).\n\n", len(sources), recognized, len(sources))
	for _, s := range sources {
		fmt.Fprintf(&out, "## %s\n\n", s.Shortname)
		fmt.Fprintf(&out, "- URL: %s\n", s.URL)
		fmt.Fprintf(&out, "- Tier: %s\n", s.Tier)
		fmt.Fprintf(&out, "- Accessed: %s\n", s.Accessed)
		out.WriteString("- Key claims:\n")
		for _, c := range s.KeyClaims {
			fmt.Fprintf(&out, "  - %s\n", strings.TrimSpace(c))
		}
		fmt.Fprintf(&out, "- Notes: %s\n\n", strings.TrimSpace(s.Notes))
	}
	return out.String()
}

func renderCompetitors(competitors []Competitor) string {
	var out strings.Builder
	out.WriteString("# Competitors\n\n")
	for _, c := range competitors {
		fmt.Fprintf(&out, "## %s\n\n", c.URL)
		fmt.Fprintf(&out, "- Angle: %s\n", strings.TrimSpace(c.Angle))
		fmt.Fprintf(&out, "- Strongest section: %s\n", strings.TrimSpace(c.StrongestSection))
		fmt.Fprintf(&out, "- Gap: %s\n\n", strings.TrimSpace(c.Gap))
	}
	return out.String()
}

func renderStandaloneSummary(b Bundle) string {
	var out strings.Builder
	fmt.Fprintf(&out, "# Research: %s\n\n", strings.TrimSpace(b.Topic))
	fmt.Fprintf(&out, "Expert lens: %s\n\n", strings.TrimSpace(b.ExpertLens))
	if len(b.KeyQuestionsAdditions) > 0 {
		out.WriteString("## Questions this research answers\n\n")
		for i, q := range b.KeyQuestionsAdditions {
			fmt.Fprintf(&out, "%d. %s\n", i+1, strings.TrimSpace(q))
		}
		out.WriteString("\n")
	}
	if len(b.AnglesMisconceptions) > 0 {
		out.WriteString("## Misconceptions this research surfaced\n\n")
		for _, m := range b.AnglesMisconceptions {
			fmt.Fprintf(&out, "- %s\n", strings.TrimSpace(m))
		}
		out.WriteString("\n")
	}
	if len(b.AnglesContrarian) > 0 {
		out.WriteString("## Contrarian / under-discussed angles this research surfaced\n\n")
		for _, c := range b.AnglesContrarian {
			fmt.Fprintf(&out, "- %s\n", strings.TrimSpace(c))
		}
		out.WriteString("\n")
	}
	out.WriteString("See sources.md and raw/ for the full research. No thesis or angle is committed here — this is background, not a brief. Run /ps-intake --from-draft= on this file to turn it into one.\n")
	return out.String()
}

var researchModeLineRe = regexp.MustCompile(`(?m)^research_mode:\s*\S+`)

func setResearchModeFull(brief string) string {
	if researchModeLineRe.MatchString(brief) {
		return researchModeLineRe.ReplaceAllString(brief, "research_mode: full")
	}
	return brief
}

var numberedItemRe = regexp.MustCompile(`(?m)^(\d+)\.\s`)

// appendNumberedSection appends newItems to the numbered list under
// heading, continuing the existing numbering and replacing intake's
// "(none listed...)" placeholder if that's all the section currently has.
// An item whose text is already present in the section is skipped -- this
// is what keeps a retried save_research_bundle call (e.g. after run-log
// append failed on an otherwise-successful save) from duplicating entries,
// since this append is otherwise not idempotent.
func appendNumberedSection(content, heading string, newItems []string) string {
	section, before, after := splitSection(content, heading)
	section = strings.ReplaceAll(section, "(none listed — expand while drafting if needed)\n", "")
	max := 0
	for _, m := range numberedItemRe.FindAllStringSubmatch(section, -1) {
		if n, convErr := strconv.Atoi(m[1]); convErr == nil && n > max {
			max = n
		}
	}
	var added strings.Builder
	n := max
	for _, item := range newItems {
		item = strings.TrimSpace(item)
		if item == "" || strings.Contains(section, item) {
			continue
		}
		n++
		fmt.Fprintf(&added, "%d. %s\n", n, item)
	}
	if added.Len() == 0 {
		return content
	}
	return before + heading + "\n\n" + strings.TrimLeft(strings.TrimRight(section, "\n")+"\n"+added.String(), "\n") + after
}

// appendBulletSection appends newItems as "- " bullets under heading, same
// idempotency-on-retry behavior as appendNumberedSection.
func appendBulletSection(content, heading string, newItems []string) string {
	section, before, after := splitSection(content, heading)
	var added strings.Builder
	for _, item := range newItems {
		item = strings.TrimSpace(item)
		if item == "" || strings.Contains(section, item) {
			continue
		}
		fmt.Fprintf(&added, "- %s\n", item)
	}
	if added.Len() == 0 {
		return content
	}
	return before + heading + "\n\n" + strings.TrimLeft(strings.TrimRight(section, "\n")+"\n"+added.String(), "\n") + after
}

// splitSection locates heading (an exact "## ..." line) in content and
// returns the section body between it and the next "## " heading (or EOF),
// plus the content before and after that body so callers can splice a
// modified section back in.
//
// intake's own Validate() only guarantees a case-insensitive "contrarian"
// substring appears somewhere in angles_md -- not this exact heading text --
// so a valid, Validate-passing angles.md can lack "## Misconceptions" or
// "## Contrarian / under-discussed" outright. Rather than erroring and
// losing the research that was about to be appended, a missing heading is
// treated as an empty section freshly created at the end of the document.
func splitSection(content, heading string) (section, before, after string) {
	idx := strings.Index(content, heading)
	if idx < 0 {
		return "", strings.TrimRight(content, "\n") + "\n\n", ""
	}
	bodyStart := idx + len(heading)
	rest := content[bodyStart:]
	rest = strings.TrimPrefix(rest, "\n")
	nextIdx := strings.Index(rest, "\n## ")
	if nextIdx < 0 {
		return rest, content[:idx], ""
	}
	return rest[:nextIdx], content[:idx], rest[nextIdx:]
}

func appendResearchChangelog(dir string, b Bundle) error {
	line := intake.FormatChangelogLine(fmt.Sprintf("sources=%d competitors=%d", len(b.Sources), len(b.Competitors)))
	return intake.AppendChangelogLine(dir, "research-changelog.md", line)
}
