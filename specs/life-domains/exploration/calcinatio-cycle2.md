# Calcinatio Report — Cycle 2
## Dialectical critique of the life-domains exploration corpus

**Status:** Calcinatio — applying refining fire to cycle-1 Chesed output
**Created:** 2026-04-18
**Source documents:**
- `tipheret-principle.md`
- `first-slice-candidates.md`
- `orientation-layer-architecture.md`
- `attendant-work-shapes.md`
- `spec.md` (core spec for context)

---

## 1. What Holds Up Well

### The Tipheret principle itself is the strongest thing here

The tipheret-principle.md document is doing real conceptual work. The distinction between pruning (Geburah acting alone) and distillation (Tipheret) is crisp and genuinely useful. The formulation "the full vision wearing its smallest clothes" does what it needs to do: it names a quality that's hard to reach but easy to recognize when you see it.

The process description — full Chesed, then calcinatio, then wait for structure to emerge, then build that — is sound as a methodology for system design, and the failure modes section (imbalanced Chesed = perpetual refinement, imbalanced Geburah = premature contraction) names real traps with real symptoms.

The "wait for structure to emerge, don't force resolution" guidance at the end of the document is the most mature thing in the entire corpus. It correctly identifies that premature resolution is harder to fix than no resolution, and that agents should escalate rather than choose a form when the right form isn't visible. This should be load-bearing in subsequent specs.

### The attendant's work shapes are well-concretized

The scenarios in `attendant-work-shapes.md` — school pickup vigil, garage context hold, morning regard discharge — do the actual spec work. Abstract architectural documents say "vigil is watching something in motion." This document shows what that means on a Tuesday with competing demands. The gap between these two levels of description is where most system designs fail; this document closes it.

The spectrum of opera scales (micro vigil check through substantial integration) is particularly valuable. "Opera scale freely" is a principle; the appendix gives it operational texture that can actually guide the attendant when building.

### The three-layer separation is architecturally sound

The decision to keep domain implementations independent, shaping only the orientation layer without touching core, is the right call. The orientation-layer doc's identification that every mechanism operates through existing slots (role file slot, skills slot, MO slot, config slot) is coherent and internally consistent. The "vocabulary, not grammar" framing is clean.

### The double-fire calcinatio is a genuine contribution

The too-much/too-little fire is not just a restatement of existing calcinatio. It addresses a structural asymmetry that the coding domain doesn't surface: most calcinatio in azer work is "does this artifact meet quality criteria?" — which is effectively only a too-little fire. An artifact either meets the bar or it doesn't. An attendant action can fail in both directions. The double-fire names something real and the orientation-layer doc's suggestion that it might generalize back to coding work is worth taking seriously.

### The regard-scoping problem is honestly named

The work-shapes doc correctly identifies regard scoping as the hardest design problem and doesn't pretend to have solved it. The three discharge criteria (landscape shift, day-structure boundary, accumulated-material threshold) are honest partial guidance rather than false certainty. "What might surprise us" section in orientation-layer-architecture.md is similarly honest about what can only be learned from operation.

---

## 2. What Needs Refinement

### The Tipheret test: the recommended first slice doesn't quite pass it

The first-slice-candidates.md recommendation is morning attunement session plus medication vigil, 5-6 days, full build plan. Apply the Tipheret test: does this slice express the full cognitive prosthesis vision in its smallest form?

The answer is: partially. The morning session exercises the attendant concept and the regard shape. The medication vigil exercises one vigil. But the first slice as specified does not exercise:

