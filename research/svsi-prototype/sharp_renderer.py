#!/usr/bin/env python3
"""
Sharp SVSI renderer — 22 major arcana only, temperature-controlled selection,
correspondence-driven rendering.
"""
import sys, os, math, random
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity
from sentence_transformers import SentenceTransformer

model = SentenceTransformer('all-MiniLM-L6-v2')
def embed(texts):
    return model.encode(texts)

RESET = "\033[0m"
ANSI = {
    "white": "\033[97m", "black": "\033[90m", "red": "\033[91m",
    "blue": "\033[94m", "yellow": "\033[93m", "green": "\033[92m",
    "orange": "\033[33m", "violet": "\033[35m", "indigo": "\033[34m",
    "gold": "\033[93;1m", "crimson": "\033[31m",
}

# ============================================================
# 22 MAJOR ARCANA with correspondence decompositions
# ============================================================
# Each has: name, description (for embedding), and correspondences
# split into major (always) and minor (stochastic)

MAJORS = [
    {
        "name": "The Fool", "symbol": "○", "hebrew": "א",
        "desc": "New beginnings, leaping into unknown, pure potential, innocence, spontaneity, the leap of faith",
        "major_corr": {"element": "air", "dynamic": "ascending", "density": "sparse", "color": "yellow"},
        "minor_corr": {"shape": "spiral", "symmetry": "asymmetric", "quality": "liberation", "motion": "erratic"},
    },
    {
        "name": "The Magician", "symbol": "☿", "hebrew": "ב",
        "desc": "Skill, will made manifest, tools mastered, focused intent, channeling power, communication",
        "major_corr": {"element": "air", "dynamic": "concentrating", "density": "dense", "color": "yellow"},
        "minor_corr": {"shape": "line", "symmetry": "bilateral", "quality": "clarity", "planet": "mercury"},
    },
    {
        "name": "High Priestess", "symbol": "☽", "hebrew": "ג",
        "desc": "Hidden knowledge, intuition, mystery, veil between seen and unseen, inner wisdom, reflection",
        "major_corr": {"element": "water", "dynamic": "stillness", "density": "sparse", "color": "blue"},
        "minor_corr": {"shape": "crescent", "symmetry": "bilateral", "quality": "mystery", "planet": "moon"},
    },
    {
        "name": "The Empress", "symbol": "♀", "hebrew": "ד",
        "desc": "Abundance, fertility, nurturing, creative growth, natural expansion, generative force",
        "major_corr": {"element": "earth", "dynamic": "expanding", "density": "dense", "color": "green"},
        "minor_corr": {"shape": "circle", "symmetry": "radial", "quality": "receptivity", "planet": "venus"},
    },
    {
        "name": "The Emperor", "symbol": "♈", "hebrew": "ה",
        "desc": "Authority, structure, control, stability, order imposed on chaos, the father, command",
        "major_corr": {"element": "fire", "dynamic": "stable", "density": "dense", "color": "red"},
        "minor_corr": {"shape": "square", "symmetry": "fourfold", "quality": "structure", "modality": "hot_dry"},
    },
    {
        "name": "Hierophant", "symbol": "♉", "hebrew": "ו",
        "desc": "Teaching, tradition, institutional wisdom, bridging divine and mundane, doctrine",
        "major_corr": {"element": "earth", "dynamic": "stable", "density": "moderate", "color": "orange"},
        "minor_corr": {"shape": "cross", "symmetry": "bilateral", "quality": "endurance", "modality": "cold_dry"},
    },
    {
        "name": "The Lovers", "symbol": "♊", "hebrew": "ז",
        "desc": "Choice, relationship, duality resolved, union of opposites, alignment of values",
        "major_corr": {"element": "air", "dynamic": "oscillating", "density": "moderate", "color": "orange"},
        "minor_corr": {"shape": "hexagram", "symmetry": "bilateral", "quality": "duality", "modality": "hot_wet"},
    },
    {
        "name": "The Chariot", "symbol": "♋", "hebrew": "ח",
        "desc": "Willpower, directed force, triumph through control, moving forward with determination",
        "major_corr": {"element": "water", "dynamic": "ascending", "density": "dense", "color": "orange"},
        "minor_corr": {"shape": "triangle", "symmetry": "bilateral", "quality": "force", "modality": "cold_wet"},
    },
    {
        "name": "Strength", "symbol": "♌", "hebrew": "ט",
        "desc": "Inner strength, patience, gentle power, endurance, compassion taming raw force",
        "major_corr": {"element": "fire", "dynamic": "flowing", "density": "moderate", "color": "green"},
        "minor_corr": {"shape": "circle", "symmetry": "bilateral", "quality": "endurance", "modality": "hot_dry"},
    },
    {
        "name": "The Hermit", "symbol": "♍", "hebrew": "י",
        "desc": "Solitude, inner light, withdrawal for wisdom, the lamp in darkness, introspection",
        "major_corr": {"element": "earth", "dynamic": "contracting", "density": "sparse", "color": "green"},
        "minor_corr": {"shape": "point", "symmetry": "asymmetric", "quality": "mystery", "modality": "cold_dry"},
    },
    {
        "name": "Wheel of Fortune", "symbol": "♃", "hebrew": "כ",
        "desc": "Cycles, fortune turning, change of circumstances, karma, the wheel of fate",
        "major_corr": {"element": "fire", "dynamic": "rotating", "density": "moderate", "color": "violet"},
        "minor_corr": {"shape": "circle", "symmetry": "radial", "quality": "transformation", "planet": "jupiter"},
    },
    {
        "name": "Justice", "symbol": "♎", "hebrew": "ל",
        "desc": "Balance, fairness, truth, consequences, weighing evidence, clear-eyed assessment",
        "major_corr": {"element": "air", "dynamic": "stillness", "density": "moderate", "color": "green"},
        "minor_corr": {"shape": "line", "symmetry": "bilateral", "quality": "clarity", "modality": "hot_wet"},
    },
    {
        "name": "Hanged Man", "symbol": "▽", "hebrew": "מ",
        "desc": "Surrender, seeing differently, suspension, sacrifice for wisdom, letting go, new perspective",
        "major_corr": {"element": "water", "dynamic": "suspended", "density": "sparse", "color": "blue"},
        "minor_corr": {"shape": "triangle", "symmetry": "bilateral", "quality": "receptivity", "modality": "cold_wet"},
    },
    {
        "name": "Death", "symbol": "♏", "hebrew": "נ",
        "desc": "Transformation, endings that enable beginnings, clearing away, profound change, rebirth",
        "major_corr": {"element": "water", "dynamic": "dissolving", "density": "moderate", "color": "indigo"},
        "minor_corr": {"shape": "crescent", "symmetry": "asymmetric", "quality": "transformation", "modality": "cold_wet"},
    },
    {
        "name": "Temperance", "symbol": "♐", "hebrew": "ס",
        "desc": "Moderation, blending, patience, the middle way, careful mixing of elements, balance",
        "major_corr": {"element": "fire", "dynamic": "flowing", "density": "moderate", "color": "blue"},
        "minor_corr": {"shape": "triangle", "symmetry": "bilateral", "quality": "balance", "modality": "hot_dry"},
    },
    {
        "name": "The Devil", "symbol": "♑", "hebrew": "ע",
        "desc": "Bondage, materialism, shadow, chains that can be removed, confronting limitation",
        "major_corr": {"element": "earth", "dynamic": "contracting", "density": "saturated", "color": "indigo"},
        "minor_corr": {"shape": "pentagon", "symmetry": "bilateral", "quality": "structure", "planet": "saturn"},
    },
    {
        "name": "The Tower", "symbol": "♂", "hebrew": "פ",
        "desc": "Sudden disruption, structures falling, liberation through destruction, revelation, shock",
        "major_corr": {"element": "fire", "dynamic": "explosive", "density": "dense", "color": "red"},
        "minor_corr": {"shape": "line", "symmetry": "asymmetric", "quality": "liberation", "planet": "mars"},
    },
    {
        "name": "The Star", "symbol": "✦", "hebrew": "צ",
        "desc": "Hope, inspiration, renewal, clarity after storm, pouring forth, guiding light",
        "major_corr": {"element": "air", "dynamic": "radiating", "density": "sparse", "color": "violet"},
        "minor_corr": {"shape": "hexagram", "symmetry": "radial", "quality": "clarity", "modality": "hot_wet"},
    },
    {
        "name": "The Moon", "symbol": "♓", "hebrew": "ק",
        "desc": "Illusion, the unconscious, uncertainty, dreams, what shifts in moonlight, deception",
        "major_corr": {"element": "water", "dynamic": "oscillating", "density": "moderate", "color": "violet"},
        "minor_corr": {"shape": "crescent", "symmetry": "bilateral", "quality": "mystery", "planet": "moon"},
    },
    {
        "name": "The Sun", "symbol": "☉", "hebrew": "ר",
        "desc": "Joy, success, clarity, vitality, everything illuminated, childlike confidence",
        "major_corr": {"element": "fire", "dynamic": "radiating", "density": "dense", "color": "gold"},
        "minor_corr": {"shape": "circle", "symmetry": "radial", "quality": "clarity", "planet": "sun"},
    },
    {
        "name": "Judgement", "symbol": "🜂", "hebrew": "ש",
        "desc": "Awakening, reckoning, answering a call, rising to purpose, final assessment, decisive",
        "major_corr": {"element": "fire", "dynamic": "ascending", "density": "moderate", "color": "red"},
        "minor_corr": {"shape": "triangle", "symmetry": "radial", "quality": "transformation", "modality": "hot_dry"},
    },
    {
        "name": "The World", "symbol": "♄", "hebrew": "ת",
        "desc": "Completion, integration, wholeness, accomplished purpose, fulfillment, the dance",
        "major_corr": {"element": "earth", "dynamic": "rotating", "density": "saturated", "color": "indigo"},
        "minor_corr": {"shape": "circle", "symmetry": "fourfold", "quality": "unity", "planet": "saturn"},
    },
]

