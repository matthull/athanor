#!/usr/bin/env python3
"""
End-to-end SVSI demo: text → embedding → arcana match → visual render

Quick and dirty — just show the pipeline working on real inputs.
"""
import sys, os, math
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

# Try sentence transformers, fall back to TF-IDF
try:
    from sentence_transformers import SentenceTransformer
    print("Using sentence-transformers (all-MiniLM-L6-v2)")
    model = SentenceTransformer('all-MiniLM-L6-v2')
    def embed(texts):
        return model.encode(texts)
    USE_SBERT = True
except Exception as e:
    print(f"sentence-transformers failed ({e}), falling back to TF-IDF")
    from sklearn.feature_extraction.text import TfidfVectorizer
    USE_SBERT = False

# ============================================================
# Minimal arcana codebook — just enough for demo
# ============================================================
ARCANA = [
    {"name": "The Fool", "symbol": "○", "hebrew": "א", "color": "\033[93m",
     "shape": "spiral", "quality": "ascending scattered light beginning",
     "desc": "New beginnings, leaping into the unknown, pure potential, innocence, spontaneity"},
    {"name": "The Magician", "symbol": "☿", "hebrew": "ב", "color": "\033[93m",
     "shape": "vertical", "quality": "focused channeling precise skill",
     "desc": "Skill, will made manifest, tools mastered, focused intent, communication of power"},
    {"name": "High Priestess", "symbol": "☽", "hebrew": "ג", "color": "\033[94m",
     "shape": "crescent", "quality": "hidden reflective deep intuitive",
     "desc": "Hidden knowledge, intuition, mystery, the veil between seen and unseen, inner wisdom"},
    {"name": "The Empress", "symbol": "♀", "hebrew": "ד", "color": "\033[92m",
     "shape": "circle", "quality": "fertile expanding abundant nurturing",
     "desc": "Abundance, fertility, nurturing, creative growth, natural expansion, generative force"},
    {"name": "The Emperor", "symbol": "♈", "hebrew": "ה", "color": "\033[91m",
     "shape": "square", "quality": "structured commanding rigid authority",
     "desc": "Authority, structure, control, stability, order imposed on chaos, the father"},
    {"name": "Hierophant", "symbol": "♉", "hebrew": "ו", "color": "\033[91m",
     "shape": "cross", "quality": "teaching bridging institutional tradition",
     "desc": "Teaching, tradition, institutional wisdom, the bridge between divine and mundane"},
    {"name": "The Lovers", "symbol": "♊", "hebrew": "ז", "color": "\033[33m",
     "shape": "twin", "quality": "choosing connecting dual relationship",
     "desc": "Choice, relationship, duality resolved, union of opposites, alignment of values"},
    {"name": "The Chariot", "symbol": "♋", "hebrew": "ח", "color": "\033[33m",
     "shape": "chevron", "quality": "directed willpower moving victory",
     "desc": "Willpower, directed force, triumph through control, moving forward with determination"},
    {"name": "Strength", "symbol": "♌", "hebrew": "ט", "color": "\033[32m",
     "shape": "lemniscate", "quality": "enduring gentle sustained patience",
     "desc": "Inner strength, patience, gentle power, endurance, compassion taming raw force"},
    {"name": "The Hermit", "symbol": "♍", "hebrew": "י", "color": "\033[32m",
     "shape": "point", "quality": "solitary inward illuminating withdrawn",
     "desc": "Solitude, inner light, withdrawal for wisdom, the lamp in darkness, introspection"},
    {"name": "Wheel", "symbol": "♃", "hebrew": "כ", "color": "\033[35m",
     "shape": "wheel", "quality": "cyclical turning fortune change",
     "desc": "Cycles, fortune turning, change of circumstances, the wheel of fate, karma"},
    {"name": "Justice", "symbol": "♎", "hebrew": "ל", "color": "\033[32m",
     "shape": "balance", "quality": "equilibrium weighing fair truth",
     "desc": "Balance, fairness, truth, consequences, weighing evidence, clear-eyed assessment"},
    {"name": "Hanged Man", "symbol": "▽", "hebrew": "מ", "color": "\033[94m",
     "shape": "inverted", "quality": "suspended reversed sacrifice letting-go",
     "desc": "Surrender, seeing differently, suspension, sacrifice for wisdom, letting go"},
    {"name": "Death", "symbol": "♏", "hebrew": "נ", "color": "\033[36m",
     "shape": "scythe", "quality": "transforming ending clearing rebirth",
     "desc": "Transformation, endings that enable beginnings, clearing away, profound change"},
    {"name": "Temperance", "symbol": "♐", "hebrew": "ס", "color": "\033[94m",
     "shape": "hourglass", "quality": "blending moderate flowing patience",
     "desc": "Moderation, blending, patience, the middle way, careful mixing of elements"},
    {"name": "The Devil", "symbol": "♑", "hebrew": "ע", "color": "\033[34m",
     "shape": "pentagon", "quality": "binding material constraint shadow",
     "desc": "Bondage, materialism, shadow, chains that can be removed, confronting limitation"},
    {"name": "The Tower", "symbol": "♂", "hebrew": "פ", "color": "\033[91m",
     "shape": "lightning", "quality": "shattering sudden breaking liberation",
     "desc": "Sudden disruption, structures falling, liberation through destruction, revelation"},
    {"name": "The Star", "symbol": "✦", "hebrew": "צ", "color": "\033[35m",
     "shape": "star", "quality": "radiating hopeful clear renewal",
     "desc": "Hope, inspiration, renewal, clarity after storm, pouring forth, guiding light"},
    {"name": "The Moon", "symbol": "♓", "hebrew": "ק", "color": "\033[35m",
     "shape": "crescents", "quality": "dreamlike shifting uncertain illusion",
     "desc": "Illusion, the unconscious, uncertainty, dreams, what shifts in moonlight"},
    {"name": "The Sun", "symbol": "☉", "hebrew": "ר", "color": "\033[33m",
     "shape": "radiant", "quality": "bright clear joyful success",
     "desc": "Joy, success, clarity, vitality, everything illuminated, childlike confidence"},
    {"name": "Judgement", "symbol": "🜂", "hebrew": "ש", "color": "\033[91m",
     "shape": "flames", "quality": "awakening rising decisive calling",
     "desc": "Awakening, reckoning, answering a call, rising to purpose, final assessment"},
    {"name": "The World", "symbol": "♄", "hebrew": "ת", "color": "\033[34m",
     "shape": "mandorla", "quality": "complete integrated whole accomplished",
     "desc": "Completion, integration, wholeness, the dance of accomplished purpose, fulfillment"},
]

