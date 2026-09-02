# Formulae System Specification

**Status:** Design complete, implementation pending
**Created:** 2026-05-04
**Origin:** Opus under collaborative-workshop MO — artifex-directed via tempering
**Extends:** `specs/job-system/spec.md` (the tzurot concept subsumes jobs)
**Related:** `shared/jobs/coder/JOB.md` § The coding dyad (first formula extraction candidate)

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[B:artifex]` | Direct from the artifex during shaping |
| `[D:reason]` | Design decision with rationale |
| `[E:existing]` | Extends existing system element |
| `[P:pending]` | Designed but not yet implemented |

---

## Overview

Formulae are compound tzurot — athanor-specific collaboration patterns that prescribe which jobs collaborate and how. Where a job defines a single professional perspective (simple tzurah), a formula defines a collaboration topology across multiple perspectives (compound tzurah). `[B:artifex]`

**Why formulae exist.** Wave 2 observation established that azers don't spontaneously collaborate — they resist leaving their bench and try to do everything themselves. The marut-as-coordinator model (Wave 3) partially addresses this by having the marut dispatch specialists after discharge. Formulae complete the picture by giving azers explicit, readable collaboration patterns at boot — the azer knows who to bring in because the formula tells them. `[E:existing — job-system/spec.md § Wave 2 observations]`

**The tzurot concept.** Jobs and formulae are collectively "tzurot" (צורות — forms). A tzurah gives shape to work. Jobs are simple tzurot (one professional perspective). Formulae are compound tzurot (collaboration patterns across perspectives). `[B:artifex]`

---

## Domain Dictionary

Extends the job-system spec's domain dictionary.

| Term | What it is | New or extends |
|------|-----------|----------------|
| **Tzurah** (pl. tzurot) | A form that gives shape to work. The genus that contains both jobs and formulae. Pronounced "tsoo-RAH." | New `[B:artifex]` |
| **Formula** (pl. formulae) | A compound tzurah — an athanor-specific collaboration pattern that prescribes which jobs collaborate and how. Lives in `$ATHANOR/formulae/<name>/FORMULA.md`. The primary azer reads the formula and creates collaborators via `ath collaborate`. | New `[B:artifex]` |
| **Simple tzurah** | A job — one professional perspective applied to work. Global (defined in `shared/jobs/`, available to all athanors). | New (classification of existing concept) |
| **Compound tzurah** | A formula — collaboration topology across multiple perspectives. Athanor-specific (defined per instance, reflects how that athanor's work flows). | New `[B:artifex]` |
| **Dispatch hierarchy** | The marut's inscription decision order: formula → job → generalist. Start with "is there a formula for this?" `[B:artifex]` | New |

---

## Design Decisions (Artifex-Directed)

These decisions are settled. They came directly from the artifex in the MO tempering (2026-05-04). `[B:artifex]`

1. **Formulae live in `$ATHANOR/formulae/<name>/`** as markdown files with frontmatter, mirroring JOB.md structure.
2. **Formulae are azer-readable** — the azer reads the formula and creates collaborators via `ath collaborate`.
3. **The opus carries a `formula:` field** in frontmatter.
4. **Dispatch hierarchy: formula > job > generalist.** The marut's inscription decision starts with "is there a formula for this?"
5. **Formulae are athanor-specific, not global.** Each athanor defines its own collaboration patterns based on how its work flows. Jobs remain global (professional perspectives are universal).
6. **Jobs and formulae are collectively "tzurot."** `ath tzurot` lists both.

---

## Formula File Format

### Location

```
$ATHANOR/formulae/<name>/FORMULA.md
```

Unlike jobs (global, symlinked from `shared/jobs/`), formulae live directly in the athanor instance directory. A work domain's software development formulae differ from a life-domain's formulae. `[B:artifex]`

**Directory creation is on demand.** `ath init` does not create the `formulae/` directory. It's created when the first formula is added (manually or by a future `ath formula create` command). `[D:on-demand — avoids empty dirs in athanors that don't use formulae yet]`

### FORMULA.md structure

A `FORMULA.md` has four body sections, modeled on the JOB.md structure:

| Section | Purpose | Analogous JOB.md section |
|---------|---------|--------------------------|
| **Identity paragraph** | What this collaboration pattern is and why it exists | Identity paragraph |
| **When this formula applies** | Routing guidance for marut and inscribers | When this role is needed |
| **The topology** | Who participates, how they relate, the essential character of each role | What you care about |
| **Inscription guidance** | How the primary azer brings in collaborators via `ath collaborate` | Your instinct |

### Frontmatter

```yaml
---
summary: Implementation with independent code review
when:
  - "code needs to be written or modified"
  - "any implementation work that will ship to production"
---
```

| Field | Required | Type | Purpose |
|-------|----------|------|---------|
| `summary` | Yes | string | One-line description for `ath tzurot` listing |
| `when` | Yes | list | Conditions that trigger this formula — routing guidance for marut |

The frontmatter is deliberately minimal. The collaboration design lives in the body as essentialist prose, not structured data. This follows the job system's principle: perspectives and patterns over procedures. `[D:minimal-frontmatter — same reasoning as JOB.md]`

**No `model:` field.** The model comes from the primary azer's job definition, not the formula. The formula layers collaboration guidance on top of professional identity; it doesn't change what the professional identity is. `[D:model-from-job]`

**No `jobs:` list in frontmatter.** The topology section in the body describes which jobs participate. A flat list in frontmatter loses the relational structure (who inscribes whom, sequencing, rationale). If tooling needs structured job data later, it can be added without breaking the format. `[D:no-structured-jobs — YAGNI, body carries this]`

**Uppercase `FORMULA.md`.** Matches the `JOB.md` convention — uppercase signals "this is a role/form definition file, not a regular document." `[D:naming-matches-JOB.md]`

**`FormulaeDir` uses Latin plural, `FormulaFile` uses singular.** This parallels `JobsDir`/`JobFile` (plural dir, singular file). The directory uses the Latin plural "formulae" rather than English "formulas" because the system vocabulary already uses the Latin form. `[D:naming-latin-plural-consistent-with-vocabulary]`

### Design principles for formula authors

**Essentialist, not procedural.** A formula describes a collaboration pattern, not a workflow. "The reviewer checks your work at natural checkpoints" is a pattern. "Step 1: write code. Step 2: send to reviewer. Step 3: wait for approval." is a procedure. Formulae are the former. `[E:existing — same principle as JOB.md]`

**The primary azer has agency.** The formula tells the azer about the collaboration topology and who to inscribe. The azer decides the specifics — what intent to provide, when exactly to inscribe, how to adapt the pattern to the current work. The formula shapes; the azer acts. `[E:existing — geas supremacy]`

**Formulae don't replace job identity.** The azer's professional identity comes from their job. The formula layers collaboration guidance on top. A coder reading a coding-dyad formula remains a coder — they just know they should inscribe a reviewer. `[D:formula-layers-on-job]`

---

## Opus Frontmatter

The `formula:` field is optional in opus frontmatter. When present, the mustered azer reads the formula at boot alongside their job definition. `[B:artifex]`

```yaml
---
status: charged
inscribed: 2026-05-04
magnum_opus: some-mo
job: coder
formula: coding-dyad
---
# Implement the auth token refresh logic

## Intent

The auth token refresh logic handles concurrent requests gracefully...
```

The `job:` field remains required — every azer has a professional identity. The `formula:` field is optional — not all work needs a prescribed collaboration pattern. When the inscriber specifies a formula, they also specify the primary azer's job. `[D:job-still-required]`

### Dispatch hierarchy at inscription time

```
Work identified
  │
  ├─ Appropriate formula exists? ──→ Inscribe with job + formula
  │                                    (formula prescribes collaboration topology)
  │
  ├─ No formula, appropriate job? ──→ Inscribe with job only
  │                                    (azer uses judgment for collaboration)
  │
  └─ No formula, no specific job ──→ Inscribe with job: general
                                       (full latitude)
```

The marut's inscription decision starts with "is there a formula for this kind of work?" Formula selection happens first because formulae encode operational learning about what collaboration patterns produce the best outcomes for specific work types in this athanor. `[B:artifex]`

---

## Boot Prompt Wiring

### How `ath muster` wires the formula

When an opus has a `formula:` field, `ath muster` adds the formula to the azer's boot prompt after the job definition:

```
Read AGENTS.md, then read MO, then read azer.md, then read JOB.md
  [, then read JOB.local.md], then read FORMULA.md.
  Your opus is at <path>. Read it and execute.
```

The formula is read AFTER the job — it layers collaboration guidance on top of professional identity. The reading order establishes context: system vocabulary → goal → role → professional identity → collaboration pattern → specific briefing. `[D:read-order — formula after job, before opus]`

### Implementation: `formulaBootClause`

Follows the same pattern as `jobBootClause` in `muster_cmd.go`:

```go
// formulaBootClause returns ", then read <path>" for the boot prompt
// if a formula is specified, or empty string if no formula is set.
func formulaBootClause(instDir, formula string) (string, error) {
    if formula == "" {
        return "", nil
    }
    formulaPath := filepath.Join(instDir, "formulae", formula, "FORMULA.md")
    if _, err := os.Stat(formulaPath); err != nil {
        return "", fmt.Errorf(
            "formula definition not found: %s", formulaPath,
        )
    }
    return fmt.Sprintf(", then read %s", formulaPath), nil
}
```

The boot prompt assembly in `runMusterOpus`, `musterInscribedOpus`, and `runMusterIntent` appends the formula clause after the job clause. `[P:pending]`

### Intent mode (`ath muster --intent`)

Intent mode builds the boot prompt without an opus file, so there's no frontmatter to read `formula:` from. Add a `--formula` flag to `ath muster` for this case:

```bash
ath muster my-mo impl-auth --intent "Implement auth token refresh" \
  --job coder --formula coding-dyad
```

When `--formula` is provided, `runMusterIntent` appends the formula clause to the boot prompt, same as opus mode. The `--formula` flag is optional in all modes. In opus mode, `--formula` overrides the opus frontmatter (same precedence pattern as `--job` and `--model`). `[P:pending]`

### Reading formula from opus frontmatter

```go
// ReadOpusFormula reads the formula field from an opus file's YAML frontmatter.
func ReadOpusFormula(path string) string {
    return readFrontmatterField(path, "formula:")
}
```

Follows the same pattern as `ReadOpusJob`. `[P:pending]`

---

## CLI Changes

### `ath inscribe` — add `--formula` flag `[P:pending]`

Optional flag. When provided, adds `formula:` to opus frontmatter. Validated against `$ATHANOR/formulae/<name>/FORMULA.md` in the resolved instance directory.

```bash
ath inscribe my-app my-mo \
  --job coder \
  --formula coding-dyad \
  --intent "Implement the auth token refresh logic"
```

The `--formula` flag is optional. `--job` remains required. The formula is validated at inscription time — inscribing with a nonexistent formula is an error. `[D:validate-at-inscription]`

**`buildOpusContent` update.** The function signature gains a `formula` parameter. The `formula:` field is emitted in frontmatter after `job:`, only when non-empty:

```go
func buildOpusContent(datestamp, moName, job, formula, intent, collaborationContext string) string {
    // ... existing frontmatter ...
    if job != "" {
        fm.WriteString(fmt.Sprintf("job: %s\n", job))
    }
    if formula != "" {
        fm.WriteString(fmt.Sprintf("formula: %s\n", formula))
    }
    // ... rest unchanged ...
}
```

All callers (`runInscribe`, `runCollaborate`) pass the formula value through. `[P:pending]`

### `ath collaborate` — add `--formula` flag `[P:pending]`

Optional flag, same validation as `ath inscribe`. The `--formula` value is written into the **collaborator's** opus frontmatter — it tells the collaborator which collaboration pattern to follow in their own work. The formula applies to the inscribed azer, not as a back-reference to the caller's pattern.

Most collaborators inscribed via a formula won't themselves carry a formula — the primary azer's formula tells them to inscribe a reviewer, and the reviewer just does review with their job identity. But when a collaborator's own work warrants a collaboration pattern (e.g., inscribing a parallel implementation track that itself follows a coding dyad), the caller specifies it:

```bash
ath collaborate my-mo \
  --job coder \
  --formula coding-dyad \
  --intent "Implement the frontend auth flow in parallel"
```

### `ath tzurot` — new command `[P:pending]`

Lists both jobs and formulae. Jobs come from the global registry (`shared/jobs/`). Formulae come from the current athanor's instance directory (`$ATHANOR/formulae/`).

```
$ ath tzurot

Jobs (global):
  assessor — Landscape survey and opus inscription
  coder — Write, modify, and review software
  general — Novel, exploratory, cross-cutting work
  investigator — Trace problems to root cause
  project-manager — Witness satisfaction audit and progress assessment
  qa-specialist — Prevent surprises in software behavior
  solution-architect — Periodic holistic architecture and spec-reality review

Formulae (my-app):
  coding-dyad — Implementation with independent code review
```

If `$ATHANOR` is not set, only jobs are listed (formulae require an athanor context). If the athanor has no `formulae/` directory or it's empty, the formulae section is omitted. `[D:graceful-absence]`

### `ath formulae` — list/show formulae `[P:pending]`

Mirrors `ath jobs`:
- No args: list all formulae with summaries and when-to-invoke entries
- One arg: show detail for a specific formula

Operates on the athanor identified by `$ATHANOR` or `--athanor` flag. Error if neither is set.

---

## Formula Validation

Unlike jobs (validated against the global shared registry), formulae are validated against the athanor instance directory:

```go
// FormulaFile is the formula definition filename.
const FormulaFile = "FORMULA.md"

// FormulaeDir is the subdirectory containing formulae in an athanor instance.
const FormulaeDir = "formulae"

// ValidateFormula checks that a named formula exists in the athanor instance.
func ValidateFormula(instDir, formulaName string) error {
    formulaPath := filepath.Join(instDir, FormulaeDir, formulaName, FormulaFile)
    if _, err := os.Stat(formulaPath); err != nil {
        if os.IsNotExist(err) {
            available, _ := ListFormulae(instDir)
            return fmt.Errorf("unknown formula %q (available: %s)",
                formulaName, formatJobList(available))
        }
        return fmt.Errorf("checking formula %q: %w", formulaName, err)
    }
    return nil
}

// ListFormulae returns the names of all formulae in an athanor instance.
func ListFormulae(instDir string) ([]string, error) {
    formulaeDir := filepath.Join(instDir, FormulaeDir)
    entries, err := os.ReadDir(formulaeDir)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil
        }
        return nil, fmt.Errorf("listing formulae: %w", err)
    }
    var names []string
    for _, e := range entries {
        if e.IsDir() {
            fmd := filepath.Join(formulaeDir, e.Name(), FormulaFile)
            if _, err := os.Stat(fmd); err == nil {
                names = append(names, e.Name())
            }
        }
    }
    return names, nil
}
```

`[P:pending]`

---

## Example: The Coding Dyad

The coder JOB.md currently contains a "coding dyad" section (lines 34-45) that prescribes collaboration with a peer reviewer. This is a formula embedded in a job — the canonical extraction candidate. `[B:artifex]`

### The formula: `$ATHANOR/formulae/coding-dyad/FORMULA.md`

```yaml
---
summary: Implementation with independent code review
when:
  - "code needs to be written or modified"
  - "any implementation work that will ship to production"
---
# Coding Dyad

Code review through independent fresh context. The implementer and reviewer are
peers in separate crucibles, each a coder by job — the split is structural, not
hierarchical. Builder bias is physics: as you accumulate implementation context,
your ability to see your own mistakes degrades. The reviewer's value comes from
freshness, not seniority.

## When this formula applies

Any time code is being written or modified that will ship. The coding dyad is the
default formula for implementation work in this athanor. Skip it only for trivial
changes where the implementation is smaller than the overhead of inscription (a
one-line config change, a typo fix).

## The topology

Two coders — one implementing, one reviewing. Both hold the coder job identity.
The reviewer lives in their own crucible with fresh context, communicating via
whisper. The relationship is peer and egalitarian: the reviewer brings genuine
perspective, not rubber stamps.

## Inscription guidance

At the start of implementation, use `ath collaborate` to inscribe a peer coder as
your reviewer:

    ath collaborate <mo> --job coder \
      --intent "Review [description]. Check construction quality, pattern
      consistency, architecture alignment. Whisper back findings."

The reviewer checks your work at natural checkpoints:
- **Green step** — when tests pass for a chunk of work, whisper the reviewer
- **PR creation** — the reviewer examines the full change from fresh context
- **Major stage boundaries** — any significant shift in what you're building

When you're the reviewer: bring your own perspective — construction quality,
pattern consistency, architecture alignment, blast radius. Push back when
something doesn't hold up. The builder benefits from genuine dialogue, not rubber
stamps.
```

### What changes in coder JOB.md

The "coding dyad" section (lines 34-45) is replaced with a lighter formula-aware section:

```markdown
## Your collaboration

When your opus carries a `formula:` field, follow the formula's inscription
guidance. The coding-dyad formula is the standard formula for implementation
work — it pairs you with a peer reviewer in a separate crucible.

When no formula is specified, your instinct should still be to seek independent
review for any non-trivial implementation. Builder bias is physics, not a
discipline problem — a craftsman who knows this seeks fresh perspective out of
respect for the craft, not insecurity.
```

The essential insight ("builder bias is physics") stays in the job. The collaboration pattern (who to inscribe, when, how) moves to the formula. `[D:insight-stays-pattern-moves]`

---

## Guidance File Updates

### AGENTS.md vocabulary additions `[P:pending]`

Add to the vocabulary table:

| Term | What it is |
|------|-----------|
| **Tzurah** (pl. tzurot) | A form that gives shape to work. Jobs are simple tzurot (one professional perspective). Formulae are compound tzurot (collaboration patterns across perspectives). `ath tzurot` lists both. |
| **Formula** (pl. formulae) | A compound tzurah — an athanor-specific collaboration pattern that prescribes which jobs collaborate and how. Lives in `$ATHANOR/formulae/<name>/FORMULA.md`. The primary azer reads the formula and creates collaborators via `ath collaborate`. |

### azer.md updates `[P:pending]`

In § Mise en Place, update step 2 to be formula-aware:

> 2. **Inscribe collaborators.** If your opus has a `formula:` field, read `formulae/<formula>/FORMULA.md` — it prescribes the collaboration topology for your work. Follow its inscription guidance to bring in collaborators via `ath collaborate`. If no formula is specified, use your judgment: consider what specialist perspectives would strengthen the outcome. Run `ath tzurot` to see what jobs and formulae are available.

### marut.md updates `[P:pending]`

In the coordination section, update the dispatch guidance:

> When inscribing work, follow the dispatch hierarchy: (1) if an appropriate formula exists for this kind of work in this athanor, use it — specify both `--job` and `--formula`; (2) if appropriate job(s) exist but no formula fits, use the job; (3) generalist azer as fallback. Run `ath tzurot` to see available jobs and formulae.

### opus.md updates `[P:pending]`

In § Frontmatter, add `formula:` to the charged opus example:

```yaml
---
status: charged
inscribed: 2026-03-23
job: coder
formula: coding-dyad       # optional — collaboration pattern from $ATHANOR/formulae/
magnum_opus: some-mo
---
```

---

## Impact on Existing Components

| Component | Change | Scope |
|-----------|--------|-------|
| `shared/AGENTS.md` | Add tzurah/tzurot, formula/formulae vocabulary | Vocabulary table |
| `shared/azer.md` | Formula-aware mise en place | § Mise en Place step 2 |
| `shared/marut.md` | Formula-first dispatch hierarchy | § Coordination |
| `shared/opus.md` | `formula:` in frontmatter docs | § Frontmatter |
| `shared/jobs/coder/JOB.md` | Extract coding dyad → formula reference | § The coding dyad → § Your collaboration |
| `specs/job-system/spec.md` | Cross-reference to this spec, tzurot concept | Brief addition |
| `internal/athanor/home.go` | `FormulaFile`, `FormulaeDir` constants; `ValidateFormula`, `ListFormulae`, `ReadOpusFormula` | New functions |
| `internal/cli/inscribe_cmd.go` | `--formula` flag, frontmatter generation | Flag + `buildOpusContent` |
| `internal/cli/collaborate_cmd.go` | `--formula` flag | Flag |
| `internal/cli/muster_cmd.go` | `--formula` flag, `formulaBootClause`, boot prompt assembly in all three modes | New flag + function + integration |
| New: `internal/cli/tzurot_cmd.go` | `ath tzurot` command | New file |
| New: `internal/cli/formulae_cmd.go` | `ath formulae` command | New file |
| `SYSTEM-MAP.md` | Instance structure shows `formulae/` | Diagram update |

---

## Implementation Sequence

Ordered by dependency and value.

**Phase 1: Specification and vocabulary** (this opus)
1. This spec document
2. AGENTS.md vocabulary additions
3. Job-system spec cross-reference

**Phase 2: Infrastructure** (coder opus)
1. `athanor` package: constants, `ValidateFormula`, `ListFormulae`, `ReadFormulaInfo`, `ReadOpusFormula`
2. `inscribe_cmd.go`: `--formula` flag, `buildOpusContent` update
3. `collaborate_cmd.go`: `--formula` flag
4. `muster_cmd.go`: `formulaBootClause`, boot prompt integration
5. `tzurot_cmd.go` and `formulae_cmd.go`: new listing commands
6. `root.go`: register new commands
7. Tests: unit tests for formula parsing, acceptance test extensions

**Phase 3: First formula and guidance updates** (after CLI ships)
1. Create `coding-dyad` formula in my-app's `formulae/` directory
2. Update coder JOB.md: extract coding dyad → formula reference
3. Update azer.md: formula-aware mise en place
4. Update marut.md: formula-first dispatch
5. Update opus.md: `formula:` frontmatter documentation
6. Update SYSTEM-MAP.md: instance structure

**Phase 4: Observation**
1. Does the marut use formulae in dispatch decisions?
2. Do azers read and follow formulae?
3. Do the collaboration topologies in trails reflect the formula patterns?
4. Are there work patterns in this athanor that need new formulae?

---

## Relationship to Existing Mechanisms

### Formulae vs. JOB.local.md

Jobs have per-athanor layers via `JOB.local.md` — additive guidance loaded after `JOB.md`. Formulae don't need this mechanism because they are already athanor-specific. A formula IS the local layer for collaboration patterns. `[D:no-formula-local — already per-athanor]`

### Formulae vs. the coding dyad in JOB.md

The coder JOB.md's "coding dyad" section is the ur-formula — a collaboration pattern currently embedded in a job definition. Extracting it to a formula:
- **Separates concerns:** The job says WHO you are (a coder). The formula says WHO ELSE should be involved (a peer reviewer) and HOW.
- **Enables variation:** Different athanors might want different review patterns. A security-focused athanor might have a "secure-coding-triad" formula (implementer + reviewer + security specialist). The coder job stays the same; the formula changes.
- **Makes the pattern selectable:** The marut can choose to apply the coding dyad or not, depending on the work. A trivial config change might not warrant a reviewer. With the pattern in the job, it always applies.

### Formulae vs. marut coordination

The marut-as-coordinator model (Wave 3) has the marut dispatching specialists after an azer discharges. Formulae complement this:
- **Formulae act proactively.** The azer inscribes collaborators at the start, per the formula. The marut doesn't need to notice and react.
- **Marut acts reactively.** When an azer discharges without following its formula (or when no formula exists), the marut catches the gap.
- **Both are needed.** Formulae handle known patterns. The marut handles emergent situations and backstops.

### Formulae and `ath sync`

`ath sync` reconciles shared components (role files, job symlinks). Formulae are not shared — they live directly in the athanor instance. `ath sync` does not touch them. No changes to sync logic needed. `[D:no-sync-for-formulae]`

---

## Edge Cases

**Deleted formula, charged opus.** An opus references `formula: coding-dyad` but the formula directory has been deleted. At muster time, `formulaBootClause` returns an error and `ath muster` refuses to launch. This is the correct behavior — mustering with a missing formula is a broken state. The marut should either recreate the formula, remove the `formula:` field from the opus, or inscribe a new opus without the formula. Discharged opera that reference deleted formulae are harmless — the reference is historical. `[D:fail-at-muster — same as missing job]`

**Formula references a missing job.** A formula's topology section says "inscribe a coder" but the athanor lacks the `coder` job symlink. The formula loads fine (it's just markdown), but when the azer runs `ath collaborate --job coder`, the job validation fails with a clear error. This is the existing failure mode for missing jobs and needs no special handling. The fix is `ath sync` to reconcile job symlinks. `[D:existing-failure-mode — no special handling]`

