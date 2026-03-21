# T1: Tmux Runner — Foundation Layer

**Dependencies:** None (this is the foundation)
**Files to create:** `internal/tmux/tmux.go`, `internal/tmux/testmain_test.go`
**Files to modify:** None

---

## Goal

Create the `Runner` struct that wraps `exec.Command("tmux", ...)` with UTF-8 flag, socket support, and structured error wrapping. This is the foundation every other file depends on — all tmux interaction goes through `Runner.run()`.

Also create the `TestMain` setup for an isolated tmux socket so integration tests never touch the user's running tmux.

---

## What to implement

### `internal/tmux/tmux.go`

#### Sentinel errors

Define these package-level sentinel errors:

```go
var (
    ErrNoServer        = errors.New("no tmux server running")
    ErrSessionNotFound = errors.New("session not found")
    ErrIdleTimeout     = errors.New("agent not idle before timeout")
)
```

`[G:tmux.go:42-47]` — Gastown defines these same sentinels. We only need the three above for whisper's scope.

#### Runner struct

```go
type Runner struct {
    socket string // tmux socket name (-L flag), empty = default
}
```

Constructor:

```go
func NewRunner() *Runner              // uses default socket (empty string)
func NewRunnerWithSocket(s string) *Runner  // for tests + explicit socket targeting
```

`[G:tmux.go:168-204]` — Gastown's `Tmux` struct. We call ours `Runner` to avoid confusion with the `tmux` package name.

#### The `run` method (unexported)

```go
func (r *Runner) run(args ...string) (string, error)
```

This is the core. Every tmux command goes through it.

1. Build args: prepend `-u` (UTF-8 mode, always). If `r.socket != ""`, prepend `-L <socket>` before the subcommand.
2. Execute via `exec.Command("tmux", allArgs...)`, capture stdout and stderr in `bytes.Buffer`.
3. On success: return `strings.TrimSpace(stdout)`, nil.
4. On error: call `wrapError(err, stderr, args)` to classify the error.

`[G:tmux.go:209-228]` — The exact pattern. Copy the structure, not the code.

**Why `-u`:** Forces UTF-8 regardless of locale. Without it, tmux silently corrupts Unicode characters on some systems. `[G:tmux.go:208 comment]`

#### Error classification (`wrapError`)

```go
func (r *Runner) wrapError(err error, stderr string, args []string) error
```

Classify stderr into sentinel errors:

| stderr contains | Return |
|----------------|--------|
| `"no server running"` | `ErrNoServer` |
| `"error connecting to"` | `ErrNoServer` |
| `"server exited unexpectedly"` | `ErrNoServer` |
| `"session not found"` | `ErrSessionNotFound` |
| `"can't find session"` | `ErrSessionNotFound` |
| (other non-empty stderr) | `fmt.Errorf("tmux %s: %s", args[0], stderr)` |
| (empty stderr) | `fmt.Errorf("tmux %s: %w", args[0], err)` |

