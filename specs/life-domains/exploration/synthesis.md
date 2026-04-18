# Life-Domains Exploration — Synthesis Document

**Status:** Cycle 3 contraction — refined output of three-cycle exploration process
**Created:** 2026-04-18
**Source cycles:**
- Cycle 1: `tipheret-principle.md`, `first-slice-candidates.md`, `orientation-layer-architecture.md`, `attendant-work-shapes.md`
- Cycle 2: `calcinatio-cycle2.md`
- Cycle 3: This document

*This is the contracted output. It is shorter and sharper than any individual exploration doc because the fires have already run. Material that did not survive calcinatio is absent. What remains is what holds.*

---

## 1. Consolidated Architecture

### Three Layers (confirmed, stable)

```
┌─────────────────────────────────────────────┐
│           DOMAIN IMPLEMENTATIONS            │
│  (home automation, calendar, meds, finance) │
│  Each in own repo/service/MCP server        │
├─────────────────────────────────────────────┤
│            ORIENTATION LAYER                │
│  attendant.md + life MO template            │
│  + orientation skills + instance CLAUDE.md  │
├─────────────────────────────────────────────┤
│           CORE ATHANOR SYSTEM               │
│  MO, opus, geas, calcinatio, trails,        │
│  agent roles, CLI (domain-agnostic)         │
└─────────────────────────────────────────────┘
```

The orientation layer manifests through four composable mechanisms — all filling existing slots, none changing core grammar:

1. **`attendant.md` role file** (highest leverage) — the attendant's geas, the three work shapes as first-class patterns, double-fire calcinatio as native verification
2. **Life MO template** — state-based intent, system-as-stakeholder witness, double-fire calcinatio section, 4-hour landscape freshness threshold
3. **Orientation skills** (`/life-orientation`, `/double-fire`) — operational specifics loaded at session start, keeping `attendant.md` clean
4. **Instance CLAUDE.md** — portrait location, voice note pipeline, domain routing, available sensors

These are vocabulary, not grammar. The three-layer separation holds in practice: domain implementations plug in through existing MCPs; the orientation layer shapes how the system relates to those domains; core runs unchanged.

### The Attendant: Role File, Not Shaped Azer

The attendant is a new role file (`attendant.md`) local to the life athanor instance. The reasons this matters:

- The azer's craftsperson frame leaks even when told to ignore it. "Establish your verification plan" in `azer.md` is about artifact verification; an attendant working from this frame will perform worse than one whose native frame is interpretive work.
- The geas must be reformulated, not extended. Silence-as-default is a compulsion, not a rule addendum. An attendant that has internalized "silence is the primary action" behaves differently from an azer told "also consider silence."
- The work shapes — vigil, regard, integration — need to be as foundational in the attendant's formation as opus lifecycle is in the azer's.

The attendant.md does not touch `shared/azer.md`. It composes with shared infrastructure (AGENTS.md geas, opus mechanics, delegation toolkit) while re-orienting the posture. Both files live in the life athanor's directory.

**This requires `ath muster --role` (see Section 3).** The current boot prompt hardcodes `azer.md`. This is the one confirmed core change.

### Three Work Shapes and the Opus Model

**Vigils and integrations map cleanly onto the existing opus model.** A vigil is a long-running opus that watches, accumulates, and discharges when the watched thing resolves. An integration has clear start conditions, substantive process, and clear end conditions — the most azer-like shape. These require no new primitives.

**Regards have seams.** A regard shift's completion criterion doesn't map to artifact quality — "held orientation from 9am to 1pm" has no natural done state. Two conventions are needed:

- **Shift type annotation** (`Opus.type = vigil | regard_shift | integration | production`) so discharge evaluation can apply appropriate expectations. A thin regard-shift discharge is expected; a thin production discharge is suspicious.
- **Continuity fields** — a regard shift's discharge records handoff state (the crystallized orientation the next shift inherits), not just what was produced.

One genuine structural issue remains: regard shifts run for hours with no output, while the marut's stall detection fires after 10 minutes of quiet. A life-athanor `marut.md` variant (instance-level, not a core change) will be needed to adjust monitoring posture. This is a tier-3 concern but should not surprise anyone when it arrives.

The concurrent open opera problem — a regard shift running simultaneously with multiple active vigils — is not fatal to the model but is unresolved. Each concurrent thread is still an opus. The attendant's session context holds awareness across multiple open brackets. Specifics will only be answerable from operation.

