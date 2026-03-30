package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/matthull/athanor/internal/tmux"
)

// FindingType classifies what patrol found in a pane.
type FindingType string

const (
	FindingPermission FindingType = "permission"
	FindingExhausted  FindingType = "exhausted"
)

// Finding represents something patrol detected in a tmux pane.
type Finding struct {
	Type   FindingType
	Pane   tmux.Pane
	Detail string
}

func runPatrol(args []string) int {
	fs := flag.NewFlagSet("patrol", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	jsonFlag := fs.Bool("json", false, "output as JSON lines")
	excludeFlag := fs.String("exclude", "", "pane target to exclude (e.g. self)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	r := tmux.NewRunner()
	findings, err := scan(r, *excludeFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if len(findings) == 0 {
		if !*jsonFlag {
			fmt.Println("No findings.")
		}
		return 0
	}

	for _, f := range findings {
		if *jsonFlag {
			fmt.Printf(`{"type":%q,"pane":%q,"session":%q,"detail":%q}%s`,
				f.Type, f.Pane.Target, f.Pane.Session, f.Detail, "\n")
		} else {
			fmt.Printf("%-12s %-30s %s\n", f.Type, f.Pane.Target, f.Detail)
		}
	}

	return 0
}

// regexps for permission prompt detection — same patterns as check.go
var (
	patrolPermMenuCursor = regexp.MustCompile(`❯\s+\d+\.`)
	patrolPermQuestion   = regexp.MustCompile(`Do you want to (proceed|make this edit)|Would you like to proceed`)
)

// scan inspects all tmux panes and returns findings.
func scan(r *tmux.Runner, exclude string) ([]Finding, error) {
	panes, err := r.ListPanes()
	if err != nil {
		return nil, fmt.Errorf("list panes: %w", err)
	}

	var findings []Finding
	for _, p := range panes {
		if exclude != "" && p.Target == exclude {
			continue
		}

		lines, err := r.CapturePaneLines(p.Target, 20)
		if err != nil {
			continue // pane may have died between list and capture
		}

		if f := detectPermissionFinding(p, lines); f != nil {
			findings = append(findings, *f)
			continue // one finding per pane, permission takes priority
		}

		// Only check for stalls in panes running Claude
		if isClaudeProcess(p.Command) {
			if f := detectExhaustedFinding(p, lines); f != nil {
				findings = append(findings, *f)
			}
		}
	}

	return findings, nil
}

// detectPermissionFinding checks the bottom 10 lines for co-occurrence of
// menu cursor and permission question text.
func detectPermissionFinding(p tmux.Pane, lines []string) *Finding {
	bottom10 := bottomLines(lines, 10)
	hasCursor := false
	hasQuestion := false
	for _, line := range bottom10 {
		if patrolPermMenuCursor.MatchString(line) {
			hasCursor = true
		}
		if patrolPermQuestion.MatchString(line) {
			hasQuestion = true
		}
		if hasCursor && hasQuestion {
			return &Finding{
				Type:   FindingPermission,
				Pane:   p,
				Detail: "permission prompt blocking",
			}
		}
	}
	return nil
}

// detectExhaustedFinding checks the bottom 5 lines for the context limit message.
func detectExhaustedFinding(p tmux.Pane, lines []string) *Finding {
	bottom5 := bottomLines(lines, 5)
	for _, line := range bottom5 {
		if strings.Contains(line, "Context limit reached") {
			return &Finding{
				Type:   FindingExhausted,
				Pane:   p,
				Detail: "context limit reached",
			}
		}
	}
	return nil
}

// isClaudeProcess returns true if the command looks like a Claude Code process.
// Matches "claude" or version-number patterns like "1.0.5".
var reVersionProcess = regexp.MustCompile(`^\d+\.\d+\.\d+`)

func isClaudeProcess(cmd string) bool {
	return cmd == "claude" || reVersionProcess.MatchString(cmd)
}

// bottomLines returns the last n lines, or all if fewer.
func bottomLines(lines []string, n int) []string {
	if len(lines) <= n {
		return lines
	}
	return lines[len(lines)-n:]
}
