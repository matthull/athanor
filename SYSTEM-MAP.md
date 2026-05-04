# Athanor System Map

A matter-of-fact map of what's implemented, where it lives, and how it connects. Not a spec, not philosophy — just what's here.

Last refreshed: 2026-05-04

---

## Two Trees: Source and Runtime

The system spans two directory trees.

**Source repo** (`~/code/athanor/`): The `ath` CLI, agent role definitions, job definitions, specs. Everything that gets built or symlinked into the runtime.

```
~/code/athanor/
├── cmd/ath/              # CLI binary entry point
├── internal/             # CLI implementation (cli/, athanor/, tmux/)
├── shared/               # Agent role files + job definitions (symlinked into instances)
│   ├── AGENTS.md, azer.md, marut.md, perceiver.md, attendant.md, opus.md, muster.md
│   └── jobs/             # Job registry: assessor, coder, general, investigator,
│                         #   project-manager, qa-specialist, solution-architect
├── specs/                # System spec, CLI spec, kadmon runbook, feature specs
├── scripts/systemd/      # Service/timer unit templates
└── hooks/                # Claude Code hooks shipped with the project
```

**Runtime home** (`~/athanor/`, or `$ATHANOR_HOME`): Where athanor instances live. Initialized as a git repo for trail durability. Each instance is a named directory under `athanors/`.

```
~/athanor/
└── athanors/
    ├── musashi/           # Work domain — software development
    ├── athanor-architect/ # Meta domain — evolving the athanor system itself
    └── maneframe/         # Life domain — personal systems
```

---

## Athanor Instance Structure

Each instance in `athanors/<name>/` has the same shape:

```
athanors/<name>/
├── athanor.yml            # Config: name, project path
├── AGENTS.md → shared/    # Symlinked role files (azer.md, marut.md, etc.)
├── jobs/ → shared/jobs/   # Symlinked job definitions (simple tzurot)
├── formulae/              # Athanor-specific collaboration patterns (compound tzurot)
│   └── <name>/
│       └── FORMULA.md     # Formula definition: topology + inscription guidance
└── magna-opera/
    └── <mo-name>/
        ├── <mo-name>.md   # Magnum Opus definition (intent, witnesses, calcinatio)
        └── opera/          # Opera files — the trail of work done under this MO
```

Role files and job definitions are symlinks from the source repo's `shared/` directory. Formulae are athanor-specific (not symlinked) — each athanor defines its own collaboration patterns. `ath init` creates the instance; `ath sync` reconciles symlinks when the source repo adds new files.

---

## The `ath` CLI

**Lifecycle** — standing up and tearing down athanor infrastructure:
- `ath init` creates an instance. `ath sync` reconciles symlinks.
- `ath kindle` launches a presence-driven role (marut, perceiver, attendant) in a tmux crucible.
- `ath reforge` kills and relaunches a marut. `ath quiesce` gracefully shuts down.

**Work management** — creating and assigning work:
- `ath inscribe` creates an opus (unit of work) under a Magnum Opus with a required job role.
- `ath muster` launches an azer in a new tmux crucible for an opus.
- `ath collaborate` combines inscribe + muster for peer azers (used from within a crucible).
- `ath craft` creates a lightweight interactive session for artifex + azer collaboration.
- `ath craft-mo` creates a new Magnum Opus interactively.

**Observation** — seeing system state:
- `ath status` shows athanor health (tmux sessions, charged opera, running roles).
- `ath dashboard` gives an at-a-glance system overview across all athanors.
- `ath opera` lists opera with their status. `ath trail` shows collaboration topology.
- `ath patrol` scans tmux panes for permission prompts and stalled sessions.
- `ath check` checks crucible health. `ath services` checks service dependencies.

**Communication** — inter-agent messaging:
- `ath whisper send` sends a message to a tmux target.
- `ath whisper idle` waits for a target to become idle.
- `ath whisper wait-and-send` combines both — wait for idle, then send.

---

## Scheduled Services (systemd)

These run on timers and provide the autonomous glue between interactive sessions. Unit files live in `~/.config/systemd/user/`; templates in `scripts/systemd/`.

| Service | Interval | What it does |
|---------|----------|--------------|
| `athanor-liveness` | 5 min | Checks that kindled roles (maruts, perceivers, attendants) have active tmux sessions. Relaunches if missing. |
| `agentic-inbox` | 1 min | Processes the agentic domain inbox — routes incoming signals to the right athanor. |
| `musashi-inbox` | 1 min | Processes the musashi domain inbox. |
| `attunement-intake` | 30 min | Invokes the perceiver to process attunement signals (sensor data, state changes). |
| `voice-notes-process` | 5 min | Downloads and processes voice recordings from Google Drive. |
| `triage-router` | path-activated | Routes signals from the triage inbox to appropriate domains. |
| `headless-log-cleanup` | daily | Cleans up headless Claude session logs older than 48 hours. |
| `telegram-mcp` | always-on | Telegram MCP server (SSE singleton) — the escalation and notification channel. |
| `yadm-sync` | 10 min | Syncs dotfiles via yadm. |

---

## Claude Code Integration

MCP servers (telegram, Slack, Keep, Linear, etc.) are auto-discovered — not listed here. The project's `CLAUDE.md` contains a **skill routing table** that maps task types to the right skills to load — this is the primary way agents orient to the tools available for a given kind of work.

---

## How Things Connect

```
  Source Repo (~/code/athanor/)
    │
    ├── builds ──→ ath CLI (installed to ~/.local/bin/ath)
    │
    └── symlinks ──→ Athanor Instances (~/athanor/athanors/*)
                       │
                       ├── tmux sessions ──→ Crucibles (maruts, azers, perceivers, attendants)
                       │                       │
                       │                       ├── whisper ──→ inter-agent messages
                       │                       ├── telegram-mcp ──→ artifex escalation
                       │                       └── egregore-mcp ──→ Slack, meetings
                       │
                       └── opera/ ──→ trail of work (git-tracked)
                       
  Systemd Timers
    ├── liveness ──→ checks tmux sessions, relaunches maruts
    ├── inboxes ──→ route incoming signals to athanors
    ├── attunement ──→ perceiver intake cycle
    └── voice-notes ──→ transcription pipeline
```

---

## Known Gaps

- **`ath tzurot` / `ath formulae`** — CLI commands for listing tzurot (jobs + formulae) and formula detail are designed (`specs/formulae-system/spec.md`) but not yet implemented. `ath jobs` exists for listing jobs. Use `ls $ATHANOR/formulae/` for formulae until CLI ships.
- **`--formula` flag** — The `--formula` flag on `ath inscribe`, `ath collaborate`, and `ath muster` is designed but not yet implemented. Formulae can still be read manually by azers — the boot prompt just won't include them automatically yet.

---

## Keeping This Map Current

This document is a snapshot. The "Last refreshed" date at the top tells you how stale it is. If you're working on system structure and notice drift, update it.
