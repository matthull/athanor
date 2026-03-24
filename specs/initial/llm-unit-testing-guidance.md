# Research Request: LLM Unit Testing Guidance — Skills, Prompts, and Overmocking Mitigations

**Created:** 2026-03-14
**Status:** Open

## Goal

Find and evaluate **all significant approaches** to improving LLM-generated unit tests — with particular focus on reducing overmocking, improving test reliability, and producing tests that actually verify behavior rather than testing the test setup.

Scope covers three tiers:
1. **General LLM unit testing guidance** — framework-agnostic prompt engineering, system prompts, or skills that make any LLM write better tests
2. **Rails/RSpec-specific guidance** — Claude Code skills, CLAUDE.md patterns, or prompt strategies targeting Rails testing
3. **Popular Rails Claude Code skill sets** — any widely-used Rails skill collection is likely to contain testing guidance worth extracting

The outcome is actionable: specific skills to install, prompt patterns to adopt, or guidance to incorporate into our existing testing skills.

## Context

### The Core Problem

Claude (and LLMs generally) writes tests with a strong bias toward overmocking. When asked to test a service object, Claude's default behavior is:
1. Mock every collaborator with `instance_double`
2. Stub every method on those doubles
3. Assert that the stubs were called
4. Call this a "passing test"

This produces tests that verify the test setup, not the code. The tests pass even if the actual implementation is completely broken. Worse, they're brittle — any refactor that changes internal method calls breaks the tests, even if behavior is preserved.

**Specific failure modes we see repeatedly:**
- **Self-mocking:** Mocking methods on the object under test, then asserting those mocks were called
- **Mock trains:** `allow(a).to receive(:b).and_return(c)` chains 3-4 levels deep
- **Stubbing implementation details:** Testing HOW (which internal methods are called) instead of WHAT (what the observable behavior is)
- **Association mocking:** `allow(user).to receive(:account).and_return(double(...))` instead of using real model associations
- **Over-specified mocks:** `expect(service).to receive(:call).with(exactly(...)).once` when a simple `allow` + behavior assertion would suffice
- **Missing integration reality:** Tests that mock the database, then assert database state — impossible to fail, impossible to catch bugs

### What We Already Have

We maintain custom Claude Code skills and guidance files. The research should evaluate external solutions **against these** — what gaps do they fill, what do they do better, what conflicts?

**`unit-testing` skill** (framework-agnostic):
- Rule 1: Never use the real clock (pin/freeze time)
- Rule 2: Never use real randomness (seed or mock)
- Rule 3: Isolate external state (env vars, filesystem, network, DB)
- Rule 4: Tests must be independent (no test-order dependencies)
- Flakiness pattern catalog (10 patterns with fixes)
- Test quality gates checklist

**Mocking guidelines** (Rails-specific, 3 separate files):
- Core principle: Mock at system boundaries only
- "Quick Rule: If you type `allow(service)` where service is the object under test, STOP"
- Good/bad examples for services, controllers, models, background jobs
- Red flags: mocking associations, `allow_any_instance_of`, >3 mocks per test, more mock setup than test code
- 80/20 rule: 80% real objects, 20% mocks (external APIs, file storage, slow operations)
- Self-mocking anti-pattern with before/after examples

**Testing strategy** (layered approach):
- Unit tests: RSpec for models, services, jobs (minimal mocking)
- Integration tests: Request specs (mock external APIs only)
- Visual tests: Storybook stories
- E2E: Browser automation

**Rails skill** (testing mandates):
- Fixture builder constraint (NEVER edit .yml fixtures, use fixture_builder.rb)
- Fixtures for typical data, factories for edge cases only
- Test every public method at its layer
- Never `allow_any_instance_of`
- Request specs over controller specs
- Prefer `let` over instance variables

