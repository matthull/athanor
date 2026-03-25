package cli

import (
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantPos   []string
		wantFlags []string
	}{
		{
			name:      "flags before positional",
			args:      []string{"--project", "/tmp", "myname"},
			wantPos:   []string{"myname"},
			wantFlags: []string{"--project", "/tmp"},
		},
		{
			name:      "flags after positional",
			args:      []string{"myname", "--project", "/tmp"},
			wantPos:   []string{"myname"},
			wantFlags: []string{"--project", "/tmp"},
		},
		{
			name:      "boolean flag after positional",
			args:      []string{"myname", "--force"},
			wantPos:   []string{"myname"},
			wantFlags: []string{"--force"},
		},
		{
			name:      "multiple flags mixed",
			args:      []string{"opus.md", "--dir", "/tmp", "--model", "opus", "--athanor", "bugsnag"},
			wantPos:   []string{"opus.md"},
			wantFlags: []string{"--dir", "/tmp", "--model", "opus", "--athanor", "bugsnag"},
		},
		{
			name:      "flag with equals",
			args:      []string{"myname", "--project=/tmp"},
			wantPos:   []string{"myname"},
			wantFlags: []string{"--project=/tmp"},
		},
		{
			name:      "only positional",
			args:      []string{"myname"},
			wantPos:   []string{"myname"},
			wantFlags: nil,
		},
		{
			name:      "only flags",
			args:      []string{"--force"},
			wantPos:   nil,
			wantFlags: []string{"--force"},
		},
		{
			name:      "empty",
			args:      []string{},
			wantPos:   nil,
			wantFlags: nil,
		},
		{
			name:      "short flag",
			args:      []string{"target", "-f", "/tmp/msg.txt"},
			wantPos:   []string{"target"},
			wantFlags: []string{"-f", "/tmp/msg.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
