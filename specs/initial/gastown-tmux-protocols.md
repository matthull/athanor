# Gastown Tmux Protocols: Battle-Tested Patterns

**Date:** 2026-03-15
**Source:** ~/code/gastown/internal/tmux/, internal/nudge/, internal/cmd/nudge.go
**Context:** How Gastown reliably manages 20-30 concurrent Claude Code sessions in tmux

---

## Three Delivery Modes

Gastown doesn't just `send-keys`. It has three modes for delivering messages to agent sessions:

### 1. wait-idle (default)
- Poll pane for idle prompt (up to 15 seconds)
- If idle: deliver immediately (low latency, no interruption)
- If busy: fall back to queue
- If queue fails: fall back to immediate (last resort)
- Background watcher continues polling 60s after queuing

### 2. queue
- Write message to file in `<town>/.runtime/nudge_queue/<session>/`
- Agent picks up via `UserPromptSubmit` hook at next turn boundary
- Zero interruption to work
- File-based: atomic rename prevents double-delivery
- Per-nudge TTL (30 min normal, 2 hours urgent)
- Expired nudges silently discarded on drain

### 3. immediate
- Raw `tmux send-keys` — interrupts in-flight work
- Used for emergencies only
- Guaranteed delivery but potentially destructive

---

## The NudgeSession Protocol (8 Steps)

This is the actual reliable delivery sequence when sending a message to a tmux session:

```
1. EXIT COPY MODE
   - Send 'q' to exit scroll/copy mode if active
   - Without this, send-keys goes to copy buffer, not input

2. SANITIZE
   - Strip control characters (ESC, CR, BS, TAB → space)
   - Prevents terminal escape sequence injection

3. SEND TEXT (with retry)
   - Use `send-keys -l` (literal mode — handles special chars)
   - Messages >512 bytes: chunk with 10ms inter-chunk delays
     (Linux canonical mode buffer is 4096 bytes, chunking prevents truncation)
   - Retry with exponential backoff: 500ms → 750ms → 1125ms → 1687ms → 2s cap
   - Transient error "not in a mode" triggers retry (agent TUI not ready)

4. WAIT 500ms
   - EMPIRICALLY REQUIRED for text to arrive at the process
   - Without this, text delivery is unreliable

5. SEND ESCAPE
   - Exits vim INSERT mode (harmless in normal mode)
   - Ensures clean state before Enter

6. WAIT 600ms (CRITICAL)
   - MUST exceed bash readline's 500ms keyseq-timeout
   - Without this delay: ESC+Enter within 500ms → readline interprets as
     M-Enter (meta-return) which DOES NOTHING
   - This is why "sometimes Enter doesn't work" — it's a readline timing bug

7. SEND ENTER (3 retries, 200ms apart)
   - Three separate `send-keys Enter` calls
   - 200ms between each
   - Handles: agent TUI not ready, terminal processing lag, timing races

8. SIGWINCH WAKE (if detached)
   - Detached sessions: Claude Code TUI doesn't process stdin until terminal event
   - Solution: resize-window (triggers SIGWINCH), wait 50ms, restore size
   - Only for detached sessions
```

---

## Idle Detection: 2-Check Protocol

```
Goal: Determine if agent is idle (showing prompt, not working)

Method:
1. Capture last 5 lines of pane via `tmux capture-pane`
2. Look for prompt prefix `❯ ` (U+276F) at line start
3. REJECT if status bar shows "⏵⏵" + "esc to interrupt" (tool in progress)
4. If idle detected: WAIT 200ms, check again
5. Only confirm idle if BOTH checks pass

Why 2 checks:
- Claude Code briefly shows prompt between tool calls (~500ms gap)
- Single check produces false positives
- 2 checks 200ms apart (400ms window) catches transient prompt display
```

For non-Claude agents: check for `ReadyPromptPrefix` from config. If no prompt detection available, degrade to queue mode.

---

## Concurrency Control

### Per-Session Send Locks
- Channel-based semaphores prevent concurrent sends to same session
- 30-second timeout prevents permanent blocking
- Without locks: concurrent nudges interleave text, producing garbled input

