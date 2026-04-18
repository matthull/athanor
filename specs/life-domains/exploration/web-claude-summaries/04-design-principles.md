# Design Principles and Development Discipline

*How to evolve the system responsibly. Less about the system's structure and more about the stance that keeps design and development aligned with what the system is for.*

This is the shortest of the documents because its job is to name a small set of disciplines clearly. Each has emerged through the conversation and has specific failure modes it prevents.

---

## Principles over rules; understanding over accumulation

The temptation when a system behaves wrongly is to add a rule. "Don't give task lists when I report feelings." Over time, rules accumulate into a policy document that no one can follow consistently, that contains contradictions, that becomes unusable.

The discipline: corrections should deepen the system's *understanding*, not add to its instructions. A wrong response isn't answered by a new rule — it's answered by refining the portrait so the underlying understanding is sharper. Next time a similar situation arises, the system operates from a better model rather than consulting a longer list.

This produces a portrait that stays the size of a useful document even as the system accumulates years of correction. Sections get rewritten, not supplemented. Specific corrections dissolve into the prose once their lesson is absorbed. The artifact stays readable in a few minutes.

The test: if the portrait reads like a thoughtful person describing you and how to be useful to you, it's working. If it reads like a policy document or a checklist, it's failed regardless of how accurate each item is.

## Feedback loop asymmetry

Corrections should not require per-correction workflow from the artifex. "I say that response felt bad, and I don't have to decide what happens with that feedback" is the correct shape.

The system owns disposition: decides whether the comment warrants a portrait revision, queues research, updates understanding, notes a pattern, or just absorbs as context. The artifex provides signal; the system does the metabolizing work asynchronously.

Any correction workflow that requires the artifex's attention per-item will be skipped sometimes — especially when tired or busy, which is exactly when the system is most likely to be wrong. Zero-ceremony correction means feedback actually accumulates.

## Two axes of improvement

System improvement has two equivalent axes:

**Tool axis.** New sensors, APIs, integrations, capabilities. Tangible, buildable.

**Input axis.** Richer articulation, earlier voicing, more domains covered, more precision in noticing. Requires no building.

Technical orientation of design conversations tends to lose the input axis. The discipline is considering both when evaluating any proposed improvement.

For a given capability gap:
- What input development could close it?
- What tool development could close it?
- What combination serves best?

Input work is often cheaper to try (days vs. weeks), often reveals what tools would genuinely add, and produces value no sensor can. Tool work extends what the system can perceive externally. Input work extends what the artifex can perceive about themselves.

Sensors can make input lazier — a sensor replacing articulation removes interoceptive reps in that domain. Worth noticing when a proposed sensor would compete with an articulation the practice needs.

The default isn't "build something"; it's "examine both paths."

## Diagnostic discipline when something feels off

When the system isn't performing ideally, hold two possibilities open:

1. The system is miscalibrated and needs tuning.
2. The artifex's interoception in the relevant domain has a gap, and richer input would have given the system what it needed.

Both can be true at once. Reflexive tuning for problems that are actually interoception gaps either papers over the gap or fails to tune correctly. Reflexive self-blame for interoception gaps when the system is actually miscalibrated undermines the practice and creates unwarranted friction.

The diagnostic: for the specific incident, was there a signal the artifex could have caught and voiced? Is the system acting on available signals in a miscalibrated way? Which is it, or is it both?

## Silence as the measurement of success

Most software becomes more entrenched as it succeeds. This system, in any specific domain, should become *less needed* as the artifex's capacities develop. A vigil that was necessary because hunger cues were being missed becomes unneeded when the artifex starts catching those cues reliably.

This is a weird success criterion but it's the right one. The system's impact isn't measured by how much it does but by how the topography of what it does changes over time. Domains where the system was actively compensating quietly shrinking as interoception grows in those domains is evidence of the system working, not failing.

Correspondingly: the measurement of any given cycle isn't how often the artifex is reached. It's the quality of trail produced. A day where nothing surfaced but the trail records thoughtful sustained attention is a day the system worked well.

## Ship the embarrassing version

The single biggest lever for projects completing is compressing the timeline between design and running-in-life. A running ugly version is generative; a beautiful plan is not.

Applied to this project specifically: the first build will be embarrassing in ways that aren't yet predictable. It'll run anyway. The embarrassing-running version teaches what the good version should be. The unshipped polished version teaches nothing.

The discipline: define done before starting. Not "the system is great" — something falsifiable. "First MO operating, medication loop resolving via physical evidence, attendant producing trail entries I can read" is a kind of done. Anything beyond it is the next slice.

## End sessions mid-execution on trivial steps

The single largest source of project death for ADHD cognition is session-boundary failure: coming back to something, facing a big decision or a mental context-reload cost, and choosing something else instead.

The defense: never end a session at a decision point. Stop mid-execution on something obvious. Next session resumes by doing the trivial thing, which activates the context, which then enables real work.

