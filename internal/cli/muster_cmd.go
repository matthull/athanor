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
	)

	positional, flagArgs := splitArgs(args)

	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: opus file required")
		fmt.Fprintln(os.Stderr, "usage: ath muster <opus-file> [--dir <path>] [--model <model>]")
		return 2
	}
	opusArg := positional[0]

	fs := flag.NewFlagSet("muster", flag.ContinueOnError)
	fs.StringVar(&dir, "dir", "", "working directory for the azer (e.g. worktree path)")
	fs.StringVar(&model, "model", "", "model override for the azer")
	fs.StringVar(&crucName, "name", "", "crucible name override")
	fs.StringVar(&athName, "athanor", "", "athanor name (if $ATHANOR not set)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	// Resolve instance directory from $ATHANOR or --athanor flag
	instDir := os.Getenv("ATHANOR")
	if athName != "" {
		home, err := athanor.Home()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		instDir = athanor.InstanceDir(home, athName)
	}
	if instDir == "" {
		fmt.Fprintln(os.Stderr, "error: $ATHANOR not set and --athanor not provided")
		return 2
	}

	// Read instance config
	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Resolve opus file path
	opusPath := opusArg
	if !filepath.IsAbs(opusPath) {
		// Search across all MO opera directories for the opus file
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
		// Strip date prefix (YYYY-MM-DD-)
		if len(base) > 11 && base[4] == '-' && base[7] == '-' && base[10] == '-' {
			base = base[11:]
		}
		crucName = "azer-" + base
	}

	// Resolve working directory
	if dir == "" {
		dir = cfg.Project
	}
	if dir == "" {
		h, _ := os.UserHomeDir()
		dir = h
	}

	// Resolve model
	if model == "" {
		model = cfg.EffectiveAzerModel()
	}

	// Build the azer boot prompt — include MO context if available
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
		"cd %s && ATHANOR=%s claude --model %s --permission-mode auto %q",
		dir, instDir, model, bootPrompt,
	)

	// Create tmux window and launch
	r := tmux.NewRunner()
	if err := r.NewWindow(crucName, dir); err != nil {
		fmt.Fprintf(os.Stderr, "error creating crucible: %v\n", err)
		return 1
	}

	if err := r.SendKeysLiteral(crucName, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error launching azer: %v\n", err)
		return 1
	}
	if err := r.SendKeys(crucName, "Enter"); err != nil {
		fmt.Fprintf(os.Stderr, "error launching azer: %v\n", err)
		return 1
	}

	fmt.Printf("Azer mustered in crucible %q\n", crucName)
	fmt.Printf("  Opus: %s\n", opusPath)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", dir)

	return 0
}