RESET = "\033[0m"

# ============================================================
# Test inputs — a few real MO sections
# ============================================================
INPUTS = [
    ("Artifex witness (Alignment)",
     "The primary witness. The athanor exists to reduce his executive function load. Not holding the system in his head. Athanors that behave as expected without correction. Trails he can read without effort. The system improving from evidence, not from his intervention. His absence-of-concern is the target state."),

    ("Artifex witness (Harvest)",
     "He feels secure — not anxious that the productivity delta will be discovered awkwardly. He feels compensated — the value is visibly reciprocated. He feels protected — the gains don't isolate him. He feels portable — the system is his IP, his career asset."),

    ("Kelly witness (Analytics)",
     "Kelly PM — direct end user; her Omni sniff test is the final confidence layer; notified in team channel each time a workstream's events are live."),

    ("Horse Owners witness",
     "Horse owners see a demo of their horse's care history and have a genuine emotional reaction — excitement, relief, I wish my stable did this. The demo is compelling enough that the reaction is unambiguous."),

    ("Calcinatio: Behavioral audit",
     "Behavioral audit is the highest-value fire. Reading trails for patterns — calcinatio skipped, geas satisfied superficially, escalations absent, opera procedural rather than witness-oriented — reveals misalignment no automated check can find."),

    ("Calcinatio: Customer outcome",
     "Customer outcome fidelity is the gate. Seismic API: case-sensitive exact string matching, atomic rejection, silent cardinality violations. Broken mappings equal customer outcome problem, not UI problem."),

    ("Intent: Alignment",
     "The athanor system's actual behavior progressively aligns with the core intents in the spec. The artifex feels EF load reducing. Agents apply calcinatio without being told. Trails are legible. This is a maintained condition, not a build goal."),

    ("Intent: Horse Owner Demo",
     "Horse owners see a demo and have a genuine emotional reaction — excitement, relief. The demo is compelling enough that the reaction is unambiguous: either they light up or they don't."),

    ("Tempering: Stabilization",
     "Current focus: Stabilization before early-next-week release. Feature is functionally complete. Bar for new opera is high — only customer-facing risks."),
]


