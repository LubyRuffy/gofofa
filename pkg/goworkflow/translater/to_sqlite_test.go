package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad_to_sqlite(t *testing.T) {
	assert.Panics(t, func() {
		workflowast.NewParser().MustParse(`to_sqlite()`)
	})

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "tbl",
	"dsn": "",
	"fields": "",
})
`,
		workflowast.NewParser().MustParse(`to_sqlite("tbl")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "tbl1",
	"dsn": "a.sqlite3",
	"fields": "",
})
`,
		workflowast.NewParser().MustParse(`to_sqlite("tbl1", "a.sqlite3")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "tbl1",
	"dsn": "",
	"fields": "a,b,c",
})
`,
		workflowast.NewParser().MustParse(`to_sqlite("tbl1", "", "a,b,c")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "tbl1",
	"dsn": "",
	"fields": "a,b,c",
})
`,
		workflowast.NewParser().MustParse(`to_sqlite("tbl1", "?a=b", "a,b,c")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "tbl1",
	"dsn": "a.sqlite",
	"fields": "a,b,c",
})
`,
		workflowast.NewParser().MustParse(`to_sqlite("tbl1", "a.sqlite", "a,b,c")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "tbl1",
	"dsn": "a.sqlite?a=b",
	"fields": "a,b,c",
})
`,
		workflowast.NewParser().MustParse(`to_sqlite("tbl1", "a.sqlite?a=b", "a,b,c")`))

	assert.Equal(t,
		`ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "tbl1",
	"dsn": "/a/b/c/a.sqlite?a=b",
	"fields": "a,b,c",
})
`,
		workflowast.NewParser().MustParse(`to_sqlite("tbl1", "/a/b/c/a.sqlite?a=b", "a,b,c")`))

}
