# Architecture Decisions: 2026-03-14

Decisions and agreements from the Gastown analysis + autonomy architecture discussion.

---

## Agreements (Decided)

### A1: Adopt GUPP "Physics Not Politeness" Framing
**Decision:** Add propulsion metaphor to autonomous mode injections.
**Rationale:** Clarifies intent for human system directors. Low cost, potential upside for LLM compliance. "You are executing a workflow. If a phase is incomplete, execute it. This is not optional."
**Action:** Update `autonomy-profile.txt` templates with GUPP-style language.

### A2: Session-Surviving Identity — Plan Doc Is Sufficient
**Decision:** No investment in formal agent identity system (like Beads agent beads).
**Rationale:** Plan doc updated at least once per phase serves as durable identity. Agents figure out where to pick up. Refine reactively as issues arise.
**Action:** None needed now.

### A3: Beads — Watch, Don't Adopt
**Decision:** Do not adopt Beads for work state tracking.
**Rationale:** At 1-3 agent scale, markdown + git is sufficient. Beads' value inflects at 5-10+ agents. Adoption cost (8-12 hours) not justified by marginal benefit.
**Revisit when:** Agent count crosses ~5 concurrent, or cross-repo dependencies emerge.
**Action:** None now. Monitor Beads development.

### A4: Machine-Readable Workflow State — Convention, Not Schema
**Decision:** Don't invest in structured schema (TOML/YAML). Use incremental checklist convention in markdown results blocks.
**Rationale:** Gets the main benefit (partial completion tracking within phases) without infrastructure change. POC-able on next `/orchestrate` run.
**POC verification:** Agents resume cleanly after compaction without redoing completed work units.
**Action:** Update `/orchestrate` SKILL.md to document incremental results convention.

Example:
```markdown
### IMPLEMENT
**Results:**
- [x] WU-1: Add notification preferences model (commit: abc123)
- [x] WU-2: Slack App Home UI (commit: def456)
- [ ] WU-3: Preference persistence API
- [ ] WU-4: Integration tests
Confidence: (pending)
Gate outcome: (pending)
```

### A5: Orchestrate vs NDI — ~80% Equivalent
**Decision:** Our phase system is functionally equivalent to Gastown's NDI for our scale. Key gaps (partial completion, eager state persistence) fixable with conventions (A4).
**The real gap:** Not "how do you resume?" but "how do you resume without a human initiating it?" That's GUPP — and that's what the autonomy architecture addresses.
**Action:** Implement A4 convention. Auto-resumption addressed by agent roles + GUPP bridge in autonomy architecture.

---

## Decisions Pending

### D1: Custom Permission System — Tiered by Role (BLOCKER RESOLVED)
**Options:**
1. `dangerouslySkipPermissions` + custom hooks for ALL roles (replace builtin entirely)
2. Layer custom gating ON TOP of builtin (belt and suspenders) — different tiers per role
3. Keep builtin, add specific gating hooks for high-risk operations only

**Current leaning:** Option 2 with sandboxing insight — sandboxable roles (implementer, fixer) get sandbox + `dangerouslySkipPermissions` + hooks. Non-sandboxable roles (supervisor, browser-qa) get builtin + hooks.

**Blocker RESOLVED (2026-03-14):** Sandbox blocks Docker by default, BUT `excludedCommands: ["docker", "docker compose"]` cleanly exempts Docker from sandbox while keeping everything else sandboxed. This means implementer agents CAN be sandboxed — Docker runs outside sandbox with hook-based gating, all other commands stay OS-level sandboxed.

**Sandbox findings (from Claude Code docs investigation):**
- Filesystem: writes restricted to cwd + configured `allowWrite` paths. OS-level enforcement (bubblewrap on Linux).
- Network: blocked by default, allowed domains configurable. Localhost blocked unless explicitly allowed.
- Docker: explicitly incompatible — use `excludedCommands` to exempt.
- Git: works fine in sandbox.
- `dangerouslySkipPermissions` is orthogonal to sandbox — sandbox provides OS-level enforcement, skip-perms removes prompt friction. Can combine safely.

**Recommended config for implementer/fixer role:**
```json
{
  "sandbox": {
    "enabled": true,
    "excludedCommands": ["docker", "docker compose"],
    "filesystem": {
      "allowWrite": ["//tmp/**"]
    }
  },
  "defaultMode": "bypassPermissions"
}
```

Docker commands then gated by PreToolUse hooks (no-push-to-master, no-merge-pr, etc.) instead of builtin permissions.

**Action:** Ready to prototype. Build for fixer role as first POC.

### D2: Supervisor Scope — MVP vs Full
**Options:**
1. MVP: tmux health check + timeout enforcement + Telegram escalation (small tool)
2. Full: Hierarchical supervision with patrol loops, learning from patterns, cross-skill coordination

**Current leaning:** MVP first. Solves the specific pain (agent goes off course, operator has to notice manually). Full version when we have more concurrent agents.

**Action:** Build MVP supervisor as part of Phase D of AppSignal rollout.

---

## Open Threads (Parked for Future Discussion)

### T1: Agent Roles as Foundation
Strong connective tissue — roles compose naturally with permission system, gating hooks, supervisor behavior, and injections. Foundation is trivial (JSON files + spawn script). But full integration with wtp, session injection, and supervisor is more work.

**Next step:** Define the JSON schema. Build for one role (fixer) as POC.

### T2: Code Coverage Gating
Concrete example of executable gating. Feasibility depends on having fast coverage checks. Would need to run coverage diff as part of the hook, which might be too slow for PreToolUse.

**Alternative:** Gate at PR creation time (check coverage in CI, block merge if regressed). Less novel but more practical.

### T3: Custom Permission Go Tool
If we go with Option 2 (tiered), the non-sandboxable roles need fine-grained permission logic. A Go tool (`perms-gate`) called by hooks would be fast and testable. Shell scripts work for POC but Go is better for production.

**Defer until:** We've validated the tiered approach with shell script hooks first.

---

## Research Findings (Reference)

### Gastown Alignment Mechanisms (Ranked by Relevance to Us)

1. **Executable gating via hooks** — HIGH relevance. Strictly better than prompting for hard constraints.
2. **GUPP propulsion framing** — MEDIUM relevance. Clarifies intent, uncertain LLM compliance benefit.
3. **Molecules (externalized workflow)** — LOW-MEDIUM relevance. Our task-plan.md + convention change gets most of the benefit.
4. **Hierarchical supervision** — MEDIUM relevance. Solves a real pain point, but start with flat supervisor.
5. **Persistent agent identity** — LOW relevance. Plan doc is sufficient at our scale.
6. **Beads (structured work state)** — LOW relevance at current scale. Revisit at 5+ agents.
7. **Wisps (ephemeral orchestration state)** — LOW relevance. `/tmp/` files suffice for now.
8. **Nondeterministic Idempotence** — MEDIUM relevance. Convention change (A4) closes most of the gap.
9. **Multi-phase prompt composition** — LOW relevance. Our session injection system is already similar.
10. **Session compaction optimization** — LOW relevance. Our PreCompact hooks handle this.
