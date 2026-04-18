# The Athanor Adaptation

*What transfers from the Athanor spec to the life-layer domain, what needs new thinking, and how the life-layer relates to the existing athanor infrastructure.*

This is bridging material between the Athanor spec and the life-layer architecture. It assumes familiarity with both.

---

## What transfers directly

The following Athanor primitives map onto the life-layer domain essentially unchanged. The philosophical stance and mechanical discipline that make them work in the coding domain make them work here.

**Magnum Opus structure.** The MO as a top-level goal container, with intent, witnesses, calcinatio, and abundant satisfaction as load-bearing sections. Life-layer MOs have the same shape.

**Witnesses as experience-definers.** Intent expressed as change in witness experience, not as technical specification. The life-layer's primary witness is the artifex; witness concerns drive calcinatio fires and satisfaction criteria.

**Geas as behavioral compulsion.** The orientation that produces rigorous work through internalized felt obligation rather than through rule-following. The life-layer attendant's geas carries this discipline, with content adapted to the prosthesis framing.

**Calcinatio as refining fires.** Concerns the work must survive. Structurally identical; specific fires differ (see below). The double-fire structure for the life-layer (too-much, too-little) is an elaboration within this framework, not a departure.

**The trail.** Discharged opera as the persistent record, sessions as ephemeral, tinctures as concentrated extracts. The "waste nothing" discipline. The trail as the primary evidence of work and the ground truth against which the marut assesses MO health.

**The lightning path principle.** Success as the best next step from where you are. No special claim for interrupted work. Sunk cost rejected. This is philosophically central to the life-layer and addresses the ADHD-specific pattern where the sunk-cost-of-unfinished-projects becomes its own cognitive burden.

**The marut as supervisor.** Relentlessly ensuring the goal advances, watching agent health and progress, narrating dispatches to the artifex. The marut's shape works for life-layer MOs the same way it works for coding ones. What the marut watches for is different (attendant health, calibration health, coverage across the topography of support) but the role is the same.

**The homunculus as the artifex's familiar.** Direct interactive assistance, non-autonomous. Works identically in the life domain — when the artifex wants direct conversational engagement rather than the autonomous system, the homunculus is that surface.

**The workshop spectrum.** The artifex walking in and out freely, with autonomy as orthogonal to role. The life-layer attendant operates autonomously when the artifex is away and collaboratively when they're at the bench — same agent, same geas, different autonomy level.

**The assay pattern.** Confirming materia before acting, no assumptions, surface gaps as new opera or escalations. The life-layer assay asks "do I have the context and access to attend to this domain well? Is the portrait fresh enough? Are recent voice notes and sensor events current?"

**The principle that opera must be actionable now.** No draft opera, no queued opera. Life-layer opera inscribe only when action (or bracketed non-action) is warranted. Most of the time, nothing is inscribed because nothing warrants a bracketed act.

**Opus structure (intent/boundary/context).** The three-part shape scales well to life-layer opera, including small ones.

**"Read once, write at discharge."** The opus is a commitment-and-record artifact, not a scratchpad. Works the same way for attendant opera.

## What needs adaptation

The following Athanor primitives work but need adjustment for the life-layer domain.

### Azer replaced by attendant

Azer is a craftsperson. The metaphor carries a specific posture (skilled hands, focused execution, producing artifacts) that doesn't fit the attendant's work.

The attendant is the new agent class for life-layer work. See `02-system-architecture.md` for the full treatment. Key differences from azer:

- Posture is witness/steward/companion rather than craftsperson.
- Output is often interpretive rather than productive.
- Completion bias manifests as "done for now" too early rather than "declared done without evidence."
- Geas orients toward quality of attention rather than quality of artifact.
- Work rhythms are different (vigils span time, regards are continuous with imposed boundaries, integrations metabolize material).

The attendant is one class, not three. Variations in work shape don't require separate classes. This matches the spec's principle that autonomy is orthogonal to role.

### Units of work

Azer opera naturally scope to tasks with production outputs (a PR, a spec, a feature). Attendant opera don't have those natural boundaries — the work is sustained attention, not production.

