package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_toint(t *testing.T) {
	assert.Equal(t,
		`ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cast(this, <{a:int64}>) ",
})
`,
		pipeparser.NewParser().Parse(`to_int("a")`))

	assertPipeCmd(t, `to_int("a")`, `{"a":"2"}`, `{"a":2}
`)
}
