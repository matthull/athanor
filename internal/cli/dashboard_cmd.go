package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/matthull/athanor/internal/athanor"
	"github.com/matthull/athanor/internal/tmux"
)

// ── Dashboard data model ─────────────────────────────────────────────

type dashboardSnapshot struct {
	Timestamp string             `json:"timestamp"`
	Athanors  []dashboardAthanor `json:"athanors"`
}

type dashboardAthanor struct {
	Name    string        `json:"name"`
	Project string        `json:"project,omitempty"`
	MOs     []dashboardMO `json:"magna_opera"`
}

type dashboardMO struct {
	Name       string          `json:"name"`
	Goal       string          `json:"goal"`
	MarutName  string          `json:"marut_name"`
	MarutState string          `json:"marut_state"`
	InFlight   []dashboardOpus `json:"in_flight,omitempty"`
	Waiting    []dashboardOpus `json:"waiting,omitempty"`
	Recent     []dashboardOpus `json:"recent,omitempty"`
	TrailSize  int             `json:"trail_size"`
}

type dashboardOpus struct {
	Name      string `json:"name"`
	Date      string `json:"date"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	AzerName  string `json:"azer_name,omitempty"`
	AzerState string `json:"azer_state,omitempty"`
	Outcome   string `json:"outcome,omitempty"`
}

// ── Command entry point ──────────────────────────────────────────────

func runDashboard(args []string) int {
	var (
		watch    bool
		interval int
		jsonOut  bool
	)

	fs := flag.NewFlagSet("dashboard", flag.ContinueOnError)
	fs.BoolVar(&watch, "watch", false, "continuously refresh the dashboard")
	fs.IntVar(&interval, "interval", 5, "refresh interval in seconds (with --watch)")
	fs.BoolVar(&jsonOut, "json", false, "output as JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if !watch {
		snap := collectDashboard(home)
		if jsonOut {
			printDashboardJSON(snap)
		} else {
			printDashboard(snap)
		}
		return 0
	}

	for {
		snap := collectDashboard(home)
		fmt.Print("\033[H\033[2J")
		if jsonOut {
			printDashboardJSON(snap)
		} else {
			printDashboard(snap)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// ── Snapshot collection ──────────────────────────────────────────────

func collectDashboard(home string) dashboardSnapshot {
	now := time.Now()
	snap := dashboardSnapshot{
		Timestamp: now.Format("15:04:05"),
	}

	instances, err := athanor.ListInstances(home)
	if err != nil || len(instances) == 0 {
		return snap
	}

	r := tmux.NewRunner()
	windows, _ := r.ListWindows()

	for _, name := range instances {
		instDir := athanor.InstanceDir(home, name)
		av := collectAthanorView(instDir, name, r, windows, now)
		if len(av.MOs) > 0 {
			snap.Athanors = append(snap.Athanors, av)
		}
	}

	return snap
}

func collectAthanorView(instDir, name string, r *tmux.Runner, windows []string, now time.Time) dashboardAthanor {
	cfg, _ := athanor.ReadConfig(instDir)
	project := ""
	if cfg != nil {
		project = cfg.Project
	}

	av := dashboardAthanor{
		Name:    name,
		Project: project,
	}

	mos, _ := athanor.ListMagnaOpera(instDir)
	legacy := athanor.HasLegacyMagnumOpus(instDir)

	for _, mo := range mos {
		operaDir := athanor.OperaPath(instDir, mo)
		operaEntries, _ := os.ReadDir(operaDir)
		mv := collectMOView(instDir, name, mo, legacy, r, windows, operaEntries, operaDir, now)
		av.MOs = append(av.MOs, mv)
	}

	return av
}

func collectMOView(instDir, athName, moName string, legacy bool, r *tmux.Runner, windows []string, operaEntries []os.DirEntry, operaDir string, now time.Time) dashboardMO {
	moPath := athanor.MagnumOpusPath(instDir, moName)
	moContent, _ := os.ReadFile(moPath)
	goal := extractSection(string(moContent), "Goal")

	var crucible string
	if legacy {
		crucible = athanor.MarutCrucibleName(athName, "")
	} else {
		crucible = athanor.MarutCrucibleName(athName, moName)
	}

	marutState := "dead"
	if sliceContains(windows, crucible) {
		state, _, _ := r.CheckCrucible(crucible)
		marutState = state.String()
	}

	mv := dashboardMO{
		Name:       moName,
		Goal:       goal,
		MarutName:  crucible,
		MarutState: marutState,
	}

	cutoff := now.Add(-48 * time.Hour)

	for _, e := range operaEntries {
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		opusPath := filepath.Join(operaDir, e.Name())
		data, err := os.ReadFile(opusPath)
		if err != nil {
			continue
		}
		content := string(data)

		opusMO := readFrontmatterField(content, "magnum_opus")
		if opusMO != "" && opusMO != moName {
			continue
		}

		status := readFrontmatterField(content, "status")
		opusFilename := strings.TrimSuffix(e.Name(), ".md")
		date, strippedName := splitOpusFilename(opusFilename)

		title := extractTitle(content)
		if title == "" {
			title = strippedName
		}

		ov := dashboardOpus{
			Name:   strippedName,
			Date:   date,
			Title:  title,
			Status: status,
		}

		switch status {
		case "charged":
			azerName := "azer-" + strippedName
			if sliceContains(windows, azerName) {
				state, _, _ := r.CheckCrucible(azerName)
				ov.AzerName = azerName
				ov.AzerState = state.String()
				mv.InFlight = append(mv.InFlight, ov)
			} else {
				mv.Waiting = append(mv.Waiting, ov)
			}

		case "discharged":
			mv.TrailSize++
			dischargedDate := readFrontmatterField(content, "discharged")
			checkDate := dischargedDate
			if checkDate == "" {
				checkDate = date
			}
			if isRecentDate(checkDate, cutoff) {
				ov.Outcome = extractSection(content, "Outcome")
				mv.Recent = append(mv.Recent, ov)
			}

		case "assessed":
			mv.TrailSize++
		}
	}

	return mv
}

// ── Rendering ────────────────────────────────────────────────────────

func printDashboard(snap dashboardSnapshot) {
	fmt.Printf("ATH DASHBOARD  %s\n\n", snap.Timestamp)

	if len(snap.Athanors) == 0 {
		fmt.Println("  No athanor instances found.")
		fmt.Println("  Create one with: ath init <name>")
		return
	}

	first := true
	for _, ath := range snap.Athanors {
		for _, mo := range ath.MOs {
			if !first {
				fmt.Println("───")
				fmt.Println()
			}
			first = false

			// Header
			header := ath.Name
			if mo.Name != ath.Name {
				header = ath.Name + "/" + mo.Name
			}
			marut := fmt.Sprintf("marut %s %s", stateIcon(mo.MarutState), mo.MarutState)
			if ath.Project != "" {
				fmt.Printf("%s  %s  %s\n", header, marut, shortenPath(ath.Project))
			} else {
				fmt.Printf("%s  %s\n", header, marut)
			}

			// Goal
			if mo.Goal != "" {
				fmt.Printf("  %s\n", mo.Goal)
			}

			// In flight
			if len(mo.InFlight) > 0 {
				fmt.Println()
				fmt.Println("  IN FLIGHT")
				for _, op := range mo.InFlight {
					icon := stateIcon(op.AzerState)
					suffix := ""
					if op.AzerState != "active" {
						suffix = fmt.Sprintf(" [%s]", op.AzerState)
					}
					fmt.Printf("    %s %s%s\n", icon, op.Title, suffix)
				}
			}

			// Waiting
			if len(mo.Waiting) > 0 {
				fmt.Println()
				fmt.Println("  WAITING")
				for _, op := range mo.Waiting {
					fmt.Printf("    ○ %s\n", op.Title)
				}
			}

			// Recent
			if len(mo.Recent) > 0 {
				fmt.Println()
				fmt.Println("  RECENT")
				for _, op := range mo.Recent {
					dateLabel := formatShortDate(op.Date)
					suffix := ""
					if op.Outcome != "" {
						suffix = " — " + op.Outcome
					}
					fmt.Printf("    ✓ %s%s (%s)\n", op.Title, suffix, dateLabel)
				}
			}

			// Trail
			active := len(mo.InFlight) + len(mo.Waiting)
			if mo.TrailSize > 0 || active > 0 {
				fmt.Printf("\n  trail: %d discharged", mo.TrailSize)
				if active > 0 {
					fmt.Printf(", %d active", active)
				}
				fmt.Println()
			}
		}
	}
}

func printDashboardJSON(snap dashboardSnapshot) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(snap)
}

// ── Helpers ──────────────────────────────────────────────────────────

// stateIcon returns a unicode indicator for a crucible state.
func stateIcon(state string) string {
	switch state {
	case "active":
		return "●"
	case "idle":
		return "◐"
	case "permission":
		return "⚠"
	case "exhausted":
		return "✖"
	case "dead":
		return "✖"
	default:
		return "?"
	}
}

// extractSection returns the first paragraph from a markdown ## heading.
func extractSection(content, heading string) string {
	marker := "## " + heading
	idx := strings.Index(content, marker)
	if idx < 0 {
		return ""
	}

	rest := content[idx+len(marker):]
	if nl := strings.IndexByte(rest, '\n'); nl >= 0 {
		rest = rest[nl+1:]
	} else {
		return ""
	}

	for strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	}

	var parts []string
	for _, line := range strings.Split(rest, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "# ") {
			break
		}
		parts = append(parts, trimmed)
	}

	result := strings.Join(parts, " ")
	if len(result) > 120 {
		result = result[:117] + "..."
	}
	return result
}

// extractTitle returns the text from the first # heading (not ## or deeper).
func extractTitle(content string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") && !strings.HasPrefix(trimmed, "## ") {
			title := strings.TrimPrefix(trimmed, "# ")
			title = strings.TrimPrefix(title, "Opus: ")
			return title
		}
	}
	return ""
}

// readFrontmatterField extracts a field value from YAML frontmatter (--- delimited).
func readFrontmatterField(content, field string) string {
	if !strings.HasPrefix(content, "---") {
		return ""
	}
	end := strings.Index(content[3:], "---")
	if end < 0 {
		return ""
	}
	fm := content[3 : 3+end]
	prefix := field + ":"
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

// splitOpusFilename extracts date and stripped name from "YYYY-MM-DD-name" format.
func splitOpusFilename(filename string) (date, name string) {
	if len(filename) >= 11 && filename[4] == '-' && filename[7] == '-' && filename[10] == '-' {
		return filename[:10], filename[11:]
	}
	return "", filename
}

// isRecentDate checks if a YYYY-MM-DD date string is on or after the cutoff.
func isRecentDate(dateStr string, cutoff time.Time) bool {
	if len(dateStr) < 10 {
		return false
	}
	t, err := time.Parse("2006-01-02", dateStr[:10])
	if err != nil {
		return false
	}
	return !t.Before(cutoff)
}

// formatShortDate converts "2026-03-25" to "Mar 25".
func formatShortDate(dateStr string) string {
	if len(dateStr) < 10 {
		return dateStr
	}
	t, err := time.Parse("2006-01-02", dateStr[:10])
	if err != nil {
		return dateStr
	}
	return t.Format("Jan 2")
}

// sliceContains returns true if slice contains the target string.
func sliceContains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// shortenPath replaces the home directory prefix with ~.
func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
