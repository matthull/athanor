# T4: Idle Detection — 2-Check Protocol

**Dependencies:** T1 (Runner — for `CapturePaneLines`), T3 (capture.go — for `MatchesPromptPrefix`, `IsStatusBarBusy`)
**Files to create:** `internal/tmux/idle.go`, `internal/tmux/idle_test.go`
**Files to modify:** None

---

## Goal

Implement the 2-check idle detection protocol adapted from Gastown's `WaitForIdle`. This polls a tmux pane, checks for Claude Code's idle prompt, and requires 2 consecutive idle readings 200ms apart to confirm genuine idle state (filtering out transient prompt flashes between tool calls).

---

## The problem this solves

Claude Code briefly shows its prompt (`❯`) between tool calls — a ~500ms gap. A single idle check produces false positives: the agent appears idle but is actually mid-execution. The 2-check protocol creates a 400ms observation window that catches these transient states.

`[G:tmux.go:2474-2477]` — "Two polls 200ms apart confirms genuine idle state."

---

## What to implement

### `internal/tmux/idle.go`

#### Constants

```go
const (
    idlePollInterval      = 200 * time.Millisecond // [G:tmux.go:2538]
    requiredConsecutiveIdle = 2                      // [G:tmux.go:2479]
    idleCaptureLines       = 5                      // [G:tmux.go:2483]
)
```

#### `WaitForIdle(r *Runner, target string, timeout time.Duration) error`

`[G:tmux.go:2470-2541]` — Direct adaptation of Gastown's WaitForIdle.

```go
func WaitForIdle(r *Runner, target string, timeout time.Duration) error {
    consecutiveIdle := 0
    deadline := time.Now().Add(timeout)

    for time.Now().Before(deadline) {
        // Step 1: Capture last 5 lines
        lines, err := r.CapturePaneLines(target, idleCaptureLines)
        if err != nil {
            // Terminal errors — no point polling further
            if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrNoServer) {
                return err
            }
            // Transient capture error — reset and retry
            consecutiveIdle = 0
            time.Sleep(idlePollInterval)
            continue
        }

        // Step 2: Check status bar — if "⏵⏵" + "esc to interrupt" visible,
        // agent is actively running a tool call. NOT idle.
        if IsStatusBarBusy(lines) {
            consecutiveIdle = 0
            time.Sleep(idlePollInterval)
            continue
        }

        // Step 3: Scan for prompt prefix (❯)
        promptFound := false
        for _, line := range lines {
            if strings.TrimSpace(line) == "" {
                continue
            }
            if MatchesPromptPrefix(line) {
                promptFound = true
                break
            }
        }

        // Step 4-6: Count consecutive idle readings
        if promptFound {
            consecutiveIdle++
            if consecutiveIdle >= requiredConsecutiveIdle {
                return nil // Confirmed idle
            }
        } else {
            consecutiveIdle = 0
        }

        time.Sleep(idlePollInterval)
    }

    return ErrIdleTimeout
}
```

**Protocol walkthrough:**

1. **Capture last 5 lines** of target pane via `Runner.CapturePaneLines`. 5 lines is enough to see the prompt + status bar. `[G:tmux.go:2483]`

2. **Check status bar first.** If lines contain `"⏵⏵"` AND `"esc to interrupt"`, Claude Code is actively running a tool call. NOT idle — regardless of whether the prompt prefix is also visible. Reset consecutive count, sleep 200ms, continue. `[G:tmux.go:2499-2513]`

3. **Scan for prompt prefix.** Look for `❯` (U+276F) at line start in any of the captured lines. Claude Code renders a status bar below the prompt line, so the prompt may not be the last non-empty line. `[G:tmux.go:2515-2528]`

4. **If prompt found:** increment consecutive idle count. `[G:tmux.go:2530-2531]`

5. **If consecutive count reaches 2:** return nil (confirmed idle). `[G:tmux.go:2532-2534]`

6. **If prompt NOT found:** reset consecutive count to 0. `[G:tmux.go:2535-2536]`

7. **Sleep 200ms**, repeat from step 1. `[G:tmux.go:2538]`

**Error handling:**
- `ErrSessionNotFound` / `ErrNoServer` → return immediately (session is gone, no recovery possible)
- Other capture errors → transient, reset count and keep polling
- Timeout → return `ErrIdleTimeout`

---

## Test requirements

### Unit test: `internal/tmux/idle_test.go`

The idle detection loop itself requires real tmux (integration test in T6). But we can unit test the decision logic by extracting it or testing the components:

```go
func TestIdleDetectionLogic(t *testing.T) {
    t.Parallel()

    t.Run("status bar busy prevents idle", func(t *testing.T) {
        t.Parallel()
        lines := []string{
            "some output",
            "⏵⏵ Running tool                    esc to interrupt",
            "",
        }
        if !IsStatusBarBusy(lines) {
            t.Error("should detect busy status bar")
        }
    })

    t.Run("prompt without busy status is idle candidate", func(t *testing.T) {
        t.Parallel()
        lines := []string{
            "previous output",
            "❯ ",
            "  ○ 3 functions | 12 subtasks",
        }
        if IsStatusBarBusy(lines) {
            t.Error("should not be busy")
        }
        found := false
        for _, line := range lines {
            if MatchesPromptPrefix(line) {
                found = true
                break
            }
        }
        if !found {
            t.Error("should find prompt prefix")
        }
    })

    t.Run("no prompt means not idle", func(t *testing.T) {
        t.Parallel()
        lines := []string{
            "Compiling...",
            "Building module graph",
            "50% complete",
        }
        if IsStatusBarBusy(lines) {
            t.Error("should not be busy")
        }
        found := false
        for _, line := range lines {
            if MatchesPromptPrefix(line) {
                found = true
                break
            }
        }
        if found {
            t.Error("should not find prompt prefix in build output")
        }
    })
}
```

### Integration tests (deferred to T6)

T6 will test:
- Start a `sleep` command in a pane (busy), verify `WaitForIdle` waits
- Let sleep finish (prompt returns), verify idle detected
- Verify 2-check doesn't false-positive on brief prompt flash

---

## Done-when criteria

1. `WaitForIdle()` polls with 200ms interval and requires 2 consecutive idle readings
2. Status bar busy (`"⏵⏵"` + `"esc to interrupt"`) resets the consecutive count
3. `ErrSessionNotFound` / `ErrNoServer` cause immediate return (no polling)
4. Transient capture errors reset count but continue polling
5. Returns `ErrIdleTimeout` when timeout expires
6. Returns `nil` on confirmed idle (2 consecutive prompt detections without busy status bar)
7. Unit tests for decision logic pass
8. `go vet ./internal/tmux/...` passes
