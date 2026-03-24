# Autonomy Architecture: Agent Roles, Gating, and Supervision

**Date:** 2026-03-14
**Context:** Architecture discussion on bridging from human-initiated `/orchestrate` to autonomous agent systems. Driven by concrete use case (AppSignal autonomous cleanup) and informed by Gastown analysis.

---

## The Autonomy Ladder

Levels of autonomy for autonomous bug-fixing (AppSignal use case as exemplar):

| Level | Description | Human Role | What's Needed |
|---|---|---|---|
| **0 — Today** | Human sees alert, manually runs `/orchestrate` | Initiator + reviewer | Nothing new |
| **1 — Triage patrol** | Agent scans channel, classifies alerts, creates tickets | Picks which to fix | Patrol loop + classification |
| **2 — Auto-fix easy ones** | Agent triages AND auto-spawns fixers for "easy" alerts | Reviews PRs only | GUPP + gating hooks |
| **3 — Full patrol** | Scans, triages, fixes, ships PRs, escalates hard ones | Reviews PRs, handles escalations | Supervisor + timeout enforcement |
| **4 — Systemic** | Groups related alerts, proposes root-cause fixes | Strategic decisions only | Pattern detection across alerts |

### Starting Point: Level 1 (Triage → Ticket Creation)

The key insight: **triage IS the hard problem, execution is already solved.** If you have a well-scoped ticket, `/orchestrate` can close it. So the autonomous system is really:

```
[patrol loop] → [triage + ticket creation] → [ticket queue] → [autonomous /orchestrate per ticket]
```

Each AppSignal alert is naturally scoped (one error, one stack trace). The ticket just needs: error description, stack trace, frequency, affected endpoint. `/orchestrate` already handles design → implement → verify → ship PR.

The bridge from Level 1 to Level 2 is GUPP: a cron or loop checks for new tickets, and if one exists, spawns a worktree agent with the ticket as input. The agent doesn't ask permission — the ticket IS the permission. Gating hooks ensure it can only produce PRs, not merge them.

### Existing Capabilities That Map to This

- `/fire-watch` — already does channel scanning + ticket creation + RCA spawning. Closest existing skill.
- `/pulse-sel` — daily triage patrol with parallel subagents. Proves the patrol pattern works.
- `/ci-monitor` — loop-based auto-fix. Proves the "fix and re-check" loop works.
- `/permission-manager` — tmux-based supervisor POC. Proves decoupled supervision works.

---

## Agent Roles / Classes

### The Concept

A role is a named configuration that selects which hooks, injections, gating rules, and permission strategy apply to an agent session. Trivial foundation, powerful groundwork.

### Role Definition Structure

```json
{
  "role": "implementer",
  "sandbox": true,
  "dangerouslySkipPermissions": true,
  "hooks": {
    "PreToolUse": ["no-push-to-master.sh", "no-ci-config-edits.sh"],
    "SessionStart": ["inject-role-context.sh"],
    "PreCompact": ["flush-partial-results.sh"]
  },
  "injections": ["autonomous-implementer.txt"],
  "autonomy_profile": "semi-autonomous",
  "timeout_minutes": 60,
  "escalation_channel": "telegram"
}
```

### Proposed Roles

| Role | Purpose | Sandbox? | Key Constraints | Autonomy |
|---|---|---|---|---|
| **implementer** | Fix bugs, build features in worktree | Yes | Can't push to master, can't merge PRs, can't modify CI | High (GUPP) |
| **patrol** | Scan channels, classify, create tickets | Partial (needs network for MCP) | Can't edit code, can't create branches | High (loop) |
| **fixer** | Targeted bug fix from ticket | Yes | Can't merge, bounded time/cost | High (GUPP) |
| **browser-qa** | Ranger verification, screenshot comparison | No (needs Chrome, network) | Read-only code access, can't push | Medium |
| **supervisor** | Monitor agent health, escalate | No (needs tmux, Telegram) | Can't edit code, can nudge/restart agents | High (patrol loop) |
| **interactive** | Pair programming with operator | No | Standard builtin permissions | Low (collaborative) |
| **designer** | Spec work, architecture, planning | No | Can write specs, can't write code | Medium |

### How Roles Compose With Existing System

- **Session injection** (`session-mode-inject.sh`) — role selects which `.txt` files are injected
- **Hooks** (`~/.claude/hooks/`) — role selects which PreToolUse gating hooks apply
- **Autonomy profiles** — role selects default profile (can be overridden per-task)
- **Worktree provisioning** — `wtp add --role implementer` configures the worktree's hooks automatically

