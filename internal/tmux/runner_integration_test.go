package tmux

import (
	"errors"
	"strings"
	"testing"
)

func TestRunnerBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r)

	// Verify we can query the session
	out, err := r.DisplayMessage(session, "#{session_name}")
	if err != nil {
		t.Fatalf("display-message: %v", err)
	}
	if !strings.Contains(out, session) {
		t.Errorf("expected session name %q in output %q", session, out)
	}
}

func TestRunnerSessionNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	r := newTestRunner(t)
	// Ensure the test tmux server is running so we get "session not found"
	// instead of "no server running".
	_ = createTestSession(t, r)
	// Use SendKeys with a session-qualified target (trailing colon) so tmux
	// returns "can't find session" rather than "can't find pane".
	// display-message -p silently succeeds on nonexistent targets in modern tmux.
	err := r.SendKeys("nonexistent-session-xyz:", "")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("expected ErrSessionNotFound, got %v", err)
	}
}
