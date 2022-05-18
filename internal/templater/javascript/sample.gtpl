{{- if (or .SecurityParameters.Query .SecurityParameters.Header .SecurityParameters.Cookie) }}
var token = "my secret token";
{{ end -}}
{{ if .BasicAuth }}
var username = "username";
var password = "********";
{{ end -}}

var url = "{{ .URL }}{{ .Path }}{{ if (or .QueryParamsString .SecurityParameters.Query) }}?{{ end }}{{ .QueryParamsString }}
{{- if (and .QueryParamsString .SecurityParameters.Query) -}}
&
{{- end }}
{{- range .SecurityParameters.Query }}
{{- (urlencode .Name) }}=" + token + "&
{{- end}}";

{{- if (or .Parameters.Cookie .SecurityParameters.Cookie) }}
$cookies = "
{{- range .Parameters.Cookie -}}
    {{ (escape .Name) }}={{ (escape .Value) }};
{{- end -}}
{{- range .SecurityParameters.Cookie -}}
    {{ (escape .Name) }}=" . {{ (escape (converttoken .Value)) }} . ";
{{- end -}}";
{{- end }}

{{- if (or .BodyString (index .Additionals "customRequestBody")) }}
$data = {{ if index .Additionals "customRequestBody" -}}
{{ index .Additionals "customRequestBody" -}};
{{ else if .BodyString -}}
"{{ (escape .BodyString) }}";
{{- end -}}
{{- end }}

var request = new XMLHttpRequest();
request.open(url);
{{- range .Parameters.Header }}
request.setRequestHeader("{{ (escape .Name) }}: {{ (escape .Value) }}")
{{- end -}}
{{- range .SecurityParameters.Header }}
request.setRequestHeader("{{ (escape .Name) }}: " + {{ (converttoken .Value) }}}})
{{- end}}



curl_setopt($curl, CURLOPT_CUSTOMREQUEST, "{{ .HTTPVerb }}");
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
{{ if (or .BodyString (index .Additionals "customRequestBody")) -}}
curl_setopt($curl, CURLOPT_POSTFIELDS, $data);
{{ end }}
{{- if (or .Parameters.Header .SecurityParameters.Header) -}}
curl_setopt($curl, CURLOPT_HTTPHEADER, $headers);
{{ end }}
{{- if (or .Parameters.Cookie .SecurityParameters.Cookie) -}}
curl_setopt($curl, CURLOPT_COOKIE, $cookies);
{{ end }}
{{- if .BasicAuth -}}
curl_setopt($ch, CURLOPT_USERPWD, $username . \":\" . $password);
{{ end -}}
var_dump(curl_exec($curl)); // Dumps the response
curl_close($curl);