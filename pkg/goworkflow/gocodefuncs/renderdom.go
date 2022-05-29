package gocodefuncs

import (
	"fmt"
	"os"
	"strings"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/gammazero/workerpool"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/net/context"
)

// renderURLDOM 生成单个url的domhtml
func renderURLDOM(p Runner, u string, timeout int) (string, error) {
	p.Debugf("render url dom: %s", u)

	var html string
	err := chromeActions(u, p.Debugf, timeout, chromedp.ActionFunc(func(ctx context.Context) error {
		node, err := dom.GetDocument().Do(ctx)
		if err != nil {
			return err
		}
		html, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
		return err
	}))
	if err != nil {
		return "", fmt.Errorf("renderURLDOM failed(%w): %s", err, u)
	}

	return html, err
}

// RenderDOM 动态渲染指定的URL，拼凑HTML
func RenderDOM(p Runner, params map[string]interface{}) *FuncResult {
	var err error
	var options screenshotParam
	if err = mapstructure.Decode(params, &options); err != nil {
		panic(fmt.Errorf("screenShot failed: %w", err))
	}

	if options.URLField == "" {
		options.URLField = "url"
	}
	if options.SaveField == "" {
		options.SaveField = "rendered_html"
	}
	if options.Timeout == 0 {
		options.Timeout = 30
	}

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

				var html string
				html, err = renderURLDOM(p, u, options.Timeout)

				// 不管是否成功都先把数据写入
				line, err = sjson.Set(line, options.SaveField, html)
				if err != nil {
					return
				}
				_, err = f.WriteString(line + "\n")
				if err != nil {
					return
				}
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
		OutFile: fn,
	}
}
