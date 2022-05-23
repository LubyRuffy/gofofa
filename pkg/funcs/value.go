package funcs

import "github.com/lubyruffy/gofofa/pkg/pipeparser"

func valueHook(fi *pipeparser.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "yield ` + fi.Params[0].RawString() + `",
})`
}

func init() {
	pipeparser.RegisterFunction("value", valueHook) // 取值
}
