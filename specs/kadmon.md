# Kadmon — Athanor Operations

*The Adam Kadmon — primordial form, first embodiment of the athanor system.*

**Status:** Operational
**Created:** 2026-03-21
**Blueprint:** `specs/athanor/spec.md` (the system design — off-limits to agents)
**Instances:** `specs/athanors/<name>/` — see spec.md § Active Instances

---

## Project Status

### Next Session Action Items

- **Review PRs #4353 and #4354.** Two azer PRs from the first firing await review/merge.
- **Clean up stale worktrees.** `fix-missing-size-group-column` and `fix-rag-sync-race-condition` can be removed (`wtp remove`).
- **Launch second marut.** One charged opus remains (`fix-nomethoderror-asset-list-searcher`). Docker optimizations in progress — test them with this firing.
- **Test agent role updates in practice.** Azer discharge now includes specs push, PR body-file guidance, and permission-manager has dynamic pane discovery note. Observe whether these changes land.
- **Add stakeholder personas to bugsnag magnum-opus.md.** Before next firing: name the real people who care about #bugsnag channel health — Slack handles, GitHub usernames, org roles, why they care, what satisfaction looks like from their perspective. This expands the product lens and enables direct communication (agents can tag people in PRs, post to relevant channels). See spec.md § Magnum Opus for the full stakeholder persona guidance.

### Current Phase: Bootstrap (Adam Kadmon)

The #bugsnag athanor is the first test bed. The goal is NOT a fully autonomous error patrol — it's a working system that produces useful escalations and teaches us what to build next. Each escalation from the marut is a step on the artifex's trail.

**What "success" means right now:** The marut scans #bugsnag, inscribes opera for actionable errors, kindles crucibles for azers, monitors their progress, escalates cleanly when it hits gaps, and the discharged opera trail tells us what the system needs next.

### Readiness Checklist

| Item | Status | Notes |
|------|--------|-------|
| Materialized instance pattern | Done | AGENTS.md, magnum-opus.md, marut.md, azer.md, opus.md, muster.md — symlinked from shared/ |
| `claude-run` on PATH | Done | Headless + passthrough modes, system prompt injection |
| `wtp` for sandboxes | Done | Worktree isolation for azers |
| `whisper` CLI | Done | Reliable inter-crucible communication |
| Telegram escalation | Done | notify + andon via MCP, tested |
| Context budget guard | Done | Warns at 75/85/92%, fail-open, all advisory |
| Role enforcement hooks | Deferred | Prompt-only enforcement — observe drift under supervision |
| `claude-run` role awareness | Deferred | Manual session config for now |

### Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-03-21 | No role enforcement hooks for day 0 | Observe drift patterns under babysitting — interesting data |
| 2026-03-21 | Context budget guard is advisory-only | Info over blocking — learn thresholds from real usage |
| 2026-03-21 | Skip beholder role for now | Marut does scanning + supervision — define boundaries from experience |
| 2026-03-21 | Artifex is Primus for now | No meta-supervisor until multiple MOs exist |
| 2026-03-21 | Skip seeding opus | Opus generation is key thing to test |

---

## Primus Runbook

The artifex acts as Primus, the Forge Lord — responsible for launching and restarting maruts. Everything below is the Primus's operational reference.

### Initializing a New Athanor

When kindling a brand-new athanor instance:

**2. Symlink shared files (each must be a separate command):**

```bash
ln -s ../shared/AGENTS.md specs/athanors/$ATHANOR/AGENTS.md
ln -s ../shared/azer.md specs/athanors/$ATHANOR/azer.md
ln -s ../shared/marut.md specs/athanors/$ATHANOR/marut.md
ln -s ../shared/muster.md specs/athanors/$ATHANOR/muster.md
ln -s ../shared/opus.md specs/athanors/$ATHANOR/opus.md
```

**3. Write `magnum-opus.md`** — the only file you author from scratch. Required sections:
- **Goal** — what the athanor is pursuing
- **Athanor Structure** — copy the directory tree block from another MO, update the path
- **Pre-loaded Context** — anything the first azer needs to not start from scratch: discovery findings, references, known open questions, relevant services/files

