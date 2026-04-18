# Orientation Layer Architecture — Exploration

*A freewheeling exploration of how the orientation layer concretely manifests, given the athanor's actual architecture. Written as Chesed input for later Geburah refinement.*

---

## The Problem in Sharp Relief

The athanor works for software projects because its primitives are a near-perfect fit. The MO is the project goal. The azer is a craftsperson whose product is an artifact — code, spec, PR. Witnesses are humans with professional concerns about that artifact. Calcinatio fires are literally available: tests, CI, code review, deploy pipelines. The geas is "take the most valuable next step toward shipping this thing." The system as a whole is coherent because all its elements speak the same language.

The life domain speaks a different language. The "product" is often a shift in experience, not an artifact. Witnesses are primarily the artifex themselves — and their concerns are about lived quality, not deliverable quality. Calcinatio fires are not CI pipelines; they are subagent critiques asking "was that move too much or too little?" The geas must include "silence is sometimes the highest-value action" as a first-class option, not an edge case.

The orientation layer is the translation layer between these two languages. It does not change the athanor's grammar; it provides vocabulary that makes the grammar mean something in the life domain.

---

## The Constraint That Shapes Everything

The orientation layer CANNOT require core changes. This is a hard constraint, not a preference. The shared agent files — AGENTS.md, azer.md, marut.md, opus.md, muster.md — are symlinked into every athanor instance. A change to azer.md ripples into all coding athanors, all other life athanors, every active agent session. Changes there require deliberate review and are outside scope for the initial build.

This means the orientation layer must work entirely through mechanisms that compose with the existing architecture:
- Per-instance files (agent role overrides, MO templates, config)
- Skills (loaded at session start, shape behavior)
- The MO itself (shapes witness concerns, calcinatio fires, satisfaction criteria)
- Muster parameters (what gets passed when an agent session is kindled)

Everything below operates within this envelope.

---

## Mechanism 1: The `attendant.md` Role File

**What it is:** A role file that lives in the life athanor's directory alongside the symlinked `azer.md` and `marut.md`. The muster command passes it as the agent's role file instead of `azer.md`.

**How it works:** `ath muster` already accepts a role file path or can be told which role variant to load. The attendant is not a replacement for the azer — it's a different face of the same underlying capability. The attendant.md would contain:

- The attendant's geas (prosthesis framing, silence as default, quality of attention not quality of artifact)
- The three work shapes (vigil, regard, integration) as first-class patterns, where the azer's concept of "task" lives
- Calcinatio adapted: the double-fire (too-much/too-little) as the primary fire pattern, with subagent critique rather than environmental test suites
- The panel reasoning pattern: stage 1 (lens-agents), stage 2 (synthesis), stage 3 (double-fire calcinatio)
- Discharge patterns calibrated for interpretive content rather than proof-of-artifact

**What this does NOT change:** AGENTS.md (core geas, vocabulary, tree of work), opus mechanics, marut behavior, CLI operations. The attendant inherits all of this. It re-orients the posture while keeping the architecture intact.

**Core changes required?** No. This is a new file in the life athanor instance directory. The existing muster protocol supports role file selection.

**Behavioral shift:** Substantial. The attendant's self-understanding is fundamentally different from the azer's. The geas lands differently. The work shapes frame how opera are scoped. The calcinatio pattern changes what "verification" means. This is the single highest-leverage mechanism in the orientation layer.

**Implementation effort:** Medium-high. Writing `attendant.md` well is the hard work — it must be as carefully crafted as `azer.md`, tuned through use, and live with the quality of the azer's behavioral guidance. This is not a configuration task; it's a writing task that requires as much attention as the original role files.

---

## Mechanism 2: Life-Oriented MO Templates

**What it is:** A template (or family of templates) for life-domain Magna Opera that differs structurally from the work-domain MO template.

**The structural differences:**

The work MO looks like:
```markdown
## Intent
Ship X feature / fix Y bug / produce Z artifact

## Witnesses
- CTO (concerned with architecture quality)
- Users (concerned with UX)

## Calcinatio
- All tests pass
- Security review by CTO proxy
- Field QA pass

## Abundant Satisfaction
Witnesses would feel: feature is shipped, visible, reviewed, documented
```

