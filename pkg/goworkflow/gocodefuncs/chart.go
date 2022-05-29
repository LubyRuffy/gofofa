package gocodefuncs

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
)

type chartParams struct {
	Type  string
	Title string
}

// GenerateChart 生成报表。每一个json行格式必须有value和count字段，对应name和value之，比如：{"value":"US","count":435}
func GenerateChart(p Runner, params map[string]interface{}) *FuncResult {
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
