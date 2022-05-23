package funcs

import "github.com/lubyruffy/gofofa/pkg/pipeparser"

func dropHook(fi *pipeparser.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "drop ` + fi.Params[0].RawString() + `",
})`
}

func init() {
	pipeparser.RegisterFunction("drop", dropHook) // 删除字段
}
