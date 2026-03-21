package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// testSocket is the isolated tmux socket name for this test run.
var testSocket string

func TestMain(m *testing.M) {
	testSocket = fmt.Sprintf("whisper-test-%d", os.Getpid())
	code := m.Run()
	// Kill the test tmux server
	_ = exec.Command("tmux", "-L", testSocket, "kill-server").Run()
	os.Exit(code)
}

// newTestRunner creates a Runner targeting the isolated test socket.
func newTestRunner(t *testing.T) *Runner {
	t.Helper()
	return NewRunnerWithSocket(testSocket)
}

// createTestSession creates an ephemeral tmux session for a test.
// Returns the session name. The session is killed via t.Cleanup.
func createTestSession(t *testing.T, r *Runner) string {
	t.Helper()
	name := fmt.Sprintf("wt-%s", t.Name())
	// Create a detached session running a simple shell.
	// Retry once on "server exited unexpectedly" — this transient error occurs
	// when a previous test's cleanup killed the last session and the server is
	// shutting down while we try to connect to the same socket.
	args := []string{"-u", "-L", testSocket, "new-session", "-d", "-s", name, "-x", "200", "-y", "50"}
	out, err := exec.Command("tmux", args...).CombinedOutput()
	if err != nil && strings.Contains(string(out), "server exited unexpectedly") {
		time.Sleep(50 * time.Millisecond)
		out, err = exec.Command("tmux", args...).CombinedOutput()
	}
	if err != nil {
		t.Fatalf("create test session %q: %v\n%s", name, err, out)
	}
	t.Cleanup(func() {
		_ = exec.Command("tmux", "-L", testSocket, "kill-session", "-t", name).Run()
	})
	return name
}
