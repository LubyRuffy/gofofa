package translater

import (
	"bytes"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"text/template"
)

// render_dom(<urlField:"url">,[saveField:"dom_html"],[timeout:30])
func renderdomHook(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("screenshot").Parse(`RenderDOM(GetRunner(), map[string]interface{}{
	"urlField": "{{.URLField}}",
	"saveField": "{{.SaveField}}",
	"timeout": {{.TimeOut}},
})`)

	urlField := "url"
	if len(fi.Params) > 0 {
		if v := fi.Params[0].RawString(); len(v) > 0 {
			urlField = v
		}
	}
	saveField := "dom_html"
	if len(fi.Params) > 1 {
		if v := fi.Params[1].RawString(); len(v) > 0 {
			saveField = v
		}
	}
	timeOut := 30
	if len(fi.Params) > 2 {
		timeOut = int(fi.Params[2].Int64())
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, struct {
		URLField  string
		SaveField string
		TimeOut   int
	}{
		URLField:  urlField,
		SaveField: saveField,
		TimeOut:   timeOut,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("render_dom", renderdomHook) // screenshot 网页截图
}
