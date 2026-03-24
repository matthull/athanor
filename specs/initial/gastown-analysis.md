# Gastown Analysis: Lessons for Our Architecture

**Date:** 2026-03-14
**Source:** Steve Yegge's "Welcome to Gas Town" article + ~/code/gastown repo exploration
**Context:** Architecture discussion comparing Gastown's multi-agent orchestration patterns against our `/orchestrate` system. Goal: identify what to steal, what to skip, and what to watch.

---

## What Gastown Is

Gas Town is a multi-agent orchestration system for 20-30+ concurrent Claude Code instances. It's Steve Yegge's fourth orchestrator (v4, written in Go, ~67 internal packages). It uses tmux as primary UI, Beads (Git-backed JSONL issue tracker) as its data plane, and Dolt SQL as backend.

### Key Concepts

- **Town**: HQ directory (`~/gt`), manages all projects
- **Rigs**: Individual project repos under Gas Town management
- **Polecats**: Ephemeral per-rig workers that spin up on demand, produce MRs, then get decommissioned
- **Refinery**: Merge queue agent — intelligently merges all changes to main, one at a time
- **Witness**: Per-rig supervisor that watches polecats and helps unstick them
- **Deacon**: Daemon beacon — runs patrol loop, propagates heartbeats downward
- **Dogs**: Deacon's helpers for maintenance and handyman work
- **Boot**: Special dog that checks on the Deacon every 5 minutes
- **Crew**: Long-lived per-rig agents for interactive work (design, back-and-forth)
- **Mayor**: Main agent you talk to — concierge and chief-of-staff
- **Convoys**: Work-order units that wrap related work for tracking

### MEOW Stack (Molecular Expression of Work)

1. **Beads**: Atomic work units (issues) in Git-backed JSONL
2. **Epics**: Beads with children (parallel by default, explicit dependencies for sequencing)
3. **Molecules**: Chained workflows of Beads — arbitrary shapes, loops, gates, Turing-complete
4. **Formulas**: TOML source form for workflows, "cooked" into protomolecules, instantiated into wisps/mols
5. **Wisps**: Ephemeral Beads — in database but not persisted to Git. Burned after use.
6. **Protomolecules**: Templates — immutable, embedded in binary, instantiated as wisps

---

## Fundamental Differences: Gastown vs Our System

| Dimension | Gastown | Our Architecture |
|---|---|---|
| **Scale** | 20-30 concurrent agents | 1-3 agents, mostly single |
| **Optimization target** | Throughput ("fish in barrels") | Attentional cost over time |
| **Error tolerance** | High — "some fish escape, some bugs fixed 2-3 times" | Low — spec compliance matters |
| **Operator stance** | PM (idea compiler) | Collaborative engineer |
| **Work style** | Vibe coding, chaos accepted | Deliberate, verified |
| **Verification** | Minimal (Refinery handles merge conflicts) | Multi-layered (spec-check, qa-plan, ranger, task-review) |
| **Identity model** | Agents are persistent Beads; sessions are cattle | Agent identity ≈ conversation context + plan doc |
| **State persistence** | Beads in Git (JSONL + Dolt SQL) | Markdown files in task dirs |

Gastown runs a factory. We run a craft workshop. Both valid — but alignment problems differ.

---

## What Gastown Does Better Than Us

### 1. Executable Gating > Prompt Injection

**The Pattern:** Claude Code hooks as executable gates that BLOCK operations, not just prompt injection that SUGGESTS behavior.

**How Gastown Does It:**
- PreToolUse hooks intercept dangerous commands before execution
- Role-specific overrides: Witness cannot create persistent molecules (hook blocks `bd mol pour`)
- Not a suggestion in the prompt — the CLI itself rejects the command
- Hook configuration cascades: DefaultBase() → on-disk hooks-base.json → role-specific overrides → on-disk per-target overrides

**Our Current Approach:** "Injection in depth" — repeat behavioral instructions at multiple prompt injection points. Works for guidance ("prefer X over Y") but unreliable for hard prohibitions ("never do X").

**What To Steal:** Use hooks for hard constraints, prompting for soft guidance. Examples:
- Implementer agent: PreToolUse hook blocks `git push` to master, blocks modifications to CI config
- Fixer agent: PreToolUse hook blocks `gh pr merge`
- Patrol agent: PreToolUse hook blocks file writes outside designated dirs

**Assessment:** High impact, low effort. This is the single most actionable insight.

### 2. Identity Survives Sessions (Agents ≠ Sessions)

