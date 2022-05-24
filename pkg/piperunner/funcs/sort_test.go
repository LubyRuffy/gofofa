package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_sort(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"sort a\",\n})\n",
		pipeast.NewParser().Parse(`sort("a")`))

	gf := gorunner.GoFunction{}
	gf.Register("ZqQuery", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := zqQuery(p, params)
		p.(*TestRunner).LastFile = fn
	})

	assertPipeCmdByTestRunner(t, &gf, `sort("a")`, `{"a":2}
{"a":1}`, `{"a":1}
{"a":2}
`)

	assertPipeCmdByTestRunner(t, &gf, `sort()`, "1\n2\n1\n", `1
1
2
`)

}
