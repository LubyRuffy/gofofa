package piperunner

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"html/template"
	"os"
	"reflect"
	"sync"
	"time"
)

var (
	functions sync.Map
)

// PipeTask 每一个pipe执行的任务统计信息
type PipeTask struct {
	Name    string        // pipe name
	Content string        // raw content
	Outfile string        // tmp file
	Cost    time.Duration // time costs
}

// Close remove tmp outfile
func (p *PipeTask) Close() {
	os.Remove(p.Outfile)
}

// PipeRunner pipe运行器
type PipeRunner struct {
	content  string
	tasks    []PipeTask
	LastFile string // 最后生成的文件名
}

// RegisterWorkflow 注册workflow
// 第一个参数是workflow名称；
// 第二个参数是workflow转换为函数调用字符串的函数
// 第三个参数是底层函数的名称
// 第四个参数是一个回调函数，参数是传递的参数，返回值是生成的文件名
// 第三四个参数可以留空值，表明只注册到语法解析器中去
func RegisterWorkflow(workflow string, transFunc pipeast.FunctionTranslateHook,
	funcName string, funcBody func(*PipeRunner, map[string]interface{}) string) {

	// 解析器的函数注册
	if len(workflow) > 0 {
		pipeast.RegisterFunction(workflow, transFunc)
	}

	// 注册底层函数
	if len(funcName) > 0 {
		// 执行并且自动生成任务队列
		functions.Store(funcName, func(p *PipeRunner, params map[string]interface{}) {
			logrus.Debug(funcName+" params:", params)

			s := time.Now()
			fn := funcBody(p, params)

			p.addPipe(PipeTask{
				Name:    funcName,
				Content: fmt.Sprintf("%v", params),
				Outfile: fn,
				Cost:    time.Since(s),
			})
		})
	}
}

// Close remove tmp outfile
func (p *PipeRunner) Close() {
	for _, task := range p.tasks {
		task.Close()
	}
}

func (p *PipeRunner) addPipe(pt PipeTask) {
	p.tasks = append(p.tasks, pt)
	p.LastFile = pt.Outfile

	logrus.Debug(pt.Name+" write to file: ", pt.Outfile)
}

// Run run pipelines
func (p *PipeRunner) Run() error {
	var err error

	p.tasks = nil

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

// DumpTasks tasks dump to html
func (p *PipeRunner) DumpTasks() string {
	t, _ := template.New("tasks").Parse(`   
<html>
<head>
    <title>gofofa tasks</title>
</head>
<body>
	<h1>gofofa tasks</h1>
	{{range .}}
		<ul>
			<li>{{ .Name }}</li>
			<li>{{ .Content }}</li>
			<li><a href="{{ .Outfile }}">{{ .Outfile }}</a></li>
			<li>{{ .Cost }}</li>
		</ul>
	{{end}}
</body>
</html>`)
	var out bytes.Buffer
	err := t.Execute(&out, p.tasks)
	if err != nil {
		panic(err)
	}

	return out.String()
}

// New create pipe runner
func New(content string) *PipeRunner {
	return &PipeRunner{
		content: content,
	}
}
