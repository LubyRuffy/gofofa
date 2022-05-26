package goworkflow

import (
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gammazero/workerpool"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/lubyruffy/gofofa/pkg/fzq"
	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Artifact 过程中生成的文件
type Artifact struct {
	FilePath string // 文件路径
	FileName string // 文件路径
	FileSize int    // 文件大小
	FileType string // 文件类型
	Memo     string // 备注，比如URL等
}

// funcResult 返回的结构
type funcResult struct {
	OutFile   string // 往后传递的文件
	Artifacts []*Artifact
}

//type innerFunction func(*PipeRunner, map[string]interface{}) (string, []string)

func removeField(p *PipeRunner, params map[string]interface{}) *funcResult {
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

	return &funcResult{
		OutFile: fn,
	}
}

type fetchFofaParams struct {
	Query  string
	Size   int
	Fields string
}

func fetchFofa(p *PipeRunner, params map[string]interface{}) *funcResult {
	var err error
	var options fetchFofaParams
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

	return &funcResult{
		OutFile: fn,
	}
}

type chartParams struct {
	Type  string
	Title string
}

// 每一个json行格式必须有value和count字段，对应name和value之，比如：{"value":"US","count":435}
func generateChart(p *PipeRunner, params map[string]interface{}) *funcResult {
	var err error
	var options chartParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	var keys []string
	barItems := make([]opts.BarData, 0)
	pieItems := make([]opts.PieData, 0)
	//lineItems := make([]opts.LineData, 0)

	err = utils.EachLine(p.GetLastFile(), func(line string) error {
		value := gjson.Get(line, "value")
		count := gjson.Get(line, "count")
		if !value.Exists() || !count.Exists() {
			return fmt.Errorf(`chart data is invalid: "value" and "count" field is needed`)
		}
		keys = append(keys, gjson.Get(line, "value").String())
		barItems = append(barItems, opts.BarData{Name: gjson.Get(line, "value").String(), Value: gjson.Get(line, "count").Int()})
		pieItems = append(pieItems, opts.PieData{Name: gjson.Get(line, "value").String(), Value: gjson.Get(line, "count").Int()})
		//lineItems = append(lineItems, opts.LineData{Name: gjson.Get(line, "value").String(), Value: gjson.Get(line, "count").Int()})
		return nil
	})
	if err != nil {
		panic(err)
	}

	var chartRender render.Renderer
	switch options.Type {
	case "bar":
		chart := charts.NewBar()
		chart.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{Title: options.Title, Left: "center"}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		)
		chart.AddSeries("data", barItems)
		chartRender = chart
	case "pie":
		chart := charts.NewPie()
		chart.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{Title: options.Title, Left: "center"}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		)
		chart.AddSeries("data", pieItems)
		chartRender = chart
	//case "line":
	//	chart := charts.NewLine()
	//	chart.SetGlobalOptions(
	//		charts.WithTitleOpts(opts.Title{Title: options.Title, Left: "center"}),
	//		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	//	)
	//	chart.AddSeries("data", lineItems)
	//	chartRender = chart
	default:
		panic("unknown chart type: [" + options.Type + "]")
	}

	f, err := utils.WriteTempFile(".html", func(f *os.File) error {
		return chartRender.Render(f)
	})

	if err != nil {
		panic(fmt.Errorf("generateChart error: %w", err))
	}

	return &funcResult{
		Artifacts: []*Artifact{{
			FilePath: f,
			FileName: filepath.Base(f),
			FileType: "",
		}},
	}
}

type zqQueryParams struct {
	Query string `json:"query"`
}

