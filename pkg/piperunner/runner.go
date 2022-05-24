package piperunner

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
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
	tasks    []corefuncs.PipeTask
	LastFile string         // 最后生成的文件名
	FofaCli  *gofofa.Client // fofa客户端
}

// Close remove tmp outfile
func (p *PipeRunner) Close() {
	for _, task := range p.tasks {
		task.Close()
	}
}

// AddWorkflow 添加一次任务的日志
func (p *PipeRunner) AddWorkflow(pt corefuncs.PipeTask) {
	p.tasks = append(p.tasks, pt)

	// 可以不写文件
	if len(pt.Outfile) > 0 {
		p.LastFile = pt.Outfile

		logrus.Debug(pt.Name+" write to file: ", pt.Outfile)
	}
}

// Run run pipelines
func (p *PipeRunner) Run() error {
	var err error

	p.tasks = nil

	i := interp.New(interp.Options{})
	_ = i.Use(stdlib.Symbols)

	exports := interp.Exports{
		"this/this": {
			"GetRunner": reflect.ValueOf(func() *PipeRunner {
				return p
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
	}).Parse(`   
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
