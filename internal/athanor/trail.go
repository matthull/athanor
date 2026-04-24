package athanor

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// OpusFrontmatter holds parsed YAML frontmatter from an opus file.
type OpusFrontmatter struct {
	Status     string
	Inscribed  string
	Discharged string
	MagnumOpus string
	Job        string
}

// CollabSignals holds collaboration signals extracted from an opus body.
type CollabSignals struct {
	WhisperMentions     int
	InscribeMentions    int
	CollaborateMentions int
	MusterMentions      int
	PeerReferences      int // mentions of "peer", "collaborator", "dyad"
	DyadEvidence        bool
}

// OpusAnalysis combines frontmatter and collaboration signals for one opus.
type OpusAnalysis struct {
	Path        string
	Name        string
	Athanor     string
	MO          string
	Frontmatter OpusFrontmatter
	Signals     CollabSignals
}

// TrailReport aggregates collaboration topology across opera.
type TrailReport struct {
	TotalOpera      int
	DischargedOpera int
	ChargedOpera    int

	// Job diversity
	JobCounts map[string]int

	// Collaboration signals (aggregated)
	TotalWhisperMentions     int
	TotalInscribeMentions    int
	TotalCollaborateMentions int
	TotalMusterMentions      int
	TotalPeerReferences      int
	OperaWithDyadEvidence    int

	// Per-opus detail
	Opera []OpusAnalysis
}

// ParseOpusFrontmatter extracts YAML frontmatter fields from an opus file's content.
// Returns zero-value OpusFrontmatter if no valid frontmatter found.
func ParseOpusFrontmatter(content string) OpusFrontmatter {
	var fm OpusFrontmatter

	if !strings.HasPrefix(content, "---") {
		return fm
	}

	end := strings.Index(content[3:], "---")
	if end < 0 {
		return fm
	}

	frontmatter := content[3 : 3+end]
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if k, v := parseFMLine(line, "status:"); v != "" {
			_ = k
			fm.Status = v
		} else if _, v := parseFMLine(line, "inscribed:"); v != "" {
			fm.Inscribed = v
		} else if _, v := parseFMLine(line, "discharged:"); v != "" {
			fm.Discharged = v
		} else if _, v := parseFMLine(line, "magnum_opus:"); v != "" {
			fm.MagnumOpus = v
		} else if _, v := parseFMLine(line, "job:"); v != "" {
			fm.Job = v
		}
	}

	return fm
}

// parseFMLine checks if a line starts with the given prefix and returns the value.
func parseFMLine(line, prefix string) (string, string) {
	if strings.HasPrefix(line, prefix) {
		return prefix, strings.TrimSpace(strings.TrimPrefix(line, prefix))
	}
	return "", ""
}

// Patterns for detecting collaboration signals in opus body text.
var (
	whisperPattern     = regexp.MustCompile(`(?i)\bwhisper\b`)
	inscribePattern    = regexp.MustCompile(`(?i)\binscrib(e|ed|ing)\b`)
	collaboratePattern = regexp.MustCompile(`(?i)\bcollaborat(e|ed|ing|or|ors|ion)\b`)
	musterPattern      = regexp.MustCompile(`(?i)\bmuster(ed|ing)?\b`)
	peerPattern        = regexp.MustCompile(`(?i)\b(peer|collaborator|dyad)\b`)
	dyadPattern        = regexp.MustCompile(`(?i)\b(dyad|independent (review|perspective|verification)|fresh (eyes|context|perspective))\b`)
)

// ScanCollabSignals scans opus body text (below the frontmatter delimiter)
// for collaboration signals.
func ScanCollabSignals(body string) CollabSignals {
	var cs CollabSignals

	cs.WhisperMentions = len(whisperPattern.FindAllString(body, -1))
	cs.InscribeMentions = len(inscribePattern.FindAllString(body, -1))
	cs.CollaborateMentions = len(collaboratePattern.FindAllString(body, -1))
	cs.MusterMentions = len(musterPattern.FindAllString(body, -1))
	cs.PeerReferences = len(peerPattern.FindAllString(body, -1))
	cs.DyadEvidence = dyadPattern.MatchString(body)

	return cs
}

// ExtractBody returns the text below the YAML frontmatter delimiter.
// If no frontmatter, returns the entire content.
func ExtractBody(content string) string {
	if !strings.HasPrefix(content, "---") {
		return content
	}
	end := strings.Index(content[3:], "---")
	if end < 0 {
		return content
	}
	// Skip past the closing ---
	bodyStart := 3 + end + 3
	if bodyStart >= len(content) {
		return ""
	}
	return content[bodyStart:]
}

// AnalyzeOpus parses a single opus file and returns its analysis.
func AnalyzeOpus(path, athanorName, moName string) (OpusAnalysis, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return OpusAnalysis{}, err
	}
	content := string(data)

	name := strings.TrimSuffix(filepath.Base(path), ".md")

	return OpusAnalysis{
		Path:        path,
		Name:        name,
		Athanor:     athanorName,
		MO:          moName,
		Frontmatter: ParseOpusFrontmatter(content),
		Signals:     ScanCollabSignals(ExtractBody(content)),
	}, nil
}

// AnalyzeTrail walks all opera across specified athanors and MOs,
// producing an aggregated TrailReport.
// If athanorName is empty, scans all athanors.
// If moName is empty, scans all MOs within each athanor.
func AnalyzeTrail(home, athanorName, moName string) (*TrailReport, error) {
	report := &TrailReport{
		JobCounts: make(map[string]int),
	}

	var instances []string
	if athanorName != "" {
		instances = []string{athanorName}
	} else {
		var err error
		instances, err = ListInstances(home)
		if err != nil {
			return nil, err
		}
	}

	for _, inst := range instances {
		instDir := InstanceDir(home, inst)

		var mos []string
		if moName != "" {
			mos = []string{moName}
		} else {
			var err error
			mos, err = ListMagnaOpera(instDir)
			if err != nil {
				continue
			}
		}

		for _, mo := range mos {
			operaDir := OperaPath(instDir, mo)
			entries, err := os.ReadDir(operaDir)
			if err != nil {
				continue
			}

			for _, e := range entries {
				if !strings.HasSuffix(e.Name(), ".md") {
					continue
				}
				opusPath := filepath.Join(operaDir, e.Name())

				analysis, err := AnalyzeOpus(opusPath, inst, mo)
				if err != nil {
					continue
				}

				report.Opera = append(report.Opera, analysis)
				report.TotalOpera++

				switch analysis.Frontmatter.Status {
				case "discharged", "assessed":
					report.DischargedOpera++
				case "charged":
					report.ChargedOpera++
				}

				// Job diversity
				job := analysis.Frontmatter.Job
				if job == "" {
					job = "(none)"
				}
				report.JobCounts[job]++

				// Aggregate signals
				report.TotalWhisperMentions += analysis.Signals.WhisperMentions
				report.TotalInscribeMentions += analysis.Signals.InscribeMentions
				report.TotalCollaborateMentions += analysis.Signals.CollaborateMentions
				report.TotalMusterMentions += analysis.Signals.MusterMentions
				report.TotalPeerReferences += analysis.Signals.PeerReferences
				if analysis.Signals.DyadEvidence {
					report.OperaWithDyadEvidence++
				}
			}
		}
	}

	return report, nil
}
