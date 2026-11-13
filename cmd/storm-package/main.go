// storm-package is a maintainer build tool, not a recovery command.
package main

import (
	"os"

	"github.com/Storm-Harbor/storm-harbor-control-plane/internal/release"
)

func main() { os.Exit(release.Run(os.Args[1:], os.Stdout, os.Stderr)) }
