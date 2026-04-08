#!/usr/bin/env python3
"""
SVSI Prototype — Stochastic Vector-Sigil Isometry

Explores multiple techniques for mapping semantic text vectors
to deterministic visual sigils, where semantic similarity produces
visual similarity.

Techniques explored:
1. TF-IDF → PCA → Braille dot matrix
2. TF-IDF → Random Projection → Braille dot matrix
3. TF-IDF → PCA → Block element density field
4. TF-IDF → PCA → Radial geometric sigil
5. TF-IDF → SimHash bits → Symbol composition
"""

import sys
import os
import hashlib
import math
from collections import defaultdict

import numpy as np
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.decomposition import PCA
from sklearn.random_projection import GaussianRandomProjection
from sklearn.metrics.pairwise import cosine_similarity

# Import test inputs
sys.path.insert(0, os.path.dirname(__file__))
from test_inputs import WITNESSES, CALCINATIO, INTENTS, TEMPERING, get_all_inputs


# ============================================================
# BRAILLE RENDERING UTILITIES
# ============================================================

# Braille character: 2x4 dot grid, Unicode U+2800 base
# Dot positions map to bits:
#   0  3
#   1  4
#   2  5
#   6  7
BRAILLE_BASE = 0x2800
BRAILLE_DOTS = [0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80]

def values_to_braille_char(values, threshold=0.0):
    """Convert 8 float values to a single Braille character.
    Values above threshold get dots."""
    code = BRAILLE_BASE
    for i, v in enumerate(values[:8]):
        if v > threshold:
            code |= BRAILLE_DOTS[i]
    return chr(code)

def vector_to_braille_grid(vec, rows=4, cols=6, threshold=0.0):
    """Convert a vector to a grid of Braille characters.
    Each Braille char encodes 8 values, so we need rows*cols*8 dimensions."""
    needed = rows * cols * 8
    # Pad or truncate
    if len(vec) < needed:
        vec = np.concatenate([vec, np.zeros(needed - len(vec))])
    else:
        vec = vec[:needed]

    lines = []
    idx = 0
    for r in range(rows):
        line = ""
        for c in range(cols):
            chunk = vec[idx:idx+8]
            line += values_to_braille_char(chunk, threshold)
            idx += 8
        lines.append(line)
    return lines


# ============================================================
# BLOCK ELEMENT UTILITIES
# ============================================================

# Block density characters from empty to full
BLOCK_CHARS = " ░▒▓█"
VBLOCK_CHARS = " ▁▂▃▄▅▆▇█"

def value_to_block(v, min_v=0.0, max_v=1.0):
    """Map a float value to a block density character."""
    normalized = (v - min_v) / (max_v - min_v + 1e-10)
    normalized = max(0.0, min(1.0, normalized))
    idx = int(normalized * (len(BLOCK_CHARS) - 1))
    return BLOCK_CHARS[idx]

def vector_to_density_field(vec, rows=6, cols=12):
    """Convert a vector to a density field using block elements."""
    needed = rows * cols
    if len(vec) < needed:
        vec = np.concatenate([vec, np.zeros(needed - len(vec))])
    else:
        vec = vec[:needed]

    # Normalize to 0-1
    mn, mx = vec.min(), vec.max()
    if mx > mn:
        vec = (vec - mn) / (mx - mn)

    lines = []
    idx = 0
    for r in range(rows):
        line = ""
        for c in range(cols):
            line += value_to_block(vec[idx])
            idx += 1
        lines.append(line)
    return lines


# ============================================================
# RADIAL SIGIL UTILITIES
# ============================================================

