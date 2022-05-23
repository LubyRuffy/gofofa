package funcs

import (
	"github.com/lubyruffy/gofofa/pkg/fzq"
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
	"github.com/mitchellh/mapstructure"
)

func zqHook(fi *pipeparser.FuncInfo) string {
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": ` + fi.Params[0].String() + `,
})`
}

type zqQueryParams struct {
	Query string `json:"query"`
}

func zqQuery(p *piperunner.PipeRunner, params map[string]interface{}) string {
	var err error
	var options zqQueryParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	name := piperunner.WriteTempJSONFile(nil)
	err = fzq.ZqQuery(options.Query, p.LastFile, name)
	if err != nil {
		panic(err)
	}

	return name
}

func init() {
	piperunner.RegisterWorkflow("zq", zqHook, "ZqQuery", zqQuery) // grep匹配再新增字段
}
