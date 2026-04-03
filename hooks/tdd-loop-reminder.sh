#!/bin/bash
# PostToolUse hook (Edit|Write): reminds about TDD empirical loop discipline
# when code files are modified. Reference: azer.md § Calcinatio Applied: Code Work.
#
# Throttled to once per 120 seconds per session to avoid noise.

INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

[ -z "$FILE_PATH" ] && exit 0

# Match common code file extensions
if [[ "$FILE_PATH" =~ \.(rb|py|js|jsx|ts|tsx|go|rs|java|kt|swift|c|cpp|h|hpp|cs|ex|exs|erl|hs|ml|clj|scala|sh|bash|zsh|lua|r|R|pl|pm|php|vue|svelte|dart|zig|nim|cr|v|sql|graphql|gql|proto|tf|hcl)$ ]]; then

  SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // empty')
  [ -z "$SESSION_ID" ] && exit 0

  SESSION_DIR="/tmp/claude-sessions/${SESSION_ID}"
  mkdir -p "$SESSION_DIR"

  LAST_FILE="${SESSION_DIR}/last-tdd-reminder"
  NOW=$(date +%s)
  LAST=$(cat "$LAST_FILE" 2>/dev/null || echo 0)
  ELAPSED=$((NOW - LAST))

  if [ "$ELAPSED" -ge 120 ]; then
    cat <<'ENDJSON'
{
  "hookSpecificOutput": {
    "hookEventName": "PostToolUse",
    "additionalContext": "TDD LOOP: Red \u2192 green \u2192 refactor. Maintain an empirical loop \u2014 never commit code you haven't observed working. Define success empirically before or alongside writing the implementation (azer.md \u00a7 Calcinatio Applied: Code Work)."
  }
}
ENDJSON
    echo "$NOW" > "$LAST_FILE"
  fi
fi
