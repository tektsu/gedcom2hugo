package cmd

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tektsu/gedcom"
)

type sourceCitationResponse struct {
	Individuals map[string]bool `json:"individuals"`
}

func newCitationResponse() *sourceCitationResponse {
	return &sourceCitationResponse{
		Individuals: make(map[string]bool),
	}
}

type sourceCitationResponses map[string]*sourceCitationResponse

type sourceResponse struct {
	ID        string                  `json:"id"`
	Author    string                  `json:"author"`
	Title     string                  `json:"title"`
	File      []string                `json:"file"`
	RefNum    int                     `json:"refnum"`
	Ref       string                  `json:"ref"`
	Note      string                  `json:"note"`
	Citations sourceCitationResponses `json:"citations"`
}

type sourceResponses map[string]*sourceResponse

func newSourceResponses() (sourceResponses, error) {
	responses := make(sourceResponses)
	return responses, nil
}

func newSourceResponsesFromGedcom(gc *gedcom.Gedcom) (sourceResponses, error) {
	responses, err := newSourceResponses()
	if err != nil {
		return responses, err
	}
	err = responses.addAll(gc.Source)
	if err != nil {
		return responses, err
	}
	return responses, nil
}

func (s sourceResponses) add(source *gedcom.SourceRecord) (*sourceResponse, error) {

	response := &sourceResponse{
		ID:        strings.ToLower(source.Xref),
		Author:    source.Author,
		Title:     source.Title,
		Ref:       source.GetReferenceString(),
		Citations: make(sourceCitationResponses),
	}

	if _, ok := s[response.ID]; ok {
		return response, fmt.Errorf("In creating source record [%+v], id is already used: [%+v]", source, s[response.ID])
	}

	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(source.Xref, 1)
	v, err := strconv.Atoi(matches[0])
	if err != nil {
		panic(fmt.Sprintf("Error converting [%s] to integer", matches[0]))
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

	s[response.ID] = response

	return response, nil
}

func (s sourceResponses) addAll(sources []*gedcom.SourceRecord) error {

	for _, source := range sources {
		_, err := s.add(source)
		if err != nil {
			return err
		}
	}

	return nil
}
