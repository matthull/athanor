# Code Exploration Tooling Specification

**Status:** Draft
**Created:** 2026-04-08
**Research:** `specs/code-exploration/research/`

---

## Traceability Legend

| Tag | Meaning |
|-----|---------|
| `[R:nvim-tmux]` | Research: nvim-tmux-review-tooling.md |
| `[R:practitioners]` | Research: advanced-review-practitioners.md |
| `[R:design]` | Research: interactive-exploration-design.md |
| `[D:reason]` | Decision made during spec |
| `[E:existing]` | Extends existing pattern |
| `[U:topic]` | Unbacked — needs validation |

---

## Overview

Build an AI-guided interactive code exploration system for terminal-first workflows (tmux + neovim + Claude Code). The system serves Pipeline 2 from the design research — eventual comprehension decoupled from shipping. The developer explores code with a knowledgeable AI guide that curates views, generates diagrams, and provides context, while the developer reads actual code with their own eyes. `[R:design]`

**Current state:** GitHub PR review and asking an azer to explain things.

**Target state:** A `/explore` skill that pre-analyzes a codebase area, generates architecture diagrams, and guides interactive exploration sessions — with individual tools validated and composable independently.

**Primary goal:** Holistic comprehension — the developer deeply understands how a chunk of the system works, how the pieces connect, and why it's shaped the way it is. Identifying problem areas is a secondary benefit, not the driver. The tool serves comprehension, not audit. `[D:artifex-direction]`

---

## Domain Dictionary

| Term | What it is |
|------|-----------|
| **Exploration session** | A focused, interactive session where the developer explores code guided by AI. Not review — exploration for comprehension. `[R:design]` |
| **Pre-analysis** | Automated reconnaissance before the session starts — mapping files, dependencies, recent changes, generating diagrams. `[R:design]` |
| **Session conductor** | The glue layer that sequences individual tools into a coherent exploration. The `/explore` skill. `[R:nvim-tmux]` |
| **Terminal-native diagram** | Architecture/dependency diagrams rendered as ASCII/Unicode art in the terminal, not requiring a browser or image viewer. `[R:nvim-tmux]` |
| **Harness calibration** | Using exploration findings to improve the autonomous shipping pipeline — the feedback loop between Pipeline 2 and Pipeline 1. `[R:design]` |

---

## Tiers

### Tier 1: Individual Tool Stack

Install, configure, and validate each tool independently with integration tests before composing them. Each tool must provide standalone value.

#### Tool: mermaid-ascii

**What:** Go binary that renders Mermaid syntax to ASCII/Unicode art in the terminal. `[R:nvim-tmux]`

**Why:** Terminal-native architecture diagrams without browser/image dependencies. Usable in any Claude Code session, any tmux pane.

**Validation test:** Generate 2-3 diagrams of varying complexity (flowchart, sequence, class) and confirm readable output.

**Validation results:** Installed via `go install github.com/AlexanderGrooff/mermaid-ascii@latest`. Small branching diagrams (3-7 nodes) render cleanly. Subgraphs and long chains degrade — nodes duplicate, layout flattens. Large comprehensive diagrams are not viable.

**Design constraint:** Diagrams must be focused illustrations, not comprehensive maps. Each diagram answers one question with 3-7 nodes. If you feel inclined to generate a diagram with more than 7 nodes, stop — instead, break it into multiple smaller diagrams embedded in textual narrative with code links (file:line references). The narrative carries the big picture; the diagrams illustrate specific relationships within it. This constraint plays to mermaid-ascii's strengths and produces better comprehension than any single large diagram could. `[D:small-focused-diagrams]`

**Status:** Installed, validated, included with constraints

#### Tool: difftastic

**Status:** Deferred `[D:focus-on-what-exists]`

Structural diff tool. Serves "what changed?" — not the primary question. Revisit when exploration sessions need evolution/history understanding.

#### Tool: CodeScene MCP Server

**Status:** Deferred `[D:focus-on-what-exists]`

Code health scoring. Serves "how healthy is this?" — audit-oriented, not comprehension-oriented. Revisit for Pipeline 1 merge gates or when exploration needs quantitative health data.

#### Tool: dependency-cruiser

**Status:** Deferred `[D:claude-traces-imports]`

JS/TS only. Claude can trace imports for any language during pre-analysis. Revisit if pre-analysis import tracing proves insufficient.

#### Tool: treesitter-based analyzers (Serena, etc.)

**Status:** Deferred `[D:mixed-reports]`

Mixed real-world utility reports. Revisit if Claude's native code reading proves insufficient for structural understanding.

#### Tool: nvr (neovim-remote)

**What:** Lightweight CLI for sending commands to a running neovim instance via its RPC socket. `[D:simpler-than-mcp]`

**Why:** Enables Claude to open files, navigate to lines, create splits, and add annotations in the developer's neovim — turning exploration from "Claude describes code" to "Claude shows you code."

**Validation results:** Installed via `pipx install neovim-remote`. All core operations validated: open file at line (`--remote +LINE`), navigate (`-c 'LINE'`), vsplit (`-c 'vsplit +LINE file'`), query state (`--remote-expr`), virtual text annotations (extmarks via lua), floating windows for longer explanations. All work reliably over the `/tmp/nvim-explore.sock` socket.

**Status:** Installed, validated, included.

#### Tool: neovim annotations (extmarks + floating windows)

**What:** Neovim's native extmark API for inline virtual text and floating windows, driven via nvr. `[D:tested-and-validated]`

