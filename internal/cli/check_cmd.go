package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/matthull/athanor/internal/tmux"
)

var reLastActive = regexp.MustCompile(`last_active:(\d{2}-\d{2} \d{2}:\d{2})`)

func runCheck(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "error: crucible name required")
		fmt.Fprintln(os.Stderr, "usage: ath check <crucible-name>")
		return 2
	}
	crucible := args[0]

	r := tmux.NewRunner()

	// Resolve bare window names to fully-qualified session:window targets.
	target := crucible
	if !strings.Contains(target, ":") {
		resolved, err := r.FindWindow(target)
		if err != nil {
			fmt.Fprintln(os.Stderr, "crucible not found")
			return 2
		}
		target = resolved
	}

	// Capture last 25 lines of the crucible pane.
	lines, err := r.CapturePaneLines(target, 25)
	if err != nil {
		fmt.Fprintln(os.Stderr, "crucible not found")
		return 2
	}

	// Extract last_active timestamp from the status line if present.
	lastActive := extractLastActive(lines)
	if lastActive != "" {
		fmt.Printf("last_active: %s\n\n", lastActive)
	}

	// Display captured pane content for human interpretation.
	for _, line := range lines {
		fmt.Println(line)
	}

	return 0
}

// extractLastActive scans lines for a last_active:<timestamp> pattern
// from the Claude Code status bar and returns the timestamp value.
// Returns the last match since the status bar is at the bottom of the pane.
func extractLastActive(lines []string) string {
	var result string
	for _, line := range lines {
		if m := reLastActive.FindStringSubmatch(line); m != nil {
			result = m[1]
		}
	}
	return result
}
