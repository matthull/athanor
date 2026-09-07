package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runCleanup(args []string) int {
	var athName string

	positional, flagArgs := splitArgs(args)

	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: crucible name required")
		fmt.Fprintln(os.Stderr, "usage: ath cleanup <crucible-name> [--athanor <name>]")
		return 2
	}
	crucible := positional[0]

	fs := flag.NewFlagSet("cleanup", flag.ContinueOnError)
	fs.StringVar(&athName, "athanor", "", "athanor name (if $ATHANOR not set)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	// Derive athanor name from flag or $ATHANOR
	if athName == "" {
		if instDir := os.Getenv("ATHANOR"); instDir != "" {
			athName = athanor.AthanorName(instDir)
		}
	}
	if athName == "" {
		fmt.Fprintln(os.Stderr, "error: $ATHANOR not set and --athanor not provided")
		return 2
	}

	session := athanor.SessionName(athName)
	target := session + ":" + crucible

	r := tmux.NewRunner()
	if err := r.KillWindow(target); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Crucible %q cleaned up.\n", crucible)
	return 0
}
