# Research Request: SKAN Agentic Ontology of Work (AOW) and Related Ontology Engineering

**Created:** 2026-03-29
**Status:** Open

## Goal

Map the SKAN Agentic Ontology of Work (AOW) framework in detail — its concepts, structure, and terminology — and identify what related work it points to in ontology engineering, organizational ontology, and work-theory domains. We're building an agentic orchestration system with a deeply considered domain vocabulary and want to understand how AOW's conceptual model compares, where it converges and diverges with ours, and what adjacent intellectual traditions (ontology engineering, activity theory, workflow ontologies, etc.) might sharpen our own design.

The specific point of interest: AOW appears to have a concept of "opus" and "magnum opus" in the context of agentic work, which directly parallels our system's core work-unit hierarchy. We need to understand their framing to see if there are insights we're missing, terminology collisions we should be aware of, or conceptual tools we should adopt.

## Context

We are building **the athanor** — an agentic orchestration system that serves as executive function infrastructure for a neurodivergent operator. The system accepts goals stated in stakeholder-value terms, breaks them into units of work called **opera** (singular: opus), and drives them to completion through autonomous AI agents under supervision.

### Our domain model (relevant excerpts)

**Core work-unit hierarchy:**
- **Magnum Opus** — Top-level goal container. Can be task-based (has a terminal state) or state-based (maintained continuously). Each MO is a directory containing its goal document and an `opera/` subdirectory. Intent is expressed as a change in the witness's experience: "Customers with mature Seismic setups feel the integration works for them, not against them."
- **Opus** — A single unit of work. Goal + geas (binding behavioral compulsion) + context. Charged to one agent at a time. Discharged when the intent is realized. Inscribed only when immediately actionable — no draft or queued opera.
- **Trail** — The sequence of discharged opera under a Magnum Opus. The authoritative record of what actually happened, not what was planned.

**Key design principles:**
- Opera are commitments to act, not plans. An opus is inscribed only when immediately actionable.
- The system proceeds inevitably toward conclusion through architectural properties, not individual agent reliability.
- Context accumulation IS progress — the system advances by reducing uncertainty.
- "Waste nothing" — at every boundary (discharge, session end), accumulated context gets released somewhere useful.
- "Abundant satisfaction" — completion standard is "would the stakeholder feel this is completely handled?" not "were requirements met?"
- Plans are materia (raw material the system processes), not first-order infrastructure.

**Vocabulary philosophy:**
We use distinctive names (drawn from alchemical, D&D monster manual, and esoteric traditions) because the naming IS the transmutation of institutional concepts. "Geas" activates reasoning about binding obligation; "acceptance criteria" activates PM-as-practiced including every way PM theater goes wrong. The vocabulary strips institutional scar tissue and lets agents reason from underlying principles.

| Term | Concept |
|------|---------|
| Artifex | The human operator who directs the system |
| Azer | Worker agent — skilled, bounded, charged with an opus |
| Marut | Supervisor agent — one per Magnum Opus, ensures goal advances |
| Opus | Single unit of work: intent + geas + context |
| Magnum Opus | Top-level goal container |
| Geas | Binding behavioral compulsion shaping agent behavior |
| Trail | Sequence of discharged opera — the authoritative record |
| Crucible | The agent's context window / tmux session |
| Materia | Anything in the world an agent can access via tools |

**Related work we've already identified:**
- **Foundry** (Aron Greenspan / Inixiative) — An independently developed agent orchestration system with striking conceptual parallels. Discovered 2026-03-23.

## Key Questions

1. **What is the SKAN Agentic Ontology of Work (AOW)?** What organization created it, when, and what problem does it solve? What is its stated purpose and scope?

2. **What are AOW's core concepts and taxonomy?** Map the full conceptual model — entities, relationships, hierarchies. Especially: how does it define "work," "task," "goal," "opus" (if it uses that term), and any magnum-opus-like higher-order goal structures?

3. **How does AOW model work decomposition and lifecycle?** How does work flow from high-level goals to actionable units? What states does a work item pass through? How does completion/fulfillment work?

4. **What does AOW say about agent roles and supervision?** Does it define worker/supervisor/orchestrator distinctions? How does it model the human-agent boundary?

5. **What ontology engineering traditions and frameworks does AOW reference or build on?** Follow the citation trail — what foundational work in ontology engineering, upper ontologies (DOLCE, BFO, SUMO), organizational ontology, or activity theory does it draw from?

6. **What adjacent frameworks exist in the "ontology of work" space?** Specifically:
   - Work ontologies from knowledge management / semantic web traditions
   - Activity Theory (Engestrom, Vygotsky) as applied to agent systems
   - OASIS/WfMC workflow ontologies
   - Enterprise ontology (Dietz, DEMO)
   - Organizational ontology research
   - Any agent-specific work ontologies from the multi-agent systems literature

7. **Where do AOW's concepts converge with and diverge from our model?** Produce a mapping table where possible: AOW concept → athanor concept, noting semantic differences.

