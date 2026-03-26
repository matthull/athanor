package tmux

import (
	"regexp"
	"strings"
)

// CrucibleState represents the health state of a crucible (tmux window).
type CrucibleState int

const (
	StateActive     CrucibleState = iota // tool call in progress
	StateIdle                            // waiting for input
	StatePermission                      // permission prompt blocking
	StateExhausted                       // context limit reached
	StateDead                            // crucible not found
)

// String returns a human-readable label for the state.
func (s CrucibleState) String() string {
	switch s {
	case StateActive:
		return "active"
	case StateIdle:
		return "idle"
	case StatePermission:
		return "permission"
	case StateExhausted:
		return "exhausted"
	case StateDead:
		return "dead"
	default:
		return "unknown"
	}
}

// regexps for permission prompt detection
var (
	rePermMenuCursor = regexp.MustCompile(`❯\s+\d+\.`)
	rePermQuestion   = regexp.MustCompile(`Do you want to (proceed|make this edit)|Would you like to proceed`)
)

// CheckCrucible inspects a crucible's tmux window and returns its health state
// with a short description. Detection priority:
//  1. Window doesn't exist → StateDead
//  2. Permission prompt in bottom 10 lines → StatePermission
//  3. Context limit in last 5 lines → StateExhausted
//  4. Status bar busy → StateActive
//  5. Prompt prefix visible → StateIdle
//  6. Fallback → StateActive
func (r *Runner) CheckCrucible(target string) (CrucibleState, string, error) {
	// 1. Check window exists
	windows, err := r.ListWindows()
	if err != nil {
		return StateDead, "crucible not found", err
	}
	found := false
	for _, w := range windows {
		if w == target {
			found = true
			break
		}
	}
	if !found {
		return StateDead, "crucible not found", nil
	}

	// 2. Capture last 20 lines
	lines, err := r.CapturePaneLines(target, 20)
	if err != nil {
		return StateDead, "capture failed", err
	}

	// 3. Check bottom 10 for permission prompt (both patterns must match)
	bottom10 := bottomN(lines, 10)
	if detectPermission(bottom10) {
		return StatePermission, "permission prompt blocking", nil
	}

	// 4. Check last 5 for context exhausted
	bottom5 := bottomN(lines, 5)
	if detectExhausted(bottom5) {
		return StateExhausted, "context limit reached", nil
	}

	// 5. Check status bar busy
	if IsStatusBarBusy(lines) {
		return StateActive, "tool call in progress", nil
	}

	// 6. Check for idle prompt
	for _, line := range lines {
		if MatchesPromptPrefix(line) {
			return StateIdle, "waiting for input", nil
		}
	}

	// 7. Fallback — something is happening
	return StateActive, "tool call in progress", nil
}

// detectPermission returns true if both the menu cursor and question patterns
// appear in the given lines (co-occurrence guards against false positives).
func detectPermission(lines []string) bool {
	hasCursor := false
	hasQuestion := false
	for _, line := range lines {
		if rePermMenuCursor.MatchString(line) {
			hasCursor = true
		}
		if rePermQuestion.MatchString(line) {
			hasQuestion = true
		}
		if hasCursor && hasQuestion {
			return true
		}
	}
	return false
}

// detectExhausted returns true if any line contains the context limit message.
func detectExhausted(lines []string) bool {
	for _, line := range lines {
		if strings.Contains(line, "Context limit reached") {
			return true
		}
	}
	return false
}

// bottomN returns the last n elements of a slice, or the whole slice if shorter.
func bottomN(lines []string, n int) []string {
	if len(lines) <= n {
		return lines
	}
	return lines[len(lines)-n:]
}
