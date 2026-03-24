# Handoff Evolution: From Blueprint to Goal-Based Sessions

**Date:** 2026-03-15
**Context:** Audit of existing handoffs + Gastown session boundary analysis reveals we may be at a local maximum with our handoff system. The system was designed before automated quality gates existed.

---

## Current State (Audited)

### Handoff Specificity Spectrum (from 5 handoff audit)

| Handoff | Type | Prescriptiveness | Latitude | When Used |
|---|---|---|---|---|
| T002 Route Relocation | Blueprint | 95% | Very Low | Structural refactor |
| T004 HMAC Signing | Blueprint | 90% | Low | Security-critical |
| T005 Mailer | Goal+Constraints | 65% | Medium | Feature logic |
| T001 Migration+Model | Specification | 60% | Medium-High | Data model |
| SUP-260 Allowlist | Design+Goal | 40% | High | Architecture |

**Finding:** We already use a risk-adjusted spectrum. Blueprints for security/structural, goals for features. This is appropriate and should continue.

### Template Composition (consistent across all projects)

- 25-35% actual task specification + acceptance criteria
- 65-75% template boilerplate (core mandates, context loading, verification loops)
- Template is identical across all handoffs (sourced from core-task-template.md)

### SAL-89 vs Proof-Recommender: Identical Structure

No variation in methodology between projects. Same template, same reference style, same verification loops. Only domain differs.

---

## The Local Maximum

The handoff system solved a real problem (agents going off rails with too much freedom) and has worked well. But three developments may have shifted the optimal point:

### 1. `/code-review` refactor loop now exists

**Before:** Quality = f(handoff specificity). No automated quality gate → detailed handoffs were the only defense against bad code.

**After:** Quality = f(quality gates). `/code-review` runs 10 specialist legs in a loop until P0/P1 clean. The refactor step catches issues regardless of how prescriptive the input was.

**Implication:** Handoff specificity matters less for quality. We can give more latitude without sacrificing output quality.

### 2. Template boilerplate taxes every session

65-75% of every handoff is the same core mandates, context loading instructions, and verification loops. Each fresh implementor session reads 300 lines to get 80 lines of task-specific content. With 4 work units = 4 sessions = the implementor reads the same core mandates 4 times.

**Implication:** Mandates should be injected via session context (hooks, CLAUDE.md), not embedded in every handoff. The handoff should be mostly task-specific content.

### 3. One work unit per session throws away accumulated context

Each implementor session starts cold — re-reads the codebase, re-discovers patterns, re-loads domain knowledge. Everything learned in session 1 is gone by session 2. This is expensive and slow.

**Implication:** Longer-lived implementor sessions that handle multiple work units would accumulate useful codebase context. The goldilocks zone is somewhere between "one handoff per session" and "one session until crash."

---

## Proposed Direction

### Multi-Unit Implementor Sessions

Instead of one handoff → one session → one work unit:

```
Task lead gives implementor:
  - List of work units (goals + acceptance criteria, NOT full handoffs)
  - Pattern references ("follow ApiClient for this, SeismicSettings for that")
  - Quality gates (tests pass, /code-review clean)
  - Context budget ("work until ~60% context consumed, then checkpoint")

Implementor session:
  1. Read work unit goals
  2. Implement work unit 1
  3. Run tests → /code-review loop → fix until clean
  4. Checkpoint (update task plan with [x] for completed unit)
  5. Implement work unit 2 (WITH accumulated context from unit 1)
  6. Run tests → /code-review loop → fix until clean
  7. Checkpoint
  8. Continue until: work complete OR context budget reached
  9. Report back to task lead with: completed units, remaining work, key context
```

### Context Budget vs Fixed Handoffs

The session boundary becomes dynamic:
- **Stop when:** context budget reached (~60% of window), OR all work units complete
- **Checkpoint:** finish current TDD cycle, update task plan, report to task lead
- **Resume:** new session reads task plan, sees completed units ([x]), continues from next

This is how Gastown works — agents do as many molecule steps as context permits, then handoff. The plan (our task-plan.md) tracks what's done.

### What The Handoff Template Becomes

Handoffs shrink dramatically. Core mandates move to session injection (already have infrastructure). Template becomes:

```markdown
## Work Unit: {name}

**Goal:** {one sentence}
**Acceptance criteria:**
- {testable criterion 1}
- {testable criterion 2}

**Pattern references:**
- Follow {existing_file} for {what aspect}

**Verification:** Tests pass, /code-review clean (P0/P1 = 0)
```

That's it. 20-30 lines instead of 300. The core mandates, TDD discipline, escalation protocol — all injected via session context, not per-handoff.

---

## Audit Required: Core Handoff Template

Before restructuring, audit the actual template to identify which restrictions are still needed vs superseded:

