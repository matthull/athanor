# Athanor CLI

Go CLI for the athanor agent orchestration system. Binary: `ath`, module: `github.com/matthull/athanor`.

## Skills Routing

| When you are... | Load skill |
|-----------------|------------|
| Implementing Go code, extracting patterns from Gastown | `/go-cli` |
| Writing or reviewing tests | `/go-cli` (testing strategy section) + `/unit-testing` |
| Working on Claude Code infrastructure (hooks, skills, settings) | `/agentic-dev` |
| Designing athanor system architecture | `/agentic-architecture` |

## Workflow

When `make check` passes, commit and push to master. No branches, no PRs — this is a single-developer project. Don't hold changes.

## Build and Test

```bash
make check          # fmt + vet + lint + test (all tests including integration)
make test           # all tests including integration (~30s)
go test ./... -short # unit tests only (~1s)
make install        # build and install to ~/.local/bin/ath
```

### Testing Strategy: Two Layers

This project uses a two-layer testing model. Both are required before committing.

**Layer 1 — Unit tests** (`go test ./... -short`): Pure logic — arg parsing, config, templates, status parsing. Fast, run constantly.

**Layer 2 — Acceptance tests** (`go test ./internal/cli/ -run TestATHFullLifecycle -timeout 120s`): Installed `ath` binary against real tmux. Exercises the full lifecycle: init, kindle, muster, status, opera, whisper, cleanup, reforge, quiesce. Requires tmux and `make install` first.

**When writing new commands:** Always add both a unit test for the parsing logic AND extend `TestATHFullLifecycle` with the new command's happy path and error cases. The acceptance test is what catches environmental bugs (tmux behavior, flag ordering, state management).

See `/go-cli` § Testing Strategy for the full rationale and patterns.

## Architecture

```
cmd/ath/              # Binary entry point
internal/cli/         # Command implementations + arg parsing
internal/athanor/     # Home directory, instance, config logic
internal/tmux/        # Tmux interaction library (from whisper)
```

### Key patterns

- **`splitArgs()`** for all commands with `<positional> [--flags]` — Go's `flag.Parse` stops at the first non-flag arg. See `/go-cli` § CLI argument parsing.
- **`tmux.Runner`** struct over global functions — explicit socket, testable
- **Sentinel errors** (`ErrNoServer`, `ErrSessionNotFound`) for control flow
- **Exact window matching** via `ListWindows()`, not `DisplayMessage` (tmux does prefix matching)

## Athanor Home

`~/athanor/` (or `$ATHANOR_HOME`) contains all athanor state:
- `shared/` — role files symlinked into each instance
- `athanors/<name>/` — per-instance directories with config, magna-opera, opera
- Initialized as a git repo for trail durability

## The Workshop Model

The athanor is a workshop, not a dark factory. The artifex walks in and out freely. The system supports a full spectrum from autonomous to hands-on — this is a core design requirement, not an edge case.

**How the artifex enters the workshop:**

- **Talk to the marut.** Switch to any marut's tmux window and direct it: "I need to investigate this stuck PR." The marut inscribes an opus and musters an azer. This is the natural entry when the work relates to an active MO.
- **`ath craft <athanor> <name> [<mo>]`** — Creates a lightweight opus and kindles an interactive crucible in one motion. The fast path for ad hoc collaborative work.
- **Walk up to any azer.** Switch to any azer's tmux window and start talking. Every crucible is interactive. The artifex can sit down at any bench at any time.

**How the artifex leaves:**

- Tell the azer to discharge the opus (outcome, reflection, proof captured to the trail).
- Or tell it to continue autonomously — the marut resumes supervision.
- Or just leave — the marut monitors for stalls and handles cleanup.

**The rule:** any work the artifex does should be inside an athanor, not in a raw Claude Code session. Raw sessions are invisible — no trail, no opera, no dashboard presence. If you find yourself opening a standalone Claude session for project work, that's a signal to use `ath craft` or talk to a marut instead.

## CLI Spec

The full CLI specification is at `specs/cli/spec.md`. The system spec is at `specs/spec.md`. Operations runbook is at `specs/kadmon.md`.
