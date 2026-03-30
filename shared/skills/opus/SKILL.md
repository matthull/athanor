---
name: opus
description: Opus lifecycle operations — inscription, witness audit, MO shaping. Each operation uses dialectical calcinatio (see /calcinatio) with a specialized subagent. Load when inscribing opera, auditing witness definitions, or shaping Magna Opera.
---

# Opus Operations

Tool-skills for the quality-sensitive acts in the opus lifecycle. Each operation fires a specialized subagent, then the calling agent reviews the output via dialectical calcinatio. **Load `/calcinatio` first** — it defines the protocol these operations depend on.

**When to use this skill:**
- Inscribing a new opus (any agent — azer at discharge, marut when queue is empty, artifex shaping work)
- Auditing witness definitions for specificity and actionability
- Shaping or refining a Magnum Opus

**What this skill does NOT cover:**
- Opus execution (that's the azer's domain)
- Discharge protocol (that's in `opus.md`)
- Operational dispatch (that's the marut's domain)

**Opus-specific convergence criteria:** Would a fresh azer reading this opus understand the witness's need and derive appropriate calcinatio? For Magna Opera, the additional test: does the artifex feel this captures their intent?

**Opus-specific elicitation:** For Magna Opera, ALWAYS elicit the artifex. An MO shapes all downstream work — the artifex's judgment is non-optional. For opera, escalate if you have ANY doubts about intent, witness identification, scope, or boundary — the cost of pausing is low; the cost of a mis-shaped opus is an entire azer context wasted.

---

## `/opus inscribe` — Opus Inscription

Inscribe a new opus by firing an inscription subagent with your context.

### What you (the calling agent) do:

**1. Gather context to pass:**

- **MO intent and witnesses** — the subagent needs to know WHO the work serves and WHAT "abundantly satisfied" means. Copy or summarize the relevant parts of the MO.
- **What you found** — the gap, the finding, the need. What evidence led you to believe an opus should exist? Be specific: "I found X, Y, Z. The witness needs A but currently experiences B."
- **Trail context** — brief summary of relevant prior opera. What's been tried? What's been learned? What's already handled?
- **Environment state** — what's the current situation? Deployed state, in-flight work, blockers, relevant recent events.

**2. Fire the subagent:**

```
Fire a subagent (model: sonnet) with:
- Your gathered context above
- Instruction: "Read the inscription protocol at [path to resources/inscription-protocol.md].
  Using the context provided, produce an opus file following the protocol.
  Return ONLY the opus file content."
```

The subagent path depends on where the skill is installed. Use the path relative to this SKILL.md file: `resources/inscription-protocol.md`.

**3. Review via dialectical calcinatio** (see protocol above):

Check specifically:
- Is the intent a witness experience delta? ("Witness experiences X but needs Y")
- Is the boundary clear? (What does the agent do vs. what does the operator do?)
- Is there procedure disguised as intent? (Steps, how-to, implementation details → those belong in Context, not Goal)
- Are the witnesses named specifically enough that the azer can derive calcinatio from them?
- Does the calcinatio section identify the right fires for these witnesses?

**4. Commit:**
- Write the opus file to `magna-opera/<mo-name>/opera/YYYY-MM-DD-<descriptive-name>.md`
- Set `status: charged` in YAML frontmatter

---

## `/opus audit-witnesses` — Witness Audit

Audit whether witness definitions in an MO are specific enough for agents to derive calcinatio and serve them effectively.

### What you do:

**1. Gather:**
- The MO's witness definitions (from the MO document)
- Optionally: examples of opera that struggled with witness clarity (from the trail)

**2. Fire subagent** with instruction to read `resources/witness-audit-protocol.md` and assess the witness definitions.

**3. Dialectical calcinatio.** The subagent returns a report: which witnesses are well-defined, which are vague, what's missing, proposed improvements. Review against your knowledge of the actual people/teams/stakeholders involved.

**4. Act on findings** — update the MO's witness section, or escalate to the artifex if witness definitions need human judgment.

---

## `/opus shape-mo` — Magnum Opus Shaping

Shape or refine a Magnum Opus document. Used when kindling a new MO or when the trail reveals the MO needs refinement.

### What you do:

**1. Gather:**
- The artifex's raw intent (what they said they want)
- Existing MO document (if refining)
- Trail evidence (if refining — what opera revealed about MO quality)
- Domain context — what project, what environment, what constraints

**2. Fire subagent** with instruction to read `resources/mo-shaping-protocol.md` and produce or refine the MO document.

**3. Dialectical calcinatio.** Check: Is the intent at the right altitude? Are witnesses named and specific? Is "abundantly satisfied" assessable? Would a marut know when to declare this done?

**4. Involve the artifex — MANDATORY for MO shaping.** Creation or refinement of a Magnum Opus is a big event — it shapes all downstream work. The artifex must be involved at one of these levels (prefer the highest feasible):

   - **Collaborative** — the artifex is present in the session, shaping the MO together. The ideal.
   - **Consulted** — present the refined draft to the artifex via escalation and **wait for their response** before committing. Use when the artifex isn't in the session.
   - **Informed** — if the MO is straightforward and clearly derived from explicit artifex intent (e.g., the artifex just said "I want X"), commit the draft and immediately notify the artifex with the full MO so they can review and request changes. The floor, not the default.

   Do not commit an MO silently. An MO the artifex hasn't seen is a system risk.

**5. Commit** — write or update the MO document.
