package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_uniq(t *testing.T) {
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"uniq \",\n})\n",
		workflowast.NewParser().MustParse(`uniq()`))
	assert.Equal(t,
		"ZqQuery(GetRunner(), map[string]interface{}{\n    \"query\": \"uniq -c\",\n})\n",
		workflowast.NewParser().MustParse(`uniq(true)`))

}
