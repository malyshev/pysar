package mcpserver

import (
	"fmt"
	"path/filepath"

	"pysar/internal/editorial"
	"pysar/internal/research"
)

// resolveAnchoredPass is the common "is this piece-anchored MCP call
// allowed to proceed" dance every such tool needs: find the named pass,
// resolve piecePath to its directory, build the piece's current
// produced-artifact State from what's actually on disk, and run the pass's
// own precondition. Shared so a future piece-anchored tool doesn't
// hand-roll another copy of this block -- and so a precondition change
// only has one place to update instead of one per tool.
//
// State is built from cheap, targeted file-existence checks (brief.md ->
// Stake+Brief, draft.md -> Draft) rather than a generic "scan every known
// marker file" abstraction -- add one more check here if a future pass
// needs a further artifact, the same way this one was added for staff-edit.
func (s *Server) resolveAnchoredPass(passName, piecePath string) (pieceDir string, err error) {
	pass := findPass(passName)
	if pass == nil {
		return "", fmt.Errorf("%s pass is not registered", passName)
	}
	pieceDir, found := research.ResolvePieceDir(s.baseDir, piecePath)
	state := editorial.NewState(editorial.SurfaceBlog)
	if found {
		state = state.WithProduced(editorial.ArtifactStake).WithProduced(editorial.ArtifactBrief)
		// Only look for further piece artifacts once the piece itself is
		// confirmed real -- when !found, pieceDir has fallen back to
		// s.baseDir (ResolvePieceDir's not-found case), and checking for
		// draft.md there would pick up an unrelated file at the project
		// root instead of correctly reporting no real piece was found.
		if fileExists(filepath.Join(pieceDir, "draft.md")) {
			state = state.WithProduced(editorial.ArtifactDraft)
		}
	}
	if _, err := editorial.Run(pass, state); err != nil {
		return "", err
	}
	return pieceDir, nil
}
