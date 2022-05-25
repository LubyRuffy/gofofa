package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

func forkHook(fi *workflowast.FuncInfo) string {
	return `Fork(` + fi.Params[0].String() + `)`
}

func init() {
	register("fork", forkHook) // fork
}
