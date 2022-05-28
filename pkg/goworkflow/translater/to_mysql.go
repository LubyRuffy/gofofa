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

	field := ""
	if len(fi.Params) > 1 {
		field = fi.Params[1].RawString()
	}

	dsn := ""
	if len(fi.Params) > 2 {
		dsn = fi.Params[2].RawString()
	}

	return `ToSql(GetRunner(), map[string]interface{}{
	"driver": "mysql",
	"table": "` + table + `",
	"fields": "` + field + `",
	"dsn": "` + dsn + `",
})`
}

func init() {
	register("to_mysql", toMysqlHook)
}
