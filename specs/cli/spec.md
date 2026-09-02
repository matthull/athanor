# Athanor CLI (`ath`) Specification

**Status:** Draft
**Created:** 2026-03-24
**Source:** Whisper CLI (github.com/matthull/whisper) — renamed and extended

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[S:spec.md]` | From athanor system spec |
| `[S:kadmon.md]` | From athanor operations doc |
| `[S:marut.md]` | From marut role definition |
| `[S:muster.md]` | From muster protocol |
| `[D:reason]` | Design decision with rationale |
| `[E:whisper]` | Extends existing whisper functionality |
| `[O:observation]` | From operational experience |

---

## Overview

The athanor CLI (`ath`) is the operational backbone of the athanor system. It replaces the standalone `whisper` CLI, absorbing its inter-crucible communication capabilities and adding the full operational command set that maruts and the artifex need to run athanors.

**Why this exists:** The marut currently constructs multi-step bash sequences (tmux, wtp, claude launch commands) for every operational action. Each step triggers permission prompts, wastes context, and introduces error surface. A single `ath muster` replaces 3-4 fragile bash commands with one reliable operation. `[O:observation]`

**What it replaces:**
- `whisper` binary → `ath whisper` subcommand (full backward compatibility)
- Manual `tmux new-window` + `tmux send-keys` sequences → `ath kindle`, `ath muster`
- Manual `symlink + mkdir` init sequences → `ath init`
- Manual `rg` + `tmux list-windows` status checking → `ath status`

**The binary is `ath`** — short, fast to type, easy to complete. The project is "athanor CLI" but nobody types that. `[D:ergonomics]`

---

## Domain Dictionary

| Term | Definition |
|------|-----------|
| **Athanor home** | `~/athanor/` (or `$ATH_HOME`). Root directory for all athanor state — shared components, instance directories, global config. Lives outside any project repo. |
| **Instance** | A single athanor scoped to one Magnum Opus. Directory at `~/athanor/athanors/<name>/`. Contains config, opera, and symlinked shared components. |
| **Shared components** | Role files, geas, protocols shared across all instances. Live in the source repo's `shared/` directory and are symlinked directly into each instance. Repo location: `$ATHANOR_REPO` (default `~/code/athanor`). |
| **Crucible** | A tmux window where an agent session runs. Named `marut-<athanor>` or `azer-<opus>`. |
| **Kindle** | Launch a presence-driven role for an athanor — set kindled state, create crucible, start session. Default role: marut. |
| **Muster** | Launch an azer (or other opus-scoped role) for an opus — create crucible, start session, verify. Supports `--role` flag (default: azer). |
| **Reforge** | Kill a dead session and spawn fresh in the same crucible. |
| **Quiesce** | Graceful shutdown — clear kindled state, verify no charged opera remain, kill sessions. |
| **Kindled state** | A flag in MO state files indicating a presence-driven role should be kept alive. Discoverable via `rg`. Set by `ath kindle`, cleared by `ath quiesce`. |

---

## Directory Structure

```
~/athanor/                          # ATH_HOME
├── config.yml                      # Global config (defaults, preferences)
└── athanors/                       # All instances
    ├── my-project/
    │   ├── athanor.yml             # Instance config (project path, models, etc.)
    │   ├── AGENTS.md               → $ATHANOR_REPO/shared/AGENTS.md
    │   ├── magnum-opus.md          # Authored per-instance (the goal)
    │   ├── marut.md                → $ATHANOR_REPO/shared/marut.md
    │   ├── azer.md                 → $ATHANOR_REPO/shared/azer.md
    │   ├── opus.md                 → $ATHANOR_REPO/shared/opus.md
    │   ├── muster.md               → $ATHANOR_REPO/shared/muster.md
    │   └── opera/                  # All opera (YAML frontmatter for status)
    └── another-project/
        └── ...
