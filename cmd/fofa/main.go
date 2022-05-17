package main

import (
	"fmt"
	"github.com/lubyruffy/gofofa/cmd/fofa/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown" // goreleaser fill
)

func main() {
	app := &cli.App{
		Name:                   "fofa",
		Usage:                  fmt.Sprintf("fofa client on Go %s, commit %s, built at %s", version, commit, date),
		Version:                version,
		UseShortOptionHandling: true,
		Flags:                  cmd.GlobalOptions,
		Before:                 cmd.BeforAction,
		Commands:               cmd.GlobalCommands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
