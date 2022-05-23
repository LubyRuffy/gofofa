package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
)

func intHook(fi *pipeast.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cast(this, <{` + fi.Params[0].RawString() + `:int64}>) ",
})`
}

func init() {
	piperunner.RegisterWorkflow("to_int", intHook, "", nil) // 将某个字段转换为int类型
}
