// Package cli provides the command-line interface for the athanor system.
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
	case "whisper":
		return runWhisper(os.Args[2:])
	case "init":
		return runInit(os.Args[2:])
	case "kindle":
		return runKindle(os.Args[2:])
	case "muster":
		return runMuster(os.Args[2:])
	case "status":
		return runStatus(os.Args[2:])
	case "reforge":
		return runReforge(os.Args[2:])
	case "cleanup":
		return runCleanup(os.Args[2:])
	case "quiesce":
		return runQuiesce(os.Args[2:])
	case "opera":
		return runOpera(os.Args[2:])
	case "completion":
		return runCompletion(os.Args[2:])
	case "version":
		fmt.Printf("ath %s (commit: %s, built: %s)\n", Version, Commit, BuildTime)
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
	fmt.Fprintf(os.Stderr, `ath — athanor agent orchestration CLI

Usage:
  ath init <name> [--project <path>]     Create a new athanor instance
  ath kindle <name>                       Launch a marut for an athanor
  ath reforge <name>                      Kill and relaunch a marut
  ath muster <opus-file> [--dir <path>]   Launch an azer for an opus
  ath cleanup <crucible>                  Clean up after a discharged opus
  ath quiesce <name>                      Graceful shutdown of an athanor
  ath status [<name>]                     Show athanor health
  ath opera [<name>]                      List opera with status

  ath whisper send <target> <message>     Send a message to a tmux target
  ath whisper idle <target>               Wait for target to become idle
  ath whisper wait-and-send <target> <msg> Wait for idle, then send

  ath completion zsh                      Generate zsh completion script
  ath version                             Print version info

`)
}
