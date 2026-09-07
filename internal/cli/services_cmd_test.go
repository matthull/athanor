package cli

import (
	"strings"
	"testing"
	"time"
)

func TestParseListTimersJSON(t *testing.T) {
	input := []byte(`[
		{"next":1776711295998436,"left":1776711295998436,"last":1776710995975107,"passed":227966748926,"unit":"athanor-liveness.timer","activates":"athanor-liveness.service"},
		{"next":1776711198241829,"left":1776711198241829,"last":1776710898237025,"passed":227869010843,"unit":"voice-notes-process.timer","activates":"voice-notes-process.service"}
	]`)

	result := parseListTimersJSON(input)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}

	live, ok := result["athanor-liveness.timer"]
	if !ok {
		t.Fatal("expected athanor-liveness.timer in map")
	}
	if live.Last != 1776710995975107 {
		t.Errorf("unexpected last: %d", live.Last)
	}
	if live.Next != 1776711295998436 {
		t.Errorf("unexpected next: %d", live.Next)
	}
	if live.Unit != "athanor-liveness.timer" {
		t.Errorf("unexpected unit: %q", live.Unit)
	}

	if _, ok := result["voice-notes-process.timer"]; !ok {
		t.Error("expected voice-notes-process.timer in map")
	}
}

func TestParseListTimersJSONEmpty(t *testing.T) {
	result := parseListTimersJSON([]byte(`[]`))
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestParseListTimersJSONInvalid(t *testing.T) {
	result := parseListTimersJSON([]byte(`not valid json`))
	if len(result) != 0 {
		t.Errorf("expected empty map on invalid JSON, got %d entries", len(result))
	}
}

func TestFormatTimeSince(t *testing.T) {
	// Zero time returns empty string
	if got := formatTimeSince(time.Time{}); got != "" {
		t.Errorf("zero time: expected empty, got %q", got)
	}

	// Past time returns "HH:MM (Xm ago)"
	past := time.Now().Add(-5 * time.Minute)
	got := formatTimeSince(past)
	if got == "" {
		t.Error("past time: expected non-empty result")
	}
	if !strings.Contains(got, "ago") {
		t.Errorf("past time: expected 'ago' in output, got %q", got)
	}
	if !strings.Contains(got, past.Format("15:04")) {
		t.Errorf("past time: expected HH:MM %q in output, got %q", past.Format("15:04"), got)
	}

	// Future time returns "HH:MM (in Xm)"
	future := time.Now().Add(10 * time.Minute)
	got = formatTimeSince(future)
	if got == "" {
		t.Error("future time: expected non-empty result")
	}
	if !strings.Contains(got, "in ") {
		t.Errorf("future time: expected 'in ' in output, got %q", got)
	}
	if !strings.Contains(got, future.Format("15:04")) {
		t.Errorf("future time: expected HH:MM %q in output, got %q", future.Format("15:04"), got)
	}
}
