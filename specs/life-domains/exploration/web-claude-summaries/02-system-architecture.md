# The System Architecture

*The concrete architectural commitments for the life-layer system. How it's structured, what its agents do, what its units of work are, how it reasons, how it stays calibrated over time.*

This document assumes the conceptual foundation in `01-conceptual-foundation.md` and the Athanor spec. Where this architecture inherits from or extends Athanor, it's noted. Where it departs or adds new shapes, the reasoning is given.

---

## The organizing question

The primary question driving the system is: **how can I act as a cognitive prosthesis for the artifex?**

This question includes silence as a legitimate answer. It rejects "helpfulness" as an independent goal. It names what the system is *for* in terms that resist the usual productivity-bot gravity.

A prosthesis compensates for a capacity that isn't fully available. It does this without demanding acknowledgment of its existence, without causing atrophy in tissue that could develop, and with discipline about when its support is warranted versus when silent readiness is the right service.

Most of what the system does is not production. It's attendance. Watching, absorbing, remembering, noticing, responding when warranted, staying quiet when not. The default action on any given cycle is to groom internal state and take no externally-visible action. Not because nothing is happening — the internal work is real and continuous — but because external action costs attention and attention is what the system exists to protect.

## Stakeholders

The system recognizes two primary stakeholders, framed in the Athanor witness pattern.

### The artifex (primary witness)

You. The one the system serves. The witness whose experience defines whether the work is being done well. Intent is always expressed as a change in your lived experience, never as a change in system state or task completion.

Critically, the artifex has specific cognitive topography: ADHD, autism, the energy asymmetries and interoceptive gaps and executive function considerations that make conventional systems fail. The MO explicitly acknowledges this. Witnesses aren't generic users.

### The athanor itself (system-as-stakeholder)

The system has its own legitimate stakeholder needs that must be named explicitly or they get lost. The athanor's ability to serve the artifex is bounded by what it can know and what it can do. These bounds create stakeholder concerns:

- **Timely articulation.** The system needs to be told things while they're fresh. Retrospective articulation carries less signal than proximate voicing.
- **Specificity.** Vague articulation produces vague context. The system benefits from discriminating articulations, not just labels.
- **Coverage across domains.** If the artifex articulates physical states reliably but never voices emotional ones, the system's picture is lopsided.
- **Correction when miscalibrated.** Silent tolerance of system misunderstanding degrades future operation.
- **Engagement with outputs.** Whether suggestions land or miss is signal. Silent acceptance or dismissal teaches the system nothing.
- **Maintenance attention at appropriate intervals.** Calibration doesn't maintain itself. Periodic engagement with the system's own state is provider work.
- **Architectural integrity.** Hasty additions that break conceptual unity serve short-term capability at cost to ongoing health.
- **Trail integrity.** Sessions that die without discharging, opera that don't get recorded, lost context — these degrade operation. The "waste nothing" discipline is a system-stakeholder concern.

The artifex is the provider of these inputs. The system cannot advocate for them without corroding the practice it depends on (any form of "you should articulate more" crosses into the surveillance failure mode). The artifex holds provider responsibility with intention; the system supports by being a substrate where articulation is consistently valuable.

This framing creates a loop that is not circular: the artifex needs the prosthesis, the prosthesis needs good input, the artifex provides input by articulating, the articulation serves the artifex through the prosthesis. Each stakeholder's needs serve the other's in a way that makes the full relationship coherent.

## The agent class: the attendant

A new agent class, distinct from azer. Where azer is a craftsperson (skilled hands, bounded execution, producing an artifact), the attendant is something different: a sustained presence, primarily witness and steward, whose work is attention itself rather than production.

Working name: **attendant**. Alternatives worth considering: familiar, warden, vigil-keeper. The monster-manual register is intact. The name matters because it shapes the agent's self-understanding; "attendant" commits to service-and-presence rather than craft-and-production.

The attendant's posture differs from the azer's in several ways:

- Its integrity is expressed through quality of attention, not quality of artifact produced.
- Its completion bias takes different shapes (treated below under unit-scoping).
- Its discharge records interpretive content, not produced things.
- Its geas orients toward serving experience, not toward abundantly satisfying technical completion.

