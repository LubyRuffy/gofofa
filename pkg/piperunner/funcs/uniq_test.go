package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_uniq(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"uniq \",\n})\n",
		pipeast.NewParser().Parse(`uniq()`))
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"uniq -c\",\n})\n",
		pipeast.NewParser().Parse(`uniq(true)`))

	assertPipeCmd(t, `uniq()`, "1\n2\n1\n", "1\n2\n1\n")
	assertPipeCmd(t, `uniq()`, "1\n1\n2\n", "1\n2\n")
	assertPipeCmd(t, `uniq(true)`, "1\n2\n1\n", `{"value":1,"count":1}
{"value":2,"count":1}
{"value":1,"count":1}
`)
	assertPipeCmd(t, `uniq(true)`, "1\n1\n2\n", `{"value":1,"count":2}
{"value":2,"count":1}
`)

	// 先sort再uniq
	assertPipeCmd(t, `sort() | uniq(true)`, "1\n2\n1\n", `{"value":1,"count":2}
{"value":2,"count":1}
`)

}
