package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_fofa(t *testing.T) {
	assert.Equal(t,
		`FetchFofa(GetRunner(), map[string]interface{} {
    "query": "host=\"https://fofa.info\"",
    "size": 1,
    "fields": "domain",
})
`,
		workflowast.NewParser().MustParse(`fofa("host=\"https://fofa.info\"", "domain", 1)`))
}
