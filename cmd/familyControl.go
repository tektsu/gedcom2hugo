package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tektsu/gedcom"
)

func newFamilyControl(api *apiControl) *familyControl {
	ic := &familyControl{
		api:           api,
		citationCount: 0,
		citationIndex: make(map[string]int),
	}

	return ic
}

func (fc *familyControl) addPhotos() error {
	for _, o := range fc.family.Object {
		if o.File.Form != "jpg" && o.File.Form != "png" {
			continue
		}
		p := fc.api.addPhotoForFamily(o, fc.response)
		fc.response.Photos = append(fc.response.Photos, p)
	}

	return nil
}

func (fc *familyControl) addCitations(citations []*gedcom.CitationRecord) []int {
	fc.api.addFamilyCitations(fc.response.ID, citations)

	var citationList []int
	for _, citation := range citations {
		indexKey := fmt.Sprintf("%s:%s", citation.Source.Xref, citation.Page)
		var citationNumber int
		var exists bool
		if citationNumber, exists = fc.citationIndex[indexKey]; !exists {
			fc.citationCount++
			citationNumber = fc.citationCount
			fc.citationIndex[indexKey] = citationNumber
			fc.response.Citations[citationNumber] = &citationResponse{
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

func (fc *familyControl) newEventResponse(event *gedcom.EventRecord) (*eventResponse, error) {

	response := &eventResponse{
		Name:      event.Tag,
		Tag:       event.Tag,
		Value:     event.Value,
		Type:      event.Type,
		Date:      event.Date,
		Place:     event.Place.Name,
		Citations: fc.addCitations(event.Citation),
	}
	if response.Tag == "EVEN" && response.Type != "" {
		response.Name = response.Type
	}
	name, exists := tagTable[response.Name]
	if exists {
		response.Name = name
	}

	for _, note := range event.Note {
		response.Notes = append(response.Notes, note.Note)
	}

	response.Citations = append(response.Citations, fc.addCitations(event.Place.Citation)...) // Append place citations to the event

	return response, nil
}

func (fc *familyControl) addEvents() error {
	for _, r := range fc.family.Event {

		eventResponse, err := fc.newEventResponse(r)
		if err != nil {
			return err
		}
		if eventResponse.Tag == "MARR" {
			fc.response.Ref.Married = eventResponse.Date
		}
		fc.response.Events = append(fc.response.Events, eventResponse)
	}

	return nil
}
