package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_render_dom(t *testing.T) {
	assert.Equal(t,
		`RenderDOM(GetRunner(), map[string]interface{}{
	"urlField": "host",
	"saveField": "dom_html",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`render_dom("host")`))

	assert.Equal(t,
		`RenderDOM(GetRunner(), map[string]interface{}{
	"urlField": "url",
	"saveField": "dom_html",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`render_dom()`))

	assert.Equal(t,
		`RenderDOM(GetRunner(), map[string]interface{}{
	"urlField": "url",
	"saveField": "dom_html",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`render_dom("")`))

	assert.Equal(t,
		`RenderDOM(GetRunner(), map[string]interface{}{
	"urlField": "host",
	"saveField": "html",
	"timeout": 30,
})
`,
		workflowast.NewParser().MustParse(`render_dom("host", "html")`))

	assert.Equal(t,
		`RenderDOM(GetRunner(), map[string]interface{}{
	"urlField": "host",
	"saveField": "html",
	"timeout": 1,
})
`,
		workflowast.NewParser().MustParse(`render_dom("host", "html", 1)`))

}
