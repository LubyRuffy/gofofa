package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
)

// uniq 参数可选
func uniqHook(fi *pipeast.FuncInfo) string {
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
	piperunner.RegisterWorkflow("uniq", uniqHook, "", nil) // 排序
}
