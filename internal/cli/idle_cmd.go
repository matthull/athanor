package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/matthull/athanor/internal/tmux"
)

func runIdle(args []string) int {
	var timeout time.Duration

	fs := flag.NewFlagSet("idle", flag.ContinueOnError)
	fs.DurationVar(&timeout, "timeout", 15*time.Second, "max wait time")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()

	if len(remaining) < 1 {
		fmt.Fprintln(os.Stderr, "error: target required")
		return 2
	}
	target := remaining[0]

	r := tmux.NewRunner()
	err := tmux.WaitForIdle(r, target, timeout)
	if err == nil {
		return 0
	}

	if errors.Is(err, tmux.ErrIdleTimeout) {
		fmt.Fprintln(os.Stderr, "timeout")
		return 1
	}
	if errors.Is(err, tmux.ErrSessionNotFound) || errors.Is(err, tmux.ErrNoServer) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	return 1
}
