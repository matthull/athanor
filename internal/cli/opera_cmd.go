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

	fs := flag.NewFlagSet("opera", flag.ContinueOnError)
	fs.StringVar(&athName, "athanor", "", "athanor name (if $ATHANOR not set)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()

	// Resolve instance directory
	var instDir string
	if len(remaining) > 0 {
		athName = remaining[0]
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
		fmt.Fprintln(os.Stderr, "usage: ath opera [<name>]")
		return 2
	}

	operaDir := filepath.Join(instDir, athanor.OperaDir)
	entries, err := os.ReadDir(operaDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No opera found.")
			return 0
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	hasOpera := false
	fmt.Printf("%-12s %-12s %s\n", "STATUS", "DATE", "OPUS")
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		hasOpera = true
		status := readOpusStatus(filepath.Join(operaDir, e.Name()))
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

		fmt.Printf("%-12s %-12s %s\n", status, date, opusName)
	}

	if !hasOpera {
		fmt.Println("No opera found.")
	}

	return 0
}