# ============================================================
# SHAPE VOCABULARIES
# ============================================================

SHAPES = {
    "point":    {"chars": ["·", "∙", "•", "◦"], "weight": 1},
    "line":     {"chars": ["│", "─", "╱", "╲"], "weight": 2},
    "triangle": {"chars": ["△", "▽", "◁", "▷"], "weight": 3},
    "square":   {"chars": ["□", "■", "▪", "◻"], "weight": 4},
    "pentagon": {"chars": ["⬠", "⬡", "⏣"], "weight": 5},
    "hexagram": {"chars": ["✡", "✦", "⬢", "⬡"], "weight": 6},
    "circle":   {"chars": ["○", "◎", "●", "◉"], "weight": 4},
    "spiral":   {"chars": ["◌", "◍", "⊛", "⊙"], "weight": 3},
    "cross":    {"chars": ["✚", "✛", "╋", "┼"], "weight": 4},
    "crescent": {"chars": ["☽", "☾", "◗", "◖"], "weight": 2},
}

ELEM_SYMBOLS = {"fire": "△", "water": "▽", "air": "◇", "earth": "□"}
DYNAMIC_SYMBOLS = {
    "ascending": "↑", "descending": "↓", "radiating": "✦",
    "concentrating": "◎", "flowing": "≈", "stillness": "·",
    "rotating": "⟳", "oscillating": "⇌", "explosive": "✴",
    "expanding": "◌", "contracting": "◍", "suspended": "∘",
    "dissolving": "∽", "stable": "═",
}


