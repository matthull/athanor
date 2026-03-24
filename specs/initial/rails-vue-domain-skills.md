# Research Request: Rails/Vue Domain Skills for Claude Code

**Created:** 2026-03-08
**Status:** Open

## Goal

Identify and evaluate authoritative Claude Code skills (published on GitHub, Anthropic marketplace, or similar) for the **Ruby on Rails** and **Vue.js** ecosystems. We need concrete install-or-skip recommendations for each candidate, with overlap analysis against our existing custom skills.

**This is NOT a general software design principles research.** We've already researched generic refactoring/clean-code skills (elifiner/refactoring, ratacat/claude-skills, matsengrp/claude-code-agents, boxclawai/skills) and found them useful but insufficient. See "Prior Research" below. This request is specifically about **stack-specific domain skills** — Rails patterns, Vue patterns, testing patterns for these frameworks.

## Context

### Project Stack
- **Backend:** Ruby on Rails 7.2, PostgreSQL, Sidekiq, Elasticsearch
- **Frontend:** Vue 2.7 (migrating to Vue 3), Composition API (`<script setup lang="ts">`), Bootstrap 4 CSS, Storybook
- **Testing:** RSpec (fixtures via fixture_builder gem, VCR for API cassettes), Vitest (JS utilities only — components use Storybook)
- **Infrastructure:** Docker Compose development, GitHub Actions CI, Heroku deployment

### What We Already Have (Custom Skills)

We have custom Claude Code skills that encode project-specific conventions. The research should evaluate external skills **against these** — complement, replace, or extend?

**`ruby-rails` skill** — Progressive disclosure (5 resource files for models, controllers, services, serialization, migrations). Critical mandates:
- Fixture builder constraint (NEVER edit .yml fixtures, use fixture_builder.rb)
- One spec file per class
- Test every public method at its layer
- Verification loops (TDD → scoped → console)
- Anti-patterns: no allow_any_instance_of, no raw SQL, no Rails.env.production?

**`vue` skill** — Vue 2.7 Composition API patterns:
- `<script setup lang="ts">` mandatory
- v-model compatibility layer (Vue 2 `value`/`@input` + Vue 3 `modelValue`/`@update:modelValue`)
- Deprecated library blocklist (vue-easy-toast, vue-select, vue-pagination-2, vue-clickaway, v-scroll-sync, vue-bootstrap)
- Storybook mandatory for all components (Vitest only for utilities)
- CSS design tokens, no hardcoded values, no !important
- Bootstrap 4 (not vue-bootstrap)

**`unit-testing` skill** — Framework-agnostic reliability principles (clock mocking, randomness isolation, external state, test independence, flakiness catalog)

**`api-integration` skill** — Documentation extraction before coding, VCR testing patterns, security requirements, error handling by status code

### Guidance Library (Being Deprecated → Skills)

We also have ~35 guidance files covering rails/, vue/, testing/, devops/, security/, etc. These are being migrated into skills. External skills that cover the same ground well could accelerate this migration. Key guidance areas:
- Rails: API patterns, backend testing, background jobs, CRUD controllers, fixture testing, models, safe migrations, service objects
- Vue: Code quality, component patterns, frontend testing, styling, Vue 3 migration rules
- Testing: API client patterns, integration QA, mocking guidelines (3 files), testing strategy, VCR/WebMock patterns

### Prior Research (Do Not Re-Research)

From a parallel project's implementation skills audit (2026-03-07), we already researched and evaluated these **generic design principles skills**:

| Skill | Verdict |
|-------|---------|
| `elifiner/refactoring` | High quality, Fowler-based. Install. Analysis-first workflow. |
| `ratacat/claude-skills` clean-code | Clean Code based. Useful but lacks decision rules for when NOT to apply. |
| `matsengrp/claude-code-agents` clean-code-reviewer | Best over-mocking treatment found. Worth cherry-picking. |
| `boxclawai/skills` code-smells-refactoring | Comprehensive Fowler smell catalog. Reference, not enforcement. |
| `obra/superpowers` | Community workflow skills. Evaluate against our pipeline, don't install directly. |

**Skip these categories entirely — already covered:**
- Generic refactoring / clean code skills
- Generic software design principles
- Workflow/pipeline orchestration skills
- General-purpose coding assistant skills

## Key Questions

### Q1: Rails-Specific Skills (Highest Priority)

1. Are there official or semi-official Rails team Claude Code skills? (DHH, Rails core team, 37signals/Basecamp)
2. Are there high-quality community Rails skills from known consultancies (thoughtbot, Evil Martians, Shopify, GitHub)?
3. For each candidate found:
   - What conventions does it encode? (List the major rules/mandates)
   - How does it overlap with our `ruby-rails` skill? (Be specific — which mandates match, which conflict, which are additive)
   - Does it cover areas our skill doesn't? (e.g., ActiveJob patterns, ActionCable, Turbo/Hotwire, caching strategies, Rails 7.2-specific patterns)
   - Install recommendation: global install, project-only, cherry-pick specific resources, or skip

### Q2: Vue-Specific Skills

