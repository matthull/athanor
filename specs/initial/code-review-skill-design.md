# Code Review Skill Design

**Date:** 2026-03-14
**Context:** Upgrading from basic `/code-review` command to a proper skill with skill discovery integration, Gastown-inspired review legs, and project-aware domain loading.

---

## Current State

`/code-review` is a **command** (`~/.claude/commands/code-review.md`) that launches up to 6 parallel `code-review-expert` subagents. It has smart file-type-based agent selection but:

- No skill discovery (doesn't load domain skills like `ruby-on-rails`, `vue`)
- Generic prompts (not surgical "Look for" checklists)
- Missing review dimensions (no wiring, resilience, or smells legs)
- No presets (always same agents, no gate/full/contextual tiers)
- Can't be composed into `/orchestrate` as cleanly as a skill

Also relevant:
- `/task-review` command — reviews code against handoff checklist (good, keep separate)
- `/pr-review` command (project) — Musashi-specific PR review (excellent, project-specific — should be loaded as domain context, not replaced)
- `/review` command — delegates to Gemini (good, keep as second-opinion tool)
- `/spec-check` skill — spec compliance (cleanly scoped, not code quality)

## Design Goals

1. **Skill, not command** — full skill with SKILL.md + resources/ for review criteria
2. **Skill discovery integrated** — auto-load domain skills based on changed files
3. **Gastown-quality leg prompts** — concrete "Look for" checklists + "Questions to answer"
4. **Presets** — `gate` (fast), `full` (comprehensive), `contextual` (auto-select based on diff)
5. **Project domain injection** — load project-specific review criteria (from `pr-review.md`, CLAUDE.md conventions)
6. **Composable** — usable standalone AND as part of `/orchestrate` VERIFY phase

## Architecture

```
/code-review [target] [--preset=gate|full|contextual]
  │
  ├── Step 1: Parse target (PR, branch, staged changes, file glob)
  ├── Step 2: Analyze diff (file types, change scope, risk signals)
  ├── Step 3: Skill discovery (catalog scan → load domain skills)
  ├── Step 4: Select legs (based on preset + file types + risk signals)
  ├── Step 5: Compose leg prompts (base template + domain guidance injection)
  ├── Step 6: Launch parallel review agents (one per leg)
  ├── Step 7: Synthesize findings (deduplicate, cross-reference, prioritize)
  └── Step 8: Output report
```

## Review Legs

### Analysis Legs (code understanding)

#### correctness
**Focus:** Logical correctness and edge case handling

**Look for:**
- Logic errors and bugs
- Off-by-one errors
- Null/nil/undefined handling gaps
- Unhandled edge cases (empty arrays, zero values, missing keys)
- Race conditions in concurrent code
- Dead code or unreachable branches
- Incorrect assumptions in comments vs actual code behavior
- Integer overflow/underflow potential
- Incorrect boolean logic (De Morgan's law violations, precedence errors)
- State mutation where immutability was expected

**Questions to answer:**
- Does the code do what it claims to do?
- What inputs could cause unexpected behavior?
- Are all code paths tested or obviously correct?
- What would break if called with unexpected arguments?

#### performance
**Focus:** Performance bottlenecks and efficiency

**Look for:**
- O(n^2) or worse algorithms where O(n) is possible
- N+1 query patterns (ActiveRecord eager loading missing)
- Unnecessary allocations in loops or hot paths
- Missing database indexes for query patterns
- Blocking operations where async is expected
- Memory leaks or unbounded growth (arrays that only grow)
- Excessive string concatenation in loops
- Missing pagination for potentially large datasets
- Redundant database queries (same data fetched multiple times)
- Heavy computation in request cycle that should be background job

**Questions to answer:**
- What happens at 10x, 100x current scale?
- Are there obvious optimizations being missed?
- Is the N+1 query risk mitigated?
- Should any of this be a background job?

**Domain injection:** When `ruby-on-rails` loaded, add: "Check for includes/preload on associations. Check for counter_cache opportunities. Check for unscoped queries on large tables."

#### security
**Focus:** Security vulnerabilities and attack surface

**Look for:**
- Input validation gaps (params not validated server-side)
- Authentication/authorization bypasses (missing `before_action :authenticate_user!`)
- Injection vulnerabilities (SQL via raw queries, XSS via unescaped output, command injection)
- Sensitive data exposure (secrets in logs, PII in error responses, credentials in code)
- Hardcoded secrets or credentials
- Mass assignment vulnerabilities (missing strong_params, using `permit!`)
- Path traversal in file operations
- CSRF protection gaps
- Insecure direct object references (accessing records without ownership check)
- Missing rate limiting on sensitive endpoints

**Questions to answer:**
- What can a malicious user do with this code?
- What data could be exposed if this fails?
- Are there defense-in-depth gaps?
- Is authorization checked server-side, not just client-side?

**Domain injection:** When `ruby-on-rails` loaded, add: "Check strong_params usage — never `permit!`. Check `before_action` chain includes authentication. Check ActiveRecord queries for SQL injection via string interpolation. Check jbuilder responses don't leak internal fields."

#### elegance
**Focus:** Design clarity and abstraction quality

**Look for:**
- Unclear abstractions or misleading naming
- Functions doing too many things (violating Single Responsibility)
- Missing or over-engineered abstractions
- Coupling that should be loose (component knowing about parent internals)
- Dependencies flowing wrong direction (model depending on controller)
- Magic numbers/strings without named constants
- Inconsistent patterns within the same module
- Reinventing existing utilities (dayjs, lodash, ActiveSupport methods)
- God classes/modules accumulating responsibilities

**Questions to answer:**
- Would a new team member understand this?
- Does the structure match the problem domain?
- Is the complexity justified by the requirements?
- Are abstractions at the right level?

#### resilience
**Focus:** Error handling and failure modes

**Look for:**
- Swallowed errors or empty rescue/catch blocks
- Missing error propagation (errors caught but not reported)
- Unclear or generic error messages ("Something went wrong")
- Missing retry/backoff logic for external service calls
- Missing timeout handling on HTTP requests
- Resource cleanup on failure (database connections, file handles)
- Partial failure states (half-completed multi-step operations)
- Missing fallback behavior for degraded services
- Unhelpful error logging (no context, no stack trace)
- Error handling that masks the original cause

**Questions to answer:**
- What happens when external services fail?
- Can the system recover from partial failures?
- Are errors actionable for operators (can they diagnose from the error message)?
- Are resources properly cleaned up in error paths?

**Domain injection:** When `ruby-on-rails` loaded, add: "Check for `rescue => e` without re-raising or reporting. Check for `reportError` usage from errorHandler.ts in Vue components. Check for transaction safety on multi-model operations."

#### style
**Focus:** Convention compliance and consistency

**Look for:**
- Naming convention violations (snake_case Ruby, camelCase JS)
- Formatting inconsistencies
- Import organization issues (relative paths instead of TypeScript aliases)
- Comment quality (obvious comments, outdated comments)
- Missing YARD docs on new classes/methods
- Log message quality and appropriate levels
- Inconsistent patterns within the codebase (new code diverging from conventions)

**Questions to answer:**
- Does this match the rest of the codebase?
- Would the linter approve? (If not, why wasn't it caught?)
- Is the code self-documenting where possible?

**Domain injection:** When project CLAUDE.md loaded, add: "Check TypeScript aliases vs relative paths. Check CSS for !important usage. Check Vue components for `scoped` directive. Check for hardcoded colors instead of sass variables."

#### smells
**Focus:** Anti-patterns and technical debt

**Look for:**
- Long methods (>50 lines is suspicious, >100 is a smell)
- Deep nesting (>3 levels of conditionals/loops)
- Shotgun surgery patterns (changing one thing requires touching many files)
- Feature envy (method uses more of another class's data than its own)
- Data clumps (same group of params passed together repeatedly)
- Primitive obsession (using strings/numbers where value objects would clarify)
- Copy-paste code (DRY violations without justification)
- TODO/FIXME accumulation (adding without timeline)
- Boolean parameters (method behavior switching on a flag)
- Long parameter lists (>4 params suggests missing object)

**Questions to answer:**
- What will cause pain during the next change?
- What would you refactor if you owned this code?
- Is technical debt being added or paid down?

### Verification Legs (implementation quality)

#### wiring
**Focus:** Installed-but-not-wired gaps

**Look for:**
- New dependency in Gemfile/package.json but never imported/required in code
- SDK/library added but old implementation still in use (added dayjs but still using moment)
- Config/env var defined but never loaded or referenced in code
- New route defined but no controller action implemented
- Migration adds column but no model attribute usage
- Feature flag added to flipper_features.yml but never checked in code
- CSS class defined but never applied in templates
- New component file created but never imported/rendered

**Questions to answer:**
- Is every new dependency actually used?
- Are there old patterns that should have been replaced by the new dependency?
- Is there dead config that suggests incomplete migration?
- Does every new route have a working endpoint?

**Domain injection:** When `ruby-on-rails` loaded, add: "Check db/seed_data/flipper_features.yml for new features — are they checked with `Flipper.enabled?` in code? Check Gemfile additions against actual `require` statements."

#### commit-discipline
**Focus:** Commit quality and atomicity

**Look for:**
- Giant "WIP" or "fix" commits mixing multiple unrelated changes
- Poor commit messages ("stuff", "update", "asdf", "fix things")
- Unatomic commits (feature + refactor + bugfix in same commit)
- Commits that touch 20+ files across different features
- Missing context in commit messages (no "why" explanation)
- Debug code committed (binding.pry, console.log, debugger)

**Questions to answer:**
- Could this history be bisected effectively?
- Would a reviewer understand the progression?
- Are commits atomic (one logical change each)?

#### test-quality
**Focus:** Test meaningfulness, not just coverage

**Look for:**
- Weak assertions (only checking `!= nil`, `expect(response).to be_successful` without body check)
- Missing negative test cases (happy path only, no error cases)
- Tests that can't fail (mocked so heavily the test is meaningless)
- Testing implementation details instead of behavior
- Missing boundary testing (edge values, empty inputs, max values)
- Flaky test indicators (sleep, Time.now, database ordering without ORDER BY)
- Overly specific assertions that break on cosmetic changes
- Missing test isolation (test depends on state from previous test)
- Using `FactoryBot.create` when fixtures would suffice (per project conventions)

**Questions to answer:**
- Do these tests actually verify behavior?
- Would a bug in the implementation cause a test failure?
- Are edge cases and error paths tested?
- Are tests following project conventions (fixtures for units, factories for edge cases)?

**Domain injection:** When `unit-testing` and `ruby-on-rails` loaded, add: "Check for `FactoryBot.create` outside request specs (should use fixtures). Check test structure follows Arrange-Act-Assert with newline separations. Check for `travel_to` or `freeze_time` when testing time-dependent behavior."

---

## Presets

### `gate` — Fast, blocker-focused (for automated flow)
**Legs:** wiring, security, smells, test-quality
**When to use:** Automated pre-merge gate, quick PR checks, CI integration
**Expected time:** ~30 seconds

### `full` — Comprehensive (for major features)
**Legs:** All 10 legs
**When to use:** Major feature PRs, architectural changes, security-sensitive code
**Expected time:** ~2-3 minutes

### `security` — Security-focused (for sensitive changes)
**Legs:** security, resilience, correctness, wiring
**When to use:** Auth changes, payment code, data handling, API endpoints
**Expected time:** ~1 minute

### `refactor` — Quality-focused (for cleanup work)
**Legs:** elegance, smells, style, commit-discipline
**When to use:** Refactoring PRs, tech debt cleanup
**Expected time:** ~1 minute

### `contextual` (default) — Auto-select based on diff
**Logic:**
```
Always include: correctness, security, test-quality (baseline)
If Rails files changed: + performance, wiring
If Vue files changed: + style, elegance
If >500 lines changed: + smells, commit-discipline
If migration files: + wiring, performance
If config/auth files: + security (upgraded to full security leg)
If test files only: test-quality only
If docs only: skip (or just style)
```

---

## Skill Discovery Integration

### Step 3: Discover and Load Domain Skills

```
1. Get changed file paths from diff
2. Scan skill catalog in system context
3. For each skill, evaluate relevance:

Skill discovery:
- ruby-on-rails: YES — app/models/ and app/controllers/ files changed
- vue: YES — app/javascript/src/components/ files changed
- api-integration: NO — no external API changes
- unit-testing: YES — spec/ files changed
- supabase-sql: NO — no database policy work
- react-native-expo: NO — not a mobile project

4. Load selected skills via Skill tool
5. Extract domain-specific review criteria from loaded skills
6. Inject into leg prompts as "Domain injection" sections
```

### Domain Injection Template

Each leg prompt gets appended with:
```markdown
## Project-Specific Review Criteria

The following domain guidance is loaded for this review:

### From ruby-on-rails skill:
{extracted relevant review criteria}

### From vue skill:
{extracted relevant review criteria}

### From project CLAUDE.md:
{extracted coding conventions relevant to this leg}
```

---

## Synthesis Step

After all legs complete, synthesize:

```markdown
# Code Review: {target}

## Executive Summary
(Overall assessment, merge recommendation: APPROVE / REQUEST CHANGES / NEEDS DISCUSSION)

## Critical Issues (P0 — Must fix before merge)
### [{Type}] {Issue title}
- **File:** {path:line}
- **Found by:** {which legs — higher confidence if multiple}
- **Impact:** {what breaks or is at risk}
- **Fix:** {specific, actionable recommendation}

## Major Issues (P1 — Should fix before merge)
...

## Minor Issues (P2 — Nice to fix)
...

## Wiring Check
(Dependencies added but not wired, configs defined but not used)

## Test Quality Assessment
(Assertion strength, negative cases, edge coverage)

## Strengths
(What's done well — important for balanced feedback)

## Observations
(Non-blocking notes, patterns noticed, suggestions for future)
```

### Deduplication Rules
- Same issue found by multiple legs → merge into single entry, note all legs that found it
- Issues found by 3+ legs → auto-promote to next severity level
- Conflicting recommendations between legs → flag as "Needs discussion" with both perspectives

---

## Integration Points

### Standalone Usage
```
/code-review                              # contextual preset, staged/uncommitted changes
/code-review --preset=full                # all 10 legs
/code-review --preset=gate                # fast, blockers only
/code-review --pr=123                     # review a specific PR
/code-review src/components/             # review specific directory
```

### In /orchestrate VERIFY Phase
- Replaces ad-hoc code-review-expert agent spawning
- Preset selected by orchestrator based on task characteristics
- Results feed into `verify-results.md` alongside spec-check and qa-plan

### With /gemini-review (Second Opinion)
- Run `/code-review` first (Claude-based, parallel, structured)
- If high-risk or uncertain → `/review` sends to Gemini for second opinion with full context
- Compare findings for high-confidence assessment

---

## Migration Path

### Phase 1: Create the Skill
1. Create `~/.claude/skills/code-review/SKILL.md` with the architecture above
2. Create `~/.claude/skills/code-review/resources/` with:
   - `leg-prompts.md` — all 10 leg prompt templates
   - `presets.md` — preset definitions
   - `synthesis-template.md` — synthesis output template
3. Keep existing `/code-review` command as backwards-compatible alias

### Phase 2: Skill Discovery Integration
1. Add Step 3 (skill discovery) to the skill
2. Add domain injection to leg prompt templates
3. Test with a real PR: verify domain-specific findings appear

### Phase 3: Gastown-Quality Prompts
1. Replace generic leg descriptions with "Look for" + "Questions to answer" format
2. Add project-specific items from CLAUDE.md conventions
3. Test: compare review quality before/after on same PR

### Phase 4: Presets and Contextual Selection
1. Implement preset system
2. Implement contextual auto-selection logic
3. Default to `contextual` preset

### Phase 5: Orchestrate Integration
1. Wire into `/orchestrate` VERIFY phase as composable skill
2. Orchestrator selects preset based on task characteristics
3. Results flow into verify-results.md
