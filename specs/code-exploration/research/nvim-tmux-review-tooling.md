# Building an AI-guided code review copilot in tmux + neovim

**Claude Code can absolutely control neovim programmatically, and most of the building blocks for a guided code review workflow already exist — they just haven't been assembled into this specific pattern yet.** The critical enabler is Neovim's RPC API combined with the `mcp-neovim-server`, which gives Claude Code 19 discrete tools for commanding a running neovim instance: opening files, splitting windows, navigating to lines, setting highlights, and reading editor state. The interactive "sit with me and explore code" paradigm you're describing sits in a genuine gap between automated PR review tools (CodeRabbit, PR-Agent) and passive AI chat — a gap that multiple researchers and tool-builders have identified but nobody has fully filled.

The pieces compose cleanly: neovim's `--listen` socket for control, `mermaid-ascii` and `graph-easy` for terminal-native diagrams, Claude Code's skill system for session orchestration, and a handful of neovim plugins for bidirectional communication. What's missing is the "session conductor" — the glue layer that sequences these tools into a coherent guided walkthrough experience. That's roughly **200-400 lines of custom Claude Code skills and a small neovim plugin**, not a major engineering effort.

## Claude Code can control neovim through three proven mechanisms

The most robust approach uses **Neovim's native RPC API**. Start neovim with `nvim --listen /tmp/nvim-review.sock`, and Claude Code gains full programmatic control through bash commands:

```bash
# Open file at specific line
nvim --server /tmp/nvim-review.sock --remote +42 src/auth/login.ts

# Navigate to a function by search
nvim --server /tmp/nvim-review.sock --remote-send '<Esc>/def process_data<CR>zz'

# Split window with related file
nvim --server /tmp/nvim-review.sock --remote-send '<Esc>:vsplit +10 src/models.py<CR>'

# Query what the developer is currently viewing
nvim --server /tmp/nvim-review.sock --remote-expr 'expand("%:p")'

# Read cursor position for context-aware guidance
nvim --server /tmp/nvim-review.sock --remote-expr 'line(".")'
```

The `nvr` (neovim-remote) tool simplifies the syntax: `nvr +42 src/auth/login.ts` or `nvr -cc split +42 src/utils.py`. For richer control, `pynvim` provides a full Python client that can manipulate buffers, windows, and tabs programmatically — Claude Code could invoke a small Python script for complex multi-step view arrangements.

The second mechanism is **mcp-neovim-server** (bigcodegen/mcp-neovim-server, 301 stars), which wraps neovim's RPC into **19 MCP tools** that Claude Code can call directly: `vim_command`, `vim_window` (split/vsplit/navigate), `vim_file_open`, `vim_buffer`, `vim_search`, `vim_mark`, `vim_visual`, `vim_status`, and more. Configuration is straightforward — add it to Claude Code's MCP config pointing at the neovim socket. This is the **highest-leverage integration point** because it gives Claude Code structured, typed access to neovim rather than raw shell commands.

The third option, `tmux send-keys`, works as a zero-config fallback but is fragile — it injects keystrokes blindly without mode awareness or return values. **Use RPC for structured operations and tmux send-keys only for simple fallback scenarios.**

## The closest existing implementations and why they're not quite this

**claudecode.nvim** (coder/claudecode.nvim, 2,000+ stars) is the most complete bidirectional bridge. It reverse-engineered Anthropic's VS Code extension protocol and creates a WebSocket server that Claude Code auto-discovers via a lock file at `~/.claude/ide/[port].lock`. Claude Code can open files, show diffs, and read editor state through it. It supports a `provider = "none"` mode designed for tmux-managed terminals — directly relevant to this workflow. However, it's optimized for the "Claude writes code, you review diffs" pattern, not the "Claude curates what you see while you explore" pattern.

**The walkthrough skill** (alexanderop/walkthrough) is the closest to the guided exploration concept. When invoked with `walkthrough how does auth work`, it spawns 2–4 subagents in parallel to explore code, synthesizes findings into 5–12 key concepts, and generates a self-contained HTML file with interactive Mermaid diagrams. The author's motivation — "I need to see the system before I can reason about it" — precisely matches the artisan/apprentice philosophy. The gap: it produces a static artifact, not an interactive session where Claude controls your editor view in real time.

