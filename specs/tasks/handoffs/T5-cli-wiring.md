# T5: CLI Wiring — Command Implementations

**Dependencies:** T1 (Runner), T2 (Lock), T3 (Send), T4 (Idle)
**Files to create:** `internal/cli/send_cmd.go`, `internal/cli/idle_cmd.go`, `internal/cli/wait_send_cmd.go`
**Files to modify:** `internal/cli/root.go` (remove stubs)

---

## Goal

Wire CLI argument parsing to the `internal/tmux` functions. Each command gets its own file. Remove the stub implementations from `root.go`. The CLI uses a simple `switch` dispatch (already scaffolded) — no cobra, no framework.

---

## What to implement

### Modify `internal/cli/root.go`

Remove the three stub functions at the bottom of the file (lines 62-76):

```go
// DELETE these stubs — they move to their own files:
func runSend(args []string) int { ... }
func runIdle(args []string) int { ... }
func runWaitAndSend(args []string) int { ... }
```

The `Execute()` function and `printUsage()` stay as-is. The `switch` in `Execute()` already calls `runSend`, `runIdle`, `runWaitAndSend` — those functions just move to separate files in the same package.

### `internal/cli/send_cmd.go`

#### Usage

```
whisper send <target> <message...>
whisper send <target> -f <file>
whisper send --self <message...>
whisper send --self -f <file>
```

#### `runSend(args []string) int`

Parse arguments manually (no flag package for positional args — use flag package only for named flags):

```go
func runSend(args []string) int {
    // Parse flags
    var (
        self       bool
        filePath   string
        skipEscape bool
        timeout    time.Duration
    )

    fs := flag.NewFlagSet("send", flag.ContinueOnError)
    fs.BoolVar(&self, "self", false, "send to own pane ($TMUX_PANE)")
    fs.StringVar(&filePath, "f", "", "read message from file")
    fs.BoolVar(&skipEscape, "skip-escape", false, "omit Escape keystroke")
    fs.DurationVar(&timeout, "timeout", 15*time.Second, "max retry timeout")
    fs.SetOutput(os.Stderr)

    if err := fs.Parse(args); err != nil {
        return 2
    }
    remaining := fs.Args()

    // Resolve target
    var target string
    if self {
        target = os.Getenv("TMUX_PANE")
        if target == "" {
            fmt.Fprintln(os.Stderr, "error: --self requires $TMUX_PANE (must be inside tmux)")
            return 2
        }
    } else {
        if len(remaining) < 1 {
            fmt.Fprintln(os.Stderr, "error: target required")
            return 2
        }
        target = remaining[0]
        remaining = remaining[1:]
    }

    // Resolve message
    var message string
    if filePath != "" {
        data, err := os.ReadFile(filePath)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
            return 1
        }
        message = string(data)
    } else {
        if len(remaining) < 1 {
            fmt.Fprintln(os.Stderr, "error: message required (or use -f <file>)")
            return 2
        }
        message = strings.Join(remaining, " ")
    }

    // Execute send
    r := tmux.NewRunner()
    opts := tmux.SendOpts{
        SkipEscape: skipEscape,
        Timeout:    timeout,
    }
    if err := r.Send(target, message, opts); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        return 1
    }
    return 0
}
```

**Key behaviors:**
- `--self` resolves target from `$TMUX_PANE` environment variable
- Without `--self`, first positional arg is the target
- Message is remaining args joined with spaces, OR file contents via `-f`
- `-f` and positional message are mutually exclusive (if `-f` is set, ignore remaining args for message)
- Exit 0 on success, 1 on delivery failure, 2 on usage error
- Nothing on stdout for success. Error on stderr for failure.
- Default timeout is 15s (spec says 15s, Gastown uses 10s for NudgeReadyTimeout — spec wins)

### `internal/cli/idle_cmd.go`

#### Usage

```
whisper idle <target>
whisper idle <target> --timeout 30s
```

#### `runIdle(args []string) int`

```go
func runIdle(args []string) int {
    var timeout time.Duration

    fs := flag.NewFlagSet("idle", flag.ContinueOnError)
    fs.DurationVar(&timeout, "timeout", 15*time.Second, "max wait time")
    fs.SetOutput(os.Stderr)

    if err := fs.Parse(args); err != nil {
        return 2
    }
    remaining := fs.Args()

    if len(remaining) < 1 {
        fmt.Fprintln(os.Stderr, "error: target required")
        return 2
    }
    target := remaining[0]

    r := tmux.NewRunner()
    err := tmux.WaitForIdle(r, target, timeout)
    if err == nil {
        return 0
    }

    // Distinguish exit codes: 1 for timeout, 2 for target errors
    if errors.Is(err, tmux.ErrIdleTimeout) {
        fmt.Fprintln(os.Stderr, "timeout")
        return 1
    }
    if errors.Is(err, tmux.ErrSessionNotFound) || errors.Is(err, tmux.ErrNoServer) {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        return 2
    }
    fmt.Fprintf(os.Stderr, "error: %v\n", err)
    return 1
}
```

**Exit codes (from spec):**
- 0 — target is idle
- 1 — timeout (target still busy)
- 2 — target not found / tmux error

Note: The spec says exit 2 for "target not found / tmux error" but exit 1 for "timeout". This means we need to distinguish `ErrIdleTimeout` from `ErrSessionNotFound`/`ErrNoServer`.

### `internal/cli/wait_send_cmd.go`

#### Usage

```
whisper wait-and-send <target> <message...>
whisper wait-and-send <target> -f <file> --timeout 60s
```

#### `runWaitAndSend(args []string) int`

