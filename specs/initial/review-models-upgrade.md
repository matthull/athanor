# Review Models: Gastown vs Ours + Upgrade Plan

**Date:** 2026-03-14
**Context:** Side-by-side comparison of Gastown's review formulas against our review skills. Goal: identify concrete upgrades to our review processes by adapting Gastown's structured, multi-dimensional review patterns.

---

## Current State: Our Review Skills

### What We Have

| Skill | Reviews | Dimensions | Structure |
|---|---|---|---|
| `/code-review` | Code changes | 6 parallel agents: architecture, quality, security, performance, testing, docs | Convoy-like (parallel `code-review-expert` agents). Best structured of our review skills. |
| `/task-review` | Uncommitted code vs handoff | Handoff compliance: critical directives, quality checklist, anti-patterns, security, testing | Single agent, checklist-driven. Focused on handoff conformance, not general code quality. |
| `/spec-check` | Implementation vs spec | Per-requirement PASS/FAIL with evidence | Single agent, very focused. Cleanly scoped — only spec compliance, explicitly excludes code quality. |
| `/review` | Anything (delegates to Gemini) | Whatever Gemini decides | Thin wrapper around Gemini CLI. No structured dimensions. |
| `/qa-plan` | Feature from docs/code | QA classification: [TEST], [BROWSER], [CONSOLE], [HUMAN] | Single agent, classifies verification items by execution channel. |

### What We DON'T Have

1. **PRD/spec review** — No structured review of specs before implementation. We write specs and go straight to implementation. No convoy of analysts questioning requirements, gaps, ambiguity, feasibility, scope.
2. **Plan review** — No structured review of implementation plans before execution. We write task plans and execute them. No reviewers checking completeness, sequencing, risk, scope-creep, testability.
3. **Design review** — No structured design exploration from multiple dimensions. We do design in conversation, ad-hoc.
4. **Presets / quality tiers** — `/code-review` always launches the same agents. No "gate" preset (fast, blockers only) vs "full" preset (comprehensive).

---

## Gastown's Review Models (Verbatim Analysis)

### 1. Code Review Convoy (10 legs)

**Legs split into two categories:**

Analysis legs (read and analyze):
- **correctness** — Logic errors, bugs, edge cases, race conditions, off-by-one
- **performance** — O(n^2) where O(n) possible, N+1 queries, memory leaks, blocking in async
- **security** — OWASP Top 10, injection, auth bypass, SSRF, hardcoded secrets
- **elegance** — Unclear abstractions, coupling, SOLID violations, reinventing utilities
- **resilience** — Swallowed errors, missing retry/backoff, resource cleanup, circuit breakers
- **style** — Convention compliance, comment quality, log levels, import organization
- **smells** — Long methods (>50 lines), deep nesting (>3 levels), shotgun surgery, god classes

Verification legs (check implementation quality):
- **wiring** — Dependencies added but not used, SDK added but old implementation remains, config defined but never loaded
- **commit-discipline** — Giant WIP commits, poor messages, unatomic changes
- **test-quality** — Weak assertions (only checking != nil), missing negative cases, tests that can't fail, flaky indicators

**Key innovations over our `/code-review`:**
1. **"Wiring" leg** — Catches "installed-but-not-wired" gaps. This is genuinely novel. We don't check for this at all.
2. **"Resilience" leg** — Separate from correctness and security. Focused on failure modes, partial failures, recovery paths.
3. **"Smells" leg** — Explicit code smell detection with specific thresholds (>50 lines, >3 nesting levels).
4. **Presets** — `gate` (4 legs, fast), `full` (10 legs), `security-focused` (4 legs), `refactor` (4 legs). Our `/code-review` is always the same.
5. **Structured output template** — Every leg produces: Summary → Critical (P0) → Major (P1) → Minor (P2) → Observations. With file:line references.
6. **Synthesis step** — Deduplicates across legs, notes which legs found same issue (higher confidence), prioritizes by impact and effort.

### 2. PRD Review Convoy (6 legs)

- **requirements** — Are success criteria defined? Can someone write a test from this PRD?
- **gaps** — What's completely absent? Auth, multi-tenancy, data migration, backwards compat, admin tooling, deprecation
- **ambiguity** — Vague language ("fast", "reasonable"), undefined terms, contradictions, "should" vs "must"
- **feasibility** — Features assuming capabilities the system doesn't have, hard technical problems, missing prerequisites
- **scope** — Missing out-of-scope statements, inevitable feature requests, MVP definition, phasing
- **stakeholders** — Unstated users, ops/admin perspective, conflicting user needs, launch coordination

**Synthesis produces:**
- Executive summary (2-3 sentences on PRD health)
- Critical questions (must answer before build)
- Important but non-blocking
- Confidence assessment table (requirements completeness, feasibility, scope clarity, ambiguity level, overall readiness — each H/M/L)
- Next steps checklist

