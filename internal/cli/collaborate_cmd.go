package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/matthull/athanor/internal/athanor"
)

type collaborateArgs struct {
	moName  string
	intent  string
	job     string
	formula string
}

func parseCollaborateArgs(args []string) (*collaborateArgs, error) {
	positional, flagArgs := splitArgs(args)

	var ca collaborateArgs

	fs := flag.NewFlagSet("collaborate", flag.ContinueOnError)
	fs.StringVar(&ca.intent, "intent", "", "natural language intent for the opus (required)")
	fs.StringVar(&ca.job, "job", "", "job role for the peer azer (required)")
	fs.StringVar(&ca.formula, "formula", "", "collaboration formula for the peer azer's own work (optional)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return nil, err
	}

	if len(positional) < 1 {
		return nil, fmt.Errorf("MO name required")
	}
	ca.moName = positional[0]

	if ca.intent == "" {
		return nil, fmt.Errorf("--intent is required")
	}

	if ca.job == "" {
		return nil, fmt.Errorf("--job is required (use \"general\" if no specific role fits)")
	}

	if err := athanor.ValidateJob(ca.job); err != nil {
		return nil, err
	}

	return &ca, nil
}

// callingCrucible returns the tmux window name of the current pane,
// or empty string if not running inside tmux.
func callingCrucible() string {
	out, err := exec.Command("tmux", "display-message", "-p", "#{window_name}").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func runCollaborate(args []string) int {
	ca, err := parseCollaborateArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, "usage: ath collaborate <mo> --job <job> --intent <text>")
		return 2
	}

	// Infer athanor from $ATHANOR
	instDir := os.Getenv("ATHANOR")
	if instDir == "" {
		fmt.Fprintln(os.Stderr, "error: $ATHANOR not set (collaborate must be called from a crucible)")
		return 2
	}

	// Validate MO exists
	moPath := athanor.MagnumOpusPath(instDir, ca.moName)
	if _, err := os.Stat(moPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: magnum opus not found: %s\n", moPath)
		return 1
	}

	// Validate formula (optional) — formulae live in the athanor instance
	if ca.formula != "" {
		if err := athanor.ValidateFormula(instDir, ca.formula); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	}

	// Build opus filename
	datestamp := time.Now().Format("2006-01-02")
	slug := slugify(ca.intent, maxSlugLen)
	if slug == "" {
		fmt.Fprintln(os.Stderr, "error: could not generate opus name from intent")
		return 1
	}
	opusFilename := fmt.Sprintf("%s-%s.md", datestamp, slug)

	// Ensure opera directory exists
	operaDir := athanor.OperaPath(instDir, ca.moName)
	if err := os.MkdirAll(operaDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating opera dir: %v\n", err)
		return 1
	}
	opusPath := filepath.Join(operaDir, opusFilename)

	// Build collaboration context
	caller := callingCrucible()
	var collabContext string
	if caller != "" {
		collabContext = fmt.Sprintf("Inscribed by: %s. Whisper back to %s when complete.", caller, caller)
	}

	// Build opus content
	content := buildOpusContent(datestamp, ca.moName, ca.job, ca.formula, ca.intent, collabContext)

	if err := os.WriteFile(opusPath, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing opus: %v\n", err)
		return 1
	}

	// Always muster
	return musterInscribedOpus(instDir, opusPath, slug)
}
