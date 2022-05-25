package translater

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"text/template"
)

func rmHook(fi *workflowast.FuncInfo) string {
	if len(fi.Params) == 0 {
		panic(fmt.Errorf("rm must has field params"))
	}

	tmpl, _ := template.New("rm").Parse(`RemoveField(GetRunner(), map[string]interface{}{
   "fields": {{ . }},
})`)
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, fi.Params[0].String())
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("rm", rmHook) // 删除字段
}
