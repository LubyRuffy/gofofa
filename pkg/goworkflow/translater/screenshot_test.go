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
	"saveField": "screenshot_filepath",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`screenshot("host")`))

	assert.Equal(t,
		`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "url",
	"saveField": "screenshot_filepath",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`screenshot()`))

	assert.Equal(t,
		`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "url",
	"saveField": "screenshot_filepath",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`screenshot("")`))

	assert.Equal(t,
		`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "host",
	"saveField": "sc_filepath",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`screenshot("host", "sc_filepath")`))

	assert.Equal(t,
		`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "host",
	"saveField": "sc_filepath",
	"timeout": 1,
})
`,
		workflowast.NewParser().MustParse(`screenshot("host", "sc_filepath", 1)`))

}
