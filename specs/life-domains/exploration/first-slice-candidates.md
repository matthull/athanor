# First Slice Candidates — Life-Layer System

**Status:** Exploration — Chesed phase
**Created:** 2026-04-18
**Purpose:** Evaluate candidate first implementations that prove the three-layer architecture works, deliver real EF value from day one, and can be built in days not weeks. Input to the contraction (Geburah) phase that will select the actual first slice.

---

## The Evaluation Frame

Before examining candidates, the criteria need to be operationalized so the evaluation has teeth.

**Proves the three-layer architecture** means: the core athanor system runs unchanged, there's a recognizable orientation layer shaping how the system relates to life domains (not just vanilla athanor with a different MO topic), and there's at least one domain implementation providing real-world signal. A candidate that just bolts an MCP onto an existing athanor and calls it done hasn't proven architecture — it's bypassed the architecture question.

**Delivers real EF value from day one** means: the first morning or first week, something in the artifex's life is measurably different because this exists. Not "useful in principle" — actually experienced as load reduction or loop-closure that wasn't happening before.

**Buildable in days** means: with existing infrastructure (the MCP servers already available, the athanor CLI already working, the Go codebase healthy), a focused effort produces a running system. "Days" is literal — five to ten days of focused work, not weeks of infrastructure-first building.

**Exercises the orientation layer** means: the attendant class is present (or at minimum a meaningful orientation variant of an azer), the geas reflects the prosthesis framing, the work shapes (vigil/regard/integration) are operationally present in some form — not just referenced in documentation.

**Infrastructure already available for free:**
- Gmail MCP: `mcp__gmail__*` — read, search, send, label, filter
- Google Calendar MCP: `mcp__google-calendar__*` — list events, create, get freebusy, search
- Google Keep MCP: `mcp__keep-mcp-pipx__*` — notes, lists, labels
- Telegram MCP: `mcp__telegram-mcp__*` — notify, andon, topics
- Calendar MCP (village): `~/code/village/src/calendar_mcp/` — existing FastMCP Python server
- Slack MCP (egregore): already in use for work athanor
- Linear MCP: already in use for work

The fact that Calendar, Gmail, Keep, and Telegram are **already wired in** is not a footnote — it means several candidates can be attempted with zero new MCP server work. This is the central infrastructure advantage.

---

## Candidate 1: Morning Attunement Session

### What it does for the artifex

Every morning, the artifex opens a tmux window and has a short structured interaction with an attendant-in-residence. The attendant has already done landscape prep: pulled today's calendar, scanned for anything that needs prep or attention, checked for unread emails with time pressure, held overnight vigils about open concerns from yesterday. The session surfaces: what today looks like, what's been held overnight, what needs a decision or action before the day starts.

The artifex's contribution is voluntary and low-friction: voice a short check-in ("low energy, might need to cancel afternoon") or just read and dismiss. The attendant adapts its regard based on what it hears and what it observes from the landscape.

This isn't a daily briefing tool — it's the opening of a working relationship. The attendant shows up with context already loaded; the artifex shows up and orients. The session ends with the artifex knowing what matters today, what's been handled, and what's been set aside with explicit intention.

### Cognitive functions addressed

- **Prospective memory:** The attendant surfaces time-sensitive things approaching today, cross-referencing calendar and email, injecting cues at the right moment rather than surfacing everything always.
- **Working memory:** Overnight material is consolidated and held; the artifex doesn't have to carry yesterday's open loops through sleep.
- **Salience filtering:** The attendant's regard does the pre-salience work — the morning session presents what actually matters today, not an unfiltered calendar dump.
- **Temporal horizon:** Making "today" concretely present, with time-until events rendered viscerally rather than spatially.
- **Context reinstatement:** Handing back yesterday's unfinished threads with context restored, not just a list.

### Tipheret expression

The full vision: a continuous attendant that maintains portrait, holds vigils across time, surfaces things at the right moment, never intrudes without earning it, becomes gradually more tuned to the artifex's rhythms over months.

The first slice expression: a single morning ritual. One session per day. The attendant runs a morning landscape query (Calendar + Gmail), holds zero or one active vigil from the prior session, surfaces a single consolidated landscape view, accepts a one-sentence check-in. Trail entry produced. Session complete.

The full vision is present in this slice: the attendant class is real, the geas is real, the prosthesis is actually functioning (the artifex would have missed things without it). Everything else — longer vigil chains, deeper portrait, ambient channel — is elaboration of this core act.

