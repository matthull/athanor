package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

func runMuster(args []string) int {
	var (
		worktreePath string
		model        string
		crucName     string
		athName      string
		intent       string
		job          string
		formula      string
	)

	positional, flagArgs := splitArgs(args)

	fs := flag.NewFlagSet("muster", flag.ContinueOnError)
	fs.StringVar(&worktreePath, "worktree-path", "", "working directory for the azer — use for git worktrees or other workspace outside the athanor's project dir")
	fs.StringVar(&model, "model", "", "model override for the azer")
	fs.StringVar(&crucName, "name", "", "crucible name override")
	fs.StringVar(&athName, "athanor", "", "athanor name (if $ATHANOR not set)")
	fs.StringVar(&intent, "intent", "", "natural language intent (launches autonomous azer)")
	fs.StringVar(&job, "job", "", "job role for the azer (overrides opus frontmatter in opus mode, required source in intent mode)")
	fs.StringVar(&formula, "formula", "", "collaboration formula (overrides opus frontmatter in opus mode; optional in intent mode)")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	// Dispatch based on mode: --intent means intent-driven autonomous azer
	if intent != "" {
		return runMusterIntent(positional, intent, worktreePath, model, crucName, athName, job, formula)
	}
	return runMusterOpus(positional, worktreePath, model, crucName, athName, job, formula)
}

// runMusterOpus launches an azer for a pre-inscribed opus file.
// Usage: ath muster <opus-file> [--worktree-path <path>] [--model <model>] [--athanor <name>]
func runMusterOpus(positional []string, worktreePath, model, crucName, athName, job, formula string) int {
	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: opus file required")
		fmt.Fprintln(os.Stderr, "usage: ath muster <opus-file> [--worktree-path <path>] [--model <model>]")
		fmt.Fprintln(os.Stderr, "       ath muster <mo> <name> --intent <text>")
		return 2
	}
	opusArg := positional[0]

	// Resolve instance directory from $ATHANOR or --athanor flag
	instDir, err := resolveInstDir(athName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitCode(err)
	}

	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Resolve opus file path
	opusPath := opusArg
	if !filepath.IsAbs(opusPath) {
		mos, _ := athanor.ListMagnaOpera(instDir)
		found := false
		for _, mo := range mos {
			candidate := filepath.Join(athanor.OperaPath(instDir, mo), opusArg)
			if _, err := os.Stat(candidate); err == nil {
				opusPath = candidate
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "error: opus file not found: %s\n", opusArg)
			return 1
		}
	} else if _, err := os.Stat(opusPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: opus file not found: %s\n", opusPath)
		return 1
	}

	// Derive crucible name from opus filename if not overridden
	if crucName == "" {
		base := filepath.Base(opusPath)
		base = strings.TrimSuffix(base, ".md")
		if len(base) > 11 && base[4] == '-' && base[7] == '-' && base[10] == '-' {
			base = base[11:]
		}
		crucName = "azer-" + base
	}

	if worktreePath == "" {
		worktreePath = cfg.Project
	}
	if worktreePath == "" {
		h, _ := os.UserHomeDir()
		worktreePath = h
	}

	// Build the azer boot prompt
	moName := athanor.ReadOpusMO(opusPath)

	// Resolve job: --job flag overrides opus frontmatter
	if job == "" {
		job = athanor.ReadOpusJob(opusPath)
	}
	// Resolve formula: --formula flag overrides opus frontmatter
	if formula == "" {
		formula = athanor.ReadOpusFormula(opusPath)
	}
	// Model priority: --model flag > job frontmatter > athanor.yml default
	if model == "" {
		model = athanor.ReadJobModel(instDir, job)
	}
	if model == "" {
		model = cfg.EffectiveAzerModel()
	}

	jobClause, err := jobBootClause(instDir, job)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	formulaClause, err := formulaBootClause(instDir, formula)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	var bootPrompt string
	if moName != "" {
		moPath := athanor.MagnumOpusPath(instDir, moName)
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s, then read %s/azer.md%s%s. Your opus is at %s. Read it and execute.",
			instDir, moPath, instDir, jobClause, formulaClause, opusPath,
		)
	} else {
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s/azer.md%s%s. Your opus is at %s. Read it and execute.",
			instDir, instDir, jobClause, formulaClause, opusPath,
		)
	}

	claudeArgs := fmt.Sprintf(
		"%scd %s && ATHANOR=%s claude --permission-mode auto --model %q %q",
		envSourcePrefix(instDir), worktreePath, instDir, model, bootPrompt,
	)

	athSessionName := athanor.SessionName(athanor.AthanorName(instDir))
	if err := launchCrucible(athSessionName, crucName, worktreePath, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Azer mustered in crucible %q\n", crucName)
	fmt.Printf("  Opus: %s\n", opusPath)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", worktreePath)

	return 0
}

