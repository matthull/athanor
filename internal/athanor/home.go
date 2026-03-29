// Package athanor provides operations on the athanor home directory (~~/athanor/).
package athanor

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultHome is the default athanor home directory.
	DefaultHome = "~/athanor"

	// AthanorsDir is the subdirectory containing athanor instances.
	AthanorsDir = "athanors"

	// SharedDir is the subdirectory containing shared components (agent roles, protocols).
	SharedDir = "shared"

	// MagnaOperaDir is the subdirectory containing magna opera (top-level goals).
	MagnaOperaDir = "magna-opera"
)

// SharedFiles are the component files symlinked from shared/ into each instance.
var SharedFiles = []string{
	"AGENTS.md",
	"azer.md",
	"marut.md",
	"muster.md",
	"opus.md",
}

// Home resolves the athanor home directory path.
// Checks $ATHANOR_HOME first, then falls back to ~/athanor.
func Home() (string, error) {
	if h := os.Getenv("ATHANOR_HOME"); h != "" {
		return expandHome(h)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	return filepath.Join(home, "athanor"), nil
}

// InstanceDir returns the path to a named athanor instance.
func InstanceDir(home, name string) string {
	return filepath.Join(home, AthanorsDir, name)
}

// SharedPath returns the path to the shared components directory.
func SharedPath(home string) string {
	return filepath.Join(home, SharedDir)
}

// ListInstances returns the names of all athanor instances.
func ListInstances(home string) ([]string, error) {
	dir := filepath.Join(home, AthanorsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("listing instances: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// EnsureHome creates the athanor home directory structure if it doesn't exist.
func EnsureHome(home string) error {
	dirs := []string{
		home,
		filepath.Join(home, AthanorsDir),
		filepath.Join(home, SharedDir),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", d, err)
		}
	}
	return nil
}

// expandHome expands a leading ~ to the user's home directory.
func expandHome(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}
	if path[0] != '~' {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, path[1:]), nil
}
