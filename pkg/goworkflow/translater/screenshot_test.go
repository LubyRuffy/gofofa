package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_screenshot(t *testing.T) {
	assert.Equal(t,
		`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "host",
	"timeout": 30,
	"quality": 80,
})
`,
		workflowast.NewParser().MustParse(`screenshot("host")`))

	assert.Equal(t,
		`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "url",
	"timeout": 30,
	"quality": 80,
})
`,
		workflowast.NewParser().MustParse(`screenshot()`))

	assert.Equal(t,
		`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "url",
	"timeout": 30,
	"quality": 80,
})
`,
		workflowast.NewParser().MustParse(`screenshot("")`))

}
