package tmux

import "testing"

func TestSanitize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain text", "hello world", "hello world"},
		{"preserves newlines", "line1\nline2", "line1\nline2"},
		{"tab to space", "hello\tworld", "hello world"},
		{"strips ESC", "hello\x1bworld", "helloworld"},
		{"strips CR", "hello\rworld", "helloworld"},
		{"strips BS", "hello\x08world", "helloworld"},
		{"strips DEL", "hello\x7fworld", "helloworld"},
		{"strips null", "hello\x00world", "helloworld"},
		{"preserves quotes", `he said "hello"`, `he said "hello"`},
		{"preserves backticks", "run `cmd`", "run `cmd`"},
		{"preserves special chars", "!~#$%^&*()", "!~#$%^&*()"},
		{"preserves unicode", "hello 世界 🌍", "hello 世界 🌍"},
		{"mixed control chars", "a\x01b\x02c\x1bd\x7fe\tf\ng", "abcde f\ng"},
		{"empty string", "", ""},
		{"only control chars", "\x01\x02\x03\x1b\x7f", ""},
		{"preserves prompt char", "❯ hello", "❯ hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Sanitize(tt.input)
			if got != tt.want {
				t.Errorf("Sanitize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
