package tmux

import "strings"

// Pane represents a tmux pane with its coordinates and process info.
type Pane struct {
	// Target is the tmux target string: "session:window.pane"
	Target string
	// Session name
	Session string
	// Window index
	Window string
	// Pane index
	Pane string
	// Command is the current foreground process in the pane
	Command string
}

// ListPanes returns all panes across all sessions with their coordinates and
// current command. Each pane is identified by "session:window.pane" format.
func (r *Runner) ListPanes() ([]Pane, error) {
	out, err := r.run("list-panes", "-a", "-F",
		"#{session_name}\t#{window_index}\t#{pane_index}\t#{pane_current_command}")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}

	var panes []Pane
	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) < 4 {
			continue
		}
		p := Pane{
			Session: parts[0],
			Window:  parts[1],
			Pane:    parts[2],
			Command: parts[3],
		}
		p.Target = p.Session + ":" + p.Window + "." + p.Pane
		panes = append(panes, p)
	}
	return panes, nil
}
