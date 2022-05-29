package gocodefuncs

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"os"
)

// GenData 生成数据
func GenData(p Runner, params map[string]interface{}) *FuncResult {
	var fn string
	var err error
	fn, err = utils.WriteTempFile("", func(f *os.File) error {
		_, err = f.WriteString(params["data"].(string))
		return err
	})
	if err != nil {
		panic(fmt.Errorf("genData failed: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
