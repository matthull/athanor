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
// Note: jobs/ is NOT in this list — it uses per-file symlinks to allow
// athanor-specific JOB.local.md layers. See syncJobs in instance.go.
var SharedDirs = []string{}

// JobFile is the shared job definition filename.
const JobFile = "JOB.md"

// JobLocalFile is the optional per-athanor job layer filename.
// Purely additive guidance loaded after JOB.md — no frontmatter, no overrides.
const JobLocalFile = "JOB.local.md"

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

// JobInfo holds the structured frontmatter from a JOB.md file.
type JobInfo struct {
	Name    string
	Summary string
	When    []string
	Model   string
}

// ReadJobInfo reads structured frontmatter from a JOB.md file.
func ReadJobInfo(jobsDir, jobName string) (JobInfo, error) {
	info := JobInfo{Name: jobName}
	path := filepath.Join(jobsDir, jobName, "JOB.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return info, fmt.Errorf("reading job %q: %w", jobName, err)
	}
	content := string(data)
	if !hasYAMLFrontmatter(content) {
		return info, nil
	}
	fm := extractFrontmatter(content)
	info.Summary = extractField(fm, "summary:")
	info.Model = extractField(fm, "model:")
	info.When = extractListField(fm, "when:")
	return info, nil
}

// ListJobInfos returns JobInfo for all available jobs in the shared registry.
func ListJobInfos() ([]JobInfo, error) {
	shared, err := SharedPath()
	if err != nil {
		return nil, err
	}
	jobsDir := filepath.Join(shared, JobsDir)
	names, err := listJobNames(jobsDir)
	if err != nil {
		return nil, err
	}
	var infos []JobInfo
	for _, name := range names {
		info, err := ReadJobInfo(jobsDir, name)
		if err != nil {
			// Skip unreadable jobs but still list them
			infos = append(infos, JobInfo{Name: name})
			continue
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// ReadJobDetail reads the full content of a JOB.md file.
func ReadJobDetail(jobName string) (JobInfo, string, error) {
	shared, err := SharedPath()
	if err != nil {
		return JobInfo{}, "", err
	}
	jobsDir := filepath.Join(shared, JobsDir)
	info, err := ReadJobInfo(jobsDir, jobName)
	if err != nil {
		return info, "", err
	}
	path := filepath.Join(jobsDir, jobName, "JOB.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return info, "", err
	}
	return info, string(data), nil
}

// StripFrontmatter removes YAML frontmatter from markdown content.
func StripFrontmatter(content string) string {
	if !hasYAMLFrontmatter(content) {
		return content
	}
	end := findEndOfFrontmatter(content[3:])
	if end < 0 {
		return content
	}
	// Skip past closing "---" and any trailing newline
	rest := content[3+end+3:]
	for len(rest) > 0 && (rest[0] == '\n' || rest[0] == '\r') {
		rest = rest[1:]
	}
	return rest
}

func hasYAMLFrontmatter(content string) bool {
	return len(content) >= 3 && content[:3] == "---"
}

func extractFrontmatter(content string) string {
	end := findEndOfFrontmatter(content[3:])
	if end < 0 {
		return ""
	}
	return content[3 : 3+end]
}

func findEndOfFrontmatter(s string) int {
	// Find the next "---" that starts a line
	for i := 0; i < len(s); i++ {
		if i == 0 || s[i-1] == '\n' {
			if len(s[i:]) >= 3 && s[i:i+3] == "---" {
				return i
			}
		}
	}
	return -1
}

func extractField(fm, prefix string) string {
	for _, line := range splitLines(fm) {
		trimmed := trimSpace(line)
		if len(trimmed) > len(prefix) && trimmed[:len(prefix)] == prefix {
			return trimSpace(trimmed[len(prefix):])
		}
	}
	return ""
}

func extractListField(fm, prefix string) []string {
	lines := splitLines(fm)
	var result []string
	inList := false
	for _, line := range lines {
		trimmed := trimSpace(line)
		if trimmed == prefix {
			inList = true
			continue
		}
		if inList {
			if len(trimmed) >= 2 && trimmed[0] == '-' && trimmed[1] == ' ' {
				item := trimSpace(trimmed[2:])
				// Strip surrounding quotes
				if len(item) >= 2 && item[0] == '"' && item[len(item)-1] == '"' {
					item = item[1 : len(item)-1]
				}
				result = append(result, item)
			} else {
				break // End of list
			}
		}
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// listJobNames returns sorted job directory names from a jobs directory.
func listJobNames(jobsDir string) ([]string, error) {
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
