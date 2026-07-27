// Package export copies a piece's most-refined revision out of
// .pysar/pieces/<slug>/ to a single Markdown file at the project root --
// the one step in the pipeline that produces something outside the piece's
// own directory. Deliberately named "export", not "publish": this project
// has consistently excluded the POC's publish.md/Status-enum/SEO ceremony
// (dec-20260718-20a08f83 and its precedents), and this package does none of
// that -- it is a plain file copy with no side effects on the piece itself.
package export

import (
	"strings"

	"pysar/internal/validation"
)

// Bundle is the structured export payload the MCP tool validates before
// writing. Unlike the editorial-stage bundles (draft/staffedit/sharpen/
// humanize), there is no revised_md or checks field -- export doesn't
// produce new prose, it only copies whichever file the piece already has.
type Bundle struct {
	// PiecePath is the existing piece to export, relative to the project
	// root -- resolved via research.ResolvePieceDir, same as every other
	// piece-anchored tool.
	PiecePath string `json:"piece_path"`
}

// Validate checks only that piece_path was given -- there's no prose
// content in this bundle for draft.ValidateContent's citation-integrity
// rules to apply to (the exported file was already validated when its
// owning stage saved it).
func Validate(b Bundle) error {
	var missing []string
	if strings.TrimSpace(b.PiecePath) == "" {
		missing = append(missing, "piece_path")
	}
	if len(missing) > 0 {
		return &validation.Error{Kind: "export bundle", Missing: missing}
	}
	return nil
}
