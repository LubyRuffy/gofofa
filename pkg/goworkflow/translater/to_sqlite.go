package translater

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"strings"
)

// to_sqlite 生成sql文件或者入库
func toSqilteHook(fi *workflowast.FuncInfo) string {
	if len(fi.Params) == 0 {
		panic(fmt.Errorf("to_sqlite should has table name as first param"))
	}
	table := fi.Params[0].RawString()

	dsn := ""
	if len(fi.Params) > 1 {
		dsn = fi.Params[1].RawString()
		dsn = strings.ReplaceAll(dsn, "\\", "/") // windows文件
	}

	field := ""
	if len(fi.Params) > 2 {
		field = fi.Params[2].RawString()
	}
	return `ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "` + table + `",
	"dsn": "` + dsn + `",
	"fields": "` + field + `",
})`
}

func init() {
	register("to_sqlite", toSqilteHook)
}
