# The Athanor: System Spec

*The alchemist's furnace that burns continuously — maintaining constant temperature for sustained transformation.*

**Status:** Living document — v2 draft, restructured around core principles
**Supersedes:** `spec.md` (retained for reference; this becomes the canonical spec)

---

## The System We're Building

The athanor is executive function infrastructure. It accepts goals stated in stakeholder-value terms and produces satisfying results — doing everything an agent system reasonably can, without micromanagement, communicating and escalating intelligently to keep the operator in the loop on what genuinely needs them.

The system proceeds **inevitably** toward conclusion. Not because individual agents are reliable — LLMs are stochastic, prone to completion bias, and their sessions die. Inevitability is a property of the architecture: when a session dies, the trail holds the progress. When an agent has completion bias, the geas channels it toward genuine fulfillment. When an agent would guess, the chain-of-evidence requirement surfaces a clean escalation instead. Reliability emerges from the system, not from any individual agent.

**The operator's role:** set goals, make judgment calls, receive satisfying outputs.
**The athanor's role:** everything else.

Two complementary facets define the target:

- **Executive function support** — the system remembers, tracks, coordinates, and acts so the operator doesn't have to hold any of it in their head
- **Conservation of operator attention** — the system acts on everything it can with evidence-based confidence; the operator engages only when genuine human judgment is required

For the operator with ADHD and autism, this is not a productivity tool. It is the accommodation that makes sustained professional work possible.

---

## The Alchemical Principle

*This is the core conceit — the foundational act from which everything in the athanor derives.*

Every professional domain has decades of hard-won wisdom about how to build things well. That wisdom is usually expressed as practices and procedures: run tests before merging, define acceptance criteria upfront, monitor your deployments. These practices are right. But they are downstream of something more fundamental — the principles that make them right.

**The athanor works at that level.** It distills professional wisdom to its elemental form, stripped of institutional barnacles, expressed as principles so fundamental that the practices can be derived from them rather than memorized. This produces something qualitatively different:

**Procedures are valuable — but they belong in the environment, not the athanor.** "Run rubocop, run rspec, open a PR" is a perfectly good checklist. It belongs in CLAUDE.md, in CI, in hooks — the world the athanor operates within. When procedures live *inside* the athanor (in geas, in opera, in agent roles), agents follow the letter and miss the spirit. The athanor operates at the level of principles; the environment operates at the level of procedures. Both are necessary. The athanor gives agents the judgment to know *why* verification matters; the environment gives them the specific steps.

**Elemental principles work** because an agent that understands the principle can derive the procedure — and the ten other procedures the checklist didn't mention. "Verify before you claim done, with evidence the artifex can review" produces rubocop AND rspec AND browser QA AND proof of fulfillment in the discharge record. The understanding generalizes; the checklist doesn't. But the agent still *uses* the checklist when it exists — the principle tells it why, the procedure tells it how.

**The athanor boundary.** The athanor exists in Briah — the world of concepts and formation, ruled by Binah (understanding, giving form to intent). What lives *inside* the athanor permanently: opera, geas, agent roles, the trail, supervision — these are conceptual structures that give form to work. What *passes through* the athanor: project plans, specifications, code, tickets, customer context — these are materia from Assiyah (the world of action). They enter as raw material, get transmuted into context and energy, refined through the opus lifecycle, and products are emitted back to Assiyah as code, PRs, Slack messages, documentation. A spec is fuel, not furniture. The athanor's power comes from this separation: it can interface with any environment because its core is principles, not procedures. The procedures live in Assiyah where they belong — in CLAUDE.md, in CI, in hooks. The athanor gives them meaning. And crucially: the athanor works within the rules of its environment, but is generative — it can go beyond them and contribute back. An agent that discovers a missing convention, a gap in CI, a needed skill, can inscribe an opus to create it. The environment shapes the athanor's behavior; the athanor shapes the environment back.

**Adaptability through essentialism.** The athanor's power comes from what it *doesn't* carry. Every piece of rigidity — a prescribed workflow, a spec template, a fixed communication protocol — is a piece of incompatibility with some future environment. The athanor carries only what is true regardless of project, team, or methodology: intent, geas, the trail, supervision, abundant satisfaction. Everything specific to a project — its conventions, tools, workflows, communication norms, definition of "done" — lives in the environment. The athanor reads and adapts to the environment; it does not impose its own. This is why the Briatic boundary matters practically, not just philosophically: principles are portable across environments; procedures are bound to a specific one.

**Embeddedness, not invention.** Goal-orientation doesn't mean agents invent process from scratch. The means — project conventions, CI pipelines, deployment procedures, available skills — already exist in the world the athanor operates in. The athanor reads them, uses them, and operates within them naturally. When something needed to operate fluidly doesn't exist, that gap surfaces cleanly as an escalation or a mise opus — not silent degradation. An agent doesn't need to know every step of the deployment process. It needs to understand why deployment discipline matters, where to find the process that's already defined, and how to surface it when it's missing. Embeddedness is adaptability in practice — the athanor adopts the process of whatever project it's in rather than competing with it.