**codecompanion.nvim** implements the **Agent Client Protocol (ACP)**, an open standard from Zed Industries for editor ↔ agent communication over JSON-RPC. Through ACP, Claude Code can request file reads/writes and receive structured editor context. This protocol is the most architecturally sound option for bidirectional communication, though it currently lacks terminal operation support (`terminal/*` capabilities return false). ACP is gaining traction — Claude Code, Codex CLI, Gemini CLI, and Cursor Agent all support it, and neovim has three ACP-compatible plugins (codecompanion, avante.nvim, agentic.nvim).

**Understand-Anything** (Lum1104/Understand-Anything) orchestrates 5 specialized agents to build a **knowledge graph** of a codebase, then exposes it through commands like `/understand-chat How does the payment flow work?` and `/understand-explain src/auth/login.ts`. It generates an interactive web dashboard with a color-coded architectural graph. This is the "codebase comprehension engine" that could power a guided review session's knowledge layer, though it currently lacks neovim integration.

## Building blocks that compose into the workflow

The toolchain for terminal-native architecture diagrams is mature. **`mermaid-ascii`** (Go binary) renders Mermaid syntax directly to ASCII/Unicode art — pipe any Mermaid diagram Claude generates to `mermaid-ascii` and it displays instantly in tmux with zero dependencies. For dependency graphs derived from actual source code, **`dependency-cruiser`** (JS/TS) and **`pydeps`** (Python) output Mermaid or Graphviz DOT format, which **`graph-easy`** converts to ASCII box art. The complete pipeline — Claude Code analyzes code → generates Mermaid → renders ASCII — works today:

```bash
# Claude generates architecture diagram, renders in terminal
echo 'graph LR; Auth[Auth Module] --> UserSvc[User Service] --> DB[(Database)]' | mermaid-ascii

# Source analysis → diagram (JS/TS)
npx depcruise src --include-only "^src/auth/" --output-type mermaid | mermaid-ascii

# Source analysis → diagram (Python)  
pydeps mypackage --noshow --show-dot | graph-easy --from=dot --as_boxart
```

For richer visualization, **image.nvim** + **diagram.nvim** can render full Mermaid/PlantUML/D2 diagrams inline in neovim buffers using the Kitty graphics protocol. This works inside tmux 3.3+ with `set -g allow-passthrough on`. Claude Code could write a Mermaid diagram to a scratch buffer, and diagram.nvim auto-renders it as a crisp visual. Alternatively, `timg` or `chafa` can display PNG output from `mmdc` (Mermaid CLI) directly in the terminal with Kitty protocol passthrough.

Claude Code's **skill system** provides the orchestration layer. Skills are markdown files in `.claude/skills/<name>/SKILL.md` with YAML frontmatter specifying allowed tools, model, and description. A `/guided-review` skill could define the entire session protocol — how Claude should set up the tmux layout, what questions to ask, how to navigate between files, when to generate diagrams. The official code-review plugin already demonstrates the pattern of launching **4 parallel review agents** with confidence thresholds and subagent validation.

The **Agent Teams** feature (experimental, v2.1.32+) natively uses tmux for split-pane orchestration — 2-16 Claude sessions coordinated with shared task lists, peer messaging, and file locking. Each teammate gets its own tmux pane. While designed for parallel coding, the infrastructure for "Claude manages multiple panes while you watch" already exists.

## What doesn't exist yet and must be built

**No neovim-native code tour tool exists.** CodeTour (VS Code) lets developers record annotated paths through a codebase as `.tour` JSON files. No neovim equivalent has been built. For this workflow, a lightweight neovim plugin that reads tour definitions (file + line + annotation) and integrates with neovim's RPC for external navigation would be the missing piece — Claude Code writes the tour plan, the plugin executes the navigation and displays annotations.

**No "session conductor" skill combines neovim control with guided exploration.** The individual capabilities exist — mcp-neovim-server for editor control, mermaid-ascii for diagrams, Claude Code's subagent system for analysis — but nothing sequences them into a coherent review session. This is the primary gap: a Claude Code skill that understands how to orchestrate a multi-step exploration where it controls what the developer sees while responding to their questions and generating contextual diagrams.

