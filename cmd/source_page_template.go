package cmd

import (
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tektsu/gedcom"
)

// sourcePageTemplate is the tmplate used to generaate a source web page.
const sourcePageTemplate string = `---
url: "/{{ .ID }}/"
categories:
  - Source
title: "Source: {{ if .Title }}{{ .Title }}{{ end }}"
{{ if .RefNum }}refnum: "{{ .RefNum }}"{{ end }}
---
<table class="page_title_table">
<tr><th class="page_title">{{ .Title }}</th></tr>
</table>

<div class="sourceref">
<p class="sourceref"><strong>Citation:</strong> {{ .Ref }}</p>
</div>

<table id="source">
<tr><th>Field</th><th>Data</th></tr>
{{ if .Type }}<tr><td>Type</td><td>{{ .Type }}</td></tr>{{ end }}
{{ if .Periodical }}<tr><td>Periodical or Book Name</td><td>{{ .Periodical }}</td></tr>{{ end }}
{{ if .Author }}<tr><td>Author</td><td>{{ .Author }}</td></tr>{{ end }}
{{ if .Title }}<tr><td>Title</td><td>{{ .Title }}</td></tr>{{ end }}
{{ if .Abbr }}<tr><td>Short Title</td><td>{{ .Abbr }}</td></tr>{{ end }}
{{ if .Publication }}<tr><td>Publication Facts</td><td>{{ .Publication }}</td></tr>{{ end }}
{{ if .Volume }}<tr><td>Volume</td><td>{{ .Volume }}</td></tr>{{ end }}
{{ if .Page }}<tr><td>Page</td><td>{{ range .Page }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .Film }}<tr><td>Film Reference</td><td>{{ range .Film }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .Date }}<tr><td>Document Date</td><td>{{ range .Date }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .Place }}<tr><td>Place</td><td>{{ range .Place }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .Repository }}<tr><td>Repository</td><td>{{ range .Repository }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .FileNumber }}<tr><td>File Number</td><td>{{ range .FileNumber }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .File }}<tr><td>Files</td><td>{{ range .File }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .URL }}<tr><td>URL</td><td>{{ range .URL }}<a href="{{ . }}" target="_blank"></a>{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .DocLocation }}<tr><td>Document Location</td><td>{{ range .DocLocation }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .DateViewed }}<tr><td>Date Viewed</td><td>{{ range .DateViewed }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .MediaType }}<tr><td>Media Type</td><td>{{ .MediaType }}</td></tr>{{ end }}
{{ if .Submitter }}<tr><td>Submitter</td><td>{{ range .Submitter }}{{ . }}<br />{{ end }}</td></tr>{{ end }}
{{ if .Notes }}<tr><td>Notes</td><td>{{ .Notes }}</td></tr>{{ end }}
</table>

{{ if .Text }}
### Notes

{{ .Text }}
{{ end }}
`

// sourceTmplData is the data sent to the sourcePageTemplate for display.
type sourceTmplData struct {
	ID          string
	Author      string
	Abbr        string
	Publication string
	Text        string
	Title       string
	Type        string
	File        []string
	FileNumber  []string
	Place       []string
	Date        []string
	DateViewed  []string
	URL         []string
	DocLocation []string
	RefNum      int
	Ref         string
	Periodical  string
	Volume      string
	MediaType   string
	Repository  []string
	Submitter   []string
	Page        []string
	Film        []string
	Notes       string
}

// newSourceTmplData builds a sourceTmplData structure from a
// gedcom.SourceRecord.
func newSourceTmplData(s *gedcom.SourceRecord) *sourceTmplData {

	d := &sourceTmplData{
		ID:          s.Xref,
		Author:      s.Author,
		Title:       s.Title,
		Abbr:        s.Abbr,
		Publication: s.Publication,
		Text:        s.Text,
		Type:        s.Type,
		File:        s.File,
		FileNumber:  s.FileNumber,
		Place:       s.Place,
		Date:        s.Date,
		DateViewed:  s.DateViewed,
		URL:         s.URL,
		DocLocation: s.DocLocation,
		Periodical:  s.Periodical,
		Volume:      s.Volume,
		MediaType:   s.MediaType,
	}

	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(s.Xref, 1)
	d.Ref = fmt.Sprintf("%s. ", matches[0])
	v, err := strconv.Atoi(matches[0])
	if err != nil {
		panic(fmt.Sprintf("Error converting [%s] to integer", matches[0]))
	}
	d.RefNum = v

	// Copy in the arrays
	if len(s.Object) > 0 {
		r := regexp.MustCompile("^.*Roots/")
		m, _ := regexp.Compile("^/")
		for _, o := range s.Object {
			name := r.ReplaceAllString(o.File.Name, "")
			if m.MatchString(name) {
				name = filepath.Base(name)
			}
			d.File = append(d.File, name)
		}
	}
	if len(s.FileNumber) > 0 {
		d.FileNumber = make([]string, len(s.FileNumber))
		copy(d.FileNumber, s.FileNumber)
	}
	if len(s.Place) > 0 {
		d.Place = make([]string, len(s.Place))
		copy(d.Place, s.Place)
	}
	if len(s.Date) > 0 {
		d.Date = make([]string, len(s.Date))
		copy(d.Date, s.Date)
	}
	if len(s.DateViewed) > 0 {
		d.DateViewed = make([]string, len(s.DateViewed))
		copy(d.DateViewed, s.DateViewed)
	}
	if len(s.URL) > 0 {
		d.URL = make([]string, len(s.URL))
		copy(d.URL, s.URL)
	}
	if len(s.DocLocation) > 0 {
		d.DocLocation = make([]string, len(s.DocLocation))
		copy(d.DocLocation, s.DocLocation)
	}
	if len(s.Repository) > 0 {
		d.Repository = make([]string, len(s.Repository))
		copy(d.Repository, s.Repository)
	}
	if len(s.Submitter) > 0 {
		d.Submitter = make([]string, len(s.Submitter))
		copy(d.Submitter, s.Submitter)
	}
	if len(s.Page) > 0 {
		d.Page = make([]string, len(s.Page))
		copy(d.Page, s.Page)
	}
	if len(s.Film) > 0 {
		d.Film = make([]string, len(s.Film))
		copy(d.Film, s.Film)
	}
	if len(s.Note) > 0 {
		for _, n := range s.Note {
			d.Notes += n.Note + "\n\n"
		}
	}
	esc := strings.TrimRight(template.HTMLEscapeString(d.Notes), "\n")
	d.Notes = strings.Replace(esc, "\n", "<br>", -1)

	// Build the reference string.
	d.Ref = s.GetReferenceString()

	return d
}
