#!/usr/bin/env python3
"""
Quality-based SVSI renderer.

Architecture:
  text → embedding → arcana profile (78-dim) → quality decomposition → render

The renderer doesn't know about tarot. It knows about qualities:
  elements, modalities, geometry, color, dynamics, density, symmetry.

Each arcanum decomposes into a weighted quality vector.
The text's quality profile = weighted sum of its arcana qualities.
The renderer maps qualities to visual primitives.
"""
import sys, os, math, random
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

try:
    from sentence_transformers import SentenceTransformer
    model = SentenceTransformer('all-MiniLM-L6-v2')
    def embed(texts):
        return model.encode(texts)
except Exception:
    print("Need sentence-transformers"); sys.exit(1)


# ============================================================
# QUALITY VOCABULARY — the renderer's input language
# ============================================================

QUALITIES = [
    # Elemental
    "fire", "water", "air", "earth",
    # Modalities (Aristotelian)
    "hot", "cold", "wet", "dry",
    # Dynamics
    "expansion", "contraction", "ascending", "descending",
    "radiating", "concentrating", "flowing", "stillness",
    "rotation", "oscillation", "explosive", "dissolving",
    # Geometric
    "point", "line", "triangle", "square", "pentagon",
    "hexagram", "circle", "spiral", "cross", "crescent",
    # Density
    "sparse", "moderate", "dense", "saturated",
    # Symmetry
    "asymmetric", "bilateral", "radial", "fourfold",
    # Chromatic (simplified from GD scales)
    "white", "black", "red", "blue", "yellow", "green",
    "orange", "violet", "indigo", "gold",
    # Abstract
    "balance", "force", "receptivity", "transformation",
    "structure", "liberation", "unity", "duality",
    "mystery", "clarity", "endurance", "swiftness",
]

Q_IDX = {q: i for i, q in enumerate(QUALITIES)}
N_QUAL = len(QUALITIES)


def q_vec(**kwargs):
    """Build a quality vector from keyword weights."""
    v = np.zeros(N_QUAL)
    for k, w in kwargs.items():
        if k in Q_IDX:
            v[Q_IDX[k]] = w
    return v


# ============================================================
# 78 ARCANA → quality decompositions
# ============================================================