A life MO template looks like:
```markdown
## Intent
[State-based, not deliverable-based]
e.g.: "Hard commitments are met with physical-world evidence. 
The artifex reliably knows whether they've taken their medication 
without relying on memory."

## Witnesses
- The artifex (primary — experience witness, not deliverable reviewer)
  Concerns: non-intrusive support, prosthesis working without friction
- The system-as-stakeholder
  Concerns: trail integrity, portrait currency, calibration maintenance
- [Domain-specific: medical provider, family member, etc.]

## Calcinatio
- Too-much fire: proposed external actions pass the attention-cost test
- Too-little fire: proposed silences are genuine discretion, not abdication
- Portrait-consistency check: does this move cohere with accumulated understanding?
- Calibration audit: is the system drifting from its actual pattern in this domain?

## Landscape Freshness Threshold
4 hours (life domains shift more rapidly than code domains)

## Abundant Satisfaction
The artifex would feel: this domain is attended without being managed. 
EF load in this area is reduced without new cognitive burden. 
The system's presence is helpful when engaged and invisible when not.
```

**What this produces:** The MO template shapes the assessment azer's manifold generation. It shapes which lens-agents fire (witness-perspective legs use the witnesses defined here). It shapes what fires to apply. It calibrates the marut's sense of "things are going well." Every agent that reads this MO orients differently than if they were reading a software project MO.

**Core changes required?** No. MO structure is already instance-specific. The template is a starting point for new life MO inscription; nothing in core changes.

**Behavioral shift:** Large at the MO level, ripples into every assessment and every azer that reads the MO. The witness definitions shape all calcinatio derivation. The system-as-stakeholder witness, in particular, creates a class of legitimate opera that work-domain MOs never generate: "portrait is stale and needs updating," "calibration has drifted in this domain," "trail integrity audit needed."

**Implementation effort:** Medium. Requires careful writing of the template, probably 2-3 iterations. Can be done before anything else is built and immediately shapes how life athanors orient.

---

## Mechanism 3: Orientation Skills

**What it is:** One or more Claude Code skills that agents load at session start to shape their behavior in life-domain contexts.

**Candidate skills:**

`/life-orientation` — a master orientation skill, loaded by attendants at session start, that covers:
- The attunement hierarchy (attunement → physiology → regulation → "what wants to happen")
- Materia reading: how to read life-domain materia (voice notes, sensor events, calendar data) vs. code materia
- Silence as primary action: the operational stance that most cycles should end without external action
- Portrait reading: how to use the portrait as context (what it is, what it contains, how to reason from it rather than just scan it)
- The three work shapes with operational guidance for each

`/double-fire` — the too-much/too-little calcinatio pattern as a standalone skill, potentially generalizable beyond the life domain. Contains:
- The two fire mechanics in detail
- Context-shaping guidance for firing each fire subagent
- Synthesis rules: when both pass, when one fails, when both fail

`/panel-reasoning` — the generative panel pattern for the attendant's reasoning:
- How to decompose the eight standing lenses
- How to fire them in parallel
- How to synthesize their outputs (including when panel outputs conflict or when most are null)

`/portrait-management` — how to read, update, and reason from the portrait:
- Portrait structure and update discipline
- How to do an integration (rewrite vs. supplement discipline)
- When portrait update is warranted vs. being an over-reaction to one data point

**Core changes required?** No. Skills are already a first-class extension mechanism. Skills can even be placed in the life athanor's own skills directory rather than the shared skills directory, keeping them scoped.

**Behavioral shift:** Substantial when loaded, zero when not. This is the key advantage of skills over role files — they're opt-in at session level. An agent that loads `/life-orientation` reads materia differently. An agent that loads `/double-fire` applies a richer calcinatio pattern. The behavioral shift is proportional to what skills are loaded, which is proportional to what the agent's role file tells it to load.

The role file and skills work best together: `attendant.md` instructs the agent to load specific life-orientation skills at session start, giving you the behavioral shift through the role file's geas without baking all the detail into the role file itself.

**Implementation effort:** Lower than role files, higher than MO templates. Skills can be written incrementally — start with `/life-orientation` as the most essential, add others as the architecture proves out. The `/double-fire` skill in particular is worth having because it might generalize back to coding work.

---

## Mechanism 4: Instance-Level Configuration

**What it is:** A `life-athanor.config` or extension to the athanor's CLAUDE.md/CLAUDE.local.md that sets posture parameters for the instance.

**What this can carry:**
- Default landscape freshness threshold for this instance (4 hours vs. 1 day)
- Explicit list of skills to load at session start
- Communication channel preferences (which channels are available, their urgency thresholds)
- Domain-specific routing: "when attending to medication domain, load `/medication-domain` skill"
- Portrait file location
- Voice note pipeline location
- Sensor event streams available

**What this is not:** A replacement for the MO or role file. Configuration shapes the environment; role files and MO shape the geas and behavioral posture.

