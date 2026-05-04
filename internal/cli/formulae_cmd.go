package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
)

// runFormulae handles the "ath formulae" command.
// No positional args: list all formulae with summaries and when-to-invoke entries.
// One positional arg: show detail for a specific formula.
// Operates on $ATHANOR or --athanor flag.
func runFormulae(args []string) int {
	var athName string
	fs := flag.NewFlagSet("formulae", flag.ContinueOnError)
	fs.StringVar(&athName, "athanor", "", "athanor name (if $ATHANOR not set)")
	fs.SetOutput(os.Stderr)

	positional, flagArgs := splitArgs(args)
	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	instDir, err := resolveInstDir(athName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitCode(err)
	}

	if len(positional) == 0 {
		return listFormulae(instDir)
	}
	return showFormula(instDir, positional[0])
}

func listFormulae(instDir string) int {
	infos, err := athanor.ListFormulaInfos(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if len(infos) == 0 {
		fmt.Fprintln(os.Stderr, "no formulae found")
		return 1
	}

	fmt.Println("Available formulae:")
	for _, info := range infos {
		summary := info.Summary
		if summary == "" {
			summary = "(no summary)"
		}
		fmt.Println()
		fmt.Printf("  %s — %s\n", info.Name, summary)
		for _, w := range info.When {
			fmt.Printf("    - %s\n", w)
		}
	}
	fmt.Println()
	return 0
}

func showFormula(instDir, name string) int {
	if err := athanor.ValidateFormula(instDir, name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	info, content, err := athanor.ReadFormulaDetail(instDir, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Structured header
	fmt.Printf("Formula: %s\n", info.Name)
	if info.Summary != "" {
		fmt.Printf("Summary: %s\n", info.Summary)
	}
	if len(info.When) > 0 {
		fmt.Println("When to invoke:")
		for _, w := range info.When {
			fmt.Printf("  - %s\n", w)
		}
	}
	fmt.Println()
	fmt.Print(athanor.StripFrontmatter(content))
	return 0
}
