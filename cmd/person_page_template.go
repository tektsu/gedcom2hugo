package cmd

import (
	"github.com/tektsu/gedcom"
)

const personPageTemplate string = `---
title: "{{ .Name.Full }}"
url: "/{{ .ID | ToLower }}/"
categories:
  - Person
{{ if .LastNames }}lastnames:
  {{ range .LastNames }}- {{ . }}{{ end }}
{{- end }}
---

<div id="person">

<div id="personal_info">
<table class="personal_info_table">
<tr><th>Name</th><td>{{ .Name.Full }}{{ if .Name.SourcesInd }}<sup>{{ range .Name.SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}</td></tr>
<tr><th>Sex</th><td>{{ .Sex }}</td></tr>
</table>
</div>

{{ if .ParentsFamily }}
<div id="parents">
{{ range .ParentsFamily}}
<table class="parents_family">
<tr><th colspan="2">Parent's Family</th></tr>
<tr><th>Father</th><td>
	{{- if .Father.ID -}}
		<a href="/{{ .Father.ID | ToLower }}/">{{ .Father.Name }}</a>
		{{- if .Father.SourcesInd -}}
			<sup>{{ range .Father.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
<tr><th>Mother</th><td>
	{{- if .Mother.ID -}}
		<a href="/{{ .Mother.ID | ToLower }}/">{{ .Mother.Name }}</a>
		{{- if .Mother.SourcesInd -}}
			<sup>{{ range .Mother.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
{{ $length := len .Children }} {{ if gt $length 0 }}
<tr><th>Siblings</th><td>
{{ range .Children }}
	<a href="/{{ .ID | ToLower }}/">{{ .Name }}</a>
	{{- if .SourcesInd -}}
		<sup>{{ range .SourcesInd }} [{{ . }}]{{ end }}</sup>
	{{- end -}} ({{ .Sex }})<br />
{{ end }}
</td></tr>
{{ end }}
</table>
{{ end }}
</div>
{{ end }}

{{ if .Sources -}}
<div id="sources">
<table class="sources_table">
<tr><th colspan="2">Sources</th></tr>
{{ range $i, $s := .Sources -}}
<tr><td>{{ add $i 1 }}.</td><td><a href="/s{{ $s.RefNum }}">{{ $s.Ref }}</a></td></tr>
{{ end -}}
</table>
</div>
{{- end }}

</div>
`

type personTmplData struct {
	ID            string
	Name          *personName
	Aliases       []*personName
	LastNames     []string
	Sex           string
	Sources       []*sourceRef
	ParentsFamily []*personFamily
}

func newPersonTmplData(person *gedcom.IndividualRecord) *personTmplData {

	count := 0 // Citation Counter

	id := person.Xref
	data := &personTmplData{
		ID:  id,
		Sex: person.Sex,
	}

	for i, n := range person.Name {
		var s []*sourceRef
		var name *personName

		lastNames := make(map[string]bool)
		count, s, name = newPersonNameWithCitations(count, n)
		for _, source := range s {
			data.Sources = append(data.Sources, source)
		}
		lastNames[name.Last] = true

		if i == 0 {
			data.Name = name
		} else {
			data.Aliases = append(data.Aliases, name)
		}
		data.LastNames = make([]string, 0, len(lastNames))
		for l := range lastNames {
			data.LastNames = append(data.LastNames, l)
		}
	}

	for _, fr := range person.Parents {
		if fr.Family != nil {
			var sources []*sourceRef
			var f *personFamily

			count, sources, f = newPersonFamily(count, fr)
			for _, s := range sources {
				data.Sources = append(data.Sources, s)
			}
			data.ParentsFamily = append(data.ParentsFamily, f)
		}
	}

	return data
}
