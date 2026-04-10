# Role System Specification

**Status:** Concept
**Created:** 2026-04-10
**Origin:** Gastown comparison during viewpoint-based calcinatio work

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[B:context]` | Business requirement |
| `[D:reason]` | Decision |
| `[E:existing]` | Extends existing |
| `[U:topic]` | Unbacked — needs validation |

---

## Overview

The athanor system currently has two agent roles: **azer** (general craftsman) and **marut** (supervisor). Everything else is ad hoc — subagents are briefed per-task with natural language, and their "role" exists only in the prompt that spawned them.

Gastown takes the opposite approach: it defines specific agent templates (polecat, witness, dog) with predefined CLAUDE.md files, TOML-defined formula legs with hardcoded perspectives, and a framework for composing task-specific behaviors from authored templates. `[B:gastown-comparison]`

Both approaches have strengths:

| | Athanor (current) | Gastown |
|---|---|---|
| **Agent definition** | Two general roles (azer, marut) | Multiple specific templates (polecat, witness, dog) |
| **Subagent perspectives** | Derived generatively at runtime | Authored statically in TOML formulas |
| **Flexibility** | High — azer adapts to anything | Lower — new task types need new templates |
| **Consistency** | Variable — depends on prompt quality | High — same template, same behavior |
| **Onboarding a new task type** | Zero setup, relies on agent judgment | Requires authoring formula + template |

The concept: **a role system that gives structure without restricting the general craftsman model.** Defined roles as a menu the azer can draw from — for dispatching subagents with specific expertise, or for adopting a lens when the work demands it. `[B:artifex-direction]`

---

## Domain Dictionary

**Role** — A defined professional identity with specific concerns, evaluation criteria, and domain expertise. Not a permanent assignment — a lens that can be adopted or dispatched as needed. Analogous to a "job" in Final Fantasy X's sphere grid / job system. `[D:FFX-analogy]`

**Role Adoption** — When an azer takes on a role's perspective for a portion of its work. The azer remains an azer (general craftsman) but channels the role's expertise. Like equipping a job class. `[U:adoption-mechanics]`

**Role Dispatch** — When an azer spawns a subagent briefed with a role's full context. The subagent operates as that role for its bounded task. This is the primary use case — viewpoint-based manifold calcinatio with role-defined perspectives. `[E:viewpoint-based-manifold-calcinatio]`

**Role Definition** — The authored artifact that captures a role's identity: what it cares about, what it evaluates, what expertise it brings, what fires it applies. Lives as a file in the athanor or project. `[U:format-and-location]`

**Role Menu** — The set of available roles for a given scope (athanor-level, project-level, or MO-level). Not all roles are relevant to all work — the menu is context-dependent. `[U:discovery-mechanism]`

---

## The Philosophical Delta

The key insight from comparing athanor and gastown:

**Gastown's power comes from predefined composition.** When you author a convoy formula with 10 review legs, you're crystallizing expertise about *which perspectives matter for code review*. The template is the knowledge. The agent just executes it. This produces consistency — every code review hits the same 10 dimensions.

**Athanor's power comes from generative derivation.** The azer reads the MO, reads the witnesses, reads the work, and *derives* what perspectives matter. This produces adaptability — every task gets perspectives tailored to its nature. But it also produces variance — the derivation quality depends on the agent's judgment.

**The role system bridges the gap.** Defined roles give the azer a vocabulary of perspectives to draw from — not replacing generative derivation, but enriching it. Instead of "I need a security perspective" (vague), the azer dispatches a subagent with the "security auditor" role (concrete, consistent, authored with expertise). Instead of "I need to think about performance" (ad hoc), the azer adopts the "performance analyst" role which comes with specific evaluation criteria.

**Roles are authored expertise, generatively composed.** The role definitions are static (like gastown's TOML). The *selection and composition* of which roles to apply is dynamic (like athanor's current approach). This is the synthesis: gastown's consistency in the role definitions, athanor's flexibility in which roles to use when. `[D:synthesis-of-approaches]`

---

## How Roles Interact with Existing System

### Viewpoint-Based Manifold Calcinatio

Roles are a natural source of perspectives for viewpoint-based manifold calcinatio. The three-source derivation process gains a fourth consideration:

1. **Witness-derived** — unchanged
2. **Research-derived** — unchanged
3. **Task-intrinsic** — now informed by the role menu: "which defined roles are relevant to this task?"
4. **Role-menu selection** — browse available roles, select those whose concerns intersect the current work

The role definition provides the subagent's briefing — instead of improvising "evaluate this from a security perspective," the azer passes the security auditor role's full context: what it evaluates, how it evaluates, what fires it applies, what it considers a finding. `[E:viewpoint-based-manifold-calcinatio]`

### Assessment Opera

Assessment already fires witness-perspective legs and materia-stimulus legs. Role-defined perspectives could become a third class of assessment leg — "which roles in the menu would see something the witness legs miss?" `[U:assessment-integration]`

### /code-review

Currently has 10 hardcoded legs. Each could become a role definition, gaining:
- Reusability beyond code review (the "security auditor" role is useful in spec review too)
- Customizability per project (a fintech project's security auditor role has different concerns than a consumer app's)
- Composability (mix code-review roles with project-specific roles for a richer review) `[U:code-review-migration]`

---

## Role Definition Shape

A role definition should capture enough for a subagent to operate with that role's expertise, without over-constraining how it works. `[U:needs-prototyping]`

```
# Role: [Name]

