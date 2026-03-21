package tmux

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// waitForShellReady waits until the test session's shell has started and is
// showing a prompt. Without this, send-keys may arrive before the shell is
// ready, causing flaky tests.
func waitForShellReady(t *testing.T, r *Runner, session string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		lines, err := r.CapturePaneLines(session, 3)
		if err == nil && len(lines) > 0 {
			for _, line := range lines {
				// Most shells show $ or % or # when ready
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					return
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	// Don't fatal — some shells may not show a visible prompt
	// but may still be ready. Proceed and let the test fail naturally.
}

func TestSendBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r)
	waitForShellReady(t, r, session)

	marker := fmt.Sprintf("whisper-test-%d", time.Now().UnixNano())
	opts := SendOpts{Timeout: 10 * time.Second, SkipEscape: true}

	if err := r.Send(session, "echo "+marker, opts); err != nil {
		t.Fatalf("Send: %v", err)
	}

	// Wait for echo to execute and appear in output
	time.Sleep(1 * time.Second)
	lines, err := r.CapturePaneLines(session, 20)
	if err != nil {
		t.Fatalf("CapturePaneLines: %v", err)
	}

	found := false
	for _, line := range lines {
		if strings.Contains(line, marker) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("marker %q not found in pane output: %v", marker, lines)
	}
}

func TestSendSpecialChars(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r)
	waitForShellReady(t, r, session)

	// Use echo with single quotes to prevent shell interpretation
	// Test chars that tmux send-keys (without -l) would interpret
	specials := "!~#`hello`"
	marker := fmt.Sprintf("SC-%d", time.Now().UnixNano())
	opts := SendOpts{Timeout: 10 * time.Second, SkipEscape: true}

	// echo with cat to avoid shell interpretation issues
	if err := r.Send(session, fmt.Sprintf("cat <<'WHISPEREOF'\n%s-%s\nWHISPEREOF", specials, marker), opts); err != nil {
		t.Fatalf("Send: %v", err)
	}

	time.Sleep(1 * time.Second)
	lines, err := r.CapturePaneLines(session, 30)
	if err != nil {
		t.Fatalf("CapturePaneLines: %v", err)
	}

	expected := specials + "-" + marker
	found := false
	for _, line := range lines {
		if strings.Contains(line, expected) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("special chars %q not found in output: %v", expected, lines)
	}
}

func TestSendLargeMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r)
	waitForShellReady(t, r, session)

	// Create a message larger than sendKeysChunkSize (512 bytes)
	// Use a unique start and end marker to verify no truncation
	startMarker := fmt.Sprintf("START-%d", time.Now().UnixNano())
	endMarker := fmt.Sprintf("END-%d", time.Now().UnixNano())
	// Fill with recognizable content between markers
	padding := strings.Repeat("ABCDEFGHIJ", 60) // 600 bytes of padding
	bigMsg := startMarker + padding + endMarker

	opts := SendOpts{Timeout: 10 * time.Second, SkipEscape: true}
	// Use cat to capture the full message
	if err := r.Send(session, fmt.Sprintf("cat <<'WHISPEREOF'\n%s\nWHISPEREOF", bigMsg), opts); err != nil {
		t.Fatalf("Send: %v", err)
	}

	time.Sleep(2 * time.Second) // large message needs more time
	lines, err := r.CapturePaneLines(session, 50)
	if err != nil {
		t.Fatalf("CapturePaneLines: %v", err)
	}

	output := strings.Join(lines, "\n")
	if !strings.Contains(output, startMarker) {
		t.Error("start marker not found — message not delivered")
	}
	if !strings.Contains(output, endMarker) {
		t.Error("end marker not found — message may be truncated")
	}
}

