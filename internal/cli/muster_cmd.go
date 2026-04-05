package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runMuster(args []string) int {
	var (
		dir      string
		model    string
		crucName string
		athName  string
		intent   string
	)

	positional, flagArgs := splitArgs(args)

	fs := flag.NewFlagSet("muster", flag.ContinueOnError)
	fs.StringVar(&dir, "dir", "", "working directory for the azer (e.g. worktree path)")
	fs.StringVar(&model, "model", "", "model override for the azer")
	fs.StringVar(&crucName, "name", "", "crucible name override")
	fs.StringVar(&athName, "athanor", "", "athanor name (if $ATHANOR not set)")
	fs.StringVar(&intent, "intent", "", "natural language intent (launches autonomous azer)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	// Dispatch based on mode: --intent means intent-driven autonomous azer
	if intent != "" {
		return runMusterIntent(positional, intent, dir, model, crucName, athName)
	}
	return runMusterOpus(positional, dir, model, crucName, athName)
}

// runMusterOpus launches an azer for a pre-inscribed opus file.
// Usage: ath muster <opus-file> [--dir <path>] [--model <model>] [--athanor <name>]
func runMusterOpus(positional []string, dir, model, crucName, athName string) int {
	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: opus file required")
		fmt.Fprintln(os.Stderr, "usage: ath muster <opus-file> [--dir <path>] [--model <model>]")
		fmt.Fprintln(os.Stderr, "       ath muster <mo> <name> --intent <text>")
		return 2
	}
	opusArg := positional[0]

	// Resolve instance directory from $ATHANOR or --athanor flag
	instDir, err := resolveInstDir(athName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitCode(err)
	}

	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Resolve opus file path
	opusPath := opusArg
	if !filepath.IsAbs(opusPath) {
		mos, _ := athanor.ListMagnaOpera(instDir)
		found := false
		for _, mo := range mos {
			candidate := filepath.Join(athanor.OperaPath(instDir, mo), opusArg)
			if _, err := os.Stat(candidate); err == nil {
				opusPath = candidate
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "error: opus file not found: %s\n", opusArg)
			return 1
		}
	} else if _, err := os.Stat(opusPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: opus file not found: %s\n", opusPath)
		return 1
	}

	// Derive crucible name from opus filename if not overridden
	if crucName == "" {
		base := filepath.Base(opusPath)
		base = strings.TrimSuffix(base, ".md")
		if len(base) > 11 && base[4] == '-' && base[7] == '-' && base[10] == '-' {
			base = base[11:]
		}
		crucName = "azer-" + base
	}

	if dir == "" {
		dir = cfg.Project
	}
	if dir == "" {
		h, _ := os.UserHomeDir()
		dir = h
	}

	if model == "" {
		model = cfg.EffectiveAzerModel()
	}

	// Build the azer boot prompt
	moName := athanor.ReadOpusMO(opusPath)
	var bootPrompt string
	if moName != "" {
		moPath := athanor.MagnumOpusPath(instDir, moName)
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s, then read %s/azer.md. Your opus is at %s. Read it and execute.",
			instDir, moPath, instDir, opusPath,
		)
	} else {
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s/azer.md. Your opus is at %s. Read it and execute.",
			instDir, instDir, opusPath,
		)
	}

	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %q %q",
		dir, instDir, model, bootPrompt,
	)

	athSessionName := athanor.SessionName(athanor.AthanorName(instDir))
	if err := launchCrucible(athSessionName, crucName, dir, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Azer mustered in crucible %q\n", crucName)
	fmt.Printf("  Opus: %s\n", opusPath)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", dir)

	return 0
}

// runMusterIntent launches an autonomous azer that inscribes its own opus from intent.
// Usage: ath muster <mo> <name> --intent <text> [--dir <path>] [--model <model>] [--athanor <name>]
func runMusterIntent(positional []string, intent, dir, model, crucName, athName string) int {
	if len(positional) < 2 {
		fmt.Fprintln(os.Stderr, "error: MO name and crucible name required with --intent")
		fmt.Fprintln(os.Stderr, "usage: ath muster <mo> <name> --intent <text>")
		return 2
	}
	moName := positional[0]
	sessionName := positional[1]

	instDir, err := resolveInstDir(athName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitCode(err)
	}

	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Validate the MO exists
	moPath := athanor.MagnumOpusPath(instDir, moName)
	if _, err := os.Stat(moPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: magnum opus not found: %s\n", moPath)
		return 1
	}

	if dir == "" {
		dir = cfg.Project
	}
	if dir == "" {
		h, _ := os.UserHomeDir()
		dir = h
	}

	if model == "" {
		model = cfg.EffectiveAzerModel()
	}

	if crucName == "" {
		crucName = "azer-" + sessionName
	}

	bootPrompt := fmt.Sprintf(
		"Read %s/AGENTS.md, then read %s, then read %s/azer.md. You are an autonomous azer. Your intent: %s. Inscribe an opus for this work under the %s MO using /opus inscribe (dialectical calcinatio), then execute it. Discharge when complete.",
		instDir, moPath, instDir, intent, moName,
	)

	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %q %q",
		dir, instDir, model, bootPrompt,
	)

	athSessionName := athanor.SessionName(athanor.AthanorName(instDir))
	if err := launchCrucible(athSessionName, crucName, dir, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Azer mustered (intent) in crucible %q\n", crucName)
	fmt.Printf("  Intent: %s\n", intent)
	fmt.Printf("  MO: %s\n", moName)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", dir)

	return 0
}

// resolveInstDir resolves the athanor instance directory from $ATHANOR or a name.
func resolveInstDir(athName string) (string, error) {
	instDir := os.Getenv("ATHANOR")
	if athName != "" {
		home, err := athanor.Home()
		if err != nil {
			return "", err
		}
		instDir = athanor.InstanceDir(home, athName)
	}
	if instDir == "" {
		return "", fmt.Errorf("$ATHANOR not set and --athanor not provided")
	}
	return instDir, nil
}

// exitCode returns 2 for usage-like errors, 1 otherwise.
func exitCode(err error) int {
	if strings.Contains(err.Error(), "not set") || strings.Contains(err.Error(), "not provided") {
		return 2
	}
	return 1
}

// launchCrucible creates a tmux window and sends the claude launch command.
func launchCrucible(session, crucName, dir, claudeArgs string) error {
	r := tmux.NewRunner()
	if err := r.EnsureSession(session); err != nil {
		return fmt.Errorf("ensuring tmux session: %w", err)
	}
	if err := r.NewWindow(session, crucName, dir); err != nil {
		return fmt.Errorf("creating crucible: %w", err)
	}
	target := session + ":" + crucName
	if err := r.SendKeysLiteral(target, claudeArgs); err != nil {
		return fmt.Errorf("launching azer: %w", err)
	}
	if err := r.SendKeys(target, "Enter"); err != nil {
		return fmt.Errorf("launching azer: %w", err)
	}
	return nil
}
