package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runCraftMO(args []string) int {
	var (
		dir   string
		model string
	)

	positional, flagArgs := splitArgs(args)

	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: athanor name required")
		fmt.Fprintln(os.Stderr, "usage: ath craft-mo <athanor> [--dir <path>] [--model <model>]")
		return 2
	}
	athName := positional[0]

	fs := flag.NewFlagSet("craft-mo", flag.ContinueOnError)
	fs.StringVar(&dir, "dir", "", "working directory")
	fs.StringVar(&model, "model", "", "model override")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	instDir := athanor.InstanceDir(home, athName)
	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
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

	bootPrompt := fmt.Sprintf(
		"Read %s/AGENTS.md, then read %s/azer.md. You are an interactive azer working alongside the artifex to create a new Magnum Opus for this athanor. Use /opus to inscribe the MO collaboratively — the artifex will guide the shaping. When done, the artifex will signal discharge.",
		instDir, instDir,
	)

	crucName := "azer-craft-mo"
	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %s %q",
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

	fmt.Printf("Craft-MO session started in crucible %q\n", crucName)
	fmt.Printf("  Athanor: %s\n", athName)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", dir)

	return 0
}