def vector_to_radial_sigil(vec, size=9):
    """Convert a vector to a radial/geometric sigil.
    Uses the vector to determine angles, radii, and which cells
    to fill in a grid centered on the middle."""
    needed = 32
    if len(vec) < needed:
        vec = np.concatenate([vec, np.zeros(needed - len(vec))])
    vec = vec[:needed]

    # Normalize
    mn, mx = vec.min(), vec.max()
    if mx > mn:
        vec = (vec - mn) / (mx - mn)

    center = size // 2
    grid = [[' ' for _ in range(size)] for _ in range(size)]

    # Place center marker
    grid[center][center] = '◉'

    # Use pairs of values as (angle, radius) to place marks
    symbols = '○◈◇◆●◎⊕⊗⊙⊛✦✧⬡⬢'
    for i in range(0, min(len(vec)-1, 28), 2):
        angle = vec[i] * 2 * math.pi
        radius = vec[i+1] * (center - 0.5) + 0.5
        sym_idx = int(vec[i] * len(symbols)) % len(symbols)

        x = int(center + radius * math.cos(angle))
        y = int(center + radius * math.sin(angle))

        if 0 <= x < size and 0 <= y < size and grid[y][x] == ' ':
            grid[y][x] = symbols[sym_idx]

    # Add symmetry: mirror some elements
    for i in range(16, min(len(vec)-1, 24), 2):
        angle = vec[i] * 2 * math.pi
        radius = vec[i+1] * (center - 0.5) + 0.5

        for mirror_angle in [angle, -angle, math.pi - angle, math.pi + angle]:
            x = int(center + radius * math.cos(mirror_angle))
            y = int(center + radius * math.sin(mirror_angle))
            if 0 <= x < size and 0 <= y < size and grid[y][x] == ' ':
                grid[y][x] = '·'

    # Connect some points with box-drawing
    return [''.join(row) for row in grid]


# ============================================================
# SIMHASH UTILITIES
# ============================================================

def simhash(text, bits=64):
    """Compute SimHash — locality-sensitive hash for text.
    Similar texts produce hashes that differ in few bits."""
    v = [0] * bits
    tokens = text.lower().split()

    for i in range(len(tokens)):
        # Use word and bigrams
        for ngram in [tokens[i]] + ([tokens[i] + " " + tokens[i+1]] if i < len(tokens)-1 else []):
            h = int(hashlib.md5(ngram.encode()).hexdigest(), 16)
            for j in range(bits):
                if h & (1 << j):
                    v[j] += 1
                else:
                    v[j] -= 1

    fingerprint = 0
    for j in range(bits):
        if v[j] > 0:
            fingerprint |= (1 << j)
    return fingerprint

def simhash_distance(h1, h2, bits=64):
    """Hamming distance between two SimHashes."""
    return bin(h1 ^ h2).count('1')

# Curated symbol vocabulary organized by visual weight/category
SYMBOL_VOCAB = {
    'light': '·∙∘○◦◌⊙⊚',
    'medium': '◈◇◆●◎⊕⊗⊛',
    'heavy': '★✦✧⬡⬢▣▤▥',
    'line': '╱╲╳│─┼╋╸',
    'arc': '╭╮╯╰◜◝◞◟',
    'dot': '⠁⠂⠄⡀⠈⠐⠠⢀',
    'planetary': '☉☽♂♀♃♄☿♆',
    'alchemical': '△▽◇○□⊕⊗⊙',
}

def simhash_to_symbols(fingerprint, bits=64):
    """Map a SimHash fingerprint to a composed symbol grid."""
    # Split hash into 8-bit chunks
    chunks = []
    for i in range(0, bits, 8):
        chunk = (fingerprint >> i) & 0xFF
        chunks.append(chunk)

    # Each chunk selects from a different symbol category
    categories = list(SYMBOL_VOCAB.values())
    grid_size = 4
    grid = []

    for r in range(grid_size):
        row = ""
        for c in range(grid_size):
            idx = r * grid_size + c
            if idx < len(chunks):
                chunk = chunks[idx % len(chunks)]
                cat = categories[idx % len(categories)]
                sym_idx = chunk % len(cat)
                row += cat[sym_idx]
            else:
                row += ' '
        grid.append(row)

    return grid


# ============================================================
# COMPOUND BRAILLE SIGIL — symmetry + structure
# ============================================================

