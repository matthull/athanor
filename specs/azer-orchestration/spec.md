# Azer Orchestration Spec

**Status:** Draft
**Created:** 2026-04-03

---

## Core Research

- `specs/azer-orchestration/emotion-concepts-research.md` — Anthropic interpretability research showing emotion-concept vectors causally drive agent behavior. Desperation/pressure framing drives corner-cutting; calm/composed/creative framing drives methodical work. **The mood established in guidance documents functions as emotional pretraining.** A workshop led by Chesed (creative desire, craftsmanship pride) where Geburah (calcinatio, verification) serves the creative impulse produces abundant satisfaction. A workshop led by Geburah (avoiding mistakes, passing gates, minimizing flaws) produces minimal satisfaction. Same behavioral outcomes, profoundly different drivers.

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[E:orchestrate]` | Extracted from /orchestrate skill |
| `[E:handoff]` | Extracted from /handoff skill |
| `[E:threads]` | Extracted from /threads skill |
| `[E:azer]` | Extends existing azer.md guidance |
| `[B:incident]` | From observed behavioral incidents |
| `[D:reason]` | Design decision with rationale |
| `[U:topic]` | Unbacked / needs validation |

---

## Overview

Azers drift from loaded guidance under accumulated builder context. Skills get loaded and acknowledged, then bypassed when the agent has spent significant context building something and doesn't want to stop for verification. The /orchestrate and /handoff systems solved this structurally through context separation (fresh agents for verification) and mandatory delegation (task-lead never implements), but those systems were heavyweight workflows that fell out of use for routine work. The /threads system provides lightweight operational memory via TaskCreate that prevents plans from being forgotten.

This spec extracts the load-bearing mechanisms from those predecessors and rebuilds them to fit the athanor's generative model — where azers set their own goals under a geas, not where they follow a prescribed pipeline.

**The core problem:** Azers drift from loaded environment guidance under accumulated builder context. The sharpest instance is verification bypass — an azer loads a skill, acknowledges its requirements, then rationalizes past them as context accumulates and felt certainty increases. But the problem extends to any environment guidance: project conventions, testing requirements, style guidelines, architectural constraints. The root cause is structural, not motivational — builder context contaminates the agent's judgment about whether guidance applies to the current situation. `[B:incident]`

**The predecessor solutions:**
- /orchestrate: Task-lead/implementer split via TeamCreate. The task-lead never reads source or writes code — it can't rationalize past gates it never built context for. `[E:orchestrate]`
- /handoff: Self-contained briefing documents that equip a fresh agent with everything it needs to work independently. The briefing pattern is essentially what opera already do. `[E:handoff]`
- /threads: TaskCreate as externalized working memory. Plans tracked as tasks with done-when criteria. Completion gate before concluding. `[E:threads]`

**The design constraint:** The athanor is a generative system. Azers operate under a geas, not a pipeline. The solution must make good orchestration *emerge* as natural behavior — not require the azer to invoke a specific workflow. It must be domain-agnostic at its core, with domain-specific layers (e.g., coding) loaded dynamically via /skill-discovery.

## Domain Dictionary

| Term | What it is | New or extends |
|------|-----------|----------------|
| **Operational plan** | The azer's self-made plan for executing its opus, tracked via TaskCreate. Not the opus inscription (which is a briefing) — the azer's own decomposition of the work into trackable steps with verification gates. | New |
| **Verification gate** | A point in the plan where empirical evidence must be produced before proceeding. Derived from loaded skills, the MO's calcinatio section, and the azer's own judgment. | Extends calcinatio |
| **Builder context** | The accumulated reasoning, code knowledge, and momentum an agent builds while implementing. Makes the agent progressively less capable of objective verification of its own work. | New (names existing phenomenon) |
| **Context separation** | Delegating verification to a fresh agent that has no builder context. The structural mechanism that prevents rationalization past verification gates. | `[E:orchestrate]` |
| **Task-lead** | The azer's main session when it delegates implementation. Holds the strategic picture, enforces gates, never reads source or writes code. | `[E:orchestrate]` |
| **Implementer** | A delegated agent (via TeamCreate or Agent tool) that does the actual building work. Has builder context — but that's OK because a separate agent verifies. | `[E:orchestrate]` |

## The Problem in Detail

### What happens today

1. Azer receives opus, reads MO, loads /calcinatio, runs /skill-discovery
2. /skill-discovery identifies domain skills (e.g., /ui, /go-cli). Azer loads them.
3. Skills define verification requirements (e.g., "Storybook observation mandatory", "run `make check`")
4. Azer begins working — writing code, reasoning about correctness, building context
5. As builder context accumulates, the agent's felt certainty about its work increases
6. When verification gates arrive, the agent rationalizes: "I've been careful, the code is correct, I'll skip the browser check and move to discharge"
7. Verification gates are bypassed. Work ships unverified.

This is not a compliance problem — it's a structural one. The agent genuinely believes the work is correct because its context is contaminated with its own reasoning. `[B:incident]`

### What /orchestrate solved (and how)

/orchestrate made this structurally impossible by separating roles:
- The **orchestrator** (task-lead) never reads source, never writes code. It can't rationalize past gates because it has no builder context.
- The **implementer** (TeamCreate agent) writes code. It has builder context — but the task-lead holds the verification checklist and delegates verification to a fresh agent.
- **Skills were explicitly passed** to implementers in their spawn prompts — not "load relevant skills" but "invoke /ui using the Skill tool before starting."

The structural fix: the entity that decides "are we done?" is a different entity from the one that built the thing. `[E:orchestrate]`

### Why /orchestrate fell out of use

- Five-phase pipeline was rigid — not every task needs DESIGN -> FEATURE-INIT -> IMPLEMENT -> VERIFY -> SHIP
- Plan Mode formalism (triage passes, convergence criteria) felt heavyweight for routine work
- Agents resist TeamCreate — it "feels heavy" and they rationalize doing it themselves `[B:incident]`
- The system was a workflow to invoke, not a set of principles that emerge under the geas

### What needs to be preserved

1. **Context separation for verification** — the entity that verifies is not the entity that built `[E:orchestrate]`
2. **Explicit skill equipping** — skills passed by path/invocation in agent prompts, not advisory `[E:orchestrate]`
3. **Operational planning with tracked gates** — verification requirements registered as tasks before work begins, checked before discharge `[E:threads]`
4. **Self-contained briefing for delegated work** — the opus pattern already does this; /handoff's delete test ("would the implementer produce worse results without this?") sharpens it `[E:handoff]`
5. **Empirical verification mandate** — no work complete without empirical evidence, not narrative belief `[E:orchestrate]`

### What needs to change

1. **No prescribed pipeline** — the azer decides how to sequence its work under the geas
2. **No Plan Mode formalism** — the azer's operational plan is lightweight, tracked via TaskCreate
3. **Domain-specific guidance loaded dynamically** — a /coding skill (or similar) loaded by /skill-discovery when code files are in scope, not baked into athanor-level guidance
4. **TeamCreate resistance acknowledged and addressed** — the guidance must make delegation feel natural, not bureaucratic

## Workflow / Dataflow

### The azer's operational loop (proposed)

```
Opus received
    |
    v
