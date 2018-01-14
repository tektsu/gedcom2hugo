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
		ID:        strings.ToLower(source.Xref),
		Author:    source.Author,
		Title:     source.Title,
		Ref:       source.GetReferenceString(),
		Citations: make(sourceCitationResponses),
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