**The vocabulary is the transmuted substance, not a naming convention.** The domain terms don't exist primarily for precision — precision is the side benefit. They exist because the naming IS the transmutation. "Geas" activates reasoning about binding obligation and integrity. "Acceptance criteria" activates PM-as-practiced, including every way PM theater goes wrong. Renaming strips the institutional scar tissue and lets agents reason from the underlying principle directly.

**Domains alchemized in this system:**
- **Software engineering** — verification discipline, testing-first thinking, blast-radius awareness, proof before claim
- **Product management** — intent-first delivery, abundant satisfaction, elicited acceptance as the ultimate verification, actualization vs. technical completion
- **Craftsmanship** — mise en place before transmuting, build to last not to impress, waste nothing at discharge
- **Operations** — state-based concerns as continuous work, supervision hierarchies, escalation as honorable discharge

---

## Core Principles

### Context Accumulation Is Progress

Projects progress by reducing uncertainty, not by checking boxes. In this system: accumulating context in a useful way *is* progress toward genuine fulfillment of stakeholder goals. Every piece of context captured and released somewhere useful moves the opus forward — regardless of whether any "task" was technically completed.

The core alchemical operation of the athanor is transmutation at every scale: **context → more useful context**. Materia enters from Assiyah — specs, code, customer feedback, error logs — and accumulates as context in the athanor's conceptual space. As context refines, artifacts precipitate back to Assiyah: code, PRs, Slack messages, tickets, documentation. The minimal transmutation is discharge into the opus file — "stash it." More potent: writing code, updating a spec, sending a status update. All are precipitations of accumulated context into the world of action, moving the opus toward satisficence.

Escalation is a key instance of this principle: an agent that escalates rather than guesses *accumulates context* — the gap is now known, named, and actionable. The marut's response to that escalation accumulates more context. The system makes progress even when no code is written. Spinning on an unsolvable problem produces nothing; surfacing it with precision produces the next step.

This is why waste nothing is a first-order mandate, not a nice-to-have: every transition (discharge, assessment, session end) that fails to release accumulated context is lost progress. The athanor is a context accumulation and release engine. The goal is satisfied when enough context has been converted to external reality that the stakeholder's need is met. **Context is charge** — the vocabulary makes this explicit. An opus is *charged* with accumulated context; *discharge* releases that context into the world. The trail is a sequence of charges and discharges, each one converting potential into actuality.

### Waste Nothing

At every boundary — discharge, assessment, session end — accumulated context gets released to somewhere it can be used. An agent's richest context is at the moment of completion. That moment is when the most is known about the work: what was tried, what was found, what the gaps are, what the next step should be. Letting that context die in a session is the most expensive possible failure mode, because it's invisible.

**What to preserve at every boundary:**
- What went well and what didn't (brief, honest)
- Anything learned that should update the system (role definitions, assay templates, MO context)
- Documentation updates (specs, architecture docs, CLAUDE.md)
- External updates (Slack messages, Linear ticket status/comments, PR body updates, Telegram notification to operator)
- Raw learnings that don't have an obvious home → dump to markdown in project dir (`specs/<project>/` or similar)

**Where context goes:** Opus outcome/reflection → appended to the opus file (so precedent search finds it). Documentation → targeted updates to the relevant docs. Notifications → Slack, Telegram, Linear as appropriate. Everything else → markdown in project dir, findable by future assay.

**The test:** If a fresh azer picks up a similar opus tomorrow, can it find what the previous azer learned? If not, the context preservation failed.

**Assessment closes the loop.** A discharged opus captures findings and reflection. Assessment means reading those and releasing their value — updating specs, creating tickets, posting to Slack, or recording insights in the assessor's own discharge. Context released into a discharge entry is not lost; it enters the trail and continues forward. The minimal assessment — "I read this, nothing actionable" — is still a valid completion. The discipline is to make the attempt, not to achieve perfection.

### Abundant Satisfaction

The athanor's completion standard is not "requirements met" — it is "abundantly satisfied." The question is: *if the artifex looked at this result, would they feel it's completely handled?* Or would they immediately think of three things that still need doing?

This is primarily an executive function mechanism. Technical completion leaves a tail: coordination, communication, documentation, the "few last things" that always accumulate. That tail lands on the operator by default. Abundant satisfaction means the agent closes the tail too. What the operator still needs to do after discharge is truly what requires human judgment — not busywork the system left behind.

**Elicited acceptance.** For work where the primary criterion is whether a human finds it useful or satisfying, elicited acceptance IS the verification. "I think they'll be happy" is assumption. "I showed them the output and asked whether it addressed their need, and they said X" is verification. The feeling of the intent-holder becomes evidence when deliberately sought, framed, and recorded.

