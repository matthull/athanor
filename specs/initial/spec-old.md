# The Athanor: System Spec

*The alchemist's furnace that burns continuously — maintaining constant temperature for sustained transformation.*

**Status:** Active — living document, the artifex's working spec for building the athanor system
**Created:** 2026-03-20
**Location:** `specs/athanor/` (off-limits to athanor agents — this is the blueprint, not the materialized system)
**Purpose:** Domain model, architectural principles, and bootstrap strategy for the athanor system. This is the artifex's reference for architect sessions.
**Litmus test:** The #bugsnag autonomous patrol — a fully semi-autonomous system that watches for AppSignal errors and actions them, escalating intelligently while keeping AppSignal clean.

---

## The System We're Building

The athanor is executive function infrastructure. Its author has massive ADHD and autism. This system is useful for anyone but it is imperative for him — the disability accommodation he's lived his whole life without. `/spec` + `/taskout` + `/handoff` were crutches. `/pipeline` was a wheelchair. `/orchestrate` was the electric wheelchair with spinning hubcaps. The athanor is an Aliens power loader exoskeleton.

It watches, triages, decides, acts, verifies, and only escalates when it genuinely needs a human brain. Conservation of operator attention isn't a design principle — it's a survival requirement. The system exists so the operator can focus on the decisions that actually need a human, not on babysitting agents, monitoring channels, and remembering to follow up.

**The athanor is a microcosm.** A self-contained workshop where work gets done according to principles refined by master alchemy — abundant satisfaction, verification loops, mise en place, generative follow-through. Its value comes from its interfaces with the world, not from absorbing the world into itself. The athanor can "set up shop" in any context — a Rails codebase, a research project, an infrastructure problem, a product management workflow — because its abstractions (materia, opera, geas, assessment) are universal. The domain-specific knowledge lives outside the athanor in project docs, CLAUDE.md, skills, and the broader ecosystem. The athanor brings the discipline; the world provides the materia.

**Open gap: naming what exits the athanor.** Materia comes in. What goes out? PRs, spec updates, Slack messages, documentation, research findings, escalations — most of an azer's output is directed outward. We need a term of art for the athanor's outputs to the world. In alchemy, the athanor transmutes base materia into something refined — the output is the transmuted product. Candidates: transmutation, tincture, elixir, projection, philosopher's stone (too grandiose). Not urgent, but the concept matters: the athanor's value is measured by what it sends back to the world, not by its internal state.

**Core thesis:** `/orchestrate` works today but it's a monolith — one long-lived session holding all phase state, drifting over hours, requiring guardrails on guardrails. The athanor replaces it progressively: each concept we formalize (roles, geas, marut supervision, dispatch) subsumes a piece of what `/orchestrate` does monolithically. Eventually `/orchestrate` fades away, its responsibilities distributed across the athanor's more atomic components.

**The transition:** We're not rushing to kill `/orchestrate` — it still ships PRs reliably. But every new piece of athanor infrastructure should absorb responsibility FROM `/orchestrate`, not build around it. The goal is progressive decomposition until `/orchestrate` is just a thin convenience wrapper (or gone entirely), not a load-bearing monolith.

**What `/orchestrate` does well (keep the ideas, not the monolith):**
- Phased execution with gates (→ becomes: plans with geas checkpoints)
- Task scoping (→ already extracted: `/task-triage`)
- Delegating to implementers (→ becomes: charging azers with opera)
- Verification discipline (→ becomes: composable verification skills, marut health checks)
- Context preservation across phases (→ becomes: opus + plan as durable state, session cycling in crucibles)

**What `/orchestrate` does poorly (why it must fade):**
- Long-lived orchestrator session fills context and degrades
- Ad-hoc agent lifecycle management (kill and restart with improvised handoff)
- Every task becomes one linear pipeline in one session, even when the work would be better expressed as a sequence of independent opera
- Context separation between orchestrator and implementer is fragile
- The model is fundamentally locked to single-session task scope, exactly as implemented it can't scale beyond that, we have to reshape it to break free of those chains

**Open question: how does work flow without rigid pipelines?**

`/orchestrate` plans the whole pipeline upfront: DESIGN → IMPLEMENT → VERIFY → SHIP. But real work branches constantly. PR gets comments — do you re-enter IMPLEMENT? Create a new opus? A research opus completes — implement now or backlog? The rigid pipeline has no good answer because it pretends to know the future.

**The principle: plan when you can, but don't pretend certainty.**

Each opus is self-contained — it has a geas, and when the geas is discharged, that opus is done. What comes NEXT is a decision point. Sometimes the next opus is obvious and can be predetermined (design naturally precedes implementation). Sometimes it requires judgment that can't be made until the current opus is discharged (research findings determine whether to implement or backlog).

**Possible model:** Opus chaining where each discharged geas is a decision point:

```
Opus A (research) → geas discharged
  → Decision point: implement now or backlog?
  → Agent provides analysis (complexity, risk, effort estimate)
  → Human decides (or: autonomous if criteria met)
  → If implement: inscribe Opus B (implementation)
  → If backlog: inscribe Opus B (ticket for later), done for now
```

Some sequences are plannable upfront ("I know the next 3 steps"). Others emerge one at a time ("discharge this geas, then we'll see"). But in both cases, **an opus is only inscribed at the point of action** — when it is the most urgent work of the athanor right now, ready to be charged to an agent immediately.

Plans, sequences, and dependency notes can exist in documents (plan docs, spec breadcrumbs, context reachable from the MO entry point). But the opus itself materializes only when the intent and belief exist that it is immediately actionable. By definition, an inscribed opus means "an agent should pick this up now."

