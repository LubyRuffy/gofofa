package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_to_mysql(t *testing.T) {
	assert.Panics(t, func() {
		workflowast.NewParser().MustParse(`to_mysql()`)
	})

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "mysql",
	"table": "tbl",
	"dsn": "",
	"fields": "",
})
`,
		workflowast.NewParser().MustParse(`to_mysql("tbl")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "mysql",
	"table": "tbl1",
	"dsn": "root:my-secret-pw@tcp(127.0.0.1:3306)/aaa",
	"fields": "",
})
`,
		workflowast.NewParser().MustParse(`to_mysql("tbl1", "root:my-secret-pw@tcp(127.0.0.1:3306)/aaa")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "mysql",
	"table": "tbl1",
	"dsn": "",
	"fields": "a,b,c",
})
`,
		workflowast.NewParser().MustParse(`to_mysql("tbl1", "", "a,b,c")`))

}
