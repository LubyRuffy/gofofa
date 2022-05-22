package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/pipeparser"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/urfave/cli/v2"
)

var (
	pipelineFile string
)

type pipeTask struct {
	content string // raw content
	outfile string // tmp file
}

// Close remove tmp outfile
func (p *pipeTask) Close() {
	os.Remove(p.outfile)
}

// PipeRunner pipe运行器
type PipeRunner struct {
	content      string
	tasks        []pipeTask
	lastFile     string
	lastFileSize int64 // 最后写入文件的大小
}

// NewPipeRunner create pipe runner
func NewPipeRunner(content string) *PipeRunner {
	return &PipeRunner{
		content: content,
	}
}

// Close remove tmp outfile
func (p *PipeRunner) Close() {
	for _, task := range p.tasks {
		task.Close()
	}
}

func (p *PipeRunner) addPipe(pt pipeTask) {
	p.tasks = append(p.tasks, pt)
	p.lastFile = pt.outfile
}

func read(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

func (p *PipeRunner) eachLine(filename string, size int64, f func(line string) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// check valid
	if size > 0 {
		var s os.FileInfo
		s, err = file.Stat()
		if err != nil {
			panic(err)
		}
		if s.Size() != size {
			panic("size is not same")
		}
	}

	reader := bufio.NewReader(file)
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		f(string(line))
	}
	return nil
}

func (p *PipeRunner) writeTempJSONFile(writeF func(f *os.File)) string {
	var f *os.File
	var err error
	f, err = os.CreateTemp(os.TempDir(), "gofofa_pipeline_")
	if err != nil {
		panic(fmt.Errorf("create outFile %s failed: %w", outFile, err))
	}
	defer f.Close()

	writeF(f)

	return f.Name()
}

type fetchFofaParams struct {
	Query  string
	Size   int
	Fields string
}

func (p *PipeRunner) fetchFofa(params map[string]interface{}) {
	logrus.Debug("fetchFofa params:", params)

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
	res, err = fofaCli.HostSearch(options.Query, options.Size, fields)
	if err != nil {
		panic(err)
	}

	pt := pipeTask{
		content: fmt.Sprintf("%v", params),
		outfile: p.writeTempJSONFile(func(f *os.File) {
			w := outformats.NewJSONWriter(f, fields)
			if err = w.WriteAll(res); err != nil {
				panic(err)
			}
		}),
	}
	p.addPipe(pt)
	logrus.Debug("write to file:", pt.outfile)
}

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

func (p *PipeRunner) addField(params map[string]interface{}) {
	logrus.Debug("addField params:", params)
	if len(p.lastFile) == 0 {
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
	p.eachLine(p.lastFile, p.lastFileSize, func(line string) error {
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
		content: fmt.Sprintf("%v", params),
		outfile: p.writeTempJSONFile(func(f *os.File) {
			content := strings.Join(newLines, "\n")
			n, err := f.WriteString(content)
			if err != nil {
				panic(err)
			}
			if n != len(content) {
				panic("write string failed")
			}
			p.lastFileSize = int64(n)
		}),
	}
	p.addPipe(pt)
	logrus.Debug("write to file:", pt.outfile, " size:", p.lastFileSize)
}

func (p *PipeRunner) removeField(params map[string]interface{}) {
	logrus.Debug("removeField params:", params)
	if len(p.lastFile) == 0 {
		panic(errors.New("removeField need input pipe or file"))
	}

	fields := strings.Split(params["fields"].(string), ",")

	var newLines []string
	p.eachLine(p.lastFile, p.lastFileSize, func(line string) error {
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
		content: fmt.Sprintf("%v", params),
		outfile: p.writeTempJSONFile(func(f *os.File) {
			content := strings.Join(newLines, "\n")
			n, err := f.WriteString(content)
			if err != nil {
				panic(err)
			}
			if n != len(content) {
				panic("write string failed")
			}
			p.lastFileSize = int64(n)
		}),
	}
	p.addPipe(pt)
	logrus.Debug("write to file:", pt.outfile, " size:", p.lastFileSize)
}

// Run run pipelines
func (p *PipeRunner) Run() error {
	var err error

	i := interp.New(interp.Options{})
	_ = i.Use(stdlib.Symbols)

	err = i.Use(interp.Exports{
		"this/this": {
			"FetchFofa":   reflect.ValueOf(p.fetchFofa),
			"RemoveField": reflect.ValueOf(p.removeField),
			"AddField":    reflect.ValueOf(p.addField),
		},
	})
	if err != nil {
		panic(err)
	}

	// i.ImportUsed()
	i.Eval(`import (
		. "this/this"
		)`)
	_, err = i.Eval(p.content)

	return err
}

// pipeline subcommand
var pipelineCmd = &cli.Command{
	Name:                   "pipeline",
	Usage:                  "fofa data pipeline",
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Usage:       "load pipeline file",
			Destination: &pipelineFile,
		},
	},
	Action: pipelineAction,
}

// pipelineAction pipeline action
func pipelineAction(ctx *cli.Context) error {
	// valid same config
	var pipelineContent string
	if len(pipelineFile) > 0 {
		v, err := os.ReadFile(pipelineFile)
		if err != nil {
			return err
		}
		pipelineContent = string(v)
	}
	if v := ctx.Args().First(); len(v) > 0 {
		if len(pipelineContent) > 0 {
			return errors.New("file and content only one is allowed")
		}

		pipelineContent = pipeparser.NewParser().Parse(v)
	}

	pr := NewPipeRunner(pipelineContent)
	err := pr.Run()
	if err != nil {
		return err
	}

	pr.eachLine(pr.lastFile, pr.lastFileSize, func(line string) error {
		fmt.Println(line)
		return nil
	})

	return nil
}

func fofaHook(fi *pipeparser.FuncInfo) string {
	tmpl, err := template.New("fofa").Parse(`FetchFofa(map[string]interface{} {
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

func cutHook(fi *pipeparser.FuncInfo) string {
	tmpl, err := template.New("cut").Parse(`RemoveField(map[string]interface{}{
    "fields": {{ . }},
})`)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, fi.Params[0].String())
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func grepAddHook(fi *pipeparser.FuncInfo) string {
	tmpl, err := template.New("grep_add").Parse(`AddField(map[string]interface{}{
    "from": map[string]interface{}{
        "method": "grep",
        "field": {{ .Field }},
        "value": {{ .Value }},
    },
    "name": {{ .Name }},
})`)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, struct {
		Field string
		Value string
		Name  string
	}{
		Field: fi.Params[0].String(),
		Value: fi.Params[1].String(),
		Name:  fi.Params[2].String(),
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	pipeparser.RegisterFunction("fofa", fofaHook)
	pipeparser.RegisterFunction("cut", cutHook)
	pipeparser.RegisterFunction("grep_add", grepAddHook)
}
