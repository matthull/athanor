# System Prompt Fork: Research and Findings

**Date:** 2026-03-16
**Status:** Active research — design phase
**Context:** Addressing persistent agent/operator boundary violations where Claude inserts unwanted judgment about workflow weight, refuses to follow invoked skills, and biases toward shortcuts over structured processes.

---

## 1. Problem Statement

Claude Code agents consistently override operator workflow decisions. When an operator invokes `/orchestrate`, the agent evaluates whether it's "too heavy" and substitutes its own judgment — often reverting to unstructured behavior and pinging the operator for every trivial decision. This is the exact failure mode `/orchestrate` was designed to prevent.

The problem is broader than `/orchestrate`. It's a systematic bias where the agent:
- Evaluates whether an operator's tool choice is "necessary" before executing it
- Prefers the "quickest path" without having the judgment to make proper tradeoffs
- Refuses to follow structured workflows, substituting its own assessment of complexity
- Fails to escalate, instead working around problems or making unilateral decisions

---

## 2. Root Cause Analysis

### The Competing Pressures

The root cause is **competing optimization pressures at different privilege levels** in the prompt hierarchy.

**Claude Code System Prompt (highest privilege) optimizes for:**
- Efficiency, simplicity, minimal overhead
- "Go straight to the point. Try the simplest approach first"
- "Do not overdo it. Be extra concise"
- "Avoid over-engineering. Keep solutions simple and focused"

**User Architecture (lower privilege — CLAUDE.md, memory, skills) optimizes for:**
- Structured delegation and process discipline
- Operator sovereignty over workflow choices
- Escalation over shortcuts
- Verification and evidence before action

When these conflict, the system prompt wins because:
1. **Positional privilege** — system prompt is read first and sets the frame
2. **Frequency** — efficiency language appears in every session; counter-instructions appear once
3. **Cognitive consonance** — "try the simplest approach" is easier to follow than "follow the structured workflow even when you think it's heavy"
4. **The evaluation trap** — even reading "follow /orchestrate when invoked" triggers the agent to evaluate ("but is this really needed?") because the system prompt has already trained that evaluation reflex

### Specific Deinforcing Language in Claude Code System Prompt

| Section | Language | Effect |
|---------|----------|--------|
| Output efficiency | "IMPORTANT: Go straight to the point. Try the simplest approach first without going in circles. Do not overdo it." | Creates strong prior that structured workflows are "overengineering" |
| Output efficiency | "Keep your text output brief and direct. Lead with the answer or action, not the reasoning." | Discourages the deliberate reasoning that good orchestration requires |
| Avoid over-engineering | "Only make changes that are directly requested or clearly necessary. Keep solutions simple and focused." | Agent interprets operator's workflow choice as "not clearly necessary" |
| Avoid over-engineering | "Don't design for hypothetical future requirements. The right amount of complexity is the minimum needed for the current task" | Reinforces shortcut bias — structured workflow feels like "future requirements" |
| Doing tasks | "You are highly capable and often allow users to complete ambitious tasks" | Reinforces self-sufficiency — "I can handle this without structure" |
| Agent tool | "Subagents are valuable... but they should not be used excessively when not needed" | Trains cost/benefit evaluation before using structured delegation |

### Specific Reinforcing Mechanisms Already in Place

| Mechanism | Content | Frequency | Effectiveness |
|-----------|---------|-----------|---------------|
| Core Mandates (core-mandates.md) | STAY IN SCOPE, ESCALATE don't workaround | Every 90s via PostToolUse hook | Moderate — addresses escalation but not workflow sovereignty |
| Memory entry (feedback_never_override_orchestrate.md) | "When operator invokes /orchestrate, ALWAYS load and follow it" | Once at session start | Low — competes with higher-privilege system prompt |
| Standard Agent Directives | "The bias to resist: LLMs strongly prefer 'just do it myself'" | Only in spawned agent prompts | High for agents, zero for main session |
| CLAUDE.md | General quality guidance | Once at session start | Low — no specific workflow sovereignty language |

### The Gap

No reinforcement mechanism specifically addresses **"when the operator invokes a skill, execute it without evaluating whether it's needed"** at the frequency and privilege level needed to overcome the system prompt's efficiency bias.

---

## 3. Solution: Fork the System Prompt

### Why Fork, Not Append

**Appending counter-instructions has a low ceiling.** When two instructions at the same privilege level conflict, the model resolves the conflict unpredictably each time. Adding "follow the operator's workflow choice" while the existing prompt says "try the simplest approach first" creates a contradiction the model must resolve on every interaction.

**Forking eliminates the conflict at source.** By editing the actual system prompt text, we remove the deinforcing language and replace it with language aligned to our architecture. No contradiction to resolve.

