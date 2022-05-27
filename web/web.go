package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/lubyruffy/gofofa"
	"github.com/lubyruffy/gofofa/pkg/goworkflow"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/sirupsen/logrus"
)

//go:embed public
var webFs embed.FS
var FofaCli *gofofa.Client

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(webFs, "public/index.html"))
	tmpl.Execute(w, "")
}

func parse(w http.ResponseWriter, r *http.Request) {
	// fofa(`title=test`) & to_int(`port`) & sort(`port`) & [cut(`port`) | cut("ip")]
	w.Header().Set("Content-Type", "application/json")

	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  true,
			"result": fmt.Sprintf("workflow parsed err: %v", err),
		})
		return
	}

	// 输入源
	sourceWorkflow := []string{
		"load", "fofa",
	}
	// 终止
	finishWorkflow := []string{
		"chart", "to_excel",
	}

	ast := workflowast.NewParser()
	realCode, err := ast.Parse(string(code))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  true,
			"result": fmt.Sprintf("workflow parsed err: %v", err),
		})
		return
	}
	var calls []string
	for _, fi := range ast.CallList {
		calls = append(calls, fi.Name)
	}

	graphCode, err := ast.ParseToGraph(string(code), func(name, s string) string {
		for _, src := range sourceWorkflow {
			if src == name {
				return `[("` + s + `")]`
			}
		}
		for _, src := range finishWorkflow {
			if src == name {
				return `[["` + s + `"]]`
			}
		}
		return `["` + s + `"]`
	}, "graph LR\n")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  true,
			"result": fmt.Sprintf("workflow parsed err: %v", err),
		})
		return
	}
	logrus.Println(graphCode)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": false,
		"result": map[string]interface{}{
			"realCode":  realCode,
			"graphCode": graphCode,
			"calls":     calls,
		},
	})
}

func run(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	workflow, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  true,
			"result": fmt.Sprintf("workflow parsed err: %v", err),
		})
		return
	}

	tm := globalTaskMonitor.new()
	go func() {
		var code string
		ast := workflowast.NewParser()
		code, err = ast.Parse(string(workflow))
		if err != nil {
			tm.addMsg("run err: " + err.Error())
		}

		p := goworkflow.New(goworkflow.WithHooks(&goworkflow.Hooks{
			OnWorkflowFinished: func(pt *goworkflow.PipeTask) {
				tm.addMsg(fmt.Sprintf("workflow finished: %s, %d", pt.Name, pt.CallID))
			},
			OnLog: func(level logrus.Level, format string, args ...interface{}) {
				tm.addMsg(fmt.Sprintf("[%s] %s", level.String(), fmt.Sprintf(format, args...)))
			},
		}))
		p.FofaCli = FofaCli
		_, err = p.Run(code)
		if err != nil {
			tm.addMsg("run err: " + err.Error())
		}

		tm.html = p.DumpTasks(true)
		tm.addMsg("<finished>")
		tm.finish()
	}()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  false,
		"result": tm.taskId,
	})
}

func fetchMsg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tid := r.FormValue("tid")

	t, ok := globalTaskMonitor.m.Load(tid)
	task := t.(*taskInfo)
	if !ok {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  true,
			"result": fmt.Sprintf("no task found"),
		})
		return
	}
	var msgs []string
	s := time.Now()
	for {
		log.Println(time.Since(s))
		info, ok := task.receiveMsg()
		if !ok {
			break
		}
		msgs = append(msgs, info)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": false,
		"result": map[string]interface{}{
			"msgs": msgs,
			"html": task.html,
		},
	})
}

// Start 启动服务器
func Start(addr string) error {
	// 默认首页
	http.HandleFunc("/", handler)
	http.HandleFunc("/parse", parse)
	http.HandleFunc("/run", run)
	http.HandleFunc("/fetchMsg", fetchMsg)
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		fn := filepath.Base(r.FormValue("url"))
		f := filepath.Join(os.TempDir(), fn)
		switch filepath.Ext(fn) {
		case ".sql", ".xlsx":
			w.Header().Set("Content-Disposition", "attachment; filename="+fn)
		}
		if len(r.FormValue("dl")) > 0 {
			w.Header().Set("Content-Disposition", "attachment; filename="+fn)
		}
		http.ServeFile(w, r, f)
	})

	// 静态资源
	http.Handle("/public/", http.StripPrefix("/",
		http.FileServer(http.FS(webFs))))

	logrus.Println("listen at: ", addr)
	return http.ListenAndServe(addr, nil)
}
