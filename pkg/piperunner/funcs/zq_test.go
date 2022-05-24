package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_zq(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"a\",\n})\n",
		pipeast.NewParser().Parse(`zq("a")`))

	gf := gorunner.GoFunction{}
	gf.Register("ZqQuery", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := zqQuery(p, params)
		p.(*TestRunner).LastFile = fn
	})
	assertPipeCmdByTestRunner(t, &gf, `zq("a")`, `{"a":1}`, "{\"a\":1}\n")
}
