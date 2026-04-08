# The Artisan's Workbench

## A System Design for Eventual Comprehension Consistency in Agentic Software Development

**Created:** 2026-04-05
**Status:** Synthesis — consolidating three research passes into actionable design
**Context:** Terminal-first (tmux + neovim), Claude Code, swarm dispatch, ADHD-informed work sessions

---

## The Core Insight

Two pipelines. Fully decoupled. Both non-negotiable.

### Pipeline 1: Autonomous Shipping

The agent swarm produces, verifies, and ships code without human intervention. It is converging toward full autonomy on the question "is this code correct and not awful?" The harness — CLAUDE.md, TDD, spec-driven development, automated AI review, CI gates — is designed to be **sufficient on its own**. Human PR review is vestigial and actively being phased out in favor of automated merge rights with quality gates.

**Design target:** This pipeline should work as if no human will ever look at the code. Every quality assurance mechanism must be automated and self-correcting. If the harness can't catch a class of error without human eyes, the harness is broken and must be improved until it can.

### Pipeline 2: Eventual Comprehension

The developer periodically, opportunistically, deeply explores chunks of the codebase to build and maintain a mental model of what the system has become. This is not review — it's exploration. It is decoupled from shipping. Code ships before the developer understands it. The developer *will* understand it, eventually, on their own schedule.

**Design target:** When the developer sits down for a focused session, they should be able to achieve deep understanding of any area of *any* codebase — their own, a teammate's, a new project — in the minimum possible time, aided by AI that curates what they see, provides context, generates diagrams, and answers questions — while the developer reads actual code with their own eyes. This is a pure exploration tool designed for maximum insight during Adderall-fueled hyperfocus sessions.

### The Relationship Between Pipelines

Pipeline 2 calibrates Pipeline 1. When the developer explores and finds the code is good, that's a strong signal the harness is working. When they find something wrong — not a bug (tests would catch that) but architectural drift, unnecessary complexity, inelegant patterns, subtle wrongness — that's a harness bug. The finding feeds back into Pipeline 1 as a harness improvement: a new guideline, a new automated check, a new test pattern.

Over time, Pipeline 1 gets better because Pipeline 2 keeps finding its blind spots. Pipeline 2 gets less frequent because Pipeline 1 keeps getting better. But Pipeline 2 never goes to zero, because:

- The codebase evolves and new areas need comprehension
- Models change and may introduce new failure modes
- Architecture is an aesthetic judgment that resists full automation
- "Simple, explicit, well-organized code" is a quality that requires human eyes to verify
- Comprehension debt is real and compounds silently

---

## Pipeline 1: Autonomous Shipping — Detailed Design

### Principle: Design as if no human will review

The harness must handle every quality dimension that can be automated:

**Correctness** — TDD with spec-driven test generation. Property-based testing (Hypothesis/fast-check) for invariant verification. Mutation testing (Stryker/mutmut) to validate test adequacy. The agent must not be able to delete or weaken tests.

**Conventions** — Linters, formatters, CLAUDE.md coding guidelines enforced by pre-commit hooks and CI. Cross-model review (generate with one model, review with another) for pattern compliance.

**Architecture** — Architecture fitness functions as automated tests. Dependency rules enforced by dependency-cruiser or similar. Module boundary violations caught automatically. CodeScene Code Health as a merge gate — block merges that degrade health below threshold.

**Security** — Static analysis (Semgrep, CodeQL) in CI. AI security review as a dedicated review agent. Explicit rules in CLAUDE.md about auth patterns, input validation, secrets handling.

**Self-correction** — The Ralph Wiggum Loop principle: agents must verify their own work against the spec before declaring completion. Multiple review agents with confidence thresholds. Non-zero exit codes from quality hooks prevent task completion.

**Drift detection** — Automated tracking of codebase health metrics over time. Complexity trends, coupling metrics, duplication rates. Alerts when metrics trend in wrong direction — these are harness alarms, not code alarms.

### Merge Automation Progression

Current state: Human reviews each PR (lightweight, checking process compliance more than code).

