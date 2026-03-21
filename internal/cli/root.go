// Package cli provides the command-line interface for whisper.
package cli

import (
	"fmt"
	"os"
)

// Build info, set via ldflags.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// Execute runs the root command and returns an exit code.
func Execute() int {
	if len(os.Args) < 2 {
		printUsage()
		return 2
	}

	switch os.Args[1] {
	case "send":
		return runSend(os.Args[2:])
	case "idle":
		return runIdle(os.Args[2:])
	case "wait-and-send":
		return runWaitAndSend(os.Args[2:])
	case "version":
		fmt.Printf("whisper %s (commit: %s, built: %s)\n", Version, Commit, BuildTime)
		return 0
	case "--help", "-h", "help":
		printUsage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		return 2
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `whisper — reliable message delivery to tmux sessions

Usage:
  whisper send <target> <message>       Send a message to a tmux target
  whisper send <target> -f <file>       Send file contents to a tmux target
  whisper send --self <message>         Send to own pane ($TMUX_PANE)
  whisper idle <target>                 Wait for target to become idle
  whisper wait-and-send <target> <msg>  Wait for idle, then send
  whisper version                       Print version info

Options:
  --skip-escape     Omit Escape keystroke (for non-Claude agents)
  --timeout <dur>   Max wait time (default: 15s)

`)
}

// Stubs — will be implemented in separate files.

func runSend(args []string) int {
	fmt.Fprintln(os.Stderr, "send: not yet implemented")
	return 1
}

func runIdle(args []string) int {
	fmt.Fprintln(os.Stderr, "idle: not yet implemented")
	return 1
}

func runWaitAndSend(args []string) int {
	fmt.Fprintln(os.Stderr, "wait-and-send: not yet implemented")
	return 1
}
