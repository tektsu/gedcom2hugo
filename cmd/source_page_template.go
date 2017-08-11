package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
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
`

func newSourceData(cx *cli.Context, source *gedcom.SourceRecord) (sourceData, error) {

	data := sourceData{
		ID:          source.Xref,
		Author:      source.Author,
		Title:       source.Title,
		Abbr:        source.Abbr,
		Publication: source.Publication,
		Text:        source.Text,
		Type:        source.Type,
		File:        source.File,
		FileNumber:  source.FileNumber,
		Place:       source.Place,
		Date:        source.Date,
		DateViewed:  source.DateViewed,
		URL:         source.URL,
		DocLocation: source.DocLocation,
		Periodical:  source.Periodical,
		Volume:      source.Volume,
		MediaType:   source.MediaType,
	}

	// Build the reference string.
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(source.Xref, 1)
	data.Ref = fmt.Sprintf("%s. ", matches[0])
	v, err := strconv.Atoi(matches[0])
	if err != nil {
		panic(fmt.Sprintf("Error converting [%s] to integer", matches[0]))
	}
	data.RefNum = v

	// Copy in the arrays
	if len(source.File) > 0 {
		data.File = make([]string, len(source.File))
		copy(data.File, source.File)
	}
	if len(source.FileNumber) > 0 {
		data.FileNumber = make([]string, len(source.FileNumber))
		copy(data.FileNumber, source.FileNumber)
	}
	if len(source.Place) > 0 {
		data.Place = make([]string, len(source.Place))
		copy(data.Place, source.Place)
	}
	if len(source.Date) > 0 {
		data.Date = make([]string, len(source.Date))
		copy(data.Date, source.Date)
	}
	if len(source.DateViewed) > 0 {
		data.DateViewed = make([]string, len(source.DateViewed))
		copy(data.DateViewed, source.DateViewed)
	}
	if len(source.URL) > 0 {
		data.URL = make([]string, len(source.URL))
		copy(data.URL, source.URL)
	}
	if len(source.DocLocation) > 0 {
		data.DocLocation = make([]string, len(source.DocLocation))
		copy(data.DocLocation, source.DocLocation)
	}
	if len(source.Repository) > 0 {
		data.Repository = make([]string, len(source.Repository))
		copy(data.Repository, source.Repository)
	}
	if len(source.Submitter) > 0 {
		data.Submitter = make([]string, len(source.Submitter))
		copy(data.Submitter, source.Submitter)
	}
	if len(source.Page) > 0 {
		data.Page = make([]string, len(source.Page))
		copy(data.Page, source.Page)
	}
	if len(source.Film) > 0 {
		data.Film = make([]string, len(source.Film))
		copy(data.Film, source.Film)
	}

	var refs []string
	if data.Author != "" {
		refs = append(refs, data.Author)
	}
	if data.Title != "" {
		refs = append(refs, fmt.Sprintf("\"%s\"", data.Title))
	}
	if len(data.Date) > 0 {
		refs = append(refs, data.Date[0])
	}
	if data.Periodical != "" {
		refs = append(refs, fmt.Sprintf("%s", data.Periodical))
	}
	if data.Volume != "" {
		refs = append(refs, fmt.Sprintf("%s", data.Volume))
	}
	if len(data.Page) > 0 {
		refs = append(refs, fmt.Sprintf("p %s", data.Page[0]))
	}
	if len(data.Film) > 0 {
		refs = append(refs, data.Film[0])
	}
	if len(data.Place) > 0 {
		refs = append(refs, data.Place[0])
	}
	if data.Type != "" {
		refs = append(refs, data.Type)
	}
	if len(data.URL) > 0 {
		refs = append(refs, fmt.Sprintf("[%s]", data.URL[0]))
	}
	if len(data.FileNumber) > 0 {
		refs = append(refs, data.FileNumber[0])
	}
	data.Ref = strings.Join(refs, ", ")

	return data, nil
}
