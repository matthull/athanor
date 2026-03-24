I want to rethink /orchestrate some. Here are various concerns I want to weave together into a roadmap forward.

Talking to someone today i found they had taken a few shots and building a Vue component and kept going back to the drawing board.
In the end i felt like the issue was they didn't have a verification loop, claude would go build it blind and they'd find it was
broken in various ways. Fundamental issue there is what we're all struggling with figuring out expectations for LLM agent, this dev 
expected claude to be able to build a simple component blind with a detailed spec but it was fucking up basic things over and over. Our
layered verification loops address that, essentially providing automated steering, in case of vue components it's the loop where agent
loads storybook in chrome MCP and checks constantly TDD style while developing. This issue is solved by handoff system which people need
to understand in depth i think (the need for verification loop steering) before graduating to higher-level orchestration. There's this key process
of focusing on validation/verification baked deeply into your workflow at every level that is more important than e2e orchestration, and I'd like
to make that a bit more of a first-order concern. It's close now but doesn't hurt to promote it for explicit focus e.g. during /orchestrate triage step
have high level testing methodology defined along with goals, have that be explicit focus for /spec  (e.g. build test tooling that is missing), etc.
Injection in depth more consistently for that concern.

Next up: /orchestrate tries to use one long-lived orchestrator session after planning phase. We keep adding guardrails to keep it on track, and i think
all the guardrails are needed but not enough for such a long-lived session. I think we need to have structured phase handoffs so there is not necessarily 
one long lived orchestrator. There _could_ be one, for instance a master orchestrator that spawns an indepdent tmux session with claude per phase, that session
following our normal orchestrator+TeamCreate model but updating a handoff and ending after that phase. But then you could also do /direct-handoffs, or human operator
could start a new session per phase... this illustrates the idea of abstracting the task/project flow from orchestration mechanisms. The transition of triage phase to 
execution currently models this, but we abstract away the phase documentation and transition mechanisms from the mechanics. Then everyone can use their favorite orchestration
tools but probably start with none at all just manually administering phase handoffs, using /compact or /clear, load last phase handoff, continue. Protocol and orchestration mechanism
separated.

This means we need to focus on how we document what is to be done, what was done... it has to be done in a fairly agnostic way or a pluggable way. Like standard concept of phase documentation
plan doc, handoffs, etc. then phase definitions can fill in some stuff to include. Underneath this we should have a detailed logging to go back to if needed... the protocl for that should also be 
separate from implementation. We can probably default to markdown docs in task dir but 

Structurally to prevent sprawl we should probably define a concept of a task dir. Where it is can be pluggable but we need a place to collect all the docs for a workflow. Also we have to tweak our terminiology a bit, currently we say "task" for task handoffs, but need a name for
the overall workflow we are executing, task makes sense for that too sense it is short lived (-1 day) and goal oriented. need to consider that some. "Workflow" doesn't really work because often that means an abstraction e.g. a repeatable series of steps not a single execution. Help me brainstorm terminiology.
