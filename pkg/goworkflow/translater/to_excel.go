package translater

import (
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// to_excel 生成execl文件
func excelHook(fi *workflowast.FuncInfo) string {
	return `ToExcel(GetRunner(), map[string]interface{}{
})`
}

func init() {
	register("to_excel", excelHook) // excel
}
