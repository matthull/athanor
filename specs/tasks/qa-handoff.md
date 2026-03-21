# QA Handoff: Whisper CLI Remaining Work

**Date:** 2026-03-21
**Context:** Implementation complete and shipped (6c4a8f9 on master). `make check` passes. This handoff covers gaps identified during verification.

## Prerequisites

- Start Claude Code in a **fresh tmux pane** (new shell sources .zshenv, gets go/bin on PATH)
- Verify: `which goimports` should return `~/go/bin/goimports` without any workarounds
- If it doesn't, the .zshenv PATH fix didn't take effect — investigate before proceeding

## Task 1: Add TestSendSelf integration test (AC 9)

**File:** `~/code/whisper/internal/tmux/integration_test.go`

Add a test that verifies sending to a pane by its ID (the `--self` path):
1. Create test session via `createTestSession`
2. Call `waitForShellReady`
3. Get pane ID: `r.DisplayMessage(session, "#{pane_id}")` — returns e.g. `%42`
4. `r.Send(paneID, "SELF_MARKER_xxx", SendOpts{})` — send using pane ID as target
5. Poll capture-pane for the marker

Follow existing patterns in integration_test.go (testing.Short skip, unique session name, deadline polling).

**Verify:** `go test -v -run TestSendSelf ./internal/tmux/`

## Task 2: Commit outstanding changes

Files to stage and commit:
- `internal/tmux/integration_test.go` (the new test)
- `.claude/settings.local.json` (permissions for Go project)

The gastown-reference skill update is at `~/.claude/skills/gastown-reference/SKILL.md` (outside this repo, no commit needed).

**Commit message:** `Add TestSendSelf integration test, project settings`
**Push to master.**

## Task 3: Run full make check

`make check` must pass with all tests including the new one.

## Task 4: Operator smoke tests

### AC 7: Idle detection false positive check
1. `make install` (installs to ~/.local/bin/whisper)
2. Start a Claude Code session in a tmux pane, note pane ID (`echo $TMUX_PANE`)
3. Give Claude a multi-tool-call task (e.g., "read these 5 files: ...")
4. While Claude is actively working, from another terminal:
   ```
   whisper idle %<pane> --timeout 10s
   echo $?
   ```
5. Expected: exit code 1 (timeout — Claude is busy, not idle)
6. Failure: exit code 0 during inter-tool-call prompt flash

### AC 11: make install
```
make install
whisper version
whisper send --help
```

## Done-when

- TestSendSelf passes
- `make check` green (all tests)
- Changes committed and pushed to master
- AC 7 smoke test performed (pass or documented finding)
- AC 11 verified (whisper binary at ~/.local/bin/whisper)
