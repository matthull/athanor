package cli

import (
	"testing"

	"github.com/matthull/athanor/internal/tmux"
)

func TestDetectPermissionFinding(t *testing.T) {
	t.Parallel()
	pane := tmux.Pane{Target: "test:0.0", Session: "test", Window: "0", Pane: "0", Command: "claude"}

	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			"both patterns in bottom 10",
			[]string{
				"some output",
				"more output",
				"Do you want to proceed with this action?",
				"  ❯ 1. Allow",
				"    2. Deny",
			},
			true,
		},
		{
			"patterns outside bottom 10 not detected",
			// 12 lines total, patterns in first 2 (outside bottom 10)
			append(
				[]string{
					"Do you want to proceed?",
					"  ❯ 1. Allow",
				},
				"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
			),
			false,
		},
		{
			"cursor only",
			[]string{"  ❯ 1. Some item"},
			false,
		},
		{
			"question only",
			[]string{"Do you want to proceed?"},
			false,
		},
		{
			"empty",
			[]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := detectPermissionFinding(pane, tt.lines)
			if (got != nil) != tt.want {
				t.Errorf("detectPermissionFinding() returned finding=%v, want finding=%v", got != nil, tt.want)
			}
			if got != nil && got.Type != FindingPermission {
				t.Errorf("finding type = %q, want %q", got.Type, FindingPermission)
			}
		})
	}
}

func TestDetectExhaustedFinding(t *testing.T) {
	t.Parallel()
	pane := tmux.Pane{Target: "test:0.0", Session: "test", Window: "0", Pane: "0", Command: "claude"}

	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			"context limit in bottom 5",
			[]string{"output", "Context limit reached", ""},
			true,
		},
		{
			"context limit outside bottom 5",
			append(
				[]string{"Context limit reached"},
				"a", "b", "c", "d", "e",
			),
			false,
		},
		{
			"no context limit",
			[]string{"all good", "still running"},
			false,
		},
		{
			"empty",
			[]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := detectExhaustedFinding(pane, tt.lines)
			if (got != nil) != tt.want {
				t.Errorf("detectExhaustedFinding() returned finding=%v, want finding=%v", got != nil, tt.want)
			}
			if got != nil && got.Type != FindingExhausted {
				t.Errorf("finding type = %q, want %q", got.Type, FindingExhausted)
			}
		})
	}
}

func TestIsClaudeProcess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		cmd  string
		want bool
	}{
		{"claude", true},
		{"1.0.5", true},
		{"12.34.56", true},
		{"zsh", false},
		{"bash", false},
		{"vim", false},
		{"node", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.cmd, func(t *testing.T) {
			t.Parallel()
			if got := isClaudeProcess(tt.cmd); got != tt.want {
				t.Errorf("isClaudeProcess(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}

func TestBottomLines(t *testing.T) {
	t.Parallel()

	t.Run("returns last n", func(t *testing.T) {
		t.Parallel()
		lines := []string{"a", "b", "c", "d", "e"}
		got := bottomLines(lines, 3)
		if len(got) != 3 || got[0] != "c" || got[1] != "d" || got[2] != "e" {
			t.Errorf("bottomLines returned %v, want [c d e]", got)
		}
	})

	t.Run("fewer than n", func(t *testing.T) {
		t.Parallel()
		lines := []string{"a", "b"}
		got := bottomLines(lines, 5)
		if len(got) != 2 {
			t.Errorf("bottomLines returned len=%d, want 2", len(got))
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := bottomLines([]string{}, 5)
		if len(got) != 0 {
			t.Errorf("bottomLines returned len=%d, want 0", len(got))
		}
	})
}
