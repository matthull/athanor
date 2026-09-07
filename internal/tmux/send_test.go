package tmux

import (
	"errors"
	"testing"
)

func TestIsTransientSendKeysError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"exact not in a mode", errors.New("not in a mode"), true},
		{"tmux error containing not in a mode", errors.New("tmux send-keys: not in a mode"), true},
		{"session not found", errors.New("session not found"), false},
		{"no server running", errors.New("no server running"), false},
		{"generic error", errors.New("something went wrong"), false},
		{"wrapped not in a mode", errors.New("error: not in a mode: target %5"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isTransientSendKeysError(tt.err)
			if got != tt.want {
				t.Errorf("isTransientSendKeysError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestChunkCount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		msgLen int
		want   int
	}{
		{"empty", 0, 0},
		{"single byte", 1, 1},
		{"exactly one chunk", sendKeysChunkSize, 1},
		{"one byte over", sendKeysChunkSize + 1, 2},
		{"exactly two chunks", sendKeysChunkSize * 2, 2},
		{"two and a bit", sendKeysChunkSize*2 + 1, 3},
		{"large message 2000", 2000, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := chunkCount(tt.msgLen)
			if got != tt.want {
				t.Errorf("chunkCount(%d) = %d, want %d", tt.msgLen, got, tt.want)
			}
		})
	}
}

// chunkCount returns the number of chunks sendMessageToTarget would produce
// for a message of the given length. Exported only in test for validation.
func chunkCount(msgLen int) int {
	if msgLen == 0 {
		return 0
	}
	if msgLen <= sendKeysChunkSize {
		return 1
	}
	n := msgLen / sendKeysChunkSize
	if msgLen%sendKeysChunkSize != 0 {
		n++
	}
	return n
}