Near-term: Auto-merge when all gates pass. Human reviews only flagged PRs (failed gates, anomalous changes, high-risk areas).

Target state: Full auto-merge with quality gates. Human never reviews for shipping purposes. All human code interaction happens in Pipeline 2.

### The Harness Improvement Loop

Every automated check that exists today was born from a human observation during Pipeline 2. The loop:

1. Developer explores code in Pipeline 2
2. Developer notices a quality issue the harness didn't catch
3. Developer determines *why* the harness missed it
4. Developer adds/modifies a harness component (new test pattern, new CLAUDE.md guideline, new linter rule, new fitness function)
5. Harness now catches that class of issue automatically
6. Pipeline 1 improves; Pipeline 2 finds fewer issues of that type

This is the "code gardener" role made precise: the developer's output is harness improvements, not code edits.

---

## Pipeline 2: Eventual Comprehension — Detailed Design

### Principle: Exploration, not review

Pipeline 2 is not triggered by PRs, not coupled to shipping, and not systematic. It is a general-purpose tool for deeply understanding code you didn't write, fast. Use cases include:

**Reviewing your own swarm's output:**
- "Walk me through what the agents built for the payments feature this sprint"
- "I haven't looked at the auth module in three weeks — what changed?"
- "The complexity metrics are trending up in the API layer — show me why"

**Reviewing teammates' work:**
- "Sarah's been building the reporting module — walk me through her architecture"
- "I need to review Jake's PR but it touches 40 files across 3 modules — orient me"

**Getting up to speed on something new:**
- "I just joined this project — show me how the data flows from API to database"
- "We're integrating a new payment processor — walk me through the vendor SDK and our adapter layer"
- "Show me the test architecture — how are things organized, what patterns are used?"

**Harness calibration:**
- "Show me the last 20 files the swarm produced — I want to spot-check quality patterns"
- "Walk me through the most complex thing the agents built this month"

### Future Iteration: Comprehension Delta Map

> **Deferred.** An interesting concept worth revisiting once the core exploration tool is built and in daily use: a visualization of "what I've explored vs. what's changed" — tracking comprehension gaps across the codebase over time. Involves state management and tracking infrastructure that would be premature to build now. Only implement if actual exploration sessions reveal a real problem with knowing where to focus. The simpler approach — just ask Claude "what's changed the most that I haven't looked at?" using git history — may be sufficient for a long time.

### The Exploration Session

A focused, high-engagement, time-boxed session. The developer sits down (or ideally, works from a standing desk or comfortable chair — minimal screen time means making screen time count).

**Session initiation:**
```
/explore "walk me through what happened in the API layer this sprint"
```
or
```
/explore "I'm new to this repo — show me the architecture starting from the main entry point"
```
or
```
/explore "deep dive into Sarah's new caching implementation"
```
or
```
/explore "spot-check the last 10 files the swarm produced — show me patterns"
```

**What Claude does before the developer sees anything:**

1. **Pre-analysis** (subagents, runs in seconds):
   - Maps the relevant area: files changed, dependency graph, architecture diagram
   - Identifies what's new, what's modified, what's deleted
   - Runs the comprehension delta check: "You last explored this area 18 days ago. 14 files changed, 3 new files, 1 deleted. Key changes: new Redis caching layer, refactored user service interface."
   - Pulls any relevant findings from automated Pipeline 1 review: "AI review flagged one concern about error handling in cache_service.ts — confidence 72%, below threshold for blocking"
   - Generates ASCII architecture diagram of the area via mermaid-ascii

2. **Session setup** (neovim via RPC):
   - Opens the architecture diagram in a scratch buffer
   - Opens the 3-5 most important entry point files in splits
   - Positions cursor at the key function/class
   - Has the full context loaded and ready for questions

3. **Guided exploration** (interactive):
   - Developer reads code, jumps around, asks questions
   - Claude navigates to related files on request
   - Claude generates diagrams on the fly ("show me the data flow for this request")
   - Claude explains design decisions inferred from commit messages and code comments
   - Claude highlights areas where the code diverges from established patterns
   - Developer's cursor position is tracked — Claude stays context-aware

