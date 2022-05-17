package cmd

import (
	"errors"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"strings"
)

var (
	fieldString string // fieldString
	size        int    // fetch size
	format      string // out format
	outFile     string // out file
	deductMode  string // deduct Mode
)

// search 子命令
var searchCmd = &cli.Command{
	Name:                   "search",
	Usage:                  "fofa host search",
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "fields",
			Aliases:     []string{"f"},
			Value:       "ip,port",
			Usage:       "visit fofa website for more info",
			Destination: &fieldString,
		},
		&cli.StringFlag{
			Name:        "format",
			Value:       "csv",
			Usage:       "can be csv/json/xml",
			Destination: &format,
		},
		&cli.StringFlag{
			Name:        "outFile",
			Aliases:     []string{"o"},
			Usage:       "if not set, wirte to stdout",
			Destination: &outFile,
		},
		&cli.IntFlag{
			Name:        "size",
			Aliases:     []string{"s"},
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
	Action: SearchAction,
}

// SearchAction search action
func SearchAction(ctx *cli.Context) error {
	// valid same config
	query := ctx.Args().First()
	if len(query) == 0 {
		return errors.New("fofa query cannot be empty")
	}
	fields := strings.Split(fieldString, ",")
	if len(fields) == 0 {
		return errors.New("fofa fields cannot be empty")
	}

	// gen output
	var outTo io.Writer
	if len(outFile) > 0 {
		var f *os.File
		var err error
		if f, err = os.Create(outFile); err != nil {
			return fmt.Errorf("create outFile %s failed: %w", outFile, err)
		}
		outTo = f
		defer f.Close()
	} else {
		outTo = os.Stdout
	}

	// gen writer
	var writer outformats.OutWriter
	switch format {
	case "csv":
		writer = outformats.NewCSVWriter(outTo)
	case "json":
		writer = outformats.NewJSONWriter(outTo, fields)
	case "xml":
		writer = outformats.NewXMLWriter(outTo, fields)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	// do search
	res, err := fofaCli.HostSearch(query, size, fields)
	if err != nil {
		return err
	}

	// output
	if err = writer.WriteAll(res); err != nil {
		return err
	}
	return nil
}
