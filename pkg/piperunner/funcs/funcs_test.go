package funcs

import (
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//
//func assertPipeCmd(t *testing.T, pipeCmd string, jsonData string, except string) {
//	p := piperunner.New()
//
//	// write json to file
//	f, err := os.CreateTemp(os.TempDir(), "piperunner_")
//	assert.Nil(t, err)
//	defer f.Close()
//	_, err = f.WriteString(jsonData)
//	assert.Nil(t, err)
//
//	p.LastFile = f.Name()
//
//	// run
//	_, err = p.Run(pipeast.NewParser().Parse(pipeCmd))
//	assert.Nil(t, err)
//
//	res, err := os.ReadFile(p.LastFile)
//	assert.Nil(t, err)
//
//	assert.Equal(t, except, string(res))
//}
//
//func assertPipeCmdError(t *testing.T, pipeCmd string, jsonData string, errStr string) {
//	p := piperunner.New()
//
//	// write json to file
//	f, err := os.CreateTemp(os.TempDir(), "piperunner_")
//	assert.Nil(t, err)
//	defer f.Close()
//	_, err = f.WriteString(jsonData)
//	assert.Nil(t, err)
//
//	p.LastFile = f.Name()
//
//	// run
//	_, err = p.Run(pipeast.NewParser().Parse(pipeCmd))
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), errStr)
//}

type TestRunner struct {
	LastFile string
}

func (tr TestRunner) GetLastFile() string {
	return tr.LastFile
}
func (tr TestRunner) GetFofaCli() *gofofa.Client {
	return nil
}
func (tr TestRunner) AddWorkflow(pt *corefuncs.PipeTask) {
}
func (tr TestRunner) GetWorkflows() []*corefuncs.PipeTask {
	return nil
}

func assertPipeCmdByTestRunner(t *testing.T, gf *gorunner.GoFunction,
	workflow string, testData string, except string) {
	var err error
	r := &TestRunner{}

	// 写入数据文件
	r.LastFile = WriteTempFile("", func(f *os.File) {
		f.WriteString(testData)
	})

	// 填充函数
	err = gf.Register("GetRunner", func() *TestRunner {
		return r
	})
	assert.Nil(t, err)

	assert.Nil(t, err)

	// 执行代码
	gr := gorunner.New(gorunner.WithFunctions(gf))
	code := pipeast.NewParser().Parse(workflow)
	_, err = gr.Run(code)
	assert.Nil(t, err)

	var data []byte
	data, err = os.ReadFile(r.GetLastFile())
	assert.Equal(t, except, string(data))
}
