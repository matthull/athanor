package cli

import (
	"strings"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "simple intent",
			input:  "Verify tab scaffold PR",
			maxLen: 50,
			want:   "verify-tab-scaffold-pr",
		},
		{
			name:   "intent with numbers and special chars",
			input:  "Verify tab scaffold PR #4718",
			maxLen: 50,
			want:   "verify-tab-scaffold-pr-4718",
		},
		{
			name:   "collapses multiple hyphens",
			input:  "Fix the   widget--loader!!!",
			maxLen: 50,
			want:   "fix-the-widget-loader",
		},
		{
			name:   "trims leading and trailing hyphens",
			input:  "  --hello world--  ",
			maxLen: 50,
			want:   "hello-world",
		},
		{
			name:   "truncates on word boundary",
			input:  "this is a very long intent that should be truncated to fit",
			maxLen: 20,
			want:   "this-is-a-very-long",
		},
		{
			name:   "empty string",
			input:  "",
			maxLen: 50,
			want:   "",
		},
		{
			name:   "only special characters",
			input:  "!@#$%",
			maxLen: 50,
			want:   "",
		},
		{
			name:   "single word within limit",
			input:  "fix",
			maxLen: 50,
			want:   "fix",
		},
		{
			name:   "zero maxLen means no truncation",
			input:  "this should not be truncated at all even though it is very long",
			maxLen: 0,
			want:   "this-should-not-be-truncated-at-all-even-though-it-is-very-long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slugify(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("slugify(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestParseInscribeArgs(t *testing.T) {
	t.Run("full args with all flags", func(t *testing.T) {
		args := []string{"my-athanor", "my-mo", "--intent", "Fix the widget", "--job", "qa-specialist", "--name", "fix-widget", "--muster"}
		ia, err := parseInscribeArgs(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ia.athanorName != "my-athanor" {
			t.Errorf("athanorName = %q, want %q", ia.athanorName, "my-athanor")
		}
		if ia.moName != "my-mo" {
			t.Errorf("moName = %q, want %q", ia.moName, "my-mo")
		}
		if ia.intent != "Fix the widget" {
			t.Errorf("intent = %q, want %q", ia.intent, "Fix the widget")
		}
		if ia.job != "qa-specialist" {
			t.Errorf("job = %q, want %q", ia.job, "qa-specialist")
		}
		if ia.name != "fix-widget" {
			t.Errorf("name = %q, want %q", ia.name, "fix-widget")
		}
		if !ia.muster {
			t.Error("muster = false, want true")
		}
	})

	t.Run("minimal args with job", func(t *testing.T) {
		args := []string{"my-athanor", "my-mo", "--intent", "Do something", "--job", "general"}
		ia, err := parseInscribeArgs(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ia.athanorName != "my-athanor" {
			t.Errorf("athanorName = %q, want %q", ia.athanorName, "my-athanor")
		}
		if ia.job != "general" {
			t.Errorf("job = %q, want %q", ia.job, "general")
		}
		if ia.muster {
			t.Error("muster = true, want false")
		}
	})

	t.Run("missing job errors", func(t *testing.T) {
		args := []string{"my-athanor", "my-mo", "--intent", "Do something"}
		_, err := parseInscribeArgs(args)
		if err == nil {
			t.Fatal("expected error for missing --job")
		}
		if !strings.Contains(err.Error(), "--job is required") {
			t.Errorf("error = %q, want mention of --job", err.Error())
		}
	})

	t.Run("invalid job errors", func(t *testing.T) {
		args := []string{"my-athanor", "my-mo", "--intent", "Do something", "--job", "nonexistent-job"}
		_, err := parseInscribeArgs(args)
		if err == nil {
			t.Fatal("expected error for invalid job")
		}
		if !strings.Contains(err.Error(), "unknown job") {
			t.Errorf("error = %q, want mention of unknown job", err.Error())
		}
	})

	t.Run("flags before positionals", func(t *testing.T) {
		args := []string{"--intent", "Fix it", "--job", "general", "my-athanor", "my-mo"}
		ia, err := parseInscribeArgs(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ia.athanorName != "my-athanor" {
			t.Errorf("athanorName = %q, want %q", ia.athanorName, "my-athanor")
		}
		if ia.intent != "Fix it" {
			t.Errorf("intent = %q, want %q", ia.intent, "Fix it")
		}
	})

	t.Run("one positional errors with MO message", func(t *testing.T) {
		args := []string{"my-athanor", "--intent", "Fix it"}
		_, err := parseInscribeArgs(args)
		if err == nil {
			t.Fatal("expected error for missing MO name")
		}
		if err.Error() != "MO name required" {
			t.Errorf("error = %q, want %q", err.Error(), "MO name required")
		}
	})

	t.Run("no positionals errors with both names", func(t *testing.T) {
		args := []string{"--intent", "Fix it"}
		_, err := parseInscribeArgs(args)
		if err == nil {
			t.Fatal("expected error for missing positionals")
		}
		if !strings.Contains(err.Error(), "athanor name") {
			t.Errorf("error = %q, want mention of athanor name", err.Error())
		}
	})

	t.Run("missing intent errors", func(t *testing.T) {
		args := []string{"my-athanor", "my-mo"}
		_, err := parseInscribeArgs(args)
		if err == nil {
			t.Fatal("expected error for missing --intent")
		}
		if !strings.Contains(err.Error(), "--intent") {
			t.Errorf("error = %q, want mention of --intent", err.Error())
		}
	})

	t.Run("no args errors", func(t *testing.T) {
		_, err := parseInscribeArgs([]string{})
		if err == nil {
			t.Fatal("expected error for no args")
		}
	})
}

func TestBuildOpusContent(t *testing.T) {
	t.Run("basic content without job", func(t *testing.T) {
		content := buildOpusContent("2026-04-22", "my-mo", "", "Fix the widget", "")
		if !strings.Contains(content, "status: charged") {
			t.Error("missing 'status: charged' in frontmatter")
		}
		if !strings.Contains(content, "inscribed: 2026-04-22") {
			t.Error("missing inscribed date")
		}
		if !strings.Contains(content, "magnum_opus: my-mo") {
			t.Error("missing magnum_opus")
		}
		if strings.Contains(content, "job:") {
			t.Error("should not contain job field when empty")
		}
		if !strings.Contains(content, "# Fix the widget") {
			t.Error("missing heading")
		}
		if !strings.Contains(content, "## Intent") {
			t.Error("missing Intent section")
		}
	})

	t.Run("content with job", func(t *testing.T) {
		content := buildOpusContent("2026-04-22", "my-mo", "qa-specialist", "Verify it", "")
		if !strings.Contains(content, "job: qa-specialist") {
			t.Error("missing job field")
		}
	})

	t.Run("content with collaboration context", func(t *testing.T) {
		content := buildOpusContent("2026-04-22", "my-mo", "", "Fix it", "Inscribed by: azer-foo. Whisper back to azer-foo when complete.")
		if !strings.Contains(content, "## Collaboration") {
			t.Error("missing Collaboration section")
		}
		if !strings.Contains(content, "Inscribed by: azer-foo") {
			t.Error("missing collaboration context text")
		}
	})
}

func TestParseCollaborateArgs(t *testing.T) {
	t.Run("full args with job", func(t *testing.T) {
		args := []string{"my-mo", "--intent", "Verify the widget", "--job", "qa-specialist"}
		ca, err := parseCollaborateArgs(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ca.moName != "my-mo" {
			t.Errorf("moName = %q, want %q", ca.moName, "my-mo")
		}
		if ca.intent != "Verify the widget" {
			t.Errorf("intent = %q, want %q", ca.intent, "Verify the widget")
		}
		if ca.job != "qa-specialist" {
			t.Errorf("job = %q, want %q", ca.job, "qa-specialist")
		}
	})

	t.Run("minimal args with job", func(t *testing.T) {
		args := []string{"my-mo", "--intent", "Do something", "--job", "general"}
		ca, err := parseCollaborateArgs(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ca.moName != "my-mo" {
			t.Errorf("moName = %q, want %q", ca.moName, "my-mo")
		}
		if ca.job != "general" {
			t.Errorf("job = %q, want %q", ca.job, "general")
		}
	})

	t.Run("missing job errors", func(t *testing.T) {
		args := []string{"my-mo", "--intent", "Do something"}
		_, err := parseCollaborateArgs(args)
		if err == nil {
			t.Fatal("expected error for missing --job")
		}
		if !strings.Contains(err.Error(), "--job is required") {
			t.Errorf("error = %q, want mention of --job", err.Error())
		}
	})

	t.Run("missing MO errors", func(t *testing.T) {
		args := []string{"--intent", "Do something"}
		_, err := parseCollaborateArgs(args)
		if err == nil {
			t.Fatal("expected error for missing MO name")
		}
		if !strings.Contains(err.Error(), "MO name required") {
			t.Errorf("error = %q, want mention of MO name", err.Error())
		}
	})

	t.Run("missing intent errors", func(t *testing.T) {
		args := []string{"my-mo"}
		_, err := parseCollaborateArgs(args)
		if err == nil {
			t.Fatal("expected error for missing --intent")
		}
		if !strings.Contains(err.Error(), "--intent") {
			t.Errorf("error = %q, want mention of --intent", err.Error())
		}
	})
}