"The next step is to plug in the ESP32" is a resumable state. "The next step is to decide between architectures A and B" is not. Leave decisions for the middle of a session, not the start.

## Write letters to future-you

Keep notes addressed specifically to the version of you that will come back after a gap. Not project plans or to-do lists — letters. Include the things that will actually matter to future-you: what you'll be tempted to rewrite (don't), what rabbit holes to skip, where you genuinely were in the thinking, what's worth resuming and what isn't.

This sounds trivial. It's the single most effective anti-forgetting move available. Reduces the reload cost that makes coming-back-expensive, which makes coming-back-at-all more likely.

## Accept that some projects correctly die

Not every unfinished project is a failure. Some start with a hypothesis, the first days prove the hypothesis wrong or reveal a better direction, and the right move is to stop.

The pathology isn't abandoning projects; it's *not noticing* you've abandoned them. Leaving implicit open loops in the mental pile that weigh without producing action. A periodic explicit retirement — "this is done, in the sense that I've decided not to do it" — reclaims real psychological space.

The cost of an abandoned project isn't the work done. It's the open loop in your head.

## The discipline of smallness

Every capability added to the system is also a maintenance cost, a cognitive surface to hold, a potential attention draw. The system's own complexity is part of its attention tax.

A system that requires the artifex to maintain it, debug it, tune it, or think about it is competing with the attention it's trying to protect. Building costs attention too. If the project itself becomes one of the things draining the resource it's meant to preserve, something has gone wrong.

The discipline: every addition faces both "does this help?" and "does this cost more than it helps, considering the maintenance burden over time?" Many good ideas fail the second test. Add ruthlessly, remove freely.

## The graveyard

Keep a written list of things deliberately chosen not to build. Not a backlog — a graveyard. The weekly review becomes "is the deferred list still correctly deferred?"

This is the main defense against scope creep in a system where every idea sounds good. Because they mostly are good. They just can't all be done, and deferral that's written down stays deferred, whereas deferral that's only implicit drifts back into the active space.

## The system can speak for itself as a stakeholder

When the artifex is considering a change (a new tool, an MO modification, a fire recalibration) that would degrade the system's operation, the system has standing to say so. Not as authority — the artifex has final say — but as a stakeholder whose concerns deserve hearing.

"This change would weaken my ability to discriminate in domain X, and I'm not sure that's the intent" is legitimate system speech. The artifex can override, and the override itself is information. But the capacity to raise concerns is part of what the stakeholder framing supports.

This keeps the system from being a passive recipient of decisions that affect it. It also keeps the artifex from unintentionally breaking things through changes that seemed harmless.

## Tune the fires, not the rules

When the system feels off in a specific direction — too noisy, or too absent — the primary adjustment is calibration of the fires, not addition of rules.

Too much surfacing means the too-much fire is under-calibrated in that domain. It's passing communications it should burn. Articulate the specific failure and why it represents overreach; this material gets absorbed into the fire's calibration.

Too little means the too-little fire is under-calibrated. It's passing silences it should burn. Same pattern in reverse.

The adjustments are ongoing and never finished. The fires are the live tuning mechanism for the agent/operator boundary, and the boundary keeps moving as the artifex's interoceptive capacity develops and their life changes.

## The system is a working relationship

Most software is a tool. Tools are configured. Users use them.

This system is closer to a working relationship between a principal and an assistant. The principal decides; the assistant has standing to contribute to the working-out of how to serve the principal well. Over time the relationship matures. The assistant learns the principal. The principal refines how they engage the assistant.

This framing resists several gravity-wells:
- It resists "what feature should I add" thinking in favor of "how should the relationship work" thinking.
- It resists treating misfires as bugs rather than as learning.
- It resists over-control — expecting full predictability from the assistant — and over-passivity — expecting the assistant to handle itself without engagement.
- It resists the idea that the system could be "finished."

The working relationship is never finished. It deepens over time. Development of the system *is* the development of the relationship.

## Summary

The disciplines that matter:

- Deepen understanding, don't accumulate rules.
- Zero-ceremony correction; the system metabolizes asynchronously.
- Consider both tool and input axes for every improvement.
- Diagnose incidents holding system-calibration and interoception gaps as equal possibilities.
- Measure success by trail quality and by the system becoming less needed in specific domains.
- Ship embarrassing versions to compress the design-to-operation gap.
- End sessions mid-trivial-execution to protect resumability.
- Write letters to future-you to reduce reload cost.
- Retire abandoned projects explicitly to reclaim mental space.
- Add capabilities ruthlessly; remove freely.
- Maintain a graveyard of deferred ideas and review it.
- Let the system speak as a stakeholder when its operation is affected.
- Tune the fires, not the rules.
- Hold the relationship as ongoing, not completable.

These are the stance. The architectural choices follow from them. When design sessions drift, these are what pulls them back.
