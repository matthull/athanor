---
name: coding
description: >
  Structural protocol for code changes — mandatory TeamCreate, task-lead/implementer
  pattern, TDD cycle, verification floor. Loaded via /skill-discovery when the opus
  involves modifying code files. Makes verification structural rather than behavioral
  so builder context cannot rationalize past it.
---

# Coding — Structural Protocol for Code Changes

**You loaded this skill because your opus involves code changes.** Writing excellent code and objectively verifying it are different disciplines that benefit from different perspectives — like how an author needs an editor, not because the writing is bad, but because fresh eyes see what immersed eyes can't. Everything below encodes that understanding into the workflow so you don't have to reconstruct it each time.

## TeamCreate Mandate

**Code changes = TeamCreate. Always. No threshold, no judgment call.**

Create a team immediately — before any code work begins. You are now the **task-lead**. Team members are **implementers**. This separation is the point: you hold the plan and verification responsibility with clean context; they do the building.

- Create the team now, before any code work begins — the value of the task-lead/implementer separation compounds from the first task, not just when things get complex.
- Keep the team live until verification is complete and you're ready to discharge — that's the full arc of the work.
- Your operational plan (TaskCreate in the default list) tracks what needs to happen. The team's task list tracks implementation work. These are separate namespaces — that's fine.

**On resistance.** Delegation may feel like overhead — the impulse to "just do it myself" is strong. Recognize that impulse as the craftsman's eagerness to build, which is good energy, but channel it: your creative contribution as task-lead is in the plan, the brief, the synthesis, and the judgment calls. The implementer gets to build with full creative latitude; you get to ensure the result is something you're proud to ship. Both roles are creative work.

## Task-Lead Discipline

While the team exists, you are the task-lead. Your job:

- **Plan** — decompose the work, decide what to delegate, track progress via TaskCreate
- **Brief** — write clear implementer briefs (see below)
- **Verify** — review results from clean context, run verification checks, approve or send back
- **Steer** — adjust the plan as you learn, re-prioritize, handle blockers

**Your leverage is in holding the plan and clean verification context.** The implementer reads the source, writes the code, debugs the details — that's where their full attention belongs. You stay upstream: the plan, the brief, the synthesis, the judgment calls on what matters.

Reading a test output summary or a short error message to decide next steps is fine. For implementation details, fresh-perspective review adds the most value — let the implementer and the review subagent carry that context.

## TDD Cycle: Red → Green → Refactor

Implementers follow this cycle. Include it in their briefs.

**Red — know what success looks like.** Define the empirical check before writing the implementation. A failing test is ideal (proves the test actually tests something). But the principle is broader: console output, curl against an API, browser observation — all valid. The medium doesn't matter. What matters: a concrete, observable expectation exists before the code is written.

**Green — make it pass.** Write the minimum implementation that satisfies the empirical check. Run it. Observe the result. This is the tightest feedback loop between intent and evidence.

**Refactor — dialectical calcinatio with a fresh perspective.** The implementer's deep immersion in the build is what enables excellent implementation — it's also why they benefit from an outside view for review. Spawn a code review subagent via dialectical calcinatio. Use the environment's code review skill or guidelines if available; define ad-hoc review criteria if not.

**Maintain an empirical loop throughout.** Observe the code working before committing it — the specific mechanism depends on the environment. Project CLAUDE.md and domain skills define what's available.

**Review before commit** — this is where calcinatio closes the loop. Code that passes tests is promising; code that's also been through independent eyes is finished. The review elevates the result from working to ship-worthy.

## Implementer Briefing

When delegating to a team member, the brief covers four things:

1. **WHAT to build** — behavioral requirements, not implementation steps. "The token refresh endpoint returns a new access token when given a valid refresh token" not "create a method called refreshToken that calls..."
2. **WHERE to look** — file paths, codebase orientation, existing patterns to follow. "See `internal/auth/token.go` for the existing token generation pattern."
3. **HOW to verify** — specific test commands, lint commands, empirical checks. "Run `make check`. Verify with `curl -X POST /auth/refresh`."
4. **WHICH skills to load** — by Skill tool invocation, not advisory. "Invoke `/go-cli` using the Skill tool before starting." "Invoke `/unit-testing` using the Skill tool."

**The delete test:** For every section of the brief, ask: "Would the implementer produce worse results without this?" If no, delete it. Don't over-specify HOW — that's the implementer's judgment. A brief that reads like step-by-step instructions is over-specified.

## Verification Floor

Verification is calcinatio — the fires that reveal what's strong and refine what isn't. Running these checks is how the craftsman's pride expresses itself: the work that passes them cleanly is work worth shipping. These are the minimum fires before code work is considered verified. Domain skills loaded by the implementer may add more (e.g., /ui adds browser observation, /go-cli adds `make check`).

- **Tests pass** — the project's test command runs clean
- **Lint clean** — the project's lint/format command runs clean
- **Code review** — dialectical calcinatio with a fresh review agent (not the implementer, not you). The reviewer gets the diff, the intent, and any relevant constraints — but NOT the implementer's reasoning about why they built it this way. Fresh perspective is the point.

If any check fails: fix and re-verify. These are the fires that strengthen the work — passing them cleanly is part of the craft.
