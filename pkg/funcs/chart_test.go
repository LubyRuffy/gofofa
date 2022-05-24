package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_chart(t *testing.T) {
	assert.Equal(t,
		"GenerateChart(GetRunner(), map[string]interface{}{\n    \"type\": \"pie\",\n    \"title\": \"\",\n})\n",
		pipeast.NewParser().Parse(`chart("pie")`))
	assert.Equal(t,
		"GenerateChart(GetRunner(), map[string]interface{}{\n    \"type\": \"bar\",\n    \"title\": \"a\",\n})\n",
		pipeast.NewParser().Parse(`chart("bar","a")`))
	assert.Equal(t,
		"GenerateChart(GetRunner(), map[string]interface{}{\n    \"type\": \"line\",\n    \"title\": \"a\",\n})\n",
		pipeast.NewParser().Parse(`chart("line","a")`))
}
