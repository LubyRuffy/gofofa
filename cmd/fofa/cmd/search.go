package cmd

import (
	"errors"
	"fmt"
	"github.com/LubyRuffy/gofofa"
	"github.com/LubyRuffy/gofofa/pkg/outformats"
	"github.com/sirupsen/logrus"
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
	inFile      string // in file
	deductMode  string // deduct Mode
	fixUrl      bool   // each host fix as url, like 1.1.1.1,80 will change to http://1.1.1.1
	urlPrefix   string // each host fix as url, like 1.1.1.1,80 will change to http://1.1.1.1
	full        bool   // search result for over a year
	batchSize   int    // amount of data contained in each batch, only for dump
	json        bool   // out format as json for short
	uniqByIP    bool   // group by ip
)

// search subcommand
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
		&cli.BoolFlag{
			Name:        "fixUrl",
			Value:       false,
			Usage:       "each host fix as url, like 1.1.1.1,80 will change to http://1.1.1.1",
			Destination: &fixUrl,
		},
		&cli.StringFlag{
			Name:        "urlPrefix",
			Value:       "",
			Usage:       "prefix of url, default is http://, can be redis:// and so on ",
			Destination: &urlPrefix,
		},
		&cli.BoolFlag{
			Name:        "full",
			Value:       false,
			Usage:       "search result for over a year",
			Destination: &full,
		},
		&cli.BoolFlag{
			Name:        "uniqByIP",
			Value:       false,
			Usage:       "uniq by ip",
			Destination: &uniqByIP,
		},
	},
	Action: SearchAction,
}

func fieldIndex(fields []string, fieldName string) int {
	for i, f := range fields {
		if f == fieldName {
			return i
		}
	}
	return -1
}

func hashField(fields []string, fieldName string) bool {
	for _, f := range fields {
		if f == fieldName {
			return true
		}
	}
	return false
}

func hasBodyField(fields []string) bool {
	return hashField(fields, "body")
}

// SearchAction search action
func SearchAction(ctx *cli.Context) error {
	// valid same config
	for _, arg := range ctx.Args().Slice() {
		if arg[0] == '-' {
			return errors.New(fmt.Sprintln("there is args after fofa query:", arg))
		}
	}

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
	if hasBodyField(fields) && format == "csv" {
		logrus.Warnln("fields contains body, so change format to json")
		writer = outformats.NewJSONWriter(outTo, fields)
	} else {
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
	}

	// do search
	res, err := fofaCli.HostSearch(query, size, fields, gofofa.SearchOptions{
		FixUrl:    fixUrl,
		UrlPrefix: urlPrefix,
		Full:      full,
		UniqByIP:  uniqByIP,
	})
	if err != nil {
		return err
	}

	// output
	if err = writer.WriteAll(res); err != nil {
		return err
	}
	return nil
}
