package main

import (
	"os"

	"github.com/Storm-Harbor/storm-harbor-control-plane/internal/cli"
)

var version = "dev"
var revision = "unknown"

func main() {
	os.Exit(cli.Execute(os.Args[1:], os.Stdout, os.Stderr, cli.Options{Version: version, Revision: revision}))
}