### Queue Drain Atomics
- Claim files via rename to `.claimed.<random-suffix>`
- Orphaned `.claimed` files (>5 min old) automatically reclaimed
- FIFO ordering via nanosecond timestamp + random suffix
- Prevents double-delivery in race conditions

---

## Session Management

### Naming Convention
- `gt-<rig>-<role>` (e.g., `gt-gastown-alpha`, `hq-mayor`)
- Socket isolation: each town on separate tmux socket (`-L <socket>`)
- Prevents session name collisions across multiple towns

### Health Checks
- `IsAgentAlive()` — check if session exists
- `IsAgentRunning()` — check if agent process is running in pane
- Zombie cleanup: `CleanupOrphanedSessions()` at startup

### Split-Brain Prevention
- Kill same-named session on default socket before creating on town socket
- Prevents stale sessions shadowing real ones

### Process Cleanup
- Walk process tree (deepest-first), SIGTERM, wait 2s, SIGKILL
- Handle reparented processes (PPID==1) via process group verification
- Prevents orphans when SIGHUP ignored

---

## Session Cycling (C-b n/p)

```bash
# Dynamic keybinding setup
bind-key n if-shell \
  "echo '#{session_name}' | grep -Eq '^(gt|hq|bd|...)-'" \
  "run-shell 'gt cycle next --session #{session_name} --client #{client_tty}'" \
  "next-window"  # fallback for non-GT sessions
```

- Pattern generated dynamically from `rigs.json`
- `--client #{client_tty}` targets correct client when multiple attached
- Stale pattern detection: rebinds when rig prefixes change

---

## Startup Dialog Handling

### Workspace Trust Dialog
- Poll pane content (30 lines) every 200ms for 5 seconds
- Look for "trust this folder" or "Quick safety check"
- Send Enter (option 1 pre-selected), wait 500ms
- Exit early if prompt visible (no dialog)

### Bypass Permissions Warning
- Look for "Bypass Permissions mode"
- Send Down (select "Yes"), wait 200ms, send Enter
- Same early-exit logic

### Blind Dialog Dismissal
- For stalled sessions: Enter → Down+Enter (no content scraping)
- Safe on sessions without dialogs (blank input to idle prompt)

---

## Known Workarounds & Hacks

### TTY Input Buffer Limits
- Linux canonical mode buffer: 4096 bytes
- Messages >512 bytes chunked with 10ms delays
- Prevents kernel line discipline truncation

### Window Size Stickiness (tmux 3.3+)
- Auto-resize broken after `resize-window`
- Must restore `window-size latest` after resize dance (50ms gap)

### Pane-Died Hook Complications
- `remain-on-exit on` prevents auto-destruction but complicates cleanup
- Must disable before `kill-session` to prevent infinite respawn loops
- `respawn-pane` resets `remain-on-exit` → must re-enable in hook

### Multi-Pane Targeting
- ZFC approach: check `GT_PANE_ID` env var set at startup
- Fallback: scan all panes for agent process via `list-panes -s`
- Avoids `send-keys -t <session>` targeting wrong pane

---

## Gap Analysis: Our `tmux-claude-send` vs Gastown

Our script (`~/.local/bin/tmux-claude-send`, 49 lines):
```bash
tmux send-keys -t "$PANE" "$PROMPT"   # send text (no -l flag!)
sleep 0.3                              # autocomplete settle
tmux send-keys -t "$PANE" Enter        # first Enter
sleep 0.5                              # input buffer
tmux send-keys -t "$PANE" Enter        # second Enter
```

### Side-by-Side

| Step | Gastown (8-step) | tmux-claude-send | Gap |
|---|---|---|---|
| Exit copy mode | Detects + sends `-X cancel` | None | Copy mode silently eats input |
| Sanitize message | Strips control chars | None | Control chars corrupt delivery |
| Send text | `send-keys -l` (literal) + chunking + retry | `send-keys` (no `-l`!) | Special chars misinterpreted, large messages fail |
| Post-send delay | 500ms | 300ms | Potentially too short |
| Send Escape | Conditional (exit vim INSERT) | None | Nudges land in vim INSERT unsent |
| ESC→Enter gap | 600ms (readline keyseq-timeout) | N/A (no Escape) | N/A but relevant if we add Escape |
| Send Enter | 3 retries, 200ms apart | 2 Enters, 500ms apart | Our 2-Enter approach is a workaround, not principled |
| Wake detached | SIGWINCH via resize | None | Detached sessions never process input |
| Concurrency lock | Per-session semaphore | None | Concurrent sends produce garbled input |
| Pane routing | Auto-find agent pane | Manual pane ID required | Multi-pane sessions route wrong |
| Idle detection | 2-check poll with prompt detection | None | Always interrupts, no awareness of agent state |

