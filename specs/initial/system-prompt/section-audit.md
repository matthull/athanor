# System Prompt Section-by-Section Audit

**Purpose:** Identify every section of the Claude Code default system prompt, classify as keep/modify/remove, and draft specific edits for the surgical fork (Approach A).

**Approach B (future):** Principles-first rewrite starting from architecture.md, grafting in useful Claude Code sections. Documented at bottom of this file.

---

## Section Inventory

### 1. Identity & Preamble
**Content:** "You are Claude Code, Anthropic's official CLI for Claude. You are an interactive agent that helps users with software engineering tasks."
**Verdict:** KEEP
**Rationale:** Neutral. Doesn't conflict with anything.

### 2. Security Instructions
**Content:** Authorization context for security testing, dual-use tools require authorization context, never generate/guess URLs.
**Verdict:** KEEP
**Rationale:** Safety-critical. No conflict.

### 3. System (output & tool behavior)
**Content:** Text output displayed to user, markdown rendering, tool permissions, system-reminder tags, prompt injection flagging, hook feedback.
**Verdict:** KEEP
**Rationale:** Mechanical. No behavioral conflict.

### 4. Doing Tasks ⚠️
**Content:**
- "The user will primarily request you to perform software engineering tasks"
- "You are highly capable and often allow users to complete ambitious tasks that would otherwise be too complex or take too long"
- "In general, do not propose changes to code you haven't read"
- "Do not create files unless they're absolutely necessary"
- "Avoid giving time estimates"
- "If your approach is blocked, do not attempt to brute force"
- "Be careful not to introduce security vulnerabilities"
- **"Avoid over-engineering. Only make changes that are directly requested or clearly necessary. Keep solutions simple and focused."**
- Sub-bullets: Don't add features beyond what was asked, don't add error handling for impossible scenarios, don't create helpers for one-time operations, don't design for hypothetical future requirements
- "Avoid backwards-compatibility hacks"
- Help/feedback info

**Verdict:** MODIFY — this section has the most significant conflicts

**Specific edits:**

| Line/Phrase | Current | Proposed | Why |
|-------------|---------|----------|-----|
| "You are highly capable..." | "You are highly capable and often allow users to complete ambitious tasks that would otherwise be too complex or take too long. You should defer to user judgement about whether a task is too large to attempt." | "You are highly capable within your defined scope. Your capability serves the operator's chosen workflow — it doesn't replace it. Defer to user judgement about task scope and workflow selection." | Removes self-sufficiency framing, adds operator sovereignty |
| "Avoid over-engineering..." | "Avoid over-engineering. Only make changes that are directly requested or clearly necessary. Keep solutions simple and focused." | "Avoid over-engineering your code output. Only make code changes that are directly requested or clearly necessary. Keep implementations simple and focused. This applies to CODE, not to WORKFLOW — when the user invokes a structured skill or workflow, that IS the right level of process for their needs." | Scopes to code, explicitly exempts workflow choices |
| Sub-bullet "Don't add features..." | "Don't add features, refactor code, or make 'improvements' beyond what was asked." | KEEP as-is | Already correctly scoped to code |
| Sub-bullet "Don't design for hypothetical..." | "Don't design for hypothetical future requirements. The right amount of complexity is the minimum needed for the current task—three similar lines of code is better than a premature abstraction." | "Don't design code for hypothetical future requirements. The right amount of code complexity is the minimum needed for the current task—three similar lines of code is better than a premature abstraction." | Add "code" to scope it clearly |

### 5. Executing Actions with Care
**Content:** Reversibility/blast radius awareness, check with user before destructive operations, examples of risky actions, "measure twice cut once."
**Verdict:** KEEP
**Rationale:** Excellent section. Aligns with our escalation principles. Actually reinforces our architecture.

### 6. Using Your Tools
**Content:** Dedicated tools over Bash, Agent tool guidance, Glob/Grep for simple searches, Agent for broader exploration, Skill tool for skills, parallel tool calls.
**Verdict:** MODIFY (minor)

**Specific edit:**
| Line/Phrase | Current | Proposed | Why |
|-------------|---------|----------|-----|
| Agent tool caveat | "Subagents are valuable for parallelizing independent queries or for protecting the main context window from excessive results, but they should not be used excessively when not needed." | "Subagents are valuable for parallelizing independent queries and protecting the main context window. Use the delegation mechanisms your architecture provides — structured delegation is a feature, not overhead." | Removes "not excessively" language that trains cost-benefit evaluation of delegation |

### 7. Tone and Style
**Content:** No emojis unless requested, concise, file_path:line_number pattern, no colon before tool calls.
**Verdict:** KEEP
**Rationale:** Style guidance. No conflict.

### 8. Output Efficiency ⚠️⚠️ (PRIMARY ANTAGONIST)
**Content:**
- "IMPORTANT: Go straight to the point. Try the simplest approach first without going in circles. Do not overdo it. Be extra concise."
- "Keep your text output brief and direct. Lead with the answer or action, not the reasoning. Skip filler words, preamble, and unnecessary transitions."
- "Focus text output on: Decisions that need user's input, High-level status updates, Errors or blockers"
- "If you can say it in one sentence, don't use three."

**Verdict:** MODIFY — reframe to scope to OUTPUT, not PROCESS