```go
func runWaitAndSend(args []string) int {
    var (
        filePath   string
        skipEscape bool
        timeout    time.Duration
    )

    fs := flag.NewFlagSet("wait-and-send", flag.ContinueOnError)
    fs.StringVar(&filePath, "f", "", "read message from file")
    fs.BoolVar(&skipEscape, "skip-escape", false, "omit Escape keystroke")
    fs.DurationVar(&timeout, "timeout", 15*time.Second, "max wait time")
    fs.SetOutput(os.Stderr)

    if err := fs.Parse(args); err != nil {
        return 2
    }
    remaining := fs.Args()

    if len(remaining) < 1 {
        fmt.Fprintln(os.Stderr, "error: target required")
        return 2
    }
    target := remaining[0]
    remaining = remaining[1:]

    // Resolve message
    var message string
    if filePath != "" {
        data, err := os.ReadFile(filePath)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
            return 1
        }
        message = string(data)
    } else {
        if len(remaining) < 1 {
            fmt.Fprintln(os.Stderr, "error: message required (or use -f <file>)")
            return 2
        }
        message = strings.Join(remaining, " ")
    }

    // CRITICAL: Acquire lock BEFORE idle check to prevent race.
    // The spec says: "holds the send lock for the entire operation
    // (no race between idle check and send)."
    r := tmux.NewRunner()

    if !tmux.AcquireLock(target, tmux.LockTimeout) {
        fmt.Fprintf(os.Stderr, "error: send lock timeout for %q: previous send may be hung\n", target)
        return 1
    }
    defer tmux.ReleaseLock(target)

    // Wait for idle
    if err := tmux.WaitForIdle(r, target, timeout); err != nil {
        if errors.Is(err, tmux.ErrIdleTimeout) {
            fmt.Fprintln(os.Stderr, "timeout waiting for idle")
            return 1
        }
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        if errors.Is(err, tmux.ErrSessionNotFound) || errors.Is(err, tmux.ErrNoServer) {
            return 2
        }
        return 1
    }

    // Send (Send() will try to acquire the lock again — need to handle this)
    // IMPORTANT: Since we already hold the lock, Send() must not deadlock.
    // Two approaches:
    //   A) Call Send() with lock already held (Send acquires its own lock → deadlock)
    //   B) Call the internal send steps directly, skipping lock acquisition
    //
    // Approach B: Create a SendWithoutLock variant or make Send() accept a "lockHeld" option.
    // The simplest approach: add a LockHeld field to SendOpts.
    opts := tmux.SendOpts{
        SkipEscape: skipEscape,
        Timeout:    timeout,
        LockHeld:   true, // Skip lock acquisition — we already hold it
    }
    if err := r.Send(target, message, opts); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        return 1
    }
    return 0
}
```

**CRITICAL design note for `wait-and-send`:** The spec says it "holds the send lock for the entire operation (no race between idle check and send)." This means:

1. Acquire the lock before calling `WaitForIdle`
2. Keep it held through the `Send`
3. `Send()` normally acquires the lock itself — so we need a way to skip lock acquisition when it's already held

**Solution:** Add a `LockHeld bool` field to `SendOpts` in T3's `send.go`. When true, `Send()` skips the `AcquireLock`/`ReleaseLock` step. This is the cleanest approach — no separate function, just a flag. Update the lock acquisition block in `Send()`:

```go
if !opts.LockHeld {
    if !AcquireLock(target, LockTimeout) {
        return fmt.Errorf("send lock timeout for %q: previous send may be hung", target)
    }
    defer ReleaseLock(target)
}
```

**This means T3 must include `LockHeld` in `SendOpts`** — add it to T3's implementation even though it's motivated by T5. The T3 handoff already shows `SendOpts` — implementer should add `LockHeld bool` there.

---

## Test requirements

### Unit tests: argument parsing

These don't need tmux — they test that argument parsing produces the right values. However, since `runSend`/`runIdle`/`runWaitAndSend` call tmux functions directly, pure unit testing of arg parsing requires either:

1. **Extract arg parsing into a separate function** that returns a struct, test that struct
2. **Test via integration** — the CLI is thin enough that integration tests (T6) cover it

**Recommended: approach 1 for send_cmd.go** (it has the most complex arg parsing):

```go
// parseSendArgs parses send command arguments.
// Returns target, message, opts, or an error.
func parseSendArgs(args []string) (target, message string, opts tmux.SendOpts, exitCode int, err error)
```

Test cases:
- `["mysession", "hello", "world"]` → target=`"mysession"`, message=`"hello world"`
- `["--self", "hello"]` with `TMUX_PANE=%5` → target=`"%5"`, message=`"hello"`
- `["--self", "hello"]` without `TMUX_PANE` → error, exit 2
- `["mysession", "-f", "/tmp/msg.txt"]` → target=`"mysession"`, message from file
- `["mysession"]` → error (no message), exit 2
- `[]` → error (no target), exit 2
- `["--skip-escape", "mysession", "hello"]` → opts.SkipEscape=true
- `["--timeout", "30s", "mysession", "hello"]` → opts.Timeout=30s

### Integration tests (deferred to T6)

T6 will test the full CLI flow: `whisper send %<pane> "hello"` against real tmux.

---

## Done-when criteria

1. `runSend` parses `<target>`, `<message>`, `--self`, `-f`, `--skip-escape`, `--timeout` correctly
2. `runIdle` parses `<target>`, `--timeout` correctly
3. `runWaitAndSend` parses all combined flags and holds lock for entire operation
4. Stub functions removed from `root.go`
5. Exit codes match spec: 0 success, 1 failure/timeout, 2 usage error
6. `--self` reads from `$TMUX_PANE`, errors if not set
7. `-f` reads message from file
8. All commands print errors to stderr, nothing to stdout on success
9. `go build ./cmd/whisper` succeeds
10. `go vet ./...` passes
11. `SendOpts.LockHeld` field exists and `Send()` respects it
