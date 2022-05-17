package cmd

import (
	"github.com/lubyruffy/gofofa"
	"github.com/urfave/cli/v2"
)

var (
	fofaCli *gofofa.Client
)

var (
	fofaURL string // fofa url
)

// GlobalCommands global commands
var GlobalCommands = []*cli.Command{
	searchCmd,
	accountCmd,
}

// IsValidCommand valid command name
func IsValidCommand(cmd string) bool {
	if len(cmd) == 0 {
		return false
	}
	if cmd[0] == '-' {
		return false
	}

	for _, command := range GlobalCommands {
		if command.Name == cmd {
			return true
		}
	}
	return false
}

// GlobalOptions global options
var GlobalOptions = []cli.Flag{
	&cli.StringFlag{
		Name:        "fofaURL",
		Aliases:     []string{"u"},
		Value:       gofofa.FofaURLFromEnv(),
		Usage:       "format: <url>/?email=<email>&key=<key>&version=<v2>",
		Destination: &fofaURL,
	},
}

// BeforAction generate fofa client
func BeforAction(context *cli.Context) error {
	var err error
	fofaCli, err = gofofa.NewClient(fofaURL)
	if err != nil {
		return err
	}

	if len(deductMode) > 0 {
		fofaCli.DeductMode = gofofa.ParseDeductMode(deductMode)
	}

	return nil
}
