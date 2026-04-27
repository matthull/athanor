# Job System Specification

**Status:** Implemented (Wave 1 complete, Wave 2 in progress)
**Created:** 2026-04-21
**Updated:** 2026-04-27
**Origin:** Interactive opus, athanor-alignment MO — artifex-directed
**Supersedes:** `specs/role-system/spec.md` (concept stage — this spec evolves and replaces it)
**Related:** `specs/azer-orchestration/spec.md` (diagnosis of the drift problem this solves)

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[B:artifex]` | Direct from the artifex during shaping |
| `[D:reason]` | Design decision with rationale |
| `[E:existing]` | Extends existing system element |
| `[E:gastown]` | Prior art from Gastown TOML roles/formulae |
| `[U:topic]` | Unbacked — needs validation |
| `[V:voice-note]` | From artifex voice note (transcribed) |
| `[I:implemented]` | Verified as implemented |
| `[O:observed]` | Observed in production behavior |

---

## Overview

Context drift in long-lived agent sessions causes skills not to load and critical tasks like QA to get skipped. The root cause is structural: prescriptive guidance (skills, verification gates, workflow requirements) erodes as builder context accumulates in large context windows. `[E:azer-orchestration/spec.md]`

The job system addresses this by introducing **job roles** — perspective-first specializations that an azer can adopt. Each job is a markdown file (`jobs/<name>/JOB.md`) that defines a professional identity: what you care about, what tools to reach for, what your instincts are. A job-role azer boots into a narrow perspective in a tight ~100k token context. The perspective stays salient because the context is small and the identity constrains what the azer attends to. `[B:artifex]`

**Jobs are the flexible prescriptive layer of the system.** They're where crystallized operational learning lives — "load these skills", "follow this verification sequence", "apply these fires". The general-purpose azer remains for novel work. Jobs are an extra tool, not a confinement. `[B:artifex]`

### The collaboration principle

Upstream of jobs is a Chesed-level principle that applies to all azers, job-role or general: **seek collaboration. A craftsman who finds the right collaborator produces better work than one who does everything alone.** `[B:artifex]`

This is not "don't do it yourself" (Geburah — restrictive). It is "get help, your work will be better for it" (Chesed — expansive). An azer that inscribes a QA specialist isn't following a rule against doing QA — it's recognizing that a fresh specialist with the right prescriptions will produce better verification than it can with accumulated builder context.

The job system makes collaboration easier by providing defined specialists to call on. But the impulse to collaborate is upstream of jobs — a general-purpose azer doing exploratory work should also inscribe peers when it recognizes that another perspective would produce better outcomes. `[B:artifex]`

### The inscription-as-spawning model

Azers create collaborators by inscribing opera with job roles and mustering peers. This replaces TeamCreate for most use cases. `[B:artifex]`

| TeamCreate (current) | Inscription spawning (proposed) |
|---|---|
| Subagent lives in parent's context budget | Peer azer gets its own crucible |
| Parent is task-lead, subagent is implementer | Peers collaborate via whisper |
| No trail — subagent work is invisible | Full trail — each specialist discharges its own opus |
| One big context holding everything | Many small contexts, each narrowly focused |

TeamCreate may remain useful for truly mechanical subtasks (search 30 files, run a computation) where a full opus is overhead. But for anything that benefits from trail visibility or prescriptive guidance, inscription wins.

**No re-specialization.** Azers are lightweight and disposable (~100k tokens). If the work shifts, don't re-specialize — inscribe and muster a new azer with the right job. The cost of a fresh context is low; the cost of a contaminated context is high. `[B:artifex]`

**What changes:**
- A geas-level collaboration principle: seek help, inscribe peers, your work is better for it
- Azers can adopt defined job roles that prescribe their perspective, skills, and workflow
- Context windows shrink (~100k tokens, discharge early) — jobs are scoped tight because they're narrow
- Orchestration becomes more dynamic — azers inscribe and muster peer azers, marut facilitates
- Inscription replaces TeamCreate as the primary spawning mechanism
- A decision point exists at inscription time: job-role azer or general-purpose azer

**What stays:**
- Foundational roles (marut, azer, perceiver, attendant)
- The geas, the MO, the Tree of Work
- The opus as briefing
- The trail and discharge protocol

---

## Domain Dictionary

| Term | What it is | New or extends |
|------|-----------|----------------|
| **Job** | A defined specialization an azer can adopt. A `JOB.md` in `shared/jobs/<name>/` that defines the professional perspective — what you care about, what tools to reach for, what your instincts are. Analogous to FF job system — equippable professional identity. Inspired by Gastown roles/formulae but perspective-first, not process-first. | New `[B:artifex]` `[E:gastown]` |
| **Job-role azer** | An azer mustered with a specific job. Boots into the job's narrow perspective. The job's prescriptions are mandatory, not advisory. Still an azer under the geas — the job constrains perspective, not agency. | New `[B:artifex]` |
| **General-purpose azer** | Today's azer model. No job role — operates under the geas with full latitude. Used when no defined job fits the work, or when the work is novel/exploratory. The fallback, not the default. | `[E:existing]` |
| **Job definition** | A `JOB.md` file that defines a job's perspective — identity, concerns, tools, instincts. Lives in `jobs/<name>/JOB.md`. Markdown is the prompt injection — read directly into the azer's context at boot. | New `[E:gastown]` |
| **Job boundary** | Not a wall but a professional instinct. A job-role azer naturally thinks in its specialty; when work of a different shape appears, its instinct is to seek collaboration rather than stretch. A QA azer that finds a typo might fix it (good judgment); one that finds a systemic issue inscribes a fixer (professional instinct). The motivation is quality, not compliance. | New `[B:artifex]` |
| **Job registry** | The set of available job definitions in `shared/jobs/*/JOB.md`. Discoverable at muster time. Grows as the system learns what patterns recur. | New `[E:gastown]` |
| **Inscription-time job assignment** | The job role is defined at inscription time — when the opus is created. The inscriber (marut, artifex, or another azer) specifies which job applies. | New `[B:artifex]` |

---

## The Problem This Solves

### Why skills drift

The current model loads prescriptive guidance (skills, AGENTS.md, azer.md) into a general-purpose context that then accumulates builder context. By the time the azer is deep in implementation, the prescriptions have been buried under hundreds of messages of reasoning. The agent genuinely believes it's following guidance — but its felt certainty about its own work overrides the loaded requirements. `[E:azer-orchestration/spec.md]`

**The azer-orchestration spec** diagnosed this precisely and proposed structural fixes within the current model (task-lead/implementer split, explicit skill equipping, operational planning with tracked gates). Those fixes work but rely on the azer choosing to apply them — which is itself subject to drift. `[E:azer-orchestration/spec.md]`

### Why jobs fix it

Jobs attack the root cause: **the context is too large and too general.**

1. **Smaller context (~100k tokens):** A job-role azer discharges at ~100k, not ~200k. The prescriptions loaded at boot are still within recent context when discharge triggers.
2. **Narrower perspective:** A QA azer doesn't hold implementation context. A code-implementation azer doesn't hold assessment context. The job boundary prevents the accumulation of irrelevant builder context that buries prescriptions.
3. **Mandatory prescriptions at boot:** The job definition specifies required skills and workflow. These aren't advisory ("consider loading /coding") — they're structural ("this job loads /coding and /unit-testing at boot").
4. **Collaboration-driven handoff:** When a job-role azer encounters work outside its specialty, its professional instinct is to seek a collaborator rather than stretch. This keeps the context narrow naturally — not through restriction, but through the craftsman's recognition that a specialist produces better work.

### The prescriptive layer

Jobs are where all prescriptive wisdom lives that currently erodes. `[B:artifex]`

| What erodes today | Where it lives with jobs |
|---|---|
| "Load /coding before writing code" | Code-implementation job: `required_skills` |
| "Run make check before discharge" | Code-implementation job: `discharge_gates` |
| "Apply dialectical calcinatio on all outputs" | QA job: `workflow` |
| "Don't skip browser observation" | UI-QA job: `verification_requirements` |
| "Update the spec after findings" | Spec-writer job: `workflow` |

---

## Discharge as Relay Handoff

The discharge process changes to maintain momentum. Currently, azer.md treats follow-up inscription as optional ("You may inscribe follow-up opera"). With jobs, inscription becomes a structural step in the discharge sequence — executed while the azer has peak context about what needs to happen next. `[B:artifex]`

**Revised discharge sequence:**

1. **Define discharge.** Draft the outcome, evidence, and reflection. You now know what you accomplished and what's left.
2. **Inscribe + muster immediate follow-up.** While you have peak context: inscribe opera for any immediate follow-on work, assign job roles, and muster the azers. The next wave is already booting by the time you finish. Use lightweight inscription (single fresh-context subagent, one pass — the job definition provides the structure).
3. **Write discharge.** Finalize the discharge record, including what you inscribed. "Follow-up: inscribed QA opus for auth refactor, mustered qa-specialist."

**Why this order matters:**
- **Peak context:** At discharge time you know better than anyone what needs to happen next. Waiting for marut to rediscover this wastes that knowledge.
- **Continuous momentum:** The next azer starts working while you're still writing the discharge record. No gap waiting for marut to notice, assess, and inscribe.
- **Clean trail:** The discharge record references what was inscribed, creating an explicit chain in the trail.

**Not every discharge inscribes.** If there's no obvious follow-up, or if the next step needs marut's judgment, just discharge normally. The marut's monitoring loop catches the tapering and re-energizes. The point is: when follow-up IS obvious, don't defer it — inscribe now while the knowledge is fresh. `[D:discharge-inscription-optional-but-structural]`

---

## Workflow / Dataflow

### The decision point

At inscription time, the inscriber (marut, artifex, or another azer) decides whether the opus needs a job-role azer or a general-purpose azer. `[B:artifex]`

```
Work identified
  │
  ├─ Defined job fits? ──→ Inscribe with job role ──→ Muster job-role azer
  │                                                     (boots with job prescriptions)
  │
  └─ No defined job ────→ Inscribe without job ───→ Muster general-purpose azer
                                                     (today's model)
```

The job is specified in the opus frontmatter:
```yaml
---
status: charged
job: qa-specialist      # or omitted for general-purpose
magnum_opus: some-mo
---
```

`[D:opus-frontmatter — lightweight, discoverable, consistent with existing YAML pattern]`

### Lightweight inscription

With frequent inscription as the norm, the inscription process must be lightweight. `[B:artifex]`

**Current model:** Dialectical calcinatio with an inscription subagent — iterative exchange to shape a well-formed opus. Appropriate for complex, ambiguous work.

**New model:** A single fresh-context subagent inscribes in one pass. The job definition carries the structural weight (prescriptions, workflow, gates), so the opus is just a briefing — what to look at and why. `[B:artifex]`

For a job-role opus, inscription might be as simple as:

```yaml
---
status: charged
job: qa-specialist
magnum_opus: some-mo
inscribed: 2026-04-21
---
# Review authentication refactor

The auth token refresh logic was rewritten in `internal/auth/`.
Tests pass but no independent review has been done.
Key concern: token expiration edge cases under concurrent requests.
```

The job fills in everything else at boot. The opus says WHAT to look at; the job says HOW to look.

For general-purpose opera (no job), the inscription can be heavier — these are the novel, ambiguous cases where more upfront shaping pays off. But even here, the dialectical process is optional, not mandatory. The marut and individual azer can work things out at runtime. `[B:artifex]`

### Subagents vs. peer azers

Not all collaboration requires inscription. The distinction: `[D:collaboration-taxonomy]`

| | Subagent (Agent tool) | Peer azer (inscription) |
|---|---|---|
| **Ownership** | Part of MY opus | Their own opus |
| **Trail** | Invisible — value in my discharge | Trail-worthy — independent discharge |
| **Context** | Bounded exchange within my context | Own fresh context with job prescriptions |
| **Lifecycle** | Ephemeral — fires and returns | Full opus lifecycle |
| **Use when** | Bounded exchanges (dialectical calcinatio, code search, quick research, file exploration) | Independent work products that benefit from job prescriptions, context separation, or trail visibility |

**Heuristic:** If the work is part of your opus and the value flows back into your discharge, use a subagent. If the work is independently valuable, benefits from its own prescriptions, or needs context separation from your builder context, inscribe a peer.

### Dynamic orchestration — the organic model

Work flows organically, with the marut as backstop when momentum drops. `[B:artifex]`

**A day in the life:**

1. **Marut kick-starts.** Morning — the marut musters an assessor azer. `[B:artifex]`
2. **Assessor generates initial work.** The assessor surveys the landscape, produces a report, and directly musters a couple of azers for the most immediate actions. The assessor discharges.
3. **Work self-propagates.** Those azers work, and in the course of their work inscribe and muster collaborators — a coder inscribes a QA specialist, a researcher inscribes a spec-writer to capture findings. Work flows organically for a while.
4. **Work tapers off.** Azers discharge, no new work is being generated. The system is quiet.
5. **Marut notices and re-energizes.** The marut sees the tapering — either musters a few azers directly for obvious next steps, or musters another assessor to survey the new landscape and generate the next wave.
6. **Cycle repeats.** Organic self-organizing phases punctuated by marut interventions when momentum drops.

**Key properties:**

- **Any azer can muster directly.** No bottleneck through marut. Azers inscribe opera and run `ath muster` themselves. The marut sees the new crucible appear and adds it to monitoring. `[B:artifex]`
- **Work generation is distributed.** Assessor generates initial work, azers generate follow-on work, marut is the backstop — not the sole generator of all work.
- **Marut's primary trigger is "momentum dropped."** Not a scheduled cadence (though landscape freshness thresholds still apply), but noticing when the system is quiet and needs re-energizing.
- **The marut's unique value is its persistent perspective.** It sees across all active azers, across the trail, across time. Individual azers come and go; the marut sees the whole picture and intervenes when things need steering.

```
Marut (persistent, monitoring)
  │
  ├─ Kick-start ──→ Muster assessor
  │                    │
  │                    ├─ Assessor musters azer-A, azer-B
  │                    │    │
  │                    │    ├─ azer-A musters azer-C (QA)
  │                    │    ├─ azer-B musters azer-D (research)
  │                    │    └─ ... organic propagation ...
  │                    │
  │                    └─ Assessor discharges
  │
  ├─ Monitor ──→ Nudge, suggest collaboration, flag stalls
  │
  └─ Momentum drops ──→ Re-energize: muster more azers or muster assessor
```

### Context lifecycle

```
Boot (~0 tokens)
  │  Load job definition, required skills, opus
  │
Working (~0–100k tokens)
  │  Narrow, focused work within job boundary
  │
Discharge signal (~100k tokens)
  │  Job prescriptions still salient — verification gates hit
  │
Discharge
  │  Hand off to next job-role azer or back to marut
```

`[B:artifex — ~100k token discharge target]`

---

## Job Definition Structure

Evolved from the initial Gastown TOML concept toward **perspective-first markdown**. Jobs are about identity, not process — the perspective naturally produces the right behavior. `[B:artifex]` `[E:gastown]`

Each job is a directory in `jobs/` with a `JOB.md` file (mirroring how skills use `SKILL.md`). The markdown IS the prompt injection — it's read directly into the azer's context at boot.

### Job file structure

A `JOB.md` has four sections:

| Section | Purpose |
|---------|---------|
| **Identity** (opening paragraph) | Who you are, what your craft is — the professional perspective |
| **When this role is needed** | System-level guidance for when to invoke this job — helps marut and other azers know when to inscribe it |
| **What you care about** | Your concerns, your lens — what you evaluate, what matters to you |
| **Your tools** | `/skill-discovery` + pointers to relevant skills for this perspective |
| **Your instinct** | Professional instinct — boundary behavior in Chesed framing (when to collaborate vs. do it yourself) |

**No workflow steps, no gates, no rigid boundary rules.** The perspective drives the behavior. A QA specialist that genuinely identifies as "my craft is preventing surprises" will naturally verify before discharging. Gates are process on top of identity — and process is what drifts.

**Exception:** The assessor job includes a `## Your process` section because the assessment formula has a deliberate structure (survey → generate → synthesize) that adds genuine value. This is the exception, not the pattern. `[D:assessor-process-exception]`

See `shared/jobs/qa-specialist/JOB.md` for a representative example.

---

## How Jobs Interact with Existing System

### Jobs and the geas

The geas remains supreme. A job constrains perspective but not agency — if a QA azer's geas says "the most valuable thing right now is to fix this critical bug myself", the geas wins. But the job boundary creates a strong default: "I'm a QA specialist, fixing is out of scope, I should spawn a fixer." The job makes the right behavior the easy behavior. `[D:geas-supremacy]`

### Jobs and opera

The opus carries the job assignment in frontmatter. `ath inscribe` enforces `--job` at inscription time (validated against the registry). The mustered azer reads its opus, sees the `job:` field, and reads `jobs/<job>/JOB.md`. `[I:implemented]`

**Gap (2026-04-27):** `ath muster` does not yet inject the job definition path into the boot prompt — loading is behavioral (azer.md instruction), not structural. An opus is inscribed to fix this.

### Jobs and marut — the shepherd model

Marut shifts from dispatcher to shepherd/facilitator. `[B:artifex]`

| | Marut today | Marut with jobs |
|---|---|---|
| **Primary activity** | Assess → inscribe → muster (drives the loop) | Monitor → nudge → facilitate (keeps work flowing) |
| **Work generation** | Marut generates all opera through assessment | Azers inscribe peers; marut inscribes when it sees gaps |
| **Coordination** | Serial — one assessment cycle at a time | Parallel — multiple azers self-organizing, marut watches the whole |
| **Interventions** | Scheduled (landscape freshness threshold) | Responsive — notices overlap, suggests collaboration, flags stalls |

**The marut's operational loop:**

1. **Observe** — what are all active azers doing? Tmux observation, whisper monitoring, trail reads.
2. **Notice** — is anything stuck? Are two azers doing overlapping work? Did someone finish implementation without inscribing QA? Is context getting full somewhere?
3. **Nudge** — whisper suggestions to azers: "you two are working on overlapping areas, consider collaborating", "your implementation looks done, consider inscribing a qa-specialist", "you're at 80k tokens, start thinking about discharge". `[B:artifex]`
4. **Inscribe when needed** — the marut still inscribes directly when it sees gaps nobody else has noticed, or when an azer escalates and needs help figuring out what to do next.

The marut doesn't order — it suggests. Azers retain agency under the geas. But the marut sees the whole picture that no individual azer can see, and that perspective is what makes its nudges valuable.

**Assessment is a job.** The formal assessment protocol (currently a large section in AGENTS.md) becomes the `assessor` job definition. When the marut decides a landscape survey is needed (staleness, no active azers, significant trail accumulation), it inscribes an assessment opus with `job: assessor`. The assessor azer boots with prescriptions for manifold generation, landscape reports, and opus inscription — all the heavyweight assessment machinery. The marut doesn't do assessment itself; it decides *when* assessment is needed and inscribes one. `[B:artifex]` `[D:assessment-as-job]`

**The marut remains the storyteller.** Offloading assessment doesn't make the marut leaner — with more specialized azers to monitor, more collaboration to facilitate, and a richer narrative to maintain, the marut may get *heavier*. It sends dispatches to the artifex at meaningful turns, maintains awareness across all active azers, and holds the big-picture context that no individual azer can see. The marut is a full presence-driven singleton with a full context window. `[B:artifex]`

### Jobs and skills

Jobs specify required skills — this is the key difference from today's advisory model. `[B:artifex]`
- **Today:** "Run /skill-discovery" → azer loads skills → skills drift
- **With jobs:** Job definition says `required_skills = ["/coding", "/unit-testing"]` → loaded at boot → survive in tight context

Skills remain independent — jobs compose them. A job is not a super-skill; it's an identity that prescribes which skills to load and how to use them. `[D:jobs-compose-skills]`

### Jobs and the Tree of Work

Jobs live at the **Yesod level** — foundation, the bridge between the environment (Malkuth) and the opus dyad (Netzach/Hod). They're the structural mechanism that ensures the environment's prescriptions (skills, workflow, verification) actually reach the working agent. `[D:tree-placement — validate through use]`

---

## The Job Catalog

### Implemented `[I:implemented]`

Seven jobs are defined in `jobs/` and symlinked into each athanor instance as `jobs/`:

| Job | Scope | Key prescriptions |
|-----|-------|-------------------|
| `coder` | Write, modify, and review code | Shift-left quality, `/code-review`, `/unit-testing`, builder bias awareness |
| `qa-specialist` | Review and verify work products | `/calcinatio`, empirical verification, no-surprises lens |
| `assessor` | Landscape survey and opus inscription | Structured process (survey → generate → synthesize → inscribe), manifold generation |
| `general` | Novel, exploratory, cross-cutting work | Breadth-first perspective, loneliness signal for collaboration, scout not fallback |
| `investigator` | Trace problems to root cause | Evidence over narrative, `/rca`, communication closes the loop |
| `solution-architect` | Periodic holistic eye on MO architecture | Spec-reality alignment, architectural fitness, `/spec-check`, `/explore`, proposal-ready findings |
| `project-manager` | Witness satisfaction audit and progress assessment | Witness coverage, hill-chart status, tunnel vision detection, `/project-management` |

### Candidate — from voice notes `[V:voice-note]`

The artifex identified recurring work patterns that suggest new job definitions (2026-04-24):

| Job | Scope | Origin signal |
|-----|-------|---------------|
| `environment-warden` | Local environment health — when an azer is stuck debugging env issues (missing deps, broken tools, config drift), this specialist takes over instead of the coder going in circles. | Recording 182 (2026-04-24) |

### Candidate — from original spec `[U:needs-validation]`

| Job | Scope | Key prescriptions |
|-----|-------|-------------------|
| `spec-writer` | Shape specifications | `/spec`, brownfield discovery, viewpoint review |
| `content-editor` | Edit documents, specs, role files | Style consistency, completeness, witness clarity |
| `researcher` | Deep investigation with citations | `/research`, source verification, citation format |
| `coordinator` | Facilitate multi-azer collaboration | Inter-azer protocol, status synthesis, bridge comms |

---

## Dynamic Spawning Model

Resolved through the organic model (see § Dynamic orchestration). Summary of answers: `[B:artifex]`

- **Any azer can spawn any other azer.** Direct `ath muster` — no bottleneck through marut. The marut monitors and intervenes when needed, not gates.
- **Context is communicated via the opus.** The inscribing azer writes the opus (lightweight for job-role opera — the job provides structure). The mustered peer reads its opus at boot like any azer.
- **Spawned peers are independent.** They have their own opus, their own discharge, their own trail entry. They may collaborate with the spawner via whisper, but they don't "report to" the spawner — they report to the trail and the MO.
- **The marut sees everything.** New crucibles appear, the marut adds them to monitoring. No explicit registration needed — the marut observes tmux.

---

## Implementation Status

### Wave 1: The foundation — COMPLETE `[I:implemented]`

**1. Job registry.** `jobs/` created with 5 job definitions: `coder`, `qa-specialist`, `assessor`, `general`, `investigator`. Symlinked as `jobs/` into each athanor instance (athanor-architect, musashi, maneframe) as of 2026-04-27.

**2. azer.md updated.** Collaboration principle (Chesed framing), discharge relay sequence, job awareness (`job:` field → read `jobs/<name>/JOB.md`), ~100k context threshold, subagent vs. peer azer guidance — all deployed.

**3. AGENTS.md updated.** Assessment section slimmed (points to assessor job), collaboration principle in core geas, job/job-registry vocabulary added.

**4. marut.md updated.** Shepherd model operational loop (observe → notice → nudge → inscribe), assessment delegated to assessor job, multi-azer monitoring guidance.

**5. CLI enforcement.** `ath inscribe` requires `--job` flag (validated against registry). `ath collaborate` likewise. Opus frontmatter carries `job:` field.

**6. Job definition path references.** Agent directives (AGENTS.md, azer.md) updated to reference `jobs/<name>/JOB.md` (relative to athanor instance dir, where the symlink lives). Previously said `jobs/` which was unreachable from instance dir.

### Wave 2: Observation and iteration — IN PROGRESS

The system is running with these changes. Key observations and open work:

**Known gaps (opera inscribed 2026-04-27):**

- **Structural job injection in muster.** `ath muster` does NOT inject the job definition into the boot prompt — job loading is purely behavioral (azer.md tells the azer to read it). An opus is inscribed to make muster read the `job:` field from opus frontmatter and add the JOB.md path to the boot prompt structurally.
- **Idempotent ath init/sync.** `ath init` is create-only. `SharedFiles` in `home.go` is stale (missing `attendant.md`, `perceiver.md`) and has no concept of shared directories (like `jobs/`). An opus is inscribed for an idempotent sync mechanism.

**Observation questions (from original spec, still open):**

- Do azers actually inscribe peers? Or do they still try to do everything themselves?
- Do job prescriptions survive the context? (Core hypothesis — untested at scale)
- Does the marut's shepherd loop work? Or does it need more structure?
- Does the discharge relay create momentum? Or do gaps still appear?

### Wave 3: Marut as active coordinator `[V:voice-note]`

The artifex's voice notes (recording 182, 2026-04-24) signal a potential shift in the collaboration model. The current spec assumes azers will spontaneously inscribe collaborators (Chesed pull). Observation suggests this is the weakest link — azers resist creating collaborators and try to do everything themselves.

**The proposed evolution:** Rather than relying on azers to self-collaborate, lean into the marut as the coordination point:

- **Role discovery.** The marut does something analogous to `/skill-discovery` but for jobs — it knows the full job registry and actively matches work-in-progress to the right specialist roles.
- **Active dispatch.** The marut watches azers and catches when they're doing work that another job would do better (e.g., a coder debugging local env issues → inscribe an environment-warden). Not preventing them from going outside their lane, but noticing when a specialist would be more effective.
- **Generative collaboration.** The ideal: marut catches a coder troubleshooting local env issues, inscribes an environment-warden collaborator, whispers the coder to stop and work with the specialist instead.
- **Shorter marut context.** More aggressive context cycling for the marut itself, so it gets fresh perspective on the workshop state more frequently.

This is upstream of formulae — a stepping stone to see how much coordination emerges from a marut with job-registry awareness before formalizing into formulae. `[V:voice-note]`

**Model selection.** The artifex notes (recording 182-(1), 2026-04-24) that marut should be opus-class but many job-role azers (especially implementation roles) could run on sonnet. Comparison testing needed, but the direction is clear: reserve opus for roles that need deep judgment (marut, assessor), use sonnet for focused execution (coder, qa-specialist). `[V:voice-note]`

### Wave 4: Periodic/patrol roles `[V:voice-note]`

The artifex (recording 175-(5), 2026-04-23) describes a concept of **patrols** — time-based agent dispatches defined per MO with intent + job role. This extends jobs into scheduled work (Slack channel monitoring, environment health checks, progress audits). See separate patrols/functions concept — not yet specced but the job system provides the foundation (patrol = scheduled muster with a job role).

---

## Test Scenarios

Behavioral scenarios to observe once the system runs: `[D:observe-in-practice]`

1. **Job prescriptions survive.** A coder azer loads `/code-review` and `/unit-testing` at boot and still follows their guidance at discharge time (~100k tokens later). Evidence: discharge record shows verification gates hit.
2. **Collaboration happens.** An azer inscribes a peer with a job role instead of doing everything itself. Evidence: trail shows multiple small specialized opera rather than one large general one.
3. **Discharge relay works.** An azer inscribes follow-up opera before writing its discharge, and the next azer is already running when the discharge is written. Evidence: timestamps in trail show overlap.
4. **Marut shepherds effectively.** The marut notices when momentum drops and re-energizes. Evidence: trail doesn't show long gaps between waves of work.
5. **General-purpose fallback works.** When no job fits, a general-purpose azer is mustered and operates as today's model. Evidence: novel/exploratory work still proceeds without a defined job.
6. **Context stays narrow.** Job-role azers discharge around ~100k tokens, not 200k+. Evidence: context percentages in discharge records.

---

## Out of Scope

- CLI implementation details for job-aware muster (`ath` changes belong in CLI spec)
- Specific job definition contents (each job gets its own definition file, refined through use)
- Changes to perceiver/attendant roles (jobs are for opus-scoped agents)

---

## Retrospective

### What shipped vs. spec

- **Job definition structure held up.** The four-section format (identity, when needed, what you care about, your instinct) plus `## Your tools` works well. The assessor exception (`## Your process`) was warranted — confirmed by use. `[O:observed]`
- **Two unplanned jobs emerged.** `general` was elevated from "fallback" to a real job with its own identity (breadth-first scout, loneliness signal). `investigator` emerged from a `general` opus that revealed a distinct diagnostic perspective worth naming. Both validated the "jobs grow from observed patterns" design. `[O:observed]`
- **`ath inscribe --job` enforcement works.** Every opus gets a job assignment. `general` serves the "no specific role fits" case cleanly. `[I:implemented]`
- **Symlink distribution was a gap.** Job definitions lived in the source repo but weren't reachable from athanor instance dirs. Fixed 2026-04-27 with `jobs/` symlink, but exposed the broader problem: no idempotent sync mechanism for shared components.

### Open questions

- [ ] Which jobs proved most valuable? Which weren't used? (Need more production trail data)
- [ ] Did the ~100k context limit work? Too tight? Too loose?
- [ ] Does the marut-as-coordinator model (Wave 3) reduce the collaboration gap, or does it need formulae?
- [ ] Does structural job injection in muster boot prompt improve prescription survival vs. behavioral loading?
- [ ] What's the right model tier per job? (opus for marut/assessor, sonnet for coder/qa — needs testing)

### Patterns for writing good job definitions

Emerging from the five implemented definitions:

- **Identity paragraph is load-bearing.** "Your craft is X" frames everything. If this paragraph is generic, the whole job is generic.
- **"When this role is needed" helps the inscriber, not the azer.** It's system-level routing guidance — the marut or peer azer reads it to decide when to inscribe this job.
- **"Your instinct" is the collaboration hinge.** This is where you define the boundary behavior — when to do it yourself vs. inscribe a specialist. Chesed framing (generous, not restrictive) works better than rules.
- **Avoid procedures.** The assessor's process section is the exception. If you're writing steps, you're writing a skill, not a job.
