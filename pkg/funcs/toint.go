package funcs

import "github.com/lubyruffy/gofofa/pkg/pipeparser"

func intHook(fi *pipeparser.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cast(this, <{` + fi.Params[0].RawString() + `:int64}>) ",
})`
}

func init() {
	pipeparser.RegisterFunction("to_int", intHook) // 将某个字段转换为int类型
}
