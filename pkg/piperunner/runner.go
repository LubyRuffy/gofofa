package piperunner

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/sirupsen/logrus"
	"html/template"
	"strings"
)

// PipeRunner pipe运行器
type PipeRunner struct {
	content  string
	Tasks    []*corefuncs.PipeTask // 执行的所有workflow
	LastTask *corefuncs.PipeTask   // 最后执行的workflow
	LastFile string                // 最后生成的文件名
	FofaCli  *gofofa.Client        // fofa客户端

	*gorunner.Runner
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
func New() *PipeRunner {
	r := &PipeRunner{}
	var err error

	// 注册底层函数
	var gf gorunner.GoFunction
	err = gf.Register("GetRunner", func() *PipeRunner {
		return r
	})
	if err != nil {
		panic(err)
	}
	err = gf.Register("Fork", func(pipe string) error {
		forkRunner := New()
		forkRunner.LastFile = r.LastFile // 从这里开始分叉
		r.LastTask.Children = append(r.LastTask.Children, forkRunner)
		_, err := forkRunner.Run(pipeast.NewParser().Parse(pipe))
		return err
	})
	if err != nil {
		panic(err)
	}

	corefuncs.Range(func(key, value any) bool {
		gf.Register(key.(string), value)
		return true
	})
	r.Runner = gorunner.New(gorunner.WithFunctions(&gf))
	return r
}
