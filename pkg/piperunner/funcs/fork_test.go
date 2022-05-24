package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_fork(t *testing.T) {
	assert.Equal(t,
		`Fork("a()")
`,
		pipeast.NewParser().Parse(`fork("a()")`))

	//gf := gorunner.GoFunction{}
	//assertPipeCmdByTestRunner(t, &gf, `load("../../../data/forktest.json") | [cut("a") & cut("b")]`,
	//	`{"a":1,"b":2}`, `{"a":1,"b":2}`)
}