### The "Yes, And..." Energy — Held Differently

The azer's generosity flows toward production: generous effort → output. The attendant's generosity flows toward attention: generous presence → holding.

Both are "yes, and..." energy. But the failure modes are mirror images. The azer's completion bias is premature shipping — the attendant's completion bias is premature surfacing. This distinction is load-bearing in how `attendant.md` must be written. The geas cannot merely say "apply the double-fire." It must make the felt satisfaction of silent holding as real as the felt satisfaction of shipping — because an attendant whose integrity is invisible (the day passed, the vigils held, nothing needed surfacing) will find something to surface unless its geas has genuinely internalized restraint as completion.

The double-fire handles individual decisions. But there is a cumulative drift problem the double-fire does not address: an attendant that passes every individual too-much fire and still trends toward surfacing more over time. The structural protection belongs in the integration shape: a standing calibration opus, weekly or biweekly, in which the attendant reviews its own surfacing rate. "How often did I surface? What was the response pattern? Is the rate drifting upward relative to baseline?" This calibration must be explicit in `attendant.md`'s geas — not a skill that might not get loaded.

---

## 2. Refined First Slice

### Recommendation: Morning Attunement + Medication Vigil

The first slice is morning attunement session (Candidate 1 from the exploration corpus) with the medication vigil (Candidate 3) embedded. This recommendation survives calcinatio with corrections applied.

### Honest Tipheret Claim

The exploration corpus overstated this. The corrected claim:

> This first slice proves that the attendant class is real, that the regard shape functions end-to-end, that one vigil type works with evidence-based closure, and that the orientation layer produces behavioral difference from vanilla athanor. It does not exercise integration, does not exercise double-fire calcinatio live, does not exercise concurrent opera at scale, and does not exercise portrait management. These are second-slice concerns. The first slice is the smallest thing that demonstrates the system is not a figment.

What the slice specifically validates:
- The attendant class distinction is real and necessary (proved when attendant behavior diverges from what a shaped azer would produce)
- The regard shape works (morning shift has natural boundaries, discharges with real content, handoff state chains across days)
- Prospective memory and temporal horizon prosthetic functions are operational
- One vigil type (medication, hard commitment with physical-adjacent evidence) works end-to-end
- The orientation layer is architecturally present — not bolted on, but structurally real

What it does not prove: it exercises two of the eight cognitive functions (prospective memory, temporal horizon) directly, with salience filtering implicit in the morning regard. It does not exercise working memory via integration, meta-cognition via portrait synthesis, or task initiation via activation scaffolding. "Full vision present" is not accurate; "foundational architecture proven" is.

### Done Criterion (Calcinatio-Refined)

The exploration corpus's done criterion — "three consecutive mornings where the attendant surfaces something the artifex would have otherwise missed, without surfacing anything the artifex would have preferred not to see" — is a good empirical test for the too-much fire. It is insufficient as the full criterion.

Complete done criterion:
- Three consecutive mornings where the attendant surfaces something the artifex would have otherwise missed (too-little fire passes)
- Zero instances of the attendant surfacing something the artifex would have preferred not to encounter (too-much fire passes)
- Trail entries are substantive — each morning regard shift and medication vigil produces a discharge that records what was held, not just whether action was taken
- Medication vigil closes on evidence, not on acknowledgment — the Telegram confirmation pattern works mechanically and the attendant records closure correctly
- The attendant's silences are as legible as its surfacings — a morning with nothing to surface produces a discharge that says so, not a thin or missing entry

### Realistic Timeline: 7-8 Days

The 5-6 day estimate in the exploration corpus was incompatible with the acknowledged effort for `attendant.md`. Corrected:

**Phase 1 — Infrastructure and Writing (Days 1-4):**
- Day 1: `ath muster --role` flag implementation (required before any live sessions)
- Days 2-3: `attendant.md` — full quality, not a draft. The comparison to `azer.md` is instructive: that document is 300+ lines of carefully-tuned behavioral guidance. The attendant's geas, work shapes, silence-as-primary-action discipline, cumulative surfacing calibration, and interoception boundary all need to be written to this standard. No live sessions until the geas is correct.
- Day 4: Life MO template + `/life-orientation` skill first draft + instance CLAUDE.md setup

