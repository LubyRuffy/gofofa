package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_rm(t *testing.T) {
	assert.Equal(t,
		`RemoveField(GetRunner(), map[string]interface{}{
   "fields": "title",
})
`,
		workflowast.NewParser().MustParse(`rm("title")`))

}