# Major Arcana (22) — each with quality weights and description for embedding
MAJOR = [
    ("The Fool", "New beginnings, leaping into unknown, pure potential, innocence, spontaneity",
     q_vec(air=0.9, ascending=0.8, sparse=0.6, spiral=0.7, yellow=0.7,
           asymmetric=1.0, liberation=0.8, swiftness=0.5)),
    ("The Magician", "Skill, will made manifest, tools mastered, focused intent, channeling power",
     q_vec(air=0.5, fire=0.5, hot=0.6, concentrating=0.8, line=0.7, yellow=0.6,
           bilateral=0.7, dense=0.6, force=0.7, clarity=0.8, swiftness=0.6)),
    ("High Priestess", "Hidden knowledge, intuition, mystery, veil between seen and unseen",
     q_vec(water=0.8, cold=0.5, wet=0.6, stillness=0.9, crescent=0.8, blue=0.7,
           bilateral=0.8, sparse=0.6, mystery=0.9, receptivity=0.8)),
    ("The Empress", "Abundance, fertility, nurturing, creative growth, generative force",
     q_vec(earth=0.5, water=0.5, wet=0.7, expansion=0.9, circle=0.7, green=0.9,
           radial=0.6, dense=0.8, receptivity=0.6, endurance=0.5)),
    ("The Emperor", "Authority, structure, control, stability, order imposed on chaos",
     q_vec(fire=0.6, hot=0.5, dry=0.7, stillness=0.6, square=0.9, red=0.8,
           fourfold=0.9, dense=0.7, structure=0.9, force=0.6)),
    ("The Hierophant", "Teaching, tradition, institutional wisdom, bridging divine and mundane",
     q_vec(earth=0.6, dry=0.5, cross=0.8, orange=0.5, bilateral=0.7,
           moderate=0.7, structure=0.7, unity=0.5, endurance=0.6)),
    ("The Lovers", "Choice, relationship, duality resolved, union of opposites, alignment",
     q_vec(air=0.7, hot=0.4, wet=0.4, oscillation=0.5, hexagram=0.6, orange=0.6,
           bilateral=0.9, moderate=0.5, duality=0.9, balance=0.6)),
    ("The Chariot", "Willpower, directed force, triumph through control, moving forward",
     q_vec(water=0.5, fire=0.4, hot=0.6, ascending=0.7, triangle=0.5, orange=0.5,
           bilateral=0.8, dense=0.6, force=0.8, swiftness=0.7)),
    ("Strength", "Inner strength, patience, gentle power, endurance, compassion taming force",
     q_vec(fire=0.6, hot=0.4, flowing=0.5, circle=0.5, green=0.5, yellow=0.4,
           bilateral=0.6, moderate=0.7, endurance=0.9, balance=0.6)),
    ("The Hermit", "Solitude, inner light, withdrawal for wisdom, lamp in darkness",
     q_vec(earth=0.7, cold=0.6, dry=0.5, contracting=0.7, point=0.9, green=0.4,
           asymmetric=0.7, sparse=0.8, mystery=0.6, clarity=0.5, contraction=0.7)),
    ("Wheel of Fortune", "Cycles, fortune turning, change, karma, the turning wheel",
     q_vec(fire=0.3, water=0.3, rotation=0.9, circle=0.8, violet=0.7,
           radial=0.8, moderate=0.5, transformation=0.7, swiftness=0.5)),
    ("Justice", "Balance, fairness, truth, consequences, weighing evidence, assessment",
     q_vec(air=0.8, hot=0.3, dry=0.6, stillness=0.6, balance=0.9, green=0.5,
           bilateral=0.9, moderate=0.6, clarity=0.8, structure=0.5)),
    ("The Hanged Man", "Surrender, seeing differently, suspension, sacrifice for wisdom",
     q_vec(water=0.9, cold=0.5, wet=0.7, stillness=0.7, descending=0.5,
           triangle=0.5, blue=0.8, bilateral=0.6, sparse=0.5, receptivity=0.9,
           transformation=0.6)),
    ("Death", "Transformation, endings that enable beginnings, clearing, profound change",
     q_vec(water=0.6, cold=0.7, dissolving=0.8, crescent=0.5, blue=0.4, indigo=0.5,
           asymmetric=0.6, moderate=0.6, transformation=0.9, liberation=0.5)),
    ("Temperance", "Moderation, blending, patience, the middle way, careful mixing",
     q_vec(water=0.5, fire=0.4, wet=0.5, flowing=0.8, triangle=0.4, blue=0.5,
           bilateral=0.8, moderate=0.8, balance=0.8, duality=0.5, endurance=0.5)),
    ("The Devil", "Bondage, materialism, shadow, chains, confronting limitation",
     q_vec(earth=0.8, cold=0.4, dry=0.7, contraction=0.7, pentagon=0.6, indigo=0.7,
           bilateral=0.6, saturated=0.7, structure=0.6, force=0.5)),
    ("The Tower", "Sudden disruption, structures falling, liberation through destruction",
     q_vec(fire=0.9, hot=0.9, explosive=0.9, line=0.6, red=0.9,
           asymmetric=0.8, dense=0.5, liberation=0.8, force=0.7, transformation=0.6)),
    ("The Star", "Hope, inspiration, renewal, clarity after storm, guiding light",
     q_vec(air=0.6, water=0.4, wet=0.4, radiating=0.9, hexagram=0.7, violet=0.5,
           radial=0.9, sparse=0.6, clarity=0.7, receptivity=0.5)),
    ("The Moon", "Illusion, unconscious, uncertainty, dreams, what shifts in moonlight",
     q_vec(water=0.8, cold=0.6, wet=0.8, oscillation=0.7, crescent=0.9, violet=0.5,
           bilateral=0.5, moderate=0.5, mystery=0.9, duality=0.5)),
    ("The Sun", "Joy, success, clarity, vitality, everything illuminated, confidence",
     q_vec(fire=0.7, hot=0.7, radiating=0.8, circle=0.7, gold=0.9, orange=0.6,
           radial=0.8, dense=0.6, clarity=0.9, force=0.5)),
    ("Judgement", "Awakening, reckoning, answering call, rising to purpose, assessment",
     q_vec(fire=0.8, hot=0.6, ascending=0.8, triangle=0.6, red=0.6, gold=0.4,
           radial=0.5, moderate=0.6, transformation=0.7, clarity=0.6, force=0.5)),
    ("The World", "Completion, integration, wholeness, dance of accomplished purpose",
     q_vec(earth=0.5, water=0.3, fire=0.3, air=0.3, rotation=0.5, circle=0.6,
           indigo=0.6, fourfold=0.7, saturated=0.8, unity=0.9, balance=0.5,
           endurance=0.6)),
]

