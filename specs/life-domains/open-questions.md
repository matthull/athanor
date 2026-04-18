# Life Domains — Open Questions for Artifex Review

**Created:** 2026-04-18
**Context:** Three-cycle exploration process complete (manifold emanation → dialectical calcinatio → synthesis). These questions survived as genuine decisions requiring the artifex's judgment. Add your responses inline — the azer will incorporate them into the core spec.

---

## 1. Privacy Model for Life-Domain Trail Data

Trail entries from morning regard shifts will contain: email content, calendar details, medication status, emotional observations, family context. Integration discharges will contain intimate pattern material.

**Questions:**
- Where should the life athanor's trail directory live? Excluded from rclone/Syncthing sync?
- Does the portrait live at the same location or separately?
- Which categories of discharge should be excluded from backup scope?
- Should the attendant flag when a discharge contains sensitive material, or treat everything as sensitive by default?

**Your input:**


---

## 2. `ath muster --role` Core Change Approval

The boot prompt in `muster_cmd.go` hardcodes `azer.md`. The attendant requires loading `attendant.md` instead. The clean fix is a `--role` flag (default: `azer`, backward-compatible). The calcinatio confirmed every workaround is architecturally worse.

This IS a core change, and the MO constraint says core changes require your explicit review.

**Approve adding `--role` flag to `ath muster` (and `ath craft`)?**

**Your input:**


---

## 3. Interoception Geas Strength

The attendant must "support without replacing" interoception. Two formulations:

**Strong:** "The attendant does not comment on the artifex's somatic state under any circumstances. Observations about physical state are held for integration review, never surfaced directly."

**Softer:** "The attendant holds observations about physical state for integration. Pattern observations queue for integration — they do not surface directly unless the artifex has already named the pattern themselves."

The strong version is safer (cleaner boundary). The softer version allows the attendant to reflect back what you've already articulated ("you mentioned feeling depleted earlier — that context is held"). 

**Which formulation, or a third option?**

**Your input:**


---

## 4. Graveyard — What to Formally Retire

Candidates for explicit retirement (deliberate "not building this, here's why"):

- [ ] **Continuous biometric monitoring** — ambient surveillance of physical state via sensors
- [ ] **Auto-vigil creation from voice notes** — vigils created without attendant judgment
- [ ] **Cross-athanor awareness (life <> work)** — in first or second slice
- [ ] **Therapeutic register** — emotional guidance, psychological pattern analysis, somatic interpretation

**Confirm which to retire and add a sentence on why for each. Add any others that should be retired:**

**Your input:**


---

## 5. Family Witness Framing

Alex and Vivian are indirect witnesses. Not blocking for the first slice, but shapes second-slice work (kid logistics vigils).

**The question:** Does the attendant ever act on their behalf (e.g., "Alex has a dentist appointment, she needs a ride"), or does it only serve your relationship to them (e.g., "you committed to driving Alex to the dentist")?

The distinction matters for how vigils are scoped and what the attendant's obligations are.

**Your input:**


---

## 6. First Slice Confirmation

The synthesis recommends: **Morning attunement session + medication vigil** as the first slice.

- Morning attunement: daily working-relationship opener. Calendar landscape, overnight vigil results, check-in. Uses existing Calendar + Gmail + Keep + Telegram MCPs.
- Medication vigil: recurring daily vigil with Telegram confirmation as loop-closure evidence.

**Done criterion:** Three consecutive mornings where the attendant surfaces something that would otherwise have been missed, without surfacing anything unwanted. Trail entries are substantive. Medication vigil closes on evidence. Silences are as legible as surfacings.

**What this proves:** Attendant class is real, orientation layer works, regard shape functions end-to-end, prospective memory and temporal horizon prosthetics are operational.

**What this does NOT prove:** Full panel reasoning, all 8 cognitive functions, portrait management, integration work shape, concurrent opera at scale. Those are second-slice territory.

**Confirm this is the right first slice? Any adjustments to the done criterion?**

**Your input:**


---

## 7. Chart/Transit Concept

You mentioned the chart/transit concept from the TaskWarrior era as reaching toward defining areas of life (houses) and how energy flows between them (transits). This connects to the Tipharet-as-life-principle territory.

**Questions:**
- Should the chart/transit vocabulary survive into the life-layer system?
- Does "house" map to something concrete (domains? MOs? areas of the landscape the regard holds?) or stay metaphorical?
- Is this first-slice territory or later exploration?

**Your input:**


---

## 8. Tipharet as Core Athanor Principle

You identified Tipharet as a supernal concept that should inform the whole athanor system, not just this MO. The middle pillar of the tree — currently implicit, should be explicit.

**Questions:**
- Should this be proposed as an addition to the core spec (`specs/spec.md`) and/or AGENTS.md?
- If so, should it be an opus under the athanor-alignment MO, or does it need your direct shaping first?
- The exploration doc (`tipheret-principle.md`) and seed doc (`tipharet-life-principle.md`) are starting points — do either capture the territory well enough to build from?

**Your input:**