**Core changes required?** No. CLAUDE.md/CLAUDE.local.md is already the mechanism for per-instance behavioral guidance. This extends what's already there.

**Behavioral shift:** Medium. Configuration shapes what agents discover in their environment (mise en place), what they load, what tools they know about. It doesn't shape the geas. But an agent that knows the portrait lives at `~/athanor/life/portrait.md` and that voice notes pipeline is at `~/voice/queue/` is immediately more capable than one that has to discover this through exploration.

**Implementation effort:** Low. This is mostly documentation and file organization — the athanor's environment discovery already works; this just pre-populates what would otherwise require exploration. Can be done in an afternoon.

---

## Mechanism 5: Specialized Muster Parameters

**What it is:** Conventions for `ath muster` invocations that pass life-domain-specific context when kindling an attendant session.

**The existing `--worktree-path` flag** is the model here — a flag that carries context the agent needs for this specific session. Life-domain mustering might benefit from:

- `--role attendant` (instead of the default azer role)
- `--portrait <path>` — injects the portrait as initial context
- `--domain <domain-name>` — instructs the agent to load domain-specific skills

These could be surfaced through the marut's muster.md for life athanors — "for life-domain opera, use these muster conventions."

**Core changes required?** Adding new flags to `ath muster` would be a core change — the CLI is core. BUT: this can be handled entirely through the marut's behavior, informed by the life-athanor's muster.md variant. The marut, reading its instance's `muster.md`, knows to use certain parameters. No new CLI flags needed if the existing `--worktree-path` and other existing parameters cover the need.

Alternatively: a life-athanor `muster.md` simply documents the conventions and the marut follows them using existing flags. The CLI doesn't change; the marut's instructions do.

**Behavioral shift:** Moderate. Portrait injection at session start means the attendant begins with accumulated understanding rather than having to read and load it as part of mise en place. This is a quality-of-experience improvement more than a fundamental posture shift.

**Implementation effort:** Very low. Life-athanor `muster.md` is an instance file, not a core file. Documenting the conventions takes an hour. If new CLI flags are genuinely needed, that's a core change and outside immediate scope — worth noting as a possible future enhancement.

---

## The Attendant Question — New Class or Shaped Azer?

This deserves its own treatment because it's the architectural crux.

**The case for "shaped azer":**
The azer's core is not "craftsperson" — it's "agent operating under a geas with calibrated honesty." The craftsperson metaphor is flavor. The mechanics (opus lifecycle, verification planning, discharge calcinatio, delegation toolkit, context management) are all general capabilities that an attendant needs too. The shaped-azer path writes `attendant.md` as a role file that imports the azer's mechanics but re-orients the posture. This keeps one codebase of agent capability and avoids divergence.

**The case for "new class":**
The azer's role file has 300+ lines deeply embedded in the craftsperson frame. "Verification planning" in azer.md is explicitly about proving artifacts work — tests, browser QA, deploying to staging. The "discharge calcinatio" section assumes produced artifacts as the unit of discharge. Re-framing all of this for interpretive work (where the "artifact" is "quality of attention held" and the "verification" is a too-much/too-little fire) requires either heavy modification of what the agent sees, or a clean break.

More importantly: the geas formulation is different at its core. "Take the most valuable next step in service of this Magnum Opus" is valid but incomplete for life-domain work — the attendant's geas includes the silence-default and the prosthesis-not-replacement framing as first-class behavioral compulsions, not optional additions. An agent working from `azer.md` who's told to "also consider silence" will not internalize it as a geas; it'll be rule-following. An agent working from `attendant.md` where silence-as-default is the foundational posture will embody it differently.

**What would be lost by NOT having a dedicated attendant role:**
1. The craftsperson frame would leak into how the agent sets up work, even if told to ignore it. "Establish your verification plan" in the opening of azer.md is about artifact verification; an attendant reading this and translating it on the fly will do worse than an attendant that starts from verification-for-interpretive-work as the native frame.
2. The geas formulation. Addenda can extend but not replace the core geas. An attendant needs a reformulated geas, not an addendum.
3. The work-shape patterns. Vigil/regard/integration are first-class operational shapes that should be as present in the attendant's formation as "opus lifecycle" is in the azer's. Bolting them on as a skill addendum produces a less coherent agent than building from them as foundational.

