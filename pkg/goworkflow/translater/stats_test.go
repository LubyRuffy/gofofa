package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_stats(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"sort | uniq -c | sort count\",\n})\n",
		workflowast.NewParser().MustParse(`stats()`))
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"yield a | sort | uniq -c | sort count\",\n})\n",
		workflowast.NewParser().MustParse(`stats("a")`))
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"yield a | sort | uniq -c | sort count | tail 2\",\n})\n",
		workflowast.NewParser().MustParse(`stats("a", 2)`))

}
