package input

import (
	"bytes"
	"errors"
	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/funcs"
	"github.com/mitchellh/mapstructure"
	"os"
	"strings"
	"text/template"
)

type fetchFofaParams struct {
	Query  string
	Size   int
	Fields string
}

func fetchFofa(p corefuncs.Runner, params map[string]interface{}) (string, []string) {
	var err error
	var options fetchFofaParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	if len(options.Query) == 0 {
		panic(errors.New("fofa query cannot be empty"))
	}
	if len(options.Fields) == 0 {
		panic(errors.New("fofa fields cannot be empty"))
	}

	fields := strings.Split(options.Fields, ",")

	var res [][]string
	res, err = p.GetFofaCli().HostSearch(options.Query, options.Size, fields)
	if err != nil {
		panic(err)
	}

	return funcs.WriteTempFile(".json", func(f *os.File) {
		w := outformats.NewJSONWriter(f, fields)
		if err = w.WriteAll(res); err != nil {
			panic(err)
		}
	}), nil
}

func fofaHook(fi *pipeast.FuncInfo) string {
	tmpl, err := template.New("fofa").Parse(`FetchFofa(GetRunner(), map[string]interface{} {
    "query": {{ .Query }},
    "size": {{ .Size }},
    "fields": {{ .Fields }},
})`)
	if err != nil {
		panic(err)
	}
	var size int64 = 10
	fields := "`host,title`"
	if len(fi.Params) > 1 {
		fields = fi.Params[1].String()
	}
	if len(fi.Params) > 2 {
		size = fi.Params[2].Int64()
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		Query  string
		Size   int64
		Fields string
	}{
		Query:  fi.Params[0].String(),
		Fields: fields,
		Size:   size,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	corefuncs.RegisterWorkflow("fofa", fofaHook, "FetchFofa", fetchFofa)
}
