package translater

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// to_sqlite 生成sql文件或者入库
func toSqilteHook(fi *workflowast.FuncInfo) string {
	if len(fi.Params) == 0 {
		panic(fmt.Errorf("to_sqlite should has table name as first param"))
	}
	table := fi.Params[0].RawString()

	field := ""
	if len(fi.Params) > 1 {
		field = fi.Params[1].RawString()
	}
	return `ToSql(GetRunner(), map[string]interface{}{
	"driver": "sqlite3",
	"table": "` + table + `",
	"fields": "` + field + `",
})`
}

func init() {
	register("to_sqlite", toSqilteHook)
}
