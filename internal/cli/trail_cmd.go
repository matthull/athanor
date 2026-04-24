package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/matthull/athanor/internal/athanor"
)

func runTrail(args []string) int {
	positional, flagArgs := splitArgs(args)

	var moFilter string
	var jsonOutput bool

	fs := flag.NewFlagSet("trail", flag.ContinueOnError)
	fs.StringVar(&moFilter, "mo", "", "filter by magnum opus name")
	fs.BoolVar(&jsonOutput, "json", false, "output as JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	var athName string
	if len(positional) > 0 {
		athName = positional[0]
	}

	report, err := athanor.AnalyzeTrail(home, athName, moFilter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if jsonOutput {
		return printTrailJSON(report)
	}
	return printTrailHuman(report, athName, moFilter)
}

func printTrailJSON(report *athanor.TrailReport) int {
	type jsonReport struct {
		TotalOpera      int            `json:"total_opera"`
		DischargedOpera int            `json:"discharged_opera"`
		ChargedOpera    int            `json:"charged_opera"`
		JobCounts       map[string]int `json:"job_counts"`
		Collaboration   struct {
			WhisperMentions     int     `json:"whisper_mentions"`
			InscribeMentions    int     `json:"inscribe_mentions"`
			CollaborateMentions int     `json:"collaborate_mentions"`
			MusterMentions      int     `json:"muster_mentions"`
			PeerReferences      int     `json:"peer_references"`
			OperaWithDyad       int     `json:"opera_with_dyad_evidence"`
			DyadCoverage        float64 `json:"dyad_coverage_pct"`
		} `json:"collaboration"`
	}

	jr := jsonReport{
		TotalOpera:      report.TotalOpera,
		DischargedOpera: report.DischargedOpera,
		ChargedOpera:    report.ChargedOpera,
		JobCounts:       report.JobCounts,
	}
	jr.Collaboration.WhisperMentions = report.TotalWhisperMentions
	jr.Collaboration.InscribeMentions = report.TotalInscribeMentions
	jr.Collaboration.CollaborateMentions = report.TotalCollaborateMentions
	jr.Collaboration.MusterMentions = report.TotalMusterMentions
	jr.Collaboration.PeerReferences = report.TotalPeerReferences
	jr.Collaboration.OperaWithDyad = report.OperaWithDyadEvidence
	if report.DischargedOpera > 0 {
		jr.Collaboration.DyadCoverage = float64(report.OperaWithDyadEvidence) / float64(report.DischargedOpera) * 100
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jr); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
		return 1
	}
	return 0
}

func printTrailHuman(report *athanor.TrailReport, athName, moFilter string) int {
	if report.TotalOpera == 0 {
		fmt.Println("No opera found.")
		return 0
	}

	// Header
	scope := "all athanors"
	if athName != "" {
		scope = athName
	}
	if moFilter != "" {
		scope += "/" + moFilter
	}
	fmt.Printf("ATH TRAIL — %s\n", scope)
	fmt.Println("═══════════════════════════════════════════════════")

	// Opera counts
	fmt.Printf("\nOpera: %d total (%d discharged, %d charged)\n",
		report.TotalOpera, report.DischargedOpera, report.ChargedOpera)

	// Job diversity
	fmt.Println("\nJob Diversity:")
	type jobCount struct {
		name  string
		count int
	}
	var jobs []jobCount
	for name, count := range report.JobCounts {
		jobs = append(jobs, jobCount{name, count})
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].count > jobs[j].count })
	for _, jc := range jobs {
		pct := float64(jc.count) / float64(report.TotalOpera) * 100
		fmt.Printf("  %-20s %3d (%4.1f%%)\n", jc.name, jc.count, pct)
	}

	// Collaboration signals
	fmt.Println("\nCollaboration Signals (across all opera):")
	fmt.Printf("  Inscribe mentions:    %d\n", report.TotalInscribeMentions)
	fmt.Printf("  Collaborate mentions: %d\n", report.TotalCollaborateMentions)
	fmt.Printf("  Muster mentions:      %d\n", report.TotalMusterMentions)
	fmt.Printf("  Whisper mentions:     %d\n", report.TotalWhisperMentions)
	fmt.Printf("  Peer references:      %d\n", report.TotalPeerReferences)

	// Dyad coverage
	fmt.Println("\nDyad Coverage:")
	fmt.Printf("  Opera with dyad evidence: %d / %d discharged",
		report.OperaWithDyadEvidence, report.DischargedOpera)
	if report.DischargedOpera > 0 {
		pct := float64(report.OperaWithDyadEvidence) / float64(report.DischargedOpera) * 100
		fmt.Printf(" (%.1f%%)", pct)
	}
	fmt.Println()

	return 0
}
