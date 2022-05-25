package goworkflow

import (
	"bytes"
	"html/template"
	"strings"
)

// DumpTasks tasks dump to html
func (p *PipeRunner) DumpTasks() string {
	t, _ := template.New("tasks").Funcs(template.FuncMap{
		"safeURL": func(u string) template.URL {
			u = strings.ReplaceAll(u, "\\", "/")
			return template.URL(u)
		},
		"GetTasks": func(p *PipeRunner) []*PipeTask {
			return p.GetWorkflows()
		},
	}).Parse(`   
<html>
<head>
    <title>gofofa tasks</title>
</head>
<body>
	<h1>gofofa tasks</h1>
	{{range .}}
		{{ template "task.tmpl" . }}
	{{end}}
</body>
</html>

{{ define "task.tmpl" }}
<ul>
	<li>{{ .Name }}</li>
	<li>{{ .Content }}</li>

	{{ if gt (len .Outfile) 0 }}
	<li><a href="{{ .Outfile | safeURL }}">{{ .Outfile }}</a></li>
	{{ end }}

	{{ range .Artifacts }}
	<li> generate files:
		<ul>
			<li><a href="{{ .FilePath | safeURL }}">{{ .FilePath }}</a> | {{ .FileType }} | {{ .Memo }}</li>
		</ul>
	</li>
	{{ end }}
	<li>{{ .Cost }}</li>

	{{ range .Children }}
	<li> fork children:
		{{ range . | GetTasks }}
			{{ template "task.tmpl" . }}
		{{ end }}
	</li>
	{{ end }}
</ul>
{{ end }}
`)
	var out bytes.Buffer
	err := t.Execute(&out, p.Tasks)
	if err != nil {
		panic(err)
	}

	return out.String()
}