func zqQuery(p *PipeRunner, params map[string]interface{}) *funcResult {
	var fn string
	var err error
	var options zqQueryParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	fn, err = utils.WriteTempFile(".json", nil)
	if err != nil {
		panic(err)
	}

	err = fzq.ZqQuery(options.Query, p.GetLastFile(), fn)
	if err != nil {
		panic(fmt.Errorf("zqQuery error: %w", err))
	}

	return &funcResult{
		OutFile: fn,
	}
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

func addField(p *PipeRunner, params map[string]interface{}) *funcResult {

	var err error
	var options addFieldParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	if options.Value == nil && options.From == nil {
		panic(fmt.Errorf("addField failed: neithor value nor from"))
	}

	var addValue string
	var addRegex *regexp.Regexp

	if options.Value != nil {
		addValue = *options.Value
	}

	var newLines []string
	utils.EachLine(p.GetLastFile(), func(line string) error {
		var newLine string
		if options.Value != nil {
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

	var fn string
	fn, err = utils.WriteTempFile(".json", func(f *os.File) error {
		content := strings.Join(newLines, "\n")
		_, err := f.WriteString(content)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("addField error: %w", err))
	}

	return &funcResult{
		OutFile: fn,
	}
}

type loadFileParams struct {
	File string
}

func loadFile(p *PipeRunner, params map[string]interface{}) *funcResult {
	var err error
	var options loadFileParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("loadFile failed: %w", err))
	}

	if len(options.File) == 0 {
		panic(errors.New("load file cannot be empty"))
	}

	var path string
	//path, _ = os.Getwd()
	path, _ = filepath.Abs(options.File)

	if _, err = os.Stat(path); err != nil {
		panic(fmt.Errorf("loadFile failed: %w", err))
	}

	//return path, nil

	fn, err := utils.WriteTempFile(".json", func(f *os.File) error {
		var bytesRead []byte
		bytesRead, err = ioutil.ReadFile(options.File)
		if err != nil {
			panic(err)
		}
		_, err = f.Write(bytesRead)
		return err
	})

	if err != nil {
		panic(fmt.Errorf("loadFile error: %w", err))
	}

	return &funcResult{
		OutFile: fn,
	}
}

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

func flatArray(p *PipeRunner, params map[string]interface{}) *funcResult {
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

	return &funcResult{
		OutFile: fn,
	}
}

type screenshotParam struct {
	URLField string `json:"urlField"` // url的字段名称，默认是url
	Timeout  int    `json:"timeout"`  // 整个浏览器操作超时
	Quality  int    `json:"quality"`  // 截图质量
}

func screenshotURL(u string, options *screenshotParam) (string, int, error) {
	logrus.Debugf("screenshot url: %s", u)

	var err error
	// prepare the chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
	)

	allocCtx, bcancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer bcancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(logrus.Debugf))
	ctx, cancel = context.WithTimeout(ctx, time.Duration(options.Timeout)*time.Second)
	defer cancel()

	// run task list
	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(u),
		chromedp.FullScreenshot(&buf, 80),
	)
	if err != nil {
		return "", 0, fmt.Errorf("screenShot failed(%w): %s", err, u)
	}

	var fn string
	fn, err = utils.WriteTempFile(".png", func(f *os.File) error {
		_, err = f.Write(buf)
		return err
	})

	return fn, len(buf), err
}

// 截图
func screenShot(p *PipeRunner, params map[string]interface{}) *funcResult {
	var err error
	var options screenshotParam
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("screenShot failed: %w", err))
	}

	if options.Timeout == 0 {
		options.Timeout = 30
	}
	if options.URLField == "" {
		options.URLField = "url"
	}

	var artifacts []*Artifact
	var fn string

	wp := workerpool.New(5)
	err = utils.EachLine(p.GetLastFile(), func(line string) error {
		wp.Submit(func() {
			u := gjson.Get(line, options.URLField).String()
			if len(u) == 0 {
				return
			}
			if !strings.Contains(u, "://") {
				u = "http://" + u
			}
			var size int
			fn, size, err = screenshotURL(u, &options)
			if err != nil {
				logrus.Warnf("screenshotURL failed: %s", err)
				return
			}

			artifacts = append(artifacts, &Artifact{
				FilePath: fn,
				FileSize: size,
				FileType: "image/png",
				FileName: filepath.Base(fn),
				Memo:     u,
			})
		})

		return nil
	})
	if err != nil {
		panic(fmt.Errorf("screenShot error: %w", err))
	}
	wp.StopWait()

	return &funcResult{
		Artifacts: artifacts,
	}
}
