package cmd

import (
	"fmt"

	"github.com/iand/gedcom"
	"github.com/urfave/cli"
)

const personPageTemplate string = `+++
title = "{{ .Name.Full }}"
url: "/{{ .ID }}/"
categories = [
	"Person"
]
+++
# {{ .Name.Full }}

Sex: {{ .Sex  }}
`

type personName struct {
	Full      string
	Last      string
	LastFirst string
}

type personData struct {
	ID      string
	Name    personName
	Aliases []personName
	Sex     string
}

func newPersonData(cx *cli.Context, people *personIndex, person *gedcom.IndividualRecord) (personData, error) {

	id := person.Xref
	data := personData{
		ID:  id,
		Sex: person.Sex,
	}
	for i, n := range person.Name {
		given, family := extractNames(n.Name)
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
	}

	return data, nil
}
