package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// to_mysql 生成sql文件
func toMysqlHook(fi *workflowast.FuncInfo) string {
	field := ""
	if len(fi.Params) > 0 {
		field = fi.Params[0].RawString()
	}
	return `ToMysql(GetRunner(), map[string]interface{}{
	"fields": "` + field + `",
})`
}

func init() {
	register("to_mysql", toMysqlHook)
}
