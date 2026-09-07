package tmux

import "testing"

func TestCrucibleStateString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		state CrucibleState
		want  string
	}{
		{StateActive, "active"},
		{StateIdle, "idle"},
		{StatePermission, "permission"},
		{StateExhausted, "exhausted"},
		{StateDead, "dead"},
		{CrucibleState(99), "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.state.String(); got != tt.want {
				t.Errorf("CrucibleState(%d).String() = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestDetectPermission(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			"both patterns present",
			[]string{
				"Do you want to proceed with this action?",
				"  ❯ 1. Allow",
				"    2. Deny",
			},
			true,
		},
		{
			"make this edit variant",
			[]string{
				"Do you want to make this edit?",
				"  ❯ 1. Yes",
			},
			true,
		},
		{
			"would you like to proceed variant",
			[]string{
				"Would you like to proceed?",
				"  ❯ 3. Cancel",
			},
			true,
		},
		{
			"only menu cursor no question",
			[]string{
				"  ❯ 1. Some item",
			},
			false,
		},
		{
			"only question no cursor",
			[]string{
				"Do you want to proceed?",
			},
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
			"cursor with extra whitespace",
			[]string{
				"Do you want to proceed?",
				"❯   42. Some long option",
			},
			true,
		},
		{
			"unrelated output",
			[]string{
				"Building project...",
				"Tests passed: 42",
				"❯ ",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := detectPermission(tt.lines)
			if got != tt.want {
				t.Errorf("detectPermission(%v) = %v, want %v", tt.lines, got, tt.want)
			}
		})
	}
}

func TestDetectExhausted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			"context limit present",
			[]string{"some output", "Context limit reached", ""},
			true,
		},
		{
			"context limit with surrounding text",
			[]string{"  Context limit reached. Start a new conversation."},
			true,
		},
		{
			"no context limit",
			[]string{"some output", "all good"},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := detectExhausted(tt.lines)
			if got != tt.want {
				t.Errorf("detectExhausted(%v) = %v, want %v", tt.lines, got, tt.want)
			}
		})
	}
}

func TestBottomN(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lines []string
		n     int
		want  int // expected length
	}{
		{"more lines than n", []string{"a", "b", "c", "d", "e"}, 3, 3},
		{"fewer lines than n", []string{"a", "b"}, 5, 2},
		{"exact match", []string{"a", "b", "c"}, 3, 3},
		{"empty slice", []string{}, 5, 0},
		{"nil slice", nil, 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := bottomN(tt.lines, tt.n)
			if len(got) != tt.want {
				t.Errorf("bottomN(len=%d, %d) returned len=%d, want %d", len(tt.lines), tt.n, len(got), tt.want)
			}
		})
	}

	// Verify bottomN returns the actual last elements
	t.Run("returns last elements", func(t *testing.T) {
		t.Parallel()
		lines := []string{"a", "b", "c", "d", "e"}
		got := bottomN(lines, 2)
		if got[0] != "d" || got[1] != "e" {
			t.Errorf("bottomN returned %v, want [d e]", got)
		}
	})
}

func TestDetectionPriority(t *testing.T) {
	t.Parallel()
	// Permission should be detected even when idle prompt and busy indicators are present.
	// This tests the priority ordering of the detection logic.
	tests := []struct {
		name  string
		lines []string
		want  CrucibleState
	}{
		{
			"permission takes priority over idle",
			[]string{
				"❯ ",
				"Do you want to proceed?",
				"  ❯ 1. Allow",
				"    2. Deny",
			},
			StatePermission,
		},
		{
			"permission takes priority over busy",
			[]string{
				"⏵⏵ Running tool — esc to interrupt",
				"Do you want to proceed?",
				"  ❯ 1. Allow",
			},
			StatePermission,
		},
		{
			"exhausted takes priority over idle",
			[]string{
				"❯ ",
				"Context limit reached",
			},
			StateExhausted,
		},
		{
			"exhausted takes priority over busy",
			[]string{
				"⏵⏵ Running tool — esc to interrupt",
				"Context limit reached",
			},
			StateExhausted,
		},
		{
			"busy takes priority over idle",
			[]string{
				"❯ ",
				"⏵⏵ Running tool — esc to interrupt",
			},
			StateActive,
		},
		{
			"idle when only prompt visible",
			[]string{
				"some output",
				"❯ ",
			},
			StateIdle,
		},
		{
			"active fallback when no patterns match",
			[]string{
				"some random output",
				"building...",
			},
			StateActive,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test detection logic without tmux by calling the individual
			// detection functions in the same priority order as CheckCrucible.
			bottom10 := bottomN(tt.lines, 10)
			bottom5 := bottomN(tt.lines, 5)

			var got CrucibleState
			switch {
			case detectPermission(bottom10):
				got = StatePermission
			case detectExhausted(bottom5):
				got = StateExhausted
			case IsStatusBarBusy(tt.lines):
				got = StateActive
			case hasPromptPrefix(tt.lines):
				got = StateIdle
			default:
				got = StateActive
			}

			if got != tt.want {
				t.Errorf("detection priority for %q: got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// hasPromptPrefix checks if any line matches the prompt prefix (mirrors CheckCrucible logic).
func hasPromptPrefix(lines []string) bool {
	for _, line := range lines {
		if MatchesPromptPrefix(line) {
			return true
		}
	}
	return false
}
