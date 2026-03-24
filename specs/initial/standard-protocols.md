# Standard Protocols

**Status:** Concept — formalizing what already exists implicitly
**Date:** 2026-03-18
**Origin:** Observed that the orchestration system has accumulated several reusable behavioral patterns, but they're scattered across docs, described inconsistently, and there's no explicit concept of "protocol" as a composable unit.

## What is a Standard Protocol?

A **protocol** is a reusable behavioral pattern that any skill, phase, or agent can adopt. It's not a skill (skills produce artifacts or guide work). It's not a hook (hooks are triggered by tool events). A protocol is a **way of operating** — a pattern of behavior that's been proven useful and formalized so it can be applied consistently.

Protocols are composable. A skill can adopt multiple protocols. A phase can require certain protocols. An orchestrator can mandate protocols across all delegated work.

### Protocol vs. Skill vs. Hook

| Concept | What it is | Example |
|---------|-----------|---------|
| **Skill** | Produces an artifact or guides a specific type of work | `/task-triage` produces a task shape |
| **Hook** | Fires on a tool event, executes a command | `session-mode-inject.sh` injects context per-turn |
| **Protocol** | A behavioral pattern adopted by skills/phases/agents | "Escalation protocol" — how to halt and report when blocked |

### What makes something a protocol?

1. **Reusable across contexts** — applies to multiple skills, phases, or agent types
2. **Behavioral, not artifactual** — defines how to act, not what to produce
3. **Composable** — can be adopted alongside other protocols without conflict
4. **Testable** — you can observe whether it's being followed or not

## Existing Protocols (Currently Implicit)

These patterns already exist in the system but aren't formally recognized as protocols. Formalizing them means: naming them, defining them precisely, and making them explicitly adoptable.

### Escalation Protocol

**Where it lives now:** Core mandates, `resources/standard-agent-directives.md`

**Pattern:** When blocked by something outside your scope: halt → report (what, why, what's needed) → route to resolver → do not proceed until resolved.

**Adopted by:** All agents (via core mandates), but inconsistently applied.

### Phase Results Protocol

**Where it lives now:** Memory (`feedback_mise_en_place.md`), orchestrate skill

**Pattern:** On phase completion, append structured results (Produced, Confidence, Gate outcome) to the plan document. The plan IS the durable state — no separate manifest.

**Adopted by:** `/orchestrate` phases.

### Injection in Depth

**Where it lives now:** `/agentic-architecture` skill, architecture.md

**Pattern:** Every component is an injection point for core principles. When creating/modifying any component, ask "which mandates should be woven in?"

**Adopted by:** `/skill-creator`, `/handoff`, skill authoring generally.

### Author/Editor

**Where it lives now:** `/agentic-architecture` skill (Quality section), ad-hoc practice

**Pattern:** Periodically spawn a review subagent with fresh context to evaluate current session's work against its own standards.

**Adopted by:** Ad-hoc — not systematically applied.

### Anti-Workaround

**Where it lives now:** `resources/standard-agent-directives.md`

**Pattern:** When a tool or method doesn't work, escalate — don't silently substitute an alternative. The workaround might "work" but bypass the reason the original approach was specified.

**Adopted by:** All delegated agents (via standard directives).

## Proposed New Protocols

### Skill Self-Audit Protocol

**Defined in:** `specs/orchestrate-next/skill-self-audit.md`

**Pattern:** For skills that produce structured artifacts with enumerable requirements — after the artifact is produced in the main context, spawn a fresh-context subagent that loads the same skill in audit mode and evaluates the artifact against the skill's own requirements. Reports gaps back.

**Why it's a protocol, not a skill:** It's a behavioral pattern that any artifact-producing skill can adopt. The skill defines what to audit; the protocol defines when and how.

**Adoption:** Skill opts in by including an `## Audit Mode` section defining audit criteria and how to evaluate artifacts against them.

## Open Questions

1. **Where do protocols live?** Options: (a) each protocol gets its own doc in `specs/orchestrate-next/`, (b) protocols get a dedicated directory, (c) protocols are sections in a single doc. Leaning toward (a) for now — individual docs are discoverable and can evolve independently.

2. **How are protocols adopted?** Currently implicit ("this skill follows the escalation protocol"). Should adoption be explicit in skill metadata? In the SKILL.md frontmatter? In the skill body? Explicit adoption helps auditing ("does this skill follow the protocols it claims to?").

3. **Protocol compliance checking?** The self-audit protocol is itself a mechanism for checking compliance with other protocols. Meta, but practical — a self-audit could check "did the escalation protocol get followed?" as one of its criteria.

4. **Relationship to core mandates?** Some protocols (escalation, anti-workaround) are currently embedded in core mandates. Are core mandates just "mandatory protocols"? Or is there a meaningful distinction? Tentative answer: core mandates are protocols that are always adopted by every agent. Standard protocols are opt-in per skill/phase.
