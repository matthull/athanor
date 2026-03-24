# Context Budget Guard Specification

**Priority:** HIGH — Immediate Value
**Status:** Draft (skeleton)
**Created:** 2026-03-20
**Source:** Gastown architecture analysis session (2026-03-20)

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[R:gastown]` | Derived from gastown codebase analysis |
| `[B:bottleneck]` | Addresses the critical orchestrator context bottleneck |
| `[D:reason]` | Design decision with rationale |
| `[E:existing]` | Extends existing musashi infrastructure |
| `[U:topic]` | Unbacked — needs validation |

---

## Overview

Add a context budget guard to musashi's agentic infrastructure. When a Claude Code
session approaches context exhaustion, the guard detects it and injects escalating
warnings — culminating in a hard gate that blocks further work and forces graceful
shutdown or handoff.

This directly addresses the **critical bottleneck** identified in `open-work-items.md`:
long-running `/orchestrate` sessions fill context, drift, and become unreliable.
The guard ensures agents either hand off cleanly or stop before producing degraded output.
`[B:bottleneck]`

## Why This Is High Priority / Immediate Value

- **Zero new infrastructure required.** Claude Code already has `UserPromptSubmit` hooks
  and JSONL session transcripts with token usage. This is a shell script + config. `[E:existing]`
- **Solves the #1 reliability problem today.** Context exhaustion is the single biggest
  source of degraded agent output in musashi. `[B:bottleneck]`
- **Proven pattern.** Gastown's `context-budget-guard.sh` is battle-tested. We're adapting,
  not inventing. `[R:gastown]`
- **Composable.** Works independently of `/orchestrate` — improves ALL Claude Code sessions
  (interactive, autonomous, orchestrated). `[D:composability]`

## Domain Dictionary

- **Context budget:** The percentage of the model's context window currently consumed
  by the session transcript (input + cache tokens vs max window).
- **Guard:** A `UserPromptSubmit` hook script that runs before each turn, reads token
  usage from the session transcript, and injects warnings or blocks execution.
- **Hard gate:** Guard exits non-zero, blocking the user turn entirely. The agent
  cannot proceed without addressing context exhaustion.
- **Soft gate:** Guard exits 0 but injects a strong recommendation message.
- **Handoff:** Ending the current session and spawning a fresh one, with work state
  preserved in durable stores (git, task-plan.md, specs/).

## Workflow / Dataflow

```
User/Agent submits prompt
  → UserPromptSubmit hook fires
    → context-budget-guard.sh runs
      → Reads ~/.claude/projects/<cwd>/sessions/<id>.jsonl
      → Extracts token usage from last assistant message
      → Calculates % of max context window
      → Below WARN threshold? → exit 0 (silent)
      → Above WARN? → inject warning message, exit 0
      → Above SOFT_GATE? → inject strong recommendation, exit 0
      → Above HARD_GATE? → inject block message, exit non-zero (blocks turn)
