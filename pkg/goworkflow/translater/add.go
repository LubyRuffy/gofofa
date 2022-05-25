package translater

import (
	"bytes"
	"text/template"

	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

func addHook(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("add").Parse(`AddField(GetRunner(), map[string]interface{}{
    "value": {{ .Value }},
    "name": {{ .Name }},
})`)

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, struct {
		Name  string
		Value string
	}{
		Name:  fi.Params[0].String(),
		Value: fi.Params[1].String(),
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("add", addHook) // grep匹配再新增字段
}
