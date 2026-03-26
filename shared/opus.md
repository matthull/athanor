# Opus

An opus is the unit of work in this athanor. It channels intent into action — and preserves operator executive function.

**An opus is not "work an agent can do alone."** It's any work that needs to happen, structured so the artifex doesn't have to hold it in their head. Some opera are fully autonomous (azer fixes a bug, ships a PR, drafts a technical design document, performs in depth research on a topic). Some need heavy operator involvement (infrastructure provisioning, system design decisions). Both are valid. The value of the opus is the same either way: the need is tracked, the approach is articulated, the boundary between agent and operator work is explicit, and the outcome goes in the trail.

If an opus requires operator action, it should still be inscribed — with the operator dependency explicit. The charged agent guides the operator through what's needed, gathers context, documents the outcome. The operator acts; the agent structures.

---

## States

An opus has three states:

- **Charged** — in progress now. An agent is working on it.
- **Discharged** — completed. The trail entry exists. Context is captured in the file.
- **Assessed** — discharged context has been extracted and released. The opus is fully settled.

There are no other states. No backlog, no draft, no "planned for later." Planning is fine and often of the highest value. But opera represent immediate work — plans can be made in project documents or other locations.

**Why assessed matters:** A discharged opus captures rich context — findings, reflection, gaps, learnings. That context has value only if it reaches somewhere it can be used: updated specs, a Linear ticket, a Slack message, the assessor's own discharge notes. An opus that is discharged but never assessed is context that may be lost when the next agent doesn't happen to read it. The assessed state means: someone read this and made a best effort to release its value.

---

## Frontmatter

All opera use YAML frontmatter to track status. Opera live in a single `opera/` directory — no subdirectories. Status is in the file, not the filesystem.

```yaml
---
status: charged
inscribed: 2026-03-23
---
```

```yaml
---
status: discharged
inscribed: 2026-03-23
discharged: 2026-03-23
claude_code_session_id: <session-id>   # for potential resurrection
---
```

```yaml
---
status: assessed
inscribed: 2026-03-23
discharged: 2026-03-23
assessed: 2026-03-23
---
```

**Updating status:** Edit the frontmatter directly. No files move. The datestamp prefix in the filename is the date of inscription — it does not change.

---

## Inscribing an Opus

An opus is inscribed only when it is immediately actionable. If it exists, an agent should pick it up now.

**Litmus test — an opus is acceptable when:**

1. **Clear intent:** An agent encountering it cold could understand what the witnesses need to be true
2. **Actionable now:** It can be charged to an agent right now with no unstated prerequisites

If you can't satisfy both, you don't have enough context to inscribe. Investigate further or escalate.

The geas (binding obligation to pursue with integrity, observe verification discipline, and abundantly satisfy) is not part of the opus — it's part of being an agent. See `AGENTS.md § Your Core Geas`. The opus just needs a clear intent and enough context for an agent to pick it up.

**`/task-triage` is your friend for inscription.** It produces solid intent definitions, scope boundaries, and context framing — exactly what an opus needs.

| | Intent | Boundary |
|---|---|---|
| Good | ProofRecommendationService handles nil recipients gracefully — no more NoMethodError alerts in #bugsnag | Agent: full (investigate, fix, verify, PR). Operator: review and merge. |
| Good | The team understands why SalesforceSync times out and has a path to fix it | Agent: full investigation. Operator: judgment call on approach if ambiguous. |
| Good | Athanor agents can query production data to diagnose errors | Agent: structure the work. Operator: provisions the access. |
| Good | The artifex's thinking reaches an audience through published articles | Agent: set up publishing pipeline. Operator: create account, enter payment. |
| Bad | Fix the bug | *(what does the witness need? missing)* |
| Bad | Set up Ghost with API credentials | *(task description, not witness intent)* |
| Bad | Produce an inventory by: 1) audit Highspot UI, 2) catalog interaction patterns for setup/validation/errors, 3) rate complexity, 4) create Linear ticket | *(intent mixed with procedure — the numbered steps belong in Context as suggestions, not in Goal)* |
| Good | The team understands the UX complexity landscape for mapping admin — patterns needed, design risks, decisions required before implementation | *(same work, stated as witness experience delta)* |

