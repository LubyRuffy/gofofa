package gocodefuncs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gammazero/workerpool"
	"github.com/lubyruffy/gofofa/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/net/context"
)

type screenshotParam struct {
	URLField  string `json:"urlField"`  // url的字段名称，默认是url
	Timeout   int    `json:"timeout"`   // 整个浏览器操作超时
	Quality   int    `json:"quality"`   // 截图质量
	SaveField string `json:"saveField"` // 保存截图地址的字段
}

//chromeActions 完成chrome的headless操作
func chromeActions(u string, logf func(string, ...interface{}), timeout int, actions ...chromedp.Action) error {
	var err error
	// prepare the chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("incognito", true), // 隐身模式
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.IgnoreCertErrors,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.WindowSize(1024, 768),
	)

	allocCtx, bcancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer bcancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(logf))
	ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	realActions := []chromedp.Action{
		chromedp.Navigate(u),
	}
	realActions = append(realActions, actions...)
	// run task list
	err = chromedp.Run(ctx,
		realActions...,
	)

	return err
}

func screenshotURL(p Runner, u string, options *screenshotParam) (string, int, error) {
	p.Debugf("screenshot url: %s", u)

	var buf []byte
	err := chromeActions(u, p.Debugf, options.Timeout, chromedp.CaptureScreenshot(&buf))
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