```

Shared agent definitions live in the source repo (`$ATHANOR_REPO`, default `~/code/athanor`) under `shared/`. Each instance symlinks directly to the repo — no intermediate directory.

**Key change from current system:** Athanors move from `specs/athanors/<name>/` (inside a project's specs repo) to `~/athanor/athanors/<name>/` (standalone). This eliminates the cross-worktree sync problem — all agents access the same filesystem path regardless of which worktree they run in. No more `git -C specs pull/push` dance. `[O:observation]` `[S:kadmon.md]`

**The athanor home is a git repository.** Trail durability is a core principle — opera, discharge records, and shared components changes are version-controlled. `[S:spec.md]`

---

## Workflow / Dataflow

### Artifex launches an athanor

```
ath init myproject --project ~/code/myproject
  → creates ~/athanor/athanors/myproject/
  → symlinks shared components
  → creates opera/ directory
  → writes template athanor.yml and magnum-opus.md

[artifex edits magnum-opus.md with the goal]

ath kindle myproject
  → reads athanor.yml for project path, marut model
  → creates tmux window "marut-myproject"
  → sets $ATHANOR=~/athanor/athanors/myproject/
  → launches claude --model sonnet with marut boot prompt
```

### Marut musters an azer

```
[marut decides an opus needs an azer]
[marut calls wtp add branch-name if code work — project-specific, not ath's job]

ath muster fix-nil-error.md --worktree-path /path/to/worktree [--role azer]
  → reads $ATHANOR to find the instance
  → reads opus file for crucible naming
  → creates tmux window "azer-fix-nil-error"
  → launches claude --model opus with role-specific boot prompt (default: azer)
  → runs ath whisper idle to verify launch
```

### Agent communicates

```
ath whisper send azer-fix-nil-error "Status check — are you making progress?"
ath whisper idle azer-fix-nil-error --timeout 30s
ath whisper wait-and-send azer-fix-nil-error "Your opus has been updated" --timeout 60s
```

---

## Instance Configuration

### athanor.yml (per-instance)

```yaml
name: my-project
project: /home/matt/code/my-app     # Working directory for agents
marut_model: sonnet                  # Model for marut sessions
azer_model: opus                     # Model for azer sessions
```

`project` is optional — assessment-only or non-code athanors may not have one.

Model defaults: marut=sonnet, azer=opus. Override per-instance when needed. `[D:cost-vs-capability — maruts do mechanical ops, azers do craft work]`

### config.yml (global)

```yaml
home: ~/athanor                      # Redundant but explicit
defaults:
  marut_model: sonnet
  azer_model: opus
```

---

## Command Reference

### Instance Lifecycle

#### `ath init <name> [--project <path>]`

Create a new athanor instance.

**What it does:**
1. Creates `~/athanor/athanors/<name>/`
2. Creates `opera/` subdirectory
3. Symlinks all shared components from the source repo (`$ATHANOR_REPO/shared/`)
4. Writes `athanor.yml` with name and project path (if provided)
5. Writes template `magnum-opus.md` with placeholder sections
6. Prints: "Athanor initialized. Edit magnum-opus.md, then `ath kindle <name>`."

**The template `magnum-opus.md`** includes the required sections (Goal, Abundant Satisfaction, Witnesses, Pre-loaded Context) with `[TODO]` placeholders. `[S:kadmon.md]`

**Exit codes:** 0=success, 1=already exists or error, 2=usage error

#### `ath kindle <name> [--role <role>] [--mo <mo-name>]`

Launch a presence-driven role for an athanor. Sets kindled state so the liveness timer keeps it alive.

**Default role:** marut. Other presence-driven roles (perceiver, attendant) use `--role`.

**What it does:**
1. Reads `~/athanor/athanors/<name>/athanor.yml`
2. Validates the role file exists (e.g., `marut.md`, `perceiver.md`, `attendant.md`)
3. **Sets kindled state** for this role in the MO's state — marks it as "should be running"
4. If a session already exists for this role: **no-op** (idempotent — the liveness timer uses this)
5. Creates tmux window: `<role>-<name>` (e.g., `marut-my-project`, `perceiver-life`)
6. Launches claude session in the window:
   - Working directory: `athanor.yml → project` (or `~` if no project)
   - Model: per-role config (default: sonnet)
   - Environment: `ATHANOR=~/athanor/athanors/<name>/`
   - Prompt: role-specific boot prompt
7. Verifies launch with `ath whisper idle <role>-<name>`

**Idempotency is critical.** The liveness timer calls `ath kindle` every 5 minutes for each kindled role. If the session is already running, kindle is a no-op. This makes the timer implementation trivial — just call kindle for everything that should be alive.

**The boot prompt is encoded in the CLI**, not in agent specs. This means prompt improvements are CLI releases, not spec edits. `[D:prompt-as-infrastructure]`

**Exit codes:** 0=role launched and verified idle (or already running), 1=launch failed, 2=usage error

#### `ath reforge <name>`

Kill a dead marut session and spawn fresh. The crucible endures; the session is reforged.

**What it does:**
1. Sends Ctrl-C + `exit` to the marut crucible to kill the old session
2. Waits briefly for cleanup
3. Re-launches with the reforge prompt variant (includes "check opera/ for in-progress work and the trail")
4. Verifies with `ath whisper idle`

**When to use:** Context exhaustion, crash, drift. The marut's context exhaustion protocol (marut.md) tells it to self-reforge when possible, but sometimes the artifex needs to force it. `[S:marut.md]`

**Exit codes:** 0=reforged and verified, 1=reforge failed, 2=crucible not found

#### `ath quiesce <name> [--role <role>]`

Graceful shutdown of an athanor or a specific kindled role.

**Without `--role` — quiesce the whole athanor:**
1. Clears kindled state for all roles (so the liveness timer stops resurrecting them)
2. Checks for charged opera — warns if any exist
3. Checks for active azer crucibles — warns if any exist
4. Kills all presence-driven role crucibles (marut, perceiver, attendant, etc.)
5. Prints summary: "Athanor <name> quiesced. Trail: N discharged opera."

**With `--role` — quiesce a single role:**
1. Clears kindled state for that role only
2. Kills that role's crucible
3. Other kindled roles continue running

**Does NOT** delete the instance directory or opera. The trail persists.

**Exit codes:** 0=quiesced, 1=error, 2=refused (active azers — use `--force` to override)

#### `ath status [<name>]`

Show athanor health at a glance.

**Without name — all athanors:**
```
ATHANOR     TYPE         MARUT    AZERS  OPERA (C/D/A)
my-project  state-based  active   1      1/5/3
analytics   task-based   idle     0      0/2/2
blogging    task-based   -        0      0/0/0
```

**With name — detailed view:**
```
Athanor: my-project
Project: /home/matt/code/my-app
Marut: active (marut-my-project)
Azers:
  azer-fix-nil-error (active)
Opera:
  Charged:    fix-nil-error-in-proof-service
  Discharged: fix-rag-sync-race-condition, mute-sso-credential
  Assessed:   fix-missing-size-group-column
```

**How it works:**
- Lists instances from `~/athanor/athanors/`
- Checks tmux for crucible windows matching `marut-<name>` / `azer-*`
- Reads opera YAML frontmatter with `rg` for status counts
- Type detection: reads `magnum-opus.md` for task-based vs state-based signals

**Exit codes:** 0=success, 1=error

### Primus: Liveness Timer

A single athanor-wide systemd timer (every 5 minutes) that implements the first primus responsibility: ensuring all kindled presence-driven roles have active sessions.

**How it works:**
1. Scans all MO directories for kindled state (via `rg` or equivalent)
2. For each kindled role: checks tmux for an active session
3. Missing session? → calls `ath kindle <athanor> --role <role>` (which is idempotent)
4. Already running? → no-op

**Implementation:** Shell script invoked by a systemd user timer. No Go code needed — just file reads and tmux checks. Replaces per-role timers (e.g., the current `attunement-intake.timer`) with one universal heartbeat.

**Kindled state storage:** Stored in MO state files, discoverable via `rg`. Format TBD — could be YAML frontmatter in the MO file, a sidecar state file, or a dedicated section. The constraint: it must be grep-discoverable so the timer script stays simple.

**Relationship to other primus responsibilities:** The permission stall detection timer is also primus responsibility but operates independently. Future primus implementations (stall detection, health dashboards) may be folded into this timer or run as separate services. `[D:incremental-primus — automate what's proven, keep the rest manual]`

### Opus-Scoped Agent Management (called by marut)

#### `ath muster <opus-file> [--worktree-path <path>] [--model <model>] [--name <name>] [--role <role>]`

Launch an agent for a charged opus. The marut's primary dispatch command.

**Context detection:** Reads `$ATHANOR` to determine which athanor instance this is for. If `$ATHANOR` is not set, requires `--athanor <name>`. `[D:agent-context-via-env]`

**Arguments:**
- `<opus-file>` — Opus filename (resolved relative to `$ATHANOR/opera/`) or absolute path
- `--worktree-path <path>` — Working directory for the agent. Default: project dir from athanor.yml. Use this to point to a git worktree or any other temporary workspace outside the athanor's project dir. The name "worktree" is generic — it applies to any alternate workspace, not just `git worktree`.
- `--model <model>` — Model override. Default: `athanor.yml → azer_model` (default: opus)
- `--name <name>` — Crucible name override. Default: `azer-<opus-slug>` derived from filename
- `--role <role>` — Agent role file to use. Default: `azer`. Allows mustering non-azer opus-scoped roles (e.g., attendant for a specific opus). The role file (`<role>.md`) must exist in the athanor instance. `[D:backward-compatible — azer default preserved]`

**What it does:**
1. Resolves opus file path
2. Derives crucible name from opus filename (strips date prefix and .md)
3. Creates tmux window with crucible name
4. Launches claude session:
   - Working directory: `--worktree-path` value or project dir
   - Model: `--model` value or athanor.yml default
   - Environment: `ATHANOR=$ATHANOR`
   - Prompt: role-specific boot prompt (reads AGENTS.md, `<role>.md`, opus file, execute)
5. Verifies launch with `ath whisper idle <crucible-name>`
6. Prints: crucible name and verification status

**Sandbox creation is NOT part of muster.** The marut handles sandbox creation (e.g., `wtp add`) before calling `ath muster`, then passes the worktree path via `--worktree-path`. This keeps project-specific tooling out of the athanor CLI. `[D:separation — ath is athanor ops, wtp is project ops]`

**Exit codes:** 0=azer launched and verified, 1=launch failed, 2=usage error

#### `ath cleanup <crucible-name>`

Clean up after a discharged opus. Called by the marut during post-discharge.

**What it does:**
1. Kills the tmux window if still open
2. Prints confirmation

**Does NOT** handle worktree removal or Docker cleanup — those are project-specific. The marut handles those separately. `[D:separation]`

**Exit codes:** 0=cleaned up, 1=error, 2=crucible not found (not an error if already gone — idempotent)

### Communication (absorbs whisper)

All existing whisper functionality moves under `ath whisper`. The protocol, timing, and behavior are identical — this is a rename, not a rewrite. `[E:whisper]`

#### `ath whisper send <target> <message>`

Reliable message delivery to a tmux target. Uses the 8-step send protocol.

**Flags:**
- `-f <file>` — Read message from file instead of argument
- `--self` — Send to own pane
- `--skip-escape` — Omit Escape keystroke (non-Claude agents)
- `--timeout <duration>` — Max retry timeout (default: 15s)

**Exit codes:** 0=delivered, 1=failed, 2=usage error

#### `ath whisper idle <target> [--timeout <duration>]`

Wait for a target to become idle (showing prompt, not running tools). Uses 2-check idle detection.

**Exit codes:** 0=idle, 1=timeout, 2=target not found

#### `ath whisper wait-and-send <target> <message> [--timeout <duration>]`

Atomic idle-then-send. Combines idle detection with message delivery.

**Flags:** Same as `send` + `idle` combined.

**Exit codes:** 0=success, 1=timeout/failed, 2=usage error

### Opera Management

#### `ath opera [<name>]`

List opera for an athanor with status.

**Context:** Uses `<name>` argument, or `$ATHANOR` env var, or prompts.

**Output:**
```
STATUS      DATE        OPUS
charged     2026-03-24  fix-nil-error-in-proof-service
discharged  2026-03-23  fix-rag-sync-race-condition
discharged  2026-03-23  mute-sso-credential-not-found
assessed    2026-03-22  fix-missing-size-group-column
```

#### `ath inscribe <athanor> <mo> --job <job> --intent <text> [--name <name>] [--muster]`

Create an opus file with correct path, date prefix, and YAML frontmatter. Generates filename from date + slugified intent (or `--name` override). Writes `status: charged`, `inscribed`, `magnum_opus`, and `job` fields. Intent becomes the opus body under a heading. The `--job` flag is required and validated against the job registry (`shared/jobs/<job>/JOB.md`); use `general` when no specific role fits.

With `--muster`, immediately launches an azer for the opus after creation.

**Exit codes:** 0=created (prints opus path), 1=runtime error, 2=usage error

#### `ath collaborate <mo> --job <job> --intent <text>`

Peer-collaboration shorthand designed to be called from a running crucible. Infers athanor from `$ATHANOR` env var, creates the opus (same as `ath inscribe` internally), adds collaboration context ("Inscribed by: \<calling crucible\>. Whisper back to \<calling crucible\> when complete."), and musters the peer immediately. The `--job` flag is required; use `general` when no specific role fits.

**Exit codes:** 0=created and mustered (prints crucible name), 1=runtime error, 2=usage error

### Utility

#### `ath version`

Print version, commit, build time (same as current whisper).

#### `ath services [--json]`

List all athanor systemd service dependencies with their current status. Serves dual purpose: health check (are all required services running?) and dependency declaration (the command's implementation IS the canonical list).

**Output (text):**
```
ATH SERVICES

UNIT                          STATUS    LAST TRIGGER          NEXT
athanor-liveness.timer        active    10:31 (4m ago)        10:36
attunement-intake.timer       active    10:30 (5m ago)        11:00
voice-notes-process.timer     active    10:30 (5m ago)        10:35

All services healthy.
```

**Output (JSON with `--json`):**
```json
{
  "timestamp": "10:35:01",
  "services": [
    {
      "unit": "athanor-liveness.timer",
      "description": "Keeps kindled presence-driven roles alive (5-min cadence)",
      "status": "active",
      "last_trigger": "Sun 2026-04-20 10:31:19 CDT",
      "next_trigger": "Sun 2026-04-20 10:36:19 CDT"
    }
  ],
  "all_healthy": true
}
```

**Canonical service list** (defined in Go code, one place):
- `athanor-liveness.timer` — keeps kindled presence-driven roles alive (5-min cadence)
- `attunement-intake.timer` — processes voice note signals into attunement corpus (30-min cadence)
- `voice-notes-process.timer` — transcribes voice notes from GDrive (5-min cadence)

Adding a new service dependency means adding it to the `athanorServices` slice in `internal/cli/services_cmd.go`.

**Color coding** (when stdout is a terminal and `NO_COLOR` is not set):
- Green: active
- Red: inactive/failed
- Yellow: unknown

**Exit codes:** 0=all services active, 1=any service inactive/failed, 2=usage error

#### `ath completion zsh`

Generate zsh completion script. Install with:

```bash
ath completion zsh > ~/.zsh/completions/_ath
```

---

## Zsh Completion

Zsh completion is a first-class requirement, not an afterthought. `[D:ergonomics — the artifex uses these commands frequently]`

**What completes:**

| Context | Completion source |
|---------|------------------|
| `ath <TAB>` | All subcommands (init, kindle, muster, status, etc.) |
| `ath kindle <TAB>` | Athanor instance names (from `~/athanor/athanors/`) |
| `ath reforge <TAB>` | Athanor instance names |
| `ath quiesce <TAB>` | Athanor instance names |
| `ath status <TAB>` | Athanor instance names |
| `ath opera <TAB>` | Athanor instance names |
| `ath inscribe <TAB>` | Athanor instance names, then MO names, then flags (`--intent`, `--job`, `--name`, `--muster`) |
| `ath inscribe --job <TAB>` | Job names from `shared/jobs/` |
| `ath collaborate <TAB>` | MO names (from `$ATHANOR`), then flags (`--intent`, `--job`) |
| `ath collaborate --job <TAB>` | Job names from `shared/jobs/` |
| `ath muster <TAB>` | Opus files from `$ATHANOR/opera/` (charged only) |
| `ath cleanup <TAB>` | Active tmux windows matching `azer-*` pattern |
| `ath whisper send <TAB>` | Active tmux windows/panes |
| `ath whisper idle <TAB>` | Active tmux windows/panes |

**Athanor name completion** reads directory names from `~/athanor/athanors/`. This is the critical one — the artifex will type `ath kindle m<TAB>` and get `my-project`. `[D:athanor-name-completion-is-primary]`

**Implementation:** `ath completion zsh` outputs a zsh completion function. Dynamic completions (athanor names, tmux targets, opus files) use shell commands at completion time, not static lists.

---

## Design Decisions

### No `--permission-mode auto`

`--permission-mode auto` is not used. The feature is not reliably available across all Anthropic account types. Agents use default Claude Code permission mode and rely on the project's `.claude/settings.json` allowlists for operational smoothness. `[O:observation — auto mode unavailable on main Anthropic account]`

### Prompt templates are CLI infrastructure

The boot prompts for marut and azer sessions are encoded in the `ath` binary, not in the agent spec files. When a prompt needs updating, it's a CLI release — not a spec edit that requires hot-patching running sessions. The prompts reference the spec files (AGENTS.md, marut.md, etc.) but the boot sequence itself is infrastructure. `[D:prompt-as-infrastructure]`

**Marut boot prompt:**
```
Read $ATHANOR/AGENTS.md, then read $ATHANOR/magnum-opus.md, then read
$ATHANOR/marut.md, then read $ATHANOR/muster.md. You are the marut for
this athanor. Start /loop 5m and begin your operational cycle.
```

**Marut reforge prompt:**
```
Read $ATHANOR/AGENTS.md, then read $ATHANOR/magnum-opus.md, then read
$ATHANOR/marut.md, then read $ATHANOR/muster.md. You are the marut for
this athanor. Check opera/ for in-progress work and the trail. Start
/loop 5m and resume your operational cycle.
```

**Azer boot prompt:**
```
Read $ATHANOR/AGENTS.md, then read $ATHANOR/azer.md. Your opus is at
$ATHANOR/opera/<opus-file>.md. Read it and execute.
```

### Sandbox creation is external

`ath muster` does NOT create worktrees, Docker environments, or any project-specific sandbox. The marut handles sandbox creation separately (e.g., `wtp add`), then passes the path to `ath muster --worktree-path <path>`. This keeps `ath` project-agnostic. `[D:separation]`

### `$ATHANOR` as agent context

Agent-facing commands (`muster`, `cleanup`, `opera`, `inscribe`) read `$ATHANOR` from the environment to know which instance they're operating on. `ath kindle` sets this when launching sessions. Human-facing commands (`init`, `kindle`, `reforge`, `status`) take the athanor name as a positional argument. All commands support `--athanor <name>` as an explicit override.

### Athanor home outside projects

`~/athanor/` lives outside any project repository. Athanors are conceptual infrastructure — they orchestrate work across projects, not within them. The cross-worktree sync problem (the biggest friction from the first firing) disappears: all agents read/write `~/athanor/` directly regardless of which worktree they run in. `[O:observation]` `[S:spec.md]`

---

## Claude Code Permissions Management

File path permissions are a foundational operational concern for the athanor system. With `~/athanor/` living outside any project, and agents working across repos and worktrees, Claude Code's file permissions will block operations silently unless configured upfront. This is athanor infrastructure — keeping the forge burning means keeping permissions out of the way. `[O:observation]` `[D:permissions-as-infrastructure]`

### `ath init` permissions setup

When `ath init` creates a new instance (or when `ath kindle` launches a session), the CLI should **verify and report** whether the required paths are accessible in the user's Claude Code settings. It does not modify settings autonomously — it tells the artifex what to add.

**Paths that must be readable/writable by all athanor agents:**
- `~/athanor/` — the athanor home (opera, shared components, config)
- `/tmp/` — scratch space for message files, temp artifacts
- Each project path referenced by an athanor's `athanor.yml → project`

**`ath init` behavior:**
1. After creating the instance, check whether `~/athanor/` and the project path (if any) are in the user's Claude Code `settings.json` allowed paths
2. If not, print a clear remediation message:
   ```
   [permissions] ~/athanor/ is not in Claude Code allowed paths.
   Add to ~/.claude/settings.json → allowedPaths:
     "~/athanor/"
   ```
3. This is advisory, not blocking — the init still succeeds

**`ath kindle` behavior:**
1. Before launching, verify paths. If key paths are missing, warn but proceed (file read/write outside the project dir may still be blocked by Claude Code's file access controls — ensure `.claude/settings.json` allowlists cover athanor paths).

**Future (Phase 3):** `ath perms check` command that audits all athanor instances and reports which paths need adding. `ath perms fix` that generates the settings.json patch.

### Artifex responsibility

The artifex should proactively add all important project directories under `~/code/` to Claude Code's global settings with read/write permissions, plus `~/athanor/` and `/tmp/`. This is a one-time setup that prevents recurring friction across all athanor operations.

---

## Migration

### GitHub repo rename

Rename the GitHub repo: `github.com/matthull/whisper` → `github.com/matthull/athanor`. GitHub will set up automatic redirects from the old URL. `[D:clean-break]`

### From whisper

1. Rename `~/code/whisper/` to `~/code/athanor/`
2. Update Go module path: `github.com/matthull/whisper` → `github.com/matthull/athanor`
3. Update GitHub remote: `git remote set-url origin git@github.com:matthull/athanor.git`
4. Rename `cmd/whisper/` to `cmd/ath/`
5. Move current `send`, `idle`, `wait-and-send` under a `whisper` subcommand
6. `make install` produces `~/.local/bin/ath` instead of `~/.local/bin/whisper`
7. Optionally: install `whisper` as a shell alias to `ath whisper` for backward compatibility during transition

### From specs/athanors/ to ~/athanor/

1. Create `~/athanor/` directory structure
2. Move each instance (`specs/athanors/<name>/`) → `~/athanor/athanors/<name>/`
3. Update symlinks to point to source repo shared location (`$ATHANOR_REPO/shared/`)
5. Write `athanor.yml` for each instance
6. Initialize `~/athanor/` as a git repo
7. Update `marut.md` and `muster.md` references from `$ATHANOR` (still works — just points to new path)
8. Update `kadmon.md` launch commands to use `ath kindle` / `ath reforge`

---

## Phasing

### Phase 1 — MVP (absorb whisper + core ops)

| Command | Priority | Notes |
|---------|----------|-------|
| `ath whisper send` | Must | Direct port from whisper |
| `ath whisper idle` | Must | Direct port from whisper |
| `ath whisper wait-and-send` | Must | Direct port from whisper |
| `ath init` | Must | Create instance with symlinks and template |
| `ath kindle` | Must | Launch marut — the primary new capability |
| `ath muster` | Must | Launch azer — the key friction reducer |
| `ath status` | Must | Athanor health at a glance |
| `ath version` | Must | Already exists in whisper |
| `ath completion zsh` | Must | With athanor name completion |

**Phase 1 deliverable:** The artifex can `ath init`, `ath kindle`, and the marut can `ath muster` instead of constructing bash sequences. Whisper works unchanged under new name.

### Phase 2 — Full lifecycle

| Command | Priority | Notes |
|---------|----------|-------|
| `ath reforge` | Should | Kill + relaunch marut |
| `ath cleanup` | Should | Post-discharge crucible cleanup |
| `ath quiesce` | Should | Graceful shutdown |
| `ath opera` | Should | List opera with status |
| `ath inscribe` | Should | Create opus with required job, intent, optional auto-muster |
| `ath collaborate` | Should | Inscribe + muster peer from crucible (friction parity with Agent tool) |
| `ath services` | Should | List and check systemd service dependencies |

### Phase 3 — Enhanced operations

| Capability | Notes |
|------------|-------|
| Enhanced status (context %, stall detection) | Requires pane capture + heuristics |
| `ath trail` | Rich trail summary with timing and outcomes |
| `ath perms` | Start/check permission manager (may be unnecessary with auto perms) |
| Cross-athanor status dashboard | For when multiple athanors run simultaneously |

---

## Test Scenarios

| Scenario | Behavior to verify |
|----------|-------------------|
| `ath init` creates correct structure | Directory, symlinks, template files all present and valid |
| `ath init` with `--project` writes athanor.yml | Config file has correct project path |
| `ath init` on existing name | Refuses with clear error (no silent overwrite) |
| `ath kindle` launches marut in tmux | Window created, claude session started, verified idle |
| `ath kindle` with missing magnum-opus.md | Refuses with clear error |
| `ath kindle` with template-only magnum-opus.md | Warns that goal hasn't been written |
| `ath muster` from agent context | Reads `$ATHANOR`, creates crucible, launches azer, verifies |
| `ath muster --worktree-path` points azer to worktree | Azer session starts in specified directory |
| `ath muster` derives crucible name from opus | `2026-03-24-fix-nil-error.md` → `azer-fix-nil-error` |
| `ath whisper send` delivers message | Identical behavior to current whisper send |
| `ath whisper idle` detects idle | Identical behavior to current whisper idle |
| `ath status` lists all athanors | Shows all instances with marut/azer/opera counts |
| `ath status <name>` shows detail | Correct crucible detection and opera breakdown |
| `ath services` lists all services | Shows all canonical service deps with status |
| `ath services --json` produces JSON | Valid JSON with timestamp, services array, all_healthy |
| `ath services` exit code reflects health | 0 if all active, 1 if any inactive/failed |
| Zsh completion for athanor names | `ath kindle <TAB>` completes from instance directory names |
| Zsh completion for opus files | `ath muster <TAB>` completes from charged opera |
| Zsh completion for crucibles | `ath cleanup <TAB>` completes from `azer-*` tmux windows |

---

## Future: Multi-Model Support

Currently the athanor runs exclusively on Claude via Claude Max accounts. As automation scales, Anthropic's terms around automated usage may become a constraint. The `ath` CLI should eventually support launching agents on different LLMs — particularly for simple automated tasks (permission scanning, status checks, basic triage) that don't need Claude's full capability.

**Near-term candidate:** Gemini — we have a license with excess tokens, making it the default for non-Claude usage despite not being the cheapest/fastest for simple tasks. Could handle permission-manager duties, basic monitoring loops, and other mechanical operations.

**What this means for `ath`:** The `--model` flag on `kindle`/`muster` currently maps to Claude model names (sonnet, opus). Eventually this should support a provider prefix (e.g., `gemini:flash`, `claude:sonnet`) and the CLI handles launching the right runtime. Not Phase 1 — but the flag structure should not preclude it.

---

## Future: System Architecture Revisit

The current athanor structure is ad-hoc — built for experimentation during bootstrap. Now that the system has its legs under it, a ground-up architecture pass is needed. Key questions:

- **Where do agent roles live?** Currently in the source repo's `shared/` directory (marut.md, azer.md), symlinked into instances. Are these the right abstractions? Should roles be more composable? Should they live in the `ath` binary itself (as embedded templates) vs. as editable files?
- **What's the right split between CLI and materia?** Boot prompts are moving into the CLI (prompt-as-infrastructure). What else should migrate from markdown specs to code?
- **Instance config vs. global config vs. embedded defaults?** Currently three layers (global config.yml, instance athanor.yml, CLI defaults). Is this the right layering?
- **Opera storage and trail management.** Flat directory + YAML frontmatter works. Does it scale? Do we need indexing, search, archival?
- **Shared components versioning.** Currently symlinked and identical across all instances. When do instances need to diverge? How do we handle that?

This is a separate spec effort — not part of the CLI Phase 1. But the CLI should be designed with awareness that the underlying structure will evolve. Keep the CLI's assumptions about directory layout isolated so they're easy to change.

---

## Out of Scope

- **Project-specific tooling** — `wtp`, Docker, sandbox creation/cleanup. These stay as project tools.
- **Automated primus** — `ath` doesn't auto-launch maruts on a schedule. The artifex kindles manually.
- **Beholder role** — Channel scanning and auto-dispatch. Future capability, not part of initial CLI.
- **Cross-machine sync** — Athanor home is local. Remote sync (if needed) is a separate concern.
- **Web UI / TUI dashboard** — The Dwarf Fortress-style view is deferred. `ath status` is the text-based equivalent.
- **`whisper` backward-compat binary** — A shell alias is fine. No need to install two binaries.

---

## Retrospective

When this feature is complete, review:

### Documentation Updates
- [ ] Update `kadmon.md` to use `ath kindle` / `ath reforge` commands
- [ ] Update `marut.md` to reference `ath muster` instead of manual tmux sequences
- [ ] Update `muster.md` to use `ath` commands
- [ ] Update athanor `spec.md` implementation status table

### Workflow Improvements
- [ ] Can marut.md be simplified now that muster is one command?

### Knowledge Capture
- [ ] Document the migration path for anyone with existing athanors in specs/
- [ ] Capture any tmux edge cases discovered during implementation
