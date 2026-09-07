package tmux

import "testing"

func TestIdleDetectionLogic(t *testing.T) {
	t.Parallel()

	t.Run("status bar busy prevents idle", func(t *testing.T) {
		t.Parallel()
		lines := []string{
			"some output",
			"⏵⏵ Running tool                    esc to interrupt",
			"",
		}
		if !IsStatusBarBusy(lines) {
			t.Error("should detect busy status bar")
		}
	})

	t.Run("prompt without busy status is idle candidate", func(t *testing.T) {
		t.Parallel()
		lines := []string{
			"previous output",
			"❯ ",
			"  ○ 3 functions | 12 subtasks",
		}
		if IsStatusBarBusy(lines) {
			t.Error("should not be busy")
		}
		found := false
		for _, line := range lines {
			if MatchesPromptPrefix(line) {
				found = true
				break
			}
		}
		if !found {
			t.Error("should find prompt prefix")
		}
	})

	t.Run("no prompt means not idle", func(t *testing.T) {
		t.Parallel()
		lines := []string{
			"Compiling...",
			"Building module graph",
			"50% complete",
		}
		if IsStatusBarBusy(lines) {
			t.Error("should not be busy")
		}
		found := false
		for _, line := range lines {
			if MatchesPromptPrefix(line) {
				found = true
				break
			}
		}
		if found {
			t.Error("should not find prompt prefix in build output")
		}
	})
}
