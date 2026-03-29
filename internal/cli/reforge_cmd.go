package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
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
		fmt.Fprintln(os.Stderr, "usage: ath reforge <name> [<mo-name>]")
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

	// Resolve MO name (same logic as kindle)
	legacy := athanor.HasLegacyMagnumOpus(instDir)
	var moName string
	if len(remaining) >= 2 {
		moName = remaining[1]
	} else if legacy {
		moName = name
	} else {
		mos, _ := athanor.ListMagnaOpera(instDir)
		if len(mos) == 0 {
			fmt.Fprintln(os.Stderr, "error: no magna opera found")
			return 1
		}
		fmt.Fprintln(os.Stderr, "error: mo-name required for multi-MO athanor")
		fmt.Fprintf(os.Stderr, "usage: ath reforge %s <mo-name>\n", name)
		fmt.Fprintf(os.Stderr, "available: %s\n", strings.Join(mos, ", "))
		return 2
	}

	var crucible string
	if legacy {
		crucible = athanor.MarutCrucibleName(name, "")
	} else {
		crucible = athanor.MarutCrucibleName(name, moName)
	}

	r := tmux.NewRunner()

	// Kill old session
	_ = r.SendKeys(crucible, "C-c")
	time.Sleep(500 * time.Millisecond)
	_ = r.SendKeysLiteral(crucible, "exit")
	_ = r.SendKeys(crucible, "Enter")
	time.Sleep(1 * time.Second)
	_ = r.KillWindow(crucible)

	workDir := cfg.Project
	if workDir == "" {
		h, _ := os.UserHomeDir()
		workDir = h
	}
	model := cfg.EffectiveMarutModel()

	moPath := athanor.MagnumOpusPath(instDir, moName)
	operaDir := athanor.OperaPath(instDir, moName)
	reforgePrompt := fmt.Sprintf(
		"Read %s/AGENTS.md, then read %s, then read %s/marut.md, then read %s/muster.md. You are the marut for this athanor. Check %s for in-progress work and the trail. Start /loop 5m and resume your operational cycle.",
		instDir, moPath, instDir, instDir, operaDir,
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

	fmt.Printf("Marut reforged for %q in crucible %q\n", moName, crucible)

	return 0
}
