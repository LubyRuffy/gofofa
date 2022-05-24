package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_uniq(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"uniq \",\n})\n",
		pipeast.NewParser().Parse(`uniq()`))
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"uniq -c\",\n})\n",
		pipeast.NewParser().Parse(`uniq(true)`))

	gf := gorunner.GoFunction{}
	gf.Register("ZqQuery", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := zqQuery(p, params)
		p.(*TestRunner).LastFile = fn
	})
	assertPipeCmdByTestRunner(t, &gf, `uniq()`, "1\n2\n1\n", "1\n2\n1\n")
	assertPipeCmdByTestRunner(t, &gf, `uniq()`, "1\n1\n2\n", "1\n2\n")
	assertPipeCmdByTestRunner(t, &gf, `uniq(true)`, "1\n2\n1\n", `{"value":1,"count":1}
{"value":2,"count":1}
{"value":1,"count":1}
`)
	assertPipeCmdByTestRunner(t, &gf, `uniq(true)`, "1\n1\n2\n", `{"value":1,"count":2}
{"value":2,"count":1}
`)

	// 先sort再uniq
	assertPipeCmdByTestRunner(t, &gf, `sort() | uniq(true)`, "1\n2\n1\n", `{"value":1,"count":2}
{"value":2,"count":1}
`)

}
