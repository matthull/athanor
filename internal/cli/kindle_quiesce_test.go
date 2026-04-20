package cli

import (
	"reflect"
	"testing"
)

func TestKindleArgPatterns(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantPos   []string
		wantFlags []string
	}{
		{
			name:      "name and mo positional only (backward compat)",
			args:      []string{"myathanor", "my-mo"},
			wantPos:   []string{"myathanor", "my-mo"},
			wantFlags: nil,
		},
		{
			name:      "name positional with --role flag",
			args:      []string{"myathanor", "--role", "perceiver", "--mo", "my-mo"},
			wantPos:   []string{"myathanor"},
			wantFlags: []string{"--role", "perceiver", "--mo", "my-mo"},
		},
		{
			name:      "name and mo positional with --role flag",
			args:      []string{"myathanor", "my-mo", "--role", "perceiver"},
			wantPos:   []string{"myathanor", "my-mo"},
			wantFlags: []string{"--role", "perceiver"},
		},
		{
			name:      "flags before positionals",
			args:      []string{"--role", "attendant", "--mo", "life-domains", "myathanor"},
			wantPos:   []string{"myathanor"},
			wantFlags: []string{"--role", "attendant", "--mo", "life-domains"},
		},
		{
			name:      "name only (legacy single-MO)",
			args:      []string{"myathanor"},
			wantPos:   []string{"myathanor"},
			wantFlags: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, flags := splitArgs(tt.args)
			if !reflect.DeepEqual(pos, tt.wantPos) {
				t.Errorf("positional = %v, want %v", pos, tt.wantPos)
			}
			if !reflect.DeepEqual(flags, tt.wantFlags) {
				t.Errorf("flags = %v, want %v", flags, tt.wantFlags)
			}
		})
	}
}

func TestQuiesceArgPatterns(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantPos   []string
		wantFlags []string
	}{
		{
			name:      "name and mo with --role",
			args:      []string{"myathanor", "my-mo", "--role", "perceiver"},
			wantPos:   []string{"myathanor", "my-mo"},
			wantFlags: []string{"--role", "perceiver"},
		},
		{
			name:      "name and mo with --force (backward compat)",
			args:      []string{"myathanor", "my-mo", "--force"},
			wantPos:   []string{"myathanor", "my-mo"},
			wantFlags: []string{"--force"},
		},
		{
			name:      "name with --role and --force",
			args:      []string{"myathanor", "--role", "perceiver", "--force"},
			wantPos:   []string{"myathanor"},
			wantFlags: []string{"--role", "perceiver", "--force"},
		},
		{
			name:      "name and mo only (backward compat)",
			args:      []string{"myathanor", "my-mo"},
			wantPos:   []string{"myathanor", "my-mo"},
			wantFlags: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, flags := splitArgs(tt.args)
			if !reflect.DeepEqual(pos, tt.wantPos) {
				t.Errorf("positional = %v, want %v", pos, tt.wantPos)
			}
			if !reflect.DeepEqual(flags, tt.wantFlags) {
				t.Errorf("flags = %v, want %v", flags, tt.wantFlags)
			}
		})
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"marut", "Marut"},
		{"perceiver", "Perceiver"},
		{"attendant", "Attendant"},
		{"", ""},
	}
	for _, tt := range tests {
		got := titleCase(tt.input)
		if got != tt.want {
			t.Errorf("titleCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