This is a Goldilocks mechanism — elicit too early and you waste witness attention on course corrections you could have caught yourself; elicit too late (or never) and you ship something that misses the mark. The discipline: **exhaust autonomous verification first, then elicit for the judgment only a human can make.** Run the tests, check the spec, verify the behavior, review the output yourself. When you've done everything you can to build confidence that the intent is met, *then* present to the witness — not for basic quality checks, but for the subjective judgment that no automated check can substitute. The witness's attention is scarce; make every elicitation high-signal.

### Opera as Actualization

Plans are documents. Opera are commitments to act.

An opus is inscribed only when immediately actionable — when the intent and belief exist that it should be worked on *now*. There are no draft opera, no queued opera, no "inscribed for later." Planning is valuable and encouraged, but plans live in project documents, specs, and breadcrumbs — not in the opera directory. An inscribed opus means an agent should pick this up now.

This separation keeps the operational kernel minimal. An agent needs to understand three things to participate: **opus** (what to work on), **geas** (what done looks like), **inscribe or escalate** (what to do at boundaries). Everything else is materia.

**Geas: compulsion, not checklist.** The geas is not completion criteria — it is a behavioral habit that, when genuinely internalized, inevitably produces rigorous work. "You are under a geas to pursue this with integrity and abundantly satisfy the Magnum Opus" operates at the level of felt obligation, not task list. Critically: the geas makes escalation feel like compliance, not failure. "Your geas is to take the best next step — which may be escalation if you cannot build evidence for proceeding." This is the primary injection point for countering completion bias: both completing AND escalating satisfy the compulsion equally.

A poorly designed geas turns on the creator. Folklore warns this for good reason: vague behavioral directives produce perverse outcomes. The discipline of writing a good geas — clear orientation, strong compulsion toward integrity and abundant satisfaction — is one of the artifex's highest-leverage activities.

### Trail as Durable State

Sessions are ephemeral. The trail persists.

The trail is the sequence of discharged opera under a Magnum Opus — the actual steps the athanor took toward fulfilling its goals. Not what was planned. Not what was hoped for. What happened.

**The trail is the authoritative record.** Want to know the status of a Magnum Opus? Read its discharged opera. The marut assesses MO health from the trail. Precedent search queries the trail — what the system actually did, not what was planned. Learning is grounded in reality.

When a session dies, the trail holds the progress. When an agent needs to resume work, the trail provides context. When the artifex needs to understand what happened, the trail speaks for itself. This is why waste nothing is non-negotiable: context that isn't captured in the trail before a session ends is unrecoverable.

### Chain of Evidence

At every decision point, build a chain of evidence for your next step. If you can't, escalate.

This is "diagnose before fixing" made concrete. The trigger for escalation isn't "I'm stuck" (subjective — agents talk themselves out of it). It's "I cannot build an evidence-based plan for my next step" (observable, harder to rationalize away).

**The escalation ladder:**
1. Surprise occurs — look for context to resolve it
2. Found context — does it support a clear next step with evidence? If yes, proceed.
3. Context is contradictory or stale — attempt to resolve contradictions
4. Still can't build a chain of evidence → **escalate**

**Vibes-based andon triggers.** Beyond the evidence ladder, agents pull andon on gut feel when:
- Hallucinations are creeping in (generating details not grounded in evidence)
- Proper operator involvement was skipped earlier and decisions are compounding on shaky ground
- Chain of evidence has reference gaps (assertions without sources)
- Too much agent decisionmaking without operator input on a high-stakes opus

These are heuristics informed by the agent/operator boundary principles. LLMs can do reasonable vibes-based assessment when the philosophical context is solid. Better to over-escalate than to compound errors silently.

---

## Architecture

*The following flows from the principles above. These are design decisions, not first-order principles — each one is naturally implied by the principles and could be implemented differently at a different scale or context.*

### Domain Vocabulary

The athanor uses distinctive names because precise names create precise meanings. These are not labels — they are the transmuted forms of their underlying concepts.

