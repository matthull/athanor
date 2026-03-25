// ath is the athanor CLI — agent orchestration, communication, and lifecycle management.
package main

import (
	"os"

	"github.com/matthull/athanor/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
