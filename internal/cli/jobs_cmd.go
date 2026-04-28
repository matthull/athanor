package cli

import (
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
)

// runJobs handles the "ath jobs" command.
// No args: list all jobs with summaries.
// One arg: show detail for a specific job.
func runJobs(args []string) int {
	if len(args) == 0 {
		return listJobs()
	}
	return showJob(args[0])
}

func listJobs() int {
	infos, err := athanor.ListJobInfos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if len(infos) == 0 {
		fmt.Fprintln(os.Stderr, "no jobs found")
		return 1
	}

	fmt.Println("Available jobs:")
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

func showJob(name string) int {
	if err := athanor.ValidateJob(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	info, content, err := athanor.ReadJobDetail(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Structured header
	fmt.Printf("Job: %s\n", info.Name)
	if info.Summary != "" {
		fmt.Printf("Summary: %s\n", info.Summary)
	}
	if info.Model != "" {
		fmt.Printf("Model: %s\n", info.Model)
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