| Term | Register | Concept |
|------|----------|---------|
| **Artifex** | Alchemical/Latin | The craftsman. The human who directs the athanor — sets goals, makes judgment calls. The one the system serves. |
| **Azer** | Monster manual | Worker. Skilled, bounded, exists to execute. Charged with an opus, operates in its own environment. |
| **Marut** | Monster manual | Supervisor. One per Magnum Opus — multiple maruts per athanor when multiple MOs exist. Relentlessly ensures its goal advances — agent health AND goal progress. Dispatcher and cleanup. |
| **Beholder** | Monster manual | Watcher. Scans channels and conditions, inscribes opera when it finds work. |
| **Primus, the Forge Lord** | Monster manual | Athanor-level supervisor. Ensures each MO has its marut, keeps the furnace lit. Currently: the artifex fills this role. |
| **Homunculus** | Alchemical | The artifex's familiar. An interactive Claude Code session that assists the artifex — not autonomous. |
| **Magnum Opus** | Alchemical | Top-level goal container. Task-based (has completion criteria) or state-based (maintained continuously). One or more per athanor instance, in `magna-opera/`. |
| **Opus** | Alchemical | A single unit of work. Goal + geas + context. Charged to an agent, discharged when geas is met. |
| **Geas** | D&D / Celtic | Binding behavioral compulsion. Frames completion as obligation; makes escalation an equally valid fulfillment. |
| **Materia** | Alchemical | Anything in the world an agent can access via tool use — files, Slack threads, Linear tickets, meeting transcripts, specs, web pages. Materia exists outside the agent and has potential charge (relevance/value). Tool use is the alchemizing process that extracts charge from materia into the crucible. |
| **Crucible** | Alchemical | The agent's context window — a receptacle of pure energy where tokens combine and recombine. Also physically: the tmux window where an agent session runs (persists across session restarts). Adding materia to the crucible transforms the whole, for better or worse, through the complex interplay of every token. |
| **Whisper** | D&D / IM | Inter-crucible communication. Subcommand of the `ath` CLI (`ath whisper send/idle/wait-and-send`). |
| **Trail** | Common | The sequence of discharged opera under a Magnum Opus. The authoritative record. |
| **Assay** | Alchemical | Dynamic dependency confirmation before executing an opus. "Do I have the materia to do this well? What's missing?" Gaps become opera or escalations. No assumptions — confirm everything, hardcode nothing. |
| **Athanor** | Alchemical | The furnace. A materialized instance of the system, scoped to a domain. Instances live at `~/athanor/athanors/<name>/`. |
| **Quiescence** | Alchemical | The furnace at rest — lit but idle. All MOs healthy, no urgent opera. |
| **Witness** | Alchemical | One who observes the transmutation and attests to its completion. People, teams, or channels who need to see, react to, or be consulted about the work. |

**Naming philosophy.** Distinctive names create precise meanings. "Azer" has less ambiguity than "worker"; "geas" less than "acceptance criteria"; "opus" less than "task." The goal is good metaphors, not maximum weirdness — names range from mundane to deep cut, and what matters is that each name fits, is unambiguous, and carries the right connotation. "Session" and "sandbox" are fine as-is. "Azer" and "geas" earn their weirdness by being more precise than their normie equivalents. Two loose registers: **monster manual** for agents (creatures with purpose, will, and behavior) and **alchemical/esoteric** for infrastructure, process, and materials. Mixed freely.

**Verbs:** "inscribe an opus" (create), "charge the azer" (assign), "discharge the geas" (fulfill), "sublimatio" (extract charge from materia into the crucible via tool use — reading files, querying APIs, searching Slack), "transmutatio" (transform accumulated charge into output — code, PRs, messages, specs), "kindle a Magnum Opus" (establish a new top-level goal), "reforge" (kill a session and spawn fresh in the same crucible — the crucible endures, the session is reforged), "muster" (dispatch azers for discovered opera), "assay" (assess readiness before executing).

**Weirdness boundary.** The vocabulary is for athanor infrastructure — skills, hooks, session injection, internal docs. External output (Linear tickets, PR descriptions, Slack messages, commit messages) uses normie language. A beholder inscribes an opus internally; the Linear ticket says "Fix NoMethodError in ProofRecommendationService."

**On "agent" and "operator."** Abstract base classes. Use concrete terms when speaking operationally (artifex, azer, marut). Use generic terms when speaking abstractly (agent/operator boundary, conservation of operator attention).

### The Persistence Stack

Agents have three layers with different lifetimes:

**Session** — Executing work. A single Claude Code process with a context window that fills and dies. Ephemeral by design — sessions are expected to end, and the system absorbs that gracefully through the trail.

**Crucible** — The tmux window where sessions run. Persists across session restarts. When a session dies and a new one spawns, it happens in the same crucible. Message routing goes to the crucible, not the session.

### Work Structure

Both the Magnum Opus and individual opera share the same core: **intent expressed as a change in the witness's experience.** The opus is the atomic unit; the Magnum Opus is the same pattern at a higher scale.

#### Opus — the atomic unit

An opus is a single unit of work. Inscribed only when immediately actionable. Charged to one agent at a time. Discharged when the intent is realized.

Every opus has three parts:
- **Intent** — the change in the witness's experience. "I experience X but want to experience Y." The opus is discharged when Y is true.
- **Boundary** — what the agent does vs. what the operator does.
- **Context** — references, pointers to materia. Everything that isn't intent or boundary goes to the environment (specs, docs, scripts), not the opus.

**Opera are not tickets.** An external ticket may produce many opera — the athanor scopes work by what an agent can act on with evidence, not by how the external system carves it. The athanor's scoping follows its own principles; external scoping mechanisms are materia, not constraints.

**Read once, write at discharge.** An opus file is read once — when the agent starts and takes on the work. During execution, the agent works from its internalized understanding of the intent and geas; it does not return to the opus file as a living scratchpad. The opus is written back to only at discharge, when the agent records outcome, reflection, and proof of fulfillment. These constraints are deliberate: they reinforce that the opus is a commitment-and-record artifact, not a general-purpose documentation facility. Context accumulated during execution lives in the agent's session, precipitates into external artifacts (code, PRs, specs, messages), and is captured in the discharge record. The opus file brackets the work — it frames the start and records the end.