# Minor Arcana (56) — generated from suit × number/court
SUITS = {
    "Wands": {"element": "fire", "color": "red", "quality": "will action energy creative force",
              "q": q_vec(fire=0.8, hot=0.7, dry=0.4, force=0.6, swiftness=0.5)},
    "Cups": {"element": "water", "color": "blue", "quality": "emotion flow depth feeling connection",
             "q": q_vec(water=0.8, cold=0.4, wet=0.7, flowing=0.5, receptivity=0.6)},
    "Swords": {"element": "air", "color": "yellow", "quality": "thought communication clarity cutting truth",
               "q": q_vec(air=0.8, hot=0.4, dry=0.5, clarity=0.6, swiftness=0.5, force=0.4)},
    "Pentacles": {"element": "earth", "color": "green", "quality": "material body stability craft manifestation",
                  "q": q_vec(earth=0.8, cold=0.4, dry=0.6, structure=0.6, endurance=0.6)},
}

NUMBERS = {
    1: ("Ace", "root pure seed beginning", q_vec(point=0.8, sparse=0.6, unity=0.8, ascending=0.5)),
    2: ("Two", "duality polarity balance choice", q_vec(line=0.6, bilateral=0.8, duality=0.9, balance=0.5)),
    3: ("Three", "synthesis growth expression", q_vec(triangle=0.8, expansion=0.6, moderate=0.6)),
    4: ("Four", "stability structure foundation", q_vec(square=0.9, fourfold=0.8, structure=0.7, stillness=0.6)),
    5: ("Five", "conflict disruption change", q_vec(pentagon=0.6, oscillation=0.6, force=0.5, transformation=0.5)),
    6: ("Six", "harmony beauty balance", q_vec(hexagram=0.8, radial=0.6, balance=0.8, flowing=0.4)),
    7: ("Seven", "assessment reflection choice", q_vec(asymmetric=0.5, mystery=0.5, contraction=0.4)),
    8: ("Eight", "movement mastery speed", q_vec(circle=0.5, swiftness=0.7, dense=0.5, force=0.4)),
    9: ("Nine", "culmination near-completion strength", q_vec(dense=0.7, endurance=0.7, saturated=0.5)),
    10: ("Ten", "completion manifestation fullness", q_vec(circle=0.6, saturated=0.8, unity=0.5, descending=0.4)),
}

COURTS = {
    "Page": ("Page", "student beginning learning message", q_vec(sparse=0.5, ascending=0.5, air=0.3)),
    "Knight": ("Knight", "action movement quest pursuit", q_vec(swiftness=0.8, force=0.5, fire=0.3)),
    "Queen": ("Queen", "mastery receptive nurturing depth", q_vec(water=0.3, receptivity=0.7, dense=0.5)),
    "King": ("King", "authority command mastery expression", q_vec(fire=0.3, force=0.6, structure=0.5, saturated=0.5)),
}

