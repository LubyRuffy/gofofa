package main

import (
	"encoding/csv"
	"github.com/lubyruffy/gofofa"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var (
	fofaURL     string // fofa url
	deductMode  string // deduct Mode
	fieldString string // fieldString
	query       string // fieldString
)

// Commands 将子命令统一暴露给 main 包
var globalCommands = []*cli.Command{
	searchCmd,
}

// globalOptions
var globalOptions = []cli.Flag{
	&cli.StringFlag{
		Name:        "fofaURL",
		Aliases:     []string{"u"},
		Value:       gofofa.FofaURLFromEnv(),
		Usage:       "format: <url>/?email=<email>&key=<key>&version=<v2>",
		Destination: &fofaURL,
	},
	&cli.StringFlag{
		Name:        "deductMode",
		Value:       "DeductModeFree",
		Usage:       "DeductModeFree or DeductModeFCoin",
		Destination: &deductMode,
	},
}

var searchCmd = &cli.Command{
	Name:  "search",
	Usage: "fofa host search",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "query",
			Aliases:     []string{"q"},
			Value:       "port=80",
			Usage:       "fofa query",
			Destination: &query,
		},
		&cli.StringFlag{
			Name:        "fields",
			Aliases:     []string{"f"},
			Value:       "ip,port",
			Usage:       "visit fofa website for more info",
			Destination: &fieldString,
		},
	},
	Action: func(ctx *cli.Context) error {
		cli, err := gofofa.NewClient(fofaURL)
		if err != nil {
			return err
		}

		if len(deductMode) > 0 {
			cli.DeductMode = gofofa.ParseDeductMode(deductMode)
		}

		res, err := cli.HostSearch(query, 100, strings.Split(fieldString, ","))
		if err != nil {
			return err
		}

		writer := csv.NewWriter(os.Stdout)
		if err = writer.WriteAll(res); err != nil {
			panic(err)
		}
		return nil
	},
}
