package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/matthull/athanor/internal/athanor"
)

func runInit(args []string) int {
	var project string

	positional, flagArgs := splitArgs(args)

	if len(positional) < 1 {
		fmt.Fprintln(os.Stderr, "error: name required")
		fmt.Fprintln(os.Stderr, "usage: ath init <name> [--project <path>]")
		return 2
	}
	name := positional[0]

	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.StringVar(&project, "project", "", "working directory for agents")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(flagArgs); err != nil {
		return 2
	}

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Ensure home directory structure exists
	if err := athanor.EnsureHome(home); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Check shared components exist in source repo
	sharedPath, err := athanor.SharedPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	for _, f := range athanor.SharedFiles {
		path := fmt.Sprintf("%s/%s", sharedPath, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: shared component %q not found at %s\n", f, sharedPath)
			fmt.Fprintf(os.Stderr, "Ensure the athanor repo is checked out (or set $ATHANOR_REPO).\n")
			return 1
		}
	}

	if err := athanor.InitInstance(home, name, project); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Ensure .env.local is gitignored in the athanor home repo
	if err := athanor.EnsureEnvGitignore(home); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not update .gitignore: %v\n", err)
		// Non-fatal — the instance was created successfully
	}

	instDir := athanor.InstanceDir(home, name)
	fmt.Printf("Athanor %q initialized at %s\n", name, instDir)
	fmt.Printf("Create a magnum opus in magna-opera/, then `ath kindle %s <mo-name>`.\n", name)

	return 0
}
