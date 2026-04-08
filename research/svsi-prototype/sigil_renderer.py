#!/usr/bin/env python3
"""
Stroke-based sigil renderer.

Real sigils are paths, arms, and nodes — not filled shapes.
Correspondences determine: number of arms, angles, terminal symbols,
center node, junction types, stroke style.
"""
import sys, os, math
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity
from sentence_transformers import SentenceTransformer

model = SentenceTransformer('all-MiniLM-L6-v2')
def embed(texts): return model.encode(texts)

sys.path.insert(0, os.path.dirname(__file__))
from sharp_renderer import MAJORS, ANSI, softmax

RESET = "\033[0m"

ELEM_TIPS = {"fire": "△", "water": "▽", "air": "◇", "earth": "□"}
ELEM_NODES = {"fire": "▲", "water": "▼", "air": "◆", "earth": "■"}

# Stroke characters for different angles
STROKE = {
    "h": "─",      # horizontal
    "v": "│",      # vertical
    "d1": "╲",     # diagonal \
    "d2": "╱",     # diagonal /
    "cross": "╳",  # crossing
    "dot": "·",    # faint
}

# Center nodes by arcanum quality
CENTER_NODES = {
    "balance": "⊕", "force": "✦", "receptivity": "◎", "transformation": "✶",
    "structure": "⊞", "liberation": "✴", "unity": "◉", "duality": "⊗",
    "mystery": "⊛", "clarity": "☉", "endurance": "◈", "swiftness": "⚡",
}

# Junction nodes
JUNCTIONS = ["○", "●", "◦", "◎", "⊙", "◌", "◍"]


def draw_line(grid, r0, c0, r1, c1, char=None):
    """Draw a line on the grid using Bresenham-ish approach with Unicode."""
    dr = r1 - r0
    dc = c1 - c0
    steps = max(abs(dr), abs(dc))
    if steps == 0:
        return

    # Choose stroke character based on angle
    if char is None:
        if dr == 0:
            char = STROKE["h"]
        elif dc == 0:
            char = STROKE["v"]
        elif (dr > 0 and dc > 0) or (dr < 0 and dc < 0):
            char = STROKE["d1"]
        else:
            char = STROKE["d2"]

    for i in range(steps + 1):
        t = i / steps if steps > 0 else 0
        r = int(round(r0 + dr * t))
        c = int(round(c0 + dc * t))
        if 0 <= r < len(grid) and 0 <= c < len(grid[0]):
            if grid[r][c] == " ":
                grid[r][c] = char


