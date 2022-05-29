package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipeRunner_urlfix(t *testing.T) {
	assert.Equal(t,
		"URLFix(GetRunner(), map[string]interface{}{\n    \"urlField\": \"url\",\n})\n",
		workflowast.NewParser().MustParse(`urlfix()`))
	assert.Equal(t,
		"URLFix(GetRunner(), map[string]interface{}{\n    \"urlField\": \"host\",\n})\n",
		workflowast.NewParser().MustParse(`urlfix("host")`))

}
