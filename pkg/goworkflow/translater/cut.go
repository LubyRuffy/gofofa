package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

func cutHook(fi *workflowast.FuncInfo) string {
	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "cut ` + fi.Params[0].RawString() + `",
})`
}

func init() {
	register("cut", cutHook) // 剪出要的字段
}
