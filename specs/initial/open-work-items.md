# Open Work Items: Gastown Architecture Session

**Date:** 2026-03-14 — 2026-03-15
**Source:** Architecture discussion sessions analyzing Gastown patterns

---

## CRITICAL BOTTLENECK

**`/orchestrate` model centering on TeamCreate is obsolete.** Long-running Claude orchestrator sessions are our biggest bottleneck right now. The TeamCreate-based model requires one long-lived session to coordinate all phases — this session fills context, drifts, and becomes unreliable.

**Need:** Separate orchestration into an external process, not a long-running Claude session. Gastown solves this with session chaining (each molecule step is a discrete session, state persists in Beads/Git, `gt prime` reinjects context at each session start). Our plan docs already provide the persistence layer — what's missing is the automation that chains sessions together.

**See:** `gastown-tmux-protocols.md` § "Session Boundaries" and § "The GUPP Gap" for the detailed comparison.

---

## High Priority

### #20: Context Budget Guard — detect and gate context exhaustion
- **Status:** Spec drafted
- **Spec:** `context-budget-guard/spec.md`
- **Priority:** IMMEDIATE VALUE — shell script + hook, shippable in one session
- **Source:** Gastown `scripts/guards/context-budget-guard.sh` — proven, battle-tested pattern
- Adapts gastown's three-tier threshold system (warn 75% → soft gate 85% → hard gate 92%)
- Role-based gating: orchestrator/autonomous sessions hard-gate, interactive/implementer warn-only
- Fail-open on errors (missing jq, bad transcript → allow the turn)
- Directly addresses the CRITICAL BOTTLENECK (long-running sessions degrading at context limits)
- Zero new infrastructure — uses existing UserPromptSubmit hooks + JSONL transcripts

### #18: Integrate /code-review into TDD refactor step
- **Status:** In progress (team implementer-2 researching + implementing)
- Red/green/REFACTOR loop where refactor step calls /code-review subagent
- Implementation-phase legs: correctness, security, wiring, smells, performance, resilience, test-quality
- Fix P0/P1 → retest → loop until clean
- Evidence: First /code-review run on already-merged PR found actionable issues that would have been caught during implementation

### #15: Build send-claude-message IPC tool (first Go CLI project)
- **Status:** Open
- Abstract from tmux. Interface: `send-claude-message <session> "message"` with `--queue` flag
- Implements Gastown's 8-step protocol (copy mode exit, sanitize, send-keys -l, 500ms delay, Escape, 600ms gap, 3-retry Enter, SIGWINCH wake)
- Also supports file-based queue mode (write to inbox, agent drains via hook)
- First Go CLI tool — proves out Go workflow
- **NOTE:** We already have a Telegram → Claude Code session delivery system (via Telegram MCP + egregore). When building the IPC tool, incorporate this as a delivery channel — not just tmux send-keys but also Telegram-originated messages routed to sessions. The IPC tool should be the unified entry point for all message delivery to Claude sessions regardless of source (tmux, Telegram, file queue, future WebSocket, etc.)

### #1: Agent roles / classes — per-session hooks + role-based injection
- **Status:** Open
- Foundation for everything: roles compose with permission system, gating hooks, supervisor behavior, injections
- Role definition JSON selects: hooks, injections, gating rules, sandbox config, autonomy profile
- Proposed roles: implementer, patrol, fixer, browser-qa, supervisor, interactive, designer
- See `autonomy-architecture.md` for full design

---

## Medium Priority

### #17: Per-session inbox with drain mechanism
- **Status:** Open
- File-based per-session message queue (JSON to `/tmp/claude-sessions/{session_id}/inbox/`)
- Drained by UserPromptSubmit hook (interactive) + explicit drain calls at phase transitions (autonomous)
- Atomic rename-based claiming from Gastown pattern
- Key insight: UserPromptSubmit doesn't fire during autonomous work. Gastown solves with session chaining; we solve with explicit drain calls at phase boundaries.

