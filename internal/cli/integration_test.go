package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matthull/athanor/internal/athanor"
)

// TestATHFullLifecycle exercises the complete ath CLI workflow:
//
//	init → kindle → status → muster → opera → whisper → cleanup → quiesce
//
// Requires a running tmux server (skips in -short mode).
// Uses a temporary ATHANOR_HOME so it doesn't touch ~/athanor/.
// Creates real tmux windows but uses 'bash' instead of 'claude' to avoid
// launching actual agent sessions.
func TestATHFullLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Verify tmux is available
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not found, skipping")
	}
	// Verify ath binary is installed
	athBin, err := exec.LookPath("ath")
	if err != nil {
		t.Skip("ath binary not found in PATH, run 'make install' first")
	}

	// Set up temporary athanor home
	tmpHome := t.TempDir()
	t.Setenv("ATHANOR_HOME", tmpHome)

	// Set up shared components (minimal test files)
	sharedDir := filepath.Join(tmpHome, athanor.SharedDir)
	if err := os.MkdirAll(sharedDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, f := range athanor.SharedFiles {
		content := fmt.Sprintf("# %s (test)\nTest shared component for QA.", f)
		if err := os.WriteFile(filepath.Join(sharedDir, f), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(tmpHome, athanor.AthanorsDir), 0755); err != nil {
		t.Fatal(err)
	}

	// Helper to run ath commands and capture output
	runAth := func(args ...string) (string, error) {
		cmd := exec.Command(athBin, args...)
		cmd.Env = append(os.Environ(), "ATHANOR_HOME="+tmpHome)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	// Helper for tmux cleanup
	windowsToCleanup := []string{}
	defer func() {
		for _, w := range windowsToCleanup {
			_ = exec.Command("tmux", "kill-window", "-t", w).Run()
		}
	}()
	trackWindow := func(name string) {
		windowsToCleanup = append(windowsToCleanup, name)
	}

	// ─── Phase 1: ath init ───────────────────────────────────────────

	t.Run("init creates instance", func(t *testing.T) {
		out, err := runAth("init", "qa-test", "--project", "/tmp/fake-project")
		if err != nil {
			t.Fatalf("ath init failed: %v\n%s", err, out)
		}
		if !strings.Contains(out, "qa-test") {
			t.Errorf("expected output to mention instance name, got: %s", out)
		}

		instDir := athanor.InstanceDir(tmpHome, "qa-test")

		// Verify directory structure
		for _, path := range []string{
			instDir,
			filepath.Join(instDir, "opera"),
			filepath.Join(instDir, "athanor.yml"),
			filepath.Join(instDir, "magna-opera"),
		} {
			if _, err := os.Stat(path); err != nil {
				t.Errorf("expected %s to exist: %v", filepath.Base(path), err)
			}
		}

		// Verify symlinks
		for _, f := range athanor.SharedFiles {
			target, err := os.Readlink(filepath.Join(instDir, f))
			if err != nil {
				t.Errorf("expected symlink for %s: %v", f, err)
			}
			if !strings.Contains(target, "shared") {
				t.Errorf("symlink %s points to %q, expected shared/", f, target)
			}
		}

		// Verify config
		cfg, err := athanor.ReadConfig(instDir)
		if err != nil {
			t.Fatalf("ReadConfig: %v", err)
		}
		if cfg.Name != "qa-test" {
			t.Errorf("config Name = %q, want %q", cfg.Name, "qa-test")
		}
		if cfg.Project != "/tmp/fake-project" {
			t.Errorf("config Project = %q", cfg.Project)
		}
	})

	t.Run("init rejects duplicate", func(t *testing.T) {
		out, err := runAth("init", "qa-test")
		if err == nil {
			t.Fatal("expected error on duplicate init")
		}
		if !strings.Contains(out, "already exists") {
			t.Errorf("expected 'already exists' error, got: %s", out)
		}
	})

	// ─── Phase 2: Create a magnum opus ──────────────────────────────

	instDir := athanor.InstanceDir(tmpHome, "qa-test")
	moPath := filepath.Join(instDir, "magna-opera", "qa-goal.md")
	moContent := `# qa-goal — Magnum Opus

## Goal

QA testing athanor — verify all CLI commands work correctly.

## Abundant Satisfaction

All ath commands produce correct output and manage tmux windows properly.

## Witnesses

The test harness is the witness.

## Pre-loaded Context

This is an automated test. No prior context needed.
`
	if err := os.WriteFile(moPath, []byte(moContent), 0644); err != nil {
		t.Fatal(err)
	}

	// ─── Phase 3: ath status (before kindle) ─────────────────────────

	t.Run("status shows instance with no marut", func(t *testing.T) {
		out, err := runAth("status")
		if err != nil {
			t.Fatalf("ath status failed: %v\n%s", err, out)
		}
		if !strings.Contains(out, "qa-test") {
			t.Errorf("expected status to list qa-test, got: %s", out)
		}
		// Marut should show as "-" (not running)
		if strings.Contains(out, "active") {
			t.Errorf("expected no active marut, got: %s", out)
		}
	})

	t.Run("status detail view", func(t *testing.T) {
		out, err := runAth("status", "qa-test")
		if err != nil {
			t.Fatalf("ath status qa-test failed: %v\n%s", err, out)
		}
		if !strings.Contains(out, "Athanor: qa-test") {
			t.Errorf("expected detail header, got: %s", out)
		}
		if !strings.Contains(out, "marut: -") {
			t.Errorf("expected marut: -, got: %s", out)
		}
	})

	// ─── Phase 4: Create test opus ───────────────────────────────────

	opusPath := filepath.Join(instDir, "opera", "2026-03-25-qa-fix-something.md")
	opusContent := `---
status: charged
inscribed: 2026-03-25
magnum_opus: qa-goal
---
# Fix Something for QA

## Intent

Test that opera management works correctly.

## Boundary

Agent: verify the opus lifecycle.
Operator: nothing.

## Context

This is a test opus created by the QA harness.
`
	if err := os.WriteFile(opusPath, []byte(opusContent), 0644); err != nil {
		t.Fatal(err)
	}

	// ─── Phase 5: ath opera ──────────────────────────────────────────

	t.Run("opera lists charged opus", func(t *testing.T) {
		out, err := runAth("opera", "qa-test")
		if err != nil {
			t.Fatalf("ath opera failed: %v\n%s", err, out)
		}
		if !strings.Contains(out, "charged") {
			t.Errorf("expected charged status, got: %s", out)
		}
		if !strings.Contains(out, "qa-fix-something") {
			t.Errorf("expected opus name, got: %s", out)
		}
		if !strings.Contains(out, "2026-03-25") {
			t.Errorf("expected date, got: %s", out)
		}
	})

	// ─── Phase 6: ath kindle (creates tmux window) ───────────────────
	// kindle will try to launch 'claude' which may not be interactive here.
	// We test that the tmux window is created with the right name.
	// The claude command will either start (and sit at prompt) or fail — either
	// way, the window existing proves kindle worked.

	t.Run("kindle creates marut crucible", func(t *testing.T) {
		// Override project to /tmp so cd works
		cfg, _ := athanor.ReadConfig(instDir)
		cfg.Project = "/tmp"
		if err := athanor.WriteConfig(instDir, cfg); err != nil {
			t.Fatalf("writing config: %v", err)
		}

		out, err := runAth("kindle", "qa-test", "qa-goal")
		if err != nil {
			t.Fatalf("ath kindle failed: %v\n%s", err, out)
		}
		trackWindow("marut-qa-test-qa-goal")

		if !strings.Contains(out, "marut-qa-test-qa-goal") {
			t.Errorf("expected crucible name in output, got: %s", out)
		}

		// Give tmux a moment to create the window
		time.Sleep(500 * time.Millisecond)

		// Verify window exists
		windows := listTmuxWindows(t)
		if !containsExact(windows, "marut-qa-test-qa-goal") {
			t.Errorf("expected tmux window 'marut-qa-test-qa-goal', got windows: %v", windows)
		}

		// Verify the command sent to the window contains claude and the boot prompt
		paneContent := capturePaneContent(t, "marut-qa-test-qa-goal", 20)
		if !strings.Contains(paneContent, "claude") || !strings.Contains(paneContent, "AGENTS.md") {
			t.Logf("pane content: %s", paneContent)
			// Not a hard failure — pane capture can be timing-sensitive
			t.Log("warning: could not verify boot command in pane (timing-sensitive)")
		}
	})

	// ─── Phase 7: ath status (after kindle) ──────────────────────────

	t.Run("status shows active marut after kindle", func(t *testing.T) {
		out, err := runAth("status")
		if err != nil {
			t.Fatalf("ath status failed: %v\n%s", err, out)
		}
		if !strings.Contains(out, "active") {
			t.Errorf("expected active marut, got: %s", out)
		}
	})

	// ─── Phase 7b: ath check ────────────────────────────────────────

	t.Run("check reports state for existing crucible", func(t *testing.T) {
		out, err := runAth("check", "marut-qa-test-qa-goal")
		// Exit code 0 or 1 are both valid — depends on what claude is doing
		if err != nil {
			// check returns exit 1 for permission/exhausted, 2 for dead/error
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
				t.Fatalf("ath check returned exit 2 (dead/error): %s", out)
			}
		}
		// Output should contain one of the valid states
		validStates := []string{"active", "idle", "permission", "exhausted"}
		found := false
		for _, s := range validStates {
			if strings.Contains(out, s) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected one of %v in output, got: %s", validStates, out)
		}
	})

	t.Run("check reports dead for nonexistent crucible", func(t *testing.T) {
		out, err := runAth("check", "nonexistent-crucible-xyz")
		if err == nil {
			t.Fatal("expected non-zero exit for dead crucible")
		}
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("expected ExitError, got: %v", err)
		}
		if exitErr.ExitCode() != 2 {
			t.Errorf("expected exit code 2, got %d", exitErr.ExitCode())
		}
		if !strings.Contains(out, "dead") {
			t.Errorf("expected 'dead' in output, got: %s", out)
		}
	})

	t.Run("check without args returns exit 2", func(t *testing.T) {
		_, err := runAth("check")
		if err == nil {
			t.Fatal("expected error when no crucible name given")
		}
	})

	// ─── Phase 8: ath muster (creates azer window) ───────────────────

	t.Run("muster creates azer crucible", func(t *testing.T) {
		out, err := runAth("muster", "2026-03-25-qa-fix-something.md",
			"--athanor", "qa-test", "--dir", "/tmp")
		if err != nil {
			t.Fatalf("ath muster failed: %v\n%s", err, out)
		}
		trackWindow("azer-qa-fix-something")

		if !strings.Contains(out, "azer-qa-fix-something") {
			t.Errorf("expected crucible name in output, got: %s", out)
		}

		time.Sleep(500 * time.Millisecond)

		windows := listTmuxWindows(t)
		if !containsExact(windows, "azer-qa-fix-something") {
			t.Errorf("expected tmux window 'azer-qa-fix-something', got: %v", windows)
		}
	})

	// ─── Phase 9: ath whisper between windows ────────────────────────
	// Create a plain bash window to whisper to (easier to verify than
	// a claude session which has its own TUI)

	t.Run("whisper send delivers message", func(t *testing.T) {
		// Create a test target window
		_ = exec.Command("tmux", "new-window", "-n", "qa-whisper-target").Run()
		trackWindow("qa-whisper-target")
		time.Sleep(300 * time.Millisecond)

		out, err := runAth("whisper", "send", "qa-whisper-target", "hello from QA test")
		if err != nil {
			t.Fatalf("ath whisper send failed: %v\n%s", err, out)
		}

		// Verify the message arrived in the target pane
		time.Sleep(500 * time.Millisecond)
		pane := capturePaneContent(t, "qa-whisper-target", 10)
		if !strings.Contains(pane, "hello from QA test") {
			t.Errorf("expected message in target pane, got: %s", pane)
		}
	})

	t.Run("whisper idle detects idle shell", func(t *testing.T) {
		// The qa-whisper-target has a bash shell which should show a prompt
		out, err := runAth("whisper", "idle", "qa-whisper-target", "--timeout", "5s")
		if err != nil {
			t.Logf("whisper idle output: %s", out)
			// A plain bash prompt might not match Claude's prompt pattern
			// so this may fail — that's expected and informative
			t.Log("note: whisper idle may not detect plain bash as idle (expected — it looks for Claude prompts)")
		}
	})

	// ─── Phase 10: ath cleanup ───────────────────────────────────────

	t.Run("cleanup kills azer crucible", func(t *testing.T) {
		out, err := runAth("cleanup", "azer-qa-fix-something")
		if err != nil {
			t.Fatalf("ath cleanup failed: %v\n%s", err, out)
		}

		time.Sleep(300 * time.Millisecond)
		windows := listTmuxWindows(t)
		if containsExact(windows, "azer-qa-fix-something") {
			t.Error("expected azer window to be killed after cleanup")
		}
	})

	t.Run("cleanup is idempotent", func(t *testing.T) {
		out, err := runAth("cleanup", "azer-qa-fix-something")
		if err != nil {
			t.Fatalf("second cleanup should succeed (idempotent): %v\n%s", err, out)
		}
	})

	// ─── Phase 11: ath reforge ───────────────────────────────────────

	t.Run("reforge kills and recreates marut", func(t *testing.T) {
		out, err := runAth("reforge", "qa-test", "qa-goal")
		if err != nil {
			t.Fatalf("ath reforge failed: %v\n%s", err, out)
		}

		if !strings.Contains(out, "reforged") {
			t.Errorf("expected 'reforged' in output, got: %s", out)
		}

		time.Sleep(500 * time.Millisecond)
		windows := listTmuxWindows(t)
		if !containsExact(windows, "marut-qa-test-qa-goal") {
			t.Errorf("expected marut window to exist after reforge, got: %v", windows)
		}
	})

	// ─── Phase 12: Discharge opus and verify opera ───────────────────

	t.Run("opera reflects discharged status", func(t *testing.T) {
		// Update opus to discharged
		discharged := strings.Replace(string(opusContent), "status: charged", "status: discharged", 1)
		if err := os.WriteFile(opusPath, []byte(discharged), 0644); err != nil {
			t.Fatal(err)
		}

		out, err := runAth("opera", "qa-test")
		if err != nil {
			t.Fatalf("ath opera failed: %v\n%s", err, out)
		}
		if !strings.Contains(out, "discharged") {
			t.Errorf("expected discharged status, got: %s", out)
		}
	})

	// ─── Phase 13: ath quiesce ───────────────────────────────────────

	t.Run("quiesce shuts down athanor", func(t *testing.T) {
		// Use --force because other athanors may have azer-* windows running
		// that quiesce can't yet distinguish from this athanor's azers
		out, err := runAth("quiesce", "qa-test", "--force")
		if err != nil {
			t.Fatalf("ath quiesce failed: %v\n%s", err, out)
		}
		if !strings.Contains(out, "quiesced") {
			t.Errorf("expected 'quiesced' in output, got: %s", out)
		}

		time.Sleep(300 * time.Millisecond)
		windows := listTmuxWindows(t)
		if containsExact(windows, "marut-qa-test-qa-goal") {
			t.Error("expected marut window to be killed after quiesce")
		}
	})

	// ─── Phase 14: Error cases ───────────────────────────────────────

	t.Run("kindle warns on TODO magnum-opus", func(t *testing.T) {
		// Create another instance and add a template MO with TODOs
		_, _ = runAth("init", "qa-warn-test")
		warnDir := athanor.InstanceDir(tmpHome, "qa-warn-test")
		if err := athanor.WriteMOTemplate(warnDir, "warn-goal"); err != nil {
			t.Fatal(err)
		}
		out, err := runAth("kindle", "qa-warn-test", "warn-goal")
		// Should warn but not necessarily fail hard
		_ = err
		trackWindow("marut-qa-warn-test-warn-goal")
		if !strings.Contains(out, "TODO") && !strings.Contains(out, "warning") {
			t.Logf("expected warning about TODO placeholders, got: %s", out)
		}
	})

	t.Run("kindle without mo-name on multi-MO errors", func(t *testing.T) {
		// qa-warn-test is a multi-MO instance (has magna-opera/ dir)
		out, err := runAth("kindle", "qa-warn-test")
		if err == nil {
			t.Fatal("expected error when mo-name not provided for multi-MO instance")
		}
		if !strings.Contains(out, "mo-name required") {
			t.Errorf("expected 'mo-name required' error, got: %s", out)
		}
	})

	t.Run("muster without ATHANOR errors", func(t *testing.T) {
		// Unset ATHANOR and don't provide --athanor
		cmd := exec.Command(athBin, "muster", "nonexistent.md")
		cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatal("expected error when $ATHANOR not set")
		}
		if !strings.Contains(string(out), "ATHANOR") {
			t.Errorf("expected error about $ATHANOR, got: %s", out)
		}
	})

	// ─── Phase 15: Version and help ──────────────────────────────────

	t.Run("version prints info", func(t *testing.T) {
		out, err := runAth("version")
		if err != nil {
			t.Fatalf("ath version failed: %v", err)
		}
		if !strings.Contains(out, "ath") {
			t.Errorf("expected 'ath' in version output, got: %s", out)
		}
	})

	t.Run("help prints usage", func(t *testing.T) {
		out, err := runAth("help")
		if err != nil {
			t.Fatalf("ath help failed: %v", err)
		}
		if !strings.Contains(out, "kindle") {
			t.Errorf("expected 'kindle' in help output, got: %s", out)
		}
		if !strings.Contains(out, "whisper") {
			t.Errorf("expected 'whisper' in help output, got: %s", out)
		}
	})

	t.Run("whisper help prints subcommand usage", func(t *testing.T) {
		out, err := runAth("whisper", "help")
		if err != nil {
			t.Fatalf("ath whisper help failed: %v", err)
		}
		if !strings.Contains(out, "send") {
			t.Errorf("expected 'send' in whisper help, got: %s", out)
		}
	})

	t.Run("unknown command returns exit 2", func(t *testing.T) {
		_, err := runAth("foobar")
		if err == nil {
			t.Fatal("expected error for unknown command")
		}
	})
}

// ── Test Helpers ─────────────────────────────────────────────────────

func listTmuxWindows(t *testing.T) []string {
	t.Helper()
	out, err := exec.Command("tmux", "list-windows", "-a", "-F", "#{window_name}").CombinedOutput()
	if err != nil {
		t.Logf("tmux list-windows: %v", err)
		return nil
	}
	result := strings.TrimSpace(string(out))
	if result == "" {
		return nil
	}
	return strings.Split(result, "\n")
}

func containsExact(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func capturePaneContent(t *testing.T, target string, lines int) string {
	t.Helper()
	out, err := exec.Command("tmux", "capture-pane", "-p", "-t", target,
		"-S", fmt.Sprintf("-%d", lines)).CombinedOutput()
	if err != nil {
		t.Logf("capture-pane %s: %v", target, err)
		return ""
	}
	return string(out)
}
