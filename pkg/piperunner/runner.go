package piperunner

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	_ "github.com/lubyruffy/gofofa/pkg/piperunner/input"
	_ "github.com/lubyruffy/gofofa/pkg/piperunner/output"
	"github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"html/template"
	"reflect"
	"strings"
)

// PipeRunner pipe运行器
type PipeRunner struct {
	content  string
	Tasks    []*corefuncs.PipeTask // 执行的所有workflow
	LastTask *corefuncs.PipeTask   // 最后执行的workflow
	LastFile string                // 最后生成的文件名
	FofaCli  *gofofa.Client        // fofa客户端
}

// Close remove tmp outfile
func (p *PipeRunner) Close() {
	for _, task := range p.Tasks {
		task.Close()
	}
}

// GetWorkflows all workflows
func (p *PipeRunner) GetWorkflows() []*corefuncs.PipeTask {
	return p.Tasks
}

// AddWorkflow 添加一次任务的日志
func (p *PipeRunner) AddWorkflow(pt *corefuncs.PipeTask) {
	p.Tasks = append(p.Tasks, pt)
	p.LastTask = pt

	// 可以不写文件
	if len(pt.Outfile) > 0 {
		p.LastFile = pt.Outfile

		logrus.Debug(pt.Name+" write to file: ", pt.Outfile)
	}
}

// Run run pipelines
func (p *PipeRunner) Run() error {
	var err error

	p.Tasks = nil

	i := interp.New(interp.Options{})
	_ = i.Use(stdlib.Symbols)

	exports := interp.Exports{
		"this/this": {
			"GetRunner": reflect.ValueOf(func() *PipeRunner {
				return p
			}),
			"Fork": reflect.ValueOf(func(pipe string) error {
				forkRunner := New(pipeast.NewParser().Parse(pipe))
				forkRunner.LastFile = p.LastFile // 从这里开始分叉
				p.LastTask.Children = append(p.LastTask.Children, forkRunner)
				return forkRunner.Run()
			}),
		},
	}
	corefuncs.Range(func(key, value any) bool {
		exports["this/this"][key.(string)] = reflect.ValueOf(value)
		return true
	})

	err = i.Use(exports)
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
	t, _ := template.New("tasks").Funcs(template.FuncMap{
		"RawHtml": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
		"safeURL": func(u string) template.URL {
			u = strings.ReplaceAll(u, "\\", "/")
			return template.URL(u)
		},
		"GetTasks": func(p corefuncs.Runner) []*corefuncs.PipeTask {
			return p.GetWorkflows()
		},
	}).Parse(`   
<html>
<head>
    <title>gofofa tasks</title>
</head>
<body>
	<h1>gofofa tasks</h1>
	{{range .}}
		{{ template "task.tmpl" . }}
	{{end}}
</body>
</html>

{{ define "task.tmpl" }}
<ul>
	<li>{{ .Name }}</li>
	<li>{{ .Content }}</li>

	{{ if gt (len .Outfile) 0 }}
	<li><a href="{{ .Outfile | safeURL }}">{{ .Outfile }}</a></li>
	{{ end }}

	{{ range .GeneratedFiles }}	
	<li> generate files:
		<ul>
			<li><a href="{{ . | safeURL }}">{{ . }}</a></li>
		</ul>
	</li>
	{{ end }}
	<li>{{ .Cost }}</li>

	{{ range .Children }}
	<li> fork children:
		{{ range . | GetTasks }}
			{{ template "task.tmpl" . }}
		{{ end }}
	</li>
	{{ end }}
</ul>
{{ end }}
`)
	var out bytes.Buffer
	err := t.Execute(&out, p.Tasks)
	if err != nil {
		panic(err)
	}

	return out.String()
}

// GetFofaCli fofa client
func (p *PipeRunner) GetFofaCli() *gofofa.Client {
	return p.FofaCli
}

// GetLastFile last genrated file
func (p *PipeRunner) GetLastFile() string {
	return p.LastFile
}

// New create pipe runner
func New(content string) *PipeRunner {
	return &PipeRunner{
		content: content,
	}
}