### What needs to be built

- **Attendant agent guidance file** (`shared/attendant.md` or life-athanor-specific variant): the attendant's geas, work shapes, and morning session protocol. This is the orientation layer made concrete.
- **Life-layer MO template** for the attention-protection MO. Witnesses, fires, satisfaction criteria adapted to the attunement framing rather than production framing.
- **Morning session skill/prompt** (`specs/life-domains/skills/morning-attunement/SKILL.md`): what the attendant does in a morning session — landscape query sequence, consolidation pattern, check-in reception, vigil management.
- **`ath craft` invocation pattern**: how the artifex opens the morning session. Possibly just `ath craft life morning`.

No new MCP servers needed. Calendar + Gmail + Keep + Telegram already available.

### Estimated build effort

3-4 days:
- Day 1: attendant guidance file, first draft geas, morning skill
- Day 2: MO template, first live session, iteration on what the attendant actually does
- Day 3: trail discipline, vigil handling across sessions, discharge patterns
- Day 4: tuning from use, second and third real morning sessions

### Architecture validation

This forces the orientation layer to exist: the attendant's geas, the regard shape (morning shift), the integration of voice check-in, the landscape-as-materia-not-tasks framing. Vanilla athanor with a morning MO would produce something that felt like a morning standup tool. The orientation layer is what makes it feel like a working relationship.

**Risk:** The morning ritual fails to stick if the interaction design is wrong — too long, too much required, too structured. The embarrassing-version discipline applies hard here: a morning session that takes two minutes and surfaces one useful thing beats a beautiful morning protocol that takes ten.

---

## Candidate 2: Email Triage Vigil

### What it does for the artifex

Email is a canonical ADHD trap: the inbox as an undifferentiated pile of obligation, guilt accumulation, and hijacked attention. Opening email is high-cost because you don't know if anything matters, but not opening it is also high-cost because things do occasionally matter and knowing which is a constant background hum.

The email triage vigil flips this. An attendant watches the inbox and maintains a standing summary: what actually needs attention, what's time-sensitive, what can be deferred, what can be acted on in thirty seconds. The artifex never opens Gmail directly for triage — they check the attendant's summary, which has already done the work.

Concretely: the attendant runs inbox scans at intervals (or on-demand), applies a small number of structural categories (needs-response, time-sensitive, FYI-important, safe-to-ignore), and holds the relevant-now items in a Keep note or Telegram message. The artifex pulls this when they want to engage email, not when Gmail decides they should.

This addresses the pull-first principle directly. Email becomes a resource, not an interruption source.

### Cognitive functions addressed

- **Salience filtering:** Pre-classifying the inbox so what matters is already separated from what doesn't. The attendant's judgment replaces the artifex's scan-everything.
- **Working memory:** The current-state summary replaces the mental inventory the artifex maintains by scanning and re-scanning.
- **Prospective memory:** Time-sensitive items surface at the right moment, not when the artifex happens to open Gmail.
- **Task initiation:** Surfacing with concrete next action ("reply to this, it needs a one-line confirmation") rather than the abstract obligation of "email to deal with."

### Tipheret expression

Full vision: continuous monitoring, rich classification schema, automated actions (labels, filters, drafts), integration with calendar (calendar invite in email automatically assessed against schedule), full correspondence history in portrait.

First slice: one scan per day, four categories (act-now, read-when-able, FYI, handled), summary to a pinned Keep note. No automated actions. The attendant reads the Keep note when the artifex asks "what's in email today?" The first week is likely to produce surprising signal about what the attendant misses and what it over-surfaces — which is exactly the information needed for calibration.

### What needs to be built

- **Email triage skill**: attendant protocol for scanning, classifying, summarizing. This is mostly prompt engineering against the Gmail MCP. The calcinatio fires for this are both content-based (did I miss something?) and attention-cost-based (did I over-surface something?).
- **Keep note structure for inbox summary**: a living note the attendant maintains and the artifex pulls.
- **Integration pattern**: how the triage vigil persists across sessions — the opus lifecycle for a running email vigil.

No new MCP servers. Gmail + Keep already available. Keep's pinned note is a lightweight ambient channel — pull-first, zero push.

### Estimated build effort

2-3 days:
- Day 1: triage skill, classification categories, first scan, Keep note structure
- Day 2: first week of live use, calibration, vigil lifecycle patterns
- Day 3: discharge patterns, calcinatio fires tuned from observation

