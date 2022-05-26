package goworkflow

import (
	"bytes"
	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// 返回值表明是错误并且匹配到了错误返回
func defaultErrorHandler(t *testing.T, err error) bool {
	assert.Nil(t, err)
	return false
}

func assertPipeCmd(t *testing.T, workflow string, testData string,
	options ...interface{}) string {
	var err error
	r := New()

	errHandler := defaultErrorHandler
	if len(options) > 0 {
		errHandler = options[0].(func(t *testing.T, err error) bool)
	}

	// 写入数据文件
	if len(testData) > 0 {
		r.LastFile, err = utils.WriteTempFile("", func(f *os.File) error {
			_, err = f.WriteString(testData)
			return err
		})
		if errHandler(t, err) {
			return ""
		}
	}

	// 执行代码
	code, err := workflowast.NewParser().Parse(workflow)
	if errHandler(t, err) {
		return ""
	}

	_, err = r.Run(code)
	if errHandler(t, err) {
		return ""
	}

	var data []byte
	data, err = os.ReadFile(r.GetLastFile())
	if errHandler(t, err) {
		return ""
	}
	return string(data)
}

func assertPipeCmdByTestRunner(t *testing.T, workflow string, testData string,
	except string, options ...interface{}) {
	data := assertPipeCmd(t, workflow, testData, options...)
	assert.Equal(t, except, data)
}

func assertPipeCmdByTestRunnerError(t *testing.T,
	workflow string, testData string, errorStr string) {
	assertPipeCmdByTestRunner(t, workflow, testData, "", func(t *testing.T, err error) bool {
		if err != nil {
			assert.Contains(t, err.Error(), errorStr)
			return true
		}
		return false
	})
}

func TestNew(t *testing.T) {
	assertPipeCmdByTestRunner(t, `add("newfield", "newvalue")`,
		`{"title":"Test123"}
{"title":"123test456"}`,
		`{"title":"Test123","newfield":"newvalue"}
{"title":"123test456","newfield":"newvalue"}`)

	// chart格式错误
	assertPipeCmdByTestRunnerError(t, `chart("line","a")`,
		`{"title":"Test123"}`,
		`"value" and "count" field is needed`)
	// chart正确
	assertPipeCmdByTestRunner(t, `chart("bar","a")`,
		`{"value":"Test123","count":10}`,
		"{\"value\":\"Test123\",\"count\":10}")
	// chart正确
	assertPipeCmdByTestRunner(t, `chart("pie","a")`,
		`{"value":"Test123","count":10}`,
		"{\"value\":\"Test123\",\"count\":10}")

	assertPipeCmdByTestRunner(t, `cut("a")`, `{"a":1,"b":2}`, "{\"a\":1}\n")
	//assertPipeCmd(t, `cut("a")`, `{"a":1,"b":2}`, "{\"a\":1}\n")

	assertPipeCmdByTestRunner(t, `drop("a")`, `{"a":1,"b":2}`, "{\"b\":2}\n")
	//assertPipeCmd(t, `drop("a")`, `{"a":1,"b":2}`, "{\"b\":2}\n")

	assertPipeCmdByTestRunner(t, `flat("a")`, `{"a":[[1], "a", [[[true]]]}`, "1\n\"a\"\ntrue\n")

	assertPipeCmdByTestRunnerError(t, `flat("")`, `{"a":[[1], "a", [[[true]]]}`,
		"field cannot be empty")

	assertPipeCmdByTestRunner(t, `grep_add("title", "(?is)test", "new_title")`,
		`{"title":"Test123"}
{"title":"123test456"}`,
		`{"title":"Test123","new_title":[["Test"]]}
{"title":"123test456","new_title":[["test"]]}`)

	// 正常
	assertPipeCmdByTestRunner(t, `rm("title")`,
		`{"title":"abc","a":1}`,
		`{"a":1}
`)

	// 字段不存在
	assertPipeCmdByTestRunner(t, `rm("title")`,
		`{"a":1}`,
		`{"a":1}
`)

	// 不提供字段
	assertPipeCmdByTestRunnerError(t, `rm()`,
		`{"a":1}`,
		`rm must has field params`)

	// 提供空字段
	assertPipeCmdByTestRunnerError(t, `rm("")`,
		`{"a":1}`,
		`path cannot be empty`)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	assert.Contains(t, assertPipeCmd(t, `gen("{\"host\":\"`+ts.URL+`\"}") | screenshot("host")`, ``), "screenshot_filepath")

	assertPipeCmdByTestRunner(t, `sort("a")`, `{"a":2}
{"a":1}`, `{"a":1}
{"a":2}
`)

	assertPipeCmdByTestRunner(t, `sort()`, "1\n2\n1\n", `1
1
2
`)

	assertPipeCmdByTestRunner(t, `stats("a")`, `{"a":1}
{"a":2}
{"a":1}
`, `{"value":2,"count":1}
{"value":1,"count":2}
`)

	// stats("a", 1) 等同于 value("a") | stats("", 1)
	assertPipeCmdByTestRunner(t, `stats("a", 1)`, `{"a":1}
{"a":2}
{"a":1}
`, `{"value":1,"count":2}
`)

	//
	assertPipeCmdByTestRunner(t, `stats("", 1)`, `1
2
1
`, `{"value":1,"count":2}
`)

	assertPipeCmdByTestRunner(t, `to_int("a")`, `{"a":"2"}`, `{"a":2}
`)

	assertPipeCmdByTestRunner(t, `uniq()`, "1\n2\n1\n", "1\n2\n1\n")
	assertPipeCmdByTestRunner(t, `uniq()`, "1\n1\n2\n", "1\n2\n")
	assertPipeCmdByTestRunner(t, `uniq(true)`, "1\n2\n1\n", `{"value":1,"count":1}
{"value":2,"count":1}
{"value":1,"count":1}
`)
	assertPipeCmdByTestRunner(t, `uniq(true)`, "1\n1\n2\n", `{"value":1,"count":2}
{"value":2,"count":1}
`)

	// 先sort再uniq
	assertPipeCmdByTestRunner(t, `sort() | uniq(true)`, "1\n2\n1\n", `{"value":1,"count":2}
{"value":2,"count":1}
`)

	assertPipeCmdByTestRunner(t, `value("a")`, `{"a":1}`, "1\n")

	assertPipeCmdByTestRunner(t, `zq("a")`, `{"a":1}`, "{\"a\":1}\n")
}

func TestLoad_fork(t *testing.T) {
	ast := workflowast.NewParser().MustParse(`load("../../data/forktest.json") | [cut("a") & cut("b")]`)
	p := New()
	_, err := p.Run(ast)
	assert.Nil(t, err)
	res, err := os.ReadFile(p.LastFile)
	assert.Nil(t, err)
	assert.Equal(t, `{"a":1,"b":2}`, string(res))
	assert.Equal(t, 2, len(p.LastTask.Children))
	res, err = os.ReadFile(p.LastTask.Children[0].GetLastFile())
	assert.Nil(t, err)
	assert.Equal(t, "{\"a\":1}\n", string(res))
	res, err = os.ReadFile(p.LastTask.Children[1].GetLastFile())
	assert.Nil(t, err)
	assert.Equal(t, "{\"b\":2}\n", string(res))
}

func TestLoad_fofa(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/info/my":
			w.Write([]byte(`{"error":false,"email":"` + r.FormValue("email") + `","fcoin":10,"isvip":true,"vip_level":1}`))
		case "/api/v1/search/all":
			w.Write([]byte(`{"error":false,"size":12345678,"page":1,"mode":"extended","query":"host=\"https://fofa.info\"","results":[["fofa1.info"]]}`))
		}
	}))
	defer ts.Close()

	var err error
	code := workflowast.NewParser().MustParse(`fofa("host=\"https://fofa1.info\"", "domain", 1)`)
	p := New()
	p.FofaCli, err = gofofa.NewClient(ts.URL)
	assert.Nil(t, err)
	_, err = p.Run(code)
	assert.Nil(t, err)

	content, err := os.ReadFile(p.LastFile)
	assert.Nil(t, err)
	assert.Equal(t, `{"domain":"fofa1.info"}
`,
		string(content))
}

func TestPipeRunner_DumpTasks(t *testing.T) {
	tpl, err := template.New("tasks").Funcs(template.FuncMap{
		"HasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
	}).Parse(`{{ if HasPrefix . "aaa" }}yes{{ end }}`)
	assert.Nil(t, err)
	var out bytes.Buffer
	err = tpl.Execute(&out, "aaa")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "yes", out.String())

	p := New()
	_, err = p.Run(workflowast.NewParser().MustParse(`load("../../data/forktest.json") | [cut("a")&cut("b")]`))
	assert.Nil(t, err)
	c := p.DumpTasks(false)
	assert.Contains(t, c, "fork")
}

func TestPipeRunner_Close(t *testing.T) {
	p := New()
	_, err := p.Run(workflowast.NewParser().MustParse(`load("../../data/forktest.json") | cut("a")`))
	assert.Nil(t, err)
	c, err := os.ReadFile(p.LastFile)
	assert.Nil(t, err)
	assert.True(t, len(c) > 0)

	p.Close()
	_, err = os.ReadFile(p.LastFile)
	assert.Error(t, err)
}