The three shapes (vigil, regard, integration) give attendant work its scoping structure:

- **Vigils** scope by the lifetime of what's watched — beginning when the watching becomes warranted, ending when the watched thing resolves, transforms, or is retired.
- **Regards** scope by shifts — periods of sustained orientation, bounded by meaningful landscape shifts, day-structure boundaries, or accumulated material thresholds. Shifts are the imposed-boundary units for work that is inherently continuous.
- **Integrations** scope by the material being metabolized — a bounded piece of work processing accumulated context into refined understanding.

**Regards are the hardest scoping problem.** Continuous work with imposed boundaries is where the on-the-clock completion bias manifests. The guidance in `02-system-architecture.md` is explicit but provisional; real scoping discipline emerges through use.

### Opera scale freely

The spec treats opera as substantial bracketed units. Life-layer attendant work produces opera across a wider range — from substantial interventions down to minute-level judgment calls.

The position taken: rather than introducing a separate "micro opus" concept, let the opus concept flex. The discipline (inscription, assay, execution, discharge, trail) scales naturally; the ceremony (length, depth of evidence, communication requirements) scales with the weight of the work.

The floor for "opus" is the unit of judgment, not the unit of execution. A decision that ambient signal should shift is an opus; the tool call that changes the light is execution within the opus.

This produces a more granular trail in attention-focused domains than in production-focused ones, which is fine. More granularity means more signal for marut assessment.

### Calcinatio fires

The spec's calcinatio is "what fires the work must survive" — concerns materialized as test suites, CI, witness review. These are the environment-resident fires.

The life-layer's primary calcinatio mechanism is subagent critique — the two-fire pattern (too-much, too-little) applied to proposed moves. This is calcinatio as live critique from differentiated perspectives rather than as environmental gates.

Both mechanisms can coexist. The spec's approach (environmental fires for substantive coding work) remains valid for any work with natural testable outputs. The life-layer's approach (critique-pass fires for judgment about communication and action) serves the domain where most work is judgment without testable outputs.

### Generative panel extension

The spec's azer does reasoning solo. The life-layer attendant reasons through a panel of lens-agents representing different cognitive support domains, each asking "what support can I provide right now?" The attendant synthesizes a proposed move from the panel's outputs.

This pattern may be worth extending back into the azer's coding domain — a panel of stakeholder and technical perspectives generating candidates, an azer synthesizing — as a general improvement to how substantive Athanor work is done. The life-layer may be forcing the issue, but the pattern has domain-independent value.

## What doesn't map cleanly

A few spec elements don't have obvious life-layer analogs.

**The beholder.** The spec describes beholders as watchers that scan channels and conditions, creating work when they find something that needs doing. For the life-layer, the periodic landscape-assessment work that a beholder might do is folded into the attendant's regard shape. The attendant *is* watching continuously; dedicated beholders may not be a separate role in this domain.

Or: sensors (pill case switches, calendar integration, voice note pipelines) act as beholders — they generate events that inform the attendant's reasoning but don't themselves inscribe opera. The attendant's regard reads these events and decides what warrants becoming an opus.

**Environmental fires.** CI pipelines, test suites, explicit quality gates. These are the coding domain's standard calcinatio infrastructure. The life-layer domain has analogs (sensor evidence for medication resolution, trail coherence checks, portrait consistency reviews) but they're not the same shape. Subagent critique fires do more of the work that environmental fires do in coding.

**PRs and similar discrete deliverables.** The coding domain's natural unit boundaries. The life-layer doesn't have an equivalent. This is why unit scoping is a harder design problem and why regard shifts require explicit boundary discipline.

## Where the life-layer sits relative to the existing athanor

The natural implementation path is running the life-layer as one or more athanor instances. The spec's primitives, the trail infrastructure, the marut supervision, the workshop model all apply directly.

Rather than one "life" athanor with one MO, the structure that seems right is multiple MOs, each addressing a distinct aspect of the prosthesis work. Three natural ones:

