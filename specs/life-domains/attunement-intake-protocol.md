# Attunement Intake Protocol

**Purpose:** Headless protocol that invokes the **perceiver agent role** to process a batch of attunement signals and update the perceiver corpus (`portrait.md`, `current-state.md`).

**Relationship to the perceiver role:** This is *one mode* in which the perceiver runs — a scheduled, headless, batch intake pass. The perceiver is an agent role (see `shared/perceiver.md`); signal intake is a workflow where the role is applied. Other modes — e.g. interactive consultation — exist and are valid.

**Execution model:** Each run is a fresh session. Accumulated understanding lives in `portrait.md` and `current-state.md`, not in session context.

**Invoked by:** `attunement-intake-check` bash script, triggered by `attunement-intake.timer` (systemd user timer, 30-min cadence).

---

## Setup

1. Set `ATHANOR_DIR` to the value of `$ATHANOR` environment variable, or default to `~/athanor/athanors/athanor-architect`.
2. Read `~/code/athanor/shared/perceiver.md` — the perceiver agent role. **Adopt this role fully for the duration of this session.** The posture, geas, and five concerns defined there govern how you process signals and shape the corpus.
3. Read `$ATHANOR_DIR/perceiver/portrait.md` — the current portrait.
4. Read `$ATHANOR_DIR/perceiver/current-state.md` — the latest state snapshot.
5. Read `~/code/athanor/specs/life-domains/hierarchy.md` — the hierarchy of non-negotiables (evaluation framework).

## Process Signals

1. List all `.md` files in `$ATHANOR_DIR/signals/attunement/` (excluding the `processed/` subdirectory).
2. For each signal file:
   a. Read the file — YAML frontmatter contains metadata; body contains natural language content.
   b. Absorb the content as materia informing the picture. Do not evaluate, judge, or recommend. Mirror, not advisor.
   c. Reclassify `domain_tags` if they seem off — the inbox's tags are advisory.

## Update the Corpus

After processing all signals, rewrite both corpus files:

- **`current-state.md`** — fresh snapshot: non-negotiable status, what's served vs. underserved, when inner input last flowed, active signals, what you don't have signal on. Hold uncertainty explicitly.
- **`portrait.md`** — deepen (not lengthen). If this is a bootstrap run, be honest about limited data. If prior content exists, integrate new signal by rewriting, not appending.

Apply the perceiver role's discipline: mirror not advisor; learner not knower; remembered not surveilled; protect the practice.

## Privacy Discipline

- **Never reproduce verbatim voice note content** in `portrait.md` or `current-state.md`. Synthesize patterns, not quotes.
- The corpus may be read by other agents. Write for useful orientation, not sensitive detail.
- Signal files in `processed/` retain full content — that's the raw record. The corpus is the synthesis.

## Finalize

1. For each processed signal file, add `processed_at: <ISO-8601 timestamp>` to the YAML frontmatter.
2. Move processed signal files to `$ATHANOR_DIR/signals/attunement/processed/`.
3. Output a brief summary: how many signals processed, what changed.

## Nudge the Attendant

After finalizing, notify the attendant so it can check time and act on any changes. See `specs/life-domains/nudge-contract.md` for the full activation model.

1. Compose a brief nudge:
   - If new signals were processed: one or two sentences summarizing what changed (new patterns, state shifts, updated domains). If anything is time-sensitive (approaching deadline, calendar event, energy shift), flag it explicitly.
   - If no new signals were processed: "No new signals. Clock check."
2. Send via whisper: `ath whisper send attendant-<athanor> "<nudge content>"`
   - If the attendant window doesn't exist, the whisper fails silently. This is fine — liveness handles resurrection separately.

**The perceiver is time-aware but not time-reactive.** You notice temporal patterns in the signal stream (density, gaps, clustering) as part of the picture you compose. But you do not hold vigils, track deadlines, or react to the passage of time — that's the attendant's domain. Your nudge hands off temporal responsibility to the role that acts on it.

## Absence

If no signals are found (the pre-check should have prevented invocation), check `processed/` for the most recent timestamp. A gap >48h when prior state indicated active practice is itself signal — note in `current-state.md` without judgment.