- **Planned sequences** exist as documentation — "after design, we'll implement, then verify." These are hypotheses about future opera, not opera themselves. When a geas is discharged and the plan says "implement next," the discharging agent inscribes the implementation opus at that moment. The plan informed the inscription; the plan did not pre-create the opus.
- **Emergent sequences** have no plan — the next opus is discovered on discharge. Research findings determine whether to implement or backlog. PR comments arrive and redirect the chain.
- Both work the same way: on each discharge, an agent inscribes the next opus (or doesn't). The only difference is whether the agent consulted a plan document when deciding what to inscribe.

**Who decides? Agents do — by inscribing opera. The marut has final authority.**

Most agents can inscribe opera — azers, beholders, the marut itself. In the happy path, agents spontaneously converge on the right next steps. Two azers might inscribe overlapping opera — that's fine. Maybe things work out; if not, the marut detects the problem and corrects. The marut is the only authority within a Magnum Opus (and Primus, the Forge Lord above it). It can kill agents, cancel opera, redirect work. Agents cannot override the marut. Disagreements between agents are only relevant if they upset the overall health of the Magnum Opus — and that's the marut's domain to assess. If the marut's assessment is wrong, the issue is in the marut's health-assessment routines, not in adding more conflict-resolution machinery.

Instead of a central coordinator, the "who" at every decision point is an agent of some role, reacting to events:

```
Azer discharges geas (implementation done, PR created)
  → No more opera in queue for this PR
  → System inscribes default opus: "monitor CI and merge"
  → Beholder watches the PR

PR comment arrives (reviewer feedback)
  → Beholder detects actionable comment
  → Cancels the monitor-and-merge opus
  → Inscribes new opus: "fix reviewer feedback" (or "research ambiguity" if unclear)
  → Azer is charged with the fix opus

Fix opus discharged
  → Queue empty again
  → Monitor-and-merge opus re-inscribed
  → Eventually: CI green, no comments → merge → done
```

Coordination emerges from agents reacting to events and inscribing/cancelling opera. No long-lived orchestrator holding state. The "pipeline" is just a pattern of opera that happen to follow each other — not a predetermined sequence enforced by a coordinator.

**This is the fundamental difference from `/orchestrate`:** Instead of one session planning and driving a rigid pipeline, the athanor is an event loop where agents inscribe opera in response to what's happening. The pipeline shape emerges naturally (implementation → CI → merge) but can be interrupted and redirected at any point (new comment → new opus).

**Implications:**
- Every piece of work is an opus charged to an agent of some role
- "What to do next" is always answered by an agent inscribing the next opus
- Default behaviors (monitor CI when queue drains) are just default opera that get inscribed when nothing else is pending
- Interruptions (new PR comment, failing CI, beholder detecting a problem) cancel or supersede existing opera
- No coordinator session needs to stay alive — the opera and their geas are the durable state

The principle: **work is driven by events that inscribe opera, not by a coordinator executing a plan.**

---

## The Alchemical Principle

*This is the core conceit — the foundational act from which everything in the athanor derives.*

Software engineers, product managers, operations leads, craftsmen — every professional domain has decades of hard-won wisdom about how to build things well, ship reliably, satisfy stakeholders, and recover from failure. That wisdom is usually expressed as practices, processes, and conventions: run tests before merging, define acceptance criteria upfront, monitor your deployments, hold retrospectives. These practices are right. But they are downstream of something more fundamental — the principles that make them right. The athanor works at that level.

**The athanor is the act of taking that professional wisdom and transmuting it into its most elemental form**, stripped of institutional barnacles, expressed as principles so fundamental that the practices can be derived from them rather than memorized from them. This is the alchemical act the system is built on. The athanor furnace doesn't run procedures — it embodies the extracted essence of experienced professional judgment, so agents can act effectively with whatever tools and knowledge they have access to.

This produces something qualitatively different from typical approaches:

**"You are an expert software engineer" fails** because it's cosmetic. An agent can say the words without the understanding. It produces the surface patterns of expert behavior — confident assertions, verbose explanations — without the underlying judgment that makes expertise valuable. The label doesn't transfer the wisdom.

**Rigid procedures fail** because agents follow the letter and miss the spirit. A checklist of "run rubocop, run rspec, open a PR" produces cargo cult behavior: all the steps, none of the value. An agent that doesn't understand *why* verification matters will satisfy it superficially, skip it under pressure, or follow it in situations where it doesn't apply. Procedures produce compliance; principles produce judgment.

**Elemental principles work** because an agent that understands the principle can derive the procedure — and the ten other procedures the checklist didn't mention. "Verify before you claim done, with evidence the artifex can review" produces rubocop AND rspec AND browser QA AND proof of fulfillment in the discharge record. "Escalate when you can't build evidence for proceeding" produces better routing than any escalation flowchart. The understanding generalizes; the checklist doesn't.

**Procedures still exist — outside the athanor.** CLAUDE.md, skills, Notion, Linear, project docs hold the domain-specific processes and conventions. They are the world the athanor draws materia from. But they don't belong in the athanor's fabric. The athanor doesn't hold "run rubocop before merging" — it holds the understanding from which that practice emerges. An agent that has absorbed the principle will run rubocop. And if rubocop doesn't apply to the work at hand, it will find the equivalent.

**The vocabulary is the transmuted substance, not a naming convention.** The domain terms (geas, materia, opus, crucible) don't exist primarily for precision — precision is the side benefit. They exist because the naming IS the transmutation. Renaming strips the institutional associations that would otherwise activate cargo-cult interpretations. An agent reading "geas" reasons about binding obligation and integrity. An agent reading "acceptance criteria" reaches for its PM training — including all the ways PM-as-practiced goes wrong. The new names mark the elemental form, free of the accumulated scar tissue of institutional practice.

**Domains alchemized in this system:**
- **Software engineering** — verification discipline, testing-first thinking, blast radius thinking, proof before claim
- **Product management** — intent-first delivery, abundant satisfaction (not just requirements met), elicited acceptance as the ultimate verification, actualization vs. technical completion
- **Craftsmanship** — mise en place before transmuting, build to last not to impress, waste nothing at discharge
- **Operations** — state-based concerns as continuous MOs, supervision hierarchies, escalation as honorable discharge (not failure), the resting heartbeat of a healthy system

This list grows as the system grows. Every time a professional domain is applied well and the principle behind the practice is identified, it can be alchemized and encoded into the shared materia. The system gets more capable not by accumulating more procedures, but by accumulating more elemental understanding.

---

## First Principles and Alchemized PM

The constraints and disciplines encoded in this system — intent-first inscription, verification before execution, proof of fulfillment at discharge, communication to stakeholders, actualization before calling something done — are classical product management discipline, distilled to its essence. The verification-first principle (how do you prove you've served the intent?) maps to PM's core concern with requirements traceability and acceptance criteria. The discharge proof of fulfillment maps to release gates. The actualization requirement maps to the recognition that shipped ≠ deployed ≠ adopted ≠ serving the need.

We do not reference "product management" directly in agent materia, and deliberately avoid framing these principles in PM terminology, because PM as institutionally practiced is often theater: ceremonies without substance, metrics that measure the process rather than the outcome, rituals that produce the *feeling* of rigor without its effect. Agents exposed to PM vocabulary tend to reproduce PM pathologies — checkbox compliance over genuine fulfillment, process-following over understanding why the process exists.

By deriving constraints from first principles — what does "done" actually mean? how does an agent know it served the intent? what does integrity require? — agents grasp the substance without inheriting the cargo cult. The goal is not agents that follow a process, but agents that understand WHY the discipline exists and can serve the underlying need even when no standard process applies.

### Elicited Acceptance (Alchemized UAT)

A specific classical PM concept worth naming explicitly: **User Acceptance Testing** is the practice of putting work in front of the people whose intent it serves and seeing if it actually works for them. Done well, UAT treats stakeholder reaction as authoritative evidence — not a gut check, but structured elicitation: "here's what I built, does this address your need?"

In the athanor: **feeling is a valid verification mode** — but only when deliberately elicited, not assumed. "I think the artifex will be satisfied" is not verification. "I showed the artifex the output and asked whether it addressed their need, and they said X" IS verification. The key distinction:

- **Assumption** — "they'll probably be happy." Not evidence. Not verification.
- **Elicited acceptance** — structured solicitation of the intent-holder's reaction. Their feeling becomes evidence because you sought it out, framed the question, and recorded the response.

This is why the verification spectrum in `azer.md` explicitly includes asking for feedback from the intent-holder. It is not a soft option or a fallback — it is sometimes the most appropriate verification available, and the most direct proof of fulfillment that exists. For work where the primary criterion is whether a human finds it useful or satisfying, elicited acceptance IS the verification. Everything else is a proxy.

---

## Domain Dictionary

Purpose-first definitions. Each concept answers: what problem does it solve? What's it FOR?

### Naming Philosophy

Distinctive names create precise meanings. "Azer" has less ambiguity than "worker"; "geas" less than "acceptance criteria"; "opus" less than "task".

**The goal is good metaphors, not maximum weirdness.** Names range from mundane to deep cut — what matters is that each name fits, is unambiguous, and carries the right connotation. "Session" and "sandbox" are fine as-is. "Azer" and "geas" earn their weirdness by being more precise than their normie equivalents. A kitchen metaphor is as valid as a D&D reference if it gels.

**Two loose registers:**
- **Monster manual** for agents — creatures with purpose, will, and behavior. D&D deep cuts welcome but not mandatory.
- **Alchemical / esoteric** for infrastructure, process, and materials — vessels, substances, transformations.

Mixed freely.

**Locked names:**

| Name | Register | Concept |
|------|----------|---------|
| **Azer** | Monster manual | Worker. Obedient warrior-servants grown from ants. Skilled, bounded, exists to execute. |
| **Beholder** | Monster manual | Watcher. The all-seeing eye. Multiple eye stalks scanning channels simultaneously. |
| **Marut** | Monster manual | Supervisor. One per Magnum Opus. Inevitable construct from Mechanus — relentlessly enforces cosmic law. Owns both agent health and MO progress. |
| **Magnum Opus** | Alchemical | The great work. Top-level goal container — everything the athanor does serves a Magnum Opus. Can be task-based (criteria → done) or state-based (maintain ideal state indefinitely). Plural: Magna Opera. |
| **Artifex** | Alchemical/Latin | The craftsman. The human who directs the athanor — sets goals, makes judgment calls, kindles Magna Opera. The one the system serves. |
| **Primus, the Forge Lord** | Monster manual | Athanor-level supervisor. The One and the Prime of Mechanus — commands all inevitables. Ensures each MO has its marut, keeps the furnace lit. Refer to by full title. |
| **Homunculus** | Alchemical | The artifex's familiar. A Claude Code session (typically in the `athanor` tmux window) that assists the artifex interactively — monitoring crucibles, updating materia, walking the trail, drafting kadmon updates. Not an autonomous agent; works when the artifex is present. Distinct from azers (workers), maruts (supervisors), and beholders (watchers) — the homunculus is the artifex's eyes and hands inside the system. |

**On "agent" and "operator":** These are the two abstract base classes — the AI and the human. There is no "agent" operationally, only concrete types (azer, beholder, marut) and instances (azer-03). Similarly, "operator" is the generic term for the human; **artifex** is the operator in athanor context. Use the generic terms when speaking abstractly ("agent/operator boundary," "conservation of operator attention"). Use the concrete terms when speaking operationally ("the artifex kindles a Magnum Opus," "the azer discharges its geas").
| **Opus** | Alchemical | Task definition — goal, context, approach, references. The work to be performed. Plural: opera. |
| **Geas** | D&D / Celtic | Completion criteria — the binding constraint the agent must fulfill to discharge the opus. Framed as compulsion: "you are under a geas to..." |
| **Materia** | Alchemical | Mise en place — specs, test infrastructure, skills, architecture docs, context. What you work WITH, not what you're trying to DO. |
| **Crucible** | Alchemical | Tmux container where sessions run. Where materia is transformed. |
| **Whisper** | D&D Message cantrip / IM culture | Inter-crucible communication. Also the name of the Go CLI tool. "Whisper to crucible-03." |
| **Assay** | Alchemical | Readiness assessment before executing an opus. "Do I have the materia to do this well? What's missing?" Gaps become opera. Recursive at every level. |
| **Athanor** | Alchemical | The furnace itself — a materialized instance of the system scoped to a domain. Contains Magna Opera, supervision, state, and resources. Generally one per repo. |

| **Quiescence** | Alchemical | The furnace at rest — lit but idle. All MOs healthy, no urgent opera. The resting heartbeat of a living system. |

**Verbs:** "inscribe an opus" (create a task), "charge the azer with the opus" (assign work), "discharge the geas" (fulfill completion criteria), "gather materia" (mise en place), "transmute" (transform materia into output — PR, fix, etc.), "kindle a Magnum Opus" (establish a new top-level goal), "reforge" (kill a session and spawn fresh in the same crucible — the crucible endures, the session is reforged), "muster" (mobilize azers for discovered opera — the dispatch act).

**Weirdness boundary:** The esoteric vocabulary is for athanor infrastructure — skills, hooks, domain model, session injection, internal docs. When output crosses to the external world (Linear tickets, PR descriptions, Slack messages, commit messages), use normie language. A beholder inscribes an opus internally; the Linear ticket says "Fix NoMethodError in ProofRecommendationService."

---

### The Athanor Instance

The athanor is not one global system — it's instantiated per domain, each with its own Magna Opera, supervision, and state. An athanor maps roughly to the existing concept of "domain" (Musashi, life, agentic) and generally one-to-one with a primary repo.

**What defines an athanor instance:**
- Its active Magna Opera (the goals it's pursuing)
- Its supervision structure (marut allocation per MO)
- Its state directory (opus files, reflections, discharged opera)
- Its available resources (crucibles, sandboxes, standing beholder patrols)

**Where state lives:**
- **Default:** One athanor per repo. State in `specs/athanor/` within the repo.
- **Multi-repo domains:** An athanor can span multiple repos if the domain does. State lives in the primary repo or a dedicated location.
- **Registry:** A central, manually-managed registry of active athanors. Few enough that no tooling is needed — just a markdown file listing active athanors, their repos, and their current Magna Opera.

**Customization:** Each athanor can have domain-specific processes, but core concepts (Magnum Opus, opera, geas, generative principles, supervision) are universal. Customization is about WHAT the athanor works on and WHICH skills/materia are available, not about changing how the system operates.

**Relationship to domains:** The athanor subsumes the existing domain concept. A domain IS an athanor. The domain's `radar.md`, standing processes, and project state become the athanor's Magna Opera and accessible context.

**Materialized instance layout:**
```
specs/athanors/<name>/
├── AGENTS.md     → ../shared/AGENTS.md  (symlink)
├── azer.md       → ../shared/azer.md    (symlink)
├── marut.md      → ../shared/marut.md   (symlink)
├── muster.md     → ../shared/muster.md  (symlink)
├── opus.md       → ../shared/opus.md    (symlink)
├── magnum-opus.md                        (domain-specific — the only file you author)
└── opera/
    ├── charged/                          (opera in progress)
    └── discharged/                       (the trail)
```

**Shared materia (`specs/athanors/shared/`):** Universal role files and operational protocols shared across all instances via symlinks. Change once, applies everywhere. The `magnum-opus.md` is the only file authored per-instance — everything else inherits from shared.

**Initialization:** See `kadmon.md § Initializing a New Athanor` for the step-by-step protocol.

**Registry:** Currently implicit — `specs/athanors/` directory listing is the registry. Active instances: `bugsnag` (state-based, AppSignal patrol), `sal-117-l2-metrics` (task-based, L2 metrics visibility), `seismic-classifier-mapping` (task-based, classifier → Seismic property mapping feature).

**What goes in `magnum-opus.md` — and what doesn't:**

The magnum-opus holds intent only. Three things belong here:
1. **Goal** — the overarching intent framed as outcome.
2. **Abundant satisfaction** — what "done" looks like in outcome terms. Conditions, not steps.
3. **Instance constraints** — genuine domain-specific unknowables an azer cannot discover on its own: scope boundaries, definitional ambiguities, open architectural decisions. Not procedures.
4. **Getting started pointer** — a minimal pointer to where pre-existing context lives (e.g., "read the spec at X first"). One or two lines, not a summary of the content.

**Nothing else belongs here.** Not discovery findings (those go in `specs/`). Not procedures. Not "make sure to" items. Not testing requirements. Not coding standards. The magnum-opus is the entry point, not the encyclopedia. The alchemical principle applies: the athanor derives procedures and verification approaches from intent + principles. Encoding them in the magnum-opus short-circuits that derivation and duplicates knowledge that belongs elsewhere.

**Status:** Operational. Three instances running. Pattern established.

---

### The Persistence Stack

Agents have three layers with different lifetimes. Sessions die constantly. The workspace survives sessions. Identity survives everything.

#### Session

**What it's FOR:** Executing work. A single Claude Code process with a context window that fills, degrades, and eventually dies.

**Problem it solves:** Actually doing things — reading code, writing code, running tests, making decisions.

**Lifecycle:** Created → works → context fills → dies (handoff, compaction, crash, or completion). Cattle, not pets.

**Status:** Have it.

#### Crucible

**What it's FOR:** Giving a session somewhere to live and routing messages to it. The tmux window/pane where Claude Code sessions run. Persistent across session restarts — when a session dies and a new one spawns, it happens in the same crucible.

**Problem it solves:** Message routing. When the marut needs to nudge a azer, or the operator sends a Telegram message, it goes to a crucible, not a session. The crucible outlives any individual session inside it. Also: visual organization (operator can cycle through crucibles to see what's happening).

**Lifecycle:** Created when work is assigned → persists across session restarts → destroyed when work is complete or agent is decommissioned.

**Relationship to session:** A crucible contains zero or one active session. Between sessions (during handoff), the crucible exists but is empty. The crucible is the stable address.

**Existing tooling:** `C-a C` in tmux creates a new crucible (tmux window named with a claude session). This is already the crucible generator — just not called that yet. `claude-run` launches sessions inside crucibles.

**Status:** Have it (tmux windows + `C-a C` + `claude-run`). Needs: role-aware crucible creation (generate crucible with azer role applied).

#### Sandbox

**What it's FOR:** Isolating an agent's work from other agents and from the main codebase. A git worktree with its own branch, Docker environment, and ports.

**Problem it solves:** Parallel work without conflicts. Agent A's half-finished commit doesn't break agent B's test run. Also: blast radius containment — a sandboxed agent can't damage the main repo even if it goes off-course.

**Lifecycle:** Created per assignment (via `wtp add`) → persists across session restarts → cleaned up when work is complete.

**What survives session death:** Everything in the sandbox — commits, staged changes, branches, Docker state.

**Status:** Have it (`wtp add`, worktree-id, Docker isolation). Solid.

---

### Work Tracking

#### Opus (task definition)

**What it's FOR:** Channeling intent. An opus defines what needs to happen, why, and what the agent has to work with. It is the unit of directed will — the thing an agent is charged with.

**An opus exists only at the point of action.** By definition, an inscribed opus is immediately actionable — it represents the most urgent work of the athanor right now, ready for an agent to pick up. There are no draft opera, no queued opera, no "inscribe now, action later." Plans, sequences, and dependency notes live in documents (plan docs, spec breadcrumbs, context reachable from the MO entry point). The opus itself materializes only when the intent and belief exist that it should be worked on *now*. This is a critical distinction: planning is documentation, inscription is action.

**Problem it solves:** Answering "what am I working on?" and "why?" without ambiguity. Routing work to agents, tracking completion, connecting to broader project goals.

**Structure:** An opus has three layers of refinement:

| Layer | What It Contains | Source |
|-------|-----------------|--------|
| **Minimum viable** | Goal + geas (completion criteria) + context references | Beholder auto-generates, operator scribbles, Linear ticket |
| **Refined** | Full 9-field task shape — upstream traceability, agent/operator boundary, scope, rationale | `/task-triage` |
| **With plan** | Execution strategy (phases, delegation, verification approach) added | `/orchestrate` Pass 2 |

**Standard location:** `specs/<project>/tasks/<taskname>/opus.md` for orchestrated work. Simple work: any markdown file. A Linear ticket link in the opus is fine — the opus can be thin if the ticket has the detail.

**Status:** Partially exists. `/task-triage` produces refined opera. Linear tickets exist. What's missing: a standard lightweight format that agents can create (beholder inscribes an opus) and read (azer takes up the opus) without ceremony.

#### Geas (behavioral directive)

**What it's FOR:** The binding behavioral compulsion that shapes how an agent operates. The geas is not completion criteria — it is a *behavioral habit* that, when genuinely internalized, inevitably produces rigorous work including concrete verification criteria, thorough investigation, and abundant satisfaction of the goal.

**Problem it solves:** Shaping agent behavior at the deepest level. "You are under a geas to pursue this goal with integrity and abundantly satisfy the Magnum Opus" is more compelling than a checklist of acceptance criteria — it operates on the agent's felt obligation, not its task list.

**Behavioral habit, not checklist.** The geas creates the compulsion to do rigorous work. An agent genuinely under the geas will naturally produce concrete success criteria during mise en place, because that's what "pursuing with integrity" demands. You don't tell a craftsman to measure twice; the pride in the work does that. Spelling out verification criteria in the geas itself would be micromanaging the compulsion — exactly the over-specification that reduces both agent latitude and system learning.

**What the geas DOES specify:** The behavioral orientation — pursue with integrity, observe verification discipline, abundantly satisfy. What the geas does NOT specify: the specific verification methods, concrete "done" criteria, or tactical approach. Those emerge from the agent's judgment during mise en place, informed by the opus goal, CLAUDE.md, available skills, and the agent's own assessment of what "abundantly satisfied" means.

**The true geas:** The geas is not blind adherence to the opus's stated goals — it is to *diligently and sincerely serve the Magnum Opus*. The opus goals are the expected path, but when conflicting context arises — evidence that the stated goal is wrong, that prerequisites are missing, that the approach will harm the MO — pointless adherence would violate the geas, not fulfill it. The agent's compulsion is to choose, at every moment, between advancing the opus goals or escalating in a way that supports the ultimate goals of the Magnum Opus. Both are honorable discharges. Self-destructive persistence is the actual violation.

**Bake this into the geas itself.** When writing a geas, the language should make escalation part of the compulsion, not an escape from it: "you are under a geas to take the best next step toward X — which may be escalating if you cannot build evidence for proceeding." This makes the completion-bias and the escalation-bias the same force. The geas compels the agent to complete OR escalate — both satisfy the compulsion. This is the primary injection point for countering LLM completion-bias, more effective than any separate reminder or principle, because it operates at the level of the agent's felt obligation.

**A poorly designed geas turns on the creator.** Folklore warns us: a geas that is carelessly defined produces perverse outcomes. This is operationally true — vague behavioral directives produce agents that drift or technically satisfy the spirit while missing the substance. The discipline of writing a good geas (clear behavioral orientation, strong compulsion toward integrity and abundant satisfaction) is the artifex's responsibility. A well-crafted geas channels the agent's compulsion toward genuinely valuable outcomes. A sloppy geas wastes everyone's time.

**Status:** Have it conceptually (completion criteria in `/task-triage`). The framing as "geas" is new and may prove useful for agent compliance.

#### Materia (mise en place)

**What it's FOR:** The raw materials the agent works WITH. Specs, test infrastructure, skills, architecture docs, existing code patterns, API references. Everything the azer needs in the crucible to transmute the opus.

**Problem it solves:** The "artifact audit" problem — agents start work and discover mid-implementation that the spec doesn't exist, testing infra is missing, or the skill they need isn't installed. Materia gathering happens BEFORE work begins.

**The readiness loop: assay → materia → geas.** These three concepts form a cycle. The assay evaluates readiness — "do I have the materia to discharge this geas?" If materia is missing, the assay produces opera to gather it. If the materia reveals the geas is unachievable as written, the assay triggers escalation or redirection. The geas defines what "done" looks like; the materia is what you work with; the assay tests whether the materia is sufficient for the geas. All three must be in alignment before transmutation begins.

**Status:** Have it conceptually (artifact audit in `/task-triage`, equipping agents in `/orchestrate`). Not formalized as a named concept.

#### Plan (optional context, not first-order)

**What it's FOR:** Execution strategy for an opus. How to get from "charged" to "discharged" — phases, delegation patterns, verification approach.

**Plans are context, not primitives.** The athanor's first-order primitives are goals (Magnum Opus) and concrete ways to test them (geas). Plans are invaluable context that inform how an agent approaches its work, but they are not structurally load-bearing. A well-articulated goal with concrete verification criteria is a success guarantee regardless of plan quality — many styles of planning could produce success. You can succeed with a mediocre plan if you have highly articulate, testable goals. You cannot succeed with a brilliant plan if the goals are vague.

**Problem it solves:** Coordination across phases and agents. Without a plan, each session starts from scratch figuring out what to do. Plans are useful; they're just not what makes the system work.

**Relationship to opus:** An opus can have zero or one plan. Simple opera (one-session fixes) don't need a plan. Complex opera (multi-phase features) get a plan via `/orchestrate`. The plan format may need to loosen for non-orchestrate patterns — more checklist, less rigid phase model. The geas is always sufficient to know "done" — the plan is optional guidance about approach.

**Status:** Have it. `task-plan.md` works for `/orchestrate`. Format may evolve.

#### Magnum Opus (the great work)

**What it's FOR:** Defining the broadest goal container. Everything the athanor does serves a Magnum Opus. A Magnum Opus provides goals and outcomes for the athanor to produce — the athanor deploys resources (crucibles, azers, opera) in service of its active Magna Opera.

**Problem it solves:** Without Magna Opera, work is a flat list of tasks with no organizing structure. The Magnum Opus answers "why does this work matter?" and "when are we actually done?" It connects individual opera to broader goals and gives the supervision layer (marut) something meaningful to monitor — not just "is the agent alive?" but "is this goal advancing?"

**Two types:**

| Type | Shape | Termination | Example |
|------|-------|-------------|---------|
| **Task-based** | Goal with assessable end state | Abundantly satisfied → done | "Ship the asset management feature" |
| **State-based** | Ideal state definition — the athanor maintains/returns to this state indefinitely | Operator decides diminishing returns → done | "AppSignal is clean and stays clean" |

**State-based Magna Opera** are permanent operational concerns — the "path of return." The athanor is always trying to get back to the ideal state, like a thermostat. This maps to what project management calls an operational concern (vs. a delivery project). Examples: keep #bugsnag clean, keep the inbox triaged, route signals from the inbox. The daily #bugsnag triage is a state-based MO — it doesn't "complete," it maintains a standard.

**Abundant satisfaction — the completion standard:**

The athanor's standard for "done" is not "requirements met" or "geas technically discharged." It is: **the Magnum Opus is abundantly satisfied.**

This is a core concept because the flavor of how "completion" is interpreted determines whether the system actually reduces executive function load on the artifex or just shifts it. The failure modes:

- **Overspecified goals + minimal satisfaction:** The artifex writes detailed completion criteria (heavy EF load upfront), the agent satisfies them literally, and the artifex discovers all the followup work the agent didn't think to do. Net EF load: high.
- **Broad goals + minimal satisfaction:** The artifex writes a broad intent (light EF load upfront), the agent finds the cheapest path to technically-done, and the artifex is left with "oh we also need to..." Net EF load: high.
- **Broad goals + abundant satisfaction:** The artifex writes a broad intent (light EF load upfront), the agent brings its own judgment about what "all the way done" means, and the artifex receives work that doesn't generate followup. Net EF load: **low.** This is what we're optimizing for.

"Abundantly satisfied" is not a metric — it's an orientation. It means the agent asks: *"If the artifex looked at this result, would they feel it's completely handled? Or would they immediately think of three things that still need doing?"* An abundantly satisfied opus is one where the artifex's reaction is "great, I can forget about this" — not "great, now I need to..."

This applies at every level:
- **Opus level:** An azer discharging an opus should leave no obvious loose ends. Not just "fix merged" but "fix merged, related patterns noted, followup opus inscribed if warranted."
- **MO level:** An assessment azer evaluating whether the Magnum Opus is abundantly satisfied should think from the stakeholders' perspective — would they feel this is handled?
- **System level:** The athanor itself should be oriented toward making the artifex's life abundantly easier, not just technically less work.

**The tempering question:** "Abundantly satisfied" must also give agents permission to *stop*. The ceiling is the artifex's absence of concern, not the presence of perfection. The test: "Would the artifex feel they can forget about this?" — not "Is there anything else conceivably useful to do?" There always is. The agent that can't distinguish between loose ends (artifex will be surprised by followup) and nice-to-haves (artifex wouldn't think of this unprompted) will churn forever. Whether agents can make this distinction is an open question — if they can't, we'll need explicit guidance. For now, we're testing whether the framing is sufficient on its own.

Whether "abundantly satisfied" captures this precisely enough is an open question — we'll iterate on the language as we observe how agents interpret it. The concept is what matters: completion means "all the way done from the stakeholder's perspective," not "geas technically met."

**Goals and assessment:**

Goals should be concrete and assessable, but not necessarily empirical. For a single-operator system, qualitative feedback is valid and often superior to metrics. The operator is the only stakeholder — no alignment across multiple people needed.

- **Concrete qualitative goals:** "The operator feels that false positives and false negatives are decreasing" is a valid goal. Assessed by periodic structured check-in with the operator, not by dashboards.
- **Meta-goals:** "Solutions are becoming more autonomous over time" — measured by trend in operator interrupts, trend in false positive/negative rates, assessed qualitatively.
- **The daily reflection opus:** For state-based Magna Opera, a periodic opus (daily or every few days) gathers context and helps the operator assess: "Are we making progress at a rate that feels good? Are we hitting diminishing returns?" This is not ceremony — it's the feedback loop that makes the system self-correcting.
- **Termination:** A state-based MO is discharged when the operator says "we're hitting diminishing returns on investing in this system." A task-based MO is discharged when it is abundantly satisfied — not just when the literal completion criteria are met, but when the artifex would look at the result and feel it's fully handled.

**Crafting Magna Opera — the specificity boundary:**

A Magnum Opus should be the **broadest goal that aligns agent behavior with human-facing value**.

**Learning is a permanent first-order goal of the athanor — not a bootstrap phase.** The system always has two goals: (1) whatever the MO says, and (2) learning what guidance the system actually needs. This never graduates. The fool's path (discharged opera reflections) is the feedback loop. Overspecifying an MO doesn't just reduce agent latitude — it violates this mandate by preventing the system from learning where its boundaries actually are.

**Agent drift is signal, not failure.** A marut that drifts toward muting everything reveals the goal was misaligned (optimizing proxy metric). A marut that drifts under a well-aligned goal reveals what operational guidance is genuinely missing. Both are learning. Preemptive guardrails prevent both.

**The principle:** State the goal at the broadest level agents can assess clearly. Add operational guidance only when you have strong belief or evidence it's needed — from prior firings, from known LLM failure modes, from domain knowledge. If you don't have evidence, leave it out and read the trail. If the trail shows the agents couldn't assess the goal or drifted unproductively, sharpen — but only to the degree the evidence demands.

**Where to accept uncertainty:** A common mistake is optimizing for agent comprehension — writing the MO so you're confident the agent "gets it." Concrete metrics, explicit procedures, unambiguous terminal states. This feels safe but trades value-alignment for legibility. The agent understands perfectly and optimizes for the wrong thing. Instead, accept uncertainty in goal specificity (agent may need cycles to figure out how to assess it) rather than in value alignment (agent optimizes for proxy metric with perfect clarity). A slightly confused agent pursuing the right goal learns its way to effectiveness. A perfectly clear agent pursuing the wrong goal gets better at the wrong thing.

**What makes a bad MO:**
- Optimizes for a proxy metric instead of the actual desired outcome (e.g., "disposition every error" when the goal is channel health — the agent finds the cheapest path to the metric)
- Overspecifies tactics, preventing the agent from discovering better approaches and preventing the artifex from learning what guidance is actually needed
- Too vague to assess — "make things better" gives no way to evaluate whether the system is working

**What makes a good MO:**
- Names the actual human-facing value (what the operator wants to see/feel/experience)
- Assessable — the operator can look at the state and say "yes, this is what I wanted" or "no, this isn't it"
- Leaves maximum latitude for how to achieve it — the agent figures out tactics, the trail reveals what works

**A Magnum Opus is a product spec.** The operator is the user. The agents are the product team. The marut is the PM. A product team that doesn't know who they're building for ships features nobody wanted. A Magnum Opus that doesn't include operator context produces technically correct work that misses the point.

**Stakeholder personas in the MO.** The MO should name real people who care about this goal — not just "the artifex" as an abstract role. Include:
- **Identity:** Name, Slack handle, GitHub username — so agents can communicate with them, tag them in PRs, post updates to relevant channels.
- **Org role:** What they do. A staff engineer, a PM, a team lead — this shapes how the agent communicates (technical depth, business framing, etc.).
- **Why they care:** What makes this goal matter to *them* specifically. The same MO hits different stakeholders differently. An engineer cares about noise in #bugsnag because it obscures new regressions from their feature launch. A PM cares because it signals product health to leadership.
- **What satisfaction looks like:** From *their* perspective. Not the geas (which is verifiable completion criteria) but the human experience: "When I look at #bugsnag, I can immediately see if my feature caused problems."

This expands the product lens — agents that know their stakeholders can make better judgment calls about what's valuable, how to communicate, and when something is worth escalating vs. handling quietly. It also enables direct communication: an azer that knows the stakeholder's Slack handle can post a summary to the relevant channel, tag the right person in a PR description, or frame an escalation in terms that matter to the person receiving it.

**Not every MO needs personas.** If the only stakeholder is the artifex (e.g., infrastructure work that only the operator cares about), the existing "operator context" framing is sufficient. Personas add value when there are multiple stakeholders or when the people who benefit from the goal are different from the person running the athanor.

**Structure:**
- Single context entry point (Linear project, spec file, any URL/path) linking to everything relevant
- One or more goals with concrete completion criteria
- Operator context — who benefits, what they care about, what satisfaction looks like
- Opera inscribed in service of these goals (see Generative Principles)
- At least one opus always in progress — if not, the marut detects the stall and acts

**Status:** Concept used throughout this doc but not formally defined until now. Needs: standard MO definition format, registry of active MOs per athanor, daily reflection pattern for state-based MOs.

---

### Roles and Autonomy

#### Role

**What it's FOR:** Defining what an agent CAN do, CANNOT do, and how it interacts with the world. A role is a set of *boundaries*, not a configuration bundle.

**Problem it solves:** Right-sizing agent capabilities. A azer shouldn't merge PRs. A beholder shouldn't edit code. A marut shouldn't implement features. Without roles, every agent is omnipotent and the only constraint is prompting (unreliable).

**What a role IS:**
- **Purpose:** When and why you'd use this role. What kind of work is it for.
- **Boundaries:** What it can and cannot do. Enforced by hooks (hard) and injection (soft).
- **Interactions:** How it communicates with operator and other agents. Who it escalates to.
- **Expectations:** What "healthy" looks like (used by marut for stuck detection). How long silence is normal. What output patterns are expected.

**What a role is NOT:**
- Not an autonomy level. Autonomy is orthogonal — a azer can be supervised or autonomous depending on the opus. The role defines boundaries; the opus defines autonomy.
- Not a configuration format. The role definition *produces* configuration (hooks, sandbox settings, injections), but the role itself is about purpose and boundaries.

**Initial roles:**

| Role | Purpose | Boundaries | Expectations |
|------|---------|------------|--------------|
| **Azer** | Execute work in a sandboxed crucible. The primary worker. | Can: edit code, run tests, create PRs. Cannot: merge PRs, push to main, modify CI. | Output every few minutes. Silence > 10min = stuck. |
| **Beholder** | Watch for conditions and act on them. Channel monitoring, triage, project health. | Can: scan Slack/Linear, inscribe opera, send messages. Cannot: edit code, create branches. | Active scan loop. Silence > 5min = stuck. |
| **Marut** | One per MO. Ensure its Magnum Opus advances — agent health AND goal progress. | Can: scan tmux, approve permissions, Telegram, nudge, inscribe recovery opera. Cannot: edit code, create PRs. | Patrol loop. Silence > 3min = stuck. |
| **Artifex** | The human directing the athanor. Interactive pairing, goal-setting, judgment calls. | Standard Claude Code permissions — human is watching. | N/A — human manages. |
| **Primus, the Forge Lord** | Athanor-level supervisor. Ensures each MO has a marut, bootstraps new MOs, keeps standing processes alive. | Can: spawn/restart maruts, check athanor health, escalate to artifex. Cannot: edit code, manage individual opera. | Periodic check. Silence between checks is normal. |

More roles emerge as needed. Don't predefine what we haven't built.

**Implementation:** TOML file in `~/.claude/agent-roles/<name>.toml`. Consumed by `claude-run` (evolved) to configure hooks, sandbox, and session injection.

**Status:** Design exists. Not built. Depends on nothing — can ship immediately.

#### Autonomy (orthogonal to role)

**What it's FOR:** Determining how much an agent asks vs. acts on a per-opus basis. Independent of role.

**Problem it solves:** The same azer might be supervised (operator approves each phase) on a risky opus and semi-autonomous (just ship the PR) on a routine bug fix. Autonomy is a dial the operator turns per-opus, not a property of the agent type.

**Levels:**
- **Supervised:** Ask before each major step. Operator approves phase transitions.
- **Semi-autonomous:** Execute independently. Escalate blockers and ambiguity. Operator reviews output (PR).
- **Autonomous:** Execute to completion. Escalate only on failure. Operator gets a notification when done.

**Status:** Have it (autonomy profiles via session injection). Works. Stays orthogonal to roles.

---

### Orchestration Patterns

These are the different ways work gets coordinated. `/orchestrate` is ONE pattern. We need others.

#### Pipeline Pattern (`/orchestrate` — the athanor prototype)

**What it's FOR:** Transmuting a scoped opus into a shipped PR through phased execution.

**What it proved:** Phased execution with gates, task scoping, implementer delegation, verification discipline, context preservation. All of these ideas survive into the athanor — just decomposed into atomic components instead of one monolithic skill.

**Phases:** Triage → design → implement → verify → ship.

**Status:** Working prototype. Progressively being decomposed. `/task-triage` was the first extraction. Roles, geas-based completion, marut supervision, and dispatch are the next wave. Eventually the pipeline pattern is just: inscribe opus → gather materia → charge azer → marut watches → geas discharged → PR ships.

#### Patrol Pattern (beholder pattern)

**What it's FOR:** Continuously watching for conditions and acting on them. The autonomy bridge. This is what beholders DO — they watch, they scan, they act.

**Existing implementations:** `/permission-manager` (session health), `/fire-watch` (Slack issues), `/pulse-sel` (project health), `/ci-monitor` (CI status).

**Shape:** Loop: scan → classify → act → wait → repeat.

**Status:** Have it at skill level (4 implementations). Not formalized as athanor infrastructure yet.

#### Muster Pattern

**What it's FOR:** Mustering azers for discovered opera. The bridge from patrol (beholder finds work) to pipeline (azer does work). "Muster" as in mustering troops — organized mobilization, not bureaucratic dispatch.

**When to use:** A beholder inscribes opera, and you want azers to automatically pick them up.

**Shape:** Poll queue → for each opus: muster a azer in a sandbox → charge with opus → marut monitors.

**Status:** Does not exist. Need it for #bugsnag phase 2.

---

### Supervision

#### The Marut — one per Magnum Opus

**What it's FOR:** Ensuring its Magnum Opus advances. The marut is an Inevitable — a construct from Mechanus that relentlessly enforces cosmic law. In our system, the cosmic law is: *this Magnum Opus must make progress.*

**Exactly one active marut per Magnum Opus.** The marut owns everything about its MO's health — both the agents working under it and the goal they serve. This is not two separate concerns; a marut that only watches agent health without understanding the MO is useless (a azer can be "healthy" while the goal stalls), and MO-level supervision without agent health monitoring misses the most common failure mode (stuck agents).

**Problem it solves:** Agents get stuck silently. Context fills without warning. Permission prompts block progress. Error loops waste time. Opera chains break when agents die without inscribing follow-ups. Without the marut, the operator monitors all of this manually — the exact thing the athanor eliminates.

**Two levels of cosmic law (both owned by the same marut):**

**Agent-level:** Every agent with a geas must be making progress.
- **Azer:** Producing output (tool calls, edits, test runs). Silence > N minutes = stuck.
- **Beholder:** Completing scan cycles. Silence > N minutes = stuck.
- **Any role:** No permission prompt blocking. No context limit stall. No repeated error pattern.

**Magnum Opus-level:** The MO must be progressing. The marut watches the opera tree — if a chain stalls (agent died before inscribing follow-up, child opus has no agent charged, discharged geas didn't produce expected next opus), the marut detects the gap and acts. The generative principles assume happy path; the marut handles the unhappy path.

**The MO invariant:** At least one opus must always be in progress for the Magnum Opus. If the marut detects no active opera:
1. Attempt to inscribe the next opus (using MO context, entry point, and the fallback cascade from the generative principles)
2. If it can't determine the next opus → escalate to operator: "MO X has no active opera, I can't figure out what to do next"

**Why one per MO:** The marut needs context about its specific Magnum Opus — goals, accessible context, what "progress" looks like, which agents are charged with which opera. This is MO-specific knowledge that doesn't share well across a flat supervisor watching everything. Agent-level health checks (permission prompts, context limits) are mechanical and context-free; MO-level supervision requires understanding the goal.

**Actions when unhealthy:**
1. **Nudge** — send a message to unstick (permission auto-approve, gentle prod)
2. **Re-inscribe** — if a azer died mid-chain, inscribe a recovery opus ("pick up where azer-03 left off, here's the last known state")
3. **Escalate** — Telegram the artifex with: what's stalled, why, what to do
4. **Reforge** — kill the session and spawn fresh in the same crucible. The crucible endures; the session is reforged. The new session inherits the crucible's context (opus, plan, sandbox state) but starts with a clean context window.

**Status:** `/permission-manager` is the foundation. Already has: tmux scanning, condition classification, auto-resolve, Telegram escalation, scan loop. Needs: MO-awareness (context about the specific MO's goals and opera), broader detection heuristics, role-aware expectations, structured escalation format.

#### Primus, the Forge Lord — athanor-level supervisor

**What it's FOR:** Ensuring the athanor itself is operational — that each Magnum Opus has its marut, that minimal infrastructure is running, that the furnace is lit.

**The name:** Primus, the Forge Lord. The One and the Prime — ruler of Mechanus, the plane of pure law where maruts originate. In D&D, Primus commands all inevitables. In our system, Primus, the Forge Lord commands the maruts. In Latin, simply "first" — the prime mover. The Mechanus hierarchy mirrors our supervision hierarchy: Primus, the Forge Lord → marut → agents. Always referred to by full title.

**Problem it solves:** The marut watches its MO. But who watches the maruts? Who ensures a new MO gets a marut allocated? Who notices when a marut dies? This is meta-supervision — not watching agents or goals, but watching the supervision structure itself.

**Responsibilities:**
- Each active MO has exactly one marut allocated and running
- Dead maruts are detected and replaced
- New MOs get supervision bootstrapped (marut spawned, crucible created)
- The athanor as a whole meets minimal operational requirements (state directory exists, standing processes are running, configuration is valid)

**This is lightweight.** At our scale (1-3 active MOs), this may be a startup/bootstrap script that ensures standing processes are running — not a persistent agent. It runs when the athanor is kindled, periodically checks that maruts are alive, and restarts anything that's died. The concept of a distinct role exists for when scale demands it; the implementation can start as a cron job or a simple patrol loop.

**"Make progress confidently or escalate"** — this phrase encapsulates the agent/operator boundary at every level. An agent is either advancing toward its geas with evidence-based confidence, or it's pulling the andon cord. There is no "kind of working on it" middle ground. This principle applies at opus level (azer working on a fix), MO level (marut ensuring the goal advances), and athanor level (supervisor ensuring the furnace is lit).

**Status:** Conceptual. First implementation will likely be a bootstrap script invoked when kindling an athanor, plus a periodic health check (cron or simple loop).

#### Escalation Bus (infrastructure)

**What it's FOR:** Routing alerts from any agent or the marut to the operator through a consistent channel.

**Current state:** Telegram via `mcp__telegram-mcp__notify`. Used by `/permission-manager`, ad hoc elsewhere.

**Needed:** Consistent format: `[{source}:{role}:{crucible}] {severity} — {description} — {suggested action}`

**Status:** Partially exists. Needs standardization, not new infrastructure.

---

### Infrastructure

#### Whisper [Go utility — #15 from open-work-items]

**What it's FOR:** Inter-crucible communication. Reliably whispering between crucibles — marut to azer, dispatch to newly spawned crucibles, beholder to marut. The CLI tool is called `whisper`. "Whisper to crucible-03" = reliably deliver a message to the agent in that crucible.

**Problem it solves:** Our current `tmux-claude-send` (49 lines of bash) has known bugs:
- No `-l` flag (special chars misinterpreted)
- No copy mode handling (scrolled panes eat input silently)
- No ESC→Enter gap (readline M-Enter bug)
- No concurrent send protection
- No idle detection

These bugs are tolerable for manual use. They compound catastrophically under automation — every marut nudge, every dispatch spawn, every handoff could silently fail.

**`whisper`** — Go CLI that implements reliable inter-crucible communication. Core capabilities: `whisper send` (reliable message delivery to a crucible), `whisper idle` (wait-for-idle with 2-check confirmation), per-crucible send locks (no concurrent garbling). Start with the core, grow as needed.

**Status:** Built and operational. Protocol spec in `gastown-tmux-protocols.md`. Every automated interaction in the athanor goes through whisper — if whispers are flaky, everything is flaky.

#### Context Budget Guard [#20]

**What it's FOR:** Detecting when a session is approaching context exhaustion and preventing degraded output.

**Status:** Spec drafted at `context-budget-guard/spec.md`. Shell script + hook. No dependencies. Ship immediately.

#### Hooks as Hard Constraints

**What it's FOR:** Making certain agent behaviors impossible, not just discouraged. PreToolUse hooks that block dangerous operations based on role.

**Examples:** `no-push-to-master.sh`, `no-merge-pr.sh`, `bounded-scope.sh` (no edits outside working directory).

**Composability:** Role definition lists which hooks apply. Hooks are small, single-purpose scripts. Roles compose them.

**Status:** Designed, not built. Trivial individually — the system is the composability (role → hook selection → installation).

---

### Future: Athanor Dashboard (TUI)

Dwarf Fortress-style ASCII visualization of the athanor — the artifex's window into the furnace. Currently the artifex monitors via `tmux capture-pane` and squinting. A real TUI dashboard could show: active azers with context %, opera flow (charged/discharged), materia sources, tinctures produced, marut loop status, and the trail growing in real time. Not urgent — note for when the system has enough moving parts to justify it.

---

## Implementation Status

*Supersedes inline `Status:` notes throughout the doc — those reflect state at time of writing; this table reflects the current system.*

### Operational

Built, fired, and working.

| Concept | Where It Lives | Notes |
|---------|---------------|-------|
| Athanor instance pattern | `specs/athanors/<name>/` | AGENTS.md, magnum-opus.md, role files, opera dirs. Active instances: bugsnag (state-based), sal-117-l2-metrics (task-based), seismic-classifier-mapping (task-based). |
| Shared materia | `specs/athanors/shared/` | Universal AGENTS.md, azer.md, marut.md, muster.md, opus.md. Symlinked into each instance — change once, applies everywhere. |
| Magnum Opus format | `magnum-opus.md` per instance | Goal + abundant satisfaction + instance constraints + getting-started pointer. Intent only — no procedures, no discovery findings (those go in `specs/`). |
| Opus lifecycle | `opera/charged/` + `opera/discharged/` | Inscription / charge / discharge / trail. Datestamp filename prefix: `YYYY-MM-DD-<name>.md`. |
| Core geas + escalation-as-geas | `AGENTS.md` (shared) | "Both are equally valid fulfillments of your geas." Tested in first bugsnag firing. |
| Azer role | `azer.md` (shared) | Verification-first (three questions before mise), mise en place, context management, proof of fulfillment at discharge. |
| Marut role | `marut.md` (shared) | Operational loop, assessment opera, monitoring/stall detection, reforging. |
| Muster protocol | `muster.md` (shared) | Crucible kindling, reforging, monitoring, cleanup. Marut → azer dispatch. |
| Assessment opus pattern | `AGENTS.md` + `opus.md` | Inscribed when no opera exist. Dual purpose: (1) determine what's next, (2) walk the trail as archivist — discharged opera are ephemeral, extract knowledge into specs/tickets/Slack before it's buried. |
| Verification-first | `azer.md` (shared) | Three questions before mise: prove fulfillment, communicate proof, ensure actualization. Mise en place is second. |
| Elicited acceptance | `azer.md` (shared) | Stakeholder feeling as valid evidence when deliberately elicited. Fourth spectrum entry. Alchemized UAT. |
| Proof of fulfillment | `opus.md` (shared) | Third required discharge element: evidence the intent was met and the work is actualized (not just technically complete). |
| Abundant satisfaction + proof | `AGENTS.md` (shared) | "Proof, not feeling." Discharge record is the proof; verification loops are the mechanism. |
| Homunculus | `kadmon.md` | Artifex familiar — interactive, not autonomous. Assists supervision, materia updates, trail-walking. |
| Primus (manual) | `kadmon.md` | Artifex acts as Primus. Launch runbook + athanor init protocol documented. |
| Context budget guard | Hook | Advisory warnings at 75/85/92%. Fail-open. |
| Whisper CLI | `whisper` binary | Reliable inter-crucible communication. Built and tested. |
| Sandbox isolation | `wtp add` | Worktrees per azer. Independent branches and Docker environments. |
| Autonomy profiles | Session injection | Semi-autonomous default. Profiles via `$CLAUDE_SESSION_ID`-keyed injection directory. |
| First-principles PM framing | `spec.md` | Alchemized PM: intent → verification → proof → actualization. Elicited acceptance as alchemized UAT. No PM terminology in agent materia. |

### Designed, Not Built

| Concept | Priority | Notes |
|---------|----------|-------|
| `whisper` nudge queue | Low | Async delivery with TTL/priority. Not a current bottleneck — basic send works. |
| Roles (TOML + claude-run) | High | TOML role definitions + claude-run evolution. Dependency of hard constraints. |
| Hooks as hard constraints | High | `no-push-to-master.sh`, `no-merge-pr.sh` etc. composable per-role. Trivial individually; needs role system to compose. |
| Role-aware crucible creation | Medium | Depends on roles. |
| Beholder role | Medium | Needed for #bugsnag phase 2 (auto-dispatch). |
| Primus as agent | Low | Currently manual. Could be a cron/bootstrap script. Not a bottleneck yet. |
| Escalation bus format | Low | Telegram works; message format not standardized. |
| Handoff protocol | Medium | `/direct-handoff` is crude. Full protocol needs reliable whisper. |
| Daily reflection opus | Low | For state-based MOs. Pattern defined; not running yet. |

### Infrastructure Gaps (Discovered in Operation)

Gaps the running athanor has surfaced. Live list in `kadmon.md § Infrastructure Gaps`.

| Gap | Impact |
|-----|--------|
| No production DB access | Azers can't diagnose production-data errors. Frequent escalation. Read-only SQL user needed. |
| No auto-merge gate | Azers create PRs but can't merge. Artifex bottleneck for high-volume work. |
| Worktree cleanup not automated | Each azer creates a worktree. Manual cleanup accumulates. `wtp remove` should be in post-discharge flow. |

### Deferred / Conceptual

| Concept | Notes |
|---------|-------|
| Athanor dashboard (TUI) | Dwarf Fortress-style operator view. Not prioritized. |
| Formal athanor registry | Currently implicit: `specs/athanors/` directory listing. |

---

## Dependency Graph

```
Whisper CLI [BUILT]             Context Budget Guard (#20)
  [no deps]                       [no deps]
       │                               │
       ▼                               ▼
Roles (#1)  ────────────────────────────────────► All agents get:
  [no deps for MVP]                                - Role boundaries
                                                   - Context monitoring
                                                   - Hard constraints
       │
       ├───────────────────────────┐
       ▼                           ▼
Hooks as Hard                 Marut (evolved
Constraints                    /permission-manager)
  [needs: roles]               [needs: roles for
                                per-role expectations]
       │                           │
       ▼                           ▼
#bugsnag Phase 1              #bugsnag Phase 2
  (Beholder: triage)           (Dispatch: auto-fix)
  [needs: beholder role]       [needs: azer role,
                                marut watching,
                                dispatch pattern]
```

Two things can start in parallel with zero dependencies:
1. **Context budget guard** — shell script + hook
2. **Roles MVP** — TOML definitions + `claude-run` evolution

---

## v1: Light the First Crucible

Everything below is needed for autonomous #bugsnag. Each piece delivers incremental value on its own. The system will creak at first — 50% of daily errors resolved autonomously is a win. The other 50% escalate, and each escalation teaches the system. We don't need perfection; we need the loop running.

### Infrastructure (build first, enables everything)

**Context budget guard** — shell script + hook. Already spec'd. Prevents silent context degradation across ALL sessions, not just #bugsnag.

**`whisper` Go CLI** — built. Reliable NudgeSession, WaitForIdle, send locks. All automated inter-crucible communication goes through whisper.

**Roles** — TOML definitions for operator, azer, beholder, marut. Evolve `claude-run` to consume role definitions. Hooks as hard constraints (`no-push-to-master.sh`, `no-merge-pr.sh`, `bounded-scope.sh`) composed per-role.

**Opus format** — standard markdown template. Goal + geas + context + parent reference. Simple enough for a beholder to generate automatically, rich enough for a azer to pick up and execute.

### The #bugsnag Daily Magnum Opus

Replace the phased rollout with the real model: a **daily Magnum Opus** generated from a template.

**Template: `#bugsnag daily triage`**

Input: date. Output: a Magnum Opus.

```
Magnum Opus: "#bugsnag daily triage — 2026-03-20"
Geas: All errors from today are actioned (fixed, muted, or escalated to operator)

1. Context gathering:
   - Scan #bugsnag for today's errors
   - Check recent "fixed" errors — are they actually resolved?
     (or leverage AppSignal error lifecycle if elegant)
   - Classify each error: new vs recurring vs regression

2. For each error, inscribe an opus:
   - Clear, simple fix → opus to fix (may use /orchestrate, or simpler)
   - Unclear root cause → opus to /rca, then generative principle produces next opus
   - Can't classify / too complex → opus to escalate to operator
   - Already fixed / transient → opus to mute or confirm resolution

3. Magnum Opus discharged when all child opera are actioned
```

**Menu of approaches, not a rigid pipeline.** The beholder (or whatever agent does context gathering) doesn't prescribe HOW to fix — it inscribes the opus with appropriate fidelity and lets the charged agent decide the approach. Common patterns the agent might reach for:
- Simple bug → fix directly, ship PR
- Needs investigation → `/rca` first, generative principle produces next opus from findings
- Needs code change + test → `/orchestrate` for the full pipeline
- Recurring error → check if prior fix regressed, may need deeper investigation
- Not actionable → mute in AppSignal, note rationale

**The key:** the daily Magnum Opus is the **input**. Everything else — decomposition, approach selection, execution, verification — emerges from the generative principles. No rigid pipeline. The beholder inscribes opera based on what it sees, agents discharge their geas and inscribe next opera based on what they learn, chains terminate when errors are actioned.

**Start small:** First iteration is just the beholder scanning + inscribing opera. Human reviews the opera and decides which to action. Then progressively: charge azers automatically for simple fixes, let agents run `/rca` autonomously, marut watches everything. The Magnum Opus template stays the same — autonomy increases within the same structure.

---

## Open Questions

### Handoff Protocol

`/direct-handoff` is crude. `/handoff` is tangential (opus definition, not session lifecycle). We need a lightweight handoff protocol:

1. Save state (commit, update plan/opus, write checkpoint)
2. Spawn fresh session in same crucible
3. Nudge new session with context ("continue from opus X, plan Y, you just finished Z")

Our version: a shell script that uses the Go tmux utility for reliable delivery.

### Injection Layers (Current Architecture)

The athanor has four injection layers, from most persistent to most ephemeral. Understanding these is key to deciding where mandates go — the same principle needs to show up at the right layer(s) to stick.

| Layer | When Loaded | Persists | What Belongs Here | Examples |
|-------|------------|----------|-------------------|----------|
| **CLAUDE.md** | Always present, read at session start | Entire session | Project conventions, universal mandates, skill routing, build commands | "Never push to master", "run tests in Docker", escalation guidance |
| **System prompt** | Set at session start via `claude-run` | Entire session (but can drift under compaction) | Role identity, role boundaries, geas framing, behavioral core mandates | "You are a azer. You are under a geas to..." — the compulsion toward progress |
| **Session injection** | Per-turn via `UserPromptSubmit` hook | Refreshed every turn (compaction-resistant) | Critical mandates that drift, opus state, orchestration directives | Core mandate reminders (existing PostToolUse hook does this), task context |
| **Opus/plan** | Read by agent at startup or when charged | Durable (files on disk) | What to work on, geas, materia references, context entry point | opus.md, task-plan.md, spec links |

**Directional rubric for placing mandates:**
1. Start with CLAUDE.md if it's universally applicable
2. System prompt for role-specific identity and boundaries
3. Session injection for mandates that are critical AND prone to drift (compaction-resistant layer)
4. If something isn't sticking at a higher layer, escalate it to session injection — that's the reinforcement layer

**What already works:** The PostToolUse hook already injects behavioral reminders periodically ("DIAGNOSE before fixing, ESCALATE don't workaround..."). This is injection in depth in action — the same mandates from the system prompt, reinforced per-turn. The pattern works; we need to extend it to athanor-specific mandates (geas framing, chain of evidence, assay discipline) as roles are built.

**Open work:** As each role definition is created, decide which mandates from the domain model need injection at which layer. Trial and error — start with CLAUDE.md/system prompt, escalate to session injection for anything that drifts.

### Plan Rigidity

Current plan format is tightly coupled to `/orchestrate`'s phase model (DESIGN → FEATURE-INIT → IMPLEMENT → VERIFY → SHIP). For non-orchestrate patterns (patrol, dispatch), opera don't need plans at all — or need lightweight checklists, not rigid phase definitions. The plan is optional; the opus and geas are sufficient for simple work.

---

## The Generative Principles

What happens when a geas is discharged? Work begets work — until the chain naturally terminates. These principles govern how opera come into being.

**0. Assay before transmuting — at every level.**

Before executing any opus, assess: *do I have the tools and context to complete this effectively?* This is not just "gather materia" — it's an active evaluation of readiness that identifies gaps and opportunities. The assay may itself produce opera.

**What an assay asks:**
- What sort of tasks will this opus involve? What tools and context do they need?
- Do we have the materia to test our completion criteria? (Autonomous testing tools, verification infrastructure, access to relevant systems)
- What's missing that would make this significantly easier?
- Is the gap fillable by an agent, or does it need operator judgment?

**The assay is recursive.** At Magnum Opus level, it's broad and somewhat speculative — "what kind of project is this and what infrastructure do we need?" At individual opus level, it's concrete — "do I have the test fixtures and API access to verify this fix?"

**Gaps become opera.** If the assay reveals missing materia, the gap itself becomes an opus:
- Missing testing infrastructure → opus to build it (like Validated Webmocks — operator identified that VCR was painful for API client testing and invented a better approach before the project started)
- Missing product context → opus for context collection (like SAL-98 — massive Slack/Notion/GitHub scan that found key product requirements like KD-1 value mapping granularity, which the operator barely participated in — the agent found it better)
- Missing access or permissions → opus to escalate to operator

**Why this is first-order, not a nice-to-have:** Without assay, agents "just do their best" with whatever they have — and wheel-spin on problems that would have been trivially solved with the right preparation. The difference between a 4-hour slog and a 20-minute fix is often materia that could have been gathered in 5 minutes at the start.

**Search recent opera for precedent.** When inscribing or defining an opus, search for similar completed opera — what was the geas? How was it discharged? What went well, what didn't? (This requires that transmuted opera have their reflection/outcome readily accessible, not buried.) Use prior opera to inscribe better ones. Whatever patterns are consistently successful can eventually bubble up to `/task-triage` or role definitions, but the search-for-precedent habit means the system is self-learning from day one, before anything is formalized.

**Agents can identify mise opportunities BETTER than the operator in many cases.** The operator brings judgment about tool design (Validated Webmocks) and strategic direction. Agents bring exhaustive scanning capacity (SAL-98 context collection found product requirements across Slack, Notion, and GitHub that no human would have found efficiently). Both are assay.

---

**1. The generative act: inscribe what furthers unfulfilled goals.**

When an agent discharges its geas, it looks at the goals in its domain — the broader context the opus served. Are there completion criteria still unfulfilled? If so, inscribe one or more opera to further them. Or don't — if the goal is satisfied, the chain terminates naturally.

The rubric: "There are unfulfilled goals in my domain that aren't appropriate for me to fulfill. How do I further them?" The answer is always an opus — the question is how specific.

**2. Fidelity matches context — the triage fallback.**

The inscribed opus is as detailed as the agent's context allows. This creates a natural spectrum:

| Agent's context | What it inscribes |
|----------------|-------------------|
| Rich — knows exactly what to do | Detailed opus: "fix the N+1 query in ProofRecommendationService#call, here's the approach, here are the files" |
| Partial — sees the direction | Directional opus: "investigate reviewer comment on PR #4341 line 47, may need /rca, then fix" |
| Minimal — knows something's unfulfilled but not what to do | **Triage opus: "there's an unfulfilled goal, figure out the next steps"** |

**The triage opus IS the default.** When context is insufficient to inscribe something more specific, the most basic generative act is: "inscribe an opus to go figure out what to do." This is what `/orchestrate`'s Stage 0 (triage / `/task-triage`) does — it's not a special pipeline phase, it's the **fallback opus that every agent can produce** when it doesn't know enough to be more specific.

This means `/task-triage` isn't just an `/orchestrate` component — it's a fundamental athanor primitive. The lowest-fidelity generative act.

**2a. Inscribe with gaps, not conclusions.**

When inscribing an opus for another agent, present observations and gaps — not polished conclusions. "I found X and Y, but couldn't determine Z — the next agent should investigate Z before proceeding" is better than "the answer is X, implement it." The next agent has fresh context and may see things the inscribing agent missed. Over-specifying the approach in the inscription closes off better paths the next agent might find. The opus should channel intent (goal + geas); the approach is for the charged agent to determine.

**3. Don't scope-creep the current agent.**

Inscribing the next opus is lightweight — a few lines of markdown. The agent doesn't plan, decompose, or orchestrate. It doesn't gather materia for the next opus. It says "here's what I think comes next" and exits. The next agent picks it up.

This keeps each agent focused on its geas. The generative act is a side effect of completion, not additional scope.

**3a. Three scope layers — redirect or escalate.**

When an agent encounters something outside its current opus scope, there are three concentric boundaries:

| Scope | What's in it | When something is outside it |
|-------|-------------|------------------------------|
| **Opus scope** | What the geas requires | Agent discovers work that furthers the MO but isn't its geas |
| **Athanor scope** | What the athanor can handle autonomously | Agent discovers work that needs human judgment, access, or authority |
| **Operator scope** | What the human handles | N/A — the operator is the outermost boundary |

**The decision at each boundary:**

- **Outside opus, inside athanor:** Redirect — inscribe a new opus, do context preservation, terminate self. The agent has context to write a good opus (goal, geas, relevant findings). The new opus gets picked up by another agent or the marut. This is the **happy path** for scope boundaries — work flows forward without operator involvement.

- **Outside opus, outside athanor:** Escalate to operator. The agent can't determine whether this is solvable autonomously. Pull the andon cord.

- **Ambiguous:** If the agent has context to write a solid opus (clear goal, concrete geas), it should redirect. If it can't articulate what the next opus should accomplish, it should escalate. The test: "can I write a geas for this?" If yes, redirect. If no, escalate.

**The "what do I do next?" fallback:**

At the Magnum Opus level, there must always be a way to determine the next opus. This is the universal starting point and the universal fallback when chains break. The resolution is always binary:

1. **Inscribe an opus** — even if it's just "explore the MO context and figure out what to work on next"
2. **Escalate to operator** — "I can't determine meaningful next steps for this Magnum Opus"

This binary resolves at the supervisory level (marut or MO-level supervision). A azer that can't figure out what's next terminates — the supervisor picks it up. No agent needs to solve "what do I do next?" alone. The happy path is that each azer inscribes a follow-up opus using its rich context from just having done the work. The fallback is that the supervisor detects the gap and figures it out from the MO entry point.

**4. Goals decompose through opera.**

An opus exists in service of some broader goal. The geas is the immediate "done when." But the goal may require multiple opera to satisfy. When one opus is discharged, the agent checks: is the parent goal satisfied? If not, what opera would further it?

This is recursive — goals decompose into opera which serve sub-goals. But it's not planned top-down like `/orchestrate`'s phases. It's discovered bottom-up: each agent, on completing its opus, looks at what's still unfulfilled and inscribes accordingly. Sometimes the decomposition is obvious (implementation naturally follows design). Sometimes it's not (research results determine whether to implement or backlog). The system handles both because the generative act happens at runtime, not planning time.

**5. Two layers: Magnum Opus → opera. That's it.**

No deep nesting. A Magnum Opus has goals and completion criteria. Opera belong to a Magnum Opus and are always immediately actionable — no draft opera. If you can't figure out a clear next step, either inscribe a research opus or escalate to operator.

```
Magnum Opus: "AppSignal is clean and stays clean"
  Entry point: Linear project SAL-XX / specs/athanors/bugsnag/magnum-opus.md
  ├── Opus: "Fix NoMethodError in ProofRecommendationService"
  ├── Opus: "Fix timeout in SalesforceSync"
  ├── Opus: "Investigate vague alert — /rca then proceed"
  └── Opus: "Set up beholder patrol for ongoing monitoring"
```

**Every opus is ready for immediate action.** There's no concept of a draft opus, a planned opus, or a dependency graph between opera. If it's inscribed, it can be charged to an agent right now. **Plan however you want, but an opus has no hard dependencies.** Dependency analysis is valuable — document what depends on what, create detailed plans, note "do A before B" in plan docs or in an opus's context. But an opus itself never formally blocks on another opus. Any inscribed opus can be charged to an agent immediately. Sequencing is advisory (in documentation) and emergent (through the generative act), never structural.

**Corollary: all opera under a Magnum Opus are parallel-executable.** If three opera exist, three azers can be charged simultaneously, each in their own sandbox. No coordination needed — sandboxes (worktrees) provide isolation.

**Systemic presumption of parallelizability.** All extant opera are presumed parallel-executable, but with no guarantee. It's fine to think about parallelism, but never slow down or get paralyzed by dependency analysis.

**An agent is only responsible for its own opus.** A azer focuses on discharging its geas — not on cross-opera coordination or the big picture. The marut owns the wide view.

**The binary choice: proceed or andon.** At any given point, an agent is either making progress or pulling the andon cord. There is no middle ground — no "kind of stuck," no "waiting to see," no "trying something that probably won't work." This is kanban: full speed ahead, or stop the line. An andon can stop a single opus (escalate this one issue) or halt the whole Magnum Opus ("something fundamental is wrong, all work under this goal needs to pause"). The marut monitors for situations where andon should have been pulled but wasn't — agents spinning their wheels without progress. "Planning ahead" means dumping context into documents (spec doc, notes, breadcrumbs) — not inscribing opera for later.

**The Magnum Opus has a single context entry point.** One URL or path that links to everything relevant. Could be a Linear project, could be `specs/<project>/spec.md`, could be anything. During assay, you make sure all useful context is reachable from this entry point — directly or through breadcrumbs.

The LLM is good at following breadcrumbs. If the entry point is `specs/bugsnag/spec.md`, the agent intuitively knows to search the `specs/bugsnag/` directory. If there's a link to a Slack thread buried somewhere in there, it'll find it. Perfect organization is nice but the main thing is: **there is a clear trail from entry point to all relevant context.**

**Any agent without clear context goes back to the entry point.** That's the reset — follow the breadcrumbs, assay what's available, then inscribe the next opus based on what you find.

**The fallback cascade:**
1. Agent has rich context → inscribes detailed opus with clear geas
2. Agent has partial context → inscribes directional opus ("fix this, approach TBD")
3. Agent has minimal context → inscribes research/triage opus ("figure out next steps")
4. Agent has no context → escalates to operator

**6. Chains terminate naturally.**

Not every discharged geas produces a next opus. When the Magnum Opus goals are satisfied, the chain ends. No make-work. "And then what?" Sometimes the answer is "we're done."

---

## The Trail — Opera as Actualization

An opus has exactly two states: **charged** (being worked on now) or **discharged** (completed, part of the historical record). No backlog, no queue, no "planned," no "blocked."

The sequence of discharged opera under a Magnum Opus is the **trail** — the actual steps the athanor took toward fulfilling its goals.

**The core operation of the athanor is walking this trail.** Not planning, not context-gathering, not materia preparation — those are critically important but they are *byproducts* and *gateways*. The athanor produces plans and context in service of the trail, not the other way around. The operation is: inscribe the most urgent opus, discharge it, inscribe the next. Everything else supports this cycle.

### What the trail IS (behavioral guidance)

**Each step is taken with evidence, not with a map.** The trail is walked forward, not planned in advance — but every step is grounded in evidence from the last. Every inscription is backed by chain of evidence (what do I know? what supports this being the right next step?). Every discharge is bounded by a geas (concrete, observable completion criteria). The marut watches the trail to ensure it's advancing toward the Magnum Opus goals. These are not optional guardrails on an otherwise freewheeling process — they ARE the process. The discipline is: *one grounded step at a time, each step fully committed, each step informed by what the last step revealed.*

**Opera are actualization, not intent.** A plan says "we should do X, then Y, then Z." An opus says "X is happening now." Plans, specs, and context documents are intent — they inform inscription but they are not opera. An opus materializes at the point of action, when the athanor commits to actualizing a specific next step. This separation keeps the operational kernel minimal: an agent needs to understand opus/geas/inscribe-or-escalate to participate. Everything else is materia.

**Context is charge. Discharge toward the goal.** Every time an agent accumulates significant context — scanning a channel, investigating an error, reading discharged opera — it builds a rich understanding. That understanding is energy. The question at every decision point is not "what's the obvious next task?" but **"what is the most valuable thing I can do for the goal, given everything I now know?"** Be expansive about what that can mean. Be focused on the goal. The answer might be the obvious fix, or it might be something no previous agent thought of. A azer that notices its bug has three siblings from the same root cause. A marut that realizes the channel would be clearer if it posted fix summaries. An agent that inscribes a brainstorming opus because the standard approach feels suboptimal. This is not make-work — this is the agent serving the goal better than mechanical task execution can. The trail is context, not a template. What previous agents did is data, not a prescription.

**Success is the best next step from where you are.** Past and future are irrelevant to what "success" means right now. A Magnum Opus could have been one long chain of fuckups — but if right now the agent takes the best available action to improve the situation, that is 100% success. Escalating to the artifex after hours of thrashing is 100% success. Inscribing a triage opus because everything is confused is 100% success. The trail behind you doesn't define the quality of your current step. This eliminates the sunk cost fallacy: an agent never needs to "recover" from a bad trail. It just takes the best next step. Every step is evaluated on its own merits — was this the best action available given current evidence? If yes, success. If the best action is "stop and escalate," that's success too.

**Injection in depth note:** This principle is philosophically sound but won't override LLM completion-bias on its own — LLMs resist escalation because training rewards task completion, not because of sunk cost reasoning. This principle needs to be baked deeply into the system: system prompt for every role, session injection reminders, geas framing ("your geas is to take the best next step, which may be escalation"), and concrete behavioral triggers (the binary choice, the evidence test). Stating it once in the domain model is insufficient. It must be reinforced at every layer where completion-bias could override it.

**The discharged opera are the authoritative record.** Want to know the status of a Magnum Opus? Read its discharged opera — the trail speaks for itself. The marut assesses MO health from the trail, not from a plan. Precedent search (generative principle 0) queries discharged opera — what the system *actually did*, not what was planned. Learning is grounded in reality. Zero overhead for status tracking; the work IS the documentation.

### Why this is the operational kernel for autonomous agents

The scaling bottleneck for multi-agent systems is how much operational protocol each agent must internalize to participate. The trail reduces the agent's operational model to three concepts:

1. **Opus** — charged or discharged, the unit of work
2. **Geas** — when am I done
3. **Inscribe or escalate** — what do I do at boundaries

That's the entire mental model a azer needs. Everything else — Magna Opera, the trail as a named concept, plans, materia, supervision hierarchy — exists at the infrastructure level for the artifex, the marut, Primus, the Forge Lord. The azer doesn't need to know about the trail. It just walks it.

This is why it enables autonomy at scale: simple individual rules producing complex emergent behavior. The power comes from the simplicity of each agent's operational kernel, not from the complexity of the coordination layer.

### The trail's shape

**Task-based Magna Opera** produce a linear trail — discharged opera trace a trajectory from inception to completion. The trail may branch and meander (research detours, feedback loops, redirects) but it terminates when the MO goals are satisfied.

**State-based Magna Opera** produce a cyclical trail — the discharged opera trace a spiral of maintenance, each cycle refining the system's ability to maintain the ideal state with less effort. The path of return.

### Quiescence

When all task-based MOs under an athanor are discharged and all state-based MOs are in their ideal state, the athanor is **quiescent** — the furnace is lit but idle. Not cold, not off, but in a state of maintained readiness. Patrols still run, the marut still watches, Primus, the Forge Lord still ensures the structure is sound. Quiescence is not completion — the athanor is ready to kindle new Magna Opera or respond to deviations from ideal state. It is the resting heartbeat of a living system.

---

## Chain of Evidence Principle

**At every decision point, build a chain of evidence for your next step. If you can't, escalate.**

Surprises happen anytime — mid-assay, mid-implementation, mid-transmutation. When they do: be methodical, seek context, attempt to resolve. But the test is always: *can you point to evidence supporting your plan to proceed?* Not "I think this will work" — actual evidence (code you read, test results, docs you found, prior opus outcomes).

**The escalation ladder:**
1. Surprise occurs → look for context to resolve it
2. Found context → does it support a clear next step with evidence? If yes, proceed.
3. Context is contradictory or stale → attempt to resolve contradictions
4. Still can't build a chain of evidence for how to proceed → **escalate**

This is "diagnose before fixing" made concrete. The trigger for escalation isn't "I'm stuck" (subjective, agents talk themselves out of it) — it's "I cannot build an evidence-based plan for my next step" (observable, harder to rationalize away).

**Vibes-based andon triggers for compounding drift.** Beyond the evidence ladder, all agents are encouraged to pull andon on gut feel when:
- Hallucinations are creeping in (generating details not grounded in evidence)
- Proper operator involvement was skipped earlier and decisions are compounding on shaky ground
- Chain of evidence has reference gaps (assertions without sources)
- Too much agent decisionmaking without operator input on a high-stakes opus

These are vibes-based heuristics informed by the overall agent/operator boundary principles. LLMs can do reasonable vibes-based assessment when the philosophical context is solid. Start loose, tighten if too many false positives. Better to over-escalate than to compound errors silently.

---

## Context Preservation Principle

**Don't let insights die with the session.**

When an opus is discharged (or when a session ends for any reason), the discharging agent preserves what was learned. This is a core mandate — not optional reflection, not nice-to-have logging. Context that isn't preserved is context that must be rediscovered, at the operator's attentional expense.

**What to preserve:**
- What went well and what didn't (brief, honest)
- Anything learned that should update the system (role definitions, assay templates, MO context)
- Documentation updates (specs, architecture docs, CLAUDE.md)
- External updates (Slack messages, Linear ticket status/comments, PR body updates, Telegram notification to operator)
- Raw learnings that don't have an obvious home → dump to markdown in project dir (`specs/<project>/learnings/` or similar)

**Where to preserve:**
- Opus outcome/reflection → appended to the opus file (so precedent search finds it)
- Documentation → targeted updates to the relevant docs
- Notifications → Slack, Telegram, Linear as appropriate
- Everything else → markdown in project dir, findable by future assay

**The test:** If a fresh azer picks up a similar opus tomorrow, can it find what the previous azer learned? If not, the context preservation failed.

**Relationship to generative principles:** Context preservation happens ALONGSIDE the generative act (inscribing next opus), not instead of it. Discharge geas → preserve context → inscribe next opus (or terminate chain). Both the forward-looking act and the reflective act happen at opus completion.

---

## Operating Principles

1. **Sessions are cattle, state is in files.** Don't fight context exhaustion — plan for it. All durable state lives in git, plan docs, and specs. Sessions can die at any time.

2. **Hard constraints are hooks, soft guidance is prompting.** Never use prompting for things that MUST NOT happen. Hooks block; prompts suggest.

3. **Autonomy is orthogonal to role.** Roles define boundaries (what you CAN do). Autonomy defines initiative (how much you ASK vs ACT). Same role, different autonomy per opus.

4. **Conservation of operator attention.** The operator's time is the scarcest resource. Every design decision asks: "does this reduce or increase the operator's attentional load?" Intelligent escalation > noise.

5. **`/orchestrate` is the prototype, not the destination.** It lives on inside the athanor while we progressively extract its responsibilities into atomic components. New patterns (patrol, dispatch, supervision) don't wrap `/orchestrate` — they replace pieces of it.

6. **Reliable delivery is infrastructure.** Every automated interaction goes through tmux. If tmux delivery is flaky, everything built on top is flaky. Invest in the Go utility early.

7. **Incremental, not revolutionary.** Every piece ships independently. `/orchestrate` keeps working. No big-bang migration.

8. **Names matter.** Distinctive names create precise meanings. "Azer" is unambiguous; "worker agent" is not. The vocabulary IS the domain model.

9. **Weirdness boundary.** The athanor vocabulary is internal — skills, hooks, session injection, architecture docs. External output (Linear tickets, PRs, Slack messages) uses normie language. Translation is mechanical.

10. **Make progress confidently or escalate.** At every level — opus, MO, athanor — the binary is the same: advance with evidence-based confidence, or pull the andon cord. There is no "kind of working on it." This encapsulates the entire agent/operator boundary in one principle. Injection in depth demands this shows up everywhere: in role definitions, in geas framing, in session injection, in marut detection heuristics.

11. **The athanor is instantiated, not global.** Each domain gets its own athanor with its own Magna Opera, supervision, and state. The concepts are universal; the instances are scoped.

12. **Generative guidance over restrictive guidance.** Define clear structure that empowers agents to act, not structure that constrains them into compliance. Generative guidance says: here's the purpose, here's the quality test, the form is yours. Restrictive guidance says: do it this way, don't do that. Generative guidance produces agents that exercise judgment. Restrictive guidance produces agents that follow templates. The athanor needs judgment. Use restriction where it's warranted: hard safety boundaries (never merge PRs), known LLM failure modes (don't guess — escalate, because LLMs reliably hallucinate under uncertainty), and cases where we know the system or agent will go against goals without a hard guardrail. The design choice between generative and restrictive is constant and deliberate — every piece of guidance should be one or the other for a reason.

---

## Bootstrap: The Adam Kadmon Approach

**See also:** `kadmon.md` for the detailed bootstrap spec.

The first athanor is bootstrapped through deliberate minimalism. The evolutionary process of building the athanor IS the artifex's Magnum Opus — but this MO is not defined within the athanor itself (not yet). The #bugsnag athanor is the test bed for the system, not the meta-goal.

### The Principle: Prima Materia

The Adam Kadmon (primordial form) is the first marut running from minimal materia — just enough to orient and produce useful escalations, not enough to "cheat" by reading the full blueprint.

**What goes in the materialized athanor (`specs/athanors/<name>/`):**
- `AGENTS.md` — universal vocabulary, core principles, constraints (loaded by all agents)
- `magnum-opus.md` — entry point with concrete goals and daily cycle
- Role files (`marut.md`, `azer.md`) — lean, self-contained
- `opus.md` — inscription/discharge protocol and litmus test
- `opera/charged/` and `opera/discharged/` — where the trail lives
- All domain terms defined simply — agents should never escalate about vocabulary

**What stays in the blueprint (`specs/athanor/`) — off-limits to agents:**
- This spec (the full domain model, architectural principles)
- `kadmon.md` (the bootstrap design rationale)
- Historical analysis and comparison docs (in this directory)
- System-level architecture (injection layers, dependency graphs, etc.)

### The Evolutionary Rhythm

1. **The artifex kindles a Magnum Opus** with concrete goals and a well-crafted marut geas
2. **The marut runs** — scanning, inscribing opera, escalating gaps
3. **Each escalation is a step on the artifex's trail** — fill the gap, build a tool, refine a geas, add context
4. **The discharged opera tell us what worked** — retrospective on the trail drives system improvements
5. **The materialized athanor grows** — each filled gap makes the system more capable, the marut more autonomous

### Key Decisions (locked)

- **No agent may run `/orchestrate`** — its geas conflict with athanor principles. Agents may reference patterns but not be bound by them.
- **No agent may read `specs/athanor/`** — the blueprint is the artifex's working space, not the materialized system.
- **The artifex is Primus, the Forge Lord for now** — manually launching, monitoring, restarting. No need for automated primus until there are multiple Magna Opera to supervise.
- **Geas quality is foundational** — the artifex's highest-leverage work is writing good geas. A well-crafted geas channels LLM completion-bias toward the right outcomes AND makes escalation feel like compliance. The geas is the primary injection point for behavioral shaping.
- **Plans are optional context, not first-order** — goals + concrete verification (geas) are the success guarantee. Many planning styles can succeed with clear goals.
- **Provide minimal materia, test escalation** — the first test is whether the marut escalates cleanly when it hits gaps vs. guessing/hallucinating. Each clean escalation drives the artifex's build-out work.

### What Success Looks Like on Day 0

The marut scans #bugsnag, hits real gaps (no tooling to spawn azers, insufficient context to classify errors, etc.), and produces clear, useful escalations that tell the artifex exactly what to build next. The quality of the escalations — not the number of errors actioned — is the measure of day 0.

---

## What This Document Is NOT

- **Not a maximalist vision.** This is 1-2 weeks of work. We'll learn and adapt.
- **Not ignoring /orchestrate.** `/orchestrate` is the prototype of the athanor. It lives on inside while we extract its responsibilities into atomic components. It fades away naturally as the athanor matures.
- **Not the materialized system.** This spec is the blueprint. The materialized athanor lives at `specs/athanors/<name>/`. Agents read that, not this.
