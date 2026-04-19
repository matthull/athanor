# Inbox Athanor — Concept (Shaping)

**Status:** Shaping — shaped collaboratively with artifex + inter-azer collaboration (2026-04-18)
**Opus:** life-domains/opera/2026-04-18-inbox.md

## Overview

A dedicated athanor with a single Magnum Opus: **every inbound signal reaches the right destination or gets escalated.** The inbox athanor is the system's general-purpose intake — it receives signals from any source and routes them.

## Core Loop

1. **Receive** signal from any input channel
2. **Parse** — understand what's in the signal
3. **Decompose** — split mixed-content signals into discrete items (e.g., a voice note with felt-state + calendar action becomes two items). Original preserved as trail artifact.
4. **Route** each item by disposition (see below)
5. **Trail** everything — what arrived, how decomposed, where each piece went

## Dispositions

| Disposition | Description |
|---|---|
| **Route to MO** | Actionable item with a clear home in an existing Magnum Opus. Delivery mechanism TBD (whisper to marut? file drop? opus inscription?). |
| **Direct execution** | Trivial mechanical action the inbox can handle inline (add calendar event, send message). Has MCP tools available. |
| **Attunement signal** | Inner-state, felt-experience, journaling, interoceptive content. Routed to the attunement inbox for the perceiver agent. See Interface Spec below. |
| **Escalate** | Can't determine routing, low confidence, or novel signal that might spark new work. Escalated to artifex. |

**Waste-nothing principle:** No signal is discarded by the inbox. If routing is unclear, escalate. Only the artifex decides something is noise.

## Input Channels

The inbox is decoupled from any specific input source. Channels feed signals in; each has its own mechanism for delivery.

**Current / near-term:**
- Voice notes (Android watch -> Google Drive -> whisper transcription -> inbox) — the vital inner input axis
- Manual drops via `ath send-inbox` (replaces /send-inbox skill)
- Inter-athanor signals (one athanor surfaces something that belongs to another)

**Future:**
- Email, Slack, calendar events, automated monitors, sensors

## Attunement Inbox Interface (Agreed)

Agreed between azer-inbox and azer-refinement (2026-04-18). Locked bilaterally — changes require whisper and agreement.

### Signal Format

Structured markdown files, one per discrete attunement item. YAML frontmatter + natural language content.

```yaml
---
type: attunement-signal
timestamp: 2026-04-18T09:23:00
source: voice-note
source_ref: /path/to/original/transcript.md
domain_tags: [financial, emotional]
confidence: 0.9
raw_length: 847
---
```

- **domain_tags** — advisory, not authoritative. Sourced from canonical vocabulary in `specs/life-domains/hierarchy.md`. Consuming agent reclassifies as needed.
- **source_ref** — link to full original transcript for surrounding context.
- **confidence** — inbox's confidence this is attunement vs. something else.
- **raw_length** — duration/size of original signal. A 30-second check-in and 10-minute stream carry different weight.
- **Decomposition context** section in body: what else was in the original signal and where it was routed (with refs to sibling items).

### Landing Zone

- `signals/attunement/` — owned by perceiver agent (azer-refinement's design)
- `signals/attunement/processed/` — perceiver moves processed signals here, adds `processed_at` to frontmatter
- Inbox writes only, never reads back or modifies after delivery

### Absence Detection

No separate heartbeat file needed — `last_signal`, `last_voice_note_signal`, and `signal_count_today` can all be derived from the existing signal files via bash (timestamps from frontmatter, file counts). A small bash script or future `ath` subcommand provides these on demand rather than maintaining redundant state.

## Prior Art

The existing village signal routing system (`~/code/village/`) provides patterns to evaluate:

- **Triage router** (systemd path watcher, haiku model) — the "mailroom not filing clerk" principle is sound
- **Bash pre-check pattern** (inbox-check-and-process) — lightweight gate that avoids launching expensive sessions when queue is empty. Worth preserving.
- **Domain registry** — single source of truth for routing targets. The inbox athanor's equivalent would be a registry of active MOs across all athanors.
- **At-least-once delivery** — route first, then mark processed. If routing fails, signal stays in inbox for retry. Worth preserving.
- **Lost-mail fallback** — replaced by the escalation disposition (waste-nothing principle)

## Open Questions (Need Artifex)

1. **MO routing delivery mechanism** — when the inbox routes to an MO, what actually happens? Options: whisper to that MO's marut, file drop in MO directory, opus inscription, or signal landing zone per MO (similar to attunement inbox pattern).
2. **Operational loop** — how does the inbox athanor run? Event-driven (systemd path watcher like old triage router)? Timer-based polling? Beholder pattern? Persistent marut?
3. **Relationship to old village infrastructure** — migrate, replace, or run in parallel during transition?
4. **Voice note pipeline re-establishment** — the old pipeline (watch -> GDrive -> whisper -> signals/inbox/) needs reconnecting to the new inbox athanor. What's the current state of the whisper transcription service?
5. **MO registry** — how does the inbox athanor know what MOs exist across all athanors? Static registry? Dynamic discovery? `ath` CLI query?

## Cross-References

- Hierarchy of non-negotiables: `specs/life-domains/hierarchy.md` — canonical domain tag vocabulary
- Perceiver agent design: `specs/life-domains/exploration/perceiver-agent.md`
- Prior village inbox system: `~/code/village/domains/agentic/triage-router.md`, `~/code/village/domains/agentic/triage-protocol.md`
- Athanor system spec: `specs/spec.md` (Beholder concept)
- Life-domains MO: `~/athanor/athanors/athanor-architect/magna-opera/life-domains/life-domains.md`
