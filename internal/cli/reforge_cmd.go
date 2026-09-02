package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	session := athanor.SessionName(name)

	// Rename old crucible to avoid name collision with the new one.
	// This must happen before creating the new window (same name).
	// Idempotent — no-op if the old window doesn't exist.
	oldTarget := session + ":" + crucible
	dyingName := crucible + "-dying"
	_ = r.RenameWindow(oldTarget, dyingName)

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
		"%scd %s && ATHANOR=%s claude --permission-mode auto --model %q %q",
		envSourcePrefix(instDir), workDir, instDir, model, reforgePrompt,
	)

	if err := r.EnsureSession(session); err != nil {
		fmt.Fprintf(os.Stderr, "error ensuring tmux session: %v\n", err)
		return 1
	}

	if err := r.NewWindow(session, crucible, workDir); err != nil {
		fmt.Fprintf(os.Stderr, "error creating crucible: %v\n", err)
		return 1
	}

	target := session + ":" + crucible
	if err := r.SendKeysLiteral(target, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error launching marut: %v\n", err)
		return 1
	}
	if err := r.SendKeys(target, "Enter"); err != nil {
		fmt.Fprintf(os.Stderr, "error launching marut: %v\n", err)
		return 1
	}

	// Kill old crucible last. Even if this propagates SIGHUP to our caller
	// (when reforge is invoked from within the old marut window), all
	// important work — new window created, claude launched — is already done.
	_ = r.KillWindow(session + ":" + dyingName)

	fmt.Printf("Marut reforged for %q in crucible %q\n", moName, crucible)

	return 0
}
