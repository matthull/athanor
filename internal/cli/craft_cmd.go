package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runCraft(args []string) int {
	var (
		dir    string
		model  string
		intent string
	)

	positional, flagArgs := splitArgs(args)

	if len(positional) < 2 {
		fmt.Fprintln(os.Stderr, "error: athanor name and session name required")
		fmt.Fprintln(os.Stderr, "usage: ath craft <athanor> <session-name> [<mo-name>] [--dir <path>] [--model <model>] [--intent <text>]")
		return 2
	}
	athName := positional[0]
	sessionName := positional[1]

	// Optional: mo-name as third positional
	var moName string
	if len(positional) >= 3 {
		moName = positional[2]
	}

	fs := flag.NewFlagSet("craft", flag.ContinueOnError)
	fs.StringVar(&dir, "dir", "", "working directory")
	fs.StringVar(&model, "model", "", "model override")
	fs.StringVar(&intent, "intent", "", "natural language intent for the session")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	instDir := athanor.InstanceDir(home, athName)
	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Resolve MO for opera path — if no moName given, try to find a single MO
	if moName == "" {
		mos, _ := athanor.ListMagnaOpera(instDir)
		if len(mos) == 1 {
			moName = mos[0]
		}
	}

	// Create the opus file — or defer if no MO exists yet (bootstrap mode)
	datestamp := time.Now().Format("2006-01-02")
	opusFilename := fmt.Sprintf("%s-%s.md", datestamp, sessionName)
	bootstrap := moName == ""

	var opusPath string
	if !bootstrap {
		operaDir := athanor.OperaPath(instDir, moName)
		if err := os.MkdirAll(operaDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating opera dir: %v\n", err)
			return 1
		}
		opusPath = filepath.Join(operaDir, opusFilename)

		// Build opus content
		intentText := "(Shape collaboratively with the artifex)"
		if intent != "" {
			intentText = intent
		}

		frontmatter := fmt.Sprintf("---\nstatus: charged\ninscribed: %s\ninteractive: true\n", datestamp)
		frontmatter += fmt.Sprintf("magnum_opus: %s\n", moName)
		frontmatter += "---\n"

		opusContent := frontmatter + fmt.Sprintf("# %s\n\n## Intent\n\n%s\n\n## Boundary\n\n- **Agent:** Interactive — work alongside the artifex\n- **Operator:** Present and directing\n\n## Context\n\n(Add context as the session progresses)\n", sessionName, intentText)

		if err := os.WriteFile(opusPath, []byte(opusContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error creating opus: %v\n", err)
			return 1
		}
	}

	// Resolve working directory
	if dir == "" {
		dir = cfg.Project
	}
	if dir == "" {
		h, _ := os.UserHomeDir()
		dir = h
	}

	// Resolve model
	if model == "" {
		model = cfg.EffectiveAzerModel()
	}

	// Build boot prompt
	var bootPrompt string
	if bootstrap {
		// No MO exists — bootstrap mode
		intentClause := ""
		if intent != "" {
			intentClause = fmt.Sprintf(" The artifex's intent: %s.", intent)
		}
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s/azer.md. You are an interactive azer working alongside the artifex in a bootstrap session — no Magnum Opus exists yet.%s The artifex will guide the session. When an MO is created during this session, create an opus file under its opera/ directory with the session name %q so discharge can be recorded. If no MO is created, the session's charge is accepted as lost — the artifex was present and nothing is lost that isn't in their head.",
			instDir, instDir, intentClause, opusFilename,
		)
	} else {
		moPath := athanor.MagnumOpusPath(instDir, moName)
		intentClause := ""
		if intent != "" {
			intentClause = fmt.Sprintf(" The artifex's intent: %s.", intent)
		}
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s, then read %s/azer.md. You are an interactive azer working alongside the artifex. Your opus is at %s — read it, then shape the intent collaboratively.%s The artifex will guide the session. When done, the artifex will signal discharge.",
			instDir, moPath, instDir, opusPath, intentClause,
		)
	}

	crucName := "azer-" + sessionName
	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %s %q",
		dir, instDir, model, bootPrompt,
	)

	// Create tmux window and launch
	r := tmux.NewRunner()
	if err := r.NewWindow(crucName, dir); err != nil {
		fmt.Fprintf(os.Stderr, "error creating crucible: %v\n", err)
		return 1
	}

	if err := r.SendKeysLiteral(crucName, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error launching azer: %v\n", err)
		return 1
	}
	if err := r.SendKeys(crucName, "Enter"); err != nil {
		fmt.Fprintf(os.Stderr, "error launching azer: %v\n", err)
		return 1
	}

	fmt.Printf("Craft session started in crucible %q\n", crucName)
	if bootstrap {
		fmt.Println("  Mode: bootstrap (no MO — opus deferred until MO is created)")
	} else {
		fmt.Printf("  Opus: %s\n", opusPath)
	}
	if intent != "" {
		fmt.Printf("  Intent: %s\n", intent)
	}
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", dir)

	return 0
}
