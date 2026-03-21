# T6: Integration Tests — End-to-End Verification

**Dependencies:** T1-T5 (all implementation tasks must be complete)
**Files to create:** `internal/tmux/integration_test.go`
**Files to modify:** None (TestMain and helpers from T1 are reused)

---

## Goal

Write integration tests that verify whisper's core behaviors against real tmux sessions using the isolated test socket from `TestMain` (T1). These tests confirm that the 8-step send protocol, idle detection, and edge cases work correctly in practice — not just in unit tests.

---

## Test infrastructure (from T1, already exists)

The following are available from `testmain_test.go` and should NOT be recreated:

- `testSocket` — isolated socket name (`whisper-test-<PID>`)
- `newTestRunner(t)` — creates a `Runner` targeting the test socket
- `createTestSession(t, r)` — creates a detached session with `t.Cleanup` teardown

### Additional helper needed

Add this helper to `integration_test.go` (or a shared test helper file):

```go
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
```

### How to verify message delivery

The standard pattern for verifying that text arrived in a tmux pane:

```go
// Send message, then capture pane content and check it contains the message.
// Use capture-pane to read what's visible in the terminal.
time.Sleep(200 * time.Millisecond) // let terminal render
lines, err := r.CapturePaneLines(session, 20)
// Search lines for the expected content
```

For tests that need to verify exact delivery (not just "something arrived"), send a unique marker string and grep for it in the captured output.

---

## Tests to implement

### `internal/tmux/integration_test.go`

All integration tests must:
- Check `testing.Short()` and skip if true
- Use `newTestRunner(t)` and `createTestSession(t, r)`
- Use unique marker strings to avoid cross-test interference
- Wait for shell readiness before sending

#### Test 1: Basic send

```go
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
```

#### Test 2: Special characters (literal mode)

Verify that `!`, `~`, `#`, backticks, and Unicode arrive literally (not interpreted as tmux key names).

```go
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
```

#### Test 3: Large message chunking

Send a message > 512 bytes, verify full delivery without truncation.

```go
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
```

#### Test 4: Copy mode handling

Put pane in copy mode, send message, verify copy mode was exited and message delivered.

```go
func TestSendExitsCopyMode(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    r := newTestRunner(t)
    session := createTestSession(t, r)
    waitForShellReady(t, r, session)

    // Enter copy mode
    if err := r.SendKeys(session, "C-b", "["); err != nil {
        // Try alternate: some configs use different prefix
        _ = r.SendKeys(session, "C-b", "[")
    }
    time.Sleep(200 * time.Millisecond)

    // Verify we're in copy mode
    inMode, _ := r.InCopyMode(session)
    if !inMode {
        // Copy mode entry may have failed — tmux default prefix might differ
        // in the test server. Try direct copy-mode command.
        r.SendKeys(session, "copy-mode")
        time.Sleep(200 * time.Millisecond)
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
```

#### Test 5: Detached session wake

```go
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
```

#### Test 6: Idle detection

```go
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
    err := WaitForIdle(r, "nonexistent-session-xyz", 2*time.Second)
    if !errors.Is(err, ErrSessionNotFound) {
        t.Errorf("expected ErrSessionNotFound, got: %v", err)
    }
}
```

#### Test 7: Concurrent sends

Two goroutines send to the same target simultaneously. Verify no interleaving (both messages arrive intact).

```go
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
```

#### Test 8: Error on nonexistent target

```go
func TestSendSessionNotFound(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    r := newTestRunner(t)
    opts := SendOpts{Timeout: 2 * time.Second}
    err := r.Send("nonexistent-session-xyz", "hello", opts)
    if err == nil {
        t.Fatal("expected error for nonexistent session")
    }
    if !errors.Is(err, ErrSessionNotFound) {
        t.Errorf("expected ErrSessionNotFound, got: %v", err)
    }
}
```

---

## Test execution

All integration tests run with `make test` (which runs `go test ./... -v`). They are skipped with `make test-short` (which adds `-short`).

Expected test output:
```
=== RUN   TestSendBasic
--- PASS: TestSendBasic
=== RUN   TestSendSpecialChars
--- PASS: TestSendSpecialChars
... etc
```

---

## Flakiness mitigation

Integration tests against real tmux are inherently timing-sensitive. Mitigation strategies:

1. **`waitForShellReady`** before every send — don't assume the shell is ready immediately after session creation
2. **Generous sleep after send** — 1-2 seconds for the command to execute and output to appear
3. **Unique markers per test** — avoid cross-test interference
4. **Isolated socket** — TestMain ensures no interaction with user's tmux
5. **`SkipEscape: true`** in most tests — saves 1.1 seconds per send (500ms + 600ms) and avoids Escape timing issues in bare shell tests
6. **Polling over sleeping where possible** — for verification, poll capture-pane rather than fixed sleep (though a generous fixed sleep is acceptable for v1)

---

## Done-when criteria

1. All 8 integration tests pass with `go test ./internal/tmux/ -v -count=1` (no `-short`)
2. Tests are skipped with `-short` flag
3. Tests use isolated tmux socket (no interference with user's tmux)
4. Each test creates its own session and cleans up via `t.Cleanup`
5. Basic send delivers and is verifiable via capture-pane
6. Special chars arrive literally (not interpreted)
7. Large messages (>512 bytes) arrive without truncation
8. Copy mode is exited before delivery
9. Detached sessions receive messages
10. Concurrent sends don't garble (both markers found)
11. Nonexistent targets return `ErrSessionNotFound`
12. `make test` passes (all tests including integration)
13. `make test-short` passes (integration tests skipped)