**How Gastown Does It:**
- Agent identity is a Bead (persistent, Git-backed) with: role pointer, mail inbox, hook (current work), orchestration state
- Sessions are ephemeral cattle thrown at persistent work
- When session dies, agent's work/state/identity persist. New session picks up hook and continues.
- `gt seance` — agents can communicate with their predecessors via Claude Code's `/resume`

**Our Current Approach:** Plan doc (`task-plan.md`) updated at least once per phase. Agents figure out where to pick up by reading it.

**Assessment:** We already have this at a functional level. Plan doc serves as durable identity. Agents resume well enough. Refine reactively as issues arise, don't over-invest.

### 3. GUPP: Physics Over Politeness

**The Pattern:** Gastown Universal Propulsion Principle — "If there is work on your hook, YOU MUST RUN IT."

**Key Design Elements:**
- Framing is physics, not social obligation: "You are a piston in a steam engine"
- Agents cannot choose to be idle — always a next action (work hook, escalate, or exit)
- "The Idle Polecat Heresy" — explicit, bolded prohibition patterns early in templates
- GUPP Nudge — 30-60 second startup poke if agent doesn't self-start
- Hierarchical heartbeats: Daemon → Boot → Deacon → Witness → Polecats

**How It Works Technically:**
```
Session starts → gt prime --hook injects context
  → Detects work on hook
  → Outputs "## AUTONOMOUS WORK MODE" header
  → "This is physics, not politeness. Gas Town is a steam engine — you are a piston."
  → Agent executes immediately, no waiting for confirmation
```

**Our Current Approach:** Autonomous mode injection via `autonomy-profile.txt`. Less aggressive — agents sometimes still wait for confirmation.

**What To Steal:** The framing. Add GUPP-like language to autonomous mode: "If a phase is incomplete, execute it. This is not optional. Do not wait for confirmation unless the task plan specifies operator approval." Value is clarity for human operators even if LLM compliance benefit is uncertain.

**Assessment:** Medium impact, low effort. Can't hurt, clarifies intent.

### 4. Molecules: Workflow as Data, Not Prompt

**How Gastown Does It:**
- Work defined as chains of Beads (the molecule) — immutable workflow definitions
- Agent walks the chain, one step at a time, checking off issues
- Workflow survives crashes, compactions, restarts — it's in the database, not the prompt
- Each step produces an activity feed entry automatically
- Nondeterministic Idempotence (NDI): nondeterministic path, deterministic outcome

**Our Current Approach:** `task-plan.md` with phases. Phase results appended on completion. Free-form prose, not structured.

**Gap Analysis (detailed):**

Where we're equivalent to NDI:
- Single source of truth (task-plan.md)
- Deterministic skip logic (has results block? skip)
- Operator-initiated resumption
- Recorded gate decisions

Where NDI has real advantages:
1. **Partial completion within phases** (HIGH pain) — if IMPLEMENT crashes after 2/4 work units, we re-run whole phase. Fix: incremental checklist convention in results blocks.
2. **Context flush contract** (MEDIUM) — no guarantee agent writes results before context dies. Fix: write partial results eagerly after each work unit.
3. **Unstructured results** (LOW at our scale) — human reads every plan. Machine-parseable matters at 20 agents, not 3.
4. **Retry history** (LOW at our scale) — operator is watching, retries visible in conversation.

**What To Steal (without infrastructure change):**
Instead of writing results once at phase end, write incrementally:
```markdown
### IMPLEMENT
**Results:**
- [x] WU-1: Add notification preferences model (branch: sal-89, commit: abc123)
- [x] WU-2: Slack App Home UI (branch: sal-89, commit: def456)
- [ ] WU-3: Preference persistence API
- [ ] WU-4: Integration tests
Confidence: (pending)
Gate outcome: (pending)
```

**POC verification:** Try on next `/orchestrate` run. Verify agents resume cleanly after compaction without redoing completed work units.

**Assessment:** Convention change, not infrastructure change. High value for crash recovery.

### 5. Hierarchical Supervision (Witness/Deacon Pattern)

**How Gastown Does It:**
- Cascading health check: Daemon → Boot → Deacon → Witness → Polecats
- Each level monitors level below, nudges/restarts as needed
- Deacon patrol runs every few minutes
- Exponential backoff when no work found
- Multi-vector escalation routed by severity (MEDIUM → bead+mail, HIGH → +email, CRITICAL → +SMS)
- Auto-re-escalation if stale (default 4h)

