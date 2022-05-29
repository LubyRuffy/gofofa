package goworkflow

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gammazero/workerpool"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lubyruffy/gofofa/pkg/fzq"
	"github.com/lubyruffy/gofofa/pkg/outformats"
	"github.com/lubyruffy/gofofa/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// FuncResult 返回的结构
type FuncResult struct {
	OutFile   string // 往后传递的文件
	Artifacts []*Artifact
}

//type innerFunction func(*PipeRunner, map[string]interface{}) (string, []string)

func removeField(p *PipeRunner, params map[string]interface{}) *FuncResult {
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

type fetchFofaParams struct {
	Query  string
	Size   int
	Fields string
}

func fetchFofa(p *PipeRunner, params map[string]interface{}) *FuncResult {
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

	return &FuncResult{
		OutFile: fn,
	}
}

type chartParams struct {
	Type  string
	Title string
}

// 每一个json行格式必须有value和count字段，对应name和value之，比如：{"value":"US","count":435}
func generateChart(p *PipeRunner, params map[string]interface{}) *FuncResult {
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
			charts.WithInitializationOpts(opts.Initialization{AssetsHost: "/public/assets/libs/echarts/"}),
		)
		chart.AddSeries("data", barItems)
		chartRender = chart
	case "pie":
		chart := charts.NewPie()
		chart.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{Title: options.Title, Left: "center"}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
			charts.WithInitializationOpts(opts.Initialization{AssetsHost: "/public/assets/libs/echarts/"}),
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
		panic(fmt.Errorf("unknown chart type: [" + options.Type + "]"))
	}

	f, err := utils.WriteTempFile(".html", func(f *os.File) error {
		return chartRender.Render(f)
	})

	if err != nil {
		panic(fmt.Errorf("generateChart error: %w", err))
	}

	return &FuncResult{
		Artifacts: []*Artifact{{
			FilePath: f,
			FileName: filepath.Base(f),
			FileType: "chart_html",
		}},
	}
}

type zqQueryParams struct {
	Query string `json:"query"`
}

func zqQuery(p *PipeRunner, params map[string]interface{}) *FuncResult {
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

	return &FuncResult{
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

func addField(p *PipeRunner, params map[string]interface{}) *FuncResult {

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

	return &FuncResult{
		OutFile: fn,
	}
}

type loadFileParams struct {
	File string
}

func loadFile(p *PipeRunner, params map[string]interface{}) *FuncResult {
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

	return &FuncResult{
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

func flatArray(p *PipeRunner, params map[string]interface{}) *FuncResult {
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

type screenshotParam struct {
	URLField  string `json:"urlField"`  // url的字段名称，默认是url
	Timeout   int    `json:"timeout"`   // 整个浏览器操作超时
	Quality   int    `json:"quality"`   // 截图质量
	SaveField string `json:"saveField"` // 保存截图地址的字段
}

func screenshotURL(p *PipeRunner, u string, options *screenshotParam) (string, int, error) {
	p.Debugf("screenshot url: %s", u)

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

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(p.Debugf))
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
func screenShot(p *PipeRunner, params map[string]interface{}) *FuncResult {
	var err error
	var options screenshotParam
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("screenShot failed: %w", err))
	}

	if options.URLField == "" {
		options.URLField = "url"
	}
	if options.SaveField == "" {
		options.SaveField = "screenshot_filepath"
	}
	if options.Timeout == 0 {
		options.Timeout = 30
	}

	var artifacts []*Artifact

	wp := workerpool.New(5)
	var fn string
	fn, err = utils.WriteTempFile(".json", func(f *os.File) error {
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
				var sfn string
				sfn, size, err = screenshotURL(p, u, &options)
				if err != nil {
					p.Warnf("screenshotURL failed: %s", err)
					f.WriteString(line + "\n")
					return
				}

				// 不管是否成功都先把数据写入
				line, err = sjson.Set(line, options.SaveField, sfn)
				if err != nil {
					return
				}
				_, err = f.WriteString(line + "\n")
				if err != nil {
					return
				}

				artifacts = append(artifacts, &Artifact{
					FilePath: sfn,
					FileSize: size,
					FileType: "image/png",
					FileName: filepath.Base(fn),
					Memo:     u,
				})
			})

			return nil
		})
		if err != nil {
			return err
		}
		wp.StopWait()
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("screenShot error: %w", err))
	}

	return &FuncResult{
		OutFile:   fn,
		Artifacts: artifacts,
	}
}

