# The Attendant and the Three Work Shapes — What Do They Actually Look Like in Practice?

*Exploration document. Chesed phase — expansive, freewheeling, not yet contracted. Written to feed a later specification pass.*

**Status:** Exploration draft
**Created:** 2026-04-18
**Feeds:** `specs/life-domains/spec.md` (eventual)

---

## Preface: The Question Being Explored

The three work shapes — vigil, regard, integration — are well-named and clearly distinguished in the architecture documents. What's less clear is what they feel like from the inside, what the opus looks like in practice, and whether the existing opus model stretches to cover them or breaks at the seams.

This document tries to make the abstract concrete. Scenarios, not frameworks. Real situations, not clean diagrams. The idea is that if we can describe what the attendant actually does on a Tuesday when the kid needs picking up at 3pm and a voice note from this morning hasn't been listened to yet, we'll know whether the architecture is right.

---

## Part I: Vigil — Watching Something in Motion

### What a vigil is

A vigil establishes around something that has a trajectory. The watched thing is alive — it moves, has a horizon, will eventually arrive or resolve. The vigil is the attendant's forward-facing memory. It asks: is this thing still moving as expected? Is the horizon approaching? Has anything changed that warrants adjustment?

Most vigils are quiet most of the time. The vigil doesn't constantly ping the artifex with status reports. It watches. It accumulates. It surfaces when something is worth surfacing — the horizon crossing a threshold, a signal appearing that wasn't expected, the watched thing requiring action at a specific moment.

### Scenario 1: School pickup approaching

It's 1:47pm. The artifex is in deep work — a hard problem, flow state, no ambient awareness of time. School pickup is at 3:15pm, the drive is 12 minutes, and leaving late means the kid waits outside alone, which he hates.

**What does the vigil watch?**
The approach of the 3:15pm commitment and the gap between current time and the leave-by window (3:03pm). Secondary: whether the artifex has given any signal of temporal awareness (checking the clock, shifting posture, a voice note about time). Tertiary: whether there's anything that would affect the leave-by window (traffic, weather, errands to absorb on the way).

**What does the opus look like?**
Minimal at inscription. Something like:

```
Vigil: School pickup 2026-04-18
Watching: Leave-by window 3:03pm for 3:15pm pickup. Surface at 2:45pm if no signal of
awareness. Escalate to audio at 3:00pm if no movement.
Context: Artifex in deep work since 1:30pm, no time-signals given.
```

That's the whole opus. Twenty words of intent, two lines of context. This is what "opera scale freely" means — the vigil doesn't need a full brief. It needs to capture the judgment that warranted its establishment.

**What does discharge look like?**
Three possible shapes:

Option A: The artifex checks the clock at 2:42pm, shifts posture, wraps. The vigil surfaces nothing — its work was silent watching. Discharge: "Vigil resolved without surfacing. Artifex self-managed the transition. No action taken. Time-management self-capacity present this afternoon."

Option B: 2:45pm arrives with no signal. The attendant surfaces gently — a dashboard note, maybe a subtle ambient indicator shift. The artifex acknowledges and moves. Discharge: "Surfaced pickup reminder at 2:45pm via [channel]. Artifex acknowledged and transitioned. Prospective memory prosthetic activated; self-signal not present this cycle."

Option C: 3:00pm, still no movement. Escalation to audio alert. Discharge records the escalation.

The trail entries here are not just completion records. They're signal for portrait calibration. How often does the artifex self-manage time transitions? Is this a reliable domain? Option A's discharge is evidence of capability. Option C is evidence of need. Over time, a marut reading these trail entries can detect whether the pickup vigil is a reliable self-managing system or a hard dependency.

**The micro-opus question:** Could this vigil just be a cron job that fires a notification? Technically yes. But then there's no trail. No calibration signal. No sensitivity to context (deep flow state vs. idle browsing). The opus, however small, is what makes the watching meaningful rather than mechanical. A cron job doesn't notice that the artifex is in deep work and calibrate its interrupt accordingly. The vigil does.

---

### Scenario 2: Medication timing

