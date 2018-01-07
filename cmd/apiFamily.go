package cmd

import (
	"github.com/tektsu/gedcom"
)

type familyResponse struct {
	ID        string                         `json:"id"`
	Pedigree  string                         `json:"pedigree"`
	AdoptedBy string                         `json:"adoptedby"`
	Events    []*eventResponse               `json:"events"`
	Mother    *individualReferenceResponse   `json:"mother"`
	Father    *individualReferenceResponse   `json:"father"`
	Children  []*individualReferenceResponse `json:"children"`
}

func (ic *individualControl) newFamilyResponse(flr *gedcom.FamilyLinkRecord) (*familyResponse, error) {

	if flr.Family == nil {
		return nil, nil
	}

	response := &familyResponse{
		ID:        flr.Family.Xref,
		Pedigree:  flr.Pedigree,
		AdoptedBy: flr.AdoptedBy,
	}

	if flr.Family.Husband != nil {
		father, err := ic.newIndividualReferenceResponse(flr.Family.Husband)
		if err != nil {
			return response, err
		}
		response.Father = father
	}

	if flr.Family.Wife != nil {
		mother, err := ic.newIndividualReferenceResponse(flr.Family.Wife)
		if err != nil {
			return response, err
		}
		response.Mother = mother
	}

	for _, i := range flr.Family.Child {
		child, err := ic.newIndividualReferenceResponse(i.Person)
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
			family, err := ic.newFamilyResponse(fr)
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
			family, err := ic.newFamilyResponse(fr)
			if err != nil {
				return err
			}
			ic.response.Family = append(ic.response.Family, family)
		}
	}

	return nil
}
