package funcs

import (
	"bytes"
	"errors"
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"github.com/lubyruffy/gofofa/pkg/piperunner"
	"github.com/tidwall/sjson"
	"os"
	"strings"
	"text/template"
)

func rmHook(fi *pipeparser.FuncInfo) string {
	tmpl, err := template.New("rm").Parse(`RemoveField(GetRunner(), map[string]interface{}{
   "fields": {{ . }},
})`)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, fi.Params[0].String())
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func removeField(p *piperunner.PipeRunner, params map[string]interface{}) string {
	if len(p.LastFile) == 0 {
		panic(errors.New("removeField need input pipe or file"))
	}

	fields := strings.Split(params["fields"].(string), ",")

	return piperunner.WriteTempJSONFile(func(f *os.File) {
		piperunner.EachLine(p.LastFile, func(line string) error {
			var err error
			newLine := line
			for _, field := range fields {
				newLine, err = sjson.Delete(newLine, field)
				if err != nil {
					panic(err)
				}
			}
			_, err = f.WriteString(newLine + "\n")
			if err != nil {
				panic(err)
			}
			return nil
		})
	})
}

func init() {
	piperunner.RegisterWorkflow("rm", rmHook, "RemoveField", removeField) // 删除字段
}
