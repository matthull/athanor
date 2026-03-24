# Orchestrate Logging & Capture — Design Input

**Status:** Research complete, design pending
**Created:** 2026-03-12
**Source:** Pipeline→Orchestrate migration audit session (2026-03-11)

---

## Problem Statement

`/orchestrate` has no logging or capture mechanism. `/pipeline` had `specs/{project}/pipeline-log.md` and a REFLECT phase. The gap means learnings from orchestrated runs are lost — the same mistakes repeat, and there's no data to improve the system.

---

## What /pipeline Did

### Pipeline Log (`specs/{project}/pipeline-log.md`)

Per-run structured entry appended after SHIP:

```markdown
## Pipeline Run: {date} — {task summary}

**Complexity:** {assessed complexity}
**Phases executed:** {list}
**Outcome:** {PR URL or status}

### What worked
- {1-3 bullet points}

### What didn't
- {1-3 bullet points, or "Nothing notable"}

### Process changes suggested
- {Specific, actionable suggestions for pipeline improvement, or "None"}
```

### Learning Propagation Table

When a learning is captured, push it to all downstream artifacts an implementer would encounter:

| Learning type | Push to |
|---------------|---------|
| Architectural constraint | `architecture.md`, relevant handoff templates, relevant skills, memory |
| New convention | `architecture.md`, handoff templates, memory, domain plugin |
| Workflow pattern | Handoff templates, pipeline skill, domain plugin |
| Tool behavior | Handoff templates, memory, MCP strategy doc |

Key insight: "If the learning only lands in `architecture.md` or `MEMORY.md`, an agent following a handoff template will never see it. Propagate to the point of use."

### How It Ran

- Orchestrator did REFLECT directly (no subagent) — lightweight
- Ran immediately after SHIP, before CI results (CI monitor was async)
- Created the log file if it didn't exist
- Small fixes to pipeline skill applied immediately; larger changes noted for discussion

---

## What /orchestrate Needs (Different From /pipeline)

### Capture vs. Reflect — Separate Concerns

**Capture** = recording what happened during a run (continuous, lightweight, non-blocking).
**Reflect** = processing learnings and propagating them to downstream artifacts (heavy, ad-hoc, can happen later).

/pipeline conflated these into a single REFLECT phase at the end. The problems:
1. REFLECT only fires if the sequence completes cleanly — crashes, context loss, and abandoned sessions lose all learnings
2. Forcing reflection at the end of every run adds overhead when there's nothing notable to reflect on
3. The valuable capture (what happened, what went wrong) should happen throughout the run, not reconstructed from memory at the end

**Design principle:** Capture continuously. Reflect on demand.

### What Matters Most to Capture

The highest-signal learnings are about **interrupt quality** — the calibration of when the agent stops to ask vs. proceeds autonomously:

#### False Positive Interrupts (FP)
Agent interrupted the operator when it didn't need to. Examples:
- Asked for approval on a decision that followed established patterns
- Gated a phase transition when the profile said not to
- Asked a question it could have answered from available context
- Requested confirmation on a routine action

**Why this matters:** Each unnecessary interrupt erodes trust and slows the workflow. Accumulated FPs push operators toward `autonomous` profile even when `review-gated` would be safer.

#### False Negative Interrupts (FN)
Agent proceeded when it should have asked. Examples:
- Made an architectural decision it wasn't equipped for
- Chose between trade-offs without surfacing options
- Proceeded past a novel pattern without flagging it
- Went ahead with something that needed human judgment
- Worked around a blocker instead of escalating

**Why this matters:** Each FN is a silent quality regression. The operator discovers the problem later (in review, in production, or when the next run hits the same issue). FNs are harder to detect than FPs — the operator has to notice what *didn't* happen.

### Other Capture Categories

Beyond FP/FN interrupts:

- **Phase transitions** — what started, what completed, duration
- **Gate decisions** — what was reviewed, approved/rejected, rationale
- **Agent summaries** — what implementer reported, what verifier found (structured, not full output)
- **Blockers and resolutions** — what broke, how it was fixed, root cause
- **Verification verdicts** — spec-check result, QA plan classification, ranger results
- **Tool/environment issues** — MCP failures, permission prompt friction, docker problems

### Capture Format (Proposed)