```

## Gastown Reference Implementation

Gastown's `scripts/guards/context-budget-guard.sh` provides the proven pattern: `[R:gastown]`

**Key design choices to preserve:**
- **Three-tier thresholds:** warn (75%), soft gate (85%), hard gate (92%)
- **Role-based differentiation:** infrastructure roles hard-gate, worker roles warn-only
- **Fail-open on errors:** If transcript parsing fails, jq missing, etc. → allow the turn
- **Environment variable configuration:** All thresholds overridable via env vars
- **Token calculation:** Sum of `input_tokens` + `cache_creation_input_tokens` + `cache_read_input_tokens`
  from the last assistant message's `usage` object in the JSONL transcript

## Adaptation for Musashi

### What transfers directly `[R:gastown]`
- Guard script structure and threshold logic
- JSONL transcript reading approach
- Fail-open error handling
- Environment variable configuration pattern

### What needs adaptation `[D:musashi-context]`

| Gastown | Musashi | Notes |
|---------|---------|-------|
| Roles: mayor, deacon, crew, polecat | Roles: TBD — interactive, autonomous, orchestrator, implementer | Map to musashi's session types |
| `gt handoff` (Go CLI, bead state) | `/direct-handoff` skill or manual session cycling | Musashi doesn't have `gt handoff` yet |
| Daemon-driven compaction (CompactorDog) | Claude Code's built-in auto-compaction | Musashi relies on CC's native compaction |
| Bead checkpoint recovery | task-plan.md phase results + git state | Musashi's persistence is plan-doc-based |

### Role Mapping (Needs Refinement) `[U:roles]`

| Musashi Session Type | Gating Behavior | Rationale |
|---------------------|-----------------|-----------|
| Interactive (user at keyboard) | Warn only | User can decide what to do |
| `/orchestrate` orchestrator | Hard gate | Most valuable to protect — context drift here cascades |
| TeamCreate implementer | Warn → soft gate | Work is scoped; can push partial progress |
| `/autonomous` | Hard gate | No human to intervene if quality degrades |
| Subagent (Agent tool) | N/A — CC manages these | Subagents have their own context windows |

### Handoff Mechanism `[D:incremental]`

**MVP:** The guard injects a message telling the agent to wrap up. For interactive sessions,
this is sufficient — the user sees the warning. For autonomous/orchestrated sessions, the
agent should:

1. Save current progress (commit, update task-plan.md phase results)
2. Log the handoff to the session log
3. Exit cleanly

**Future:** Integrate with `/direct-handoff` for automatic session chaining.

## Configuration `[R:gastown]` `[D:musashi-defaults]`

```bash
# Thresholds (percentage of max context)
GT_CONTEXT_BUDGET_WARN=0.75
GT_CONTEXT_BUDGET_SOFT_GATE=0.85
GT_CONTEXT_BUDGET_HARD_GATE=0.92
GT_CONTEXT_BUDGET_MAX_TOKENS=200000

# Role-based gating
GT_CONTEXT_BUDGET_HARD_GATE_ROLES=orchestrator,autonomous
GT_CONTEXT_BUDGET_WARN_ONLY_ROLES=interactive,implementer

# Escape hatch
GT_CONTEXT_BUDGET_DISABLE=0
```

**Open question:** How does the guard know the current session's role? Options: `[U:role-detection]`
1. Environment variable set by session startup (e.g., `/orchestrate` sets `GT_ROLE=orchestrator`)
2. Read from session injection file (`/tmp/claude-sessions/{id}/role.txt`)
3. Default to warn-only, let specific workflows opt into hard gating

## Transcript Location `[E:existing]`

Claude Code stores session transcripts at:
```
~/.claude/projects/<cwd-with-dashes>/sessions/<session-id>.jsonl
```

The guard needs to find the current session's transcript. Approach TBD — gastown
reads the most recent `.jsonl` in the directory. `[U:transcript-discovery]`

## Test Scenarios

| Scenario | Expected Behavior |
|----------|-------------------|
| Session at 50% context | Guard silent, exit 0 |
| Session crosses 75% | Warning injected, exit 0 |
| Session crosses 85% | Strong recommendation injected, exit 0 |
| Orchestrator session crosses 92% | Hard gate, exit non-zero, turn blocked |
| Interactive session crosses 92% | Warning only, exit 0 |
| jq not installed | Fail open, exit 0 |
| Transcript file not found | Fail open, exit 0 |
| Malformed JSONL | Fail open, exit 0 |
| `GT_CONTEXT_BUDGET_DISABLE=1` | Guard skipped entirely |

## Out of Scope

- **Automatic session chaining** — that's a separate, larger effort (the `/direct-handoff` evolution)
- **Daemon-driven compaction** — gastown's CompactorDog pattern; musashi uses CC's built-in compaction
- **Checkpoint/crash recovery** — gastown's bead checkpointing; musashi's plan-doc persistence is sufficient for now
- **Subagent context monitoring** — CC manages Agent tool subagents internally

## Implementation Notes

This is a shell script + Claude Code settings.json hook entry. Estimated scope:
one guard script, one settings.json update, one test pass. Could ship in a single
focused session.

## Product FYI

N/A — internal infrastructure, no product team impact.

## Retrospective

When complete, review:

### Documentation Updates
- [ ] Add guard to `system-map.md` in village agentic domain
- [ ] Update `open-work-items.md` to mark this resolved
- [ ] Document the role-detection approach chosen

### Workflow Improvements
- [ ] Does `/orchestrate` need to set a role env var at startup?
- [ ] Should `/autonomous` set one too?
- [ ] Consider a `/claude-code-config` update to make hook installation a skill

### Knowledge Capture
- [ ] JSONL transcript parsing patterns (reusable for other guards)
- [ ] Threshold tuning learnings after real-world usage