**Why:** Enables Claude to annotate specific lines with brief explanations (virtual text at end of line) and present longer explanations as floating windows positioned near relevant code. Three highlight groups: DiagnosticInfo (blue, general), DiagnosticWarn (yellow, key points), DiagnosticHint (teal, supplementary).

**Design constraint:** Max 3-5 annotations per file. Clear annotations before moving to a new file. Close floating windows after the developer has read them. `[D:annotation-discipline]`

**Status:** Validated, included.

#### Layout: Side-by-side (Claude left, neovim right)

**What:** Vertical tmux split — Claude Code on the left (~35% width), neovim on the right (~65% width). Both get full terminal height.

**Why:** Developer reads code on the right, conversation on the left. Both visible simultaneously without switching. Claude Code pane is narrow enough to encourage concise messages but wide enough to be readable. `[D:tested-layout]`

**Setup command:**
```bash
tmux split-window -h -l 65% "nvim --listen /tmp/nvim-explore.sock $(pwd)"
```

**Status:** Validated, included.

#### Tool: mcp-neovim-server

**Status:** Deferred `[D:nvr-sufficient]`

Full MCP integration with 19 tools. nvr covers all current needs. Revisit if nvr proves limiting.

---

### Tier 2: The `/explore` Skill

**What:** A Claude Code skill that orchestrates pre-analysis and interactive guided exploration.

**Invocation:**
```
/explore "walk me through the seismic classifier mapping in musashi"
/explore "what changed in the API layer this sprint"  
/explore "deep dive into the auth module"
```

**Responsibilities:**
1. **Pre-analysis phase** — spawn subagents to map the target area:
   - File inventory and structure
   - Dependency relationships
   - Recent changes (git log analysis)
   - Architecture diagram generation (mermaid-ascii)
   - Key entry points and abstractions identified

2. **Briefing** — present findings to the developer:
   - ASCII architecture diagram
   - Summary of area scope, key files, recent changes
   - Suggested exploration path (developer can override)

3. **Interactive exploration** — respond to developer's questions with full context:
   - Navigate related files on request
   - Generate diagrams on the fly
   - Explain design decisions from commit messages and code structure
   - Highlight pattern divergences or interesting areas
   - Track what's been explored vs. what remains

4. **Wrap-up** — capture session output:
   - Summary of what was explored
   - Harness improvement ideas (if any)
   - Open questions for follow-up

**Design constraints:**
- Works on any codebase — not tied to a specific language or project `[D:generality]`
- Pre-analysis should complete in seconds, not minutes `[D:adhd-friendly]`
- Developer leads, Claude follows — but Claude can suggest "look at this next" `[R:design]`
- No neovim control required at Tier 2 — all via conversation `[D:incremental]`

**Validation test:** Run a real exploration session on the musashi project's seismic classifier mapping solution. The session should produce genuine comprehension, not just a report.

---

### Tier 3: Guided Neovim Sessions (stretch)

**What:** Claude controls neovim via mcp-neovim-server during exploration sessions.

**Additions over Tier 2:**
- Claude opens files and navigates to specific lines in the developer's neovim
- Split windows showing related files
- Diagrams rendered in scratch buffers
- Cursor position tracking for context awareness

**Status:** TBD — pursue if Tiers 1 and 2 go well.

---

## Workflow

```
Developer invokes: /explore "<topic>"
         │
         ▼
┌─────────────────┐
│  Pre-analysis   │  Subagents fan out:
│  (parallel)     │  - file mapping
│                 │  - dependency graph
│                 │  - git history
│                 │  - diagram generation
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Briefing       │  Present: diagram, summary,
│                 │  suggested path
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Interactive     │  Developer asks questions,
│  exploration    │  Claude provides context,
│  (loop)         │  generates diagrams,
│                 │  navigates codebase
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Wrap-up        │  Summary, harness ideas,
│                 │  open questions
└─────────────────┘
```

---

## Test Scenarios

| Scenario | What to verify |
|----------|---------------|
| mermaid-ascii renders a flowchart | Readable ASCII output from Mermaid syntax |
| mermaid-ascii renders a complex diagram | Handles multiple nodes, edges, subgraphs |
| difftastic shows structural diff | Ignores formatting, highlights logic changes |
| `/explore` pre-analysis completes in <30s | Subagents return useful results quickly |
| `/explore` generates accurate architecture diagram | Diagram reflects actual code structure |
| `/explore` handles unfamiliar codebase | Works on a project the developer hasn't seen |
| `/explore` session on musashi seismic classifier | Produces genuine comprehension of the mapping solution |

---

## Out of Scope

- Pipeline 1 (autonomous shipping) improvements — separate work
- Auto-merge with quality gates — separate from exploration tooling
- Comprehension delta map — deferred per design doc `[R:design]`
- Full CodeTour-style recorded tours — premature `[R:nvim-tmux]`
- Browser-based visualization — terminal-first only

---

## Retrospective

When this feature is complete, review:

### Tool Stack
- [ ] Which tools provided genuine value vs. added complexity?
- [ ] Any tools we should have included but didn't?
- [ ] Installation/configuration friction points

### Skill Design
- [ ] Did pre-analysis produce useful context or just noise?
- [ ] What was the right balance of Claude-leads vs. developer-leads?
- [ ] How long did useful sessions last? Were they time-boxed?
- [ ] Did the skill work across different codebase types/languages?

### Workflow
- [ ] Did exploration sessions produce harness improvement ideas?
- [ ] Was the feedback loop to Pipeline 1 natural or forced?
- [ ] What would make this a daily tool vs. occasional?

---

## Product FYI

N/A — single-developer tooling project, no product team.
