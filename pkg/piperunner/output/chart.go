package output

import (
	"bytes"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/lubyruffy/gofofa/pkg/pipeast"
	"github.com/lubyruffy/gofofa/pkg/piperunner/corefuncs"
	"github.com/lubyruffy/gofofa/pkg/piperunner/funcs"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"os"
	"text/template"
)

// chart 生成报表
// chart("pie")
// 第一个参数是报表类型
func chart(fi *pipeast.FuncInfo) string {
	tmpl, _ := template.New("chart").Parse(`GenerateChart(GetRunner(), map[string]interface{}{
    "type": {{ .Type }},
    "title": "{{ .Title }}",
})`)

	typeStr := "bar"
	if len(fi.Params) > 0 {
		typeStr = fi.Params[0].String()
	}
	title := ""
	if len(fi.Params) > 1 {
		title = fi.Params[1].RawString()
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, struct {
		Type  string
		Title string
	}{
		Type:  typeStr,
		Title: title,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

type chartParams struct {
	Type  string
	Title string
}

// 每一个json行格式必须有value和count字段，对应name和value之，比如：{"value":"US","count":435}
func generateChart(p corefuncs.Runner, params map[string]interface{}) (string, []string) {
	var err error
	var options chartParams
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(err)
	}

	var keys []string
	barItems := make([]opts.BarData, 0)
	pieItems := make([]opts.PieData, 0)
	//lineItems := make([]opts.LineData, 0)

	funcs.EachLine(p.GetLastFile(), func(line string) error {
		keys = append(keys, gjson.Get(line, "value").String())
		barItems = append(barItems, opts.BarData{Name: gjson.Get(line, "value").String(), Value: gjson.Get(line, "count").Int()})
		pieItems = append(pieItems, opts.PieData{Name: gjson.Get(line, "value").String(), Value: gjson.Get(line, "count").Int()})
		//lineItems = append(lineItems, opts.LineData{Name: gjson.Get(line, "value").String(), Value: gjson.Get(line, "count").Int()})
		return nil
	})

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

	f := funcs.WriteTempFile(".html", func(f *os.File) {
		if err = chartRender.Render(f); err != nil {
			panic(err)
		}
	})

	return "", []string{f}
}

func init() {
	corefuncs.RegisterWorkflow("chart", chart, "GenerateChart", generateChart)
}
