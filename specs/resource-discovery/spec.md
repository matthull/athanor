# Resource Discovery & Skill Visibility for Athanor Agents

**Status:** Draft
**Created:** 2026-03-25
**Project:** athanor

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[B:session-2026-03-24]` | From artifex shaping session 2026-03-24 |
| `[B:seismic-observation]` | Triggered by seismic-classifier-mapping athanor observations |
| `[E:CLAUDE.md]` | Extends existing musashi CLAUDE.md patterns |
| `[E:CLAUDE.local.md]` | Extends existing musashi CLAUDE.local.md |
| `[D:reason]` | Design decision with rationale |

---

## Overview

Athanor agents operating in musashi have access to a rich ecosystem of skills, documentation, and operational context — but can't reliably find or use most of it. The seismic-classifier-mapping athanor produced excellent technical work but missed project management concerns (UX complexity risk, status snapshots, witness notifications) because agents didn't load skills like `/nfr-product-lens`, `/project-management`, or `/qa-plan` that would have surfaced those concerns naturally. `[B:seismic-observation]`

This spec addresses three gaps:

1. **Skill discoverability** — athanor agents don't know what skills exist or when to use them
2. **Risk-oriented development principles** — the artifex's project management philosophy isn't available to agents in a distilled form
3. **Skill usage tracking** — no visibility into which skills agents actually load, making it impossible to detect gaps

All changes live in the environment (CLAUDE.local.md, hooks, scripts) — not in the athanor itself. This follows the adaptability-through-essentialism principle: the athanor reads the environment; enriching the environment enriches the athanor. `[D:environment-not-athanor]`

---

## Domain Dictionary

**Skill** — A Claude Code slash command (e.g., `/pm-spec`, `/rca`) that injects domain-specific expertise into an agent's context. Skills live in `~/.claude/skills/` (global) and `.claude/skills/` (project-specific). Project skills override global skills of the same name.

**Skill discovery** — The process by which an agent identifies which skills are relevant to its current task. Currently relies on `/skill-discovery` command or the routing table in CLAUDE.md. Both require the agent to actively seek them out.

**Resource** — Broader than skills: includes meeting transcripts, specs, Linear projects, Slack channels, prd-lens docs, process docs — anything an agent might need to do its work well.

**Witness concerns** — Already defined in CLAUDE.local.md. The stated needs of each team member from any work produced. Drives what "abundantly satisfied" means for each project.

---

## Problem Statement

### Current State

An athanor azer working on a musashi project:

1. Reads its opus (intent + context)
2. Reads AGENTS.md, azer.md (athanor protocol)
3. Reads CLAUDE.md (project conventions — includes routing table and skill catalog)
4. May or may not run `/skill-discovery`
5. Loads skills it thinks are relevant based on the opus content

### Gaps

**Skill selection is ad-hoc.** The CLAUDE.md routing table is organized by workflow ("When you are designing a feature...") and domain ("When you are writing Vue..."). But athanor agents often need *cross-cutting* skills — an assessment azer needs `/project-management` + `/nfr-product-lens` + `/qa-plan` simultaneously. The routing table doesn't surface these combinations. `[B:session-2026-03-24]`

**Personal skills are invisible.** The artifex has 56 global skills, many highly relevant to athanor work (PM, product, risk). These aren't in CLAUDE.md's routing table (which covers project-specific skills). An athanor agent has no way to discover `/nfr-product-lens` or `/pm-spec` without running `/skill-discovery` — which it may not think to do for assessment work. `[B:session-2026-03-24]`

**Risk-oriented development philosophy is undocumented.** The artifex's approach (unknowns first, UX complexity as project risk, adoption-first thinking, NFRs as product requirements) exists in `specs/project-lifecycle/process.md` and scattered prd-lens docs, but isn't distilled for agent consumption. Agents produce technically excellent work that misses project-level concerns because the philosophy isn't in their context. `[B:seismic-observation]`

**No skill usage tracking.** When an athanor agent completes an opus, there's no record of which skills it loaded. This means:
- Can't detect skill gaps (agent didn't load a skill it should have)
- Can't assess skill effectiveness (skill was loaded but didn't help)
- Can't improve the routing/discovery system based on actual usage `[B:session-2026-03-24]`

---

## Deliverables

### D1: Skill Directory in CLAUDE.local.md

Add a curated skill directory to CLAUDE.local.md organized by *purpose* (not by workflow or domain). This makes personal skills visible to all agents working in musashi, including athanor agents.

**Intent:** An athanor agent doing assessment work sees "Risk detection & project management" and naturally loads `/nfr-product-lens` or `/project-management`. `[B:session-2026-03-24]`

**Structure:** Group by purpose with brief descriptions:
- Understanding stakeholder goals & shaping requirements
- Risk detection & project management
- Executing work
- Communication & context

**Scope:** Personal/global skills only — project-specific skills are already covered by CLAUDE.md's routing table and skill catalog. `[D:no-duplication]`

**Boundary:**
- Agent: Draft the directory based on research from this session
- Operator: Review and adjust emphasis/grouping

---

### D2: Risk-Oriented Development Principles in CLAUDE.local.md

Add a distilled section covering the artifex's project management philosophy. Not duplicating `process.md` — distilling the *detection heuristics* an agent can apply.

**Intent:** Agents internalize patterns like "UX complexity is a project risk, flag it early" without needing to read the full process doc. `[B:seismic-observation]`

**Source material:**
- `specs/project-lifecycle/process.md` — adoption-first lifecycle, risk-ordered work `[E:process.md]`
- prd-lens docs (seismic-classifier-mapping, proof-recommender-slack, proof-recommender-salesforce) — risk-first sequencing patterns `[E:prd-lens]`
- `/nfr-product-lens` skill — NFRs as product requirements pattern
- `/project-management` skill — certainty breakdown, gap analysis

**Key principles to distill:**
- Unknowns before knowns (sequence by uncertainty)
- UX complexity is a project risk, not a design task
- Adoption-first (done = customers using it, not code shipped)
- NFRs are product requirements (security, performance, observability, cost)
- Stalls are signals (patterns, not instances)

**Boundary:**
- Agent: Draft the section
- Operator: Review — this is the artifex's philosophy, needs to feel right

---

### D3: Skill Usage Tracking

Record which skills each agent session loads, so the artifex can detect gaps and improve discovery.

**Intent:** After the seismic-classifier-mapping athanor runs, the artifex can see "assessment azer loaded `/research` but not `/project-management` or `/nfr-product-lens`" and know why witness concerns weren't surfaced. `[B:session-2026-03-24]`

**Approach options (needs investigation):**

1. **Hook-based tracking** — A PostToolUse hook on the Skill tool that logs skill invocations to a session file or the opus discharge record. Automatic, no agent cooperation needed.

2. **Discharge convention** — Add `skills_loaded` to opus discharge frontmatter. Requires agent cooperation (soft enforcement). Simpler infrastructure.

3. **Both** — Hook captures actual usage; discharge convention captures agent's self-report. Delta between them is interesting signal.

**Key questions:**
- Can a PostToolUse hook on the Skill tool capture which skill was invoked? (Need to check hook input schema)
- Where should usage data live? Per-session log? Per-opus discharge? Centralized file?
- Is this athanor-specific or global (all Claude Code sessions)?

**Boundary:**
- Agent: Investigate hook feasibility, propose approach
- Operator: Decide scope (athanor-only vs global) and storage location

---

### D4: Assessment Skill Loading Guidance

Ensure assessment opera specifically prompt skill loading for cross-cutting concerns. This may be a small addition to AGENTS.md's assessment section or a note in the assessment opus template.

**Intent:** An assessment azer reads the assessment opus template and sees explicit guidance to load risk/PM/product skills, not just technical ones. `[B:session-2026-03-24]`

**Approach:** Minimal — a Context line in the assessment opus template pointing at the skill directory and naming the key assessment-relevant skills. Follows the "pour into Context" principle from opus.md.

**Boundary:**
- Agent: Draft the addition
- Operator: Review for alignment with athanor essentialism

---

## Task Breakdown

| Task | Deliverable | Dependencies | Complexity |
|------|-------------|-------------|------------|
| T1: Skill directory | D1 | None | Low — content is drafted, needs final edit |
| T2: Risk principles | D2 | None | Low — distill from known sources |
| T3: Skill tracking investigation | D3 | None | Medium — hook feasibility unknown |
| T4: Skill tracking implementation | D3 | T3 | Medium — depends on approach chosen |
| T5: Assessment skill guidance | D4 | D1 (needs skill directory to reference) | Low |

T1 and T2 are independent and ready now. T3 is investigation. T4 depends on T3. T5 depends on T1.

---

## Out of Scope

- Changes to shared athanor components beyond the assessment template Context line (D4)
- Automated skill recommendation (too complex for now — start with visibility)
- Skill effectiveness measurement (need usage data first)
- Changes to the musashi project CLAUDE.md routing table or skill catalog
- New skills — this spec is about making existing skills findable

---

## Test Scenarios

- An athanor assessment azer working on seismic-classifier-mapping loads `/project-management` and `/nfr-product-lens` during its assessment run
- After an opus discharge, the artifex can determine which skills were loaded during that opus
- A fresh athanor azer reading CLAUDE.local.md can find the skill directory and identify relevant skills for its task type without running `/skill-discovery`

---

## Retrospective

When complete, review:

### Documentation Updates
- [ ] Did the skill directory actually get used by athanor agents?
- [ ] Were the risk principles actionable or too abstract?
- [ ] Did agents discover skills they wouldn't have found otherwise?

### Workflow Improvements
- [ ] Should skill loading become part of the azer mise en place protocol?
- [ ] Does the assessment opus template need further Context enrichment?
- [ ] Should `/skill-discovery` be run automatically for athanor agents?

### Knowledge Capture
- [ ] Skill usage patterns from tracking data
- [ ] Skills that agents should have loaded but didn't (gap analysis)
