package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// account subcommand
var accountCmd = &cli.Command{
	Name:  "account",
	Usage: "fofa account information",
	Action: func(ctx *cli.Context) error {
		fmt.Println(fofaCli.Account)
		return nil
	},
}
