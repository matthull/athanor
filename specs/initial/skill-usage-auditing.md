# Skill Usage Auditing — Research Notes

**Date:** 2026-03-31
**Source:** orthogonality opus, interactive session with artifex
**Status:** Research findings, not yet actioned

## Problem

No visibility into which skills agents actually load and use. Can't tell if routing tables are working, if skills are being discovered, or which skills provide value in practice.

## Current Mitigations

- **Discharge convention (bandaid):** `opus.md` Reflection section now asks azers to note which skills they loaded. Relies on agent compliance — will drift.

## Hook-Based Approach (Recommended)

The Claude Code hook system supports `PostToolUse` hooks with a `matcher` field that can match on tool name. The `Skill` tool is invoked every time an agent loads a skill via `/skill-name`.

### Feasibility

Confirmed in `~/.claude/settings.json`:
- `PostToolUse` hooks already exist (e.g., `Edit|Write` matcher for session log nudge)
- `Skill(name)` patterns already appear in the allowlist, confirming the tool name format
- A `Skill` matcher on PostToolUse would fire every time any skill is loaded

### Design

```json
{
  "matcher": "Skill",
  "hooks": [
    {
      "type": "command",
      "command": "~/.claude/hooks/skill-usage-log.sh"
    }
  ]
}
```

The hook script would:
1. Receive tool input (includes skill name) via stdin or environment
2. Append `timestamp | session_id | skill_name` to a log file (e.g., `~/egregore/activity/skill-usage.log`)
3. Exit 0 (non-blocking)

### What the hook receives

Need to verify: what data does a PostToolUse hook receive? Likely the tool input JSON which for Skill calls includes `{"skill": "calcinatio"}`. The hook script would parse the skill name from this.

### Open Questions

- What exact format does PostToolUse pass to the hook command? (stdin? env vars? args?)
- Should the log be per-session or global?
- Should it feed into the session log via egregore-mcp, or stay as a separate lightweight file?
- Is there value in also logging which skills are loaded at session start vs. mid-session?

### Value

- Track which skills are actually picked up by agents
- Identify routing table gaps (skills that should be loaded but aren't)
- Identify dead skills (never loaded by anyone)
- Correlate skill usage with opus quality over time
