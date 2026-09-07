package cli

import (
	"fmt"
	"os"
)

// runWhisper routes whisper subcommands (send, idle, wait-and-send).
func runWhisper(args []string) int {
	if len(args) < 1 {
		printWhisperUsage()
		return 2
	}

	switch args[0] {
	case "send":
		return runSend(args[1:])
	case "idle":
		return runIdle(args[1:])
	case "wait-and-send":
		return runWaitAndSend(args[1:])
	case "--help", "-h", "help":
		printWhisperUsage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown whisper command: %s\n", args[0])
		printWhisperUsage()
		return 2
	}
}

func printWhisperUsage() {
	fmt.Fprintf(os.Stderr, `ath whisper — reliable message delivery to tmux sessions

Usage:
  ath whisper send <target> <message>       Send a message to a tmux target
  ath whisper send <target> -f <file>       Send file contents to a tmux target
  ath whisper send --self <message>         Send to own pane ($TMUX_PANE)
  ath whisper idle <target>                 Wait for target to become idle
  ath whisper wait-and-send <target> <msg>  Wait for idle, then send

Options:
  --skip-escape     Omit Escape keystroke (for non-Claude agents)
  --timeout <dur>   Max wait time (default: 15s)

`)
}
