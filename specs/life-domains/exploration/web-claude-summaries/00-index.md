# Life-Layer System Documents

Four markdown files, produced from an extended design conversation. Each covers a distinct aspect of the project and can be read independently, though they reference each other.

## Reading order

**For first-read or when re-orienting:** Read all four in order. The sequence moves from conceptual to concrete to operational.

**For refreshing context before a work session:** `04-design-principles.md` is the shortest and pulls the working stance into focus quickly.

**For orienting a future LLM context with the project:** Paste `02-system-architecture.md` as the primary reference; add `01-conceptual-foundation.md` if the LLM needs the why, add `03-athanor-adaptation.md` if Athanor context is missing.

**For evolving the system:** `04-design-principles.md` is the governing document for how to think about changes.

## The documents

### 01-conceptual-foundation.md

Why this project exists. Domain-independent thinking about attention, cognitive prosthetics, and the specific shape of ADHD cognition that makes conventional systems fail. The energy-asymmetry principle. The cognitive function catalog (prospective memory, working memory, salience filtering, temporal horizon, task initiation, interoception, context reinstatement, meta-cognition). The interoception deep-dive. The structure-channel-binding separation.

This is the document that would outlive specific architectural choices. If the system got reimplemented on entirely different infrastructure, this stays valid.

### 02-system-architecture.md

The concrete architectural commitments. The attendant as new agent class. The three shapes of work (vigils, regards, integrations). Opera scaling. The panel/synthesis/calcinatio reasoning pattern. The double-fire calcinatio (too-much, too-little). Silence as default and grooming as first-class work. The system-as-stakeholder framing. Unit scoping as the hard design problem.

This is the document referred to when actually building or when orienting a session with the project.

### 03-athanor-adaptation.md

What transfers from the Athanor spec to the life-layer domain, what needs new thinking, and how the life-layer relates to the existing athanor infrastructure. Bridging material — reads naturally if you know Athanor; less relevant if you don't.

Includes thoughts on spec changes that might eventually follow if life-layer patterns prove out (attendant class as alternative to azer, generative panel pattern, system-as-stakeholder, double-fire calcinatio).

### 04-design-principles.md

The stance that keeps design and development aligned. Principles over rules. Feedback loop asymmetry. Two axes of improvement (tools and inputs). Diagnostic discipline. Silence as success measurement. Ship-the-embarrassing-version. Session-boundary discipline. The graveyard. The working-relationship framing.

Shortest document. Highest frequency of use over the long arc of the project.

## Scope notes

These documents are not implementation specifications. They cover what the system is, why it's structured that way, and what disciplines govern its evolution — not how specifically to build the first version.

The "first slice" content from the conversation (ESP32 pill case switch, Home Assistant, AppDaemon, LLM loop, ambient light, dashboard) is intentionally not captured in these documents. That work is closer to execution planning and should live in whatever project tracking the first build uses — likely an athanor MO once one is inscribed.

## A caveat about confidence

The architecture described is the landed position of an extended thinking conversation, not a validated working system. Specific commitments (the three work shapes, the double-fire pattern, the lens panel) have conceptual coherence but haven't been proven against actual operation. The documents represent the right starting position, not the final position.

The architecture is expected to refine through use. Unit scoping for regards in particular cannot be perfected in design; it has to be tuned against real trail output. Several other commitments are similarly provisional. The commitment to letting the architecture refine through use is itself one of the design principles (in doc 4).
