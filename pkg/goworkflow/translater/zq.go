package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

func zqHook(fi *workflowast.FuncInfo) string {
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": ` + fi.Params[0].String() + `,
})`
}

func init() {
	register("zq", zqHook) // grep匹配再新增字段
}