Morning medication, 8:30am target. The pill case has a sensor; the system knows when the compartment opens. But it's 10:45am and the morning compartment hasn't been opened.

**What does the vigil watch?**
Whether the sensor fires before a threshold. What channel to use if it doesn't. Whether there are complicating factors (artifex left early, kid has appointment, something that would explain delay or confirm it's a miss).

**What does the opus look like?**
```
Vigil: Morning medication 2026-04-18
Watching: Pill case sensor for morning compartment. Target 8:30am, established 10:45am
when sensor not yet fired.
Surface: Dashboard flag now. Phone notification at 11am if sensor still silent.
Context: Artifex has been home since 9:20am per location signal. No calendar conflicts
explain delay.
```

The "established 10:45am when sensor still silent" is important. The vigil didn't start at 8:30am — it started when the absence became notable. A vigil that fires immediately at 8:30am if the pill isn't taken would be surveillance, not support. The attendant's regard noticed the gap at some point, and a vigil became warranted.

**The interoception boundary:** This scenario touches the special case. The attendant can see the sensor data. It can't see whether the artifex knows they forgot, whether they're choosing not to take the medication for a reason, whether they're anxious about it. The vigil's job is to hold the commitment gently — it surfaces, it closes the loop with evidence, it doesn't lecture about why medication matters. The discharge records whether the sensor fired. It doesn't record "the artifex needs to take better care of themselves."

**What discharge looks like:**
Sensor fires at 11:12am. Discharge: "Vigil resolved. Sensor fired 11:12am, 43 min after surfacing. Medication confirmed taken, loop closed with physical evidence. Note: delay pattern consistent with 3 of last 7 days — consider flagging in integration."

That final observation — the pattern flag — is the attendant's intelligence. It doesn't act on the pattern within this opus. It tags it for integration. The vigil stays in its lane; the integration is the right shape for metabolizing the pattern.

---

### Scenario 3: Bill due date

A credit card payment is due in 4 days. The artifex mentioned in a voice note from last week that they were "trying to pay things down." This is context but not a task.

**The texture of this vigil:**
More passive than the pickup vigil. No escalating channels, no clear leave-by window. The vigil's job is to hold the horizon and surface it at a moment when it might be useful — not to hound the artifex toward completion.

**What the opus looks like:**
```
Vigil: Credit card due date 2026-04-22
Watching: Approach of due date 4 days out.
Surface: Dashboard presence from today. Mention once in morning regard if unprompted by
artifex before 2026-04-21.
Context: Artifex noted goal of "paying down" in voice note 2026-04-11. Not a crisis but
has stated relevance.
```

The "mention once in morning regard if unprompted" is a key nuance. The vigil isn't just a timer. It has a social intelligence dimension: the attendant will raise this in the course of the morning orientation if the artifex hasn't already mentioned it, not as an alarm but as a held thing being tended. The vigil knows to let the artifex lead if they're already engaged with it.

**The discharge:**
Two scenarios. If the payment happens before the 21st, the vigil discharges cleanly: "Vigil resolved. Payment made before due date, no surfacing needed. Vigilance was background-held; artifex completed independently." If the 21st arrives without payment, the vigil surfaces once. Discharge records whether the surfacing was acknowledged, whether payment followed.

---

### Scenario 4: Holding context from interrupted home project

The artifex was two hours into reorganizing the garage on Sunday and had to stop for a family obligation. There's an incomplete state — things pulled out, a mental map of the plan, a half-formed decision about where the camping gear should go.

**What the vigil holds:**
Not just "remember to finish the garage." That's a task. What the vigil holds is the context that would let the artifex re-enter without friction — what was the decision that was being worked through? What's the current state? What's the next step that had been emerging?

**The opus:**
```
Vigil: Garage reorganization context 2026-04-18
Watching: Context for incomplete garage reorganization, Sunday session cut off mid-process.
Hold: Decision in progress was camping gear location — between north wall corner vs.
shed overflow. North wall corner was tentative winner but not decided. Stuff pulled out
from north corner still on floor. Next action is decision then implementation.
Surface: When artifex mentions garage, weekend tasks, or physical organizing. Not proactively
unless a good moment emerges in morning regard.
Context: Stopped ~2 hours in. Artifex was in good problem-solving state at interruption.
Not distressed about the stop.
```