4. **Session wrap-up**:
   - Any harness improvement ideas are captured as todos
   - If developer found a quality issue: immediately opens a harness improvement task
   - Brief summary of what was explored and key takeaways

### ADHD-Friendly Design Principles

- **High engagement**: Claude actively manages the view, poses questions, generates diagrams. This is not passive reading — it's interactive exploration with a knowledgeable guide.
- **Clear entry point**: "Here's what's most interesting to explore right now" eliminates the paralysis of "where do I start?"
- **Time-boxed**: 30-60 minute sessions. Clock visible. When time's up, session wraps with a summary.
- **Variable depth**: Developer can go deep on one function or skim across a module. Claude adjusts depth to follow the developer's attention.
- **Tangible output**: Every session produces either "harness is calibrated, comprehension updated" or "found issue, harness improvement filed." Never "I just read code for an hour and don't know if it mattered."
- **No guilt about gaps**: You'll get to everything eventually. Gaps are normal and expected. The tool works just as well on code from yesterday as code from last month.

---

## The Technical Stack

### Pipeline 1 (Autonomous Shipping)

**Already in place or immediately implementable:**
- Claude Code with robust CLAUDE.md (coding guidelines, TDD enforcement, review loops)
- Swarm dispatch via /batch or worktrees + tmux
- Automated test suite (unit, integration, property-based)
- CI/CD with linters, formatters, type checking
- Claude Code /review plugin (4 parallel review agents)
- Git hooks for pre-commit quality checks

**To add:**
- CodeScene MCP server for Code Health scoring as merge gate
- Mutation testing (Stryker/mutmut) as CI step to validate test quality
- Automated complexity/coupling trend tracking (time series, alerting on drift)
- Cross-model review for high-risk changes
- Architecture fitness function tests (dependency rules, module boundaries)
- Automated merge with quality gates (replacing human PR review)

### Pipeline 2 (Eventual Comprehension)

**The Exploration Environment:**
- tmux session with 3 panes: neovim (main), Claude Code (conversation), scratch (diagrams/output)
- neovim with `--listen` socket for RPC control
- mcp-neovim-server connected to Claude Code for editor control
- diffview.nvim for multi-file navigation
- sidekick.nvim or claudecode.nvim for bidirectional awareness
- mermaid-ascii for terminal-native architecture diagrams
- difftastic for structural diffs (cutting through AI reformatting noise)
- lazygit for git history exploration

**The `/explore` Skill** (to build):
- Custom Claude Code slash command (~200 lines of skill markdown)
- Pre-analysis phase: spawns subagents for area mapping, dependency graphing, change analysis
- Session setup phase: opens neovim with relevant files via RPC
- Interactive phase: responds to developer questions with full codebase context
- Wrap-up phase: captures harness improvement ideas as todos
- Works on any codebase — own projects, teammates' projects, new repos

### Shared Infrastructure

- **CodeScene** (or similar): Serves both pipelines. Pipeline 1 uses it as a merge gate. Pipeline 2 uses it to identify areas worth exploring (low health scores, high churn).
- **Git history**: Pipeline 1 reads it for change context in review. Pipeline 2 reads it to orient the developer ("what changed in this area recently?").
- **CLAUDE.md / harness configuration**: Pipeline 1 consumes it. Pipeline 2 produces improvements to it.

---

## What This Is Not

**Not vibe coding.** Every line of code will be understood by a human. Just not necessarily before it ships, and not in order.

**Not "trust the plan."** The plan matters, but implementation quality is independently verified by automated systems AND by eventual human exploration.

**Not "trust the harness."** The harness is designed to be trustworthy, but it is continuously verified by human exploration sessions that test its effectiveness.

**Not traditional code review.** No human reviews code for the purpose of approving it to ship. Humans explore code for the purpose of understanding it and calibrating the system that produces it.

**Not a dark software factory.** The developer maintains deep understanding of the entire codebase. They just do it on their own schedule rather than at the moment of shipping.

---

