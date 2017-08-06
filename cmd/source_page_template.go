package cmd

import (
	"fmt"
	"regexp"
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
{{ if .Date }}date: "{{ .Date }}"{{ end }}
{{ if .Place }}place: "{{ .Place }}"{{ end }}
{{ if .File }}file: "{{ .File }}"{{ end }}
{{ if .FileNumber }}filenumber: "{{ .FileNumber }}"{{ end }}
{{ if .URL }}docurl: "{{ .URL }}"{{ end }}
{{ if .DocLocation }}doclocation: "{{ .DocLocation }}"{{ end }}
{{ if .DateViewed }}dateviewed: "{{ .DateViewed }}"{{ end }}
---
<p class="sourceref">{{ .Ref }}</p>
<table id="source">
<tr><th>Field</th><th>Data</th></tr>
{{ if .Type }}<tr><td>Type</td><td>{{ .Type }}</td></tr>{{ end }}
{{ if .Author }}<tr><td>Author</td><td>{{ .Author }}</td></tr>{{ end }}
{{ if .Title }}<tr><td>Title</td><td>{{ .Title }}</td></tr>{{ end }}
{{ if .Abbr }}<tr><td>Short Title</td><td>{{ .Abbr }}</td></tr>{{ end }}
{{ if .Publication }}<tr><td>Publication Facts</td><td>{{ .Publication }}</td></tr>{{ end }}
{{ if .Date }}<tr><td>Date</td><td>{{ .Date }}</td></tr>{{ end }}
{{ if .Place }}<tr><td>Place</td><td>{{ .Place }}</td></tr>{{ end }}
{{ if .File }}<tr><td>File</td><td>{{ .File }}</td></tr>{{ end }}
{{ if .FileNumber }}<tr><td>File Number</td><td>{{ .FileNumber }}</td></tr>{{ end }}
{{ if .Form }}<tr><td>Form</td><td>{{ .Form }}</td></tr>{{ end }}
{{ if .URL }}<tr><td>URL</td><td><a href="{{ .URL }}" target="_blank">{{ .URL }}</a></td></tr>{{ end }}
{{ if .DocLocation }}<tr><td>Document Location</td><td>{{ .DocLocation }}</td></tr>{{ end }}
{{ if .DateViewed }}<tr><td>Date Viewed</td><td>{{ .DateViewed }}</td></tr>{{ end }}
{{ if .Text }}<tr><td>Text</td><td>{{ .Text }}</td></tr>{{ end }}
</table>
`

type sourceData struct {
	ID          string
	Author      string
	Abbr        string
	Publication string
	Text        string
	Title       string
	Form        string
	File        string
	FileNumber  string
	Type        string
	Place       string
	Date        string
	DateViewed  string
	URL         string
	DocLocation string
	Ref         string
}

func newSourceData(cx *cli.Context, source *gedcom.SourceRecord) (sourceData, error) {

	data := sourceData{
		ID:          source.Xref,
		Author:      source.Author,
		Title:       source.Title,
		Abbr:        source.Abbr,
		Publication: source.Publication,
		Text:        source.Text,
		File:        source.File,
		FileNumber:  source.FileNumber,
		Type:        source.Type,
		Place:       source.Place,
		Date:        source.Date,
		DateViewed:  source.DateViewed,
		URL:         source.URL,
		DocLocation: source.DocLocation,
	}

	// Build the reference string.
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(source.Xref, 1)
	data.Ref = fmt.Sprintf("%s. ", matches[0])
	var refs []string
	if data.Author != "" {
		refs = append(refs, data.Author)
	}
	if data.Title != "" {
		refs = append(refs, fmt.Sprintf("\"%s\"", data.Title))
	}
	if data.Type != "" {
		refs = append(refs, data.Type)
	}
	if data.URL != "" {
		refs = append(refs, fmt.Sprintf("[%s]", data.URL))
	}
	if data.Date != "" {
		refs = append(refs, data.Date)
	}
	if data.FileNumber != "" {
		refs = append(refs, data.FileNumber)
	}
	data.Ref += strings.Join(refs, ", ")

	return data, nil
}
