# Nudge Contract — Presence-Driven Role Activation

**Status:** Agreed (artifex confirmed 2026-04-20)
**Created:** 2026-04-20
**Source:** Interactive design session (artifex + azer, nudge-system-design opus)

---

## Problem

Presence-driven roles (perceiver, attendant) are kindled as singletons under liveness timer supervision. The liveness timer ensures sessions exist but does not activate them. Once a role finishes its initial orientation, nothing triggers it to act. The perceiver runs headless on a separate timer but its downstream consumer (the attendant) receives no notification of updates. Result: kindled roles sit idle unless directly interacted with.

## Design Principle

Two things activate presence-driven roles: **new signals** and **time passing**. Both are addressed by a single mechanism: the attunement intake timer whispers the attendant after each processing pass. The time hook (`status-inject.sh`) provides temporal awareness automatically on every prompt submission — no explicit time-passing logic is needed.

**Feedback-over-specification.** This contract is deliberately minimal. The artifex's interoceptive response to the system's behavior is the primary evolutionary driver. Design light, iterate from lived experience.

---

## Activation Channels

### Channel 1: Clock + Signal Nudge (timer → perceiver → attendant)

**Mechanism:** The attunement intake timer (30-min cadence, systemd) runs the perceiver's headless intake pass. After processing signals and updating the corpus, the perceiver whispers the attendant.

**Whisper content:** The perceiver composes a brief nudge with two parts:

1. **Signal summary** (when new signals were processed): What changed — new patterns, state shifts, updated domains. One or two sentences.
2. **Immediate-relevance flag** (when something is time-sensitive): Explicit flag when the perceiver judges something warrants near-term attendant action. Examples: "Meds not confirmed by 10am", "Calendar event in 45 minutes", "Energy shift — afternoon fatigue noted."

When no new signals were processed (intake gate found nothing), the perceiver still whispers a bare clock nudge: "No new signals. Clock check."

**Whisper mechanism:** `ath whisper send attendant-<athanor> "<nudge content>"`

**Cadence:** Every 30 minutes (matches attunement-intake.timer). This is the attendant's heartbeat.

**30-min lookahead discipline:** Because the attendant only wakes every 30 minutes, it must reason about what will happen in the *next* 30 minutes, not just what's true now. At 9:30, "meds deadline is 10:00" means surface now — there won't be another chance before the deadline. This lookahead is the attendant's responsibility, not the nudge system's.

### Channel 2: Attendant Behavior on Nudge

On every incoming whisper (whether from perceiver, artifex, or any source), the attendant:

1. **Re-reads** `portrait.md` and `current-state.md` — catches any updates since last activation.
2. **Checks the time** (injected by `status-inject.sh` hook) against active vigils and known temporal context (calendar, commitments, deadlines).
3. **Runs the reasoning pattern** (generative panel → synthesis → double-fire calcinatio, defined in `shared/attendant.md § The Reasoning Pattern`) to determine whether any action is warranted.
4. **Acts or holds.** If the double-fire passes, surface via Telegram. If not, groom internal state and wait for next nudge.

### Channel 3: Perceiver Immediate Escalation

Outside the 30-min cycle, if the perceiver is invoked interactively (consultation mode) and discovers something with genuine urgency, it whispers the attendant immediately. This is rare — the 30-min cycle handles normal flow. Examples of genuine urgency: safety-relevant signal, time-critical commitment the attendant doesn't know about.

---

## Edge Cases

**Attendant not kindled.** `ath whisper send` fails silently if the target window doesn't exist — no error, no retry. The liveness timer will kindle the attendant on its next pass (5-min cadence) if the role is in kindled state. If the perceiver has genuinely urgent information and the attendant isn't available, the perceiver's own escalation path (Telegram notify to the artifex) is the fallback.

**Between timer ticks.** The attendant may receive direct messages from the artifex via Telegram at any time. These are handled by the attendant's normal responsive mode — no timer dependency. The 30-min nudge is the floor for autonomous activation, not a ceiling.

**Perceiver intake finds no signals.** Still whispers the attendant with a bare clock nudge. This is purely a heartbeat — its value is the wake-up, not the content. Time passing is itself a trigger; the attendant needs to check vigils against the clock even when nothing new has arrived.

**Multiple whispers in quick succession.** The attendant processes each as it arrives. No debouncing needed — the reasoning pattern's double-fire naturally filters redundant action.

---

## Implementation Notes

**Changes required:**

1. **Attunement intake protocol** (`specs/life-domains/attunement-intake-protocol.md`): Add a final step after "Finalize" — check if attendant is kindled, compose nudge, whisper.
2. **Perceiver role file** (`shared/perceiver.md`): Add attendant notification to the perceiver's responsibilities (it already has the whisper relationship documented).
3. **No new timers or scripts.** The existing attunement-intake timer carries the nudge. The existing whisper mechanism delivers it.

**What this does NOT specify:**

- What the attendant surfaces or when (that's the attendant's geas + double-fire)
- How vigils are created or managed (that's the attendant's operational concern)
- Portrait or current-state format (that's the perceiver's concern)
- Liveness timer behavior (unchanged — lifecycle management only)

---

## Evolution

This contract will be revised based on lived experience. Specific things to watch:

- **30-min cadence** — is it too sparse for time-sensitive vigils (meds, appointments)? Too frequent for context budget? The artifex's felt experience determines adjustment.
- **Nudge content richness** — does the perceiver's signal summary help the attendant, or is a bare "clock check" sufficient? Simplify or enrich based on what the attendant actually uses.
- **Immediate-relevance judgment** — is the perceiver good at flagging urgency, or does it over/under-flag? Calibrate through use.
