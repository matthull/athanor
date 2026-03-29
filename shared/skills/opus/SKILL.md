---
name: opus
description: Opus lifecycle operations — inscription, witness audit, MO shaping. Fires specialized subagents for quality-sensitive acts using the calcinatio loop (context-holder/specialist refinement). Load when inscribing opera, auditing witness definitions, or shaping Magna Opera.
---

# Opus Operations

Tool-skills for the quality-sensitive acts in the opus lifecycle. Each operation fires a specialized subagent with focused expertise, then the calling agent reviews the output against its rich context. This is the calcinatio loop — neither side alone produces quality; the interplay does.

**When to use this skill:**
- Inscribing a new opus (any agent — azer at discharge, marut when queue is empty, artifex shaping work)
- Auditing witness definitions for specificity and actionability
- Shaping or refining a Magnum Opus

**What this skill does NOT cover:**
- Opus execution (that's the azer's domain)
- Discharge protocol (that's in `opus.md`)
- Operational dispatch (that's the marut's domain)

---

## The Calcinatio Loop

*This is the general protocol for any azer-tool interaction where focused expertise matters. It applies to all operations in this skill and to other tool-skills (code review, research, etc.).*

**The problem it solves:** You have rich context but can't do a focused task well inline — your context biases you, or the task requires a different cognitive mode. A subagent has fresh focus but might produce poor output because the context handoff was lossy.

**The protocol:**

```
1. SHAPE context for handoff
   - What does the subagent need to do its job?
   - What's the minimum context that preserves the important signal?
   - What would bias or distract the subagent? Leave that out.

2. FIRE the subagent
   - Include: the shaped context + instruction to load the relevant protocol file
   - The subagent reads the protocol (not you — keep your context lean)
   - The subagent produces an artifact (opus, audit report, MO draft)

3. REVIEW the artifact against your rich context
   - You are the calcinatio fire here — you burn away what doesn't survive
     contact with the reality you know
   - Does the output align with what you know about the situation?
   - Did the subagent miss something your context reveals?
   - Is the intent truly a witness experience delta, or did it slide into procedure?

4. REFINE if needed
   - Fire again with corrective guidance: "You missed X context" / "The witness
     is actually Y, not Z" / "This is procedural — the intent should be..."
   - The subagent's first pass is EXPECTED to be imperfect
   - The review is EXPECTED to catch gaps
   - Refinement is the normal path, not a failure case

5. ELICIT the artifex (when present or when in doubt)
   - If the artifex is present (collaborative mode): present the refined artifact
     for their reaction. Their judgment is the highest-fidelity fire.
   - If the artifex is absent: review the artifact honestly — do you have ANY
     doubts about intent, witness identification, scope, or boundary? If yes,
     escalate to the artifex before committing. The cost of pausing is low;
     the cost of a mis-shaped opus is an entire azer context wasted.
   - For Magna Opera: ALWAYS elicit the artifex. An MO shapes all downstream
     work — the artifex's judgment is non-optional here.

6. CONVERGE
   - When the context-holder (and artifex, if involved) has nothing substantive
     to add, it's done
   - Commit the artifact (write the opus file, update the MO, etc.)
```

**Convergence criteria:** The context-holder reviews the output and has no corrections that would change how an agent acts on it. Cosmetic preferences don't count. The test: would a fresh azer reading this opus understand the witness's need and derive appropriate calcinatio? For Magna Opera, the additional test: does the artifex feel this captures their intent?

**Why this works:** Multiple fires from different angles. The specialist's expertise burns away procedural thinking and witness-blindness. The context-holder's knowledge burns away context-poor assumptions. The artifex's judgment burns away misalignment with true intent. What survives all fires is refined.

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

**3. Review via calcinatio loop** (see protocol above):

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

**3. Calcinatio loop.** The subagent returns a report: which witnesses are well-defined, which are vague, what's missing, proposed improvements. Review against your knowledge of the actual people/teams/stakeholders involved.

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

**3. Calcinatio loop.** Check: Is the intent at the right altitude? Are witnesses named and specific? Is "abundantly satisfied" assessable? Would a marut know when to declare this done?

**4. Involve the artifex — MANDATORY for MO shaping.** Creation or refinement of a Magnum Opus is a big event — it shapes all downstream work. The artifex must be involved at one of these levels (prefer the highest feasible):

   - **Collaborative** — the artifex is present in the session, shaping the MO together. The ideal.
   - **Consulted** — present the refined draft to the artifex via escalation and **wait for their response** before committing. Use when the artifex isn't in the session.
   - **Informed** — if the MO is straightforward and clearly derived from explicit artifex intent (e.g., the artifex just said "I want X"), commit the draft and immediately notify the artifex with the full MO so they can review and request changes. The floor, not the default.

   Do not commit an MO silently. An MO the artifex hasn't seen is a system risk.

**5. Commit** — write or update the MO document.
