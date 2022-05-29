package gocodefuncs

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/tidwall/gjson"
	"github.com/xuri/excelize/v2"
)

// ToExcel 写excel文件
func ToExcel(p Runner, params map[string]interface{}) *FuncResult {
	var fn string
	var err error

	fn, err = utils.WriteTempFile(".xlsx", nil)
	if err != nil {
		panic(fmt.Errorf("toExcel failed: %w", err))
	}

	f := excelize.NewFile()
	defer f.Close()

	lineNo := 2
	err = utils.EachLine(p.GetLastFile(), func(line string) error {
		v := gjson.Parse(line)
		colNo := 'A'
		v.ForEach(func(key, value gjson.Result) bool {
			// 设置第一行
			if lineNo == 2 {
				err = f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", colNo, lineNo-1), key.Value())
			}

			// 写值
			err = f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", colNo, lineNo), value.Value())
			colNo++
			if err != nil {
				panic(fmt.Errorf("SetCellValue failed: %w", err))
			}
			return true
		})
		lineNo++
		return err
	})
	if err != nil {
		panic(fmt.Errorf("toExcel failed: %w", err))
	}

	err = f.SaveAs(fn)
	if err != nil {
		panic(fmt.Errorf("toExcel failed: %w", err))
	}

	return &FuncResult{
		Artifacts: []*Artifact{
			{
				FilePath: fn,
				FileType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			},
		},
	}
}
