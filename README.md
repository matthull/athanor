# The Athanor

*The alchemist's furnace that burns continuously — maintaining constant temperature for sustained transformation.*

The athanor is an agent orchestration system built on Claude Code. It accepts goals stated in stakeholder-value terms and produces satisfying results — doing everything an agent system reasonably can, without micromanagement, communicating and escalating intelligently to keep the operator informed of what genuinely needs them.

For its operator (who has ADHD and autism), this is not a productivity tool. It is executive function infrastructure — the accommodation that makes sustained professional work possible.

## Motivation

LLM agents are stochastic. Sessions die. Agents have completion bias. Individual agents are unreliable. But a *system* of agents can be reliable if the architecture absorbs these failures gracefully.

The athanor makes progress **inevitable** — not through individual agent reliability, but through:
- **The trail** — discharged opera form a persistent record. When a session dies, the trail holds the progress.
- **The geas** — a behavioral compulsion that channels completion bias toward genuine fulfillment, making escalation feel like compliance rather than failure.
- **Chain of evidence** — agents proceed with evidence or escalate. No guessing, no hallucinating past gaps.

## Key Design Components

**Opera** — Units of work. Each opus has intent (what's true for witnesses when done), boundary (agent vs. operator responsibility), and context. Inscribed when immediately actionable, discharged with outcome, proof, and reflection.

**Geas** — Binding behavioral compulsion, not a checklist. "Pursue this with integrity and abundantly satisfy the goal" produces verification discipline, communication, and clean escalation without enumerating every procedure.

**Agents** — Azers (workers) execute opera in isolated environments. Maruts (supervisors) keep the operational loop turning — dispatching azers, monitoring progress, cleaning up, and narrating significant turns to the operator.

**The Trail** — The sequence of discharged opera. The authoritative record of what actually happened, not what was planned. Status, precedent, and learning all come from the trail.

**Adaptability through essentialism** — The athanor carries only principles, never procedures. Procedures belong in the environment (CLAUDE.md, CI, hooks). This means the athanor can interface with any project without conflicting with existing workflows.

**Abundant satisfaction** — The completion standard is "would the operator feel this is completely handled?" not "were the requirements technically met?" Agents close the tail — coordination, communication, documentation — so the operator engages only where genuine human judgment is required.

## Repository Structure

```
cmd/ath/             -- CLI entry point
internal/
  cli/               -- command implementations
  athanor/           -- home directory, instance, config logic
  tmux/              -- tmux interaction library
specs/
  spec.md            -- system specification (start here)
  cli/spec.md        -- ath CLI specification
  kadmon.md          -- operational runbook (first athanor instance)
  issues.md          -- known issues discovered in operation
```

The athanor home (`~/athanor/`) lives outside this repo — it contains shared materia and all athanor instances. This repo is the CLI source and system spec.

The system spec at [`specs/spec.md`](specs/spec.md) is the canonical reference for principles, architecture, vocabulary, and design decisions.

## The `ath` CLI

The `ath` binary is the operational backbone. Install with `make install` (puts it at `~/.local/bin/ath`).

### Setup

```bash
# One-time: create the athanor home and populate shared materia
mkdir -p ~/athanor/shared ~/athanor/athanors
# Copy role files (AGENTS.md, marut.md, azer.md, muster.md, opus.md) into ~/athanor/shared/

# Optional: install zsh completions
ath completion zsh > ~/.zsh/completions/_ath
```

### Workflow: Spin up a new athanor

```bash
# 1. Create the instance
ath init bugsnag --project ~/code/musashi

# 2. Edit the magnum opus — define the goal, witnesses, and context
vim ~/athanor/athanors/bugsnag/magnum-opus.md

# 3. Kindle the marut — it reads its materia and starts the operational loop
ath kindle bugsnag

# 4. Check on it
ath status
ath status bugsnag
```

### Workflow: Marut musters an azer

From within a running marut session (or manually):

```bash
# Muster an azer for a charged opus
ath muster 2026-03-25-fix-nil-error.md --athanor bugsnag --dir ~/code/musashi-worktree

# Check opera status
ath opera bugsnag

# Clean up after the azer finishes
ath cleanup azer-fix-nil-error
```

### Workflow: Communication between agents

```bash
# Send a message to a crucible
ath whisper send marut-bugsnag "Status check — are you making progress?"

# Wait for an agent to be idle, then send
ath whisper wait-and-send azer-fix-nil-error "Your opus has been updated" --timeout 60s
```

### Workflow: Lifecycle management

```bash
# Reforge a crashed/exhausted marut (kills session, relaunches fresh)
ath reforge bugsnag

# Graceful shutdown
ath quiesce bugsnag

# Force shutdown (even with active azers)
ath quiesce bugsnag --force
```

### All commands

Run `ath help` for the full command list, or `ath whisper help` for whisper subcommands.

## Status

The athanor is in active use with multiple concurrent instances (task-based and state-based). The core loop works: opera get inscribed, azers execute them, maruts supervise, the trail accumulates, and the system learns. Infrastructure is being built incrementally as gaps surface in operation.