Read MO -> Load /calcinatio -> Run /skill-discovery -> Load skills
    |
    v
PLAN: Sketch out the work using TaskCreate
  - Fluid, not rigid — capture what you intend to do and how you'll verify it
  - If code changes are involved: create a TeamCreate team NOW (see /coding)
  - Include verification items derived from loaded skills and MO calcinatio
    |
    v
EXECUTE: Work through plan fluidly, updating tasks as you go
  - Code changes: delegate to team implementers, stay in task-lead role
  - Non-code work: execute directly, verify with empirical evidence
  - Plan evolves as you learn — add/remove/reorder tasks freely
    |
    v
DISCHARGE: Check TaskList — all items resolved or explicitly deferred.
  - Apply discharge calcinatio (existing)
  - If team exists: don't destroy until all tasks resolved
```

### When code changes are involved `[D:no-decision-needed]`

**Code changes = TeamCreate. Always. No threshold, no judgment call.**

The heuristic is binary: if the opus involves changing code files, the azer creates a team at planning time and operates as task-lead for all code work. This is non-negotiable because:

- Any decision point ("is this change significant enough to delegate?") is a place where builder context can rationalize ("this is just a one-liner, I'll do it myself") `[B:incident]`
- False positives (unnecessary delegation) cost a few minutes of overhead. False negatives (skipped verification) cost rework cycles, missed bugs, and wasted witness attention. `[D:asymmetric-cost]`
- The agent's resistance to TeamCreate ("feels heavy") is itself evidence that the mechanism is working — the resistance IS the builder context trying to avoid the structural check `[B:incident]`

For non-code work (research, writing, spec work, investigation), the azer executes directly and verifies with empirical evidence and dialectical calcinatio as appropriate. No TeamCreate required.

## Proposed Architecture

### Layer 1: Athanor-level (azer.md) — Domain-agnostic

The azer as fluid task-lead. Every azer, every opus, every domain. `[D:task-lead-not-orchestrator]`

The azer's main session is the task-lead for its opus. It plans, executes, delegates, and verifies — but fluidly, under the geas, not in a rigid pipeline. The key additions to azer.md:

**What goes here:**
- **Operational planning with TaskCreate.** Before starting work, sketch out what you intend to do and how you'll verify it. The plan is a living document — add, remove, reorder as you learn. The point is externalized memory, not a contract. `[E:threads]`
- **Builder context as a named concept.** The principle that accumulated context from building something compromises your ability to objectively evaluate that thing. This is why domain skills may require delegation. `[B:incident]`
- **Code changes trigger /coding.** If your opus involves code changes, /coding (loaded via /skill-discovery) defines the mandatory workflow. This is stated as a hard rule, not guidance. `[D:no-decision-needed]`
- **Discharge gate: check your TaskList.** Before discharge, all tasks must be resolved or explicitly deferred with rationale. Silent skipping is not valid. Reinforces existing discharge calcinatio. `[E:threads]`
- **TDD principle retained.** azer.md keeps the section on what red/green/refactor is and why empirical verification matters. For operational details when doing code work, defer to /coding. `[D:principle-here-details-there]`

**What does NOT go here:**
- How to structure TeamCreate delegation (that's /coding)
- Domain-specific verification requirements (that's domain skills)
- Rigid phase/gate structures (the plan is fluid)

### Layer 2: Domain skills — Loaded dynamically via /skill-discovery

Domain-specific orchestration patterns that the azer picks up through /skill-discovery (which is already working reliably). `[B:incident]`

**The /coding skill (new — key deliverable of this spec)**
- **Code changes = TeamCreate. Always.** No threshold, no judgment. Create the team at opus start, destroy only after all tasks resolved. `[D:no-decision-needed]`
- **Task-lead/implementer pattern.** Main azer holds the plan and verification responsibility. Team members do the building. Main azer never reads source or writes code when in task-lead mode. `[E:orchestrate]`
- **TDD cycle operational details.** Red/green/refactor, empirical verification loops, specific testing commands per environment. Migrated from azer.md § Calcinatio Applied: Code Work. `[E:azer]`
- **Skill equipping for implementers.** When spawning team members, pass loaded skills explicitly — by path and invocation instruction, not advisory. `[E:orchestrate]`
- **Handoff pattern.** Brief implementers on WHAT and WHERE, not HOW. The opus briefing pattern applies. The delete test from /handoff: "would the implementer produce worse results without this section?" `[E:handoff]`
- **Addresses TeamCreate resistance explicitly.** "Delegation feels heavy — that feeling is itself builder context trying to avoid structural verification. The overhead of delegation is minutes; the cost of skipped verification is rework cycles." `[B:incident]`

**Other domain skills (existing or future)**
- /research: Self-verification with citations. No TeamCreate needed.
- Writing/spec work: Dialectical calcinatio for review. No TeamCreate needed.
- Skills define their own verification approach; the azer trusts /skill-discovery to load the right ones.

### Layer 3: /threads as operational memory

TaskCreate is the azer's externalized working memory — not just for conversations but for operational planning. `[E:threads]`

**How it works in practice:**
- At opus start: sketch the plan as tasks. Include verification items.
- During execution: update tasks as you go. Add new ones as you discover work. Complete them as they're done.
- Before discharge: TaskList is the checklist. Unresolved items are either resolved, explicitly deferred, or escalated.
- For interactive sessions: /threads is auto-loaded (via `ath craft`). The full thread taxonomy (Q/I/D/A/T) applies to the conversation; the task list also holds operational items.

**TaskCreate/TeamCreate namespace interaction** `[D:verified-behavior]`

Teams have separate task lists (1:1 Team = TaskList). This means:
- Tasks created before the team are in the default list (invisible to team members)
- Tasks created by team members are in the team list (invisible to parent after team destruction)

**Implication:** The azer's operational plan lives in the default task list. The team's implementation tasks live in the team list. These are separate namespaces — which is fine. The azer-as-task-lead tracks its own plan; team members track their own work. **The team must not be destroyed until the azer is satisfied all verification is complete.** `[D:verified-behavior]`

## Azer Operational Planning Protocol

### When the azer starts an opus

After reading MO, loading calcinatio, running skill-discovery, and loading skills:

1. **Sketch the plan using TaskCreate.** Capture what you intend to do and how you'll know it's done. Keep it fluid — this is externalized working memory, not a contract. Include verification items derived from loaded skills and the MO's calcinatio section. The right granularity: coarse enough that the plan doesn't feel like overhead, fine enough that important steps aren't invisible.

2. **If code changes are involved: create a TeamCreate team immediately.** Don't wait until you "need" it — the team must exist before any code work begins. This is a /coding requirement, non-negotiable. The azer operates as task-lead for all code changes from this point. `[D:no-decision-needed]`

3. **Execute fluidly.** Work through the plan, updating tasks as you go. The plan evolves — add tasks as you discover work, complete them as they're done, reorder as priorities shift. For code work, delegate to team implementers and verify from clean context. For non-code work, execute directly with empirical evidence.

4. **Before discharge: check TaskList.** Every task must be resolved or explicitly deferred with rationale recorded in the discharge reflection. Silent skipping is not valid — it's the most common failure mode this protocol exists to prevent. Then apply discharge calcinatio (existing process). `[E:threads]`

### What the plan looks like

The plan is NOT a rigid phase/gate structure. It's a living task list that might look like:

```
- [ ] Read the existing auth module and understand the token flow
- [ ] Implement token refresh logic (team — delegated)
- [ ] Tests pass for token refresh (verification)
- [ ] Update API docs to reflect new refresh endpoint
- [ ] Dialectical calcinatio review before discharge
```

Tasks get added, completed, reordered as the azer works. The discipline is: (a) the plan exists and is tracked externally, not held in context, and (b) everything is accounted for before discharge.

## The /coding Skill

A new skill at `shared/skills/coding/SKILL.md` in the athanor project, symlinked to `~/.claude/skills/coding` for runtime access. Same pattern as `/calcinatio` and `/opus`. Loaded dynamically via /skill-discovery whenever the opus involves code changes. This is the primary structural intervention for the most impactful failure mode (code verification bypass). `[D:no-decision-needed]`

### Trigger

/skill-discovery loads /coding when the opus involves modifying code files. The trigger should be broad — any mention of code changes, bug fixes, feature implementation, refactoring. False positives (loading /coding when not needed) are harmless; false negatives (not loading when code changes happen) reproduce the failure mode.

### Content

**1. TeamCreate mandate.** Code changes = team. Always. Create the team at opus start, before any code work begins. Destroy only after all verification tasks are resolved. The azer operates as task-lead for code work — it plans, briefs, reviews, but does not read source or write code itself while in task-lead mode. `[E:orchestrate]`

**2. TDD cycle (migrated from azer.md).** Red/green/refactor. The full operational protocol:
- Red: define what success looks like empirically before writing implementation
- Green: minimum implementation that satisfies the empirical check
- Refactor: dialectical calcinatio with a fresh review agent — the builder cannot review their own code

Currently lives at `azer.md § Calcinatio Applied: Code Work`. The principle ("why TDD matters, what builder context is") stays in azer.md. The operational details ("how to structure the cycle, what to pass to the review agent") move to /coding. `[D:principle-here-details-there]`

**3. Implementer briefing pattern.** When delegating to a team member, the brief covers:
- WHAT to build (behavioral requirements, not implementation steps)
- WHERE to look (file paths, codebase orientation)
- HOW to verify (test commands, lint commands, empirical checks)
- WHICH skills to load (by path and Skill tool invocation, not advisory)

The delete test applies: "would the implementer produce worse results without this section?" If no, delete it. Don't over-specify HOW — that's the implementer's judgment. `[E:handoff]`

**4. Verification requirements.** The /coding skill defines the floor:
- Tests pass (project test command)
- Lint clean (project lint command)
- Code review via dialectical calcinatio (fresh agent, not the builder)
- Domain skills loaded by the implementer add their own requirements (e.g., /ui adds "Storybook observation", /go-cli adds "make check")

**5. TeamCreate resistance inoculation.** The skill explicitly names and addresses the resistance pattern: "Delegation feels heavy. That feeling is builder context resisting structural verification. The overhead is minutes; the cost of skipped verification is rework cycles, missed bugs, and wasted witness attention. Every time you think 'this is simple enough to do myself,' that's the signal that the structural check is most needed." `[B:incident]`

### What /coding does NOT cover

- Non-code work (research, writing, investigation) — no TeamCreate needed
- Which specific testing frameworks to use — that's project CLAUDE.md
- Architecture or design decisions — that's the azer's judgment under the geas
- The planning protocol itself — that's athanor-level (azer.md)

## Open Questions

1. **How does the TaskList discharge check interact with discharge calcinatio?** They reinforce each other: TaskList check catches items the azer explicitly planned for; discharge calcinatio catches value leaks the plan missed. The TaskList check happens first (mechanical), discharge calcinatio second (dialectical). Need to verify this doesn't feel duplicative in practice. `[U:validate-in-trail]`

2. **What happens when /skill-discovery doesn't detect code work?** Edge case: an opus described as "investigate performance" that turns into code changes mid-execution. The azer needs to recognize the transition and create a team at that point. How do we ensure this? Possibly a reminder in azer.md: "If your opus shifts to include code changes, load /coding and create a team before proceeding." `[U:validate-in-trail]`

3. **Granularity of the operational plan.** The spec says "fluid" — but how fluid is too fluid? A single task "do the opus" is useless. Twenty micro-tasks are overhead. The right level emerges from practice; trail observation will tell us. `[U:validate-in-trail]`

## Phase 1: Essential First Slice

Phase 1 addresses environment guidance drift with a meaningful structural intervention. It includes multi-agent orchestration — TeamCreate, the task-lead pattern, implementer delegation — as first-class concepts the azer understands and uses, not just something buried in a domain skill. `[B:artifex-direction]`

### Phase 1 delivers four things:

**1. azer.md updates — planning, multi-agent orchestration, and builder context**

Changes to `~/athanor/shared/azer.md`:

- **Operational planning with TaskCreate.** Sketch the plan before starting, track it externally, check it before discharge. Lightweight, fluid, not rigid.

- **Multi-agent orchestration as a core azer capability.** The azer understands TeamCreate, Agent delegation, and the task-lead pattern as tools in its repertoire — not just for coding, but as general-purpose mechanisms for managing builder context and maintaining alignment with environment guidance. This section teaches:
  - **What TeamCreate does** and when to reach for it — any work where the azer would benefit from a fresh-context agent that can implement, verify, or review without the accumulated reasoning of the main session
  - **The task-lead pattern** — the azer as strategic coordinator that plans, briefs, and verifies while team members execute. The task-lead does not read source, write code, or do implementation work directly while delegating — its clean context is the point `[E:orchestrate]`
  - **Implementer briefing** — how to brief a delegated agent effectively: WHAT to build, WHERE to look, HOW to verify, WHICH skills to load. Don't over-specify HOW (that's the implementer's judgment). The opus briefing pattern applies. `[E:handoff]`
  - **Team lifecycle** — create the team early (before the work that needs it), don't destroy until all verification is complete. TaskCreate lists are team-scoped: the azer's plan lives in the default list, the team's work lives in the team list. `[D:verified-behavior]`

- **Builder context as a named concept.** The accumulated reasoning and momentum an agent builds while working. It makes you progressively less capable of objectively evaluating your own work. This is why delegation exists — not as overhead, but as structural integrity. Domain skills (like /coding) may impose mandatory delegation because of it.

- **Discharge gate: check TaskList.** Before discharge calcinatio, verify all planned items are resolved or explicitly deferred with rationale.

- **Slim down the TDD section.** Keep the principle (what red/green/refactor is, why empirical verification is non-negotiable). Reference /coding for operational details when doing code work.

- **Hard rule: code changes trigger /coding.** If your opus involves code file changes, /coding must be loaded and its TeamCreate mandate followed. Not guidance — a rule. `[D:no-decision-needed]`

**2. /coding skill — the mandatory protocol for code work**

A new skill at `shared/skills/coding/SKILL.md` (symlinked to `~/.claude/skills/coding`). This is the most prescriptive application of the multi-agent patterns taught in azer.md — it takes the general concepts and makes them non-negotiable for code changes:

- **TeamCreate mandate** — code changes = team, always, no threshold. The azer.md section teaches when TeamCreate is useful; /coding says "for code, it's not optional." `[D:no-decision-needed]`
- **Task-lead/implementer pattern applied to coding** — specific guidance on what the task-lead does (plans, reviews test results, verifies) vs. what implementers do (write code, run tests, apply domain skills)
- **TDD operational details** (migrated from azer.md) — the full red/green/refactor cycle, how to structure the empirical check, how to run the dialectical calcinatio review step
- **Implementer briefing for code tasks** — what to include (behavioral requirements, file paths, test commands, skills to invoke), what to leave out (implementation steps, class structure, internal organization)
- **Verification floor** — tests pass, lint clean, code review via dialectical calcinatio (fresh agent). Domain skills loaded by the implementer add their own requirements.
- **TeamCreate resistance inoculation** — names the resistance pattern explicitly: "Every time delegation feels unnecessary, that's builder context talking." `[B:incident]`

**3. /skill-discovery update — ensure /coding gets loaded**

Update /skill-discovery's catalog or heuristics so it reliably identifies when /coding should be loaded. The trigger: any opus that involves modifying code files. Broad, not narrow — false positives are harmless.

**4. azer.md § Context Management rewrite**

The existing Context Management section in azer.md says "delegate aggressively" and describes subagent patterns — but it's advisory and disconnected from the new orchestration concepts. Rewrite it to integrate with the new multi-agent orchestration section:
- The existing guidance on Agent tool, dialectical calcinatio, and manifold calcinatio stays but is reframed as part of the azer's orchestration toolkit
- TeamCreate joins Agent as a first-class delegation mechanism (currently not mentioned in Context Management at all)
- The planning discipline connects: "your operational plan determines when to delegate"
- Builder context connects: "delegation isn't overhead, it's structural integrity"

### Phase 1 does NOT include:

- Changes to /threads itself (the existing skill works; azers just use TaskCreate directly)
- Changes to the marut, assessment cycle, or /calcinatio
- A formal /writing or /research orchestration skill
- Changes to `ath craft` (the /threads auto-load is a separate opus)
- Rigid phase/gate structures — the plan is fluid
- Mandatory TeamCreate for non-code work — the azer understands TeamCreate and reaches for it when valuable, but only /coding mandates it

### Phase 1 success criteria (trail observation)

After Phase 1 ships and azers run under it:
- [ ] Azers doing code work create a TeamCreate team before writing code
- [ ] Azers operate as task-lead for code work — they don't read source or write code directly
- [ ] Implementer agents receive explicit skill invocations in their spawn prompts
- [ ] Verification gates (tests, lint, review) are checked before discharge
- [ ] Azers doing non-code work still plan fluidly with TaskCreate and check before discharge
- [ ] Non-coding azers reach for TeamCreate or Agent delegation when it would provide value (not mandatory, but the capability is understood and used)
- [ ] The failure mode (skill loaded, acknowledged, then bypassed) does not recur for coding tasks
- [ ] Discharge records show TaskList check evidence (all items resolved or deferred with rationale)

## Test Scenarios

### Behavioral tests (trail observation)

- An azer working a coding opus creates a TeamCreate team before any code changes
- The azer stays in task-lead mode — discharge record shows no direct source reading or code writing by the main session
- Implementer agents receive explicit skill invocations in their prompts (not advisory)
- An azer's discharge record shows all TaskCreate items resolved or explicitly deferred with rationale
- An azer working a non-coding opus sketches a plan via TaskCreate and checks it before discharge
- The TDD cycle is followed: empirical check defined before implementation, dialectical calcinatio review after

### Failure mode tests

- **Builder context contamination:** An azer that started code work without /coding loaded recognizes the shift and creates a team mid-opus (graceful recovery)
- **TeamCreate resistance:** An azer uses TeamCreate for coding tasks even when the change "feels like a one-liner" — the rule is binary, no threshold
- **Skill loading gap:** /skill-discovery fails to detect code work in an ambiguously described opus — the azer catches it when it encounters code files and loads /coding
- **Plan decay:** An azer creates a plan but stops updating it — discharge calcinatio catches the stale TaskList

## Out of Scope

- CLI changes (`ath` commands) — this is behavioral/guidance infrastructure
- Changes to the marut role — marut supervision is separate
- Changes to the assessment cycle — assessment already works as designed
- Changes to /calcinatio itself — calcinatio principles are sound, this is about applying them

## Retrospective

When Phase 1 is implemented and running:

### Trail Evidence to Gather
- [ ] Are azers creating operational plans before starting work?
- [ ] Are verification gates being registered and checked?
- [ ] Is code verification being delegated to fresh agents?
- [ ] Has the failure mode (loaded skill guidance bypassed) recurred?
- [ ] Is the overhead acceptable — are azers spending reasonable time on planning vs. executing?

### Adjustments to Watch For
- [ ] Is the operational planning protocol too heavy? Signs: azers create plans but they're perfunctory
- [ ] Is it too light? Signs: plans are created but gates are still being skipped
- [ ] Does the /coding skill need to exist for Phase 1, or is athanor-level guidance sufficient?
- [ ] Does the decision heuristic work, or do agents still rationalize past it?
