package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_flat(t *testing.T) {
	assert.Equal(t,
		`FlatArray(GetRunner(), map[string]interface{}{
    "field": "a",
})
`,
		pipeast.NewParser().Parse(`flat("a")`))

	gf := gorunner.GoFunction{}
	gf.Register("FlatArray", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := flatArray(p, params)
		p.(*TestRunner).LastFile = fn
	})
	assertPipeCmdByTestRunner(t, &gf, `flat("a")`, `{"a":[[1], "a", [[[true]]]}`, "1\n\"a\"\ntrue\n")

	//assertPipeCmdError(t, `flat("")`, `{"a":[[1], "a", [[[true]]]}`, "field cannot be empty")
}
