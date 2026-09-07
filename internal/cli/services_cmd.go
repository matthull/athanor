package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type serviceEntry struct {
	Unit        string
	Description string
}

var athanorServices = []serviceEntry{
	{"athanor-liveness.timer", "Keeps kindled presence-driven roles alive (5-min cadence)"},
	{"attunement-intake.timer", "Processes voice note signals into attunement corpus (30-min cadence)"},
	{"voice-notes-process.timer", "Transcribes voice notes from GDrive (5-min cadence)"},
}

type serviceStatus struct {
	Unit        string    `json:"unit"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	LastTrigger string    `json:"last_trigger"`
	NextTrigger string    `json:"next_trigger"`
	lastTime    time.Time `json:"-"`
	nextTime    time.Time `json:"-"`
}

type servicesSnapshot struct {
	Timestamp  string          `json:"timestamp"`
	Services   []serviceStatus `json:"services"`
	AllHealthy bool            `json:"all_healthy"`
}

// listTimerEntry mirrors one element of `systemctl --user list-timers --output=json`.
// next/last are Unix microsecond timestamps; zero means unset (e.g., never triggered).
type listTimerEntry struct {
	Unit string `json:"unit"`
	Next int64  `json:"next"`
	Last int64  `json:"last"`
}

func runServices(args []string) int {
	var jsonOut bool
	fs := flag.NewFlagSet("services", flag.ContinueOnError)
	fs.BoolVar(&jsonOut, "json", false, "output as JSON")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	snap := collectServices()

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(snap)
	} else {
		printServices(snap)
	}

	if snap.AllHealthy {
		return 0
	}
	return 1
}

func collectServices() servicesSnapshot {
	snap := servicesSnapshot{
		Timestamp:  time.Now().Format("15:04:05"),
		AllHealthy: true,
	}

	units := make([]string, 0, len(athanorServices))
	for _, svc := range athanorServices {
		units = append(units, svc.Unit)
	}
	timerTimes := queryTimerTimes(units)

	for _, svc := range athanorServices {
		status, _ := checkServiceStatus(svc.Unit)

		ss := serviceStatus{
			Unit:        svc.Unit,
			Description: svc.Description,
			Status:      status,
		}

		if entry, ok := timerTimes[svc.Unit]; ok {
			if entry.Last > 0 {
				ss.lastTime = time.UnixMicro(entry.Last)
				ss.LastTrigger = ss.lastTime.Format("Mon 2006-01-02 15:04:05 MST")
			}
			if entry.Next > 0 {
				ss.nextTime = time.UnixMicro(entry.Next)
				ss.NextTrigger = ss.nextTime.Format("Mon 2006-01-02 15:04:05 MST")
			}
		}

		if status != "active" {
			snap.AllHealthy = false
		}

		snap.Services = append(snap.Services, ss)
	}

	return snap
}

func checkServiceStatus(unit string) (string, error) {
	out, err := exec.Command("systemctl", "--user", "is-active", unit).Output()
	status := strings.TrimSpace(string(out))
	if status == "" {
		return "unknown", err
	}
	return status, nil
}

// queryTimerTimes runs `systemctl --user list-timers <units...> --no-pager --output=json`
// and returns a map from unit name to its timer entry. Monotonic timers (OnUnitActiveSec=)
// don't populate NextElapseUSecRealtime in `systemctl show`, but list-timers always does.
// On any error, returns an empty map (graceful degradation — services still display).
func queryTimerTimes(units []string) map[string]listTimerEntry {
	if len(units) == 0 {
		return map[string]listTimerEntry{}
	}
	args := append([]string{"--user", "list-timers"}, units...)
	args = append(args, "--no-pager", "--output=json")
	out, err := exec.Command("systemctl", args...).Output()
	if err != nil {
		return map[string]listTimerEntry{}
	}
	return parseListTimersJSON(out)
}

// parseListTimersJSON decodes the JSON array emitted by list-timers into a map
// keyed by unit name. Invalid JSON yields an empty map.
func parseListTimersJSON(data []byte) map[string]listTimerEntry {
	result := map[string]listTimerEntry{}
	var entries []listTimerEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return result
	}
	for _, e := range entries {
		result[e.Unit] = e
	}
	return result
}

// formatTimeSince renders a time as "HH:MM (Xm ago)" or "HH:MM (in Xm)".
// Zero time returns "".
func formatTimeSince(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	hhmm := t.Format("15:04")
	delta := time.Since(t)
	if delta < 0 {
		return fmt.Sprintf("%s (in %s)", hhmm, roundDuration(-delta))
	}
	return fmt.Sprintf("%s (%s ago)", hhmm, roundDuration(delta))
}

func roundDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	default:
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
}

func colorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func colorWrap(text, status string) string {
	if !colorEnabled() {
		return text
	}
	switch status {
	case "active":
		return "\033[32m" + text + "\033[0m"
	case "inactive", "failed":
		return "\033[31m" + text + "\033[0m"
	default:
		return "\033[33m" + text + "\033[0m"
	}
}

func printServices(snap servicesSnapshot) {
	fmt.Println("ATH SERVICES")
	fmt.Println()
	fmt.Printf("%-30s %-10s %-22s %s\n", "UNIT", "STATUS", "LAST TRIGGER", "NEXT")

	for _, svc := range snap.Services {
		lastFmt := formatTimeSince(svc.lastTime)
		nextFmt := formatTimeSince(svc.nextTime)
		if lastFmt == "" {
			lastFmt = "—"
		}
		if nextFmt == "" {
			nextFmt = "—"
		}
		statusPad := fmt.Sprintf("%-10s", svc.Status)
		fmt.Printf("%-30s %s %-22s %s\n",
			svc.Unit,
			colorWrap(statusPad, svc.Status),
			lastFmt,
			nextFmt,
		)
	}

	fmt.Println()
	unhealthy := 0
	for _, svc := range snap.Services {
		if svc.Status != "active" {
			unhealthy++
		}
	}
	if unhealthy == 0 {
		fmt.Println("All services healthy.")
	} else {
		fmt.Printf("%d service(s) unhealthy — check failed units above.\n", unhealthy)
	}
}
