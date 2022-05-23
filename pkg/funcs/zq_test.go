package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_zq(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"a\",\n})\n",
		pipeast.NewParser().Parse(`zq("a")`))

	assertPipeCmd(t, `zq("a")`, `{"a":1}`, "{\"a\":1}\n")
}