def vector_to_compound_sigil(vec, size=7):
    """Create a more sigil-like form by applying symmetry operations
    to a projected vector. This is the 'corralling' step — taking
    raw projection and giving it recognizable structure."""
    needed = size * (size // 2 + 1)  # Only need half, mirror the rest
    if len(vec) < needed:
        vec = np.concatenate([vec, np.zeros(needed - len(vec))])
    vec = vec[:needed]

    # Normalize to 0-1
    mn, mx = vec.min(), vec.max()
    if mx > mn:
        vec = (vec - mn) / (mx - mn)

    # Build half-grid, then mirror for bilateral symmetry
    half_w = size // 2 + 1
    grid = [[' ' for _ in range(size)] for _ in range(size)]

    # Symbol selection based on value ranges
    sym_ranges = [
        (0.0, 0.15, ' '),
        (0.15, 0.3, '·'),
        (0.3, 0.45, '∘'),
        (0.45, 0.6, '○'),
        (0.6, 0.75, '◈'),
        (0.75, 0.85, '◆'),
        (0.85, 1.01, '●'),
    ]

    idx = 0
    for r in range(size):
        for c in range(half_w):
            v = vec[idx] if idx < len(vec) else 0
            idx += 1

            sym = ' '
            for lo, hi, s in sym_ranges:
                if lo <= v < hi:
                    sym = s
                    break

            # Place on left side
            grid[r][c] = sym
            # Mirror to right side
            mirror_c = size - 1 - c
            if mirror_c != c:
                grid[r][mirror_c] = sym

    return [''.join(row) for row in grid]


# ============================================================
# APPROACH 6: Braille with 4-fold symmetry (mandala-like)
# ============================================================

def vector_to_mandala(vec, quadrant_rows=3, quadrant_cols=4):
    """Create a mandala-like Braille pattern with 4-fold symmetry.
    Only one quadrant is computed from the vector; the rest are mirrored."""
    needed = quadrant_rows * quadrant_cols * 8
    if len(vec) < needed:
        vec = np.concatenate([vec, np.zeros(needed - len(vec))])
    vec = vec[:needed]

    # Build one quadrant of Braille chars
    quadrant = []
    idx = 0
    for r in range(quadrant_rows):
        row = []
        for c in range(quadrant_cols):
            chunk = vec[idx:idx+8]
            char = values_to_braille_char(chunk, threshold=0.0)
            row.append(char)
            idx += 8
        quadrant.append(row)

    # Compose full mandala: mirror horizontally and vertically
    full_rows = quadrant_rows * 2
    full_cols = quadrant_cols * 2
    lines = []

    for r in range(full_rows):
        line = ""
        qr = r if r < quadrant_rows else (full_rows - 1 - r)
        for c in range(full_cols):
            qc = c if c < quadrant_cols else (full_cols - 1 - c)
            line += quadrant[qr][qc]
        lines.append(line)

    return lines


# ============================================================
# MAIN: Run all approaches on all inputs
# ============================================================

def compute_similarity_matrix(vectors, labels):
    """Compute and display cosine similarity for key pairs."""
    sim = cosine_similarity(vectors)
    return sim

def print_sigil_comparison(title, sigils_by_key, keys, labels):
    """Print sigils side by side for comparison."""
    print(f"\n{'='*72}")
    print(f"  {title}")
    print(f"{'='*72}")

    # Print in rows of 3
    for batch_start in range(0, len(keys), 3):
        batch_keys = keys[batch_start:batch_start+3]
        batch_labels = [labels[k] for k in batch_keys]

        # Print labels
        col_width = 22
        header = ""
        for lbl in batch_labels:
            header += f"  {lbl[:col_width-2]:<{col_width-2}}"
        print(f"\n{header}")
        print("  " + "─" * (col_width * len(batch_labels)))

        # Get max rows
        max_rows = max(len(sigils_by_key[k]) for k in batch_keys)

        for r in range(max_rows):
            line = ""
            for k in batch_keys:
                sigil = sigils_by_key[k]
                if r < len(sigil):
                    line += f"  {sigil[r]:<{col_width-2}}"
                else:
                    line += f"  {'':<{col_width-2}}"
            print(line)
    print()


def main():
    all_inputs = get_all_inputs()
    keys = [k for k, _ in all_inputs]
    texts = [d['text'] for _, d in all_inputs]
    labels = {k: d['label'] for k, d in all_inputs}
    types = {k: d['type'] for k, d in all_inputs}

    print("SVSI Prototype — Stochastic Vector-Sigil Isometry")
    print(f"Total inputs: {len(all_inputs)}")
    print(f"  Witnesses: {sum(1 for _,d in all_inputs if d['type']=='witness')}")
    print(f"  Calcinatio: {sum(1 for _,d in all_inputs if d['type']=='calcinatio')}")
    print(f"  Intents: {sum(1 for _,d in all_inputs if d['type']=='intent')}")
    print(f"  Tempering: {sum(1 for _,d in all_inputs if d['type']=='tempering')}")

    # ---- VECTORIZE ----
    print("\n--- Vectorizing with TF-IDF ---")
    vectorizer = TfidfVectorizer(
        max_features=500,
        stop_words='english',
        ngram_range=(1, 2),
        sublinear_tf=True,
    )
    tfidf_matrix = vectorizer.fit_transform(texts).toarray()
    print(f"TF-IDF matrix shape: {tfidf_matrix.shape}")

    # ---- SIMILARITY ANALYSIS ----
    print("\n--- Cosine Similarity: Key Comparisons ---")
    sim = cosine_similarity(tfidf_matrix)

    # Find indices for key comparisons
    key_to_idx = {k: i for i, k in enumerate(keys)}

    comparisons = [
        ("Matt/Alignment vs Matt/Harvest", "artifex_alignment", "artifex_harvest"),
        ("Matt/Alignment vs Matt/Anvil", "artifex_alignment", "artifex_anvil"),
        ("Matt/Harvest vs Matt/Anvil", "artifex_harvest", "artifex_anvil"),
        ("Kelly/Seismic vs Kelly/Analytics", "kelly_seismic", "kelly_analytics"),
        ("Kelly/Seismic vs Kelly/SF", "kelly_seismic", "kelly_sf"),
        ("Tom/Seismic vs Tom/SF", "tom_seismic", "tom_sf"),
        ("Leadership/Seismic vs Leadership/SF", "leadership_seismic", "leadership_sf"),
        ("Leadership/Seismic vs Leadership/Harvest", "leadership_seismic", "leadership_harvest"),
        ("Matt/Alignment vs Horse Owners", "artifex_alignment", "horse_owners"),
        ("Matt/Alignment vs Behavioral Audit", "artifex_alignment", "behavioral_audit"),
        ("Behavioral Audit vs Trail Legibility", "behavioral_audit", "trail_legibility"),
        ("Behavioral Audit vs Customer Outcome", "behavioral_audit", "customer_outcome"),
        ("Alignment Intent vs Anvil Intent", "alignment_intent", "anvil_intent"),
        ("Alignment Intent vs Seismic Intent", "alignment_intent", "seismic_intent"),
        ("Horse Intent vs Horse Owners (witness)", "horse_intent", "horse_owners"),
    ]

    print(f"\n{'Comparison':<48} {'Similarity':>10}")
    print("─" * 60)
    for label, k1, k2 in comparisons:
        i1, i2 = key_to_idx[k1], key_to_idx[k2]
        s = sim[i1, i2]
        bar = '█' * int(s * 20) + '░' * (20 - int(s * 20))
        print(f"  {label:<46} {s:>6.3f}  {bar}")

    # ---- PROJECTIONS ----

    # PCA to various dimensions — capped at n_samples - 1
    max_comp = tfidf_matrix.shape[0] - 1  # 37 for 38 samples
    print(f"\n--- Projecting (max PCA components: {max_comp}) ---")

    n48 = min(48, max_comp)
    pca_48 = PCA(n_components=n48, random_state=42)
    proj_pca_48 = pca_48.fit_transform(tfidf_matrix)
    print(f"PCA-{n48} explained variance: {pca_48.explained_variance_ratio_.sum():.3f}")

    n96 = min(96, max_comp)
    pca_96 = PCA(n_components=n96, random_state=42)
    proj_pca_96 = pca_96.fit_transform(tfidf_matrix)

    n192 = min(192, max_comp)
    pca_192 = PCA(n_components=n192, random_state=42)
    proj_pca_192 = pca_192.fit_transform(tfidf_matrix)

    # Random projection (JL lemma)
    rp = GaussianRandomProjection(n_components=192, random_state=42)
    proj_rp = rp.fit_transform(tfidf_matrix)

    # ============================================================
    # APPROACH 1: TF-IDF → PCA-192 → Braille 4×6 grid
    # ============================================================
    print("\n" + "=" * 72)
    print("  APPROACH 1: PCA → Braille Dot Matrix (4×6)")
    print("=" * 72)

    sigils = {}
    for i, k in enumerate(keys):
        sigils[k] = vector_to_braille_grid(proj_pca_192[i], rows=4, cols=6, threshold=0.0)

    # Show Matt across MOs
    matt_keys = [k for k in keys if k.startswith('artifex_')]
    print_sigil_comparison("Matt across MOs (should look similar)", sigils, matt_keys, labels)

    # Show Kelly across MOs
    kelly_keys = [k for k in keys if k.startswith('kelly_')]
    print_sigil_comparison("Kelly across MOs (should look similar)", sigils, kelly_keys, labels)

    # Show different types
    type_sample = ['artifex_alignment', 'behavioral_audit', 'alignment_intent', 'alignment_tempering']
    print_sigil_comparison("Different types from same MO (should look different)", sigils, type_sample, labels)

    # Show calcinatio principles
    calc_keys = list(CALCINATIO.keys())[:6]
    print_sigil_comparison("Calcinatio principles (varied)", sigils, calc_keys, labels)

    # ============================================================
    # APPROACH 2: Random Projection → Braille (JL lemma)
    # ============================================================
    print("\n" + "=" * 72)
    print("  APPROACH 2: Random Projection → Braille (JL lemma)")
    print("=" * 72)

    sigils_rp = {}
    for i, k in enumerate(keys):
        sigils_rp[k] = vector_to_braille_grid(proj_rp[i], rows=4, cols=6, threshold=0.0)

    print_sigil_comparison("Matt across MOs (RP)", sigils_rp, matt_keys, labels)
    print_sigil_comparison("Different types same MO (RP)", sigils_rp, type_sample, labels)

    # ============================================================
    # APPROACH 3: PCA → Block Element Density Field
    # ============================================================
    print("\n" + "=" * 72)
    print("  APPROACH 3: PCA → Block Element Density Field (6×12)")
    print("=" * 72)

    sigils_block = {}
    for i, k in enumerate(keys):
        vec = proj_pca_192[i] if len(proj_pca_192[i]) >= 72 else np.concatenate([proj_pca_192[i], np.zeros(72)])
        sigils_block[k] = vector_to_density_field(vec, rows=6, cols=12)

    print_sigil_comparison("Matt across MOs (blocks)", sigils_block, matt_keys, labels)
    print_sigil_comparison("Different types same MO (blocks)", sigils_block, type_sample, labels)

    # ============================================================
    # APPROACH 4: PCA → Radial Geometric Sigil
    # ============================================================
    print("\n" + "=" * 72)
    print("  APPROACH 4: PCA → Radial Geometric Sigil (9×9)")
    print("=" * 72)

    sigils_radial = {}
    for i, k in enumerate(keys):
        sigils_radial[k] = vector_to_radial_sigil(proj_pca_48[i], size=9)

    print_sigil_comparison("Matt across MOs (radial)", sigils_radial, matt_keys, labels)
    print_sigil_comparison("Different types same MO (radial)", sigils_radial, type_sample, labels)

    # ============================================================
    # APPROACH 5: SimHash → Symbol Composition
    # ============================================================
    print("\n" + "=" * 72)
    print("  APPROACH 5: SimHash → Symbol Composition (4×4)")
    print("=" * 72)

    sigils_sim = {}
    hashes = {}
    for k, d in all_inputs:
        h = simhash(d['text'])
        hashes[k] = h
        sigils_sim[k] = simhash_to_symbols(h)

    print_sigil_comparison("Matt across MOs (SimHash)", sigils_sim, matt_keys, labels)

    # SimHash distances
    print("\n  SimHash Hamming Distances (lower = more similar):")
    print(f"  {'Comparison':<48} {'Distance':>8}")
    print("  " + "─" * 58)
    for label, k1, k2 in comparisons[:8]:
        dist = simhash_distance(hashes[k1], hashes[k2])
        bar = '█' * (32 - dist) + '░' * dist
        print(f"  {label:<48} {dist:>4}/64  {bar}")

    # ============================================================
    # APPROACH 6: PCA → Compound Sigil with bilateral symmetry
    # ============================================================
    print("\n" + "=" * 72)
    print("  APPROACH 6: PCA → Bilateral Symmetric Sigil (7×7)")
    print("=" * 72)

    sigils_compound = {}
    for i, k in enumerate(keys):
        sigils_compound[k] = vector_to_compound_sigil(proj_pca_48[i], size=7)

    print_sigil_comparison("Matt across MOs (symmetric)", sigils_compound, matt_keys, labels)
    print_sigil_comparison("Kelly across MOs (symmetric)", sigils_compound, kelly_keys, labels)
    print_sigil_comparison("Different types same MO (symmetric)", sigils_compound, type_sample, labels)
    print_sigil_comparison("All intents (symmetric)", sigils_compound, list(INTENTS.keys()), labels)

    # ============================================================
    # APPROACH 7: PCA → Mandala (4-fold symmetric Braille)
    # ============================================================
    print("\n" + "=" * 72)
    print("  APPROACH 7: PCA → 4-fold Symmetric Mandala (Braille)")
    print("=" * 72)

    sigils_mandala = {}
    for i, k in enumerate(keys):
        sigils_mandala[k] = vector_to_mandala(proj_pca_192[i], quadrant_rows=3, quadrant_cols=4)

    print_sigil_comparison("Matt across MOs (mandala)", sigils_mandala, matt_keys, labels)
    print_sigil_comparison("Kelly across MOs (mandala)", sigils_mandala, kelly_keys, labels)
    print_sigil_comparison("Tom across MOs (mandala)", sigils_mandala,
                          [k for k in keys if k.startswith('tom_')], labels)
    print_sigil_comparison("Leadership across MOs (mandala)", sigils_mandala,
                          [k for k in keys if k.startswith('leadership_')], labels)
    print_sigil_comparison("Different types same MO (mandala)", sigils_mandala, type_sample, labels)
    print_sigil_comparison("All intents (mandala)", sigils_mandala, list(INTENTS.keys()), labels)
    print_sigil_comparison("All calcinatio (mandala)", sigils_mandala, list(CALCINATIO.keys())[:6], labels)

    # ============================================================
    # SIMILARITY PRESERVATION ANALYSIS
    # ============================================================
    print("\n" + "=" * 72)
    print("  SIMILARITY PRESERVATION ANALYSIS")
    print("=" * 72)

    print("\n  Do similar inputs produce similar outputs?")
    print("  Comparing TF-IDF cosine similarity vs visual similarity proxy")
    print("  (Visual similarity = cosine similarity of projected vectors)\n")

    proj_sim = cosine_similarity(proj_pca_192)

    # Group pairs by expected similarity
    same_person = [
        ("artifex_alignment", "artifex_harvest"),
        ("artifex_alignment", "artifex_anvil"),
        ("artifex_harvest", "artifex_anvil"),
        ("kelly_seismic", "kelly_analytics"),
        ("kelly_seismic", "kelly_sf"),
        ("tom_seismic", "tom_sf"),
        ("leadership_seismic", "leadership_sf"),
    ]

    same_type_diff_content = [
        ("artifex_alignment", "horse_owners"),
        ("artifex_alignment", "appnigma"),
        ("kelly_seismic", "tom_seismic"),
        ("behavioral_audit", "customer_outcome"),
        ("alignment_intent", "seismic_intent"),
    ]

    cross_type = [
        ("artifex_alignment", "behavioral_audit"),
        ("artifex_alignment", "alignment_intent"),
        ("behavioral_audit", "alignment_intent"),
        ("horse_owners", "horse_intent"),
        ("kelly_seismic", "confidence_chain"),
    ]

    for group_name, pairs in [
        ("Same person across MOs", same_person),
        ("Same type, different content", same_type_diff_content),
        ("Cross-type (witness vs calcinatio vs intent)", cross_type),
    ]:
        tfidf_sims = []
        proj_sims = []
        for k1, k2 in pairs:
            i1, i2 = key_to_idx[k1], key_to_idx[k2]
            tfidf_sims.append(sim[i1, i2])
            proj_sims.append(proj_sim[i1, i2])

        print(f"  {group_name}:")
        print(f"    TF-IDF cosine sim (mean): {np.mean(tfidf_sims):.3f}")
        print(f"    Projected cosine sim (mean): {np.mean(proj_sims):.3f}")
        print(f"    Correlation preserved: {'YES' if np.mean(proj_sims) > 0.3 else 'WEAK' if np.mean(proj_sims) > 0.1 else 'NO'}")
        print()


if __name__ == '__main__':
    main()
