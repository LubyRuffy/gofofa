package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// sort 参数可选
func sortHook(fi *workflowast.FuncInfo) string {
	// 调用zq
	field := ""
	if len(fi.Params) > 0 {
		field = fi.Params[0].RawString()
	}
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "sort ` + field + `",
})`
}

func init() {
	register("sort", sortHook) // 排序
}
