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

	positional, flagArgs := splitArgs(args)

	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: athanor name required")
		fmt.Fprintln(os.Stderr, "usage: ath quiesce <name> [<mo-name>] [--force]")
		return 2
	}
	name := positional[0]
	var moName string
	if len(positional) >= 2 {
		moName = positional[1]
	}

	fs := flag.NewFlagSet("quiesce", flag.ContinueOnError)
	fs.BoolVar(&force, "force", false, "force shutdown even with active azers")
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

	// Check for active azers
	windows, _ := r.ListWindows()
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

	// Kill marut crucible(s)
	if moName != "" {
		// Kill specific MO's marut
		crucible := athanor.MarutCrucibleName(name, moName)
		_ = r.KillWindow(crucible)
		fmt.Printf("Marut for %q quiesced.\n", moName)
	} else {
		// Kill all maruts for this athanor (both legacy and multi-MO patterns)
		_ = r.KillWindow(athanor.MarutCrucibleName(name, "")) // legacy
		mos, _ := athanor.ListMagnaOpera(instDir)
		for _, mo := range mos {
			_ = r.KillWindow(athanor.MarutCrucibleName(name, mo))
		}
	}

	// Kill azers if forced
	if force {
		for _, a := range activeAzers {
			_ = r.KillWindow(a)
		}
	}

	// Count discharged opera for summary
	_, discharged := countOpera(instDir)
	fmt.Printf("Athanor %q quiesced. Trail: %d discharged opera.\n", name, discharged)

	return 0
}
