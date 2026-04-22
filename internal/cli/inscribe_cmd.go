package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/matthull/athanor/internal/athanor"
)

const maxSlugLen = 50

// collapseHyphens matches two or more consecutive hyphens.
var collapseHyphens = regexp.MustCompile(`-{2,}`)

// slugify converts a human-readable string into a filename-safe slug.
// Examples:
//
//	"Verify tab scaffold PR #4718" → "verify-tab-scaffold-pr-4718"
//	"Fix the   widget--loader!!!"  → "fix-the-widget-loader"
//
// The result is lowercased, non-alphanumeric characters are replaced with
// hyphens, consecutive hyphens are collapsed, and the result is truncated
// to maxLen runes on a word boundary.
func slugify(s string, maxLen int) string {
	s = strings.ToLower(s)

	// Replace non-alphanumeric chars with hyphens
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	s = b.String()

	// Collapse multiple hyphens
	s = collapseHyphens.ReplaceAllString(s, "-")

	// Trim leading/trailing hyphens
	s = strings.Trim(s, "-")

	// Truncate on a hyphen boundary if too long (rune-safe)
	runes := []rune(s)
	if maxLen > 0 && len(runes) > maxLen {
		s = string(runes[:maxLen])
		// Don't cut mid-word — trim back to last hyphen
		if idx := strings.LastIndex(s, "-"); idx > 0 {
			s = s[:idx]
		}
		s = strings.TrimRight(s, "-")
	}

	return s
}

type inscribeArgs struct {
	athanorName string
	moName      string
	intent      string
	job         string
	name        string
	muster      bool
}

func parseInscribeArgs(args []string) (*inscribeArgs, error) {
	positional, flagArgs := splitArgs(args)

	var ia inscribeArgs

	fs := flag.NewFlagSet("inscribe", flag.ContinueOnError)
	fs.StringVar(&ia.intent, "intent", "", "natural language intent for the opus (required)")
	fs.StringVar(&ia.job, "job", "", "job role for the azer (required)")
	fs.StringVar(&ia.name, "name", "", "opus filename override (default: slugified intent)")
	fs.BoolVar(&ia.muster, "muster", false, "immediately muster an azer for this opus")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return nil, err
	}

	if len(positional) < 1 {
		return nil, fmt.Errorf("athanor name and MO name required")
	}
	if len(positional) < 2 {
		return nil, fmt.Errorf("MO name required")
	}
	ia.athanorName = positional[0]
	ia.moName = positional[1]

	if ia.intent == "" {
		return nil, fmt.Errorf("--intent is required")
	}

	if ia.job == "" {
		return nil, fmt.Errorf("--job is required (use \"general\" if no specific role fits)")
	}

	if err := athanor.ValidateJob(ia.job); err != nil {
		return nil, err
	}

	return &ia, nil
}

func runInscribe(args []string) int {
	ia, err := parseInscribeArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, "usage: ath inscribe <athanor> <mo> --job <job> --intent <text> [--name <name>] [--muster]")
		return 2
	}

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	instDir := athanor.InstanceDir(home, ia.athanorName)

	// Validate MO exists
	moPath := athanor.MagnumOpusPath(instDir, ia.moName)
	if _, err := os.Stat(moPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: magnum opus not found: %s\n", moPath)
		return 1
	}

	// Build opus filename
	datestamp := time.Now().Format("2006-01-02")
	slug := ""
	if ia.name != "" {
		slug = slugify(ia.name, maxSlugLen)
	} else {
		slug = slugify(ia.intent, maxSlugLen)
	}
	if slug == "" {
		fmt.Fprintln(os.Stderr, "error: could not generate opus name from intent (try --name)")
		return 1
	}
	opusFilename := fmt.Sprintf("%s-%s.md", datestamp, slug)

	// Ensure opera directory exists
	operaDir := athanor.OperaPath(instDir, ia.moName)
	if err := os.MkdirAll(operaDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating opera dir: %v\n", err)
		return 1
	}
	opusPath := filepath.Join(operaDir, opusFilename)

	// Build opus content
	content := buildOpusContent(datestamp, ia.moName, ia.job, ia.intent, "")

	if err := os.WriteFile(opusPath, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing opus: %v\n", err)
		return 1
	}

	if ia.muster {
		// Muster an azer for this opus
		return musterInscribedOpus(instDir, opusPath, slug)
	}

	fmt.Println(opusPath)
	return 0
}

// buildOpusContent generates the markdown content for an inscribed opus.
func buildOpusContent(datestamp, moName, job, intent, collaborationContext string) string {
	// YAML frontmatter
	var fm strings.Builder
	fm.WriteString("---\n")
	fm.WriteString("status: charged\n")
	fm.WriteString(fmt.Sprintf("inscribed: %s\n", datestamp))
	fm.WriteString(fmt.Sprintf("magnum_opus: %s\n", moName))
	if job != "" {
		fm.WriteString(fmt.Sprintf("job: %s\n", job))
	}
	fm.WriteString("---\n")

	// Body
	var body strings.Builder
	body.WriteString(fm.String())
	body.WriteString(fmt.Sprintf("# %s\n\n", intent))
	body.WriteString(fmt.Sprintf("## Intent\n\n%s\n", intent))

	if collaborationContext != "" {
		body.WriteString(fmt.Sprintf("\n## Collaboration\n\n%s\n", collaborationContext))
	}

	return body.String()
}

// musterInscribedOpus musters an azer for a freshly inscribed opus file.
func musterInscribedOpus(instDir, opusPath, slug string) int {
	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	crucName := "azer-" + slug
	model := cfg.EffectiveAzerModel()

	worktreePath := cfg.Project
	if worktreePath == "" {
		h, _ := os.UserHomeDir()
		worktreePath = h
	}

	moName := athanor.ReadOpusMO(opusPath)
	var bootPrompt string
	if moName != "" {
		moPath := athanor.MagnumOpusPath(instDir, moName)
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s, then read %s/azer.md. Your opus is at %s. Read it and execute.",
			instDir, moPath, instDir, opusPath,
		)
	} else {
		bootPrompt = fmt.Sprintf(
			"Read %s/AGENTS.md, then read %s/azer.md. Your opus is at %s. Read it and execute.",
			instDir, instDir, opusPath,
		)
	}

	claudeArgs := fmt.Sprintf(
		"cd %s && ATHANOR=%s claude --model %q %q",
		worktreePath, instDir, model, bootPrompt,
	)

	athSessionName := athanor.SessionName(athanor.AthanorName(instDir))
	if err := launchCrucible(athSessionName, crucName, worktreePath, claudeArgs); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Opus inscribed and azer mustered in crucible %q\n", crucName)
	fmt.Printf("  Opus: %s\n", opusPath)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Working dir: %s\n", worktreePath)

	return 0
}
