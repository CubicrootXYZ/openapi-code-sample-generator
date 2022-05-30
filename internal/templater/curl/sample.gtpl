{{- if (or .SecurityParameters.Header .SecurityParameters.Cookie .SecurityParameters.Query .SecurityParameters.Path) -}}
TOKEN="my secure token"
{{ end -}}
{{ if .BasicAuth -}}
USERNAME="user"
PASSWORD="******"

{{ end -}}
curl {{ .URL }}{{ .Path }}{{ if (or .QueryParamsString .SecurityParameters.Query) }}?{{ end }}{{ .QueryParamsString }}
{{- if (and .QueryParamsString .SecurityParameters.Query) -}}
&
{{- end }}
{{- range .SecurityParameters.Query }}
{{- (urlencode .Name) }}=${TOKEN}&
{{- end}} -X {{ .HTTPVerb }}

{{- if .BasicAuth }} -u ${USERNAME}:${PASSWORD}
{{- end }} 

{{- range $header := .Parameters.Header }} -H "{{ (escape $header.Name) }}: {{ (escape $header.Value) }}"
{{- end}}

{{- range $header := .SecurityParameters.Header }} -H "{{ (escape $header.Name) }}: {{ (escape $header.Value) }}"
{{- end}}

{{- if .Formatting.Format}} -H "Content-Type: {{ .Formatting.Format}} {{ if .Formatting.FormData.OuterBoundary }}boundary={{ .Formatting.FormData.OuterBoundary }}{{ end }}"
{{- end}}

{{- range $cookie := .Parameters.Cookie }} -b "{{ (escape $cookie.Name) }}={{ (escape $cookie.Value) }}"
{{- end}}

{{- range $cookie := .SecurityParameters.Cookie }} -b "{{ (escape $cookie.Name) }}={{ (escape $cookie.Value) }}"
{{- end}}

{{- if .BodyString }} -d "{{ (escape .BodyString) }}"
{{- end}}