def build_minor_arcana():
    """Generate 56 minor arcana from suit × (number + court)."""
    cards = []
    for suit_name, suit in SUITS.items():
        for num, (name, desc, num_q) in NUMBERS.items():
            card_name = f"{name} of {suit_name}"
            card_desc = f"{name} of {suit_name}. {desc}. {suit['quality']}."
            card_q = 0.5 * suit["q"] + 0.5 * num_q  # blend suit and number qualities
            cards.append((card_name, card_desc, card_q))
        for court_name, (name, desc, court_q) in COURTS.items():
            card_name = f"{court_name} of {suit_name}"
            card_desc = f"{court_name} of {suit_name}. {desc}. {suit['quality']}."
            card_q = 0.5 * suit["q"] + 0.5 * court_q
            cards.append((card_name, card_desc, card_q))
    return cards

ALL_ARCANA = MAJOR + build_minor_arcana()  # 22 + 56 = 78


# ============================================================
# QUALITY-DRIVEN RENDERER
# ============================================================

# Map qualities to visual primitives
SHAPE_CHARS = {
    "point": ["·", "∙", "•"],
    "line": ["│", "─", "║", "═"],
    "triangle": ["△", "▽", "◁", "▷"],
    "square": ["□", "■", "▪", "▫"],
    "pentagon": ["⬠", "⬡"],
    "hexagram": ["✡", "✦", "⬡", "⬢"],
    "circle": ["○", "◎", "●", "◉"],
    "spiral": ["◌", "◍", "⊛", "⊙"],
    "cross": ["✚", "✛", "╋", "┼"],
    "crescent": ["☽", "☾", "◗", "◖"],
}

COLOR_ANSI = {
    "white": "\033[97m", "black": "\033[90m", "red": "\033[91m",
    "blue": "\033[94m", "yellow": "\033[93m", "green": "\033[92m",
    "orange": "\033[33m", "violet": "\033[35m", "indigo": "\033[34m",
    "gold": "\033[93m",
}
RESET = "\033[0m"

SYMMETRY_OPS = {
    "asymmetric": lambda g: g,  # no mirror
    "bilateral": lambda g: mirror_h(g),
    "radial": lambda g: mirror_hv(g),
    "fourfold": lambda g: mirror_hv(g),
}

def mirror_h(half_lines):
    """Mirror lines horizontally."""
    result = []
    for line in half_lines:
        result.append(line + line[::-1])
    return result

def mirror_hv(quarter_lines):
    """Mirror both horizontally and vertically."""
    full_h = mirror_h(quarter_lines)
    return full_h + full_h[::-1]


