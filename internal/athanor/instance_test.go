package athanor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("effective models use defaults", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{Name: "test"}
		if got := cfg.EffectiveMarutModel(); got != DefaultMarutModel {
			t.Errorf("EffectiveMarutModel() = %q, want %q", got, DefaultMarutModel)
		}
		if got := cfg.EffectiveAzerModel(); got != DefaultAzerModel {
			t.Errorf("EffectiveAzerModel() = %q, want %q", got, DefaultAzerModel)
		}
	})

	t.Run("effective models use overrides", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{Name: "test", MarutModel: "opus", AzerModel: "sonnet"}
		if got := cfg.EffectiveMarutModel(); got != "opus" {
			t.Errorf("EffectiveMarutModel() = %q, want %q", got, "opus")
		}
		if got := cfg.EffectiveAzerModel(); got != "sonnet" {
			t.Errorf("EffectiveAzerModel() = %q, want %q", got, "sonnet")
		}
	})
}

func TestWriteReadConfig(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	cfg := &Config{
		Name:    "bugsnag",
		Project: "/home/matt/code/musashi",
	}

	if err := WriteConfig(dir, cfg); err != nil {
		t.Fatalf("WriteConfig: %v", err)
	}

	got, err := ReadConfig(dir)
	if err != nil {
		t.Fatalf("ReadConfig: %v", err)
	}

	if got.Name != cfg.Name {
		t.Errorf("Name = %q, want %q", got.Name, cfg.Name)
	}
	if got.Project != cfg.Project {
		t.Errorf("Project = %q, want %q", got.Project, cfg.Project)
	}
}

func TestInitInstance(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	home := filepath.Join(tmp, "athanor")

	// Set up home and shared components
	if err := EnsureHome(home); err != nil {
		t.Fatal(err)
	}
	sharedDir := SharedPath(home)
	for _, f := range SharedFiles {
		if err := os.WriteFile(filepath.Join(sharedDir, f), []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Init the instance
	if err := InitInstance(home, "test-project", "/home/matt/code/test"); err != nil {
		t.Fatalf("InitInstance: %v", err)
	}

	instDir := InstanceDir(home, "test-project")

	// Verify directory exists
	if info, err := os.Stat(instDir); err != nil || !info.IsDir() {
		t.Fatalf("instance directory not created")
	}

	// Verify opera/ exists
	if info, err := os.Stat(filepath.Join(instDir, OperaDir)); err != nil || !info.IsDir() {
		t.Fatalf("opera/ directory not created")
	}

	// Verify symlinks
	for _, f := range SharedFiles {
		path := filepath.Join(instDir, f)
		target, err := os.Readlink(path)
		if err != nil {
			t.Errorf("expected symlink for %s: %v", f, err)
			continue
		}
		expectedTarget := filepath.Join("..", "..", SharedDir, f)
		if target != expectedTarget {
			t.Errorf("symlink %s -> %q, want %q", f, target, expectedTarget)
		}
	}

	// Verify athanor.yml
	cfg, err := ReadConfig(instDir)
	if err != nil {
		t.Fatalf("ReadConfig: %v", err)
	}
	if cfg.Name != "test-project" {
		t.Errorf("config Name = %q, want %q", cfg.Name, "test-project")
	}
	if cfg.Project != "/home/matt/code/test" {
		t.Errorf("config Project = %q", cfg.Project)
	}

	// Verify magnum-opus.md exists and has template content
	moPath := filepath.Join(instDir, "magnum-opus.md")
	moData, err := os.ReadFile(moPath)
	if err != nil {
		t.Fatalf("magnum-opus.md not created: %v", err)
	}
	if len(moData) == 0 {
		t.Error("magnum-opus.md is empty")
	}

	// Verify duplicate init fails
	if err := InitInstance(home, "test-project", ""); err == nil {
		t.Error("expected error on duplicate init, got nil")
	}
}

func TestValidateMagnumOpus(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := ValidateMagnumOpus(dir); err == nil {
			t.Error("expected error for missing file")
		}
	})

	t.Run("template with TODOs", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "magnum-opus.md"), []byte("# Test\n\n[TODO] fill this in"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ValidateMagnumOpus(dir); err == nil {
			t.Error("expected error for template with TODOs")
		}
	})

	t.Run("real content passes", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "magnum-opus.md"), []byte("# Test\n\nReal goal content here."), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ValidateMagnumOpus(dir); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
