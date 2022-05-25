package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_grepAdd(t *testing.T) {
	assert.Equal(t,
		"AddField(GetRunner(), map[string]interface{}{\n    \"from\": map[string]interface{}{\n        \"method\": \"grep\",\n        \"field\": \"title\",\n        \"value\": \"(?is)test\",\n    },\n    \"name\": \"new_title\",\n})\n",
		workflowast.NewParser().MustParse(`grep_add("title", "(?is)test", "new_title")`))

}