- Integration (the third work shape is absent)
- Double-fire calcinatio as a live pattern (there's no natural calcinatio moment in a morning session or medication vigil — nothing to apply too-much/too-little fire to)
- The portrait (explicitly out of scope)
- Concurrent opera (the spec notes this as a model gap but the first slice doesn't surface it as a test)

This is not fatal. A first slice doesn't have to exercise every element. But the Tipheret framing makes a specific claim — "the full vision is entirely present in this slice" — and that claim overstates it. The medication vigil expresses prospective-memory support; the morning session expresses salience filtering and temporal horizon. Six of the eight cognitive functions are not exercised. "Nothing is compromised" is not accurate.

The honest framing is: this first slice proves that the attendant class can function, that the regard shape is real, that one vigil type works, and that the orientation layer has actual behavioral effect. That's enough. The full vision doesn't need to be present in the first slice. But the spec text shouldn't claim it is when it isn't — that claim will be load-bearing and it's soft.

### The 5-6 day estimate is optimistic given what actually needs to be built

Day 1-2: attendant.md is a "medium-high effort writing task" per orientation-layer-architecture.md, estimated at "1-2 days of careful writing." The first-slice plan allots Day 1 to this plus a live morning session. That's not realistic — attendant.md, written to azer.md quality, is a multi-day writing project on its own. The comparison is instructive: azer.md is over 300 lines of carefully-tuned behavioral guidance that has been revised through months of use. Producing a comparable document in one day produces something thinner than what the architecture requires.

The 5-6 day estimate also assumes the life MO template, the morning skill, the medication vigil lifecycle design, Telegram integration, and trail discipline all work from first-draft. Real systems don't work this way. The "first live session" on Day 2 will surface problems with the attendant.md draft that require another iteration before the sessions are useful.

A more honest estimate: 8-10 days for a first slice that genuinely works, or 5-6 days for a first slice that runs but is substantially miscalibrated and needs another 3-4 days before the done criterion is met.

### The "yes, and..." balance — the surfacing-too-much failure mode is named but not protected against

The attendant-work-shapes.md document identifies "proportion" as an integrity criterion and "completion bias" as the failure mode. The orientation-layer doc names the helpful-bot gravity well. These are the right names for the right problem.

But the documents don't provide a structural protection against the cumulative version of this failure. Individual decisions are handled by the double-fire calcinatio: each proposed action passes too-much fire before it surfaces. But what about the attendant that applies double-fire conscientiously to each individual decision and still ends up surfacing more over time as it gets more familiar with the landscape and finds more things worth surfacing?

The too-much fire is applied per-action. The cumulative surfacing rate is not checked anywhere. An attendant could pass every individual too-much fire and still be, in aggregate, taking more and more of the artifex's attention. This is a real failure mode — the attendant getting incrementally more helpful in a way that costs more than it gives — and it's not addressed structurally.

The integration shape is the natural place for this check: an integration opus could include a calibration review ("am I surfacing at the right rate relative to baseline?"). But this needs to be explicit in the attendant's geas or in the `/life-orientation` skill, not implicit in the double-fire.

### The "regard shift as opus" tension isn't fully resolved

The work-shapes doc correctly identifies that regard shifts break the azer opus model's completion assumption. The proposed resolution — "the opus model bends without breaking" — is too easy. A regard shift that runs from 9am to 1pm is held by a single Claude Code session. Sessions have context limits. A 4-hour regard shift with continuous attention will hit context pressure. The document doesn't address what happens when a long regard shift runs into context exhaustion: does it discharge early? Does a new session pick it up? How does the handoff state work mechanically?

This isn't fatal to the architecture, but calling it "resolved" by convention (shift type annotation, continuity fields) is premature. These are naming conventions, not solutions to the underlying problem that a session is bounded.

---

## 3. What Doesn't Survive

### The claim that `ath muster` supports role file selection — it doesn't

The orientation-layer architecture doc states: "The muster command passes it as the agent's role file instead of azer.md." This is the foundational claim for the attendant mechanism. It does not currently hold.

Reading `muster_cmd.go`: the boot prompt is hardcoded:
```
"Read %s/AGENTS.md, then read %s, then read %s/azer.md. Your opus is at %s. Read it and execute."
```

There is no `--role` flag. There is no way to substitute a different role file through the existing CLI. The architecture document says "Core changes required? No." This is incorrect — getting an attendant to load `attendant.md` instead of `azer.md` requires either:
1. A `--role` flag in `ath muster` (a CLI change = core change)
2. The attendant.md including an instruction to load itself at the start of a session (not how boot prompts work)
3. Creating the `attendant.md` inside the instance directory and having the boot prompt auto-discover it (would require changing the boot prompt logic)

The workaround is not complicated — the marut can construct a custom boot command that references `attendant.md` instead of `azer.md`. But this would bypass the standard `ath muster` command entirely, which is a worse architecture outcome than adding a `--role` flag to `ath muster`. The orientation-layer doc's confidence that no core changes are needed is wrong on this specific point.

This is the most consequential finding from the architecture boundary test. Everything else about the orientation layer works through existing mechanisms. This one does not.

### The "concurrent open opera" claim isn't examined

The work-shapes doc identifies that an attendant may have concurrent open opera (regard shift plus active vigils running simultaneously). It notes this as a "model gap" but asserts "the opus model still works; each thread is still an opus." This claim isn't examined.

The actual problem: the marut's stall detection is based on "no output in 10 minutes → nudge." An attendant running a morning regard shift has no output for an hour, which is correct. The orientation-layer doc acknowledges this: "A life-athanor marut.md variant may be needed." But this contradicts the "no core changes" framing — a marut.md variant is in the instance directory, so it doesn't touch core. However, the marut code itself lives in the marut's behavior, not in CLI code. The question of whether the marut's core instructions can be varied enough through an instance file is unexamined.

### The Tipheret principle's own formulation fails its test when applied to system design

The tipheret-principle.md is a 2,000-word exploration document. The distilled form of the principle is in its last three paragraphs. Everything before that is the Chesed phase elaboration — correct in spirit, too long to serve as agent guidance. If the principle is "express the full ambition through the smallest surface area," the principle's own document should demonstrate that. It doesn't. This isn't a major criticism — it's an exploration document, not a specification — but the irony is worth naming for cycle 3.

---

## 4. Cross-Document Tensions

### The "new class vs. shaped azer" question is answered differently in different documents

The orientation-layer doc argues for a new role file as the pragmatic resolution: write `attendant.md` as a fresh document that composes with shared infrastructure. This is the right call architecturally.

But the first-slice-candidates doc wavers — it uses "attendant-in-residence" as a concept but doesn't commit to whether the first morning session requires a full attendant.md or whether a well-prompted azer with the right MO would do. The evaluation matrix scores Candidate 1 "High" on "Exercises Orientation Layer" but the description could be achieved by a vanilla azer with a life-domain MO. The orientation layer's behavioral shift comes from attendant.md, not from the morning session design alone.

The tension: if the first slice doesn't build attendant.md (too much to build in Days 1-2 properly), what's running is a shaped azer, not an attendant. That's fine but it should be named honestly. The architecture is only "proven" when the attendant's geas is genuinely different from the azer's, not when the azer is doing morning stuff.

### "No core changes required" vs. the muster boot prompt

Addressed above. The orientation-layer doc says no core changes needed; the actual CLI code shows the boot prompt hardcodes `azer.md`. These cannot both be true.

### First slice scope: 5-6 days vs. "medium-high effort" attendant.md

The orientation-layer doc says attendant.md is "1-2 days of careful writing and iteration." The first-slice plan allocates Day 1 for it. These estimates are incompatible unless Day 1 means "start a draft" and a subsequent revision happens after Day 2's live session — but the plan doesn't say this.

### The "done criterion" and the Tipheret principle

The first-slice doc's done criterion is: "three consecutive mornings where the attendant surfaces something the artifex would have otherwise missed, without surfacing anything the artifex would have preferred not to see."

This is a good empirical criterion. But it measures one fire (the regard shape surfacing correctly). It doesn't test the too-little fire (did the attendant sit on something it should have surfaced?), nor does it test trail quality, nor vigil lifecycle, nor any integration work. The Tipheret framing in the same doc says "the full vision is entirely present" — but the done criterion doesn't check for the full vision. The done criterion and the Tipheret claim are misaligned.

---

## 5. Missing Elements

### Privacy for life-domain data — the spec names it but the exploration docs ignore it

The core spec states: "Privacy: life domain data stays local or in trusted services — no exposure through trails, logs, or shared systems." This is listed as a constraint. None of the four exploration documents address it.

This matters because:
- The morning session surfaces email content, calendar items, medication status. The trail records this. Where does the trail live? Is the athanor's trail directory excluded from any sync services?
- The portrait (when it exists) holds sensitive longitudinal data: energy patterns, emotional observations, health context. How is this protected?
- The attendant's integration discharges (Scenario 9: meal pattern, Scenario 10: kid frustration notes) contain intimate material. If these trail entries sync somewhere, that's a privacy breach.

The fact that the village repo privacy model isn't discussed in any exploration doc is a gap. The spec constraint is present but uninstantiated. Cycle 3 needs to answer: what is the trail directory, where does it live, and what protection does it have?

### Alex and Vivian — the family witness is absent from exploration

The core spec lists the family witness as a prospective reviewer (spec.md: "Is the family witness seeing downstream effects (follow-through)?"). The exploration docs focus entirely on the artifex.

The first-slice framing says "the attendant serves the artifex" and that's correct. But the family witness is relevant to at least:
- Kid logistics vigils (Vivian's pickup schedule is relevant; Alex's awareness affects household coordination)
- Hard commitments that affect others (doctor appointments, school events)
- The meta-question of whether the system creates visible downstream value that family members can observe

The work-shapes doc has a kid logistics scenario (school pickup) but the scenario only considers the artifex's experience. Alex isn't in it. The family witness frame — present in the core spec — is absent from all four exploration docs. This is a gap, not a fatal flaw, but the second slice will hit it.

### The graveyard — deliberate non-builds are not named

There is no list of things explicitly decided not to build. The out-of-scope section of the core spec lists several things, but that's different from a graveyard — a living record of ideas that were genuinely considered and deliberately set aside with rationale.

Candidates for the graveyard that appear in the exploration corpus and should be formally retired:
- Continuous monitoring / ambient surveillance of the artifex's biometrics
- Automatic vigil creation from voice notes without attendant judgment
- Cross-athanor awareness (life ↔ work) in the first slice
- Therapeutic register: the system should not provide emotional guidance, therapeutic framing, or analysis of the artifex's psychological patterns

These aren't out of scope by omission. They're candidates that were considered and should be retired with a reason, so future readers know they were examined and rejected.

### Interoception — named but not operationalized

The core spec mentions interoception as a special case (support without replacing). The first-slice-candidates doc mentions it in Candidate 5 (voice note routing). The work-shapes doc mentions it twice in passing.

But nowhere in the four documents is "support without replacing interoception" operationalized into guidance. What does the attendant actually do differently when interoception is involved vs. other domains? The distinction is named but the behavior isn't specified. The risk: an attendant built without this guidance will fall into the over-interpretation trap (detecting body signals and commenting on them) because nothing in its geas tells it not to.

The difference that needs to be specified: the attendant can hold what was articulated (voice note content), but it cannot and should not interpret or infer somatic experience from behavioral proxies. "You've been sitting for 90 minutes" is observable. "You're probably dysregulated" is not for the attendant to say. This distinction needs to be in attendant.md's geas, not just named as a principle in a spec document.

---

## 6. Recommendations for Cycle 3

### Priority 1: Resolve the muster architecture gap before committing to a build plan

The most load-bearing unresolved issue: `ath muster` does not support role file selection. Before cycle 3 contracts a build plan, it must decide:

**Option A:** Add `--role` flag to `ath muster` (small core change, clean solution). Boot prompt parameterized: `...read %s/<role>.md...`. This is the right architecture — a week of work, backwards-compatible, and opens the door to future role variants. The resistance to core changes is a useful discipline, not an absolute. This one earns its place.

**Option B:** The marut constructs a custom boot command that names `attendant.md` explicitly, bypassing `ath muster`. Architecturally messier — it means the marut's instructions in the life instance have to know the boot prompt format. This couples the marut to CLI implementation details.

**Option C:** The attendant.md is written to load itself when the default boot prompt is used — a CLAUDE.md instruction that redirects to attendant.md instead of azer.md. This might work but is fragile and non-obvious.

Option A is the right answer. Cycle 3 should plan it as a pre-first-slice infrastructure task.

### Priority 2: Write attendant.md as the core work of the first slice, not a Day 1 deliverable

Cycle 3 should separate the first-slice build into two phases:
- **Phase 1 (3-4 days):** Write attendant.md to full quality. Not a draft — a document as carefully tuned as azer.md. No live sessions until the geas is correct. The life MO template can be developed in parallel.
- **Phase 2 (3-4 days):** First live morning sessions. Medication vigil. Trail discipline. Iteration on what the attendant actually does vs. what was specified.

This reframes the timeline to 7-8 days with realistic quality expectations, which is more honest than 5-6 days that produces something needing substantial revision.

### Priority 3: Add the cumulative-surfacing-rate check to the attendant's geas

The double-fire handles per-action decisions. The integration shape should include a standing calibration opus: weekly or biweekly, the attendant reviews its own surfacing behavior. How many times did it surface something? What was the artifex's response pattern? Is the rate drifting up over time? This review pattern needs to be in the attendant's geas, not in a skill that might not get loaded.

This is the structural protection against the cumulative too-much failure mode that the current exploration docs don't provide.

### Priority 4: Contract the privacy model

Before any trail entries about medication, emotional observations, or family context are written to disk, the privacy model needs to be specified:
- Trail directory location and sync exclusion status
- Portrait file access permissions
- Which categories of integration discharge are excluded from backup systems
- Whether the attendant has any obligation to signal when an integration contains sensitive material

This is a one-afternoon specification task that prevents a class of problems that are hard to fix retroactively.

### Priority 5: Write the graveyard

One page, deliberately. For each retired candidate: what it was, why it was considered, why it was set aside. This document protects against rediscovering bad ideas in cycle 5, and it's a forcing function — things that haven't been formally retired are still in play.

### Priority 6: Ground the Tipheret claim for the first slice

The first-slice spec should stop claiming the full vision is present in the first slice. The honest Tipheret claim for the morning session + medication vigil first slice is:

"This slice proves that the attendant class is real, that the regard shape functions, that one vigil type works end-to-end with evidence-based closure, and that the orientation layer produces behavioral difference from vanilla athanor. It does not exercise integration, does not exercise double-fire calcinatio live, does not exercise portrait management. These are second-slice concerns. The first slice is the smallest thing that demonstrates the system is not a figment."

That claim is true, defensible, and sufficient. The overstated claim creates a mismatch between design aspiration and build reality.

### For the work-shapes doc: answer the context-exhaustion question

The regard-shift-as-opus model needs a concrete answer for context pressure:
- What happens when a long regard shift runs out of context?
- Is there a discharge-on-context-pressure discipline?
- How does the handoff state ensure continuity across sessions?

This can probably be addressed through the attendant.md's session architecture guidance rather than requiring new opus primitives. But it needs to be addressed before the first live session, not discovered when a morning regard shift runs long.

---

## Summary Assessment

The four documents represent a genuinely strong Chesed phase. The concepts are right. The architecturally sound elements — three-layer separation, attendant as new role file, double-fire calcinatio, the concrete work-shape scenarios — should survive into the specification.

The fires burn away: the false confidence that no core changes are needed (muster does need a role flag), the overstated Tipheret claim for the first slice, the optimistic timeline, and the gap on cumulative-surfacing drift.

The too-little fire finds: privacy model absent, family witness absent, interoception not operationalized into geas guidance, no graveyard.

The architecture is coherent and buildable. The path to cycle 3 is clear: close the muster gap, write attendant.md at full quality, add the cumulative-surfacing discipline, contract the privacy model. Then build.

---

*This document is calcinatio output — it applies refining fire, not generative expansion. Material here enters the spec only after explicit artifex review. The recommendations section is the highest-confidence output; the cross-document tensions section names what needs explicit resolution before building starts.*
