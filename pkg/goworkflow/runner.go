package goworkflow

import (
	"fmt"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/coderunner"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/translater"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
	"reflect"
	"strings"
	"time"
)

// PipeTask 每一个pipe执行的任务统计信息
type PipeTask struct {
	Name      string        // pipe name
	Content   string        // raw content
	Outfile   string        // tmp json file 统一格式
	Artifacts []*Artifact   // files to archive 非json格式的文件不往后进行传递
	Cost      time.Duration // time costs
	Children  []*PipeRunner // fork children
}

// Close remove tmp outfile
func (p *PipeTask) Close() {
	os.Remove(p.Outfile)
}

// Hooks 消息通知
type Hooks struct {
	OnWorkflowFinished func(pt *PipeTask)
}

// PipeRunner pipe运行器
type PipeRunner struct {
	content  string
	hooks    *Hooks         // 消息通知
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

// GetFofaCli fofa client
func (p *PipeRunner) GetFofaCli() *gofofa.Client {
	return p.FofaCli
}

// GetLastFile last genrated file
func (p *PipeRunner) GetLastFile() string {
	return p.LastFile
}

type RunnerOption func(*PipeRunner)

// WithHooks user defined hooks
func WithHooks(hooks *Hooks) RunnerOption {
	return func(r *PipeRunner) {
		r.hooks = hooks
	}
}

// 核心函数
func (p *PipeRunner) fork(pipe string) error {
	forkRunner := New(WithHooks(p.hooks))
	forkRunner.LastFile = p.LastFile // 从这里开始分叉
	p.LastTask.Children = append(p.LastTask.Children, forkRunner)
	code, err := workflowast.NewParser().Parse(pipe)
	if err != nil {
		return err
	}
	_, err = forkRunner.Run(code)
	return err
}

func (p *PipeRunner) genData(s string) error {
	var fn string
	var err error
	fn, err = utils.WriteTempFile("", func(f *os.File) error {
		_, err = f.WriteString(s)
		return err
	})
	if err != nil {
		return err
	}

	pt := &PipeTask{
		Name:    "gen",
		Content: fmt.Sprintf("%v", s),
		Outfile: fn,
	}
	p.AddWorkflow(pt)
	return err
}

// 自动补齐url
func (p *PipeRunner) urlFix(field string) error {
	var fn string
	var err error

	if len(field) == 0 {
		return fmt.Errorf("urlfix must has a field")
	}

	fn, err = utils.WriteTempFile("", func(f *os.File) error {
		return utils.EachLine(p.GetLastFile(), func(line string) error {
			v := gjson.Get(line, field).String()
			if !strings.Contains(v, "://") {
				v = "http://" + gjson.Get(line, field).String()
			}
			line, err := sjson.Set(line, field, v)
			if err != nil {
				return err
			}
			_, err = f.WriteString(line + "\n")
			return err
		})
	})
	if err != nil {
		return err
	}

	pt := &PipeTask{
		Name:    "gen",
		Content: fmt.Sprintf("%v", field),
		Outfile: fn,
	}
	p.AddWorkflow(pt)
	return err
}

// New create pipe runner
func New(options ...RunnerOption) *PipeRunner {
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
	err = gf.Register("Fork", r.fork)
	err = gf.Register("gen", r.genData)
	err = gf.Register("urlfix", r.urlFix)
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
		{"Screenshot", screenShot},
	}
	for i := range funcs {
		funcName := funcs[i][0].(string)
		funcBody := funcs[i][1].(func(*PipeRunner, map[string]interface{}) *funcResult)
		gf.Register(funcName, func(p *PipeRunner, params map[string]interface{}) {
			logrus.Debug(funcName+" params:", params)

			s := time.Now()
			result := funcBody(p, params)

			pt := &PipeTask{
				Name:      funcName,
				Content:   fmt.Sprintf("%v", params),
				Outfile:   result.OutFile,
				Artifacts: result.Artifacts,
				Cost:      time.Since(s),
			}
			p.AddWorkflow(pt)
			if p.hooks != nil {
				p.hooks.OnWorkflowFinished(pt)
			}
		})
	}

	logrus.Debug("ast support workflows:", translater.Translators)

	r.gocodeRunner = coderunner.New(coderunner.WithFunctions(&gf))
	for _, opt := range options {
		opt(r)
	}
	return r
}