def softmax(x, temperature=1.0):
    """Softmax with temperature. Lower temp = sharper."""
    x = np.array(x) / temperature
    x = x - x.max()
    e = np.exp(x)
    return e / e.sum()


def render_sigil(primary, secondary, tertiary, size=9):
    """Render a sigil from three arcana's correspondences.

    Primary dominates the center and overall structure.
    Secondary modulates the mid-ring.
    Tertiary provides corner/edge accents.
    """
    p_corr = {**primary["major_corr"], **primary["minor_corr"]}
    s_corr = {**secondary["major_corr"], **secondary["minor_corr"]}
    t_corr = {**tertiary["major_corr"], **tertiary["minor_corr"]}

    # Visual properties from primary
    p_color = ANSI.get(p_corr["color"], "")
    s_color = ANSI.get(s_corr["color"], "")
    t_color = ANSI.get(t_corr["color"], "")

    p_shape = SHAPES.get(p_corr["shape"], SHAPES["circle"])
    s_shape = SHAPES.get(s_corr["shape"], SHAPES["circle"])
    t_shape = SHAPES.get(t_corr.get("shape", "point"), SHAPES["point"])

    p_sym = p_shape["chars"][0]
    p_alt = p_shape["chars"][1] if len(p_shape["chars"]) > 1 else p_sym
    s_sym = s_shape["chars"][0]
    t_sym = t_shape["chars"][0]

    p_elem = ELEM_SYMBOLS.get(p_corr["element"], "·")
    s_elem = ELEM_SYMBOLS.get(s_corr["element"], "·")
    t_elem = ELEM_SYMBOLS.get(t_corr.get("element", ""), "·")

    p_dyn = DYNAMIC_SYMBOLS.get(p_corr["dynamic"], "·")
    s_dyn = DYNAMIC_SYMBOLS.get(s_corr["dynamic"], "·")

    density = p_corr["density"]
    symmetry = p_corr.get("symmetry", s_corr.get("symmetry", "bilateral"))

    half = size // 2

    # Build the sigil grid
    grid = [[" " for _ in range(size)] for _ in range(size)]

    for r in range(size):
        for c in range(size):
            dist = math.sqrt((r - half)**2 + (c - half)**2) / (half + 0.5)
            is_center = dist < 0.25
            is_inner = 0.25 <= dist < 0.55
            is_mid = 0.55 <= dist < 0.8
            is_outer = dist >= 0.8

            if is_center:
                # Primary dominates center
                if dist < 0.1:
                    grid[r][c] = ("p_main", p_sym)
                else:
                    grid[r][c] = ("p_dyn", p_dyn)
            elif is_inner:
                # Primary shape + element
                if (r + c) % 2 == 0:
                    grid[r][c] = ("p_elem", p_elem)
                else:
                    grid[r][c] = ("p_shape", p_alt)
            elif is_mid:
                # Secondary arcana
                if (r + c) % 3 == 0:
                    grid[r][c] = ("s_sym", s_sym)
                elif (r + c) % 3 == 1:
                    grid[r][c] = ("s_dyn", s_dyn)
                else:
                    grid[r][c] = ("s_elem", s_elem)
            elif is_outer:
                # Tertiary accents at edges/corners
                if r == 0 or r == size-1 or c == 0 or c == size-1:
                    if (r + c) % 2 == 0:
                        grid[r][c] = ("t_sym", t_sym)
                    else:
                        grid[r][c] = ("t_elem", t_elem)
                else:
                    grid[r][c] = ("bg", " ")

    # Apply symmetry
    if symmetry in ("bilateral", "fourfold", "radial"):
        for r in range(size):
            for c in range(half + 1):
                mc = size - 1 - c
                if mc != c:
                    grid[r][mc] = grid[r][c]
    if symmetry in ("fourfold", "radial"):
        for r in range(half + 1):
            mr = size - 1 - r
            if mr != r:
                for c in range(size):
                    grid[mr][c] = grid[r][c]

    # Render with colors
    lines = []
    for r in range(size):
        line = ""
        for c in range(size):
            tag, ch = grid[r][c] if isinstance(grid[r][c], tuple) else ("bg", grid[r][c])
            if tag.startswith("p_"):
                line += f"{p_color}{ch}{RESET}"
            elif tag.startswith("s_"):
                line += f"{s_color}{ch}{RESET}"
            elif tag.startswith("t_"):
                line += f"{t_color}{ch}{RESET}"
            else:
                line += ch
        lines.append(line)

    return lines


