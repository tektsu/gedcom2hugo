package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tektsu/gedcom"
)

type individualResponse struct {
	ID            string                    `json:"id"`
	Sex           string                    `json:"sex"`
	Birth         string                    `json:"birth"`
	Death         string                    `json:"death"`
	Name          *individualNameResponse   `json:"name"`
	Aliases       []*individualNameResponse `json:"aliases"`
	Events        []*eventResponse          `json:"events"`
	Attributes    []*eventResponse          `json:"attributes"`
	ParentsFamily []*familyResponse         `json:"parentsfamily"`
	Family        []*familyResponse         `json:"family"`
	TopPhoto      *photoResponse            `json:"topphoto"`
	Photos        []*photoResponse          `json:"photos"`
	Citations     citationResponses         `json:"citations"`
	//LastNames     []string
}

type individualResponses map[string]*individualResponse

func newIndividualResponses() (individualResponses, error) {
	responses := make(individualResponses)
	return responses, nil
}

func (r individualResponses) add(individual *gedcom.IndividualRecord, iccb citationCallback, photocb photoCallback) (*individualResponse, error) {

	citationCount := 0
	citationIndex := make(map[string]int)

	response := &individualResponse{
		ID:        strings.ToLower(individual.Xref),
		Sex:       individual.Sex,
		Citations: make(citationResponses),
	}
	if individual.Sex != "M" && individual.Sex != "F" {
		individual.Sex = "U"
	}

	if _, ok := r[response.ID]; ok {
		return response, fmt.Errorf("In creating individual record [%+v], id is already used: [%+v]", individual, r[response.ID])
	}

	// Callback function for citations
	ccb := func(citations []*gedcom.CitationRecord) []int {
		iccb(response.ID, citations)

		var citationList []int
		for _, citation := range citations {
			indexKey := fmt.Sprintf("%s:%s", citation.Source.Xref, citation.Page)
			var citationNumber int
			var exists bool
			if citationNumber, exists = citationIndex[indexKey]; !exists {
				citationCount++
				citationNumber = citationCount
				citationIndex[indexKey] = citationNumber
				response.Citations[citationNumber] = &citationResponse{
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

	// Add in the person's names.
	for i, n := range individual.Name {

		nameResponse, err := newIndividualNameResponse(n, ccb)
		if err != nil {
			return response, err
		}

		if i == 0 {
			response.Name = nameResponse
		} else {
			response.Aliases = append(response.Aliases, nameResponse)
		}
	}
	response.Name.Citations = append(response.Name.Citations, ccb(individual.Citation)...) // Append general citations to the primary name
	sort.Ints(response.Name.Citations)

	// Add in personal attributes.
	for _, r := range individual.Attribute {

		if r.Tag == "SSN" { // Skip social security number.
			continue
		}
		eventResponse, err := newEventResponse(r, ccb)
		if err != nil {
			return response, err
		}
		response.Attributes = append(response.Attributes, eventResponse)
	}

	// Add in personal events.
	for _, r := range individual.Event {

		if r.Tag == "Photo" {
			continue
		}
		eventResponse, err := newEventResponse(r, ccb)
		if err != nil {
			return response, err
		}
		if eventResponse.Name == "Birth" {
			response.Birth = eventResponse.Date
		}
		if eventResponse.Name == "Death" {
			response.Death = eventResponse.Date
		}
		response.Events = append(response.Events, eventResponse)
	}

	// Add in the parent's families.
	for _, fr := range individual.Parents {
		if fr.Family != nil {
			family, err := newFamilyResponse(fr, ccb)
			if err != nil {
				return response, err
			}
			response.ParentsFamily = append(response.ParentsFamily, family)
		}
	}

	// Add in the families.
	for _, fr := range individual.Family {
		if fr.Family != nil {
			family, err := newFamilyResponse(fr, ccb)
			if err != nil {
				return response, err
			}
			response.Family = append(response.Family, family)
		}
	}

	// Add in photos
	for _, o := range individual.Object {

		if o.File.Form != "jpg" && o.File.Form != "png" {
			continue
		}
		p := photocb(o, response)
		response.Photos = append(response.Photos, p)
	}

	// Get Top photo
	if individual.Photo != nil {
		p := photocb(individual.Photo, response)
		response.TopPhoto = p
	}

	r[response.ID] = response

	return response, nil
}

func (r individualResponses) addAll(individuals []*gedcom.IndividualRecord, iccb citationCallback, photocb photoCallback) error {

	for _, individual := range individuals {
		_, err := r.add(individual, iccb, photocb)
		if err != nil {
			return err
		}
	}

	return nil
}
