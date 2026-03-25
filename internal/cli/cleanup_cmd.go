package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/tmux"
)

func runCleanup(args []string) int {
	fs := flag.NewFlagSet("cleanup", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()

	if len(remaining) < 1 {
		fmt.Fprintln(os.Stderr, "error: crucible name required")
		fmt.Fprintln(os.Stderr, "usage: ath cleanup <crucible-name>")
		return 2
	}
	crucible := remaining[0]

	r := tmux.NewRunner()
	if err := r.KillWindow(crucible); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Crucible %q cleaned up.\n", crucible)
	return 0
}
