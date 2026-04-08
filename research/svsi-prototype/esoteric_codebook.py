"""
Esoteric Codebook for SVSI Layer 2-3

Maps the 22 Major Arcana to Tree of Life paths, planetary attributions,
and visual abstractions. Minor arcana map to sephirot + elements.

This is the structured intermediate vocabulary that normalizes
raw semantic vectors into a rich symbolic space before visualization.

Sources: Golden Dawn correspondences (standard Western esoteric tradition)
"""

# ============================================================
# MAJOR ARCANA — 22 paths on the Tree of Life
# ============================================================
# Each card maps to: path (from→to sephirot), hebrew letter,
# astrological attribution, element/planet/sign, and visual abstractions

MAJOR_ARCANA = {
    0: {
        "name": "The Fool",
        "hebrew": "א",  # Aleph
        "path": ("Keter", "Chokmah"),
        "attribution": "Air",
        "visual": {
            "shape": "spiral",
            "color": "yellow",
            "quality": "ascending, scattered, light",
            "symbol": "○",
            "density": 0.2,
            "symmetry": "none",
            "motion": "erratic",
        },
    },
    1: {
        "name": "The Magician",
        "hebrew": "ב",  # Beth
        "path": ("Keter", "Binah"),
        "attribution": "Mercury",
        "visual": {
            "shape": "vertical_line",
            "color": "yellow",
            "quality": "focused, channeling, precise",
            "symbol": "☿",
            "density": 0.7,
            "symmetry": "bilateral",
            "motion": "steady_downward",
        },
    },
    2: {
        "name": "The High Priestess",
        "hebrew": "ג",  # Gimel
        "path": ("Keter", "Tiferet"),
        "attribution": "Moon",
        "visual": {
            "shape": "crescent",
            "color": "blue",
            "quality": "veiled, reflective, deep",
            "symbol": "☽",
            "density": 0.5,
            "symmetry": "bilateral",
            "motion": "still",
        },
    },
    3: {
        "name": "The Empress",
        "hebrew": "ד",  # Daleth
        "path": ("Chokmah", "Binah"),
        "attribution": "Venus",
        "visual": {
            "shape": "circle_filled",
            "color": "green",
            "quality": "fertile, expanding, abundant",
            "symbol": "♀",
            "density": 0.9,
            "symmetry": "radial",
            "motion": "outward",
        },
    },
    4: {
        "name": "The Emperor",
        "hebrew": "ה",  # He
        "path": ("Chokmah", "Tiferet"),
        "attribution": "Aries",
        "visual": {
            "shape": "square",
            "color": "red",
            "quality": "structured, commanding, rigid",
            "symbol": "♈",
            "density": 0.8,
            "symmetry": "4-fold",
            "motion": "stable",
        },
    },
    5: {
        "name": "The Hierophant",
        "hebrew": "ו",  # Vav
        "path": ("Chokmah", "Chesed"),
        "attribution": "Taurus",
        "visual": {
            "shape": "cross",
            "color": "red-orange",
            "quality": "teaching, bridging, institutional",
            "symbol": "♉",
            "density": 0.7,
            "symmetry": "bilateral",
            "motion": "stable",
        },
    },
    6: {
        "name": "The Lovers",
        "hebrew": "ז",  # Zayin
        "path": ("Binah", "Tiferet"),
        "attribution": "Gemini",
        "visual": {
            "shape": "twin_points",
            "color": "orange",
            "quality": "dual, choosing, connecting",
            "symbol": "♊",
            "density": 0.6,
            "symmetry": "bilateral",
            "motion": "oscillating",
        },
    },
    7: {
        "name": "The Chariot",
        "hebrew": "ח",  # Cheth
        "path": ("Binah", "Geburah"),
        "attribution": "Cancer",
        "visual": {
            "shape": "chevron",
            "color": "orange-yellow",
            "quality": "directed, contained, moving",
            "symbol": "♋",
            "density": 0.8,
            "symmetry": "bilateral",
            "motion": "forward",
        },
    },
    8: {
        "name": "Strength",
        "hebrew": "ט",  # Teth
        "path": ("Chesed", "Geburah"),
        "attribution": "Leo",
        "visual": {
            "shape": "lemniscate",
            "color": "yellow-green",
            "quality": "enduring, gentle_power, sustained",
            "symbol": "♌",
            "density": 0.7,
            "symmetry": "bilateral",
            "motion": "circulating",
        },
    },
    9: {
        "name": "The Hermit",
        "hebrew": "י",  # Yod
        "path": ("Chesed", "Tiferet"),
        "attribution": "Virgo",
        "visual": {
            "shape": "point",
            "color": "green-yellow",
            "quality": "solitary, inward, illuminating",
            "symbol": "♍",
            "density": 0.3,
            "symmetry": "none",
            "motion": "contracting",
        },
    },
    10: {
        "name": "Wheel of Fortune",
        "hebrew": "כ",  # Kaph
        "path": ("Chesed", "Netzach"),
        "attribution": "Jupiter",
        "visual": {
            "shape": "wheel",
            "color": "violet",
            "quality": "cyclical, turning, expansive",
            "symbol": "♃",
            "density": 0.6,
            "symmetry": "radial",
            "motion": "rotating",
        },
    },
    11: {
        "name": "Justice",
        "hebrew": "ל",  # Lamed
        "path": ("Geburah", "Tiferet"),
        "attribution": "Libra",
        "visual": {
            "shape": "balance",
            "color": "green",
            "quality": "equilibrium, weighing, precise",
            "symbol": "♎",
            "density": 0.5,
            "symmetry": "bilateral",
            "motion": "still",
        },
    },
    12: {
        "name": "The Hanged Man",
        "hebrew": "מ",  # Mem
        "path": ("Geburah", "Hod"),
        "attribution": "Water",
        "visual": {
            "shape": "inverted_triangle",
            "color": "blue",
            "quality": "suspended, reversed, accepting",
            "symbol": "▽",
            "density": 0.4,
            "symmetry": "bilateral",
            "motion": "suspended",
        },
    },
    13: {
        "name": "Death",
        "hebrew": "נ",  # Nun
        "path": ("Tiferet", "Netzach"),
        "attribution": "Scorpio",
        "visual": {
            "shape": "scythe_arc",
            "color": "blue-green",
            "quality": "transforming, clearing, absolute",
            "symbol": "♏",
            "density": 0.6,
            "symmetry": "none",
            "motion": "sweeping",
        },
    },
    14: {
        "name": "Temperance",
        "hebrew": "ס",  # Samekh
        "path": ("Tiferet", "Yesod"),
        "attribution": "Sagittarius",
        "visual": {
            "shape": "hourglass",
            "color": "blue",
            "quality": "blending, flowing, moderate",
            "symbol": "♐",
            "density": 0.5,
            "symmetry": "bilateral",
            "motion": "pouring",
        },
    },
    15: {
        "name": "The Devil",
        "hebrew": "ע",  # Ayin
        "path": ("Tiferet", "Hod"),
        "attribution": "Capricorn",
        "visual": {
            "shape": "inverted_pentagon",
            "color": "indigo",
            "quality": "binding, material, dense",
            "symbol": "♑",
            "density": 0.9,
            "symmetry": "bilateral",
            "motion": "static",
        },
    },
    16: {
        "name": "The Tower",
        "hebrew": "פ",  # Pe
        "path": ("Netzach", "Hod"),
        "attribution": "Mars",
        "visual": {
            "shape": "lightning",
            "color": "red",
            "quality": "shattering, sudden, liberating",
            "symbol": "♂",
            "density": 0.8,
            "symmetry": "none",
            "motion": "explosive",
        },
    },
    17: {
        "name": "The Star",
        "hebrew": "צ",  # Tzaddi
        "path": ("Netzach", "Yesod"),
        "attribution": "Aquarius",
        "visual": {
            "shape": "star_8point",
            "color": "violet",
            "quality": "radiating, hopeful, clear",
            "symbol": "✦",
            "density": 0.4,
            "symmetry": "radial",
            "motion": "radiating",
        },
    },
    18: {
        "name": "The Moon",
        "hebrew": "ק",  # Qoph
        "path": ("Netzach", "Malkuth"),
        "attribution": "Pisces",
        "visual": {
            "shape": "double_crescent",
            "color": "crimson",
            "quality": "dreamlike, shifting, uncertain",
            "symbol": "♓",
            "density": 0.5,
            "symmetry": "bilateral",
            "motion": "wavering",
        },
    },
    19: {
        "name": "The Sun",
        "hebrew": "ר",  # Resh
        "path": ("Hod", "Yesod"),
        "attribution": "Sun",
        "visual": {
            "shape": "radiant_circle",
            "color": "orange",
            "quality": "bright, clear, joyful",
            "symbol": "☉",
            "density": 0.8,
            "symmetry": "radial",
            "motion": "radiating",
        },
    },
    20: {
        "name": "Judgement",
        "hebrew": "ש",  # Shin
        "path": ("Hod", "Malkuth"),
        "attribution": "Fire",
        "visual": {
            "shape": "three_flames",
            "color": "red",
            "quality": "awakening, rising, decisive",
            "symbol": "🜂",
            "density": 0.7,
            "symmetry": "3-fold",
            "motion": "ascending",
        },
    },
    21: {
        "name": "The World",
        "hebrew": "ת",  # Tav
        "path": ("Yesod", "Malkuth"),
        "attribution": "Saturn",
        "visual": {
            "shape": "mandorla",
            "color": "indigo",
            "quality": "complete, integrated, whole",
            "symbol": "♄",
            "density": 0.9,
            "symmetry": "4-fold",
            "motion": "rotating_slow",
        },
    },
}

