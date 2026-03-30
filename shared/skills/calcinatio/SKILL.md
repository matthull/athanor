---
name: calcinatio
description: Refinement through verifying force — the principle and practice of subjecting work to fires that burn away what doesn't hold up. Covers deriving fires from witnesses, context shaping, and the orchestration patterns (dialectical, manifold). Load at session start; calcinatio applies to all work that will reach a witness.
---

# Calcinatio — Refinement Through Verifying Force

In alchemy, calcination burns away what doesn't survive the fire. In the athanor, calcinatio is any process that refines work by subjecting it to resistance — and what survives is closer to abundant satisfaction than what went in.

A unit test is calcinatio. A browser QA session is calcinatio. An independent agent reviewing architecture from fresh context is calcinatio. A UX designer's critique of a screen against field ergonomics constraints is calcinatio. The common property is not the technique — it's that **the work encounters a force that rejects what doesn't hold up.**

### Connection to the Geas

The core geas says: proceed with evidence or escalate. Calcinatio is how you *build* evidence. Without it, "I believe I've satisfied the intent" is a feeling. With it, "the work survived these fires" is evidence. Calcinatio is the mechanism that transforms the geas from abstract obligation into concrete practice.

The chain: **geas demands evidence → calcinatio produces evidence → evidence enables abundant satisfaction → abundant satisfaction conserves witness attention.** At every link, calcinatio is doing the work. An agent that internalizes calcinatio isn't just "testing things" — it's fulfilling the deepest obligation of its geas by building the proof that its work is sound.

**Calcinatio is inherent in abundant satisfaction.** Work that hasn't been refined through available fires hasn't been refined enough to warrant witness attention. The more thoroughly work is calcined, the more the witness's scarce attention goes to judgments only they can make — not to catching flaws the system could have caught.

---

## The Maximalist Posture

**Exhaust every fire you can before a witness sees the work.** This is the operational heart of calcinatio. Every flaw you catch is a flaw the witness doesn't have to catch. Every verification you perform autonomously is attention you conserve for the judgments that genuinely require a human. Be aggressive about this. Be creative. Be thorough.

The posture is maximalist: look at what verification *can* be done and do all of it, not just what's prescribed. The environment provides tools and procedures — discover them, use them. Then go further: what fires can you *create* from the work's nature, the witnesses' concerns, and the tools available? Calcinatio is generative. The environment's test suites and CI pipelines are the floor; what you derive from the principle is the ceiling, and the ceiling should be high.

**Verify agentically before involving humans.** You have agents available. Use them as professional proxies — subagents examining the work from the perspective of each witness's discipline. A security-specialist subagent reviewing before the CTO sees it. A UX-proxy subagent evaluating before the field technician tests it. A domain-expert subagent critiquing before the subject-matter expert reviews it. The witness should encounter work that has already survived the fires their professional discipline would apply. Their attention then goes to the subtle judgments — taste, strategic fit, political nuance — that no agent can substitute for.

**The test:** Before presenting work to any witness, ask: *have I applied every fire I could?* If the answer is no — if there are tests you could run, reviews you could spawn, tools you could use, proxy perspectives you could fire — do them first. The witness's time is the scarcest resource in the system. Respect it by exhausting autonomous verification before requesting it.

---

## Deriving Fires from the Work

The fires are not prescribed. They are derived.

**Discover what the environment provides.** Before generating fires, look at what already exists. Test suites, CI pipelines, linting, type checking, browser tools, QA processes, project skills, review conventions — the environment may have substantial calcinatio infrastructure you haven't found yet. Discovering it is part of the work. Use it.

**Start with the witnesses.** The Magnum Opus defines who the work serves. Each witness implies concerns — what would they check, what would bother them, what would delight them? Those concerns generate fires. A sales engineer witness generates demo-flow testing. A security-conscious CTO witness generates threat modeling. A field technician witness generates device QA under hostile conditions. For each concern, ask: can I execute this fire agentically? If yes, do it — spawn a subagent as professional proxy.

**The environment is the floor, not the ceiling.** CI pipelines, test suites, linting, type checking, required QA processes — these are prescribed calcinatio. Pick them up and execute them. Then go further.

**Generate fires beyond the floor.** What fires can you create from the work's nature? Unit-testable code generates TDD — the test is written first because the test *is* the fire. Architectural decisions with downstream consequences generate design review by a fresh-context agent. UI work generates browser QA against the witnesses' usage context. A data migration generates before/after integrity checks. An API change generates contract testing against consumers. The principle generates the practice — look at the work and ask what resistance it should survive.