`[G:tmux.go:231-253]` — The exact classification logic. We skip `ErrSessionExists` and `ErrInvalidSessionName` (whisper doesn't create sessions). We also skip the `"no current target"` → `ErrNoServer` mapping — keep it simple.

#### Exported helper methods on Runner

These thin wrappers expose specific tmux commands that other packages need:

```go
// SendKeys sends raw key names (e.g., "Enter", "Escape") to a target.
func (r *Runner) SendKeys(target string, keys ...string) error {
    args := append([]string{"send-keys", "-t", target}, keys...)
    _, err := r.run(args...)
    return err
}

// SendKeysLiteral sends literal text (send-keys -l) to a target.
func (r *Runner) SendKeysLiteral(target, text string) error {
    _, err := r.run("send-keys", "-t", target, "-l", text)
    return err
}

// DisplayMessage runs display-message -p and returns the formatted output.
func (r *Runner) DisplayMessage(target, format string) (string, error) {
    return r.run("display-message", "-p", "-t", target, format)
}

// CapturePaneLines captures the last N lines of a pane as a slice.
func (r *Runner) CapturePaneLines(target string, lines int) ([]string, error) {
    out, err := r.run("capture-pane", "-p", "-t", target, "-S", fmt.Sprintf("-%d", lines))
    if err != nil {
        return nil, err
    }
    if out == "" {
        return nil, nil
    }
    return strings.Split(out, "\n"), nil
}

// ResizeWindow resizes the window containing target to width x.
func (r *Runner) ResizeWindow(target string, width int) error {
    _, err := r.run("resize-window", "-t", target, "-x", fmt.Sprintf("%d", width))
    return err
}

// GetWindowWidth returns the current width of the window containing target.
func (r *Runner) GetWindowWidth(target string) (int, error) {
    out, err := r.run("display-message", "-p", "-t", target, "#{window_width}")
    if err != nil {
        return 0, err
    }
    var w int
    if _, err := fmt.Sscanf(strings.TrimSpace(out), "%d", &w); err != nil {
        return 0, fmt.Errorf("parse window width %q: %w", out, err)
    }
    return w, nil
}

// SetWindowOption sets a window option on the target.
func (r *Runner) SetWindowOption(target, option, value string) error {
    _, err := r.run("set-option", "-w", "-t", target, option, value)
    return err
}

// IsSessionAttached returns true if the session has any clients attached.
func (r *Runner) IsSessionAttached(target string) bool {
    out, err := r.run("display-message", "-t", target, "-p", "#{session_attached}")
    return err == nil && strings.TrimSpace(out) == "1"
}

// CancelCopyMode sends cancel to exit copy/scroll mode.
func (r *Runner) CancelCopyMode(target string) error {
    _, err := r.run("send-keys", "-t", target, "-X", "cancel")
    return err
}

// InCopyMode returns true if the target pane is in copy/scroll mode.
func (r *Runner) InCopyMode(target string) (bool, error) {
    out, err := r.run("display-message", "-p", "-t", target, "#{pane_in_mode}")
    if err != nil {
        return false, err
    }
    return strings.TrimSpace(out) == "1", nil
}
```

These keep the rest of the codebase from calling `run` directly (it's unexported) while staying thin enough that they're not worth testing independently — they're tested through integration tests.

### `internal/tmux/testmain_test.go`

```go
package tmux

import (
    "fmt"
    "os"
    "os/exec"
    "testing"
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
    // Create a detached session running a simple shell
    args := []string{"-u", "-L", testSocket, "new-session", "-d", "-s", name, "-x", "200", "-y", "50"}
    cmd := exec.Command("tmux", args...)
    if out, err := cmd.CombinedOutput(); err != nil {
        t.Fatalf("create test session %q: %v\n%s", name, err, out)
    }
    t.Cleanup(func() {
        _ = exec.Command("tmux", "-L", testSocket, "kill-session", "-t", name).Run()
    })
    return name
}
```

`[G:testmain_test.go]` — Gastown's pattern. Key differences:
- We use `"whisper-test-PID"` prefix instead of `"gt-test-PID"`.
- We don't use `SetDefaultSocket` (Gastown's global state pattern) — instead we pass the socket explicitly via `newTestRunner`.
- We provide `createTestSession` as a test helper that other test files will use.

---

## Test requirements

### Unit test: `internal/tmux/tmux_test.go`

Test `wrapError` classification:

```go
func TestWrapError(t *testing.T) {
    r := NewRunner()
    tests := []struct {
        name   string
        stderr string
        want   error
    }{
        {"no server", "no server running on /tmp/...", ErrNoServer},
        {"connection error", "error connecting to /tmp/...", ErrNoServer},
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
```

**Important:** `wrapError` needs to be testable. Since it's unexported, the test lives in the same package (`package tmux`, not `package tmux_test`).

### Integration test: `internal/tmux/runner_integration_test.go`

```go
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
    _, err := r.DisplayMessage("nonexistent-session-xyz", "#{session_name}")
    if !errors.Is(err, ErrSessionNotFound) {
        t.Errorf("expected ErrSessionNotFound, got %v", err)
    }
}
```

---

## Done-when criteria

1. `Runner` struct exists with `NewRunner()` and `NewRunnerWithSocket()` constructors
2. `run()` method prepends `-u` and `-L socket` flags, captures stdout/stderr, classifies errors
3. `wrapError()` maps stderr strings to `ErrNoServer` and `ErrSessionNotFound` sentinels
4. All exported helper methods compile and delegate to `run()`
5. `TestMain` creates an isolated tmux socket per test run
6. `createTestSession` helper creates/cleans up test sessions
7. `go vet ./internal/tmux/...` passes
8. Unit tests for error classification pass
9. Integration tests pass (create session, query it, verify error classification against real tmux)
