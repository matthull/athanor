# Gastown Formulas Audit: Workflow Definitions Analysis

**Date:** 2026-03-14
**Source:** ~/code/gastown/internal/formula/formulas/ (43 TOML files) + internal/formula/*.go (execution model)
**Purpose:** Identify workflow patterns worth adopting from Gastown's formula system

---

## Formula System Overview

Gastown formulas are TOML-based declarative workflow blueprints. They define the shape of work (steps, dependencies, gates, variables) but delegate execution to agents. Four formula types:

| Type | Purpose | Execution Model | Example |
|---|---|---|---|
| **Workflow** | Sequential steps with dependencies | Topological sort, execute ready steps | `shiny`, `beads-release`, `mol-polecat-work` |
| **Convoy** | Parallel specialist legs + synthesis | All legs run simultaneously, synthesis combines | `code-review`, `design`, `mol-prd-review` |
| **Expansion** | Template that generates steps | Placeholder substitution when composed | `rule-of-five` |
| **Aspect** | Cross-cutting concerns | Pre/post advice woven into target steps | `security-audit` |

Formulas compose via `extends` (inherit + modify) and `compose.expand` (wrap a step with an expansion template). Three-tier resolution: project → town → system defaults (embedded in binary).

---

## Complete Formula Inventory

### Core Engineering Pipelines

**`shiny`** — Base engineering workflow
- Steps: design → implement → self-review → test → submit
- Variables: `{{feature}}` (required)
- This is their equivalent of our feature-implementation sequence

**`shiny-enterprise`** — Extended pipeline
- Extends `shiny`
- Composes: `rule-of-five` expansion on `implement` step
- Composes: `security-audit` aspect
- Result: design → [5-pass implement] → self-review → test → [security scan] → submit

**`shiny-secure`** — Maximum quality pipeline
- Extends `shiny-enterprise` with additional security gates

### Release Orchestration

**`beads-release`** — 16-step release process
- Steps: git status check → version bump → changelog → build → test → tag → publish → verify → daemon restart
- Variables: `{{version}}` (required semantic version)
- Acceptance criteria at each step with concrete verification commands

**`gastown-release`** — Similar structure for Gas Town releases

### Convoy Formulas (Parallel Specialist Analysis)

**`code-review`** — 10 parallel review legs
- Legs: correctness, performance, security, elegance, resilience, style, smells, wiring, commit-discipline, test-quality
- Each leg has a `focus` field defining its analysis lens
- Synthesis step combines all findings into consolidated review
- Agent: per-leg assignment possible

**`design`** — 6 parallel design dimensions
- Legs: API & interface, data modeling, UX/interaction, scalability, security, integration
- Synthesis produces design doc with options and trade-offs

**`mol-prd-review`** — 6 parallel PRD review dimensions
- Legs: requirements completeness, gap analysis, ambiguity detection, feasibility, scope validation, stakeholder impact
- Synthesis produces critical questions for human decision

**`mol-plan-review`** — 5 parallel plan review dimensions
- Legs: completeness, sequencing, risk assessment, scope-creep detection, testability
- Synthesis produces must-fix items and go/no-go recommendation

### Expansion Formulas

**`rule-of-five`** — Iterative refinement template (Jeffrey Emanuel pattern)
- Template steps:
  1. `{target}.draft` — Initial attempt, breadth-first, don't aim for perfection
  2. `{target}.refine-1` — Focus: CORRECTNESS (facts, logic, edge cases)
  3. `{target}.refine-2` — Focus: CLARITY (readability, naming, structure)
  4. `{target}.refine-3` — Focus: EDGE CASES (error handling, boundary conditions)
  5. `{target}.refine-4` — Focus: EXCELLENCE (optimization, elegance, production-readiness)
- Each refinement pass needs the previous one
- Applied via `compose.expand` rules in parent formulas

### Aspect Formulas

**`security-audit`** — Cross-cutting security concern
- Pre-scan advice applied before `implement` step
- Post-scan advice applied after `submit` step
- Focus: OWASP Top 10, dependency vulnerabilities, auth/authz

### Mega-Workflow

**`mol-idea-to-plan`** — Full idea-to-beads pipeline (10+ steps)
1. Intake (draft PRD from vague idea)
2. Parallel PRD review convoy (6 legs)
3. Human clarification gate
4. Parallel design convoy (6 dimensions)
5-7. PRD alignment rounds (3 rounds × 2 parallel polecats)
8-10. Plan self-review rounds (3 rounds × 2 parallel polecats)
11. Create beads from plan
12-14. Verify beads (3 passes)
- Variables: `{{problem}}` (required), `{{context}}` (optional)
- Total: 6+ iterative review rounds before implementation begins

### Agent Lifecycle Formulas

**`mol-polecat-work`** — Polecat worker lifecycle
- Steps: load-context → branch → implement → self-review → build-check → pre-verify → submit-mr → await-verdict → handle-verdict
- Acceptance: "MUST have at least 1 commit" (hard gate on commits)
- FIX_NEEDED loop: up to 3 fix attempts before escalating
- Self-cleans on MERGED verdict

**`mol-polecat-review-pr`** — Self-cleaning PR review worker
- Load → review → decide → file findings as beads → `gt done` → disappear

**`mol-polecat-code-review`** — Self-cleaning code review worker
- Similar to PR review but focused on code quality

### Infrastructure Formulas

**`mol-boot-triage`** — Daemon watchdog triage
- Observe → decide → act → cleanup → exit
- Degraded mode: log error but continue

**`mol-gastown-boot`** — Mayor bootstrap
- Verification-gated steps: daemon, deacon, witnesses, refineries (parallel)

**`mol-convoy-feed`** — Dog feeds stranded convoys
- Load → check capacity → dispatch → report

**`mol-dog-doctor`** — Health probe on Dolt server
- Connectivity → resources → orphans → backups

### Durability Proof

**`towers-of-hanoi`** — Pre-computed Hanoi solution as molecule
- Agent mechanically closes steps
- On crash/restart: query `bd mol current`, resume from unfinished steps
- Proves: state is entirely in the molecule graph (bead system)

---

## Execution Model (Go Implementation)

### Core Types

```go
type Formula struct {
    Name        string
    Type        FormulaType  // "workflow", "convoy", "expansion", "aspect"
    Steps       []Step       // Workflow steps with dependencies
    Legs        []Leg        // Convoy parallel units
    Synthesis   *Synthesis   // Post-leg synthesis step
    Vars        map[string]Var     // Template variables
    Extends     []string     // Parent formula names
    Compose     *ComposeRules      // Expansion/aspect composition
    Pour        bool         // If true, steps materialized as sub-wisps
}

type Step struct {
    ID          string
    Title       string
    Description string       // Execution instructions (can be very long)
    Needs       []string     // Prerequisite step IDs
    Parallel    bool         // Can run concurrently with parallel siblings
    Acceptance  string       // Exit criteria
}
```

### Execution Planning

- **TopologicalSort()** — Kahn's algorithm for dependency ordering
- **ReadySteps(completed)** — Returns steps with all dependencies satisfied
- **ParallelReadySteps()** — Groups ready steps into parallel sets + sequential batches
- Cycle detection via DFS (prevents circular dependencies)

### Composition Resolution

```
Formula with extends/compose
  → Load parent(s) recursively (3-tier resolution)
  → Merge vars (child overrides parent)
  → Append steps (parent first, then child)
  → Apply compose.expand rules (substitute target step with expansion template)
  → Validate merged result (unique IDs, valid dependencies, no cycles)
```

### Health & Provisioning

- Formulas embedded in binary, provisioned to `.beads/formulas/`
- Checksums tracked in `.installed.json`
- Safe updates: skip user-modified files, update outdated + untracked
- Health check statuses: ok, outdated, modified, missing, new, untracked, error

---

## Patterns Worth Stealing

### 1. Convoy Model (Parallel Specialist Analysis) — HIGH VALUE

**What it is:** Fan out to N specialists in parallel, each with a narrow focus, then synthesize findings.

**Gastown examples:**
- Code review: 10 parallel legs (correctness, performance, security, elegance, etc.)
- Design: 6 parallel legs (API, data, UX, scale, security, integration)
- PRD review: 6 parallel legs (requirements, gaps, ambiguity, feasibility, scope, stakeholders)
- Plan review: 5 parallel legs (completeness, sequencing, risk, scope-creep, testability)

**How we can use it:**
- **VERIFY phase:** Run spec-check, qa-plan, task-review in parallel instead of sequentially. We already have the infrastructure (Agent tool with `run_in_background`). Just need the synthesis step.
- **Code review:** Our `/code-review` skill uses `code-review-expert` agent type but could fan out to 6 parallel aspects (architecture, quality, security, performance, testing, docs) — same as Gastown's convoy but using our Agent tool.
- **Design exploration:** For complex features, fan out to multiple design agents each exploring a different dimension (data model, API shape, UX flow, performance implications).

**Implementation:** No infrastructure change needed. Convention: when a skill says "convoy review", spawn N parallel subagents with focused prompts, collect results, synthesize. Already doing this in `/pulse-sel`.

### 2. Rule of Five (Iterative Self-Refinement) — MEDIUM VALUE

**What it is:** Before submitting work, the agent self-reviews 4 times with different focus areas:
1. Draft (breadth-first, imperfect)
2. Correctness pass (facts, logic, edge cases)
3. Clarity pass (readability, naming, structure)
4. Edge cases pass (error handling, boundaries)
5. Excellence pass (optimization, elegance, production-ready)

**How we can use it:**
- NOT for all work (too expensive). As an opt-in quality tier.
- Task plan declares quality level: `quality: standard` (single pass) vs `quality: refined` (Rule of Five)
- Good for: complex business logic, security-sensitive code, architectural changes, specs
- Skip for: routine bug fixes, simple feature additions, config changes

**Implementation:** Add as a convention to handoff templates. When quality tier is "refined", the handoff instructions say "implement in 5 passes: draft, correctness, clarity, edge cases, excellence." No infrastructure change.

### 3. Composable Workflow Fragments — MEDIUM VALUE

**What it is:** Small formulas compose into larger ones via extends + expand. Security audit wraps any step. Rule of Five wraps any step. Base pipeline gets enhanced by composition, not modification.

**How we can use it:**
- Our sequences are monolithic (`feature-implementation.md`). Could decompose into composable fragments.
- A task plan could say: "use standard sequence + security review aspect + refined quality tier for implement"
- The orchestrator selects fragments based on task characteristics (security-sensitive? complex? routine?)

**Implementation:** Design-time decision, not infrastructure. The orchestrator's planning phase already selects modes and profiles. Adding "quality tier" and "aspects" (e.g., security review) is a task-plan-level concern.

### 4. Explicit Verification Commands in Step Descriptions — LOW VALUE (we already do this)

**What it is:** Each step includes concrete bash commands for verification:
```bash
git log origin/main..HEAD --oneline     # Verify commits exist
bd blocked --json                       # Verify dependencies correct
```

**Assessment:** We already embed verification methods in completion criteria. Confirmed alignment.

---

## Patterns We Don't Need

| Pattern | Why Not |
|---|---|
| Polecat lifecycle management | Scale-dependent (20+ workers). Our worktree management is simpler. |
| Dog infrastructure tasks | Support roles for large-scale orchestration. Overkill at 1-3 agents. |
| Three-tier formula resolution | We have two-tier skill resolution (project → user). Sufficient. |
| Towers of Hanoi durability proof | Interesting but not applicable to our work. |
| TOML formula syntax | Our markdown-based task plans are more human-readable. The TOML is Gastown's implementation detail, not a transferable pattern. |
| Embedded formula provisioning | We don't compile skills into binaries. Our skill files are direct. |

---

## Comparison: Gastown Formulas vs Our Phase System

| Dimension | Gastown Formulas | Our Task Plans |
|---|---|---|
| **Format** | TOML with structured types | Markdown with conventions |
| **Step tracking** | Beads (JSONL in Git) | Results blocks in task-plan.md |
| **Dependencies** | Explicit `needs` arrays with cycle detection | Implicit phase ordering |
| **Parallelism** | First-class (convoy type, parallel flag) | Ad-hoc (subagents with `run_in_background`) |
| **Composition** | Extends + expand + aspects | Monolithic sequences |
| **Variables** | Explicit `[vars]` with types and defaults | Inline in goal/criteria prose |
| **Verification** | Acceptance field + embedded commands | Completion criteria + confidence plan |
| **Quality tiers** | Composition (shiny → shiny-enterprise → shiny-secure) | Single tier (could add) |
| **Review model** | Convoy (N parallel specialists → synthesis) | Sequential (spec-check → qa-plan → task-review) |
| **Self-refinement** | Rule of Five (4 editorial passes) | None (single-pass implement) |

### Key Gaps Our System Has:
1. No parallel verification convoy (run verification aspects in parallel)
2. No self-refinement passes (Rule of Five)
3. No composable quality tiers (standard vs refined vs secure)
4. No explicit dependency DAGs (phases are implicitly sequential)
5. No aspect composition (security review as cross-cutting concern)

### Key Advantages Our System Has:
1. Human-readable markdown (vs TOML that needs tooling to read)
2. Simpler mental model (5 named phases vs arbitrary DAGs)
3. Flexible prose in completion criteria (vs structured acceptance fields)
4. Operator can edit plans in any text editor
5. No compilation/provisioning step required
