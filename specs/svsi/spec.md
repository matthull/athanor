# SVSI: Stochastic Vector-Sigil Isometry

**Status:** Draft (from interactive shaping session 2026-04-07)
**Created:** 2026-04-07
**Prototype Code:** `research/svsi-prototype/`
**Prototype Outputs:** `research/svsi-prototype/*-output.txt`
**Best current renderer:** `research/svsi-prototype/sigil_renderer.py` (stroke-based)
**Correspondence provenance:** Golden Dawn / Western Hermetic synthesis (see § Correspondence Sources)
**Key tuning parameter:** softmax temperature = 0.3 for sharp arcana selection

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[D:session]` | Decision from 2026-04-07 shaping session |
| `[D:prototype]` | Validated by prototype code |
| `[D:theory]` | Theoretical basis, not yet validated |
| `[U:topic]` | Unbacked — needs further exploration |

---

## Overview

SVSI is a system for generating deterministic visual sigils from arbitrary text, where **semantic similarity in the input produces visual similarity in the output**. The core insight: Western esoteric correspondence tables (tarot, kabbalah, astrology) function as a **structured semantic codebook** — a hand-curated, centuries-refined normalization layer that compresses high-dimensional text semantics into a finite vocabulary of archetypes with rich, pre-existing visual traditions.

The system serves two purposes:

1. **Athanor visualization** — rendering the state of Magna Opera, opera, and agent sessions as living sigils in the terminal, giving the artifex an at-a-glance feel for "are things moving along? Is confidence increasing? Are we approaching satisfaction?"

2. **A potentially novel ML primitive** — esoteric correspondence tables as structured semantic codebooks for embedding systems, offering interpretable dimensionality reduction with cross-modal grounding built in.

### The Core Property

**Identical semantics → identical sigil.** Deterministic. `[D:session]`

**Similar semantics → similar sigil, with meaningful frequency.** Stochastic, not guaranteed. The mapping preserves semantic distance often enough that an experienced viewer recognizes visual families. `[D:session]`

**Dissimilar semantics → visually distinct sigils.** Different archetypes, different shapes, different correspondences. `[D:prototype]`

---

## Domain Dictionary

| Term | Definition |
|------|-----------|
| **SVSI** | Stochastic Vector-Sigil Isometry — the overall system and its core property |
| **Sigil** | A deterministic visual glyph generated from text input. Not arbitrary — semantically grounded. Human-recognizable form, not raw data visualization |
| **Correspondence** | An atom of the esoteric vocabulary. The fundamental visual/conceptual primitives: elements (fire, water, air, earth), modalities (hot, cold, wet, dry), geometric forms, chromatic values, dynamics, densities. The term comes from Western esoteric practice where it describes exactly this: the cross-domain mapping atoms |
| **Arcanum** (pl. arcana) | A pre-composed bundle of correspondences. The 22 Major Arcana of the tarot serve as the primary semantic codebook — each is a cluster center in archetype space with established correspondence decompositions |
| **Correspondence table** | The mapping from arcana to correspondences. E.g., Temperance → {fire, water, balance, flowing, bilateral, blue, triangle, moderate}. Sourced from Golden Dawn and broader Western Hermetic tradition |
| **Major correspondence** | A correspondence that always fires for a given arcanum. Definitional. E.g., "earth" for Pentacles. Weight: 1.0 |
| **Minor correspondence** | A correspondence that fires probabilistically. Associative. E.g., "playfulness" for 2 of Pentacles. Weight determined by stochasticity dial |
| **Stochasticity dial** | A parameter (0.0–1.0) controlling how many minor correspondences are emitted. 0.0 = only major correspondences (fully deterministic). 1.0 = all correspondences (maximum variation) |
| **Arcana profile** | A 22-dimensional vector representing a text's resonance with each Major Arcanum. Computed via embedding similarity |
| **Quality profile** | The correspondence vector derived by decomposing the arcana profile through the correspondence table. This is the renderer's actual input |

---

## Architecture — The Pipeline

The system is a multi-stage pipeline. Each stage is a specialized transformation that normalizes the signal for the next stage. `[D:session]`

```
STAGE 1: SEMANTIC EMBEDDING
  raw text → dense vector (384-dim via sentence-transformers)
  