// runMusterIntent launches an autonomous azer that inscribes its own opus from intent.
// Usage: ath muster <mo> <name> --intent <text> [--worktree-path <path>] [--model <model>] [--athanor <name>]
func runMusterIntent(positional []string, intent, worktreePath, model, crucName, athName, job, formula string) int {
	if len(positional) < 2 {
		fmt.Fprintln(os.Stderr, "error: MO name and crucible name required with --intent")
		fmt.Fprintln(os.Stderr, "usage: ath muster <mo> <name> --intent <text>")
		return 2
	}
	moName := positional[0]
	sessionName := positional[1]

	instDir, err := resolveInstDir(athName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitCode(err)
	}

	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Validate the MO exists
	moPath := athanor.MagnumOpusPath(instDir, moName)
	if _, err := os.Stat(moPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: magnum opus not found: %s\n", moPath)
		return 1
	}

	if worktreePath == "" {
		worktreePath = cfg.Project
	}
	if worktreePath == "" {
		h, _ := os.UserHomeDir()
		worktreePath = h
	}

	// Model priority: --model flag > job frontmatter > athanor.yml default
	if model == "" {
		model = athanor.ReadJobModel(instDir, job)
	}
	if model == "" {
		model = cfg.EffectiveAzerModel()
	}

	if crucName == "" {
		crucName = "azer-" + sessionName
	}

	jobClause, err := jobBootClause(instDir, job)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	formulaClause, err := formulaBootClause(instDir, formula)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	bootPrompt := fmt.Sprintf(
		"Read %s/AGENTS.md, then read %s, then read %s/azer.md%s%s. You are an autonomous azer. Your intent: %s. Inscribe an opus for this work under the %s MO using /opus inscribe (dialectical calcinatio), then execute it. Discharge when complete.",
		instDir, moPath, instDir, jobClause, formulaClause, intent, moName,
	)

	claudeArgs := fmt.Sprintf(
		"%scd %s && ATHANOR=%s claude --permission-mode auto --model %q %q",
		envSourcePrefix(instDir), worktreePath, instDir, model, bootPrompt,
	)

	athSessionName := athanor.SessionName(athanor.AthanorName(instDir))
	if err := launchCrucible(athSessionName, crucName, worktreePath, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Azer mustered (intent) in crucible %q\n", crucName)
	fmt.Printf("  Intent: %s\n", intent)
	fmt.Printf("  MO: %s\n", moName)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", worktreePath)

	return 0
}

// resolveInstDir resolves the athanor instance directory from $ATHANOR or a name.
func resolveInstDir(athName string) (string, error) {
	instDir := os.Getenv("ATHANOR")
	if athName != "" {
		home, err := athanor.Home()
		if err != nil {
			return "", err
		}
		instDir = athanor.InstanceDir(home, athName)
	}
	if instDir == "" {
		return "", fmt.Errorf("$ATHANOR not set and --athanor not provided")
	}
	return instDir, nil
}

// exitCode returns 2 for usage-like errors, 1 otherwise.
func exitCode(err error) int {
	if strings.Contains(err.Error(), "not set") || strings.Contains(err.Error(), "not provided") {
		return 2
	}
	return 1
}

// formulaBootClause returns ", then read <path>" for the boot prompt if a formula
// is specified, or an empty string if no formula is set. Returns an error if the
// formula is specified but the FORMULA.md file does not exist at the expected path.
// The clause is appended after the job clause so the formula layers collaboration
// guidance on top of professional identity.
func formulaBootClause(instDir, formula string) (string, error) {
	if formula == "" {
		return "", nil
	}
	formulaPath := filepath.Join(instDir, athanor.FormulaeDir, formula, athanor.FormulaFile)
	if _, err := os.Stat(formulaPath); err != nil {
		return "", fmt.Errorf("formula definition not found: %s", formulaPath)
	}
	return fmt.Sprintf(", then read %s", formulaPath), nil
}

// jobBootClause returns ", then read <path>" for the boot prompt if a job is specified,
// or an empty string if no job is set. Returns an error if the job is specified but
// the JOB.md file does not exist at the expected path.
// If a JOB.local.md layer exists alongside JOB.md, it is appended to the clause.
func jobBootClause(instDir, job string) (string, error) {
	if job == "" {
		return "", nil
	}
	jobPath := filepath.Join(instDir, "jobs", job, athanor.JobFile)
	if _, err := os.Stat(jobPath); err != nil {
		return "", fmt.Errorf("job definition not found: %s (run 'ath sync' to update job symlinks)", jobPath)
	}
	clause := fmt.Sprintf(", then read %s", jobPath)

	// Include athanor-specific layer if present
	localPath := filepath.Join(instDir, "jobs", job, athanor.JobLocalFile)
	if _, err := os.Stat(localPath); err == nil {
		clause += fmt.Sprintf(", then read %s", localPath)
	}
	return clause, nil
}

// envSourcePrefix returns a shell fragment that sources the instance's .env.local
// if it exists. The fragment is safe to prepend to any command — it's a no-op when
// the file is absent. This enables per-athanor secrets like CLAUDE_CODE_OAUTH_TOKEN.
func envSourcePrefix(instDir string) string {
	envPath := filepath.Join(instDir, ".env.local")
	return fmt.Sprintf("test -f %s && . %s; ", envPath, envPath)
}

// launchCrucible creates a tmux window and sends the claude launch command.
func launchCrucible(session, crucName, dir, claudeArgs string) error {
	r := tmux.NewRunner()
	if err := r.EnsureSession(session); err != nil {
		return fmt.Errorf("ensuring tmux session: %w", err)
	}
	if err := r.NewWindow(session, crucName, dir); err != nil {
		return fmt.Errorf("creating crucible: %w", err)
	}
	target := session + ":" + crucName
	if err := r.SendKeysLiteral(target, claudeArgs); err != nil {
		return fmt.Errorf("launching azer: %w", err)
	}
	if err := r.SendKeys(target, "Enter"); err != nil {
		return fmt.Errorf("launching azer: %w", err)
	}
	return nil
}
