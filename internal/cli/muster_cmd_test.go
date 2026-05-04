package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matthull/athanor/internal/athanor"
)

func TestFormulaBootClause(t *testing.T) {
	t.Run("empty formula returns empty clause", func(t *testing.T) {
		t.Parallel()
		clause, err := formulaBootClause(t.TempDir(), "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if clause != "" {
			t.Errorf("clause = %q, want empty", clause)
		}
	})

	t.Run("valid formula returns clause", func(t *testing.T) {
		t.Parallel()
		instDir := t.TempDir()
		formulaDir := filepath.Join(instDir, athanor.FormulaeDir, "coding-dyad")
		if err := os.MkdirAll(formulaDir, 0755); err != nil {
			t.Fatal(err)
		}
		formulaPath := filepath.Join(formulaDir, athanor.FormulaFile)
		if err := os.WriteFile(formulaPath, []byte("# Coding Dyad"), 0644); err != nil {
			t.Fatal(err)
		}

		clause, err := formulaBootClause(instDir, "coding-dyad")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.HasPrefix(clause, ", then read ") {
			t.Errorf("clause = %q, want prefix ', then read '", clause)
		}
		if !strings.Contains(clause, formulaPath) {
			t.Errorf("clause = %q, want path %q", clause, formulaPath)
		}
	})

	t.Run("missing formula errors", func(t *testing.T) {
		t.Parallel()
		_, err := formulaBootClause(t.TempDir(), "nonexistent")
		if err == nil {
			t.Fatal("expected error for missing formula")
		}
		if !strings.Contains(err.Error(), "formula definition not found") {
			t.Errorf("error = %q, want 'formula definition not found'", err.Error())
		}
	})
}
