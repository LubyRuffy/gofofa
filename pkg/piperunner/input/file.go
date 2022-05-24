package input

import (
	"bytes"
	"errors"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/mitchellh/mapstructure"
	"os"
	"path/filepath"
	"text/template"
)

type loadFileParams struct {
	File string
}

func loadFile(p corefuncs.Runner, params map[string]interface{}) (string, []string) {
	var err error
	var options loadFileParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	if len(options.File) == 0 {
		panic(errors.New("load file cannot be empty"))
	}

	var path string
	//path, _ = os.Getwd()
	path, _ = filepath.Abs(options.File)

	if _, err = os.Stat(path); err != nil {
		panic(err)
	}

	return path, nil

	//fn := funcs.WriteTempFile(".json", func(f *os.File) {
	//	var bytesRead []byte
	//	bytesRead, err = ioutil.ReadFile(options.File)
	//	if err != nil {
	//		panic(err)
	//	}
	//	_, err = f.Write(bytesRead)
	//	if err != nil {
	//		panic(err)
	//	}
	//})
	//
	//return fn, nil
}

func loadHook(fi *pipeast.FuncInfo) string {
	tmpl, err := template.New("load").Parse(`LoadFile(GetRunner(), map[string]interface{} {
    "file": {{ .File }},
})`)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		File string
	}{
		File: fi.Params[0].String(),
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	corefuncs.RegisterWorkflow("load", loadHook, "LoadFile", loadFile)
}
