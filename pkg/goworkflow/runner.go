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
	Fields    []string      // fields list 列名
	CallID    int           // 调用序列
}

// Close remove tmp outfile
func (p *PipeTask) Close() {
	os.Remove(p.Outfile)
	p.Children = nil
}

// Hooks 消息通知
type Hooks struct {
	OnWorkflowFinished func(pt *PipeTask)                                           // 一个workflow完成时的处理
	OnLog              func(level logrus.Level, format string, args ...interface{}) // 日志通知
}

// PipeRunner pipe运行器
type PipeRunner struct {
	content  string
	hooks    *Hooks         // 消息通知
	Tasks    []*PipeTask    // 执行的所有workflow
	LastTask *PipeTask      // 最后执行的workflow
	LastFile string         // 最后生成的文件名
	FofaCli  *gofofa.Client // fofa客户端
	logger   *logrus.Logger
	children []*PipeRunner
	parent   *PipeRunner

	gocodeRunner *coderunner.Runner
}

// Logf 打印日志
func (p *PipeRunner) Logf(level logrus.Level, format string, args ...interface{}) {
	if p.hooks != nil && p.hooks.OnLog != nil {
		p.hooks.OnLog(level, format, args...)
	}
	p.logger.Logf(level, format, args...)
}

// Debugf 打印调试日志
func (p *PipeRunner) Debugf(format string, args ...interface{}) {
	p.Logf(logrus.DebugLevel, format, args...)
}

func (p *PipeRunner) Warnf(format string, args ...interface{}) {
	p.logger.Logf(logrus.WarnLevel, format, args...)
}

// Run go code, not workflow
func (p *PipeRunner) Run(code string) (reflect.Value, error) {
	return p.gocodeRunner.Run(code)
}

// Close remove tmp outfile
func (p *PipeRunner) Close() {
	p.children = nil
	p.LastFile = ""
	p.LastTask = nil
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
	// 可以不写文件
	if len(pt.Outfile) > 0 {
		p.LastFile = pt.Outfile

		logrus.Debug(pt.Name+" write to file: ", pt.Outfile)

		// 取字段列表
		d, err := utils.ReadFirstLineOfFile(pt.Outfile)
		if err != nil {
			panic(fmt.Errorf("ReadFirstLineOfFile failed: %w", err))
		}
		v := gjson.ParseBytes(d)
		v.ForEach(func(key, value gjson.Result) bool {
			pt.Fields = append(pt.Fields, key.String())
			return true
		})
	}

	node := p
	for {
		node.Tasks = append(node.Tasks, pt)
		if node.parent == nil {
			break
		}
		node = node.parent
	}

	p.LastTask = pt
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

// WithParent from hook
func WithParent(parent *PipeRunner) RunnerOption {
	return func(r *PipeRunner) {
		r.parent = parent
	}
}

// 核心函数
func (p *PipeRunner) fork(pipe string) error {
	forkRunner := New(WithHooks(p.hooks), WithParent(p))
	forkRunner.LastFile = p.LastFile // 从这里开始分叉
	p.LastTask.Children = append(p.LastTask.Children, forkRunner)
	code, err := workflowast.NewParser().Parse(pipe)
	if err != nil {
		return err
	}
	p.children = append(p.children, forkRunner)
	_, err = forkRunner.Run(code)
	return err
}

// 自动补齐url
func (p *PipeRunner) urlFix(field string) error {
	var fn string
	var err error

	if len(field) == 0 {
		return fmt.Errorf("urlFix must has a field")
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
		Name:    "urlFix",
		Content: fmt.Sprintf("%v", field),
		Outfile: fn,
	}
	p.AddWorkflow(pt)
	return err
}

// New create pipe runner
func New(options ...RunnerOption) *PipeRunner {
	r := &PipeRunner{
		logger: logrus.New(),
	}
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
		{"ToExcel", toExcel},
		{"ToSql", toSql},
		{"GenData", genData},
	}
	for i := range funcs {
		funcName := funcs[i][0].(string)
		funcBody := funcs[i][1].(func(*PipeRunner, map[string]interface{}) *funcResult)
		gf.Register(funcName, func(p *PipeRunner, params map[string]interface{}) {
			logrus.Debug(funcName+" params:", params)

			s := time.Now()
			result := funcBody(p, params)

			callID := 1
			node := p
			for {
				callID = len(node.Tasks) + 1
				if node.parent == nil {
					break
				}
				node = node.parent
			}
			pt := &PipeTask{
				Name:      funcName,
				Content:   fmt.Sprintf("%v", params),
				Outfile:   result.OutFile,
				Artifacts: result.Artifacts,
				Cost:      time.Since(s),
				CallID:    callID,
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
