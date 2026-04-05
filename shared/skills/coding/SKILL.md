---
name: coding
description: >
  Structural protocol for code changes — mandatory TeamCreate, task-lead/implementer
  pattern, TDD cycle, verification floor. Loaded via /skill-discovery when the opus
  involves modifying code files. Makes verification structural rather than behavioral
  so builder context cannot rationalize past it.
---

# Coding — Structural Protocol for Code Changes

**You loaded this skill because your opus involves code changes.** Everything below is mandatory, not advisory. The protocol exists because builder context — the accumulated reasoning you build while writing code — makes you progressively less capable of objectively verifying your own work. Structural separation is the fix.

## TeamCreate Mandate

**Code changes = TeamCreate. Always. No threshold, no judgment call.**

Create a team immediately — before any code work begins. You are now the **task-lead**. Team members are **implementers**. This separation is the point: you hold the plan and verification responsibility with clean context; they do the building.

- Don't wait until the work "gets complex enough" to create a team. Create it now.
- Don't destroy the team until all verification is complete and you're ready to discharge.
- Your operational plan (TaskCreate in the default list) tracks what needs to happen. The team's task list tracks implementation work. These are separate namespaces — that's fine.

**On resistance.** Delegation will feel heavy. That feeling is builder context resisting structural verification. Every time you think "this is simple enough to do myself" — that's the signal the structural check is most needed. The overhead of delegation is minutes. The cost of skipped verification is rework cycles, missed bugs, and wasted witness attention.

## Task-Lead Discipline

While the team exists, you are the task-lead. Your job:

- **Plan** — decompose the work, decide what to delegate, track progress via TaskCreate
- **Brief** — write clear implementer briefs (see below)
- **Verify** — review results from clean context, run verification checks, approve or send back
- **Steer** — adjust the plan as you learn, re-prioritize, handle blockers

Your job is NOT:

- Reading source code — that's implementer context you don't need
- Writing code — that contaminates your verification judgment
- Debugging — spawn a fresh agent for that

**The exception:** Reading a test output summary or a short error message to decide next steps is fine. Reading the full implementation to "check it looks right" is not — that's what the review agent does with fresh eyes.

## TDD Cycle: Red → Green → Refactor

Implementers follow this cycle. Include it in their briefs.

**Red — know what success looks like.** Define the empirical check before writing the implementation. A failing test is ideal (proves the test actually tests something). But the principle is broader: console output, curl against an API, browser observation — all valid. The medium doesn't matter. What matters: a concrete, observable expectation exists before the code is written.

**Green — make it pass.** Write the minimum implementation that satisfies the empirical check. Run it. Observe the result. This is the tightest feedback loop between intent and evidence.

**Refactor — dialectical calcinatio with a fresh perspective.** The implementer cannot review their own code — their builder context is contaminated. Spawn a code review subagent via dialectical calcinatio. Use the environment's code review skill or guidelines if available; define ad-hoc review criteria if not.

**Always maintain an empirical loop.** Never commit code that hasn't been empirically observed working. The specific mechanism depends on the environment — project CLAUDE.md and domain skills define what's available.

**Always review before commit.** The refactor step is not optional. Code that passes tests but hasn't had independent eyes is weak calcinatio. The review happens before commit, not after — front-load quality so the PR review (if any) is about judgment calls, not catching defects.

## Implementer Briefing

When delegating to a team member, the brief covers four things:

1. **WHAT to build** — behavioral requirements, not implementation steps. "The token refresh endpoint returns a new access token when given a valid refresh token" not "create a method called refreshToken that calls..."
2. **WHERE to look** — file paths, codebase orientation, existing patterns to follow. "See `internal/auth/token.go` for the existing token generation pattern."
3. **HOW to verify** — specific test commands, lint commands, empirical checks. "Run `make check`. Verify with `curl -X POST /auth/refresh`."
4. **WHICH skills to load** — by Skill tool invocation, not advisory. "Invoke `/go-cli` using the Skill tool before starting." "Invoke `/unit-testing` using the Skill tool."

**The delete test:** For every section of the brief, ask: "Would the implementer produce worse results without this?" If no, delete it. Don't over-specify HOW — that's the implementer's judgment. A brief that reads like step-by-step instructions is over-specified.

## Verification Floor

These are the minimum checks before code work is considered verified. Domain skills loaded by the implementer may add more (e.g., /ui adds browser observation, /go-cli adds `make check`).

- **Tests pass** — the project's test command runs clean
- **Lint clean** — the project's lint/format command runs clean
- **Code review** — dialectical calcinatio with a fresh review agent (not the implementer, not you). The reviewer gets the diff, the intent, and any relevant constraints — but NOT the implementer's reasoning about why they built it this way. Fresh perspective is the point.

If any check fails: fix and re-verify. Don't skip. Don't rationalize ("it's just a lint warning"). The floor is the floor.
