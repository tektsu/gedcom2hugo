package cmd

import (
	"fmt"
	"strconv"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

const personPageTemplate string = `---
title: "{{ .Name.Full }}"
url: "/{{ .ID }}/"
categories:
  - Person
{{ if .LastNames }}lastnames:
  {{ range .LastNames }}- {{ . }}{{ end }}
{{- end }}
name:
  full: "{{ .Name.Full }}"
  last: "{{ .Name.Last }}"
  lastfirst: "{{ .Name.LastFirst }}"
  {{ if .Name.SourcesInd -}}
  sources:
  {{ range .Name.SourcesInd }}  - {{ . }}
  {{ end }}
  {{- end }}
{{ if .Sex }}sex: "{{ .Sex }}"{{ end }}
{{ if .Sources }}sources:
  {{ range .Sources }}-
    ref: {{ .Ref }}
    refnum: {{ .RefNum }}
  {{ end }}
{{- end }}
---
{{ "personbody" | shortcode }}
`

func newPersonData(cx *cli.Context, people *personIndex, person *gedcom.IndividualRecord) (personData, error) {

	cc := 0 // Citation Counter

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

		for _, c := range n.Citation {
			r, err := strconv.Atoi(c.Source.Xref[1:len(c.Source.Xref)])
			if err != nil {
				panic(err)
			}
			cc++
			name.SourcesInd = append(name.SourcesInd, cc)
			ref := sl[r]
			data.Sources = append(data.Sources, sourceRef{
				RefNum: r,
				Ref:    ref,
			})

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
