package cli

import (
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/tmux"
)

func runCheck(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "error: crucible name required")
		fmt.Fprintln(os.Stderr, "usage: ath check <crucible-name>")
		return 2
	}
	crucible := args[0]

	r := tmux.NewRunner()
	state, desc, err := r.CheckCrucible(crucible)
	if err != nil && state != tmux.StateDead {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	fmt.Printf("%-12s%s\n", state, desc)

	switch state {
	case tmux.StateActive, tmux.StateIdle:
		return 0
	case tmux.StatePermission, tmux.StateExhausted:
		return 1
	default:
		return 2
	}
}
