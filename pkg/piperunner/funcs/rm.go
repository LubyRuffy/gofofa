package funcs

import (
	"bytes"
	"errors"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/tidwall/sjson"
	"os"
	"strings"
	"text/template"
)

func rmHook(fi *pipeast.FuncInfo) string {
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

func removeField(p corefuncs.Runner, params map[string]interface{}) (string, []string) {
	if len(p.GetLastFile()) == 0 {
		panic(errors.New("removeField need input pipe or file"))
	}

	fields := strings.Split(params["fields"].(string), ",")

	return WriteTempFile(".json", func(f *os.File) {
		EachLine(p.GetLastFile(), func(line string) error {
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
	}), nil
}

func init() {
	corefuncs.RegisterWorkflow("rm", rmHook, "RemoveField", removeField) // 删除字段
}