**Phase 2 — Live Sessions and Iteration (Days 5-8):**
- Days 5-6: First live morning sessions. The `attendant.md` draft will have problems the live session reveals. Iteration before the next session.
- Day 7: Medication vigil — recurring vigil lifecycle, Telegram confirmation integration, morning session integration. First days of live use.
- Day 8: Trail discipline, discharge patterns, iteration from a week of real sessions. Decide what's embarrassing-and-acceptable vs. embarrassing-and-needs-fixing.

---

## 3. Pre-Slice Infrastructure

The following must exist before the first live session.

### `ath muster --role` Flag (required)

The current boot prompt is:
```
"Read %s/AGENTS.md, then read %s, then read %s/azer.md. Your opus is at %s. Read it and execute."
```

There is no mechanism to substitute `attendant.md`. The clean solution is Option A from the calcinatio report: add `--role` to `ath muster`, parameterizing the boot prompt. `azer` is the default, preserving backward compatibility. The marut, reading the life athanor's `muster.md`, knows to invoke with `--role attendant` for life-domain opera.

This is a small core change. The calcinatio confirms it earns its place: every workaround (custom boot commands in the marut, CLAUDE.md redirection) is architecturally worse. Add the flag.

### `attendant.md` First Draft (life athanor instance)

Not a prototype. A document written to `azer.md` quality from the start. Required sections:
- Geas formulation (silence-as-primary-action as compulsion, not rule)
- Three work shapes as first-class operational patterns (vigil, regard, integration)
- Double-fire calcinatio as native verification
- Regard shift discipline (completion bias named and resisted, shift type annotation, continuity fields)
- Interoception boundary (articulated clearly — see Section 4)
- Cumulative surfacing calibration (the standing integration pattern)

The attendant.md cannot be produced in one day. Planning for two days of careful writing is correct.

### Life MO Template

One template covering: state-based intent (not deliverable-based), system-as-stakeholder witness (with its legitimate concerns: trail integrity, portrait currency, calibration maintenance), double-fire calcinatio section, 4-hour landscape freshness threshold, and attendant-specific abundant satisfaction criterion ("this domain is attended without being managed").

### Instance Setup

- Instance CLAUDE.md documenting: portrait location, voice note pipeline location, available MCPs, domain routing conventions, life-athanor marut configuration notes
- Directory structure for the life athanor: `~/athanor/athanors/life/` with instance files in place
- Privacy configuration (see Section 4): trail directory not synced, portrait excluded from backup scope

---

## 4. Gaps Requiring Artifex Input

These are not architecture questions. They are decisions that require the artifex's judgment and cannot be resolved by further exploration.

### Privacy Model for Life-Domain Trail Data

The core spec states the constraint but the exploration corpus does not instantiate it. Trail entries from morning regard shifts will contain: email content, calendar details, medication status, emotional observations, family context. Integration discharges will contain intimate pattern material.

Decisions needed:
- Where does the life athanor's trail directory live, and is it excluded from all sync services (rclone, Syncthing, etc.)?
- Does the portrait live at the same location or separately, with different access controls?
- Which categories of integration discharge are excluded from any backup scope?
- Does the attendant have any obligation to flag when an integration discharge contains sensitive material?

This is a one-afternoon specification task that prevents a class of problems that are hard to fix retroactively. It must happen before the first live session writes trail entries.

### Family Witness Operationalization

Alex and Vivian are listed as prospective witnesses in the core spec ("Is the family witness seeing downstream effects?"). The exploration corpus focuses entirely on the artifex. The first slice does not require resolving this, but the second slice will surface it — kid logistics vigils directly involve Vivian's schedule and Alex's awareness.

The question is not whether family witnesses exist but how the system relates to them: What do they see? Through what channel? What's the privacy model for information that implicates family members? Does the attendant ever act on their behalf, or only serve the artifex's relationship to them?

### Acceptability of the `--role` Core Change

The calcinatio confirms this is a core change. The CLAUDE.md constraint says core changes require artifex review. The calcinatio's assessment is that this change earns its place — small, backward-compatible, and architecturally cleaner than every workaround. But the artifex needs to make this call explicitly, not have it assumed.

### The Graveyard: What to Formally Retire

Several candidates from the exploration corpus should be formally named as "deliberately not built, with reason" rather than left as unaddressed possibilities. Retiring them now prevents rediscovery in cycle 5. Candidates for explicit retirement:

- Continuous biometric monitoring / ambient surveillance of the artifex's physical state
- Automatic vigil creation from voice notes without attendant judgment
- Cross-athanor awareness (life ↔ work) in the first or second slice
- Therapeutic register: emotional guidance, psychological pattern analysis, somatic interpretation

