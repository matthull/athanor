# Whisper: Reliable Inter-Session Communication CLI

**Status:** Draft
**Created:** 2026-03-21
**Location:** `specs/athanor/whisper/`
**Language:** Go
**Replaces:** `~/.local/bin/tmux-claude-send` (49 lines of bash)

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[G:tmux.go]` | From Gastown `internal/tmux/tmux.go` (battle-tested, 2500+ lines) |
| `[G:protocols]` | From `specs/athanor/gastown-tmux-protocols.md` (our analysis) |
| `[D:reason]` | Design decision |
| `[E:existing]` | Extends existing code/pattern |
| `[B:context]` | Business/architectural requirement |

---

## Overview

Whisper is a Go CLI that reliably delivers messages to Claude Code sessions running in tmux. It replaces our current `tmux-claude-send` bash script, which has known bugs that are tolerable for manual use but compound catastrophically under automation.

**Why this matters:** Every automated interaction in the athanor is a whisper. The marut nudging azers, dispatch spawning crucibles, handoff chaining sessions — all depend on reliable message delivery to tmux panes. If whispers are flaky, everything built on top is flaky. `[B:domain-model.md § Infrastructure]`

**Core thesis:** Take Gastown's battle-tested 8-step NudgeSession protocol, implement it as a standalone CLI, and make it the single entry point for all message delivery to Claude Code sessions. Not 39K lines — but built on the same protocol that runs 20-30 concurrent agents reliably.

---

## Domain Dictionary

| Term | What it is |
|------|-----------|
| **Whisper** | The CLI tool itself. Also the act of delivering a message. "Whisper to pane %5" = deliver a message reliably. |
| **Target** | A tmux pane identifier. Can be a session name (resolves to the active pane), a pane ID (`%5`), or a `session:window.pane` specifier. |
| **Nudge** | Gastown's term for the same concept. We use "whisper" externally but the protocol is the NudgeSession protocol. |
| **Idle** | A Claude Code session showing its input prompt (`>` prefix), not actively running tools. Confirmed by 2-check protocol. |
| **Send lock** | Per-target mutex preventing concurrent sends from garbling input. |

---

## What's Wrong with tmux-claude-send

The current script (`~/.local/bin/tmux-claude-send`):

```bash
tmux send-keys -t "$PANE" "$PROMPT"   # no -l flag!
sleep 0.3
tmux send-keys -t "$PANE" Enter
sleep 0.5
tmux send-keys -t "$PANE" Enter
```

**Known bugs** `[G:protocols § Gap Analysis]`:

| Bug | Impact | Whisper fix |
|-----|--------|-------------|
| No `-l` flag on send-keys | `!`, `~`, `#` interpreted as tmux key names | Always use `send-keys -l` (literal mode) |
| No copy mode handling | Scrolled panes silently eat all input | Detect and exit copy mode before sending |
| No ESC→Enter gap | Readline interprets ESC+Enter within 500ms as M-Enter (meta-return) which does nothing | 600ms gap between Escape and Enter |
| No concurrent send protection | Multiple senders garble input | Per-target send locks with 30s timeout |
| No idle detection | Always interrupts, no awareness of agent state | 2-check idle detection protocol |
| No message sanitization | Control chars corrupt delivery | Strip ESC, CR, BS, DEL; TAB→space |
| No chunking for large messages | Messages >512 bytes can be truncated by TTY buffer | Chunk with 10ms inter-chunk delays |
| No retry on transient errors | "not in a mode" errors during TUI startup cause permanent failure | Exponential backoff retry |
| No detached session wake | Detached sessions never process input | SIGWINCH via resize-window trick |

---

## CLI Interface

### `whisper send`

Deliver a message to a tmux target reliably.

```
whisper send <target> <message>
whisper send <target> -f <file>
whisper send --self <message>
whisper send --self -f <file>
```

**Arguments:**
- `<target>` — tmux pane/session identifier (e.g., `%5`, `mysession`, `mysession:0.1`)
- `<message>` — text to deliver (remaining args joined with spaces)
- `-f <file>` — read message from file instead of args
- `--self` — send to own pane (uses `$TMUX_PANE`)

**Options:**
- `--skip-escape` — omit Escape keystroke + 600ms delay (for agents where Escape cancels generation, e.g. Gemini CLI) `[G:tmux.go NudgeOpts.SkipEscape]`
- `--timeout <duration>` — max time to wait for send-keys retry (default: 15s) `[G:tmux.go constants.NudgeReadyTimeout]`

