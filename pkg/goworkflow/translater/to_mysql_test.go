package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_to_msyql(t *testing.T) {
	assert.Equal(t,
		`ToMysql(GetRunner(), map[string]interface{}{
	"fields": "",
})
`,
		workflowast.NewParser().MustParse(`to_mysql()`))

	assert.Equal(t,
		`ToMysql(GetRunner(), map[string]interface{}{
	"fields": "a,b,c",
})
`,
		workflowast.NewParser().MustParse(`to_mysql("a,b,c")`))

}
