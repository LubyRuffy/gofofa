package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/fzq"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/mitchellh/mapstructure"
)

func zqHook(fi *pipeast.FuncInfo) string {
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": ` + fi.Params[0].String() + `,
})`
}

type zqQueryParams struct {
	Query string `json:"query"`
}

func zqQuery(p corefuncs.Runner, params map[string]interface{}) (string, []string) {
	var err error
	var options zqQueryParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	name := WriteTempFile(".json", nil)
	err = fzq.ZqQuery(options.Query, p.GetLastFile(), name)
	if err != nil {
		panic(err)
	}

	return name, nil
}

func init() {
	corefuncs.RegisterWorkflow("zq", zqHook, "ZqQuery", zqQuery) // grep匹配再新增字段
}