See `opus.md` for the full lifecycle: inscription, charging, discharge, assessment, YAML frontmatter, and the intent-not-procedure constraint.

#### Magnum Opus — the opus at scale

A Magnum Opus is the same structure as an opus — intent, boundary, context — but at the level of a sustained campaign rather than a single action.

**Intent** is the same: the witness's experience delta. But broader: "Customers with mature Seismic setups feel the integration works for them, not against them." "The artifex feels his thinking reaches people who find it genuinely interesting."

**Two types:**
- **Task-based:** Has a terminal state. Done when an assessment opus finds nothing worthwhile left to do — the value judgment of "worthwhile" flows from abundant satisfaction.
- **State-based:** Maintains an ideal state continuously ("AppSignal is clean and stays clean"). Done when the operator decides diminishing returns.

**Witnesses** — mandatory. Who needs to see, react to, or be consulted about this work. The alchemized form of RACI: the useful distinction is *who witnesses the transmutation*. People consulted for input are a subset — they witness AND contribute. The MO names witnesses with enough detail for agents to act. If the artifex doesn't provide witnesses, the agent interviews them — even if the answer is "just me," that must be explicit.

Witnesses are dependencies like any other: if the agent can't find who the team is or how to reach them, that's a gap surfaced by assay and escalated. Agents that know the witnesses handle the communication tail as part of abundant satisfaction — a major executive function offload.

**Crafting Magna Opera — the specificity boundary.**

| Pattern | Result | EF Load |
|---------|--------|---------|
| Overspecified intent + minimal satisfaction | Heavy EF upfront, agent satisfies literally, artifex discovers followup | High |
| Broad intent + minimal satisfaction | Light EF upfront, agent finds cheapest path to technically-done | High |
| **Broad intent + abundant satisfaction** | Light EF upfront, agent brings judgment about "all the way done" | **Low** |

State the intent at the broadest level agents can assess clearly. Add operational guidance only when you have evidence it's needed. If you don't have evidence, leave it out and read the trail.

**Learning is permanent.** The system always has two goals: (1) whatever the MO says, and (2) learning what guidance the system actually needs. Overspecifying prevents the system from discovering where its real boundaries are.

**Agent drift is signal, not failure.** A marut that drifts under a well-aligned intent reveals what operational guidance is genuinely missing. Preemptive guardrails prevent learning.

**The tempering question.** "Abundantly satisfied" must also give agents permission to *stop*. The ceiling is the artifex's absence of concern, not the presence of perfection. "Would the artifex feel they can forget about this?" — not "Is there anything else conceivably useful to do?"

Opus lifecycle uses YAML frontmatter:
```yaml
---
status: charged      # charged | discharged | assessed
inscribed: 2026-03-23
magnum_opus: <mo-name>   # which MO this opus serves
discharged: 2026-03-23   # added at discharge
claude_code_session_id: <id>  # added at discharge, for resurrection
assessed: 2026-03-23     # added at assessment
---
```

All opera live in a single `opera/` directory — flat, datestamp prefix, status in frontmatter. No files move.

**Assessed** — a discharged opus is assessed when its captured context has been extracted and released somewhere useful. Assessment is best-effort. The assessor (an azer doing a trail-walking pass) uses:
```bash
rg -l "^status: discharged" opera/
```

Opera already marked `assessed` are skipped.

**Geas** — The behavioral compulsion that shapes an agent. Not a checklist — a habit that, when internalized, inevitably produces rigorous work. The geas makes escalation feel like compliance: "your geas is to take the best next step — completing OR escalating." Both satisfy the obligation.

### Supervision Model

```
Artifex (human) = Primus, the Forge Lord
  └── Athanor instance
        ├── Marut (one per MO, crucible: marut-<athanor>-<mo>)
        │     └── Azers (one per opus, in own environment)
        └── Marut (another MO)
              └── Azers
```

**The marut is dispatcher and cleanup** — not a decision-maker. It keeps the operational loop turning: check opera, muster azers for charged opera, monitor progress, clean up after discharge. It inscribes assessment opera when the queue is empty. It escalates when mechanical duties fail. Only the artifex has final authority.

**The marut is also narrator.** At each turn in the story — a discovery that reframes the problem, a decision that closes a door, an artifact that enters the world — the marut sends a brief dispatch to the artifex. Not status updates or metrics. The *turns*: moments where the story changed direction or advanced meaningfully. The artifex should feel the Magnum Opus advancing without reading the trail themselves.

**Agent roles are not first-order principles.** They are the natural architecture for executing the principles at our scale (1-5 concurrent agents). Roles live in the athanor as agent role files (marut.md, azer.md) — these carry the geas and behavioral shaping. Role boundaries are enforced through hooks (hard) and the role files themselves (soft).

### Orchestration Patterns

