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
| **Athanor home** | `~/athanor/` (or `$ATH_HOME`). Root directory for all athanor state — shared materia, instance directories, global config. Lives outside any project repo. |
| **Instance** | A single athanor scoped to one Magnum Opus. Directory at `~/athanor/athanors/<name>/`. Contains config, opera, and symlinked shared materia. |
| **Shared materia** | Role files, geas, protocols shared across all instances. Live at `~/athanor/shared/` and are symlinked into each instance. |
| **Crucible** | A tmux window where an agent session runs. Named `marut-<athanor>` or `azer-<opus>`. |
| **Kindle** | Launch a marut for an athanor — create crucible, start session. |
| **Muster** | Launch an azer for an opus — create crucible, start session, verify. |
| **Reforge** | Kill a dead session and spawn fresh in the same crucible. |
| **Quiesce** | Graceful shutdown — verify no charged opera remain, kill the marut. |

---

## Directory Structure

```
~/athanor/                          # ATH_HOME
├── config.yml                      # Global config (defaults, preferences)
├── shared/                         # Shared materia (symlinked into instances)
│   ├── AGENTS.md
│   ├── azer.md
│   ├── marut.md
│   ├── muster.md
│   └── opus.md
└── athanors/                       # All instances
    ├── bugsnag/
    │   ├── athanor.yml             # Instance config (project path, models, etc.)
    │   ├── AGENTS.md               → ../shared/AGENTS.md
    │   ├── magnum-opus.md          # Authored per-instance (the goal)
    │   ├── marut.md                → ../shared/marut.md
    │   ├── azer.md                 → ../shared/azer.md
    │   ├── opus.md                 → ../shared/opus.md
    │   ├── muster.md               → ../shared/muster.md
    │   └── opera/                  # All opera (YAML frontmatter for status)
    └── sal-117/
        └── ...
```

