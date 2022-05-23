package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/funcs"
	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
	"github.com/mitchellh/mapstructure"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
	"text/template"
)

var (
	pipelineFile string
)

// pipeline subcommand
var pipelineCmd = &cli.Command{
	Name:                   "pipeline",
	Usage:                  "fofa data pipeline",
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Usage:       "load pipeline file",
			Destination: &pipelineFile,
		},
	},
	Action: pipelineAction,
}

// pipelineAction pipeline action
// 基本逻辑是：命令行的query是一个pipeline模式，每一个pipeline的workflow都要转换成底层可以执行的代码
// 也就是说注册一个pipeline可以支持的命令，需要：一）注册底层函数；二）注册pipeline的函数到底层函数调用的代码转换器
func pipelineAction(ctx *cli.Context) error {

	funcs.Load()
	piperunner.RegisterWorkflow("fofa", fofaHook, "FetchFofa", fetchFofa)

	// valid same config
	var pipelineContent string
	if len(pipelineFile) > 0 {
		v, err := os.ReadFile(pipelineFile)
		if err != nil {
			return err
		}
		pipelineContent = string(v)
	}
	if v := ctx.Args().First(); len(v) > 0 {
		if len(pipelineContent) > 0 {
			return errors.New("file and content only one is allowed")
		}
		pipelineContent = pipeparser.NewParser().Parse(v)
	}

	pr := piperunner.New(pipelineContent)
	err := pr.Run()
	if err != nil {
		return err
	}

	err = piperunner.EachLine(pr.LastFile, func(line string) error {
		fmt.Println(line)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type fetchFofaParams struct {
	Query  string
	Size   int
	Fields string
}

func fetchFofa(p *piperunner.PipeRunner, params map[string]interface{}) string {
	var err error
	var options fetchFofaParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	if len(options.Query) == 0 {
		panic(errors.New("fofa query cannot be empty"))
	}
	if len(options.Fields) == 0 {
		panic(errors.New("fofa fields cannot be empty"))
	}

	fields := strings.Split(options.Fields, ",")

	var res [][]string
	res, err = fofaCli.HostSearch(options.Query, options.Size, fields)
	if err != nil {
		panic(err)
	}

	return piperunner.WriteTempJSONFile(func(f *os.File) {
		w := outformats.NewJSONWriter(f, fields)
		if err = w.WriteAll(res); err != nil {
			panic(err)
		}
	})
}

func fofaHook(fi *pipeparser.FuncInfo) string {
	tmpl, err := template.New("fofa").Parse(`FetchFofa(GetRunner(), map[string]interface{} {
    "query": {{ .Query }},
    "size": {{ .Size }},
    "fields": {{ .Fields }},
})`)
	if err != nil {
		panic(err)
	}
	var size int64 = 10
	fields := "`host,title`"
	if len(fi.Params) > 1 {
		fields = fi.Params[1].String()
	}
	if len(fi.Params) > 2 {
		size = fi.Params[2].Int64()
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		Query  string
		Size   int64
		Fields string
	}{
		Query:  fi.Params[0].String(),
		Fields: fields,
		Size:   size,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}