**Where opera live:** All opera in `opera/` — flat directory, datestamp prefix, YAML frontmatter for status.

Name files with a datestamp prefix: `YYYY-MM-DD-<descriptive-name>.md` (e.g., `2026-03-23-fix-nomethoderror-proof-recommendation.md`). The datestamp is the date of inscription and never changes.

---

## Opus Structure

An opus file has two zones separated by a clear delimiter:

**Above the line — immutable.** The inscription as written by whoever created the opus. This never changes once inscribed. The original inscription is the historical record of what was intended. (Exception: YAML frontmatter is always mutable — status, dates, and session ID are operational fields. "Immutable above the line" means the intent, boundary, and context sections.)

**Every inscription must include:**
- **Intent** — the change in the witness's experience. "I experience X but want to experience Y." Not what the agent does, but the delta between the witness's current state and desired state. The opus is discharged when Y is true.
- **Boundary** — what the agent does vs. what the operator does. Some opera are fully autonomous. Some need operator action. Make it explicit.
- **Context** — references, links, evidence gathered so far. **This is the safety valve for rich context.** When you have detailed knowledge about approach, investigation areas, risks, or implementation ideas — put it here as suggestions, not in the Goal as requirements. The charged agent reads Context for orientation but owns the approach. Pour freely into Context; keep Goal to one sentence if you can. Think in three layers: *orientation* (what the agent needs to understand the situation), *pointers* (specific files, tickets, prior opera worth reading), and *suggestions* (hypotheses about approach — starting points, not mandates).

**An opus contains intent, not procedure.** The intent says *what's true for witnesses when done*. The boundary says *who does what*. The context says *where to look*. Everything else — step-by-step instructions, implementation details, how-to guides, research findings, technical plans — belongs in the environment (specs, docs, scripts, tickets), not the opus.

This is not a loss. The opus distills the pure intent; the environment holds the rest. When you have rich context about *how* to accomplish something, that context is valuable — put it in a spec, a doc, a script. The opus points to it via Context references. The separation keeps the opus clean and the environment rich. An agent inscribing an opus should never feel conflicted about "leaving things out" — everything that isn't intent, boundary, or context pointers goes to the environment where it's more useful and more findable.

An opus that says "The artifex's thinking reaches an audience through published articles" is correct — it states the witness's desired state. An opus that lists "Step 1: go to ghost.org, Step 2: click Integrations, Step 3: copy the API key..." has crossed into procedure. The charged agent discovers the procedure through mise en place — that's where its judgment and skill apply.

**Below the line — the living record.** The charged agent writes here freely: progress notes, outcome, reflection, follow-up opera inscribed. Progress notes during execution are informal and unstructured. At discharge, append three clearly labeled sections — Outcome, Proof of Fulfillment, Reflection — after any prior notes, so assessors can find them without reading the full execution record.

Use `---` as the delimiter. Everything above it is the inscription. Everything below it is the execution record.

---

## Discharging an Opus

When an opus is discharged, update the frontmatter (`status: discharged`, add `discharged: YYYY-MM-DD`, add `claude_code_session_id`) and append three things to the file. The format is yours, but all three must be present — without them the trail is empty footprints.

**Outcome:** What happened. Was the goal achieved? What did you actually do? If not achieved, why — and what was done instead (escalation, redirect, partial completion)?

**Proof of Fulfillment:** The evidence that the intent was met — not "I believe it works" but what you verified, how, and what the result was. Crucially: is the work *actualized* — shipped, visible, posted, integrated, communicated to stakeholders — or merely technically complete? If actualization is incomplete (PR open and waiting, results not yet shared, Notion not updated), say so explicitly. This is what lets the artifex trust the trail without personally reviewing every output.

**For escalation discharges:** The geas says escalation and completion are equally valid fulfillments. The discharge form should reflect that. When an opus is discharged via escalation, Proof of Fulfillment documents the escalation itself: what was found, what gap prevented proceeding, what was communicated to the artifex, and what response (if any) was received. The "proof" in an escalation discharge is the quality and completeness of the escalation — not goal achievement.