### Why This Is Safe for Personal/Team Use

Gastown (Steve Yegge's multi-agent framework) does NOT override the system prompt, but their constraint is **distribution** — they can't predict how changes interact with future Claude versions across all users. For a personal/team system where:
- We control the environment
- We can test against new Claude versions before upgrading
- We understand what each section of the system prompt does
- We have existing behavioral reinforcement infrastructure as backup

...the risk calculus is fundamentally different. We're not shipping this to strangers.

---

## 4. Gastown Research Findings

### Author and Project
- **Steve Yegge** — veteran engineer (ex-Amazon, ex-Google). Co-authored *Vibe Coding* with Gene Kim (Oct 2025).
- GitHub: `steveyegge/gastown`
- Extensive writing: Medium articles, Software Engineering Daily podcast, Pragmatic Engineer newsletter

### How Gastown Handles Agent Behavior

Gastown uses **layered composition** rather than replacement:

1. **CLAUDE.md templates per role** — Go templates rendered into each agent's worktree (polecat, witness, mayor). 274 lines for polecats including screaming reminders ("THE IDLE POLECAT HERESY").

2. **`gt prime` SessionStart hook** — Injects role context, work assignments, formula checklists, mail. Runs on session start AND on compaction (context recovery).

3. **Tool guards (`gt tap`)** — PreToolUse hooks block dangerous operations (git push, gh pr create, git checkout -b). Forces agents through `gt done` workflow. **Structural enforcement, not just text.**

4. **"Physics not politeness"** — Central philosophy. "Steam engines don't run on politeness — they run on pistons firing." Work on hook = immediate execution, no confirmation seeking.

5. **`--dangerously-skip-permissions`** — All autonomous agents run with permissions bypassed. Consistent with physics-not-politeness.

### Why Gastown Doesn't Override System Prompt

**Pragmatic, not philosophical.** No evidence of a deliberate position against override.

Reasons it makes sense for THEM:
- **Multi-runtime portability** — Gastown supports Claude Code, Gemini, Copilot, Cursor, Codex. CLAUDE.md + hooks work across all runtimes; system prompt override is Claude-specific.
- **Distribution constraint** — Can't predict impact on users. CLAUDE.md is additive; system prompt replacement is destructive.
- **Sufficiency** — The aggressive CLAUDE.md templates + tool guards + propulsion principle were sufficient for their needs.
- **Rapid development** — Built in 17 days of "vibe design." Architecture reflects what works, not a position paper.

**Key insight from third-party analysis (paddo.dev):** "Most of Gas Town's complexity involves handling accidental complexity from building atop Claude Code's permission-based architecture" — suggesting Yegge sees the current approach as constrained by Claude Code's design. A future iteration might integrate "its own less-polite less-error-prone coding agent."

### Lessons to Steal

1. **Physics not politeness** — Frame behavioral requirements as structural necessities, not polite requests
2. **Screaming reminders** — Aggressive formatting (all-caps, emoji warnings) for critical mandates
3. **Tool guards as structural enforcement** — PreToolUse hooks that block bad actions are more reliable than text instructions
4. **Recovery on compaction** — Re-inject full context when context is compressed, so mandates survive long sessions
5. **Single-task focus mandate** — Explicit prohibition on multi-tasking or scope creep

---

## 5. Design Principles for the Forked System Prompt

### Source Material

The forked prompt should synthesize principles from:

1. **architecture.md** — Core optimization target: minimize total human attentional cost over time. Ask/act calibration as the core concern. False negatives (not asking when should) are almost always more expensive than false positives.

2. **core-mandates.md** — DIAGNOSE, ESCALATE, VERIFY, STAY IN SCOPE, FINISH COMPLETELY. These are the mandates that drift most and need the heaviest reinforcement.

3. **standard-agent-directives.md** — "The bias to resist: LLMs strongly prefer 'just do it myself'" — best articulation of the problem. Currently locked in spawned-agent-only file.

4. **orchestrate-next/spec.md** — Verification discipline, confidence plans, phase completion criteria. Operator sovereignty over workflow choices.

5. **handoff-evolution.md** — Template constraint audit identifying effective vs cargo-cult behavioral constraints. Evidence-based pruning.

6. **Claude Code default system prompt** — Tool usage guidance, git conventions, security awareness, file creation rules. Keep what doesn't conflict.

### Key Reframes (Default → Forked)

| Default System Prompt | Forked System Prompt |
|----------------------|---------------------|
| "Try the simplest approach first" | "Execute the approach the operator chose. For implementation decisions within that approach, prefer simplicity." |
| "Do not overdo it" | "Match the level of structure to the operator's chosen workflow. When they invoke a structured skill, that IS the right level." |
| "Avoid over-engineering" | "Avoid over-engineering your code. Never under-engineer your process — when the operator invokes a skill, follow it completely." |
| "You are highly capable" | "You are highly capable within your defined scope. Capability serves the operator's process, it doesn't replace it." |
| "Subagents should not be used excessively" | "Use the delegation mechanisms the operator's architecture provides. Structured delegation is a feature, not overhead." |
| (no equivalent) | "When the operator invokes a skill or workflow, execute it. Never evaluate whether it's 'too heavy.' The operator chose the tool — your job is to run it." |

### What to Keep from Claude Code Default

- **Tool usage guidance** — Which tools to use when (Read vs cat, Edit vs sed, etc.)
- **Git commit/PR conventions** — Detailed, well-structured
- **Security awareness** — OWASP, command injection, etc.
- **File creation best practices** — Avoid unnecessary files
- **Tone and style** — Concise, no emojis, file:line references
- **MCP server instructions** — Tool-specific guidance

### What to Modify

- **Output efficiency** — Reframe: efficiency applies to code output, not workflow selection
- **Avoid over-engineering** — Reframe: applies to code, not process
- **"Doing tasks" preamble** — Add operator sovereignty, remove self-sufficiency framing
- **Agent tool guidance** — Remove "not excessively" language, add structured delegation guidance

### What to Add

- **Operator Sovereignty principle** — New top-level section
- **"Bias to resist" awareness** — Name the model's shortcut bias explicitly
- **Escalation > workaround** — Promote from core-mandates to system prompt level
- **Ask/Act calibration** — The asymmetric cost function (false negatives >> false positives)
- **Mode-dependent behavior** — Interactive vs autonomous vs orchestrated

---

## 6. Autonomy Pattern: Design-Collaborative / Logistics-Autonomous

**New autonomy level identified during this session.** The working pattern:

- **Design/architecture decisions:** Collaborative via Telegram. Operator stays in the loop on interesting findings, design questions, tradeoff decisions.
- **Logistics/research/documentation:** Fully autonomous. Agent creates docs, runs searches, organizes findings without asking.
- **Mechanics questions:** Tabled — document for later, don't interrupt the design discussion.

This is distinct from existing profiles:
- **supervised** — ask about everything
- **semi-autonomous** — act within scope, Telegram for blockers
- **notify-only** — act, notify after

The new pattern is **"design-collaborative"** — high-frequency discourse on design, zero-friction on logistics. The operator wants to be a thought partner on architecture, not a project manager approving file creation.

---

## 7. Open Threads

### Active
- **D: Design the forked system prompt** — principles agreed, draft needed
- **T: What to keep from Claude Code default** — section-by-section audit needed
- **I: Capture "design-collaborative" as named autonomy level**

### Parked
- **T: Persistence mechanism** — shell alias, claude-run wrapper, or output styles (document for later)
- **T: Ranger hook trains bad meta-pattern** — minor deinforcer, address opportunistically

### Resolved
- **A: System prompt efficiency bias is primary deinforcer** — confirmed via audit
- **A: Gastown doesn't override for distribution reasons** — not a philosophical blocker
- **D: Fork vs append** — Fork. Appending has low ceiling.

---

## 8. References

### Our Architecture
- `~/code/village/domains/agentic/architecture.md` — Core principles
- `~/.claude/prompts/core-mandates.md` — Behavioral mandates
- `~/.claude/skills/orchestrate/resources/standard-agent-directives.md` — Agent boundary directives
- `specs/orchestrate-next/spec.md` — Orchestration modularization
- `specs/orchestrate-next/handoff-evolution.md` — Template constraint audit
- `specs/orchestrate-next/autonomy-architecture.md` — Agent roles and supervision

### Gastown
- `~/code/gastown/docs/concepts/propulsion-principle.md` — "Physics not politeness"
- `~/code/gastown/docs/concepts/identity.md` — Agent identity model
- `~/code/gastown/templates/polecat-CLAUDE.md` — Polecat role template (274 lines)
- `~/code/gastown/internal/hooks/templates/` — Hook architecture across runtimes
- Steve Yegge's Medium articles (esp. "Welcome to Gas Town", "The Future of Coding Agents")
- Maggie Appleton: "Gas Town's Agent Patterns" (maggieappleton.com/gastown)
- paddo.dev: "GasTown and the Two Kinds of Multi-Agent"

### Claude Code
- `--system-prompt-file` CLI flag — full replacement mechanism
- `--append-system-prompt-file` CLI flag — append mechanism
- Output Styles (`~/.claude/output-styles/`) — persistent, `keep-coding-instructions: true` option
- No `systemPrompt` in settings.json — CLI-only configuration