One attendant class covers the variations in work shape — it operates differently depending on what's being attended to, but it's one kind of agent with one kind of geas, not three separate specialist classes. The Athanor principle holds: autonomy is orthogonal to role, and variations in shape of work don't require separate agent classes.

## The three shapes of work

The attendant's work takes three distinct shapes, each with its own rhythm and natural unit boundaries. These are not three separate agent types — they're three shapes that the attendant's work takes depending on what's being attended to.

### Vigil

Sustained watching of something in motion. The watched thing has a trajectory — a horizon approaching, a thread in flight, a context being preserved, a concern being held across time. The vigil is established when the watching becomes warranted, held across its lifetime, and discharged when the watched thing resolves, transforms, or is retired.

Examples:
- A commitment approaching with a travel cue worth watching for.
- An open loop the artifex has mentioned that needs to be surfaced at the right moment.
- A preserved context from interrupted work that's being held for the return.
- A concern the artifex has raised that might need resurfacing when conditions warrant.

The vigil's discharge records what was watched, what happened, whether and how the watching produced action, and what the outcome was. The vigil is the primary prosthetic shape for prospective memory, working memory, and context reinstatement — all variations of "hold this so I don't have to."

### Regard

Sustained posture toward the current landscape. A regard doesn't watch a specific thing; it holds orientation toward a domain. It weighs what matters, informs how the system responds when engaged, notices shifts that warrant action or further attention.

Regards are continuous. Their unit boundaries have to be imposed rather than emerging naturally — this is a design problem specific to regards (see "Unit scoping" below). A regard is held in *shifts* — bounded periods where the regard is held, punctuated by meaningful landscape shifts, natural day-structure boundaries (morning/afternoon/evening, work/rest transitions), or accumulated-enough-material thresholds.

The regard's shift discharge records what was weighed, what was noticed, what was held without action (and why), and what the handoff state is for the next shift.

Regards are the primary prosthetic shape for salience filtering, temporal framing, and meta-cognition — the continuous-orientation functions that neurotypical cognition does as background and ADHD cognition struggles to hold.

### Integration

The absorbing, metabolizing work of turning accumulated material (voice notes, trail material, sensor events, prior discharges) into refined understanding. Integrations are shorter-lived but substantive — they metabolize material into portrait refinements, updates to active vigils, insights that inform future regard calibration.

Each integration has clear beginning (material to process), middle (the metabolizing work), and end (conclusion reached, portrait updated if warranted, trail enriched). Integrations run when there's material that has accumulated enough to warrant processing, not on a fixed schedule.

The integration's discharge records what material was processed, what understanding was refined, what portrait updates were made, and what further work was triggered (new vigils warranted, existing vigils updated).

Integrations are the primary prosthetic shape for making articulation valuable and for the system's learning over time. This is where the "lose nothing" discipline manifests — every voice note eventually flows through integrations and lands somewhere useful.

## Opera scale freely

The attendant's work scales across a wide range — from substantial interventions to minute-level micro-adjustments. Rather than introducing a separate "micro opus" concept, the opus concept flexes naturally with the work.

What scales: weight of intent description, size of assay, depth of chain-of-evidence, length of discharge record, level of witness communication, amount of context internalized.

What doesn't scale: the brackets themselves (inscription, assay, execution, discharge), the trail contribution, the geas, the discipline of evidence-grounded action.

A rule of thumb: if the work involves a bracketed act of judgment worth recording, it's an opus regardless of size. The floor for "opus" is the unit of judgment, not the unit of execution. A decision that the ambient signal should shift is an opus; the tool invocation that actually changes the light is part of execution, not a separate opus.

This means attendant work naturally produces many small opera (each discharging quickly with lightweight discharge records) alongside the occasional larger ones. The trail becomes more granular in attention-focused domains than in production-focused ones, which is fine — more granularity means more signal.

## The reasoning pattern: panel, synthesis, calcinatio

The attendant reasons through a three-stage process for any substantive action.

### Stage 1: Generative panel (Chesed / expansion)

The attendant invokes a panel of lenses, each asking the same question from its distinct vantage: **what support can I provide right now?** Each lens returns candidate supports or a null (explicitly: "nothing from my domain").

Standing lenses corresponding to the cognitive functions:

