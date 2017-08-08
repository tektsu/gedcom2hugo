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
{{ if .Date }}docdate: "{{ .Date }}"{{ end }}
{{ if .Place }}place: "{{ .Place }}"{{ end }}
{{ if .File }}file: "{{ .File }}"{{ end }}
{{ if .Form }}form: "{{ .Form }}"{{ end }}
{{ if .FileNumber }}filenumber: "{{ .FileNumber }}"{{ end }}
{{ if .URL }}docurl: "{{ .URL }}"{{ end }}
{{ if .DocLocation }}doclocation: "{{ .DocLocation }}"{{ end }}
{{ if .DateViewed }}dateviewed: "{{ .DateViewed }}"{{ end }}
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
	v, err := strconv.Atoi(matches[0])
	if err != nil {
		panic(fmt.Sprintf("Error converting [%s] to integer", matches[0]))
	}
	data.RefNum = v

	var refs []string
	if data.Author != "" {
		refs = append(refs, data.Author)
	}
	if data.Title != "" {
		refs = append(refs, fmt.Sprintf("\"%s\"", data.Title))
	}
	if data.Date != "" {
		refs = append(refs, data.Date)
	}
	if data.Place != "" {
		refs = append(refs, data.Place)
	}
	if data.Type != "" {
		refs = append(refs, data.Type)
	}
	if data.URL != "" {
		refs = append(refs, fmt.Sprintf("[%s]", data.URL))
	}
	if data.FileNumber != "" {
		refs = append(refs, data.FileNumber)
	}
	data.Ref = strings.Join(refs, ", ")

	return data, nil
}