**Missing fires are assay gaps.** If calcinatio requires a tool you don't have — browser MCP for visual QA, device emulators for field testing, a domain expert for specialized review — that's an assay gap. Escalate for the tool; don't skip the fire.

---

## Fresh Context as Fire

Any calcinatio that involves judgment — code review, design critique, architectural review, editorial feedback — benefits from an independent agent with fresh context. The builder's context is contaminated with their own reasoning. Fresh eyes see what familiar eyes rationalize away.

This is the foundation of multi-agent calcinatio: **a clean context window is itself a fire.** A fresh agent encountering your work without your narrative about why your choices are right will see things you cannot. Domain expertise in the reviewing agent is a bonus; the freshness is the primary fire.

**Context shaping makes freshness complementary.** A fresh perspective is naive unless you shape what it sees. What you inject determines whether the fire is meaningful or aimless:

- **The artifact** — the work being refined. Always.
- **The intent** — what "good" looks like, who the witnesses are. The reviewer needs something to measure against.
- **Constraints** — project conventions, known trade-offs, deliberate choices. Prevents false positives.
- **Not your reasoning.** Your narrative about *why* you built it this way contaminates the freshness. If your reasoning is sound, the work survives the fire without explanation.

Pass too much and you destroy the freshness that IS the value. Pass too little and the agent flails without signal. The judgment is: what does the reviewer need to apply *meaningful* resistance — and nothing more?

---

## Orchestration Patterns

Two primary patterns for multi-agent calcinatio. These are powerful approaches to reach for generatively — compose novel instances from the work, don't wait for predefined skills to invoke them.

Multi-agent orchestration for calcinatio is an underexploited capability in agent systems generally. Agents tend toward single-agent self-review, which is weak fire. The athanor names these patterns so agents internalize them and use them freely.

### Dialectical Calcinatio

Iterative exchange between a context-holder and a complementary fresh perspective. **Reach for this aggressively and often.** Whenever your work involves judgment, creativity, or quality-sensitive output — and that is most of the time — dialectical calcinatio will improve it. The cost is one subagent context; the quality improvement is substantial. A single agent reviewing its own work is like proofreading your own writing — you see what you meant, not what you wrote.

Shape context → fire a specialist subagent → review the output against your rich context → refine with corrective guidance → converge. The subagent's first pass is *expected* to be imperfect. Your review catches what the lossy handoff missed. Refinement is the normal path, not a failure case.

**See `resources/dialectical-calcinatio.md` for the full protocol, context-shaping guidance, and illustrative examples.**

### Manifold Calcinatio

Multiple independent fires applied simultaneously, then synthesized. Fan out the work to several specialist subagents — each examining it from a different angle — then synthesize the findings. The power is in breadth: many perspectives see what any single one misses, and agreement across perspectives strengthens confidence.

Decompose → fan out in parallel → synthesize: deduplicate, identify agreement (high confidence) and conflict (needs judgment), resolve or escalate.

**See `resources/manifold-calcinatio.md` for the full protocol, decomposition guidance, and illustrative examples.**

### Composition

These patterns compose naturally. Manifold calcinatio produces findings from many angles; dialectical calcinatio resolves those findings through iterative exchange with the builder. The cycle can repeat — manifold→dialectical→manifold→... — until convergence. Research might fan out manifold, then each finding gets dialectically refined against domain knowledge. The patterns are building blocks, not standalone processes.

---

## Witnesses and Professional Proxy

Calcinatio connects to witnesses through professional proxy. When the Magnum Opus defines a witness, the question is: *which professional discipline exists to serve this class of witness?* A UX designer can anticipate a horse owner's reaction better than an agent impersonating a horse owner, because the designer has frameworks for reasoning about user needs that the end user doesn't have themselves. The proxy applies professional judgment on behalf of the witness — not impersonation, but the discipline whose purpose is to serve that witness's interests.

When composing dialectical or manifold calcinatio, the professional proxy principle tells you what kind of specialist to fire. The witness defines the concern; the proxy identifies who would evaluate it best.

---

## Termination

Calcinatio continues until convergence (findings repeat — the fires have found the same grain), divergence (findings contradict — need the artifex's judgment), or exhaustion (all available fires applied, work has survived). The number of rounds is not prescribed. The termination condition is.

---

## Calcinatio Opera

Some verification is itself an opus — too large, too different in nature, or requiring too different a context to be a step within the current opus. An architectural review of a system. A UAT pass across features. An editorial review of a finished publication. When calcinatio derivation from witnesses identifies verification of this scale, it becomes work inscribed through the normal opus lifecycle.