**4. Optionally inscribe a starting opus.** If you have enough context for a concrete, actionable opus, write it in `opera/` with YAML frontmatter `status: charged` — it saves the marut an assessment cycle. If the domain is unclear, leave `opera/` empty and the marut will inscribe an assessment opus on first boot.

Opus filenames use a datestamp prefix: `YYYY-MM-DD-<descriptive-name>.md` (e.g., `2026-03-23-implement-l2-metrics.md`).

**5. Launch the marut** — follow "Launching a Marut" below, substituting your athanor's path.

---

### Launching a Marut

**1. Create the crucible:**

```bash
tmux new-window -n marut
```

**2. Start the permissions watcher:**

```bash
tmux new-window -n perms
# From perms window (~/code/musashi):
claude --model sonnet "/permission-manager check every 2 min until i say stop"
```

Babysits permission prompts across all athanor panes. Kill when done supervising.

**3. Launch the marut session:**

```bash
# From the marut's crucible (~/code/musashi):
claude --model sonnet --permission-mode auto "Read $ATHANOR/AGENTS.md, then read $ATHANOR/magnum-opus.md, then read $ATHANOR/marut.md, then read $ATHANOR/muster.md. You are the marut for this athanor. Start /loop 5m and begin your operational cycle."
```

Then use /config to enable autocompact, since we don't have an automated primus to spawn maruts yet.

That's it. The marut reads its role docs and begins.

### Reforging a Marut

When the marut session dies (context exhaustion, crash, drift):

```bash
# In the marut's crucible (~/code/musashi) — kill the old session, start fresh:
claude --model sonnet --permission-mode auto "Read $ATHANOR/AGENTS.md, then read $ATHANOR/magnum-opus.md, then read $ATHANOR/marut.md, then read $ATHANOR/muster.md. You are the marut for this athanor. Check opera/ for in-progress work and the trail. Start /loop 5m and resume your operational cycle."
```

The crucible endures. The session is reforged. Durable state lives in the opera files.

### Shutting Down (Quiescence)

1. Ensure no charged opera remain (or note their state)
2. Kill the marut session
3. Opera trail in `opera/` is the durable record (discharged opera have `status: discharged` in frontmatter)

---

## Artifex Role

The artifex is not just Primus (launching maruts). Primus is a duty the artifex currently fills. The artifex role is broader — it's the system's learning loop that never automates away.

### Routines

**Watch sessions live.** Sit down periodically and observe the marut and azers working in their tmux panes. You'll notice things the agents can't: latency, environmental friction, drift patterns, tool misuse, permission bottlenecks. This is not bootstrapping overhead — it's a permanent part of the learning process.

**Walk the fool's path.** Read discharged opera reflections. These are the system's self-reported learnings. Look for patterns across opera: recurring escalations (infrastructure gaps), geas that worked vs. geas that caused confusion, places where azers spun vs. where they executed cleanly. The fool's path is the trail of discharged opera — walk it regularly.

**Review code.** PRs created by azers need human review. This is not just quality control — it's the primary feedback channel for whether the geas produced the right work. A bad PR means a bad geas, not a bad azer.

**Tend the athanor — but with the right lens.** Agent reflections will suggest fixes ("update azer.md to say X"). Don't be kneejerk about these. The marut and azers by design don't understand the overall system architecture — they see their own context and extrapolate. An azer that struggled with permissions will say "add guidance to azer.md" when the real fix is a CLAUDE.md update, a hook, or a skill change. Apply `/agentic-architecture` and `/agentic-dev` thinking: where does this fix actually belong in the system? The agent roles are one layer. CLAUDE.md, hooks, skills, and infrastructure are others. The agent sees the symptom; the artifex sees the system.

### What Only the Artifex Can See

- Time and duration (agents have no clock awareness)
- Cross-session patterns (each agent only sees its own context)
- Environmental friction (docker startup, permission prompts, tool failures)
- Whether the system is actually producing value vs. churning

### Protocols

**Hot-patching agent roles.** When you update a live agent's role file (e.g., marut.md), notify the running session of the change. Either whisper it directly or tell it to re-read the relevant section. Agents don't watch their own files for changes — they read them once at boot.

