package cmd

import (
	"fmt"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

const personPageTemplate string = `---
title: "{{ .Name.Full }}"
url: "/{{ .ID }}/"
categories:
  - Person
{{ if .LastNames }}lastnames:
{{ range .LastNames }}  - {{ . }}{{ end }}
{{- end }}
{{ if .Sex }}sex: "{{ .Sex }}"{{ end }}
---
# {{ .Name.Full }}

Sex: {{ .Sex  }}
`

type personName struct {
	Full      string
	Last      string
	LastFirst string
}

type personData struct {
	ID        string
	Name      personName
	Aliases   []personName
	LastNames []string
	Sex       string
}

func newPersonData(cx *cli.Context, people *personIndex, person *gedcom.IndividualRecord) (personData, error) {

	id := person.Xref
	data := personData{
		ID:  id,
		Sex: person.Sex,
	}

	lastNames := make(map[string]bool)
	for i, n := range person.Name {
		given, family := extractNames(n.Name)
		lastNames[family] = true
		name := personName{
			Last:      family,
			Full:      fmt.Sprintf("%s %s", given, family),
			LastFirst: fmt.Sprintf("%s, %s", family, given),
		}
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

	return data, nil
}
