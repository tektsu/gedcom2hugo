package cmd

import (
	"fmt"

	"github.com/tektsu/gedcom"
)

type individualNameResponse struct {
	First     string `json:"first"`
	Last      string `json:"last"`
	Full      string `json:"full"`
	LastFirst string `json:"lastfirst"`
	Citations []int  `json:"citations"`
}

func newIndividualNameResponse(name *gedcom.NameRecord, ccb citationSubCallback) (*individualNameResponse, error) {

	given, family := extractNames(name.Name)
	response := &individualNameResponse{
		First:     given,
		Last:      family,
		Full:      fmt.Sprintf("%s %s", given, family),
		LastFirst: fmt.Sprintf("%s, %s", family, given),
		Citations: ccb(name.Citation),
	}

	return response, nil
}
