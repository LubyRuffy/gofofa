package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_drop(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"drop a\",\n})\n",
		pipeast.NewParser().Parse(`drop("a")`))

	gf := gorunner.GoFunction{}
	gf.Register("ZqQuery", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := zqQuery(p, params)
		p.(*TestRunner).LastFile = fn
	})
	assertPipeCmdByTestRunner(t, &gf, `drop("a")`, `{"a":1,"b":2}`, "{\"b\":2}\n")
	//assertPipeCmd(t, `drop("a")`, `{"a":1,"b":2}`, "{\"b\":2}\n")
}
