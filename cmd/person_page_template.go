package cmd

import (
	"github.com/tektsu/gedcom"
)

// personPageTemplate is the template used to build a person web page.
const personPageTemplate string = `---
title: "{{ .Name.Full }}{{ if or .Birth .Death }} ({{ .Birth }} - {{ .Death }}){{ end }}"
url: "/{{ .ID | ToLower }}/"
categories:
  - Person
{{ if .LastNames }}lastnames:
  {{ range .LastNames }}- {{ . }}{{ end }}
{{- end }}
{{ if .TopPhoto }}portrait: {{ .TopPhoto.File }}{{end}}
---

<div id="person">

<div id="page_title">
<table class="page_title_table">
<tr>
{{- if .TopPhoto }}
<th class="page_title"{{ if ne .TopPhoto.Width 0 }} style="width:{{ min .TopPhoto.Width 200 }}px"{{ end }}><img src="/images/photos/{{ .TopPhoto.File }}" class="portrait" /></th>
{{- end }}
<th class="page_title" style="width:auto;">{{ .Name.Full }}{{ if or .Birth .Death }}<br />({{ .Birth }} - {{ .Death }}){{ end }}</th>
</tr>
</table>
</div>

<div id="personal_info">
<table class="personal_info_table">
<tr><th colspan="2" class="table_header">Personal Information</th></tr>
<tr><th>Name</th><td class="sex_{{ .Sex }}">
{{- .Name.Full }}
{{- if .Name.SourcesInd }}<sup>{{ range .Name.SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}
{{- if .SourcesInd }}<sup>{{ range .SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}
</td></tr>
<tr><th class="attrib_heading">Sex</th><td class="sex_{{ .Sex }}">{{ .Sex }}</td></tr>
{{ if not .Living }}
{{ range .Attributes }}
<tr><th class="attrib_heading">{{ .Name }}</th><td>
{{ if .Date }}{{ .Date }}: {{ end }}
{{ if .Value }}{{ .Value }} {{ end }}
{{ if .Place }}{{ .Place }} {{ end }}
{{- if .SourcesInd }}<sup>{{ range .SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}
</td></tr>
{{ end }}
{{ end }}
</table>
</div>

{{ if not .Living }}
{{ $len := len .Events }}{{ if gt $len 0 }}
<div id="personal_events">
<table class="personal_event_table">
<tr><th colspan="2" class="table_header">Life Events</th></tr>
{{ range .Events }}
<tr><th class="attrib_heading">{{ .Name }}</th><td>
{{ if .Date }}{{ .Date }}: {{ end }}
{{ if .Type }}{{ .Type }} {{ end }}
{{ if .Place }}{{ .Place }} {{ end }}
{{- if .SourcesInd }}<sup>{{ range .SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}
</td></tr>
{{ end }}
</table>
</div>
{{ end }}
{{ end }}

{{ if .ParentsFamily }}
<div id="parents">
{{ range .ParentsFamily }}
<table class="parents_family">
<tr><th colspan="2" class="table_header">Parent's Family</th></tr>
<tr><th class="attrib_heading">Father</th><td class="sex_M">
	{{- if .Father -}}
		<a href="/{{ .Father.ID | ToLower }}/">{{ .Father.Name }}</a>
		{{- if .Father.SourcesInd -}}
			<sup>{{ range .Father.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
<tr><th class="attrib_heading">Mother</th><td class="sex_F">
	{{- if .Mother -}}
		<a href="/{{ .Mother.ID | ToLower }}/">{{ .Mother.Name }}</a>
		{{- if .Mother.SourcesInd -}}
			<sup>{{ range .Mother.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
{{ $events := len .Events }} {{ if gt $events 0 }}
{{ range $index, $event := .Events }}
<tr><th class="attrib_heading">{{ .Name  }}</th><td>
{{ if .Date }}{{ .Date }}: {{ end }}
{{ if .Type }}{{ .Type }} {{ end }}
{{ if .Place }}{{ .Place }} {{ end }}
{{- if .SourcesInd }}<sup>{{ range .SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}
</td></tr>
{{ end }}
{{ end }}
{{ $length := len .Children }} {{ if gt $length 0 }}
{{ range $index, $child := .Children }}
<tr>{{ if eq $index 0 }}<th class="attrib_heading" rowspan="{{ $length }}">Children</th>{{ end }}
<td class="sex_{{ $child.Sex }}">{{ if ne .ID $.ID }}<a href="/{{ $child.ID | ToLower }}/">{{ end }}{{ $child.Name }}{{ if ne .ID $.ID }}</a>{{ end }}
{{- if $child.SourcesInd -}}<sup>{{ range $child.SourcesInd }} [{{ . }}]{{ end }}</sup>{{- end -}}
</td></tr>
{{ end }}
{{ end }}
</table>
{{ end }}
</div>
{{ end }}

{{ if .Family }}
<div id="family">
{{ range .Family}}
<table class="family">
<tr><th colspan="2" class="table_header">Family</th></tr>
{{ if ne .Father.ID $.ID }}
<tr><th class="attrib_heading">Spouse</th><td class="sex_M">
	{{- if .Father -}}
		<a href="/{{ .Father.ID | ToLower }}/">{{ .Father.Name }}</a>
		{{- if .Father.SourcesInd -}}
			<sup>{{ range .Father.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
{{ end }}
{{ if ne .Mother.ID $.ID }}
<tr><th class="attrib_heading">Spouse</th><td class="sex_F">
	{{- if .Mother -}}
		<a href="/{{ .Mother.ID | ToLower }}/">{{ .Mother.Name }}</a>
		{{- if .Mother.SourcesInd -}}
			<sup>{{ range .Mother.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
{{ end }}
{{ $events := len .Events }} {{ if gt $events 0 }}
{{ range $index, $event := .Events }}
<tr><th class="attrib_heading">{{ .Name }}</th><td>
{{ if .Date }}{{ .Date }}: {{ end }}
{{ if .Type }}{{ .Type }} {{ end }}
{{ if .Place }}{{ .Place }} {{ end }}
{{- if .SourcesInd }}<sup>{{ range .SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}
</td></tr>
{{ end }}
{{ end }}
{{ $length := len .Children }} {{ if gt $length 0 }}
{{ range $index, $child := .Children }}
<tr>{{ if eq $index 0 }}<th class="attrib_heading" rowspan="{{ $length }}">Children</th>{{ end }}
<td class="sex_{{ $child.Sex }}"><a href="/{{ $child.ID | ToLower }}/">{{ $child.Name }}</a>
{{- if $child.SourcesInd -}}<sup>{{ range $child.SourcesInd }} [{{ . }}]{{ end }}</sup>{{- end -}}
</td></tr>
{{ end }}
{{ end }}
</table>
{{ end }}
</div>
{{ end }}

{{ if .Photos }}
<div id="photos">
<table class="photos_table">
<tr><th colspan="2" class="table_header">Photos</th></tr>
</table>
{{ "load-photoswipe" | shortcode }}
{{ "gallery" | shortcode }}
{{ range .Photos }}
{{ openShortcode }}figure link="/images/photos/{{ .File }}" caption="{{ .Title }}"{{ closeShortcode}}
{{ end }}
{{ "/gallery" | shortcode }}
<br />
</div>
{{ end }}

{{ if .Sources -}}
<div id="sources">
<table class="sources_table">
<tr><th colspan="2" class="table_header">Sources</th></tr>
{{ range $i, $s := .Sources -}}
<tr><th class="source_heading">{{ add $i 1 }}.</th><td><a href="/s{{ $s.RefNum }}">{{ $s.Ref }}</a>{{ if $s.Detail }}, {{ $s.Detail }}{{ end }}</td></tr>
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
	Living        bool
	Sources       []*sourceRef
	ParentsFamily []*personFamily
	Family        []*personFamily
	SourcesInd    []int
	Attributes    []*eventRef
	Events        []*eventRef
	Birth         string
	Death         string
	TopPhoto      *photoRef
	Photos        []*photoRef
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
		ID:     id,
		Sex:    person.Sex,
		Living: people[person.Xref].Living,
	}
	if data.Sex != "M" && data.Sex != "F" {
		data.Sex = "U"
	}

	// appendSources is the callback method sent to any function which might
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

	// Get top-level citations
	data.SourcesInd = appendSources(sourcesFromCitations(person.Citation))

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

	// Add in personal attributes
	for _, a := range person.Attribute {

		if a.Tag == "SSN" { // Skip social security number
			continue
		}
		event := newEventRef(a, appendSources)
		data.Attributes = append(data.Attributes, event)
	}

	// Add in personal events
	for _, a := range person.Event {

		event := newEventRef(a, appendSources)
		if event.Name == "Photo" {
			continue
		}
		if !data.Living {
			if event.Name == "Birth" {
				data.Birth = event.Date
			}
			if event.Name == "Death" {
				data.Death = event.Date
			}
		}
		data.Events = append(data.Events, event)
	}

	// Add in the person's parents.
	for _, fr := range person.Parents {
		if fr.Family != nil {
			family := newPersonFamily(fr, appendSources)
			data.ParentsFamily = append(data.ParentsFamily, family)
		}
	}

	// Add in the person's family.
	for _, fr := range person.Family {
		family := newPersonFamily(fr, appendSources)
		data.Family = append(data.Family, family)
	}

	// Add in photos
	for _, o := range person.Object {

		if o.File.Form != "jpg" && o.File.Form != "png" {
			continue
		}
		p := newPhotoRef(o, person)
		data.Photos = append(data.Photos, p)
	}

	// Get Top photo
	if person.Photo != nil {
		p := newPhotoRef(person.Photo, person)
		data.TopPhoto = p
	}

	return data
}