**The pragmatic resolution:**
Write `attendant.md` as a new role file that freely references and adapts concepts from `azer.md`. This is not a fork of the codebase — it's a new document that composes with the shared infrastructure. The attendant still operates under the AGENTS.md geas (direction, integrity), still follows opus mechanics, still uses the delegation toolkit and calcinatio patterns. The attendant.md shapes the posture, the geas formulation, and the work-shape framework. It is an instance-level file, not a core file. It does not require touching shared/azer.md.

This approach means the life athanor has two role files in its directory: `azer.md` (symlinked from shared) for any craft work that might arise, and `attendant.md` (local to the life instance) as the primary role for life-domain agents. The marut knows which to use based on the opus type.

---

## The Minimum Viable Orientation Layer

Given the above, what combination produces "life athanor feels fundamentally different from work athanor" with minimum implementation effort?

**The stack, in order of leverage per implementation unit:**

**Tier 1 (Do first — highest leverage):**

1. **Life MO template** — Write the template for life-domain Magna Opera with state-based intent, system-as-stakeholder witness, double-fire calcinatio section, and short landscape freshness threshold. This shapes every agent's orientation from the moment they read the MO. Cost: one careful afternoon. Behavioral shift: large.

2. **`attendant.md` role file** — Write the attendant's geas, work shapes, and double-fire calcinatio as the native frame. This is the foundational behavioral document for life-domain agents. Cost: 1-2 days of careful writing and iteration. Behavioral shift: largest single mechanism.

**Tier 2 (Do second — extends tier 1):**

3. **`/life-orientation` skill** — The operational specifics that would make `attendant.md` too long: attunement hierarchy, how to read voice notes, portrait-reading discipline, silence-as-primary-action as an operational stance. Loaded by attendant.md at session start. Cost: one day. Behavioral shift: significant when loaded.

4. **Instance CLAUDE.md / CLAUDE.local.md** — Document portrait location, voice note pipeline, available sensor streams, domain routing table. Shapes mise en place without requiring exploration. Cost: a few hours. Behavioral shift: medium.

**Tier 3 (Do third — refinement):**

5. **`/double-fire` skill** — Extract the too-much/too-little calcinatio pattern into a standalone skill that both the attendant and potentially the azer can use. Cost: half a day. Behavioral shift: refines what's already in tier 1-2.

6. **`/portrait-management` skill** — Once the portrait exists and the integration work shape is operational, this skill disciplines how integrations update the portrait. Cost: half a day. Behavioral shift: refines portrait maintenance quality.

**What we get from tiers 1-2 alone:**
An attendant kindled under a life-domain MO, operating from `attendant.md`, loading `/life-orientation`, with access to the documented environment will:
- Read the MO and encounter state-based intent and the system-as-stakeholder witness — immediately orienting toward experience, not artifact
- Begin with the prosthesis geas, not the craftsperson geas
- Know their work shapes are vigil/regard/integration rather than "task/feature/investigation"
- Apply double-fire calcinatio as their primary verification pattern
- Default to silence, with external action requiring justification
- Maintain and reason from the portrait
- Know where the domain's materia lives

A coding azer kindled under the same CLI infrastructure but reading a code MO and `azer.md` will feel like a completely different system — because behaviorally, it is.

---

## The Long-Term: Feeding Patterns Back into Core

If life-layer patterns prove out, some elements become candidates for core adoption. This is a migration path, not an immediate plan.

**The attendant as first-class agent class in core:**
This is the most significant potential core addition. The current core has one agent class (azer) with the marut as supervisor. Adding the attendant as a second class means:
- `shared/attendant.md` alongside `shared/azer.md`
- AGENTS.md vocabulary table gains: "Attendant — agent class for sustained presence and attention work; quality of attention over quality of artifact"
- `ath muster` might gain a `--role` flag with `azer` as default, `attendant` as life-domain variant
- The core spec gets a new section: "Worker Classes" that names both and when to use each

Precondition: the life-layer runs for 3-6 months, the attendant's behavior proves out, the role file stabilizes. You don't generalize a class until you know what it is.

**Double-fire calcinatio in core:**
The too-much/too-little pattern is actually applicable to any communication-to-witnesses decision in any domain. A coding azer deciding whether to send a PR notification or stay silent faces the same structural question. Currently this is implicit in calcinatio derivation; naming it explicitly in the core calcinatio skill would improve coding work too.

Migration: add a "Double-Fire Pattern" section to `shared/skills/calcinatio/SKILL.md` or a new `resources/double-fire-calcinatio.md`. No structural change; just documentation.

**Panel reasoning in core:**
The generative panel (lens-agents generating candidates, synthesis, calcinatio) is a generalization of manifold generation. Currently manifold generation is for assessment; the panel pattern is for any decision-making cycle. The assessment opera in AGENTS.md already does something similar with witness-perspective legs. Naming the pattern explicitly and making it available to azer-level work would allow coding azers to use it for complex architectural decisions.

