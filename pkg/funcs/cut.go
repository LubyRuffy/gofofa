package funcs

import "github.com/lubyruffy/gofofa/pkg/pipeparser"

func cutHook(fi *pipeparser.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cut ` + fi.Params[0].RawString() + `",
})`
}

func init() {
	pipeparser.RegisterFunction("cut", cutHook) // 剪出要的字段
}
