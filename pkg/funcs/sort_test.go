package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_sort(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"sort a\",\n})\n",
		pipeparser.NewParser().Parse(`sort("a")`))

	assertPipeCmd(t, `sort("a")`, `{"a":2}
{"a":1}`, `{"a":1}
{"a":2}
`)

}
