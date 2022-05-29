package gocodefuncs

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gammazero/workerpool"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type screenshotParam struct {
	URLField  string `json:"urlField"`  // url的字段名称，默认是url
	Timeout   int    `json:"timeout"`   // 整个浏览器操作超时
	Quality   int    `json:"quality"`   // 截图质量
	SaveField string `json:"saveField"` // 保存截图地址的字段
}

func screenshotURL(p Runner, u string, options *screenshotParam) (string, int, error) {
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

// ScreenShot 截图
func ScreenShot(p Runner, params map[string]interface{}) *FuncResult {
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
