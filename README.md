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
specs/
  spec.md          -- full system specification (start here)
  kadmon.md        -- operational runbook and bootstrap procedures
  issues.md        -- known issues discovered in operation
  whisper/         -- whisper CLI spec (inter-agent communication tool)
  cli/             -- ath CLI spec (athanor management tool)
```

The system spec at [`specs/spec.md`](specs/spec.md) is the canonical reference for principles, architecture, vocabulary, and design decisions.

## Also in This Repo

**Whisper** — A Go CLI for reliable inter-agent communication via tmux. Three commands: `whisper send` (8-step delivery protocol), `whisper idle` (2-check idle detection), `whisper wait-and-send` (combined). Built and operational.

## Status

The athanor is in active use with multiple concurrent instances (task-based and state-based). The core loop works: opera get inscribed, azers execute them, maruts supervise, the trail accumulates, and the system learns. Infrastructure is being built incrementally as gaps surface in operation.