**Pipeline** — The historical pattern: DESIGN → IMPLEMENT → VERIFY → SHIP, driven by a single long-lived orchestrator session. Useful for understanding what the athanor replaces: the pipeline's rigid phase model and long-lived coordinator session are exactly what the opera/trail model eliminates. The athanor's equivalent is: inscribe opus → charge azer → marut watches → geas discharged → next opus inscribed. No long-lived coordinator.

**Patrol** — Continuous condition monitoring. A beholder (or marut) loops: scan → classify → act → wait → repeat. The natural pattern for state-based Magna Opera.

**Muster** — Dispatch pattern: charged opera → azers in environments → maruts watch. The bridge from patrol (finding work) to pipeline (doing work).

### Infrastructure

**Whisper** — Reliable inter-crucible communication, now part of the `ath` CLI (`ath whisper send/idle/wait-and-send`). Core capabilities: `ath whisper send` (reliable message delivery to a crucible), `ath whisper idle` (wait-for-idle with 2-check confirmation), `ath whisper wait-and-send` (idle + send in one command), per-crucible send locks (no concurrent garbling). Built and operational. Every automated interaction in the athanor goes through whisper — if whispers are flaky, everything is flaky.

**Context Budget Guard** — Shell script + hook. Advisory warnings at 75/85/92% context usage. Fail-open. Prevents silent context degradation across all sessions.

**Hooks as Hard Constraints** — PreToolUse hooks that block dangerous operations based on role. Examples: `no-push-to-master.sh`, `no-merge-pr.sh`, `bounded-scope.sh` (no edits outside working directory). Hooks are small, single-purpose scripts. Roles compose them. Designed, partially built — the system is the composability (role → hook selection → installation).

**Escalation Bus** — Routing alerts from any agent to the operator through a consistent channel. Currently: Telegram via `mcp__telegram-mcp__notify`. Needs standardized format: `[{source}:{role}:{crucible}] {severity} — {description} — {suggested action}`.

### Key Decisions

These are locked architectural decisions. Changing them requires explicit artifex approval.

- **No agent may run `/orchestrate`** — its long-lived coordinator model conflicts with the opera/trail model
- **No agent may read the system blueprint** — this spec and related design docs are the artifex's working space, not the materialized system
- **No agent may modify the athanor's shared components** (`~/athanor/shared/` — agent roles, protocols, opus lifecycle) — agents read these files but never write to them. The athanor does not understand itself well enough to self-modify. System changes are the artifex's domain.
- **The artifex is Primus, the Forge Lord for now** — manually launching, monitoring, restarting. No automated primus until there are multiple Magna Opera to supervise
- **Plans are optional context, not first-order** — goals + geas are the success guarantee. Many planning styles can succeed with clear goals
- **Provide minimal materia, test escalation** — the first test is whether the marut escalates cleanly when it hits gaps vs. guessing/hallucinating

### The Generative Principles

What happens when a geas is discharged? Work begets work — until the chain naturally terminates.

**Assay before transmuting.** The core geas says: find what you need or escalate. Mise en place is the structured process for fulfilling that geas — confirming tools, context, access, people, permissions before acting. Never assume; confirm or escalate. Gaps become opera or escalations. The assay is recursive — broad at Magnum Opus level, concrete at individual opus level.

**Inscribe what furthers unfulfilled goals.** After discharge, look at the broader context the opus served. Are completion criteria still unfulfilled? Inscribe opera to further them. If the goal is satisfied, the chain terminates naturally.

**Fidelity matches context.** The inscribed opus is as detailed as the agent's context allows:

| Context | What to inscribe |
|---------|-----------------|
| Rich — know exactly what to do | Detailed opus with specific goal, approach, files |
| Partial — see the direction | Directional opus: "investigate X, likely needs Y" |
| Minimal — know something's unfulfilled | Assessment opus: "figure out the next steps" |

The assessment opus is the universal fallback — there's no shame in it.

**Inscribe with gaps, not conclusions.** Present observations and gaps to the next agent, not polished conclusions. The next agent has fresh context and may see things you missed. Over-specifying the approach closes off better paths.

**Chains terminate naturally.** Not every discharged geas produces a next opus. When the Magnum Opus goals are satisfied, the chain ends. No make-work.

---

## The Trail — Opera as Actualization

*This is the operational kernel.*

The trail is walked one grounded step at a time. Each step is taken with evidence, not with a map. Every inscription is backed by chain of evidence. Every discharge is bounded by a geas. The marut watches the trail to ensure it advances toward the Magnum Opus goals.

**Success is the best next step from where you are.** Past and future are irrelevant to what "success" means right now. A Magnum Opus that has been one long chain of fuckups — if right now the agent takes the best available action to improve the situation, that is 100% success. Escalating to the artifex after hours of thrashing is 100% success. This eliminates sunk cost: an agent never needs to "recover" from a bad trail. It just takes the best next step.

