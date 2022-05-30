package cmd

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

// host subcommand
var hostCmd = &cli.Command{
	Name:                   "host",
	Usage:                  "fofa host",
	UseShortOptionHandling: true,
	Action:                 hostAction,
}

// hostAction stats action
func hostAction(ctx *cli.Context) error {
	// valid same config
	host := ctx.Args().First()
	if len(host) == 0 {
		return errors.New("fofa host cannot be empty")
	}

	// do search
	res, err := fofaCli.HostStats(host)
	if err != nil {
		return err
	}

	fmt.Println("Host:\t\t", res.Host)
	fmt.Println("IP:\t\t", res.IP)
	fmt.Println("ASN:\t\t", res.ASN)
	fmt.Println("ORG:\t\t", res.ORG)
	fmt.Println("Country:\t", res.Country)
	fmt.Println("CountryCode:\t", res.CountryCode)
	fmt.Println("Ports:\t\t", res.Ports)
	fmt.Println("Protocols:\t", strings.Join(res.Protocols, ","))
	fmt.Println("Categories:\t", strings.Join(res.Categories, ","))
	fmt.Println("Products:\t", strings.Join(res.Products, ","))
	fmt.Println("UpdateTime:\t", res.UpdateTime)

	return nil
}
