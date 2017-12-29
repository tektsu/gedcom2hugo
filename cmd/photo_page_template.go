package cmd

// sourcePageTemplate is the tmplate used to generaate a source web page.
const photoPageTemplate string = `---
url: "/{{ .ID }}/"
categories:
  - Photo
---
<figure>
	<img src="/images/photos/{{ .File }}" alt="{{ .Title }}" />
	<figcaption><strong>{{ .Title }}</strong>
{{ $length := len .Persons }} {{ if gt $length 0 }}
<br /><br />People linked to this photo:<br />
{{ range $index, $person := .Persons }}
<a href="/{{ $person.ID | ToLower }}/">{{ $person.Name }}</a><br />
{{ end }}
{{ end }}
	</figcaption>
</figure>
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