def render_sigil(top_arcana, input_label):
    """Render a sigil from top-3 matched arcana."""
    # The sigil is composed from the visual properties of matched arcana
    a1, s1 = top_arcana[0]
    a2, s2 = top_arcana[1]
    a3, s3 = top_arcana[2]

    c1, c2, c3 = a1["color"], a2["color"], a3["color"]
    sym1, sym2, sym3 = a1["symbol"], a2["symbol"], a3["symbol"]
    h1, h2, h3 = a1["hebrew"], a2["hebrew"], a3["hebrew"]

    # Build a 7-line sigil from the three arcana
    # Primary arcana dominates center, secondary flanks, tertiary corners

    lines = []

    # Line 1: tertiary accents
    lines.append(f"    {c3}{sym3}       {sym3}{RESET}")
    # Line 2: secondary
    lines.append(f"  {c2}{sym2}   {c1}{h1}{RESET}   {c2}{sym2}{RESET}")
    # Line 3: primary dominant
    lines.append(f"  {c1}{sym1} {sym1} {sym1} {sym1} {sym1}{RESET}")
    # Line 4: center — hebrew letters of all three
    lines.append(f"  {c2}{h2} {c1}  {sym1}  {RESET}{c2}{h2}{RESET}")
    # Line 5: primary
    lines.append(f"  {c1}{sym1} {sym1} {sym1} {sym1} {sym1}{RESET}")
    # Line 6: secondary
    lines.append(f"  {c2}{sym2}   {c3}{h3}{RESET}   {c2}{sym2}{RESET}")
    # Line 7: tertiary
    lines.append(f"    {c3}{sym3}       {sym3}{RESET}")

    return lines


def render_sigil_v2(top_arcana, input_label):
    """Alternative: denser, more abstract sigil using Braille + symbols."""
    a1, s1 = top_arcana[0]
    a2, s2 = top_arcana[1]
    a3, s3 = top_arcana[2]

    c1, c2, c3 = a1["color"], a2["color"], a3["color"]

    # Use similarity scores to modulate density
    # Higher score = denser Braille
    def score_to_braille(s):
        if s > 0.5: return "⣿"
        if s > 0.4: return "⣶"
        if s > 0.3: return "⡶"
        if s > 0.2: return "⠶"
        return "⠄"

    b1, b2, b3 = score_to_braille(s1), score_to_braille(s2), score_to_braille(s3)

    lines = []
    lines.append(f"  {c3}{b3}{a3['symbol']}{b3}{RESET} {c1}{b1}{b1}{b1}{RESET} {c2}{b2}{a2['symbol']}{b2}{RESET}")
    lines.append(f"  {c1}{b1}{b1}{a1['symbol']}{b1}{b1}{a1['symbol']}{b1}{b1}{RESET}")
    lines.append(f"  {c2}{b2}{RESET}{c1} {a1['hebrew']} {a1['symbol']} {a1['hebrew']} {RESET}{c2}{b2}{RESET}")
    lines.append(f"  {c1}{b1}{b1}{a1['symbol']}{b1}{b1}{a1['symbol']}{b1}{b1}{RESET}")
    lines.append(f"  {c2}{b2}{a2['symbol']}{b2}{RESET} {c1}{b1}{b1}{b1}{RESET} {c3}{b3}{a3['symbol']}{b3}{RESET}")

    return lines


