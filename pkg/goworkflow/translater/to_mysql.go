package translater

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// to_mysql 生成sql文件或者入库
func toMysqlHook(fi *workflowast.FuncInfo) string {
	if len(fi.Params) == 0 {
		panic(fmt.Errorf("to_mysql should has table name as first param"))
	}
	table := fi.Params[0].RawString()

	dsn := ""
	if len(fi.Params) > 1 {
		dsn = fi.Params[1].RawString()
	}

	field := ""
	if len(fi.Params) > 2 {
		field = fi.Params[2].RawString()
	}
	return `ToSql(GetRunner(), map[string]interface{}{
	"driver": "mysql",
	"table": "` + table + `",
	"dsn": "` + dsn + `",
	"fields": "` + field + `",
})`
}

func init() {
	register("to_mysql", toMysqlHook)
}