def render_from_qualities(q_profile, size=9, stochastic=0.0, seed=None):
    """Render a sigil from a quality profile vector.

    stochastic: 0.0 = fully deterministic, 1.0 = maximum randomness
    """
    rng = random.Random(seed)

    # Extract dominant qualities
    top_indices = np.argsort(q_profile)[::-1]
    top_quals = [(QUALITIES[i], q_profile[i]) for i in top_indices if q_profile[i] > 0.1]

    # Determine primary visual properties from quality weights
    # Color: weighted blend of chromatic qualities
    color_quals = {c: q_profile[Q_IDX[c]] for c in
                   ["white","black","red","blue","yellow","green","orange","violet","indigo","gold"]
                   if q_profile[Q_IDX[c]] > 0.05}
    primary_color = max(color_quals, key=color_quals.get) if color_quals else "white"
    secondary_color = sorted(color_quals, key=color_quals.get, reverse=True)[1] if len(color_quals) > 1 else primary_color

    # Shape: highest geometric quality
    shape_quals = {s: q_profile[Q_IDX[s]] for s in SHAPE_CHARS.keys() if q_profile[Q_IDX[s]] > 0.05}
    primary_shape = max(shape_quals, key=shape_quals.get) if shape_quals else "circle"

    # Symmetry
    sym_quals = {s: q_profile[Q_IDX[s]] for s in ["asymmetric","bilateral","radial","fourfold"] if q_profile[Q_IDX[s]] > 0.05}
    symmetry = max(sym_quals, key=sym_quals.get) if sym_quals else "bilateral"

    # Density
    dens_quals = {d: q_profile[Q_IDX[d]] for d in ["sparse","moderate","dense","saturated"] if q_profile[Q_IDX[d]] > 0.05}
    density = max(dens_quals, key=dens_quals.get) if dens_quals else "moderate"
    density_val = {"sparse": 0.25, "moderate": 0.5, "dense": 0.75, "saturated": 0.95}[density]

    # Elemental balance
    elements = {e: q_profile[Q_IDX[e]] for e in ["fire","water","air","earth"]}

    # Dynamics
    dynamics = {d: q_profile[Q_IDX[d]] for d in
                ["expansion","contraction","ascending","descending","radiating",
                 "concentrating","flowing","stillness","rotation","oscillation",
                 "explosive","dissolving"] if q_profile[Q_IDX[d]] > 0.05}
    primary_dynamic = max(dynamics, key=dynamics.get) if dynamics else "stillness"

    # Abstract qualities for character selection
    abstracts = {a: q_profile[Q_IDX[a]] for a in
                 ["balance","force","receptivity","transformation","structure",
                  "liberation","unity","duality","mystery","clarity","endurance","swiftness"]
                 if q_profile[Q_IDX[a]] > 0.05}

    # --- BUILD THE SIGIL ---

    c1 = COLOR_ANSI.get(primary_color, "")
    c2 = COLOR_ANSI.get(secondary_color, "")

    # Select shape characters (with optional stochastic variation)
    shapes = SHAPE_CHARS.get(primary_shape, ["◉"])
    if stochastic > 0 and rng.random() < stochastic:
        # Sometimes pick from secondary shape
        if len(shape_quals) > 1:
            alt_shape = sorted(shape_quals, key=shape_quals.get, reverse=True)[1]
            shapes = SHAPE_CHARS.get(alt_shape, shapes)

    sym_main = shapes[0]
    sym_alt = shapes[1] if len(shapes) > 1 else shapes[0]
    sym_small = shapes[-1] if len(shapes) > 1 else "·"

    # Element symbols for accents
    elem_syms = {"fire": "△", "water": "▽", "air": "◇", "earth": "□"}
    top_elem = max(elements, key=elements.get)
    sec_elem = sorted(elements, key=elements.get, reverse=True)[1]
    e1 = elem_syms[top_elem]
    e2 = elem_syms[sec_elem]

    # Dynamic modifiers
    if primary_dynamic in ["ascending", "radiating", "expansion"]:
        flow_char = "↑" if primary_dynamic == "ascending" else "✦" if primary_dynamic == "radiating" else "◌"
    elif primary_dynamic in ["descending", "concentrating", "contraction"]:
        flow_char = "↓" if primary_dynamic == "descending" else "◎" if primary_dynamic == "concentrating" else "◍"
    elif primary_dynamic in ["rotation", "oscillation"]:
        flow_char = "⟳" if primary_dynamic == "rotation" else "⇌"
    elif primary_dynamic == "explosive":
        flow_char = "✴"
    elif primary_dynamic == "flowing":
        flow_char = "≈"
    else:
        flow_char = "·"

    # Density determines fill
    if density_val > 0.8:
        fill = sym_main
        bg = sym_alt
    elif density_val > 0.6:
        fill = sym_main
        bg = sym_small
    elif density_val > 0.4:
        fill = sym_small
        bg = " "
    else:
        fill = "·"
        bg = " "

    # Build based on symmetry type
    half = size // 2

    if symmetry == "radial" or symmetry == "fourfold":
        # Build one quadrant, mirror 4 ways
        quad = []
        for r in range(half + 1):
            row = ""
            for c in range(half + 1):
                dist = math.sqrt((r - half)**2 + (c - half)**2) / half
                if dist < 0.3:
                    ch = sym_main
                elif dist < 0.6:
                    if (r + c) % 2 == 0:
                        ch = fill
                    else:
                        ch = flow_char if dist < 0.5 else bg
                elif dist < 0.85:
                    ch = e1 if (r * c) % 3 == 0 else bg
                else:
                    ch = e2 if r == 0 or c == 0 else bg

                if stochastic > 0 and rng.random() < stochastic * 0.3:
                    ch = sym_small if ch != " " else " "

                row += ch
            quad.append(row)

        # Mirror to full grid
        lines = []
        for r in range(size):
            row = ""
            for c in range(size):
                qr = r if r <= half else size - 1 - r
                qc = c if c <= half else size - 1 - c
                row += quad[qr][qc]
            lines.append(f"  {c1}{row}{RESET}")

    elif symmetry == "bilateral":
        # Build left half, mirror right
        lines = []
        for r in range(size):
            row = ""
            for c in range(half + 1):
                dist_center = abs(r - half) / half
                dist_edge = c / half

                if dist_center < 0.3 and dist_edge < 0.3:
                    ch = sym_main
                elif dist_center < 0.5:
                    ch = fill if (r + c) % 2 == 0 else flow_char
                elif dist_center < 0.8:
                    ch = e1 if c % 2 == 0 else bg
                else:
                    ch = e2 if c == 0 else bg

                if stochastic > 0 and rng.random() < stochastic * 0.3:
                    ch = sym_small if ch != " " else " "

                row += ch
            # Mirror
            full = row + row[-2::-1]
            lines.append(f"  {c1}{full}{RESET}")

    else:  # asymmetric
        lines = []
        for r in range(size):
            row = ""
            for c in range(size):
                dist = math.sqrt((r - half)**2 + (c - half)**2) / half
                angle = math.atan2(r - half, c - half)

                if dist < 0.3:
                    ch = sym_main
                elif dist < 0.7:
                    ch = fill if math.sin(angle * 3 + dist * 5) > 0 else bg
                else:
                    ch = e1 if math.sin(angle * 2) > 0.5 else bg

                row += ch
            lines.append(f"  {c1}{row}{RESET}")

    return lines, {
        "primary_color": primary_color,
        "secondary_color": secondary_color,
        "primary_shape": primary_shape,
        "symmetry": symmetry,
        "density": density,
        "primary_dynamic": primary_dynamic,
        "top_element": top_elem,
        "top_qualities": [(q, w) for q, w in top_quals[:8]],
    }