# ============================================================
# TEST INPUTS
# ============================================================

INPUTS = [
    ("Artifex (Alignment)",
     "The primary witness. The athanor exists to reduce his executive function load. Not holding the system in his head. Athanors that behave as expected without correction. Trails he can read without effort. The system improving from evidence."),
    ("Artifex (Harvest)",
     "He feels secure — not anxious that the productivity delta will be discovered. He feels compensated — the value is reciprocated. He feels protected. He feels portable — the system is his IP, his career asset."),
    ("Artifex (Anvil)",
     "The artifex sits down to work and finds a capable, oriented partner already there — genuinely oriented to the codebase. EF overhead shifts to the azer. The artifex directs judgment, not logistics."),
    ("Kelly (Analytics)",
     "Kelly PM — direct end user; her Omni sniff test is the final confidence layer; notified in team channel each time a workstream's events are live."),
    ("Horse Owners",
     "Horse owners see a demo of their horse's care history and have a genuine emotional reaction — excitement, relief, I wish my stable did this."),
    ("Calcinatio: Audit",
     "Behavioral audit is the highest-value fire. Reading trails for patterns — calcinatio skipped, geas satisfied superficially — reveals misalignment no automated check can find."),
    ("Calcinatio: Customer",
     "Customer outcome fidelity is the gate. Seismic API: case-sensitive exact string matching, atomic rejection. Broken mappings equal customer outcome problem."),
    ("Intent: Alignment",
     "The athanor system's actual behavior progressively aligns with the core intents. Agents apply calcinatio without being told. Trails are legible. A maintained condition."),
    ("Intent: Horse Demo",
     "Horse owners see a demo and have a genuine emotional reaction — excitement, relief. The reaction is unambiguous: either they light up or they don't."),
    ("Tempering: Stabilize",
     "Stabilization before release. Feature is functionally complete. Bar for new opera is high — only customer-facing risks."),
]


