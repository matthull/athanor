# Research Request: Deterministic Code Quality Guardrails for LLM Coding

**Created:** 2026-03-14
**Status:** Open

## Goal

Find **deterministic techniques** — beyond standard linting — that enforce code quality on LLM-generated code. The operating principle: if a rule can be checked by a tool rather than a prompt, it should be. LLMs need *more aggressive* guardrails than human coding because they fail silently and consistently on the same patterns.

We're NOT looking for:
- Code smell detection / heuristic analysis (Reek, Flog, Flay — already explored)
- Standard linter configuration advice (we have RuboCop + ESLint configured)
- General code quality philosophy

We ARE looking for:
- **Techniques that go beyond what standard linters offer** — pattern matching, type-level enforcement, contract checking, custom AST analysis, property-based verification
- **Hook/gate architectures** that integrate deterministic checks into LLM coding workflows at the right points
- **Custom lint rules people are writing** specifically because LLMs need them (anti-patterns that humans rarely produce but LLMs frequently do)
- **Emerging tools** (Semgrep, oxlint, Biome, etc.) being used as deterministic layers
- **Cross-file and project-wide analysis** techniques that catch issues single-file linting misses

## Context

### The Determinism Gap

We have prompt-level guidance (Claude Code skills) encoding rules like "never mock the object under test" and "never use `allow_any_instance_of`." These are probabilistic — Claude follows them ~90% of the time. We also have standard linting (RuboCop, ESLint) running as PostToolUse hooks that auto-correct and block on violations.

The gap is between these layers. Rules that are **too semantic for standard linting** but **too important to leave to prompt adherence**. Examples:

| Rule | Standard linter can catch? | Currently enforced? |
|------|--------------------------|-------------------|
| `allow_any_instance_of` in tests | Maybe (rubocop-rspec?) | No — prompt only |
| Self-mocking (mocking the object under test) | No | Prompt only |
| Missing assertions in test (`it` block with no `expect`) | Maybe | No |
| Deprecated library import (project-specific) | Yes (custom ESLint rules) | Yes — 4 custom rules |
| `console.log` left in production code | Yes (`no-console`) | Yes — warn |
| N+1 query introduced by new association usage | No (runtime only?) | No |
| Dead code / unused method after refactor | Partially | No |
| Type error in TypeScript | Yes (`tsc`) | No — not in hooks |
| Cross-file breakage (removed method still called elsewhere) | Yes (`tsc`, sort of) | No — not in hooks |
| Test that passes but tests nothing meaningful | No | No |

### What We Already Have

**Hook architecture (PostToolUse, fires on every file Claude writes):**
```
Claude writes file → PostToolUse hook fires →
  Ruby: rubocop -A (auto-fix) → rubocop (check) → block on violations
  JS/TS/Vue: eslint --fix → eslint (check) → block on violations
  Ruby source: find spec → block if spec missing
```

**Custom enforcement already built:**
- 2 custom RuboCop cops (trivial: remove redundant `require`, remove `RSpec.` prefix)
- 4 custom ESLint local rules (block deprecated library usage: BootstrapVue components/imports, Pinia stores/imports)
- Spec-required hook with `.specignore` for exceptions
- Behavioral mandate injection every turn (prompt-level, not deterministic)

**Standard linting already configured:**
- RuboCop: style, metrics (AbcSize, CyclomaticComplexity, MethodLength, ClassLength), performance, Rails, RSpec
- ESLint: TypeScript strict, Vue essential, Prettier, JSDoc
- Both run inside the write-hook pipeline with auto-correct → check → block

### Project Stack
- Ruby on Rails 7.2, Ruby 3.3.6, PostgreSQL, Sidekiq
- Vue 2.7 → 3 migration, TypeScript, Composition API
- RSpec (fixture_builder, WebMock, VCR), Vitest
- Docker Compose dev, GitHub Actions CI
- Claude Code as primary dev tool

## Key Questions

### Q1: Custom Deterministic Rules for LLM Anti-Patterns (Highest Priority)

LLMs produce specific anti-patterns that human developers rarely do. What deterministic rules are people writing to catch them?

1. **Test quality enforcement:**
   - Is there a rubocop-rspec cop for `allow_any_instance_of`? What about excessive `instance_double` usage? Empty `it` blocks?
   - Can Semgrep patterns catch self-mocking (mocking methods on the subject under test)?
   - Are there ESLint rules or Semgrep patterns for Vitest/Jest anti-patterns (e.g., `expect` never called in a test)?

