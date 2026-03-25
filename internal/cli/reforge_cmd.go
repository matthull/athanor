package cli

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runReforge(args []string) int {
	fs := flag.NewFlagSet("reforge", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()

	if len(remaining) < 1 {
		fmt.Fprintln(os.Stderr, "error: athanor name required")
		fmt.Fprintln(os.Stderr, "usage: ath reforge <name>")
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

	crucible := fmt.Sprintf("marut-%s", name)
	r := tmux.NewRunner()

	// Kill old session: send Ctrl-C then exit
	_ = r.SendKeys(crucible, "C-c")
	time.Sleep(500 * time.Millisecond)
	_ = r.SendKeysLiteral(crucible, "exit")
	_ = r.SendKeys(crucible, "Enter")
	time.Sleep(1 * time.Second)

	// Kill the window if still alive and recreate
	_ = r.KillWindow(crucible)

	workDir := cfg.Project
	if workDir == "" {
		h, _ := os.UserHomeDir()
		workDir = h
	}
	model := cfg.EffectiveMarutModel()

	// Build the reforge prompt
	reforgePrompt := fmt.Sprintf(
		"Read %s/AGENTS.md, then read %s/magnum-opus.md, then read %s/marut.md, then read %s/muster.md. You are the marut for this athanor. Check opera/ for in-progress work and the trail. Start /loop 5m and resume your operational cycle.",
		instDir, instDir, instDir, instDir,
	)

	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %s --permission-mode auto %q",
		workDir, instDir, model, reforgePrompt,
	)

	if err := r.NewWindow(crucible, workDir); err != nil {
		fmt.Fprintf(os.Stderr, "error creating crucible: %v\n", err)
		return 1
	}

	if err := r.SendKeysLiteral(crucible, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error launching marut: %v\n", err)
		return 1
	}
	if err := r.SendKeys(crucible, "Enter"); err != nil {
		fmt.Fprintf(os.Stderr, "error launching marut: %v\n", err)
		return 1
	}

	fmt.Printf("Marut reforged in crucible %q\n", crucible)

	return 0
}
