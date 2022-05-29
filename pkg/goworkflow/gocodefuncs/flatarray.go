package gocodefuncs

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"os"
)

type flatParams struct {
	Field string
}

func jsonArrayEnum(node gjson.Result, f func(result gjson.Result) error) error {
	if node.IsArray() {
		for _, child := range node.Array() {
			err := jsonArrayEnum(child, f)
			if err != nil {
				return err
			}
		}
	} else {
		return f(node)
	}
	return nil
}

// FlatArray 打平一个Array数据内容
func FlatArray(p Runner, params map[string]interface{}) *FuncResult {
	var err error
	var options flatParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	if len(options.Field) == 0 {
		panic(fmt.Errorf("flatArray: field cannot be empty"))
	}

	var fn string
	fn, err = utils.WriteTempFile(".json", func(f *os.File) error {
		return utils.EachLine(p.GetLastFile(), func(line string) error {
			for _, item := range gjson.Get(line, options.Field).Array() {
				err = jsonArrayEnum(item, func(result gjson.Result) error {
					_, err := f.WriteString(result.Raw + "\n")
					return err
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
	})
	if err != nil {
		panic(fmt.Errorf("flatArray error: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
