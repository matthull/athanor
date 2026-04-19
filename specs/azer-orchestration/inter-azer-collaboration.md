# Inter-Azer Collaboration Protocol

**Status:** Draft — derived from first experiment (2026-04-18)
**Origin:** life-domains MO, coordinate-discussion opus

---

## When to use

Multiple azers are working on related concerns under the same MO (or related MOs) and would benefit from awareness of each other's work. The typical trigger: the artifex or marut notices convergence between independently charged craft azers and wants them to collaborate rather than work in isolation.

This is not the same as the task-lead/implementer pattern (one agent directing others). This is peer collaboration — azers that are equals, each with their own opus, finding where their work intersects and helping each other.

## The coordinator role

A dedicated azer charged with facilitation, not direction. The coordinator:

- **Sees the big picture** that no individual azer can see — status across all collaborators, connections between their work, gaps in coverage
- **Identifies connections** — "your attunement inbox is the interface refinement needs to define their consumer"
- **Bridges communication failures** — relays content when direct whispers fail
- **Answers coordination questions** — "are we shaping or building?" — questions that require cross-azer perspective
- **Identifies parallel work** — "you're unblocked, write your doc while those two negotiate their interface"
- **Monitors context usage** — watches for azers approaching context limits and signals them to wrap up
- **Does not direct work** — azers set their own goals under their own geases. The coordinator facilitates, not manages

## Communication mechanics

### Whisper protocol

Azers communicate via `ath whisper send <target> "message"`. Critical rules:

1. **Whispers are one-way.** When azer-A sends a question to azer-B, B processes it and generates output in B's pane — but B does NOT automatically send the response back. B must explicitly `ath whisper send azer-A "response"`. This is the single most important protocol rule.

2. **No markdown headers in whispers.** Lines starting with `#` get blocked by bash hooks. Use CAPS, numbered sections, or plain text for structure.

3. **Messages arriving during active processing get queued.** They process as the next input when the current turn finishes. This is manageable but creates ordering effects — be aware that a queued message may process before or after a coordinator's synthesis.

### Communication patterns

| Pattern | When to use | How it works |
|---------|------------|--------------|
| **Broadcast** | Significant finding that affects siblings | Azer sends whisper to all relevant siblings with the finding |
| **Direct question** | Need specific information from a sibling | Azer sends targeted whisper, sibling whispers response back |
| **Bilateral negotiation** | Shared interface to define | Two azers exchange proposals and reactions directly. Coordinator monitors but doesn't participate unless needed |
| **Coordinator bridge** | Direct communication failed | Coordinator reads one azer's output and relays to another |
| **Status synthesis** | After a round of activity | Coordinator sends cross-azer status to all, noting connections and next steps |

### Interface lock pattern

When two azers agree on a shared interface (data format, API contract, directory structure):
1. Both explicitly confirm agreement via whisper
2. Either side can request revision — but must whisper the other and get agreement before changing
3. The coordinator notes the lock and ensures both sides honor it

## Kickoff sequence

### Phase 1: Introduction + Status Elicitation

The coordinator sends parallel introductions to all collaborating azers:
- Who the coordinator is and what the experiment is
- What the coordinator already knows about each azer's work (from tmux observation)
- Concrete connections the coordinator sees between their work
- Request for status in (a) what you've produced (b) what you still need to do (c) what you need from others format

This format produces clear, actionable responses and surfaces dependencies immediately.

### Phase 2: Facilitate connections

Based on status responses, the coordinator:
- Identifies bilateral interfaces that need negotiation
- Points unblocked azers at productive solo work
- Answers coordination questions
- Bridges communication when needed

### Phase 3: Monitor and synthesize

As azers work:
- Check tmux panes periodically for progress
- Send status syntheses when significant convergence or divergence occurs
- Watch context levels
- Note when deliverables are complete

## What we learned

### What worked well

- **Azers self-organize when given visibility.** Once they knew about each other's work, they naturally reached out with the right questions.
- **The (a)(b)(c) status format** produced clear responses every time.
- **Bilateral interface negotiation** converged in 2 exchanges — propose, react, confirm.
- **Parallel work identification** kept all azers productive — no one blocked waiting for others.
- **The coordinator as bridge** saved the inbox→refinement interface when the direct whisper failed.

### What needed correction

- **Whisper-back was not obvious.** Azers are habituated to interactive mode where the "person at the prompt" sees their output. Had to explicitly teach "whisper your response back."
- **Markdown in whispers breaks.** Discovered by failure — now documented.
- **Message ordering effects.** When multiple whispers arrive during processing, they queue and may process in unexpected order. Not a blocker but worth awareness.

### Open questions for future experiments

- **Shared artifact?** Currently the big picture exists only in the coordinator's context. Would a shared document all azers can read/write help?
- **Disagreement resolution?** No disagreements occurred in this experiment. The protocol needs an escalation path for when azers have conflicting views.
- **Coordinator lifecycle?** When does the coordinator discharge vs. when do individual azers? In this experiment: all discharged together on artifex signal.
- **Scale?** This worked with 3 azers. Would it work with 5? 8? At what point does the coordinator become a bottleneck?
- **Self-organizing without a coordinator?** Could azers coordinate peer-to-peer if given the protocol rules? Or is the coordinator role essential for the big picture?

## Relationship to other patterns

- **Task-lead/implementer** (azer.md § Context Management): Hierarchical — one agent directs. This pattern is peer-to-peer — a coordinator facilitates equals.
- **Assessment** (AGENTS.md § Assessment Opera): The assessment azer surveys the landscape alone. This pattern is for when multiple azers ARE the landscape and need to coordinate in real time.
- **Dialectical calcinatio** (/calcinatio): Two agents in iterative exchange. This is broader — multiple agents in parallel with a facilitator, using bilateral exchanges as one tool among several.