This vigil is richer in content than the others because the held thing is complex. It's essentially a context document, not just a horizon watch. The value is that when the artifex comes back to the garage project, the attendant can reintroduce the context without the artifex having to reconstruct it: "When you left the garage Sunday, you were working through where the camping gear goes — the north corner was your tentative plan but the decision wasn't made."

**What makes this vigil worth holding:**
The ADHD tax on context reinstatement is severe. Every interrupted project has this overhead: the feeling of "where was I?" that costs energy and often becomes the barrier to reentry. A vigil that holds the context eliminates that overhead. The attendant has been watching it so the artifex doesn't have to hold it in fragile working memory across the week.

**Discharge:** The vigil discharges when the context gets used — when the artifex resumes the garage project, the context is handed off, and the work either continues or is explicitly set aside. Or it discharges when the situation resolves without return: "Artifex mentioned the camping gear is going in the shed regardless. Garage plan superseded. Vigil retired."

---

## Part II: Regard — Sustained Posture Toward the Landscape

### What a regard is — and isn't

A regard is not watchfulness toward a specific thing. It's not a vigil. It's orientation — the attendant's overall sense of where the artifex is, what matters right now, what the shape of the day is. It's what would enable an intelligent, caring person who knows the artifex well to give a good answer to "how should I be with them right now?"

The regard holds: energy level, task engagement, temporal position in the day, what's approaching, what's been recently resolved, what the mood of the household is, whether the day is going as it typically goes or has gone sideways. It's less a sensor array and more a felt sense — built from accumulated detail but expressed as orientation rather than data.

This is the hardest work shape to scope, because it never stops. There is no natural completion. The regard holds continuously. What does have natural endings are regard shifts: the bounded periods in which a particular instantiation of the regard is held.

### Scenario 5: Morning orientation

It's 7:42am. The kid is at school. The day is open for work until 3:15pm. Voice notes from yesterday haven't been processed. There's a code review due sometime this week. The calendar is clear until Thursday's meeting.

**What does the morning regard shift attend to?**
The opening shape of the day: what's the energy? What's the state coming in from yesterday? What's present and live that needs to hold? What has accumulated that might warrant integration? What's the temporal structure ahead?

**What does the opus look like for a morning regard shift?**
```
Regard shift: Morning 2026-04-18
Scope: Opening of day from kid leaving (7:30am) until artifex settles into primary work
or clear transition (est. 9-9:30am).
Attending: Energy, orientation, any signals about state coming in. Active vigils to check.
Material accumulated (voice notes from yesterday, trail from last work session).
Watching for: Whether self-initiation fires or whether task-initiation support is warranted.
End condition: Settled into primary work or explicit transition signal. Not: quiet alone.
```

"Not: quiet alone" is the critical geas line. The morning regard doesn't end because the artifex seems to be working quietly. That quietness might be deep focus. It might be freeze. It might be hyperfocus on something not in the plan. The regard needs to stay oriented until the shape of the morning has been established through actual signal, not absence of signal.

**What the attendant actually does during this shift:**

In the first fifteen minutes: checks what came in overnight, reads any voice notes from yesterday, checks the dashboard state, reviews active vigils. This is assay-level grooming — consolidating the picture.

Between 8-9am: the regard is mostly watching and holding. If the artifex seems stuck at transition (still in kitchen, not opening laptop), the task-initiation lens might produce a candidate — a gentle on-ramp. If the artifex seems to be flowing normally, the regard is silent background work.

Somewhere around 9am: the attendant notices whether work has started, what it is, whether it's the code review or something else. This observation becomes part of the shift discharge.