**Key innovation:** The output is QUESTIONS FOR THE HUMAN, not assertions. The review's job is to surface what's missing/unclear, not fill in gaps. This is exactly right — the review should generate questions that force the spec author to think harder.

### 3. Plan Review Convoy (5 legs)

- **completeness** — Requirements with no corresponding plan step, missing infra/setup/monitoring/rollback steps
- **sequencing** — Steps depending on things not yet built, missing blocking relationships, parallelization opportunities
- **risk** — High uncertainty steps, external dependencies, hard-to-reverse changes, "we'll figure it out" language
- **scope-creep** — Steps not necessary for MVP, gold-plating, bundled refactors, premature abstraction
- **testability** — Steps with no test plan, features hard to test in isolation, missing smoke tests post-deploy

**Per-leg output:** PASS / PASS WITH NOTES / FAIL verdict + must-fix + should-fix + observations

**Synthesis produces:**
- Overall GO / GO WITH FIXES / NO-GO recommendation
- Leg verdicts table
- Must-fix items (blocking)
- Should-fix items
- Next steps

**Key innovation:** GO/NO-GO gate. The review produces a binary recommendation, not just a list of issues. This forces the reviewer to commit to a judgment.

### 4. Design Convoy (6 legs)

- **api** — Interface design, CLI ergonomics, naming, consistency with existing patterns
- **data** — Data model, storage format, schema design, migration strategy, data lifecycle
- **ux** — Mental model, workflow integration, learning curve, error experience, discoverability
- **scale** — Scale dimensions, bottlenecks, complexity, caching opportunities, degradation modes
- **security** — Trust boundaries, attack surface, threat model, sensitive data, defense-in-depth
- **integration** — Existing components touched, dependencies, backwards compatibility, feature flags

**Per-leg output:** Summary → Key considerations → Options explored (with pros/cons/effort for each) → Recommendation → Constraints → Open questions → Integration points

**Synthesis produces:**
- Executive summary
- Proposed design (overview, key components, interface, data model)
- Trade-offs and decisions (decisions made, open questions, trade-offs)
- Risks and mitigations
- Implementation plan (Phase 1 MVP, Phase 2 Polish, Phase 3 Future)

---

## Gap Analysis: What We Should Steal

### Priority 1: PRD/Spec Review (We Have Nothing)

This is the biggest gap. We go from spec to implementation without structured review. A `/spec-review` skill using the convoy pattern would catch issues before they become expensive to fix.

**Proposed dimensions (adapted from Gastown for our context):**

| Leg | Focus | Why It Matters For Us |
|---|---|---|
| **requirements** | Are completion criteria testable? Can someone write a test from this spec? | Our specs sometimes have vague criteria that lead to ambiguous implementations |
| **gaps** | What's missing? Auth, data migration, backwards compat, admin tooling | We've been burned by missing migration plans, auth gaps |
| **ambiguity** | Vague language, undefined terms, "should" vs "must" | Our CLAUDE.md warns about this but we don't systematically check |
| **feasibility** | Hard technical problems, missing prerequisites, external dependencies | Catches issues before commit to implementation |
| **scope** | What's in/out? MVP definition? Phasing? | Scope creep is a real problem — explicit boundary review helps |

Drop `stakeholders` leg — we're a small team, stakeholder analysis is overkill.

**Output:** Questions for the spec author (human), confidence assessment table, GO/NO-GO recommendation.

### Priority 2: Plan Review (We Have Nothing)

After `/orchestrate` generates a task plan, it goes straight to execution. A `/plan-review` step would catch:
- Missing test plans
- Bad sequencing (dependency errors)
- Scope creep in the plan
- Risk not addressed early

**Proposed dimensions (adapted for our context):**

| Leg | Focus | Why It Matters For Us |
|---|---|---|
| **completeness** | All spec requirements covered? Missing setup/migration/test steps? | We've had plans that forgot migration steps |
| **sequencing** | Dependencies correct? Schema before code? Feature flags before feature? | Our migration discipline (CLAUDE.md) requires careful sequencing |
| **risk** | High uncertainty? External dependencies? Hard-to-reverse changes? | Catches "we'll figure it out" before it bites |
| **scope-creep** | Steps not in spec? Gold-plating? Bundled refactors? | Keeps implementation focused |
| **testability** | Can we verify each step worked? What does QA look like? | Feeds into our `/qa-plan` step |

**Output:** GO / GO WITH FIXES / NO-GO + per-leg verdicts + must-fix items.

### Priority 3: Code Review Upgrades (We Have Something, But Gaps)

Our `/code-review` already does parallel agents. Upgrades:

