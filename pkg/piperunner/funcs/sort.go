package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
)

// sort 参数可选
func sortHook(fi *pipeast.FuncInfo) string {
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
	corefuncs.RegisterWorkflow("sort", sortHook, "", nil) // 排序
}
