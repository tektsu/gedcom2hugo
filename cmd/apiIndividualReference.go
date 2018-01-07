package cmd

import (
	"fmt"
	"strings"

	"github.com/tektsu/gedcom"
)

type individualReferenceResponse struct {
	ID        string `json:"id"`
	Sex       string `json:"sex"`
	Name      string `json:"name"`
	Citations []int  `json:"citations"`
}

func (ic *individualControl) newIndividualReferenceResponse(i *gedcom.IndividualRecord) (*individualReferenceResponse, error) {
	individual := &individualReferenceResponse{
		ID:        strings.ToLower(i.Xref),
		Sex:       i.Sex,
		Citations: ic.addCitations(i.Name[0].Citation),
	}
	given, family := extractNames(i.Name[0].Name)
	individual.Name = fmt.Sprintf("%s %s", given, family)

	return individual, nil
}
