package cmd

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
)

// count subcommand
var countCmd = &cli.Command{
	Name:  "count",
	Usage: "fofa query results count",
	Action: func(ctx *cli.Context) error {
		query := ctx.Args().First()
		if len(query) == 0 {
			return errors.New("fofa query cannot be empty")
		}
		size, err := fofaCli.HostSize(query)
		if err != nil {
			return err
		}
		fmt.Println(size)
		return nil
	},
}
