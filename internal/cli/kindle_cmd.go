package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runKindle(args []string) int {
	fs := flag.NewFlagSet("kindle", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()

	if len(remaining) < 1 {
		fmt.Fprintln(os.Stderr, "error: athanor name required")
		fmt.Fprintln(os.Stderr, "usage: ath kindle <name>")
		return 2
	}
	name := remaining[0]

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	instDir := athanor.InstanceDir(home, name)
	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Validate magnum-opus.md has real content
	if err := athanor.ValidateMagnumOpus(instDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		// Warn but don't block — artifex may be testing
	}

	crucible := fmt.Sprintf("marut-%s", name)
	workDir := cfg.Project
	if workDir == "" {
		h, _ := os.UserHomeDir()
		workDir = h
	}
	model := cfg.EffectiveMarutModel()

	// Build the marut boot prompt
	bootPrompt := fmt.Sprintf(
		"Read %s/AGENTS.md, then read %s/magnum-opus.md, then read %s/marut.md, then read %s/muster.md. You are the marut for this athanor. Start /loop 5m and begin your operational cycle.",
		instDir, instDir, instDir, instDir,
	)

	// Build the claude command
	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %s --permission-mode auto %q",
		workDir, instDir, model, bootPrompt,
	)

	// Create tmux window and launch
	r := tmux.NewRunner()
	if err := r.NewWindow(crucible, workDir); err != nil {
		fmt.Fprintf(os.Stderr, "error creating crucible: %v\n", err)
		return 1
	}

	// Send the claude launch command to the new window
	if err := r.SendKeysLiteral(crucible, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error launching marut: %v\n", err)
		return 1
	}
	if err := r.SendKeys(crucible, "Enter"); err != nil {
		fmt.Fprintf(os.Stderr, "error launching marut: %v\n", err)
		return 1
	}

	fmt.Printf("Marut kindled in crucible %q\n", crucible)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", workDir)
	fmt.Printf("  Instance: %s\n", instDir)

	return 0
}
