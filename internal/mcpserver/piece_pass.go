package mcpserver

import (
	"fmt"

	"pysar/internal/editorial"
	"pysar/internal/research"
)

// resolveAnchoredPass is the common "is this piece-anchored MCP call
// allowed to proceed" dance every such tool needs: find the named pass,
// resolve piecePath to its directory, build the piece's current
// produced-artifact State from what ResolvePieceDir already found on disk,
// and run the pass's own precondition. Shared so a future piece-anchored
// tool (a third one, after research and draft) doesn't hand-roll a fourth
// copy of this block -- and so a precondition change only has one place to
// update instead of one per tool.
func (s *Server) resolveAnchoredPass(passName, piecePath string) (pieceDir string, err error) {
	pass := findPass(passName)
	if pass == nil {
		return "", fmt.Errorf("%s pass is not registered", passName)
	}
	pieceDir, found := research.ResolvePieceDir(s.baseDir, piecePath)
	state := editorial.NewState(editorial.SurfaceBlog)
	if found {
		state = state.WithProduced(editorial.ArtifactStake).WithProduced(editorial.ArtifactBrief)
	}
	if _, err := editorial.Run(pass, state); err != nil {
		return "", err
	}
	return pieceDir, nil
}
