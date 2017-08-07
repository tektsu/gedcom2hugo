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
{{ range .LastNames }}  - {{ . }}{{ end }}
{{- end }}
{{ if .Sex }}sex: "{{ .Sex }}"{{ end }}
---
# {{ .Name.Full }}{{ if .Name.SourcesInd }}{{ range .Name.SourcesInd }} <span class="citref">[{{ . }}]</span>{{ end }}{{ end }}

Sex: {{ .Sex  }}

{{ if .Sources -}}
<div id="sources">
{{ range $i, $s := .Sources }}{{ add $i 1 }}. {{ $s.Ref }}{{ end }}
</div>
{{- end }}
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
			ref := sourceList[r]
			data.Sources = append(data.Sources, SourceRef{
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