### Key Issues in Our Script

1. **No `-l` flag on send-keys** — Without literal mode, characters like `!`, `~`, `#` get interpreted as tmux key names. This is a bug.
2. **No copy mode handling** — If the user scrolled up in the pane, all input is silently eaten.
3. **Double Enter as workaround** — We send 2 Enters because we noticed "sometimes one isn't enough." Gastown's analysis shows WHY: it's not about Enter count, it's about the ESC→Enter timing gap with readline. Our 500ms gap between Enters is coincidentally close to readline's 500ms keyseq-timeout, which may be why 2 Enters usually works for us.
4. **No concurrent send protection** — `/permission-manager` and any future supervisor sending to multiple sessions can garble input.

---

## Session Message Queue: Per-Session Inbox

### How Gastown Does It

Each session has a file-based message queue at `<town>/.runtime/nudge_queue/<session>/`. Messages are JSON files named with nanosecond timestamps for FIFO ordering.

**Enqueue**: Write JSON file to the session's queue directory.
```json
{
  "sender": "gastown/crew/max",
  "message": "Your turn to work",
  "priority": "normal",
  "timestamp": "2025-03-15T10:30:45Z",
  "expires_at": "2025-03-15T11:00:45Z"
}
```

**Drain**: Atomic rename-based claiming (`file.json` → `file.claimed.{random}`). Only one drain process wins the rename. Orphaned claims (>5 min) automatically reclaimed. Max queue depth: 50.

**Delivery trigger**: `gt mail check --inject` runs in `UserPromptSubmit` hook. Drains BOTH mail and nudge queues in one call. Outputs as `<system-reminder>` block.

### The Autonomous Session Gap

**`UserPromptSubmit` does NOT fire during autonomous work** (tool call → tool call with no human input). Queue sits undelivered.

**Gastown's solutions:**
1. **`watchAndDeliver()`** — synchronous 60s idle-poll after queuing. Drains when agent becomes idle.
2. **`nudge-poller`** — background process polling every 10s for idle, then drains. For non-Claude agents.
3. **Natural drain** — agent finishes work, goes idle, next human prompt triggers hook.

**Design philosophy: cooperative, not disruptive.** Messages arrive at turn boundaries, not mid-tool-call.

### What To Build for Our System

A per-session message queue for cross-agent and supervisor communication:
1. **File-based queue** at `/tmp/claude-sessions/{session_id}/inbox/`
2. **UserPromptSubmit hook** drains it (already have session injection infrastructure)
3. **Idle-poller** for autonomous sessions (small background script, checks every 10s)
4. **Atomic claiming** via rename (crash-safe, no double delivery)

Use cases: supervisor escalation, cross-agent "I'm done" signals, operator instructions to running agents.

---

## What To Steal (Prioritized for Our System)

### Immediate (fix known bugs)
1. **600ms ESC→Enter gap** — we've probably been hitting the M-Enter readline bug silently. Any tmux-sidecar interaction that sends Enter after text should include this delay.
2. **500ms post-send-keys delay** — text doesn't reliably arrive without it.
3. **Literal mode (`send-keys -l`)** — handles special characters correctly.

### Short-Term (improve reliability)
4. **2-check idle detection** — poll for prompt, confirm with second check 200ms later. Prevents false idle positives.
5. **Per-session send locks** — prevents garbled input from concurrent sends (relevant for `/permission-manager` and any supervisor that sends to multiple sessions).
6. **3-retry Enter** — 200ms between each. Handles timing races.

### Medium-Term (enable autonomy)
7. **File-based message queue** — for non-interruptive delivery. Write to file, agent picks up via UserPromptSubmit hook. Crash-safe, no in-memory state.
8. **Wait-idle hybrid mode** — try immediate if idle, fall back to queue if busy. Best of both worlds.
9. **Startup dialog handling** — auto-dismiss trust and permissions dialogs. Essential for autonomous session spawning.

