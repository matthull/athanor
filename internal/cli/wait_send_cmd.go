package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/matthull/whisper/internal/tmux"
)

func runWaitAndSend(args []string) int {
	var (
		filePath   string
		skipEscape bool
		timeout    time.Duration
	)

	fs := flag.NewFlagSet("wait-and-send", flag.ContinueOnError)
	fs.StringVar(&filePath, "f", "", "read message from file")
	fs.BoolVar(&skipEscape, "skip-escape", false, "omit Escape keystroke")
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
	remaining = remaining[1:]

	// Resolve message
	var message string
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
			return 1
		}
		message = string(data)
	} else {
		if len(remaining) < 1 {
			fmt.Fprintln(os.Stderr, "error: message required (or use -f <file>)")
			return 2
		}
		message = strings.Join(remaining, " ")
	}

	// Acquire lock BEFORE idle check to prevent race between idle check and send.
	// The spec says: "holds the send lock for the entire operation."
	r := tmux.NewRunner()

	if !tmux.AcquireLock(target, tmux.LockTimeout) {
		fmt.Fprintf(os.Stderr, "error: send lock timeout for %q: previous send may be hung\n", target)
		return 1
	}
	defer tmux.ReleaseLock(target)

	// Wait for idle
	if err := tmux.WaitForIdle(r, target, timeout); err != nil {
		if errors.Is(err, tmux.ErrIdleTimeout) {
			fmt.Fprintln(os.Stderr, "timeout waiting for idle")
			return 1
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if errors.Is(err, tmux.ErrSessionNotFound) || errors.Is(err, tmux.ErrNoServer) {
			return 2
		}
		return 1
	}

	// Send with LockHeld=true since we already hold the lock
	opts := tmux.SendOpts{
		SkipEscape: skipEscape,
		Timeout:    timeout,
		LockHeld:   true,
	}
	if err := r.Send(target, message, opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}
