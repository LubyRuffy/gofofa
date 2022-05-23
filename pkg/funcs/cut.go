package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
)

func cutHook(fi *pipeast.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cut ` + fi.Params[0].RawString() + `",
})`
}

func init() {
	piperunner.RegisterWorkflow("cut", cutHook, "", nil) // 剪出要的字段
}