## The Manufacturing Parallel, Made Precise

A modern factory floor:

- **The production line** (Pipeline 1) is designed to produce zero-defect output. Statistical Process Control monitors the process continuously. Automated inspection catches defects in real time. The line is *designed to work without human inspectors*.

- **Quality auditors** (Pipeline 2) periodically walk the floor, pull samples, examine output, and verify process integrity. Their job is not catching defects — the automated systems do that. Their job is **verifying that the automated systems are still working correctly** and **understanding the process deeply enough to improve it**.

- When an auditor finds a defect, it's a **process alarm**, not just a bad unit. The response is to fix the process, not just the unit.

- Auditors don't inspect every unit. They don't inspect on a fixed schedule. They inspect based on risk, on process changes, on their own judgment about where attention is most needed. But they never stop inspecting entirely, because processes drift.

The developer is the quality auditor. The harness is the production line. Code exploration sessions are process audits. Harness improvements are process fixes.

---

## Implementation Roadmap

> **NOTE:** These phases are directional only — a rough sketch of the order things might happen. Actual scope, timing, and priorities will be defined separately based on what we learn as we go. Don't treat this as a plan.

### Phase 1: Strengthen Pipeline 1 (weeks 1-2)

- Add CodeScene MCP server to Claude Code config
- Add mutation testing to CI for critical modules
- Set up complexity/coupling metric tracking with trend alerting
- Write architecture fitness function tests for key module boundaries
- Begin tracking which PRs get human-modified vs. merged as-is (data for confidence in automation)

### Phase 2: Build the Exploration Environment (weeks 2-4)

- Install and configure mcp-neovim-server
- Set up the 3-pane tmux layout with neovim socket
- Install mermaid-ascii, difftastic, lazygit
- Write a basic `/explore` skill that does pre-analysis and opens relevant files
- Test with a few manual exploration sessions, iterate on the skill

### Phase 3: Decouple Shipping from Comprehension (weeks 4-6)

- Implement auto-merge with quality gates for low-risk changes
- Gradually expand auto-merge scope based on confidence data from Phase 1
- Fully decouple: shipping happens without human involvement for passing PRs
- Human interaction moves entirely to Pipeline 2

### Phase 4: Refine and Compound (ongoing)

- Every exploration session improves the harness
- Every harness improvement makes the next exploration session find fewer issues
- Exploration sessions become more efficient as the `/explore` skill learns from patterns
- If a real problem emerges with knowing where to focus, revisit the Comprehension Delta Map concept
- The factory gets better. The auditor gets more strategic.

---

## Open Questions

1. **What's the right level of pre-analysis vs. interactive exploration?** Too much pre-analysis and the session feels like reading a report. Too little and the developer wastes hyperfocus time orienting. The balance probably varies by session type (new codebase vs. reviewing recent changes vs. spot-checking quality).

2. **How does the `/explore` skill handle different codebase scales?** Exploring a 500-file mobile app is different from exploring a 5,000-file Rails monolith. The pre-analysis phase needs to be smart about scope and not overwhelm with information.

3. **What triggers a mandatory exploration vs. opportunistic?** Some changes might be important enough to explore before the next natural session — new architectural patterns, security-critical changes, new external integrations. How are these flagged without recreating the PR review bottleneck?

4. **How does this work in a team context?** Comprehension is distributed across a team. Does each developer maintain their own exploration practice? Is there value in shared exploration sessions?

5. **What's the right interaction model during exploration?** Should Claude lead ("let me show you the next important thing") or follow ("I'll wait for your questions")? Probably depends on the session type and developer's energy. The skill should support both modes.

6. **How do exploration findings get prioritized as harness improvements?** Not every finding justifies a harness change. Some are one-off oddities, some are critical patterns. Need a lightweight triage process that doesn't become its own overhead.

7. **Can exploration sessions be partially asynchronous?** Claude does pre-analysis and generates a "briefing" that the developer reads on their phone before sitting down at the terminal. The actual code exploration is synchronous but the context-gathering happens earlier. This would make the focused session even more efficient.