def main():
    arcana_texts = [m["desc"] for m in MAJORS]
    input_texts = [t for _, t in INPUTS]

    print("Embedding 22 major arcana...")
    arcana_emb = embed(arcana_texts)
    input_emb = embed(input_texts)

    sims = cosine_similarity(input_emb, arcana_emb)

    print(f"\n{'='*72}")
    print(f"  SHARP SVSI — 22 Major Arcana, Temperature=0.5")
    print(f"{'='*72}")

    profiles = []

    for i, (label, text) in enumerate(INPUTS):
        scores = sims[i]
        # Sharp selection via softmax with low temperature
        weights = softmax(scores, temperature=0.3)
        profiles.append(weights)

        top_idx = np.argsort(scores)[::-1]
        top3 = top_idx[:3]

        p = MAJORS[top3[0]]
        s = MAJORS[top3[1]]
        t = MAJORS[top3[2]]

        print(f"\n{'─'*72}")
        print(f"  {label}")
        print(f"  ① {p['symbol']} {p['name']:<18} ({scores[top3[0]]:.2f})"
              f"  ② {s['symbol']} {s['name']:<18} ({scores[top3[1]]:.2f})"
              f"  ③ {t['symbol']} {t['name']:<18} ({scores[top3[2]]:.2f})")

        # Show correspondences
        pc = p["major_corr"]
        print(f"  Correspondences: {pc['element']}, {pc['dynamic']}, {pc['density']}, {pc['color']}"
              f" | shape:{p['minor_corr']['shape']}, {p['minor_corr']['quality']}")

        sigil = render_sigil(p, s, t, size=9)
        print()
        for line in sigil:
            print(f"    {line}")

    # ── Similarity comparison ──
    print(f"\n{'='*72}")
    print(f"  ARCANA-PROFILE SIMILARITY (22-dim, sharp weights)")
    print(f"{'='*72}")

    prof_arr = np.array(profiles)
    psim = cosine_similarity(prof_arr)

    pairs = [
        (0, 1, "Artifex/Alignment vs Artifex/Harvest"),
        (0, 2, "Artifex/Alignment vs Artifex/Anvil"),
        (1, 2, "Artifex/Harvest vs Artifex/Anvil"),
        (0, 3, "Artifex vs Kelly"),
        (0, 4, "Artifex vs Horse Owners"),
        (0, 5, "Artifex vs Calcinatio:Audit"),
        (4, 8, "Horse Owners witness vs Horse Intent"),
        (5, 6, "Calcinatio:Audit vs Calcinatio:Customer"),
        (7, 9, "Intent:Alignment vs Tempering:Stabilize"),
        (3, 6, "Kelly vs Calcinatio:Customer"),
    ]

    print(f"\n  {'Pair':<48} {'Sim':>6}")
    print(f"  {'─'*60}")
    for i1, i2, lbl in pairs:
        s = psim[i1, i2]
        bar = '█' * int(s * 20) + '░' * (20 - int(s * 20))
        print(f"  {lbl:<48} {s:.3f} {bar}")

    print()


if __name__ == '__main__':
    main()
