Hi @{{ .ReviewerUsername }}, please review the following:
**{{ .Description }}**

{{ range $i, $mr := .MergeRequests -}}
- {{ $mr.RepoName }}|{{ $mr.TargetBranch }}: {{ $mr.Link }}
{{ end -}}
{{ if .Footer }}
{{ .Footer }}
{{ end }}
Thank you! :capoo-thanks:
