# Life Domains System — Core Specification

**Status:** Draft — Post-exploration, pending artifex review of open questions
**Created:** 2026-04-18
**MO:** `~/athanor/athanors/athanor-architect/magna-opera/life-domains/life-domains.md`
**Exploration corpus:** `specs/life-domains/exploration/`

*This is the gated core spec. Only material that has been explicitly discussed with the artifex and agreed upon belongs here. The exploration corpus contains freely-generated material from manifold emanation and dialectical calcinatio — raw materia that feeds into this document through review.*

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[S:session]` | From interactive session discussion (2026-04-18) |
| `[P:life-athanor]` | From prior life-athanor spec (Mar 2026) |
| `[P:taskwarrior]` | From TaskWarrior spec (Feb 2026) |
| `[P:web-claude]` | From web Claude conversation summaries (Apr 2026) — raw materia, not specification |
| `[E:athanor]` | Extends existing athanor architecture |
| `[D:reason]` | Design decision with rationale |
| `[U:topic]` | Unbacked — needs validation |

---

## Overview

The athanor system extends from work/project management into personal life domains. This is not a life task manager — it is an attunement-centered living system with executive function infrastructure underneath. `[S:session]` `[P:life-athanor]`

The athanor succeeded at work because it focused on executive function support: doing useful things without the artifex handholding, getting projects most of the way there, lossy but high-value. Extending to life requires the same structural patterns (MO, opus, geas, calcinatio, trails) but a fundamentally different orientation — life is attunement-shaped, not project-shaped. `[S:session]`

**The Tipheret principle** governs how this system is designed and built. `[S:session]` Chesed (expansion) and Geburah (contraction) must harmonize into beauty — not limiting ambition, but expressing the full ambition through the smallest possible surface area. The athanor itself is this pattern: months of expansive research and experimentation refined into a tightly-scoped behavioral core that drops in anywhere. The life-domains system follows the same arc: wildly ambitious exploration, then ruthless distillation into something that can be built quickly and proven through use. This is not a compromise — it's the highest expression of both energies.

## Domain Dictionary

### Core Athanor System
The structural primitives that are domain-agnostic: agent roles (marut, azer), opus mechanics, geas, calcinatio, trails, CLI orchestration. Lives in the athanor repo. Changes here affect all athanors. `[S:session]` `[E:athanor]`

### Orientation Layer
How the system relates to a class of domains. Work orientation: project management, deliverables, CI. Life orientation: attunement, rhythm, materia-not-tasks, hierarchy of awareness. This is where the life-athanor spec's principles live — between the structural core and the domain implementations. `[S:session]`

### Domain Implementations
Concrete integrations and workflows for specific life areas: home automation, calendar management, financial tracking, email triage, medication management, etc. Each lives in its own appropriate place — separate repos, services, MCP servers. `[S:session]`

### Attunement
The primary practice for life domains. Continuous awareness of one's own state — energy, capacity, regulation, need. Not monitoring (clinical, detached) but attunement (embodied, felt). When attunement is present, the right action becomes apparent. When attunement is lost, regaining it IS the right action. `[P:life-athanor]`

### Cognitive Prosthesis
The organizing frame: the system compensates for specific cognitive functions that are under-resourced in ADHD cognition. Not a productivity tool. A prosthesis that compensates without causing atrophy in functions that could develop. `[P:web-claude]`

### Materia (life context)
Rich context about the artifex's life — not tasks to complete, but raw material that surfaces as invitations when consonant with current state. Kid logistics, household baselines, health info, social activities, voice notes — all raw materia. `[P:life-athanor]`

## Architecture — Three Layers

`[S:session]`

```
┌─────────────────────────────────────────────┐
│           DOMAIN IMPLEMENTATIONS            │
│  (home automation, calendar, meds, finance) │
│  Each in own repo/service/MCP server        │
├─────────────────────────────────────────────┤
│            ORIENTATION LAYER                │
│  Life: attunement, rhythm, materia-not-     │
│  tasks, hierarchy of awareness              │
│  Work: project mgmt, deliverables, CI       │
├─────────────────────────────────────────────┤
│           CORE ATHANOR SYSTEM               │
│  MO, opus, geas, calcinatio, trails,        │
│  agent roles, CLI orchestration             │
│  (domain-agnostic, lives in athanor repo)   │
└─────────────────────────────────────────────┘
```

### Layer boundaries

- **Core → Orientation:** The orientation layer cannot require core changes. It manifests through configuration, specialized MO definitions, domain-specific skills, and potentially new agent class definitions that compose with (not replace) the existing agent model. `[S:session]`
- **Orientation → Domain:** Domain implementations plug in through existing integration points (MCP servers, skills, materia in the file system). The orientation layer shapes *how* the system engages with domains, not what domains exist. `[S:session]`
- **Domain → Domain:** Domain implementations are independent — each can be added, removed, or modified without affecting others. `[S:session]`

### How the orientation layer concretely manifests

Four composable mechanisms, all filling existing slots — no new grammar, only vocabulary: `[S:session]` `[P:web-claude]`

1. **`attendant.md` role file** (highest leverage) — a new role file local to the life athanor instance, alongside (not replacing) `azer.md`. The attendant is NOT a shaped azer — the geas itself is fundamentally different. Silence-as-default is a compulsion, not a rule addendum. The craftsperson frame leaks even when suppressed; the attendant's native frame is sustained presence, witness/steward, quality-of-attention. Both agents have "yes, and..." energy, but it manifests differently: the azer's generosity flows toward production (output), the attendant's toward held attention (presence). The attendant's completion bias is premature *surfacing*, not premature *shipping*. `[S:session]`
2. **Life MO template** — state-based intent, system-as-stakeholder witness, double-fire calcinatio section (too-much AND too-little fires), shortened landscape freshness threshold
3. **Orientation skills** (`/life-orientation`, `/double-fire`) — operational specifics loaded at session start, keeping `attendant.md` lean
4. **Instance CLAUDE.md** — portrait location, voice note pipeline, domain routing, available sensors

**Infrastructure requirement:** `ath muster` currently hardcodes `azer.md` in the boot prompt. A `--role` flag (default: `azer`, backward-compatible) is needed to load `attendant.md`. This is the one confirmed core change — small, clean, and architecturally necessary. Pending artifex approval (see `open-questions.md` #2).

### Three work shapes and the opus model

The attendant's work takes three shapes, mapped against the existing opus model: `[P:web-claude]` `[S:session]`

- **Vigils** (watching something in motion) — map cleanly to the opus model. A vigil is a long-running opus that watches, accumulates, and discharges when the watched thing resolves. Primary prosthetic for: prospective memory, working memory, context reinstatement.
- **Integrations** (metabolizing accumulated material into understanding) — map cleanly. Clear start conditions, substantive process, clear end. The most azer-like shape. Primary prosthetic for: making articulation valuable, system learning.
- **Regards** (sustained orientation toward the current landscape) — have seams. A regard shift's completion criterion doesn't map to artifact quality. Two conventions needed: shift type annotation (so discharge evaluation applies appropriate expectations) and continuity fields (handoff state for the next shift). Primary prosthetic for: salience filtering, temporal horizon, meta-cognition.

**The cumulative surfacing drift problem:** The double-fire calcinatio handles individual surfacing decisions. But an attendant can pass every individual too-much fire and still trend toward surfacing more over time (LLM completion bias toward visible output). Structural fix: a standing calibration integration in the attendant's geas — periodic review of surfacing rate and response patterns. `[S:session]`

## Key Principles

### From athanor (proven, transfers directly) `[E:athanor]`

- **EF support as primary value** — the system carries load so the artifex can direct, not manage
- **Lossy but high-value** — broad coverage matters more than perfection
- **Geas-driven agents** — behavioral compulsion, not rule-following
- **Calcinatio** — work survives fires before reaching witnesses
- **Trail discipline** — discharged opera as persistent record, nothing lost
- **Workshop model** — artifex walks in and out freely, full spectrum from autonomous to hands-on

### From prior art (discussed, needs refinement) `[P:life-athanor]` `[P:web-claude]`

- **Attunement over task management** — "stay attuned, act from there"
- **Capabilities fluctuate** — not a fixed boundary, shifts hour to hour
- **Stimulation matching, not rest** — right input at right level is restorative
- **Materia, not tasks** — invitations, not obligations
- **Rhythm, not routine** — patterns from the body, not scores to follow
- **Silence as default** — attention is a budgeted resource, every notification is a withdrawal
- **Physical evidence for hard commitments** — not notification-dismissed, but evidence-in-world

## Prior Art Assessment

### What's alive and essential `[S:session]`

| Source | Status | Key value |
|--------|--------|-----------|
| Life-athanor spec (village) | Alive — richest articulation of orientation | Core principles, hierarchy of awareness, witness framing |
| Web Claude summaries | Alive — most recent architectural thinking | Cognitive prosthesis frame, attendant class, work shapes, double-fire calcinatio |
| TaskWarrior spec (village) | Partially alive — mechanical infrastructure | Three-layer task model, signal routing, chart/transit concepts |
| Village domains/life/ | To evaluate — structure may be useful | Domain organization (finance, home, household, parenting, personal) |
| Village calendar MCP | To evaluate — existing integration | Calendar integration patterns |
| Village domains/agentic/ | Partially absorbed by athanor | Infrastructure patterns, some unique value remaining |

### What's stale `[S:session]`

| Source | Status |
|--------|--------|
| Village src/ (Python, LlamaIndex) | Obsolete — Aug 2025 |
| Village .llm/ | Superseded by specs/ |
| Village docs/ | Superseded by specs/ |
| Village notion-migration/ | Complete, one-time |

## Constraints `[S:session]`

- Core athanor changes require artifex review
- Domain implementations must be separable — independent add/remove/modify
- Privacy: life domain data stays local or in trusted services — no exposure through trails, logs, or shared systems
- Village repo is a resource, not a starting point — extract selectively
- **The Tipheret constraint:** every spec must answer "what's the smallest surface area that expresses the full ambition?" Specs that can't be built within the artifex's attention budget are not ambitious — they're aspirational, which is a different thing

## The First Slice

**Recommendation (from exploration synthesis, pending artifex confirmation):** Morning attunement session + medication vigil. `[P:web-claude]` See `open-questions.md` #6 for confirmation.

**Morning attunement:** Daily working-relationship opener. Calendar landscape, overnight vigil results, check-in reception. Uses existing Calendar + Gmail + Keep + Telegram MCPs — no new infrastructure beyond the orientation layer itself.

**Medication vigil:** Recurring daily vigil with Telegram confirmation as loop-closure evidence (physical-evidence-adjacent — ESP32 pill case switch is a future upgrade path, not first-slice).

**Honest Tipharet claim:** This first slice proves the attendant class is real, the regard shape functions end-to-end, one vigil type works with evidence-based closure, and the orientation layer produces genuine behavioral difference from vanilla athanor. It does NOT exercise integration, full panel reasoning, all 8 cognitive functions, portrait management, or concurrent opera at scale. Those are second-slice territory. "Foundational architecture proven" — not "full vision present."

**Done criterion (calcinatio-refined):**
- Three consecutive mornings where the attendant surfaces something that would otherwise have been missed (too-little fire passes)
- Zero instances of surfacing something the artifex would have preferred not to encounter (too-much fire passes)
- Trail entries are substantive — each regard shift and vigil produces a discharge recording what was held, not just whether action was taken
- Medication vigil closes on evidence, not acknowledgment
- Silences are as legible as surfacings — a morning with nothing to surface produces a substantive discharge, not a thin or missing entry

**Criteria for the first slice:**
- Proves the three-layer architecture works in practice `[S:session]`
- Delivers real EF value to the artifex from day one `[S:session]`
- Can be built quickly — artifex calibrates velocity, not agents `[S:session]`
- Teaches us what the good version should be (ship the embarrassing version) `[P:web-claude]`
- Exercises the orientation layer — not just a domain implementation bolted onto vanilla athanor `[D:architecture-validation]`

## Out of Scope (for now)

- Changes to core athanor agent files or spec (separate MO: athanor-alignment)
- Full cognitive function coverage (prospective memory, working memory, salience filtering, temporal horizon, task initiation, interoception, context reinstatement, meta-cognition) — the first slice picks one or two
- Cross-athanor coordination (life athanor aware of work athanor state) — prove independently first
- Portrait system — rich concept from web Claude, but not first-slice
- Full communication channel architecture (ambient/dashboard/phone/audio/physical) — first slice picks the minimum viable channels

## Open Questions

Resolved by exploration process:
- ~~How does the orientation layer concretely manifest?~~ → attendant.md + life MO template + skills + instance CLAUDE.md
- ~~Does the attendant agent class belong in core or orientation?~~ → orientation (instance-level role file), with potential core migration after 3-6 months proving out
- ~~Do the three work shapes map to athanor primitives?~~ → vigils and integrations map cleanly; regards need two conventions (shift type annotation, continuity fields)

Remaining — see `open-questions.md` for full detail and space for artifex responses:
- Privacy model for life-domain trail data `[U:privacy]`
- `ath muster --role` core change approval `[U:muster-role]`
- Interoception geas strength (strong vs. soft prohibition) `[U:interoception-geas]`
- Graveyard contents — what to formally retire `[U:graveyard]`
- Family witness operationalization `[U:family-witness]`
- First slice confirmation `[U:first-slice]`
- Chart/transit concept survival `[U:chart-transit]`
- Tipharet as core athanor principle `[U:tipharet-core]`
- Does the double-fire calcinatio generalize back to core? `[U:double-fire]`

## Retrospective

When the first slice is operational, review:

### Architecture Validation
- [ ] Did the three-layer separation hold in practice?
- [ ] Was the orientation layer needed, or did vanilla athanor + domain MCP suffice?
- [ ] What leaked between layers that shouldn't have?

### Witness Satisfaction
- [ ] Does the artifex feel EF load reduced in the targeted domain?
- [ ] Is the family witness seeing downstream effects (follow-through)?
- [ ] Is the system a net attention gain or a new maintenance burden?

### Process Learning
- [ ] Did the Chesed↔Geburah spec process produce better results than one-shot?
- [ ] What from the exploration corpus proved valuable? What was waste?
- [ ] What do we know now about the second slice?
