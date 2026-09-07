package cli

import (
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
)

func runSync(args []string) int {
	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	var names []string
	if len(args) > 0 {
		names = []string{args[0]}
	} else {
		// Sync all instances
		var err error
		names, err = athanor.ListInstances(home)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		if len(names) == 0 {
			fmt.Println("No athanor instances found.")
			return 0
		}
	}

	exitCode := 0
	for _, name := range names {
		if err := athanor.SyncInstance(home, name); err != nil {
			fmt.Fprintf(os.Stderr, "error syncing %s: %v\n", name, err)
			exitCode = 1
			continue
		}
		fmt.Printf("Synced %s\n", name)
	}

	return exitCode
}
