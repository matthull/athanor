# Athanor Technical Architecture

*How the system physically runs — processes, sessions, lifecycle, infrastructure.*

**Status:** Work in progress (2026-03-30)

This document describes the runtime architecture of the athanor system. The conceptual system design lives in `spec.md`. This document covers: how agents become processes, how those processes are managed, how the CLI orchestrates them, and the infrastructure that keeps the system running.

---

## Runtime Model

The athanor runs as a constellation of Claude Code processes inside tmux sessions. There is no daemon, no server, no central process. The `ath` CLI is the orchestration surface — it creates tmux windows, launches Claude sessions into them, and provides observation and communication tools.

### Process Lifecycle by Agent Type

| Agent | Launch | Session Mode | Permission Mode | Lifecycle | Termination |
|-------|--------|-------------|-----------------|-----------|-------------|
| **Marut** | `ath kindle` | Interactive | Default | Long-lived, restartable via `ath reforge` | Manual (`kill-window`) or context exhaustion → reforge |
| **Azer** | `ath muster` | Interactive | Default | Medium-lived, bound to opus | Self-discharge → marut cleanup |
| **Homunculus** | Manual | Interactive | Default | Persistent, artifex's workbench | Manual |
| **Craft session** | `ath craft` | Interactive | Default | Artifex-directed | Manual discharge |

---

## The `ath` CLI as Session Launcher

`ath` is the de facto Claude session launcher. All agent sessions are created through it (except the homunculus, which is manual). The launch mechanism is tmux keystroke injection:

1. `ath` creates a tmux window (`new-window -n <name>`)
2. Sends `cd <project-dir>` via `SendKeysLiteral`
3. Sends the full `claude` command with flags and boot prompt
4. Sends `Enter`

This means `ath` doesn't exec Claude directly — it operates through tmux. The Claude process is a child of the shell in the tmux pane, not of `ath`. This is deliberate: `ath` can exit after launch, and the agent session persists independently.

### Environment Variables

- `ATHANOR=<instance-dir>` — set in the shell before launch, tells the agent which athanor instance it belongs to

---

## Tmux Architecture

### Session Layout

Sessions are organized by project/concern (not by agent type):

```
athanor-dev      # athanor system development
├── homunculus   # artifex's workbench
├── artifex      # direct artifex work
├── azer-*       # active azer windows
└── marut-*      # marut supervisor windows

my-app           # product development
├── artifex
├── marut-my-app-*
└── azer-*
```

### Naming Conventions

| Entity | Window Name Pattern |
|--------|-------------------|
| Marut | `marut-<athanor>` or `marut-<athanor>-<mo>` |
| Azer | `azer-<opus-name>` |
| Homunculus | `homunculus` |

### No Socket Isolation (Current)

All sessions share the default tmux socket. Known gap vs. the Gastown predecessor, which used per-town socket isolation (`-L <socket>`).

---

## `ath patrol` — System Scanner

Scans all tmux panes for conditions that need attention. Currently detects:

| Finding Type | Detection | Source |
|-------------|-----------|--------|
| `permission` | Co-occurrence of menu cursor (`❯ N.`) and question text in bottom 10 lines | All panes |
| `exhausted` | "Context limit reached" in bottom 5 lines | Claude-process panes only |

Extensible — new `FindingType` values and detection functions drop in without changing the scan loop.

**Output modes:**
- Human-readable table (default)
- JSON lines (`--json`) for machine consumption
- `--exclude <pane>` to skip a specific pane (prevents self-scanning)

---

## Inter-Agent Communication

### Whisper Protocol

`ath whisper send <target> <message>` — reliable message delivery to a tmux target. The full protocol (`internal/tmux/send.go`):

1. Acquire lock (prevents concurrent sends to same target)
2. Exit copy mode if active
3. Sanitize message (escape special characters)
4. Chunk and send with retry
5. Release lock

Used by azers to notify maruts of discharge, by maruts to direct azers, and by the artifex to communicate with any agent.

### Push Notification Escalation

Via Claude's built-in `PushNotification` tool (desktop + mobile push):
- Standard notifications — concise message under 200 chars
- `ANDON:` prefix — urgent intervention request

---

## Observation Infrastructure

### `ath status` / `ath dashboard`

Reads athanor state from:
- Tmux window presence (is the marut/azer alive?)
- `CheckCrucible()` health detection (active/idle/permission/exhausted/dead)
- Opera directory YAML frontmatter (charged/discharged/assessed counts)

### `ath check <crucible>`

Direct health check of a single crucible. Uses the same detection logic as patrol but targeted at a named window rather than scanning all panes.

---

## Key Implementation Patterns

### Internal Package Structure

```
cmd/ath/              # Binary entry point (shim)
internal/cli/         # Command implementations — one file per command
internal/athanor/     # Home directory, instance, config, MO/opera logic
internal/tmux/        # Tmux interaction library
```

### tmux.Runner

All tmux interaction goes through `tmux.Runner`. Methods include:
- `ListPanes()` — all panes across all sessions with coordinates and process info
- `ListWindows()` — all window names across all sessions
- `CapturePaneLines(target, n)` — last N lines of a pane
- `CheckCrucible(target)` — health state detection
- `SendKeys` / `SendKeysLiteral` — keystroke injection
- `Send(target, message, opts)` — full reliable send protocol with locking

### Detection Logic

Permission and stall detection use the same patterns in both `CheckCrucible()` (single-target) and `ath patrol` (all-panes scan):

- **Permission:** Co-occurrence guard — both `❯\s+\d+\.` (menu cursor) AND question text must appear in the same bottom-10 block. Prevents false positives from captured output showing another pane's content.
- **Exhausted:** Simple string match on "Context limit reached" in bottom-5 lines, filtered to Claude-process panes only.

---

## Infrastructure Dependencies

| Component | Purpose | Location |
|-----------|---------|----------|
| tmux | Session/window management, process hosting | System package |
| Claude Code CLI | Agent runtime | `claude` on PATH |
| systemd (user scope) | Timer-based automation | `~/.config/systemd/user/` |
| PushNotification | Artifex escalation | Built-in Claude Code tool |
| Egregore MCP | Session logging, Slack access | MCP server |

---

## Obsolete / Superseded

| Component | Superseded By |
|-----------|--------------|
| `claude-run` | `ath` CLI (direct claude launch via tmux) |
| `claude-perm-scan` | `ath patrol` (Go implementation) |
| `claude-stall-scan` | `ath patrol` (Go implementation) |
| `kadmon.md` | This document + evolved operational practice |