### Architecture validation

Email triage lives cleanly in the domain implementation layer — it's a Gmail + Keep integration. The orientation layer manifests in the attendant's posture toward the inbox: not "process all email" but "maintain the artifex's relationship with their inbox as a low-cost, high-value channel." The geas shapes the attendant to be protective of the artifex's attention budget, not to be a helpful email bot.

**Risk:** Classification accuracy on day one will be wrong in both directions — things that matter missed, things that don't matter surfaced. The embarrassing version ships anyway and fails informatively. The risk is that early misclassification creates more anxiety than no system would (the "now I have to validate the attendant's work" problem). Mitigated by: starting with a single daily scan, not continuous monitoring, and keeping the first slice bias toward under-surfacing rather than over-surfacing.

---

## Candidate 3: Medication Management with Physical Evidence

### What it does for the artifex

Medication compliance is a prospective memory problem with extremely high stakes. Missing a dose isn't just inconvenient — it affects cognitive function in ways that compound (worse cognition → less likely to remember meds → worse cognition). Conventional reminder systems fail because they produce notification-dismissed, which is not the same as meds-taken.

The physical evidence approach: a hardware switch on the pill case (or the actual pill case position in a predictable physical location) generates a binary signal — meds taken or not. The attendant watches this signal. Loop closes on evidence, not on acknowledgment.

Without hardware (the embarrassing version): the loop closes on an explicit confirmation interaction ("meds taken, 8:47am") sent via Telegram. The attendant holds the vigil, the artifex provides the evidence-signal, the attendant records closure. Hardware comes later if the pattern proves worth automating.

What the attendant adds beyond a simple reminder: holds the vigil without nagging, surfaces the missed-dose risk at the right contextual moment rather than on a rigid schedule, integrates with the temporal horizon concern (missed meds + late afternoon = flag that needs attending, not a second reminder), and records the pattern over time for health awareness.

### Cognitive functions addressed

- **Prospective memory:** The hardest possible version of this problem — daily, critical, easily overridden by morning hyperfocus or distraction. The attendant holds this as its highest-priority vigil.
- **Temporal horizon:** Meds at the wrong time of day have different impacts. The attendant's regard on the temporal landscape knows that 10am unconfirmed is different from 2pm unconfirmed.
- **Meta-cognition:** Pattern-tracking over time ("you've missed three doses in the last two weeks, all on mornings when you had early calls") without surveillance — this is retrospective calibration material, not a running audit.

### Tipheret expression

Full vision: ESP32 on pill case, MQTT to home automation, automatic loop closure, ambient light shifts from "meds reminder" to "meds confirmed" state, no phone interaction required.

First slice: Telegram confirmation message, attendant vigil, daily check in morning session, trail entry. The artifex sends a Telegram message ("meds") and the vigil closes. If not closed by 10am, the morning session surfaces it. That's the whole thing.

The full vision is entirely present: loop closes on physical-adjacent evidence (the act of taking and confirming), the attendant holds the vigil without nagging, the system knows. The ESP32 is just automating the confirmation message — the logic is identical.

### What needs to be built

- **Medication vigil pattern**: attendant protocol for a recurring daily vigil. This is mostly lifecycle design — how a vigil that recurs daily works (does it inscribe a new opus each day? does one vigil persist and get refreshed?).
- **Telegram → vigil closure integration**: the attendant watches for confirmation signals. Telegram MCP already provides this as incoming messages.
- **Morning session integration**: meds vigil state surfaces in morning session if unconfirmed.

Minimal infrastructure. Telegram already available. The vigil lifecycle is a core attendant pattern question that needs to be answered regardless of which first slice is chosen.

### Estimated build effort

2 days:
- Day 1: recurring vigil lifecycle design, medication vigil skill, Telegram confirmation protocol
- Day 2: integration with morning session, trail discipline, first week of live use

### Architecture validation

This is the clearest test of a hard commitment pattern — the thing the prior art (web Claude summaries) calls out as most important to get right. Physical evidence as loop-closure discipline, not notification-acknowledgment. The orientation layer manifests in the attendant's framing of this as a health vigil rather than a reminder system — the geas produces different behavior than a reminder agent's geas would.

**Risk:** This is the most emotionally loaded candidate. Medication management has psychological weight — missed doses carry shame, and a system that surfaces misses can feel like surveillance even if it's genuinely trying to help. The attendant's geas needs to be explicit about the anti-surveillance posture here. First version should err strongly toward understated.

