package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.0.1-dev" // Dev version string until release tagging exists

func main() {
	rootCmd := &cobra.Command{
		Use:   "pysar",
		Short: "An author-directed editorial engine for your writing projects",
		Long: `Pysar — an author-directed editorial engine for your writing projects.

Bring your take (an idea or a draft). Pysar helps you shape it into a ship-ready
piece and a platform-scoped publish checklist — without forcing pipeline jargon
on you, and without posting on your behalf.

This binary is the console CLI. Agentic slash commands (ps-* / /ps) are a
separate host-agent surface when installed; they are not invoked by typing pysar.`,
		Example: `  pysar init              # scaffold a Claude Code project in .
  pysar init --claude ./my-piece
  pysar --version`,
		Version: version,
	}
	rootCmd.SetVersionTemplate("pysar {{.Version}}\n")

	rootCmd.AddCommand(newInitCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newInitCmd() *cobra.Command {
	var claude, cursor, codex bool

	cmd := &cobra.Command{
		Use:   "init [dir]",
		Short: "Scaffold a writing project for a host agent (default: Claude Code)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			switch {
			case cursor:
				fmt.Fprintln(os.Stderr, "pysar init --cursor: not yet supported")
				os.Exit(1)
			case codex:
				fmt.Fprintln(os.Stderr, "pysar init --codex: not yet supported")
				os.Exit(1)
			default:
				scaffoldClaude(dir)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&claude, "claude", false, "scaffold for Claude Code (default)")
	cmd.Flags().BoolVar(&cursor, "cursor", false, "scaffold for Cursor (not yet supported)")
	cmd.Flags().BoolVar(&codex, "codex", false, "scaffold for Codex (not yet supported)")
	cmd.MarkFlagsMutuallyExclusive("claude", "cursor", "codex")

	return cmd
}

func scaffoldClaude(dir string) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "pysar init: "+err.Error())
		os.Exit(1)
	}
	fmt.Printf("pysar init: Claude Code project initialized at %s\n", dir)
}