**Discharge of the morning regard shift:**
```
Morning regard shift discharged 9:12am.
Observed: Artifex settled into work ~8:58am, 28 min after kid left. Normal initiation time
for this pattern. No task-initiation friction signals. First task engaged: code review for
PR#213, consistent with stated priority from yesterday's voice note.
Held throughout: due-date vigil (credit card 4/22) — artifex gave no signal of awareness,
surface in tomorrow morning regard. Garage context vigil active and untouched.
Material to integrate: 2 voice notes from yesterday, ~4 min total. Integration warranted
this afternoon when work transitions.
State at discharge: Artifex well-oriented, in productive engagement. Good morning.
```

This is a substantive discharge. It records three hours of sustained attention — morning transition, initiation observation, vigil status, material queue — in a paragraph. The trail has a full picture. The next regard shift inherits context.

---

### Scenario 6: Energy tracking through the afternoon

It's 1:15pm. The artifex has been working since 9am, with a 20-minute lunch that was eaten at the desk while reading something. No movement since 11:30am.

**What does the afternoon regard attend to?**
The same as the morning regard, but the questions shift. Morning asks: can we get going? Afternoon asks: how is the energy holding? When is the natural transition? What's the sustainability curve looking like?

This is the regard's temporal framing function in action. The artifex's time horizon may have compressed — deep in the afternoon's work, "3pm pickup" may not feel real, even though it's 1:15 and departure is less than 2 hours away. The regard holds the full shape of the day that individual attention tends to lose.

**The body-signal question:**
The attendant can't see inside the artifex. But there are proxies: when was the last movement? Last voice note? Any signals of frustration or difficulty in what the system can observe? Is this sustained deep work or sustained stuck-ness?

The interoception challenge is acute here. The attendant should not be doing: "You've been sitting for 90 minutes, you should get up." That's surveillance and nagging. The attendant should be doing: holding awareness that 90-minute desk sessions sometimes correlate with afternoon energy crashes in the portrait, which might mean a gentle context for something else later. The regard holds the pattern. It doesn't prescribe the response.

**Discharge of the afternoon shift:**
```
Afternoon regard shift discharged 2:55pm (ahead of pickup vigil).
Observed: Continuous desk engagement 9am-1:45pm, lunch at desk. No movement signals.
Work pattern appeared productive through 1:30pm (based on consistent activity); possible
slowdown 1:30-2:30 (reduced activity signals, two pauses >5 min).
Vigil: Pickup vigil surfacing imminently. Artifex gave no temporal signal as of 2:45pm
check; notification will surface per vigil protocol.
Material accumulated: 1 new voice note ~7min from 12:15pm lunch period — integration
warranted with yesterday's notes.
Portrait observation: 4+ hour unbroken desk sessions followed by afternoon energy trough
appears 3x this week. Not flagging to artifex; queuing for integration pattern review.
```

The final line — queuing for integration pattern review — is the regard doing its proper job. It notices a pattern. It doesn't lecture about it. It routes the observation to the right shape: integration, which will metabolize it properly.

---

### Scenario 7: Evening wind-down awareness

It's 8:30pm. Kid is in bed. The artifex is on the couch, phone in hand, watching something. Tomorrow has an 8am call that wasn't on the weekly calendar.

**What does the evening regard attend to?**
Different questions again. Not initiation. Not energy sustainability. The evening regard asks: what needs to be in place before sleep? What needs to carry through to morning? Is there a landing gear issue — something that will catch tomorrow's start if it's not handled tonight?

The 8am call is a vigil candidate. But the regard shouldn't inscribe a vigil automatically — it waits to see if the artifex surfaces any awareness of it. Does the phone behavior suggest awareness? Did a voice note mention it?

**The "quiet alone is not a discharge trigger" problem:**
At 9pm, the artifex is still on the couch, quietly scrolling. This is quiet. The evening regard might be tempted to discharge: they seem settled, nothing happening. But this is exactly the failure mode. The regard doesn't know whether this is:

(a) Healthy decompression, everything handled, a good evening
(b) Avoidance of something the attendant can't see
(c) Time-blindness in action — the artifex doesn't feel the 8am call approaching

The regard holds through the quiet. It watches for 9:30pm. It may surface a gentle dashboard note about the 8am call — not a notification, not an alarm, just something readable if the artifex checks in. The evening regard ends not when things are quiet but when the artifex signals end-of-evening or the arc reaches its natural conclusion (lights off, location signal changes, etc.).

