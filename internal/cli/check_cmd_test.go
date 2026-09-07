package cli

import "testing"

func TestExtractLastActive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			"status line with last_active",
			[]string{
				"some output",
				"[CC] 32% context | claude-sonnet-4-20250514 | last_active:04-27 11:42",
			},
			"04-27 11:42",
		},
		{
			"last_active at different position",
			[]string{
				"[CC] 80% context | claude-opus-4-20250514 | last_active:04-26 09:15",
				"some other output",
			},
			"04-26 09:15",
		},
		{
			"no last_active present",
			[]string{
				"some output",
				"building...",
				"❯ ",
			},
			"",
		},
		{
			"empty lines",
			[]string{},
			"",
		},
		{
			"nil lines",
			nil,
			"",
		},
		{
			"returns last match when multiple present",
			[]string{
				"[CC] 32% context | model | last_active:04-27 10:00",
				"some output in between",
				"[CC] 45% context | model | last_active:04-27 11:30",
			},
			"04-27 11:30",
		},
		{
			"ignores non-timestamp last_active patterns",
			[]string{
				"last_active:foo bar",
				"[CC] 32% context | model | last_active:04-27 11:42",
			},
			"04-27 11:42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := extractLastActive(tt.lines)
			if got != tt.want {
				t.Errorf("extractLastActive() = %q, want %q", got, tt.want)
			}
		})
	}
}
