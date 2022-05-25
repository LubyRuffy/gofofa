package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_chart(t *testing.T) {
	assert.Equal(t,
		"GenerateChart(GetRunner(), map[string]interface{}{\n    \"type\": \"pie\",\n    \"title\": \"\",\n})\n",
		workflowast.NewParser().MustParse(`chart("pie")`))
	assert.Equal(t,
		"GenerateChart(GetRunner(), map[string]interface{}{\n    \"type\": \"bar\",\n    \"title\": \"a\",\n})\n",
		workflowast.NewParser().MustParse(`chart("bar","a")`))
	assert.Equal(t,
		"GenerateChart(GetRunner(), map[string]interface{}{\n    \"type\": \"line\",\n    \"title\": \"a\",\n})\n",
		workflowast.NewParser().MustParse(`chart("line","a")`))
}