1. **Add "wiring" leg** — "installed-but-not-wired" detection. This catches a real class of bugs we don't look for: dependency added to Gemfile/package.json but old implementation still used.

2. **Add "resilience" leg** — Separate from security. Focused on error handling, partial failures, recovery paths. Our codebase has patterns where errors are swallowed or recovery is missing.

3. **Add presets** — `gate` (fast: wiring, security, smells, test-quality), `full` (all), `security` (security, resilience, correctness, wiring), `refactor` (elegance, smells, style, commit-discipline). Currently we always run the same set.

4. **Better synthesis** — Our consolidation is decent but could adopt Gastown's deduplication pattern: note which legs found the same issue (higher confidence signal).

### Priority 4: Design Convoy (Nice to Have)

For complex features, fan out design exploration to multiple dimensions. Already possible with our Agent tool — just needs a skill definition. Lower priority because we do design in conversation already and it works okay for our scale.

---

## Implementation Plan

### Phase 1: `/spec-review` (New Skill) — HIGH VALUE

Create a new skill that reviews specs before implementation begins. Integrates into `/orchestrate` DESIGN phase.

**Structure:**
```
Invoke: /spec-review specs/feature/spec.md
Execution: 5 parallel subagents (requirements, gaps, ambiguity, feasibility, scope)
Output: questions for human + confidence table + GO/NO-GO
Integration: DESIGN phase gate in /orchestrate
```

**Leg prompts:** Adapt directly from Gastown's `mol-prd-review.formula.toml` — the "Look for" and "Questions to answer" sections are immediately usable. Swap Gastown-specific references (beads, polecats) for our terminology (specs, task plans).

### Phase 2: `/plan-review` (New Skill) — MEDIUM VALUE

Review task plans before execution. Integrates between DESIGN and IMPLEMENT phases.

**Structure:**
```
Invoke: /plan-review specs/tasks/task-name/task-plan.md
Execution: 5 parallel subagents (completeness, sequencing, risk, scope-creep, testability)
Output: GO / GO WITH FIXES / NO-GO + per-leg verdicts + must-fix
Integration: Post-DESIGN gate in /orchestrate
```

### Phase 3: `/code-review` Upgrades — MEDIUM VALUE

Add to existing skill:
1. Wiring leg definition
2. Resilience leg definition
3. Preset system (`/code-review --preset=gate` vs `/code-review --preset=full`)
4. Enhanced synthesis with cross-leg deduplication

### Phase 4: Design Convoy — LOW VALUE (defer)

Only if we find ourselves doing complex design work where conversation-based design isn't cutting it.

---

## Gastown Review Patterns: Reference Material

### Standard Leg Prompt Template (Adaptable)

Every Gastown review leg follows this structure:
```
**Look for:**
- Specific thing to check (with examples)
- Another specific thing (with anti-pattern description)
- ...

**Questions to answer:**
- Key question that forces judgment (not just listing)
- Another question framed from a specific perspective
- ...
```

This is much better than our current approach of giving agents a general focus area. The "Look for" list makes the review systematic. The "Questions to answer" forces judgment rather than just enumeration.

### Standard Output Template (Adaptable)

```markdown
# [Dimension] Review

## Summary
(1-2 paragraph overview)

## Critical Issues (P0 - Must fix)
- Issue with file:line reference
- Impact description
- Suggested fix

## Major Issues (P1 - Should fix)
- ...

## Minor Issues (P2 - Nice to fix)
- ...

## Observations
(Non-blocking notes)
```

### Synthesis Template (Adaptable)

```markdown
# Review Synthesis

## Executive Summary
(Overall assessment, recommendation)

## [Verdict Table]
| Dimension | Verdict | Key Finding |
|---|---|---|

## Must Fix Before Proceeding
### [Issue]
- Found by: [which legs]
- Problem: [what's wrong]
- Required fix: [what to change]

## Should Fix
- ...

## Observations
- ...

## Next Steps
- [ ] ...
```

### Key Design Principles From Gastown Reviews

1. **Each leg has a narrow focus** — not "review this code" but "review this code for EXACTLY these specific patterns"
2. **"Look for" lists are concrete** — not "check for security issues" but "check for input validation gaps, auth bypasses, injection vulnerabilities, hardcoded secrets, SSRF, deserialization vulns"
3. **Questions force judgment** — "What would a new team member not understand?" is better than "Is the code readable?"
4. **Output is structured and prioritized** — P0/P1/P2 with file:line references. Actionable, not advisory.
5. **Synthesis deduplicates and cross-references** — Issues found by multiple legs are higher confidence. Note which legs found them.
6. **Reviews produce QUESTIONS, not assertions** (for spec/PRD review) — The reviewer's job is to surface what's missing, not fill in gaps.
7. **GO/NO-GO gate on plan review** — Forces reviewer to commit to a binary judgment, not just list concerns.
