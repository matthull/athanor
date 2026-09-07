package tmux

import (
	"errors"
	"strings"
	"time"
)

// Idle detection constants. [G:tmux.go:2474-2479,2483,2538]
const (
	idlePollInterval        = 200 * time.Millisecond // Time between consecutive polls
	requiredConsecutiveIdle = 2                      // Polls needed to confirm idle
	idleCaptureLines        = 5                      // Lines to capture per poll
)

// WaitForIdle polls a tmux pane until the Claude Code agent is confirmed idle.
// Idle requires 2 consecutive readings (200ms apart) showing the prompt prefix
// without a busy status bar — filtering out transient prompt flashes during
// inter-tool-call gaps (~500ms). [G:tmux.go:2470-2541]
func WaitForIdle(r *Runner, target string, timeout time.Duration) error {
	consecutiveIdle := 0
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		lines, err := r.CapturePaneLines(target, idleCaptureLines)
		if err != nil {
			// Terminal errors — session is gone, no recovery possible.
			if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrNoServer) {
				return err
			}
			// Transient capture error — reset and retry.
			consecutiveIdle = 0
			time.Sleep(idlePollInterval)
			continue
		}

		// Check status bar first: if busy, agent is running a tool call.
		// NOT idle regardless of prompt visibility. [G:tmux.go:2499-2513]
		if IsStatusBarBusy(lines) {
			consecutiveIdle = 0
			time.Sleep(idlePollInterval)
			continue
		}

		// Scan for prompt prefix in captured lines. [G:tmux.go:2515-2528]
		promptFound := false
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			if MatchesPromptPrefix(line) {
				promptFound = true
				break
			}
		}

		// Count consecutive idle readings. [G:tmux.go:2530-2536]
		if promptFound {
			consecutiveIdle++
			if consecutiveIdle >= requiredConsecutiveIdle {
				return nil // Confirmed idle
			}
		} else {
			consecutiveIdle = 0
		}

		time.Sleep(idlePollInterval)
	}

	return ErrIdleTimeout
}