STAGE 2: ESOTERIC MAPPING ("reverse tarot reading")
  dense vector → arcana profile (22-dim)
  Method: embed arcana descriptions, compute cosine similarity
  Sharp selection via softmax with temperature control
  
STAGE 3: CORRESPONDENCE DECOMPOSITION
  arcana profile → correspondence vector
  Each arcanum decomposes into major + minor correspondences
  Stochasticity dial controls minor correspondence emission
  
STAGE 4: RENDERING
  correspondence vector → terminal sigil
  Correspondences drive visual primitives:
    - Shape/silhouette (from geometric correspondence)
    - Symmetry type (from symmetry correspondence)
    - Fill density (from density correspondence)
    - Color (from chromatic correspondence)
    - Motion/animation (from dynamic correspondence)
    - Element symbols (from elemental correspondence)
```

### Pipeline Properties

- **No LLM in the hot path.** Embedding model + cosine similarity + lookup tables + deterministic rendering. `[D:session]`
- **LLM role is bounded.** Optionally used once at setup to enrich correspondence tables, or to do initial "reverse tarot reading" for complex texts. Results are cached. `[D:session]`
- **Offline-capable.** With a local embedding model (e.g., `all-MiniLM-L6-v2` via sentence-transformers), the entire pipeline runs without network. `[D:prototype]`
- **Fast.** Embedding: ~50ms. Similarity + decomposition: <1ms. Rendering: <1ms. `[D:prototype]`

---

## Stage 1: Semantic Embedding

**Input:** Raw text (MO section, opus inscription, session transcript excerpt)
**Output:** Dense vector (384 dimensions for MiniLM-L6-v2)

### Embedding Model

The prototype uses `all-MiniLM-L6-v2` from HuggingFace sentence-transformers. `[D:prototype]`

Key findings from prototype:
- TF-IDF is **not semantic enough** — Matt as witness across MOs gets only 0.066 cosine similarity because the vocabulary differs per MO. `[D:prototype]`
- Sentence-transformers captures actual semantic similarity — same-person witnesses across MOs get 0.869+ in arcana-profile space. `[D:prototype]`

### Considerations

- **Model size vs. quality tradeoff.** MiniLM is fast but small. Larger models (e.g., `all-mpnet-base-v2`) might produce better semantic separation. `[U:embedding-model-comparison]`
- **Local vs. API.** Local model preferred for offline capability and zero cost. HuggingFace transformers library handles this. `[D:session]`

---

## Stage 2: Esoteric Mapping

**Input:** Dense text embedding
**Output:** 22-dimensional arcana profile

### The Codebook: 22 Major Arcana

The 22 Major Arcana serve as the primary semantic codebook. Each arcanum has:
- A **description** (embedded for similarity matching)
- **Major correspondences** (always fire)
- **Minor correspondences** (probabilistically fire)

The prototype validates that the arcana produce semantically sensible mappings: `[D:prototype]`

| Input Text | Top Arcana | Assessment |
|------------|-----------|------------|
| Artifex witness (Alignment) | Justice, High Priestess, Judgement | Weighing evidence, hidden knowledge, assessment — apt for a verification-oriented witness |
| Artifex witness (Harvest) | Magician, Strength, Sun | Skill manifest, enduring power, success — apt for career/value concerns |
| Horse Owners | Hanged Man, Sun, Judgement | Seeing differently, joy, unambiguous reaction — apt for emotional demo |
| Calcinatio: Audit | Justice, Moon, Judgement | Assessment, pattern-reading, reckoning — apt for behavioral audit |
| Tempering: Stabilization | Temperance, Tower, World | Moderation, structural risk, completion — apt for pre-release stabilization |

### Selection Sharpness

Raw cosine similarity produces diffuse weights across all 22 arcana. Softmax with temperature control sharpens the selection: `[D:prototype]`

- **Temperature 0.3** produces sharp selection where top-3 arcana dominate
- **Temperature 1.0** produces diffuse selection (everything contributes)
- Temperature is a tuning parameter, not a user-facing control

### Minor Arcana Role

The 56 Minor Arcana (4 suits × 14 cards) are **not** part of the primary codebook. `[D:session]`

**Rationale:** Minor arcana are compositional (suit element × sephirotic number), making them similar to each other. Averaging over 78 cards kills discriminative power — prototype showed all quality profiles converging to 0.978+ similarity. `[D:prototype]`

**Future role:** Minor arcana could be used for stochastic enrichment — after primary arcana are selected, a randomly chosen associated minor arcanum adds variation to the correspondence output. This matches the major/minor correspondence weighting model. `[D:session]`

### Similarity Preservation

Validated by prototype — arcana-profile similarity for key pairs: `[D:prototype]`

| Pair | Similarity | Expected |
|------|-----------|----------|
| Artifex/Alignment vs Artifex/Harvest | 0.869–0.926 | High (same person) |
| Horse Owners witness vs Horse Intent | 0.980–0.989 | Very high (same content) |
| Calcinatio:Audit vs Calcinatio:Customer | 0.847–0.971 | High (same type) |
| Artifex vs Kelly | 0.603–0.959 | Lower (different person) |

---

## Stage 3: Correspondence Decomposition

**Input:** Arcana profile (22-dim)
**Output:** Correspondence vector (set of weighted correspondences)

### The Correspondence Vocabulary

Correspondences are the renderer's input language. They are the atoms of the visual system. `[D:session]`

**Elemental:** fire, water, air, earth `[D:session]`
**Modalities (Aristotelian):** hot, cold, wet, dry `[D:session]`
**Geometric forms:** point, line, triangle, square, pentagon, hexagram, circle, spiral, cross, crescent `[D:session]`
**Chromatic:** white, black, red, blue, yellow, green, orange, violet, indigo, gold (simplified from Golden Dawn King/Queen color scales) `[D:session]`
**Dynamics:** expansion, contraction, ascending, descending, radiating, concentrating, flowing, stillness, rotation, oscillation, explosive, dissolving `[D:session]`
**Density:** sparse, moderate, dense, saturated `[D:session]`
**Symmetry:** asymmetric, bilateral, radial, 3-fold, 4-fold `[D:session]`
**Abstract:** balance, force, receptivity, transformation, structure, liberation, unity, duality, mystery, clarity, endurance, swiftness `[D:session]`

Total: ~60 correspondences in prototype. The vocabulary is extensible. `[D:prototype]`

### Major vs. Minor Correspondences

Each arcanum maps to correspondences at two weight levels: `[D:session]`

- **Major correspondences:** Always emitted. Definitional for the arcanum. E.g., Temperance → {fire, water, balance, flowing, bilateral, blue, triangle, moderate}
- **Minor correspondences:** Probabilistically emitted based on the stochasticity dial. Associative. E.g., Temperance → {patience, dualism, hourglass, moderation}

This two-level system provides simple, tunable stochasticity: `[D:session]`
- Moving more correspondences into "major" → reduces variation → more deterministic
- Moving more into "minor" → increases variation → more generative
- Per-arcanum weighting is also possible: "earth" and "duality" at 100% for 2 of Pentacles, "dialectic" and "playfulness" at 40%

### Correspondence Sources

The correspondence tables draw from Golden Dawn tradition and broader Western Hermeticism. Sources available in training data include Liber 777 (Crowley's master correspondence table), the Golden Dawn color scales, and standard tarot-astrology-kabbalah mappings. `[D:session]`

Tables can also be extracted from texts or refined through use. `[D:session]`

**Provenance note:** The prototype's correspondence tables are a synthesis from Golden Dawn tradition (primary), with standard tarot-astrology-kabbalah attributions. Not strictly Thoth, not strictly Rider-Waite — a pragmatic synthesis for the purpose of semantic normalization. The specific table used matters for reproducibility: two different correspondence systems will produce incompatible sigils. The prototype's tables in `esoteric_codebook.py` and `sharp_renderer.py` are the current canonical source. `[D:session]`

**Why Western esotericism specifically?** Several properties make Hermetic correspondence tables unusually well-suited as a semantic codebook: `[D:session]`
- **Stable over centuries** — the core attributions haven't changed significantly since the Golden Dawn systematization
- **Multi-dimensional** — each arcanum maps simultaneously to planet, element, number, color, path, hebrew letter. One mapping gives you grounding across multiple modalities.
- **Abstract/universal** — the arcana describe patterns (transformation, balance, directed force) not specific objects. They apply equally to psychology, chemistry, warfare, business.
- **Large community of practice** — centuries of use means the correspondences have been stress-tested across domains by many practitioners
- **Finite and structured** — 22 arcana is small enough to be a practical codebook while covering the full range of archetypal patterns

---

## Stage 4: Rendering

**Input:** Correspondence vector
**Output:** Terminal sigil (Unicode art with optional ANSI color)

### Rendering Approach: Stroke-Based Sigils

Real sigils are **paths, arms, and nodes** — not filled shapes. The final prototype (`sigil_renderer.py`) validates this approach: correspondences drive a stroke-based renderer that produces distinctive silhouettes recognizable at a glance. `[D:prototype]`

The renderer constructs sigils from:
- **Arms:** Lines radiating from center, drawn with box-drawing characters (│─╱╲). Number and angles determined by geometric correspondence.
- **Nodes:** Symbols placed at center, junctions, and arm endpoints. Center node from abstract quality (⊕ balance, ✦ force, ◎ receptivity, ✶ transformation). Tip nodes from element (△ fire, ▽ water, ◇ air, □ earth).
- **Secondary structure:** Cross-bars, junction nodes, or connecting arcs from the secondary arcanum's correspondences.
- **Tertiary accents:** Small marks near arm tips from the tertiary arcanum.
- **Dynamic modifiers:** Motion marks (· ascending dots, ≈ flowing waves, ⟳ rotation marks, ∗ explosive scatter) from the dynamic correspondence.

This produces sigils that are **immediately distinguishable by silhouette:**

```
CRESCENT (Priestess)    STAFF (Justice)    INVERTED △ (Hanged Man)    BRANCHING △ (Temperance)

     △                     ◇                  △     △                     □
   △ ▽ △                   ◇                  △ ··· ▽                     △
   ▽ │ ▽                   │                  ▽╲· ╱╱▽                     │
    ◇│◇                    │                    ╲◎╱                       ◇
     ⊛                     │                     │                      ≈≈│≈≈
                           ◗                     │                        ⊕
                           │                     ▽                     □╱╱ ◇╲□
                           │                     △                     △  ≈  △
                           │
                           ◇
                           ◇
