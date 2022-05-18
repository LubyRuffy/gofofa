package cmd

import (
	"github.com/lubyruffy/gofofa"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
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
	countCmd,
	statsCmd,
}

// IsValidCommand valid command name
func IsValidCommand(cmd string) bool {
	if len(cmd) == 0 {
		return false
	}
	if cmd[0] == '-' {
		switch cmd {
		case "--help", "-help", "-h", "--version", "-version", "-v":
			// 自带的配置
			return true
		default:
			for _, option := range GlobalOptions {
				for _, name := range option.Names() {
					if cmd == "--"+name || cmd == "-"+name {
						return true
					}
				}
			}
		}
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
	&cli.BoolFlag{
		Name:  "verbose",
		Usage: "print more information",
	},
}

// BeforAction generate fofa client
func BeforAction(context *cli.Context) error {
	var err error

	// not any command, and no query for default command(search)
	if len(os.Args) == 1 {
		return nil
	}

	if context.Bool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	fofaCli, err = gofofa.NewClient(fofaURL)
	if err != nil {
		return err
	}

	if len(deductMode) > 0 {
		fofaCli.DeductMode = gofofa.ParseDeductMode(deductMode)
	}

	return nil
}
