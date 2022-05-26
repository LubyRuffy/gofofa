package translater

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// to_mysql 生成sql文件
func toMysqlHook(fi *workflowast.FuncInfo) string {
	if len(fi.Params) == 0 {
		panic(fmt.Errorf("to_mysql should has table name as first param"))
	}
	table := fi.Params[0].RawString()

	field := ""
	if len(fi.Params) > 1 {
		field = fi.Params[1].RawString()
	}
	return `ToMysql(GetRunner(), map[string]interface{}{
	"fields": "` + field + `",
	"table": "` + table + `",
})`
}

func init() {
	register("to_mysql", toMysqlHook)
}
