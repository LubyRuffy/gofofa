package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_add(t *testing.T) {
	assert.Equal(t,
		"AddField(GetRunner(), map[string]interface{}{\n    \"value\": \"newvalue\",\n    \"name\": \"newfield\",\n})\n",
		workflowast.NewParser().MustParse(`add("newfield", "newvalue")`))
}
