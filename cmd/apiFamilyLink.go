package cmd

import (
	"github.com/tektsu/gedcom"
)

type familyLinkResponse struct {
	ID        string                       `json:"id"`
	Pedigree  string                       `json:"pedigree"`
	AdoptedBy string                       `json:"adoptedby"`
	Events    []*eventResponse             `json:"events"`
	Mother    *individualReferenceResponse `json:"mother"`
	Father    *individualReferenceResponse `json:"father"`
	Children  individualReferenceResponses `json:"children"`
}

func (ic *individualControl) newFamilyLinkResponse(flr *gedcom.FamilyLinkRecord) (*familyLinkResponse, error) {

	if flr.Family == nil {
		return nil, nil
	}

	response := &familyLinkResponse{
		ID:        flr.Family.Xref,
		Pedigree:  flr.Pedigree,
		AdoptedBy: flr.AdoptedBy,
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
		child, err := ic.api.getIndividualIndexEntry(i.Person.Xref)
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