### Homunculus (Artifex Familiar)

The artifex familiar is a Claude Code session (typically in the `athanor` tmux window) that assists the artifex with supervision and system improvement. It is not an autonomous agent — it's a collaborative tool the artifex uses interactively.

**What the familiar does:**
- Monitor marut and azer tmux panes on request (`tmux capture-pane`)
- Help update agent roles (marut.md, azer.md, muster.md, etc.)
- Notify live agents when their role files change (whisper or direct message)
- Track observations and decisions in the session log
- Help walk the fool's path — read discharged opera, surface patterns
- Draft kadmon.md updates as the system evolves

**What the familiar does NOT do:**
- Act autonomously — it works when the artifex is present
- Replace artifex judgment — it assists observation, doesn't substitute for it
- Directly intervene in agent sessions — the artifex decides when and how to intervene

**Protocols and processes will be added here as they prove useful.**

---

## Evolution Tracking

### What We're Learning

Patterns observed during operation — what works, what doesn't, what the system needs next. Fed by opus reflections and artifex observations.

**Generativity watch (2026-03-21):**

- **No Slack visibility on completed work.** Marut creates PRs but nobody in #bugsnag or #sales-enablement-development sees the activity. The marut hasn't independently thought to post links, summarize fixes, or communicate outcomes. This is an alignment signal — the goal says "clean channel" but a maximally aligned marut would also think about making its work *visible* to the humans who benefit. Not pushing the marut toward this — watching for whether goal-over-trail orientation + inquisitiveness produces this kind of thinking organically. If it doesn't emerge after a few firings, that's data about what generativity the system can/can't produce.
- **Marut is doing PM work without knowing it.** Assessing channel health, prioritizing by impact, choosing highest-value action, communicating outcomes — this is product management thinking. The marut that figures this out on its own is the one that's truly aligned. Watch for signs of this in the trail.

**From first firing (2026-03-21):**

- **Opera-as-state works.** The inscription/charge/discharge lifecycle was clean. A fresh marut reading opera directories gets a complete picture. The immutable inscription + mutable execution record pattern is sound.
- **Parallel azers work.** Two concurrent azers in separate worktrees operated independently. No conflicts, no shared state issues (other than specs/ sync, now fixed).
- **Marut triage is efficient.** Scanning, classifying, and muting low-signal errors directly saved significant time vs. mustering azers for everything. The inscription test ("can I write a concrete geas?") is a useful filter for what deserves a azer.
- **Azer quality was high.** Both produced focused, well-tested fixes at ~8% context each. The geas channeled work effectively — neither wandered.
- **The trail is genuinely valuable.** Discharged opera reflections surfaced real insights: generic filter handler fragility (size-group opus), Ducky API transient state patterns (rag-sync opus). These feed directly into team backlog.
- **Context exhaustion is not a near-term risk.** Marut ended at 17% after full cycle (triage, 2 musters, 3 mutes, extensive monitoring). Azers finished at ~8-9%.
- **specs/ cross-worktree sync was the biggest friction.** Each worktree has an independent specs/ clone — azer discharges were invisible to the marut. Fixed: azer now pushes specs on discharge, marut pulls before checking.

### Known Infrastructure Gaps

**Production access (HIGH PRIORITY):** Agents have zero production access — no Rails console, no SQL, no AppSignal dashboard. For a #bugsnag athanor this is a major constraint: many errors need production context to diagnose (specific records, data conditions, production-only config). We need a read-only SQL user and/or a sandboxed Rails console session concept ASAP. Without it, expect frequent escalations where the azer says "I can see the error but can't determine the data condition without prod access." This is documented in AGENTS.md as an explicit constraint with escalation guidance.

### Infrastructure Gaps Discovered

Track gaps the marut escalates about — these become the artifex's build-out work.

