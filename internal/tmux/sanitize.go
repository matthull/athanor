package tmux

import "strings"

// Sanitize strips control characters that corrupt tmux send-keys -l delivery.
// TAB is replaced with a space (avoids shell tab completion), newlines are
// preserved, all other chars < 0x20 and DEL (0x7f) are stripped.
// Printable characters including quotes, backticks, and Unicode are preserved.
// [G:tmux.go:1285-1303]
func Sanitize(msg string) string {
	var b strings.Builder
	b.Grow(len(msg))
	for _, r := range msg {
		switch {
		case r == '\t':
			b.WriteRune(' ')
		case r == '\n':
			b.WriteRune(r)
		case r < 0x20:
			continue
		case r == 0x7f:
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
