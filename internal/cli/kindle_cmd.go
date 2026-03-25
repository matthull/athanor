package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintln(os.Stderr, "usage: ath kindle <name> [<mo-name>]")
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

	// Resolve MO name
	legacy := athanor.HasLegacyMagnumOpus(instDir)
	var moName string
	if len(remaining) >= 2 {
		moName = remaining[1]
	} else if legacy {
		moName = name // legacy: use athanor name
	} else {
		// Multi-MO: mo-name required
		mos, _ := athanor.ListMagnaOpera(instDir)
		if len(mos) == 0 {
			fmt.Fprintln(os.Stderr, "error: no magna opera found — create one in magna-opera/")
			return 1
		}
		fmt.Fprintln(os.Stderr, "error: mo-name required for multi-MO athanor")
		fmt.Fprintf(os.Stderr, "usage: ath kindle %s <mo-name>\n", name)
		fmt.Fprintf(os.Stderr, "available: %s\n", strings.Join(mos, ", "))
		return 2
	}

	// Validate MO has real content
	if err := athanor.ValidateMO(instDir, moName); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	// Crucible naming: legacy = marut-<name>, multi-MO = marut-<name>-<mo>
	var crucible string
	if legacy {
		crucible = athanor.MarutCrucibleName(name, "")
	} else {
		crucible = athanor.MarutCrucibleName(name, moName)
	}

	workDir := cfg.Project
	if workDir == "" {
		h, _ := os.UserHomeDir()
		workDir = h
	}
	model := cfg.EffectiveMarutModel()

	// Build the marut boot prompt — reference specific MO
	moPath := athanor.MagnumOpusPath(instDir, moName)
	bootPrompt := fmt.Sprintf(
		"Read %s/AGENTS.md, then read %s, then read %s/marut.md, then read %s/muster.md. You are the marut for this athanor. Start /loop 5m and begin your operational cycle.",
		instDir, moPath, instDir, instDir,
	)

	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %s --permission-mode auto %q",
		workDir, instDir, model, bootPrompt,
	)

	r := tmux.NewRunner()
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

	fmt.Printf("Marut kindled for %q in crucible %q\n", moName, crucible)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", workDir)
	fmt.Printf("  Instance: %s\n", instDir)

	return 0
}
