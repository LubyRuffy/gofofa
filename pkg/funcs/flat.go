package funcs

import (
	"bytes"
	"errors"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"os"
	"text/template"
)

// flat 将数组分成多条一维记录并且展开
// flat("a")
// 第一个参数是字段名称
// 注意：空值会移除
func flat(fi *pipeast.FuncInfo) string {
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

func jsonArrayEnum(node gjson.Result, f func(result gjson.Result)) {
	if node.IsArray() {
		for _, child := range node.Array() {
			jsonArrayEnum(child, f)
		}
	} else {
		f(node)
	}
}

type flatParams struct {
	Field string
}

func flatArray(p *piperunner.PipeRunner, params map[string]interface{}) string {
	var err error
	var options flatParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	if len(options.Field) == 0 {
		panic(errors.New("field cannot be empty"))
	}

	return piperunner.WriteTempJSONFile(func(f *os.File) {
		piperunner.EachLine(p.LastFile, func(line string) error {
			for _, item := range gjson.Get(line, options.Field).Array() {
				jsonArrayEnum(item, func(result gjson.Result) {
					f.WriteString(result.Raw + "\n")
				})
			}
			return nil
		})
	})
}

func init() {
	piperunner.RegisterWorkflow("flat", flat, "FlatArray", flatArray)
}
