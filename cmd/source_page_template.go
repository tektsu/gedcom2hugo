package cmd

import (
	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

const sourcePageTemplate string = `---
url: "/{{ .ID }}/"
categories:
  - Source
{{ if .Title }}title: "{{ .Title }}"{{ end }}
{{ if .Date }}date: "{{ .Date }}"{{ end }}
{{ if .DateViewed }}dateviewed: "{{ .DateViewed }}"{{ end }}
{{ if .File }}file: "{{ .File }}"{{ end }}
{{ if .FileNumber }}filenumber: "{{ .FileNumber }}"{{ end }}
{{ if .Type }}type: "{{ .Type }}"{{ end }}
{{ if .Place }}place: "{{ .Place }}"{{ end }}
{{ if .URL }}docurl: "{{ .URL }}"{{ end }}
{{ if .DocLocation }}doclocation: "{{ .DocLocation }}"{{ end }}
---
# {{ .Title }}
`

type sourceData struct {
	ID          string
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
}

func newSourceData(cx *cli.Context, source *gedcom.SourceRecord) (sourceData, error) {

	data := sourceData{
		ID:          source.Xref,
		Title:       source.Title,
		File:        source.File,
		FileNumber:  source.FileNumber,
		Type:        source.Type,
		Place:       source.Place,
		Date:        source.Date,
		DateViewed:  source.DateViewed,
		URL:         source.URL,
		DocLocation: source.DocLocation,
	}

	return data, nil
}
