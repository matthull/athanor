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
	case "inscribe":
		return runInscribe(os.Args[2:])
	case "collaborate":
		return runCollaborate(os.Args[2:])
	case "craft":
		return runCraft(os.Args[2:])
	case "craft-mo":
		return runCraftMO(os.Args[2:])
	case "status":
		return runStatus(os.Args[2:])
	case "reforge":
		return runReforge(os.Args[2:])
	case "check":
		return runCheck(os.Args[2:])
	case "cleanup":
		return runCleanup(os.Args[2:])
	case "quiesce":
		return runQuiesce(os.Args[2:])
	case "view":
		return runView(os.Args[2:])
	case "opera":
		return runOpera(os.Args[2:])
	case "patrol":
		return runPatrol(os.Args[2:])
	case "dashboard":
		return runDashboard(os.Args[2:])
	case "services":
		return runServices(os.Args[2:])
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
  ath init <name> [--project <path>]        Create a new athanor instance
  ath kindle <name> [<mo-name>] [--role <r>] Launch a presence-driven role
  ath reforge <name> [<mo-name>]           Kill and relaunch a marut
  ath inscribe <athanor> <mo> --intent <text> [--job <job>] [--muster]
                                           Create an opus file
  ath collaborate <mo> --intent <text> [--job <job>]
                                           Inscribe + muster a peer azer (from crucible)
  ath muster <opus-file> [--worktree-path <path>]    Launch an azer for an opus
  ath muster <mo> <name> --intent <text>   Launch autonomous azer from intent
  ath craft <athanor> <mo> <name>           Interactive session with the artifex
  ath craft-mo <athanor>                   Create a new Magnum Opus interactively
  ath check <crucible>                      Check crucible health
  ath cleanup <crucible>                   Clean up after a discharged opus
  ath quiesce <name> [<mo-name>] [--role <r>] [--force] Graceful shutdown
  ath status [<name>]                      Show athanor health
  ath view <athanor> <mo> [<opus>]          Open MO or opus in $EDITOR
  ath opera [<name>] [--mo <mo-name>]      List opera with status
  ath patrol [--json] [--exclude <pane>]  Scan panes for prompts/stalls
  ath dashboard [--watch] [--json]        At-a-glance system overview
  ath services [--json]                    Check athanor service dependencies

  ath whisper send <target> <message>     Send a message to a tmux target
  ath whisper idle <target>               Wait for target to become idle
  ath whisper wait-and-send <target> <msg> Wait for idle, then send

  ath completion zsh                      Generate zsh completion script
  ath version                             Print version info

`)
}