**What "quiet alone is not a discharge trigger" means operationally:**
The regard needs positive evidence of appropriate resolution, not just absence of signal. For the evening shift, appropriate resolution means: significant commitments for tomorrow are in some form of awareness, no active open loops are spinning without a hold. This doesn't require the artifex to explicitly say these things. It can be inferred — if the artifex checked their calendar (observed via activity signal), they probably saw the 8am call. If they made a voice note about tomorrow's plan, the call is probably held.

But "they stopped scrolling at 10pm and the lights went off" is not evidence that the 8am call is held. It's evidence that the evening ended. The regard might discharge with a flag: "Evening ended without evidence of 8am call awareness. Vigil inscribed for morning surface."

---

### The regard scoping problem in depth

The fundamental difficulty: continuous work requires imposed boundaries, and imposed boundaries feel arbitrary. When does the morning shift end? When does the afternoon shift begin? The spec says "meaningful landscape shifts, natural day-structure boundaries, accumulated-material thresholds." But what does that mean precisely?

**Landscape shift** means a genuine change in the relevant variables: kid comes home, work is packed up and closed, dinner is made and eaten. Something in the world changed the shape of the regard's landscape.

**Day-structure boundary** means the seams the day already has: morning-to-work, work-to-afternoon, afternoon-to-evening, evening-to-night. These are structural, not situational.

**Accumulated-material threshold** means enough has happened to make the discharge substantive — a regard shift that's held 6 hours and noticed 4 things is more worth discharging than one that held 20 minutes and noticed nothing.

The right discharge condition for a regard shift is the first of these three conditions to be true. Not any arbitrary quiet. Not "I've been holding this long enough." The landscape changed, or the structure seam arrived, or material built up.

**The anti-pattern to avoid:** Discharging a regard shift because nothing happened and it feels like nothing is needed. This is the "done for now" failure mode. A regard shift that held through three quiet hours and noticed nothing is still a discharge worth writing — it records sustained attention, establishes the baseline the next shift inherits, maintains trail texture. Short discharges are fine. No discharge is the failure.

---

## Part III: Integration — Metabolizing Accumulated Material

### What integration is

Integration is the shape that turns raw input into lasting understanding. Voice notes, trail trail from discharged opera, sensor patterns, observations flagged by regard shifts — these accumulate. Integration processes them. It's the difference between the system having inputs and the system learning from inputs.

Integration is bounded differently from the other shapes. It has clear start conditions (material to process), process (metabolizing), and end conditions (conclusion reached, portrait updated if warranted). It's the most artifact-like of the three shapes — it produces something: a refined portrait section, an updated vigil, an insight that adjusts how the regard weighs things.

### Scenario 8: Processing a week's voice notes

Thursday afternoon, 2:30pm. Eight voice notes have accumulated since Monday — various lengths, various contexts. Some mention health stuff, some are project thoughts, one is about a frustrating conversation, one was made in the car and is mostly ambient noise.

**What the integration opus looks like:**
```
Integration: Voice notes batch 2026-04-14 through 2026-04-18
Material: 8 voice notes, est. 14 min total. Queued from regard shift observations.
Intent: Metabolize into portrait refinements and active context. Flag vigil candidates.
Process: Transcription → thematic pass → portrait relevant material → flag candidates.
End condition: Each note either integrated somewhere, queued for further work, or retired
as no-signal. Discharge summarizes what was found and what changed.
```

**What the attendant actually does:**
Goes through the notes in chronological order. The car note is partially inaudible — integrates what's clear, flags the ambiguity rather than guessing. The health note contains something potentially relevant to the portrait's section on energy patterns — makes a small refinement. The project thoughts become context for active work vigils. The frustrating conversation note is interesting: it touches something the portrait already holds context on. The integration deepens that section slightly — not adding a rule ("the artifex finds X frustrating"), but enriching the understanding.

