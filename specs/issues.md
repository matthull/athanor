# Athanor — Known Issues

Tracked issues discovered in operation. Move to GitHub Issues when the repo is set up.

---

## Whisper: Long strings with quotes trigger permissions prompts

**Discovered:** 2026-03-24
**Severity:** Annoying — blocks automated muster/nudge flows

When `whisper send` (now `ath whisper send`) delivers a long message containing quotes, the bash command triggers a permissions prompt in Claude Code. The bash-allowlist-coach hook sees the complex quoted string and can't whitelist it.

**Reproduction:** Any `whisper send <target> "<long message with quotes and special chars>"` where the message is multi-sentence with embedded quotes.

**Example that triggered it:**
```
whisper send azer-e2e-pr "STOP waiting for the API key. New direction from the artifex: no standalone Anthropic API key available (Max account, not API access). Redesign the compose workflow to use Claude Code where possible — make it a Claude Code skill or use 'claude -p ...' CLI subprocess instead of direct API calls..."
```

**Workaround:** Use `-f` flag to send from a temp file instead of inline string. The azer.md already has guidance about writing PR bodies to temp files for similar reasons (`gh pr create` heredoc rejection).

**Proper fix candidates:**
1. Have `ath whisper send` automatically use temp file delivery when message exceeds a length threshold
2. Teach the bash-allowlist-coach to recognize `whisper send` / `ath whisper send` as safe regardless of payload
3. Both — auto-tempfile for robustness + allowlist for ergonomics

---

## Context budget guard: false alarm at 88% followed by 5%

**Discovered:** 2026-03-24
**Severity:** Confusing — triggers premature discharge behavior

The context budget guard hook reported 88% context usage, then on the very next message reported 5%. Autocompact is off, so compression is not an explanation. The cause is unknown.

**Impact:** The artifex nearly discharged a session with plenty of context remaining. False urgency undermines trust in the signal.

**Note:** The context measurement system is known to be inaccurate — it's best-effort. Fixing this specific issue may mean exploring a better overall approach to context measurement rather than patching the current hook.
