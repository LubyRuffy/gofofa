package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_flat(t *testing.T) {
	assert.Equal(t,
		`FlatArray(GetRunner(), map[string]interface{}{
    "field": "a",
})
`,
		pipeparser.NewParser().Parse(`flat("a")`))

	assertPipeCmd(t, `flat("a")`, `{"a":[[1], "a", [[[true]]]}`, "1\n\"a\"\ntrue\n")

	assertPipeCmdError(t, `flat("")`, `{"a":[[1], "a", [[[true]]]}`, "field cannot be empty")
}
