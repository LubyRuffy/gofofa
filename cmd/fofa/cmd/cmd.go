package cmd

import (
	"errors"
	"fmt"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var (
	fofaCli *gofofa.Client
)

var (
	fofaURL     string // fofa url
	deductMode  string // deduct Mode
	fieldString string // fieldString
	query       string // fieldString
	size        int    // fetch size
	outFormat   string // out format
)

// GlobalCommands global commands
var GlobalCommands = []*cli.Command{
	searchCmd,
	accountCmd,
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

// account 子命令
var accountCmd = &cli.Command{
	Name:  "account",
	Usage: "fofa account information",
	Action: func(ctx *cli.Context) error {
		fmt.Println(fofaCli.Account)
		return nil
	},
}

// search 子命令
var searchCmd = &cli.Command{
	Name:  "search",
	Usage: "fofa host search",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "query",
			Aliases:     []string{"q"},
			Value:       "",
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
		&cli.StringFlag{
			Name:        "outFormat",
			Aliases:     []string{"o"},
			Value:       "csv",
			Usage:       "can be csv/json/xml",
			Destination: &outFormat,
		},
		&cli.IntFlag{
			Name:        "size",
			Value:       100,
			Usage:       "if DeductModeFree set, select free limit size automatically",
			Destination: &size,
		},
		&cli.StringFlag{
			Name:        "deductMode",
			Value:       "DeductModeFree",
			Usage:       "DeductModeFree or DeductModeFCoin",
			Destination: &deductMode,
		},
	},
	Action: func(ctx *cli.Context) error {
		// valid same config
		if len(query) == 0 {
			return errors.New("fofa query cannot be empty")
		}
		fields := strings.Split(fieldString, ",")
		if len(fields) == 0 {
			return errors.New("fofa fields cannot be empty")
		}

		// gen writer
		outF := os.Stdout
		var writer outformats.OutWriter
		switch outFormat {
		case "csv":
			writer = outformats.NewCSVWriter(outF)
		case "json":
			writer = outformats.NewJSONWriter(outF, fields)
		case "xml":
			writer = outformats.NewXMLWriter(outF, fields)
		default:
			return fmt.Errorf("unknown outFormat: %s", outFormat)
		}

		// do search
		res, err := fofaCli.HostSearch(query, 100, fields)
		if err != nil {
			return err
		}

		// output
		if err = writer.WriteAll(res); err != nil {
			return err
		}
		return nil
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
