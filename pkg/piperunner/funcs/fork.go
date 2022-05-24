package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
)

func forkHook(fi *pipeast.FuncInfo) string {
	return `Fork(` + fi.Params[0].String() + `)`
}

func init() {
	corefuncs.RegisterWorkflow("fork", forkHook, "", nil) // fork
}
