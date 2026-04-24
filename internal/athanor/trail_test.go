package athanor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseOpusFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    OpusFrontmatter
	}{
		{
			name: "full frontmatter",
			content: `---
status: discharged
inscribed: 2026-04-22
discharged: 2026-04-22
magnum_opus: content-performance-dashboards-v2
job: coder
---
# Some Opus Title
`,
			want: OpusFrontmatter{
				Status:     "discharged",
				Inscribed:  "2026-04-22",
				Discharged: "2026-04-22",
				MagnumOpus: "content-performance-dashboards-v2",
				Job:        "coder",
			},
		},
		{
			name: "charged with job",
			content: `---
status: charged
inscribed: 2026-04-24
magnum_opus: collaborative-workshop
job: qa-specialist
---
# Intent
`,
			want: OpusFrontmatter{
				Status:     "charged",
				Inscribed:  "2026-04-24",
				MagnumOpus: "collaborative-workshop",
				Job:        "qa-specialist",
			},
		},
		{
			name: "no job field",
			content: `---
status: discharged
inscribed: 2026-04-22
discharged: 2026-04-22
magnum_opus: some-mo
---
`,
			want: OpusFrontmatter{
				Status:     "discharged",
				Inscribed:  "2026-04-22",
				Discharged: "2026-04-22",
				MagnumOpus: "some-mo",
			},
		},
		{
			name:    "no frontmatter",
			content: "# Just a heading\nSome content.",
			want:    OpusFrontmatter{},
		},
		{
			name:    "empty content",
			content: "",
			want:    OpusFrontmatter{},
		},
		{
			name:    "unclosed frontmatter",
			content: "---\nstatus: charged\n",
			want:    OpusFrontmatter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseOpusFrontmatter(tt.content)
			if got != tt.want {
				t.Errorf("ParseOpusFrontmatter() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestExtractBody(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "normal opus",
			content: `---
status: discharged
---
# Title
Body text here.`,
			want: "\n# Title\nBody text here.",
		},
		{
			name:    "no frontmatter",
			content: "Just plain text.",
			want:    "Just plain text.",
		},
		{
			name:    "frontmatter only",
			content: "---\nstatus: charged\n---",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractBody(tt.content)
			if got != tt.want {
				t.Errorf("ExtractBody() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestScanCollabSignals(t *testing.T) {
	tests := []struct {
		name string
		body string
		want CollabSignals
	}{
		{
			name: "rich collaboration signals",
			body: `## Outcome
Inscribed a qa-specialist to verify the auth module. Used whisper to coordinate with
the peer coder working on frontend. The collaborator provided fresh perspective on
the API design. Mustered an azer for the follow-up work. The dyad between implementation
and review produced better results.`,
			want: CollabSignals{
				WhisperMentions:     1,
				InscribeMentions:    1,
				CollaborateMentions: 1,
				MusterMentions:      1,
				PeerReferences:      3, // "peer" + "collaborator" + "dyad"
				DyadEvidence:        true,
			},
		},
		{
			name: "no signals",
			body: `## Outcome
Fixed the bug in the login flow. Tests pass. PR created.`,
			want: CollabSignals{},
		},
		{
			name: "whisper heavy",
			body: `Sent whisper to marut. Received whisper back. Used whisper to coordinate.`,
			want: CollabSignals{
				WhisperMentions: 3,
			},
		},
		{
			name: "independent review counts as dyad",
			body: `Fired an independent review subagent before discharge.`,
			want: CollabSignals{
				DyadEvidence: true,
			},
		},
		{
			name: "fresh context counts as dyad",
			body: `Inscribed a coder with fresh context to review.`,
			want: CollabSignals{
				InscribeMentions: 1,
				DyadEvidence:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScanCollabSignals(tt.body)
			if got.WhisperMentions != tt.want.WhisperMentions {
				t.Errorf("WhisperMentions = %d, want %d", got.WhisperMentions, tt.want.WhisperMentions)
			}
			if got.InscribeMentions != tt.want.InscribeMentions {
				t.Errorf("InscribeMentions = %d, want %d", got.InscribeMentions, tt.want.InscribeMentions)
			}
			if got.CollaborateMentions != tt.want.CollaborateMentions {
				t.Errorf("CollaborateMentions = %d, want %d", got.CollaborateMentions, tt.want.CollaborateMentions)
			}
			if got.MusterMentions != tt.want.MusterMentions {
				t.Errorf("MusterMentions = %d, want %d", got.MusterMentions, tt.want.MusterMentions)
			}
			if got.PeerReferences != tt.want.PeerReferences {
				t.Errorf("PeerReferences = %d, want %d", got.PeerReferences, tt.want.PeerReferences)
			}
			if got.DyadEvidence != tt.want.DyadEvidence {
				t.Errorf("DyadEvidence = %v, want %v", got.DyadEvidence, tt.want.DyadEvidence)
			}
		})
	}
}

func TestAnalyzeTrail(t *testing.T) {
	// Build a temp athanor home with test opera
	tmpHome := t.TempDir()
	instDir := filepath.Join(tmpHome, AthanorsDir, "test-ath")
	operaDir := filepath.Join(instDir, MagnaOperaDir, "test-mo", "opera")
	moDir := filepath.Join(instDir, MagnaOperaDir, "test-mo")
	if err := os.MkdirAll(operaDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Write a minimal MO file so ListMagnaOpera finds it
	if err := os.WriteFile(filepath.Join(moDir, "test-mo.md"), []byte("# Test MO"), 0644); err != nil {
		t.Fatal(err)
	}

	// Discharged opus with collaboration signals
	opus1 := `---
status: discharged
inscribed: 2026-04-22
discharged: 2026-04-22
magnum_opus: test-mo
job: coder
---
# Fix auth

## Outcome
Inscribed a qa-specialist peer for independent review. Used whisper to coordinate.
`
	// Charged opus with no signals
	opus2 := `---
status: charged
inscribed: 2026-04-24
magnum_opus: test-mo
job: qa-specialist
---
# QA auth flow

## Intent
Verify the auth module works correctly.
`
	// Discharged opus without job
	opus3 := `---
status: discharged
inscribed: 2026-04-20
discharged: 2026-04-21
magnum_opus: test-mo
---
# Old opus without job field

## Outcome
Did some work.
`

	for name, content := range map[string]string{
		"2026-04-22-fix-auth.md": opus1,
		"2026-04-24-qa-auth.md":  opus2,
		"2026-04-20-old-opus.md": opus3,
		"not-an-opus.txt":        "ignored", // non-.md file
	} {
		if err := os.WriteFile(filepath.Join(operaDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a second athanor for filter exclusion testing
	instDir2 := filepath.Join(tmpHome, AthanorsDir, "other-ath")
	operaDir2 := filepath.Join(instDir2, MagnaOperaDir, "other-mo", "opera")
	moDir2 := filepath.Join(instDir2, MagnaOperaDir, "other-mo")
	if err := os.MkdirAll(operaDir2, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moDir2, "other-mo.md"), []byte("# Other MO"), 0644); err != nil {
		t.Fatal(err)
	}
	otherOpus := `---
status: discharged
inscribed: 2026-04-23
discharged: 2026-04-23
magnum_opus: other-mo
job: assessor
---
# Assessment
## Outcome
Did assessment work.
`
	if err := os.WriteFile(filepath.Join(operaDir2, "2026-04-23-assessment.md"), []byte(otherOpus), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("analyze all", func(t *testing.T) {
		report, err := AnalyzeTrail(tmpHome, "", "")
		if err != nil {
			t.Fatalf("AnalyzeTrail: %v", err)
		}
		if report.TotalOpera != 4 {
			t.Errorf("TotalOpera = %d, want 4", report.TotalOpera)
		}
		if report.DischargedOpera != 3 {
			t.Errorf("DischargedOpera = %d, want 3", report.DischargedOpera)
		}
		if report.ChargedOpera != 1 {
			t.Errorf("ChargedOpera = %d, want 1", report.ChargedOpera)
		}
		if report.JobCounts["coder"] != 1 {
			t.Errorf("JobCounts[coder] = %d, want 1", report.JobCounts["coder"])
		}
		if report.JobCounts["qa-specialist"] != 1 {
			t.Errorf("JobCounts[qa-specialist] = %d, want 1", report.JobCounts["qa-specialist"])
		}
		if report.JobCounts["(none)"] != 1 {
			t.Errorf("JobCounts[(none)] = %d, want 1", report.JobCounts["(none)"])
		}
		if report.TotalWhisperMentions != 1 {
			t.Errorf("TotalWhisperMentions = %d, want 1", report.TotalWhisperMentions)
		}
		if report.TotalInscribeMentions != 1 {
			t.Errorf("TotalInscribeMentions = %d, want 1", report.TotalInscribeMentions)
		}
		if report.OperaWithDyadEvidence != 1 {
			t.Errorf("OperaWithDyadEvidence = %d, want 1", report.OperaWithDyadEvidence)
		}
		// Per-MO job diversity
		moKey := "test-ath/test-mo"
		if report.MOJobCounts[moKey] == nil {
			t.Fatalf("MOJobCounts[%s] is nil", moKey)
		}
		if report.MOJobCounts[moKey]["coder"] != 1 {
			t.Errorf("MOJobCounts[%s][coder] = %d, want 1", moKey, report.MOJobCounts[moKey]["coder"])
		}
		if report.MOJobCounts[moKey]["qa-specialist"] != 1 {
			t.Errorf("MOJobCounts[%s][qa-specialist] = %d, want 1", moKey, report.MOJobCounts[moKey]["qa-specialist"])
		}
	})

	t.Run("filter by athanor excludes other", func(t *testing.T) {
		report, err := AnalyzeTrail(tmpHome, "test-ath", "")
		if err != nil {
			t.Fatalf("AnalyzeTrail: %v", err)
		}
		if report.TotalOpera != 3 {
			t.Errorf("TotalOpera = %d, want 3 (should exclude other-ath)", report.TotalOpera)
		}
		if report.JobCounts["assessor"] != 0 {
			t.Errorf("JobCounts[assessor] = %d, want 0 (other-ath should be excluded)", report.JobCounts["assessor"])
		}
	})

	t.Run("filter by athanor and mo excludes other", func(t *testing.T) {
		report, err := AnalyzeTrail(tmpHome, "test-ath", "test-mo")
		if err != nil {
			t.Fatalf("AnalyzeTrail: %v", err)
		}
		if report.TotalOpera != 3 {
			t.Errorf("TotalOpera = %d, want 3", report.TotalOpera)
		}
	})

	t.Run("filter other athanor only", func(t *testing.T) {
		report, err := AnalyzeTrail(tmpHome, "other-ath", "")
		if err != nil {
			t.Fatalf("AnalyzeTrail: %v", err)
		}
		if report.TotalOpera != 1 {
			t.Errorf("TotalOpera = %d, want 1", report.TotalOpera)
		}
		if report.JobCounts["assessor"] != 1 {
			t.Errorf("JobCounts[assessor] = %d, want 1", report.JobCounts["assessor"])
		}
	})

	t.Run("nonexistent athanor returns empty", func(t *testing.T) {
		report, err := AnalyzeTrail(tmpHome, "nonexistent", "")
		if err != nil {
			t.Fatalf("AnalyzeTrail: %v", err)
		}
		if report.TotalOpera != 0 {
			t.Errorf("TotalOpera = %d, want 0", report.TotalOpera)
		}
	})

	t.Run("nonexistent mo returns empty", func(t *testing.T) {
		report, err := AnalyzeTrail(tmpHome, "test-ath", "nonexistent-mo")
		if err != nil {
			t.Fatalf("AnalyzeTrail: %v", err)
		}
		if report.TotalOpera != 0 {
			t.Errorf("TotalOpera = %d, want 0", report.TotalOpera)
		}
	})
}