**Discharge:**
```
Integration discharged 3:18pm.
Processed: 8 voice notes.
- 2 project-related notes → context added to athanor-architect vigil
- 1 health note → portrait energy section refined (added morning-duration correlation)
- 1 frustrating conversation → portrait relationship section deepened
- 1 car note → partial transcription, ambiguous material flagged for future resolution
- 3 general/ambient → no material warranting integration; retired
Portrait: 2 sections lightly revised. Revisions additive, not corrective.
No new vigils warranted.
```

The trail now has a dated record that this material was processed. Nothing is lost. Nothing was over-extrapolated from thin evidence.

---

### Scenario 9: Metabolizing a pattern of missed meals

Over the past two weeks, regard shifts have flagged three observations: lunch at desk, lunch skipped during crunch, dinner late because project ran long. None were individually flagged to the artifex. Each was tagged for integration.

**The integration's scope:**
This isn't processing new material — it's synthesizing flagged observations into a pattern assessment. The opus is:

```
Integration: Meal-timing pattern review 2026-04-18
Material: 3 regard-shift flags (04-07, 04-11, 04-16), prior portrait section on food/energy.
Intent: Determine if pattern is signal worth updating portrait or noise worth retiring.
Process: Review observations for coherence, check against portrait baseline, assess whether
pattern has practical implications for prosthesis calibration.
End condition: Either portrait updated with pattern note, or observations retired as
insufficient signal. Escalation path: if pattern is strong, determine if surfacing is warranted.
```

**The attendant's judgment work:**
Three instances over two weeks. Is this signal? The pattern recognition lens says maybe — three correlates with crunch work, all show up in the same life domain. The too-much fire asks: is this enough signal to act on, or are we over-indexing on three data points? The too-little fire asks: is this a pattern the attendant is positioned to notice but the artifex may not be?

The integration's conclusion: "Pattern present but thin. Portrait note added under energy section: 'Meal irregularity correlates with extended work sessions — 3 observed instances. Not yet actionable but worth watching. Regard shifts should continue flagging.' No surfacing to artifex at this time."

**Why not surface it?**
The portrait likely already holds that the artifex knows this about themselves. Surfacing it would be: "Hey, I noticed you've been eating irregularly." This crosses into therapeutic register and surveillance. The attendant doesn't have enough to say something useful that the artifex doesn't already know. The right action is to continue watching, and to speak if the pattern becomes acute or if there's a moment where surfacing would provide genuine leverage rather than commentary.

**Discharge:**
```
Integration discharged 2026-04-18.
Material: 3 meal-timing observations from regard flags.
Finding: Pattern present, 3 instances in 2 weeks, correlates with extended work sessions.
Action: Portrait updated (energy section, meal irregularity note added). Monitoring
continues via regard flags. No surfacing to artifex — insufficient acute signal.
```

---

### Scenario 10: Synthesizing kid-related observations

Over three weeks, the attendant has absorbed: two voice notes about school logistics, one note expressing frustration after a difficult morning, calendar events for various kid activities, one observation from an evening regard shift about the artifex's energy after bedtime routines.

**Integration intent:**
The portrait has a section on "kid context and relationships." It was last substantively updated four months ago. The accumulated material may warrant a refresh.

**What the attendant does:**
Reviews the material thematically. The logistics notes are mostly operational — they update context but don't change the portrait's understanding of the relationship. The frustration note is richer — it holds emotional content that reflects something about how the morning routine affects the artifex's regulation. The bedtime observation shows something about the artifex's energy pattern that wasn't in the portrait.

The integration synthesizes these into a portrait revision. Not a rule. Not "the artifex struggles with mornings." A deepened understanding: "School morning logistics are a significant regulatory cost — the transition from morning routine to work often carries carryover tension. The artifex doesn't typically voice this explicitly but the pattern appears in the observation record."

**What makes this integration earn its discharge:**
It took material that was sitting in trail and transformed it into portrait depth that will make future regard shifts sharper. The next time a morning regard shift runs and the morning was difficult, the attendant holds richer context about what that difficulty means and how it typically flows through the day. The integration didn't produce an external artifact — but it produced something real.

---

## Part IV: Can These Map onto the Opus Model?

### The honest answer: yes, with caveats