# ============================================================
# MAIN
# ============================================================

INPUTS = [
    ("Artifex witness (Alignment)",
     "The primary witness. The athanor exists to reduce his executive function load. Not holding the system in his head. Athanors that behave as expected without correction. Trails he can read without effort. The system improving from evidence, not from his intervention."),
    ("Artifex witness (Harvest)",
     "He feels secure — not anxious that the productivity delta will be discovered awkwardly. He feels compensated — the value is visibly reciprocated. He feels protected — the gains don't isolate him. He feels portable — the system is his IP, his career asset."),
    ("Horse Owners witness",
     "Horse owners see a demo of their horse's care history and have a genuine emotional reaction — excitement, relief, I wish my stable did this. The demo is compelling enough that the reaction is unambiguous."),
    ("Calcinatio: Behavioral audit",
     "Behavioral audit is the highest-value fire. Reading trails for patterns — calcinatio skipped, geas satisfied superficially, escalations absent — reveals misalignment no automated check can find."),
    ("Calcinatio: Customer outcome",
     "Customer outcome fidelity is the gate. Seismic API: case-sensitive exact string matching, atomic rejection. Broken mappings equal customer outcome problem, not UI problem."),
    ("Intent: Alignment",
     "The athanor system's actual behavior progressively aligns with the core intents. Agents apply calcinatio without being told. Trails are legible. This is a maintained condition, not a build goal."),
    ("Tempering: Stabilization",
     "Current focus: Stabilization before release. Feature is functionally complete. Bar for new opera is high — only customer-facing risks."),
]


