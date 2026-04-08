#!/usr/bin/env python3
"""
Fluid SVSI renderer — shape mask determines overall silhouette,
not just fill pattern. The bounding shape varies by correspondences.
"""
import sys, os, math
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity
from sentence_transformers import SentenceTransformer

model = SentenceTransformer('all-MiniLM-L6-v2')
def embed(texts): return model.encode(texts)

RESET = "\033[0m"

# Reuse MAJORS from sharp_renderer
sys.path.insert(0, os.path.dirname(__file__))
from sharp_renderer import MAJORS, SHAPES, ELEM_SYMBOLS, DYNAMIC_SYMBOLS, ANSI, softmax


def shape_mask(shape, size, element, dynamic):
    """Return a boolean mask determining which cells are 'inside' the sigil.
    The shape correspondence determines the silhouette."""
    half = size // 2
    mask = [[False]*size for _ in range(size)]

    if shape == "circle":
        for r in range(size):
            for c in range(size):
                dist = math.sqrt((r-half)**2 + (c-half)**2)
                mask[r][c] = dist <= half + 0.3
    elif shape == "triangle":
        # Upward if fire/ascending, downward if water/descending
        up = element in ("fire", "air") or dynamic in ("ascending", "radiating", "explosive")
        for r in range(size):
            row = r if up else (size - 1 - r)
            width = int((size - row) * size / (size + 1)) // 2
            center = half
            for c in range(size):
                mask[r][c] = abs(c - center) <= width
    elif shape == "square":
        # Diamond (rotated 45°) feels less boxy
        for r in range(size):
            for c in range(size):
                mask[r][c] = abs(r - half) + abs(c - half) <= half + 0.5
    elif shape == "hexagram":
        # Six-pointed star: two overlapping triangles
        for r in range(size):
            for c in range(size):
                # Triangle up
                row_up = r
                width_up = int((size - row_up) * size / (size + 1)) // 2
                in_up = abs(c - half) <= width_up
                # Triangle down
                row_dn = size - 1 - r
                width_dn = int((size - row_dn) * size / (size + 1)) // 2
                in_dn = abs(c - half) <= width_dn
                mask[r][c] = in_up or in_dn
    elif shape == "cross":
        arm = max(1, half // 2)
        for r in range(size):
            for c in range(size):
                mask[r][c] = abs(r - half) <= arm or abs(c - half) <= arm
    elif shape == "crescent":
        for r in range(size):
            for c in range(size):
                dist1 = math.sqrt((r-half)**2 + (c-half)**2)
                dist2 = math.sqrt((r-half)**2 + (c-half-2)**2)
                mask[r][c] = dist1 <= half + 0.3 and dist2 > half - 1
    elif shape == "line":
        # Tall narrow
        width = max(1, half // 2)
        for r in range(size):
            for c in range(size):
                mask[r][c] = abs(c - half) <= width
    elif shape == "pentagon":
        for r in range(size):
            for c in range(size):
                angle = math.atan2(r - half, c - half)
                max_r = half * (1 + 0.15 * math.cos(5 * angle))
                dist = math.sqrt((r-half)**2 + (c-half)**2)
                mask[r][c] = dist <= max_r
    elif shape == "spiral":
        for r in range(size):
            for c in range(size):
                dist = math.sqrt((r-half)**2 + (c-half)**2)
                angle = math.atan2(r-half, c-half)
                spiral_r = half * (0.3 + 0.7 * ((angle + math.pi) / (2 * math.pi)))
                mask[r][c] = abs(dist - spiral_r) < 1.5 or dist < 1.5
    elif shape == "point":
        for r in range(size):
            for c in range(size):
                dist = math.sqrt((r-half)**2 + (c-half)**2)
                mask[r][c] = dist <= half * 0.5
    else:
        # Default: circle
        for r in range(size):
            for c in range(size):
                dist = math.sqrt((r-half)**2 + (c-half)**2)
                mask[r][c] = dist <= half + 0.3

    return mask


def render_fluid(primary, secondary, tertiary, size=11):
    """Render with fluid shape mask."""
    p_corr = {**primary["major_corr"], **primary["minor_corr"]}
    s_corr = {**secondary["major_corr"], **secondary["minor_corr"]}
    t_corr = {**tertiary["major_corr"], **tertiary["minor_corr"]}

    # Get shape mask from primary
    mask = shape_mask(p_corr["shape"], size, p_corr["element"], p_corr["dynamic"])

    # Symmetry
    symmetry = p_corr.get("symmetry", "bilateral")
    half = size // 2

    # Apply symmetry to mask
    if symmetry in ("bilateral", "fourfold", "radial"):
        for r in range(size):
            for c in range(half + 1):
                mc = size - 1 - c
                mask[r][mc] = mask[r][mc] or mask[r][c]
                mask[r][c] = mask[r][mc]
    if symmetry in ("fourfold", "radial"):
        for r in range(half + 1):
            mr = size - 1 - r
            for c in range(size):
                mask[mr][c] = mask[mr][c] or mask[r][c]
                mask[r][c] = mask[mr][c]

    # Visual elements
    p_shape_chars = SHAPES.get(p_corr["shape"], SHAPES["circle"])["chars"]
    s_shape_chars = SHAPES.get(s_corr["shape"], SHAPES["circle"])["chars"]

    p_elem = ELEM_SYMBOLS.get(p_corr["element"], "·")
    s_elem = ELEM_SYMBOLS.get(s_corr["element"], "·")
    t_elem = ELEM_SYMBOLS.get(t_corr.get("element", ""), "·")

    p_dyn = DYNAMIC_SYMBOLS.get(p_corr["dynamic"], "·")
    s_dyn = DYNAMIC_SYMBOLS.get(s_corr["dynamic"], "·")

    p_sym = p_shape_chars[0]
    p_alt = p_shape_chars[1] if len(p_shape_chars) > 1 else p_sym
    s_sym = s_shape_chars[0]
    t_sym = SHAPES.get(t_corr.get("shape", "point"), SHAPES["point"])["chars"][0]

    density = p_corr["density"]
    dense_val = {"sparse": 0.3, "moderate": 0.5, "dense": 0.7, "saturated": 0.9}.get(density, 0.5)

    p_color = ANSI.get(p_corr["color"], "")
    s_color = ANSI.get(s_corr["color"], "")
    t_color = ANSI.get(t_corr.get("color", "white"), "")

    # Build grid
    grid = []
    for r in range(size):
        row = ""
        for c in range(size):
            if not mask[r][c]:
                row += " "
                continue

            dist = math.sqrt((r-half)**2 + (c-half)**2) / (half + 0.5)

            # Edge detection: is this cell on the boundary of the mask?
            is_edge = False
            for dr, dc in [(-1,0),(1,0),(0,-1),(0,1)]:
                nr, nc = r+dr, c+dc
                if 0 <= nr < size and 0 <= nc < size:
                    if not mask[nr][nc]:
                        is_edge = True
                        break
                elif True:  # out of bounds = outside mask
                    is_edge = True
                    break

            if is_edge:
                # Border: tertiary arcana
                row += f"{t_color}{t_elem}{RESET}"
            elif dist < 0.25:
                # Core: primary symbol
                row += f"{p_color}{p_sym}{RESET}"
            elif dist < 0.45:
                # Inner: primary dynamic + element
                if (r + c) % 2 == 0:
                    row += f"{p_color}{p_dyn}{RESET}"
                else:
                    row += f"{p_color}{p_elem}{RESET}"
            elif dist < 0.65:
                # Mid: secondary arcana
                if (r + c) % 3 == 0:
                    row += f"{s_color}{s_sym}{RESET}"
                elif (r + c) % 3 == 1:
                    row += f"{s_color}{s_dyn}{RESET}"
                else:
                    row += f"{s_color}{s_elem}{RESET}"
            else:
                # Outer fill
                if dense_val > 0.6:
                    row += f"{t_color}{p_alt}{RESET}" if (r+c) % 2 == 0 else f"{t_color}{s_elem}{RESET}"
                elif dense_val > 0.4:
                    row += f"{t_color}{t_sym}{RESET}" if (r+c) % 3 == 0 else " "
                else:
                    row += " " if (r+c) % 2 == 0 else f"{t_color}·{RESET}"

        grid.append(row)

    return grid


INPUTS = [
    ("Artifex (Alignment)",
     "The primary witness. The athanor exists to reduce his executive function load. Not holding the system in his head. Athanors that behave as expected. Trails he can read without effort."),
    ("Artifex (Harvest)",
     "He feels secure, compensated, protected, portable. The value is reciprocated. The system is his IP, his career asset."),
    ("Artifex (Anvil)",
     "The artifex sits down and finds a capable partner oriented to the codebase. EF overhead shifts to the azer. The artifex directs judgment, not logistics."),
    ("Kelly (Analytics)",
     "Kelly PM — direct end user; her Omni sniff test is the final confidence layer; notified when events are live."),
    ("Horse Owners",
     "Horse owners see a demo of their horse's care history and have a genuine emotional reaction — excitement, relief."),
    ("Calcinatio: Audit",
     "Behavioral audit is the highest-value fire. Reading trails for patterns reveals misalignment no automated check can find."),
    ("Calcinatio: Customer",
     "Customer outcome fidelity is the gate. Case-sensitive exact string matching, atomic rejection. Broken mappings equal customer problem."),
    ("Intent: Alignment",
     "The athanor system progressively aligns with spec. Agents apply calcinatio without being told. Trails are legible. Maintained condition."),
    ("Intent: Horse Demo",
     "Horse owners see a demo and have a genuine emotional reaction. Either they light up or they don't."),
    ("Tempering: Stabilize",
     "Stabilization before release. Feature functionally complete. Bar for new opera is high — only customer-facing risks."),
]


def main():
    arcana_texts = [m["desc"] for m in MAJORS]
    input_texts = [t for _, t in INPUTS]

    arcana_emb = embed(arcana_texts)
    input_emb = embed(input_texts)
    sims = cosine_similarity(input_emb, arcana_emb)

    print(f"\n{'='*72}")
    print(f"  FLUID SVSI — Shape masks from correspondences")
    print(f"{'='*72}")

    for i, (label, text) in enumerate(INPUTS):
        scores = sims[i]
        top_idx = np.argsort(scores)[::-1][:3]

        p = MAJORS[top_idx[0]]
        s = MAJORS[top_idx[1]]
        t = MAJORS[top_idx[2]]

        pc = {**p["major_corr"], **p["minor_corr"]}

        print(f"\n{'─'*72}")
        print(f"  {label}")
        print(f"  ① {p['name']:<18}  ② {s['name']:<18}  ③ {t['name']}")
        print(f"  shape:{pc['shape']}  sym:{pc.get('symmetry','bilateral')}"
              f"  elem:{pc['element']}  dyn:{pc['dynamic']}  density:{pc['density']}")

        sigil = render_fluid(p, s, t, size=11)
        print()
        for line in sigil:
            print(f"    {line}")

    print()


if __name__ == '__main__':
    main()