**Opus has formula but no job.** Not possible — `job:` is required and enforced by `ath inscribe`. If an opus file is manually created with `formula:` but no `job:`, `ath muster` reads an empty job and proceeds without a job clause (existing behavior for missing job field). The formula still loads. This is a degenerate case, not a design concern. `[D:job-required-enforced]`

---

## Test Scenarios `[P:pending]`

The project requires both unit tests and acceptance test extensions (see project CLAUDE.md § Testing Strategy).

### Unit tests

| Test | What it verifies |
|------|-----------------|
| `TestReadOpusFormula` | Reads `formula:` from opus frontmatter; returns empty string when absent |
| `TestValidateFormula` | Returns nil for existing formula; returns error with available list for missing formula; handles missing `formulae/` directory gracefully |
| `TestListFormulae` | Lists formula directories containing FORMULA.md; skips dirs without FORMULA.md; returns nil for missing `formulae/` dir |
| `TestFormulaBootClause` | Returns clause for valid formula; returns empty string for empty formula name; returns error for missing formula file |
| `TestParseInscribeArgsFormula` | Parses `--formula` flag; omitting `--formula` is valid; `--formula` without `--job` still requires `--job` |
| `TestBuildOpusContentFormula` | Emits `formula:` in frontmatter when provided; omits when empty; placement is after `job:` |

### Acceptance test extensions (`TestATHFullLifecycle`)

| Scenario | What it exercises |
|----------|-----------------|
| `inscribe --formula` | Creates opus with `formula:` in frontmatter |
| `muster` with formula opus | Boot prompt includes formula clause |
| `ath tzurot` | Lists both jobs and formulae |
| `ath formulae` | Lists formulae; shows detail for one |
| `collaborate --formula` | Creates collaborator opus with `formula:` in frontmatter |

---

## Out of Scope

- **Global formulae.** No global formulae identified. If a formula proves universally useful, it can be copied between athanors — but the artifex explicitly said no global formulae at this time. `[B:artifex]`
- **Formula inheritance or composition.** A formula is a standalone document. No mechanism for one formula to extend or include another. If needed, it can be added later without breaking the format.
- **Formula creation wizard (`ath formula create`).** Formulae are markdown files — create them with an editor or azer. CLI automation can come later if the creation frequency warrants it.
- **Automated formula selection.** The marut selects formulae through judgment, not pattern matching. Heuristic selection (match opus intent keywords to formula `when:` conditions) is a possible future enhancement but not designed here.