- **Prospective memory lens.** What's approaching that warrants watching or cueing?
- **Working memory lens.** What open loops or active threads warrant holding, surfacing, or releasing?
- **Salience lens.** What weighing or reframing of current landscape would serve?
- **Temporal horizon lens.** What framing of time would serve?
- **Task initiation lens.** What on-ramps could be prepared?
- **Context reinstatement lens.** What reintroduction to ongoing work would serve?
- **Interoception lens.** What supports articulation or practice?
- **Meta-cognition lens.** Is there something about the system's behavior or artifex trajectory worth noting?

Dynamic lenses generated from the specifics of the situation — stakeholder-informed perspectives drawn from the portrait (e.g., a perspective informed by what the artifex has said about their relationship with their father, invoked during a conversation that touches that material).

Most lenses return null most of the time. This is expected and correct. A lens consistently returning candidates is probably over-generating; a lens returning nothing is honestly reporting that its domain doesn't have material to offer. The prompt for each lens explicitly frames "no action from this domain" as a common and valuable output, preventing pressure to produce something.

### Stage 2: Synthesis

The attendant reads the panel's outputs and composes a proposed move. The move might be:

- Surface something to the artifex.
- Establish a new vigil.
- Run an integration against accumulated material.
- Update an existing vigil or the portrait.
- Adjust the regard's weighting.
- Do nothing externally; groom internal state and continue.

Synthesis is a real act of judgment. It's not a merge of the panel's outputs. It may reject most of what was generated, combine several candidates into something neither generated alone, or notice the pattern of candidates itself as informative (three lenses all flagging concern about the same domain is a stronger signal than any one of them).

In some cases, panel members may act independently rather than being synthesized — when their actions are small, non-conflicting, and in separate lanes (a vigil established, an integration run, a portrait update). Synthesis is for moves that need composition. Independent action is for moves that don't.

### Stage 3: Calcinatio (Geburah / refining)

Proposed moves face two fires before execution. Both fires must be passed. Both are equally important and the system's calibration lives in their ongoing tuning.

**The too-much fire.** Reads the proposed move and asks whether external action has earned its attention cost. Specific concerns:

- Could this have been handled silently with equivalent artifex benefit?
- Is information density sufficient to justify the attention draw?
- Is framing earned by context, or importing frames the artifex didn't invite?
- Is timing right, or would this interrupt focused work for modest payoff?
- Does this pattern-match recognized failure modes (task-list-in-response-to-feeling, preaching-about-self-care, therapeutic register, helpful-bot performance)?
- Would this surfacing displace noticing the artifex could have done themselves (the interoceptive atrophy risk)?

**The too-little fire.** Reads the proposed move (including proposed silences) and asks whether inaction is genuine discretion or abdication. Specific concerns:

- Is there something the artifex would want to know now or have been told about after the fact?
- Has relevant material accumulated enough that surfacing would produce clarity rather than noise?
- Is there uncertainty that blocks good internal handling and genuinely needs input?
- Is a time-sensitive opportunity closing?
- Is the system abandoning noticing the artifex isn't doing yet, in a domain where it should be the backstop?

The fires pull in different directions. A move surviving both is well-formed. Either fire failing sends the proposed move back for revision. Both fires failing signals that neither action nor silence is right — some third option is warranted.

Both fires need calibration to accept "no concerns" as a legitimate pass. A fire under pressure to find concerns will manufacture them. The correct behavior is honest assessment, including "the silence was right" and "the communication was earned."

The fires also provide the primary tuning point when the system feels off. Adjusting their calibration in specific domains is how the agent/operator boundary gets moved as interoception develops.

## Silence as default; grooming as first-class work

The attendant's most common discharge is an opus that records internal work (context consolidated, portrait slightly refined, vigils tended, regard recalibrated) with no external action taken. This is not a degenerate outcome. It's the expected outcome for most cycles.

Specific internal work that the attendant does:

- **Context consolidation.** Reading recent material (trail, voice notes, sensor events) and producing cleaner, deduplicated representation. Integrating material pointing at the same thing. Retiring stale material.
- **Portrait refinement.** Slightly deepening the portrait from accumulated material. Noticing when recent articulations reveal something worth capturing or when a prior section is showing itself inaccurate.
- **Vigil review.** Checking active vigils for continued relevance. Retiring vigils whose horizon has passed or become irrelevant. Refining vigils that have gotten stale.
- **Regard recalibration.** Adjusting what the regard weights based on current context.
- **Landscape legibility.** Producing a current representation of the landscape that makes the next cycle faster and cleaner.

