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
	ic.api.addCitations(ic.response.ID, citations)

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
	ic.individual = individual
	ic.response = &individualResponse{
		ID:        strings.ToLower(individual.Xref),
		Sex:       individual.Sex,
		Citations: make(citationResponses),
	}
	if ic.response.Sex != "M" && ic.response.Sex != "F" {
		ic.response.Sex = "U"
	}

	if _, ok := ic.api.individuals[ic.response.ID]; ok {
		return fmt.Errorf("In creating individual record [%+v], id is already used: [%+v]", individual, ic.api.individuals[ic.response.ID])
	}

	var err error

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
