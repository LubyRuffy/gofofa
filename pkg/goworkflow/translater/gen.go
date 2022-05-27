package translater

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"text/template"
)

func genHook(fi *workflowast.FuncInfo) string {
	tmpl, err := template.New("gen").Parse(`GenData(GetRunner(), map[string]interface{} {
    "data": {{ .Data }},
})`)
	if err != nil {
		panic(fmt.Errorf("genHook failed: %w", err))
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		Data string
	}{
		Data: fi.Params[0].String(),
	})
	if err != nil {
		panic(fmt.Errorf("genHook failed: %w", err))
	}
	return tpl.String()
}

func init() {
	register("gen", genHook)
}
