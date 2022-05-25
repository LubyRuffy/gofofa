package goworkflow

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/coderunner"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/translater"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/sirupsen/logrus"
	"html/template"
	"os"
	"reflect"
	"strings"
	"time"
)

// PipeTask 每一个pipe执行的任务统计信息
type PipeTask struct {
	Name           string        // pipe name
	Content        string        // raw content
	Outfile        string        // tmp json file 统一格式
	GeneratedFiles []string      // files to archive 非json格式的文件不往后进行传递
	Cost           time.Duration // time costs
	Children       []*PipeRunner // fork children
}

// Close remove tmp outfile
func (p *PipeTask) Close() {
	os.Remove(p.Outfile)
}

// PipeRunner pipe运行器
type PipeRunner struct {
	content  string
	Tasks    []*PipeTask    // 执行的所有workflow
	LastTask *PipeTask      // 最后执行的workflow
	LastFile string         // 最后生成的文件名
	FofaCli  *gofofa.Client // fofa客户端

	gocodeRunner *coderunner.Runner
}

// Run go code, not workflow
func (p *PipeRunner) Run(code string) (reflect.Value, error) {
	return p.gocodeRunner.Run(code)
}

// Close remove tmp outfile
func (p *PipeRunner) Close() {
	for _, task := range p.Tasks {
		task.Close()
	}
}

// GetWorkflows all workflows
func (p *PipeRunner) GetWorkflows() []*PipeTask {
	return p.Tasks
}

// AddWorkflow 添加一次任务的日志
func (p *PipeRunner) AddWorkflow(pt *PipeTask) {
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
		"safeURL": func(u string) template.URL {
			u = strings.ReplaceAll(u, "\\", "/")
			return template.URL(u)
		},
		"GetTasks": func(p *PipeRunner) []*PipeTask {
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
	var gf coderunner.GoFunction
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
		code, err := workflowast.NewParser().Parse(pipe)
		if err != nil {
			return err
		}
		_, err = forkRunner.Run(code)
		return err
	})
	if err != nil {
		panic(err)
	}

	funcs := [][]interface{}{
		{"RemoveField", removeField},
		{"FetchFofa", fetchFofa},
		{"GenerateChart", generateChart},
		{"ZqQuery", zqQuery},
		{"AddField", addField},
		{"LoadFile", loadFile},
		{"FlatArray", flatArray},
	}
	for i := range funcs {
		funcName := funcs[i][0].(string)
		funcBody := funcs[i][1].(func(*PipeRunner, map[string]interface{}) (string, []string))
		gf.Register(funcName, func(p *PipeRunner, params map[string]interface{}) {
			logrus.Debug(funcName+" params:", params)

			s := time.Now()
			fn, gfs := funcBody(p, params)

			p.AddWorkflow(&PipeTask{
				Name:           funcName,
				Content:        fmt.Sprintf("%v", params),
				Outfile:        fn,
				GeneratedFiles: gfs,
				Cost:           time.Since(s),
			})
		})
	}

	logrus.Debug("ast support workflows:", len(translater.Translators))

	r.gocodeRunner = coderunner.New(coderunner.WithFunctions(&gf))
	return r
}
