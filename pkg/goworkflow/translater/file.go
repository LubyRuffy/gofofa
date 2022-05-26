package translater

import (
	"bytes"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"text/template"
)

func loadHook(fi *workflowast.FuncInfo) string {
	tmpl, err := template.New("load").Parse(`LoadFile(GetRunner(), map[string]interface{} {
    "file": {{ .File }},
})`)
	if err != nil {
		panic(fmt.Errorf("loadHook failed: %w", err))
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		File string
	}{
		File: fi.Params[0].String(),
	})
	if err != nil {
		panic(fmt.Errorf("loadHook failed: %w", err))
	}
	return tpl.String()
}

func init() {
	register("load", loadHook)
}
