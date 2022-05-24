package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/gorunner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_rm(t *testing.T) {
	assert.Equal(t,
		`RemoveField(GetRunner(), map[string]interface{}{
   "fields": "title",
})
`,
		pipeast.NewParser().Parse(`rm("title")`))

	gf := gorunner.GoFunction{}
	gf.Register("RemoveField", func(p corefuncs.Runner, params map[string]interface{}) {
		fn, _ := removeField(p, params)
		p.(*TestRunner).LastFile = fn
	})
	assertPipeCmdByTestRunner(t, &gf, `rm("title")`,
		`{"title":"abc","a":1}`,
		`{"a":1}
`)
}
