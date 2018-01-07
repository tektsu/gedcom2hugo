package cmd

import (
	"fmt"
	"strings"

	"github.com/tektsu/gedcom"
)

type familyMemberResponse struct {
	ID        string `json:"id"`
	Sex       string `json:"sex"`
	Name      string `json:"name"`
	Citations []int  `json:"citations"`
}

func newFamilyMemberResponse(i *gedcom.IndividualRecord, ccb citationSubCallback) (*familyMemberResponse, error) {
	child := &familyMemberResponse{
		ID:        strings.ToLower(i.Xref),
		Sex:       i.Sex,
		Citations: ccb(i.Name[0].Citation),
	}
	given, family := extractNames(i.Name[0].Name)
	child.Name = fmt.Sprintf("%s %s", given, family)

	return child, nil
}

type familyResponse struct {
	ID        string                  `json:"id"`
	Pedigree  string                  `json:"pedigree"`
	AdoptedBy string                  `json:"adoptedby"`
	Events    []*eventResponse        `json:"events"`
	Mother    *familyMemberResponse   `json:"mother"`
	Father    *familyMemberResponse   `json:"father"`
	Children  []*familyMemberResponse `json:"children"`
}

func newFamilyResponse(flr *gedcom.FamilyLinkRecord, ccb citationSubCallback) (*familyResponse, error) {

	if flr.Family == nil {
		return nil, nil
	}

	response := &familyResponse{
		ID:        flr.Family.Xref,
		Pedigree:  flr.Pedigree,
		AdoptedBy: flr.AdoptedBy,
	}

	if flr.Family.Husband != nil {
		father, err := newFamilyMemberResponse(flr.Family.Husband, ccb)
		if err != nil {
			return response, err
		}
		response.Father = father
	}

	if flr.Family.Wife != nil {
		mother, err := newFamilyMemberResponse(flr.Family.Wife, ccb)
		if err != nil {
			return response, err
		}
		response.Mother = mother
	}

	for _, i := range flr.Family.Child {
		child, err := newFamilyMemberResponse(i.Person, ccb)
		if err != nil {
			return response, err
		}
		response.Children = append(response.Children, child)
	}

	for _, e := range flr.Family.Event {
		event, err := newEventResponse(e, ccb)
		if err != nil {
			return response, err
		}
		response.Events = append(response.Events, event)
	}

	return response, nil
}
