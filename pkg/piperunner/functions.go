package piperunner

import (
	"errors"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/fzq"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
	"regexp"
	"strings"
)

type addFieldFrom struct {
	Method  string `json:"method"`
	Field   string
	Value   string
	Options string
}

type addFieldParams struct {
	Name  string
	Value *string       // 可以没有，就取from
	From  *addFieldFrom // 可以没有，就取Value
}

func addField(p *PipeRunner, params map[string]interface{}) {
	logrus.Debug("addField params:", params)
	if len(p.LastFile) == 0 {
		panic(errors.New("addField need input pipe or file"))
	}

	var err error
	var options addFieldParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	var addValue string
	var addRegex *regexp.Regexp

	var newLines []string
	EachLine(p.LastFile, func(line string) error {
		var newLine string
		if options.Value != nil {
			if addValue == "" {
				addValue = *options.Value
			}
			newLine, _ = sjson.Set(line, options.Name, addValue)
		} else {
			switch options.From.Method {
			case "grep":
				if addRegex == nil {
					addRegex, err = regexp.Compile(options.From.Value)
					if err != nil {
						panic(err)
					}
				}
				res := addRegex.FindAllStringSubmatch(gjson.Get(line, options.From.Field).String(), -1)
				newLine, err = sjson.Set(line, options.Name, res)
				if err != nil {
					panic(err)
				}
			default:
				panic(errors.New("unknown from type"))
			}
		}
		newLines = append(newLines, newLine)
		return nil
	})

	pt := pipeTask{
		name:    "addField",
		content: fmt.Sprintf("%v", params),
		outfile: WriteTempJSONFile(func(f *os.File) {
			content := strings.Join(newLines, "\n")
			n, err := f.WriteString(content)
			if err != nil {
				panic(err)
			}
			if n != len(content) {
				panic("write string failed")
			}
			p.LastFileSize = int64(n)
		}),
	}
	p.addPipe(pt)
}

func removeField(p *PipeRunner, params map[string]interface{}) {
	logrus.Debug("removeField params:", params)
	if len(p.LastFile) == 0 {
		panic(errors.New("removeField need input pipe or file"))
	}

	fields := strings.Split(params["fields"].(string), ",")

	var newLines []string
	EachLine(p.LastFile, func(line string) error {
		var err error
		newLine := line
		for _, field := range fields {
			newLine, err = sjson.Delete(newLine, field)
			if err != nil {
				panic(err)
			}
		}
		newLines = append(newLines, newLine)
		return nil
	})

	pt := pipeTask{
		name:    "removeField",
		content: fmt.Sprintf("%v", params),
		outfile: WriteTempJSONFile(func(f *os.File) {
			content := strings.Join(newLines, "\n")
			n, err := f.WriteString(content)
			if err != nil {
				panic(err)
			}
			if n != len(content) {
				panic("write string failed")
			}
			p.LastFileSize = int64(n)
		}),
	}
	p.addPipe(pt)
}

type zqQueryParams struct {
	Query string `json:"query"`
}

func zqQuery(p *PipeRunner, params map[string]interface{}) {
	logrus.Debug("zqQuery params:", params)
	var err error
	var options zqQueryParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.CreateTemp(os.TempDir(), defaultPipeTmpFilePrefix)
	if err != nil {
		panic(fmt.Errorf("create tmpfile failed: %w", err))
	}
	name := f.Name()
	f.Close()

	err = fzq.ZqQuery(options.Query, p.LastFile, name)
	if err != nil {
		panic(err)
	}
	pt := pipeTask{
		name:    "zqQuery",
		content: fmt.Sprintf("%v", params),
		outfile: name,
	}
	p.addPipe(pt)
}