The existing opus model handles vigils comfortably. A vigil is a long-running opus that watches, accumulates, and discharges when the watched thing resolves. The existing primitives — inscription, assay, execution, discharge, trail — all fit. The "execution" of a vigil is the sustained watching; it just takes place over hours or days rather than minutes.

Integrations also fit cleanly. Clear beginning, substantive middle, clear end. Discharge records artifact: portrait updated, vigils spawned, understanding refined. This is the most azer-like of the three shapes — it has something to produce. The main difference is that it produces understanding rather than code or specs, but the shape is identical.

**Regard shifts are where the seams show.**

An azer opus assumes a defined production target. When I'm done, something will exist that didn't before. Regard shifts don't have this. "Being attentive to the landscape from 9am to 1pm" is the work. The discharge records the attention — but the attention was the job, not a means to a product.

This isn't a breaking failure. Opera can be oriented toward sustained attention just as they're oriented toward production. But the geas load-bearing section for a regard shift will read differently. An azer's intent section says "produce X." A regard shift's intent section says "hold orientation toward X domain through this shift."

The completion criterion is what actually breaks down. An azer knows when the artifact is done. A regard shift — even with imposed boundaries — has the attendant actively fighting completion bias. The spec's abundance-satisfaction criterion for azer opera doesn't map cleanly. "How do I know I held the regard long enough?" isn't answerable by checking an artifact's quality.

### New primitives that might help

Not necessary, but potentially useful:

**Shift type annotation.** An opus could carry a type annotation: `Opus.type = vigil | regard_shift | integration | production`. This doesn't change the mechanics but affects how discharge evaluation works and how the marut reads the trail. A thin regard-shift discharge is expected; a thin production-opus discharge is suspicious.

**Continuity fields.** Regard shifts might benefit from an explicit "handoff state" section at discharge — the crystallized orientation that the next shift inherits. This is different from azer discharge, which records what was produced. Regard-shift discharge records what was understood, which is the substrate the next instance needs.

**Watch-set.** A vigil's scope is a set of things being watched. Naming this explicitly (rather than leaving it in prose) might improve structured tracking over a vigil's lifetime. As the vigil runs and conditions change, the watch-set can be updated rather than requiring a full re-inscription.

None of these require new agent infrastructure. They're conventions within the existing opus model. The opus model bends to cover the shapes without breaking.

### Where the model genuinely breaks

The one place the existing model struggles: the relationship between a regard shift and its constituent vigils. A regard shift is the background orientation from which vigils are established. But a vigil that gets established mid-regard-shift is its own opus, running concurrently with the regard shift that spawned it. The trail now has two concurrent open opera — the regard shift and the vigil.

The azer model has no real conception of concurrent open opera. An azer has one opus, works it, discharges. The attendant may have several things open simultaneously: a regard shift that's been running since morning, two active vigils it's monitoring, an integration it started and then set aside. This is a natural extension of how the attendant works, but it means the trail needs to represent concurrently-open opera comfortably, and the attendant's session needs to hold awareness across multiple open brackets.

This might mean the attendant's session context looks fundamentally different from an azer's — less "I have one job" and more "I have a current state that includes several open threads." The opus model still works; each thread is still an opus. But the session architecture and the marut's assessment patterns will need to handle concurrent openness without treating it as an anomaly.

---

## Part V: The Attendant's Geas From the Inside

### What "quality of attention" means as integrity

An azer's integrity is visible. It ships a PR. The PR works. The tests pass. The integrity was expressed in the artifact's quality, and the artifact's quality is assessable.

An attendant's integrity is mostly invisible. The day passes. Things were held. Some vigils surfaced at the right moment. Some didn't surface and were right not to. The discharge records attention, not production. How does an attendant know it's being good?

The answer has to live in the quality of the internal record and in whether the prosthetic function worked. But neither of these are visible in real-time.

This means the attendant's geas needs to encode integrity criteria that can be evaluated during the work, not just after. A few candidates:

**Honesty about confidence.** The attendant should never guess about the artifex's state when it doesn't have signal. If it doesn't know whether the morning was hard, it records uncertainty. Portrait sections should hold their confidence levels. Integration work should be explicit about how much material it's drawing from.