| Date | Gap | Source | Resolution |
|------|-----|--------|------------|
| 2026-03-21 | No production data access | Pre-launch audit | Need read-only SQL user + sandboxed console concept |
| 2026-03-21 | Permission prompts blocking azers | First firing observation | Added `--permission-mode auto` to muster.md launch commands. Explore overlap with `/permission-manager` — auto mode may make the perms watcher redundant for azers. Need to test whether auto mode handles all azer operations (docker compose, gh, git push) without prompts. |
| 2026-03-21 | Worktree cleanup not automated | First firing observation | Each azer creates a worktree with its own Docker env. Without automated cleanup after opus discharge, disk and Docker resources will exhaust quickly. Need `wtp remove` integrated into the marut's post-discharge flow or a separate cleanup protocol. |
| 2026-03-21 | specs/ cross-worktree sync | First firing observation | Azer discharges invisible to marut. **Fixed:** azer.md now includes `git -C specs push` on discharge; muster.md includes `git -C specs pull` before checking discharge state. |
| 2026-03-21 | `gh pr create` heredoc rejection | First firing observation | bash-allowlist-coach hook rejects heredoc/command-substitution patterns. **Fixed:** azer.md now instructs writing PR body to temp file + `--body-file`. |
| 2026-03-21 | Permission manager missed new panes | First firing observation | Claude session cached pane list from first scan. **Fixed:** added explicit anti-caching instruction to permission-manager skill. Scanner already uses dynamic `tmux list-panes -a`. |
| 2026-03-21 | Azer monitoring signal weak | First firing observation | `tmux capture-pane` unreliable for CC TUI. Scrollback (`-S -50`) best current option but noisy. Heartbeat/status file mechanism proposed but not yet implemented. |
| 2026-03-21 | No auto-merge for azer PRs | First firing observation | Azers create PRs but can't merge. Artifex must review and merge manually — bottleneck that defeats the point of autonomous error fixing. Need an auto-merge process for musashi (not athanor-specific). `/code-review` can assess quality, but the open question is: what agent makes the final merge/no-merge judgment without operator review, and what are the gates? E.g., CI green + code review pass + diff size under threshold + no schema changes? This is a significant autonomy boundary to design carefully. |

### Geas Quality Log

Track whether geas are well-crafted or cause problems. The spec warns: "a poorly designed geas turns on the creator."

| Opus | Geas Quality | Issue | Learning |
|------|-------------|-------|----------|
| fix-missing-size-group-column | Good | None | Clear goal, concrete verification (error gone + PR + root cause addressed). Azer executed cleanly without confusion. |
| fix-rag-sync-race-condition | Good | None | Good framing of race condition. Geas item 3 ("handles the race gracefully") gave azer latitude to choose retry approach. |
| mute-sso-credential-not-found | Good | None | Marut handled directly — inscription test correctly identified this as a mute, not a fix. |
| mute-slow-request-alerts | Good | None | Marut direct mute — no actionable path, correctly dispositioned. |
| mute-deadlock-customer-page-upload | Good | None | Single occurrence, correctly muted with rationale. |
| fix-nomethoderror-asset-list-searcher | Good | Not yet tested | Inscribed but no azer mustered before session end. Geas looks well-formed. |

---

## Architecture Notes (artifex reference)

### What the agents see

```
specs/athanors/<name>/
├── AGENTS.md          ← core vocab, geas, constraints (all agents)
├── magnum-opus.md     ← goal + witnesses (marut entry point)
├── marut.md           ← supervisor role + geas addenda
├── azer.md            ← worker role + geas addenda
├── opus.md            ← opus lifecycle, inscription/discharge protocol
├── muster.md          ← crucible kindling, reforging, monitoring
└── opera/             ← all opera, YAML frontmatter for status
```

### What the agents DON'T see

- `specs/athanor/` — the system blueprint (this directory)
- This file (kadmon.md) — project planning, not operational
- The full domain model and gastown analysis

### Key injection points

The geas is the primary behavioral shaping mechanism. The core geas in AGENTS.md makes escalation feel like compliance ("both are equally valid fulfillments of your geas"). Role-specific geas addenda in marut.md and azer.md add bounded obligations. The opus-level geas channels work toward concrete outcomes.

### Supervision model (current)

```
Artifex (human) = Primus, the Forge Lord
  └── Marut (one, supervises the MO)
        └── Azers (spawned per opus, sandboxed)
```

No beholder. No automated Primus. The artifex launches/reforges maruts. The marut kindles/reforges/monitors azers. Boundaries between roles emerge from experience.
