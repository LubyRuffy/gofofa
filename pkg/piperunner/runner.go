package piperunner

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"os"
	"reflect"
	"sync"
)

var (
	functions sync.Map
)

type pipeTask struct {
	name    string // pipe name
	content string // raw content
	outfile string // tmp file
}

// Close remove tmp outfile
func (p *pipeTask) Close() {
	os.Remove(p.outfile)
}

// PipeRunner pipe运行器
type PipeRunner struct {
	content  string
	tasks    []pipeTask
	LastFile string // 最后生成的文件名
}

// RegisterWorkflow 注册workflow
// 第一个参数是workflow名称；
// 第二个参数是workflow转换为函数调用字符串的函数
// 第三个参数是底层函数的名称
// 第四个参数是一个回调函数，参数是传递的参数，返回值是生成的文件名
func RegisterWorkflow(workflow string, transFunc pipeparser.FunctionTranslateHook, funcName string, funcBody func(*PipeRunner, map[string]interface{}) string) {
	// 解析器的函数注册
	pipeparser.RegisterFunction(workflow, transFunc)

	// 注册底层函数
	if len(funcName) > 0 {
		// 执行并且自动生成任务队列
		functions.Store(funcName, func(p *PipeRunner, params map[string]interface{}) {
			logrus.Debug(funcName+" params:", params)

			pt := pipeTask{
				name:    funcName,
				content: fmt.Sprintf("%v", params),
				outfile: funcBody(p, params),
			}
			p.addPipe(pt)
		})
	}
}

// Close remove tmp outfile
func (p *PipeRunner) Close() {
	for _, task := range p.tasks {
		task.Close()
	}
}

func (p *PipeRunner) addPipe(pt pipeTask) {
	p.tasks = append(p.tasks, pt)
	p.LastFile = pt.outfile

	logrus.Debug(pt.name+"write to file: ", pt.outfile)
}

// Run run pipelines
func (p *PipeRunner) Run() error {
	var err error

	i := interp.New(interp.Options{})
	_ = i.Use(stdlib.Symbols)

	funcs := interp.Exports{
		"this/this": {
			"GetRunner": reflect.ValueOf(func() *PipeRunner {
				return p
			}),
		},
	}
	functions.Range(func(key, value any) bool {
		funcs["this/this"][key.(string)] = reflect.ValueOf(value)
		return true
	})

	err = i.Use(funcs)
	if err != nil {
		panic(err)
	}

	// i.ImportUsed()
	i.Eval(`import (
		. "this/this"
		)`)

	_, err = i.Eval(p.content)

	return err
}

// New create pipe runner
func New(content string) *PipeRunner {
	return &PipeRunner{
		content: content,
	}
}