```

### SVSI Property in Sigil Output

The stroke-based renderer preserves the core SVSI property — validated examples: `[D:prototype]`

- **Same primary arcanum → same structural skeleton.** Both Calcinatio:Audit and Calcinatio:Customer map to Justice → both render as vertical staffs. The secondary arcana change the junction node and tip symbols, producing family resemblance with different accent.
- **Same arcana across inputs → identical sigil.** Horse Owners witness and Horse Intent both map to Hanged Man/Sun/Judgement → identical output.
- **Different primary arcanum → different silhouette.** A crescent (Priestess) looks nothing like a staff (Justice) or an inverted triangle (Hanged Man). The eye catches this instantly.
- **Shared secondary/tertiary → shared accent.** Two sigils with different primaries but shared secondary arcana have recognizable accent similarities even with different overall forms.

### Rendering Dimensions

Correspondences map to visual primitives across multiple dimensions: `[D:session]` `[D:prototype]`

#### Shape / Silhouette (highest visual priority)

The overall silhouette is the first thing the eye sees. Geometric correspondence determines the arm structure: `[D:session]` `[D:prototype]`

- `triangle` → 3 arms (up for fire/air, down for water/earth)
- `square` → 4 arms at 90 degrees
- `hexagram` → 6 arms at 60 degrees
- `cross` → 4 cardinal arms
- `circle` → 8 short arms (ring of nodes)
- `crescent` → 3 arms curving to one side
- `line` → 2 arms (up and down only) — produces a staff/wand
- `pentagon` → 5 arms at 72 degrees
- `spiral` → 3 arms with curve
- `point` → no arms, center node only

#### Symmetry

Applied as post-process to the stroke pattern: `[D:prototype]`

- `asymmetric` → no mirroring, organic
- `bilateral` → left-right mirror (strokes mirrored, ╱↔╲ swapped)
- `radial` → all-axis mirror, mandala-like
- `fourfold` → 4-way symmetry, crystalline

#### Node Composition (replaces "Fill / Interior")

Three layers of nodes, each from a different arcanum: `[D:prototype]`
- **Center node:** Primary arcanum's abstract quality → center symbol (⊕✦◎✶⊞✴◉⊗⊛☉◈⚡)
- **Arm tips:** Primary arcanum's element → tip symbols (△▽◇□)
- **Junctions:** Secondary arcanum → mid-arm markers (○●◦◎⊙ or element symbols)
- **Accents:** Tertiary arcanum → small marks adjacent to tips

#### Color `[D:session]`

Chromatic correspondences map to ANSI terminal colors. Multiple color layers:
- Primary arcanum → dominant color
- Secondary arcanum → accent color
- Tertiary arcanum → border/edge color

Color is a significant dimension but was deprioritized in the prototype to focus on shape differentiation first.

#### Animation / Time `[D:session]`

Dynamic correspondences map to animation behaviors in a live terminal display:
- `stillness` → no animation
- `flowing` → characters shift in a wave pattern
- `radiating` → pulses outward from center
- `rotating` → pattern rotates
- `oscillating` → pattern alternates between two states
- `ascending` → upward drift
- `explosive` → periodic burst pattern

Animation is a planned dimension, not yet prototyped. It adds a temporal axis to the sigil that is particularly relevant for live athanor status display.

#### Density

How filled the sigil is:
- `sparse` → mostly empty, occasional dots
- `moderate` → balanced positive and negative space
- `dense` → mostly filled
- `saturated` → nearly solid

---

## Application: Athanor Visualization

The primary application is rendering athanor state as living sigils in the terminal. `[D:session]`

### Two-Layer Architecture

1. **Outer layer (MO-level, slow):** Updated periodically (on assessment cycles). Represents the MO's overall state — its accumulated trajectory, not a snapshot. The outer layer literally **accumulates** — glyphs deposited over time that don't clear, so visual history is baked into the current render. A fresh MO looks sparse. A mature one looks rich. A stalled one looks frozen. `[D:session]`

2. **Inner layer (activity-level, fast):** Updated in real-time from live agent activity. Shows current work — active opera, whisper messages, git commits, test results. The inner layer is **perturbation** on the deeper pattern. `[D:session]`

### MO-Type Differentiation

Different MO types need fundamentally different visual behaviors, not mode switches: `[D:session]`

- **Goal-based MOs** → sigil converges toward completion. Rings fill, geometry resolves. The sigil tightens as confidence increases.
- **State-based MOs** → sigil breathes. It's always in motion, never resolving. Health = consistent rhythm. Stall = frozen pattern.
- **Open-ended MOs** → sigil accretes organically. Each work cycle deposits new visual material. The shape is unpredictable, a record of accumulated exploration.

### Data Sources for Live Visualization

All available without LLM at render time: `[D:session]`

| Data Source | How Accessed | What It Shows |
|------------|-------------|---------------|
| Opera YAML frontmatter | `rg` / file read | Opera lifecycle state, timestamps |
| Discharge records | File read | Outcome, reflection, evidence quality |
| Tmux window state | `tmux list-windows` | Active crucibles, agent types |
| Git history | `git log` | Commit cadence, diff sizes |
| Whisper messages | Whisper log files | Inter-agent communication |
| Session logs | Session log files | Timestamped activity entries |
| Tool-call patterns | Claude session JSON | Chesed/Geburah energy balance |

### Chesed/Geburah Energy Decomposition

A live session can be classified on the expansion/refinement axis using tool-call patterns — no LLM needed: `[D:session]`

| Tool Pattern | Energy |
|---|---|
| Agent spawns, manifold generation | Chesed (expansion) |
| `go test`, `make check` | Geburah (verification) |
| Read, Grep, Glob | Investigation (Chesed-adjacent) |
| Write, Edit | Creation (Netzach/Hod) |
| Dialectical review agents | Geburah (refinement) |

### CLI Integration

The visualization would live in the `ath` CLI ecosystem. Possible commands: `[D:session]`

- `ath gaze [mo-name]` — render the MO's sigil
- `ath gaze --live` — continuous update mode
- `ath status --viz` — add sigil to existing status output

---

## Conceptual Foundations

### Esotericism as Semantic Normalization

The key insight from this exploration: **esoteric correspondence systems are, computationally, hand-curated semantic codebooks optimized for cross-domain analogical reasoning.** `[D:session]`

The 22 Major Arcana aren't "topics about specific things" — they're abstract archetypal patterns that apply equally to a business MO, a therapy session, a chemical process, or a military campaign. Each arcanum is a cluster center in abstract semantic space with explicit cross-domain mappings (psychological ↔ elemental ↔ astrological ↔ numerical ↔ chromatic ↔ geometric).

This is what makes them useful as a normalization layer: they compress arbitrary text into a structured, finite vocabulary that has visual tradition built in.

### Potential ML Primitive

The esoteric codebook concept may have applications beyond athanor visualization: `[D:session]`

- **Interpretable dimensionality reduction** — 22 labeled dimensions vs. 384 opaque ones. Each dimension has human-readable meaning.
- **Cross-modal grounding for free** — each arcanum maps to colors, numbers, geometries, elements. Text→visual without multi-modal training.
- **Robustness** — a fixed codebook can't overfit to training data. If archetypes capture universal patterns, representations should generalize better.

**Adjacent existing work:**
- Concept Bottleneck Models (task-specific, not universal)
- VQ-VAE codebooks (learned, not pre-defined)
- Zero-shot classification via text prototypes (CLIP-style)
- Topic models (learned topics, not pre-defined archetypes)

**What's potentially novel:** Using a **universal, domain-invariant, cross-modal** codebook refined over centuries as the bottleneck. `[U:novelty-validation]`

**Research question:** Compare archetypal projection against learned PCA/LDA on interpretability and cross-domain generalization benchmarks. `[U:benchmark-design]`

---

## Prototype Results Summary

### What Worked

1. **Sentence-transformers embeddings** produce meaningful semantic similarity. Same-person witnesses across MOs: 0.869+. `[D:prototype]`
2. **22 Major Arcana as codebook** normalize text into interpretable archetypal dimensions. The mappings are semantically sensible. `[D:prototype]`
3. **Correspondence decomposition** drives visual differentiation — different primary arcana produce clearly different silhouettes (crescent vs. triangle vs. line vs. radial). `[D:prototype]`
4. **Shape mask rendering** gives the sigils distinctive overall forms that are immediately distinguishable. `[D:prototype]`
5. **The pipeline is fast** — embedding + matching + rendering runs in <100ms. `[D:prototype]`

### What Didn't Work (and iterations that fixed them)

1. **TF-IDF vectorization** — too lexical, can't see semantic similarity across different vocabulary. Matt as witness: 0.066 cosine similarity. **Fix:** sentence-transformers. `[D:prototype]`
2. **78 arcana (major + minor)** — averaging over similar minor arcana kills discrimination. All profiles converge to 0.978+ similarity. **Fix:** 22 major arcana only, minor arcana reserved for stochastic enrichment. `[D:prototype]`
3. **Square grid rendering** — forces all sigils into the same bounding box. Human eyes see silhouette first; identical silhouettes mask interior differences. **Fix:** shape-mask renderer (fluid_renderer.py). `[D:prototype]`
4. **Filled-shape rendering** — even with shape masks, filled shapes (ovals, rectangles, triangles) don't look like sigils. **Fix:** stroke-based renderer (sigil_renderer.py) — arms, nodes, and paths produce actual sigil silhouettes. `[D:prototype]`
5. **Top-3 discrete rendering** — winner-take-all arcana selection loses the continuous quality profile. Two texts with 0.869 profile similarity but different top-3 look completely different. **Partially addressed** by stroke renderer where shared primary = shared skeleton. Full continuous rendering remains `[U:continuous-rendering]`. `[D:prototype]`

### Prototype Code

All in `research/svsi-prototype/`:

| File | Purpose |
|------|---------|
| `test_inputs.py` | 38 extracted MO sections (witnesses, calcinatio, intents, tempering) from 7 real MOs |
| `svsi_prototype.py` | 7 visualization approaches: Braille, block density, radial, SimHash, bilateral symmetric, mandala |
| `output.txt` | Full output from the 7-approach prototype |
| `esoteric_codebook.py` | 22 Major Arcana + 10 sephirot + elements with correspondence tables |
| `e2e_demo.py` | End-to-end pipeline: embed → match 22 arcana → render composed sigil |
| `e2e-output.txt` | Output showing arcana matches and similarity scores |
| `quality_renderer.py` | 78 arcana → correspondence decomposition → quality-driven rendering |
| `quality-output.txt` | Output showing over-averaging problem with 78 arcana |
| `sharp_renderer.py` | 22 arcana with softmax temperature → correspondence-driven 9x9 sigils |
| `sharp-output.txt` | Output showing improved differentiation with sharp selection |
| `fluid_renderer.py` | Shape-mask renderer where correspondences determine silhouette |
| `fluid-output.txt` | Output showing distinct silhouettes per primary arcanum |
| `sigil_renderer.py` | **Best result.** Stroke-based renderer — arms, nodes, paths. Produces actual sigil forms |
| `sigil-output.txt` | Output showing distinctive sigil silhouettes with SVSI property validated |

---

## Open Questions

1. **Embedding model selection.** Is MiniLM-L6-v2 sufficient, or do larger models produce meaningfully better semantic separation? `[U:embedding-model-comparison]`

2. **Correspondence table completeness.** The prototype has basic correspondences. Full Golden Dawn tables (Liber 777) would provide much richer decompositions. `[U:correspondence-enrichment]`

3. **Animation design.** Dynamic correspondences map naturally to animation, but the specific terminal animation techniques need prototyping. `[U:animation-prototype]`

4. **MO-level accumulation.** How to implement the "deposited history" outer layer — what data structure captures the visual trajectory over time? `[U:accumulation-model]`

5. **Stochasticity tuning.** What's the right default balance of major/minor correspondences? Needs user testing. `[U:stochasticity-defaults]`

6. **ML primitive viability.** Is the archetypal codebook genuinely useful as a general-purpose embedding primitive? Needs benchmarking against learned alternatives. `[U:novelty-validation]`

7. **Cross-MO witness recognition.** The same person appearing as a witness across MOs should produce recognizably similar sigils. Prototype shows 0.869 in arcana space but the rendering doesn't yet preserve this visually at that level. `[U:visual-similarity-fidelity]`

---

## Out of Scope

- Full `ath` CLI integration (separate implementation opus)
- Production rendering engine (prototype is sufficient for concept validation)
- Multi-athanor dashboard composition
- Non-terminal rendering targets (web, PDF)

---

## Retrospective

When this concept moves to implementation, review:

### Documentation Updates
- [ ] Correspondence table format and sourcing documented
- [ ] Pipeline architecture as a reference for other semantic→visual systems

### Workflow Improvements
- [ ] Could SVSI become an `ath` skill or MCP tool?
- [ ] Could the esoteric codebook be a reusable library?

### Knowledge Capture
- [ ] Benchmark archetypal projection vs. PCA/LDA if pursued as ML primitive
- [ ] Document which correspondence sources (Liber 777, GD color scales) are most useful
