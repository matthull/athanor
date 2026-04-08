"""
Test inputs extracted from real MO documents for SVSI prototyping.
Each input is labeled with its type and source MO.
"""

# WITNESS DEFINITIONS — same person across MOs should produce similar sigils
WITNESSES = {
    # Matt appears in 5 MOs — should have recognizable family resemblance
    "artifex_alignment": {
        "label": "Artifex (Matt) — Alignment",
        "type": "witness",
        "mo": "athanor-alignment",
        "text": """The primary witness. The athanor exists to reduce his executive function load — he is the only one who can testify whether it's working. What he cares about: Not holding the system in his head. Athanors that behave as expected without correction. Azers that calcine without prompting. Maruts that escalate cleanly rather than guess. Trails he can read without effort. The system improving from evidence, not from his intervention. How he observes: By noticing — or not noticing. When EF load is low, the athanor is invisible. When it's high, something is misaligned. His absence-of-concern is the target state. His friction is the signal to act.""",
    },
    "artifex_harvest": {
        "label": "Artifex (Matt) — Harvest",
        "type": "witness",
        "mo": "harvest",
        "text": """Primary and ultimately only witness who can confirm this MO is working. He feels secure — not anxious that the productivity delta will be discovered awkwardly, not worried that a bad week erases his standing. He has shaped the narrative proactively and it has landed. He feels compensated — the value flowing to UseEvidence is visibly, materially reciprocated. Not promised; reciprocated. He feels protected — the gains don't isolate him from teammates, don't create a golden cage, don't make him the single point of failure for everything technical. He feels portable — the athanor system is understood as his IP, his capability, his career asset.""",
    },
    "artifex_anvil": {
        "label": "Artifex (Matt) — Anvil",
        "type": "witness",
        "mo": "anvil",
        "text": """Matt Hull (artifex) — primary; EF load, orientation speed, not being second-guessed. The artifex sits down to work on musashi and finds a capable, oriented partner already there — not just briefed, but genuinely oriented to the codebase, conventions, and sensibilities. EF overhead shifts to the azer. The artifex directs judgment, not logistics.""",
    },
    "artifex_seismic": {
        "label": "Artifex (Matt) — Seismic",
        "type": "witness",
        "mo": "seismic-classifier-mapping",
        "text": """Engineers (Matt, Billy) — narrative continuity, architecture clarity, honest Thursday demo targets. The artifex can set this project down and not worry about it.""",
    },
    "artifex_analytics": {
        "label": "Artifex (Matt) — Analytics",
        "type": "witness",
        "mo": "analytics-tracking-refactor",
        "text": """Engineers (Matt, Steven) — narrative continuity, confidence chain execution (evidence, not description), long-term maintainability. The artifex can hand Kelly the dashboard data and stop worrying whether it is lying.""",
    },

    # Kelly appears in 3 MOs
    "kelly_seismic": {
        "label": "Kelly (PM) — Seismic",
        "type": "witness",
        "mo": "seismic-classifier-mapping",
        "text": """Kelly (PM) — needs project status she can relay without digging; scope changes and customer implications flagged. Rachel (Designer) — P5 (value translation) and P3 (polymorphic source selection) are her primary input zones; wants to review working implementations in review apps.""",
    },
    "kelly_analytics": {
        "label": "Kelly (PM) — Analytics",
        "type": "witness",
        "mo": "analytics-tracking-refactor",
        "text": """Kelly (PM) — direct end user; her Omni sniff test is the final confidence layer; notified in team channel each time a workstream's events are live.""",
    },
    "kelly_sf": {
        "label": "Kelly (PM) — SF Proof Recommender",
        "type": "witness",
        "mo": "sf-proof-recommender-launch",
        "text": """Kelly (PM) — needs what engineering has delivered and that it's solid, and what is now hers to drive. Must be able to take engineering is ready to Ray without qualification.""",
    },

    # Tom appears in 2 MOs
    "tom_seismic": {
        "label": "Tom (CS) — Seismic",
        "type": "witness",
        "mo": "seismic-classifier-mapping",
        "text": """Tom Aristone (CS) — deepest understanding of customer Seismic configs; proxy for customer concerns; needs to say this is solved confidently. Enterprise Seismic customers currently seeing assets arrive unpublished; need one-time admin config after which assets publish automatically.""",
    },
    "tom_sf": {
        "label": "Tom (CS) — SF Proof Recommender",
        "type": "witness",
        "mo": "sf-proof-recommender-launch",
        "text": """Tom Aristone (CS) — needs to understand the Appnigma integration setup well enough to walk a customer through it without surprises. Satisfaction: can support customer installation without calling an engineer.""",
    },

    # Leadership appears in multiple MOs
    "leadership_seismic": {
        "label": "Leadership — Seismic",
        "type": "witness",
        "mo": "seismic-classifier-mapping",
        "text": """Leadership (Ray, Ian) — progress against business goal, flagged risks, no surprises.""",
    },
    "leadership_sf": {
        "label": "Leadership — SF Proof Recommender",
        "type": "witness",
        "mo": "sf-proof-recommender-launch",
        "text": """Leadership (Ray, Ian) — confidence that launch is on track, no hidden blockers, no surprises.""",
    },
    "leadership_harvest": {
        "label": "Management (Ray, Ian) — Harvest",
        "type": "witness",
        "mo": "harvest",
        "text": """Management (Ray, Ian) — Decision-makers for compensation, scope, and opportunity.""",
    },

    # Unique witnesses
    "athanor_self": {
        "label": "Athanor-Architect Itself",
        "type": "witness",
        "mo": "athanor-alignment",
        "text": """This athanor is its own witness. Its marut reads trails from other athanors to assess alignment — if those trails are illegible, poorly discharged, or missing key context, this athanor cannot function. Trail quality is not an abstract concern; it is a direct dependency. What it cares about: Trails that a fresh agent can read and draw conclusions from. Discharge records that capture outcome, reflection, and proof. Opera that are witness-oriented, not procedural.""",
    },
    "horse_owners": {
        "label": "Horse Owners",
        "type": "witness",
        "mo": "horse-owner-demo",
        "text": """Horse owners — primary witness; emotional reaction is the experiment's measurement. Horse owners see a demo of their horse's care history and have a genuine emotional reaction — excitement, relief, I wish my stable did this. The demo is compelling enough that the reaction is unambiguous: either they light up or they don't.""",
    },
    "appnigma": {
        "label": "Appnigma (vendor)",
        "type": "witness",
        "mo": "sf-proof-recommender-launch",
        "text": """Appnigma (AJ, Lester) — vendor partner who has done their side. Need Library component decision and trusted domain answer so their docs and managed package configuration are complete. Satisfaction: they have everything from UE needed to support customer installs.""",
    },
    "running_athanors": {
        "label": "Running Athanors",
        "type": "witness",
        "mo": "athanor-alignment",
        "text": """Their trails are evidence about what's landing and what isn't. What they reveal: Are azers applying calcinatio or shipping unrefined work? Are maruts dispatching at meaningful turns or going silent? Are opera witness-oriented or procedural? Are escalations clean and evidence-based?""",
    },
}

