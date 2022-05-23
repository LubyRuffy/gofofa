package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_value(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"yield a\",\n})\n",
		pipeast.NewParser().Parse(`value("a")`))

	assertPipeCmd(t, `value("a")`, `{"a":1}`, "1\n")
}
