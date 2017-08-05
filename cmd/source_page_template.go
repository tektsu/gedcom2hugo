package cmd

import (
	"github.com/iand/gedcom"
	"github.com/urfave/cli"
)

const sourcePageTemplate string = `+++
title = "{{ .Title }}"
url: "/{{ .ID }}/"
categories = [
	"Source"
]
+++
# {{ .Title }}
`

type sourceData struct {
	ID    string
	Title string
}

func newSourceData(cx *cli.Context, source *gedcom.SourceRecord) (sourceData, error) {

	data := sourceData{
		ID:    source.Xref,
		Title: source.Title,
	}

	return data, nil
}
