{{- if (or .SecurityParameters.Query .SecurityParameters.Header .SecurityParameters.Cookie) -}}
var token = "my secret token";
{{- end }}
{{- if .BasicAuth }}
var username = "username";
var password = "********";
{{- end }}
var url = "{{ .URL }}{{ .Path }}{{ if (or .QueryParamsString .SecurityParameters.Query) }}?{{ end }}{{ .QueryParamsString }}
{{- if (and .QueryParamsString .SecurityParameters.Query) -}}
&
{{- end }}
{{- range .SecurityParameters.Query }}
{{- (urlencode .Name) }}=" + token + "&
{{- end}}";

var request = new XMLHttpRequest();
request.open("{{ .HTTPVerb }}", url);
{{- range .Parameters.Header }}
request.setRequestHeader("{{ (escape .Name) }}", "{{ (escape .Value) }}");
{{- end -}}
{{- range .SecurityParameters.Header }}
request.setRequestHeader("{{ (escape .Name) }}", {{ (converttoken .Value) }});
{{- end}}
{{- if .BasicAuth }}
request.setRequestHeader("Authorization", "Basic " + btoa(username + ":" + password));
{{- end }}
{{- if .Formatting.Format }}
request.setRequestHeader("Content-Type", "{{ .Formatting.Format }}
{{- if (and .Formatting.FormData .Formatting.FormData.OuterBoundary) }} boundary={{ .Formatting.FormData.OuterBoundary}}
{{- end }}");
{{- end }}

request.send({{ if index .Additionals "jsBody" }}{{ index .Additionals "jsBody" }}{{ else }}"{{ (escape .BodyString) }}"{{ end }});
console.log(request.responseText);