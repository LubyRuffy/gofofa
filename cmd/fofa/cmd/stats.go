package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"os"
	"strings"
	"sync"
)

// stats subcommand
var statsCmd = &cli.Command{
	Name:                   "stats",
	Usage:                  "fofa stats",
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "fields",
			Aliases:     []string{"f"},
			Value:       "title,country",
			Usage:       "visit fofa website for more info",
			Destination: &fieldString,
		},
		&cli.IntFlag{
			Name:        "size",
			Aliases:     []string{"s"},
			Value:       5,
			Usage:       "aggs size",
			Destination: &size,
		},
	},
	Action: statsAction,
}

// statsAction stats action
func statsAction(ctx *cli.Context) error {
	// valid same config
	query := ctx.Args().First()
	if len(query) == 0 {
		return errors.New("fofa query cannot be empty")
	}
	fields := strings.Split(fieldString, ",")
	if len(fields) == 0 {
		return errors.New("fofa fields cannot be empty")
	}

	// do search
	res, err := fofaCli.Stats(query, size, fields)
	if err != nil {
		return err
	}

	for _, obj := range res {
		color.New(color.FgBlue).Fprintln(os.Stdout, "=== ", obj.Name)
		for _, item := range obj.Items {
			color.New(color.FgHiGreen).Fprintf(os.Stdout, "%s\tassets(%d)", item.Name, item.Count)
			if item.Uniq != nil {
				for k, v := range item.Uniq {
					color.New(color.FgHiGreen).Fprintf(os.Stdout, "\t%s(%d)", k, v)
				}
				color.New(color.FgHiYellow).Fprintln(os.Stdout)
			}

			if item.Detail != nil {
				if len(item.Detail.IPDetail.Domains) > 0 {
					// 去重打印
					var rootDomains []string
					var uniqMap sync.Map
					for _, domain := range item.Detail.IPDetail.Domains {
						rootDomain, err := publicsuffix.Domain(domain)
						if err != nil {
							continue
						}
						if _, ok := uniqMap.LoadOrStore(rootDomain, true); ok {
							continue
						}
						rootDomains = append(rootDomains, rootDomain)
					}
					color.New(color.FgHiGreen).Fprintln(os.Stdout,
						fmt.Sprintf("\tdomains(%d): %s", len(rootDomains), strings.Join(rootDomains, ",")))
				}

				if item.Detail.CertDetail.Issuer != nil {
					color.New(color.FgHiGreen).Fprintln(os.Stdout, "\tvalid: ", item.Detail.CertDetail.IsValid)
					color.New(color.FgHiGreen).Fprintln(os.Stdout, "\texpired: ", item.Detail.CertDetail.IsExpired)
					color.New(color.FgHiGreen).Fprintln(os.Stdout, "\tnot_before: ", item.Detail.CertDetail.NotBefore)
					color.New(color.FgHiGreen).Fprintln(os.Stdout, "\torganization: ", strings.Join(item.Detail.CertDetail.Subject.Organizations, ","))
					color.New(color.FgHiGreen).Fprintln(os.Stdout, "\troot_domains: ", strings.Join(item.Detail.CertDetail.RootDomains, ","))
				}

				if item.Detail.ASNDetail.Org != nil {
					color.New(color.FgHiGreen).Fprintln(os.Stdout,
						fmt.Sprintf("\torg: %s", *item.Detail.ASNDetail.Org))
				}
			}
		}
	}

	return nil
}
