package athanor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHome(t *testing.T) {
	t.Run("uses ATHANOR_HOME when set", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("ATHANOR_HOME", tmp)
		got, err := Home()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != tmp {
			t.Errorf("Home() = %q, want %q", got, tmp)
		}
	})

	t.Run("falls back to ~/athanor", func(t *testing.T) {
		t.Setenv("ATHANOR_HOME", "")
		got, err := Home()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		home, _ := os.UserHomeDir()
		want := filepath.Join(home, "athanor")
		if got != want {
			t.Errorf("Home() = %q, want %q", got, want)
		}
	})
}

func TestRepoDir(t *testing.T) {
	t.Run("uses ATHANOR_REPO when set", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("ATHANOR_REPO", tmp)
		got, err := RepoDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != tmp {
			t.Errorf("RepoDir() = %q, want %q", got, tmp)
		}
	})

	t.Run("falls back to ~/code/athanor", func(t *testing.T) {
		t.Setenv("ATHANOR_REPO", "")
		got, err := RepoDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		home, _ := os.UserHomeDir()
		want := filepath.Join(home, "code", "athanor")
		if got != want {
			t.Errorf("RepoDir() = %q, want %q", got, want)
		}
	})
}

func TestEnsureHome(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	home := filepath.Join(tmp, "athanor")

	if err := EnsureHome(home); err != nil {
		t.Fatalf("EnsureHome: %v", err)
	}

	// Verify directories exist
	for _, sub := range []string{"", AthanorsDir} {
		path := filepath.Join(home, sub)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected %s to exist: %v", path, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", path)
		}
	}

	// Idempotent — calling again should not error
	if err := EnsureHome(home); err != nil {
		t.Fatalf("EnsureHome (idempotent): %v", err)
	}
}

func TestInstanceDir(t *testing.T) {
	t.Parallel()
	got := InstanceDir("/home/matt/athanor", "bugsnag")
	want := "/home/matt/athanor/athanors/bugsnag"
	if got != want {
		t.Errorf("InstanceDir() = %q, want %q", got, want)
	}
}

func TestListInstances(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	home := filepath.Join(tmp, "athanor")

	// Empty — no athanors dir
	names, err := ListInstances(home)
	if err != nil {
		t.Fatalf("ListInstances (empty): %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected 0 instances, got %d", len(names))
	}

	// Create some instances
	if err := EnsureHome(home); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"bugsnag", "blogging"} {
		if err := os.MkdirAll(filepath.Join(home, AthanorsDir, name), 0755); err != nil {
			t.Fatal(err)
		}
	}

	names, err = ListInstances(home)
	if err != nil {
		t.Fatalf("ListInstances: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 instances, got %d: %v", len(names), names)
	}
}

func TestReadJobInfo(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()

	// Create a job with full frontmatter
	coderDir := filepath.Join(tmp, "coder")
	if err := os.MkdirAll(coderDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(coderDir, "JOB.md"), []byte(`---
model: sonnet
summary: Write, modify, and review software
when:
  - "code needs to be written or modified"
  - "implementation needs review"
---
# Coder

You are a coder.
`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a job with no frontmatter
	bareDir := filepath.Join(tmp, "bare")
	if err := os.MkdirAll(bareDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bareDir, "JOB.md"), []byte("# Bare\n\nNo frontmatter.\n"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("reads full frontmatter", func(t *testing.T) {
		t.Parallel()
		info, err := ReadJobInfo(tmp, "coder")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Name != "coder" {
			t.Errorf("Name = %q, want %q", info.Name, "coder")
		}
		if info.Summary != "Write, modify, and review software" {
			t.Errorf("Summary = %q", info.Summary)
		}
		if info.Model != "sonnet" {
			t.Errorf("Model = %q, want %q", info.Model, "sonnet")
		}
		if len(info.When) != 2 {
			t.Fatalf("When has %d entries, want 2", len(info.When))
		}
		if info.When[0] != "code needs to be written or modified" {
			t.Errorf("When[0] = %q", info.When[0])
		}
		if info.When[1] != "implementation needs review" {
			t.Errorf("When[1] = %q", info.When[1])
		}
	})

	t.Run("handles no frontmatter", func(t *testing.T) {
		t.Parallel()
		info, err := ReadJobInfo(tmp, "bare")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Summary != "" {
			t.Errorf("Summary = %q, want empty", info.Summary)
		}
		if len(info.When) != 0 {
			t.Errorf("When has %d entries, want 0", len(info.When))
		}
	})

	t.Run("nonexistent job returns error", func(t *testing.T) {
		t.Parallel()
		_, err := ReadJobInfo(tmp, "nope")
		if err == nil {
			t.Error("expected error for nonexistent job")
		}
	})
}

func TestExtractListField(t *testing.T) {
	t.Parallel()
	fm := `
summary: test
when:
  - "first item"
  - "second item"
  - unquoted item
model: sonnet
`
	got := extractListField(fm, "when:")
	if len(got) != 3 {
		t.Fatalf("got %d items, want 3: %v", len(got), got)
	}
	if got[0] != "first item" {
		t.Errorf("[0] = %q", got[0])
	}
	if got[1] != "second item" {
		t.Errorf("[1] = %q", got[1])
	}
	if got[2] != "unquoted item" {
		t.Errorf("[2] = %q", got[2])
	}
}

func TestExpandHome(t *testing.T) {
	t.Run("expands tilde", func(t *testing.T) {
		got, err := expandHome("~/foo")
		if err != nil {
			t.Fatal(err)
		}
		home, _ := os.UserHomeDir()
		want := filepath.Join(home, "foo")
		if got != want {
			t.Errorf("expandHome(~/foo) = %q, want %q", got, want)
		}
	})

	t.Run("no tilde unchanged", func(t *testing.T) {
		t.Parallel()
		got, err := expandHome("/absolute/path")
		if err != nil {
			t.Fatal(err)
		}
		if got != "/absolute/path" {
			t.Errorf("expandHome(/absolute/path) = %q", got)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		got, err := expandHome("")
		if err != nil {
			t.Fatal(err)
		}
		if got != "" {
			t.Errorf("expandHome('') = %q", got)
		}
	})
}
