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

	// SharedDir is the subdirectory containing shared agent definitions in the source repo.
	SharedDir = "shared"

	// MagnaOperaDir is the subdirectory containing magna opera (top-level goals).
	MagnaOperaDir = "magna-opera"
)

// SharedFiles are the component files symlinked from shared/ into each instance.
var SharedFiles = []string{
	"AGENTS.md",
	"attendant.md",
	"azer.md",
	"marut.md",
	"muster.md",
	"opus.md",
	"perceiver.md",
}

// SharedDirs are directories symlinked from shared/ into each instance.
var SharedDirs = []string{
	"jobs",
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

// RepoDir resolves the athanor source repository path.
// Checks $ATHANOR_REPO first, then falls back to ~/code/athanor.
func RepoDir() (string, error) {
	if r := os.Getenv("ATHANOR_REPO"); r != "" {
		return expandHome(r)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	return filepath.Join(home, "code", "athanor"), nil
}

// SharedPath returns the path to the shared agent definitions in the source repo.
func SharedPath() (string, error) {
	repo, err := RepoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(repo, SharedDir), nil
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
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", d, err)
		}
	}
	return nil
}

// JobsDir is the subdirectory containing job definitions in the source repo.
const JobsDir = "jobs"

// ValidateJob checks that a named job exists in the shared job registry.
// Returns nil if shared/jobs/<name>/JOB.md exists, or an error describing
// what's wrong and listing available jobs.
func ValidateJob(jobName string) error {
	shared, err := SharedPath()
	if err != nil {
		return fmt.Errorf("resolving shared path: %w", err)
	}
	jobPath := filepath.Join(shared, JobsDir, jobName, "JOB.md")
	if _, err := os.Stat(jobPath); err != nil {
		if os.IsNotExist(err) {
			available, _ := ListJobs()
			return fmt.Errorf("unknown job %q (available: %s)", jobName, formatJobList(available))
		}
		return fmt.Errorf("checking job %q: %w", jobName, err)
	}
	return nil
}

// ListJobs returns the names of all available jobs in the shared registry.
func ListJobs() ([]string, error) {
	shared, err := SharedPath()
	if err != nil {
		return nil, err
	}
	jobsDir := filepath.Join(shared, JobsDir)
	entries, err := os.ReadDir(jobsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("listing jobs: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			jmd := filepath.Join(jobsDir, e.Name(), "JOB.md")
			if _, err := os.Stat(jmd); err == nil {
				names = append(names, e.Name())
			}
		}
	}
	return names, nil
}

func formatJobList(jobs []string) string {
	if len(jobs) == 0 {
		return "none"
	}
	result := ""
	for i, j := range jobs {
		if i > 0 {
			result += ", "
		}
		result += j
	}
	return result
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
