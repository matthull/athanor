# Assessment Enhancement Specification

**Status:** Draft (active shaping session)
**Created:** 2026-03-31
**MO:** athanor-alignment
**Opus:** 2026-03-31-assess-enhance

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[S:session]` | From the 2026-03-31 shaping session (artifex + azer) |
| `[G:gastown]` | Pattern validated in Gastown codebase |
| `[E:agents-md]` | Extends existing AGENTS.md § Assessment Opera |
| `[R:research]` | From research (pending or completed) |
| `[D:reason]` | Design decision with rationale |
| `[U:topic]` | Unbacked / needs validation |

---

## Overview

The assessment function is the athanor's most critical value-extraction moment. It is the point where the widest context is assembled — trail, environment, discharged opera, witness definitions, corpus — and the system decides what to do next.

**The core reframe:** Assessment is not defensive ("was value lost?"). Assessment is **generative and jovian** — "how do we bountifully fulfill?" The defensive/saturnian energy belongs at the discharge boundary (discharge calcinatio in `azer.md`). Assessment is where the athanor's engine runs at peak power. `[S:session]`

**The problem with the current process:** The existing assessment in AGENTS.md § Assessment Opera gives per-opus assessor subagents a broad mandate ("was charge fully converted to materia?") but this framing is loss-prevention (Saturn), not value-generation (Jupiter). The per-opus fan-out structure also over-centers discharged opera as the primary source of signal, when a Slack thread, a fresh read of the spec, or an environment change could be equally or more valuable. `[S:session]` `[E:agents-md]`

**What this spec covers:** A redesigned assessment process that centers the generative question — "given everything available, what is the most bountiful path forward?" — organized around a sephirotic hierarchy that establishes the energy flow from intent through manifestation.

---

## Domain Dictionary

### Supernal Energies (the drivers)

**Chokmah (Force/Will)** — The geas to seek abundant satisfaction. The raw directive energy that says "the world should be different in THIS way." In the assessment context: the MO intent and the drive toward fulfillment. Not all geases — specifically the seeking energy. `[S:session]`

**Binah (Understanding/Form)** — Understanding of witnesses. The receptive intelligence that shapes force into specific form. Knows WHO we serve and what fulfillment looks like from each witness's vantage. `[S:session]`

### Operative Energies (the modes)

**Chesed (Jupiter/Mercy/Expansion)** — Bountiful, generative action. The mode of assessment that generates ideas, imagines possibilities, brainstorms from rich context. Orthogonal emanation is a Chesed tool. Manifold generation (see below) is a Chesed pattern. `[S:session]`

**Geburah (Saturn-Mars/Severity/Refinement)** — Disciplined, refining action. The mode that prunes, challenges, verifies. Calcinatio in all its forms (dialectical, manifold) is Geburah energy. `[S:session]`

### Structural Positions (the path)

**Tiphereth (Beauty/Balance)** — The synthesis point. Where the assessment azer integrates Chesed-generated ideas and Geburah-refined findings into a coherent set of next actions. `[S:session]`

**Yesod (Foundation)** — The opera inscribed. Practical infrastructure — concrete, actionable work units with clear paths to manifestation. `[S:session]`

**Malkuth (Kingdom/Manifestation)** — The actual work done in the world. Code shipped, docs updated, witnesses notified. Where the next azer picks up the inscribed opus. `[S:session]`

### Formula (Briah-level concept)

**Formula** — An essential pattern you reach for in certain situations. Formulae live in Briah — the world of creation/essential patterns. Not the specific implementation (Yetzirah) and not the concrete execution (Assiyah). A formula is archetypal: it defines the *shape* of a process, not the details. You reach for a formula like you reach for a tool, but the formula is *what* to do, not *how* in a specific case. The term is native to both the alchemical register (a recipe for transmutation) and the Gastown precedent (structured multi-agent orchestration patterns). `[S:session]`

**Examples of formulae:**
- **Dialectical calcinatio** — two parties, thesis → challenge → refine → converge
- **Manifold calcinatio** — parallel independent fires → synthesize
- **Manifold generation** — parallel independent generation from stimuli → collect
- **Orthogonal emanation** — maximally divergent candidates from single intent
- **The assessment process itself** — sephirotic flow from supernals through manifestation

### Patterns

**Manifold Generation** — A formula. The generative analog of manifold calcinatio. Context is decomposed into individual stimuli; a subagent is fired per stimulus, each paired with the constant supernal context (MO intent + witness understanding). The diversity of ideas comes from the diversity of stimuli, not from forcing a single agent to think divergently. Fan-out for breadth of *ideas*, not breadth of *critique*. `[S:session]`

**Combinatorial Stimulus** — The mechanism underlying manifold generation. The supernals (MO intent + witness understanding) are constants. The pieces of loaded materia are variables. Each subagent gets constants + ONE stimulus and is asked: "What does this combination spark? What actions could bountifully serve the intent given what you see here?" `[S:session]`

**Witness-Perspective Legs** — A specific class of stimulus within manifold generation. Instead of pairing supernals with a *materia item*, pair them with a *single witness's full concern profile*. A subagent thinking purely from Witness A's vantage sparks different ideas than one thinking from Witness B's. These run alongside materia-stimulus legs in the Chesed phase. `[S:session]`

---

## The Hierarchy

The assessment process follows a sephirotic hierarchy — not a flat checklist. Energy flows from the supernal drivers through operative modes to concrete manifestation.

```
Chokmah (will to satisfy) + Binah (understanding of witnesses)
            |                          |
     produces either/both:
  Chesed (expansion)        Geburah (refinement)
            |                          |
              Tiphereth (synthesis)
                    |
              Yesod (opera inscribed)
                    |
              Malkuth (work manifested)
