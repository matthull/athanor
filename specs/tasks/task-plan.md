# Task Plan: Implement Whisper CLI

## EXECUTION PREREQUISITES — NO EXCEPTIONS
1. **Load `/orchestrate`** using the Skill tool before executing ANY phase.
   Do NOT skip this. Do NOT implement anything directly. Do NOT rationalize that the task is "simple enough"
   to do inline. The orchestrator DELEGATES — it never writes code, edits files, or implements.
2. Create a team via TeamCreate, then create TaskCreate entries for every phase.
   These are durable progress markers — if context is lost, they show exactly where to resume.
   Update task status as phases complete.

## Goal

Implement the whisper CLI: a standalone Go binary that reliably delivers messages to Claude Code sessions running in tmux. Three commands: `whisper send` (8-step protocol), `whisper idle` (2-check detection), `whisper wait-and-send` (combined).

## Completion Criteria

1. `make check` passes — fmt, vet, lint, test all green (includes unit + integration tests)
2. Integration tests verify ACs 1-6, 8-9 against real tmux sessions (isolated socket)
3. Binary installable via `make install` to `~/.local/bin/whisper`

## Agent/Operator Boundary

**Agent owns:** All implementation, all tests, running `make check`, committing, pushing to master
**Operator owns:** AC#7 manual smoke test, final `make install`, migration wrapper
**Gates at:** None — autonomous (notify-only profile)

## Triage Source

`~/code/whisper/specs/tasks/triage.md` — triaged 2026-03-21, confirmed by operator

## Phases

### 1. DESIGN — Task breakdown and handoffs

- Mode: Task-Lead-Implementer (TeamCreate) — task lead in full mode
- Mode resource: `~/.claude/skills/orchestrate/resources/task-lead-implementer-mode.md`
- Task lead reads the spec (`specs/spec.md`), the Gastown reference code (`~/code/gastown/internal/tmux/tmux.go`), loads `/gastown-reference` skill, and the existing scaffold
- Produces task breakdown + handoff documents at `specs/tasks/handoffs/`
- Handoffs should be ordered for incremental buildability:
  1. tmux runner (tmux.go) — foundation everything else depends on
  2. sanitize + lock — pure logic, independently testable
  3. send protocol (send.go, wake.go, capture.go) — the core 8-step protocol
  4. idle detection (idle.go) — 2-check protocol
  5. CLI wiring (send_cmd.go, idle_cmd.go, wait_send_cmd.go) — arg parsing, connects CLI to internal packages
  6. Integration tests — end-to-end verification against real tmux
- Completion: handoff files exist at `specs/tasks/handoffs/T*.md`
- Gate: none (autonomous)
- **Produced:** 6 handoff documents at specs/tasks/handoffs/T1-T6. Design decisions: Runner struct (not Tmux), SendOpts.LockHeld for wait-and-send, exported helpers on Runner, explicit socket passing (no global), 15s default timeout per spec.
- **Confidence:** All 6 files verified to exist via Glob. Handoffs reference correct Gastown line numbers and spec sections.
- **Gate outcome:** skipped (autonomous)

### 2. IMPLEMENT — Execute handoffs via task-lead-implementer

- Mode: Task-Lead-Implementer — implement-only (task lead retains context, fresh implementer per handoff)
- Mode resource: `~/.claude/skills/orchestrate/resources/task-lead-implementer-mode.md`
- One fresh implementer per handoff, sequential execution
- Each implementer reads the handoff + references Gastown code + applies gastown-reference skill
- Gate check per handoff: tests pass, lint clean
- Completion: all handoffs implemented, `go build` succeeds, unit tests pass
- **Produced:** T1-T6 all implemented. Files: internal/tmux/{tmux,sanitize,lock,send,wake,capture,idle}.go, internal/cli/{send_cmd,idle_cmd,wait_send_cmd}.go, tests in internal/tmux/{tmux_test,sanitize_test,lock_test,send_test,capture_test,idle_test,runner_integration_test,integration_test}.go and internal/cli/send_cmd_test.go. 9 integration tests, unit tests for sanitize (16), lock (4), error classification (6), chunking, prompt matching, status bar, idle logic, arg parsing (10).
- **Confidence:** All implementers reported tests passing. T6 ran full integration suite: 9/9 PASS in 29.5s. go vet clean across all packages. go build clean.
- **Gate outcome:** skipped (autonomous)

### 3. VERIFY — make check + integration tests

- Mode: Inline — orchestrator delegates to a single verification agent
- Agent runs `make check` (fmt + vet + lint + test) — this runs ALL tests including integration
- If integration tests need tmux (not `-short`), agent runs `make test` specifically
- Any failures: route through task lead for fix handoffs, re-implement, re-verify
- Completion: `make check` exits 0
- **Produced:** `make check` passed (fmt + vet + lint + test). Binary builds to /tmp/whisper and runs (`whisper version`). 28 tests across 2 packages all pass (cli: 10, tmux: 18). Integration tests: 30.57s. No formatting changes needed.
- **Confidence:** Full `make check` output captured. All quality gates green. Binary confirmed runnable.
- **Gate outcome:** skipped (autonomous)

### 4. SHIP — Commit and push to master

- Mode: Inline — orchestrator delegates commit + push
- Agent commits all changes with descriptive message
- Pushes to master (no branch, no PR — per project workflow)
- Runs `make check` one final time after commit to confirm clean state
- Completion: changes on master at origin, `make check` green

### 5. RETROSPECTIVE — Update gastown-reference skill

- Mode: Solo Agent
- Mode resource: `~/.claude/skills/orchestrate/resources/solo-agent-mode.md`
- Agent reads the implementation, identifies lessons learned about Go patterns, Gastown reference process, testing
- Updates `~/.claude/skills/gastown-reference/SKILL.md` with findings
- Completion: skill file modified with concrete lessons from this implementation
- Notify operator via Telegram: done, summary of what was built, what to smoke test

## Autonomy Profile

notify-only — no gates, telegram for blockers/completion
Resource: `~/.claude/skills/orchestrate/resources/autonomy-profiles/notify-only.md`

## Starting Artifacts

- Spec: `~/code/whisper/specs/spec.md` (complete, operator-confirmed)
- Triage: `~/code/whisper/specs/tasks/triage.md`
- Gastown reference: `~/code/gastown/internal/tmux/tmux.go` (lines 1191-1537, 2034-2541)
- Gastown reference skill: `~/.claude/skills/gastown-reference/SKILL.md`
- Existing scaffold: `~/code/whisper/cmd/whisper/main.go`, `~/code/whisper/internal/cli/root.go`
- Makefile: `~/code/whisper/Makefile`
- Go module: `~/code/whisper/go.mod` (github.com/matthull/whisper)

## Verification Strategy

| Criterion | Verification Method |
|-----------|-------------------|
| AC 1-6, 8-9 | Integration tests using isolated tmux socket (TestMain pattern from Gastown) |
| AC 7 | Operator manual smoke test (deferred — cannot automate) |
| AC 10 | `make check` exit code 0 |
| AC 11 | `make install` exit code 0 (operator runs post-implementation) |
| AC 12 | Diff of gastown-reference skill shows modifications |
| Code quality | `golangci-lint`, `go vet`, `gofmt` via `make check` |