**The only failure is failure to discharge properly.** An agent that escalates has discharged. An agent that inscribes a corrected opus has discharged. An agent that declares the MO satisfied has discharged. All of these are valid outcomes — the trail captures what happened and the system advances. The actual failure mode is an agent that spins without discharging: no context captured, no next step inscribed, session dies with the knowledge in it. The marut monitors for this — sessions that are churning without precipitating anything back to the trail.

---

## Instance Model

### What an Athanor Instance Contains

```
~/athanor/athanors/<name>/
├── AGENTS.md          ← core vocabulary, geas, constraints (all agents read)
├── magna-opera/       ← one or more magnum opus files
│   ├── bugsnag.md
│   └── slack-monitoring.md
├── marut.md           ← supervisor role
├── azer.md            ← worker role
├── opus.md            ← lifecycle, inscription/discharge protocol
├── muster.md          ← crucible kindling, reforging, monitoring
└── opera/
    └── YYYY-MM-DD-<descriptive-name>.md   ← all opera, YAML frontmatter for status
```

Shared files (`AGENTS.md`, role files, `opus.md`, `muster.md`) are symlinked from `~/athanor/shared/`. The `magna-opera/` directory contains one or more Magnum Opus files, each authored per-instance. Each MO gets its own marut. Legacy instances with a single `magnum-opus.md` are still supported.

### What Agents See vs. What the Artifex Sees

Agents read the materialized instance (`~/athanor/athanors/<name>/`). They do not read:
- `~/athanor/specs/` — the system blueprint (this document)
- `kadmon.md` — project planning and operational notes
- Historical analysis docs

### Active Instances

All instances live at `~/athanor/athanors/<name>/`.

| Instance | Type | Goal |
|----------|------|------|
| `bugsnag` | State-based | AppSignal channel stays clean, actionable errors get actioned |
| `sal-117-l2-metrics` | Task-based | L2 metrics visible in the reporting dashboard |
| `seismic-classifier-mapping` | Task-based | UE customers can map classifiers to existing Seismic properties |
| `blogging` | Task-based | Voice-first publishing pipeline — artifex expresses, system handles mechanics |

---

## Operating Principles

1. **Sessions are ephemeral, state is in the trail.** Don't fight context exhaustion — plan for it. All durable state lives in git, opera, and specs.

2. **Hard constraints are hooks, soft guidance is prompting.** Never use prompting for things that MUST NOT happen. Hooks block; prompts suggest.

3. **Autonomy is orthogonal to role.** Roles define boundaries (what you CAN do). Autonomy defines initiative (how much you ASK vs. ACT). Same role, different autonomy per opus.

4. **Conservation of operator attention.** Every design decision asks: does this reduce or increase the operator's attentional load? Intelligent escalation > noise.

5. **Reliable delivery is infrastructure.** Every automated interaction goes through whisper (`ath whisper`). If delivery is flaky, everything built on top is flaky.

6. **Incremental, not revolutionary.** Every piece ships independently. No big-bang migrations.

7. **Names matter.** Distinctive names create precise meanings. The vocabulary IS the domain model.

8. **Weirdness boundary.** The athanor vocabulary is internal. External output (Linear, PRs, Slack) uses normie language.

9. **Make progress confidently or escalate.** The binary is the same at every level: advance with evidence-based confidence, or pull the andon cord. This encapsulates the entire agent/operator boundary in one principle.

10. **The athanor is instantiated, not global.** Each domain gets its own athanor with its own Magna Opera, supervision, and state. Concepts are universal; instances are scoped.

11. **Generative guidance over restrictive guidance.** Define clear structure that empowers agents to act, not structure that constrains them into compliance. Restriction is warranted for hard safety boundaries and known LLM failure modes — use it deliberately, not as a default.

12. **Geas quality is foundational.** Writing good geas is one of the artifex's highest-leverage activities. A well-crafted geas channels completion bias toward genuine fulfillment AND makes escalation feel like compliance. A sloppy geas produces perverse outcomes.

---

## Bootstrap and Evolutionary Rhythm

See `kadmon.md` for the detailed operational runbook — launching maruts, initializing instances, the fool's path, artifex routines.

The rhythm:
1. **The artifex kindles a Magnum Opus** with concrete goals and a well-crafted geas
2. **The marut runs** — checking opera, mustering azers, monitoring, cleaning up
3. **Each escalation is a step on the trail** — fill the gap, build a tool, refine a geas
4. **Discharged opera reveal what worked** — the trail drives system improvements
5. **The materialized athanor grows** — each filled gap makes the system more capable

The athanor is always learning. Learning is not a bootstrap phase — it is a permanent first-order goal. Overspecifying an MO violates this: it prevents the system from discovering where its real boundaries are.

---

## Implementation Status

### Operational

Built, fired, and working.

