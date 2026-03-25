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

	// Verify magna-opera/ directory exists
	moDir := filepath.Join(instDir, MagnaOperaDir)
	if info, err := os.Stat(moDir); err != nil || !info.IsDir() {
		t.Fatalf("magna-opera/ directory not created")
	}

	// Verify duplicate init fails
	if err := InitInstance(home, "test-project", ""); err == nil {
		t.Error("expected error on duplicate init, got nil")
	}
}

func TestValidateMO(t *testing.T) {
	t.Run("missing file multi-MO", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, MagnaOperaDir), 0755); err != nil {
			t.Fatal(err)
		}
		if err := ValidateMO(dir, "nonexistent"); err == nil {
			t.Error("expected error for missing MO file")
		}
	})

	t.Run("template with TODOs", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		moDir := filepath.Join(dir, MagnaOperaDir)
		if err := os.MkdirAll(moDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(moDir, "test.md"), []byte("# Test\n\n[TODO] fill this in"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ValidateMO(dir, "test"); err == nil {
			t.Error("expected error for template with TODOs")
		}
	})

	t.Run("real content passes", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		moDir := filepath.Join(dir, MagnaOperaDir)
		if err := os.MkdirAll(moDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(moDir, "test.md"), []byte("# Test\n\nReal goal content here."), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ValidateMO(dir, "test"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("legacy fallback", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "magnum-opus.md"), []byte("# Legacy\n\nReal content."), 0644); err != nil {
			t.Fatal(err)
		}
		// ValidateMO on legacy instance should check magnum-opus.md
		if err := ValidateMO(dir, filepath.Base(dir)); err != nil {
			t.Errorf("unexpected error on legacy: %v", err)
		}
	})
}

func TestHasLegacyMagnumOpus(t *testing.T) {
	t.Run("legacy with magnum-opus.md", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "magnum-opus.md"), []byte("# MO"), 0644); err != nil {
			t.Fatal(err)
		}
		if !HasLegacyMagnumOpus(dir) {
			t.Error("expected legacy=true when magnum-opus.md exists")
		}
	})

	t.Run("multi-MO with magna-opera dir", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, MagnaOperaDir), 0755); err != nil {
			t.Fatal(err)
		}
		if HasLegacyMagnumOpus(dir) {
			t.Error("expected legacy=false when magna-opera/ exists")
		}
	})

	t.Run("magna-opera takes precedence", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, MagnaOperaDir), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "magnum-opus.md"), []byte("# MO"), 0644); err != nil {
			t.Fatal(err)
		}
		if HasLegacyMagnumOpus(dir) {
			t.Error("expected legacy=false when both exist (magna-opera takes precedence)")
		}
	})
}

func TestListMagnaOpera(t *testing.T) {
	t.Run("multi-MO", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		moDir := filepath.Join(dir, MagnaOperaDir)
		if err := os.MkdirAll(moDir, 0755); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"bugsnag.md", "slack-monitoring.md"} {
			if err := os.WriteFile(filepath.Join(moDir, name), []byte("# "+name), 0644); err != nil {
				t.Fatal(err)
			}
		}
		mos, err := ListMagnaOpera(dir)
		if err != nil {
			t.Fatalf("ListMagnaOpera: %v", err)
		}
		if len(mos) != 2 {
			t.Fatalf("expected 2 MOs, got %d", len(mos))
		}
	})

	t.Run("legacy fallback", func(t *testing.T) {
		dir := t.TempDir()
		// Rename dir to have a known basename
		instDir := filepath.Join(dir, "my-athanor")
		if err := os.MkdirAll(instDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(instDir, "magnum-opus.md"), []byte("# MO"), 0644); err != nil {
			t.Fatal(err)
		}
		mos, err := ListMagnaOpera(instDir)
		if err != nil {
			t.Fatalf("ListMagnaOpera: %v", err)
		}
		if len(mos) != 1 || mos[0] != "my-athanor" {
			t.Errorf("expected [my-athanor], got %v", mos)
		}
	})

	t.Run("empty magna-opera dir", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, MagnaOperaDir), 0755); err != nil {
			t.Fatal(err)
		}
		mos, err := ListMagnaOpera(dir)
		if err != nil {
			t.Fatalf("ListMagnaOpera: %v", err)
		}
		if len(mos) != 0 {
			t.Errorf("expected 0 MOs, got %d", len(mos))
		}
	})
}

func TestMagnumOpusPath(t *testing.T) {
	t.Run("multi-MO path", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, MagnaOperaDir), 0755); err != nil {
			t.Fatal(err)
		}
		got := MagnumOpusPath(dir, "bugsnag")
		want := filepath.Join(dir, MagnaOperaDir, "bugsnag.md")
		if got != want {
			t.Errorf("MagnumOpusPath = %q, want %q", got, want)
		}
	})

	t.Run("legacy path", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "magnum-opus.md"), []byte("# MO"), 0644); err != nil {
			t.Fatal(err)
		}
		got := MagnumOpusPath(dir, "anything")
		want := filepath.Join(dir, "magnum-opus.md")
		if got != want {
			t.Errorf("MagnumOpusPath = %q, want %q", got, want)
		}
	})
}

func TestReadOpusMO(t *testing.T) {
	t.Run("with magnum_opus field", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		content := "---\nstatus: charged\nmagnum_opus: bugsnag\n---\n# Test opus"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusMO(path)
		if got != "bugsnag" {
			t.Errorf("ReadOpusMO = %q, want %q", got, "bugsnag")
		}
	})

	t.Run("without magnum_opus field", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		content := "---\nstatus: charged\n---\n# Test opus"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusMO(path)
		if got != "" {
			t.Errorf("ReadOpusMO = %q, want empty", got)
		}
	})
}

func TestWriteMOTemplate(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := WriteMOTemplate(dir, "test-goal"); err != nil {
		t.Fatalf("WriteMOTemplate: %v", err)
	}
	path := filepath.Join(dir, MagnaOperaDir, "test-goal.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("MO template not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("MO template is empty")
	}
}