Each entry is lightweight and structured:

```markdown
### [{CATEGORY}] {description}
- **Phase:** {current phase}
- **Agent:** {which agent, if applicable}
- **Detail:** {1-3 sentences}
- **Signal:** {FP | FN | BLOCKER | DECISION | INFO}
```

For FP/FN entries specifically:

```markdown
### [FP] Asked operator to approve task list when profile was review-gated (within-phase-gates: off)
- **Phase:** DESIGN
- **Agent:** orchestrator
- **Detail:** Profile had within-phase-gates: off but orchestrator still gated task list approval. Operator had to dismiss.
- **Signal:** FP
- **Autonomy context:** profile=review-gated, trait=within-phase-gates:off
- **Fix:** Check profile traits before evaluating gates (bug in gate evaluation logic)

### [FN] Implementer chose REST over GraphQL without surfacing trade-off
- **Phase:** IMPLEMENT
- **Agent:** implementer-1
- **Detail:** Handoff said "API endpoint for X" without specifying protocol. Implementer chose REST. Codebase uses GraphQL for all other endpoints. Should have escalated to task lead.
- **Signal:** FN
- **Fix:** Implementer prompt should include "flag decisions where codebase conventions exist but handoff is ambiguous"
```

### Where Capture Goes

**Two destinations (different granularity):**

1. **`<task-dir>/run-log.md`** — per-task, detailed, all categories. Lives alongside triage.md and task-plan.md. This is the full record of what happened during this specific orchestration run.

2. **Egregore session log** — cross-project, summary entries only. Phase transitions and notable events (FP, FN, blockers). Uses existing `mcp__egregore-mcp__session_log_append` infrastructure.

The per-task log is the primary artifact. Egregore entries are for cross-session visibility (e.g., pulse scans that look for patterns across runs).

### When Capture Happens

- **Phase transitions:** Log when entering/completing each phase
- **Gate decisions:** Log immediately when a gate fires (or is skipped)
- **Agent completions:** Log when an agent reports back with summary
- **Blockers:** Log immediately when escalated
- **FP/FN signals:** Log when detected (operator feedback, or orchestrator self-assessment)

The orchestrator should capture as a background habit — not announce it, not block on it. Same philosophy as egregore session logging.

---

## Reflecting on Captured Learnings (Separate Concern)

Reflection processes captured learnings and propagates them. This can happen:
- Ad-hoc: operator triggers `/reflect` or similar
- Periodic: pulse scan checks for unprocessed learnings
- Next session start: check for CAPTURED learnings from prior sessions

The /pipeline propagation table is still valid — learnings must reach point of use. But the trigger doesn't need to be baked into the sequence as a mandatory phase.

### Processing Status

Every learning carries a status:

```
- **Status:** CAPTURED | PROCESSED
- **Propagated to:** [list of artifacts updated, filled when PROCESSED]
```

This makes unprocessed learnings queryable. A future `/reflect` skill or pulse check can scan for `Status: CAPTURED` entries and process them.

---

## Relationship to Existing Skills

- **`/fp` and `/fn` skills already exist** — they flag false positive and false negative interrupts. The capture system should integrate with these (or at minimum use the same taxonomy).
- **Egregore session logging** — already captures session-level events. Orchestrate logging adds structured phase/agent/interrupt data.
- **Task dir convention** — `<task-dir>/run-log.md` fits naturally alongside `triage.md` and `task-plan.md`.

---

## Open Questions

1. **Should the orchestrator self-assess FP/FN?** Or only capture them when the operator flags them (via `/fp`, `/fn`)? Self-assessment risks false confidence; operator-only risks missing signals when operator doesn't notice.

2. **How verbose should the per-task log be?** Full agent output summaries bloat the log. But too terse loses the detail needed for reflection. The format above (1-3 sentences per entry) may be the right balance.

3. **Should capture be in the orchestrate SKILL.md or in the feature-implementation sequence?** If in SKILL.md, it applies to all modes/sequences. If in the sequence, it's specific to feature implementation. Probably SKILL.md — capture is a universal orchestrator concern.

4. **Integration with `/reflect` (future skill):** When `/reflect` is built, it should be able to scan `run-log.md` files for unprocessed learnings and drive propagation. The capture format needs to support this (hence the Status field).
