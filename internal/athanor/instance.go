package athanor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the per-instance athanor.yml configuration.
type Config struct {
	Name       string `yaml:"name"`
	Project    string `yaml:"project,omitempty"`
	MarutModel string `yaml:"marut_model,omitempty"`
	AzerModel  string `yaml:"azer_model,omitempty"`
}

// Defaults for agent models.
const (
	DefaultMarutModel = "sonnet"
	DefaultAzerModel  = "opus"
)

// EffectiveMarutModel returns the marut model, falling back to default.
func (c *Config) EffectiveMarutModel() string {
	if c.MarutModel != "" {
		return c.MarutModel
	}
	return DefaultMarutModel
}

// EffectiveAzerModel returns the azer model, falling back to default.
func (c *Config) EffectiveAzerModel() string {
	if c.AzerModel != "" {
		return c.AzerModel
	}
	return DefaultAzerModel
}

// ReadConfig reads the athanor.yml for an instance.
func ReadConfig(instanceDir string) (*Config, error) {
	data, err := os.ReadFile(filepath.Join(instanceDir, "athanor.yml"))
	if err != nil {
		return nil, fmt.Errorf("reading athanor.yml: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing athanor.yml: %w", err)
	}
	return &cfg, nil
}

// WriteConfig writes the athanor.yml for an instance.
func WriteConfig(instanceDir string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(filepath.Join(instanceDir, "athanor.yml"), data, 0644)
}

// InitInstance creates a new athanor instance with symlinked shared components.
func InitInstance(home, name, project string) error {
	instDir := InstanceDir(home, name)

	// Check if instance already exists
	if _, err := os.Stat(instDir); err == nil {
		return fmt.Errorf("athanor %q already exists at %s", name, instDir)
	}

	// Create instance directory
	if err := os.MkdirAll(instDir, 0755); err != nil {
		return fmt.Errorf("creating instance directory: %w", err)
	}

	// Symlink shared components
	sharedDir := SharedPath(home)
	for _, f := range SharedFiles {
		src := filepath.Join(sharedDir, f)
		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf("shared component %q not found at %s (run setup first?)", f, src)
		}
		// Use relative symlink: ../shared/<file>
		relSrc := filepath.Join("..", "..", SharedDir, f)
		dst := filepath.Join(instDir, f)
		if err := os.Symlink(relSrc, dst); err != nil {
			return fmt.Errorf("symlinking %s: %w", f, err)
		}
	}

	// Write athanor.yml
	cfg := &Config{
		Name:    name,
		Project: project,
	}
	if err := WriteConfig(instDir, cfg); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	// Create magna-opera directory
	if err := os.MkdirAll(filepath.Join(instDir, MagnaOperaDir), 0755); err != nil {
		return fmt.Errorf("creating magna-opera directory: %w", err)
	}

	return nil
}

// ValidateMagnumOpus checks the legacy magnum-opus.md. Deprecated: use ValidateMO.
func ValidateMagnumOpus(instanceDir string) error {
	return ValidateMO(instanceDir, filepath.Base(instanceDir))
}

// ValidateMO checks that a specific magnum opus exists and has real content.
func ValidateMO(instanceDir, moName string) error {
	path := MagnumOpusPath(instanceDir, moName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("magnum opus %q not found at %s", moName, path)
		}
		return fmt.Errorf("reading magnum opus %q: %w", moName, err)
	}

	content := string(data)
	if strings.Contains(content, "[TODO]") {
		return fmt.Errorf("magnum opus %q still has [TODO] placeholders — fill them in before kindling", moName)
	}

	return nil
}

// HasLegacyMagnumOpus returns true if the instance uses the old single-file format.
func HasLegacyMagnumOpus(instanceDir string) bool {
	moDir := filepath.Join(instanceDir, MagnaOperaDir)
	if _, err := os.Stat(moDir); err == nil {
		return false // magna-opera/ exists, not legacy
	}
	legacyPath := filepath.Join(instanceDir, "magnum-opus.md")
	_, err := os.Stat(legacyPath)
	return err == nil
}

// MagnumOpusPath returns the filesystem path to a specific MO file.
// For legacy instances (no magna-opera/ dir), returns magnum-opus.md.
// For multi-MO instances, returns magna-opera/<moName>/<moName>.md.
func MagnumOpusPath(instanceDir, moName string) string {
	if HasLegacyMagnumOpus(instanceDir) {
		return filepath.Join(instanceDir, "magnum-opus.md")
	}
	return filepath.Join(instanceDir, MagnaOperaDir, moName, moName+".md")
}

// OperaPath returns the filesystem path to a specific MO's opera directory.
// Returns magna-opera/<moName>/opera.
func OperaPath(instanceDir, moName string) string {
	return filepath.Join(instanceDir, MagnaOperaDir, moName, "opera")
}

// ListMagnaOpera returns the names of all magna opera in an instance.
// For legacy instances, returns a single-element list with the instance name.
func ListMagnaOpera(instanceDir string) ([]string, error) {
	moDir := filepath.Join(instanceDir, MagnaOperaDir)
	entries, err := os.ReadDir(moDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Check for legacy magnum-opus.md
			if HasLegacyMagnumOpus(instanceDir) {
				return []string{filepath.Base(instanceDir)}, nil
			}
			return nil, nil
		}
		return nil, fmt.Errorf("listing magna opera: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// ReadOpusMO reads the magnum_opus field from an opus file's YAML frontmatter.
func ReadOpusMO(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return ""
	}
	end := strings.Index(content[3:], "---")
	if end < 0 {
		return ""
	}
	frontmatter := content[3 : 3+end]
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "magnum_opus:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "magnum_opus:"))
		}
	}
	return ""
}

// WriteMOTemplate writes a template magnum opus file to magna-opera/<moName>/<moName>.md.
func WriteMOTemplate(instanceDir, moName string) error {
	content := fmt.Sprintf(`# %s — Magnum Opus

## Goal

[TODO] What is this magnum opus pursuing? Be specific about the desired end state.

## Abundant Satisfaction

[TODO] What does abundant satisfaction look like? How will you know the goal is met?

## Witnesses

[TODO] Who are the stakeholders? What does success look like from their perspective?

## Marut Directives

(Leave empty for default behavior — the marut pursues the goal without constraints. Add directives to scope what the marut focuses on, avoids, or stops at.)

## Tempering

(Empty by default. Transient guidance — "the weather today." Updated by the marut during artifex conversation. Always timestamped.)

## Pre-loaded Context

[TODO] What does the first azer need to not start from scratch? Discovery findings, references, known open questions, relevant services/files.
`, moName)

	moDir := filepath.Join(instanceDir, MagnaOperaDir, moName)
	if err := os.MkdirAll(moDir, 0755); err != nil {
		return fmt.Errorf("creating MO directory: %w", err)
	}
	// Create opera subdirectory for this MO
	if err := os.MkdirAll(filepath.Join(moDir, "opera"), 0755); err != nil {
		return fmt.Errorf("creating MO opera directory: %w", err)
	}
	return os.WriteFile(filepath.Join(moDir, moName+".md"), []byte(content), 0644)
}

// MarutCrucibleName returns the crucible name for a marut.
// For legacy (single MO), pass empty moName to get "marut-<athanor>".
// For multi-MO, pass the MO name to get "marut-<athanor>-<mo>".
func MarutCrucibleName(athanorName, moName string) string {
	if moName == "" {
		return fmt.Sprintf("marut-%s", athanorName)
	}
	return fmt.Sprintf("marut-%s-%s", athanorName, moName)
}
