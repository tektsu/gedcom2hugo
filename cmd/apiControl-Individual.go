package cmd

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tektsu/gedcom"
)

func (api *apiControl) addIndividuals() error {
	for _, individual := range api.gc.Individual {
		err := api.addIndividual(individual)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *apiControl) getIndividualIndexEntry(id string) (*individualReferenceResponse, error) {
	id = strings.ToLower(id)
	var entry *individualReferenceResponse
	var ok bool
	if entry, ok = api.indIndex[id]; !ok {
		entry = &individualReferenceResponse{ID: id}
		api.indIndex[id] = entry
	}

	return entry, nil
}

func (api *apiControl) addIndividual(individual *gedcom.IndividualRecord) error {
	ic := newIndividualControl(api)
	err := ic.build(individual)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiControl) addIndividualCitations(individualID string, citations []*gedcom.CitationRecord) error {
	for _, citation := range citations {
		sourceID := strings.ToLower(citation.Source.Xref)
		var c *sourceCitationResponse
		if _, ok := api.sources[sourceID].Citations[citation.Page]; ok {
			c = api.sources[sourceID].Citations[citation.Page]
		} else {
			c = newSourceCitationResponse()
			api.sources[sourceID].Citations[citation.Page] = c
		}
		c.Individuals[individualID], _ = api.getIndividualIndexEntry(individualID)
	}

	return nil
}

func (api *apiControl) exportIndividualAPI() error {
	individualAPIDir := filepath.Join(api.cx.String("project"), "static", "api", "individual")
	err := os.MkdirAll(individualAPIDir, 0777)
	if err != nil {
		return err
	}
	for id, individual := range api.individuals {
		file := filepath.Join(individualAPIDir, strings.ToLower(id+".json"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}

		j, err := json.Marshal(individual)
		if err != nil {
			fh.Close()
			return err
		}
		_, err = fh.Write(j)
		if err != nil {
			fh.Close()
			return err
		}
		fh.Close()
	}

	return nil
}

func (api *apiControl) exportIndividualPages() error {

	const personPageTemplate = `---
title: "{{ .Ref.Name }}{{ if or .Ref.Birth .Ref.Death }} ({{ .Ref.Birth }} - {{ .Ref.Death }}){{ end }}"
url: "/{{ .ID }}/"
categories:
  - Person
{{ if .Ref.LastNames }}lastnames:
  {{ range .Ref.LastNames }}- {{ . }}{{ end }}
{{- end }}
{{ if .Ref.Photo }}portrait: {{ .Ref.Photo }}{{end}}
---
<script src="/js/jquery.min.js"></script>
<script src="/js/idrisutil.js"></script>
<script src="/js/individualdisplay.js"></script>

<link rel="stylesheet" href="/js/photoswipe.css">
<link rel="stylesheet" href="/js/default-skin/default-skin.css">
<script src="/js/photoswipe.min.js"></script>
<script src="/js/photoswipe-ui-default.min.js"></script>

<script>
$(document).ready(function(){
    individualdisplay("{{ .ID }}")
});
</script>

<div id="display"></div>

<div id="raw"></div>
`

	personDir := filepath.Join(api.cx.String("project"), "content", "person")
	err := os.MkdirAll(personDir, 0777)
	if err != nil {
		return err
	}

	for _, individual := range api.individuals {
		file := filepath.Join(personDir, individual.ID+".md")

		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fh.Close()

		tpl := template.New("person")
		tpl, err = tpl.Parse(personPageTemplate)
		if err != nil {
			return err
		}
		err = tpl.Execute(fh, individual)
	}

	return nil
}

// extractNames splits a full name into a given name and a family name.
func extractNames(name string) (string, string) {
	var given, family string

	re := regexp.MustCompile("^([^/]+) +/(.+)/(.*)$")
	names := re.FindStringSubmatch(name)
	given = names[1]
	family = names[2]

	return given, family
}