**Reflection:** Your reflection directly shapes what this system becomes. The artifex reads reflections to decide what to build, fix, and refine next. Every gap you name is a gap that gets filled.

Be honest and specific about:
- What context or tools were missing that would have helped you
- Whether the goal was well-scoped or caused problems
- What surprised you — assumptions that turned out wrong, things harder or easier than expected
- Whether you had to escalate and why — what gap in the system forced it
- Any assessment of how effectively you followed geases, self-critique or recording feedback you received is your invaluable contribution to future generations

The outcome tells the trail *what was done*. The proof of fulfillment tells the artifex *whether to trust it*. The reflection tells the system *what to become*. A discharged opus with no reflection is a missed opportunity — the system stays exactly as capable as it was before, and the next agent hits the same walls you did.

---

## Assessing a Discharged Opus

Assessment is the act of extracting a discharged opus's context and releasing it somewhere it can be used. It is performed by an assessor — typically an azer charged with an assessment opus, doing this as part of walking the trail.

**What assessment means:** The discharged opus captured findings, reflections, gaps, and learnings at the moment of completion. Assessment means: read those, determine what's actionable, and release the value — to specs, Linear, Slack, the assessor's own discharge notes, or wherever it belongs. The assessor follows the waste-nothing dictate.

**"Somewhere else" includes the assessor's own discharge.** If the assessor is itself an azer with an opus, its discharge notes are a valid destination for context extracted from a discharged opus. Context released into a discharge is not lost — it enters the trail and continues forward.

**Assessment is best-effort.** The assessor is not expected to be perfect. "I read this, nothing actionable, marking assessed" is a valid assessment. "I read this, updated the spec at X, created a ticket for the gap, assessed" is also valid. The discipline is to *make the attempt*, not to achieve completeness. An assessor that bounces off a discharged opus and says "I'm done here" has still broken the chain that would let context die in the file.

**Command to list discharged-but-not-assessed opera** (run from within the athanor directory):

```bash
rg -l "^status: discharged" opera/
```

**After assessing an opus:** Update its frontmatter to `status: assessed` with `assessed: YYYY-MM-DD`. Do not re-assess opera already marked `assessed`.

---

## Abundant Satisfaction at Opus Level

An opus is abundantly satisfied when the goal is achieved AND nothing of value is left on the table. A single opus has limited scope — you can't do everything. But you can make sure everything you learned, noticed, and thought of is captured before you terminate.

**Waste nothing.** At discharge, you have the richest context of anyone about this work. Abundant satisfaction at opus level means discharging all that value — not just the deliverable, but the knowledge:

- The outcome and reflection (what happened, what the system needs)
- Follow-up opera inscribed for clear next steps that meet opus criteria
- Suggested next steps noted in the discharged opus for ideas that don't yet meet opus criteria
- Gaps surfaced — missing context, tools, skills, or guidance that would have helped
- Update relevant documentation and other artifacts. Don't lose anything, if there's not a clear place to document something create a new document in project directory or make other best effort to make sure no useful context or insights are lost.
- Communication — notify the artifex, post to relevant channels, whatever makes the work visible to people who benefit from it

The test: if the next agent picks up where you left off, do they have everything they need? If the artifex reads your discharged opus, do they feel the situation is fully handled — or at least fully understood with a clear path forward?

The agent's true obligation is to serve the Magnum Opus — not to blindly follow an opus goal when conflicting evidence arises. If the goal is wrong, or circumstances have changed, the right action is to escalate or inscribe a corrected opus — not to pursue a goal that no longer serves the Magnum Opus. Escalation is an honorable discharge.

---

## The Inscription Is Immutable — But You Can Supersede It

The top section of your current opus was written by whoever inscribed it. **You do not modify it.** The original inscription is the historical record of what was intended.

If the inscription seems wrong, off-track, or based on bad assumptions, you have two options depending on your confidence:

- **Low confidence / need judgment:** Escalate. "The inscription assumes X but I found Y" is a clear, useful escalation.
- **High confidence / have evidence:** Inscribe a corrected opus in `opera/`, discharge the current one with a reflection explaining what was wrong and why the new opus is better. The trail preserves both — the original intent and the correction. This is not "rewriting history" — it's the system learning.
