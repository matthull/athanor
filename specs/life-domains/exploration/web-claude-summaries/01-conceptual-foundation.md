# The Conceptual Foundation

*Why this project exists and what it's trying to be. Domain-independent thinking about attention, cognitive prosthetics, and the specific shape of ADHD cognition that makes conventional systems fail.*

---

## The core problem

Every notification has a cost, and for ADHD that cost is asymmetric and non-obvious. A neurotypical person gets a ping, processes it, returns to task with perhaps thirty seconds of refocus. For ADHD cognition, that same ping can cost the entire thread of thought, plus the meta-cost of recovering what the thread was, plus potential emotional dysregulation if the notification carries implicit judgment ("you still haven't done this").

The real unit of measurement isn't "notification" — it's intrusion-energy-cost versus value-delivered. Traditional productivity systems optimize the numerator. Nothing seriously optimizes the denominator. The design gap is an order of magnitude, not a percentage.

This means conventional solutions, including those explicitly marketed to neurodivergent users, cannot suffice by being incrementally better. They operate on assumptions about attention economics that are wrong for this population. A ten-times energy asymmetry between neurotypical and ADHD attention cannot be addressed by systems built on neurotypical premises, no matter how thoughtfully adapted.

## Key principles that fall out of the energy-asymmetry frame

**Attention is a budgeted resource, not a free one.** The system has to model attention explicitly. Every notification is a withdrawal. The budget varies by time of day, medication timing, sleep debt, current cognitive load. A well-designed system refuses to spend when the budget is low unless the value is extremely high. Default: no. Notifications have to earn their intrusion.

**State is primary, events are emissions.** Most systems are event-driven — something happens, a notification fires. A better architecture maintains a continuous model of what matters right now and emits events only when the cost of not-signaling exceeds the cost of signaling. This inversion means missed events aren't lossy — the state is still there next time you check. Traditional notifications are lossy: miss it, lose the information.

**Pull-first, push-as-last-resort.** A living document that re-grooms itself based on what still matters is architecturally different from a queue that accumulates. Checking it must be cheap and not checking it must be safe. A queue builds guilt. A living document doesn't, if stale items gracefully degrade or self-resolve.

**Ambient before explicit.** Peripheral perception can register change without attending to it. A subtle light-shift or ambient signal costs nearly nothing if you don't look at it, but remains available when you do. This is the right substrate for information you want available without being intrusive. Phone-style notifications are a last tier, not a first resort.

**Channels have habituation budgets.** Any signal loses meaning when used too often. The defense is rarity — three to five ambient signals a day, at most, before the channel becomes wallpaper. This means the system must be ruthless about what earns a signal, and channels must be matched to urgency with discipline.

**Hard commitments need physical-world evidence.** Medication taken is not the same as notification dismissed. Acknowledge-and-forget is a primary failure mode. Real closure of a loop requires evidence in the physical world (the pill case opened, the meeting attended) — not a button press. This is one of the largest departures from conventional systems, and it's central.

## Cognitive prosthesis as the organizing frame

The deeper reframe: this isn't a productivity tool. It's not a notification system. It's a cognitive prosthesis for a brain whose specific functions are under-resourced compared to the neurotypical baseline.

A prosthesis compensates for a capacity that isn't fully available. This framing immediately produces better design questions than "what features should it have?" The question becomes: *which specific cognitive functions are impaired, and what shape of support does each require?*

Crucially, prostheses have to avoid causing atrophy in the tissue they support. A poorly-designed prosthesis creates dependence and weakens what remains. A well-designed one compensates without replacing the underlying function where development is still possible. This distinction matters enormously for the life-layer system and governs several downstream design choices.

## The cognitive functions at stake

Each of these points to a different shape of support and a different literature of research.

### Prospective memory

Remembering to do something in the future when a cue arrives. Neurotypical brains do background pattern-matching: seeing the keys triggers remembering to grab the library book. ADHD prospective memory is unreliable in a specific way — the intention exists, but the cue doesn't fire retrieval.

Design implication: inject cues at the right moment based on context, not on clock time. Event-based cues work much better than time-based cues for ADHD. This literature (Einstein, McDaniel) is directly applicable.

### Working memory

Neurotypical working memory holds 5-7 items fluidly. ADHD working memory is smaller and more volatile — items drop out, and when they drop, the fact that they dropped also drops. The system is literally external working memory.

Design implication: render at working-memory-compatible resolution. A dashboard that shows 20 items replicates the failure mode rather than compensating. The system surfaces 3-5 items at appropriate granularity; the rest is accessible but not presented.

### Salience filtering

