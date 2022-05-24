package cmd

import (
	"errors"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/funcs"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
)

var (
	pipelineFile    string
	pipelineTaskOut string // 导出任务列表文件
	listWorkflows   bool
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
		&cli.StringFlag{
			Name:        "taskOut",
			Aliases:     []string{"t"},
			Usage:       "output pipeline tasks",
			Destination: &pipelineTaskOut,
		},
		&cli.BoolFlag{
			Name:        "list",
			Aliases:     []string{"l"},
			Usage:       "list support workflows",
			Destination: &listWorkflows,
		},
	},
	Action: pipelineAction,
}

// pipelineAction pipeline action
// 基本逻辑是：命令行的query是一个pipeline模式，每一个pipeline的workflow都要转换成底层可以执行的代码
// 也就是说注册一个pipeline可以支持的命令，需要：一）注册底层函数；二）注册pipeline的函数到底层函数调用的代码转换器
func pipelineAction(ctx *cli.Context) error {

	if listWorkflows {
		fmt.Println(corefuncs.SupportWorkflows())
		return nil
	}

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
		pipelineContent = pipeast.NewParser().Parse(v)
	}

	pr := piperunner.New(pipelineContent)
	pr.FofaCli = fofaCli
	err := pr.Run()
	if err != nil {
		return err
	}

	err = funcs.EachLine(pr.LastFile, func(line string) error {
		fmt.Println(line)
		return nil
	})
	if err != nil {
		return err
	}

	if len(pipelineTaskOut) > 0 {
		err = ioutil.WriteFile(pipelineTaskOut, []byte(pr.DumpTasks()), 0666)
		if err != nil {
			return err
		}

		browser.OpenFile(pipelineTaskOut)
	}

	return nil
}
