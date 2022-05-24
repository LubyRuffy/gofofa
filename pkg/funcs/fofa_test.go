package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
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
		pipeast.NewParser().Parse(`fofa("host=\"https://fofa.info\"", "domain", 1)`))
}
