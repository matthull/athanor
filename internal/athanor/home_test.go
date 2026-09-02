package athanor

import (
	"os"
	"path/filepath"
	"strings"
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
	got := InstanceDir("/home/matt/athanor", "my-project")
	want := "/home/matt/athanor/athanors/my-project"
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
	for _, name := range []string{"my-project", "blogging"} {
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

func TestValidateFormula(t *testing.T) {
	t.Run("valid formula", func(t *testing.T) {
		t.Parallel()
		instDir := t.TempDir()
		formulaDir := filepath.Join(instDir, FormulaeDir, "coding-dyad")
		if err := os.MkdirAll(formulaDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(formulaDir, FormulaFile), []byte("# Coding Dyad"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ValidateFormula(instDir, "coding-dyad"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing formula errors with available list", func(t *testing.T) {
		t.Parallel()
		instDir := t.TempDir()
		// Create one valid formula so the available list isn't empty
		validDir := filepath.Join(instDir, FormulaeDir, "valid-formula")
		if err := os.MkdirAll(validDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(validDir, FormulaFile), []byte("# Valid"), 0644); err != nil {
			t.Fatal(err)
		}

		err := ValidateFormula(instDir, "nonexistent")
		if err == nil {
			t.Fatal("expected error for missing formula")
		}
		if !strings.Contains(err.Error(), "unknown formula") {
			t.Errorf("error = %q, want mention of unknown formula", err.Error())
		}
		if !strings.Contains(err.Error(), "valid-formula") {
			t.Errorf("error = %q, want mention of available formula", err.Error())
		}
	})

	t.Run("missing formulae dir reports unknown gracefully", func(t *testing.T) {
		t.Parallel()
		instDir := t.TempDir()
		err := ValidateFormula(instDir, "any-formula")
		if err == nil {
			t.Fatal("expected error for missing formula")
		}
		if !strings.Contains(err.Error(), "unknown formula") {
			t.Errorf("error = %q, want mention of unknown formula", err.Error())
		}
	})
}

func TestListFormulae(t *testing.T) {
	t.Run("missing formulae dir returns nil", func(t *testing.T) {
		t.Parallel()
		instDir := t.TempDir()
		got, err := ListFormulae(instDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Errorf("ListFormulae = %v, want nil", got)
		}
	})

	t.Run("lists formula directories with FORMULA.md", func(t *testing.T) {
		t.Parallel()
		instDir := t.TempDir()
		for _, name := range []string{"coding-dyad", "audit-triad"} {
			d := filepath.Join(instDir, FormulaeDir, name)
			if err := os.MkdirAll(d, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(d, FormulaFile), []byte("# "+name), 0644); err != nil {
				t.Fatal(err)
			}
		}
		got, err := ListFormulae(instDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("got %d formulae, want 2: %v", len(got), got)
		}
	})

	t.Run("skips directories without FORMULA.md", func(t *testing.T) {
		t.Parallel()
		instDir := t.TempDir()
		// One real formula
		realDir := filepath.Join(instDir, FormulaeDir, "real")
		if err := os.MkdirAll(realDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(realDir, FormulaFile), []byte("# Real"), 0644); err != nil {
			t.Fatal(err)
		}
		// One empty directory (no FORMULA.md)
		emptyDir := filepath.Join(instDir, FormulaeDir, "empty")
		if err := os.MkdirAll(emptyDir, 0755); err != nil {
			t.Fatal(err)
		}

		got, err := ListFormulae(instDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0] != "real" {
			t.Errorf("ListFormulae = %v, want [real]", got)
		}
	})
}

func TestReadFormulaInfo(t *testing.T) {
	t.Parallel()
	instDir := t.TempDir()

	// Formula with full frontmatter
	dyadDir := filepath.Join(instDir, FormulaeDir, "coding-dyad")
	if err := os.MkdirAll(dyadDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dyadDir, FormulaFile), []byte(`---
summary: Implementation with independent code review
when:
  - "code needs to be written or modified"
  - "any implementation work that will ship to production"
---
# Coding Dyad

Body.
`), 0644); err != nil {
		t.Fatal(err)
	}

	// Formula with no frontmatter
	bareDir := filepath.Join(instDir, FormulaeDir, "bare")
	if err := os.MkdirAll(bareDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bareDir, FormulaFile), []byte("# Bare\n\nNo frontmatter.\n"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("reads full frontmatter", func(t *testing.T) {
		info, err := ReadFormulaInfo(instDir, "coding-dyad")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Name != "coding-dyad" {
			t.Errorf("Name = %q, want %q", info.Name, "coding-dyad")
		}
		if info.Summary != "Implementation with independent code review" {
			t.Errorf("Summary = %q", info.Summary)
		}
		if len(info.When) != 2 {
			t.Fatalf("When has %d entries, want 2", len(info.When))
		}
		if info.When[0] != "code needs to be written or modified" {
			t.Errorf("When[0] = %q", info.When[0])
		}
	})

	t.Run("handles no frontmatter", func(t *testing.T) {
		info, err := ReadFormulaInfo(instDir, "bare")
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

	t.Run("nonexistent formula returns error", func(t *testing.T) {
		_, err := ReadFormulaInfo(instDir, "nope")
		if err == nil {
			t.Error("expected error for nonexistent formula")
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
