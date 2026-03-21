# Task Triage: Implement Whisper CLI

**Triaged:** 2026-03-21
**Operator:** Matt Hull

## Goal

Implement the whisper CLI: a standalone Go binary that reliably delivers messages to Claude Code sessions running in tmux, replacing the buggy `tmux-claude-send` bash script. Three commands: `whisper send` (8-step protocol), `whisper idle` (2-check detection), `whisper wait-and-send` (combined). Installed to `~/.local/bin/whisper`.

## Completion Criteria

1. `whisper send %<pane> "hello world"` delivers and exits 0 — **verified by integration test** (ephemeral tmux session, capture-pane)
2. Special chars `!~#` + backticks + Unicode arrive literally — **verified by integration test**
3. 2KB message delivers without truncation — **verified by integration test** (send, capture, count chars)
4. Copy mode pane: exits copy mode and delivers — **verified by integration test** (enter copy-mode, send, capture)
5. Detached session: delivers via SIGWINCH wake — **verified by integration test** (detached session, send, capture)
6. `whisper idle` returns 0 when idle, 1 when busy — **verified by integration test** (session running `sleep` vs shell prompt)
7. Idle detection doesn't false-positive during inter-tool-call prompt flash — **manual smoke test** against live Claude Code session (known gap, cannot automate)
8. Two concurrent sends produce non-garbled output — **verified by integration test** (goroutines, capture, check message boundaries)
9. `whisper send --self "test"` works from within tmux — **verified by integration test**
10. `make check` passes (fmt + vet + lint + test) — **verified by running it**
11. `make install` puts binary at `~/.local/bin/whisper` — **verified by running it**
12. Gastown-reference skill updated with lessons learned — **verified by checking file was modified**

## Agent/Operator Boundary

- **Agent owns:** All implementation, all tests, running `make check`, pushing to master
- **Operator owns:** AC#7 manual smoke test (live Claude Code session), final `make install`
- **Gates:** None. Agent runs autonomously through completion. Push to master when `make check` passes.

## Project Goal Link

Whisper is the reliable transport layer for the athanor. Every automated interaction — marut nudging myrmidons, dispatch spawning crucibles, handoff chaining sessions — is a whisper. Flaky message delivery makes everything built on top flaky.

## Team Goal Link

Solo operator. This unblocks the Telegram MCP reply path (currently using buggy `tmux-claude-send`) and enables reliable automated agent-to-agent communication.

## Delivery Acceleration

Eliminates a class of bugs (garbled input, missed Enter, copy-mode swallowing messages) that cause manual intervention in automated workflows. Every minute debugging a stuck nudge is a minute not building.

## Scope

### In scope
- `internal/tmux/` — tmux.go (runner), sanitize.go, lock.go, send.go, capture.go, idle.go, wake.go
- `internal/cli/` — send_cmd.go, idle_cmd.go, wait_send_cmd.go (arg parsing, flag handling)
- Unit tests for sanitize, chunking, lock
- Integration tests for send, idle, concurrent sends, copy mode, detached session, large message, special chars
- `--self`, `-f`, `--skip-escape`, `--timeout` flags
- Spec retrospective: update gastown-reference skill with lessons learned

### Out of scope
- File-based message queue — separate tool, uses whisper as transport
- Multi-pane targeting (FindAgentPane) — future
- Startup dialog dismissal — separate tool
- Cross-process send locks — v1 is in-process only
- CI infrastructure — `make check` locally, push to master
- Migration wrapper (`tmux-claude-send` -> `whisper send`) — follow-up
- Athanor doc updates (domain-model.md, open-work-items.md) — follow-up

## Why This Approach

Clean reimplementation from Gastown's battle-tested protocol rather than porting Gastown code. The protocol is well-documented in the spec with line-number traceability. A standalone CLI keeps the tool simple (no framework deps) and installable anywhere. One-shot implementation via `/orchestrate` because the spec is complete and the reference code is available — no design decisions remain.

## References
- Spec: `~/code/whisper/specs/spec.md`
- Gastown reference: `~/code/gastown/internal/tmux/tmux.go`
- Gastown reference skill: `~/.claude/skills/gastown-reference/SKILL.md`
- Existing scaffold: `~/code/whisper/cmd/whisper/main.go`, `~/code/whisper/internal/cli/root.go`
- Makefile: `~/code/whisper/Makefile`

## Mise en Place

- **Knowledge:** Spec complete, Gastown reference code read and understood, skill created
- **Tooling:** Go toolchain ready (goimports, golangci-lint, staticcheck), `make check` works, tmux available on dev machine
- **Environment:** GitHub remote created (`matthull/whisper`, private), master pushed
- **Gaps:**
  - Blockers: **None**
  - Pre-work: **None**
  - Follow-ups (operator-owned): AC#7 manual smoke test, migration wrapper, athanor doc updates
