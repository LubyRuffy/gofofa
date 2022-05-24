package corefuncs

import (
	"fmt"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var (
	functions sync.Map // 全局注册的底层函数
)

// PipeTask 每一个pipe执行的任务统计信息
type PipeTask struct {
	Name           string        // pipe name
	Content        string        // raw content
	Outfile        string        // tmp json file 统一格式
	GeneratedFiles []string      // files to archive 非json格式的文件不往后进行传递
	Cost           time.Duration // time costs
	Children       []Runner      // fork children
}

// Close remove tmp outfile
func (p *PipeTask) Close() {
	os.Remove(p.Outfile)
}

// Runner 的接口定义
type Runner interface {
	GetFofaCli() *gofofa.Client
	GetLastFile() string
	AddWorkflow(*PipeTask)
	GetWorkflows() []*PipeTask
}

// RegisterWorkflow 注册workflow
// 第一个参数是workflow名称；
// 第二个参数是workflow转换为函数调用字符串的函数
// 第三个参数是底层函数的名称
// 第四个参数是一个回调函数，参数是传递的参数，返回值是生成的文件名
// 第三四个参数可以留空值，表明只注册到语法解析器中去
func RegisterWorkflow(workflow string, transFunc pipeast.FunctionTranslateHook,
	funcName string, funcBody func(Runner, map[string]interface{}) (string, []string)) {

	// 解析器的函数注册
	if len(workflow) > 0 {
		pipeast.RegisterFunction(workflow, transFunc)
	}

	// 注册底层函数
	if len(funcName) > 0 {
		// 执行并且自动生成任务队列
		functions.Store(funcName, func(p Runner, params map[string]interface{}) {
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
}

// SupportWorkflows 手动加载，否则init不执行
func SupportWorkflows() []string {
	return pipeast.SupportWorkflows()
}

// Range 底层函数的遍历
func Range(f func(key, value any) bool) {
	functions.Range(f)
}