Neurotypical brains automatically tag incoming stimuli with "relevant-now / not-relevant-now." ADHD brains have flatter salience gradients — novel and emotionally charged items hijack attention regardless of actual importance. ADHD is driven by interest, novelty, challenge, and urgency rather than abstract importance (Dodson's interest-based nervous system).

Design implication: the system pre-does salience work. The dashboard reflects its judgment of what's relevant-now rather than presenting undifferentiated lists. Prioritization isn't a user task; it's system work.

### Temporal horizon / time perception

Neurotypical brains have fairly accurate visceral senses of time. ADHD time is more binary: now, and not-now. This is sometimes called "time blindness." Standard calendars show time spatially but don't help it *feel* temporally.

Design implication: make time-until-event viscerally present. Make not-now items actually recede from view rather than sitting in the same list as now. Don't pretend scheduled times produce the same felt urgency they would for a neurotypical user.

### Task initiation

Neurotypical brains transition from "I should do X" to "doing X" with relatively low friction. ADHD has a massive activation barrier, especially for non-interest-driven tasks.

Design implication: surfacing a task too early or without a concrete next physical action wastes a notification. Tasks benefit from activation-reducing scaffolding embedded in how they're surfaced ("open laptop, the draft is started with this subject line") rather than presented as abstract items.

### Interoception

See extensive treatment below. The support here is different from the others because interoception is a skill that can be strengthened rather than a capacity that should be externally replaced.

### Context reinstatement

Neurotypical brains return to interrupted tasks and re-load context fairly quickly. ADHD context reload is expensive and often fails — you come back to the document and genuinely don't remember what you were doing.

Design implication: capture context at interruption points so return is cheap. "You were writing section 3, your last thought was about the edge case with X, here's the sentence you were mid-writing." This is known from programmer-interruption research (Parnin & Rugaber) but almost nobody builds it into personal tools.

### Meta-cognition and self-monitoring

Neurotypical brains run background checks: "am I still on the right task? am I making progress? should I switch?" ADHD self-monitoring is weaker, which is why hyperfocus can consume hours on something unimportant.

Design implication: the system can externally run these checks, but delicately — it can feel like nagging. The right shape is probably occasional gentle interruption rather than continuous monitoring.

## Compensatory versus restorative support

Not all cognitive deficits should be compensated for the same way. The rehabilitation medicine distinction between *compensatory* (work around the deficit) and *restorative* (rebuild the function) strategies applies directly.

Most of these functions — prospective memory, working memory, context reinstatement — should be fully externalized. The brain can't do them well and fighting that is waste. The prosthesis takes over and doesn't apologize for it.

Interoception is different. It's a skill that can be strengthened, and externalization actively harms it. This makes it a special case and the most important distinction in the design space.

## Interoception as the special case

Interoception is the perception of internal bodily states — hunger, thirst, fatigue, muscle tension, the physical substrate of emotion, the subtle "something is off" signal before you consciously know what. ADHD impairs it through several compounding mechanisms:

**Attention allocation.** Interoception requires sustained, low-grade attention to quiet internal signals. ADHD attention is biased toward external novelty and high-salience stimuli. Internal signals lose the competition for attentional resources.

**Salience network function.** The insula, part of the salience network, is responsible for tagging stimuli as worth attending to. ADHD appears to involve atypical salience network function. Internal signals get weakly tagged and don't rise to awareness.

**Stimulant medication effects.** Stimulants can blunt interoceptive signals further — particularly hunger and fatigue. This is mechanistically straightforward (dopaminergic effects on appetite, alertness masking fatigue) and worth knowing about.

**Hyperfocus as active suppression.** During deep engagement, the ADHD brain can actively suppress interoceptive signals to maintain task focus. You emerge from three hours of flow genuinely surprised you're starving, stiff, and need the bathroom. The signals were present; they were suppressed.

**Developmental vocabulary gap.** Growing up with impaired interoception means fewer reps of linking sensations to labels. The learned mapping from "this feeling in my stomach" to "anxiety" or "hunger" or "needing the bathroom" is thinner than neurotypical.

**Emotional interoception.** Related to alexithymia — difficulty identifying emotions. Often in ADHD, the body signal is present but doesn't translate into an emotional label. This is a primary path to emotional dysregulation: frustration builds without being noticed until it spills over.

### Why interoception can't be fully externalized

Sensors can report "heart rate elevated" or "no meal in seven hours" or "voice tension pattern suggests frustration." These are external signals masquerading as internal ones. If the system reliably tells you when you're hungry before you notice, you develop learned dependence and the underlying interoceptive skill atrophies further. This is analogous to how GPS has measurably degraded natural navigation ability — the scaffolding replaces the skill.

The goal isn't to tell you your internal state. It's to help you build and maintain the skill of knowing your internal state. This is a different design target from all the other cognitive functions.

### What the system can do for interoception

**Prompts to check in, without revealing the answer.** The system can invite noticing without providing the answer. "It's been three hours, how's your body?" The system may know from sensors that heart rate is up and you haven't moved. It doesn't say so. It invites. After you report, it can share what the sensors show — creating a feedback loop you calibrate against.

**Catching hyperfocus suppression.** When sustained focus has gone on long enough that interoceptive suppression is likely, gentle interruption can break the suppression — not by telling you what you missed, but by creating space where the signals can emerge on their own.

**Reducing noise so signals are detectable.** Part of why interoception is hard is bad signal-to-noise ratio. When you're spending all your attention tracking open loops, the quiet internal signals get drowned out. The system's broader job of reducing cognitive load indirectly improves interoception by freeing attentional bandwidth.

**Post-hoc calibration.** "You said you felt fine at 2pm. By 4pm you were clearly dysregulated. Want to look at what the sensors showed between?" This is training data for your own perception. Over time you learn what 2pm-you was missing.

**Making space for practice.** Interoception genuinely improves with practice. The system's role might be as simple as protecting morning check-in time, noticing when practice has been skipped for a week, treating the practice itself as a hard commitment.

### Real-time journaling as interoception practice

The key operational insight: a pipeline for capturing voice notes in real time — whenever a thought arises — is itself interoception practice, done naturally.

When you voice "just grabbed a granola bar, didn't have a real meal today, carb-heavy," several things happen at once:

- The thought is externalized and stops consuming working memory.
- Rumination is prevented — unvoiced thoughts of this kind recur, each time demanding a small re-processing.
- The act of articulating reinforces the noticing. You're naming the action, the absence, the pattern.
- The system receives rich first-person context no sensor could produce.

This produces interoception practice disguised as utility. You're not doing a body scan; you're just voicing what you notice. Over time the noticing sharpens because articulation rewards precision.

Critically: this is *your* practice, not something the system prompts you into. The rhythm belongs to you. The system is the silent recipient that makes articulation feel consequential by actually using what you say. Rhythms that depend on external prompting are fragile. Rhythms that are yours persist.

### The anti-surveillance constraint

The system should not interpret your voice notes back to you in ways that feel like monitoring. "You've voiced three notes about fatigue this week, you might be burning out" sounds useful but would destroy the practice, because now you're aware that the system is watching and evaluating. The notes become performative. You'd start editing what you say.

The distinction is subtle but real: the system using voice notes as context for reasoning (responding differently because of what you said earlier) is fine — you feel *remembered*, which is different from *surveilled*. Holding your observations up as mirror, unsolicited, crosses the line.

The heuristic: remembered good, surveilled bad. The system can answer questions you explicitly ask about your patterns. It doesn't volunteer them.

## Structure-and-channel separation

One concrete architectural commitment drops out of all this thinking.

**Structure** is what a piece of information is and how it should behave: how urgent, how important, how long it remains relevant, what counts as resolution, whether it needs physical evidence or just acknowledgment, what it's related to.

**Channels** are delivery surfaces with properties: ambient light (cheap, peripheral, habituation-prone), dashboard (pull-only, high-bandwidth, zero intrusion), phone (persistent until dismissed, high-cost), audio (immediate), physical evidence sensors (completion-detecting).

**Binding** is the policy that maps structure to channel based on current context — your location, your activity, your attention budget, recent channel usage, the item's properties.

Most systems conflate these. The result is that you can't easily change how something is delivered without changing what it is, and the same piece of information can't adapt to different contexts. Separating them lets the same content be delivered ambiently on a calm day and escalated to audio on a chaotic one. It also means channels can be evolved independently of content and vice versa.

## The overall shape of what a better system looks like

- Maintains a continuous model of current state rather than reacting to events.
- Treats attention as a budgeted resource and every action as a withdrawal.
- Defaults to silence; emits only when silence would leave the artifex worse off.
- Uses ambient channels as the primary output, with escalation tiers reserved.
- Requires physical-world evidence for hard commitments (meds, appointments).
- Supports interoception through articulation rather than replacing it.
- Operates from understanding of the person rather than from rules.
- Adapts calibration as the person's interoceptive capacity develops.
- Treats context loss as a primary failure mode and prevents it structurally.

This is the stance the system is built from. The architectural choices that follow are in service of these commitments. Where conventional design would default to more features, more notifications, more structure, this design defaults to less — ruthlessly — because the cost of each addition is higher than conventional thinking accounts for.