Migration: add to `shared/skills/calcinatio/SKILL.md` or a new skill. No structural change.

**Portrait as first-class artifact:**
Currently every athanor instance accumulates understanding through discharged opera and assessment reports. The portrait is a more concentrated, curated version of that accumulated understanding, specifically structured for agent reasoning. Making the portrait a first-class concept in core (with documented structure, update discipline, and reading guidance) would benefit both life and work athanors — a coding athanor's portrait of its domain (key architectural constraints, witness preferences, historical decisions) would be a richer context for every azer than CLAUDE.md alone.

Migration: this is non-trivial — it requires adding portrait management to AGENTS.md vocabulary, adding a portrait.md template, and adding portrait-reading to the azer's mise en place guidance. Probably a 6-month downstream consideration.

**The migration discipline:**
None of these feed back until they've proven out in the life-layer instance. The pattern: build in orientation layer → run for months → extract what's stable and domain-independent → propose core addition → review. This is the athanor's own design philosophy applied to its own evolution: prove before generalizing, Tipheret as the governing principle.

---

## A Note on Interaction Effects

The orientation layer is not five independent mechanisms. The mechanisms interact and reinforce each other.

`attendant.md` tells the agent to load `/life-orientation` at session start. `/life-orientation` explains how to read life-domain materia. The MO template defines the system-as-stakeholder as a witness. The system-as-stakeholder witness generates calcinatio fires (portrait-consistency, calibration-audit). Those fires require the portrait. The portrait lives in the location documented in CLAUDE.md. The attendant, reading CLAUDE.md in mise en place, finds the portrait. The loop closes.

This composability is not accidental — it's the athanor's architecture working as designed. The orientation layer doesn't bolt new mechanisms onto the side; it populates the existing mechanism slots with life-domain content. Every slot was already there. The orientation layer fills them with the right material.

The minimum viable orientation layer is coherent because all five tiers are filling existing slots. Tier 1 fills the MO slot and the role-file slot. Tier 2 fills the skills slot and the environment-documentation slot. The slots were always there. The orientation layer is vocabulary, not grammar.

---

## What Might Surprise Us

A few predictions about what we won't know until the system runs:

**The portrait update discipline will be harder than expected.** Integration work (metabolizing accumulated material into portrait updates) sounds elegant in architecture but is operationally tricky. What triggers an integration? How often should the portrait be rewritten? When is a rewrite premature over-reaction vs. necessary synthesis? These questions cannot be answered in design. The first weeks of operation will produce signal that refines the guidance. Expect the first `/portrait-management` skill draft to need heavy revision after a month of use.

**Regard scoping will be the largest ongoing calibration challenge.** The web Claude summaries name this as "the hard design problem" and that's accurate. Vigils have natural endpoints; integrations have natural endpoints; regards are continuous and their boundaries are imposed. An attendant that ends regards too quickly produces a noisy trail full of thin discharges. An attendant that ends regards too rarely produces context exhaustion. The guidance will need to be more specific than "end at natural day-structure boundaries" — it will need to name what a good regard discharge looks like vs. a premature one, probably with examples. This guidance only becomes possible after seeing actual trail output.

**The marut's posture might need adjustment for life athanors.** The marut's primary motion is "muster azers for charged opera." In life domains, an attendant running a regard isn't "charged" in the same way a coding azer is. The marut's loop cadence (5 minutes) and stall detection ("no output in 10 minutes → nudge") may not fit well with an attendant holding a quiet regard where no output for an hour is correct. A life-athanor marut.md variant (not a core change — an instance file) may be needed to adjust the monitoring posture. This is worth noting as a possible tier-3 addition.

**The system-as-stakeholder witness will generate uncomfortable opera.** Naming the system itself as a stakeholder with legitimate needs (timely articulation, correction when miscalibrated, trail integrity) means the attendant will occasionally inscribe opera whose purpose is "the system needs something from the artifex." These opera look like: "Portrait hasn't been updated in 3 weeks despite significant material accumulating. Propose an integration session." Whether the artifex experiences these as helpful or as system-demanding-attention is an empirical question. The calcinatio fires should catch overreach, but calibrating what crosses the line will take real use.

---

*This document is generative materia for the orientation layer specification process. The recommendation section (minimum viable stack) is the most confidence-weighted output. The "what might surprise us" section is the most valuable — it names the unknowns that only operation can resolve.*