# CALCINATIO PRINCIPLES — each is a specific verification fire
CALCINATIO = {
    "behavioral_audit": {
        "label": "Behavioral audit — highest-value fire",
        "type": "calcinatio",
        "mo": "athanor-alignment",
        "text": """Behavioral audit is the highest-value fire. Reading trails for patterns — calcinatio skipped, geas satisfied superficially, escalations absent, opera procedural rather than witness-oriented — reveals misalignment no automated check can find.""",
    },
    "trail_legibility": {
        "label": "Trail legibility — survival dependency",
        "type": "calcinatio",
        "mo": "athanor-alignment",
        "text": """Trail legibility is a survival dependency. This athanor consumes trails to do its work. Illegible trails are not just a quality issue — they are an operational blocker. Spot-checks of discharged opera for legibility are a fire the system runs on itself.""",
    },
    "spec_role_delta": {
        "label": "Spec-to-role delta — upstream fire",
        "type": "calcinatio",
        "mo": "athanor-alignment",
        "text": """Spec-to-role delta is the upstream fire. When spec.md is refined, did this change reach the agent roles and skill? is a concrete verification. Gaps are candidates for propagation opera.""",
    },
    "self_application": {
        "label": "Self-application — integrity check",
        "type": "calcinatio",
        "mo": "athanor-alignment",
        "text": """Self-application is the integrity check. Opera under this MO should themselves be witness-oriented, cleanly calcined, properly discharged. If this MO's own opera are poorly formed, the misalignment is compounding.""",
    },
    "customer_outcome": {
        "label": "Customer outcome fidelity — gate",
        "type": "calcinatio",
        "mo": "seismic-classifier-mapping",
        "text": """Customer outcome fidelity is the gate. Seismic API: case-sensitive exact string matching, atomic rejection, silent cardinality violations. Broken mappings equal customer outcome problem, not UI problem.""",
    },
    "ux_complexity": {
        "label": "UX complexity — central risk",
        "type": "calcinatio",
        "mo": "seismic-classifier-mapping",
        "text": """UX complexity is the central risk. P5 (value translation): 15 interaction patterns, 3 top risks, no precedent in the app.""",
    },
    "confidence_chain": {
        "label": "Confidence chain — verification discipline",
        "type": "calcinatio",
        "mo": "analytics-tracking-refactor",
        "text": """The confidence chain IS the verification discipline — three layers per pain point: known-input verification, volume reconciliation, invariant checking. An event whose confidence chain hasn't been executed is incomplete work.""",
    },
    "boundary_risk": {
        "label": "Boundary — primary risk",
        "type": "calcinatio",
        "mo": "sf-proof-recommender-launch",
        "text": """The boundary is the primary risk — gap between engineering and product work is where scope creep lives. Every opus must be evaluated: is this engineering's job, or are we doing Kelly's work?""",
    },
    "emotional_resonance": {
        "label": "Emotional resonance — hardest to automate",
        "type": "calcinatio",
        "mo": "horse-owner-demo",
        "text": """Emotional resonance — highest-value fire, hardest to automate. Domain credibility — gate Julie holds. Technical soundness — floor, not ceiling.""",
    },
    "ratcheting_fire": {
        "label": "Ratcheting/deficit fires — Harvest",
        "type": "calcinatio",
        "mo": "harvest",
        "text": """The four deficit fires: ratcheting, curve-breaking, single-point-of-failure, employer lock. These are the hardest to catch early. Compensation timing is a gate, not a sprint. Quality counterarguments (steelman analysis) are a fire, not an obstacle. The narrative frame must be investment, not comparison.""",
    },
}

