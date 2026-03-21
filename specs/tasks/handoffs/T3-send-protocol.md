# T3: Send Protocol — The 8-Step Core

**Dependencies:** T1 (Runner), T2 (Sanitize, Lock)
**Files to create:** `internal/tmux/send.go`, `internal/tmux/wake.go`, `internal/tmux/capture.go`
**Files to modify:** None

---

## Goal

Implement the 8-step send protocol adapted from Gastown's `NudgeSessionWithOpts`. This is the core of whisper — every `whisper send` executes these steps. Also implement `WakePane`/`WakePaneIfDetached` for detached session handling.

---

## What to implement

### `internal/tmux/send.go`

#### Constants

```go
const (
    sendKeysChunkSize  = 512                  // [G:tmux.go:1314] bytes per chunk
    postSendDelay      = 500 * time.Millisecond // [G:tmux.go:1453] empirically required
    escEnterGap        = 600 * time.Millisecond // [G:tmux.go:1460] must exceed readline keyseq-timeout (500ms)
    enterRetryDelay    = 200 * time.Millisecond // [G:tmux.go:1471] between Enter retries
    enterMaxAttempts   = 3                      // [G:tmux.go:1469]
    copyModeCancelWait = 50 * time.Millisecond  // [G:tmux.go:1440]
    interChunkDelay    = 10 * time.Millisecond  // [G:tmux.go:1340]
    nudgeRetryInterval = 500 * time.Millisecond // [G:constants:54] base retry backoff
    maxRetryInterval   = 2 * time.Second        // [G:tmux.go:1386] backoff cap
)
```

**CRITICAL:** These timing values are empirically derived from production. Do not change them without evidence.

#### SendOpts

```go
type SendOpts struct {
    SkipEscape bool          // omit Escape + 600ms delay
    Timeout    time.Duration // max time for send-keys retry (default 15s)
}
```

#### `Send(target, message string, opts SendOpts) error`

The 8-step protocol. Reference `[G:tmux.go:1421-1482]` (NudgeSessionWithOpts) step by step:

**Step 1: Acquire send lock** `[G:tmux.go:1424]`
```go
if !AcquireLock(target, LockTimeout) {
    return fmt.Errorf("send lock timeout for %q: previous send may be hung", target)
}
defer ReleaseLock(target)
```

**Step 2: Exit copy mode** `[G:tmux.go:1438]`
```go
inMode, err := r.InCopyMode(target)
if err == nil && inMode {
    _ = r.CancelCopyMode(target)
    time.Sleep(copyModeCancelWait)
}
```
Errors from `InCopyMode` are non-fatal — if we can't check, proceed anyway.

**Step 3: Sanitize** `[G:tmux.go:1444]`
```go
sanitized := Sanitize(message)
```

**Step 4: Send text with chunking + retry** `[G:tmux.go:1448]`
```go
if err := r.sendMessageToTarget(target, sanitized, opts.Timeout); err != nil {
    return err
}
```

**Step 5: Post-send delay** `[G:tmux.go:1453]`
```go
time.Sleep(postSendDelay)
```

**Step 6: Send Escape** `[G:tmux.go:1458]` (skipped if `opts.SkipEscape`)
```go
if !opts.SkipEscape {
    _ = r.SendKeys(target, "Escape")
```

**Step 7: ESC→Enter gap** `[G:tmux.go:1460]` (skipped if `opts.SkipEscape`)
```go
    time.Sleep(escEnterGap)
}
```
**Why 600ms:** Bash readline's `keyseq-timeout` defaults to 500ms. ESC + Enter within that window becomes M-Enter (meta-return), which does nothing. 600ms ensures ESC is processed alone before Enter arrives.

**Step 8: Send Enter with retry** `[G:tmux.go:1469-1481]`
```go
var lastErr error
for attempt := 0; attempt < enterMaxAttempts; attempt++ {
    if attempt > 0 {
        time.Sleep(enterRetryDelay)
    }
    if err := r.SendKeys(target, "Enter"); err != nil {
        lastErr = err
        continue
    }
    // Step 8b: Wake detached sessions
    WakePaneIfDetached(r, target)
    return nil
}
return fmt.Errorf("failed to send Enter after %d attempts: %w", enterMaxAttempts, lastErr)
```