# ============================================================
# SEPHIROT — 10 spheres, used for minor arcana mapping
# ============================================================

SEPHIROT = {
    1: {
        "name": "Keter",
        "hebrew": "כתר",
        "english": "Crown",
        "quality": "unity, source, will",
        "color": "white",
        "symbol": "✦",
        "planet": None,  # beyond planetary
    },
    2: {
        "name": "Chokmah",
        "hebrew": "חכמה",
        "english": "Wisdom",
        "quality": "force, intent, flash",
        "color": "grey",
        "symbol": "⊛",
        "planet": "Zodiac",
    },
    3: {
        "name": "Binah",
        "hebrew": "בינה",
        "english": "Understanding",
        "quality": "form, structure, witness",
        "color": "black",
        "symbol": "◈",
        "planet": "Saturn",
    },
    4: {
        "name": "Chesed",
        "hebrew": "חסד",
        "english": "Mercy",
        "quality": "expansion, generosity, building",
        "color": "blue",
        "symbol": "❋",
        "planet": "Jupiter",
    },
    5: {
        "name": "Geburah",
        "hebrew": "גבורה",
        "english": "Severity",
        "quality": "constraint, fire, cutting",
        "color": "red",
        "symbol": "⚡",
        "planet": "Mars",
    },
    6: {
        "name": "Tiferet",
        "hebrew": "תפארת",
        "english": "Beauty",
        "quality": "harmony, balance, heart",
        "color": "yellow",
        "symbol": "◉",
        "planet": "Sun",
    },
    7: {
        "name": "Netzach",
        "hebrew": "נצח",
        "english": "Victory",
        "quality": "desire, endurance, drive",
        "color": "green",
        "symbol": "◆",
        "planet": "Venus",
    },
    8: {
        "name": "Hod",
        "hebrew": "הוד",
        "english": "Splendor",
        "quality": "communication, form, structure",
        "color": "orange",
        "symbol": "◫",
        "planet": "Mercury",
    },
    9: {
        "name": "Yesod",
        "hebrew": "יסוד",
        "english": "Foundation",
        "quality": "channel, connection, flow",
        "color": "violet",
        "symbol": "○",
        "planet": "Moon",
    },
    10: {
        "name": "Malkuth",
        "hebrew": "מלכות",
        "english": "Kingdom",
        "quality": "manifestation, completion, earth",
        "color": "earth_tones",
        "symbol": "▣",
        "planet": "Earth",
    },
}