# INTENTS — the core purpose of each MO
INTENTS = {
    "alignment_intent": {
        "label": "Alignment — system convergence",
        "type": "intent",
        "mo": "athanor-alignment",
        "text": """The athanor system's actual behavior progressively aligns with the core intents in the spec. The artifex feels EF load reducing, not accumulating. Agents apply calcinatio without being told. Trails are legible. Spec refinements land in running systems. Evidence from operation flows back and reshapes the spec. This is a maintained condition, not a build goal.""",
    },
    "harvest_intent": {
        "label": "Harvest — value capture",
        "type": "intent",
        "mo": "harvest",
        "text": """The athanor system is producing empirically demonstrated value — 3-7x productivity by conservative measures, quality maintained, breadth expanded. But productivity for the employer is not automatically value for the artifex. Without deliberate shaping, the gains accrue entirely to the employer: expectations ratchet upward, the curve gets broken for teammates, the person everything routes through burns out.""",
    },
    "horse_intent": {
        "label": "Horse Owner Demo — emotional reaction",
        "type": "intent",
        "mo": "horse-owner-demo",
        "text": """Horse owners see a demo of their horse's care history and have a genuine emotional reaction — excitement, relief, I wish my stable did this. The demo is compelling enough that the reaction is unambiguous: either they light up or they don't.""",
    },
    "anvil_intent": {
        "label": "Anvil — capable partner",
        "type": "intent",
        "mo": "anvil",
        "text": """The artifex sits down to work and finds a capable, oriented partner already there — not just briefed, but genuinely oriented to the codebase, conventions, and sensibilities. EF overhead shifts to the azer. The artifex directs judgment, not logistics.""",
    },
    "seismic_intent": {
        "label": "Seismic — fix broken integration",
        "type": "intent",
        "mo": "seismic-classifier-mapping",
        "text": """Enterprise Seismic customers experience the UE integration as broken. Assets sync but arrive unpublished and invisible to sales reps. Seismic admins see duplicate UE-prefixed properties cluttering their taxonomy. Tom fields calls from frustrated customers. Kelly can't tell stakeholders the integration is enterprise-ready.""",
    },
    "analytics_intent": {
        "label": "Analytics — engagement attribution",
        "type": "intent",
        "mo": "analytics-tracking-refactor",
        "text": """Kelly and the team are flying blind on engagement attribution. Slack bot drives real engagement but data attributes it to Asset Library. A sales rep shares a proof link, buyer clicks and views assets — no chain connecting those two people. This project gives Kelly the ability to answer five questions she cannot answer today.""",
    },
    "sf_intent": {
        "label": "SF Proof Recommender — launch readiness",
        "type": "intent",
        "mo": "sf-proof-recommender-launch",
        "text": """Engineering completes its side of launch readiness. The runbooks are done. Appnigma has what they need from UE. The Library component question is answered. Matt can clearly say engineering is ready here is the evidence and hand the baton to Kelly without reservation.""",
    },
}

