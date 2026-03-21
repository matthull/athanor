package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/matthull/whisper/internal/tmux"
)

// sendArgs holds parsed send command arguments.
type sendArgs struct {
	target  string
	message string
	opts    tmux.SendOpts
}

// parseSendArgs parses send command arguments and returns the parsed result.
// Returns a non-zero exitCode on error.
func parseSendArgs(args []string) (result sendArgs, exitCode int, err error) {
	var (
		self       bool
		filePath   string
		skipEscape bool
		timeout    time.Duration
	)

	fs := flag.NewFlagSet("send", flag.ContinueOnError)
	fs.BoolVar(&self, "self", false, "send to own pane ($TMUX_PANE)")
	fs.StringVar(&filePath, "f", "", "read message from file")
	fs.BoolVar(&skipEscape, "skip-escape", false, "omit Escape keystroke")
	fs.DurationVar(&timeout, "timeout", 15*time.Second, "max retry timeout")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return sendArgs{}, 2, err
	}
	remaining := fs.Args()

	// Resolve target
	var target string
	if self {
		target = os.Getenv("TMUX_PANE")
		if target == "" {
			return sendArgs{}, 2, fmt.Errorf("--self requires $TMUX_PANE (must be inside tmux)")
		}
	} else {
		if len(remaining) < 1 {
			return sendArgs{}, 2, fmt.Errorf("target required")
		}
		target = remaining[0]
		remaining = remaining[1:]
	}

	// Resolve message
	var message string
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return sendArgs{}, 1, fmt.Errorf("reading file: %w", err)
		}
		message = string(data)
	} else {
		if len(remaining) < 1 {
			return sendArgs{}, 2, fmt.Errorf("message required (or use -f <file>)")
		}
		message = strings.Join(remaining, " ")
	}

	return sendArgs{
		target:  target,
		message: message,
		opts: tmux.SendOpts{
			SkipEscape: skipEscape,
			Timeout:    timeout,
		},
	}, 0, nil
}

func runSend(args []string) int {
	parsed, exitCode, err := parseSendArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitCode
	}

	r := tmux.NewRunner()
	if err := r.Send(parsed.target, parsed.message, parsed.opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}
