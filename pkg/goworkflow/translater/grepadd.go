package translater

import (
	"bytes"
	"text/template"

	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

func grepAddHook(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("grep_add").Parse(`AddField(GetRunner(), map[string]interface{}{
    "from": map[string]interface{}{
        "method": "grep",
        "field": {{ .Field }},
        "value": {{ .Value }},
    },
    "name": {{ .Name }},
})`)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, struct {
		Field string
		Value string
		Name  string
	}{
		Field: fi.Params[0].String(),
		Value: fi.Params[1].String(),
		Name:  fi.Params[2].String(),
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("grep_add", grepAddHook) // grep匹配再新增字段
}
