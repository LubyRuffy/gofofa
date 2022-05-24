package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_stats(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"sort | uniq -c | sort count\",\n})\n",
		pipeast.NewParser().Parse(`stats()`))
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"yield a | sort | uniq -c | sort count\",\n})\n",
		pipeast.NewParser().Parse(`stats("a")`))
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"yield a | sort | uniq -c | sort count | tail 2\",\n})\n",
		pipeast.NewParser().Parse(`stats("a", 2)`))

	gf := gorunner.GoFunction{}
	gf.Register("ZqQuery", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := zqQuery(p, params)
		p.(*TestRunner).LastFile = fn
	})
	assertPipeCmdByTestRunner(t, &gf, `stats("a")`, `{"a":1}
{"a":2}
{"a":1}
`, `{"value":2,"count":1}
{"value":1,"count":2}
`)

	// stats("a", 1) 等同于 value("a") | stats("", 1)
	assertPipeCmdByTestRunner(t, &gf, `stats("a", 1)`, `{"a":1}
{"a":2}
{"a":1}
`, `{"value":1,"count":2}
`)

	//
	assertPipeCmdByTestRunner(t, &gf, `stats("", 1)`, `1
2
1
`, `{"value":1,"count":2}
`)

}
