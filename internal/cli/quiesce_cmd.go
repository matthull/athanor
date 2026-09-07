package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runQuiesce(args []string) int {
	var force bool
	var role string

	positional, flagArgs := splitArgs(args)

	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: athanor name required")
		fmt.Fprintln(os.Stderr, "usage: ath quiesce <name> [<mo-name>] [--role <role>] [--force]")
		return 2
	}
	name := positional[0]
	var moName string
	if len(positional) >= 2 {
		moName = positional[1]
	}

	fs := flag.NewFlagSet("quiesce", flag.ContinueOnError)
	fs.BoolVar(&force, "force", false, "force shutdown even with active azers")
	fs.StringVar(&role, "role", "", "quiesce a specific kindled role")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	instDir := athanor.InstanceDir(home, name)
	if _, err := athanor.ReadConfig(instDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	r := tmux.NewRunner()
	session := athanor.SessionName(name)

	// Single-role quiesce path
	if role != "" {
		if moName == "" {
			mos, _ := athanor.ListMagnaOpera(instDir)
			switch len(mos) {
			case 0:
				fmt.Fprintln(os.Stderr, "error: no magna opera found")
				return 1
			case 1:
				moName = mos[0]
			default:
				fmt.Fprintln(os.Stderr, "error: mo-name required with --role for multi-MO athanor")
				return 2
			}
		}

		if err := athanor.ClearKindled(instDir, moName, role); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		crucible := athanor.RoleCrucibleName(role, name, moName)
		_ = r.KillWindow(session + ":" + crucible)
		fmt.Printf("%s for %q quiesced.\n", capitalize(role), moName)
		return 0
	}

	// Check for active azers in this athanor's session
	windows, _ := r.ListSessionWindows(session)
	var activeAzers []string
	for _, w := range windows {
		if strings.HasPrefix(w, "azer-") {
			activeAzers = append(activeAzers, w)
		}
	}

	if len(activeAzers) > 0 && !force {
		fmt.Fprintf(os.Stderr, "error: %d active azer(s): %s\n", len(activeAzers), strings.Join(activeAzers, ", "))
		fmt.Fprintln(os.Stderr, "Use --force to override.")
		return 2
	}

	// Clear all kindled state for the targeted MO(s)
	if moName != "" {
		_ = athanor.ClearAllKindled(instDir, moName)
	} else {
		mos, _ := athanor.ListMagnaOpera(instDir)
		for _, mo := range mos {
			_ = athanor.ClearAllKindled(instDir, mo)
		}
	}

	// Kill marut crucible(s)
	if moName != "" {
		crucible := athanor.MarutCrucibleName(name, moName)
		_ = r.KillWindow(session + ":" + crucible)
		fmt.Printf("Marut for %q quiesced.\n", moName)
	} else {
		// Kill all maruts for this athanor (both legacy and multi-MO patterns)
		_ = r.KillWindow(session + ":" + athanor.MarutCrucibleName(name, "")) // legacy
		mos, _ := athanor.ListMagnaOpera(instDir)
		for _, mo := range mos {
			_ = r.KillWindow(session + ":" + athanor.MarutCrucibleName(name, mo))
		}
	}

	// Kill role-specific windows (perceiver-*, attendant-*, etc.) for this athanor
	if moName != "" {
		killRoleWindows(r, session, windows, name, moName)
	} else {
		mos, _ := athanor.ListMagnaOpera(instDir)
		for _, w := range windows {
			// Match <role>-<name> (legacy) or <role>-<name>-<mo>
			for _, mo := range mos {
				suffix := "-" + name + "-" + mo
				if strings.HasSuffix(w, suffix) && !strings.HasPrefix(w, "marut-") && !strings.HasPrefix(w, "azer-") {
					_ = r.KillWindow(session + ":" + w)
					break
				}
			}
			legacySuffix := "-" + name
			if strings.HasSuffix(w, legacySuffix) && !strings.HasPrefix(w, "marut-") && !strings.HasPrefix(w, "azer-") {
				_ = r.KillWindow(session + ":" + w)
			}
		}
	}

	// Kill azers if forced
	if force {
		for _, a := range activeAzers {
			_ = r.KillWindow(session + ":" + a)
		}
	}

	// Count discharged opera for summary
	_, discharged := countOpera(instDir)
	fmt.Printf("Athanor %q quiesced. Trail: %d discharged opera.\n", name, discharged)

	return 0
}

// capitalize reuses titleCase from kindle_cmd.go — keeping this as a local alias
// for readability in quiesce output messages.
var capitalize = titleCase

func killRoleWindows(r *tmux.Runner, session string, windows []string, name, moName string) {
	suffix := "-" + name + "-" + moName
	for _, w := range windows {
		if !strings.HasSuffix(w, suffix) {
			continue
		}
		if strings.HasPrefix(w, "marut-") || strings.HasPrefix(w, "azer-") {
			continue
		}
		_ = r.KillWindow(session + ":" + w)
	}
}
