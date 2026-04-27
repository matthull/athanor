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
	tmp := t.TempDir()
	home := filepath.Join(tmp, "athanor")
	repoDir := filepath.Join(tmp, "repo")

	// Set ATHANOR_REPO to temp repo dir
	t.Setenv("ATHANOR_REPO", repoDir)

	// Set up home
	if err := EnsureHome(home); err != nil {
		t.Fatal(err)
	}

	// Set up shared components in the repo
	sharedDir := filepath.Join(repoDir, SharedDir)
	if err := os.MkdirAll(sharedDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, f := range SharedFiles {
		if err := os.WriteFile(filepath.Join(sharedDir, f), []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range SharedDirs {
		if err := os.MkdirAll(filepath.Join(sharedDir, d), 0755); err != nil {
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

	// Verify file symlinks point to absolute repo paths
	for _, f := range SharedFiles {
		path := filepath.Join(instDir, f)
		target, err := os.Readlink(path)
		if err != nil {
			t.Errorf("expected symlink for %s: %v", f, err)
			continue
		}
		expectedTarget := filepath.Join(sharedDir, f)
		if target != expectedTarget {
			t.Errorf("symlink %s -> %q, want %q", f, target, expectedTarget)
		}
	}

	// Verify directory symlinks
	for _, d := range SharedDirs {
		path := filepath.Join(instDir, d)
		target, err := os.Readlink(path)
		if err != nil {
			t.Errorf("expected symlink for dir %s: %v", d, err)
			continue
		}
		expectedTarget := filepath.Join(sharedDir, d)
		if target != expectedTarget {
			t.Errorf("symlink %s -> %q, want %q", d, target, expectedTarget)
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
		if err := os.MkdirAll(filepath.Join(dir, MagnaOperaDir, "nonexistent"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := ValidateMO(dir, "nonexistent"); err == nil {
			t.Error("expected error for missing MO file")
		}
	})

	t.Run("template with TODOs", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		moDir := filepath.Join(dir, MagnaOperaDir, "test")
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
		moDir := filepath.Join(dir, MagnaOperaDir, "test")
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
		for _, name := range []string{"bugsnag", "slack-monitoring"} {
			moDir := filepath.Join(dir, MagnaOperaDir, name)
			if err := os.MkdirAll(moDir, 0755); err != nil {
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
		want := filepath.Join(dir, MagnaOperaDir, "bugsnag", "bugsnag.md")
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

func TestReadOpusJob(t *testing.T) {
	t.Run("with job field", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		content := "---\nstatus: charged\nmagnum_opus: bugsnag\njob: coder\n---\n# Test opus"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusJob(path)
		if got != "coder" {
			t.Errorf("ReadOpusJob = %q, want %q", got, "coder")
		}
	})

	t.Run("without job field", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		content := "---\nstatus: charged\nmagnum_opus: bugsnag\n---\n# Test opus"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusJob(path)
		if got != "" {
			t.Errorf("ReadOpusJob = %q, want empty", got)
		}
	})

	t.Run("with empty job field", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		content := "---\nstatus: charged\njob: \n---\n# Test opus"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusJob(path)
		if got != "" {
			t.Errorf("ReadOpusJob = %q, want empty", got)
		}
	})
}

// setupSyncTestRepo creates a temp home and repo with shared files/dirs,
// returning (home, sharedDir). Sets ATHANOR_REPO for the test.
func setupSyncTestRepo(t *testing.T) (string, string) {
	t.Helper()
	tmp := t.TempDir()
	home := filepath.Join(tmp, "athanor")
	repoDir := filepath.Join(tmp, "repo")
	t.Setenv("ATHANOR_REPO", repoDir)

	if err := EnsureHome(home); err != nil {
		t.Fatal(err)
	}
	sharedDir := filepath.Join(repoDir, SharedDir)
	if err := os.MkdirAll(sharedDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, f := range SharedFiles {
		if err := os.WriteFile(filepath.Join(sharedDir, f), []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range SharedDirs {
		if err := os.MkdirAll(filepath.Join(sharedDir, d), 0755); err != nil {
			t.Fatal(err)
		}
	}
	return home, sharedDir
}

func TestSyncInstance(t *testing.T) {
	t.Run("creates missing symlinks", func(t *testing.T) {
		home, sharedDir := setupSyncTestRepo(t)

		// Create a bare instance directory (no symlinks)
		name := "bare-instance"
		instDir := InstanceDir(home, name)
		if err := os.MkdirAll(instDir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := SyncInstance(home, name); err != nil {
			t.Fatalf("SyncInstance: %v", err)
		}

		// Verify all file symlinks created
		for _, f := range SharedFiles {
			target, err := os.Readlink(filepath.Join(instDir, f))
			if err != nil {
				t.Errorf("expected symlink for %s: %v", f, err)
				continue
			}
			if target != filepath.Join(sharedDir, f) {
				t.Errorf("symlink %s -> %q, want %q", f, target, filepath.Join(sharedDir, f))
			}
		}

		// Verify all dir symlinks created
		for _, d := range SharedDirs {
			target, err := os.Readlink(filepath.Join(instDir, d))
			if err != nil {
				t.Errorf("expected symlink for dir %s: %v", d, err)
				continue
			}
			if target != filepath.Join(sharedDir, d) {
				t.Errorf("symlink %s -> %q, want %q", d, target, filepath.Join(sharedDir, d))
			}
		}
	})

	t.Run("idempotent on correct symlinks", func(t *testing.T) {
		home, _ := setupSyncTestRepo(t)

		// Init creates correct symlinks
		if err := InitInstance(home, "idem-test", ""); err != nil {
			t.Fatalf("InitInstance: %v", err)
		}

		// Sync again — should succeed with no changes
		if err := SyncInstance(home, "idem-test"); err != nil {
			t.Fatalf("SyncInstance (idempotent): %v", err)
		}
	})

	t.Run("fixes stale symlinks", func(t *testing.T) {
		home, sharedDir := setupSyncTestRepo(t)

		name := "stale-test"
		instDir := InstanceDir(home, name)
		if err := os.MkdirAll(instDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create a symlink pointing to wrong target
		wrongTarget := "/tmp/nonexistent-target"
		if err := os.Symlink(wrongTarget, filepath.Join(instDir, SharedFiles[0])); err != nil {
			t.Fatal(err)
		}

		if err := SyncInstance(home, name); err != nil {
			t.Fatalf("SyncInstance: %v", err)
		}

		// Verify the stale symlink was fixed
		target, err := os.Readlink(filepath.Join(instDir, SharedFiles[0]))
		if err != nil {
			t.Fatalf("expected symlink for %s: %v", SharedFiles[0], err)
		}
		expected := filepath.Join(sharedDir, SharedFiles[0])
		if target != expected {
			t.Errorf("symlink %s -> %q, want %q", SharedFiles[0], target, expected)
		}
	})

	t.Run("errors on nonexistent instance", func(t *testing.T) {
		home, _ := setupSyncTestRepo(t)
		if err := SyncInstance(home, "nonexistent"); err == nil {
			t.Error("expected error for nonexistent instance")
		}
	})

	t.Run("refuses to overwrite regular file", func(t *testing.T) {
		home, _ := setupSyncTestRepo(t)

		name := "regular-file-test"
		instDir := InstanceDir(home, name)
		if err := os.MkdirAll(instDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create a regular file where a symlink should go
		if err := os.WriteFile(filepath.Join(instDir, SharedFiles[0]), []byte("real file"), 0644); err != nil {
			t.Fatal(err)
		}

		err := SyncInstance(home, name)
		if err == nil {
			t.Error("expected error when regular file exists at symlink target")
		}
	})
}

func TestWriteMOTemplate(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := WriteMOTemplate(dir, "test-goal"); err != nil {
		t.Fatalf("WriteMOTemplate: %v", err)
	}
	path := filepath.Join(dir, MagnaOperaDir, "test-goal", "test-goal.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("MO template not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("MO template is empty")
	}
	// Verify opera subdirectory was created
	operaDir := filepath.Join(dir, MagnaOperaDir, "test-goal", "opera")
	if info, err := os.Stat(operaDir); err != nil || !info.IsDir() {
		t.Fatalf("opera/ subdirectory not created under MO directory")
	}
}