```

`[S:session]`

**Why this hierarchy matters:** It establishes that assessment starts from WHY (the supernal energies) and flows through HOW (the operative modes) to WHAT (the concrete actions). The current flat process starts from WHAT (per-opus evaluation) and works backward. The hierarchy reverses this — intent drives action, not the other way around.

**Chesed and Geburah are both valid outputs.** The interplay of Chokmah + Binah determines which energy is needed. A landscape rich with unserved witnesses calls for Chesed (generous expansion). A landscape where recent work is unrefined calls for Geburah (calcinatio). Often both are needed in the same assessment cycle.

**The path of the sword (Tiphereth -> Yesod -> Malkuth)** is the path of direct manifestation — synthesis of findings into balanced decisions, crystallized as inscribed opera, actualized in the world. This is the practical output chain. `[S:session]`

---

## The Assessment Process (Proposed)

### Phase 1: Load Supernals

The assessment azer begins by loading the two driving energies:

**Chokmah — Load the will.** Read the MO intent. Feel the force. What change in the world does this MO demand? What does the geas to seek abundant satisfaction compel? This is not a checklist item — it's establishing the directive energy that everything else serves. `[S:session]`

**Binah — Load the understanding.** Read each witness definition. Understand their specific concerns, how they observe, how they're reached, what "bountifully served" means from their vantage. This is the receptive intelligence that shapes the will into specific forms of fulfillment. `[S:session]`

The supernals become the **constants** passed to every subsequent subagent — the shared frame that ensures all generated ideas serve the same intent and witnesses.

### Phase 2: Load Charge (Prospector Survey)

The assessment azer surveys the full landscape of available materia. This is the prospector phase — not evaluating, but discovering where the richest signal lies. `[S:session]`

**Sources to survey** (not exhaustive — the azer should discover what's available):
- Trail of discharged opera (outcomes, reflections, gaps noted)
- Environment state (repo diffs, CI status, open PRs, failing tests)
- Communication channels (Slack threads, ticket comments, meeting transcripts)
- Corpus (specs, docs, ADRs, runbooks — read with fresh post-trail eyes)
- Any other materia relevant to the MO

**Per-opus intake.** Discharged opera still get read and marked `status: assessed` during this phase. But this is intake — loading charge into the crucible — not the center of the process. A discharged opus is one rich source among many. `[E:agents-md]` `[S:session]`

**Output of this phase: the landscape report.** A structured summary of the landscape with citations to primary sources (use `/citations` skill for source traceability). The report captures what the prospector found — trail state, environment signals, channel activity, corpus condition, witness pulse — so that generation subagents receive rich context without re-gathering it. `[S:session]`

The landscape report serves three purposes:
1. **Shared base for generation legs.** Baked verbatim into every Phase 3 subagent's prompt alongside the supernals. The legwork is done once, documented once, shared to all.
2. **Depth on demand.** Citations link to primary sources (Slack URLs, Linear issues, file paths, opus paths). Generation subagents can follow citations to drill deeper on items that spark ideas.
3. **Durable artifact.** A snapshot of the athanor's state at assessment time. Gets appended to the assessment discharge record. Future assessments can read prior landscape reports to see how state evolved.

The report's structure, level of detail, and emphasis are the assessor's judgment call — shaped by what the landscape actually contains, not by a predetermined template. The only formatting directive is `/citations` for source traceability. `[D:jovian-default — avoid premature structural constraint; add structure only in response to specific observed needs]`

**Implementation note:** This phase may benefit from subagent scouts dispatched per source type (trail scout, environment scout, channel scout, corpus scout) to keep the main session lean for synthesis. See Emanation #2 (Prospector) for the detailed scout pattern. `[S:session]`

### Phase 3: Chesed — Manifold Generation

This is the core generative phase. The assessment azer decomposes the loaded materia into individual stimuli and fires a subagent per stimulus, each paired with the supernal constants.

**The formula:**

Two classes of legs fire in parallel:

**Context shaping for all legs.** Every leg receives a self-contained prompt with three layers `[S:session]` `[G:gastown]`:
   - **Supernals** (shared constant) — the MO intent and witness definitions, baked in verbatim. Not referenced by path — embedded in the prompt, because each subagent launches into fresh context with nothing else.
   - **Landscape report** (shared constant) — the Phase 2 output, baked in verbatim. Gives every leg the full picture without re-gathering. Citations link to primary sources for depth on demand.
   - **Stimulus** (unique per leg) — the one thing that makes this leg different from its siblings.

No leg knows about sibling legs. Independence at launch is what produces divergence. `[G:gastown]`

**Materia-stimulus legs** — one per richest materia item identified in the landscape report. The stimulus is ONE specific item (an opus reflection, a Slack thread, a spec section, an environment signal). The generative question: *"Flowing from this intent and these witness desires, what does this specific item spark? What actions, communications, documents, connections, opera could bountifully serve? Recommend your single best idea with rationale."*

**Witness-perspective legs** — one per witness defined in the MO. The stimulus is ONE witness's full concern profile — their concerns, how they observe, how they're reached. The generative question: *"Thinking purely from this witness's vantage — given everything in the landscape report, what would most bountifully serve them? Recommend your single best idea with rationale."*

**Collect the full idea space.** Every subagent's recommendation, unfiltered. The materia legs produce ideas sparked by specific signals; the witness legs produce ideas driven by specific stakeholder needs. Together they cover both "what does the landscape offer?" and "what do the witnesses need?" `[S:session]`

**Why this works:** A single agent with huge context generates ideas linearly — the most salient signals dominate and the long tail never surfaces. Multiple focused subagents, each with a small context + clear direction, produce specific and surprising ideas. The diversity comes from the diversity of stimuli. `[S:session]`

**Gastown validation:** This pattern is proven in Gastown's convoy formula type. The `design` formula fires 6 parallel legs, each with the same problem + a different focused lens. Key learning: each leg should "explore options, list pros/cons, and make a recommendation" — not just brainstorm, but *recommend*. The synthesis step then has N concrete recommendations to work with. `[G:gastown]`

**The Gastown `mol-idea-to-plan` pipeline** implements the full Chesed->Geburah flow: 6 parallel generation legs (Chesed), then iterative critique/refinement rounds (Geburah). Generate broadly, then constrain and refine. `[G:gastown]`

**Pure brainstorming complement.** In addition to the per-stimulus subagents, the assessment azer itself may brainstorm freely from its loaded context: *"Flowing from the supernals of witness desires and magnum opus intent, in the spirit of Chesed, generate ideas for how we can make use of the charge loaded to the crucible."* This captures ideas that emerge from the *connections between* materia items, not just individual stimuli. `[S:session]`

**Scaling.** The number of stimuli (and thus subagents) scales naturally with available charge. A light assessment cycle with one discharged opus and no environmental signals might fire 2-3 subagents. A rich cycle after major work with active Slack threads and spec changes might fire 8-10. The prospector survey (Phase 2) determines the scale. `[D:natural-scaling]`

### Phase 4: Geburah — Calcine

The generated idea space is now refined through calcinatio. The assessment azer applies the saturnian tools:

**Dialectical calcinatio** — The azer takes the most promising generated ideas and fires challenger subagents. Each challenger gets the idea + the supernals + relevant context, and attacks: "Is this actually the most valuable action? What's wrong with this recommendation? What would a witness say?"

**Manifold calcinatio** — If the idea space is large, multiple independent fires can evaluate different proposed actions simultaneously from different angles (feasibility, witness impact, geas alignment, urgency).

**The forward-pressure question** (from Emanation #4, Dialectical-Synthesis): One challenger round should be explicitly generative rather than adversarial — "What would create the *feeling of abundance*? What one or two additions would shift this from 'done' to 'generous'?" This keeps Chesed energy alive even in the Geburah phase. `[S:session]`

### Phase 5: Tiphereth — Synthesize

The assessment azer integrates the Chesed-generated ideas (surviving Geburah refinement) into a coherent set of next actions. This is the beauty/balance point — harmonizing competing priorities, resolving tensions between witness needs, choosing what actually gets inscribed.

**Outcomes** (same three as current process, but arrived at through a richer path):
1. **Inscribe one or more concrete opera** — each tracing back to a generated idea and its originating stimulus
2. **Declare the Magnum Opus abundantly satisfied** — with evidence from the full landscape survey, not just opera assessment
3. **Escalate** — if the artifex's judgment is needed

**Tension mapping** (from Emanation #1): Where witness needs conflict, surface the tension explicitly rather than resolving it silently. Conflicting witness needs are escalations, not silent compromises. `[S:session]`

### Phase 6: Yesod — Inscribe

Concrete opera are inscribed using the standard inscription process (dialectical calcinatio with inscription subagent per `azer.md`). Each opus traces to:
- The originating stimulus (which materia item sparked it)
- The supernal energy it serves (Chesed or Geburah)
- The witness(es) it bountifully serves

### Phase 7: Malkuth — Manifest

The next azer picks up the inscribed opus and manifests it in the world. This is outside the assessment process itself — it's the downstream effect.

---

## Relationship to Existing Process

The proposed process **extends**, not replaces, the current AGENTS.md § Assessment Opera. Key deltas:

| Current | Proposed |
|---------|----------|
| Centered on per-opus evaluation | Centered on supernal energies (intent + witnesses) |
| Fan-out per discharged opus | Fan-out per materia stimulus (opera are one source among many) |
| Core question: "was value lost?" (Saturn) | Core question: "what's the most bountiful path?" (Jupiter) |
| No generative phase | Explicit Chesed phase (manifold generation) |
| Synthesis from subagent reports | Synthesis at Tiphereth after Chesed + Geburah |
| Flat structure | Hierarchical energy flow (supernals → operative → manifestation) |

**What's preserved:**
- Per-opus assessment still happens (Phase 2 intake) — discharged opera are still read and marked `assessed`
- Environment scanning still happens (Phase 2 prospector survey)
- Subagent delegation for context management
- Three outcomes: inscribe, declare satisfied, escalate
- Discharge calcinatio for the assessment opus itself
- The assessment opus template (updated to reflect the new process)

---

## The Generative Pattern: Manifold Generation

### The Problem

When a single agent session has a huge loaded context and is asked "what should we do?", ideas come out linearly. The most salient signals dominate, and the long tail of valuable-but-non-obvious ideas never surfaces. More context in one session can actually *reduce* generative diversity because prominent signals crowd out subtle ones. `[S:session]`

### The Pattern

**Manifold generation** is the generative analog of manifold calcinatio:

- **Manifold calcinatio:** Fan out for breadth of *critique*. Multiple independent fires evaluate the same work from different angles.
- **Manifold generation:** Fan out for breadth of *ideas*. Multiple independent agents generate from different stimuli, each paired with shared direction.

The mechanism is **combinatorial stimulus:**
- **Constants** (passed to every subagent): The supernals — MO intent + witness understanding
- **Variables** (one per subagent): A specific materia item — a discharged opus reflection, a Slack thread, a spec section, an environment signal
- **Question** (same for every subagent): "What does this combination spark? What would bountifully serve?"

`[S:session]`

### Gastown Precedent

The convoy formula type in Gastown is the proven implementation of this pattern. Key examples:

**`design` formula** — 6 parallel legs (api, data, ux, scale, security, integration), each given the same problem + a focused lens. Each explores options, lists pros/cons, makes a recommendation. Synthesis collects all recommendations. `[G:gastown]`

**`mol-prd-review`** — 6 legs analyzing a PRD from different angles (requirements, gaps, ambiguity, feasibility, scope, stakeholders). Same input, different lens. "Identify overlapping findings across legs (these are higher-confidence)." `[G:gastown]`

**`mol-idea-to-plan`** — Full pipeline: idea → 6-leg generation (Chesed) → iterative critique/refinement (Geburah) → beads. The most complete example of all patterns working together. `[G:gastown]`

**Consensus pattern** — Explicitly documented: "Same prompt → N agents → compare. Routing goal: Diversity (multiple perspectives)." `[G:gastown]`

**Mountain-Eater Dogs** — Ephemeral agents with narrow scope + fresh context. "The Dog starts with zero state. One convoy, one question." Validates the fresh-context-per-stimulus approach. `[G:gastown]`

### Open Research

A research request has been created at `research/requests/generative-context-extraction-patterns.md` covering:

1. Established patterns in multi-agent systems for generating diverse ideas from accumulated context
2. Knowledge management frameworks for the "what should we do next?" moment (AAR, Toyota kata, retrospectives — generative versions)
3. Semantic web / ontology patterns for reasoning about potential actions from state descriptions
4. Patterns for combating mode collapse in idea generation (TRIZ, lateral thinking, computational creativity)
5. Multi-agent architectures for creative/generative tasks (swarm intelligence, blackboard architectures)
6. Whether "combinatorial stimulus" is studied and has known refinements (SCAMPER, oblique strategies, design space exploration)

`[R:research]` — Findings pending.

### SKAN Agentic Ontology of Work

Research was previously commissioned on the SKAN AOW framework, which uses "opus" and "magnum opus" in an agentic work context — directly parallel to the athanor. The research request is at `research/requests/skan-agentic-ontology-of-work.md`. Findings have been executed but not yet integrated into this spec. `[R:research]` — Findings pending integration.

---

## Orthogonal Emanations (Design Exploration Record)

Five maximally divergent assessment designs were generated via orthogonal emanation. Each was seeded on a structurally different organizing principle. The final design synthesizes elements from multiple emanations.

### Emanation #1: Witness-Centric

**Organizing unit:** One subagent per witness (not per opus).

Each "witness advocate" reads the full trail through that witness's eyes. Key innovations:
- **Roster expansion** — infer implicit witnesses the MO didn't name
- **Tension mapping** — where witnesses conflict, surface the tension rather than silently resolving
- **Witness overflow** — when a need exceeds MO scope, flag for MO revision rather than inscribing opera

**What was adopted:** Tension mapping (Phase 5). Witness understanding as a supernal driver (Binah). `[S:session]`

### Emanation #2: Context-Landscape (Prospector)

**Organizing unit:** One scout per source type (trail, specs, environment, channels, corpus).

The assessor is a prospector surveying all available materia for the richest veins. Key innovation: **follows the signal, not the opus chain.** A hot Slack thread outranks a routine discharge.

**What was adopted:** The prospector survey (Phase 2). The principle that discharged opera are one source among many. `[S:session]`

### Emanation #3: Manifold-Axis

**Organizing unit:** One subagent per target of bountiful fulfillment (intent fidelity, witness satisfaction, geas integrity, abundant satisfaction, corpus enrichment).

Key innovation: **conjunctive synthesis** — each axis is veto-capable. No averaging.

**What was adopted:** The insight that targets of bountiful fulfillment should be centered. However, the hierarchy (Chokmah+Binah producing Chesed/Geburah) replaced the flat axis structure. The targets are downstream manifestations, not independent axes. `[S:session]`

### Emanation #4: Dialectical-Synthesis

**Organizing unit:** Sequential rounds (thesis → challenge → refine).

The assessment azer builds understanding through progressive dialectical exchanges. Three rounds: corpus faithful → witness bountiful → forward pressure. Key innovation: **the forward-pressure round** — "What would create the *feeling of abundance*?" — is distinctly jovian.

**What was adopted:** The forward-pressure question (Phase 4). The principle that the azer is a continuous thread building lived understanding. `[S:session]`

### Emanation #5: Tincture-Flow

**Organizing unit:** Output products (four tinctures: witness, corpus, intent, geas).

Assessment as distillation — raw materia in, refined products out, each flowing to a specific destination. Key innovation: **the gap IS the charge** — a cloudy tincture mechanically becomes an opus.

**What was adopted:** The output-oriented thinking contributed to the Yesod phase. However, the tincture framing was superseded by the sephirotic hierarchy — witness understanding and MO intent are not "outputs" but supernal *drivers*. `[S:session]`

---

## Tools Mapping

The athanor's tools map to sephirotic energies:

| Tool/Pattern | Energy | Use in Assessment |
|---|---|---|
| Manifold generation | Chesed (Jupiter) | Phase 3 — generate ideas from stimuli |
| Pure brainstorming | Chesed (Jupiter) | Phase 3 — generate from connections |
| Orthogonal emanation | Chesed (Jupiter) | When choosing between design approaches |
| Dialectical calcinatio | Chokmah-Binah interplay | Phase 4 — thesis + understanding in dialogue |
| Manifold calcinatio | Geburah (Saturn-Mars) | Phase 4 — independent fires from multiple angles |
| Forward-pressure question | Chesed within Geburah | Phase 4 — "what would create abundance?" |
| Assessment synthesis | Tiphereth (Beauty) | Phase 5 — balance, harmonize, decide |
| Opus inscription | Yesod (Foundation) | Phase 6 — concrete actionable work units |

`[S:session]`

---

## Open Questions

1. **How does the per-opus `assessed` marking fit into the new flow?** Currently it happens in Phase 2 (intake). Is that sufficient, or should it be deferred until after the full assessment cycle? `[U:lifecycle]`

2. **What's the right number of stimuli?** The scaling principle says "it depends on available charge," but are there practical guidelines? Gastown's formulas use 5-10 legs. `[U:scaling]`

3. **Should the assessment opus template change?** The current template in AGENTS.md needs updating to reflect the new process. The template itself may need to encode the sephirotic phases. `[U:template]`

4. **How does this interact with the marut?** The marut inscribes assessment opera. Does it need to know about the new process, or does the assessment azer just execute it? `[D:marut-agnostic — the azer executes the process; the marut's role doesn't change]`

5. **What are the failure modes?** Manifold generation could produce overwhelming noise if too many stimuli are fired. The Geburah phase needs to be proportionally strong. `[U:failure-modes]`

6. **Research integration pending.** The generative-context-extraction-patterns research request may surface proven patterns that refine or replace the manifold generation approach. The SKAN AOW findings may offer structural insights. `[R:research]`

7. **How explicit should the Kabbalistic framing be in AGENTS.md?** The hierarchy is a useful design tool. The question is how much of the metaphysical framing belongs in the operational guidance vs. this spec. `[U:framing]`

---

## Workflow / Dataflow

```
                    [MO Intent]         [Witness Definitions]
                        |                       |
                    Chokmah                  Binah
                    (will)              (understanding)
                        |                       |
                        +----------+------------+
                                   |
                          [Supernal Constants]
                                   |
                    +--------------+--------------+
                    |                             |
              [Phase 2: Survey]            [Phase 2: Survey]
              Trail, Environment           Channels, Corpus
                    |                             |
                    +----------+------------------+
                               |
                    [Identified Materia Items]
                    (richest signals from landscape)
                               |
              +-----+-----+---+---+-----+-----+
              |     |     |       |     |     |
           [Sub1] [Sub2] [Sub3] [Sub4] [Sub5] [...]
           stim+  stim+  stim+  stim+  stim+
           supn   supn   supn   supn   supn
              |     |     |       |     |     |
              +-----+-----+---+---+-----+-----+
                               |
                    [Full Idea Space]
                    (N concrete recommendations)
                               |
                    [Phase 4: Geburah]
                    Dialectical + manifold calcinatio
                    + forward-pressure question
                               |
                    [Refined Ideas]
                               |
                    [Phase 5: Tiphereth]
                    Synthesize, balance, decide
                               |
              +--------+-------+-------+
              |        |               |
          [Inscribe] [Declare     [Escalate]
          opera      satisfied]
              |
          [Phase 6: Yesod]
          Concrete opera with
          clear paths to Malkuth