def main():
    # Embed arcana descriptions
    arcana_texts = [a["desc"] for a in ARCANA]
    input_texts = [text for _, text in INPUTS]
    input_labels = [label for label, _ in INPUTS]

    if USE_SBERT:
        arcana_emb = embed(arcana_texts)
        input_emb = embed(input_texts)
    else:
        from sklearn.feature_extraction.text import TfidfVectorizer
        vec = TfidfVectorizer(max_features=1000, stop_words='english', ngram_range=(1,2))
        all_texts = arcana_texts + input_texts
        all_emb = vec.fit_transform(all_texts).toarray()
        arcana_emb = all_emb[:len(arcana_texts)]
        input_emb = all_emb[len(arcana_texts):]

    # For each input, find top-3 matching arcana
    sims = cosine_similarity(input_emb, arcana_emb)

    print("\n" + "=" * 72)
    print("  SVSI END-TO-END DEMO")
    print("  text → embed → match arcana → render sigil")
    print("=" * 72)

    for i, (label, text) in enumerate(INPUTS):
        scores = sims[i]
        top_idx = np.argsort(scores)[::-1][:5]

        print(f"\n{'─' * 72}")
        print(f"  INPUT: {label}")
        print(f"  \"{text[:80]}...\"")
        print()
        print(f"  TOP ARCANA MATCHES:")
        for rank, idx in enumerate(top_idx[:5]):
            card = ARCANA[idx]
            bar = '█' * int(scores[idx] * 40) + '░' * (40 - int(scores[idx] * 40))
            print(f"    {rank+1}. {card['color']}{card['symbol']} {card['name']:<18}{RESET} "
                  f"{scores[idx]:.3f} {bar}")

        # Render sigils
        top3 = [(ARCANA[top_idx[j]], scores[top_idx[j]]) for j in range(3)]

        print(f"\n  SIGIL (composed):          SIGIL (dense):")
        s1 = render_sigil(top3, label)
        s2 = render_sigil_v2(top3, label)
        max_lines = max(len(s1), len(s2))
        for line_idx in range(max_lines):
            left = s1[line_idx] if line_idx < len(s1) else ""
            right = s2[line_idx] if line_idx < len(s2) else ""
            # Rough alignment — ANSI codes make this tricky
            print(f"  {left:<40} {right}")

    # ============================================================
    # SIMILARITY COMPARISON — do same-person witnesses cluster?
    # ============================================================
    print(f"\n{'=' * 72}")
    print(f"  ARCANA-SPACE SIMILARITY")
    print(f"  Do same-person witnesses get similar arcana profiles?")
    print(f"{'=' * 72}")

    # Compare arcana profiles (not raw embeddings)
    # Profile = the full similarity vector against all 22 arcana
    profiles = sims  # shape: (n_inputs, 22)

    pairs = [
        (0, 1, "Artifex/Alignment vs Artifex/Harvest"),
        (0, 2, "Artifex/Alignment vs Kelly/Analytics"),
        (0, 3, "Artifex/Alignment vs Horse Owners"),
        (0, 4, "Artifex/Alignment vs Calcinatio:Audit"),
        (2, 3, "Kelly/Analytics vs Horse Owners"),
        (3, 7, "Horse Owners witness vs Horse Intent"),
        (4, 5, "Calcinatio:Audit vs Calcinatio:Customer"),
        (6, 7, "Intent:Alignment vs Intent:Horse"),
    ]

    print(f"\n  {'Pair':<48} {'Profile Sim':>10}")
    print(f"  {'─' * 60}")
    for i1, i2, label in pairs:
        p1 = profiles[i1:i1+1]
        p2 = profiles[i2:i2+1]
        psim = cosine_similarity(p1, p2)[0, 0]
        bar = '█' * int(psim * 20) + '░' * (20 - int(psim * 20))
        print(f"  {label:<48} {psim:.3f}  {bar}")

    print()


if __name__ == '__main__':
    main()