1. Are there official Vue team Claude Code skills? (Evan You, Vue core team)
2. Callstack or similar consultancy skills for Vue? (Note: Callstack's `callstackincubator/agent-skills` covers React Native, not Vue — check if they have Vue-specific ones)
3. Community Vue skills with quality comparable to the React Native ones we found?
4. **Critical filter:** Must support Vue 2.7 Composition API or be version-agnostic. Pure Vue 3 skills with Composition API patterns may still be useful if the patterns transfer. Skills that assume Options API are not useful.
5. For each candidate: overlap analysis against our `vue` skill, same structure as Q1.

### Q3: Testing Skills (Rails/Vue Specific)

1. RSpec-specific skills beyond our `unit-testing` (which is framework-agnostic)?
   - Fixture patterns, factory_bot best practices, request spec conventions
   - VCR/WebMock patterns
   - RSpec performance optimization
2. Vitest or Vue Test Utils skills?
3. Storybook-specific skills? (Our `vue` skill mandates Storybook but doesn't deeply cover story patterns, interaction testing, visual regression)
4. **Specific gap to fill:** Over-mocking detection in RSpec context. `matsengrp/claude-code-agents` has the best generic treatment — is there an RSpec-specific equivalent?

### Q4: Adjacent Domain Skills

1. Sidekiq/ActiveJob skills? (Background job patterns, retry strategies, idempotency)
2. PostgreSQL skills beyond our existing `supabase-postgres-best-practices` and `pg-perf-power-user`? (Those are Supabase-oriented — we need vanilla Rails/PG patterns)
3. Elasticsearch skills? (We use ES extensively)
4. Bootstrap 4 CSS skills? (Not Bootstrap 5 — we're on 4)
5. GitHub Actions CI skills? (We have complex parallel test matrix)

### Q5: Notable CLAUDE.md Files and Patterns

1. Are there well-known Rails projects with published CLAUDE.md files worth studying? (Equivalent to Evan Bacon's minimal React Native CLAUDE.md that was valuable in the prior research)
2. Any Rails-ecosystem templates or starter repos with Claude Code configuration?

## What Would Change Our Approach

- **A high-quality official Rails skill exists** → We'd likely install it globally and make our `ruby-rails` skill a thin project extension (fixture_builder constraint, spec organization) rather than a full standalone skill.
- **A high-quality Vue 2/3 compatibility skill exists** → Could replace our v-model compatibility section and migration rules.
- **An RSpec skill with fixture/factory guidance exists** → Could replace several of our testing guidance files and accelerate the guidance→skills migration.
- **Nothing good exists for Rails/Vue** → Confirms our custom skills are the right investment. Focus shifts to enriching them with escalation protocols (Layer 2 STOP triggers) rather than looking for external sources.

## Known Assumptions to Challenge

1. **"Rails is too opinionated to need skills"** — Rails has strong conventions, but Claude still makes silent decisions about service object boundaries, migration strategies, and ActiveRecord patterns. Skills may still add value.
2. **"Vue 2.7 is too niche for published skills"** — Most Vue skills will target Vue 3. But Composition API patterns are largely the same — the differences (v-model, Teleport, Suspense) are enumerable.
3. **"Storybook skills would be React-focused"** — Storybook is framework-agnostic in principle. Check if any Storybook skills exist that aren't React-specific.
4. **"Official framework skills will be better than ours"** — Not necessarily. Our skills encode project-specific constraints (fixture_builder, no vue-bootstrap, Bootstrap 4 not 5). Official skills may be more general. The question is complementarity, not replacement.

## Output Format

**Comparison table per category** (Rails, Vue, Testing, Adjacent) with columns:
- Skill name and source
- Install method
- Key conventions it encodes (bullet list)
- Overlap with our existing skills (specific)
- Unique value-add (what it covers that we don't)
- Recommendation (install globally / project-only / cherry-pick / skip)
- Evidence tier (1-4, see below)

Plus a **summary recommendation** section with:
- Tier 1 installs (do immediately)
- Tier 2 evaluations (install in sandbox, test before committing)
- Gaps confirmed (nothing exists, need custom skill)
- Guidance files that could be retired if external skill is installed

## Recency Requirements

**Critical (2025-2026 only)** — Claude Code skills ecosystem is evolving rapidly. Skills from 2024 may use outdated formats or conventions. The Anthropic marketplace and GitHub skill repos are actively changing. Older content about Claude Code configuration is likely outdated.

Exception: foundational Rails/Vue convention resources (Rails Guides, Vue docs) are stable and acceptable regardless of age.

## Instructions for Researcher

### Footnote Format (Non-negotiable)

Use inline reference marks like `[^rails-1]` in the text body. Define each footnote at the bottom of its section (or document) with this exact format:

```
[^rails-1]: **Title or description** — [URL](url) (accessed YYYY-MM-DD). Tier N. "Exact quote if applicable."
```

Components:
- **Reference mark:** `[^kebab-case-prefix-N]` where prefix groups related sources
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

### Search Strategy Suggestions

1. **GitHub search:** `claude code skill` + framework name, `SKILL.md` in repos, `claude.md` in major Rails/Vue repos
2. **Anthropic marketplace:** Check official skill marketplace/registry if one exists
3. **Framework team repos:** Check rails/rails, vuejs/core, vuejs/vue for any CLAUDE.md or .claude/ directories
4. **Consultancy repos:** thoughtbot, Evil Martians, Shopify, Callstack, Vercel Labs — check for Claude Code skills
5. **Community aggregators:** awesome-claude-code lists, HN/Reddit threads about Claude Code skills

### Transparency Requirements
- Flag `[UNCERTAIN]` when no authoritative source exists
- Flag `[CONFLICTING]` when sources disagree (present both sides)
- Flag `[SINGLE SOURCE]` when a claim relies on only one reference
- Flag `[OUTDATED]` when newest source is >2 years old on fast-moving topics

### Research Process
- Push back if any questions seem misguided or if there's a better framing
- Distinguish between well-established findings and emerging/contested ideas
- Be direct about limitations or gaps in available evidence
- If the Claude Code skills ecosystem for Rails/Vue is genuinely sparse, say so clearly — a "nothing good exists" finding is valuable
