package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matthull/athanor/internal/athanor"
)

func runOpera(args []string) int {
	var athName string
	var moFilter string

	positional, flagArgs := splitArgs(args)

	fs := flag.NewFlagSet("opera", flag.ContinueOnError)
	fs.StringVar(&athName, "athanor", "", "athanor name (if $ATHANOR not set)")
	fs.StringVar(&moFilter, "mo", "", "filter by magnum opus name")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	// Resolve instance directory
	var instDir string
	if len(positional) > 0 {
		athName = positional[0]
	}
	if athName != "" {
		home, err := athanor.Home()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		instDir = athanor.InstanceDir(home, athName)
	} else {
		instDir = os.Getenv("ATHANOR")
	}

	if instDir == "" {
		fmt.Fprintln(os.Stderr, "error: provide athanor name or set $ATHANOR")
		fmt.Fprintln(os.Stderr, "usage: ath opera [<name>] [--mo <mo-name>]")
		return 2
	}

	// Determine which MOs to scan
	var mosToScan []string
	if moFilter != "" {
		mosToScan = []string{moFilter}
	} else {
		mos, _ := athanor.ListMagnaOpera(instDir)
		mosToScan = mos
	}

	hasOpera := false
	fmt.Printf("%-12s %-12s %-20s %s\n", "STATUS", "DATE", "MO", "OPUS")

	for _, mo := range mosToScan {
		operaDir := athanor.OperaPath(instDir, mo)
		entries, err := os.ReadDir(operaDir)
		if err != nil {
			continue
		}

		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			opusPath := filepath.Join(operaDir, e.Name())

			hasOpera = true
			status := readOpusStatus(opusPath)
			if status == "" {
				status = "unknown"
			}
			opusName := strings.TrimSuffix(e.Name(), ".md")

			// Extract date from filename (YYYY-MM-DD prefix)
			date := "-"
			if len(opusName) >= 10 && opusName[4] == '-' && opusName[7] == '-' {
				date = opusName[:10]
				opusName = opusName[11:] // Strip date prefix
			}

			fmt.Printf("%-12s %-12s %-20s %s\n", status, date, mo, opusName)
		}
	}

	if !hasOpera {
		fmt.Println("No opera found.")
	}

	return 0
}