// 写excel文件
func toExcel(p *PipeRunner, params map[string]interface{}) *FuncResult {
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

type sqlParam struct {
	Driver string `json:"driver"` // 连接字符串: db_user:password@tcp(localhost:3306)/my_db
	DSN    string `json:"dsn"`    // 连接字符串: db_user:password@tcp(localhost:3306)/my_db
	Table  string `json:"table"`  // 表名
	Fields string `json:"fields"` // 写入的列名
}

func sqliteDSNToFilePath(dsn string) string {
	fqs := strings.SplitN(dsn, "?", 2)
	return fqs[0]
}

// 写入sql数据库
func toSql(p *PipeRunner, params map[string]interface{}) *FuncResult {
	var err error
	var db *sql.DB
	var options sqlParam
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("toSql failed: %w", err))
	}

	// 打开数据库
	if len(options.DSN) > 0 {
		switch options.Driver {
		case "sqlite3":
			// 文件名进行替换，只能写在临时目录吗？？？
			fqs := strings.SplitN(options.DSN, "?", 2)
			dir := filepath.Dir(fqs[0])
			if dir != os.TempDir() {
				if len(fqs) > 1 {
					options.DSN = filepath.Join(os.TempDir(), filepath.Base(fqs[0])) + "?" + fqs[1]
				} else {
					options.DSN = filepath.Join(os.TempDir(), filepath.Base(options.DSN))
				}
			}
		}
		db, err = sql.Open(options.Driver, options.DSN)
		if err != nil {
			panic(fmt.Errorf("toSql failed: %w", err))
		}
	} else {
		switch options.Driver {
		case "sqlite3":
			fn, err := utils.WriteTempFile(".sqlite3", nil)
			if err != nil {
				panic(fmt.Errorf("toSql failed: %w", err))
			}
			options.DSN = fn
			db, err = sql.Open(options.Driver, options.DSN)
			if err != nil {
				panic(fmt.Errorf("toSql failed: %w", err))
			}
		}
	}

	// 获取数据的列
	line, err := utils.ReadFirstLineOfFile(p.LastFile)
	if err != nil {
		panic(fmt.Errorf("ReadFirstLineOfFile failed: %w", err))
	}
	fieldsWithType := utils.JSONLineFieldsWithType(string(line))
	if len(fieldsWithType) == 0 {
		return &FuncResult{}
	}

	var tableNotExist bool
	var columns []string
	if len(options.Fields) == 0 {
		if db == nil {
			// 没有配置db，从文件读取
			for _, key := range fieldsWithType {
				columns = append(columns, key[0])
			}
		} else {
			// 自动从数据库获取一次
			var rows *sql.Rows
			rows, err = db.Query(fmt.Sprintf("select * from %s limit 1", options.Table))
			if err != nil {
				// 表格不存在的错误提示
				// sqlite3: no such table: tbl
				// mysql: 1146 table doesn’t exists
				if !strings.Contains(err.Error(), "no such table") &&
					!strings.Contains(err.Error(), "table doesn’t exist") {
					panic(fmt.Errorf("toSql failed: %w", err))
				}
				tableNotExist = true
			} else {
				var cols []string
				cols, err = rows.Columns()
				if err != nil {
					panic(fmt.Errorf("toSql failed: %w", err))
				}

				for _, col := range cols {
					for _, field := range fieldsWithType {
						if strings.ToLower(col) == strings.ToLower(field[0]) {
							columns = append(columns, field[0])
						}
					}
				}
			}
		}
	} else {
		columns = strings.Split(options.Fields, ",")
	}

	// 创建表结构
	if db != nil {
		// 还没有取到列，可能是表不存在
		if columns == nil {
			for _, f := range fieldsWithType {
				columns = append(columns, f[0])
			}
		}
		// 创建表结构
		var sqlColumnDesc []string
		for _, f := range fieldsWithType {
			needField := false
			if tableNotExist {
				needField = true
			}
			for _, col := range columns {
				// 两边都有，才创建
				if col == f[0] {
					needField = true
					break
				}
			}
			if needField {
				sqlColumnDesc = append(sqlColumnDesc, fmt.Sprintf("%s %s", f[0], f[1]))
			}
		}
		if db != nil {
			sqlString := fmt.Sprintf("create table if not exists %s (%s);", options.Table, strings.Join(sqlColumnDesc, ","))
			_, err = db.Exec(sqlString)
			if err != nil {
				panic(fmt.Errorf("create table failed: %w", err))
			}
		}
	}

	if len(columns) == 0 {
		panic(fmt.Errorf("toSql failed: no columns matched"))
	}
	var columnsString = strings.Join(columns, ",")

	var fn string
	fn, err = utils.WriteTempFile(".sql", func(f *os.File) error {
		err = utils.EachLine(p.GetLastFile(), func(line string) error {
			var valueString string
			for _, field := range columns {
				var vs string
				v := gjson.Get(line, field).Value()
				switch t := v.(type) {
				case string:
					vs = `"` + utils.EscapeString(v.(string)) + `"`
				case bool:
					vs = strconv.FormatBool(v.(bool))
				case float64:
					vs = strconv.FormatInt(int64(v.(float64)), 10)
				default:
					return fmt.Errorf("toSql failed: unknown data type %v", t)
				}
				if len(valueString) > 0 {
					valueString += ","
				}
				valueString += vs
			}

			sqlLine := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)\n",
				options.Table, columnsString, valueString)
			_, err = f.WriteString(sqlLine)
			if err != nil {
				return err
			}

			if db != nil {
				_, err = db.Exec(sqlLine)
				if err != nil {
					return err
				}
			}
			return err
		})

		return err
	})

	if err != nil {
		panic(fmt.Errorf("toSql failed: %w", err))
	}

	artifacts := []*Artifact{
		{
			FilePath: fn,
			FileType: "text/sql",
		},
	}
	switch options.Driver {
	case "sqlite3":
		fn = sqliteDSNToFilePath(options.DSN)
		artifacts = append(artifacts, &Artifact{
			FilePath: fn,
			FileType: "application/vnd.sqlite3",
		})
	}

	return &FuncResult{
		Artifacts: artifacts,
	}
}

func genData(p *PipeRunner, params map[string]interface{}) *FuncResult {
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

// 自动补齐url
func urlFix(p *PipeRunner, params map[string]interface{}) *FuncResult {
	var fn string
	var err error
	field := "url"
	if len(params) > 0 {
		field = params["url"].(string)
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
