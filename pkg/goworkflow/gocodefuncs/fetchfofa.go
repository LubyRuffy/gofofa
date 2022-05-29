package gocodefuncs

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"os"
	"strings"
)

// FetchFofaParams 获取fofa的参数
type FetchFofaParams struct {
	Query  string
	Size   int
	Fields string
}

// FetchFofa 从fofa获取数据
func FetchFofa(p Runner, params map[string]interface{}) *FuncResult {
	var err error
	var options FetchFofaParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("fetchFofa failed: %w", err))
	}

	if len(options.Query) == 0 {
		panic(fmt.Errorf("fofa query cannot be empty"))
	}
	if len(options.Fields) == 0 {
		panic(fmt.Errorf("fofa fields cannot be empty"))
	}

	fields := strings.Split(options.Fields, ",")

	var res [][]string
	res, err = p.GetFofaCli().HostSearch(options.Query, options.Size, fields)
	if err != nil {
		panic(fmt.Errorf("HostSearch failed: fofa fields cannot be empty"))
	}

	var fn string
	fn, err = utils.WriteTempFile(".json", func(f *os.File) error {
		w := outformats.NewJSONWriter(f, fields)
		return w.WriteAll(res)
	})
	if err != nil {
		panic(fmt.Errorf("fetchFofa error: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
