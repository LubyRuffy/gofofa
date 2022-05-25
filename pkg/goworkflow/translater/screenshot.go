package translater

import (
	"bytes"
	"text/template"

	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

func screenshotHook(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("screenshot").Parse(`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "{{.URLField}}",
	"timeout": 30,
	"quality": 80,
})`)

	urlField := "url"
	if len(fi.Params) > 0 {
		if v := fi.Params[0].RawString(); len(v) > 0 {
			urlField = v
		}
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, struct {
		URLField string
	}{
		URLField: urlField,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("screenshot", screenshotHook) // screenshot 网页截图
}
