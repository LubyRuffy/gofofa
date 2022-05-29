package gocodefuncs

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/tidwall/sjson"
	"os"
	"strings"
)

// RemoveField 移除字段
func RemoveField(p Runner, params map[string]interface{}) *FuncResult {
	if len(p.GetLastFile()) == 0 {
		panic(fmt.Errorf("removeField need input pipe"))
	}

	fields := strings.Split(params["fields"].(string), ",")

	var fn string
	var err error
	fn, err = utils.WriteTempFile(".json", func(f *os.File) error {
		return utils.EachLine(p.GetLastFile(), func(line string) error {
			var err error
			newLine := line
			for _, field := range fields {
				newLine, err = sjson.Delete(newLine, field)
				if err != nil {
					return err
				}
			}
			_, err = f.WriteString(newLine + "\n")
			if err != nil {
				return err
			}
			return nil
		})
	})
	if err != nil {
		panic(fmt.Errorf("removeField error: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
