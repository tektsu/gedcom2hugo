package cmd

import (
	"github.com/iand/gedcom"
	"github.com/urfave/cli"
)

const personPageTemplate string = `+++
title = "{{ .Name.Name }}"
weight = {{ .AlphaWeight }}
categories = [
	"Person"
]
+++
# {{ .Name.Name }}

Sex: {{ .Sex  }}
`

type personName struct {
	Name string
}

type personData struct {
	ID          string
	Name        personName
	Aliases     []personName
	Sex         string
	AlphaWeight int64
}

func newPersonData(cx *cli.Context, people *personIndex, person *gedcom.IndividualRecord) (personData, error) {

	id := person.Xref
	data := personData{
		ID:          id,
		Sex:         person.Sex,
		AlphaWeight: (*people)[id].AlphaWeight,
	}
	for i, n := range person.Name {
		name := personName{
			Name: n.Name,
		}
		if i == 0 {
			name.Name = (*people)[id].FullName
			data.Name = name
		} else {
			data.Aliases = append(data.Aliases, name)
		}
	}

	return data, nil
}
