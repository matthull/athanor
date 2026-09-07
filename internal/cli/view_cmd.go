package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/matthull/athanor/internal/athanor"
)

func runView(args []string) int {
	positional, _ := splitArgs(args)

	if len(positional) < 2 {
		fmt.Fprintln(os.Stderr, "usage: ath view <athanor> <mo> [<opus>]")
		return 2
	}

	athName := positional[0]
	moName := positional[1]

	home, err := athanor.Home()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	instDir := athanor.InstanceDir(home, athName)

	var filePath string
	if len(positional) >= 3 {
		opusName := positional[2]
		operaDir := athanor.OperaPath(instDir, moName)
		filePath = filepath.Join(operaDir, opusName)
		if !strings.HasSuffix(filePath, ".md") {
			filePath += ".md"
		}
	} else {
		filePath = athanor.MagnumOpusPath(instDir, moName)
	}

	if _, err := os.Stat(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s not found\n", filePath)
		return 1
	}

	editor := resolveEditor()

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func resolveEditor() string {
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if _, err := exec.LookPath("nvim"); err == nil {
		return "nvim"
	}
	return "vi"
}
