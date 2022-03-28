**What I have done this week**
{{ range $jira, $updates := .UpdatesMap -}}
- [{{ $jira.Title }}]({{ $jira.Link }})
{{range $updates }}  - {{.}}
{{ end }}{{ end }}
**What I will do next working week**
{{ range $jira, $updates := .UpdatesMap -}}
- [{{ $jira.Title }}]({{ $jira.Link }})
{{ end }}