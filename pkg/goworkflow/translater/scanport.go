package translater

import (
	"bytes"
	"text/template"

	"github.com/lubyruffy/gofofa/pkg/goworkflow/workflowast"
)

// scanport 端口扫描：scanport(targets,[ports])
func scanPortHook(fi *workflowast.FuncInfo) string {
	tmpl, _ := template.New("scan_port").Parse(`ScanPort(GetRunner(), map[string]interface{}{
    "targets": "{{ .Targets }}",
    "ports": "{{ .Ports }}",
})`)

	targets := "127.0.0.1"
	if len(fi.Params) > 0 {
		targets = fi.Params[0].RawString()
	}
	ports := "22,80,443,1080,3389,8000,8080,8443"
	if len(fi.Params) > 1 {
		ports = fi.Params[1].RawString()
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, struct {
		Targets string
		Ports   string
	}{
		Targets: targets,
		Ports:   ports,
	})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	register("scan_port", scanPortHook) // grep匹配再新增字段
}