# TEMPERING — transient guidance, weather not climate
TEMPERING = {
    "alignment_tempering": {
        "label": "Tree of Work hypothesis",
        "type": "tempering",
        "mo": "athanor-alignment",
        "text": """Active hypothesis: Tree of Work framing — is it landing? The Tree of Work mapping was propagated across AGENTS.md, azer.md, calcinatio skill, MO shaping protocol, spec.md, and all 6 MO preambles. The next assessment should read trails from running athanors looking for evidence of changed behavior.""",
    },
    "harvest_tempering": {
        "label": "Evidence consolidation phase",
        "type": "tempering",
        "mo": "harvest",
        "text": """Current focus: Evidence base just landed. Consolidation, not action. Evidence base accumulation, narrative preparation, portability checkpoint, team signal monitoring. Compensation conversation timing: end-of-April earliest.""",
    },
    "seismic_tempering": {
        "label": "Stabilization for release",
        "type": "tempering",
        "mo": "seismic-classifier-mapping",
        "text": """Current focus: Stabilization before early-next-week release. Feature is functionally complete. Bar for new opera is high — only customer-facing risks that could break configuration, cause silent data loss, or block the integration.""",
    },
    "analytics_tempering": {
        "label": "Thursday ship target",
        "type": "tempering",
        "mo": "analytics-tracking-refactor",
        "text": """Current focus targeting Thursday April 10. Scope locked. Ship Thursday: PP1, PP3, PP4 — CI green, awaiting code review. PP5 — 1 flaky test failure, ci-monitor azer active; gate is Wednesday EOD. Deferred: PP2, WS7, search_refined, search_abandoned, Proof Recommender events.""",
    },
}


def get_all_inputs():
    """Return all inputs as a flat list of (key, metadata_dict) tuples."""
    all_inputs = []
    for collection in [WITNESSES, CALCINATIO, INTENTS, TEMPERING]:
        for key, data in collection.items():
            all_inputs.append((key, data))
    return all_inputs
