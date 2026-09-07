package athanor

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestRoleCrucibleName(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		athanorName string
		moName      string
		want        string
	}{
		{
			name:        "no mo — marut",
			role:        "marut",
			athanorName: "athtest",
			moName:      "",
			want:        "marut-athtest",
		},
		{
			name:        "no mo — perceiver",
			role:        "perceiver",
			athanorName: "athtest",
			moName:      "",
			want:        "perceiver-athtest",
		},
		{
			name:        "with mo — marut",
			role:        "marut",
			athanorName: "athtest",
			moName:      "featurex",
			want:        "marut-athtest-featurex",
		},
		{
			name:        "with mo — attendant",
			role:        "attendant",
			athanorName: "athtest",
			moName:      "featurex",
			want:        "attendant-athtest-featurex",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RoleCrucibleName(tc.role, tc.athanorName, tc.moName)
			if got != tc.want {
				t.Errorf("RoleCrucibleName(%q, %q, %q) = %q, want %q",
					tc.role, tc.athanorName, tc.moName, got, tc.want)
			}
		})
	}
}

func TestRoleCrucibleName_MatchesMarutPattern(t *testing.T) {
	// RoleCrucibleName("marut", ...) must match MarutCrucibleName exactly.
	cases := []struct {
		athanor string
		mo      string
	}{
		{"a", ""},
		{"a", "mo1"},
		{"athanor-x", ""},
		{"athanor-x", "my-mo"},
	}
	for _, c := range cases {
		got := RoleCrucibleName("marut", c.athanor, c.mo)
		want := MarutCrucibleName(c.athanor, c.mo)
		if got != want {
			t.Errorf("RoleCrucibleName(marut, %q, %q) = %q, MarutCrucibleName = %q — must match",
				c.athanor, c.mo, got, want)
		}
	}
}

func TestKindledDir(t *testing.T) {
	instDir := "/some/athanor/inst"
	moName := "mymo"
	want := filepath.Join(instDir, MagnaOperaDir, moName, "kindled")
	got := KindledDir(instDir, moName)
	if got != want {
		t.Errorf("KindledDir = %q, want %q", got, want)
	}
}

// moScaffold sets up a realistic instance dir with a magna-opera/<mo>/ layout.
func moScaffold(t *testing.T, moName string) string {
	t.Helper()
	instDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(instDir, MagnaOperaDir, moName), 0755); err != nil {
		t.Fatalf("scaffold: %v", err)
	}
	return instDir
}

func TestSetKindled_AndIsKindled(t *testing.T) {
	instDir := moScaffold(t, "mo1")

	if IsKindled(instDir, "mo1", "marut") {
		t.Fatalf("IsKindled returned true before set")
	}

	if err := SetKindled(instDir, "mo1", "marut"); err != nil {
		t.Fatalf("SetKindled: %v", err)
	}
	if !IsKindled(instDir, "mo1", "marut") {
		t.Fatalf("IsKindled returned false after set")
	}

	// Verify file exists on disk at the expected path.
	path := filepath.Join(KindledDir(instDir, "mo1"), "marut")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at %s: %v", path, err)
	}
}

func TestSetKindled_CreatesDir(t *testing.T) {
	// Instance dir exists but magna-opera/<mo>/kindled does NOT — SetKindled must create it.
	instDir := moScaffold(t, "mo1")
	kindled := KindledDir(instDir, "mo1")
	if _, err := os.Stat(kindled); !os.IsNotExist(err) {
		t.Fatalf("precondition: kindled dir should not exist yet")
	}

	if err := SetKindled(instDir, "mo1", "perceiver"); err != nil {
		t.Fatalf("SetKindled: %v", err)
	}
	if _, err := os.Stat(kindled); err != nil {
		t.Fatalf("kindled dir not created: %v", err)
	}
}

func TestSetKindled_Idempotent(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	if err := SetKindled(instDir, "mo1", "marut"); err != nil {
		t.Fatalf("SetKindled (1): %v", err)
	}
	if err := SetKindled(instDir, "mo1", "marut"); err != nil {
		t.Fatalf("SetKindled (2): %v", err)
	}
	if !IsKindled(instDir, "mo1", "marut") {
		t.Fatalf("IsKindled returned false after double-set")
	}
}

func TestClearKindled(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	if err := SetKindled(instDir, "mo1", "marut"); err != nil {
		t.Fatalf("SetKindled: %v", err)
	}

	if err := ClearKindled(instDir, "mo1", "marut"); err != nil {
		t.Fatalf("ClearKindled: %v", err)
	}
	if IsKindled(instDir, "mo1", "marut") {
		t.Fatalf("IsKindled returned true after clear")
	}
}

func TestClearKindled_Idempotent(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	// Clearing a role that was never set must not error.
	if err := ClearKindled(instDir, "mo1", "marut"); err != nil {
		t.Fatalf("ClearKindled on nonexistent: %v", err)
	}
	// Clearing in a dir that does not exist at all must not error.
	if err := ClearKindled(instDir, "mo-nonexistent", "marut"); err != nil {
		t.Fatalf("ClearKindled in missing dir: %v", err)
	}
}

func TestClearAllKindled(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	for _, role := range []string{"marut", "perceiver", "attendant"} {
		if err := SetKindled(instDir, "mo1", role); err != nil {
			t.Fatalf("SetKindled(%s): %v", role, err)
		}
	}

	if err := ClearAllKindled(instDir, "mo1"); err != nil {
		t.Fatalf("ClearAllKindled: %v", err)
	}

	got, err := ListKindled(instDir, "mo1")
	if err != nil {
		t.Fatalf("ListKindled: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty after ClearAllKindled, got %v", got)
	}
}

func TestClearAllKindled_MissingDir(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	// Kindled dir never created — must not error.
	if err := ClearAllKindled(instDir, "mo1"); err != nil {
		t.Fatalf("ClearAllKindled on missing dir: %v", err)
	}
	// MO dir entirely missing — must not error.
	if err := ClearAllKindled(instDir, "mo-nonexistent"); err != nil {
		t.Fatalf("ClearAllKindled on missing mo: %v", err)
	}
}

func TestListKindled(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	for _, role := range []string{"marut", "perceiver", "attendant"} {
		if err := SetKindled(instDir, "mo1", role); err != nil {
			t.Fatalf("SetKindled(%s): %v", role, err)
		}
	}

	got, err := ListKindled(instDir, "mo1")
	if err != nil {
		t.Fatalf("ListKindled: %v", err)
	}
	sort.Strings(got)
	want := []string{"attendant", "marut", "perceiver"}
	if len(got) != len(want) {
		t.Fatalf("ListKindled = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ListKindled = %v, want %v", got, want)
		}
	}
}

func TestListKindled_MissingDir(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	got, err := ListKindled(instDir, "mo1")
	if err != nil {
		t.Fatalf("ListKindled on missing dir: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}

	got, err = ListKindled(instDir, "mo-nonexistent")
	if err != nil {
		t.Fatalf("ListKindled on missing mo: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestIsKindled_MissingDir(t *testing.T) {
	instDir := moScaffold(t, "mo1")
	if IsKindled(instDir, "mo1", "marut") {
		t.Fatalf("IsKindled on missing dir returned true")
	}
	if IsKindled(instDir, "mo-nonexistent", "marut") {
		t.Fatalf("IsKindled on missing mo returned true")
	}
}
