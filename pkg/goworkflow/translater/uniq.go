package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// uniq 参数可选
func uniqHook(fi *workflowast.FuncInfo) string {
	// 调用zq
	count := ""
	if len(fi.Params) > 0 {
		count = "-c"
	}
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "uniq ` + count + `",
})`
}

func init() {
	register("uniq", uniqHook) // 排序
}
