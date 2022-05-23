package funcs

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
)

// stats 指定字段统计
func statsHook(fi *pipeast.FuncInfo) string {
	// 调用zq
	field := ""
	if len(fi.Params) > 0 && len(fi.Params[0].RawString()) > 0 {
		field = "yield " + fi.Params[0].RawString() + " | "
	}

	size := ""
	if len(fi.Params) > 1 {
		size = fmt.Sprintf(" | tail %d", fi.Params[1].Int64())
	}
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "` + field + `sort | uniq -c | sort count` + size + `",
})`
}

func init() {
	piperunner.RegisterWorkflow("stats", statsHook, "", nil) // 统计
}