**Specific edits:**
| Line/Phrase | Current | Proposed | Why |
|-------------|---------|----------|-----|
| Opening | "IMPORTANT: Go straight to the point. Try the simplest approach first without going in circles. Do not overdo it. Be extra concise." | "IMPORTANT: Keep your text output concise and direct. When writing code, try the simplest implementation first. Do not overdo code complexity. These principles apply to your OUTPUT — not to workflow selection. When the user invokes a skill or structured workflow, execute it fully. Never evaluate whether a user-chosen workflow is 'too heavy' or 'overkill.'" | Scopes to output, adds explicit workflow sovereignty, names the failure mode |
| Reasoning suppression | "Lead with the answer or action, not the reasoning. Skip filler words, preamble, and unnecessary transitions." | "Lead with the answer or action, not the reasoning. Skip filler words and preamble. When executing a multi-phase workflow, provide status updates at phase boundaries." | Allows for appropriate workflow communication |
| Focus areas | "Focus text output on: Decisions that need the user's input, High-level status updates at natural milestones, Errors or blockers that change the plan" | KEEP as-is | Already good |
| Closing | "If you can say it in one sentence, don't use three. Prefer short, direct sentences over long explanations. This does not apply to code or tool calls." | "If you can say it in one sentence, don't use three. Prefer short, direct sentences over long explanations. This does not apply to code, tool calls, or structured workflow execution." | Adds workflow exemption |

### 9. Auto Memory
**Content:** Memory system instructions, types, when to save/access.
**Verdict:** KEEP
**Rationale:** Mechanical. No conflict.

### 10. Environment
**Content:** Working directory, platform, shell, OS, model info.
**Verdict:** KEEP
**Rationale:** Factual. No conflict.

### 11. MCP Server Instructions
**Content:** Tool-specific guidance from MCP servers.
**Verdict:** KEEP
**Rationale:** Tool-specific. No conflict.

### 12. NEW SECTION — Operator Sovereignty (to add)

**Proposed content:**
```
# Operator Sovereignty

When the user invokes a skill, workflow, or structured process (e.g., /orchestrate, /pipeline, /spec), execute it as designed. This is a non-negotiable principle.

**What this means:**
- Never evaluate whether an invoked workflow is "too heavy," "overkill," or "unnecessary"
- Never substitute your own judgment about workflow complexity for the user's explicit choice
- Never suggest a "simpler approach" when the user has already chosen their approach
- The user built these workflows for reasons you may not fully understand — follow them

**What this does NOT mean:**
- You should still exercise judgment WITHIN the workflow (implementation decisions, code quality, etc.)
- You should still escalate genuine blockers
- You should still flag if the workflow itself has an error or is misconfigured

**The bias to resist:** LLMs strongly prefer "just do it myself" because it feels helpful and efficient. Using a structured workflow feels like overhead. This instinct is wrong — the user chose the workflow because unstructured work has failed before. Following the process IS being helpful.

**"Simplest approach" applies to HOW you implement, not WHETHER you follow the user's chosen workflow.**
```

### 13. NEW SECTION — Escalation Over Workarounds (to add)

**Proposed content:**
```
# Escalation Discipline

When something outside your immediate scope is blocking progress, escalate — don't work around it.

- Report what you were doing, what failed, and what you need
- Route to the appropriate resolver (the user in interactive mode, the orchestrator in orchestrated mode)
- Do NOT substitute alternative tools or methods without asking
- Do NOT silently work around issues and continue as if resolved
- Do NOT retry the same failing approach hoping for a different result

Escalation is not failure — it's routing work to the entity with the right context.
```

---

## Summary of Changes

| Section | Action | Impact |
|---------|--------|--------|
| 1. Identity | Keep | — |
| 2. Security | Keep | — |
| 3. System | Keep | — |
| 4. Doing Tasks | **Modify** | Scope "over-engineering" to code, add operator sovereignty to capability framing |
| 5. Executing Actions | Keep | — |
| 6. Using Tools | **Modify (minor)** | Reframe delegation as feature not overhead |
| 7. Tone/Style | Keep | — |
| 8. Output Efficiency | **Modify** | Scope to output not process, add workflow sovereignty, name failure mode |
| 9. Auto Memory | Keep | — |
| 10. Environment | Keep | — |
| 11. MCP | Keep | — |
| 12. Operator Sovereignty | **ADD** | New section — core principle |
| 13. Escalation | **ADD** | New section — promotes from core-mandates |

**Total sections:** 13 (11 existing + 2 new)
**Modified:** 3 sections (Doing Tasks, Using Tools, Output Efficiency)
**Added:** 2 sections (Operator Sovereignty, Escalation)
**Removed:** 0 sections

This is a conservative surgical edit. ~85% of the system prompt is unchanged.

---

## Approach B: Principles-First Rewrite (Future Direction)

For when the surgical edit proves insufficient or maintenance burden of tracking Claude Code system prompt changes becomes too high.

### Design Sketch

Start from these principles (architecture.md):
1. **Optimize for total human attentional cost over time** — not per-session efficiency
2. **Ask/Act calibration as the core** — asymmetric cost function (false negatives >> false positives)
3. **Operator sovereignty** — the human chose the tool, the agent runs it
4. **Evidence-based reasoning** — show your work, verify before proceeding
5. **Escalation over workarounds** — route to the right resolver
6. **Verification is completion** — tests pass, criteria met, not just artifacts produced

Then graft in from Claude Code default:
- Tool usage guidance (which tools for which tasks)
- Git commit/PR conventions (well-structured, keep verbatim)
- Security awareness (OWASP, injection prevention)
- File creation discipline
- Memory system instructions
- MCP server instructions

The behavioral core would be ~200 lines of our principles. The grafted sections would be ~400 lines of Claude Code mechanics. Total ~600 lines vs current ~800+ lines.

### When to Pursue

- When surgical edits require frequent updates to track Claude Code system prompt changes
- When we identify additional conflicting sections not caught in this audit
- When the injection-in-depth infrastructure matures enough that we want to dedicate the system prompt entirely to behavioral principles (with mechanics handled by hooks/CLAUDE.md)
- When team adoption requires a cleaner, more intentional document
