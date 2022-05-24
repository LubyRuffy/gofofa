package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
)

func valueHook(fi *pipeast.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "yield ` + fi.Params[0].RawString() + `",
})`
}

func init() {
	corefuncs.RegisterWorkflow("value", valueHook, "", nil) // 取值
}
