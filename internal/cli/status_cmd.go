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

func runStatus(args []string) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	remaining := fs.Args()

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if len(remaining) > 0 {
		return showInstanceStatus(home, remaining[0])
	}
	return showAllStatus(home)
}

func showAllStatus(home string) int {
	instances, err := athanor.ListInstances(home)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if len(instances) == 0 {
		fmt.Println("No athanor instances found.")
		fmt.Printf("Create one with: ath init <name>\n")
		return 0
	}

	r := tmux.NewRunner()
	fmt.Printf("%-20s %-24s %-8s %-6s %s\n", "ATHANOR", "MAGNUM OPUS", "MARUT", "AZERS", "OPERA (C/D)")
	for _, name := range instances {
		instDir := athanor.InstanceDir(home, name)

		mos, _ := athanor.ListMagnaOpera(instDir)
		if len(mos) == 0 {
			// No MOs at all
			azerCount := countAzerWindows(r, name)
			charged, discharged := countOpera(instDir)
			fmt.Printf("%-20s %-24s %-8s %-6d %d/%d\n", name, "(none)", "-", azerCount, charged, discharged)
			continue
		}

		legacy := athanor.HasLegacyMagnumOpus(instDir)
		azerCount := countAzerWindows(r, name)
		charged, discharged := countOpera(instDir)

		for i, mo := range mos {
			var crucible string
			if legacy {
				crucible = athanor.MarutCrucibleName(name, "")
			} else {
				crucible = athanor.MarutCrucibleName(name, mo)
			}
			marutStatus := "-"
			if windowExists(r, crucible) {
				marutStatus = "active"
			}

			moLabel := mo
			if legacy {
				moLabel = "(legacy)"
			}

			// Show athanor name and opera counts only on first row
			if i == 0 {
				fmt.Printf("%-20s %-24s %-8s %-6d %d/%d\n", name, moLabel, marutStatus, azerCount, charged, discharged)
			} else {
				fmt.Printf("%-20s %-24s %-8s\n", "", moLabel, marutStatus)
			}
		}
	}

	return 0
}

func showInstanceStatus(home, name string) int {
	instDir := athanor.InstanceDir(home, name)

	cfg, err := athanor.ReadConfig(instDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	r := tmux.NewRunner()

	fmt.Printf("Athanor: %s\n", name)
	if cfg.Project != "" {
		fmt.Printf("Project: %s\n", cfg.Project)
	}

	// List magna opera and their marut status
	mos, _ := athanor.ListMagnaOpera(instDir)
	legacy := athanor.HasLegacyMagnumOpus(instDir)
	if len(mos) > 0 {
		fmt.Println("Magna Opera:")
		for _, mo := range mos {
			var crucible string
			if legacy {
				crucible = athanor.MarutCrucibleName(name, "")
			} else {
				crucible = athanor.MarutCrucibleName(name, mo)
			}
			status := "-"
			if windowExists(r, crucible) {
				status = fmt.Sprintf("active (%s)", crucible)
			}
			label := mo
			if legacy {
				label = fmt.Sprintf("%s (legacy)", mo)
			}
			fmt.Printf("  %s — marut: %s\n", label, status)
		}
	} else {
		fmt.Println("Magna Opera: (none)")
	}

	// List opera
	operaDir := filepath.Join(instDir, athanor.OperaDir)
	entries, err := os.ReadDir(operaDir)
	if err == nil {
		var charged, discharged, assessed []string
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			status := readOpusStatus(filepath.Join(operaDir, e.Name()))
			opusName := strings.TrimSuffix(e.Name(), ".md")
			switch status {
			case "charged":
				charged = append(charged, opusName)
			case "discharged":
				discharged = append(discharged, opusName)
			case "assessed":
				assessed = append(assessed, opusName)
			}
		}
		fmt.Println("Opera:")
		if len(charged) > 0 {
			fmt.Printf("  Charged:    %s\n", strings.Join(charged, ", "))
		}
		if len(discharged) > 0 {
			fmt.Printf("  Discharged: %s\n", strings.Join(discharged, ", "))
		}
		if len(assessed) > 0 {
			fmt.Printf("  Assessed:   %s\n", strings.Join(assessed, ", "))
		}
		if len(charged)+len(discharged)+len(assessed) == 0 {
			fmt.Println("  (none)")
		}
	}

	return 0
}

// windowExists checks if a tmux window with the exact given name exists.
func windowExists(r *tmux.Runner, name string) bool {
	windows, err := r.ListWindows()
	if err != nil {
		return false
	}
	for _, w := range windows {
		if w == name {
			return true
		}
	}
	return false
}

// countAzerWindows counts tmux windows matching "azer-*" pattern.
// This is a simple heuristic — checks common patterns.
func countAzerWindows(r *tmux.Runner, _ string) int {
	// List all windows and count ones starting with "azer-"
	windows, err := r.ListWindows()
	if err != nil {
		return 0
	}
	count := 0
	for _, w := range windows {
		if strings.HasPrefix(w, "azer-") {
			count++
		}
	}
	return count
}

// countOpera counts opera files by status (charged, discharged).
func countOpera(instDir string) (charged, discharged int) {
	operaDir := filepath.Join(instDir, athanor.OperaDir)
	entries, err := os.ReadDir(operaDir)
	if err != nil {
		return 0, 0
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		status := readOpusStatus(filepath.Join(operaDir, e.Name()))
		switch status {
		case "charged":
			charged++
		case "discharged":
			discharged++
		}
	}
	return
}

// readOpusStatus reads the YAML frontmatter status from an opus file.
func readOpusStatus(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return ""
	}
	end := strings.Index(content[3:], "---")
	if end < 0 {
		return ""
	}
	frontmatter := content[3 : 3+end]
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "status:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "status:"))
		}
	}
	return ""
}
