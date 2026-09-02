package athanor

import (
	"os"
	"path/filepath"
	"strings"
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
		Name:    "my-project",
		Project: "/home/matt/code/my-app",
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
		for _, name := range []string{"my-project", "slack-monitoring"} {
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
		got := MagnumOpusPath(dir, "my-project")
		want := filepath.Join(dir, MagnaOperaDir, "my-project", "my-project.md")
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
		content := "---\nstatus: charged\nmagnum_opus: my-project\n---\n# Test opus"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusMO(path)
		if got != "my-project" {
			t.Errorf("ReadOpusMO = %q, want %q", got, "my-project")
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
		content := "---\nstatus: charged\nmagnum_opus: my-project\njob: coder\n---\n# Test opus"
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
		content := "---\nstatus: charged\nmagnum_opus: my-project\n---\n# Test opus"
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

func TestReadOpusFormula(t *testing.T) {
	t.Run("with formula field", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		content := "---\nstatus: charged\njob: coder\nformula: coding-dyad\n---\n# Test"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusFormula(path)
		if got != "coding-dyad" {
			t.Errorf("ReadOpusFormula = %q, want %q", got, "coding-dyad")
		}
	})

	t.Run("without formula field", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		content := "---\nstatus: charged\njob: coder\n---\n# Test"
		path := filepath.Join(dir, "test.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadOpusFormula(path)
		if got != "" {
			t.Errorf("ReadOpusFormula = %q, want empty", got)
		}
	})

	t.Run("missing file returns empty", func(t *testing.T) {
		t.Parallel()
		got := ReadOpusFormula("/nonexistent/path.md")
		if got != "" {
			t.Errorf("ReadOpusFormula = %q, want empty", got)
		}
	})
}

func TestReadJobModel(t *testing.T) {
	t.Run("reads model from JOB.md frontmatter", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		jobDir := filepath.Join(dir, "jobs", "coder")
		if err := os.MkdirAll(jobDir, 0755); err != nil {
			t.Fatal(err)
		}
		content := "---\nmodel: sonnet\n---\n# Coder\n\nYou are a coder."
		if err := os.WriteFile(filepath.Join(jobDir, "JOB.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadJobModel(dir, "coder")
		if got != "sonnet" {
			t.Errorf("ReadJobModel = %q, want %q", got, "sonnet")
		}
	})

	t.Run("returns empty for JOB.md without frontmatter", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		jobDir := filepath.Join(dir, "jobs", "general")
		if err := os.MkdirAll(jobDir, 0755); err != nil {
			t.Fatal(err)
		}
		content := "# General\n\nYou are a general-purpose azer."
		if err := os.WriteFile(filepath.Join(jobDir, "JOB.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadJobModel(dir, "general")
		if got != "" {
			t.Errorf("ReadJobModel = %q, want empty", got)
		}
	})

	t.Run("returns empty for nonexistent job", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		got := ReadJobModel(dir, "nonexistent")
		if got != "" {
			t.Errorf("ReadJobModel = %q, want empty", got)
		}
	})

	t.Run("returns empty for empty job name", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		got := ReadJobModel(dir, "")
		if got != "" {
			t.Errorf("ReadJobModel = %q, want empty", got)
		}
	})

	t.Run("reads model with other frontmatter fields", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		jobDir := filepath.Join(dir, "jobs", "qa-specialist")
		if err := os.MkdirAll(jobDir, 0755); err != nil {
			t.Fatal(err)
		}
		content := "---\nmodel: sonnet\nother_field: value\n---\n# QA Specialist"
		if err := os.WriteFile(filepath.Join(jobDir, "JOB.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		got := ReadJobModel(dir, "qa-specialist")
		if got != "sonnet" {
			t.Errorf("ReadJobModel = %q, want %q", got, "sonnet")
		}
	})
}

func TestWriteEnvTemplate(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := WriteEnvTemplate(dir); err != nil {
		t.Fatalf("WriteEnvTemplate: %v", err)
	}
	path := filepath.Join(dir, ".env.local.template")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("template not created: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("template is empty")
	}
	if !strings.Contains(content, "CLAUDE_CODE_OAUTH_TOKEN") {
		t.Error("template should mention CLAUDE_CODE_OAUTH_TOKEN")
	}
}

func TestEnsureEnvGitignore(t *testing.T) {
	t.Run("creates gitignore entry", func(t *testing.T) {
		dir := t.TempDir()
		if err := EnsureEnvGitignore(dir); err != nil {
			t.Fatalf("EnsureEnvGitignore: %v", err)
		}
		data, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
		if err != nil {
			t.Fatalf("reading .gitignore: %v", err)
		}
		if !strings.Contains(string(data), "athanors/*/.env.local") {
			t.Error("expected .env.local pattern in .gitignore")
		}
	})

	t.Run("appends to existing gitignore", func(t *testing.T) {
		dir := t.TempDir()
		existing := "# existing content\n*.tmp\n"
		if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(existing), 0644); err != nil {
			t.Fatal(err)
		}
		if err := EnsureEnvGitignore(dir); err != nil {
			t.Fatalf("EnsureEnvGitignore: %v", err)
		}
		data, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
		if err != nil {
			t.Fatalf("reading .gitignore: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "*.tmp") {
			t.Error("existing content was lost")
		}
		if !strings.Contains(content, "athanors/*/.env.local") {
			t.Error("expected .env.local pattern in .gitignore")
		}
	})

	t.Run("idempotent", func(t *testing.T) {
		dir := t.TempDir()
		if err := EnsureEnvGitignore(dir); err != nil {
			t.Fatal(err)
		}
		if err := EnsureEnvGitignore(dir); err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(filepath.Join(dir, ".gitignore"))
		content := string(data)
		// Count occurrences — should be exactly one
		count := 0
		for i := 0; i+len("athanors/*/.env.local") <= len(content); i++ {
			if content[i:i+len("athanors/*/.env.local")] == "athanors/*/.env.local" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("expected exactly 1 occurrence of pattern, got %d", count)
		}
	})
}

func TestEnvLocalPath(t *testing.T) {
	t.Parallel()
	got := EnvLocalPath("/home/matt/athanor/athanors/myath")
	want := "/home/matt/athanor/athanors/myath/.env.local"
	if got != want {
		t.Errorf("EnvLocalPath = %q, want %q", got, want)
	}
}

func TestInitInstanceCreatesEnvTemplate(t *testing.T) {
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

	if err := InitInstance(home, "env-test", ""); err != nil {
		t.Fatalf("InitInstance: %v", err)
	}

	instDir := InstanceDir(home, "env-test")
	templatePath := filepath.Join(instDir, ".env.local.template")
	if _, err := os.Stat(templatePath); err != nil {
		t.Errorf("expected .env.local.template at %s: %v", templatePath, err)
	}
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
	// Create shared job definitions for per-file sync testing
	for _, job := range []string{"coder", "general"} {
		jobDir := filepath.Join(sharedDir, JobsDir, job)
		if err := os.MkdirAll(jobDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(jobDir, JobFile), []byte("# "+job), 0644); err != nil {
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

		// Verify per-file job symlinks created
		for _, job := range []string{"coder", "general"} {
			jobDir := filepath.Join(instDir, JobsDir, job)
			fi, err := os.Stat(jobDir)
			if err != nil {
				t.Errorf("expected job directory %s: %v", job, err)
				continue
			}
			if !fi.IsDir() {
				t.Errorf("expected %s to be a real directory", jobDir)
			}
			// JOB.md should be a symlink
			target, err := os.Readlink(filepath.Join(jobDir, JobFile))
			if err != nil {
				t.Errorf("expected JOB.md symlink for %s: %v", job, err)
				continue
			}
			expected := filepath.Join(sharedDir, JobsDir, job, JobFile)
			if target != expected {
				t.Errorf("job %s JOB.md -> %q, want %q", job, target, expected)
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

	t.Run("migrates old directory symlink to per-file symlinks", func(t *testing.T) {
		home, sharedDir := setupSyncTestRepo(t)

		name := "migrate-test"
		instDir := InstanceDir(home, name)
		if err := os.MkdirAll(instDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create old-style directory symlink
		oldTarget := filepath.Join(sharedDir, JobsDir)
		if err := os.Symlink(oldTarget, filepath.Join(instDir, JobsDir)); err != nil {
			t.Fatal(err)
		}

		if err := SyncInstance(home, name); err != nil {
			t.Fatalf("SyncInstance: %v", err)
		}

		// jobs/ should now be a real directory, not a symlink
		jobsPath := filepath.Join(instDir, JobsDir)
		if _, err := os.Readlink(jobsPath); err == nil {
			t.Error("jobs/ should be a real directory after migration, but is still a symlink")
		}
		fi, err := os.Stat(jobsPath)
		if err != nil {
			t.Fatalf("jobs/ should exist: %v", err)
		}
		if !fi.IsDir() {
			t.Fatal("jobs/ should be a directory")
		}

		// JOB.md should be a per-file symlink
		target, err := os.Readlink(filepath.Join(instDir, JobsDir, "coder", JobFile))
		if err != nil {
			t.Fatalf("expected JOB.md symlink for coder: %v", err)
		}
		expected := filepath.Join(sharedDir, JobsDir, "coder", JobFile)
		if target != expected {
			t.Errorf("coder JOB.md -> %q, want %q", target, expected)
		}
	})

	t.Run("preserves JOB.local.md during sync", func(t *testing.T) {
		home, sharedDir := setupSyncTestRepo(t)

		name := "local-layer-test"
		instDir := InstanceDir(home, name)
		if err := os.MkdirAll(instDir, 0755); err != nil {
			t.Fatal(err)
		}

		// First sync to establish per-file structure
		if err := SyncInstance(home, name); err != nil {
			t.Fatalf("SyncInstance: %v", err)
		}

		// Write a JOB.local.md layer
		localContent := "# Coder layer for this athanor\nFocus on Go CLI patterns."
		localPath := filepath.Join(instDir, JobsDir, "coder", JobLocalFile)
		if err := os.WriteFile(localPath, []byte(localContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Sync again — local file should survive
		if err := SyncInstance(home, name); err != nil {
			t.Fatalf("SyncInstance (second run): %v", err)
		}

		data, err := os.ReadFile(localPath)
		if err != nil {
			t.Fatalf("JOB.local.md should survive sync: %v", err)
		}
		if string(data) != localContent {
			t.Errorf("JOB.local.md content changed: got %q", string(data))
		}

		// JOB.md symlink should still be correct
		target, err := os.Readlink(filepath.Join(instDir, JobsDir, "coder", JobFile))
		if err != nil {
			t.Fatalf("expected JOB.md symlink: %v", err)
		}
		expected := filepath.Join(sharedDir, JobsDir, "coder", JobFile)
		if target != expected {
			t.Errorf("coder JOB.md -> %q, want %q", target, expected)
		}
	})

	t.Run("syncs new job added to shared", func(t *testing.T) {
		home, sharedDir := setupSyncTestRepo(t)

		name := "new-job-test"
		instDir := InstanceDir(home, name)
		if err := os.MkdirAll(instDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Initial sync
		if err := SyncInstance(home, name); err != nil {
			t.Fatalf("SyncInstance: %v", err)
		}

		// Add a new job to shared
		newJobDir := filepath.Join(sharedDir, JobsDir, "researcher")
		if err := os.MkdirAll(newJobDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(newJobDir, JobFile), []byte("# researcher"), 0644); err != nil {
			t.Fatal(err)
		}

		// Sync again — new job should appear
		if err := SyncInstance(home, name); err != nil {
			t.Fatalf("SyncInstance (after new job): %v", err)
		}

		target, err := os.Readlink(filepath.Join(instDir, JobsDir, "researcher", JobFile))
		if err != nil {
			t.Fatalf("expected JOB.md symlink for researcher: %v", err)
		}
		expected := filepath.Join(newJobDir, JobFile)
		if target != expected {
			t.Errorf("researcher JOB.md -> %q, want %q", target, expected)
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
