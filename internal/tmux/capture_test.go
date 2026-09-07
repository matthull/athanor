package tmux

import "testing"

func TestMatchesPromptPrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		line string
		want bool
	}{
		{"prompt with text", "❯ hello", true},
		{"prompt with leading spaces", "  ❯ ", true},
		{"dollar prompt", "$ hello", false},
		{"empty string", "", false},
		{"bare prompt char", "❯", true},
		{"prompt not at start", "text ❯", false},
		{"NBSP after prompt", "❯\u00a0hello", true},
		{"leading space and NBSP", "  ❯\u00a0", true},
		{"just whitespace", "   ", false},
		{"similar unicode", "〉 hello", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := MatchesPromptPrefix(tt.line)
			if got != tt.want {
				t.Errorf("MatchesPromptPrefix(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestIsStatusBarBusy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			"busy with both indicators",
			[]string{"some output", "⏵⏵ Running tool — esc to interrupt"},
			true,
		},
		{
			"only progress indicator",
			[]string{"⏵⏵ Running tool"},
			false,
		},
		{
			"empty lines",
			[]string{},
			false,
		},
		{
			"nil lines",
			nil,
			false,
		},
		{
			"only esc to interrupt without progress",
			[]string{"esc to interrupt"},
			false,
		},
		{
			"progress and interrupt on same line with whitespace",
			[]string{"  ⏵⏵ Tool running — esc to interrupt  "},
			true,
		},
		{
			"progress on one line interrupt on another",
			[]string{"⏵⏵ Running", "esc to interrupt"},
			false,
		},
		{
			"idle prompt",
			[]string{"❯ ", ""},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := IsStatusBarBusy(tt.lines)
			if got != tt.want {
				t.Errorf("IsStatusBarBusy(%v) = %v, want %v", tt.lines, got, tt.want)
			}
		})
	}
}
