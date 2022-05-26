package goworkflow

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
)

// DumpTasks tasks dump to html
func (p *PipeRunner) DumpTasks(server bool) string {
	t, err := template.New("tasks").Funcs(template.FuncMap{
		"toFileName": func(u string) string {
			return filepath.Base(u)
		},
		"HasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"safeURL": func(u string) template.URL {
			if server {
				return template.URL("/file?url=" + filepath.Base(u))
			}
			u = strings.ReplaceAll(u, "\\", "/")
			return template.URL(u)
		},
		"GetTasks": func(p *PipeRunner) []*PipeTask {
			return p.GetWorkflows()
		},
	}).Parse(`

{{range .}}
	{{ template "task.tmpl" . }}
{{end}}

{{ define "task.tmpl" }}
<ul>
	<li>{{ .Name }}</li>
	<li>{{ .Content }}</li>

	{{ if gt (len .Outfile) 0 }}
	<li><a href="{{ .Outfile | safeURL }}" target="_blank">{{ .Outfile | toFileName }}</a></li>
	{{ end }}

	{{ range .Artifacts }}
generate files:
	<li> 
		<ul>
			<li><a href="{{ .FilePath | safeURL }}" target="_blank">
				{{ if HasPrefix .FileType "image/" }}
					<img src="{{ .FilePath | safeURL }}" height="80px">
				{{ else }}
					{{ .FilePath | toFileName }}
				{{ end }}
			</a> | {{ .FileType }} | {{ .Memo }}</li>
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
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	err = t.Execute(&out, p.Tasks)
	if err != nil {
		panic(err)
	}

	return out.String()
}
