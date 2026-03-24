# Research Request: Layered Review Loops for LLM-Generated Code

**Created:** 2026-03-14
**Status:** Open

## Goal

Find out if anyone is implementing **multi-level review loops** for LLM-generated code — specifically, review architectures that go beyond line-level linting and PR-level code review to include **architecture-level evaluation** as a continuous concern, not a one-time gate.

The hypothesis: LLMs are worse at architecture than humans but better at critique than creation. This means LLM-generated code accumulates structural problems that look fine line-by-line but degrade the system over time. Traditional code review (even by other LLMs) doesn't catch this because it's scoped to the diff. Someone must be addressing this — or if nobody is, that itself is a finding worth understanding (is it because it's hard, or because nobody's thought about it, or because the problem isn't as bad as we think?).

**Secondary question:** Are people front-loading review criteria into implementation prompts? The idea is simple — show the coding agent what the reviewer will check before it writes code. This is nearly free (just append criteria to the prompt) and should improve first-pass quality. Is anyone doing this systematically?

## Context

### The Problem We're Seeing

We use Claude Code extensively with an orchestrated pipeline: spec → task plan → implement → verify → ship. The implementation agent has domain skills loaded (Rails patterns, Vue patterns, testing rules). The verify phase runs spec-check, QA classification, and code review.

What we observe:
- **Line-level quality is good** — Claude follows style rules, writes tests, handles edge cases
- **Feature-level quality is decent** — changes work, tests pass, the PR does what the spec says
- **Architecture-level quality is subtly poor** — service boundaries are slightly wrong, data models have unnecessary coupling, abstractions are created where none are needed, patterns are applied without understanding when they're appropriate

Each individual decision is defensible. But they compound. After 20 PRs of LLM-generated code, the codebase has drifted toward complexity that no single PR introduced but that collectively makes the system harder to work with.

**This isn't caught by:**
- Linting (style/syntax only)
- Standard code review (scoped to the diff)
- Spec compliance checking (the code does what the spec says — the spec just didn't say "don't create an unnecessary service object")
- Tests (they pass — the architecture problem is about maintainability, not correctness)

### Three Review Scopes

We've identified three distinct scopes that need different review approaches:

| Scope | What's reviewed | Deterministic? | Loop speed | Current state |
|-------|----------------|---------------|-----------|---------------|
| **Code-level** | Individual lines/methods | Mostly yes (linters, types) | Per-file write | Strong — RuboCop + ESLint hooks |
| **Feature-level** | Whole PR / change set | Mixed | Per-PR | Decent — convoy-style code review (6 parallel agents) |
| **Architecture-level** | System structure, boundaries, coupling | Almost entirely non-deterministic | Per-feature? Periodic? | **Nothing** |

### The Front-Loading Insight

We've also noticed that review criteria are currently applied only post-hoc. The implementation agent doesn't know what the reviewer will check. A simple optimization: inject the reviewer's evaluation criteria into the implementation prompt as a "you will be evaluated on" section. This should be nearly free and high-ROI, but we haven't seen anyone discuss doing this systematically.

### What We Already Have (Review Skills)

We have a multi-dimensional review system modeled after "convoy-style" parallel review agents:

**`/code-review`** — 6 parallel agents: architecture, quality, security, performance, testing, docs. Each produces structured findings with P0/P1/P2 severity. Synthesis deduplicates across agents.

**`/spec-check`** — Per-requirement PASS/FAIL against the implementation diff. Cleanly scoped to spec compliance only.

**`/qa-plan`** — Classifies verification items by execution channel: [TEST], [BROWSER], [CONSOLE], [HUMAN].

**`/task-review`** — Reviews uncommitted code against task handoff quality gates.

**What we DON'T have:**
- Spec/PRD review before implementation (questions for the author, not gap-filling)
- Plan review before execution (GO/NO-GO gate)
- Architecture-level review (does this change fit the system's structural constraints?)
- Front-loading of review criteria into implementation prompts

We have a detailed internal analysis comparing our review skills against a reference implementation (Gastown) that uses 10-leg code review convoys, PRD review convoys, plan review convoys, and design convoys. The key insight from that analysis: reviews should produce QUESTIONS for humans (not assertions), force binary GO/NO-GO judgments, and cross-reference findings across review dimensions for confidence signals.

### Why Architecture Review Should Be Easy (But Maybe Isn't)

This is what puzzles us. An architecture review loop for LLM-generated code seems straightforward:

1. Before implementing a feature, describe the system's current architecture and constraints
2. After implementing, have a separate LLM pass review the changes through an architecture lens
3. Feed findings back before the next feature

LLMs are good at critique. They're good at spotting coupling, unnecessary abstraction, and pattern misapplication — when explicitly asked to look for it. The standard code review prompt just doesn't ask at the right level.

So either:
- **People are doing this and we haven't found it** — in which case, show us
- **People aren't doing this because it's harder than it looks** — in which case, what makes it hard? (Note: the hard part may not be the review loop itself — encoding tacit architectural knowledge is hard, but the review process could be 100% manual/human and still valuable. LLM tools can facilitate by building codebase maps and dependency graphs, even if they shouldn't be trusted for the architectural judgment itself.)
- **People aren't doing this because they haven't thought about it** — in which case, the finding is "green field opportunity"

### The Reinvestment Argument

Research consistently shows LLM-generated code has lower architectural quality than human-written code — this is a known finding, not a hypothesis to test. The question isn't whether the problem exists but whether anyone is systematically addressing it.

Agentic coding saves substantial implementation time. Some of that saved time should be reinvested in reviewing what was built. Architecture review is the highest-leverage place for that investment because:

1. **It's what LLMs are worst at** — line-level quality is good, architecture is where the gap is widest
2. **It's what standard review processes miss** — linting catches syntax, code review catches logic, nothing catches structural drift
3. **It compounds** — each individually-reasonable decision makes the next one harder to correct

### The Meta-Learning Argument

There's a second motivation beyond catching problems: **understanding the tool's failure modes.** Agentic coding systems are new. Without systematic review of what they produce architecturally, you can't:
- Know what kinds of structural decisions they make poorly
- Improve your skills/prompts/guardrails based on observed patterns
- Build a feedback loop that makes the system better over time
- Distinguish between "the LLM made a bad decision" and "our spec/constraints were underspecified"

This is an observability investment, not just a quality investment. You can't improve a system you're not observing.

### The Circularity Risk

Using LLMs to review architecture of LLM-generated code has a circularity problem: the same model that decided "a service object is the right abstraction here" may defend that decision when asked to review it. Mitigations:
- Use a different model for review than implementation
- Frame review with specific architectural criteria that override general "is this reasonable?" assessment
- Treat LLM review as a facilitator (generate maps, surface metrics, identify patterns) while keeping architectural judgment with humans
- At minimum, LLM review with explicit criteria is better than no review at all

## Key Questions

### Q1: Multi-Level Review Architectures (Highest Priority)

1. Is anyone implementing review systems with distinct scopes (code-level, feature-level, architecture-level) for LLM-generated code?
   - Academic research on review granularity for AI-generated code?
   - Industry blog posts or talks about layered review?
   - Open-source implementations of multi-scope review?

2. Specifically, is anyone doing **architecture-level review** of LLM-generated code?
   - Not just "is this code clean?" but "does this change fit the system?"
   - Checking for coupling introduced, abstractions created, patterns applied
   - Evaluating whether an LLM's structural decisions match the codebase's existing architecture
   - Any tools, prompts, or processes for this?

3. Are there "architecture drift detection" tools or practices?
   - Tools that track coupling metrics over time and flag increases
   - Dependency analysis that detects when new coupling is introduced
   - Anything that compares a change's structural impact against architectural intent

4. How do people handle the **accumulation problem** — where each change is locally reasonable but globally degrades the system?
   - Is this a recognized problem in LLM coding?
   - Are there periodic "architecture health" reviews?
   - Any evidence on how fast LLM-generated code degrades system structure vs. human-written code?

### Q2: Front-Loading Review Criteria

1. Is anyone systematically injecting review criteria into implementation prompts?
   - "You will be evaluated on: [checklist]" patterns
   - Sharing the reviewer's rubric with the implementer
   - Showing the coding agent what tests/checks will run against its output

2. Does front-loading actually improve first-pass quality?
   - Any A/B testing or measurement of this?
   - Research on LLM performance when given evaluation criteria upfront vs. not?
   - Prompt engineering literature on "rubric-aware generation"?

3. What's the right level of specificity?
   - Literal reviewer checklist vs. summarized criteria?
   - Domain rules only, or also process rules (commit discipline, test structure)?
   - Does too much criteria overwhelm the implementation context?

### Q3: Continuous Review vs. Gate Review

1. Are people running review checks **during** implementation, not just after?
   - Reviewer agent watching code as it's written (not just linting, but semantic review)
   - Periodic check-ins during long implementation tasks
   - "Micro-review loops" at natural breakpoints (after each file, after each task in a list)

2. What's the performance/cost trade-off?
   - LLM review on every file write is expensive — is anyone doing it? Is it worth it?
   - Are there lighter-weight "sniff tests" that run inline, with full review deferred?
   - Does the Cursor / Copilot ecosystem have inline review patterns?

3. How do continuous review findings feed back into the implementation agent?
   - Inject as system prompt context?
   - Block and require fix (like our lint hooks)?
   - Accumulate and present at a breakpoint?

### Q4: The Organizational Analog

Software organizations have architecture review boards, design reviews, and technical debt tracking. Are people recreating these for LLM-assisted development?

1. **Architecture Decision Records (ADRs)** — are people making LLMs generate and respect ADRs?
   - LLM reads ADRs before implementing → constrained by previous decisions
   - LLM generates new ADRs when making structural decisions → traceable
   - Any tools or practices for this?

2. **Technical debt tracking** — are people tracking LLM-introduced tech debt systematically?
   - Beyond TODO comments — actual debt registers
   - Automated detection of debt introduced by LLM changes

3. **"Architecture as code"** — are people encoding architectural constraints in machine-readable form?
   - ArchUnit (Java) style — architectural tests
   - Dependency rules (module A should not import module B)
   - Layer enforcement (controllers don't call repositories directly)
   - Is anyone doing this for Ruby/Rails or TypeScript?

## What Would Change Our Approach

- **Multi-level review architectures exist and work** → We'd implement a three-tier review system (code/feature/architecture) and integrate all three into our pipeline
- **Front-loading criteria is proven effective** → We'd add an "evaluation preview" section to every implementation handoff, extracted from the verify phase's skill criteria. This is nearly free to implement.
- **Architecture review during/after implementation is common** → We'd build an architecture review leg, probably as a periodic check rather than per-PR
- **"Architecture as code" is practical for Rails/TS** → We'd encode architectural constraints as tests/rules (like ArchUnit but for our stack)
- **Nobody is doing this** → Genuinely a gap. We'd design and build it, starting with the cheapest intervention (front-loading criteria) and working up to the most expensive (continuous architecture review). Worth open-sourcing the approach.
- **Architecture review produces useful meta-learning even when it doesn't find problems** → The observations about what agentic coding produces become the basis for better skills, guardrails, and orchestration. We'd build the review loop partly as an observability investment — you can't improve a system you're not observing.

## Known Assumptions to Challenge

1. **"LLMs are specifically bad at architecture"** — Research confirms LLM-generated code has lower architectural quality. But what specifically do they get wrong? Unnecessary abstraction? Wrong boundaries? Pattern misapplication? Over-engineering? Understanding the failure modes shapes what the review should look for.
2. **"Architecture problems accumulate invisibly"** — Maybe they're visible but we're not measuring the right things. What metrics would surface this? Coupling metrics? Dependency graph complexity? File/class count growth rate?
3. **"Review criteria front-loading helps"** — Maybe it just makes the prompt longer without changing behavior. Has anyone tested this?
4. **"The code review diff scope is the problem"** — Maybe the problem is that reviewers (human or LLM) don't have enough system context, not that the scope is wrong. Providing full system context to a diff-scoped review might be sufficient.
5. **"This is easy to do"** — We assume architecture review is straightforward because LLMs are good at critique. But maybe the hard part is defining WHAT to critique — the architectural constraints themselves are tacit knowledge that's hard to encode. Counter-argument: even manual/human architecture review with no LLM assistance is valuable if you're reinvesting time saved by agentic coding. The loop doesn't require automation to be worthwhile.

## Output Format

**Structured findings per question**, with emphasis on:
- Concrete implementations found (tools, repos, blog posts with actual details)
- Research findings on LLM code quality at different review scopes
- Patterns that can be immediately adopted vs. requiring significant investment
- Evidence for or against the core hypothesis (LLM code needs architecture-level review)

**Synthesis:**
1. **State of practice** — Is anyone doing this? How mature is it?
2. **What works** — Proven approaches with evidence
3. **What's missing** — Gaps in the ecosystem
4. **Our recommended path** — Based on findings, what should we build/adopt first?
5. **Open questions** — What we still don't know after research

## Recency Requirements

**Critical (2024-2026 only):**
- LLM coding review practices
- AI-assisted development workflow patterns
- Prompt engineering for code quality

**Prefer recent, older acceptable:**
- Architecture review practices (enterprise patterns are stable)
- ADR practices
- ArchUnit and "architecture as code" patterns

**Not critical:**
- General code review theory
- Software architecture principles

## Instructions for Researcher

### Footnote Format (Non-negotiable)

Use inline reference marks like `[^review-1]` in the text body. Define each footnote at the bottom of its section (or document) with this exact format:

```
[^review-1]: **Title or description** — [URL](url) (accessed YYYY-MM-DD). Tier N. "Exact quote if applicable."
```

Components:
- **Reference mark:** `[^kebab-case-prefix-N]` where prefix groups related sources (e.g., `[^arch-1]`, `[^loop-1]`, `[^front-1]`)
- **Title:** Descriptive label for the source
- **URL:** Direct link (not shortened)
- **Access date:** When you visited the page
- **Evidence tier:** Tier 1 (official docs/primary/peer-reviewed), Tier 2 (official examples/blogs), Tier 3 (verified community, multiple sources agree), Tier 4 (single blog/forum, flag uncertainty)
- **Quote:** Exact quote for key claims (strongly preferred over paraphrase)

### Evidence Requirements
- **Quote, don't paraphrase** — Include exact quotes for key findings
- **Link everything** — Every significant claim needs a footnote with source URL
- **Date your sources** — Note publication date AND access date in footnotes
- **Show evidence tiers** — Flag whether sources are Tier 1-4

### Search Strategy

1. **LLM code quality research:** "LLM code architecture quality", "AI generated code technical debt", "LLM coding structural problems", "AI pair programming code quality study"
2. **Review architectures:** "multi-level code review AI", "architecture review LLM", "code review beyond diff"
3. **Front-loading/rubric-aware generation:** "evaluation criteria in prompt", "rubric-aware code generation", "self-evaluation LLM coding"
4. **Architecture as code:** ArchUnit, dependency-cruiser (JS), Deptry (Python), packwerk (Ruby/Shopify), "architectural fitness functions"
5. **Architecture drift:** "architecture drift detection", "coupling metrics over time", "structural erosion software"
6. **Claude Code / Cursor / Copilot:** Community discussions about review workflows, multi-pass review, architecture concerns
7. **Enterprise architecture review patterns adapted for AI:** Any adaptation of architecture review boards, design authority, technical governance for AI-assisted development
8. **ADRs + LLMs:** "architecture decision records AI", "LLM respect ADR", tools that integrate ADRs into AI coding workflows

### Research Process
- **Prioritize Q1 and Q2** — Multi-level review architectures and front-loading criteria are the core questions
- If the answer to Q1 is "nobody is doing this," go deep on Q4 (organizational analogs) — maybe the patterns exist in traditional software engineering and just haven't been adapted for LLM workflows
- If you find ONE good implementation of architecture-level review for LLM code, go very deep on it — extract the full approach
- **Do NOT spend time proving the problem exists** — research already establishes that LLM-generated code has lower architectural quality. Take that as given. Focus on: what specifically goes wrong, who's addressing it, and what approaches work.
- **Look for the meta-learning angle** — is anyone using architecture review not just to catch problems, but to systematically understand and catalog agentic coding failure modes? This is the "observability investment" framing: review as a way to improve the system, not just the output.
