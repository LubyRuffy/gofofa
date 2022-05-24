package web

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

import "embed"

//go:embed public
var webFs embed.FS

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
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