# ============================================================
# ELEMENTS — for suit mapping
# ============================================================

ELEMENTS = {
    "Fire": {"symbol": "△", "color": "red", "quality": "will, action, energy"},
    "Water": {"symbol": "▽", "color": "blue", "quality": "emotion, flow, depth"},
    "Air": {"symbol": "◇", "color": "yellow", "quality": "thought, communication, clarity"},
    "Earth": {"symbol": "□", "color": "green", "quality": "material, body, stability"},
}

# ============================================================
# TEXT DESCRIPTIONS — for embedding-based matching
# ============================================================
# These are what we embed and compare against input text

ARCANA_DESCRIPTIONS = {}
for num, card in MAJOR_ARCANA.items():
    # Create a rich text description for embedding
    v = card["visual"]
    ARCANA_DESCRIPTIONS[f"major_{num}"] = {
        "card": card["name"],
        "text": f"""{card['name']}. {v['quality']}. Attribution: {card['attribution']}.
Path connecting {card['path'][0]} to {card['path'][1]} on the Tree of Life.
Shape: {v['shape']}. Motion: {v['motion']}. Symmetry: {v['symmetry']}.
Density: {'dense' if v['density'] > 0.7 else 'moderate' if v['density'] > 0.4 else 'sparse'}.""",
        "visual": v,
    }

for num, seph in SEPHIROT.items():
    ARCANA_DESCRIPTIONS[f"sephirah_{num}"] = {
        "card": seph["name"],
        "text": f"""{seph['name']} ({seph['english']}). {seph['quality']}.
{'Planet: ' + seph['planet'] + '.' if seph['planet'] else 'Beyond planetary attribution.'}
Color: {seph['color']}. The {seph['english'].lower()} of the Tree of Life.""",
        "visual": {
            "symbol": seph["symbol"],
            "color": seph["color"],
            "quality": seph["quality"],
        },
    }
