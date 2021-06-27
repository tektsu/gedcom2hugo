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

	"github.com/iand/gedcom"
)

func (api *apiControl) addSources() error {
	for _, source := range api.gc.Source {
		err := api.addSource(source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *apiControl) addSource(source *gedcom.SourceRecord) error {
	response := &sourceResponse{
		ID:          strings.ToLower(source.Xref),
		Author:      source.Originator,
		Title:       source.Title,
		Publication: source.PublicationFacts,
		Ref:         GetReferenceString(source),
		Citations:   make(sourceCitationResponses),
	}

	if _, ok := api.sources[response.ID]; ok {
		return fmt.Errorf("in creating source record [%+v], id is already used: [%+v]", source, api.sources[response.ID])
	}

	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(source.Xref, 1)
	v, err := strconv.Atoi(matches[0])
	if err != nil {
		return fmt.Errorf("error converting [%s] to integer", matches[0])
	}
	response.RefNum = v

	// Get files.
	if len(source.Media) > 0 {
		r1 := regexp.MustCompile("^.*Roots/")
		r2 := regexp.MustCompile("^.*/idris_project/sources/")
		m, _ := regexp.Compile("^/")
		for _, o := range source.Media {
			for _, i := range o.File {
				name := r1.ReplaceAllString(i.Name, "")
				name = r2.ReplaceAllString(name, "")
				if m.MatchString(name) {
					name = filepath.Base(name)
				}
				response.File = append(response.File, name)
			}
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

func GetReferenceString(s *gedcom.SourceRecord) string {
	var refs []string
	if s.Originator != "" {
		refs = append(refs, s.Originator)
	}
	if s.Title != "" {
		refs = append(refs, fmt.Sprintf("\"%s\"", s.Title))
	}

	var pubFacts, pubName, pubDate string
	pubParts := strings.Split(s.PublicationFacts, ";")
	r, _ := regexp.Compile("^ *([^:]+): (.+)$")
	for _, p := range pubParts {
		match := r.FindStringSubmatch(p)
		if len(match) > 0 {
			switch label := match[1]; label {
			case "Location":
				pubFacts = match[2]
			case "Name":
				pubName = match[2]
			case "Date":
				pubDate = match[2]
			}
		}
	}
	if pubName != "" {
		if pubFacts != "" {
			pubFacts += ", "
		}
		pubFacts += pubName
	}
	if pubDate != "" {
		if pubFacts != "" {
			pubFacts += ", "
		}
		pubFacts += pubDate
	}
	if pubFacts != "" {
		refs = append(refs, "("+pubFacts+")")
	}

	return strings.Join(refs, ", ")
}

func (api *apiControl) exportSourceAPI() error {
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
			_ = fh.Close()
			return err
		}
		_, err = fh.Write(j)
		if err != nil {
			_ = fh.Close()
			return err
		}
		_ = fh.Close()
	}

	return nil
}

func (api *apiControl) exportSourcePages() error {
	// sourcePageTemplate is the tmplate used to generaate a source web page.
	const sourcePageTemplate = `---
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
		defer func(fh *os.File) {
			_ = fh.Close()
		}(fh)

		tpl := template.New("source")
		tpl, err = tpl.Parse(sourcePageTemplate)
		if err != nil {
			return err
		}
		err = tpl.Execute(fh, source)
	}

	return nil
}
