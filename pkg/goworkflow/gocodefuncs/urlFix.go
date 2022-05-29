package gocodefuncs

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
	"strings"
)

// UrlFix 自动补齐url
func UrlFix(p Runner, params map[string]interface{}) *FuncResult {
	var fn string
	var err error
	field := "url"
	if len(params) > 0 {
		field = params["urlField"].(string)
	}
	if len(field) == 0 {
		panic(fmt.Errorf("urlFix must has a field"))
	}

	fn, err = utils.WriteTempFile("", func(f *os.File) error {
		return utils.EachLine(p.GetLastFile(), func(line string) error {
			v := gjson.Get(line, field).String()
			if !strings.Contains(v, "://") {
				v = "http://" + gjson.Get(line, field).String()
			}
			line, err := sjson.Set(line, field, v)
			if err != nil {
				return err
			}
			_, err = f.WriteString(line + "\n")
			return err
		})
	})
	if err != nil {
		panic(fmt.Errorf("urlFix failed: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
