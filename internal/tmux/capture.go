package tmux

import "strings"

// DefaultPromptPrefix is Claude Code's idle prompt character.
// ❯ is U+276F HEAVY RIGHT-POINTING ANGLE QUOTATION MARK.
// [G:tmux.go:2462]
const DefaultPromptPrefix = "❯"

// MatchesPromptPrefix reports whether a captured pane line starts with the
// Claude Code prompt prefix. Normalizes NBSP (U+00A0) → regular space because
// some terminal renderers use NBSP. [G:tmux.go:2410-2421]
func MatchesPromptPrefix(line string) bool {
	trimmed := strings.TrimSpace(line)
	// Normalize NBSP (U+00A0) → regular space
	trimmed = strings.ReplaceAll(trimmed, "\u00a0", " ")
	return strings.HasPrefix(trimmed, DefaultPromptPrefix)
}

// IsStatusBarBusy returns true if the captured lines indicate Claude Code is
// actively running a tool call. Detects the "⏵⏵" progress indicator combined
// with "esc to interrupt" text. [G:tmux.go:2500-2508]
func IsStatusBarBusy(lines []string) bool {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "⏵⏵") {
			if strings.Contains(trimmed, "esc to interrupt") {
				return true
			}
			break
		}
	}
	return false
}
