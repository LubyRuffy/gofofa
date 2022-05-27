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
			var ts []*PipeTask
			for _, wf := range p.GetWorkflows() {
				if wf.Runner != p {
					break
				}
				ts = append(ts, wf)
			}
			return ts
		},
	}).Parse(`

{{ template "task.tmpl" (GetTasks .) }}

{{ define "task.tmpl" }}
{{ range . }}
<ul>
	<li>{{ .Name }} ({{ .Content }}) </li>

	{{ if gt (len .Outfile) 0 }}
	<li><a href="{{ .Outfile | safeURL }}" target="_blank">{{ .Outfile | toFileName }}</a></li>
	{{ end }}

	{{ if gt (len .Artifacts) 0 }}
		<li>
		generate files:
		{{ range .Artifacts }}
			<ul>
				<li><a href="{{ .FilePath | safeURL }}" target="_blank">
					{{ if HasPrefix .FileType "image/" }}
						<img src="{{ .FilePath | safeURL }}" height="80px">
					{{ else if eq .FileType "chart_html"}}
						show <iframe width="640" height="480" src="{{ .FilePath | safeURL }}" frameBorder="0"></iframe>
					{{ else }}
						{{ .FilePath | toFileName }}
					{{ end }}
				</a> | {{ .FileType }} | {{ .Memo }}</li>
			</ul>
		{{ end }}
		</li>
	{{ end }}

	<li>{{ .Cost }}</li>

	{{ range .Children }}
	<li> fork children:
		{{ template "task.tmpl" (GetTasks .) }}
	</li>
	{{ end }}
</ul>
{{ end }}
{{ end }}
`)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	err = t.Execute(&out, p)
	if err != nil {
		panic(err)
	}

	return out.String()
}
