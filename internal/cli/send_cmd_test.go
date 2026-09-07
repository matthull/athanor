package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseSendArgs(t *testing.T) {
	t.Run("target and message from positional args", func(t *testing.T) {
		t.Parallel()
		result, exitCode, err := parseSendArgs([]string{"mysession", "hello", "world"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exitCode != 0 {
			t.Fatalf("expected exit 0, got %d", exitCode)
		}
		if result.target != "mysession" {
			t.Errorf("target = %q, want %q", result.target, "mysession")
		}
		if result.message != "hello world" {
			t.Errorf("message = %q, want %q", result.message, "hello world")
		}
	})

	t.Run("--self reads TMUX_PANE", func(t *testing.T) {
		// Cannot use t.Parallel with t.Setenv
		t.Setenv("TMUX_PANE", "%5")
		result, exitCode, err := parseSendArgs([]string{"--self", "hello"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exitCode != 0 {
			t.Fatalf("expected exit 0, got %d", exitCode)
		}
		if result.target != "%5" {
			t.Errorf("target = %q, want %q", result.target, "%5")
		}
		if result.message != "hello" {
			t.Errorf("message = %q, want %q", result.message, "hello")
		}
	})

	t.Run("--self without TMUX_PANE errors", func(t *testing.T) {
		// Cannot use t.Parallel with t.Setenv
		t.Setenv("TMUX_PANE", "")
		_, exitCode, err := parseSendArgs([]string{"--self", "hello"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if exitCode != 2 {
			t.Errorf("exitCode = %d, want 2", exitCode)
		}
	})

	t.Run("-f reads message from file", func(t *testing.T) {
		t.Parallel()
		tmp := filepath.Join(t.TempDir(), "msg.txt")
		if err := os.WriteFile(tmp, []byte("file content"), 0644); err != nil {
			t.Fatal(err)
		}
		// Flags must precede positional args for Go's flag package
		result, exitCode, err := parseSendArgs([]string{"-f", tmp, "mysession"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exitCode != 0 {
			t.Fatalf("expected exit 0, got %d", exitCode)
		}
		if result.target != "mysession" {
			t.Errorf("target = %q, want %q", result.target, "mysession")
		}
		if result.message != "file content" {
			t.Errorf("message = %q, want %q", result.message, "file content")
		}
	})

	t.Run("-f with missing file errors", func(t *testing.T) {
		t.Parallel()
		_, exitCode, err := parseSendArgs([]string{"-f", "/nonexistent/path", "mysession"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if exitCode != 1 {
			t.Errorf("exitCode = %d, want 1", exitCode)
		}
	})

	t.Run("no target errors", func(t *testing.T) {
		t.Parallel()
		_, exitCode, err := parseSendArgs([]string{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if exitCode != 2 {
			t.Errorf("exitCode = %d, want 2", exitCode)
		}
	})

	t.Run("target but no message errors", func(t *testing.T) {
		t.Parallel()
		_, exitCode, err := parseSendArgs([]string{"mysession"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if exitCode != 2 {
			t.Errorf("exitCode = %d, want 2", exitCode)
		}
	})

	t.Run("--skip-escape sets option", func(t *testing.T) {
		t.Parallel()
		result, exitCode, err := parseSendArgs([]string{"--skip-escape", "mysession", "hello"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exitCode != 0 {
			t.Fatalf("expected exit 0, got %d", exitCode)
		}
		if !result.opts.SkipEscape {
			t.Error("SkipEscape = false, want true")
		}
	})

	t.Run("--timeout sets duration", func(t *testing.T) {
		t.Parallel()
		result, exitCode, err := parseSendArgs([]string{"--timeout", "30s", "mysession", "hello"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exitCode != 0 {
			t.Fatalf("expected exit 0, got %d", exitCode)
		}
		if result.opts.Timeout != 30*time.Second {
			t.Errorf("Timeout = %v, want 30s", result.opts.Timeout)
		}
	})

	t.Run("default timeout is 15s", func(t *testing.T) {
		t.Parallel()
		result, _, err := parseSendArgs([]string{"mysession", "hello"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.opts.Timeout != 15*time.Second {
			t.Errorf("Timeout = %v, want 15s", result.opts.Timeout)
		}
	})
}