**Exit codes:**
- `0` — message delivered successfully
- `1` — delivery failed (target not found, lock timeout, retry exhausted)
- `2` — usage error (bad args)

**Stdout:** Nothing on success. Error description on failure.

`[D:simplicity]` No `--queue` flag in v1. Queue-based delivery (file inbox + hook drain) is a separate concern — it uses whisper send as its delivery mechanism but the queue logic lives elsewhere (likely a hook or the marut). Whisper is the reliable transport layer, not the queuing system.

### `whisper idle`

Wait for a target to become idle (showing prompt, not running tools).

```
whisper idle <target>
whisper idle <target> --timeout 30s
```

**Arguments:**
- `<target>` — tmux pane/session identifier

**Options:**
- `--timeout <duration>` — max time to wait (default: 15s)

**Exit codes:**
- `0` — target is idle
- `1` — timeout (target still busy)
- `2` — target not found / tmux error

**Stdout:** Nothing on success. "timeout" or error on failure.

### `whisper wait-and-send`

Convenience: wait for idle, then send. The common pattern for non-interruptive delivery.

```
whisper wait-and-send <target> <message>
whisper wait-and-send <target> -f <file> --timeout 60s
```

Equivalent to `whisper idle <target> && whisper send <target> <message>` but holds the send lock for the entire operation (no race between idle check and send).

**Options:** combines options from `send` and `idle`.

**Exit codes:**
- `0` — waited for idle, then delivered
- `1` — idle timeout OR delivery failed
- `2` — usage error

---

## The 8-Step Send Protocol

This is the core of whisper. Adapted from Gastown's `NudgeSessionWithOpts` (`internal/tmux/tmux.go:1421`). Every `whisper send` executes these steps in order.

### Step 1: Acquire send lock `[G:tmux.go:1424]`

Per-target channel-based semaphore. Prevents concurrent sends to the same target from interleaving text.

- Timeout: 30 seconds `[G:tmux.go nudgeLockTimeout]`
- If lock not acquired: return error "previous send may be hung"
- Lock released on function exit (defer)

**Why channel-based, not mutex:** Supports timed acquisition. A hung tmux command holding a mutex would block all future sends permanently.

### Step 2: Exit copy mode `[G:tmux.go:1438]`

If the target pane is in copy/scroll mode, input goes to the copy buffer instead of the process.

```
tmux display-message -p -t <target> "#{pane_in_mode}"
```

If result is `"1"`: send `tmux send-keys -t <target> -X cancel`, wait 50ms.

### Step 3: Sanitize message `[G:tmux.go:1285]`

Strip control characters that corrupt tmux delivery:

| Character | Action | Why |
|-----------|--------|-----|
| TAB (0x09) | Replace with space | Triggers shell tab completion |
| Newline (0x0a) | Preserve | send-keys -l handles newlines |
| All other chars < 0x20 | Strip | ESC triggers terminal sequences, CR acts as premature Enter, BS deletes chars |
| DEL (0x7f) | Strip | Deletes characters |
| Everything else | Preserve | Including quotes, backticks, Unicode |

### Step 4: Send text with chunking + retry `[G:tmux.go:1316]`

Use `tmux send-keys -l -t <target> <text>` (literal mode — handles special chars correctly).

**Chunking** (messages > 512 bytes):
- Split into 512-byte chunks
- First chunk: use retry logic (see below)
- Subsequent chunks: `send-keys -l`, then 10ms delay between chunks
- Rationale: Linux canonical mode TTY buffer is 4096 bytes. Chunking prevents kernel line discipline truncation. `[G:protocols § Known Workarounds]`

**Retry logic** (first chunk or small messages):
- Exponential backoff: 500ms → 750ms → 1125ms → 1687ms → 2000ms cap
- Retry on transient errors: `"not in a mode"` (TUI not ready during cold startup) `[G:tmux.go:1272]`
- Non-transient errors: fail immediately
- Max retries governed by timeout flag

### Step 5: Post-send delay — 500ms `[G:tmux.go:1453]`

Wait 500ms after send-keys completes. **Empirically required** for text to arrive at the process. Without this, text delivery is unreliable.

### Step 6: Send Escape `[G:tmux.go:1458]`

`tmux send-keys -t <target> Escape`

Exits vim INSERT mode if active (harmless in normal mode). Ensures clean state before Enter.

**Skipped if `--skip-escape` is set** (for agents where Escape cancels in-flight generation).

### Step 7: ESC→Enter gap — 600ms `[G:tmux.go:1460]`

**CRITICAL.** Must exceed bash readline's `keyseq-timeout` (500ms default).

