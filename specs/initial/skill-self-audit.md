# Skill Self-Audit Protocol

**Status:** Concept — standard protocol (see `standard-protocols.md`)
**Date:** 2026-03-18
**Origin:** SAL-101 triage session — observed drift in `/task-triage` output (missed "why now" emphasis despite skill template requiring it)

## Problem

Skills with specific output requirements drift when loaded into a discursive main context. The skill's guidance competes with conversation flow, operator tangents, and accumulated context for the LLM's attention. The longer the session runs, the more the skill's requirements fade into background noise. This is not a skill authoring problem — it's an operational property of how LLMs allocate attention across large contexts.

Bandaids like bold text, repetition, and "injection in depth" help but don't solve the root cause: **the same context that drifted can't reliably audit itself for drift.**

## Insight

This is the same problem that code review solves for code authoring. The author's context is optimized for creation — they're thinking about the problem, making trade-offs, building momentum. A reviewer brings fresh attention to the artifact's requirements — style guide, correctness, completeness — because they haven't spent their attention budget on the creation process.

The same principle applies to skills. The main session uses the skill for guidance during creation. A fresh session uses the same skill to audit the output against the skill's own requirements.

## Pattern: Skill Self-Audit

For any skill that produces a structured artifact with specific field/quality requirements:

1. **Main context (author):** Load skill, use it to guide discursive work with the operator. The skill shapes the conversation and helps produce the artifact.

2. **Fresh context (auditor):** Spawn a background subagent that loads the *same skill* and audits the produced artifact against the skill's requirements. The auditor has full attention on the skill's spec because it hasn't spent context on the conversation.

3. **Report:** Auditor produces a gap report — fields missing, requirements unmet, quality issues — that surfaces back to the main context or operator.

### When to trigger the audit

- **End of skill execution** — after the operator confirms the artifact, run the audit before declaring done
- **Periodically during long sessions** — for skills used across extended conversations (e.g., `/orchestrate` managing a multi-phase pipeline)
- **On demand** — operator invokes the audit explicitly ("check this against the skill's requirements")

### What the auditor gets

- The skill's SKILL.md (full, uncompeted-with attention)
- The produced artifact
- Minimal context about what was asked for (the original input/goal)
- NOT the conversation history (that's the point — fresh eyes)

## Design Considerations

### Which skills benefit from self-audit?

Skills with **structured outputs** and **enumerable requirements**:
- `/task-triage` — 9 required fields, each with quality criteria
- `/spec` — required sections, traceability tags, completeness checks
- `/handoff` — required sections, embedded mandates, scope sizing
- `/qa-plan` — coverage matrix, channel classification completeness
- `/pr-create` — PR template sections, ticket linking, description quality

Skills that are purely procedural (e.g., `/ci-monitor` — run loop until green) or conversational (e.g., `/threads`) don't benefit because there's no structured artifact to audit.

### Skill anatomy changes

Skills that support self-audit would need a section (or companion file) defining:
- **Audit criteria** — the checklist the auditor uses (may overlap with existing skill guidance, but made explicit for machine consumption)
- **Audit mode entry** — how the skill behaves when invoked as auditor vs. author (e.g., "when given an existing artifact, evaluate it against these criteria rather than creating a new one")

This could be as simple as a `## Audit Mode` section in the SKILL.md that says: "When invoked with an existing artifact, evaluate it against the following criteria and report gaps."

### Relationship to existing review patterns

- **`/spec-check`** is already a self-audit pattern — it loads a spec and checks implementation against it. The difference is it audits *code* against a *spec*, not a *skill output* against a *skill's requirements*.
- **`/code-review`** audits code against general quality standards. Skill self-audit audits skill outputs against skill-specific standards.
- **`/task-review`** audits code against handoff requirements. Same family.

The pattern is general: **any artifact with defined quality criteria can be audited by loading those criteria into a fresh context and evaluating the artifact against them.**

### Cost/benefit

- **Cost:** One subagent spawn per skill execution. Typically small context (skill + artifact), fast execution.
- **Benefit:** Catches drift that the main context can't see. Prevents the "looks done but isn't" failure mode that wastes downstream time.
- **When NOT worth it:** Quick, low-stakes skill invocations. Audit overhead should be proportional to artifact importance.

## Implementation Path

1. **Pilot with `/task-triage`** — add an `## Audit Mode` section to the skill. After operator confirms triage, spawn auditor. Evaluate whether it catches real drift.
2. **If effective, extend to `/spec` and `/handoff`** — these are the highest-leverage structured artifacts.
3. **Standardize the pattern** — add self-audit as a documented capability in `/skill-creator` so new skills can opt into it.
4. **Consider making it a standard phase in `/orchestrate`** — after any skill produces a structured artifact, optionally run self-audit before proceeding.

## Relationship to Other Concepts

- **Injection in depth** — complementary, not competing. Injection in depth reduces drift during authoring. Self-audit catches what injection in depth missed.
- **Author/editor from agentic architecture** — self-audit is the formalized version of the "periodically spawn a review subagent" guidance. This makes it a standard pattern rather than ad-hoc.
- **Phase results protocol** — self-audit could produce a confidence score that feeds into phase results, giving the orchestrator a signal about artifact quality.
