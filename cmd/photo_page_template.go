package cmd

// sourcePageTemplate is the tmplate used to generaate a source web page.
const photoPageTemplate string = `---
url: "/{{ .ID }}/"
categories:
  - Photo
lead_photo: {{ .File }}
photo_key: {{ .ID  }}
{{ $length := len .Persons }}{{ if gt $length 0 }}photo_people:
{{ range $index, $person := .Persons }}  - {{ $person.ID | ToLower }}
{{ end }}
{{- end }}
---
<div id="photopage">

<div id="photopage_title">
<p>{{ .Title }}</p>
</div>

{{ "lead_photo_center" | shortcode }}

{{ $length := len .Persons }} {{ if gt $length 0 }}
<div id="photopage_linklist">
<p class="photopage_linklist_header">People linked to this photo:</p>
{{ range $index, $person := .Persons }}
<p class="photopage_linklist_entry"><a href="/{{ $person.ID | ToLower }}/">{{ $person.Name }}</a></p>
{{ end }}
</div>
{{ end }}

</div>
`

type photoTmplData struct {
	ID      string
	File    string
	Title   string
	Persons []*personRef
}

func newPhotoTmplData(photo *photoRef) *photoTmplData {

	d := &photoTmplData{
		ID:    photo.ID,
		File:  photo.File,
		Title: photo.Title,
	}

	for _, person := range photo.Persons {
		d.Persons = append(d.Persons, person)
	}

	return d
}
