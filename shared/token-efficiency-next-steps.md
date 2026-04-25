# Token Efficiency — Next Steps

Companion to the analysis in `athanor-architect/magna-opera/collaborative-workshop/opera/2026-04-25-efficiency-findings.md`.

## Done (this session)

- Prose compression of AGENTS.md (25%), azer.md (45%), opus.md (9%). ~6,450 tokens saved.
- Restored full Tree of Work as compact sephirotic map (critic had dropped supernals).
- Calcinatio SKILL.md prose compression (21%). ~800 tokens saved.
- Marut.md Assessment Opera tightening + added lost "assessment vs. retrieval" heuristic. ~75 tokens saved.
- Confirmed calcinatio resource files are ALREADY on-demand (the 4,300 token win was pre-existing).
- Removed dyad minimum (confirmed obsolete by artifex).
- Total: **~7,250 tokens saved across shared files (~27% reduction).**

## Remaining Work

### 1. Git/PR protocols → skills (~1,300 tokens from system prompt)

The custom system prompt (`~/.claude/prompts/system-prompt.md`, ~6,400 tokens) includes full git commit protocol (~800 tokens) and PR creation protocol (~500 tokens). These are loaded on every session but used <1% of the time.

**Change:** Extract into `/commit` and `/pr` skills that load the protocol when invoked. The system prompt keeps a one-liner: "Use /commit and /pr skills for git operations."

**Risk:** Low. These are procedural, not behavioral. The procedure only needs to be in context when executing.

### 2. System prompt Behavioral Core — keep as-is (decision)

The Behavioral Core section (~1,500 tokens) overlaps with AGENTS.md geas system. Since all sessions are now athanor sessions, this is technically redundant. However: **the Core Mandates are the things that drift most in long sessions.** Having them in the system prompt, AGENTS.md, AND the periodic behavioral hook reminders provides triple reinforcement. The redundancy is earning its tokens through reliability, not unique information. Leave as-is.

### 3. Move Assessment Opera content to assessor JOB.md

AGENTS.md now has a one-line pointer. The detailed content (organic work model, assessment vs. retrieval, two layers of value extraction) should be relocated to the assessor job definition. The assessor JOB.md already covers "when to assess" and the full process, so this is mostly about the organic work model description.

### 4. Vocabulary table compression (deferred — review with artifex)

Some entries are long (Corpus Map ~100 words, Job ~70 words). Reviewed together and decided the gain (~300 tokens) isn't worth the risk of losing embedded behavioral guidance. Revisit if a future audit finds agents misusing corpus maps or job definitions.

### 5. MO file compression (per-athanor, varies)

MO files vary in size. Not part of shared infrastructure. Worth considering for large MOs with accumulated Tempering sections.

## Impact Summary

| Change | Tokens | Status |
|--------|--------|--------|
| Prose compression (AGENTS, azer, opus) | ~6,450 | Done |
| SKILL.md prose compression | ~800 | Done |
| Marut.md tightening | ~75 | Done |
| Calcinatio resources on-demand | 0 (already was) | N/A |
| Git/PR → skills | ~1,300 | Next |
| Assessment content → assessor JOB | ~0 (relocation) | Next |
| System prompt Behavioral Core | ~0 (keep) | Decision |
| Vocab table | ~0 (deferred) | Decision |
| **Total saved** | **~7,250** | |
| **Total remaining** | **~1,300** | |

Boot impact: controllable files went from ~30k to ~23k tokens. Full boot ~50k → ~43k. Boot as % of useful life: 33% → 29%.