**Our Current Approach:**
- `/permission-manager` — scans tmux for permission prompts, approves/escalates (narrow POC)
- `/ci-monitor` — watches CI, auto-fixes (specific to CI)
- `/fire-watch` — scans Slack, creates tickets (specific to launch monitoring)
- No general-purpose agent health monitoring

**What's Missing (confirmed by subagent analysis):**
- Agent health monitoring ("is this agent stuck? idle for 10 min? going in circles?")
- Off-course detection ("agent was supposed to fix a bug but is refactoring 15 files")
- Unified escalation bus (currently: Telegram, stdout, SendMessage — no consistency)
- Feedback loop (permission-manager rules are static, don't learn)

**Assessment:** NOT overkill — solves a specific painful operator problem. Decoupled layer, can leverage tmux like Gastown. Start with MVP: tmux health check + timeout enforcement + Telegram escalation bus.

---

## What We Do Better Than Gastown

### 1. Verification Discipline
Gastown optimizes for throughput and accepts slop. Our multi-layered verification (spec-check, qa-plan, task-review, ranger) ensures spec compliance. Gastown's Refinery handles merge conflicts but doesn't verify spec compliance.

### 2. Attentional Cost Optimization
Our core goal — minimize total human attentional cost over time — is more sophisticated than "keep the engine running." We think about false positive vs. false negative costs. Gastown's answer is "throw more agents at it."

### 3. Specification Leverage
Gastown treats work as "guzzoline" — fluid, lossy. We treat specs as source of truth. This produces more reliable output even if slower.

---

## Beads Assessment

**Recommendation: WATCH (don't adopt, don't skip)**

At 1-3 agents, markdown + git gives us everything we need. Beads' killer features (queryable work queue, multi-agent routing, dependency DAGs, capability-based routing) only matter at 5-10+ concurrent agents.

| Dimension | Markdown | Beads | Winner at our scale |
|---|---|---|---|
| Simplicity | Excellent | Good | Markdown |
| Querying | Poor | Excellent | Beads (but don't need it) |
| Human readability | Excellent | Poor | Markdown |
| Multi-agent visibility | Poor | Excellent | Beads (but don't need it) |
| Learning curve | None | 4-6 hours | Markdown |
| Adoption cost | None | 8-12 hours | Markdown |

**Crossover point:** ~5 concurrent agents, or cross-repo dependencies.

**One interesting concept:** Wisps (ephemeral beads that auto-delete) for patrol transient state. Our patrols generate state that pollutes git if persisted but is lost if not. Currently handled with `/tmp/` files. Not a pain point yet.

---

## Gastown's Prompt Engineering Patterns (Reference)

### Multi-Phase Prompt Composition (`gt prime`)
1. **Beacon** — Session metadata: `[GAS TOWN] polecat rust (rig: gastown) <- mayor • timestamp • assigned:gt-abc12`
2. **Role Context** — Embedded templates (polecat-CLAUDE.md, witness-CLAUDE.md)
3. **Beads Integration** — `bd prime` outputs beads workflow context
4. **Memory Injection** — Loads from beads KV store
5. **Mail Injection** — `gt mail check --inject` injects pending messages
6. **Hooked Work Detection** — Triggers AUTONOMOUS MODE
7. **Molecule Context** — Attached formula steps inline
8. **Checkpoint Recovery** — Previous session state for crash recovery

### Compaction Optimization
After compaction, `gt prime` skips full setup — brief "Context recovered" message instead of re-announcing AUTONOMOUS WORK MODE. Prevents agents from re-initializing and losing momentum.

### Session State Detection
```
1. Handoff marker (previous session handed off to me)
2. Checkpoint (previous session crashed with state)
3. Hooked work (active assignment)
4. Normal startup (no prior state)
```

### Alignment Mechanisms Summary

| Mechanism | Level | How It Works |
|---|---|---|
| **Hooks (executable)** | CLI | PreToolUse blocks dangerous commands |
| **Session Startup (context)** | Prompt | `gt prime` outputs AUTONOMOUS WORK MODE |
| **Role Templates** | Prompt | Polecat template has "THE IDLE POLECAT HERESY" |
| **Beads (immutable)** | Data | Work defined in beads, not human conversation |
| **Session State Detection** | Behavior | Compaction avoids re-initialization |
| **Property Layers** | Config | Role config merged from 4 layers (wisp → rig → town → system) |
| **Escalation Protocol** | Coordination | Multi-channel, severity-routed, auto-retry at 4h |
| **Memory Injection** | Context | Beads KV store persists agent memory |
| **Work Attribution** | Telemetry | OTEL tracks which bead/mol each agent executed |
