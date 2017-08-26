package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tektsu/gedcom"
)

const sourcePageTemplate string = `---
url: "/{{ .ID }}/"
categories:
  - Source
{{ if .Type }}type: "{{ .Type }}"{{ end }}
{{ if .Author }}author: "{{ .Author }}"{{ end }}
title: "Source: {{ if .Title }}{{ .Title }}{{ end }}"
{{ if .Abbr }}shorttitle: "{{ .Abbr }}"{{ end }}
{{ if .Publication }}pubfacts: "{{ .Publication }}"{{ end }}
{{ if .Periodical }}periodical: "{{ .Periodical }}"{{ end }}
{{ if .Volume }}volume: "{{ .Volume }}"{{ end }}
{{ if .MediaType }}mediatype: "{{ .MediaType }}"{{ end }}
{{ if .Repository }}repository:
{{ range .Repository }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .Submitter }}submitter:
{{ range .Submitter }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .Page }}page:
{{ range .Page }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .Film }}film:
{{ range .Film }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .File }}file:
{{ range .File }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .FileNumber }}filenumber:
{{ range .FileNumber }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .Date }}docdate:
{{ range .Date }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .DateViewed }}dateviewed:
{{ range .DateViewed }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .Place }}place:
{{ range .Place }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .URL }}docurl:
{{ range .URL }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .DocLocation }}doclocation:
{{ range .DocLocation }}  - "{{ . }}"
{{ end }}{{ end }}
{{ if .RefNum }}refnum: "{{ .RefNum }}"{{ end }}
{{ if .Ref }}ref: |
  {{ .Ref }}{{ end }}
---
{{ "sourcebody" | shortcode }}

{{ if .Text }}
### Notes

{{ .Text }}
{{ end }}
`

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
}

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
	if len(s.File) > 0 {
		d.File = make([]string, len(s.File))
		copy(d.File, s.File)
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

	// Build the reference string.
	var refs []string
	if d.Author != "" {
		refs = append(refs, d.Author)
	}
	if d.Title != "" {
		refs = append(refs, fmt.Sprintf("\"%s\"", d.Title))
	}
	if len(d.Date) > 0 {
		refs = append(refs, d.Date[0])
	}
	if d.Periodical != "" {
		refs = append(refs, fmt.Sprintf("%s", d.Periodical))
	}
	if d.Volume != "" {
		refs = append(refs, fmt.Sprintf("%s", d.Volume))
	}
	if len(d.Page) > 0 {
		refs = append(refs, fmt.Sprintf("p %s", d.Page[0]))
	}
	if len(d.Film) > 0 {
		refs = append(refs, d.Film[0])
	}
	if len(d.Place) > 0 {
		refs = append(refs, d.Place[0])
	}
	if d.Type != "" {
		refs = append(refs, d.Type)
	}
	if len(d.URL) > 0 {
		refs = append(refs, fmt.Sprintf("[%s]", d.URL[0]))
	}
	if len(d.FileNumber) > 0 {
		refs = append(refs, d.FileNumber[0])
	}
	d.Ref = strings.Join(refs, ", ")

	return d
}
