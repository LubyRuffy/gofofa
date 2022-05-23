package funcs

import "github.com/lubyruffy/gofofa/pkg/pipeparser"

// sort 参数可选
func sortHook(fi *pipeparser.FuncInfo) string {
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
	pipeparser.RegisterFunction("sort", sortHook) // 排序
}