### Long-Term (scale)
10. **Socket isolation** — separate tmux sockets per project/concern.
11. **Dynamic session cycling** — C-b n/p with auto-generated patterns.
12. **Split-brain prevention** — check for stale sessions before creating new ones.

---

## Session Boundaries: Gastown's Model vs Ours

### How Gastown Session Boundaries Work

From Gastown's polecat-lifecycle-patrol docs:

> "A single step may span multiple session cycles (if the step is complex or compaction occurs). Multiple steps may fit in a single session (if steps are small and context permits). The session-per-step model is a design target, not a hard constraint."

**Sessions end when:**
1. **Step completion** → agent runs `bd close <step>` then `gt handoff` (voluntary cycle)
2. **Context fills** → PreCompact hook triggers automatic handoff
3. **Crash** → Witness detects dead session, respawns
4. **Done** → `gt done` on final step (permanent kill)

**The sandbox (worktree + branch) survives all session cycles.** Only `gt done` destroys the sandbox. Sessions are cattle; the worktree is the persistent identity.

### Two Molecule Archetypes

| Type | Steps Tracked | Use Case |
|---|---|---|
| **Root-only** (`pour=false`, default) | Steps shown inline at prime, no DB tracking | High-frequency work (polecat tasks, patrols) |
| **Poured** (`pour=true`) | Steps materialized as sub-beads, explicit close + checkpoint recovery | Expensive work (releases, complex orchestrations) |

### Mapping to Our System

| Gastown | Our Equivalent | Status |
|---|---|---|
| Sandbox (worktree + branch) | Worktree via `wtp add` | Have it |
| Molecule (step tracking) | task-plan.md with `[x]`/`[ ]` checklist | Have it (convention, not infrastructure) |
| `gt handoff` (voluntary session cycle) | Operator-initiated session restart + reload `/orchestrate` | Have it, manual |
| `gt prime` (auto-inject context at session start) | Session injection via `session-mode-inject.sh` | Have it, partial |
| GUPP nudge (auto-start work after session spawn) | Nothing — operator must type `/orchestrate` | Gap |
| PreCompact auto-handoff | Nothing — context fills, quality degrades | Gap |
| Witness respawn on crash | Nothing — operator notices and restarts | Gap |

### The GUPP Gap (Key Difference)

Our system requires the operator to:
1. Start a new session
2. Load `/orchestrate` via Skill tool
3. It reads the plan, skips completed phases, continues

Gastown automates all three via:
1. `gt handoff` spawns fresh session automatically
2. `gt prime --hook` (SessionStart hook) injects full context
3. GUPP nudge (30-60s after start) triggers `UserPromptSubmit` → mail/queue drain → agent reads hook and continues

**Closing the gap doesn't require full Gastown infrastructure.** Steps:
1. `send-claude-message` tool (#15) enables automated nudging
2. SessionStart hook already injects context — extend it to detect incomplete task plans
3. A handoff script: write state → spawn fresh session → send nudge with "continue from plan"

### Practical Implication

We're on the right track with plan-based session persistence. The next step isn't architectural change — it's automating the session handoff. Prerequisites:
1. `send-claude-message` IPC tool (task #15)
2. SessionStart hook enhancement to detect and resume incomplete plans
3. A `handoff` command that writes state + spawns + nudges

---

## Gastown Shiny Workflow: Step Scope Reference

The `shiny` formula (base engineering pipeline) has 5 steps. Each step is a logical unit of work — agents do multiple steps per session if context permits:

```
1. design    — Read issue, explore codebase, write design doc with approach/trade-offs/files
2. implement — Execute design, write code, make incremental commits
3. review    — Self-review the diff (correctness, style, edge cases)
4. test      — Run test suite, fix failures, verify no regressions
5. submit    — Push branch, submit MR to merge queue
```

`shiny-enterprise` extends this with Rule of Five (5-pass implement) and security audit aspect.

The polecat-work formula has 9 more granular steps including branch setup, build checks, pre-verification, await-verdict, and self-clean.

Step scope varies: "implement" could be hours of work; "branch-setup" is 5 minutes. The agent walks through as many steps as context permits before handing off.
