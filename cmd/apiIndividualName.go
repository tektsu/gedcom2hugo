package cmd

import (
	"fmt"
	"sort"

	"github.com/tektsu/gedcom"
)

type individualNameResponse struct {
	First     string `json:"first"`
	Last      string `json:"last"`
	Full      string `json:"full"`
	LastFirst string `json:"lastfirst"`
	Citations []int  `json:"citations"`
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
