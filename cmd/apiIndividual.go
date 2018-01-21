package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/tektsu/gedcom"
)

type individualResponse struct {
	ID            string                       `json:"id"`
	Ref           *individualReferenceResponse `json:"ref"`
	Name          *individualNameResponse      `json:"name"`
	Aliases       []*individualNameResponse    `json:"aliases"`
	Events        []*eventResponse             `json:"events"`
	Attributes    []*eventResponse             `json:"attributes"`
	ParentsFamily []*familyLinkResponse        `json:"parentsfamily"`
	Family        []*familyLinkResponse        `json:"family"`
	TopPhoto      *photoResponse               `json:"topphoto"`
	Photos        []*photoResponse             `json:"photos"`
	Citations     citationResponses            `json:"citations"`
}

type individualResponses map[string]*individualResponse

func (api *apiResponse) addIndividuals() error {
	for _, individual := range api.gc.Individual {
		err := api.addIndividual(individual)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *apiResponse) addIndividual(individual *gedcom.IndividualRecord) error {
	ic := newIndividualControl(api)
	err := ic.build(individual)
	if err != nil {
		return err
	}

	return nil
}

type individualControl struct {
	api           *apiResponse
	citationCount int
	citationIndex map[string]int
	individual    *gedcom.IndividualRecord
	response      *individualResponse
}

func newIndividualControl(api *apiResponse) *individualControl {
	ic := &individualControl{
		api:           api,
		citationCount: 0,
		citationIndex: make(map[string]int),
	}

	return ic
}

func (ic *individualControl) addCitations(citations []*gedcom.CitationRecord) []int {
	ic.api.addIndividualCitations(ic.response.ID, citations)

	var citationList []int
	for _, citation := range citations {
		indexKey := fmt.Sprintf("%s:%s", citation.Source.Xref, citation.Page)
		var citationNumber int
		var exists bool
		if citationNumber, exists = ic.citationIndex[indexKey]; !exists {
			ic.citationCount++
			citationNumber = ic.citationCount
			ic.citationIndex[indexKey] = citationNumber
			ic.response.Citations[citationNumber] = &citationResponse{
				ID:        citationNumber,
				SourceID:  strings.ToLower(citation.Source.Xref),
				SourceRef: citation.Source.GetReferenceString(),
				Detail:    citation.Page,
			}
		}
		citationList = append(citationList, citationNumber)
	}

	sort.Ints(citationList)
	return citationList
}

func (ic *individualControl) build(individual *gedcom.IndividualRecord) error {
	var err error

	ic.individual = individual
	ic.response = &individualResponse{
		ID:        strings.ToLower(individual.Xref),
		Citations: make(citationResponses),
	}
	ic.response.Ref, err = ic.api.getIndividualIndexEntry(strings.ToLower(individual.Xref))
	if err != nil {
		return err
	}
	ic.response.Ref.Sex = individual.Sex
	if ic.response.Ref.Sex != "M" && ic.response.Ref.Sex != "F" {
		ic.response.Ref.Sex = "U"
	}
	given, family := extractNames(individual.Name[0].Name)
	ic.response.Ref.Name = fmt.Sprintf("%s %s", given, family)
	ic.response.Ref.LastNames = append(ic.response.Ref.LastNames, family)

	if individual.Photo != nil {
		ic.response.Ref.Photo = filepath.Base(individual.Photo.File.Name)
	}

	if _, ok := ic.api.individuals[ic.response.ID]; ok {
		return fmt.Errorf("In creating individual record [%+v], id is already used: [%+v]", individual, ic.api.individuals[ic.response.ID])
	}

	err = ic.addNames()
	if err != nil {
		return err
	}

	err = ic.addAttributes()
	if err != nil {
		return err
	}

	err = ic.addEvents()
	if err != nil {
		return err
	}

	err = ic.addParentFamilies()
	if err != nil {
		return err
	}

	err = ic.addFamilies()
	if err != nil {
		return err
	}

	err = ic.addPhotos()
	if err != nil {
		return err
	}

	err = ic.addTopPhoto()
	if err != nil {
		return err
	}

	ic.api.individuals[ic.response.ID] = ic.response

	return nil
}

func (api *apiResponse) exportIndividualAPI() error {
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

func (api *apiResponse) exportIndividualPages() error {

	const personPageTemplate string = `---
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
