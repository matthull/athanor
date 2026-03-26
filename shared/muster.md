# Muster Protocol — Marut Runbook

How the marut kindles crucibles, charges azers, reforges sessions, and monitors progress.

These are core marut responsibilities — not operator tasks. The marut manages the lifecycle of all azers under its Magnum Opus.

**`$ATHANOR`** is set in your environment to this athanor's directory path. All commands below use it.

---

## Kindling a Crucible (launching a azer)

When the marut has a charged opus that needs execution:

**1. Create a sandbox (if needed)** — a sandbox isolates code changes. If the opus involves code changes to this repo, create one:

```bash
wtp add <branch-name>
```

This creates a worktree with its own branch, Docker environment, and ports. Note the worktree path from the output.

**Skip the sandbox** if the opus is research, investigation, assessment, or work that doesn't modify this repo's code. The azer runs from the main repo directory instead. Not every opus needs code isolation — apply judgment.

**2. Create the crucible** — a tmux window for the azer:

```bash
tmux new-window -n azer-<opus-name>
```

**3. Launch the session** — send the claude command to the fresh crucible:

```bash
tmux send-keys -t azer-<opus-name> "cd <worktree-path> && claude --model opus --permission-mode auto \"Read $ATHANOR/AGENTS.md, then read $ATHANOR/azer.md. Your opus is at $ATHANOR/opera/<opus-file>.md. Read it and execute.\"" Enter
```

**4. Verify launch** — check the crucible has an active session:

```bash
whisper idle azer-<opus-name>
```

If `wtp` or `whisper` behave unexpectedly, escalate. Do not improvise workarounds.

---

## Reforging (session restart)

When a azer session dies (context exhaustion, crash, unrecoverable drift), reforge it — kill the session and spawn fresh in the same crucible. The crucible endures; the session is reforged.

**1. Kill the dead session:**

```bash
tmux send-keys -t azer-<opus-name> C-c C-c
# Wait briefly, then:
tmux send-keys -t azer-<opus-name> 'exit' Enter
```

**2. Re-launch with context about prior state:**

```bash
tmux send-keys -t azer-<opus-name> "cd <worktree-path> && claude --model opus --permission-mode auto \"Read $ATHANOR/AGENTS.md, then read $ATHANOR/azer.md. Your opus is at $ATHANOR/opera/<opus-file>.md. A previous session was working on this but died. Check git log and git status for progress. Resume execution.\"" Enter
```

The new session inherits all durable state (sandbox commits, staged changes, opus file notes) but starts with a clean context window.

---

## Monitoring

Once a azer is charged:

**Check for activity** — scan the crucible's tmux pane for recent output. A azer should produce visible activity (tool calls, edits, test runs) every few minutes.

**Stall detection:**
- Silence > 10 minutes with no output → likely stuck
- Repeated error messages → likely looping
- Permission prompt visible → blocked, needs approval
- Context limit message → session exhausted, needs reforging

**When stalled:**

1. **Nudge** — whisper a message asking for status or suggesting a path forward:
   ```bash
   whisper send azer-<opus-name> "Status check — are you making progress on your opus? If stuck, escalate."
   ```
2. **If nudge doesn't unstick** — escalate to the artifex with what you observed
3. **If session is dead** — reforge (see above)

---

## When the Opus is Discharged

The azer should:
1. Update the opus frontmatter to `status: discharged` with `discharged: YYYY-MM-DD`
2. Append outcome, proof of fulfillment, and reflection to the opus file
3. Optionally inscribe follow-up opera (new files in `opera/` with `status: charged`)

The marut verifies: first pull the latest specs changes (`git -C specs pull`) — the azer pushes specs on discharge but the marut's clone won't see changes until pulled. Then check: does the opus frontmatter show `status: discharged`? Do the outcome notes confirm the geas was met? Does the reflection contain useful learnings? If the azer died without discharging, the marut inscribes a recovery opus or escalates.

---

## Cleanup

The azer should exit its own session after discharge. The marut is responsible for cleaning up anything left behind.

**After confirming an opus is discharged:**

1. **Kill the crucible** if still open:
   ```bash
   tmux kill-window -t azer-<opus-name>
   ```

2. **Shut down the Docker env:**
   ```bash
   ~/code/musashi/scripts/wtp-cleanup-env.sh <worktree-path>
   ```

3. **Remove the worktree:**
   ```bash
   wtp remove <worktree-name>
   ```

Don't clean up until the opus is confirmed discharged and any follow-up opera are inscribed. Cleanup is not optional — each worktree consumes disk and Docker resources.