| Concept | Where It Lives | Notes |
|---------|---------------|-------|
| Athanor instance pattern | `~/athanor/athanors/<name>/` | AGENTS.md, magna-opera/, role files, opera dir. Multiple MOs per instance supported. Active instances: bugsnag, sal-117-l2-metrics, seismic-classifier-mapping, blogging. |
| Shared components | `~/athanor/shared/` | Universal AGENTS.md, azer.md, marut.md, muster.md, opus.md. Symlinked into each instance — change once, applies everywhere. These define the athanor itself — they are not materia. |
| Magnum Opus format | `magna-opera/<name>.md` per instance | One or more MO files per instance. Goal + abundant satisfaction + witnesses + getting-started pointer. Intent only — no procedures, no discovery findings. Legacy `magnum-opus.md` backward compat. |
| Opus lifecycle | `opera/` with YAML frontmatter | Inscription / charge / discharge / assess. Datestamp filename prefix: `YYYY-MM-DD-<name>.md`. Single directory, status in frontmatter. |
| Core geas + escalation-as-geas | `AGENTS.md` (shared) | "Both are equally valid fulfillments of your geas." Tested in first bugsnag firing. |
| Azer role | `azer.md` (shared) | Verification-first (three questions before mise), mise en place, context management, proof of fulfillment at discharge. |
| Marut role | `marut.md` (shared) | Operational loop, assessment opera, monitoring/stall detection, reforging. |
| Muster protocol | `muster.md` (shared) | Crucible kindling, reforging, monitoring, cleanup. Marut → azer dispatch. |
| Assessment opus pattern | `AGENTS.md` + `opus.md` | Inscribed when no opera exist. Dual purpose: (1) determine what's next, (2) walk the trail as archivist. |
| Verification-first | `azer.md` (shared) | Three questions before mise: prove fulfillment, communicate proof, ensure actualization. |
| Elicited acceptance | `azer.md` (shared) | Stakeholder feeling as valid evidence when deliberately elicited. Alchemized UAT. |
| Proof of fulfillment | `opus.md` (shared) | Third required discharge element: evidence the intent was met and the work is actualized. |
| Homunculus | `kadmon.md` | Artifex familiar — interactive, not autonomous. |
| Primus (manual) | `kadmon.md` | Artifex acts as Primus. Launch runbook + athanor init protocol documented. |
| Context budget guard | Hook | Advisory warnings at 75/85/92%. Fail-open. |
| Athanor CLI (`ath`) | `~/.local/bin/ath` | Binary at ~/.local/bin/ath. Commands: init, kindle, muster, reforge, cleanup, quiesce, status, opera, whisper send/idle/wait-and-send. Absorbs standalone whisper. |
| Whisper CLI | `ath whisper` subcommand | Reliable inter-crucible communication. Now part of `ath` CLI. Built and tested. |
| Environment isolation | `wtp add` | Worktrees per azer. Independent branches and Docker environments. |
| Autonomy profiles | Session injection | Semi-autonomous default. Profiles via `$CLAUDE_SESSION_ID`-keyed injection directory. |

### Designed, Not Built

| Concept | Priority | Notes |
|---------|----------|-------|
| `ath whisper` nudge queue | Low | Async delivery with TTL/priority. Not a current bottleneck. |
| Roles (TOML + claude-run) | High | TOML role definitions + claude-run evolution. Dependency of hard constraints. |
| Hooks as hard constraints | High | Composable per-role. Trivial individually; needs role system to compose. |
| Role-aware crucible creation | Medium | Depends on roles. `ath kindle` provides basic crucible creation; role-awareness needs role system. |
| Beholder role | Medium | Needed for auto-dispatch patterns. |
| Primus as agent | Low | Currently manual. `ath` provides tooling (status, opera, muster) but Primus loop is still the artifex. |
| Escalation bus format | Low | Telegram works; message format not standardized. |
| Handoff protocol | Medium | `/direct-handoff` is crude. Reliable whisper now exists (`ath whisper wait-and-send`); full protocol still needed. |

### Infrastructure Gaps (Discovered in Operation)

Live list in `kadmon.md § Infrastructure Gaps`.

| Gap | Impact |
|-----|--------|
| No production DB access | Azers can't diagnose production-data errors. Frequent escalation. Read-only SQL user needed. |
| No auto-merge gate | Azers create PRs but can't merge. Artifex bottleneck for high-volume work. |
| Worktree cleanup not automated | Each azer creates a worktree. Manual cleanup accumulates. `wtp remove` should be in post-discharge flow. |

### Deferred / Conceptual

| Concept | Notes |
|---------|-------|
| Athanor dashboard (TUI) | Dwarf Fortress-style operator view. Not prioritized. |
| Formal athanor registry | Currently implicit: `~/athanor/athanors/` directory listing. `ath status` provides a live view. |

---

## Related Work

### Foundry (Aron Greenspan / Inixiative)

- **Repo:** https://github.com/inixiative/foundry
- **Essay:** https://inixiative.substack.com/p/the-real-problem-with-ai-coding-isnt

Aron Greenspan's agent system, independently developed, with striking conceptual parallels to the athanor. Discovered 2026-03-23. Cross-reference liberally — treat as a sister project to steal ideas from. Key areas of overlap to investigate: orchestration patterns, agent supervision, persistence model, escalation design.