#### `sendMessageToTarget(target, text string, timeout time.Duration) error` (unexported)

`[G:tmux.go:1316-1344]`

```go
func (r *Runner) sendMessageToTarget(target, text string, timeout time.Duration) error {
    if len(text) <= sendKeysChunkSize {
        return r.sendKeysLiteralWithRetry(target, text, timeout)
    }
    for i := 0; i < len(text); i += sendKeysChunkSize {
        end := i + sendKeysChunkSize
        if end > len(text) {
            end = len(text)
        }
        chunk := text[i:end]
        if i == 0 {
            if err := r.sendKeysLiteralWithRetry(target, chunk, timeout); err != nil {
                return err
            }
        } else {
            if err := r.SendKeysLiteral(target, chunk); err != nil {
                return err
            }
        }
        if end < len(text) {
            time.Sleep(interChunkDelay)
        }
    }
    return nil
}
```

Key points:
- First chunk uses retry logic (handles startup race)
- Subsequent chunks use plain `SendKeysLiteral` (no retry — if the first chunk worked, the target is ready)
- 10ms delay between chunks to avoid overwhelming the TTY buffer
- Chunking at 512 bytes — Linux canonical mode TTY buffer is 4096 bytes, but chunking smaller is safer

#### `sendKeysLiteralWithRetry(target, text string, timeout time.Duration) error` (unexported)

`[G:tmux.go:1359-1391]`

```go
func (r *Runner) sendKeysLiteralWithRetry(target, text string, timeout time.Duration) error {
    deadline := time.Now().Add(timeout)
    interval := nudgeRetryInterval
    var lastErr error

    for time.Now().Before(deadline) {
        err := r.SendKeysLiteral(target, text)
        if err == nil {
            return nil
        }
        if !isTransientSendKeysError(err) {
            return err // non-transient — fail fast
        }
        lastErr = err
        remaining := time.Until(deadline)
        if remaining <= 0 {
            break
        }
        sleep := interval
        if sleep > remaining {
            sleep = remaining
        }
        time.Sleep(sleep)
        // Exponential backoff: 500ms → 750ms → 1125ms → 1687ms → 2s (cap)
        interval = interval * 3 / 2
        if interval > maxRetryInterval {
            interval = maxRetryInterval
        }
    }
    return fmt.Errorf("target not ready for input after %s: %w", timeout, lastErr)
}
```

#### `isTransientSendKeysError(err error) bool` (unexported)

`[G:tmux.go:1272-1278]`

```go
func isTransientSendKeysError(err error) bool {
    if err == nil {
        return false
    }
    return strings.Contains(err.Error(), "not in a mode")
}
```

Only `"not in a mode"` is retryable — this happens when the TUI hasn't initialized its input handling during cold startup. All other errors (session gone, no server) should fail immediately.

### `internal/tmux/wake.go`

#### `WakePane(r *Runner, target string)`

`[G:tmux.go:1226-1257]`

Triggers SIGWINCH by doing a resize dance:

```go
func WakePane(r *Runner, target string) {
    width, err := r.GetWindowWidth(target)
    if err != nil || width < 1 {
        return // session may be dead
    }
    // Bump width +1, wait 50ms, restore
    _ = r.ResizeWindow(target, width+1)
    time.Sleep(50 * time.Millisecond)
    _ = r.ResizeWindow(target, width)

    // Reset window-size to "latest" — resize-window sets it to "manual"
    // which permanently locks dimensions and causes dot borders when
    // a human later attaches with a different terminal size.
    // [G:tmux.go:1251-1256]
    _ = r.SetWindowOption(target, "window-size", "latest")
}
```

**Why resize-window, not resize-pane:** `[G:tmux.go:1228]` `resize-pane` doesn't work on single-pane sessions because the pane already fills the window. `resize-window` changes the window dimensions, which sends SIGWINCH to all processes.

**Why restore window-size to "latest":** `[G:tmux.go:1251]` tmux 3.3+ automatically sets `window-size` to `"manual"` whenever `resize-window` is called. This prevents the window from auto-sizing when a human attaches, causing "dots around the edges" visual artifacts.

