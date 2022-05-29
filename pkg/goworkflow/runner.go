package goworkflow

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/lubyruffy/gofofa/pkg/goworkflow/gocodefuncs"

	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/coderunner"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/translater"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/sirupsen/logrus"
)

// PipeTask 每一个pipe执行的任务统计信息
type PipeTask struct {
	Name         string                  // pipe name
	WorkFlowName string                  // workflow name
	Content      string                  // raw content
	Runner       *PipeRunner             // runner
	CallID       int                     // 调用序列
	Cost         time.Duration           // time costs
	Result       *gocodefuncs.FuncResult // 结果
	Children     []*PipeRunner           // fork children
	Fields       []string                // fields list 列名
	Error        error                   // 错误信息
}

// Close remove tmp outfile
func (p *PipeTask) Close() {
	os.Remove(p.Result.OutFile)
	p.Children = nil
}

// Hooks 消息通知
type Hooks struct {
	OnWorkflowFinished func(pt *PipeTask)                                           // 一个workflow完成时的处理
	OnWorkflowStart    func(funcName string, callID int)                            // 一个workflow完成时的处理
	OnLog              func(level logrus.Level, format string, args ...interface{}) // 日志通知
}

// PipeRunner pipe运行器
type PipeRunner struct {
	gf       *coderunner.GoFunction // 函数注册
	ast      *workflowast.Parser    // ast
	content  string                 // 运行的内容
	hooks    *Hooks                 // 消息通知
	Tasks    []*PipeTask            // 执行的所有workflow
	LastTask *PipeTask              // 最后执行的workflow
	LastFile string                 // 最后生成的文件名
	FofaCli  *gofofa.Client         // fofa客户端
	logger   *logrus.Logger
	children []*PipeRunner
	Parent   *PipeRunner

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
	p.content = code
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
	p.Tasks = nil
}

// GetWorkflows all workflows
func (p *PipeRunner) GetWorkflows() []*PipeTask {
	return p.Tasks
}

// AddWorkflow 添加一次任务的日志
func (p *PipeRunner) AddWorkflow(pt *PipeTask) {
	// 可以不写文件
	if pt.Result != nil && len(pt.Result.OutFile) > 0 {
		p.LastFile = pt.Result.OutFile

		logrus.Debug(pt.Name+" write to file: ", pt.Result.OutFile)

		// 取字段列表
		pt.Fields = utils.JSONLineFields(pt.Result.OutFile)
	}
	p.LastTask = pt

	// 把任务也加到上层所有的父节点
	node := p
	for {
		node.Tasks = append(node.Tasks, pt)
		if node.Parent == nil {
			break
		}
		node = node.Parent
	}

	if p.hooks != nil {
		if pt.Error != nil {
			p.hooks.OnLog(logrus.ErrorLevel, "task error: %v", pt.Error)
		}
		p.hooks.OnWorkflowFinished(pt)
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

// WithParent from hook
func WithParent(parent *PipeRunner) RunnerOption {
	return func(r *PipeRunner) {
		r.Parent = parent
	}
}

// WithUserFunction Function to register
func WithUserFunction(funcs ...[]interface{}) RunnerOption {
	return func(r *PipeRunner) {
		r.registerFunctions(funcs...)
	}
}

// WithAST Function to register
func WithAST(ast *workflowast.Parser) RunnerOption {
	return func(r *PipeRunner) {
		r.ast = ast
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

// registerFunctions 注册用户自定义函数，做一层PipeTask封装
func (p *PipeRunner) registerFunctions(funcs ...[]interface{}) {
	for i := range funcs {
		funcName := funcs[i][0].(string)
		funcBody := funcs[i][1].(func(gocodefuncs.Runner, map[string]interface{}) *gocodefuncs.FuncResult)
		p.gf.Register(funcName, func(runner gocodefuncs.Runner, params map[string]interface{}) {
			callID := 1
			node := p
			for {
				callID = len(node.Tasks) + 1
				if node.Parent == nil {
					break
				}
				node = node.Parent
			}
			logrus.Debug(funcName+" params:", params)
			if p.hooks != nil {
				p.hooks.OnWorkflowStart(funcName, callID)
			}
			s := time.Now()

			workflowName := ""
			if p.ast != nil {
				workflowName = p.ast.CallList[callID-1].Name
			}
			pt := &PipeTask{
				Name:         funcName,
				WorkFlowName: workflowName,
				Content:      fmt.Sprintf("%v", params),
				CallID:       callID,
				Runner:       p,
			}

			// 异常捕获
			defer func() {
				if r := recover(); r != nil {
					pt.Error = r.(error)
					pt.Cost = time.Since(s)
					p.AddWorkflow(pt)
					panic(r)
				}
			}()

			result := funcBody(p, params)
			pt.Result = result
			pt.Cost = time.Since(s)

			p.AddWorkflow(pt)
		})
	}
}

// New create pipe runner
func New(options ...RunnerOption) *PipeRunner {
	r := &PipeRunner{
		logger: logrus.New(),
		gf:     &coderunner.GoFunction{},
	}
	var err error

	// 注册底层函数
	err = r.gf.Register("GetRunner", func() *PipeRunner {
		return r
	})
	if err != nil {
		panic(err)
	}
	err = r.gf.Register("Fork", r.fork)
	if err != nil {
		panic(err)
	}

	innerFuncs := [][]interface{}{
		{"RemoveField", gocodefuncs.RemoveField},
		{"FetchFofa", gocodefuncs.FetchFofa},
		{"GenFofaFieldData", gocodefuncs.GenFofaFieldData},
		{"GenerateChart", gocodefuncs.GenerateChart},
		{"ZqQuery", gocodefuncs.ZqQuery},
		{"AddField", gocodefuncs.AddField},
		{"LoadFile", gocodefuncs.LoadFile},
		{"FlatArray", gocodefuncs.FlatArray},
		{"Screenshot", gocodefuncs.ScreenShot},
		{"ToExcel", gocodefuncs.ToExcel},
		{"ToSql", gocodefuncs.ToSql},
		{"GenData", gocodefuncs.GenData},
		{"URLFix", gocodefuncs.UrlFix},
		{"RenderDOM", gocodefuncs.RenderDOM},
	}
	r.registerFunctions(innerFuncs...)

	logrus.Debug("ast support workflows:", translater.Translators)

	r.gocodeRunner = coderunner.New(coderunner.WithFunctions(r.gf))
	for _, opt := range options {
		opt(r)
	}

	return r
}
