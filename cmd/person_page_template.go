package cmd

import (
	"fmt"
	"strconv"

	"github.com/tektsu/gedcom"
	"github.com/urfave/cli"
)

const personPageTemplate string = `---
title: "{{ .Name.Full }}"
url: "/{{ .ID | ToLower }}/"
categories:
  - Person
{{ if .LastNames }}lastnames:
  {{ range .LastNames }}- {{ . }}{{ end }}
{{- end }}
---

<div id="person">

<div id="personal_info">
<table class="personal_info_table">
<tr><th>Name</th><td>{{ .Name.Full }}{{ if .Name.SourcesInd }}<sup>{{ range .Name.SourcesInd }} [{{ . }}]{{ end }}</sup>{{ end }}</td></tr>
<tr><th>Sex</th><td>{{ .Sex }}</td></tr>
</table>
</div>

{{ if .ParentsFamily }}
<div id="parents">
{{ range .ParentsFamily}}
<table class="parents_family">
<tr><th colspan="2">Parent's Family</th></tr>
<tr><th>Father</th><td>
	{{- if .Father.ID -}}
		<a href="/{{ .Father.ID | ToLower }}/">{{ .Father.Name }}</a>
		{{- if .Father.SourcesInd -}}
			<sup>{{ range .Father.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
<tr><th>Mother</th><td>
	{{- if .Mother.ID -}}
		<a href="/{{ .Mother.ID | ToLower }}/">{{ .Mother.Name }}</a>
		{{- if .Mother.SourcesInd -}}
			<sup>{{ range .Mother.SourcesInd }} [{{ . }}]{{ end }}</sup>
		{{- end -}}
	{{- end -}}<br />
</td></tr>
{{ $length := len .Children }} {{ if gt $length 0 }}
<tr><th>Siblings</th><td>
{{ range .Children }}
	<a href="/{{ .ID | ToLower }}/">{{ .Name }}</a>
	{{- if .SourcesInd -}}
		<sup>{{ range .SourcesInd }} [{{ . }}]{{ end }}</sup>
	{{- end -}}<br />
{{ end }}
</td></tr>
{{ end }}
</table>
{{ end }}
</div>
{{ end }}

{{ if .Sources -}}
<div id="sources">
<table class="sources_table">
<tr><th colspan="2">Sources</th></tr>
{{ range $i, $s := .Sources -}}
<tr><td>{{ add $i 1 }}.</td><td><a href="/s{{ $s.RefNum }}">{{ $s.Ref }}</a></td></tr>
{{ end -}}
</table>
</div>
{{- end }}

</div>
`

func newPersonData(cx *cli.Context, people personIndex, person *gedcom.IndividualRecord) (personData, error) {

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

	for _, fr := range person.Parents {
		if fr.Family != nil {
			f := personFamily{
				ID:        fr.Family.Xref,
				Pedigree:  fr.Pedigree,
				AdoptedBy: fr.AdoptedBy,
			}
			if fr.Family.Husband != nil {
				f.Father.ID = fr.Family.Husband.Xref
				f.Father.Name = people[fr.Family.Husband.Xref].FullName
				for _, c := range fr.Family.Husband.Name[0].Citation {
					r, err := strconv.Atoi(c.Source.Xref[1:len(c.Source.Xref)])
					if err != nil {
						panic(err)
					}
					cc++

					f.Father.SourcesInd = append(f.Father.SourcesInd, cc)
					data.Sources = append(data.Sources, sourceRef{
						RefNum: r,
						Ref:    sl[r],
					})
				}
			}
			if fr.Family.Wife != nil {
				f.Mother.ID = fr.Family.Wife.Xref
				f.Mother.Name = people[fr.Family.Wife.Xref].FullName
				for _, c := range fr.Family.Wife.Name[0].Citation {
					r, err := strconv.Atoi(c.Source.Xref[1:len(c.Source.Xref)])
					if err != nil {
						panic(err)
					}
					cc++

					f.Mother.SourcesInd = append(f.Mother.SourcesInd, cc)
					data.Sources = append(data.Sources, sourceRef{
						RefNum: r,
						Ref:    sl[r],
					})
				}
			}
			for _, cr := range fr.Family.Child {
				if cr.Xref == data.ID {
					continue
				}
				child := personRef{
					ID:   cr.Xref,
					Name: people[cr.Xref].GivenName,
				}
				for _, c := range cr.Name[0].Citation {
					r, err := strconv.Atoi(c.Source.Xref[1:len(c.Source.Xref)])
					if err != nil {
						panic(err)
					}
					cc++

					child.SourcesInd = append(child.SourcesInd, cc)
					data.Sources = append(data.Sources, sourceRef{
						RefNum: r,
						Ref:    sl[r],
					})
				}
				f.Children = append(f.Children, child)
			}

			data.ParentsFamily = append(data.ParentsFamily, f)
		}
	}

	return data, nil
}