def main():
    # Embed all arcana descriptions
    arcana_texts = [desc for _, desc, _ in ALL_ARCANA]
    arcana_qualities = np.array([q for _, _, q in ALL_ARCANA])
    arcana_names = [name for name, _, _ in ALL_ARCANA]

    input_texts = [t for _, t in INPUTS]

    print("Embedding 78 arcana + inputs...")
    arcana_emb = embed(arcana_texts)
    input_emb = embed(input_texts)

    sims = cosine_similarity(input_emb, arcana_emb)

    print(f"\n{'='*72}")
    print(f"  QUALITY-BASED SVSI RENDERER")
    print(f"  text → 78 arcana → qualities → visual")
    print(f"  {len(ALL_ARCANA)} arcana, {N_QUAL} qualities")
    print(f"{'='*72}")

    all_quality_profiles = []

    for i, (label, text) in enumerate(INPUTS):
        scores = sims[i]

        # Compute quality profile: weighted sum of arcana qualities
        # Weight by similarity score (only positive contributions)
        weights = np.maximum(scores, 0)
        weights = weights / (weights.sum() + 1e-10)  # normalize
        quality_profile = arcana_qualities.T @ weights  # (n_qualities,)
        all_quality_profiles.append(quality_profile)

        # Top arcana
        top_idx = np.argsort(scores)[::-1][:5]

        print(f"\n{'─'*72}")
        print(f"  {label}")
        print(f"  Top arcana: ", end="")
        for j, idx in enumerate(top_idx[:3]):
            print(f"{arcana_names[idx]} ({scores[idx]:.2f})", end="  ")
        print()

        # Render
        sigil, props = render_from_qualities(quality_profile, size=9, stochastic=0.0)
        print(f"  Shape: {props['primary_shape']}, Sym: {props['symmetry']}, "
              f"Color: {props['primary_color']}/{props['secondary_color']}, "
              f"Density: {props['density']}, Dynamic: {props['primary_dynamic']}, "
              f"Element: {props['top_element']}")
        print(f"  Top qualities: {', '.join(f'{q}({w:.2f})' for q,w in props['top_qualities'][:6])}")
        print()
        for line in sigil:
            print(f"  {line}")

        # Also render stochastic variant
        print(f"\n  Stochastic variant (0.3):")
        sigil2, _ = render_from_qualities(quality_profile, size=9, stochastic=0.3, seed=42)
        for line in sigil2:
            print(f"  {line}")

    # Quality-space similarity
    qp = np.array(all_quality_profiles)
    qsim = cosine_similarity(qp)

    print(f"\n{'='*72}")
    print(f"  QUALITY-SPACE SIMILARITY")
    print(f"{'='*72}")

    pairs = [
        (0, 1, "Artifex/Alignment vs Artifex/Harvest"),
        (0, 3, "Artifex/Alignment vs Calcinatio:Audit"),
        (0, 2, "Artifex vs Horse Owners"),
        (3, 4, "Calcinatio:Audit vs Calcinatio:Customer"),
        (0, 5, "Artifex witness vs Intent:Alignment"),
        (5, 6, "Intent:Alignment vs Tempering:Stabilization"),
    ]

    print(f"\n  {'Pair':<48} {'Sim':>6}")
    print(f"  {'─'*58}")
    for i1, i2, label in pairs:
        s = qsim[i1, i2]
        bar = '█' * int(s * 20) + '░' * (20 - int(s * 20))
        print(f"  {label:<48} {s:.3f} {bar}")

    print()


if __name__ == '__main__':
    main()
