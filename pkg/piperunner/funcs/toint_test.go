package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_toint(t *testing.T) {
	assert.Equal(t,
		`ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cast(this, <{a:int64}>) ",
})
`,
		pipeast.NewParser().Parse(`to_int("a")`))

	gf := gorunner.GoFunction{}
	gf.Register("ZqQuery", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := zqQuery(p, params)
		p.(*TestRunner).LastFile = fn
	})
	assertPipeCmdByTestRunner(t, &gf, `to_int("a")`, `{"a":"2"}`, `{"a":2}
`)
}
