package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_drop(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"drop a\",\n})\n",
		pipeparser.NewParser().Parse(`drop("a")`))

	assertPipeCmd(t, `drop("a")`, `{"a":1,"b":2}`, "{\"b\":2}\n")
}