#### `WakePaneIfDetached(r *Runner, target string)`

`[G:tmux.go:1262-1267]`

```go
func WakePaneIfDetached(r *Runner, target string) {
    if r.IsSessionAttached(target) {
        return
    }
    WakePane(r, target)
}
```

Skip the wake on attached sessions — Claude is already processing terminal events, and the resize dance adds unnecessary latency.

### `internal/tmux/capture.go`

This file is mostly a re-export of Runner methods needed by idle detection (T4). Since `CapturePaneLines` is already on `Runner` from T1, this file contains the prompt-matching logic.

#### Prompt detection constants

```go
// DefaultPromptPrefix is Claude Code's idle prompt character.
// ❯ is U+276F HEAVY RIGHT-POINTING ANGLE QUOTATION MARK.
const DefaultPromptPrefix = "❯"
```

`[G:tmux.go:2462]` — Gastown uses `"❯ "` (with trailing space). We use just `"❯"` and do prefix matching, which handles both with and without trailing space.

#### `MatchesPromptPrefix(line string) bool`

`[G:tmux.go:2410-2421]`

```go
func MatchesPromptPrefix(line string) bool {
    trimmed := strings.TrimSpace(line)
    // Normalize NBSP (U+00A0) → regular space
    trimmed = strings.ReplaceAll(trimmed, "\u00a0", " ")
    return strings.HasPrefix(trimmed, DefaultPromptPrefix)
}
```

**Why NBSP normalization:** `[G:tmux.go:2417]` Some terminal renderers use NBSP (U+00A0) instead of regular space. Without normalization, prompt matching silently fails.

#### `IsStatusBarBusy(lines []string) bool`

Extract the status bar check from `WaitForIdle` so it's independently testable:

```go
func IsStatusBarBusy(lines []string) bool {
    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if strings.Contains(trimmed, "⏵⏵") {
            if strings.Contains(trimmed, "esc to interrupt") {
                return true
            }
            break
        }
    }
    return false
}
```

`[G:tmux.go:2500-2508]` — If lines contain `"⏵⏵"` AND `"esc to interrupt"`, Claude Code is actively running a tool call.

---

## Test requirements

### Unit tests (in `internal/tmux/`)

**`send_test.go`:**
- `TestIsTransientSendKeysError`: verify `"not in a mode"` returns true, other errors return false, nil returns false
- `TestChunking`: verify that `sendMessageToTarget` would chunk correctly — this is hard to unit test without mocking. Instead, test the chunking math: for input lengths of 1, 512, 513, 1024, 2000, verify expected chunk count. Write a helper `chunkCount(msgLen int) int` if needed, or just validate in integration tests.

**`capture_test.go`:**
- `TestMatchesPromptPrefix`: `"❯ hello"` → true, `"  ❯ "` → true, `"$ hello"` → false, `""` → false, `"❯"` → true, `"text ❯"` → false (not prefix), NBSP variant → true
- `TestIsStatusBarBusy`: lines with `"⏵⏵"` + `"esc to interrupt"` → true, lines with only `"⏵⏵"` → false, empty lines → false, lines with only `"esc to interrupt"` → false

### Integration tests (deferred to T6)

The 8-step protocol is best tested end-to-end against real tmux sessions. T6 covers: basic send, special chars, large message chunking, copy mode handling, concurrent sends, and detached session wake.

---

## Done-when criteria

1. `Send()` implements all 8 steps in order with correct timing constants
2. `sendMessageToTarget()` chunks messages > 512 bytes with 10ms inter-chunk delays
3. `sendKeysLiteralWithRetry()` uses exponential backoff (500ms → 750ms → 1125ms → 1687ms → 2s cap)
4. `isTransientSendKeysError()` only retries on `"not in a mode"`
5. `WakePane()` does resize dance: get width → bump +1 → 50ms → restore → set window-size latest
6. `WakePaneIfDetached()` skips attached sessions
7. `MatchesPromptPrefix()` handles NBSP normalization and prefix matching
8. `IsStatusBarBusy()` detects `"⏵⏵"` + `"esc to interrupt"` pattern
9. All unit tests pass
10. `go vet ./internal/tmux/...` passes