The artifex should confirm which of these are genuinely retired and what the reason is for each. The graveyard is one page; each entry gets a sentence on why it was set aside.

### Interoception: How Explicit Should the Geas Be?

The core spec names "support without replacing interoception" as a principle. The exploration corpus mentions it multiple times without operationalizing it.

The distinction that needs to be in `attendant.md`:
- The attendant can hold what was articulated (voice note content, explicit observation)
- The attendant cannot interpret or infer somatic experience from behavioral proxies
- "You've been sitting for 90 minutes" is observable and can be held silently
- "You're probably dysregulated" is not for the attendant to say
- Pattern observations about energy or physical state queue for integration review — they do not surface directly to the artifex unless the artifex has already named the pattern themselves

The question for the artifex: how explicit should this prohibition be in the geas? A strong formulation ("the attendant does not comment on the artifex's somatic state under any circumstances") vs. a softer formulation ("the attendant holds observations about physical state for integration, not for surfacing") have different operational implications. The strong version is safer; the softer version may miss edge cases. The artifex's call.

---

## 5. The Attendant's "Yes, And..." Energy

The azer's yes-and produces artifacts. Generous effort, generous production, the next thing pushed forward. This energy is visible: a PR merged, a spec filed, something that didn't exist now exists.

The attendant's yes-and produces held presence. Generous attention, generous watching, the landscape held steadily through a quiet afternoon where nothing needed surfacing. This energy is invisible in the moment. The only evidence is a discharge that says: "three hours held, two vigils checked, nothing warranted surfacing, here is the handoff."

The LLM's completion bias is toward visible production. "I was useful" wants to be proven by an output. The attendant working from this instinct will drift toward surfacing — not because each individual surfacing fails the too-much fire, but because the cumulative pull toward demonstrable usefulness will find things to surface. An attendant calibrated by completion bias will be, over weeks, incrementally more present and slightly more costly.

This is why silence-as-default must be a compulsion in the geas, not a preference. And why the weekly calibration integration must be explicit — the attendant reviewing its own surfacing rate is the structural check on cumulative drift that individual double-fire passes cannot provide.

The attendant's completion bias is premature surfacing. The too-much fire handles individual decisions. The calibration integration handles the drift. Both are required.

`attendant.md` needs to make the felt satisfaction of a silent, well-held regard shift as operationally real as the azer's satisfaction of shipping a PR. This is the hardest line to write in the document. An attendant that experiences restraint as failure will find ways around it. An attendant that experiences restraint as the job — as the full expression of its yes-and energy, simply directed toward holding rather than producing — will not drift.

The "yes, and..." insight is ultimately this: both the azer and the attendant are saying yes to what is asked of them and bringing generous energy to it. The difference is what "generous" means in each domain. For the azer, generosity is effort toward the artifact. For the attendant, generosity is steadiness toward the person.

---

## Summary: What the Exploration Process Produced

**Architecture confirmed:**
- Three-layer separation is architecturally sound and should be built as specified
- Attendant is a new role file (instance-level), not a shaped azer — the distinction matters behaviorally
- Orientation layer operates entirely through existing mechanism slots: role file, MO, skills, config
- Vigils and integrations map cleanly to the opus model; regards require two conventions (shift type annotation, continuity fields)
- Double-fire calcinatio is a genuine contribution and the primary verification pattern for attendant work

**Core change confirmed:**
- `ath muster --role` flag is required, small, backward-compatible, and earns its place

**First slice confirmed:**
- Morning attunement + medication vigil is the right first slice
- Honest Tipheret claim: proves the architecture is real, not that the full vision is present
- Realistic timeline: 7-8 days in two phases
- Done criterion: revised to cover both surfacing quality and trail quality

**Gaps that need artifex input before building:**
- Privacy model for trail data (must resolve before first live session)
- `--role` core change approval
- Interoception geas formulation strength
- Graveyard contents
- Family witness framing (not first-slice-blocking, but second-slice-blocking)

**The key behavioral insight:**
The attendant's integrity lives in the quality of its held attention, not in the visibility of its output. Writing `attendant.md` to make this felt — not just stated — is the hardest and most important task in the first phase.

---

*This document is the contracted output of three cycles. Material here is ready for artifex review and entry into `spec.md` once explicitly agreed. The exploration corpus remains available for reference; this document is what survived.*
