package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Sentinel errors for tmux operations. [G:tmux.go:42-47]
var (
	ErrNoServer        = errors.New("no tmux server running")
	ErrSessionNotFound = errors.New("session not found")
	ErrIdleTimeout     = errors.New("agent not idle before timeout")
)

// Runner wraps tmux command execution with UTF-8 and socket support.
// All tmux interaction goes through Runner.run(). [G:tmux.go:168-204]
type Runner struct {
	socket string // tmux socket name (-L flag), empty = default
}

// NewRunner creates a Runner using the default tmux socket.
func NewRunner() *Runner {
	return &Runner{}
}

// NewRunnerWithSocket creates a Runner targeting a named socket.
// Used in tests for isolation and for explicit socket targeting.
func NewRunnerWithSocket(s string) *Runner {
	return &Runner{socket: s}
}

// run executes a tmux command and returns stdout.
// All commands include -u flag for UTF-8 support regardless of locale settings.
// [G:tmux.go:209-228]
func (r *Runner) run(args ...string) (string, error) {
	// Prepend -u (UTF-8 mode) and optionally -L (socket) before the subcommand.
	allArgs := []string{"-u"}
	if r.socket != "" {
		allArgs = append(allArgs, "-L", r.socket)
	}
	allArgs = append(allArgs, args...)

	cmd := exec.Command("tmux", allArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", r.wrapError(err, stderr.String(), args)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// wrapError classifies tmux stderr into sentinel errors. [G:tmux.go:231-253]
func (r *Runner) wrapError(err error, stderr string, args []string) error {
	stderr = strings.TrimSpace(stderr)

	if strings.Contains(stderr, "no server running") ||
		strings.Contains(stderr, "error connecting to") ||
		strings.Contains(stderr, "server exited unexpectedly") {
		return ErrNoServer
	}
	if strings.Contains(stderr, "session not found") ||
		strings.Contains(stderr, "can't find session") {
		return ErrSessionNotFound
	}

	if stderr != "" {
		return fmt.Errorf("tmux %s: %s", args[0], stderr)
	}
	return fmt.Errorf("tmux %s: %w", args[0], err)
}

// SendKeys sends raw key names (e.g., "Enter", "Escape") to a target.
func (r *Runner) SendKeys(target string, keys ...string) error {
	args := append([]string{"send-keys", "-t", target}, keys...)
	_, err := r.run(args...)
	return err
}

// SendKeysLiteral sends literal text (send-keys -l) to a target.
func (r *Runner) SendKeysLiteral(target, text string) error {
	_, err := r.run("send-keys", "-t", target, "-l", text)
	return err
}

// DisplayMessage runs display-message -p and returns the formatted output.
func (r *Runner) DisplayMessage(target, format string) (string, error) {
	return r.run("display-message", "-p", "-t", target, format)
}

// CapturePaneLines captures the last N lines of a pane as a slice.
func (r *Runner) CapturePaneLines(target string, lines int) ([]string, error) {
	out, err := r.run("capture-pane", "-p", "-t", target, "-S", fmt.Sprintf("-%d", lines))
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// ResizeWindow resizes the window containing target to width x.
func (r *Runner) ResizeWindow(target string, width int) error {
	_, err := r.run("resize-window", "-t", target, "-x", fmt.Sprintf("%d", width))
	return err
}

// GetWindowWidth returns the current width of the window containing target.
func (r *Runner) GetWindowWidth(target string) (int, error) {
	out, err := r.run("display-message", "-p", "-t", target, "#{window_width}")
	if err != nil {
		return 0, err
	}
	var w int
	if _, err := fmt.Sscanf(strings.TrimSpace(out), "%d", &w); err != nil {
		return 0, fmt.Errorf("parse window width %q: %w", out, err)
	}
	return w, nil
}

// SetWindowOption sets a window option on the target.
func (r *Runner) SetWindowOption(target, option, value string) error {
	_, err := r.run("set-option", "-w", "-t", target, option, value)
	return err
}

// IsSessionAttached returns true if the session has any clients attached.
func (r *Runner) IsSessionAttached(target string) bool {
	out, err := r.run("display-message", "-t", target, "-p", "#{session_attached}")
	return err == nil && strings.TrimSpace(out) == "1"
}

// CancelCopyMode sends cancel to exit copy/scroll mode.
func (r *Runner) CancelCopyMode(target string) error {
	_, err := r.run("send-keys", "-t", target, "-X", "cancel")
	return err
}

// InCopyMode returns true if the target pane is in copy/scroll mode.
func (r *Runner) InCopyMode(target string) (bool, error) {
	out, err := r.run("display-message", "-p", "-t", target, "#{pane_in_mode}")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == "1", nil
}
