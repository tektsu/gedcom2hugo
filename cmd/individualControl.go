package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iand/gedcom"
)

func newIndividualControl(api *apiControl) *individualControl {
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
				SourceRef: GetReferenceString(citation.Source),
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

	for i := range individual.UserDefined {
		if individual.UserDefined[i].Tag == "_PHOTO" {
			ic.response.Ref.Photo = filepath.Base(individual.UserDefined[i].Value)
			break
		}
	}

	if _, ok := ic.api.individuals[ic.response.ID]; ok {
		return fmt.Errorf("in creating individual record [%+v], id is already used: [%+v]", individual, ic.api.individuals[ic.response.ID])
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

	for _, note := range individual.Note {
		ic.response.Notes = append(ic.response.Notes, note.Note)
	}

	ic.api.individuals[ic.response.ID] = ic.response

	return nil
}

func (ic *individualControl) newIndividualNameResponse(name *gedcom.NameRecord) (*individualNameResponse, error) {
	given, family := extractNames(name.Name)
	response := &individualNameResponse{
		First:     given,
		Last:      family,
		Full:      fmt.Sprintf("%s %s", given, family),
		LastFirst: fmt.Sprintf("%s, %s", family, given),
		Citations: ic.addCitations(name.Citation),
	}

	return response, nil
}

func (ic *individualControl) addNames() error {
	for i, n := range ic.individual.Name {
		nameResponse, err := ic.newIndividualNameResponse(n)
		if err != nil {
			return err
		}

		if i == 0 {
			ic.response.Name = nameResponse
		} else {
			ic.response.Aliases = append(ic.response.Aliases, nameResponse)
		}
	}
	ic.response.Name.Citations = append(ic.response.Name.Citations, ic.addCitations(ic.individual.Citation)...) // Append general citations to the primary name
	sort.Ints(ic.response.Name.Citations)

	return nil
}

func (ic *individualControl) newFamilyLinkResponse(flr *gedcom.FamilyLinkRecord) (*familyLinkResponse, error) {
	if flr.Family == nil {
		return nil, nil
	}

	response := &familyLinkResponse{
		ID:        strings.ToLower(flr.Family.Xref),
		Pedigree:  flr.Type,
	}

	if flr.Family.Husband != nil {
		father, err := ic.api.getIndividualIndexEntry(flr.Family.Husband.Xref)
		if err != nil {
			return response, err
		}
		response.Father = father
	}

	if flr.Family.Wife != nil {
		mother, err := ic.api.getIndividualIndexEntry(flr.Family.Wife.Xref)
		if err != nil {
			return response, err
		}
		response.Mother = mother
	}

	for _, i := range flr.Family.Child {
		child, err := ic.api.getIndividualIndexEntry(i.Xref)
		if err != nil {
			return response, err
		}
		response.Children = append(response.Children, child)
	}

	for _, e := range flr.Family.Event {
		event, err := ic.newEventResponse(e)
		if err != nil {
			return response, err
		}
		response.Events = append(response.Events, event)
	}

	return response, nil
}

func (ic *individualControl) addParentFamilies() error {
	for _, fr := range ic.individual.Parents {
		if fr.Family != nil {
			family, err := ic.newFamilyLinkResponse(fr)
			if err != nil {
				return err
			}
			ic.response.ParentsFamily = append(ic.response.ParentsFamily, family)
		}
	}

	return nil
}

func (ic *individualControl) addFamilies() error {
	for _, fr := range ic.individual.Family {
		if fr.Family != nil {
			family, err := ic.newFamilyLinkResponse(fr)
			if err != nil {
				return err
			}
			ic.response.Family = append(ic.response.Family, family)
		}
	}

	return nil
}

func (ic *individualControl) addPhotos() error {
	for _, o := range ic.individual.Media {
		for _, photo := range o.File {
			if photo.Format != "jpg" && photo.Format != "png" {
				continue
			}
			p := ic.api.addPhotoForIndividual(o, ic.response)
			ic.response.Photos = append(ic.response.Photos, p)
		}
	}

	return nil
}

func (ic *individualControl) addTopPhoto() error {
	if ic.individual.Media != nil {
		p := ic.api.addPhoto(ic.individual.Media[0]) // Don't use addPhotoForIndividual() here, or there will be a duplicate on the photo page
		ic.response.TopPhoto = p
	}

	return nil
}

func (ic *individualControl) newEventResponse(event *gedcom.EventRecord) (*eventResponse, error) {
	response := &eventResponse{
		Name:      event.Tag,
		Tag:       event.Tag,
		Value:     event.Value,
		Type:      event.Type,
		Date:      event.Date,
		Place:     event.Place.Name,
		Citations: ic.addCitations(event.Citation),
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

	response.Citations = append(response.Citations, ic.addCitations(event.Place.Citation)...) // Append place citations to the event

	return response, nil
}

func (ic *individualControl) addAttributes() error {
	for _, r := range ic.individual.Attribute {

		if r.Tag == "SSN" { // Skip social security number.
			continue
		}
		eventResponse, err := ic.newEventResponse(r)
		if err != nil {
			return err
		}
		ic.response.Attributes = append(ic.response.Attributes, eventResponse)
	}

	return nil
}

func (ic *individualControl) addEvents() error {
	for _, r := range ic.individual.Event {

		if r.Tag == "Photo" {
			continue
		}
		eventResponse, err := ic.newEventResponse(r)
		if err != nil {
			return err
		}
		if eventResponse.Name == "Birth" {
			ic.response.Ref.Birth = eventResponse.Date
		}
		if eventResponse.Name == "Death" {
			ic.response.Ref.Death = eventResponse.Date
		}

		ic.response.Events = append(ic.response.Events, eventResponse)
	}

	return nil
}