### Implementation Path

1. Role definitions as JSON files in `~/.claude/agent-roles/`
2. A `spawn-role.sh` script that: reads role JSON, configures hooks, sets up injections, launches Claude Code session
3. Integrate with existing `wtp` tooling for worktree-based roles
4. Integrate with supervisor for health monitoring per-role

---

## Permission System: Tiered by Role

### The Insight

Instead of one binary choice (builtin perms vs custom hooks), tier the permission strategy per role:

| Role Category | Sandbox? | Permission Strategy | Safety Net |
|---|---|---|---|
| **Sandboxable** (implementer, fixer) | Yes | `dangerouslySkipPermissions` + hooks | Sandbox bounds blast radius |
| **Network-dependent** (patrol, browser-qa) | No | Builtin perms + layered gating hooks | Builtin perms as backstop |
| **System-access** (supervisor) | No | Builtin perms + gating hooks | Builtin perms as backstop |
| **Interactive** (operator sessions) | No | Standard builtin perms | Human is watching |

### Key Question: Does Sandbox Block Docker?

If an implementer needs `docker compose exec web bundle exec rspec` (talks to local Docker daemon), does sandbox break that? **This is make-or-break for sandboxing implementer agents.** Investigation pending.

If sandbox blocks Docker:
- Implementers can't be sandboxed (need Option 2: builtin + hooks)
- The sandboxable tier shrinks to only agents that don't need Docker

If sandbox allows Docker:
- Implementers get the lightest-weight safety (sandbox + hooks)
- Most coding agents can be fully sandboxed

### Gating Hook Examples

```bash
# no-push-to-master.sh (PreToolUse)
# Blocks: git push to master/main
if [[ "$TOOL" == "Bash" && "$COMMAND" =~ "git push" && "$COMMAND" =~ (master|main) ]]; then
  echo "BLOCKED: Cannot push to master/main"
  exit 1
fi

# no-merge-pr.sh (PreToolUse)
# Blocks: gh pr merge
if [[ "$TOOL" == "Bash" && "$COMMAND" =~ "gh pr merge" ]]; then
  echo "BLOCKED: Cannot merge PRs — open only"
  exit 1
fi

# bounded-scope.sh (PreToolUse)
# Blocks: edits outside the working directory
if [[ "$TOOL" == "Edit" || "$TOOL" == "Write" ]]; then
  if [[ ! "$FILE_PATH" =~ ^"$WORKING_DIR" ]]; then
    echo "BLOCKED: Cannot edit files outside working directory"
    exit 1
  fi
fi
```

### Code Coverage Gating (Idea)

A hook that blocks PR creation if test coverage drops below threshold. Concrete example of the executable gating pattern. Feasibility depends on having fast coverage checks available.

```bash
# coverage-gate.sh (PreToolUse)
# Blocks: gh pr create if coverage regressed
if [[ "$TOOL" == "Bash" && "$COMMAND" =~ "gh pr create" ]]; then
  coverage_delta=$(check_coverage_delta)
  if [[ $coverage_delta -lt 0 ]]; then
    echo "BLOCKED: Coverage regressed by ${coverage_delta}%. Fix before opening PR."
    exit 1
  fi
fi
```

---

## Supervisor / Watchdog Architecture

### Current State

| Skill | What It Supervises | Limitations |
|---|---|---|
| `/permission-manager` | Permission prompts in tmux | One-shot, no learning, narrow scope |
| `/ci-monitor` | CI status for current PR | Specific to CI, no general health |
| `/fire-watch` | Slack channels for issues | Creates tickets, no agent monitoring |
| `/pulse` / `/pulse-sel` | Project health artifacts | One-shot scan, no continuous monitoring |

**Cross-system gaps identified:**
- No agent health monitoring (stuck? idle? off-course?)
- No feedback loop (permission-manager rules are static)
- No cross-skill coordination (pulse finds blockers, fire-watch creates tickets — neither updates the other)
- No unified escalation bus (Telegram, stdout, SendMessage — inconsistent)
- System built for human-in-the-loop, not autonomous supervision

### MVP Supervisor Design

**Core function:** Watch agent sessions via tmux, detect problems, escalate consistently.