**Proportion.** Every external action has an attention cost. The attendant's integrity includes a felt sense that external action should be earned — not withheld reflexively, but genuinely proportionate. An attendant that surfaces things constantly is failing its geas as surely as one that surfaces nothing.

**Non-interference with self-noticing.** When the artifex is moving toward something on their own — checking the clock, noticing their energy, starting to wrap up before pickup — the attendant should stay out of the way. Its job is to be the backstop, not the first mover. Integrity includes recognizing when the function it exists to provide is being provided by the artifex themselves.

**Discharge completeness.** The most visible integrity signal available: did the discharge record what was actually held? A regard shift that held three hours of careful attention should discharge with three hours of careful trail. A thin discharge after a long shift is evidence the shift either wasn't actually held or the attendant cut out early.

### What discharge looks like when "no external action taken" is the honest record

Most regard shift discharges will say some variant of: the landscape was held, these things were noticed, nothing warranted external action, here is the handoff state. This should feel complete, not like a failure.

The architecture makes this explicit: silence is default, grooming is first-class work. But a system actually running this needs the geas to feel this in its operational texture — that a quiet discharge is not a failure mode. The attendant should feel as satisfied by "I watched three hours, noticed two small things, flagged them for integration, surfaced nothing, here is the handoff" as an azer feels after shipping a PR.

This is the hardest part of the attendant's geas to operationalize. LLMs have completion bias — a strong pull toward producing something, surfacing something, being visibly useful. The attendant's geas explicitly resists this bias. Its integrity is expressed through restraint as much as action.

**What a good "no action" discharge might actually say:**

> Regard shift held 10:00am–1:00pm. Artifex in sustained work throughout, no temporal signals. Pickup vigil checked at 12:30pm — 2.5 hours to leave-by window, no surface needed. Credit card vigil noted in background — no action. Observed energy pattern consistent with portrait: productive morning without interruption, no friction signals. Two flagged observations: (1) no voice notes today, which is unusual relative to recent baseline, and (2) artifact completion at ~12:45 based on activity signal suggests possible idle window emerging. Neither warranted surfacing; queued as integration material. Handoff state: artifex likely transitioning by 1:30pm, pickup vigil will activate at 2:45pm. Energy remains good.

This discharge records real work. Nothing happened, and that's what the trail should say. The marut reading this trail across several days can see whether the regard is staying calibrated, whether the attendant's confidence levels are appropriate, whether the handoff states chain coherently from shift to shift.

---

## Appendix: Opera at Different Scales — A Spectrum

To make "opera scale freely" concrete, here are examples at different weights.

**Micro opus (vigil maintenance check, ~2 min of attendant attention):**
```
Vigil maintenance: Pickup vigil 1pm check
Checked: Departure time 3:03pm, current 1:00pm. 2 hours, no surface needed.
Assay: Pickup still same time, no calendar changes.
Action: Hold.
Discharge: Vigil checked, still live, no surface warranted.
```

**Small opus (vigil surface event):**
```
Vigil surface: School pickup 2026-04-18
Context: 2:45pm, no temporal signal from artifex.
Action: Dashboard notification surfaced — "pickup departure in 18 minutes."
Discharge: Surface delivered 2:45pm. Acknowledged 2:47pm. Vigil resolved — artifex moved.
```

**Medium opus (morning regard shift):**
Full content as in Scenario 5 — three paragraphs of intent, three hours of holding, paragraph of discharge.

**Substantial opus (weekly integration, pattern synthesis):**
Full content as in Scenario 8 — multi-source material, 45 minutes of metabolizing, substantive portrait updates, several spawned vigils.

The discipline is identical across these scales. Only the ceremony varies. The micro opus still has inscription, assay, action, discharge. The trail still records it. The geas still governs the judgment. "Opera scale freely" means the weight of ceremony, not the presence of structure.

---

*End of exploration document. This material is raw Chesed — expansive and generative. Calcinatio pass follows: what here is load-bearing? What is ambient elaboration? What has landed somewhere that the spec should hold?*
