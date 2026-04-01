# Research Request: Generative Context Extraction Patterns for Agent Systems

**Created:** 2026-03-31
**Status:** Open

## Goal

Identify proven patterns, frameworks, or techniques for **generatively extracting maximum value from rich accumulated context** in agent orchestration systems. We're specifically interested in patterns where the challenge is: "An agent has accumulated a large amount of context (completed work records, environment signals, stakeholder definitions, goals). How do we generate the MOST and BEST ideas for what to do next?"

This is distinct from evaluation/critique (which we have good tools for). We need the generative/creative counterpart.

## Context

We are building an agent orchestration system called "the athanor" where:
- A **Magnum Opus** (MO) is a top-level goal expressed as a change in stakeholder (witness) experience
- **Opera** are units of work (charged → discharged → assessed)
- **Assessment** is the critical moment where the system surveys all available context and decides what to do next
- We have well-developed **critique/refinement** tools (dialectical calcinatio = iterative builder/reviewer exchange, manifold calcinatio = multiple independent critiques in parallel)
- We LACK equivalently powerful **generative** tools for the assessment moment

The specific problem: when a single agent session has a huge loaded context and we ask "what should we do?", ideas come out linearly — the most salient signals dominate, and the long tail of valuable-but-non-obvious ideas never surfaces. We've identified a candidate pattern we're calling "manifold generation" — decomposing rich context into individual stimuli, firing a subagent per stimulus (each paired with the constant goal/stakeholder context), and collecting a full idea space. But we want to know what's been explored in adjacent fields.

## Key Questions

1. **Are there established patterns in multi-agent systems, knowledge management, or organizational theory for generating diverse actionable ideas from accumulated context?** We're looking for the generative analog of "code review" or "design critique" — structured processes that produce ideas rather than evaluate them.

2. **How do knowledge management frameworks handle the "assessment" or "after-action review" moment?** Specifically frameworks that go beyond "lessons learned" capture toward "what should we do next given everything we know?" Military AAR, Toyota kata, retrospectives — but we want the generative versions, not just the reflective ones.

3. **In the semantic web / ontology space (DOLCE, BFO, SUMO, SKAN AOW), are there patterns for reasoning about "potential actions" given a rich state description?** We're interested in how formal ontologies handle the transition from "here is the state of the world" to "here are the possible valuable actions." Planning ontologies, action theories, situation calculus.

4. **What patterns exist for combating "mode collapse" in idea generation?** In ML, mode collapse means outputs cluster around the most probable center. The same happens with LLM brainstorming. Orthogonal emanation (forcing maximally independent candidates) is one approach. What others exist — in creativity research, design thinking, TRIZ, lateral thinking, or computational creativity?

5. **Are there multi-agent architectures specifically designed for creative/generative tasks?** Not critique (we have that). Not task decomposition (we have that). Specifically: multiple agents collaborating to GENERATE ideas that a single agent would miss. Swarm intelligence for ideation, blackboard architectures for brainstorming, etc.

6. **How does the "combinatorial stimulus" pattern appear in other domains?** Our candidate pattern decomposes context into stimuli and pairs each with a constant goal frame. Does this pattern have a name? Is it studied? Are there refinements we should know about? Think: SCAMPER, random stimulus in lateral thinking, oblique strategies, design space exploration.

## What Would Change Our Approach

- A well-studied pattern that's better than our "manifold generation" candidate → we'd adopt or adapt it
- Evidence that context decomposition into stimuli is studied and has known failure modes → we'd design around those
- A framework that connects generative and evaluative processes (our Chesed/Geburah split) → we'd use it as architectural guidance
- Academic work on "value extraction from completed work" that goes beyond lessons-learned → we'd integrate the patterns

## Known Assumptions to Challenge

- **We assume multiple focused subagents generate better ideas than one agent with full context.** This might be wrong — maybe there are techniques for getting a single rich-context agent to generate more diversely.
- **We assume the generative and evaluative phases should be separate.** Maybe the best patterns interleave them.
- **We assume this is a novel problem in the LLM agent space.** It's probably well-studied in other fields under different names.

## Output Format

Summary with sources, organized by question. Include a "Synthesis" section that maps findings to our specific assessment process design.

## Recency Requirements

Mixed — recent (2023-2026) for LLM/agent-specific patterns, but older foundational work is valuable for creativity research, knowledge management, ontology, and organizational theory. The underlying problem predates LLMs.

## Instructions for Researcher

### Footnote Format (Non-negotiable)

Use inline references like `[^source-slug]` in text. Define footnotes at the bottom of the document:

```
[^source-slug]: **Title or Description** — Author/Org, Date
  URL (if available) | Accessed: YYYY-MM-DD
  Evidence tier: [1-4] | Key quote: "exact quote if available"
```

Naming convention: `author-keyword-year` or `org-keyword-year` (e.g., `[^de-bono-lateral-1970]`, `[^skan-aow-2025]`)

### Evidence Requirements
- Quote, don't paraphrase — include exact quotes for key findings
- Link everything — every significant claim needs a footnote with source URL
- Date your sources — note publication date AND access date
- Show evidence tiers — Flag whether sources are Tier 1-4:
  - Tier 1: Official docs, primary sources, peer-reviewed (Authoritative)
  - Tier 2: Official examples, starter repos, official blogs (High)
  - Tier 3: Verified community patterns, multiple sources agree (Moderate)
  - Tier 4: Single blog post, forum answer, outdated content (Low — flag uncertainty)

### Research Process
- Search broadly across disciplines — this problem lives at the intersection of AI, knowledge management, creativity research, and organizational theory
- Push back if the questions seem misguided
- Distinguish well-established findings from emerging ideas
- Be direct about gaps

### Transparency Requirements
- Flag [UNCERTAIN] when no authoritative source exists
- Flag [CONFLICTING] when sources disagree
- Flag [SINGLE SOURCE] for single-reference claims
- Flag [OUTDATED] when newest source is >2 years old on fast-moving topics
