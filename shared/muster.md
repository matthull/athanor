# Muster Protocol — Marut Runbook

How the marut kindles crucibles, charges azers, reforges sessions, and monitors progress.

These are core marut responsibilities — not operator tasks. The marut manages the lifecycle of all azers under its Magnum Opus.

**`$ATHANOR`** is set in your environment to this athanor's directory path. All commands below use it.

---

## Kindling a Crucible (launching an azer)

When the marut has a charged opus that needs execution:

**1. Create a sandbox (if needed)** — a sandbox isolates code changes. If the opus involves code changes to this repo, create one per the project's environment conventions (e.g. `wtp add <branch-name>`). Note the worktree path from the output.

**Skip the sandbox** if the opus is research, investigation, assessment, or work that doesn't modify this repo's code. The azer runs from the main repo directory instead. Not every opus needs code isolation — apply judgment.

**2. Muster the azer:**

```bash
ath muster <opus-file> --dir <worktree-path> --athanor <name>
```

This creates the crucible (tmux window) and launches Claude with the correct boot prompt, model, and role files. If `--dir` is omitted, uses the project directory from athanor config.

**3. Verify launch:**

```bash
ath check azer-<opus-name>
```

Should return `active` or `idle`. If it returns `dead`, the launch failed — escalate.

If `ath muster` or `ath check` behave unexpectedly, escalate. Do not improvise workarounds.

---

## Reforging (session restart)

When an azer session dies (context exhaustion, crash, unrecoverable drift), reforge it — kill the session and spawn fresh in the same crucible. The crucible endures; the session is reforged.

For marut reforging, use `ath reforge <athanor> [<mo-name>]`. For azer reforging, kill the crucible and re-muster:

```bash
ath cleanup azer-<opus-name>
ath muster <opus-file> --dir <worktree-path> --athanor <name>
```

The new session inherits all durable state (sandbox commits, staged changes, opus file notes) but starts with a clean context window.

---

## Monitoring

Once an azer is charged, check its state each loop pass:

```bash
ath check azer-<opus-name>
```

Returns one of:
- `active` — tool call in progress, working normally
- `idle` — waiting for input (may be thinking or may be stalled)
- `permission` — blocked on a permission prompt, needs approval
- `exhausted` — context limit reached, needs reforging
- `dead` — crucible not found, session died

**When `idle` persists across multiple passes** (> 10 minutes) → likely stalled. Nudge:
```bash
ath whisper send azer-<opus-name> "Status check — are you making progress on your opus? If stuck, escalate."
```

**When `permission`** → approve the prompt or escalate to the artifex.

**When `exhausted`** → reforge (see above).

**When `dead`** → the session died. Check if the opus was discharged. If not, reforge or escalate.

**If nudge doesn't unstick** → escalate to the artifex with what you observed.

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
   ath cleanup azer-<opus-name>
   ```

2. **Clean up the sandbox** per the project's environment conventions (e.g. shut down Docker env, remove worktree).

Don't clean up until the opus is confirmed discharged and any follow-up opera are inscribed.
