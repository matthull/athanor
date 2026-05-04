package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
)

// runTzurot handles the "ath tzurot" command — lists both jobs (global)
// and formulae (athanor-specific). Formulae require an athanor context;
// when $ATHANOR is not set and no --athanor flag is given, only jobs are listed.
func runTzurot(args []string) int {
	var athName string
	fs := flag.NewFlagSet("tzurot", flag.ContinueOnError)
	fs.StringVar(&athName, "athanor", "", "athanor name (overrides $ATHANOR for formulae listing)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	// Jobs section — always shown
	jobInfos, err := athanor.ListJobInfos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println("Jobs (global):")
	if len(jobInfos) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, info := range jobInfos {
			summary := info.Summary
			if summary == "" {
				summary = "(no summary)"
			}
			fmt.Printf("  %s — %s\n", info.Name, summary)
		}
	}

	// Formulae section — only when an athanor is resolvable.
	// Graceful absence: if $ATHANOR is not set and --athanor is not given,
	// or if the formulae/ directory is missing/empty, the section is omitted.
	instDir, athLabel := resolveTzurotAthanor(athName)
	if instDir == "" {
		return 0
	}

	formulaInfos, err := athanor.ListFormulaInfos(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: listing formulae: %v\n", err)
		return 0
	}
	if len(formulaInfos) == 0 {
		return 0
	}

	fmt.Println()
	fmt.Printf("Formulae (%s):\n", athLabel)
	for _, info := range formulaInfos {
		summary := info.Summary
		if summary == "" {
			summary = "(no summary)"
		}
		fmt.Printf("  %s — %s\n", info.Name, summary)
	}
	return 0
}

// resolveTzurotAthanor returns the instance dir and a display label for the athanor.
// Returns ("", "") when no athanor context is available. Unlike resolveInstDir, this
// never errors — graceful absence is the documented behavior.
func resolveTzurotAthanor(athName string) (instDir, label string) {
	if athName != "" {
		home, err := athanor.Home()
		if err != nil {
			return "", ""
		}
		return athanor.InstanceDir(home, athName), athName
	}
	if env := os.Getenv("ATHANOR"); env != "" {
		return env, athanor.AthanorName(env)
	}
	return "", ""
}