**Key change from current system:** Athanors move from `specs/athanors/<name>/` (inside a project's specs repo) to `~/athanor/athanors/<name>/` (standalone). This eliminates the cross-worktree sync problem — all agents access the same filesystem path regardless of which worktree they run in. No more `git -C specs pull/push` dance. `[O:observation]` `[S:kadmon.md]`

**The athanor home is a git repository.** Trail durability is a core principle — opera, discharge records, and shared materia changes are version-controlled. `[S:spec.md]`

---

## Workflow / Dataflow

### Artifex launches an athanor

```
ath init myproject --project ~/code/myproject
  → creates ~/athanor/athanors/myproject/
  → symlinks shared materia
  → creates opera/ directory
  → writes template athanor.yml and magnum-opus.md

[artifex edits magnum-opus.md with the goal]

ath kindle myproject
  → reads athanor.yml for project path, marut model
  → creates tmux window "marut-myproject"
  → sets $ATHANOR=~/athanor/athanors/myproject/
  → launches claude --model sonnet --permission-mode auto with marut boot prompt
```

### Marut musters an azer

```
[marut decides an opus needs an azer]
[marut calls wtp add branch-name if code work — project-specific, not ath's job]

ath muster fix-nil-error.md --dir /path/to/worktree
  → reads $ATHANOR to find the instance
  → reads opus file for crucible naming
  → creates tmux window "azer-fix-nil-error"
  → launches claude --model opus --permission-mode auto with azer boot prompt
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
name: bugsnag
project: /home/matt/code/musashi    # Working directory for agents
marut_model: sonnet                  # Model for marut sessions
azer_model: opus                     # Model for azer sessions
```

`project` is optional — assessment-only or non-code athanors may not have one.

Model defaults: marut=sonnet, azer=opus. Override per-instance when needed. `[D:cost-vs-capability — maruts do mechanical ops, azers do craft work]`

### config.yml (global)

```yaml
home: ~/athanor                      # Redundant but explicit
shared: ~/athanor/shared             # Shared materia location
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
3. Symlinks all shared materia from `~/athanor/shared/`
4. Writes `athanor.yml` with name and project path (if provided)
5. Writes template `magnum-opus.md` with placeholder sections
6. Prints: "Athanor initialized. Edit magnum-opus.md, then `ath kindle <name>`."

**The template `magnum-opus.md`** includes the required sections (Goal, Abundant Satisfaction, Witnesses, Pre-loaded Context) with `[TODO]` placeholders. `[S:kadmon.md]`

**Exit codes:** 0=success, 1=already exists or error, 2=usage error

#### `ath kindle <name>`

Launch a marut for an athanor. The marut is the furnace — once kindled, it runs autonomously.

**What it does:**
1. Reads `~/athanor/athanors/<name>/athanor.yml`
2. Validates magnum-opus.md exists and has content (not just template placeholders)
3. Creates tmux window: `marut-<name>`
4. Launches claude session in the window:
   - Working directory: `athanor.yml → project` (or `~` if no project)
   - Model: `athanor.yml → marut_model` (default: sonnet)
   - Flags: `--permission-mode auto` (always) `[D:all-agents-auto-perms]`
   - Environment: `ATHANOR=~/athanor/athanors/<name>/`
   - Prompt: the marut boot prompt (reads AGENTS.md, magnum-opus.md, marut.md, muster.md, starts /loop 5m)
5. Verifies launch with `ath whisper idle marut-<name>`

**The boot prompt is encoded in the CLI**, not in agent specs. This means prompt improvements are CLI releases, not spec edits. `[D:prompt-as-infrastructure]`

**Exit codes:** 0=marut launched and verified idle, 1=launch failed, 2=usage error

#### `ath reforge <name>`

Kill a dead marut session and spawn fresh. The crucible endures; the session is reforged.

**What it does:**
1. Sends Ctrl-C + `exit` to the marut crucible to kill the old session
2. Waits briefly for cleanup
3. Re-launches with the reforge prompt variant (includes "check opera/ for in-progress work and the trail")
4. Verifies with `ath whisper idle`

**When to use:** Context exhaustion, crash, drift. The marut's context exhaustion protocol (marut.md) tells it to self-reforge when possible, but sometimes the artifex needs to force it. `[S:marut.md]`

**Exit codes:** 0=reforged and verified, 1=reforge failed, 2=crucible not found

#### `ath quiesce <name>`

Graceful shutdown of an athanor.

**What it does:**
1. Checks for charged opera — warns if any exist
2. Checks for active azer crucibles — warns if any exist
3. Kills the marut crucible
4. Prints summary: "Athanor <name> quiesced. Trail: N discharged opera."

**Does NOT** delete the instance directory or opera. The trail persists.

**Exit codes:** 0=quiesced, 1=error, 2=refused (active azers — use `--force` to override)

#### `ath status [<name>]`

Show athanor health at a glance.

**Without name — all athanors:**
```
ATHANOR     TYPE         MARUT    AZERS  OPERA (C/D/A)
bugsnag     state-based  active   1      1/5/3
sal-117     task-based   idle     0      0/2/2
blogging    task-based   -        0      0/0/0
```

**With name — detailed view:**
```
Athanor: bugsnag
Project: /home/matt/code/musashi
Marut: active (marut-bugsnag)
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

### Azer Management (called by marut)

#### `ath muster <opus-file> [--dir <path>] [--model <model>] [--name <name>]`

Launch an azer for a charged opus. The marut's primary dispatch command.

**Context detection:** Reads `$ATHANOR` to determine which athanor instance this is for. If `$ATHANOR` is not set, requires `--athanor <name>`. `[D:agent-context-via-env]`

**Arguments:**
- `<opus-file>` — Opus filename (resolved relative to `$ATHANOR/opera/`) or absolute path
- `--dir <path>` — Working directory for the azer. Default: project dir from athanor.yml. Use this to point to a worktree.
- `--model <model>` — Model override. Default: `athanor.yml → azer_model` (default: opus)
- `--name <name>` — Crucible name override. Default: `azer-<opus-slug>` derived from filename

**What it does:**
1. Resolves opus file path
2. Derives crucible name from opus filename (strips date prefix and .md)
3. Creates tmux window with crucible name
4. Launches claude session:
   - Working directory: `--dir` value or project dir
   - Model: `--model` value or athanor.yml default
   - Flags: `--permission-mode auto` (always) `[D:all-agents-auto-perms]`
   - Environment: `ATHANOR=$ATHANOR`
   - Prompt: azer boot prompt (reads AGENTS.md, azer.md, opus file, execute)
5. Verifies launch with `ath whisper idle <crucible-name>`
6. Prints: crucible name and verification status

**Sandbox creation is NOT part of muster.** The marut handles sandbox creation (e.g., `wtp add`) before calling `ath muster`, then passes the worktree path via `--dir`. This keeps project-specific tooling out of the athanor CLI. `[D:separation — ath is athanor ops, wtp is project ops]`

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

#### `ath inscribe [<name>] <opus-name>`

Create a new opus from the standard template. Writes to `$ATHANOR/opera/YYYY-MM-DD-<opus-name>.md` with YAML frontmatter `status: charged` and section placeholders.

**Exit codes:** 0=created, 1=error

### Utility

#### `ath version`

Print version, commit, build time (same as current whisper).

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
| `ath inscribe <TAB>` | Athanor instance names, then opus name |
| `ath muster <TAB>` | Opus files from `$ATHANOR/opera/` (charged only) |
| `ath cleanup <TAB>` | Active tmux windows matching `azer-*` pattern |
| `ath whisper send <TAB>` | Active tmux windows/panes |
| `ath whisper idle <TAB>` | Active tmux windows/panes |

**Athanor name completion** reads directory names from `~/athanor/athanors/`. This is the critical one — the artifex will type `ath kindle b<TAB>` and get `bugsnag`. `[D:athanor-name-completion-is-primary]`

**Implementation:** `ath completion zsh` outputs a zsh completion function. Dynamic completions (athanor names, tmux targets, opus files) use shell commands at completion time, not static lists.

---

## Design Decisions

### All agents launch with `--permission-mode auto`

Every claude session launched by `ath` (marut or azer) uses `--permission-mode auto`. No opt-out. Rationale: maruts do only mechanical operations (tmux, git, rg). Azers are already launched with auto perms. Permission prompts are the single biggest operational friction — every prompt blocks an agent until the permission manager or artifex intervenes. `[O:observation]` `[S:kadmon.md]`

**Consequence:** The `/permission-manager` may become unnecessary. Monitor after rollout.

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

`ath muster` does NOT create worktrees, Docker environments, or any project-specific sandbox. The marut handles sandbox creation separately (e.g., `wtp add`), then passes the path to `ath muster --dir <path>`. This keeps `ath` project-agnostic. `[D:separation]`

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
- `~/athanor/` — the athanor home (opera, shared materia, config)
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
1. Before launching, verify paths. If key paths are missing, warn but proceed (the `--permission-mode auto` on the session handles runtime perms, but file read/write outside the project dir may still be blocked by Claude Code's file access controls).

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
2. Move `specs/athanors/shared/` → `~/athanor/shared/`
3. Move each instance (`specs/athanors/<name>/`) → `~/athanor/athanors/<name>/`
4. Update symlinks to point to new shared location
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
| `ath inscribe` | Nice | Create opus from template |

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
| `ath muster --dir` points azer to worktree | Azer session starts in specified directory |
| `ath muster` derives crucible name from opus | `2026-03-24-fix-nil-error.md` → `azer-fix-nil-error` |
| `ath whisper send` delivers message | Identical behavior to current whisper send |
| `ath whisper idle` detects idle | Identical behavior to current whisper idle |
| `ath status` lists all athanors | Shows all instances with marut/azer/opera counts |
| `ath status <name>` shows detail | Correct crucible detection and opera breakdown |
| Zsh completion for athanor names | `ath kindle <TAB>` completes from instance directory names |
| Zsh completion for opus files | `ath muster <TAB>` completes from charged opera |
| Zsh completion for crucibles | `ath cleanup <TAB>` completes from `azer-*` tmux windows |

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
- [ ] Is `/permission-manager` still needed with universal `--permission-mode auto`?
- [ ] Can marut.md be simplified now that muster is one command?
- [ ] Should `ath kindle` also start the permission watcher, or is it fully unnecessary?

### Knowledge Capture
- [ ] Document the migration path for anyone with existing athanors in specs/
- [ ] Capture any tmux edge cases discovered during implementation
