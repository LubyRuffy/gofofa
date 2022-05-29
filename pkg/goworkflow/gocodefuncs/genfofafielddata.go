package gocodefuncs

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/sjson"
	"os"
	"strings"
)

// GenFofaFieldData 根据字段生成随机数据
func GenFofaFieldData(p Runner, params map[string]interface{}) *FuncResult {
	var err error
	var options FetchFofaParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("fetchFofa failed: %w", err))
	}

	var res string
	fields := strings.Split(options.Fields, ",")
	for i := 0; i < options.Size; i++ {
		v := `{}`
		for _, f := range fields {
			var value interface{}
			switch f {
			case "host":
				value = gofakeit.FirstName() + "." + gofakeit.DomainName()
			case "domain":
				value = gofakeit.DomainName()
			case "ip":
				value = gofakeit.IPv4Address()
			case "port":
				value = gofakeit.IntRange(21, 65534)
			case "country":
				value = gofakeit.CountryAbr()
			}
			v, _ = sjson.Set(v, f, value)
		}
		res += v + "\n"
	}

	var fn string
	fn, err = utils.WriteTempFile(".json", func(f *os.File) error {
		_, err = f.WriteString(res)
		return err
	})
	if err != nil {
		panic(fmt.Errorf("fetchFofa error: %w", err))
	}

	return &FuncResult{
		OutFile: fn,
	}
}
