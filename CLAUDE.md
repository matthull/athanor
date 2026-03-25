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
- `athanors/<name>/` — per-instance directories with config, magnum-opus, opera
- Initialized as a git repo for trail durability

## CLI Spec

The full CLI specification is at `specs/cli/spec.md`. The system spec is at `specs/spec.md`. Operations runbook is at `specs/kadmon.md`.