```

---

## Test Scenarios

Since this is a process/protocol spec (not code), test scenarios are behavioral:

1. **An assessment with one discharged opus and no environmental signals** should still survey the full landscape and may generate ideas from corpus or channel sources, not just the opus.
2. **An assessment with no discharged opera** (queue-empty assessment) should still execute the full process — the prospector survey and manifold generation are not dependent on having discharged opera.
3. **An assessment where a Slack thread is the highest-value signal** should produce opera oriented around that thread, not centered on discharged opera.
4. **An assessment where witnesses conflict** should surface the tension as an escalation, not silently resolve it.
5. **An assessment where the MO is abundantly satisfied** should be able to declare this with evidence from the full landscape, not just from per-opus evaluation.
6. **The manifold generation phase should produce ideas the main session would not have generated alone** — this is the core value proposition and the hardest to test.

---

## Out of Scope

- **CLI changes** — this spec covers the conceptual process, not `ath` commands
- **Marut role changes** — the marut's operational loop (inscribing assessment opera, monitoring, cleanup) is unchanged
- **Discharge calcinatio** — already well-specified in `azer.md`, complementary to but separate from assessment
- **Opus inscription protocol** — already well-specified in `/opus` skill
- **Implementation of manifold generation as a reusable tool** — this spec uses it within assessment; generalizing it as a named pattern/skill is a follow-on concern

---

## Retrospective

When this design is implemented and tested in real assessment cycles:

### Documentation Updates
- [x] AGENTS.md § Assessment Opera — fully rewritten (commit 19703c9, 2026-04-01)
- [x] Assessment opus template — updated with manifold generation and landscape report (commit 19703c9)
- [ ] Assessment opus template — encode sephirotic phases (deferred — kept in spec pending more operational experience)
- [ ] azer.md — check for any needed updates re: assessment behavior
- [x] specs/spec.md — cross-reference added to Built table (commit 19703c9)
- [ ] marut.md — check alignment with new assessment process
- [ ] calcinatio skill — add manifold generation as peer formula alongside manifold calcinatio

### Pattern Capture
- [ ] Manifold generation pattern — document as a reusable formula alongside manifold calcinatio and orthogonal emanation
- [ ] Combinatorial stimulus — if the research finds a name, adopt it
- [ ] Sephirotic hierarchy — if useful beyond assessment, document as a general design principle
- [x] Formula as Briah-level concept — added to AGENTS.md vocabulary (commit 19703c9)

### Research Integration
- [ ] Generative context extraction patterns — execute research request and integrate findings
- [ ] SKAN AOW — integrate findings (research executed, findings not yet loaded)
- [ ] A/B test results — compare thin-context vs full-MO assessments from seismic-classifier-mapping, integrate learnings
- [ ] Any patterns discovered that refine or replace manifold generation
