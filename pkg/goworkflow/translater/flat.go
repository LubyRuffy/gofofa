package translater

import (
	"bytes"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"text/template"
)

// flat 将数组分成多条一维记录并且展开
// flat("a")
// 第一个参数是字段名称
// 注意：空值会移除
func flat(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("flat").Parse(`FlatArray(GetRunner(), map[string]interface{}{
    "field": {{ . }},
})`)
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, fi.Params[0].String())
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("flat", flat)
}
