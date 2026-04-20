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

func runKindle(args []string) int {
	var role, moFlag string

	positional, flagArgs := splitArgs(args)

	fs := flag.NewFlagSet("kindle", flag.ContinueOnError)
	fs.StringVar(&role, "role", "marut", "presence-driven role to kindle (marut, perceiver, attendant)")
	fs.StringVar(&moFlag, "mo", "", "magnum opus name (alternative to positional arg)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: athanor name required")
		fmt.Fprintln(os.Stderr, "usage: ath kindle <name> [<mo-name>] [--role <role>]")
		return 2
	}
	name := positional[0]

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

	// Resolve MO name — positional takes precedence, then --mo, then legacy fallback.
	legacy := athanor.HasLegacyMagnumOpus(instDir)
	var moName string
	switch {
	case len(positional) >= 2:
		moName = positional[1]
	case moFlag != "":
		moName = moFlag
	case legacy:
		moName = name
	default:
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

	if err := athanor.ValidateMO(instDir, moName); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	if _, err := os.Stat(filepath.Join(instDir, role+".md")); err != nil {
		fmt.Fprintf(os.Stderr, "error: role file %s.md not found in %s\n", role, instDir)
		return 1
	}

	var crucible string
	if legacy {
		crucible = athanor.RoleCrucibleName(role, name, "")
	} else {
		crucible = athanor.RoleCrucibleName(role, name, moName)
	}

	r := tmux.NewRunner()
	session := athanor.SessionName(name)

	// Idempotency — the liveness timer hits this path every tick.
	windows, _ := r.ListSessionWindows(session)
	for _, w := range windows {
		if w == crucible {
			fmt.Printf("%s already running for %q in crucible %q\n", titleCase(role), moName, crucible)
			return 0
		}
	}

	if err := athanor.SetKindled(instDir, moName, role); err != nil {
		fmt.Fprintf(os.Stderr, "error setting kindled state: %v\n", err)
		return 1
	}

	workDir := cfg.Project
	if workDir == "" {
		h, _ := os.UserHomeDir()
		workDir = h
	}
	// All presence-driven roles (marut, perceiver, attendant) use the marut model.
	// They are lightweight singletons — sonnet is appropriate. Craft work (azer) uses opus.
	model := cfg.EffectiveMarutModel()

	moPath := athanor.MagnumOpusPath(instDir, moName)
	var bootPrompt string
	if role == "marut" {
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s, then read %s/marut.md, then read %s/muster.md. You are the marut for this athanor. Start /loop 5m and begin your operational cycle.",
			instDir, moPath, instDir, instDir,
		)
	} else {
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s, then read %s/%s.md. You are the %s for this athanor. Prepare your station.",
			instDir, moPath, instDir, role, role,
		)
	}

	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %q %q",
		workDir, instDir, model, bootPrompt,
	)

	if err := launchCrucible(session, crucible, workDir, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("%s kindled for %q in crucible %q\n", titleCase(role), moName, crucible)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", workDir)
	fmt.Printf("  Instance: %s\n", instDir)

	return 0
}

// titleCase uppercases the first rune of s. Sufficient for role names
// which are always lowercase ASCII (marut, perceiver, attendant).
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
