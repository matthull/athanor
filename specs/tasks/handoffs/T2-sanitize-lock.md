# T2: Sanitize + Lock — Pure Logic Units

**Dependencies:** T1 (Runner struct must exist, but these files don't call it)
**Files to create:** `internal/tmux/sanitize.go`, `internal/tmux/lock.go`, `internal/tmux/sanitize_test.go`, `internal/tmux/lock_test.go`
**Files to modify:** None

---

## Goal

Implement two independently testable pure-logic modules:
1. **Sanitize** — strip control characters from messages before tmux delivery
2. **Lock** — per-target channel-based semaphore with timed acquisition

Both are pure Go (no tmux subprocess calls), so they get full unit test coverage with `t.Parallel()`.

---

## What to implement

### `internal/tmux/sanitize.go`

#### `Sanitize(msg string) string`

Strip/replace control characters that corrupt tmux `send-keys -l` delivery.

Rules (apply per-rune, iterating with `range`):

| Character | Action | Why |
|-----------|--------|-----|
| TAB (0x09) | Replace with single space | Triggers shell tab completion |
| Newline (0x0a) | **Preserve** | `send-keys -l` handles newlines correctly |
| All other chars < 0x20 | **Strip** | ESC (0x1b) triggers terminal escape sequences, CR (0x0d) acts as premature Enter, BS (0x08) deletes characters |
| DEL (0x7f) | **Strip** | Deletes characters |
| Everything else | **Preserve** | Including quotes, backticks, Unicode, all printable ASCII |

Implementation pattern:

```go
func Sanitize(msg string) string {
    var b strings.Builder
    b.Grow(len(msg))
    for _, r := range msg {
        switch {
        case r == '\t':
            b.WriteRune(' ')
        case r == '\n':
            b.WriteRune(r)
        case r < 0x20:
            continue
        case r == 0x7f:
            continue
        default:
            b.WriteRune(r)
        }
    }
    return b.String()
}
```

`[G:tmux.go:1285-1303]` — This is a direct translation of `sanitizeNudgeMessage`. The logic is identical; we just export it and give it a cleaner name.

### `internal/tmux/lock.go`

#### Per-target send lock

Prevents concurrent sends to the same target from interleaving text. Uses channel-based semaphores (not `sync.Mutex`) to support timed acquisition.

`[G:tmux.go:24-35, 1180-1209]` — Gastown's `sessionNudgeLocks` pattern.

**Constants:**

```go
const LockTimeout = 30 * time.Second // [G:tmux.go:35]
```

**Package-level state:**

```go
var targetLocks sync.Map // map[string]chan struct{}
```

**Functions:**

```go
// getTargetSem returns the channel semaphore for a target.
// Creates a new buffered channel of size 1 if one doesn't exist.
func getTargetSem(target string) chan struct{} {
    sem := make(chan struct{}, 1)
    actual, _ := targetLocks.LoadOrStore(target, sem)
    return actual.(chan struct{})
}

// AcquireLock attempts to acquire the per-target send lock.
// Returns true if acquired within timeout, false if timeout expired.
func AcquireLock(target string, timeout time.Duration) bool {
    sem := getTargetSem(target)
    select {
    case sem <- struct{}{}:
        return true
    case <-time.After(timeout):
        return false
    }
}

// ReleaseLock releases the per-target send lock.
func ReleaseLock(target string) {
    sem := getTargetSem(target)
    select {
    case <-sem:
    default:
        // Lock wasn't held — shouldn't happen, but don't block
    }
}
```

`[G:tmux.go:1183-1209]` — Direct translation. Key design decisions preserved:
- `sync.Map` for lock-free concurrent access to the semaphore map
- `LoadOrStore` ensures exactly one semaphore per target (race-safe)
- `ReleaseLock` has a `default` case to prevent blocking if called without holding the lock
- Channel-based (not mutex) so `AcquireLock` can timeout instead of blocking forever

**Why channel-based, not mutex:** `[G:tmux.go:29]` A hung tmux command holding a `sync.Mutex` would block all future sends to that target permanently. With channels, callers can give up after `LockTimeout` and report the problem.

---

## Test requirements

### `internal/tmux/sanitize_test.go`

All tests should use `t.Parallel()`.

```go
func TestSanitize(t *testing.T) {
    t.Parallel()
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"plain text", "hello world", "hello world"},
        {"preserves newlines", "line1\nline2", "line1\nline2"},
        {"tab to space", "hello\tworld", "hello world"},
        {"strips ESC", "hello\x1bworld", "helloworld"},
        {"strips CR", "hello\rworld", "helloworld"},
        {"strips BS", "hello\x08world", "helloworld"},
        {"strips DEL", "hello\x7fworld", "helloworld"},
        {"strips null", "hello\x00world", "helloworld"},
        {"preserves quotes", `he said "hello"`, `he said "hello"`},
        {"preserves backticks", "run `cmd`", "run `cmd`"},
        {"preserves special chars", "!~#$%^&*()", "!~#$%^&*()"},
        {"preserves unicode", "hello 世界 🌍", "hello 世界 🌍"},
        {"mixed control chars", "a\x01b\x02c\x1bd\x7fe\tf\ng", "abcde f\ng"},
        {"empty string", "", ""},
        {"only control chars", "\x01\x02\x03\x1b\x7f", ""},
        {"preserves prompt char", "❯ hello", "❯ hello"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            got := Sanitize(tt.input)
            if got != tt.want {
                t.Errorf("Sanitize(%q) = %q, want %q", tt.input, got, tt.want)
            }
        })
    }
}
```

### `internal/tmux/lock_test.go`

```go
func TestAcquireRelease(t *testing.T) {
    t.Parallel()
    target := "test-lock-basic"
    if !AcquireLock(target, time.Second) {
        t.Fatal("should acquire uncontested lock")
    }
    ReleaseLock(target)
    // Should be able to re-acquire after release
    if !AcquireLock(target, time.Second) {
        t.Fatal("should re-acquire after release")
    }
    ReleaseLock(target)
}

func TestAcquireLockTimeout(t *testing.T) {
    t.Parallel()
    target := "test-lock-timeout"
    if !AcquireLock(target, time.Second) {
        t.Fatal("first acquire should succeed")
    }
    defer ReleaseLock(target)

    // Second acquire should timeout
    start := time.Now()
    if AcquireLock(target, 50*time.Millisecond) {
        t.Fatal("should not acquire contested lock")
    }
    elapsed := time.Since(start)
    if elapsed < 40*time.Millisecond {
        t.Errorf("timeout too fast: %v", elapsed)
    }
}

func TestDifferentTargetsIndependent(t *testing.T) {
    t.Parallel()
    target1 := "test-lock-ind-1"
    target2 := "test-lock-ind-2"
    if !AcquireLock(target1, time.Second) {
        t.Fatal("should acquire target1")
    }
    defer ReleaseLock(target1)
    // target2 should be independently acquirable
    if !AcquireLock(target2, time.Second) {
        t.Fatal("should acquire target2 independently")
    }
    ReleaseLock(target2)
}

func TestReleaseWithoutAcquire(t *testing.T) {
    t.Parallel()
    // Should not panic or block
    ReleaseLock("test-lock-never-held")
}
```

---

## Done-when criteria

1. `Sanitize()` strips all control chars < 0x20 except TAB (→space) and newline (preserved)
2. `Sanitize()` strips DEL (0x7f)
3. `Sanitize()` preserves all printable ASCII, quotes, backticks, Unicode
4. `AcquireLock()` returns true immediately for uncontested lock
5. `AcquireLock()` returns false after timeout for contested lock
6. `ReleaseLock()` makes the lock re-acquirable
7. Different targets have independent locks
8. `ReleaseLock()` on unheld lock doesn't panic or block
9. All tests pass with `go test ./internal/tmux/ -run "TestSanitize|TestAcquire|TestDifferent|TestRelease" -v`
10. `go vet ./internal/tmux/...` passes
