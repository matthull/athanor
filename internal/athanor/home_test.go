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

func TestEnsureHome(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	home := filepath.Join(tmp, "athanor")

	if err := EnsureHome(home); err != nil {
		t.Fatalf("EnsureHome: %v", err)
	}

	// Verify directories exist
	for _, sub := range []string{"", AthanorsDir, SharedDir} {
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