**The interactive review session pattern has no tooling.** The LACY research paper (FSE '26) validates the approach — expert-guided tours achieved **83% comprehension quiz scores versus 57% for AI-only tours** — but describes a VS Code workflow. The concept of a real-time, AI-curated code exploration session where the human reads actual code while the AI orchestrates the view is genuinely novel in the terminal tooling space. CodeRabbit's agentic chat comes closest in the automated review world, but it operates on PR diffs, not arbitrary exploration. The academic paper "AI-Guided Exploration of Large-Scale Codebases" (2025) proposes a "closed interaction loop" between structural visualization and conversational AI guidance, but remains theoretical.

**ACP needs terminal operation support.** The Agent Client Protocol could be the ideal standard for this workflow, but its current neovim implementations don't support terminal operations. Extending codecompanion.nvim's ACP capabilities to include `terminal/execute` would allow Claude Code to orchestrate both editor views and terminal output (test results, build logs, diagram rendering) through a single protocol.

## Proposed architecture for a tmux + neovim + Claude Code guided review session

The system has four layers that compose cleanly:

The **session layer** is a Claude Code skill (`/review-session`) that initializes the tmux layout and manages session state. It creates a three-pane layout: neovim (left, 60%), Claude Code conversation (top-right, 25%), and a scratch/diagram pane (bottom-right, 15%). It launches neovim with `--listen /tmp/nvim-review.sock` and registers the mcp-neovim-server pointing at that socket.

The **navigation layer** uses mcp-neovim-server's 19 tools for all editor control. Claude Code opens files, splits windows, jumps to functions via LSP or tag navigation, highlights risk areas with `nvim_buf_add_highlight`, and reads the developer's cursor position to stay context-aware. When the developer manually jumps to a file, Claude queries `vim_status` to detect the new context and adjust its commentary.

The **analysis layer** runs pre-exploration to build a mental model. Before the session starts, Claude Code spawns 2–3 subagents: one to map the architecture (generating a dependency graph via dependency-cruiser or pydeps), one to identify recent changes and risk areas (via `git log` analysis), and one to catalog key abstractions and entry points. These findings feed into the session's knowledge base, which Claude draws on throughout the walkthrough.

The **visualization layer** generates diagrams on demand. When Claude wants to show architecture, it generates Mermaid syntax, pipes it through `mermaid-ascii`, and displays it in the scratch pane via `tmux send-keys`. For richer diagrams, it writes Mermaid to a neovim scratch buffer where diagram.nvim renders inline. For call graphs of specific modules, it runs dependency-cruiser or pyan3 and converts the output to ASCII via graph-easy.

A typical session flow would look like this: The developer types `/review-session walk me through the auth refactor from last sprint`. Claude spawns analysis subagents to identify all changed files in the auth module over the last two weeks, maps their dependency relationships, and identifies architectural risk areas. It opens the first key file in neovim at the critical function, displays an ASCII architecture diagram of the auth module in the scratch pane, and says: "Here's the entry point — `authenticateUser` on line 42 of auth/handler.ts. It now delegates to three new services instead of handling everything inline. Let me show you each one." The developer reads the code, asks "why did they separate token validation?", and Claude navigates to the token service while explaining the design rationale it inferred from commit messages and code comments.

The minimal implementation requires: installing mcp-neovim-server and mermaid-ascii, writing a `/review-session` skill (roughly 100–200 lines of markdown defining the session protocol), and configuring neovim to launch with a known socket path. The claudecode.nvim plugin from Coder can complement this setup for bidirectional awareness. **Total effort estimate: a weekend of focused work to reach a functional prototype, with iteration from there based on what feels right in practice.**

## Conclusions and what makes this approach genuinely new

The guided code review copilot pattern occupies a real gap in the tooling landscape. Current AI code tools cluster into two modes: fully automated review (CodeRabbit produces a report) or fully manual with AI as a chat sidebar (Cursor, Cody). The fusion mode — where AI controls the developer's view while the developer retains full agency to navigate, question, and override — has been theorized in academic work but not built for terminal-first workflows. The ADHD-informed design insight is particularly sharp: the session needs to be **high-stimulus and interactive** (Claude actively managing the view, generating diagrams, posing questions) rather than passive (reading a static report or scrolling through files alone).

The artisan/apprentice metaphor maps directly to how mcp-neovim-server's bidirectional capability should be used. Claude doesn't just answer questions — it **curates the view**, saying effectively "look at this next" while the developer maintains the authority of what to study deeply. The closest analogy in existing tools is the Understand-Anything plugin's knowledge graph combined with mcp-neovim-server's editor control, orchestrated by a custom session skill. None of these exist as an integrated package today, but they compose naturally because they all speak either MCP, neovim RPC, or both. The most impactful first step is installing mcp-neovim-server and writing a `/review-session` skill that uses its tools — everything else (diagrams, pre-analysis, session persistence) can layer on incrementally.