8. **What concepts does AOW (or its related work) have that we lack?** Are there useful abstractions, relationships, or patterns we haven't considered?

## What Would Change Our Approach

- **If AOW has a mature, well-grounded conceptual model for work decomposition** that covers ground we've developed intuitively — we might adopt specific concepts or refine our own vocabulary with their insights.
- **If AOW points to established ontology engineering work** that formalizes concepts we use informally (like the relationship between intent, action, and completion) — we might ground our model more rigorously.
- **If there are terminology collisions** (e.g., if "opus" means something specific and different in AOW) — we need to know for interoperability and communication purposes.
- **If adjacent frameworks offer patterns for the human-agent boundary** that we haven't considered — these would directly inform our supervision and escalation model.
- **If the ontology engineering literature has established solutions** to problems we're solving from first principles (like how to model "inevitability" of completion, or state-based vs. task-based goal types) — we should know about them.

## Known Assumptions to Challenge

1. **We assume our domain model is novel.** It probably isn't — ontology engineering and organizational theory likely have decades of formalized thinking about work decomposition, goal hierarchies, and agent roles. We may be reinventing wheels.

2. **We assume alchemical/esoteric vocabulary is uniquely ours.** AOW or other frameworks may use similar metaphorical vocabulary for similar reasons (precision through distinctiveness).

3. **We assume "opus" and "magnum opus" carry the meaning we've assigned.** In AOW or related traditions, these terms might have specific semantic commitments we're not aware of.

4. **We assume our distinction between "plans as materia" and "opera as commitments to act" is a design choice.** It might map to an established ontological distinction in the literature (e.g., planning ontologies vs. action ontologies).

5. **We assume our "trail" concept (record of what happened, not what was planned) is a practical pattern.** It might correspond to established concepts in process ontology or provenance modeling.

## Output Format

Structured findings document with:
- Summary of AOW framework with source citations
- Concept mapping table: AOW ↔ athanor concepts
- Annotated bibliography of referenced/adjacent work, organized by relevance
- Synthesis section: what insights we should adopt or investigate further
- All claims cited per footnote format below

## Recency Requirements

Prefer recent (2023-2026) for AOW and agentic ontology work specifically — this is a fast-moving area tied to the LLM agent explosion. Older foundational content (ontology engineering, activity theory, enterprise ontology) is fine and expected — these fields are mature.

## Instructions for Researcher

### Footnote Format (Non-negotiable)

Use markdown footnotes for all significant claims.

**Inline reference:**
```markdown
AOW defines work as a goal-directed transformation of state[^aow-work-def].
```

**Footnote definition (at document end):**
```markdown
[^aow-work-def]: SKAN AOW Documentation, 2025. Section: "Core Concepts". "Work is defined as..." [View Source](https://example.com/aow) (accessed 2026-03-29)
```

**Naming convention:** Short descriptive slugs: `[^org-topic]` or `[^topic-subtopic]`

**Required components per footnote:**
1. Source context — what kind of source (paper, documentation, blog post, specification)
2. Date — when published
3. Author/organization — who created it
4. Key claim — brief quote or precise paraphrase
5. Link — URL when available, with access date for web sources

**Chain of evidence:** Footnotes travel with claims. When citing a source that AOW cites, cite the primary source directly, not AOW's citation of it. Never cite intermediate documents as sources.

### Evidence Requirements
- **Quote, don't paraphrase** — Include exact quotes for key findings
- **Link everything** — Every significant claim needs a footnote with source URL
- **Date your sources** — Note publication date AND access date in footnotes
- **Show evidence tiers** — Flag whether sources are Tier 1-4:

| Tier | Source Type | Weight |
|------|-------------|--------|
| 1 | Official docs, primary sources, peer-reviewed papers | Authoritative |
| 2 | Official examples, starter repos, official blogs | High |
| 3 | Verified community patterns (multiple sources agree) | Moderate |
| 4 | Single blog post, forum answer, outdated content | Low (flag uncertainty) |

### Research Process
- Start with SKAN AOW specifically, then follow its reference trail outward
- Distinguish between AOW's own concepts and concepts it inherits from upstream ontologies
- When mapping to our concepts, be precise about semantic differences — "similar" isn't enough, explain HOW they're similar and where they diverge
- Push back if the questions seem misguided or if there's a better framing for any of this
- If AOW turns out to be thin/shallow, say so — don't inflate it. Redirect energy to the richer adjacent traditions it points to.

### Transparency Requirements
- Flag `[UNCERTAIN]` when no authoritative source exists
- Flag `[CONFLICTING]` when sources disagree (present both sides)
- Flag `[SINGLE SOURCE]` when a claim relies on only one reference
- Flag `[OUTDATED]` when newest source is >2 years old on fast-moving topics

### Special Instruction
We already know our own system well — don't spend time explaining our concepts back to us. Focus your energy on understanding the external landscape and producing actionable comparisons. The concept mapping table is the highest-value deliverable.