1. **Artifex attention is protected and cognitive load is offloaded.** State-based. The general attention-protection and executive-function support work.
2. **Hard commitments are met with evidence-based resolution.** State-based. Medication, appointments, time-anchored obligations with physical-world evidence closing loops.
3. **Interoceptive practice is supported without being managed.** State-based. The voice-note substrate, the remembered-not-surveilled discipline, the rhythm-protection concerns.

Each MO has its own witnesses (with the artifex primary in all three, but other witnesses varying — e.g., medical providers for hard commitments), its own calcinatio fires, its own satisfaction criteria. Each gets its own marut. Opera under each MO serve that MO's specific intent.

These MOs share the attendant class and the panel/synthesis/calcinatio reasoning patterns, but their specific work is differentiated by what their MO defines.

Alternatively: a single life-layer MO with broader intent, addressing all three concerns within one structure. Probably worse — the calcinatio fires for attention protection differ from those for hard-commitment completion, and mixing them loses discrimination. Separate MOs keep the concerns properly separated while the shared infrastructure (attendant, trail, tools) handles the actual work.

Whether to run these in the same athanor instance as the coding work or in separate instances is an open question. Arguments for separate: different domains, different attention patterns, different MO cadences. Arguments for same: shared attendant capabilities, shared infrastructure, no coordination overhead. Default expectation: separate athanor instance for life-layer work, with the possibility of eventual integration if the instances need to be aware of each other.

## Spec changes that might follow

If the life-layer work matures and proves out some of its patterns, parts of it may warrant feeding back into the main spec:

**The attendant class** as an alternative to (or alongside) the azer. Not all work is crafting; some is attending. The spec could name both agent classes with clear guidance about which applies to which kinds of work.

**The generative panel pattern** as a first-class reasoning structure for substantive opera. Could complement the existing solo-azer approach as an option for work that benefits from multi-perspective generation.

**Opera scaling guidance** to make explicit that opera can range from substantial to minute-level as long as the bracketing discipline is preserved. The spec implicitly allows this but doesn't call it out.

**The system-as-stakeholder pattern** in MO writing. Applies beyond the life-layer — any Athanor MO has system-stakeholder concerns (trail integrity, architectural integrity, calibration maintenance) that could be named rather than left implicit.

**The double-fire calcinatio pattern** (too-much / too-little) for decisions that can fail in both directions. Probably relevant to coding decisions too, especially around communication with witnesses — when to update, when to stay silent, when to escalate.

None of these are urgent. The life-layer should prove out on its own first before generalizing back.

## One coordination question worth flagging

If the life-layer runs as its own athanor, there's an eventual question of whether the coding-work athanor and the life-layer athanor should know about each other. Examples where they might:

- The coding athanor's marut notices work is stuck; could inform the life-layer to relax pressure if the artifex's state suggests overload.
- The life-layer attendant notices the artifex has been away from a specific coding project for a long time; could surface it gently at a good moment.
- The life-layer's regard on attention patterns could inform the coding athanor about when to escalate versus hold.

These are appealing but not urgent. Premature integration would couple the systems in ways that constrain each of them. Better to prove out both independently first; integration concerns arise only when there's concrete value that can't be achieved otherwise.

## Summary

The life-layer inherits most of Athanor's architecture intact. It adds:

- The attendant as a new agent class (alongside azer).
- The three work shapes (vigil, regard, integration) as attendant-specific work patterns.
- Opera scaling explicitly down to judgment-level granularity.
- The double-fire calcinatio pattern.
- The generative panel as the attendant's primary reasoning structure.
- The system-as-stakeholder framing alongside the artifex as primary witness.

It relies on Athanor's existing:

- MO structure and witness/calcinatio/satisfaction framework.
- Geas as behavioral compulsion.
- Trail discipline and "waste nothing" commitment.
- Lightning-path principle and rejection of sunk cost.
- Marut supervision and workshop model.
- Homunculus for direct artifex interaction.

The implementation approach: separate athanor instance for life-layer work, with multiple MOs addressing distinct prosthesis concerns, all served by attendant-class agents operating under adapted calcinatio patterns. Prove the architecture through use; generalize patterns back to the spec only when they've matured.
