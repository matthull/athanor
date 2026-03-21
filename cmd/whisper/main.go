// whisper reliably delivers messages to Claude Code sessions running in tmux.
package main

import (
	"os"

	"github.com/mattcollier/whisper/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
