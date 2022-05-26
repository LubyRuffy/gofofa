package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_to_msyql(t *testing.T) {
	assert.Panics(t, func() {
		workflowast.NewParser().MustParse(`to_mysql()`)
	})

	assert.Equal(t,
		`ToMysql(GetRunner(), map[string]interface{}{
	"fields": "",
	"table": "tbl",
})
`,
		workflowast.NewParser().MustParse(`to_mysql("tbl")`))

	assert.Equal(t,
		`ToMysql(GetRunner(), map[string]interface{}{
	"fields": "a,b,c",
	"table": "tbl1",
})
`,
		workflowast.NewParser().MustParse(`to_mysql("tbl1", "a,b,c")`))

}