def render_sigil(primary, secondary, tertiary, size=15):
    """Render a stroke-based sigil."""
    p = {**primary["major_corr"], **primary["minor_corr"]}
    s = {**secondary["major_corr"], **secondary["minor_corr"]}
    t = {**tertiary["major_corr"], **tertiary["minor_corr"]}

    grid = [[" "] * size for _ in range(size)]
    half = size // 2
    cx, cy = half, half  # center

    # ── Determine arm structure from shape correspondence ──

    shape = p["shape"]
    element = p["element"]
    density = p["density"]
    dynamic = p["dynamic"]
    symmetry = p.get("symmetry", "bilateral")
    quality = p.get("quality", "balance")

    # Arm reach based on density
    reach = {"sparse": 3, "moderate": 4, "dense": 5, "saturated": 6}.get(density, 4)

    # Number of primary arms and their angles
    if shape == "triangle":
        if element in ("fire", "air"):
            # Upward pointing: arms at -90, -210, -330 degrees
            angles = [-90, -210, -330]
        else:
            # Downward: arms at 90, 210, 330
            angles = [90, 210, 330]
    elif shape == "square":
        angles = [0, 90, 180, 270]
    elif shape == "hexagram":
        angles = [0, 60, 120, 180, 240, 300]
    elif shape == "cross":
        angles = [0, 90, 180, 270]
    elif shape == "circle":
        # Ring — arms point outward at 8 positions
        angles = [0, 45, 90, 135, 180, 225, 270, 315]
        reach = max(2, reach - 1)  # shorter arms, more of them
    elif shape == "crescent":
        # Arc — arms curve to one side
        angles = [-45, -90, -135]
    elif shape == "line":
        angles = [-90, 90]  # just up and down
        reach = reach + 1
    elif shape == "pentagon":
        angles = [90, 162, 234, 306, 18]
    elif shape == "spiral":
        # Spiral arms
        angles = [0, 120, 240]
    elif shape == "point":
        angles = []  # no arms, just center
        reach = 0
    else:
        angles = [0, 90, 180, 270]

    # ── Draw primary arms ──

    tip_sym = ELEM_TIPS.get(element, "·")
    arm_endpoints = []

    for angle_deg in angles:
        angle = math.radians(angle_deg)
        # End point
        er = int(round(cy + reach * math.sin(angle)))
        ec = int(round(cx + reach * math.cos(angle)))
        er = max(0, min(size-1, er))
        ec = max(0, min(size-1, ec))

        draw_line(grid, cy, cx, er, ec)
        arm_endpoints.append((er, ec))

        # Place tip symbol
        if 0 <= er < size and 0 <= ec < size:
            grid[er][ec] = tip_sym

    # ── Center node ──
    center_sym = CENTER_NODES.get(quality, "◉")
    grid[cy][cx] = center_sym

    # ── Secondary structure: cross-bars or rings connecting arms ──

    s_shape = s.get("shape", "line")
    s_elem = s.get("element", "air")
    s_dynamic = s.get("dynamic", "stillness")
    s_sym = ELEM_TIPS.get(s_elem, "·")

    if len(arm_endpoints) >= 2:
        if s_shape in ("hexagram", "circle", "cross"):
            # Connect adjacent arm endpoints
            for i in range(len(arm_endpoints)):
                j = (i + 1) % len(arm_endpoints)
                r0, c0 = arm_endpoints[i]
                r1, c1 = arm_endpoints[j]
                # Draw connecting stroke at 2/3 distance from center
                mr0 = int(cy + (r0 - cy) * 0.6)
                mc0 = int(cx + (c0 - cx) * 0.6)
                mr1 = int(cy + (r1 - cy) * 0.6)
                mc1 = int(cx + (c1 - cx) * 0.6)
                draw_line(grid, mr0, mc0, mr1, mc1, STROKE["dot"])
        elif s_shape in ("triangle", "square"):
            # Place junction nodes at midpoints of arms
            for er, ec in arm_endpoints:
                mr = (cy + er) // 2
                mc = (cx + ec) // 2
                if 0 <= mr < size and 0 <= mc < size and grid[mr][mc] in (STROKE["h"], STROKE["v"], STROKE["d1"], STROKE["d2"], " "):
                    grid[mr][mc] = "○"
        elif s_shape == "crescent":
            # Arc between first and last arm
            if len(arm_endpoints) >= 2:
                r0, c0 = arm_endpoints[0]
                r1, c1 = arm_endpoints[-1]
                mr = (r0 + r1) // 2
                mc = (c0 + c1) // 2
                if 0 <= mr < size and 0 <= mc < size:
                    grid[mr][mc] = "◗"
        else:
            # Default: place secondary element at arm midpoints
            for er, ec in arm_endpoints[::2]:
                mr = (cy + er) // 2
                mc = (cx + ec) // 2
                if 0 <= mr < size and 0 <= mc < size and grid[mr][mc] != center_sym:
                    grid[mr][mc] = s_sym

    # ── Tertiary accents: small marks near tips ──

    t_elem = t.get("element", "air")
    t_tip = ELEM_TIPS.get(t_elem, "·")
    t_shape = t.get("shape", "point")

    for er, ec in arm_endpoints:
        # Place tertiary accent offset from tip
        for dr, dc in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
            nr, nc = er + dr, ec + dc
            if 0 <= nr < size and 0 <= nc < size and grid[nr][nc] == " ":
                grid[nr][nc] = t_tip
                break  # just one accent per tip

    # ── Dynamic modifier: add motion marks ──

    if dynamic == "radiating":
        # Small dots around the perimeter
        for angle_deg in range(0, 360, 30):
            if angle_deg not in [int(a) % 360 for a in angles]:
                angle = math.radians(angle_deg)
                r = int(round(cy + (reach + 1) * math.sin(angle)))
                c = int(round(cx + (reach + 1) * math.cos(angle)))
                if 0 <= r < size and 0 <= c < size and grid[r][c] == " ":
                    grid[r][c] = "·"
    elif dynamic == "rotating":
        # Curved accent marks
        for er, ec in arm_endpoints:
            # Offset perpendicular to arm
            dr = er - cy
            dc = ec - cx
            pr, pc = -dc, dr  # perpendicular
            norm = math.sqrt(pr*pr + pc*pc) + 0.001
            pr, pc = pr/norm, pc/norm
            nr = int(round(er + pr))
            nc = int(round(ec + pc))
            if 0 <= nr < size and 0 <= nc < size and grid[nr][nc] == " ":
                grid[nr][nc] = "⟳" if (er + ec) % 2 == 0 else "·"
    elif dynamic == "ascending":
        # Upward tick marks along vertical arm if present
        for r in range(max(0, cy - reach), cy):
            if grid[r][cx] == " ":
                grid[r][cx] = "·"
    elif dynamic == "explosive":
        # Scatter marks
        import random
        rng = random.Random(hash(primary["name"]))
        for _ in range(4):
            r = rng.randint(0, size-1)
            c = rng.randint(0, size-1)
            if grid[r][c] == " ":
                grid[r][c] = "∗"
    elif dynamic == "flowing":
        # Wavy marks between arms
        for i in range(len(arm_endpoints)):
            j = (i + 1) % len(arm_endpoints)
            r0, c0 = arm_endpoints[i]
            r1, c1 = arm_endpoints[j]
            mr = (r0 + r1) // 2
            mc = (c0 + c1) // 2
            if 0 <= mr < size and 0 <= mc < size and grid[mr][mc] == " ":
                grid[mr][mc] = "≈"
    elif dynamic == "oscillating":
        # Double-headed marks
        for er, ec in arm_endpoints:
            dr = er - cy
            dc = ec - cx
            nr = er - (1 if dr > 0 else -1 if dr < 0 else 0)
            nc = ec - (1 if dc > 0 else -1 if dc < 0 else 0)
            if 0 <= nr < size and 0 <= nc < size and grid[nr][nc] == " ":
                grid[nr][nc] = "⇌"

    # ── Apply symmetry as post-process (only if bilateral) ──
    if symmetry == "bilateral":
        for r in range(size):
            for c in range(half):
                mc = size - 1 - c
                if grid[r][c] != " " and grid[r][mc] == " ":
                    ch = grid[r][c]
                    # Mirror stroke characters
                    if ch == STROKE["d1"]:
                        ch = STROKE["d2"]
                    elif ch == STROKE["d2"]:
                        ch = STROKE["d1"]
                    grid[r][mc] = ch
                elif grid[r][mc] != " " and grid[r][c] == " ":
                    ch = grid[r][mc]
                    if ch == STROKE["d1"]:
                        ch = STROKE["d2"]
                    elif ch == STROKE["d2"]:
                        ch = STROKE["d1"]
                    grid[r][c] = ch

    # ── Apply color ──
    p_color = ANSI.get(p["color"], "")
    s_color = ANSI.get(s.get("color", "white"), "")

    lines = []
    for row in grid:
        line = ""
        for ch in row:
            if ch in (STROKE["h"], STROKE["v"], STROKE["d1"], STROKE["d2"], STROKE["cross"], STROKE["dot"]):
                line += f"{p_color}{ch}{RESET}"
            elif ch in ELEM_TIPS.values() or ch in ELEM_NODES.values():
                line += f"{s_color}{ch}{RESET}"
            elif ch in CENTER_NODES.values() or ch in JUNCTIONS:
                line += f"{p_color}{ch}{RESET}"
            elif ch in ("·", "≈", "⟳", "∗", "⇌", "◗"):
                line += f"{s_color}{ch}{RESET}"
            else:
                line += ch
        lines.append(line)

    return lines