**Known gaps in our current guidance:**
- No structured "mock escalation ladder" (when to add a mock, in what order of preference)
- No explicit "test smell" detection checklist for code review
- No guidance on testing private methods (we say "test public interface" but don't address the common scenario where Claude wants to test a private method)
- No prompt-level strategies — our guidance is all declarative rules, not structured as LLM-optimized prompts
- Limited coverage of what GOOD tests look like beyond mocking (assertions, test naming, arrange-act-assert patterns for complex scenarios)

### Project Stack (for Rails-specific evaluation)
- Ruby on Rails 7.2, PostgreSQL, Sidekiq
- RSpec with fixture_builder gem, WebMock, VCR
- Vitest for JS utilities, Storybook for Vue components
- Docker Compose development environment

### Prior Related Research

From a 2026-03-08 research request on Rails/Vue domain skills:
- `matsengrp/claude-code-agents` clean-code-reviewer was flagged as having "best over-mocking treatment found" — **dig deeper into this specifically**
- `elifiner/refactoring` was flagged as high quality, Fowler-based, analysis-first workflow
- Generic clean code / refactoring skills were covered — **do NOT re-research those**

## Key Questions

### Q1: General LLM Unit Testing Strategies (Highest Priority)

1. What published approaches exist for making LLMs write better unit tests? This includes:
   - Claude Code skills specifically for testing
   - System prompts or prompt engineering techniques for test generation
   - Cursor rules, Windsurf rules, or other AI coding tool configurations for testing
   - Research papers or blog posts on LLM test generation quality
   - Any "AI testing manifesto" or community-developed standards

2. Specifically for overmocking:
   - Has anyone published a systematic approach to preventing LLM overmocking?
   - Are there prompt patterns that reliably reduce mocking (e.g., "write the test first with zero mocks, then add mocks only for external boundaries")?
   - Are there "mock escalation ladder" frameworks (prefer real object → prefer fixture → prefer factory → stub method → use double → use instance_double, in that order)?

3. What does the Anthropic documentation or Claude best practices say about test generation? Any official guidance on writing better tests with Claude?

### Q2: Rails/RSpec-Specific Skills and Guidance

1. Are there Claude Code skills specifically for RSpec testing?
   - Skills that encode mocking boundaries, fixture/factory guidance, test structure
   - Skills from Rails consultancies (thoughtbot, Evil Martians, Shopify)

2. **Deep dive on `matsengrp/claude-code-agents`**: Our prior research flagged the clean-code-reviewer as having the best overmocking treatment found. Get the actual content — what specific rules does it encode? How does it detect and prevent overmocking? What's the prompt structure?

3. Are there CLAUDE.md files in major Rails open-source projects that include testing guidance? Check:
   - rails/rails
   - thoughtbot repos (factory_bot, shoulda-matchers, etc.)
   - Shopify repos
   - GitLab (uses Rails + RSpec extensively)
   - Discourse
   - Mastodon
   - Any other large Rails project with Claude/AI coding tool configuration

4. Are there `.cursorrules` or similar AI coding tool config files in Rails projects that include testing guidance? These often contain the same kind of guidance as CLAUDE.md files.

### Q3: Testing Skill Architecture

1. How do high-quality testing skills structure their guidance?
   - Rules vs. examples — what ratio works best for LLM compliance?
   - Positive examples ("do this") vs. negative examples ("don't do this") — which is more effective?
   - Does progressive disclosure help? (Load basic rules always, detailed examples on demand)
   - Should testing guidance be one monolithic file or split by concern?

2. What enforcement mechanisms exist beyond declarative rules?
   - Pre-commit hooks that detect overmocking?
   - Custom linting rules (rubocop-rspec patterns)?
   - Test quality metrics that can be checked programmatically?

3. Are there any testing skills that use a "review-then-fix" pattern — write the test, then systematically review it against a checklist before committing?

### Q4: Popular Rails Claude Code Skill Sets

1. What are the most popular/starred Claude Code skill repositories for Rails?
2. For each, extract and evaluate their testing-related guidance specifically
3. Are there any "Rails starter" CLAUDE.md templates being shared in the community?
4. Check the Anthropic skill marketplace/registry (if one exists) for Rails-related skills

### Q5: Cross-Framework Testing Patterns

1. Are there testing principles from other ecosystems (Go, Rust, Elixir) that have been successfully encoded as LLM guidance and could transfer to Rails?
   - Go's table-driven tests pattern
   - Rust's `#[cfg(test)]` module pattern
   - Elixir's ExUnit patterns
2. Are there any language-agnostic "LLM testing constitutions" — sets of principles designed to constrain LLM test generation across any framework?

## What Would Change Our Approach

- **A well-tested prompt engineering pattern exists for reducing overmocking** → We'd restructure our mocking guidelines from declarative rules to LLM-optimized prompt patterns (e.g., "zero-mock-first" workflow)
- **A high-quality RSpec testing skill exists** → We'd evaluate installing it and making our guidance a thin extension
- **`matsengrp/claude-code-agents` has a transferable overmocking detection approach** → We'd extract and adapt it into our unit-testing skill
- **No one has solved this well** → Confirms this is a genuine gap. We'd invest in developing a comprehensive solution and potentially open-source it
- **Enforcement mechanisms exist (linting, hooks)** → We'd add programmatic gates rather than relying solely on prompt-level guidance
- **A "mock escalation ladder" exists somewhere** → We'd adopt or adapt it as a core testing principle

## Known Assumptions to Challenge

1. **"More examples = better compliance"** — LLMs may respond better to principles + a few examples than to exhaustive example catalogs. Research what actually works.
2. **"Declarative rules are sufficient"** — Our current approach is "here are the rules, follow them." Maybe a workflow-based approach ("Step 1: write test with no mocks. Step 2: identify what fails. Step 3: add minimal mocks only for failures at system boundaries") would be more effective.
3. **"This is a Claude-specific problem"** — Check if other LLMs (GPT-4, Gemini) have the same overmocking tendency. If so, solutions are likely transferable. If not, the solution may need to be Claude-specific prompt engineering.
4. **"Testing guidance belongs in a standalone skill"** — Maybe testing guidance is more effective when embedded in the domain skill (e.g., the Rails skill's testing section) rather than a separate testing skill loaded alongside. Check if anyone has research on this.
5. **"The problem is knowledge, not attention"** — Claude may KNOW not to overmock but lose track of those instructions in a long context. The solution might be reinforcement/repetition rather than better initial rules.

## Output Format

**Structured findings per question area**, each with:
- Key findings with evidence (quotes, links)
- Actionable recommendations
- Confidence level (high/medium/low)

Plus a **synthesis section** with:
- Priority-ordered action items (what to do first)
- Skills/resources to install immediately
- Patterns to adopt in our existing skills
- Gaps that need custom development
- Open questions for further research

## Recency Requirements

**Critical (2024-2026 only)** for:
- Claude Code skills (ecosystem is months old, older content is likely outdated)
- AI coding tool configurations (.cursorrules, CLAUDE.md, etc.)
- Prompt engineering for code generation

**Prefer recent but older acceptable** for:
- Testing principles and anti-patterns (well-established knowledge)
- Mocking philosophy (books like "Growing Object-Oriented Software" are still relevant)
- RSpec best practices (stable ecosystem)

## Instructions for Researcher

### Footnote Format (Non-negotiable)

Use inline reference marks like `[^test-1]` in the text body. Define each footnote at the bottom of its section (or document) with this exact format:

```
[^test-1]: **Title or description** — [URL](url) (accessed YYYY-MM-DD). Tier N. "Exact quote if applicable."
```

Components:
- **Reference mark:** `[^kebab-case-prefix-N]` where prefix groups related sources (e.g., `[^mock-1]`, `[^rails-skill-1]`, `[^prompt-1]`)
- **Title:** Descriptive label for the source
- **URL:** Direct link (not shortened)
- **Access date:** When you visited the page
- **Evidence tier:** Tier 1 (official docs/primary), Tier 2 (official examples/blogs), Tier 3 (verified community, multiple sources agree), Tier 4 (single blog/forum, flag uncertainty)
- **Quote:** Exact quote for key claims (strongly preferred over paraphrase)

### Evidence Requirements
- **Quote, don't paraphrase** — Include exact quotes for key findings
- **Link everything** — Every significant claim needs a footnote with source URL
- **Date your sources** — Note publication date AND access date in footnotes
- **Show evidence tiers** — Flag whether sources are Tier 1-4
- **For GitHub repos**: Include star count and last commit date as freshness indicators

### Search Strategy Suggestions

1. **GitHub search for skills:** `claude code skill rspec`, `SKILL.md rspec`, `SKILL.md testing`, `claude.md rails testing`
2. **GitHub search for configs:** `.cursorrules rails`, `.cursorrules rspec`, `claude.md rspec`
3. **Major Rails repos:** Check rails/rails, discourse/discourse, mastodon/mastodon, gitlabhq/gitlabhq, solidusio/solidus for `.claude/`, `CLAUDE.md`, `.cursorrules`
4. **Skill aggregators:** awesome-claude-code, awesome-cursorrules, community skill marketplaces
5. **Anthropic resources:** Claude documentation on testing, best practices guides, prompt engineering guides
6. **Research/blogs:** "LLM test generation quality", "AI overmocking", "Claude Code testing", "prompt engineering unit tests"
7. **`matsengrp/claude-code-agents`**: Deep dive — read the actual skill files, understand the architecture
8. **Testing philosophy sources:** "Growing Object-Oriented Software" community discussions, Martin Fowler's testing articles, thoughtbot testing blog posts — specifically as they relate to LLM guidance

### Transparency Requirements
- Flag `[UNCERTAIN]` when no authoritative source exists
- Flag `[CONFLICTING]` when sources disagree (present both sides)
- Flag `[SINGLE SOURCE]` when a claim relies on only one reference
- Flag `[OUTDATED]` when newest source is >2 years old on fast-moving topics

### Research Process
- Push back if any questions seem misguided or if there's a better framing
- Distinguish between well-established findings and emerging/contested ideas
- Be direct about limitations — if the LLM testing guidance space is sparse, say so
- **Prioritize Q1 and Q2 over Q4 and Q5** — general overmocking solutions and Rails-specific approaches are more valuable than exhaustive skill catalogs
- If you find a particularly good resource, go deep rather than cataloging many mediocre ones