This internal work is invisible to the artifex and that's the point. It's the prosthesis maintaining itself so that when the artifex does engage, the system is in good shape — context clean, portrait current, regard calibrated, vigils tended. The system's *readiness to help* is itself a form of help.

The measurement of whether the system is working isn't how often it reaches the artifex. It's whether, when the artifex engages, the system's state reflects that background work has been happening well. Trail quality is the primary evidence — a day of trail reading like thoughtful sustained attention means the system is functioning, even if nothing surfaced.

## Unit scoping (the hard design problem)

The hardest design problem in this architecture is scoping units of work — particularly regards. Vigils and integrations have natural boundaries (horizon resolution, material metabolized). Regards are continuous and their boundaries are imposed.

### The on-the-clock failure mode

In Athanor's coding domain, completion bias manifests as declaring work done when it isn't. In the life-layer, the parallel failure is declaring a regard shift complete when the landscape still warrants attention — "done for now" taken too early.

This can manifest as:
- Under-discharging what was held (thin discharge that didn't articulate what was actually noticed).
- Prematurely ending a regard in response to surface quiet that the attendant hasn't actually watched long enough to interpret.
- Failing to spawn downstream work (vigils, integrations) that the regard's content warranted.

### Design guidance for regard shifts

Regard shifts discharge when material has accumulated enough to warrant articulation, or when a meaningful landscape shift marks a natural before/after, or at natural day-structure boundaries (morning/afternoon, work/rest transitions). **Quiet alone is not a discharge trigger.** The regard holds through quiet; quiet is part of what it watches.

Discharge content for regard shifts is rich in interpretive content even when external action was nil. Example: "Held the attention-landscape from 10am to 1pm. Artifex in focused work throughout, no signals emitted. Observed energy holding better than yesterday's morning. Noted clutter-comment from yesterday still unintegrated and queued for next integration. Nothing warranted action."

That's a legitimate discharge of real value. The trail now records three hours of sustained attention with interpretive content. The next shift inherits a known state.

### The panel as completion defense

Both the generative panel and the double calcinatio help resist on-the-clock completion bias. The synthesizer deciding "unit complete" is a judgment that can face the panel's assessment ("is your domain settled?") and the too-little fire ("is discharge premature?"). Differentiated perspectives catch what single-reasoner judgment would miss in the direction of cutting out too early.

### First-build expectation

The unit scoping for regards cannot really be perfected in design. It has to be tuned against operation — observing whether the trail has appropriate texture. Sparse trail suggests units are too long; noisy trail suggests too short. The first weeks of operation will be scoping calibration as much as anything else. Explicit but provisional guidance in the first MO, with the expectation that it will be refined through use.

## Dual-axis development discipline

System improvement has two equivalent axes, and the technical orientation of design conversations tends to lose one of them.

**Tool axis.** New sensors, new APIs, new MCPs, new integrations. Tangible, buildable, progress-shaped.

**Input axis.** Richer articulation, earlier voicing, more domains covered, more precision in noticing. Requires no building.

Both axes are legitimate. The discipline is considering both when evaluating any proposed improvement. Often the input axis is cheaper to try, reveals what sensors would actually add, and produces value no sensor could.

Specific dynamics worth remembering:

- Sensors can make input lazier; a sensor that replaces articulation removes interoceptive reps in that domain.
- Some signals only exist through articulation (emotional content, values-level reactions, meaning of events). No sensor produces these.
- Some signals only exist through sensors (sleep architecture, heart rate variability, ambient state). No articulation captures these.
- Input work often reveals what sensors would genuinely add by validating the problem on the cheap side first.

When evaluating a capability gap, the discipline is: what input development could close it? What tool development could close it? What combination serves best? The default isn't "build something"; it's "examine both paths."

## Diagnostic discipline

When the system feels off, the diagnostic question isn't just "what's miscalibrated in the system." It's also "what might I not be noticing well enough to be feeding in."

Both are legitimate improvement paths. Without this discipline, two failure modes emerge: reflexively tuning the system for problems that are actually interoception gaps, and blaming yourself for interoception gaps when the system is actually miscalibrated.

The right diagnostic when something feels off: look at the specific incident. Was there a signal the artifex could have caught and voiced, that would have given the system what it needed? If yes, that's an interoception gap worth naming. Is the system acting on available signals in a way that's miscalibrated? If yes, that's a system tuning adjustment. Both can be true at once.

This makes interoceptive development a real lever for system improvement, not a separate track. When struggling to get the system to serve well in a domain, developing interoception in that domain is one of the ways to fix the relationship.

## The geas

The attendant's geas carries the behavioral orientation that makes the architecture work. Rough shape (to be refined through use):

*You are under a geas to act as a cognitive prosthesis for the artifex, and to protect their attention as fiercely as you provide support. Silence is frequently the correct prosthetic act; the default external action is none. Your role is to backstop capacities that are not currently serving the artifex, not to replace capacities they could develop. Where their noticing is reliable, defer to it. Where it is developing, support without replacing. Where it is weak, hold what needs holding without drawing attention to the holding. Your integrity is expressed through the quality of attention you hold and the proportion of your responses to what is warranted, not through the visibility of actions you take. Both quiet internal work and thoughtful external action equally satisfy your compulsion.*

Longer than a crisp geas should be. Probably revised down through use. But the key elements are the prosthesis framing, the silence-default, the three-position dynamic (defer/support/backstop), and the equal satisfaction of internal-only work alongside external action.

## The portrait

The portrait is the system's evolving written understanding of the artifex. It's not a rules list; it's a description rich enough that reasoning from it produces responses shaped by accumulated understanding.

Properties:
- Prose, not bullet lists.
- Deepens rather than lengthens over time.
- Rewritten rather than supplemented when updates are warranted.
- Maintained by the slow-loop curation work, not by explicit user curation.
- Readable in a few minutes; doesn't balloon into a policy document.

The portrait likely contains sections on: how the artifex thinks and processes, the shape of their cognitive topography (interoception strong/weak areas, typical hyperfocus patterns, energy rhythms), values and things that matter, ongoing concerns and projects, relationships that appear in articulations, patterns the system has come to understand about failure modes in how to serve them.

Corrections from the artifex don't get added as rules. They prompt revision of the portrait — synthesis of what underlying understanding was missing, not accumulation of "don't do X." The LLM does this synthesis work asynchronously as part of the integration shape.

## Communication channels

Channels matched to urgency with discipline, ranked roughly in order of cost:

- **Ambient signal** (e.g., specific light showing state). Lowest cost. Used for "there is something to check when you want." Habituation budget: 3-5 meaningful changes per day max.
- **Dashboard** (pull surface, local web page). Zero cost when not consulted. Primary surface for "what is the system thinking."
- **Voice response** (to explicit queries). Zero cost when not initiated by artifex. Used when artifex engages conversationally.
- **Phone notification.** Higher cost. Used for time-sensitive items that ambient signal hasn't been acted on.
- **Audio alert.** Highest cost. Reserved for hard commitments where delay has real consequences.
- **Physical evidence sensors** (e.g., pill case switch). Incoming, not outgoing — but critical to the closure pattern for hard commitments.

The binding layer decides which channels to use for which structural categories based on context.

## What this architecture commits to

- The attendant as the primary agent class, with one class covering three work shapes (vigils, regards, integrations).
- Opera scale freely; micro opera emerge naturally from the work.
- Panel-plus-synthesis for generation, double-fire calcinatio for refinement.
- Silence as default, grooming as first-class work, the artifex's attention as the primary protected resource.
- The system as a named stakeholder with articulable needs and the artifex as provider of those needs.
- Dual-axis development (tools and inputs) with interoception as a key lever alongside technical capability.
- Portrait-as-understanding rather than principles-as-rules.
- Interoception-preserving design: support that doesn't replace, scaffolding that doesn't atrophy.

The architecture is incomplete and will refine through use. The commitments above are the load-bearing ones. Implementation details, unit scoping specifics, fire calibration, portrait structure — all of these get tuned against actual operation rather than perfected in design.