**Hardware path:** If the Telegram confirmation pattern proves out in first weeks, the ESP32 step is mechanical — same logic, automatic confirmation source. The first slice doesn't need hardware; hardware is an upgrade path.

---

## Candidate 4: Calendar-Aware Temporal Landscape

### What it does for the artifex

Calendar-aware temporal landscape is about making time visceral. ADHD time blindness means "meeting in two hours" and "meeting in two weeks" feel roughly the same until the meeting is in twenty minutes. Standard calendar views show events spatially; they don't create temporal presence.

The temporal landscape attendant maintains a living model of the coming period — today, tomorrow, the next meaningful horizon — and surfaces it as context the artifex can check rather than events that happen-to-them. The key difference from a calendar notification: this is a state model, not an event queue. Checking it costs nothing; not checking it is safe.

Concretely: the attendant pulls calendar events for the next 48-72 hours at session start, produces a temporal summary ("light morning, substantial afternoon — 2pm meeting requires prep you haven't done, evening free"), and holds this summary as the regard's ground truth. Time-sensitive prep surfaces in the morning session. Non-time-sensitive things recede.

The prep-requirement modeling is where this gets genuinely valuable: not just "you have a meeting at 2pm" but "this meeting has a prep need that hasn't been addressed — 30 minutes of context you'll want." That's the attendant's judgment, not a calendar feature.

### Cognitive functions addressed

- **Temporal horizon:** Making the coming period viscerally present rather than spatially represented. "The next thing that matters is X in Y hours" is different from a calendar grid.
- **Prospective memory:** Prep requirements surfaced at the right moment, not when the meeting is starting.
- **Salience filtering:** The attendant's temporal summary pre-selects what's actually salient in the calendar, not all events equally.
- **Task initiation:** "This prep needs to start now to be ready" is activation-scaffolding, not a reminder.

### Tipheret expression

Full vision: continuous temporal awareness, integrations with all calendars (work + life), prep-requirement tracking, automatic escalation as deadlines approach, ambient light reflecting urgency tier, morning landscape including full family schedule.

First slice: single calendar query each morning, 48-hour summary, one prep-need surfaced if any. No escalation, no ambient channel. The morning session includes the temporal landscape as standard context. That's everything.

### What needs to be built

- **Temporal summary skill**: attendant protocol for calendar query and temporal summary production. Calendar MCP already provides this.
- **Prep-requirement heuristic**: how the attendant identifies that an event has prep needs (keyword matching, event type, duration, absence of preceding prep blocks). This is the orientation layer's domain judgment — it's not in the calendar, it's in the attendant's regard.
- **Integration with morning session**: temporal landscape becomes a standard component of morning attunement.

No new MCP servers. Calendar MCP already available.

### Estimated build effort

1-2 days standalone, or ~0 additional days if built as part of Candidate 1 (morning attunement subsumes this).

### Architecture validation

This candidate is most naturally a subsystem of Candidate 1 rather than a standalone first slice. The temporal landscape is exactly what the morning attunement session surfaces — it's the calendar-informed component of the attendant's regard. If this is the first slice, it risks being "calendar app but with an agent" rather than something that exercises the orientation layer. The full architecture value shows only when the temporal landscape is held *within* the attendant's regard structure, not just queried and summarized.

**Risk:** Low technical risk. High conceptual risk of being too thin to prove architecture — you can produce a useful temporal summary with a simple script and no orientation layer at all. The orientation layer's presence needs to be explicit in design, not emergent.

---

## Candidate 5: Voice Note Triage and Routing

### What it does for the artifex

Voice notes are a lightweight externalization method — the artifex voices a thought, it goes somewhere, working memory frees. The problem is "goes somewhere" has historically meant "disappears into an app." The village system had routing infrastructure for this; it's worth examining as a candidate because the infrastructure exists and the value is real.

Voice note routing: the artifex records a note on their phone, it arrives in a designated Keep note (or Telegram, or another inbox), the attendant reads it on its next cycle, classifies it (open loop to hold, action required, context update, emotional reflection), and routes it. Open loops become vigils. Actions enter the appropriate domain. Context updates flow to the portrait. Emotional reflections are held without surveillance.

The interoception support angle is strongest here: voice notes are how the artifex practices noticing, and a system that makes articulation feel consequential (it goes somewhere and gets used) supports the practice without replacing it.

