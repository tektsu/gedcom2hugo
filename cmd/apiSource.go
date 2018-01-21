package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tektsu/gedcom"
)

type sourceCitationResponse struct {
	Individuals map[string]*individualReferenceResponse `json:"individuals"`
	Families    map[string]*familyReferenceResponse     `json:"families"`
}

func newCitationResponse() *sourceCitationResponse {
	return &sourceCitationResponse{
		Individuals: make(map[string]*individualReferenceResponse),
		Families:    make(map[string]*familyReferenceResponse),
	}
}

type sourceCitationResponses map[string]*sourceCitationResponse

type sourceResponse struct {
	ID          string                  `json:"id"`
	Author      string                  `json:"author"`
	Title       string                  `json:"title"`
	Publication string                  `json:"publication"`
	File        []string                `json:"file"`
	RefNum      int                     `json:"refnum"`
	Ref         string                  `json:"ref"`
	Note        string                  `json:"note"`
	Citations   sourceCitationResponses `json:"citations"`
}

type sourceResponses map[string]*sourceResponse

func newSourceResponses() (sourceResponses, error) {
	responses := make(sourceResponses)
	return responses, nil
}

func (api *apiResponse) addSources() error {

	for _, source := range api.gc.Source {
		err := api.addSource(source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *apiResponse) addSource(source *gedcom.SourceRecord) error {

	response := &sourceResponse{
		ID:          strings.ToLower(source.Xref),
		Author:      source.Author,
		Title:       source.Title,
		Publication: source.Publication,
		Ref:         source.GetReferenceString(),
		Citations:   make(sourceCitationResponses),
	}

	if _, ok := api.sources[response.ID]; ok {
		return fmt.Errorf("In creating source record [%+v], id is already used: [%+v]", source, api.sources[response.ID])
	}

	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(source.Xref, 1)
	v, err := strconv.Atoi(matches[0])
	if err != nil {
		return fmt.Errorf("Error converting [%s] to integer", matches[0])
	}
	response.RefNum = v

	// Get files.
	if len(source.Object) > 0 {
		r1 := regexp.MustCompile("^.*Roots/")
		r2 := regexp.MustCompile("^.*/idris_project/sources/")
		m, _ := regexp.Compile("^/")
		for _, o := range source.Object {
			name := r1.ReplaceAllString(o.File.Name, "")
			name = r2.ReplaceAllString(name, "")
			if m.MatchString(name) {
				name = filepath.Base(name)
			}
			response.File = append(response.File, name)
		}
	}

	// Get note.
	if len(source.Note) > 0 {
		for _, n := range source.Note {
			response.Note += n.Note + "\n\n"
		}
	}

	api.sources[response.ID] = response

	return nil
}

func (api *apiResponse) exportSourceAPI() error {
	sourceAPIDir := filepath.Join(api.cx.String("project"), "static", "api", "source")
	err := os.MkdirAll(sourceAPIDir, 0777)
	if err != nil {
		return err
	}
	for id, source := range api.sources {
		file := filepath.Join(sourceAPIDir, strings.ToLower(id+".json"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}

		j, err := json.Marshal(source)
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

func (api *apiResponse) exportSourcePages() error {

	// sourcePageTemplate is the tmplate used to generaate a source web page.
	const sourcePageTemplate string = `---
url: "/{{ .ID }}/"
categories:
  - Source
title: "Source: {{ if .Title }}{{ .Title }}{{ end }}"
{{ if .RefNum }}refnum: "{{ .RefNum }}"{{ end }}
---
<script src="/js/jquery.min.js"></script>
<script src="/js/sourcedisplay.js"></script>
<script>
$(document).ready(function(){
    sourcedisplay("{{ .ID }}")
});
</script>

<div id="display"></div>

<div id="raw"></div>
`

	sourceDir := filepath.Join(api.cx.String("project"), "content", "source")
	err := os.MkdirAll(sourceDir, 0777)
	if err != nil {
		return err
	}
	for _, source := range api.sources {
		file := filepath.Join(sourceDir, strings.ToLower(source.ID+".md"))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		defer fh.Close()

		tpl := template.New("source")
		tpl, err = tpl.Parse(sourcePageTemplate)
		if err != nil {
			return err
		}
		err = tpl.Execute(fh, source)
	}

	return nil
}
