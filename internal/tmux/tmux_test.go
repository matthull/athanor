package tmux

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrapError(t *testing.T) {
	r := NewRunner()
	tests := []struct {
		name   string
		stderr string
		want   error
	}{
		{"no server", "no server running on /tmp/...", ErrNoServer},
		{"connection error", "error connecting to /tmp/...", ErrNoServer},
		{"server exited", "server exited unexpectedly", ErrNoServer},
		{"session not found", "session not found: foo", ErrSessionNotFound},
		{"cant find session", "can't find session: bar", ErrSessionNotFound},
		{"generic error", "some other error", nil}, // not a sentinel
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.wrapError(fmt.Errorf("exit 1"), tt.stderr, []string{"send-keys"})
			if tt.want != nil {
				if !errors.Is(err, tt.want) {
					t.Errorf("got %v, want %v", err, tt.want)
				}
			} else if errors.Is(err, ErrNoServer) || errors.Is(err, ErrSessionNotFound) {
				t.Errorf("should not be a sentinel error, got %v", err)
			}
		})
	}
}
