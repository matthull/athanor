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

	// Create instance directory and opera/
	if err := os.MkdirAll(filepath.Join(instDir, OperaDir), 0755); err != nil {
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

	// Write template magnum-opus.md
	if err := writeMagnumOpusTemplate(instDir, name); err != nil {
		return fmt.Errorf("writing magnum-opus template: %w", err)
	}

	return nil
}

func writeMagnumOpusTemplate(instDir, name string) error {
	content := fmt.Sprintf(`# %s — Magnum Opus

## Goal

[TODO] What is this athanor pursuing? Be specific about the desired end state.

## Abundant Satisfaction

[TODO] What does abundant satisfaction look like? How will you know the goal is met?

## Witnesses

[TODO] Who are the stakeholders? What does success look like from their perspective?

## Pre-loaded Context

[TODO] What does the first azer need to not start from scratch? Discovery findings, references, known open questions, relevant services/files.

## Athanor Structure

`+"```"+`
~/athanor/athanors/%s/
├── AGENTS.md          ← core vocabulary, geas, constraints (all agents read)
├── magnum-opus.md     ← this file
├── marut.md           ← supervisor role
├── azer.md            ← worker role
├── opus.md            ← lifecycle, inscription/discharge protocol
├── muster.md          ← crucible kindling, reforging, monitoring
├── athanor.yml        ← instance configuration
└── opera/
    └── YYYY-MM-DD-<descriptive-name>.md
`+"```"+`
`, name, name)

	return os.WriteFile(filepath.Join(instDir, "magnum-opus.md"), []byte(content), 0644)
}

// ValidateMagnumOpus checks that the magnum-opus.md exists and has real content.
// Returns nil if valid, an error describing the issue otherwise.
func ValidateMagnumOpus(instanceDir string) error {
	path := filepath.Join(instanceDir, "magnum-opus.md")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("magnum-opus.md not found — create it before kindling")
		}
		return fmt.Errorf("reading magnum-opus.md: %w", err)
	}

	content := string(data)
	if strings.Contains(content, "[TODO]") {
		return fmt.Errorf("magnum-opus.md still has [TODO] placeholders — fill them in before kindling")
	}

	return nil
}
