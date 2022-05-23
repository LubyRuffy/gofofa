package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_cut(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"cut a\",\n})\n",
		pipeparser.NewParser().Parse(`cut("a")`))

	assertPipeCmd(t, `cut("a")`, `{"a":1,"b":2}`, "{\"a\":1}\n")
}