**Candidates for removal/relocation:**
- Core mandates → already in CLAUDE.md and session injection
- TDD discipline → now in task-lead-implementer-mode.md refactor step
- Testing discipline → already in CLAUDE.md
- Verification loops → standardized, could be a skill not template content
- Skills to load → skill discovery handles this
- Context loading instructions → session injection handles this

**Candidates to keep in handoff:**
- Task-specific goal and acceptance criteria (unique per handoff)
- Pattern references (unique per handoff)
- Spec context excerpts with line numbers (unique per handoff)
- QA expectations when non-standard (unique per handoff)

---

## Template Constraint Audit

Reviewed `~/.claude/skills/handoff/templates/core.md` (257 lines). Lines 1-142 are core mandates preamble. Lines 147-257 are task-specific template.

| Section | Lines | Verdict | Why |
|---|---|---|---|
| I. Protocol Availability | 15-18 | REMOVE | Just says "protocols are system-wide" — adds nothing |
| II. Anti-Rationalization | 22-25 | DEMOTE to 1 line | Cargo cult — no evidence "don't rationalize" helps LLMs |
| III. Evidence-Based | 27-48 | RELOCATE | Already in CLAUDE.md "Show Your Work". Inject via session, not per-handoff |
| IV. Testing Discipline | 51-64 | RELOCATE | Already in CLAUDE.md + now enforced by /code-review refactor loop |
| V. Verification Principle | 67-76 | RELOCATE | Valid but belongs in session injection, not handoff template |
| VI. Collaboration Gates | 79-87 | REWRITE | "Never decide without user" wrong for autonomous agents. Task lead = authority. |
| VII. Completion Standards | 90-97 | RELOCATE | Already in CLAUDE.md core mandates |
| VIII. Documentation Clarity | 99-107 | KEEP | Right principle — handoffs describe WHAT not HOW |
| IX. Tool Efficiency | 110-116 | REMOVE | /skill-discovery handles this |
| A. STOP and Ask | 122-133 | REMOVE | Rewritten as escalation protocol in core mandates |
| B. Linguistic Reinforcement | 135-142 | REMOVE | Cargo cult prompting. No evidence it helps. |

**Result:** ~110 of 142 preamble lines are redundant with CLAUDE.md, session injection, or superseded. The unique handoff value is in the task-specific template (lines 147-257).

### What the Slim Handoff Looks Like

After relocating redundant mandates to session injection:

```markdown
# Work Unit: {name}

**Goal:** {one sentence}
**Type:** {feature | refactor | bugfix | security}

## Context
- Spec: {path with line refs}
- Pattern references: {existing code to follow}
- Skills: {auto-loaded via /skill-discovery}

## Acceptance Criteria
- [ ] {testable criterion 1}
- [ ] {testable criterion 2}

## Verification
- Targeted: {isolated spec command}
- Integration: {related component spec}
- E2E: {manual validation if needed}

## Quality Gate
Tests pass + /code-review loop clean (P0/P1 = 0)
```

~30-40 lines vs ~300 lines. Core mandates come from session context, not the handoff.

---

## What To Experiment With (Priority)

### Experiment 1: Multi-Unit Implementor Session

**On the next medium-risk feature:**
1. Task lead writes 3-4 work units as slim handoffs (30 lines each, not 300)
2. Give all work units to one implementor session
3. Implementor works through units sequentially, running /code-review refactor loop after each
4. Context budget: "work until ~60% context consumed, then finish current TDD cycle and checkpoint"
5. Compare: time to completion, code quality, number of /code-review iterations vs our current approach

**Measure:**
- Did accumulated context help? (Fewer re-reads, better pattern recognition?)
- Did quality suffer? (More P0/P1 issues in later work units?)
- Did the implementor stay on track? (No off-course behavior?)
- Was the slim handoff sufficient? (Or did it need more detail?)

### Experiment 2: Slim Handoff Template

**In parallel with experiment 1:**
1. Remove the 142-line preamble from handoff template
2. Ensure core mandates are injected via session injection (already have infrastructure)
3. Write work units using the slim format above
4. Compare: do implementors follow the mandates equally well from injection vs embedded?

---

## Risk Assessment

### What could go wrong with looser handoffs + longer sessions:
1. **Implementor goes off course** → Mitigated by /code-review refactor loop catching bad patterns
2. **Accumulated context causes drift** → Mitigated by context budget + checkpoint at TDD cycle boundary
3. **No blueprint for security-critical work** → Keep blueprint tier for security/structural (risk-adjusted spectrum continues)
4. **Core mandates not followed** → Mitigated by session injection (hooks, CLAUDE.md) — doesn't depend on handoff

### What we gain:
1. **Faster iteration** — no cold start per work unit
2. **Better pattern recognition** — implementor learns codebase across units
3. **Less orchestrator overhead** — task lead writes goals not blueprints
4. **Less template tax** — 20 lines per unit vs 300 lines per handoff
5. **More natural TDD** — refactor step has context from previous work