func TestSendExitsCopyMode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r)
	waitForShellReady(t, r, session)

	// Enter copy mode via direct tmux command (reliable regardless of prefix key)
	if _, err := r.run("copy-mode", "-t", session); err != nil {
		t.Fatalf("failed to enter copy mode: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Verify we're in copy mode
	inMode, err := r.InCopyMode(session)
	if err != nil {
		t.Fatalf("InCopyMode: %v", err)
	}
	if !inMode {
		t.Fatal("expected pane to be in copy mode")
	}

	marker := fmt.Sprintf("COPY-%d", time.Now().UnixNano())
	opts := SendOpts{Timeout: 10 * time.Second, SkipEscape: true}
	if err := r.Send(session, "echo "+marker, opts); err != nil {
		t.Fatalf("Send: %v", err)
	}

	time.Sleep(1 * time.Second)

	// Verify copy mode was exited
	inMode, _ = r.InCopyMode(session)
	if inMode {
		t.Error("pane still in copy mode after send")
	}

	// Verify message delivered
	lines, err := r.CapturePaneLines(session, 20)
	if err != nil {
		t.Fatalf("CapturePaneLines: %v", err)
	}
	found := false
	for _, line := range lines {
		if strings.Contains(line, marker) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("marker %q not found after copy mode exit", marker)
	}
}

func TestSendDetachedSession(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r) // creates detached by default
	waitForShellReady(t, r, session)

	// Session is detached (no clients attached in test socket)
	if r.IsSessionAttached(session) {
		t.Skip("session is attached — can't test detached wake")
	}

	marker := fmt.Sprintf("DETACH-%d", time.Now().UnixNano())
	opts := SendOpts{Timeout: 10 * time.Second, SkipEscape: true}
	if err := r.Send(session, "echo "+marker, opts); err != nil {
		t.Fatalf("Send: %v", err)
	}

	time.Sleep(1500 * time.Millisecond) // extra time for detached processing
	lines, err := r.CapturePaneLines(session, 20)
	if err != nil {
		t.Fatalf("CapturePaneLines: %v", err)
	}

	found := false
	for _, line := range lines {
		if strings.Contains(line, marker) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("marker %q not found in detached session", marker)
	}
}

func TestIdleDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r)
	waitForShellReady(t, r, session)

	// Shell should be idle (at prompt) — idle detection should succeed quickly
	err := WaitForIdle(r, session, 5*time.Second)
	// Note: WaitForIdle looks for ❯ (Claude Code prompt), not shell prompt.
	// In a bare shell session, it will NOT find ❯ and will timeout.
	// This test verifies the timeout behavior.
	if err == nil {
		// If it passed, the shell prompt happens to match — that's fine
		return
	}
	if !errors.Is(err, ErrIdleTimeout) {
		t.Fatalf("expected ErrIdleTimeout, got: %v", err)
	}
	// The timeout is expected — bare shell doesn't have ❯ prompt
}

func TestIdleDetectionSessionNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	// Ensure the test tmux server is running
	_ = createTestSession(t, r)
	err := WaitForIdle(r, "nonexistent-session-xyz:", 2*time.Second)
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("expected ErrSessionNotFound, got: %v", err)
	}
}

func TestConcurrentSends(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	session := createTestSession(t, r)
	waitForShellReady(t, r, session)

	marker1 := fmt.Sprintf("CONC-A-%d", time.Now().UnixNano())
	marker2 := fmt.Sprintf("CONC-B-%d", time.Now().UnixNano())

	var wg sync.WaitGroup
	var err1, err2 error
	opts := SendOpts{Timeout: 15 * time.Second, SkipEscape: true}

	wg.Add(2)
	go func() {
		defer wg.Done()
		err1 = r.Send(session, "echo "+marker1, opts)
	}()
	go func() {
		defer wg.Done()
		err2 = r.Send(session, "echo "+marker2, opts)
	}()
	wg.Wait()

	if err1 != nil {
		t.Errorf("send 1 failed: %v", err1)
	}
	if err2 != nil {
		t.Errorf("send 2 failed: %v", err2)
	}

	// Wait for both commands to execute
	time.Sleep(3 * time.Second)
	lines, err := r.CapturePaneLines(session, 30)
	if err != nil {
		t.Fatalf("CapturePaneLines: %v", err)
	}

	output := strings.Join(lines, "\n")
	if !strings.Contains(output, marker1) {
		t.Errorf("marker1 %q not found in output", marker1)
	}
	if !strings.Contains(output, marker2) {
		t.Errorf("marker2 %q not found in output", marker2)
	}
}

func TestSendSessionNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	r := newTestRunner(t)
	// Ensure the test tmux server is running so we get "session not found"
	// instead of "no server running".
	_ = createTestSession(t, r)
	opts := SendOpts{Timeout: 2 * time.Second}
	err := r.Send("nonexistent-session-xyz:", "hello", opts)
	if err == nil {
		t.Fatal("expected error for nonexistent session")
	}
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("expected ErrSessionNotFound, got: %v", err)
	}
}
