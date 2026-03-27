package cli

import (
	"strings"
	"testing"
	"time"
)

func TestExtractSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		heading string
		want    string
	}{
		{
			name: "simple goal",
			content: `# Test MO

## Goal

AppSignal channel stays clean, actionable errors get actioned.

## Abundant Satisfaction

All clear.
`,
			heading: "Goal",
			want:    "AppSignal channel stays clean, actionable errors get actioned.",
		},
		{
			name: "multi-line paragraph",
			content: `## Goal

First line of the goal
continues on second line.

## Next Section
`,
			heading: "Goal",
			want:    "First line of the goal continues on second line.",
		},
		{
			name:    "section not found",
			content: "## Other\n\nSome content.\n",
			heading: "Goal",
			want:    "",
		},
		{
			name:    "section at end of file no trailing newline",
			content: "## Outcome\n\nThe fix was deployed.",
			heading: "Outcome",
			want:    "The fix was deployed.",
		},
		{
			name:    "empty section",
			content: "## Goal\n\n## Next\n",
			heading: "Goal",
			want:    "",
		},
		{
			name:    "truncates long content",
			content: "## Goal\n\n" + strings.Repeat("x", 200) + "\n",
			heading: "Goal",
			want:    strings.Repeat("x", 117) + "...",
		},
		{
			name: "with TODO placeholder",
			content: `## Goal

[TODO] What is this magnum opus pursuing?

## Witnesses
`,
			heading: "Goal",
			want:    "[TODO] What is this magnum opus pursuing?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSection(tt.content, tt.heading)
			if got != tt.want {
				t.Errorf("extractSection(%q) = %q, want %q", tt.heading, got, tt.want)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "simple title",
			content: "# Fix NoMethodError in ProofRecommendationService\n\n## Goal\n",
			want:    "Fix NoMethodError in ProofRecommendationService",
		},
		{
			name:    "with Opus prefix",
			content: "# Opus: Assess and orient\n\n## Goal\n",
			want:    "Assess and orient",
		},
		{
			name:    "ignores ## headings",
			content: "## Goal\n\nSomething.\n",
			want:    "",
		},
		{
			name:    "title after frontmatter",
			content: "---\nstatus: charged\n---\n# Fix the bug\n\n## Goal\n",
			want:    "Fix the bug",
		},
		{
			name:    "no title",
			content: "Just some text without headings.\n",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTitle(tt.content)
			if got != tt.want {
				t.Errorf("extractTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadFrontmatterField(t *testing.T) {
	fm := `---
status: charged
inscribed: 2026-03-25
magnum_opus: bugsnag
---
# Title
`
	tests := []struct {
		field string
		want  string
	}{
		{"status", "charged"},
		{"inscribed", "2026-03-25"},
		{"magnum_opus", "bugsnag"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			got := readFrontmatterField(fm, tt.field)
			if got != tt.want {
				t.Errorf("readFrontmatterField(%q) = %q, want %q", tt.field, got, tt.want)
			}
		})
	}

	t.Run("no frontmatter", func(t *testing.T) {
		got := readFrontmatterField("# Just a title\n", "status")
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("unclosed frontmatter", func(t *testing.T) {
		got := readFrontmatterField("---\nstatus: charged\n", "status")
		if got != "" {
			t.Errorf("expected empty for unclosed frontmatter, got %q", got)
		}
	})
}

func TestSplitOpusFilename(t *testing.T) {
	tests := []struct {
		filename string
		wantDate string
		wantName string
	}{
		{"2026-03-25-fix-nil-error", "2026-03-25", "fix-nil-error"},
		{"2026-03-25-assess-and-orient", "2026-03-25", "assess-and-orient"},
		{"no-date-prefix", "", "no-date-prefix"},
		{"short", "", "short"},
		{"2026-03-25-x", "2026-03-25", "x"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			date, name := splitOpusFilename(tt.filename)
			if date != tt.wantDate {
				t.Errorf("date = %q, want %q", date, tt.wantDate)
			}
			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
		})
	}
}

func TestIsRecentDate(t *testing.T) {
	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	cutoff := now.Add(-48 * time.Hour) // 2026-03-25 12:00

	tests := []struct {
		name    string
		dateStr string
		want    bool
	}{
		{"today", "2026-03-27", true},
		{"yesterday", "2026-03-26", true},
		{"two days ago", "2026-03-25", false}, // 2026-03-25 00:00 < 2026-03-25 12:00 (cutoff)
		{"three days ago", "2026-03-24", false},
		{"invalid", "not-a-date", false},
		{"empty", "", false},
		{"short", "2026", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRecentDate(tt.dateStr, cutoff)
			if got != tt.want {
				t.Errorf("isRecentDate(%q) = %v, want %v", tt.dateStr, got, tt.want)
			}
		})
	}
}

func TestFormatShortDate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2026-03-25", "Mar 25"},
		{"2026-01-01", "Jan 1"},
		{"2026-12-31", "Dec 31"},
		{"invalid", "invalid"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := formatShortDate(tt.input)
			if got != tt.want {
				t.Errorf("formatShortDate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStateIcon(t *testing.T) {
	tests := []struct {
		state string
		want  string
	}{
		{"active", "●"},
		{"idle", "◐"},
		{"permission", "⚠"},
		{"exhausted", "✖"},
		{"dead", "✖"},
		{"unknown", "?"},
	}

	for _, tt := range tests {
		t.Run(tt.state, func(t *testing.T) {
			got := stateIcon(tt.state)
			if got != tt.want {
				t.Errorf("stateIcon(%q) = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestSliceContains(t *testing.T) {
	slice := []string{"foo", "bar", "baz"}

	if !sliceContains(slice, "bar") {
		t.Error("expected to find 'bar'")
	}
	if sliceContains(slice, "qux") {
		t.Error("did not expect to find 'qux'")
	}
	if sliceContains(nil, "foo") {
		t.Error("nil slice should not contain anything")
	}
}