Without this delay: ESC + Enter within 500ms → readline interprets as M-Enter (meta-return) → **does nothing**. This is the root cause of "sometimes Enter doesn't work." `[G:protocols § Known Workarounds]`

**Skipped if `--skip-escape` is set** (no Escape sent, so no timing issue).

### Step 8: Send Enter with retry `[G:tmux.go:1469]`

Three attempts, 200ms between each:

```
attempt 1: tmux send-keys -t <target> Enter
  if error: wait 200ms, retry
attempt 2: tmux send-keys -t <target> Enter
  if error: wait 200ms, retry
attempt 3: tmux send-keys -t <target> Enter
  if error: fail
```

On success: call WakePaneIfDetached (step 8b).

### Step 8b: Wake detached sessions `[G:tmux.go:1262]`

Detached tmux sessions don't process stdin until a terminal event triggers SIGWINCH.

If the session is NOT attached:
1. `tmux resize-window` (triggers SIGWINCH)
2. Wait 50ms
3. Restore `window-size latest` (tmux 3.3+ stickiness bug) `[G:protocols § Known Workarounds]`

---

## 2-Check Idle Detection Protocol

Used by `whisper idle` and `whisper wait-and-send`. Adapted from Gastown's `WaitForIdle` (`internal/tmux/tmux.go:2470`).

### The Problem

Claude Code briefly shows its prompt between tool calls (~500ms gap). A single idle check produces false positives — the agent appears idle but is actually mid-execution.

### The Protocol

Poll in a loop until timeout:

1. **Capture last 5 lines** of target pane via `tmux capture-pane`
2. **Check status bar** — if lines contain `"⏵⏵"` AND `"esc to interrupt"`, agent is actively running a tool. NOT idle. Reset consecutive count, sleep 200ms, continue.
3. **Check for prompt prefix** — scan captured lines for `>` (U+276F) at line start. (Claude Code's idle prompt.)
4. **If prompt found:** increment consecutive idle count
5. **If consecutive count reaches 2:** return idle (success)
6. **If prompt NOT found:** reset consecutive count
7. Sleep 200ms, repeat from step 1

**Why 2 checks:** Two polls 200ms apart create a 400ms observation window. This catches the transient prompt display between tool calls, which lasts ~500ms. Two consecutive idle readings 200ms apart confirms genuine idle state. `[G:tmux.go:2474]`

**Error handling:**
- Session not found / no tmux server → return error immediately (no retry)
- Transient capture errors → reset count, continue polling

---

## Project Structure

```
whisper/
  cmd/
    whisper/
      main.go           # CLI entry point, arg parsing
  internal/
    tmux/
      send.go           # 8-step send protocol
      idle.go           # 2-check idle detection
      lock.go           # Per-target send locks
      sanitize.go       # Message sanitization
      capture.go        # Pane capture + prompt detection
      wake.go           # Detached session wake
      tmux.go           # Low-level tmux command execution
    cli/
      send_cmd.go       # `whisper send` command
      idle_cmd.go       # `whisper idle` command
      wait_send_cmd.go  # `whisper wait-and-send` command
  go.mod
  go.sum
  Makefile
```

`[D:standalone]` Whisper is a standalone Go module, not part of a larger monorepo. Single binary, installable to `~/.local/bin/whisper`.

`[D:no-gastown-deps]` No dependency on Gastown packages. We study the protocol from their code but implement independently. The protocol is simple enough that a clean implementation is better than importing their internal packages.

---

## Reference Code

The Gastown implementation at `~/code/gastown/internal/tmux/tmux.go` is the primary reference. Key functions to study:

| Function | Line | What to learn |
|----------|------|---------------|
| `NudgeSessionWithOpts` | 1421 | The 8-step protocol, end to end |
| `sanitizeNudgeMessage` | 1285 | Character sanitization rules |
| `sendMessageToTarget` | 1316 | Chunking + retry logic |
| `sendKeysLiteralWithRetry` | (search) | Exponential backoff pattern |
| `isTransientSendKeysError` | 1272 | Which errors are retryable |
| `WaitForIdle` | 2470 | 2-check idle detection |
| `WakePaneIfDetached` | 1262 | SIGWINCH wake trick |
| `acquireNudgeLock` / `releaseNudgeLock` | 1191 / 1202 | Channel-based semaphore |
| `FindAgentPane` | 1736 | Multi-pane targeting (future, not v1) |

**How to use this reference:** Read the function, understand the protocol step, implement your own clean version. Gastown's code is embedded in a large framework with config, telemetry, and multi-rig concerns — extract the protocol logic, leave the framework.

---

## Test Strategy

### Unit tests (pure logic, no tmux)

- **Sanitization:** Known control chars stripped, printable preserved, Unicode preserved, TAB→space, newline preserved
- **Chunking:** Messages split at 512-byte boundaries, chunk count correct for various sizes
- **Argument parsing:** All CLI forms parse correctly, error on bad args
- **Lock timeout:** Simulated lock contention times out after 30s

### Integration tests (require tmux)

These tests start a real tmux session, send messages, and verify delivery.

- **Basic send:** Send a known string, capture pane content, verify it arrived
- **Literal mode:** Send special chars (`!`, `~`, `#`, backtick), verify they arrive literally
- **Large message:** Send >512 byte message, verify full delivery (no truncation)
- **Copy mode:** Put pane in copy mode, send message, verify copy mode exited and message delivered
- **Idle detection:** Start a `sleep` command (busy), verify `whisper idle` waits. Let it finish, verify idle detected.
- **Concurrent sends:** Two goroutines send to same target simultaneously, verify no garbling (messages arrive sequentially, not interleaved)
- **Detached session:** Create detached session, send message, verify delivery

`[D:test-isolation]` Integration tests create ephemeral tmux sessions with unique names, clean up after themselves. No interaction with user's running sessions.

### What NOT to test in v1

- Multi-pane targeting (FindAgentPane) — not in v1 scope
- Queue/inbox mechanism — separate concern
- Startup dialog dismissal — separate concern

---

## Scope

### In scope (v1)

- `whisper send` with full 8-step protocol
- `whisper idle` with 2-check detection
- `whisper wait-and-send` convenience command
- Per-target send locks (in-process; sufficient for single whisper binary)
- `--self` flag for self-targeting
- `-f` flag for file-based input
- `--skip-escape` flag
- `--timeout` flag
- Install to `~/.local/bin/whisper`

### Out of scope (later)

- **File-based message queue** (inbox + drain hook) — separate tool/infrastructure, uses whisper as transport
- **Multi-pane targeting** (`FindAgentPane`) — when we have multi-pane sessions
- **Startup dialog dismissal** (trust dialog, bypass permissions) — separate tool or whisper subcommand later
- **Cross-process send locks** — v1 locks are in-process (channel semaphore). Cross-process locks (flock) needed when multiple whisper processes run concurrently. Add when that's a real scenario.
- **Session lifecycle** (create, kill, health check) — different concern, different tool
- **Configurable prompt prefix** — hardcoded to Claude Code's `>` for now

---

## Acceptance Criteria

1. `whisper send %<pane> "hello world"` delivers the message and returns exit 0
2. `whisper send` with special chars (`!~#` + backticks + Unicode) delivers them literally
3. `whisper send` with a 2KB message delivers the full message without truncation
4. `whisper send` to a pane in copy mode exits copy mode and delivers
5. `whisper send` to a detached session delivers (SIGWINCH wake)
6. `whisper idle %<pane>` returns 0 when pane is idle, 1 when busy
7. `whisper idle` does not false-positive during inter-tool-call prompt flash
8. Two concurrent `whisper send` to the same target produce non-garbled output
9. `whisper send --self "test"` works from within a tmux session
10. `tmux-claude-send` can be replaced with `whisper send` as a drop-in (same args pattern for basic usage)
11. Integration tests pass in CI (tmux available in test environment)

---

## Migration

`tmux-claude-send` is used in:
- `~/.local/bin/tmux-claude-send` (the script itself)
- Telegram MCP's reply delivery (`tmux-claude-send` called by egregore)
- Ad-hoc operator use

**Migration path:**
1. Build whisper, install to `~/.local/bin/whisper`
2. Create `~/.local/bin/tmux-claude-send` wrapper that calls `whisper send` (backward compat)
3. Update Telegram MCP / egregore to call `whisper send` directly
4. Remove wrapper after all callers updated

---

## Retrospective

When whisper is complete, review:

### Documentation Updates
- [ ] Update `specs/athanor/open-work-items.md` — mark #15 complete, update name to whisper
- [ ] Update `specs/athanor/domain-model.md` — mark whisper as built
- [ ] Update `specs/athanor/gastown-tmux-protocols.md` — note gap closure

### Workflow Improvements
- [ ] Create Go CLI development skill if patterns emerge worth capturing
- [ ] Assess whether Gastown's Go project structure is a good template for future tools

### Knowledge Capture
- [ ] Document any tmux quirks discovered during implementation
- [ ] Note any protocol deviations from Gastown (and why)