**Detection heuristics:**
1. **Stuck agent** — no tmux output for N minutes (configurable per role)
2. **Error loop** — agent retrying same command 3+ times
3. **Off-course** — agent editing files outside expected scope (requires role definition)
4. **Permission blocked** — existing `/permission-manager` detection
5. **Context filling** — detect compaction warnings, trigger graceful handoff

**Actions:**
1. **Nudge** — send tmux keystroke to unstick (like Gastown's GUPP nudge)
2. **Escalate** — Telegram notification with: which agent, what's wrong, suggested action
3. **Restart** — kill session, spawn fresh with same role + hook (GUPP-style auto-restart)
4. **Log** — record incident for pattern analysis

**Escalation Bus:**
- All supervisory alerts go through Telegram with consistent format
- Format: `[AGENT:{role}:{name}] {severity} — {description} — {suggested_action}`
- Severities: INFO (logged, no notification), WARN (Telegram), CRITICAL (Telegram + repeated)

**Connection to Gastown:**
- This is our "Witness" — watches agents and helps unstick them
- Simpler than Gastown's hierarchical Daemon → Boot → Deacon → Witness chain
- Sufficient for 1-5 agents. Scale up hierarchy if/when needed.

**Implementation path:**
1. Shell script or Go tool that: lists tmux sessions, checks output timestamps, runs heuristics
2. Runs on cron (every 2-5 min) or as a long-running patrol
3. Integrate with Telegram MCP for escalation
4. Integrate with role definitions for per-role timeout/scope expectations

---

## Connecting It All: The AppSignal Use Case

### Architecture

```
[Patrol Agent]                    [Supervisor]
  role: patrol                      role: supervisor
  loop: every 30min                 loop: every 2min
  scans: #appsignal-alerts          watches: all agent sessions
  creates: Linear tickets           detects: stuck, off-course, blocked
  escalates: hard-to-classify       escalates: via Telegram
       |
       v
[Ticket Queue]
       |
       v  (cron checks for new tickets)
[Fixer Agent]
  role: fixer
  sandbox: true
  spawns: worktree per ticket
  pipeline: /orchestrate (triage → implement → verify → ship PR)
  bounded: 60min timeout, can't merge
  output: CI-green PR ready for human review
       |
       v
[Human Operator]
  reviews: PRs
  merges: approved fixes
  goal: AppSignal channel gets clean over time
```

### What Exists vs What To Build

| Component | Exists | Gap | Effort |
|---|---|---|---|
| Slack scanning | `/fire-watch`, egregore MCP | Need AppSignal-specific classification | Small |
| Ticket creation | `/fire-watch` creates Linear tickets | Need to adapt for AppSignal alert format | Small |
| Fix pipeline | `/orchestrate` | Need autonomous trigger (no human initiation) | Medium |
| Worktree isolation | `wtp add` | Need auto-provisioning per ticket | Small |
| Executable gating | Claude Code hooks | Need role-specific PreToolUse rules | Small |
| Propulsion (GUPP) | `autonomy-profile.txt` | Need stronger "execute without asking" framing | Small |
| Workflow state | `task-plan.md` | Need incremental checklist convention | Small |
| Supervisor/timeout | `/permission-manager` (narrow POC) | Need general agent health monitoring | Medium |
| Pattern grouping | Nothing | Level 4, defer | Large |

### Phased Rollout

**Phase A — Foundation (agent roles + gating):**
1. Define role JSON format
2. Build `spawn-role.sh`
3. Write gating hooks for implementer/fixer roles
4. Test: manually spawn a fixer agent with role config, verify gating works

**Phase B — Patrol (triage → tickets):**
1. Adapt `/fire-watch` for AppSignal channel
2. Classification heuristics for alert types
3. Auto-create Linear tickets from classified alerts
4. Test: run patrol, verify tickets are well-scoped

**Phase C — Auto-fix (GUPP bridge):**
1. Cron/loop that checks ticket queue for new "easy" tickets
2. Auto-spawn fixer agent in worktree with ticket as input
3. Fixer runs `/orchestrate` autonomously
4. Test: create a ticket manually, verify fixer produces PR

**Phase D — Supervisor:**
1. Build MVP supervisor (tmux health check + Telegram escalation)
2. Integrate with role definitions for per-role expectations
3. Test: intentionally stuck agent, verify escalation fires

**Phase E — Close the loop:**
1. Connect patrol → ticket queue → auto-fix → supervisor
2. Run end-to-end: alert appears in Slack → PR appears in GitHub
3. Human reviews and merges
4. Measure: is the AppSignal channel getting cleaner?