### Cognitive functions addressed

- **Working memory:** Articulation frees working memory. The system absorbs and holds; the artifex's cognitive load reduces.
- **Context reinstatement:** Voice notes about in-progress work become context-restoration material when the attendant holds them.
- **Interoception (special case):** Articulation practice — the system receives and uses what the artifex notices, making the noticing feel worthwhile.
- **Meta-cognition:** Pattern material accumulates over time without surveillance.

### Tipheret expression

Full vision: frictionless voice input from any device, semantic classification, automatic vigil creation, portrait integration, cross-domain routing, retrospective pattern insights on-demand.

First slice: Keep note as the designated inbox (the artifex already uses Keep), attendant scans the inbox note each morning session, classifies entries into three buckets (hold-as-vigil, requires-action, pure-context), produces a brief digest. One Keep note in, one attendant reading per day.

### What needs to be built

- **Voice note inbox pattern**: a designated Keep note the artifex drops things into. Already possible with existing Keep MCP.
- **Classification skill**: how the attendant reads and classifies raw voice note content.
- **Routing logic**: what happens to each category (vigil inscription, action note, portrait update).

The routing logic is the hard part — vigil inscription requires the attendant to be running, and "routing to portrait" requires a portrait to route to. This candidate has a dependency on the orientation layer being somewhat mature.

### Estimated build effort

3-4 days, but depends heavily on other infrastructure being in place first. Doesn't stand alone well as a first slice — it needs vigils to route to and a portrait to update.

### Architecture validation

Exercises the integration work shape specifically. Strong orientation layer signal: the classification must operate from the prosthesis framing (what does the artifex need held? what serves them?) rather than from a generic "organize my notes" posture. But the dependency on other infrastructure makes this a second or third slice, not a first.

---

## Candidate 6: Kid Logistics Vigilance

### What it does for the artifex

Kid logistics — school, therapy appointments, activity scheduling, medical needs — are a high-stakes prospective memory domain with high social consequences when missed. The specific ADHD failure mode here is that the logistics *are in the calendar* but they don't feel imminent until they're imminent, and prep requirements (packed bags, forms, ride arrangements) don't surface until too late.

The kid logistics attendant holds vigils on upcoming appointments, surfaces prep requirements with enough lead time, and maintains a lightweight standing awareness of the kid-logistics calendar. This differs from the general temporal landscape (Candidate 4) by being specifically focused on this domain's failure modes and prep requirements.

### Cognitive functions addressed

- **Prospective memory:** The specific failure mode this addresses — appointments known but prep missed until too late.
- **Task initiation:** Prep requirements surfaced with concrete first step, not abstract "prep for appointment."
- **Temporal horizon:** Making upcoming appointments feel temporally present rather than future-space.

### Tipheret expression

Full vision: integrated with school apps, therapy portals, activity platforms; automatic prep requirement inference; family calendar integration; coordination with other caregivers.

First slice: a dedicated Keep list of upcoming kid appointments and their prep requirements, maintained by the attendant pulling from the calendar. One weekly review vigil that surfaces the coming week's kid logistics. The attendant's job: check the calendar for kid-tagged events, identify prep needs from pattern knowledge, hold the list current.

No new MCP servers. Calendar + Keep already available.

### Estimated build effort

2-3 days — this is essentially a domain-specific variant of Candidate 4 with prep-requirement modeling focused on a specific context.

### Architecture validation

This exercises the domain implementation layer specifically — kid logistics is one domain within the life orientation layer's reach. The orientation layer's role: shaping how the attendant relates to this domain (high-stakes, consequences are social/family, prep-requirement modeling requires contextual knowledge about what "prep" means for different appointment types). The architecture question is whether this domain-specific knowledge lives in a skill file, in the MO's context section, or somewhere else.

**Risk:** The most value here comes from prep-requirement modeling, which requires the attendant to know things about what different appointment types require. That knowledge has to come from somewhere — either rich initial context in the MO, or accumulation over time via the portrait. First slice may feel thin if the prep-requirement modeling is weak.

---

## Evaluation Matrix

