package translater

import (
	"bytes"
	"text/template"

	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// screenshot(<urlField:"url">,[saveField:"screenshot_filepath"],[timeout:30])
func screenshotHook(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("screenshot").Parse(`Screenshot(GetRunner(), map[string]interface{}{
	"urlField": "{{.URLField}}",
	"saveField": "{{.SaveField}}",
	"timeout": {{.TimeOut}},
	"quality": 80,
})`)

	urlField := "url"
	if len(fi.Params) > 0 {
		if v := fi.Params[0].RawString(); len(v) > 0 {
			urlField = v
		}
	}
	saveField := "screenshot_filepath"
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
	register("screenshot", screenshotHook) // screenshot 网页截图
}