### #16: Create Go CLI development skill
- **Status:** Open
- No Go CLI skill exists anywhere (confirmed by audit)
- Built-in cli-expert is npm/Node-focused
- Need: project structure, build/test workflow, module management, single-binary distribution
- Prerequisite for building Go tools efficiently

### #4: Supervisor/watchdog agents — expand beyond /permission-manager
- **Status:** Open
- MVP: tmux health check + timeout enforcement + Telegram escalation bus
- Detects: stuck agents, error loops, off-course behavior, permission blocks, context filling
- Gastown equivalent: Witness/Deacon patrol pattern
- NOT overkill — solves specific painful operator problem

### #10: AppSignal autonomous triage → ticket creation
- **Status:** Open
- Starting point: triage that creates tickets, which can be actioned autonomously
- Architecture: patrol loop → classify alerts → create Linear tickets → existing /orchestrate handles execution
- Uses /fire-watch as foundation
- Autonomy ladder: Level 1 (triage + tickets) → Level 2 (auto-fix easy ones) → Level 3 (full patrol)

### #14: Cross-reference domain skills against code-review anti-pattern vocabulary
- **Status:** Parked (blocked by code-review skill stabilization)
- After code-review skill legs and anti-pattern categories stabilize, sweep domain skills (ruby-on-rails, vue, etc.) to cross-reference
- Bidirectional reinforcement: code review finds pattern, domain skill knows the project-specific fix

---

## Low Priority / Parked

### #7: Code coverage gating via hooks
- **Status:** Open
- PreToolUse hook that blocks PR creation if coverage drops
- Concrete example of executable gating pattern
- Alternative: gate at PR creation time in CI (less novel but more practical)

---

## Decisions Made (Resolved Threads)

| # | Decision | Resolution |
|---|---|---|
| A1 | GUPP "physics not politeness" framing | Adopt in autonomous mode injections |
| A2 | Session-surviving identity | Plan doc is sufficient, refine reactively |
| A3 | Beads assessment | WATCH — md files sufficient at 1-3 agents. Revisit at 5+ |
| A4 | Machine-readable workflow state | Incremental checklist convention (`[x]`/`[ ]`), no schema change |
| A5 | Orchestrate vs NDI | ~80% equivalent. Convention fixes close gap |
| D1 | Permission system | Tiered by role. Sandbox + bypassPerms for coding roles, builtin+hooks for network-dependent |
| A6 | Convoy model | Already intrinsic in /orchestrate modes |
| A7 | Rule of Five | Adopt as opt-in quality tier in modes/sequences |
| A8 | Composable quality tiers | Skip until need identified |
| A9 | TOML vs markdown | Markdown — our consumer is an LLM, not a compiler |
| A10 | Go for CLI tooling | Yes — single binary, no venv issues, Claude is excellent at Go |

---

## Spec Files Produced

All in `specs/orchestrate-next/`:

| File | Contents |
|---|---|
| `gastown-analysis.md` | Overall Gastown comparison, what to steal/skip/watch |
| `autonomy-architecture.md` | Agent roles, tiered permissions, supervisor, AppSignal use case |
| `decisions-2026-03-14.md` | Formal agreements + pending decisions |
| `gastown-formulas-audit.md` | All 43 formulas analyzed, 3 patterns worth stealing |
| `review-models-upgrade.md` | Review skills gap analysis + upgrade plan (spec-review, plan-review, code-review upgrades) |
| `code-review-skill-design.md` | Code review skill design doc |
| `gastown-tmux-protocols.md` | Tmux protocols, gap analysis vs tmux-claude-send, session inbox, session boundaries |
| `open-work-items.md` | This file — all open threads and work items |

---

## Skills Built

| Skill | Location | Status |
|---|---|---|
| `/code-review` (new) | `~/.claude/skills/code-review/SKILL.md` + `resources/` | Built, tested against real PR |
| `/code-review` (old command) | `~/.claude/commands/code-review.md` | Deleted |
