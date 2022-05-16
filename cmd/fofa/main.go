package main

import (
	"github.com/lubyruffy/gofofa"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	fofaCli *gofofa.Client
)

func main() {
	app := &cli.App{
		Name:                   "fofa",
		Usage:                  "fofa client on Go",
		Version:                "v0.0.1",
		UseShortOptionHandling: true,
		Flags:                  globalOptions,
		Commands:               globalCommands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