## Identity
What this role is. One sentence.

## Concerns
What this role cares about — the questions it asks, the things it notices,
the problems it's trained to find. These are the role's "fires."

## Evaluation Criteria
How this role judges quality. What "good" looks like from this perspective.
What findings look like — severity, format, evidence requirements.

## Domain Knowledge
Key concepts, frameworks, or references this role draws on.
Not exhaustive — enough to orient the agent adopting or dispatched with this role.

## Anti-Concerns
What this role explicitly does NOT evaluate. Prevents scope creep
and keeps the perspective narrow and independent.
```

**Location:** TBD — could live at multiple scopes:
- `~/.claude/roles/` — global roles available everywhere
- `<project>/.claude/roles/` — project-specific roles
- `<athanor>/roles/` — athanor-specific roles
- Role definitions in the athanor system spec's shared directory

**Discovery:** TBD — the azer needs a way to browse available roles. Could be a skill (`/roles`), a directory convention, or integration with `/skill-discovery`. `[U:discovery-mechanism]`

---

## Scope and Non-Goals

### In Scope (this spec)
- The concept and its relationship to existing athanor patterns
- Role definition shape (what a role file contains)
- How roles integrate with viewpoint-based manifold calcinatio
- The philosophical position: authored expertise, generatively composed

### Out of Scope
- Implementation details (file format, CLI commands, tooling)
- Migration of existing patterns (e.g., converting /code-review legs to roles)
- Role lifecycle (who authors roles, how they're reviewed, versioning)
- Role adoption mechanics for the azer itself (vs dispatch)
- Automated role suggestion / recommendation engine

### Open Questions

1. **Adoption vs dispatch balance.** The primary use case is dispatch (subagent with role context). Role adoption (azer channels the role) is hazier — does the azer load the role definition and shift its own perspective, or is it always better to dispatch a subagent for independence? `[U:adoption-mechanics]`

2. **Role granularity.** How specific should roles be? "Security auditor" is broad. "Authentication flow security analyst for multi-tenant SaaS" is narrow. Gastown's legs are narrow (one per review dimension). The right granularity probably varies by use case. `[U:granularity]`

3. **Composition rules.** Can roles compose? Is "security auditor + fintech compliance" a valid compound role, or should it be two separate dispatch legs? Keeping roles independent (like manifold calcinatio legs) probably produces better results. `[U:composition]`

4. **Relationship to skills.** Skills provide capabilities (how to do something). Roles provide perspectives (how to evaluate something). There may be overlap — a `/security-audit` skill and a "security auditor" role serve related purposes through different mechanisms. `[U:skill-role-relationship]`

---

## Test Scenarios

Since this is a concept spec, test scenarios describe how to validate the concept rather than automated tests:

1. **Prototype with viewpoint-based manifold calcinatio.** Author 3-5 role definitions for a concrete task (e.g., spec review). Use them as subagent briefings in a viewpoint-based review. Compare output quality vs. ad-hoc perspective derivation.
2. **Cross-task reuse.** Use the same "security auditor" role in both a code review and a spec review. Does the role definition produce relevant findings in both contexts?
3. **Project-specific roles.** Author a project-specific role (e.g., "PowerSync offline-first analyst" for musashi). Does it produce findings that generic roles miss?

---

## Retrospective

When this concept is prototyped, review:

### Documentation Updates
- [ ] Does the role definition shape need revision based on prototype experience?
- [ ] Should role system be added to the athanor spec (specs/spec.md)?
- [ ] Does viewpoint-based-manifold-calcinatio.md need updating to reference roles?

### Workflow Improvements
- [ ] Should a `/roles` skill be created for role discovery and management?
- [ ] Should `/code-review` be refactored to use role definitions?
- [ ] Should assessment opera reference the role menu?

### Knowledge Capture
- [ ] What worked and didn't work about the role definition shape?
- [ ] Which granularity level produced the best subagent output?
- [ ] How did authored roles compare to ad-hoc perspective derivation?