2. **Deprecated usage enforcement beyond imports:**
   - Pattern-based blocking of specific method calls (e.g., `Rails.env.production?`, raw SQL strings, `update_columns` without justification)
   - Blocking specific ActiveRecord patterns (e.g., `default_scope` creation, `unscoped` usage)
   - How are teams using Semgrep vs. custom lint rules for this? What's easier to write and maintain?

3. **LLM-specific code generation patterns:**
   - Placeholder/TODO detection that's more aggressive than `no-warning-comments` (catching "implement this", "add error handling here", "this should be replaced with")
   - Detecting over-abstraction (single-use helper methods, unnecessary service objects wrapping one method call)
   - Catching LLM "politeness code" (unnecessary nil checks on values that can't be nil, redundant `.present?` guards on required fields)

4. **What's the authoring experience like?** For each tool:
   - How hard is it to write a custom rule? (Show me the simplest non-trivial example)
   - How fast does it run? (Per-file latency)
   - How does it integrate with Claude Code hooks?

### Q2: Semgrep as a Deterministic Layer

Semgrep seems like the most promising "beyond standard linting" tool. Deep dive:

1. **What can Semgrep catch that RuboCop/ESLint can't?**
   - Cross-file patterns?
   - Semantic patterns (not just syntax)?
   - Multi-language patterns in a single rule?

2. **Community rule sets for Ruby/Rails:**
   - What exists in the Semgrep registry for Ruby?
   - Quality and coverage — are they actually useful or toy examples?
   - Anyone maintaining a serious Rails Semgrep ruleset?

3. **Custom Semgrep rules for LLM guardrails:**
   - How expressive is the pattern language? Can it catch things like "mocking a method that's defined on the same class"?
   - Performance: fast enough for per-file PostToolUse hooks?
   - Can it run on a single file or does it need project context?

4. **Semgrep vs. custom lint rules (RuboCop cops / ESLint rules):**
   - When should you write a Semgrep pattern vs. a native lint rule?
   - Maintenance burden comparison
   - Performance comparison

### Q3: Type-Level Enforcement

Types are deterministic contracts. How aggressively can we use them?

1. **TypeScript `tsc --noEmit` as a hook gate:**
   - Is anyone doing this? Performance on a ~500-file TS project?
   - `--incremental` — does it make this feasible for per-edit gating?
   - What's the interaction with Vue SFC type checking?

2. **Stricter TypeScript configurations for LLM code:**
   - `strict: true` sub-options worth enabling beyond defaults?
   - `noUncheckedIndexedAccess`, `exactOptionalPropertyTypes` — do these catch LLM bugs?
   - Are there TypeScript ESLint rules that approximate type checking without running `tsc`?

3. **Ruby type checking (Sorbet / RBS):**
   - Is anyone using Sorbet or RBS as a deterministic gate for LLM-generated Ruby code?
   - Practical to adopt incrementally on an existing Rails 7.2 app?
   - Is this worth the investment or is RuboCop + good tests sufficient?

### Q4: Cross-File and Project-Wide Deterministic Analysis

Single-file linting misses breakage across files. What's available?

1. **Dead code detection:**
   - Tools that find unused methods, classes, constants across a Ruby project
   - TypeScript `--noUnusedLocals` / `--noUnusedParameters` as a subset
   - Anything that runs fast enough for interactive gating, or is this CI-only?

2. **Dependency graph analysis:**
   - If Claude changes file A, what other files could break?
   - Tools that can answer this question quickly
   - Is this even practical at the hook level, or is it a CI concern?

3. **Migration safety (strong_migrations):**
   - Does this work as a deterministic gate in a hook?
   - What dangerous migration patterns does it catch?
   - Performance for single-migration analysis?

### Q5: Hook Architecture — What's State of the Art?

1. **Published Claude Code hook architectures:**
   - Are there repos with sophisticated `.claude/hooks/` setups?
   - What patterns go beyond "run linter on file write"?
   - PreToolUse vs. PostToolUse — what checks belong where?

2. **Multi-tool orchestration in hooks:**
   - Running multiple analysis tools in parallel within a single hook
   - Aggregating results from multiple tools into a single pass/fail
   - Tiered severity (auto-fix / warn-in-prompt / block)

3. **"Write-check-rewrite" loops:**
   - Does anyone have hooks that detect violations, feed them back to the LLM, and block until fixed — as an automated loop rather than requiring Claude to see the error and self-correct?
   - Is this built into any Claude Code hook pattern?

4. **Diff-only analysis:**
   - Running analysis only on the lines Claude actually changed (not the whole file)
   - Avoiding false positives from pre-existing violations in files Claude touched
   - How do teams handle the `.rubocop_todo.yml` problem (grandfathered violations triggering on files the LLM edits)?

## What Would Change Our Approach

- **Semgrep has strong Ruby patterns and runs fast** → Add as a second deterministic layer alongside RuboCop, targeting semantic patterns RuboCop can't express
- **`tsc --noEmit --incremental` is fast enough for hooks** → Add TypeScript type checking to our PostToolUse pipeline immediately
- **Custom rubocop-rspec cops exist for mocking anti-patterns** → Enable them — zero custom code needed
- **Published hook architectures exist with multi-tool orchestration** → Adopt the pattern rather than inventing our own
- **Diff-only analysis is practical** → Solve the grandfathered-violations problem that currently causes friction
- **Strong_migrations works as a hook gate** → Add migration safety checking to the pipeline
- **Nobody is doing this beyond standard linting** → We're on the frontier. Document what we build and share it.

## Known Assumptions to Challenge

1. **"Custom lint rules are hard to write and maintain"** — Maybe tools like Semgrep make this trivial. What's the actual effort for a typical rule?
2. **"Cross-file analysis is too slow for interactive gating"** — Incremental analysis and caching may make this feasible. What's the actual latency?
3. **"Type checking is a CI concern, not a hook concern"** — `tsc --incremental` might be fast enough. Has anyone measured this in a Claude Code context?
4. **"More deterministic checks = more friction"** — Maybe, but the alternative is bugs shipping. Is there a sweet spot between aggressiveness and developer velocity?
5. **"Standard linters already catch the important stuff"** — Our experience says no. The gap between "style/syntax" and "semantic correctness" is where LLM bugs live.

## Output Format

**Per question area:**
- Specific tools/techniques with evidence of use
- Performance characteristics (latency, resource usage)
- Integration difficulty (how hard to add to our hook pipeline)
- Actual rule/pattern examples (not just "this tool exists")

**Synthesis:**
1. **Immediate wins** — things to add this week, minimal effort
2. **Custom rules to write** — specific patterns, estimated effort, expected ROI
3. **Tools to evaluate** — hands-on trial needed, what to test
4. **Architecture changes** — hook pipeline modifications
5. **Not worth it** — things that sound good but don't deliver (with evidence)

## Recency Requirements

**Critical (2024-2026 only):**
- Claude Code hook patterns, AI coding tool integrations
- Semgrep community rule sets and performance data
- New-generation tools (oxlint, Biome)

**Prefer recent, older acceptable:**
- Custom lint rule development (RuboCop cops, ESLint rules)
- Type checking strategies
- Cross-file analysis approaches

## Instructions for Researcher

### Footnote Format (Non-negotiable)

Use inline reference marks like `[^lint-1]` in the text body. Define each footnote at the bottom of its section (or document) with this exact format:

```
[^lint-1]: **Title or description** — [URL](url) (accessed YYYY-MM-DD). Tier N. "Exact quote if applicable."
```

Components:
- **Reference mark:** `[^kebab-case-prefix-N]` where prefix groups related sources (e.g., `[^semgrep-1]`, `[^hook-1]`, `[^tsc-1]`)
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
- **For tools**: Include latest version number and release date

### Search Strategy

1. **Semgrep registry:** Search for Ruby, Rails, RSpec rule sets. Evaluate quality, not just existence.
2. **Claude Code hooks:** GitHub search for `.claude/hooks/` directories, "PostToolUse" hook examples, published hook architectures
3. **Custom lint rules:** "custom rubocop cop testing", "custom eslint rule LLM", "semgrep custom rule rails"
4. **Type checking performance:** "tsc incremental performance", "vue tsc hook", "typescript type checking CI speed"
5. **LLM code quality research:** Academic papers or industry blog posts on deterministic analysis of LLM-generated code
6. **Strong_migrations:** Integration patterns, hook usage
7. **Dead code analysis:** Ruby unused method detection tools, comparison of approaches
8. **AI coding tool communities:** Cursor forums, Claude Code discussions — what deterministic checks are power users running?

### Research Process
- **Prioritize Q1 and Q2** — custom rules and Semgrep are the most likely high-ROI areas
- Go deep on promising finds rather than cataloging everything surface-level
- For any tool recommendation, include a concrete example of a rule/pattern — not just "this tool can do X"
- Performance data is critical — if you can't find performance numbers, flag it
- Be skeptical of marketing claims — we want proven practice, not feature lists
