package editorial

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RunLogFilename is the append-only per-piece pass record
// (dec-20260718-6ba897d4).
const RunLogFilename = "run-log.jsonl"

// RunLogEntry is one successful pass against a piece.
type RunLogEntry struct {
	At      time.Time `json:"at"`
	Pass    string    `json:"pass"`
	Summary string    `json:"summary,omitempty"`
}

// AppendRunLog adds one JSON line to pieceDir/run-log.jsonl.
func AppendRunLog(pieceDir string, entry RunLogEntry) error {
	if entry.At.IsZero() {
		entry.At = time.Now().UTC()
	}
	path := filepath.Join(pieceDir, RunLogFilename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open run-log: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("encode run-log entry: %w", err)
	}
	return nil
}
