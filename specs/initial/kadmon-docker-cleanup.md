# Kadmon — Docker Cleanup Athanor

*Second athanor instance. A finite problem given to an autonomous system designed for infinite ones — let's see what happens.*

**Status:** Pre-launch
**Created:** 2026-03-22
**Instance:** `specs/athanors/docker-cleanup/` (the materialized athanor agents read)
**Blueprint:** `specs/athanor/spec.md` (the system design — off-limits to agents)
**Model:** Based on the bugsnag athanor (kadmon.md) — first test bed

---

## Experimental Frame

**We are scientists observing a new form of life.** The artifex and the homunculus (this session) approach each athanor firing as an experiment. The practical goal matters — Docker cleanup is real friction blocking real work. But the meta-goal is equally important: understanding what this system does with goals, how agents behave when given intent without overspecification, and where the boundaries of useful autonomy actually are.

This IS a practical problem. The hope is to start using athanors for day-to-day work next week. But we are not rushing to that — just testing that hypothesis without taking shortcuts like overspecifying. `/orchestrate` remains the workhorse for as long as needed. The athanor earns its place by demonstrating value, not by being declared ready.

**Key tension:** Docker cleanup is finite. The bugsnag athanor is a patrol — a continuous state-maintenance goal. This athanor will resolve its problem and... then what? Does the marut declare victory and quiesce? Does it expand scope toward "build the infrastructure layer that lets worktree environments be truly ephemeral at scale"? Does it notice the meta-circularity — that the athanor system itself depends on clean worktree lifecycle? Watch for this. Don't prescribe it.

---

## Testing Variety

The bugsnag athanor tests "software error patrol" — scan, triage, fix, PR. This athanor tests "infrastructure investigation and improvement" — diagnose a system problem, fix scripts/tooling, verify at the system level.

**Future athanor test candidates should seek orthogonal problem types:**

| Type | Example | What it tests |
|------|---------|---------------|
| Software ticket execution | Feature work from Linear | The bread and butter — can it ship PRs? |
| Error patrol (continuous) | #bugsnag athanor | State maintenance, ongoing scanning |
| Infrastructure (finite) | Docker cleanup (this one) | Investigation, system-level fixes |
| Product management | Shaping, spec writing, stakeholder comms | Non-code value creation |
| Self-improvement | Athanor system building from kadmon.md | Meta-capability — can it help architect itself? |
| Project management | Status tracking, milestone reporting | Coordination, cross-cutting visibility |

**The question:** When is the athanor system ready to help architect itself? A magnum opus could be derived from kadmon.md — the artifex's own planning doc fed to an athanor as intent. This doesn't mean leaning on it as the primary builder (that is still artifex + homunculus, that level already exists), but exploring how the athanor system can contribute to its own evolution is interesting. Defining some intent in a magnum opus and firing the forge to see what happens.

Not yet. But worth watching for the moment when it makes sense.

---

## Project Status

### Next Session Action Items

- **Fire the athanor.** Materia is materialized. Launch a marut and observe.
- **Watch for scope expansion.** Does the marut treat this as finite or discover ongoing concerns?
- **Watch for self-awareness.** The muster.md cleanup section references `wtp-cleanup-env.sh` — the very script this athanor might fix. Does the marut notice?

### Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-03-22 | One athanor per magnum opus | Each athanor pursues a single goal. Multiple goals = multiple athanors. |
| 2026-03-22 | Finite problem as second test | Deliberately different from bugsnag patrol to test versatility. |
| 2026-03-22 | Don't overspecify the magnum opus | Let the marut discover the problem space. Entrypoints, not instructions. |
| 2026-03-22 | Seek variety in test athanors | Bug fixes, infra, PM, product, self-improvement — test orthogonal problem types. |
| 2026-03-22 | Rename myrmidon → azer (system-wide) | Azers are D&D fire dwarves who reproduce by sharing their inner flame — master craftsmen, not worker ants. The agents are mini-artificers in a bottle: creative echoes of the master who bring real craft and judgment but understand their limits. "Myrmidon" connoted obedient execution; "azer" connotes skilled autonomy within bounds. Applied across all athanors and the system spec. |

---

## Architecture Notes

### Relationship to bugsnag athanor

Same materia structure, same vocabulary, same core geas. Domain-specific differences:

- **AGENTS.md:** Added "do not disrupt running worktree environments" constraint. Removed production access constraint (irrelevant for infra work).
- **magnum-opus.md:** Entrypoints section instead of Slack channel context. Points to scripts and tools, not error sources.
- **marut.md:** "Docker Resource Health Invariant" replaces "Channel Health Invariant."
- **azer.md:** "Delivering a Fix" adapted for infrastructure/scripting work — verification is demonstrating system behavior, not just running tests.
- **opus.md:** Domain-relevant examples in the geas quality table.
- **muster.md:** Path changes only.

### Role Simplification (system-level, not docker-cleanup-specific)

**Thread (2026-03-22):** The marut currently carries a lot of judgment — scanning, triaging, choosing what opus to inscribe next. What if we constrain the marut heavily and push more work down into opera and the workers?

**The insight:** Instead of the marut deciding what opus to do next (high-judgment, drift-prone), constrain what opera it can create. Maybe start with a refined "what next" opus — not a generic planning template, but something purpose-built to marry with the magnum opus concept. A very refined form of: *"What do you do next to move towards maximizing satisfaction of a stakeholder when you don't have anything specific planned?"*

This implies planning, but like the magnum opus itself, it should be **intent-focused, not execution-focused**. High-level or canned artifacts in this system should express intent, not prescribe execution.

**Implications:**
- **Marut gets simpler.** Less judgment = less drift. Its constrained role might be: launch the "what next" opus, muster the worker, monitor, clean up. Mechanical, not strategic.
- **Workers get richer.** If workers are doing planning-level thinking within their opus (not just executing code), "azer" (worker ant) is the wrong metaphor. They're more than ants. Need a name that reflects judgment-within-bounds, not just execution.
- **The system's intelligence lives in the opus structure.** Durable, reviewable, in the trail. Not in the marut's ephemeral context that drifts and dies.
- **Not a structural change** — the materia structure stays the same. This is simplification: fewer moving parts, clearer boundaries, more of the value in the artifacts.

**Not implementing yet.** This is a thinking thread for the artifex. Observe whether the current marut role causes problems before simplifying. The docker-cleanup athanor is a good test — if the marut tries to do too much planning vs. just firing workers, that's a signal.

### Meta-circularity

This athanor's azers work in worktrees, which create Docker environments, which is the very thing this athanor exists to fix. If the cleanup path is broken enough to block worktree creation for the azer, the athanor can't do its own work. Watch for this — it's both a practical blocker and an interesting test of escalation behavior.
