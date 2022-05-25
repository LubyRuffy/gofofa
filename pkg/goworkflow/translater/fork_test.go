package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_fork(t *testing.T) {
	assert.Equal(t,
		`Fork("a()")
`,
		workflowast.NewParser().MustParse(`fork("a()")`))

	//gf := coderunner.GoFunction{}
	//assertPipeCmdByTestRunner(t, &gf, `load("../../../data/forktest.json") | [cut("a") & cut("b")]`,
	//	`{"a":1,"b":2}`, `{"a":1,"b":2}`)
}