INPUTS = [
    ("Artifex (Alignment)",
     "The primary witness. The athanor exists to reduce his executive function load. Not holding the system in his head. Athanors that behave as expected. Trails he can read without effort."),
    ("Artifex (Harvest)",
     "He feels secure, compensated, protected, portable. The value is reciprocated. The system is his IP, his career asset."),
    ("Artifex (Anvil)",
     "The artifex sits down and finds a capable partner oriented to the codebase. EF overhead shifts to the azer. Directs judgment, not logistics."),
    ("Kelly (Analytics)",
     "Kelly PM — direct end user; her Omni sniff test is the final confidence layer; notified when events are live."),
    ("Horse Owners",
     "Horse owners see a demo of their horse's care history and have a genuine emotional reaction — excitement, relief."),
    ("Calcinatio: Audit",
     "Behavioral audit is the highest-value fire. Reading trails for patterns reveals misalignment no automated check can find."),
    ("Calcinatio: Customer",
     "Customer outcome fidelity is the gate. Case-sensitive exact string matching. Broken mappings equal customer problem."),
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

    print(f"\n{'='*60}")
    print(f"  STROKE-BASED SIGIL RENDERER")
    print(f"  Paths, arms, and nodes — not filled shapes")
    print(f"{'='*60}")

    for i, (label, text) in enumerate(INPUTS):
        scores = sims[i]
        top_idx = np.argsort(scores)[::-1][:3]

        p = MAJORS[top_idx[0]]
        s = MAJORS[top_idx[1]]
        t = MAJORS[top_idx[2]]

        pc = {**p["major_corr"], **p["minor_corr"]}

        print(f"\n{'─'*60}")
        print(f"  {label}")
        print(f"  ① {p['name']:<16} ② {s['name']:<16} ③ {t['name']}")
        print(f"  {pc['shape']} / {pc['element']} / {pc['dynamic']} / {pc.get('symmetry','bilateral')}")

        sigil = render_sigil(p, s, t, size=15)
        print()
        for line in sigil:
            # Only print if line has non-space content
            stripped = line.replace(RESET, "").replace(" ", "")
            for code in ANSI.values():
                stripped = stripped.replace(code, "")
            if stripped:
                print(f"    {line}")

    print()


if __name__ == '__main__':
    main()
