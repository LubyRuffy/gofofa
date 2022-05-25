package translater

import (
	"bytes"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"text/template"
)

func fofaHook(fi *workflowast.FuncInfo) string {
	tmpl, err := template.New("fofa").Parse(`FetchFofa(GetRunner(), map[string]interface{} {
    "query": {{ .Query }},
    "size": {{ .Size }},
    "fields": {{ .Fields }},
})`)
	if err != nil {
		panic(err)
	}
	var size int64 = 10
	fields := "`host,title`"
	if len(fi.Params) > 1 {
		fields = fi.Params[1].String()
	}
	if len(fi.Params) > 2 {
		size = fi.Params[2].Int64()
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		Query  string
		Size   int64
		Fields string
	}{
		Query:  fi.Params[0].String(),
		Fields: fields,
		Size:   size,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("fofa", fofaHook)
}