| Candidate | Proves Architecture | Delivers Day-1 EF Value | Buildable in Days | Exercises Orientation Layer | Uses Existing Infrastructure |
|-----------|--------------------|-----------------------|------------------|----------------------------|------------------------------|
| 1. Morning Attunement | High — orientation layer most visible | High — daily practice from day one | 3-4 days | High — regard shape, geas, vigil | Full (Calendar, Gmail, Keep, Telegram) |
| 2. Email Triage Vigil | Medium — classification is domain implementation, not orientation layer | High — inbox as controlled resource | 2-3 days | Medium — geas shapes posture but classification is mechanical | Full (Gmail, Keep) |
| 3. Medication Vigil | Medium — physical evidence pattern is orientation layer signal | Immediate — first morning | 2 days | Medium — anti-surveillance posture, hard-commitment framing | Full (Telegram) |
| 4. Temporal Landscape | Low standalone — subsystem of Candidate 1 | Medium — useful but thin | 1-2 days | Low standalone | Full (Calendar) |
| 5. Voice Note Routing | Medium — integration shape well-exercised | Medium — depends on other infrastructure | 3-4 days (with dependencies) | Medium | Full (Keep, Telegram) |
| 6. Kid Logistics | Low standalone — domain variant of Candidate 4 | Medium — depends on prep-requirement quality | 2-3 days | Medium | Full (Calendar, Keep) |

---

## Recommendation: Candidate 1 (Morning Attunement) as the first slice, with Candidate 3 (Medication Vigil) embedded

The morning attunement session is the clearest first slice for three reasons.

**First, it exercises all three layers simultaneously.** The core athanor runs an attendant instance (core layer). The attendant's geas, work shapes, and regard framing are the orientation layer made visible — without them, the morning session is just a calendar query. The domain implementations (Calendar, Gmail, Keep) plug in through existing MCPs as the bottom layer. The architecture is alive and all layers are present from day one.

**Second, it's a daily practice, which means rapid iteration.** A first slice that runs daily teaches the system faster than one that runs weekly or on-demand. Miscalibrations surface immediately; improvements compound daily. The embarrassing version of a morning session is genuinely usable — even a thin attendant that shows today's calendar and surfaces one open loop is EF-positive from the first morning.

**Third, it can absorb the other candidates as they mature.** Email triage vigil becomes part of what the morning session surfaces. Temporal landscape is the calendar-informed component of the morning regard. Medication vigil holds through the night and surfaces at morning check-in if unconfirmed. The morning session is the natural integration point — over time it becomes richer as each new vigil type activates, rather than requiring separate interaction patterns.

The medication vigil should be built alongside the morning session rather than separately. It's a two-day addition and it solves the most concretely high-stakes problem in the first slice. Building both together means the very first morning session includes: temporal landscape, email salience, and medication confirmation status. That's a genuinely useful first day.

**The Tipheret expression of this recommendation:** The full ambition of the life-layer system — continuous prosthetic support across all life domains, portrait-informed regard, deep interoceptive partnership — is fully present in a morning attendant that holds overnight vigils, shows what matters today, and surfaces one critical hard commitment. Nothing is compromised. Everything else is elaboration.

### First slice build plan (5-6 days)

**Days 1-2:** Attendant guidance file (geas, morning work shape, regard discipline). Life-layer MO template. First morning session script — what the attendant does step by step. Test with a live morning.

**Days 3-4:** Medication vigil pattern. Recurring vigil lifecycle. Telegram confirmation integration. Integration with morning session — meds status as standard component.

**Day 5:** Temporal landscape component. Calendar query in morning session. 48-hour summary production. First prep-need heuristic (simple: anything > 1 hour duration that wasn't in the prior session's landscape is worth flagging).

**Day 6:** Trail discipline, discharge patterns, iteration from 5-6 live morning sessions. Calcinatio fires for the morning session tuned. Decision: what's embarrassing and acceptable vs. what's embarrassing and needs fixing before the first slice is considered done.

The done criterion: three consecutive mornings where the attendant surfaces something the artifex would have otherwise missed, without surfacing anything the artifex would have preferred not to see. That's real.

---

## The Second Slice (preview)

Based on this analysis, the second slice is most likely email triage (Candidate 2). It's almost-independent infrastructure using the same attendant that the morning session runs, it delivers continuous value between morning sessions, and it proves the domain-specific classification competence of the orientation layer. Voice note routing (Candidate 5) is the third slice — it needs vigils and portrait to route into, which exist after the first two slices are operational.

Kid logistics (Candidate 6) and the full temporal landscape (Candidate 4) likely become ongoing enrichment of the morning session rather than separate slices. They deepen the morning regard rather than introducing new interaction patterns.
