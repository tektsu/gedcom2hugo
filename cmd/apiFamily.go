package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tektsu/gedcom"
)

type familyResponse struct {
	ID        string                       `json:"id"`
	Note      string                       `json:"note"`
	Ref       *familyReferenceResponse     `json:"ref"`
	Events    []*eventResponse             `json:"events"`
	Children  individualReferenceResponses `json:"children"`
	Citations citationResponses            `json:"citations"`
}

type familyResponses map[string]*familyResponse

type familyControl struct {
	api           *apiResponse
	citationCount int
	citationIndex map[string]int
	family        *gedcom.FamilyRecord
	response      *familyResponse
}

func newFamilyControl(api *apiResponse) *familyControl {
	ic := &familyControl{
		api:           api,
		citationCount: 0,
		citationIndex: make(map[string]int),
	}

	return ic
}

func (api *apiResponse) addFamilies() error {

	for _, family := range api.gc.Family {
		err := api.addFamily(family)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *apiResponse) addFamily(family *gedcom.FamilyRecord) error {

	fc := newFamilyControl(api)
	fc.family = family
	fc.response = &familyResponse{
		ID:        strings.ToLower(family.Xref),
		Citations: make(citationResponses),
	}
	ref, err := api.getFamilyIndexEntry(fc.response.ID)
	if err != nil {
		return err
	}
	fc.response.Ref = ref

	if family.Husband != nil {
		father, err := api.getIndividualIndexEntry(family.Husband.Xref)
		if err != nil {
			return err
		}
		fc.response.Ref.Husband = father
	}

	if family.Wife != nil {
		mother, err := api.getIndividualIndexEntry(family.Wife.Xref)
		if err != nil {
			return err
		}
		fc.response.Ref.Wife = mother
	}

	for _, i := range family.Child {
		child, err := api.getIndividualIndexEntry(i.Person.Xref)
		if err != nil {
			return err
		}
		fc.response.Children = append(fc.response.Children, child)
	}

	err = fc.addEvents()
	if err != nil {
		return err
	}

	// Get note.
	if len(family.Note) > 0 {
		for _, n := range family.Note {
			fc.response.Note += n.Note + "\n\n"
		}
	}

	api.families[fc.response.ID] = fc.response

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
