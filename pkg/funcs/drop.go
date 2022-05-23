package funcs

import "github.com/lubyruffy/gofofa/pkg/pipeparser"

func dropHook(fi *pipeparser.FuncInfo) string {
	//	tmpl, err := template.New("cut").Parse(`RemoveField(GetRunner(), map[string]interface{}{
	//    "fields": {{ . }},
	//})`)
	//	if err != nil {
	//		panic(err)
	//	}
	//	var tpl bytes.Buffer
	//	err = tmpl.Execute(&tpl, fi.Params[0].String())
	//	if err != nil {
	//		panic(err)
	//	}
	//	return tpl.String()

	// 调用zq
	return `ZqQuery(GetRunner(), map[string]interface{}{
    "query": "drop ` + fi.Params[0].RawString() + `",
})`
}

func init() {
	pipeparser.RegisterFunction("drop", dropHook) // 删除字段
}
