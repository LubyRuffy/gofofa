package translater

import (
	"bytes"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"text/template"
)

// chart 生成报表
// chart("pie")
// 第一个参数是报表类型
func chart(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("chart").Parse(`GenerateChart(GetRunner(), map[string]interface{}{
    "type": {{ .Type }},
    "title": "{{ .Title }}",
})`)

	typeStr := "bar"
	if len(fi.Params) > 0 {
		typeStr = fi.Params[0].String()
	}
	title := ""
	if len(fi.Params) > 1 {
		title = fi.Params[1].RawString()
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, struct {
		Type  string
		Title string
	}{
		Type:  typeStr,
		Title: title,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("chart", chart)
}
