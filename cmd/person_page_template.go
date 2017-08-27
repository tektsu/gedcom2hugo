package cmd

import (
	"github.com/tektsu/gedcom"
)

// personPageTemplate is the template used to build a person web page.
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

// personTmplData is the structure that is sent to the personPageTemplate for display.
type personTmplData struct {
	ID            string
	Name          *personName
	Aliases       []*personName
	LastNames     []string
	Sex           string
	Sources       []*sourceRef
	ParentsFamily []*personFamily
}

// sourceCB is the type of the callback function passed to various methods to
// handle source references.
type sourceCB func(s []*sourceRef) []int

// newPersonTmplData builds a personTmplData structure from a
// gedcom.IndividualRecord.
func newPersonTmplData(person *gedcom.IndividualRecord) *personTmplData {

	count := 0 // Local citation counter

	id := person.Xref
	data := &personTmplData{
		ID:  id,
		Sex: person.Sex,
	}
	if data.Sex != "M" && data.Sex != "F" {
		data.Sex = "U"
	}

	// appendSources is the callback method send to any function which might
	// produce sources. It accumulates any sources in the personTmplData
	// structure and returns to the caller a list of local source references.
	appendSources := func(s []*sourceRef) []int {
		var localRefs []int

		for _, source := range s {
			data.Sources = append(data.Sources, source)
			count++
			localRefs = append(localRefs, count)
		}

		return localRefs
	}

	// Add in the person's names.
	for i, n := range person.Name {

		lastNames := make(map[string]bool)
		name := newPersonNameWithCitations(n, appendSources)
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

	// Add in the person's parents.
	for _, fr := range person.Parents {
		if fr.Family != nil {
			family := newPersonFamily(fr, appendSources)
			data.ParentsFamily = append(data.ParentsFamily, family)
		}
	}

	return data
}
