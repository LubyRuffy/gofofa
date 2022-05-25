package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_toint(t *testing.T) {
	assert.Equal(t,
		`ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cast(this, <{a:int64}>) ",
})
`,
		workflowast.NewParser().MustParse(`to_int("a")`))

}
