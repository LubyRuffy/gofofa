package web

import (
	"fmt"
	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
	"github.com/sirupsen/logrus"
	"net/http"
	"text/template"
)

import "embed"

//go:embed public
var webFs embed.FS

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(webFs, "public/index.html"))
	code := r.FormValue("code")
	mermaid, err := workflowast.NewParser().ParseToGraph(code)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("workflow parsed err: %v", err)))
		return
	}
	tmpl.Execute(w, mermaid)
}

// Start 启动服务器
func Start(addr string) error {
	// 默认首页
	http.HandleFunc("/", handler)

	// 静态资源
	http.Handle("/public/", http.StripPrefix("/",
		http.FileServer(http.FS(webFs))))

	logrus.Println("listen at: ", addr)
	return http.ListenAndServe(addr, nil)
}
